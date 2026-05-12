package managers

import (
	"math/rand/v2"
	"sort"

	"github.com/CalebRose/SimFBA/structs"
)

// ============================================================
// SEC Conference Schedule Generator
// 16 teams | 9 conf games | 12 total games | 3 bye weeks
// 3 guaranteed annual opponents per team
// 6 flexible opponents on 4-year Group A/B rotation
// Week 14: 6 conf rivalry locks + 4 OOC locks + 6 bye teams
// ============================================================

// SEC team IDs
// Alabama(3), Arkansas(7), Auburn(10), Florida(30), Georgia(34),
// Kentucky(47), LSU(53), Mississippi State(63), Missouri(64),
// Oklahoma(78), Ole Miss(81), South Carolina(93), Tennessee(100),
// Texas(101), Texas A&M(102), Vanderbilt(120)

// secGuaranteedPairs lists all 24 guaranteed annual conference matchups.
// Every pair will be scheduled every season (home/away flips on 2-year cycle).
var secGuaranteedPairs = []SchedulerHistoryKey{
	makeHistoryKey(3, 10),    // Alabama vs Auburn      (Week 14 lock)
	makeHistoryKey(3, 100),   // Alabama vs Tennessee
	makeHistoryKey(3, 63),    // Alabama vs Mississippi State
	makeHistoryKey(7, 53),    // Arkansas vs LSU        (Week 14 lock)
	makeHistoryKey(7, 64),    // Arkansas vs Missouri
	makeHistoryKey(7, 101),   // Arkansas vs Texas
	makeHistoryKey(10, 34),   // Auburn vs Georgia
	makeHistoryKey(10, 120),  // Auburn vs Vanderbilt
	makeHistoryKey(30, 34),   // Florida vs Georgia
	makeHistoryKey(30, 93),   // Florida vs South Carolina
	makeHistoryKey(30, 47),   // Florida vs Kentucky
	makeHistoryKey(34, 93),   // Georgia vs South Carolina
	makeHistoryKey(47, 100),  // Kentucky vs Tennessee
	makeHistoryKey(47, 93),   // Kentucky vs South Carolina
	makeHistoryKey(53, 81),   // LSU vs Ole Miss
	makeHistoryKey(53, 102),  // LSU vs Texas A&M
	makeHistoryKey(63, 120),  // Mississippi State vs Vanderbilt
	makeHistoryKey(63, 81),   // Mississippi State vs Ole Miss (Week 14 lock)
	makeHistoryKey(64, 78),   // Missouri vs Oklahoma          (Week 14 lock)
	makeHistoryKey(64, 102),  // Missouri vs Texas A&M
	makeHistoryKey(78, 101),  // Oklahoma vs Texas
	makeHistoryKey(78, 81),   // Oklahoma vs Ole Miss
	makeHistoryKey(100, 120), // Tennessee vs Vanderbilt      (Week 14 lock)
	makeHistoryKey(101, 102), // Texas vs Texas A&M           (Week 14 lock)
}

// secWeek14ConfLocks: pairs locked to Week 14 as conference games.
var secWeek14ConfLocks = map[SchedulerHistoryKey]bool{
	makeHistoryKey(3, 10):    true, // Alabama vs Auburn
	makeHistoryKey(7, 53):    true, // Arkansas vs LSU
	makeHistoryKey(63, 81):   true, // Ole Miss vs Mississippi State
	makeHistoryKey(64, 78):   true, // Missouri vs Oklahoma
	makeHistoryKey(100, 120): true, // Tennessee vs Vanderbilt
	makeHistoryKey(101, 102): true, // Texas vs Texas A&M
}

// secWeek14OOC: teams with OOC rivals in Week 14 (handled by annual rivalry pass).
// These teams take a bye from conf scheduling in Week 14.
var secWeek14OOC = map[uint]bool{
	30: true, // Florida  (vs Florida State)
	34: true, // Georgia  (vs Georgia Tech)
	47: true, // Kentucky (vs Louisville)
	93: true, // South Carolina (vs Clemson)
}

// secTeamIDs for convenience
var secTeamIDs = []uint{3, 7, 10, 30, 34, 47, 53, 63, 64, 78, 81, 93, 100, 101, 102, 120}

// secFlexibleGroupA maps each team to their 6 Group A (Years 1 & 3) flexible opponents.
// Group B is derived as the remaining 12 non-guaranteed opponents minus Group A.
var secFlexibleGroupA = map[uint][]uint{
	3:   {7, 30, 64, 78, 81, 101},  // Alabama
	7:   {3, 10, 47, 81, 93, 120},  // Arkansas
	10:  {7, 47, 53, 81, 101, 102}, // Auburn
	30:  {3, 7, 10, 53, 64, 101},   // Florida
	34:  {7, 53, 64, 78, 101, 102}, // Georgia (was {7,10,53,64,101,102}; removed Auburn=10 guaranteed, added Oklahoma=78)
	47:  {3, 7, 10, 34, 63, 101},   // Kentucky
	53:  {3, 10, 30, 64, 78, 120},  // LSU
	63:  {30, 34, 47, 78, 93, 101}, // Miss State (was {3,30,34,47,93,101}; removed Alabama=3 guaranteed, added Oklahoma=78)
	64:  {3, 7, 10, 30, 47, 81},    // Missouri
	78:  {3, 7, 30, 53, 63, 102},   // Oklahoma
	81:  {3, 7, 10, 30, 47, 102},   // Ole Miss
	93:  {3, 7, 10, 53, 64, 78},    // South Carolina
	100: {7, 30, 53, 63, 64, 78},   // Tennessee (was {3,7,30,53,63,64}; removed Alabama=3 guaranteed, added Oklahoma=78)
	101: {3, 47, 53, 81, 93, 120},  // Texas (was {3,7,47,53,78,120}; removed Arkansas=7+Oklahoma=78 guaranteed, added OleMiss=81+SC=93)
	102: {3, 7, 30, 34, 47, 81},    // Texas A&M
	120: {3, 30, 34, 47, 53, 78},   // Vanderbilt (was {3,10,30,34,47,53}; removed Auburn=10 guaranteed, added Oklahoma=78)
}

// secFlexGroupB derives Group B for a team as all non-guaranteed non-GroupA opponents.
func secFlexGroupB(teamID uint) []uint {
	guaranteed := make(map[uint]bool)
	for _, key := range secGuaranteedPairs {
		if key.A == teamID {
			guaranteed[key.B] = true
		} else if key.B == teamID {
			guaranteed[key.A] = true
		}
	}
	groupA := make(map[uint]bool)
	for _, id := range secFlexibleGroupA[teamID] {
		groupA[id] = true
	}
	var groupB []uint
	for _, id := range secTeamIDs {
		if id == teamID {
			continue
		}
		if !guaranteed[id] && !groupA[id] {
			groupB = append(groupB, id)
		}
	}
	return groupB
}

// GenerateSECSchedule produces all SEC conference games for the season.
func GenerateSECSchedule(
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
	useGroupA := groupAThisSeason(season)

	teamMap := buildTeamMapFromSlice(collegeTeams)

	// Seed homecountMap from rivalry-pass home game counts so ShouldBeHome sees
	// correct context from the very first conference game assignment.
	homecountMap := make(map[uint]int, len(homeCountSeedMap))
	for id, count := range homeCountSeedMap {
		homecountMap[id] = count
	}

	// Pre-seed confGameCount from any conference games already placed before this
	// generator runs (e.g. annual rivalry games from the rivalry pass).
	const secConferenceGamesPerTeam = 9
	secTeamSet := make(map[uint]bool, len(secTeamIDs))
	for _, id := range secTeamIDs {
		secTeamSet[id] = true
	}
	confGameCount := make(map[uint]int, len(secTeamIDs))
	for _, id := range secTeamIDs {
		for oppID := range gamesPlayedAgainstOpponentsMap[id] {
			if secTeamSet[oppID] {
				confGameCount[id]++
			}
		}
	}

	// recordGame commits a game and updates all tracking maps including live confGameCount.
	recordGame := func(home, away structs.CollegeTeam, week uint) {
		markWeek(home.ID, away.ID, week, gamesPlayedByWeekMap)
		markOpponents(home.ID, away.ID, gamesPlayedAgainstOpponentsMap)
		homecountMap[home.ID]++
		confGameCount[home.ID]++
		confGameCount[away.ID]++
		g := MakeCollegeGameRecord(home, away, week, seasonID, stadiumMap, stadiumMapByID, rivalryMap)
		games = append(games, g)
	}

	// emit is a safe wrapper for week-locked games: skips duplicates and busy weeks,
	// then delegates to recordGame.
	emit := func(home, away structs.CollegeTeam, week uint) {
		if home.ID == 0 || away.ID == 0 {
			return
		}
		if alreadyScheduled(home.ID, away.ID, gamesPlayedAgainstOpponentsMap) {
			return
		}
		if gamesPlayedByWeekMap[home.ID] != nil && gamesPlayedByWeekMap[home.ID][week] {
			return
		}
		if gamesPlayedByWeekMap[away.ID] != nil && gamesPlayedByWeekMap[away.ID][week] {
			return
		}
		recordGame(home, away, week)
	}

	processed := make(map[SchedulerHistoryKey]bool)

	// Phase 1: Guaranteed pairs — Week 14 locks first
	for _, key := range secGuaranteedPairs {
		if processed[key] {
			continue
		}
		a := teamMap[key.A]
		b := teamMap[key.B]
		if a.ID == 0 || b.ID == 0 {
			continue
		}
		// If the annual rivalry loop already created this game, honour it and
		// skip — no duplicate should be emitted.
		if alreadyScheduled(key.A, key.B, gamesPlayedAgainstOpponentsMap) {
			processed[key] = true
			continue
		}
		if confGameCount[key.A] >= secConferenceGamesPerTeam || confGameCount[key.B] >= secConferenceGamesPerTeam {
			processed[key] = true
			continue
		}
		var home, away structs.CollegeTeam
		if ShouldBeHome(key.A, key.B, season, lastHomeMap, homecountMap) {
			home, away = a, b
		} else {
			home, away = b, a
		}
		var week uint
		if secWeek14ConfLocks[key] {
			emit(home, away, 14)
		} else {
			week = assignWeek(home.ID, away.ID, 4, 13, gamesPlayedByWeekMap)
			if week == 0 {
				week = assignWeek(home.ID, away.ID, 1, 3, gamesPlayedByWeekMap)
			}
			if week == 0 {
				continue
			}
			recordGame(home, away, week)
		}
		_ = week
		processed[key] = true
	}

	// Phase 2: Mark Week 14 for OOC teams (those not already marked)
	for teamID := range secWeek14OOC {
		if gamesPlayedByWeekMap[teamID] == nil {
			gamesPlayedByWeekMap[teamID] = make(map[uint]bool)
		}
		gamesPlayedByWeekMap[teamID][14] = true
	}

	// Phase 3: Flexible games — global shuffled pair approach.
	// The SEC flex-rotation lists are asymmetric: Tennessee and Vanderbilt only
	// appear in their OWN lists, never in their desired opponents' lists. A
	// per-team loop therefore fills every opponent to cap before Tennessee/
	// Vanderbilt are reached, leaving them permanently short. Collecting ALL
	// pairs from the union of every team's flex list and shuffling them gives
	// every pair (including Tennessee's) an equal chance of being scheduled
	// early; the retry loop in the manager tries different shuffle seeds.
	// confGameCount is already live from Phase 1.

	// Build global flex pair set — include a pair if EITHER side lists the other.
	allFlexPairs := make(map[SchedulerHistoryKey]bool)
	for _, teamID := range secTeamIDs {
		var opponents []uint
		if useGroupA {
			opponents = secFlexibleGroupA[teamID]
		} else {
			opponents = secFlexGroupB(teamID)
		}
		for _, oppID := range opponents {
			key := makeHistoryKey(teamID, oppID)
			if !processed[key] {
				allFlexPairs[key] = true
			}
		}
	}

	flexPairSlice := make([]SchedulerHistoryKey, 0, len(allFlexPairs))
	for k := range allFlexPairs {
		flexPairSlice = append(flexPairSlice, k)
	}
	rand.Shuffle(len(flexPairSlice), func(i, j int) {
		flexPairSlice[i], flexPairSlice[j] = flexPairSlice[j], flexPairSlice[i]
	})
	// Sort by historical play count ascending so least-played flex pairs get priority.
	// Stable sort preserves the shuffle order within ties.
	sort.SliceStable(flexPairSlice, func(i, j int) bool {
		return playCountMap[flexPairSlice[i]] < playCountMap[flexPairSlice[j]]
	})

	for _, key := range flexPairSlice {
		aID, bID := key.A, key.B
		if alreadyScheduled(aID, bID, gamesPlayedAgainstOpponentsMap) {
			continue
		}
		if confGameCount[aID] >= secConferenceGamesPerTeam || confGameCount[bID] >= secConferenceGamesPerTeam {
			continue
		}
		tA := teamMap[aID]
		tB := teamMap[bID]
		if tA.ID == 0 || tB.ID == 0 {
			continue
		}
		var home, away structs.CollegeTeam
		if ShouldBeHome(aID, bID, season, lastHomeMap, homecountMap) {
			home, away = tA, tB
		} else {
			home, away = tB, tA
		}
		w := assignWeek(home.ID, away.ID, 4, 13, gamesPlayedByWeekMap)
		if w == 0 {
			w = assignWeek(home.ID, away.ID, 1, 3, gamesPlayedByWeekMap)
		}
		if w == 0 {
			continue
		}
		recordGame(home, away, w)
	}

	// Phase 4: Rescue any SEC team still below 9 conference games.
	for _, teamID := range secTeamIDs {
		if confGameCount[teamID] >= secConferenceGamesPerTeam {
			continue
		}
		for _, oppID := range secTeamIDs {
			if oppID == teamID {
				continue
			}
			if confGameCount[teamID] >= secConferenceGamesPerTeam {
				break
			}
			if confGameCount[oppID] >= secConferenceGamesPerTeam {
				continue
			}
			if alreadyScheduled(teamID, oppID, gamesPlayedAgainstOpponentsMap) {
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
			w := assignWeek(home.ID, away.ID, 1, 13, gamesPlayedByWeekMap)
			if w == 0 {
				continue
			}
			recordGame(home, away, w)
		}
	}

	// Retry signal: return empty slice if any team is still short so the
	// manager's retry loop can try a different shuffle ordering.
	for _, id := range secTeamIDs {
		if confGameCount[id] < secConferenceGamesPerTeam {
			return []structs.CollegeGame{}
		}
	}

	return games
}
