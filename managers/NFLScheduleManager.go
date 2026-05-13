package managers

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
)

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

// nflPair is an ordered (homeID, awayID) key for deduplication.
type nflPair [2]uint

// GenerateSimNFLSchedule generates a full 17-game regular season (18 weeks, 1 bye)
// for the SimNFL league and saves all game records to the database.
func GenerateSimNFLSchedule(isTest bool) {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID := ts.NFLSeasonID
	prevSeasonID := seasonID - 1

	// Collect teams and build lookup structures.
	allTeams := GetAllNFLTeams()
	teamByID := make(map[uint]structs.NFLTeam, len(allTeams))

	// Build teamsByDiv keyed by DivisionID. Division IDs are globally unique
	// (1–8 across all NFL divisions), so no composite key is needed.
	teamsByDiv := make(map[uint][]structs.NFLTeam)
	for _, t := range allTeams {
		teamByID[t.ID] = t
		key := t.DivisionID
		teamsByDiv[key] = append(teamsByDiv[key], t)
	}

	// Build division lists per conference, sorted for canonical ordering.
	// Also derive the rotation maps dynamically so they are consistent with
	// whatever ConferenceID / DivisionID values are in the database.
	conf0Divs, conf1Divs := buildDivisionLists(teamsByDiv)
	intraRot0 := buildIntraRotation(conf0Divs)
	intraRot1 := buildIntraRotation(conf1Divs)
	interRot0 := buildInterRotation(conf0Divs, conf1Divs)

	// Build stadium map: homeTeamID -> Stadium
	stadiums := GetAllStadiums()
	stadiumByTeam := make(map[uint]structs.Stadium, len(stadiums))
	for _, s := range stadiums {
		if s.LeagueName != "NFL" {
			continue
		}
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

	// Compute home game count per team from previous season to balance 17th game H/A.
	prevHomeCount := computePrevSeasonHomeCount(prevSeasonID, allTeams)

	// Accumulate all matchups.  Key is exact (homeID, awayID) — reversals are distinct games.
	scheduledPairs := make(map[nflPair]bool)
	var allMatchups []matchup

	// currentHomeCount tracks how many home games each team has been assigned so far.
	// Steps 3 and 5 use this shared counter to guarantee exactly 1H+1A per team across
	// both steps combined, achieving the required 9H/8A or 8H/9A final balance.
	currentHomeCount := make(map[uint]int, len(allTeams))

	addMatchup := func(homeID, awayID uint, isDivisional bool) {
		key := nflPair{homeID, awayID}
		if scheduledPairs[key] {
			return
		}
		scheduledPairs[key] = true
		currentHomeCount[homeID]++
		allMatchups = append(allMatchups, matchup{homeID: homeID, awayID: awayID, isDivisional: isDivisional})
	}

	// 1. Divisional games: each pair plays twice — once at each stadium (6 games per team).
	// Only iterate divisions with exactly 4 teams to skip any stray/empty buckets.
	for _, teams := range teamsByDiv {
		if len(teams) != 4 {
			continue
		}
		for i := 0; i < 4; i++ {
			for j := i + 1; j < 4; j++ {
				addMatchup(teams[i].ID, teams[j].ID, true)
				addMatchup(teams[j].ID, teams[i].ID, true)
			}
		}
	}
	fmt.Printf("[DEBUG] After step 1 (divisional): %d matchups\n", len(allMatchups))

	// 2. Intra-conference full-division rotation (4 games per team: 2H, 2A).
	scheduleFullDivisionMatchup(teamsByDiv, conf0Divs, intraRot0, intraIdx, addMatchup)
	scheduleFullDivisionMatchup(teamsByDiv, conf1Divs, intraRot1, intraIdx, addMatchup)
	fmt.Printf("[DEBUG] After step 2 (intra-conf rotation): %d matchups\n", len(allMatchups))

	// 3. Two same-placement intra-conference games from the 2 non-rotation divisions (2 games per team: 1H, 1A).
	// currentHomeCount is passed so placement H/A decisions are visible to step 5.
	schedulePlacementGames(teamsByDiv, conf0Divs, intraRot0, intraIdx, divRank, currentHomeCount, addMatchup)
	schedulePlacementGames(teamsByDiv, conf1Divs, intraRot1, intraIdx, divRank, currentHomeCount, addMatchup)
	fmt.Printf("[DEBUG] After step 3 (placement games): %d matchups\n", len(allMatchups))

	// 4. Inter-conference full-division rotation (4 games per team: 2H, 2A).
	// Enumerating conf0 divisions covers both conferences since each cross-conf game
	// is recorded once (homeID, awayID) and the reverse is also added.
	scheduleFullInterConferenceMatchup(teamsByDiv, conf0Divs, interRot0, interIdx, addMatchup)
	fmt.Printf("[DEBUG] After step 4 (inter-conf rotation): %d matchups\n", len(allMatchups))

	// 5. 17th game: same-placement opponent from the inter-conference division played 2 years ago.
	// Uses currentHomeCount (updated by all prior steps) to guarantee each team's 17th game
	// corrects any H/A imbalance: a team at 8H hosts, a team at 9H is away.
	schedule17thGame(teamsByDiv, conf0Divs, interRot0, inter17Idx, divRank, prevHomeCount, currentHomeCount, scheduledPairs, addMatchup)
	fmt.Printf("[DEBUG] After step 5 (17th game): %d matchups (expected 272)\n", len(allMatchups))

	// Validate: every team must have exactly 17 matchups.
	validateScheduleGameCount(allMatchups, teamByID)

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
	applySimNFLByeWeekFlags(finalGames)

	// Validate home/away balance: 16 teams should have 9H/8A, 16 teams 8H/9A.
	validateHomeAwayBalance(finalGames, teamByID)

	if isTest {
		exportNFLScheduleToCSV(finalGames, teamByID, seasonID)
	} else {
		repository.CreateNFLGameRecordsBatch(db, finalGames, 250)
		GenerateWeatherForGames()
	}
}

// buildDivisionLists groups teamsByDiv keys (DivisionIDs) by conference and returns
// two sorted slices — conf0 for the lower ConferenceID, conf1 for the higher.
func buildDivisionLists(teamsByDiv map[uint][]structs.NFLTeam) (conf0, conf1 []uint) {
	confDivs := make(map[uint][]uint)
	for key, teams := range teamsByDiv {
		if len(teams) == 0 {
			continue
		}
		confID := teams[0].ConferenceID
		confDivs[confID] = append(confDivs[confID], key)
	}
	// Sort conference IDs for canonical conf0 / conf1 assignment.
	var confIDs []uint
	for id := range confDivs {
		confIDs = append(confIDs, id)
	}
	sort.Slice(confIDs, func(i, j int) bool { return confIDs[i] < confIDs[j] })
	// Sort division IDs within each conference for canonical slot ordering.
	for _, id := range confIDs {
		divList := confDivs[id]
		sort.Slice(divList, func(i, j int) bool { return divList[i] < divList[j] })
		confDivs[id] = divList
	}
	if len(confIDs) > 0 {
		conf0 = confDivs[confIDs[0]]
	}
	if len(confIDs) > 1 {
		conf1 = confDivs[confIDs[1]]
	}
	return conf0, conf1
}

// buildIntraRotation generates a 3-year round-robin rotation map for 4 same-conference
// divisions. Each division key maps to the 3 division keys it plays (one per year).
// The round-robin schedule is:
//   - Round 0: (0,1) and (2,3)
//   - Round 1: (0,2) and (1,3)
//   - Round 2: (0,3) and (1,2)
func buildIntraRotation(divs []uint) map[uint][3]uint {
	result := make(map[uint][3]uint, len(divs))
	if len(divs) < 4 {
		return result
	}
	// Pairs per round (indices into divs).
	rounds := [3][2][2]int{
		{{0, 1}, {2, 3}},
		{{0, 2}, {1, 3}},
		{{0, 3}, {1, 2}},
	}
	// Precompute: for each div index i, its opponent in each round.
	opponentIdx := [4][3]int{}
	for r, roundPairs := range rounds {
		for _, pair := range roundPairs {
			opponentIdx[pair[0]][r] = pair[1]
			opponentIdx[pair[1]][r] = pair[0]
		}
	}
	for i, div := range divs {
		result[div] = [3]uint{
			divs[opponentIdx[i][0]],
			divs[opponentIdx[i][1]],
			divs[opponentIdx[i][2]],
		}
	}
	return result
}

// buildInterRotation generates a 4-year inter-conference rotation map.
// Each division in conf0 maps to the 4 conf1 divisions it plays (one per year).
// The pattern ensures each conf0 division plays each conf1 division exactly once
// across the 4-year cycle.
func buildInterRotation(conf0, conf1 []uint) map[uint][4]uint {
	result := make(map[uint][4]uint, len(conf0))
	if len(conf0) < 4 || len(conf1) < 4 {
		return result
	}
	// For conf0 division at index i, in year j it plays conf1 division at index
	// basePattern[(i+j) % 4], where basePattern = [1,2,0,3].
	// This is derived from the standard NFL inter-conference rotation and ensures
	// all 4×4 matchup pairs occur across the cycle with no repetition.
	basePattern := [4]int{1, 2, 0, 3}
	for i, div := range conf0 {
		var rot [4]uint
		for j := 0; j < 4; j++ {
			rot[j] = conf1[basePattern[(i+j)%4]]
		}
		result[div] = rot
	}
	return result
}

// computePrevSeasonHomeCount returns a map of teamID -> number of home games played last season.
func computePrevSeasonHomeCount(prevSeasonID int, allTeams []structs.NFLTeam) map[uint]int {
	result := make(map[uint]int, len(allTeams))
	prevGames := GetNFLGamesBySeasonID(strconv.Itoa(prevSeasonID))
	for _, g := range prevGames {
		if !g.IsPreseasonGame {
			result[uint(g.HomeTeamID)]++
		}
	}
	return result
}

// validateScheduleGameCount logs a warning for any team that doesn't have exactly 17 games.
func validateScheduleGameCount(matchups []matchup, teamByID map[uint]structs.NFLTeam) {
	counts := make(map[uint]int, len(teamByID))
	for _, m := range matchups {
		counts[m.homeID]++
		counts[m.awayID]++
	}
	for id, t := range teamByID {
		if counts[id] != 17 {
			fmt.Printf("[WARN] NFL Schedule: %s %s has %d games (expected 17)\n", t.TeamName, t.Mascot, counts[id])
		}
	}
}

// validateHomeAwayBalance checks that exactly 16 teams have 9 home / 8 away games
// and the other 16 teams have 8 home / 9 away games. Any deviations are logged as warnings.
func validateHomeAwayBalance(games []structs.NFLGame, teamByID map[uint]structs.NFLTeam) {
	homeCount := make(map[uint]int, len(teamByID))
	awayCount := make(map[uint]int, len(teamByID))
	for _, g := range games {
		homeCount[uint(g.HomeTeamID)]++
		awayCount[uint(g.AwayTeamID)]++
	}

	nine8 := 0  // teams with 9H / 8A
	eight9 := 0 // teams with 8H / 9A
	for id, t := range teamByID {
		h, a := homeCount[id], awayCount[id]
		switch {
		case h == 9 && a == 8:
			nine8++
		case h == 8 && a == 9:
			eight9++
		default:
			fmt.Printf("[WARN] H/A imbalance: %s %s — %dH / %dA (expected 9/8 or 8/9)\n",
				t.TeamName, t.Mascot, h, a)
		}
	}
	fmt.Printf("[INFO] Home/away balance: %d teams at 9H/8A, %d teams at 8H/9A (need 16 each)\n", nine8, eight9)
	if nine8 != 16 || eight9 != 16 {
		fmt.Printf("[WARN] Home/away split is off — expected 16/16, got %d/%d\n", nine8, eight9)
	}
}

// exportNFLScheduleToCSV writes the generated schedule to a CSV file for inspection.
func exportNFLScheduleToCSV(games []structs.NFLGame, teamByID map[uint]structs.NFLTeam, seasonID int) {
	filename := fmt.Sprintf("nfl_schedule_season_%d.csv", seasonID)
	f, err := os.Create(filename)
	if err != nil {
		fmt.Printf("[ERROR] Could not create CSV file: %v\n", err)
		return
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	// Header
	_ = w.Write([]string{
		"Week", "HomeTeamID", "HomeTeam", "AwayTeamID", "AwayTeam",
		"TimeSlot", "Stadium", "City", "State",
		"IsConference", "IsDivisional", "HomePreviousBye", "AwayPreviousBye",
	})

	// Tally games per team for summary
	homeCount := make(map[uint]int)
	awayCount := make(map[uint]int)

	for _, g := range games {
		homeCount[uint(g.HomeTeamID)]++
		awayCount[uint(g.AwayTeamID)]++

		_ = w.Write([]string{
			strconv.Itoa(g.Week),
			strconv.Itoa(g.HomeTeamID),
			g.HomeTeam,
			strconv.Itoa(g.AwayTeamID),
			g.AwayTeam,
			g.TimeSlot,
			g.Stadium,
			g.City,
			g.State,
			strconv.FormatBool(g.IsConference),
			strconv.FormatBool(g.IsDivisional),
			strconv.FormatBool(g.HomePreviousBye),
			strconv.FormatBool(g.AwayPreviousBye),
		})
	}

	// Append per-team summary rows after a blank line
	_ = w.Write([]string{})
	_ = w.Write([]string{"TeamID", "TeamName", "HomeGames", "AwayGames", "TotalGames"})
	for id, t := range teamByID {
		h := homeCount[id]
		a := awayCount[id]
		_ = w.Write([]string{
			strconv.Itoa(int(id)),
			t.TeamName + " " + t.Mascot,
			strconv.Itoa(h),
			strconv.Itoa(a),
			strconv.Itoa(h + a),
		})
	}

	fmt.Printf("[INFO] NFL schedule exported to %s (%d games)\n", filename, len(games))
}

// computeDivisionalRanks returns teamID -> divisional rank (1=first, 4=last) using
// sorted standings from the previous season fetched per division.
func computeDivisionalRanks(prevSeasonID int, teamsByDiv map[uint][]structs.NFLTeam) map[uint]int {
	result := make(map[uint]int)
	seasonStr := strconv.Itoa(prevSeasonID)

	for _, teams := range teamsByDiv {
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
	}
	return result
}

// scheduleFullDivisionMatchup schedules all 4 games between a division and its
// intra-conference rotation counterpart.
// Each team plays all 4 opponents from the target division, hosting 2 and visiting 2.
// This yields 16 total directional matchups between the two 4-team divisions.
func scheduleFullDivisionMatchup(
	teamsByDiv map[uint][]structs.NFLTeam,
	divisions []uint,
	rotation map[uint][3]uint,
	idx int,
	add addMatchupFn,
) {
	seen := make(map[[2]uint]bool)
	for _, div := range divisions {
		targetDiv := rotation[div][idx]
		pair := [2]uint{div, targetDiv}
		rpair := [2]uint{targetDiv, div}
		if seen[pair] || seen[rpair] {
			continue
		}
		seen[pair] = true

		divTeams := teamsByDiv[div]
		targetTeams := teamsByDiv[targetDiv]
		if len(divTeams) < 4 || len(targetTeams) < 4 {
			continue
		}
		// Each team in divTeams hosts 2 opponents and visits 2.
		// Pattern: divTeams[i] hosts targetTeams[i] and targetTeams[(i+2)%4]
		//          divTeams[i] visits targetTeams[(i+1)%4] and targetTeams[(i+3)%4]
		for i := 0; i < 4; i++ {
			add(divTeams[i].ID, targetTeams[i].ID, false)
			add(divTeams[i].ID, targetTeams[(i+2)%4].ID, false)
			add(targetTeams[(i+1)%4].ID, divTeams[i].ID, false)
			add(targetTeams[(i+3)%4].ID, divTeams[i].ID, false)
		}
	}
}

// schedulePlacementGames schedules 2 same-placement intra-conference games for each team
// against teams from the 2 intra-conference divisions not covered by the main rotation.
// Each pair of divisions is processed once. Teams are paired by sorted rank index so that
// the Nth-ranked team in one division plays the Nth-ranked team in the other.
func schedulePlacementGames(
	teamsByDiv map[uint][]structs.NFLTeam,
	divisions []uint,
	rotation map[uint][3]uint,
	idx int,
	divRank map[uint]int,
	currentHomeCount map[uint]int,
	add addMatchupFn,
) {
	type divPair [2]uint
	processedPairs := make(map[divPair]bool)

	for _, div := range divisions {
		targetDiv := rotation[div][idx]
		var remaining []uint
		for _, d := range divisions {
			if d != div && d != targetDiv {
				remaining = append(remaining, d)
			}
		}
		for _, remDiv := range remaining {
			// Process each canonical (div, remDiv) pair exactly once.
			canonical := divPair{div, remDiv}
			if div > remDiv {
				canonical = divPair{remDiv, div}
			}
			if processedPairs[canonical] {
				continue
			}
			processedPairs[canonical] = true

			// Sort both divisions by rank so index i pairs same-ranked opponents.
			sortedDiv := sortTeamsByRank(teamsByDiv[div], divRank)
			sortedRem := sortTeamsByRank(teamsByDiv[remDiv], divRank)
			for i := 0; i < len(sortedDiv) && i < len(sortedRem); i++ {
				t, rt := sortedDiv[i], sortedRem[i]
				// Use the global home count so step 5 sees the correct running total.
				if currentHomeCount[t.ID] <= currentHomeCount[rt.ID] {
					add(t.ID, rt.ID, false)
				} else {
					add(rt.ID, t.ID, false)
				}
			}
		}
	}
}

// scheduleFullInterConferenceMatchup schedules all 4 games between a division and its
// inter-conference rotation counterpart.
// Each team plays all 4 opponents from the target division, hosting 2 and visiting 2.
func scheduleFullInterConferenceMatchup(
	teamsByDiv map[uint][]structs.NFLTeam,
	divisions []uint,
	rotation map[uint][4]uint,
	idx int,
	add addMatchupFn,
) {
	seen := make(map[[2]uint]bool)
	for _, div := range divisions {
		targetDiv := rotation[div][idx]
		pair := [2]uint{div, targetDiv}
		rpair := [2]uint{targetDiv, div}
		if seen[pair] || seen[rpair] {
			continue
		}
		seen[pair] = true

		divTeams := teamsByDiv[div]
		targetTeams := teamsByDiv[targetDiv]
		if len(divTeams) < 4 || len(targetTeams) < 4 {
			continue
		}
		// Same pattern as intra: each team hosts 2, visits 2.
		for i := 0; i < 4; i++ {
			add(divTeams[i].ID, targetTeams[i].ID, false)
			add(divTeams[i].ID, targetTeams[(i+2)%4].ID, false)
			add(targetTeams[(i+1)%4].ID, divTeams[i].ID, false)
			add(targetTeams[(i+3)%4].ID, divTeams[i].ID, false)
		}
	}
}

// schedule17thGame schedules 1 inter-conference game per team: same divisional rank
// as their opponent from the division played 2 years ago.
// Home/away assignment is determined by previous season home game count: a team that
// had more home games last season (≥9) will be the away team for the 17th game.
// Teams are paired by sorted rank index so missing divRank data never produces 0 games.
func schedule17thGame(
	teamsByDiv map[uint][]structs.NFLTeam,
	divisions []uint,
	rotation map[uint][4]uint,
	idx int,
	divRank map[uint]int,
	prevHomeCount map[uint]int,
	currentHomeCount map[uint]int,
	scheduledPairs map[nflPair]bool,
	add addMatchupFn,
) {
	for _, div := range divisions {
		targetDiv := rotation[div][idx]
		if len(teamsByDiv[div]) == 0 || len(teamsByDiv[targetDiv]) == 0 {
			continue
		}
		// Sort both divisions by rank so index i pairs same-ranked opponents.
		sortedDiv := sortTeamsByRank(teamsByDiv[div], divRank)
		sortedTarget := sortTeamsByRank(teamsByDiv[targetDiv], divRank)
		for i := 0; i < len(sortedDiv) && i < len(sortedTarget); i++ {
			t, rt := sortedDiv[i], sortedTarget[i]
			// Primary: cross-season balance. The team with more home games last
			// season should be the away team this season.
			homeID, awayID := t.ID, rt.ID
			if prevHomeCount[t.ID] > prevHomeCount[rt.ID] {
				// t had more home games last season → t travels
				homeID, awayID = rt.ID, t.ID
			} else if prevHomeCount[rt.ID] > prevHomeCount[t.ID] {
				// rt had more home games last season → rt travels
				homeID, awayID = t.ID, rt.ID
			}
			// else: tied prev season, keep t as home by default.

			// Safety override: if currentHomeCount demands a flip to hit the required
			// 8H or 9H season target, honour it. This preserves the H/A balance
			// guarantee regardless of what last season's counts look like.
			if currentHomeCount[awayID] < currentHomeCount[homeID] {
				homeID, awayID = awayID, homeID
			}

			// If this exact direction is already scheduled, flip.
			if scheduledPairs[nflPair{homeID, awayID}] {
				homeID, awayID = awayID, homeID
			}
			add(homeID, awayID, false)
		}
	}
}

// sortTeamsByRank returns a copy of teams sorted by divisional rank (ascending),
// with team ID as a stable tiebreaker when ranks are equal or missing.
func sortTeamsByRank(teams []structs.NFLTeam, divRank map[uint]int) []structs.NFLTeam {
	sorted := make([]structs.NFLTeam, len(teams))
	copy(sorted, teams)
	sort.Slice(sorted, func(i, j int) bool {
		ri, rj := divRank[sorted[i].ID], divRank[sorted[j].ID]
		if ri != rj {
			return ri < rj
		}
		return sorted[i].ID < sorted[j].ID
	})
	return sorted
}

// assignMatchupsToWeeks distributes all matchups across 18 weeks with the following rules:
//   - No team plays twice in the same week (creates exactly 1 bye per team).
//   - Bye weeks are restricted to weeks 5–14.
//   - Week 18 contains only divisional matchups.
//   - DAL and DET each play a divisional game in week 13.
//
// If any matchup cannot be placed, the function retries from scratch with a fresh
// shuffle, up to maxAttempts times.
func assignMatchupsToWeeks(
	matchups []matchup,
	teamByID map[uint]structs.NFLTeam,
	dalID, detID uint,
) []weekMatchup {
	const maxAttempts = 30000

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		assigned, failed := attemptWeekAssignment(matchups, teamByID, dalID, detID)
		if len(failed) == 0 {
			if attempt > 1 {
				fmt.Printf("[INFO] Week assignment succeeded on attempt %d\n", attempt)
			}
			return assigned
		}
		fmt.Printf("[INFO] Week assignment attempt %d: %d matchup(s) unplaced, retrying...\n", attempt, len(failed))
	}

	// All retries exhausted — return the last attempt's result, logging the failures.
	fmt.Printf("[WARN] Week assignment failed after %d attempts; some matchups will be missing\n", maxAttempts)
	assigned, failed := attemptWeekAssignment(matchups, teamByID, dalID, detID)
	for _, m := range failed {
		fmt.Printf("[WARN] Could not place matchup homeID=%d awayID=%d\n", m.homeID, m.awayID)
	}
	return assigned
}

// attemptWeekAssignment performs a single attempt at distributing matchups across weeks.
// It returns the successfully assigned games and any matchups that could not be placed.
func attemptWeekAssignment(
	matchups []matchup,
	teamByID map[uint]structs.NFLTeam,
	dalID, detID uint,
) (assigned []weekMatchup, failed []matchup) {
	const totalWeeks = 18

	teamWeekGames := make(map[uint]map[int]bool, len(teamByID))
	for id := range teamByID {
		teamWeekGames[id] = make(map[int]bool)
	}

	// Shuffle the input so each attempt gets a different ordering.
	work := make([]matchup, len(matchups))
	copy(work, matchups)
	rand.Shuffle(len(work), func(i, j int) { work[i], work[j] = work[j], work[i] })

	// Separate week-18-only (divisional) matchups from the rest.
	var divisionalMatchups []matchup
	var otherMatchups []matchup
	// Also collect DAL and DET divisional games for week-13 priority.
	var dalDivisional []matchup
	var detDivisional []matchup
	for _, m := range work {
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

	assigned = make([]weekMatchup, 0, len(work))

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

	// Step 2: Fill week 18 with as many divisional matchups as possible.
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

	// Step 3: Assign remaining divisional and all non-divisional games to weeks 1-17.
	remaining := append(week18Overflow, otherMatchups...)
	rand.Shuffle(len(remaining), func(i, j int) { remaining[i], remaining[j] = remaining[j], remaining[i] })

	var unassigned []matchup
	for _, m := range remaining {
		placed := false
		weekOrder := rand.Perm(totalWeeks - 1)
		for _, wi := range weekOrder {
			w := wi + 1
			if w == 18 {
				continue
			}
			if canPlace(m, w) {
				assignToWeek(m, w)
				placed = true
				break
			}
		}
		if !placed {
			unassigned = append(unassigned, m)
		}
	}

	// Second pass: relax week-18 restriction.
	var stillUnassigned []matchup
	for _, m := range unassigned {
		placed := false
		for w := 1; w <= totalWeeks; w++ {
			if canPlace(m, w) {
				assignToWeek(m, w)
				placed = true
				break
			}
		}
		if !placed {
			stillUnassigned = append(stillUnassigned, m)
		}
	}

	// Third pass: swap-and-place. For each unplaced game, find a week where one team
	// is free and the other has a moveable game, then bump that game to another slot.
	for _, m := range stillUnassigned {
		placed := false
		for w := 1; w <= totalWeeks && !placed; w++ {
			if canPlace(m, w) {
				assignToWeek(m, w)
				placed = true
				break
			}
			// Determine which of the two teams is busy in this week.
			var busyID uint
			if teamWeekGames[m.homeID][w] && !teamWeekGames[m.awayID][w] {
				busyID = m.homeID
			} else if teamWeekGames[m.awayID][w] && !teamWeekGames[m.homeID][w] {
				busyID = m.awayID
			} else {
				continue // both busy — no point trying to bump
			}
			// Find the assigned game for busyID in week w and try to move it.
			for idx := range assigned {
				ag := &assigned[idx]
				if ag.week != w || (ag.homeID != busyID && ag.awayID != busyID) {
					continue
				}
				teamWeekGames[ag.homeID][w] = false
				teamWeekGames[ag.awayID][w] = false
				bumpedM := matchup{homeID: ag.homeID, awayID: ag.awayID, isDivisional: ag.isDivisional}
				moved := false
				for altW := 1; altW <= totalWeeks; altW++ {
					if altW == w {
						continue
					}
					if canPlace(bumpedM, altW) {
						ag.week = altW
						teamWeekGames[ag.homeID][altW] = true
						teamWeekGames[ag.awayID][altW] = true
						moved = true
						break
					}
				}
				if moved && canPlace(m, w) {
					assignToWeek(m, w)
					placed = true
				} else {
					// Restore original slot — swap didn't help.
					if moved {
						teamWeekGames[ag.homeID][ag.week] = false
						teamWeekGames[ag.awayID][ag.week] = false
						ag.week = w
					}
					teamWeekGames[ag.homeID][w] = true
					teamWeekGames[ag.awayID][w] = true
				}
				break
			}
		}
		if !placed {
			failed = append(failed, m)
		}
	}

	return assigned, failed
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
			if stadium.ID == 0 {
				fmt.Printf("[WARN] No NFL stadium record found for home team ID=%d (%s %s)\n",
					ht.ID, ht.TeamName, ht.Mascot)
			}

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

// applySimNFLByeWeekFlags sets HomePreviousBye and AwayPreviousBye on game records
// where the corresponding team had no game the previous week (i.e., was on bye).
func applySimNFLByeWeekFlags(nflGames []structs.NFLGame) {
	teamWeeks := make(map[int]map[int]bool)
	for _, g := range nflGames {
		if teamWeeks[g.HomeTeamID] == nil {
			teamWeeks[g.HomeTeamID] = make(map[int]bool)
		}
		if teamWeeks[g.AwayTeamID] == nil {
			teamWeeks[g.AwayTeamID] = make(map[int]bool)
		}
		teamWeeks[g.HomeTeamID][g.Week] = true
		teamWeeks[g.AwayTeamID][g.Week] = true
	}

	for i := range nflGames {
		g := &nflGames[i]
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
