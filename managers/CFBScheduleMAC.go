package managers

import (
	"github.com/CalebRose/SimFBA/structs"
)

// ============================================================
// MAC Conference Schedule Generator
// 13 teams | 4 pods | each team plays all pod-mates every year
// Cross-pod rotation ensures every team meets at least once in 4 years
// Locked Week 14: Akron/Kent, BG/Toledo, CMU/WMU, Buffalo/UMass, Miami OH/Ohio
// Locked Week 13: BG/Sacramento State
// ============================================================
//
// MAC team IDs:
// Pod 1 (East-A): Akron(2), Buffalo(16), Kent State(46), UMass(113)
// Pod 2 (East-B): Ball State(11), Miami OH(58), Ohio(76)
// Pod 3 (West-A): Bowling Green(15), Sacramento State(164), Toledo(105)
// Pod 4 (West-B): CMU(19), EMU(28), WMU(128)

var macPods = [4][]uint{
	{2, 16, 46, 113}, // Pod 1
	{11, 58, 76},     // Pod 2
	{15, 164, 105},   // Pod 3
	{19, 28, 128},    // Pod 4
}

func macPodOf(teamID uint) int {
	for p, pod := range macPods {
		for _, id := range pod {
			if id == teamID {
				return p
			}
		}
	}
	return -1
}

// macWeek14Locks: conf game locks in Week 14
var macWeek14Locks = []SchedulerHistoryKey{
	makeHistoryKey(2, 46),   // Akron vs Kent State
	makeHistoryKey(15, 105), // Bowling Green vs Toledo
	makeHistoryKey(19, 128), // CMU vs WMU
	makeHistoryKey(16, 113), // Buffalo vs UMass
	makeHistoryKey(58, 76),  // Miami OH vs Ohio
}

// macWeek13Locks: conf game locks in Week 13
var macWeek13Locks = []SchedulerHistoryKey{
	makeHistoryKey(15, 164), // BG vs Sacramento State
}

// macCrossRotation defines which 2 additional pods each pod faces per cycle year.
// Index: [podIndex][cycleYear-1] => other pod indices to play cross games against
var macCrossRotation = [4][4][]int{
	{{1, 2}, {1, 3}, {2, 3}, {1, 2}}, // Pod 0 plays pods...
	{{0, 2}, {0, 3}, {0, 2}, {0, 3}}, // Pod 1
	{{0, 3}, {1, 3}, {0, 1}, {0, 3}}, // Pod 2
	{{1, 2}, {0, 2}, {0, 1}, {1, 2}}, // Pod 3
}

// GenerateMACSchedule produces all MAC conference games for the season.
func GenerateMACSchedule(
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
	cy := cycleYear(season)

	teamMap := buildTeamMapFromSlice(collegeTeams)

	// Seed homecountMap from rivalry-pass home game counts.
	homecountMap := make(map[uint]int, len(homeCountSeedMap))
	for id, count := range homeCountSeedMap {
		homecountMap[id] = count
	}

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
		g := MakeCollegeGameRecord(home, away, week, seasonID, stadiumMap, stadiumMapByID, rivalryMap)
		games = append(games, g)
	}

	// Phase 1: Locked games
	allLocked := make(map[SchedulerHistoryKey]bool)
	for _, key := range macWeek14Locks {
		allLocked[key] = true
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
	for _, key := range macWeek13Locks {
		allLocked[key] = true
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
		emit(home, away, 13)
	}

	// Phase 2: Intra-pod games (all pod-mates every season, non-locked)
	for _, pod := range macPods {
		for i := 0; i < len(pod); i++ {
			for j := i + 1; j < len(pod); j++ {
				a := pod[i]
				b := pod[j]
				key := makeHistoryKey(a, b)
				if allLocked[key] {
					continue // week already set
				}
				if alreadyScheduled(a, b, gamesPlayedAgainstOpponentsMap) {
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
				w := assignWeek(home.ID, away.ID, 5, 13, gamesPlayedByWeekMap)
				if w == 0 {
					w = assignWeek(home.ID, away.ID, 1, 4, gamesPlayedByWeekMap)
				}
				if w == 0 {
					continue
				}
				markOpponents(home.ID, away.ID, gamesPlayedAgainstOpponentsMap)
				homecountMap[home.ID]++
				g := MakeCollegeGameRecord(home, away, w, seasonID, stadiumMap, stadiumMapByID, rivalryMap)
				games = append(games, g)
			}
		}
	}

	// Phase 3: Cross-pod games based on rotation
	crossProcessed := make(map[SchedulerHistoryKey]bool)
	for podIdx, pod := range macPods {
		otherPodIndices := macCrossRotation[podIdx][cy-1]
		for _, oppPodIdx := range otherPodIndices {
			oppPod := macPods[oppPodIdx]
			for _, teamID := range pod {
				for _, oppID := range oppPod {
					key := makeHistoryKey(teamID, oppID)
					if crossProcessed[key] {
						continue
					}
					if alreadyScheduled(teamID, oppID, gamesPlayedAgainstOpponentsMap) {
						crossProcessed[key] = true
						continue
					}
					tA := teamMap[teamID]
					tB := teamMap[oppID]
					if tA.ID == 0 || tB.ID == 0 {
						continue
					}
					var home, away structs.CollegeTeam
					if ShouldBeHome(teamID, oppID, season, nil, homecountMap) {
						home, away = tA, tB
					} else {
						home, away = tB, tA
					}
					w := assignWeek(home.ID, away.ID, 5, 13, gamesPlayedByWeekMap)
					if w == 0 {
						w = assignWeek(home.ID, away.ID, 1, 13, gamesPlayedByWeekMap)
					}
					if w == 0 {
						crossProcessed[key] = true
						continue
					}
					crossProcessed[key] = true
					markOpponents(home.ID, away.ID, gamesPlayedAgainstOpponentsMap)
					homecountMap[home.ID]++
					g := MakeCollegeGameRecord(home, away, w, seasonID, stadiumMap, stadiumMapByID, rivalryMap)
					games = append(games, g)
				}
			}
		}
	}

	return games
}
