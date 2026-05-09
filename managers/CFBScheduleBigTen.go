package managers

import (
	"sort"

	"github.com/CalebRose/SimFBA/structs"
)

// ============================================================
// Big Ten Conference Schedule Generator
// 18 teams | no divisions | 9 conf games per team
// Protected annual rivalries (some Week 14, some mid-season)
// 5-year rolling cycle — coverage gap prioritisation
// ============================================================
//
// Big Ten team IDs:
// Illinois(40), Indiana(41), Iowa(42), Maryland(55), Michigan(59),
// Michigan State(60), Minnesota(62), Nebraska(67), Northwestern(74),
// Ohio State(77), Oregon(82), Penn State(84), Purdue(86), Rutgers(88),
// UCLA(111), USC(115), Washington(124), Wisconsin(129)

// bigTenWeek14Rivalries: pairs always played in Week 14
var bigTenWeek14Rivalries = []struct {
	key  SchedulerHistoryKey
	week uint
}{
	{makeHistoryKey(59, 77), 14},   // Michigan vs Ohio State
	{makeHistoryKey(42, 67), 14},   // Iowa vs Nebraska
	{makeHistoryKey(62, 129), 14},  // Minnesota vs Wisconsin
	{makeHistoryKey(41, 86), 14},   // Indiana vs Purdue
	{makeHistoryKey(40, 74), 14},   // Illinois vs Northwestern
	{makeHistoryKey(55, 88), 14},   // Maryland vs Rutgers
	{makeHistoryKey(111, 115), 14}, // UCLA vs USC
	{makeHistoryKey(84, 60), 14},   // Penn State vs Michigan State

}

// bigTenMidSeasonRivalries: pairs with a traditional fixed week (not Week 14)
var bigTenMidSeasonRivalries = []struct {
	key  SchedulerHistoryKey
	week uint
}{
	{makeHistoryKey(82, 124), 8},  // Oregon vs Washington
	{makeHistoryKey(59, 60), 8},   // Michigan vs Michigan State
	{makeHistoryKey(42, 129), 10}, // Iowa vs Wisconsin
	{makeHistoryKey(42, 62), 8},   // Iowa vs Minnesota
	{makeHistoryKey(67, 62), 6},   // Nebraska vs Minnesota

}

// bigTenTeamIDs for convenience
var bigTenTeamIDs = []uint{40, 41, 42, 55, 59, 60, 62, 67, 74, 77, 82, 84, 86, 88, 111, 115, 124, 129}

// bigTenProtectedSet returns a set of all protected-pair keys.
func bigTenProtectedSet() map[SchedulerHistoryKey]uint {
	m := make(map[SchedulerHistoryKey]uint)
	for _, r := range bigTenWeek14Rivalries {
		m[r.key] = r.week
	}
	for _, r := range bigTenMidSeasonRivalries {
		m[r.key] = r.week
	}
	return m
}

// GenerateBigTenSchedule produces all Big Ten conference games for the season.
// Returns an empty slice if Phase 4 cannot fill all required games, which signals
// the manager's retry loop to reseed and try again.
func GenerateBigTenSchedule(
	collegeTeams []structs.CollegeTeam,
	stadiumMap map[uint]structs.Stadium,
	stadiumMapByID map[uint]structs.Stadium,
	rivalryMap map[uint][]structs.CollegeRival,
	gamesPlayedAgainstOpponentsMap map[uint]map[uint]bool,
	gamesPlayedByWeekMap map[uint]map[uint]bool,
	playCountMap map[SchedulerHistoryKey]int,
	lastHomeMap map[uint]map[uint]bool,
	homeCountSeedMap map[uint]int,
	ts structs.Timestamp,
) []structs.CollegeGame {
	games := []structs.CollegeGame{}
	season := ts.Season
	seasonID := uint(ts.CollegeSeasonID)

	teamMap := buildTeamMapFromSlice(collegeTeams)
	protectedSet := bigTenProtectedSet()

	bigTenTeamSet := make(map[uint]bool, len(bigTenTeamIDs))
	for _, id := range bigTenTeamIDs {
		bigTenTeamSet[id] = true
	}

	// Pre-seed confGameCount from any conference games already placed before this
	// generator runs (e.g. annual rivalry games between Big Ten teammates from the
	// rivalry pass that are already in gamesPlayedAgainstOpponentsMap).
	confGameCount := make(map[uint]int)
	for _, id := range bigTenTeamIDs {
		for oppID := range gamesPlayedAgainstOpponentsMap[id] {
			if bigTenTeamSet[oppID] {
				confGameCount[id]++
			}
		}
	}

	// Seed homecountMap from rivalry-pass home game counts so ShouldBeHome sees
	// correct context from the very first conference game assignment.
	homecountMap := make(map[uint]int, len(homeCountSeedMap))
	for id, count := range homeCountSeedMap {
		homecountMap[id] = count
	}

	// recordGame commits a game and updates all tracking maps.
	recordGame := func(home, away structs.CollegeTeam, week uint) {
		markWeek(home.ID, away.ID, week, gamesPlayedByWeekMap)
		markOpponents(home.ID, away.ID, gamesPlayedAgainstOpponentsMap)
		homecountMap[home.ID]++
		confGameCount[home.ID]++
		confGameCount[away.ID]++
		g := MakeCollegeGameRecord(home, away, week, seasonID, stadiumMap, stadiumMapByID, rivalryMap)
		games = append(games, g)
	}

	// Phase 1: Lock all protected rivalries at their designated weeks.
	// Explicitly verify the week is free before recording to handle cases where
	// a pre-existing rivalry-pass game has already claimed that slot.
	for _, r := range bigTenWeek14Rivalries {
		a := teamMap[r.key.A]
		b := teamMap[r.key.B]
		if a.ID == 0 || b.ID == 0 {
			continue
		}
		if alreadyScheduled(a.ID, b.ID, gamesPlayedAgainstOpponentsMap) {
			continue
		}
		if confGameCount[a.ID] >= 9 || confGameCount[b.ID] >= 9 {
			continue
		}
		if gamesPlayedByWeekMap[a.ID] != nil && gamesPlayedByWeekMap[a.ID][r.week] {
			continue
		}
		if gamesPlayedByWeekMap[b.ID] != nil && gamesPlayedByWeekMap[b.ID][r.week] {
			continue
		}
		var home, away structs.CollegeTeam
		if ShouldBeHome(r.key.A, r.key.B, season, lastHomeMap, homecountMap) {
			home, away = a, b
		} else {
			home, away = b, a
		}
		recordGame(home, away, r.week)
	}
	for _, r := range bigTenMidSeasonRivalries {
		a := teamMap[r.key.A]
		b := teamMap[r.key.B]
		if a.ID == 0 || b.ID == 0 {
			continue
		}
		if alreadyScheduled(a.ID, b.ID, gamesPlayedAgainstOpponentsMap) {
			continue
		}
		if confGameCount[a.ID] >= 9 || confGameCount[b.ID] >= 9 {
			continue
		}
		if gamesPlayedByWeekMap[a.ID] != nil && gamesPlayedByWeekMap[a.ID][r.week] {
			continue
		}
		if gamesPlayedByWeekMap[b.ID] != nil && gamesPlayedByWeekMap[b.ID][r.week] {
			continue
		}
		var home, away structs.CollegeTeam
		if ShouldBeHome(r.key.A, r.key.B, season, lastHomeMap, homecountMap) {
			home, away = a, b
		} else {
			home, away = b, a
		}
		recordGame(home, away, r.week)
	}

	// Phase 3: Fill rotating conference games via shuffled round-robin.
	// getRoundRobinPairs returns a randomly shuffled list, preventing ID-ordering bias.
	ids := sortedTeamIDs(collegeTeams)
	pairs := getRoundRobinPairs(ids)
	// Sort by historical play count ascending so least-played pairs get priority.
	// Stable sort preserves the shuffle order from getRoundRobinPairs within ties.
	sort.SliceStable(pairs, func(i, j int) bool {
		ki := makeHistoryKey(pairs[i][0], pairs[i][1])
		kj := makeHistoryKey(pairs[j][0], pairs[j][1])
		return playCountMap[ki] < playCountMap[kj]
	})

	for _, pair := range pairs {
		a, b := pair[0], pair[1]
		pkey := makeHistoryKey(a, b)
		if _, isProtected := protectedSet[pkey]; isProtected {
			continue // already handled in Phase 1
		}
		if alreadyScheduled(a, b, gamesPlayedAgainstOpponentsMap) {
			continue
		}
		if confGameCount[a] >= 9 || confGameCount[b] >= 9 {
			continue
		}
		tA := teamMap[a]
		tB := teamMap[b]
		if tA.ID == 0 || tB.ID == 0 {
			continue
		}
		var home, away structs.CollegeTeam
		if ShouldBeHome(a, b, season, lastHomeMap, homecountMap) {
			home, away = tA, tB
		} else {
			home, away = tB, tA
		}
		w := assignWeek(home.ID, away.ID, 4, 13, gamesPlayedByWeekMap)
		if w == 0 {
			continue
		}
		recordGame(home, away, w)
	}

	// Phase 4: Validation pass — check every team's conference game count and
	// attempt to fill any remaining gaps with an expanded week window (1–14).
	// Oregon and Washington are free in Week 14 so the full 1–14 range is valid.
	// Protected pairs are skipped — they have designated weeks and should have
	// been placed in Phase 1; reschedule attempts here could violate week locks.
	for _, teamID := range bigTenTeamIDs {
		if confGameCount[teamID] >= 9 {
			continue
		}
		for _, oppID := range bigTenTeamIDs {
			if oppID == teamID {
				continue
			}
			if confGameCount[teamID] >= 9 {
				break
			}
			if confGameCount[oppID] >= 9 {
				continue
			}
			if alreadyScheduled(teamID, oppID, gamesPlayedAgainstOpponentsMap) {
				continue
			}
			pkey := makeHistoryKey(teamID, oppID)
			if _, isProtected := protectedSet[pkey]; isProtected {
				continue
			}
			tA := teamMap[teamID]
			tB := teamMap[oppID]
			if tA.ID == 0 || tB.ID == 0 {
				continue
			}
			var home, away structs.CollegeTeam
			if ShouldBeHome(teamID, oppID, season, lastHomeMap, homecountMap) {
				home, away = tA, tB
			} else {
				home, away = tB, tA
			}
			w := assignWeek(home.ID, away.ID, 1, 14, gamesPlayedByWeekMap)
			if w == 0 {
				continue
			}
			recordGame(home, away, w)
		}
	}

	// Final validation: if any team is still short, return an empty slice to
	// signal the manager's retry loop to reseed and try again.
	for _, teamID := range bigTenTeamIDs {
		if confGameCount[teamID] < 9 {
			return []structs.CollegeGame{}
		}
	}

	return games
}
