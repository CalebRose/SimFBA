package managers

import (
	"github.com/CalebRose/SimFBA/structs"
)

// ============================================================
// Pac-12 Conference Schedule Generator
// 8 teams | full round-robin (7 conf games each) | 5 OOC
// Locked games:
//   Boise State/Fresno State — Week 14
//   SDSU/Fresno State        — Week 4
//   Oregon State/Wazzu       — Week 13
// ============================================================
//
// Pac-12 team IDs:
// Boise State(13), Colorado State(25), Fresno State(33),
// Oregon State(83), San Diego State(89), Texas State(103),
// Utah State(119), Washington State(125)

var pacTwelveTeamIDs = []uint{13, 25, 33, 83, 89, 103, 119, 125}

// pac12LockedMatchups: pairs with a specific week assignment
var pac12LockedMatchups = []struct {
	key  SchedulerHistoryKey
	week uint
}{
	{makeHistoryKey(13, 33), 14},  // Boise State vs Fresno State
	{makeHistoryKey(33, 89), 4},   // SDSU vs Fresno State
	{makeHistoryKey(83, 125), 13}, // Oregon State vs Washington State
}

// GeneratePacTwelveSchedule produces all Pac-12 conference games for the season.
func GeneratePacTwelveSchedule(
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
	lockedSet := make(map[SchedulerHistoryKey]uint)
	for _, l := range pac12LockedMatchups {
		lockedSet[l.key] = l.week
	}

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
	for _, l := range pac12LockedMatchups {
		a := teamMap[l.key.A]
		b := teamMap[l.key.B]
		if a.ID == 0 || b.ID == 0 {
			continue
		}
		var home, away structs.CollegeTeam
		if ShouldBeHome(l.key.A, l.key.B, season, nil, homecountMap) {
			home, away = a, b
		} else {
			home, away = b, a
		}
		emit(home, away, l.week)
	}

	// Phase 2: Full round-robin for remaining pairs
	for i := 0; i < len(pacTwelveTeamIDs); i++ {
		for j := i + 1; j < len(pacTwelveTeamIDs); j++ {
			a := pacTwelveTeamIDs[i]
			b := pacTwelveTeamIDs[j]
			key := makeHistoryKey(a, b)
			if _, locked := lockedSet[key]; locked {
				continue
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
			w := assignWeek(home.ID, away.ID, 5, 12, gamesPlayedByWeekMap)
			if w == 0 {
				w = assignWeek(home.ID, away.ID, 3, 4, gamesPlayedByWeekMap)
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

	return games
}
