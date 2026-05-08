package managers

import (
	"github.com/CalebRose/SimFBA/structs"
)

// ============================================================
// Sun Belt Conference Schedule Generator
// 14 teams | 2 divisions of 7 | 6 intra-div + 2 cross-div = 8 conf games
// East: App State(4), Coastal Carolina(23), Georgia Southern(35),
//       Georgia State(36), James Madison(131), Marshall(54), Old Dominion(80)
// West: Arkansas State(8), Louisiana(49), Louisiana Monroe(50),
//       Louisiana Tech(51), South Alabama(92), Southern Miss(95), Troy(106)
// ============================================================

var sunBeltEast = []uint{4, 23, 35, 36, 131, 54, 80}
var sunBeltWest = []uint{8, 49, 50, 51, 92, 95, 106}

// GenerateSunBeltSchedule produces all Sun Belt conference games for the season.
func GenerateSunBeltSchedule(
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

	scheduleDiv := func(div []uint) {
		for i := 0; i < len(div); i++ {
			for j := i + 1; j < len(div); j++ {
				a := div[i]
				b := div[j]
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
					w = assignWeek(home.ID, away.ID, 4, 4, gamesPlayedByWeekMap)
				}
				if w == 0 {
					return
				}
				markOpponents(home.ID, away.ID, gamesPlayedAgainstOpponentsMap)
				homecountMap[home.ID]++
				g := CreateCollegeGameRecord(home, away, w, seasonID, stadiumMap, stadiumMapByID, rivalryMap)
				games = append(games, g)
			}
		}
	}

	// Phase 1: Intra-divisional games (all 6 opponents in division)
	scheduleDiv(sunBeltEast)
	scheduleDiv(sunBeltWest)

	// Phase 2: Cross-divisional games (2 per team, rotating)
	// Each East team plays 2 West teams per season; cycle over 4 years to meet all 7.
	// Pair each East team with 2 West opponents using (teamIndex + cycleOffset) mod 7.
	cy := cycleYear(season)
	crossProcessed := make(map[SchedulerHistoryKey]bool)
	for eIdx, eTeam := range sunBeltEast {
		for k := 0; k < 2; k++ {
			wIdx := (eIdx + (cy-1)*2 + k) % len(sunBeltWest)
			wTeam := sunBeltWest[wIdx]
			key := makeHistoryKey(eTeam, wTeam)
			if crossProcessed[key] {
				continue
			}
			if alreadyScheduled(eTeam, wTeam, gamesPlayedAgainstOpponentsMap) {
				crossProcessed[key] = true
				continue
			}
			tE := teamMap[eTeam]
			tW := teamMap[wTeam]
			if tE.ID == 0 || tW.ID == 0 {
				continue
			}
			var home, away structs.CollegeTeam
			if ShouldBeHome(eTeam, wTeam, season, nil, homecountMap) {
				home, away = tE, tW
			} else {
				home, away = tW, tE
			}
			w := assignWeek(home.ID, away.ID, 5, 13, gamesPlayedByWeekMap)
			if w == 0 {
				w = assignWeek(home.ID, away.ID, 4, 4, gamesPlayedByWeekMap)
			}
			if w == 0 {
				crossProcessed[key] = true
				continue
			}
			crossProcessed[key] = true
			markOpponents(home.ID, away.ID, gamesPlayedAgainstOpponentsMap)
			homecountMap[home.ID]++
			g := CreateCollegeGameRecord(home, away, w, seasonID, stadiumMap, stadiumMapByID, rivalryMap)
			games = append(games, g)
		}
	}

	return games
}
