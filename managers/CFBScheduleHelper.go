package managers

import (
	"math/rand/v2"
	"sort"

	"github.com/CalebRose/SimFBA/structs"
)

// SchedulerHistoryKey uniquely identifies a (teamA, teamB) pair in a canonical
// order (lower ID first) so we can look up play counts and home/away history.
type SchedulerHistoryKey struct {
	A, B uint
}

func makeHistoryKey(a, b uint) SchedulerHistoryKey {
	if a < b {
		return SchedulerHistoryKey{a, b}
	}
	return SchedulerHistoryKey{b, a}
}

// HistoryRecord stores how many times A (lower ID) hosted and B (higher ID)
// hosted this pair within the tracking window.
type HistoryRecord struct {
	AHosted uint // games where team with lower ID was home
	BHosted uint // games where team with higher ID was home
}

// BuildScheduleHistoryMaps builds two maps from a slice of historic regular-season
// games (non-spring, non-postseason):
//   - playCountMap: (teamA, teamB) -> total games played
//   - homeAwayMap:  (teamA, teamB) -> HistoryRecord
//   - lastHomeMap:  teamID -> map[opponentID]bool (true = teamID was home last time)
func BuildScheduleHistoryMaps(historicGames []structs.CollegeGame) (
	playCountMap map[SchedulerHistoryKey]int,
	homeAwayMap map[SchedulerHistoryKey]HistoryRecord,
	lastHomeMap map[uint]map[uint]bool,
) {
	playCountMap = make(map[SchedulerHistoryKey]int)
	homeAwayMap = make(map[SchedulerHistoryKey]HistoryRecord)
	lastHomeMap = make(map[uint]map[uint]bool)

	// Sort oldest to newest so "last" home tracking is correct.
	sorted := make([]structs.CollegeGame, len(historicGames))
	copy(sorted, historicGames)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].SeasonID < sorted[j].SeasonID ||
			(sorted[i].SeasonID == sorted[j].SeasonID && sorted[i].Week < sorted[j].Week)
	})

	for _, g := range sorted {
		if g.IsBowlGame || g.IsPlayoffGame || g.IsNationalChampionship || g.IsConferenceChampionship || g.IsSpringGame {
			continue
		}
		home := uint(g.HomeTeamID)
		away := uint(g.AwayTeamID)
		if home == 0 || away == 0 {
			continue
		}
		key := makeHistoryKey(home, away)
		playCountMap[key]++
		rec := homeAwayMap[key]
		if key.A == home {
			rec.AHosted++
		} else {
			rec.BHosted++
		}
		homeAwayMap[key] = rec

		// update lastHomeMap
		if lastHomeMap[home] == nil {
			lastHomeMap[home] = make(map[uint]bool)
		}
		if lastHomeMap[away] == nil {
			lastHomeMap[away] = make(map[uint]bool)
		}
		lastHomeMap[home][away] = true
		lastHomeMap[away][home] = false
	}
	return
}

// ShouldBeHome returns true if teamID should be the home team against oppID.
// Priority order:
//  1. homeCountMap: if one team has strictly fewer home games this session, they host.
//  2. lastHomeMap: if we have history, flip from last time.
//  3. Coin flip: (teamID + oppID + season) % 2 — deterministic per-pair per-season,
//     giving ~50% home/away balance for teams with no shared play history.
//
// Pass nil for homeCountMap or lastHomeMap when those signals are unavailable.
func ShouldBeHome(teamID, oppID uint, season int, lastHomeMap map[uint]map[uint]bool, homeCountMap map[uint]int) bool {
	// Priority 1: balance by home game count within the current scheduling session.
	if homeCountMap != nil {
		aHomes := homeCountMap[teamID]
		bHomes := homeCountMap[oppID]
		if aHomes < bHomes {
			return true
		}
		if bHomes < aHomes {
			return false
		}
		// Equal — fall through to history/parity
	}
	// Priority 2: flip from last known result.
	if lastHomeMap[teamID] != nil {
		if wasHome, ok := lastHomeMap[teamID][oppID]; ok {
			return !wasHome
		}
	}
	// Priority 3: coin flip — deterministic but varies by pair and season.
	// Shift right by one bit before testing parity: the plain sum of two odd IDs
	// is always even, which would produce a constant result for same-parity
	// conferences (e.g. Pac-12 where every team ID is odd). The >>1 uses the
	// second-to-last bit of the sum, which differs meaningfully even when all
	// IDs share the same parity, giving a roughly 50/50 split across pairs.
	return ((teamID+oppID+uint(season))>>1)%2 == 0
}

// assignWeek tries to place a floating matchup into an open week between minWeek
// and maxWeek, updating gamesPlayedByWeekMap. Returns 0 if no slot was found.
func assignWeek(homeID, awayID uint, minWeek, maxWeek uint, gamesPlayedByWeekMap map[uint]map[uint]bool) uint {
	for w := maxWeek; w >= minWeek; w-- {
		homeOccupied := gamesPlayedByWeekMap[homeID] != nil && gamesPlayedByWeekMap[homeID][w]
		awayOccupied := gamesPlayedByWeekMap[awayID] != nil && gamesPlayedByWeekMap[awayID][w]
		if !homeOccupied && !awayOccupied {
			if gamesPlayedByWeekMap[homeID] == nil {
				gamesPlayedByWeekMap[homeID] = make(map[uint]bool)
			}
			if gamesPlayedByWeekMap[awayID] == nil {
				gamesPlayedByWeekMap[awayID] = make(map[uint]bool)
			}
			gamesPlayedByWeekMap[homeID][w] = true
			gamesPlayedByWeekMap[awayID][w] = true
			return w
		}
	}
	// If pairing still isn't made, try to find any open week for both teams up to maxWeek to avoid leaving a game unassigned.
	// Check backwards from maxWeek to preserve the ideal conference window as much as possible, only using this as a last resort.
	// This is a fallback to prevent deadlock in edge cases where the ideal conference window is too tight to fit all games.
	if gamesPlayedByWeekMap[homeID] == nil && gamesPlayedByWeekMap[awayID] == nil {
		for w := uint(13); w >= 1; w-- {
			homeOccupied := gamesPlayedByWeekMap[homeID] != nil && gamesPlayedByWeekMap[homeID][w]
			awayOccupied := gamesPlayedByWeekMap[awayID] != nil && gamesPlayedByWeekMap[awayID][w]
			if !homeOccupied && !awayOccupied {
				if gamesPlayedByWeekMap[homeID] == nil {
					gamesPlayedByWeekMap[homeID] = make(map[uint]bool)
				}
				if gamesPlayedByWeekMap[awayID] == nil {
					gamesPlayedByWeekMap[awayID] = make(map[uint]bool)
				}
				gamesPlayedByWeekMap[homeID][w] = true
				gamesPlayedByWeekMap[awayID][w] = true
				return w
			}
		}
	}
	return 0
}

// markWeek marks both teams as playing in a given week (for pre-assigned locked games).
func markWeek(homeID, awayID, week uint, gamesPlayedByWeekMap map[uint]map[uint]bool) {
	if gamesPlayedByWeekMap[homeID] == nil {
		gamesPlayedByWeekMap[homeID] = make(map[uint]bool)
	}
	if gamesPlayedByWeekMap[awayID] == nil {
		gamesPlayedByWeekMap[awayID] = make(map[uint]bool)
	}
	gamesPlayedByWeekMap[homeID][week] = true
	gamesPlayedByWeekMap[awayID][week] = true
}

// markOpponents records that two teams have been matched against one another.
func markOpponents(homeID, awayID uint, gamesPlayedAgainstOpponentsMap map[uint]map[uint]bool) {
	if gamesPlayedAgainstOpponentsMap[homeID] == nil {
		gamesPlayedAgainstOpponentsMap[homeID] = make(map[uint]bool)
	}
	if gamesPlayedAgainstOpponentsMap[awayID] == nil {
		gamesPlayedAgainstOpponentsMap[awayID] = make(map[uint]bool)
	}
	gamesPlayedAgainstOpponentsMap[homeID][awayID] = true
	gamesPlayedAgainstOpponentsMap[awayID][homeID] = true
}

// alreadyScheduled returns true if the two teams have already been matched in
// the current season's opponent map.
func alreadyScheduled(a, b uint, gamesPlayedAgainstOpponentsMap map[uint]map[uint]bool) bool {
	if gamesPlayedAgainstOpponentsMap[a] == nil {
		return false
	}
	return gamesPlayedAgainstOpponentsMap[a][b]
}

// sortedTeamIDs returns the IDs of a team slice in ascending order.
func sortedTeamIDs(teams []structs.CollegeTeam) []uint {
	ids := make([]uint, len(teams))
	for i, t := range teams {
		ids[i] = t.ID
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return ids
}

// buildTeamMapFromSlice builds a map of ID -> CollegeTeam.
func buildTeamMapFromSlice(teams []structs.CollegeTeam) map[uint]structs.CollegeTeam {
	m := make(map[uint]structs.CollegeTeam, len(teams))
	for _, t := range teams {
		m[t.ID] = t
	}
	return m
}

// getRoundRobinPairs returns all unique (i, j) pairs from a list of IDs, i < j.
func getRoundRobinPairs(ids []uint) [][2]uint {
	var pairs [][2]uint
	for i := 0; i < len(ids); i++ {
		for j := i + 1; j < len(ids); j++ {
			pairs = append(pairs, [2]uint{ids[i], ids[j]})
		}
	}
	// Proposing a shuffle.
	rand.Shuffle(len(pairs), func(i, j int) {
		pairs[i], pairs[j] = pairs[j], pairs[i]
	})
	return pairs
}

// cycleYear maps a global season number to a 1-4 cycle year.
func cycleYear(season int) int {
	y := ((season - 1) % 4) + 1
	if y <= 0 {
		y += 4
	}
	return y
}

// selectGroupForCycleYear returns true if the team should play Group A this season.
func groupAThisSeason(season int) bool {
	cy := cycleYear(season)
	return cy == 1 || cy == 3
}

// homeForThisSeason returns whether teamID should be home against oppID given
// the cycle year and a given "base" home assignment (true = teamID hosts in Year 1).
// Years 1 & 3 keep base, Years 2 & 4 flip.
func homeForCycleYear(season int, baseIsHome bool) bool {
	cy := cycleYear(season)
	if cy == 1 || cy == 3 {
		return baseIsHome
	}
	return !baseIsHome
}

// deepCopyBoolMap creates a full deep copy of a map[uint]map[uint]bool.
// Used to snapshot shared scheduling maps before retry attempts so each seed
// can be tried in isolation and the winning attempt committed cleanly.
func deepCopyBoolMap(src map[uint]map[uint]bool) map[uint]map[uint]bool {
	dst := make(map[uint]map[uint]bool, len(src))
	for k, inner := range src {
		innerCopy := make(map[uint]bool, len(inner))
		for ik, iv := range inner {
			innerCopy[ik] = iv
		}
		dst[k] = innerCopy
	}
	return dst
}

// generateGenericRoundRobinSchedule creates a conference schedule for any conference
// using a round-robin model with a specified game count, conference window, and
// optional locked-week rivalries (intraConf only). It obeys the global
// gamesPlayedAgainstOpponentsMap and gamesPlayedByWeekMap.
//
// lockedPairs: map of (lowID, highID) -> preferred week (0 = no lock).
// conferenceGamesPerTeam: how many conf games each team plays (< n-1 means partial round-robin).
// confWindow: [minWeek, maxWeek] for floating games.
// overflowWindow: [minWeek, maxWeek] used only if confWindow exhausted.
func generateGenericRoundRobinSchedule(
	collegeTeams []structs.CollegeTeam,
	stadiumMap map[uint]structs.Stadium,
	stadiumMapByID map[uint]structs.Stadium,
	rivalryMap map[uint][]structs.CollegeRival,
	gamesPlayedAgainstOpponentsMap map[uint]map[uint]bool,
	gamesPlayedByWeekMap map[uint]map[uint]bool,
	playCountMap map[SchedulerHistoryKey]int,
	lastHomeMap map[uint]map[uint]bool,
	ts structs.Timestamp,
	lockedPairs map[SchedulerHistoryKey]uint, // key -> preferred week (0 = floating)
	conferenceGamesPerTeam int,
	minWeek, maxWeek uint,
	overflowMinWeek, overflowMaxWeek uint,
	retrySeed uint, // 0 = default pair order; >0 = shuffled for retry
) []structs.CollegeGame {
	games := []structs.CollegeGame{}
	season := ts.Season
	seasonID := uint(ts.CollegeSeasonID)

	teamMap := buildTeamMapFromSlice(collegeTeams)
	ids := sortedTeamIDs(collegeTeams)

	// Track how many conference games each team has been assigned this season.
	confGameCount := make(map[uint]int)
	// Track home-game count per team within this scheduling session so ShouldBeHome
	// can balance home/away assignments, especially for new programs with no history.
	homeCountMap := make(map[uint]int)

	// Pre-seed confGameCount from conference games already placed before this
	// function runs (e.g. annual rivalry games between conference teammates that
	// were scheduled in the rivalry pass and are already in gamesPlayedAgainstOpponentsMap).
	for _, t := range collegeTeams {
		for oppID := range gamesPlayedAgainstOpponentsMap[t.ID] {
			if opp, ok := teamMap[oppID]; ok && opp.ConferenceID != 0 && opp.ConferenceID == t.ConferenceID {
				confGameCount[t.ID]++
			}
		}
	}

	// Step 1: place locked pairs first
	type lockedGame struct {
		homeID, awayID uint
		week           uint
	}
	var locked []lockedGame
	for key, week := range lockedPairs {
		if week == 0 {
			continue
		}
		if alreadyScheduled(key.A, key.B, gamesPlayedAgainstOpponentsMap) {
			continue
		}
		a := teamMap[key.A]
		b := teamMap[key.B]
		if a.ID == 0 || b.ID == 0 {
			continue
		}
		var home, away structs.CollegeTeam
		if ShouldBeHome(key.A, key.B, season, lastHomeMap, homeCountMap) {
			home, away = a, b
		} else {
			home, away = b, a
		}
		// Skip if either team is already playing this week (pre-existing OOC or conf game).
		if gamesPlayedByWeekMap[home.ID] != nil && gamesPlayedByWeekMap[home.ID][week] {
			continue
		}
		if gamesPlayedByWeekMap[away.ID] != nil && gamesPlayedByWeekMap[away.ID][week] {
			continue
		}
		locked = append(locked, lockedGame{homeID: home.ID, awayID: away.ID, week: week})
		markWeek(home.ID, away.ID, week, gamesPlayedByWeekMap)
		markOpponents(home.ID, away.ID, gamesPlayedAgainstOpponentsMap)
		isConference := home.ConferenceID != 0 && home.ConferenceID == away.ConferenceID
		if isConference {
			confGameCount[home.ID]++
			confGameCount[away.ID]++
		}
		homeCountMap[home.ID]++
		g := CreateCollegeGameRecord(home, away, week, seasonID, stadiumMap, stadiumMapByID, rivalryMap)
		games = append(games, g)
	}

	// Step 2: floating locked pairs (week == 0 in lockedPairs but pair still required)
	for key, week := range lockedPairs {
		if week != 0 {
			continue // already handled above
		}
		if alreadyScheduled(key.A, key.B, gamesPlayedAgainstOpponentsMap) {
			continue
		}
		a := teamMap[key.A]
		b := teamMap[key.B]
		if a.ID == 0 || b.ID == 0 {
			continue
		}
		// Respect the conference game cap for floating locked pairs, same as Step 3.
		if a.ConferenceID != 0 && a.ConferenceID == b.ConferenceID &&
			(confGameCount[key.A] >= conferenceGamesPerTeam || confGameCount[key.B] >= conferenceGamesPerTeam) {
			continue
		}
		var home, away structs.CollegeTeam
		if ShouldBeHome(key.A, key.B, season, lastHomeMap, homeCountMap) {
			home, away = a, b
		} else {
			home, away = b, a
		}
		w := assignWeek(home.ID, away.ID, minWeek, maxWeek, gamesPlayedByWeekMap)
		if w == 0 {
			w = assignWeek(home.ID, away.ID, overflowMinWeek, overflowMaxWeek, gamesPlayedByWeekMap)
		}
		if w == 0 {
			continue
		}
		markOpponents(home.ID, away.ID, gamesPlayedAgainstOpponentsMap)
		isConference := home.ConferenceID != 0 && home.ConferenceID == away.ConferenceID
		if isConference {
			confGameCount[home.ID]++
			confGameCount[away.ID]++
		}
		homeCountMap[home.ID]++
		g := CreateCollegeGameRecord(home, away, w, seasonID, stadiumMap, stadiumMapByID, rivalryMap)
		games = append(games, g)
	}

	// Step 3: Fill remaining conference games
	// Build a priority list of pairs, favouring least-recently-played opponents.
	pairs := getRoundRobinPairs(ids)
	// Sort pairs by combined play count ascending
	sort.SliceStable(pairs, func(i, j int) bool {
		ki := makeHistoryKey(pairs[i][0], pairs[i][1])
		kj := makeHistoryKey(pairs[j][0], pairs[j][1])
		return playCountMap[ki] < playCountMap[kj]
	})
	// Always shuffle pairs to reduce ID-ordering bias. When retrySeed is 0 (first pass),
	// use the season number so the order is deterministic per season but not tied to
	// team IDs. Retry attempts use different seeds to vary ordering further.
	{
		seed := retrySeed
		if seed == 0 {
			seed = uint(season)
		}
		r := uint64(seed) * 2654435761
		for i := len(pairs) - 1; i > 0; i-- {
			r = r*6364136223846793005 + 1442695040888963407
			j := int(r>>33) % (i + 1)
			pairs[i], pairs[j] = pairs[j], pairs[i]
		}
	}

	for _, pair := range pairs {
		a, b := pair[0], pair[1]
		if alreadyScheduled(a, b, gamesPlayedAgainstOpponentsMap) {
			continue
		}
		if confGameCount[a] >= conferenceGamesPerTeam || confGameCount[b] >= conferenceGamesPerTeam {
			continue
		}
		tA := teamMap[a]
		tB := teamMap[b]
		if tA.ID == 0 || tB.ID == 0 {
			continue
		}
		var home, away structs.CollegeTeam
		if ShouldBeHome(a, b, season, lastHomeMap, homeCountMap) {
			home, away = tA, tB
		} else {
			home, away = tB, tA
		}
		w := assignWeek(home.ID, away.ID, minWeek, maxWeek, gamesPlayedByWeekMap)
		if w == 0 {
			w = assignWeek(home.ID, away.ID, overflowMinWeek, overflowMaxWeek, gamesPlayedByWeekMap)
		}
		if w == 0 {
			continue
		}
		markOpponents(home.ID, away.ID, gamesPlayedAgainstOpponentsMap)
		isConference := home.ConferenceID != 0 && home.ConferenceID == away.ConferenceID
		if isConference {
			confGameCount[home.ID]++
			confGameCount[away.ID]++
		}
		homeCountMap[home.ID]++
		g := CreateCollegeGameRecord(home, away, w, seasonID, stadiumMap, stadiumMapByID, rivalryMap)
		games = append(games, g)
	}

	// Retry signal: if any team is below target, return empty slice.
	// Exception: with an odd number of teams and an odd per-team target, the sum
	// (n × target) is odd and cannot be divided into whole game-pairs, so exactly
	// one team will always be one game short. Accept that result rather than
	// burning all retry slots on an unsolvable parity problem.
	// oddTeamsParity := len(ids)%2 == 1 && conferenceGamesPerTeam%2 == 1
	for _, teamID := range ids {
		if confGameCount[teamID] < conferenceGamesPerTeam {
			// if oddTeamsParity {
			// 	return games // one short is the mathematical floor; accept it
			// }
			return []structs.CollegeGame{}
		}
	}

	return games
}

// MissingGame records a conference matchup that ValidateAndRescueConference
// could not fit into any available week.
type MissingGame struct {
	ConferenceID int
	TeamAID      uint
	TeamBID      uint
	Reason       string
}

// ConferenceExpectedGames maps conference ID to the minimum number of conference
// games every team in that conference should play. Used by the validation pass.
var ConferenceExpectedGames = map[int]int{
	3:  8, // ACC (one designated team plays 8, rest play 9; floor is 8)
	4:  9, // Big Ten
	5:  9, // Big 12
	6:  7, // Pac-12
	7:  9, // SEC
	8:  8, // AAC
	9:  8, // C-USA
	10: 8, // MAC (3 intra-pod + 6 cross-pod = 9 per MAC code design)
	11: 8, // MWC
	12: 8, // Sun Belt
	14: 8, // MVFC
	15: 7, // Patriot
	16: 8, // Big Sky
	17: 7, // Ivy League
	18: 9, // SoCon
	19: 8, // SWAC
	20: 7, // Big South-OVC
	21: 8, // CAA
	23: 5, // MEAC
	24: 8, // NEC
	25: 8, // Pioneer (11 teams × 9 target = 99 slots / 2 = non-integer; one team always gets 8)
	26: 8, // Southland
	27: 7, // UAC
}

// rebalanceConferenceHomeAway performs a multi-pass greedy swap of home/away
// roles for conference games where the current home team is home-heavy and the
// current away team is away-heavy. The goal is to bring every team's
// |home − away| conference game count to ≤ 2.
//
// Neutral-site and rivalry games are never swapped. The game record (stadium,
// city, timeslot, etc.) is fully rebuilt from the new home team's data for
// every swapped pair.
func rebalanceConferenceHomeAway(
	games []structs.CollegeGame,
	teamMap map[uint]structs.CollegeTeam,
	stadiumMap map[uint]structs.Stadium,
	stadiumMapByID map[uint]structs.Stadium,
	rivalryMap map[uint][]structs.CollegeRival,
) []structs.CollegeGame {
	const maxPasses = 20
	for pass := 0; pass < maxPasses; pass++ {
		// Recount home/away at the start of each pass.
		homeCount := make(map[uint]int)
		awayCount := make(map[uint]int)
		for _, g := range games {
			if !g.IsConference || g.IsNeutral {
				continue
			}
			homeCount[uint(g.HomeTeamID)]++
			awayCount[uint(g.AwayTeamID)]++
		}

		swapMade := false
		for i, g := range games {
			if !g.IsConference || g.IsNeutral {
				continue
			}
			rivalries := rivalryMap[uint(g.HomeTeamID)]
			isAnnualRivalry := false
			for _, r := range rivalries {
				if (r.TeamOneID == uint(g.AwayTeamID) || r.TeamTwoID == uint(g.AwayTeamID)) && r.IsAnnualRivalry {
					isAnnualRivalry = true
					break
				}
			}
			if isAnnualRivalry {
				continue
			}
			homeID := uint(g.HomeTeamID)
			awayID := uint(g.AwayTeamID)

			// homeDiff > 0: current home team has more home than away conf games.
			homeDiff := homeCount[homeID] - awayCount[homeID]

			// Only bother swapping when the home team is noticeably home-heavy.
			if homeDiff <= 0 {
				continue
			}

			// Compute what the away team's imbalance would look like after
			// the swap: they gain 1 home game and lose 1 away game.
			//   postSwapAwayDiff = (homeCount[away]+1) − (awayCount[away]−1)
			//                    = homeCount[away] − awayCount[away] + 2
			// We allow the swap as long as the away team doesn't end up with
			// a home-heavy diff > 2 (i.e. worse than acceptable).
			postSwapAwayDiff := (homeCount[awayID] + 1) - (awayCount[awayID] - 1)
			if postSwapAwayDiff > 3 {
				continue
			}

			newHome := teamMap[awayID]
			newAway := teamMap[homeID]
			if newHome.ID == 0 || newAway.ID == 0 {
				continue
			}
			games[i] = CreateCollegeGameRecord(newHome, newAway, uint(g.Week), uint(g.SeasonID), stadiumMap, stadiumMapByID, rivalryMap)
			// Update running counts immediately so later games in this pass
			// see the corrected state rather than the pre-swap state.
			homeCount[homeID]--
			awayCount[homeID]++
			homeCount[awayID]++
			awayCount[awayID]--
			swapMade = true
		}
		if !swapMade {
			break
		}
	}
	return games
}

// ValidateAndRescueConference verifies every team in the conference has at least
// expectedGames conference games scheduled. For any team that is short, it tries
// to add games vs unplayed conference opponents using the full week 1–14 window.
// Returns newly created rescue games and any pairs that still cannot be scheduled.
func ValidateAndRescueConference(
	conferenceID int,
	expectedGames int,
	teams []structs.CollegeTeam,
	stadiumMap map[uint]structs.Stadium,
	stadiumMapByID map[uint]structs.Stadium,
	rivalryMap map[uint][]structs.CollegeRival,
	gamesPlayedAgainstOpponentsMap map[uint]map[uint]bool,
	gamesPlayedByWeekMap map[uint]map[uint]bool,
	seasonID uint,
	season int,
) ([]structs.CollegeGame, []MissingGame) {
	rescueGames := []structs.CollegeGame{}
	var missing []MissingGame

	if len(teams) == 0 {
		return rescueGames, missing
	}

	// Build a set of conference team IDs for fast membership checks.
	confTeamSet := make(map[uint]bool, len(teams))
	for _, t := range teams {
		confTeamSet[t.ID] = true
	}
	teamMap := buildTeamMapFromSlice(teams)

	// Count conference games already scheduled for each team.
	confGameCount := make(map[uint]int, len(teams))
	for _, t := range teams {
		for oppID := range gamesPlayedAgainstOpponentsMap[t.ID] {
			if confTeamSet[oppID] {
				confGameCount[t.ID]++
			}
		}
	}

	// Track which pairs we've already attempted to avoid double-scheduling.
	rescuedPairs := make(map[SchedulerHistoryKey]bool)

	for _, t := range teams {
		if confGameCount[t.ID] >= expectedGames {
			continue
		}
		for _, opp := range teams {
			if opp.ID == t.ID {
				continue
			}
			if confGameCount[t.ID] >= expectedGames {
				break
			}
			if alreadyScheduled(t.ID, opp.ID, gamesPlayedAgainstOpponentsMap) {
				continue
			}
			key := makeHistoryKey(t.ID, opp.ID)
			if rescuedPairs[key] {
				continue
			}
			rescuedPairs[key] = true
			if confGameCount[opp.ID] >= expectedGames {
				continue
			}

			var home, away structs.CollegeTeam
			if ShouldBeHome(t.ID, opp.ID, season, nil, nil) {
				home, away = teamMap[t.ID], teamMap[opp.ID]
			} else {
				home, away = teamMap[opp.ID], teamMap[t.ID]
			}

			w := assignWeek(home.ID, away.ID, 1, 14, gamesPlayedByWeekMap)
			if w == 0 {
				missing = append(missing, MissingGame{
					ConferenceID: conferenceID,
					TeamAID:      t.ID,
					TeamBID:      opp.ID,
					Reason:       "no open week available in 1-14 for both teams",
				})
				continue
			}

			markOpponents(home.ID, away.ID, gamesPlayedAgainstOpponentsMap)
			confGameCount[home.ID]++
			confGameCount[away.ID]++
			g := CreateCollegeGameRecord(home, away, w, seasonID, stadiumMap, stadiumMapByID, rivalryMap)
			rescueGames = append(rescueGames, g)
		}
	}

	return rescueGames, missing
}
