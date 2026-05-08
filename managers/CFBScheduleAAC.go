package managers

import (
	"github.com/CalebRose/SimFBA/structs"
)

// ============================================================
// AAC Conference Schedule Generator
// 14 teams | 8 conf games per team
// Annual locks: Army/Navy (Week 14), Memphis/UAB (Week 14)
// OOC busy Week 14: Tulane, USF, UTSA
// Week 14 bye: ECU, FAU, Charlotte, North Texas, Rice, Temple, Tulsa
// ============================================================
//
// AAC team IDs:
// Army(9), Charlotte(20), ECU(27), FAU(31), Memphis(56), Navy(65),
// North Texas(72), Rice(87), South Florida(94), Temple(99),
// Tulane(107), Tulsa(108), UAB(109), UTSA(117)

var aacTeamIDs = []uint{9, 20, 27, 31, 56, 65, 72, 87, 94, 99, 107, 108, 109, 117}

// aacWeek14ConfLocks: conference games locked to Week 14
var aacWeek14ConfLocks = []SchedulerHistoryKey{
	makeHistoryKey(9, 65),   // Army vs Navy
	makeHistoryKey(56, 109), // Memphis vs UAB
}

// aacWeek14OOC: teams busy Week 14 with OOC rivals
var aacWeek14OOC = map[uint]bool{
	107: true, // Tulane     (vs Southern Miss)
	94:  true, // South Florida (vs UCF)
	117: true, // UTSA       (vs Texas State)
}

// aacWeek14Bye: teams with a hard bye in Week 14
var aacWeek14Bye = map[uint]bool{
	27:  true, // ECU
	31:  true, // FAU
	20:  true, // Charlotte
	72:  true, // North Texas
	87:  true, // Rice
	99:  true, // Temple
	108: true, // Tulsa
}

// GenerateAACSchedule produces all AAC conference games for the season.
func GenerateAACSchedule(
	collegeTeams []structs.CollegeTeam,
	stadiumMap map[uint]structs.Stadium,
	stadiumMapByID map[uint]structs.Stadium,
	rivalryMap map[uint][]structs.CollegeRival,
	gamesPlayedAgainstOpponentsMap map[uint]map[uint]bool,
	gamesPlayedByWeekMap map[uint]map[uint]bool,
	ts structs.Timestamp,
) []structs.CollegeGame {
	games := []structs.CollegeGame{}
	season := ts.Season
	seasonID := uint(ts.CollegeSeasonID)

	teamMap := buildTeamMapFromSlice(collegeTeams)

	homecountMap := make(map[uint]int)

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
		markWeek(home.ID, away.ID, week, gamesPlayedByWeekMap)
		markOpponents(home.ID, away.ID, gamesPlayedAgainstOpponentsMap)
		homecountMap[home.ID]++
		g := CreateCollegeGameRecord(home, away, week, seasonID, stadiumMap, stadiumMapByID, rivalryMap)
		games = append(games, g)
	}

	// Mark Week 14 busy for OOC and bye teams
	for teamID := range aacWeek14OOC {
		if gamesPlayedByWeekMap[teamID] == nil {
			gamesPlayedByWeekMap[teamID] = make(map[uint]bool)
		}
		gamesPlayedByWeekMap[teamID][14] = true
	}
	for teamID := range aacWeek14Bye {
		if gamesPlayedByWeekMap[teamID] == nil {
			gamesPlayedByWeekMap[teamID] = make(map[uint]bool)
		}
		gamesPlayedByWeekMap[teamID][14] = true
	}

	// Phase 1: Week 14 conference locks
	for _, key := range aacWeek14ConfLocks {
		a := teamMap[key.A]
		b := teamMap[key.B]
		if a.ID == 0 || b.ID == 0 {
			continue
		}
		var home, away structs.CollegeTeam
		if ShouldBeHome(key.A, key.B, season, nil, homecountMap) {
			home, away = a, b
		} else {
			home, away = b, a
		}
		emit(home, away, 14)
	}

	// Phase 2: Rotating partial round-robin for remaining 8 conf games per team
	aacTeamSet := make(map[uint]bool, len(aacTeamIDs))
	for _, id := range aacTeamIDs {
		aacTeamSet[id] = true
	}
	confGameCount := make(map[uint]int)
	// Seed from all already-scheduled AAC games (annual rivalries + Phase 1)
	for _, id := range aacTeamIDs {
		for oppID := range gamesPlayedAgainstOpponentsMap[id] {
			if aacTeamSet[oppID] {
				confGameCount[id]++
			}
		}
	}

	ids := sortedTeamIDs(collegeTeams)
	pairs := getRoundRobinPairs(ids)
	// Shuffle for variety
	for i := len(pairs) - 1; i > 0; i-- {
		j := i / 2
		pairs[i], pairs[j] = pairs[j], pairs[i]
	}

	for _, pair := range pairs {
		a, b := pair[0], pair[1]
		if alreadyScheduled(a, b, gamesPlayedAgainstOpponentsMap) {
			continue
		}
		if confGameCount[a] >= 8 || confGameCount[b] >= 8 {
			continue
		}
		tA := teamMap[a]
		tB := teamMap[b]
		if tA.ID == 0 || tB.ID == 0 {
			continue
		}
		var home, away structs.CollegeTeam
		if ShouldBeHome(a, b, season, nil, homecountMap) {
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
		markOpponents(home.ID, away.ID, gamesPlayedAgainstOpponentsMap)
		homecountMap[home.ID]++
		confGameCount[home.ID]++
		confGameCount[away.ID]++
		g := CreateCollegeGameRecord(home, away, w, seasonID, stadiumMap, stadiumMapByID, rivalryMap)
		games = append(games, g)
	}

	// Final validation: if any team is still short, return an empty slice to
	// signal the manager's retry loop to reseed and try again.
	for _, teamID := range aacTeamIDs {
		if confGameCount[teamID] < 8 {
			return []structs.CollegeGame{}
		}
	}

	return games
}
