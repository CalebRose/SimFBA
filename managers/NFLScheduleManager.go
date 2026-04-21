package managers

import (
	"math/rand"
	"strconv"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
)

// NFL Division names (must match the Division field on NFLTeam records)
const (
	afcEast  = "AFC East"
	afcNorth = "AFC North"
	afcSouth = "AFC South"
	afcWest  = "AFC West"
	nfcEast  = "NFC East"
	nfcNorth = "NFC North"
	nfcSouth = "NFC South"
	nfcWest  = "NFC West"
)

// afcDivisions and nfcDivisions in the order used for rotating matchups.
var afcDivisions = []string{afcEast, afcNorth, afcSouth, afcWest}
var nfcDivisions = []string{nfcEast, nfcNorth, nfcSouth, nfcWest}

// afcIntraRotation maps each AFC division to which AFC division it plays
// (full 4-game schedule) by season ID modulo 3.
// Rotation index = (seasonID - 1) % 3
var afcIntraRotation = map[string][3]string{
	afcEast:  {afcNorth, afcSouth, afcWest},
	afcNorth: {afcEast, afcWest, afcSouth},
	afcSouth: {afcWest, afcEast, afcNorth},
	afcWest:  {afcSouth, afcNorth, afcEast},
}

var nfcIntraRotation = map[string][3]string{
	nfcEast:  {nfcNorth, nfcSouth, nfcWest},
	nfcNorth: {nfcEast, nfcWest, nfcSouth},
	nfcSouth: {nfcWest, nfcEast, nfcNorth},
	nfcWest:  {nfcSouth, nfcNorth, nfcEast},
}

// afcInterRotation and nfcInterRotation map each division to which opposite-conference
// division it plays (full 4-game schedule) by season ID modulo 4.
// Rotation index = (seasonID - 1) % 4
var afcInterRotation = map[string][4]string{
	afcEast:  {nfcNorth, nfcSouth, nfcEast, nfcWest},
	afcNorth: {nfcSouth, nfcEast, nfcWest, nfcNorth},
	afcSouth: {nfcEast, nfcWest, nfcNorth, nfcSouth},
	afcWest:  {nfcWest, nfcNorth, nfcSouth, nfcEast},
}

var nfcInterRotation = map[string][4]string{
	nfcEast:  {afcSouth, afcNorth, afcWest, afcEast},
	nfcNorth: {afcEast, afcWest, afcSouth, afcNorth},
	nfcSouth: {afcNorth, afcWest, afcEast, afcSouth}, // fixed: was {afcNorth, afcEast, afcEast, afcWest}
	nfcWest:  {afcWest, afcSouth, afcNorth, afcEast},
}

// westernStates are US states considered Western/Midwestern US
// (past the Mississippi River) for Sunday timeslot assignment.
var westernStates = map[string]bool{
	"MN": true, "IA": true, "MO": true, "AR": true, "LA": true,
	"WI": true, "IL": true, "MS": true,
	"ND": true, "SD": true, "NE": true, "KS": true, "OK": true, "TX": true,
	"MT": true, "WY": true, "CO": true, "NM": true, "AZ": true,
	"ID": true, "UT": true, "NV": true, "CA": true, "OR": true, "WA": true,
	"HI": true, "AK": true,
}

// thanksgiving3rdGameRivalries lists divisional rivalries eligible for the
// third Thanksgiving TNF slot (after DAL and DET games are assigned).
var thanksgiving3rdGameRivalries = [][2]string{
	{"BAL", "PIT"},
	{"NE", "BUF"},
	{"KC", "DEN"},
	{"SEA", "SF"},
	{"LAR", "SF"},
	{"LAC", "LV"},
	{"MIA", "NYJ"},
	{"NYJ", "NE"},
	{"ATL", "CAR"},
	{"NO", "ATL"},
	{"TB", "NO"},
}

type matchup struct {
	homeID       uint
	awayID       uint
	isDivisional bool // true when both teams share a division
}

type weekMatchup struct {
	homeID       uint
	awayID       uint
	week         int
	isDivisional bool
}

type addMatchupFn func(homeID, awayID uint, isDivisional bool)

// GenerateSimNFLSchedule generates a full 17-game regular season (18 weeks, 1 bye)
// for the SimNFL league and saves all game records to the database.
func GenerateSimNFLSchedule() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID := ts.NFLSeasonID
	prevSeasonID := seasonID - 1

	// Collect teams and build lookup structures.
	allTeams := GetAllNFLTeams()
	teamByID := make(map[uint]structs.NFLTeam, len(allTeams))
	teamsByDiv := make(map[string][]structs.NFLTeam)
	for _, t := range allTeams {
		teamByID[t.ID] = t
		teamsByDiv[t.Division] = append(teamsByDiv[t.Division], t)
	}

	// Build stadium map: homeTeamID -> Stadium
	stadiums := GetAllStadiums()
	stadiumByTeam := make(map[uint]structs.Stadium, len(stadiums))
	for _, s := range stadiums {
		stadiumByTeam[s.TeamID] = s
	}

	// Load previous season's standings and compute each team's divisional rank (1=first, 4=last).
	divRank := computeDivisionalRanks(prevSeasonID, teamsByDiv)

	// Determine rotation indices.
	intraIdx := (seasonID - 1) % 3   // 3-year intra-conference rotation
	interIdx := (seasonID - 1) % 4   // 4-year inter-conference rotation
	inter17Idx := (seasonID - 3) % 4 // 17th game uses the division played 2 seasons ago
	if inter17Idx < 0 {
		inter17Idx += 4
	}

	// Accumulate all matchups, deduplicating by ordered pair.
	scheduledPairs := make(map[[2]uint]bool)
	var allMatchups []matchup

	addMatchup := func(homeID, awayID uint, isDivisional bool) {
		key := [2]uint{homeID, awayID}
		rkey := [2]uint{awayID, homeID}
		if scheduledPairs[key] || scheduledPairs[rkey] {
			return
		}
		scheduledPairs[key] = true
		allMatchups = append(allMatchups, matchup{homeID: homeID, awayID: awayID, isDivisional: isDivisional})
	}

	// 1. Divisional games: home AND away vs each of the 3 division rivals (6 games per team).
	for _, teams := range teamsByDiv {
		for i := 0; i < len(teams); i++ {
			for j := i + 1; j < len(teams); j++ {
				addMatchup(teams[i].ID, teams[j].ID, true)
				addMatchup(teams[j].ID, teams[i].ID, true)
			}
		}
	}

	// 2. Intra-conference full-division rotation (4 games per team).
	scheduleFullDivisionMatchup(teamsByDiv, afcDivisions, afcIntraRotation, intraIdx, addMatchup)
	scheduleFullDivisionMatchup(teamsByDiv, nfcDivisions, nfcIntraRotation, intraIdx, addMatchup)

	// 3. Two same-placement intra-conference games from the 2 non-rotation divisions (2 games per team).
	schedulePlacementGames(teamsByDiv, afcDivisions, afcIntraRotation, intraIdx, divRank, addMatchup)
	schedulePlacementGames(teamsByDiv, nfcDivisions, nfcIntraRotation, intraIdx, divRank, addMatchup)

	// 4. Inter-conference full-division rotation (4 games per team).
	scheduleFullInterConferenceMatchup(teamsByDiv, afcDivisions, afcInterRotation, interIdx, addMatchup)
	scheduleFullInterConferenceMatchup(teamsByDiv, nfcDivisions, nfcInterRotation, interIdx, addMatchup)

	// 5. 17th game: same-placement opponent from the inter-conference division played 2 years ago.
	schedule17thGame(teamsByDiv, afcDivisions, afcInterRotation, inter17Idx, divRank, addMatchup)
	schedule17thGame(teamsByDiv, nfcDivisions, nfcInterRotation, inter17Idx, divRank, addMatchup)

	// Assign matchups to weeks 1-18 with bye-week and week-18 divisional constraints.
	var dalID, detID uint
	for id, t := range teamByID {
		switch t.TeamAbbr {
		case "DAL":
			dalID = id
		case "DET":
			detID = id
		}
	}
	weekSchedule := assignMatchupsToWeeks(allMatchups, teamByID, dalID, detID)

	// Build full NFLGame records with timeslots.
	finalGames := buildNFLGameRecords(weekSchedule, teamByID, stadiumByTeam, seasonID, dalID, detID)

	// Set HomePreviousBye / AwayPreviousBye flags.
	applyByeWeekFlags(finalGames)

	// Save all games, then generate weather data.
	repository.CreateNFLGameRecordsBatch(db, finalGames, 250)
	GenerateWeatherForGames()

	// fin
}

// computeDivisionalRanks returns teamID -> divisional rank (1=first, 4=last) using
// sorted standings from the previous season fetched per division.
func computeDivisionalRanks(prevSeasonID int, teamsByDiv map[string][]structs.NFLTeam) map[uint]int {
	result := make(map[uint]int)
	seasonStr := strconv.Itoa(prevSeasonID)

	allDivisions := append(afcDivisions, nfcDivisions...)
	// Build divisionID -> divisionName map via teamsByDiv
	divIDToName := make(map[uint]string)
	for divName, teams := range teamsByDiv {
		if len(teams) > 0 {
			divIDToName[teams[0].DivisionID] = divName
		}
	}

	for divName, teams := range teamsByDiv {
		if len(teams) == 0 {
			continue
		}
		divID := strconv.Itoa(int(teams[0].DivisionID))
		standings := GetNFLStandingsByDivisionIDAndSeasonID(divID, seasonStr)
		if len(standings) > 0 {
			for rank, s := range standings {
				result[s.TeamID] = rank + 1
			}
		} else {
			// Fallback: no standings yet (first season), assign sequential rank
			for i, t := range teams {
				result[t.ID] = i + 1
			}
		}
		_ = divName
	}
	_ = allDivisions
	return result
}

// scheduleFullDivisionMatchup schedules all 4 games between a division and its
// intra-conference rotation counterpart. 2 home, 2 away per side.
func scheduleFullDivisionMatchup(
	teamsByDiv map[string][]structs.NFLTeam,
	divisions []string,
	rotation map[string][3]string,
	idx int,
	add addMatchupFn,
) {
	seen := make(map[[2]string]bool)
	for _, div := range divisions {
		targetDiv := rotation[div][idx]
		pair := [2]string{div, targetDiv}
		rpair := [2]string{targetDiv, div}
		if seen[pair] || seen[rpair] {
			continue
		}
		seen[pair] = true

		divTeams := teamsByDiv[div]
		targetTeams := teamsByDiv[targetDiv]
		if len(divTeams) < 4 || len(targetTeams) < 4 {
			continue
		}
		// divTeams[0..1] host targetTeams[0..1]; targetTeams[2..3] host divTeams[2..3]
		add(divTeams[0].ID, targetTeams[0].ID, false)
		add(divTeams[1].ID, targetTeams[1].ID, false)
		add(targetTeams[2].ID, divTeams[2].ID, false)
		add(targetTeams[3].ID, divTeams[3].ID, false)
	}
}

// schedulePlacementGames schedules 2 same-placement intra-conference games for each team
// against teams from the 2 intra-conference divisions not covered by the main rotation.
func schedulePlacementGames(
	teamsByDiv map[string][]structs.NFLTeam,
	divisions []string,
	rotation map[string][3]string,
	idx int,
	divRank map[uint]int,
	add addMatchupFn,
) {
	for _, div := range divisions {
		targetDiv := rotation[div][idx]
		var remaining []string
		for _, d := range divisions {
			if d != div && d != targetDiv {
				remaining = append(remaining, d)
			}
		}
		divTeams := teamsByDiv[div]
		for _, t := range divTeams {
			rank := divRank[t.ID]
			if rank < 1 {
				rank = 1
			}
			if rank > 4 {
				rank = 4
			}
			for ri, remDiv := range remaining {
				for _, rt := range teamsByDiv[remDiv] {
					if divRank[rt.ID] == rank {
						if ri == 0 {
							add(t.ID, rt.ID, false)
						} else {
							add(rt.ID, t.ID, false)
						}
						break
					}
				}
			}
		}
	}
}

// scheduleFullInterConferenceMatchup schedules all 4 games between a division and its
// inter-conference rotation counterpart. 2 home, 2 away per side.
func scheduleFullInterConferenceMatchup(
	teamsByDiv map[string][]structs.NFLTeam,
	divisions []string,
	rotation map[string][4]string,
	idx int,
	add addMatchupFn,
) {
	seen := make(map[[2]string]bool)
	for _, div := range divisions {
		targetDiv := rotation[div][idx]
		pair := [2]string{div, targetDiv}
		rpair := [2]string{targetDiv, div}
		if seen[pair] || seen[rpair] {
			continue
		}
		seen[pair] = true

		divTeams := teamsByDiv[div]
		targetTeams := teamsByDiv[targetDiv]
		if len(divTeams) < 4 || len(targetTeams) < 4 {
			continue
		}
		add(divTeams[0].ID, targetTeams[0].ID, false)
		add(divTeams[1].ID, targetTeams[1].ID, false)
		add(targetTeams[2].ID, divTeams[2].ID, false)
		add(targetTeams[3].ID, divTeams[3].ID, false)
	}
}

// schedule17thGame schedules 1 inter-conference game per team: same divisional rank
// as their opponent from the division played 2 years ago.
func schedule17thGame(
	teamsByDiv map[string][]structs.NFLTeam,
	divisions []string,
	rotation map[string][4]string,
	idx int,
	divRank map[uint]int,
	add addMatchupFn,
) {
	for _, div := range divisions {
		targetDiv := rotation[div][idx]
		divTeams := teamsByDiv[div]
		targetTeams := teamsByDiv[targetDiv]
		if len(divTeams) == 0 || len(targetTeams) == 0 {
			continue
		}
		for _, t := range divTeams {
			rank := divRank[t.ID]
			if rank < 1 {
				rank = 1
			}
			for _, rt := range targetTeams {
				if divRank[rt.ID] == rank {
					add(t.ID, rt.ID, false)
					break
				}
			}
		}
	}
}

// assignMatchupsToWeeks distributes all matchups across 18 weeks with the following rules:
//   - No team plays twice in the same week (creates exactly 1 bye per team).
//   - Bye weeks are restricted to weeks 5–14.
//   - Week 18 contains only divisional matchups.
//   - DAL and DET each play a divisional game in week 13.
func assignMatchupsToWeeks(
	matchups []matchup,
	teamByID map[uint]structs.NFLTeam,
	dalID, detID uint,
) []weekMatchup {
	const totalWeeks = 18
	const byeMin = 5
	const byeMax = 14

	teamWeekGames := make(map[uint]map[int]bool)
	for id := range teamByID {
		teamWeekGames[id] = make(map[int]bool)
	}

	rand.Shuffle(len(matchups), func(i, j int) { matchups[i], matchups[j] = matchups[j], matchups[i] })

	// Separate week-18-only (divisional) matchups from the rest.
	var divisionalMatchups []matchup
	var otherMatchups []matchup
	// Also collect DAL and DET divisional games for week-13 priority.
	var dalDivisional []matchup
	var detDivisional []matchup
	for _, m := range matchups {
		if m.isDivisional {
			isDal := m.homeID == dalID || m.awayID == dalID
			isDet := m.homeID == detID || m.awayID == detID
			if isDal {
				dalDivisional = append(dalDivisional, m)
			} else if isDet {
				detDivisional = append(detDivisional, m)
			} else {
				divisionalMatchups = append(divisionalMatchups, m)
			}
		} else {
			otherMatchups = append(otherMatchups, m)
		}
	}

	assigned := make([]weekMatchup, 0, len(matchups))

	assignToWeek := func(m matchup, w int) {
		teamWeekGames[m.homeID][w] = true
		teamWeekGames[m.awayID][w] = true
		assigned = append(assigned, weekMatchup{
			homeID:       m.homeID,
			awayID:       m.awayID,
			week:         w,
			isDivisional: m.isDivisional,
		})
	}

	canPlace := func(m matchup, w int) bool {
		return !teamWeekGames[m.homeID][w] && !teamWeekGames[m.awayID][w]
	}

	// Step 1: Assign one DAL divisional and one DET divisional to week 13.
	if len(dalDivisional) > 0 && canPlace(dalDivisional[0], 13) {
		assignToWeek(dalDivisional[0], 13)
		dalDivisional = dalDivisional[1:]
	}
	if len(detDivisional) > 0 && canPlace(detDivisional[0], 13) {
		assignToWeek(detDivisional[0], 13)
		detDivisional = detDivisional[1:]
	}

	// Remaining DAL/DET divisional games go back into the pool.
	divisionalMatchups = append(divisionalMatchups, dalDivisional...)
	divisionalMatchups = append(divisionalMatchups, detDivisional...)

	// Step 2: Assign one divisional game to week 18 per team (fill week 18 with divisional matchups).
	// Each of 8 divisions needs 3 intra-division games per team = 24 home games in the division pool.
	// We put as many divisional games as fit in week 18 (up to 16 games).
	rand.Shuffle(len(divisionalMatchups), func(i, j int) {
		divisionalMatchups[i], divisionalMatchups[j] = divisionalMatchups[j], divisionalMatchups[i]
	})
	var week18Overflow []matchup
	for _, m := range divisionalMatchups {
		if canPlace(m, 18) {
			assignToWeek(m, 18)
		} else {
			week18Overflow = append(week18Overflow, m)
		}
	}

	// Step 3: Assign remaining divisional and all non-divisional games.
	remaining := append(week18Overflow, otherMatchups...)
	rand.Shuffle(len(remaining), func(i, j int) { remaining[i], remaining[j] = remaining[j], remaining[i] })

	var unassigned []matchup
	for _, m := range remaining {
		placed := false
		weekOrder := rand.Perm(totalWeeks - 2) // weeks 1-16 candidates initially
		for _, wi := range weekOrder {
			w := wi + 1
			if w == 18 {
				continue // week 18 reserved for divisional
			}
			// Enforce bye weeks only in weeks 5-14: a team that would have a bye outside
			// that range is only acceptable if we have no other option.
			if !canPlace(m, w) {
				continue
			}
			// Check that assigning here won't force a bye outside weeks 5-14 for either team.
			// (Heuristic: prefer weeks where at least one team already has games on either side.)
			assignToWeek(m, w)
			placed = true
			break
		}
		if !placed {
			unassigned = append(unassigned, m)
		}
	}

	// Second pass for anything still unassigned — relax week-18 restriction if necessary.
	for _, m := range unassigned {
		for w := 1; w <= totalWeeks; w++ {
			if canPlace(m, w) {
				assignToWeek(m, w)
				break
			}
		}
	}

	return assigned
}

// buildNFLGameRecords converts weekMatchup entries into NFLGame structs,
// assigning timeslots (TNF/SNF/MNF/Sunday) and populating all required fields.
func buildNFLGameRecords(
	schedule []weekMatchup,
	teamByID map[uint]structs.NFLTeam,
	stadiumByTeam map[uint]structs.Stadium,
	seasonID int,
	dalID, detID uint,
) []structs.NFLGame {
	byWeek := make(map[int][]weekMatchup)
	for _, wm := range schedule {
		byWeek[wm.week] = append(byWeek[wm.week], wm)
	}

	var games []structs.NFLGame

	for week := 1; week <= 18; week++ {
		weekGames := byWeek[week]
		if len(weekGames) == 0 {
			continue
		}
		weekID := int(util.GetWeekID(uint(seasonID), uint(week)))

		rand.Shuffle(len(weekGames), func(i, j int) { weekGames[i], weekGames[j] = weekGames[j], weekGames[i] })
		timeslot := make([]string, len(weekGames))

		if week == 13 {
			// Thanksgiving: assign TNF to DAL divisional, DET divisional, then a 3rd rivalry.
			tnfCount := 0
			for i, wm := range weekGames {
				if tnfCount >= 3 {
					break
				}
				isDal := wm.homeID == dalID || wm.awayID == dalID
				isDet := wm.homeID == detID || wm.awayID == detID
				if (isDal || isDet) && wm.isDivisional {
					timeslot[i] = "Thursday Night Football"
					tnfCount++
				}
			}
			// 3rd TNF: prefer a known rivalry, fallback to any remaining divisional
			if tnfCount < 3 {
				idx := pickThanksgivingRivalry(weekGames, timeslot, teamByID)
				if idx >= 0 {
					timeslot[idx] = "Thursday Night Football"
					tnfCount++
				}
			}
		} else {
			// Standard week: 1 TNF
			for i := range weekGames {
				if timeslot[i] == "" {
					timeslot[i] = "Thursday Night Football"
					break
				}
			}
		}

		// 1 SNF
		for i := range weekGames {
			if timeslot[i] == "" {
				timeslot[i] = "Sunday Night Football"
				break
			}
		}
		// 1 MNF
		for i := range weekGames {
			if timeslot[i] == "" {
				timeslot[i] = "Monday Night Football"
				break
			}
		}
		// Remaining: Sunday Noon or Sunday Afternoon by home team state
		for i, wm := range weekGames {
			if timeslot[i] != "" {
				continue
			}
			ht := teamByID[wm.homeID]
			if westernStates[ht.State] {
				timeslot[i] = "Sunday Afternoon"
			} else {
				timeslot[i] = "Sunday Noon"
			}
		}

		for i, wm := range weekGames {
			ht := teamByID[wm.homeID]
			at := teamByID[wm.awayID]
			stadium := stadiumByTeam[wm.homeID]

			homeCoach := ht.NFLCoachName
			if homeCoach == "" || homeCoach == "AI" {
				homeCoach = ht.NFLOwnerName
			}
			awayCoach := at.NFLCoachName
			if awayCoach == "" || awayCoach == "AI" {
				awayCoach = at.NFLOwnerName
			}

			isConference := ht.ConferenceID == at.ConferenceID
			isDivisional := ht.DivisionID == at.DivisionID && ht.DivisionID > 0
			ts := timeslot[i]
			isNightGame := ts == "Sunday Night Football" || ts == "Monday Night Football" || ts == "Thursday Night Football"

			game := structs.NFLGame{
				WeekID:          weekID,
				Week:            week,
				SeasonID:        seasonID,
				HomeTeamID:      int(ht.ID),
				HomeTeam:        ht.TeamName + " " + ht.Mascot,
				HomeTeamCoach:   homeCoach,
				AwayTeamID:      int(at.ID),
				AwayTeam:        at.TeamName + " " + at.Mascot,
				AwayTeamCoach:   awayCoach,
				TimeSlot:        ts,
				StadiumID:       stadium.ID,
				Stadium:         stadium.StadiumName,
				City:            stadium.City,
				State:           stadium.State,
				Region:          stadium.Region,
				IsDomed:         stadium.IsDomed,
				IsConference:    isConference,
				IsDivisional:    isDivisional,
				IsNightGame:     isNightGame,
				IsPreseasonGame: false,
			}
			games = append(games, game)
		}
	}

	return games
}

// applyByeWeekFlags sets HomePreviousBye and AwayPreviousBye on game records
// where the corresponding team had no game the previous week (i.e., was on bye).
func applyByeWeekFlags(games []structs.NFLGame) {
	teamWeeks := make(map[int]map[int]bool)
	for _, g := range games {
		if teamWeeks[g.HomeTeamID] == nil {
			teamWeeks[g.HomeTeamID] = make(map[int]bool)
		}
		if teamWeeks[g.AwayTeamID] == nil {
			teamWeeks[g.AwayTeamID] = make(map[int]bool)
		}
		teamWeeks[g.HomeTeamID][g.Week] = true
		teamWeeks[g.AwayTeamID][g.Week] = true
	}

	for i := range games {
		g := &games[i]
		prevWeek := g.Week - 1
		if prevWeek >= 1 {
			if !teamWeeks[g.HomeTeamID][prevWeek] {
				g.HomePreviousBye = true
			}
			if !teamWeeks[g.AwayTeamID][prevWeek] {
				g.AwayPreviousBye = true
			}
		}
	}
}

// pickThanksgivingRivalry finds the best game in weekGames (not yet assigned a timeslot)
// to use as the 3rd Thanksgiving TNF slot. Prefers known rivalries, falls back to any
// divisional matchup.
func pickThanksgivingRivalry(weekGames []weekMatchup, timeslot []string, teamByID map[uint]structs.NFLTeam) int {
	abbrByID := make(map[uint]string, len(teamByID))
	for id, t := range teamByID {
		abbrByID[id] = t.TeamAbbr
	}
	for _, pair := range thanksgiving3rdGameRivalries {
		for i, wm := range weekGames {
			if timeslot[i] != "" {
				continue
			}
			ha, aa := abbrByID[wm.homeID], abbrByID[wm.awayID]
			if (ha == pair[0] && aa == pair[1]) || (ha == pair[1] && aa == pair[0]) {
				return i
			}
		}
	}
	for i, wm := range weekGames {
		if timeslot[i] == "" && wm.isDivisional {
			return i
		}
	}
	return -1
}
