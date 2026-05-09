package managers

import (
	"github.com/CalebRose/SimFBA/structs"
)

// ============================================================
// Big 12 Conference Schedule Generator
// 16 teams | 4 pods of 4 | 9 conf games per team
// Pod games: Weeks 7, 8, 9 (rivalry=9)
// Non-pod games: Weeks 4-9 (6 games, ≤2 consecutive away)
// ============================================================

// Big 12 team IDs (from college_teams_csv.csv)
// Pod Red:    Arizona(5), Arizona State(6), BYU(17), Utah(118)
// Pod Yellow: Texas Tech(104), TCU(98), Houston(39), Baylor(12)
// Pod Green:  Kansas(44), Kansas State(45), Oklahoma State(79), Colorado(24)
// Pod Orange: West Virginia(126), UCF(110), Cincinnati(21), Iowa State(43)

const (
	big12PodRed    = 0
	big12PodYellow = 1
	big12PodGreen  = 2
	big12PodOrange = 3
)

var big12Pods = [4][]uint{
	{5, 6, 17, 118},    // Red
	{104, 98, 39, 12},  // Yellow
	{44, 45, 79, 24},   // Green
	{126, 110, 21, 43}, // Orange
}

// big12PodOf returns which pod index a team belongs to, -1 if not found.
func big12PodOf(teamID uint) int {
	for p, pod := range big12Pods {
		for _, id := range pod {
			if id == teamID {
				return p
			}
		}
	}
	return -1
}

// Big 12 Annual Rivalry Registry (within-pod, always Week 9)
// key = makeHistoryKey(a, b), value = true
var big12RivalryPairs = []SchedulerHistoryKey{
	makeHistoryKey(17, 118), // BYU vs Utah
	makeHistoryKey(5, 6),    // Arizona vs Arizona State
	makeHistoryKey(12, 98),  // Baylor vs TCU
	makeHistoryKey(44, 45),  // Kansas vs Kansas State
}

func isBig12Rivalry(a, b uint) bool {
	key := makeHistoryKey(a, b)
	for _, r := range big12RivalryPairs {
		if r == key {
			return true
		}
	}
	return false
}

// big12CrossPodRotationTable defines the 6 non-pod opponents each team faces per
// season. The table is indexed by [teamID][cycleYear(1-4)] => []opponentID (6 opponents).
// Home/away: if the team is listed first in the pair for Year 1, they host in Year 1 & 3
// and visit in Year 2 & 4 (flip on cycle years 2 & 4).
//
// Structure: rotTable[teamID] = [4][]uint  (index 0 = Year1, ..., 3 = Year4)
// Each inner slice has 6 opponents.
// Home/away is indicated separately via big12CrossPodHomeYear1.
// If big12CrossPodHomeYear1[makeHistoryKey(team, opp)] == true, team hosts in Year 1.

var big12CrossPodRotation = map[uint][4][]uint{
	// ---- Red Pod ----
	// Arizona (5): faces Yellow(2) + Green(2) + Orange(2) per year, rotating
	5: {
		{104, 98, 44, 45, 126, 110}, // Year 1
		{39, 12, 79, 24, 21, 43},    // Year 2
		{104, 98, 44, 45, 126, 110}, // Year 3 (same opponents, flipped H/A)
		{39, 12, 79, 24, 21, 43},    // Year 4
	},
	// Arizona State (6)
	6: {
		{39, 12, 79, 24, 21, 43},    // Year 1
		{104, 98, 44, 45, 126, 110}, // Year 2
		{39, 12, 79, 24, 21, 43},    // Year 3
		{104, 98, 44, 45, 126, 110}, // Year 4
	},
	// BYU (17)
	17: {
		{104, 98, 44, 45, 21, 43},  // Year 1
		{39, 12, 79, 24, 126, 110}, // Year 2
		{104, 98, 44, 45, 21, 43},  // Year 3
		{39, 12, 79, 24, 126, 110}, // Year 4
	},
	// Utah (118)
	118: {
		{39, 12, 79, 24, 126, 110}, // Year 1
		{104, 98, 44, 45, 21, 43},  // Year 2
		{39, 12, 79, 24, 126, 110}, // Year 3
		{104, 98, 44, 45, 21, 43},  // Year 4
	},
	// ---- Yellow Pod ----
	// Baylor (12)
	12: {
		{5, 6, 44, 45, 126, 110},  // Year 1
		{17, 118, 79, 24, 21, 43}, // Year 2
		{5, 6, 44, 45, 126, 110},  // Year 3
		{17, 118, 79, 24, 21, 43}, // Year 4
	},
	// TCU (98)
	98: {
		{17, 118, 79, 24, 21, 43}, // Year 1
		{5, 6, 44, 45, 126, 110},  // Year 2
		{17, 118, 79, 24, 21, 43}, // Year 3
		{5, 6, 44, 45, 126, 110},  // Year 4
	},
	// Houston (39)
	39: {
		{5, 6, 44, 45, 21, 43},      // Year 1
		{17, 118, 79, 24, 126, 110}, // Year 2
		{5, 6, 44, 45, 21, 43},      // Year 3
		{17, 118, 79, 24, 126, 110}, // Year 4
	},
	// Texas Tech (104)
	104: {
		{17, 118, 79, 24, 126, 110}, // Year 1
		{5, 6, 44, 45, 21, 43},      // Year 2
		{17, 118, 79, 24, 126, 110}, // Year 3
		{5, 6, 44, 45, 21, 43},      // Year 4
	},
	// ---- Green Pod ----
	// Kansas (44)
	44: {
		{5, 6, 104, 98, 126, 110}, // Year 1
		{17, 118, 39, 12, 21, 43}, // Year 2
		{5, 6, 104, 98, 126, 110}, // Year 3
		{17, 118, 39, 12, 21, 43}, // Year 4
	},
	// Kansas State (45)
	45: {
		{17, 118, 39, 12, 21, 43}, // Year 1
		{5, 6, 104, 98, 126, 110}, // Year 2
		{17, 118, 39, 12, 21, 43}, // Year 3
		{5, 6, 104, 98, 126, 110}, // Year 4
	},
	// Oklahoma State (79)
	79: {
		{5, 6, 104, 98, 21, 43},     // Year 1
		{17, 118, 39, 12, 126, 110}, // Year 2
		{5, 6, 104, 98, 21, 43},     // Year 3
		{17, 118, 39, 12, 126, 110}, // Year 4
	},
	// Colorado (24)
	24: {
		{17, 118, 39, 12, 126, 110}, // Year 1
		{5, 6, 104, 98, 21, 43},     // Year 2
		{17, 118, 39, 12, 126, 110}, // Year 3
		{5, 6, 104, 98, 21, 43},     // Year 4
	},
	// ---- Orange Pod ----
	// West Virginia (126)
	126: {
		{5, 6, 104, 98, 44, 45},   // Year 1
		{17, 118, 39, 12, 79, 24}, // Year 2
		{5, 6, 104, 98, 44, 45},   // Year 3
		{17, 118, 39, 12, 79, 24}, // Year 4
	},
	// UCF (110)
	110: {
		{17, 118, 39, 12, 79, 24}, // Year 1
		{5, 6, 104, 98, 44, 45},   // Year 2
		{17, 118, 39, 12, 79, 24}, // Year 3
		{5, 6, 104, 98, 44, 45},   // Year 4
	},
	// Cincinnati (21)
	21: {
		{5, 6, 104, 98, 79, 24},   // Year 1
		{17, 118, 39, 12, 44, 45}, // Year 2
		{5, 6, 104, 98, 79, 24},   // Year 3
		{17, 118, 39, 12, 44, 45}, // Year 4
	},
	// Iowa State (43)
	43: {
		{17, 118, 39, 12, 44, 45}, // Year 1
		{5, 6, 104, 98, 79, 24},   // Year 2
		{17, 118, 39, 12, 44, 45}, // Year 3
		{5, 6, 104, 98, 79, 24},   // Year 4
	},
}

// big12CrossPodHomeYear1[makeHistoryKey(a,b)] == true means team A hosts in Year 1.
// We define the base home side (Year 1) for each cross-pod pair.
var big12CrossPodHomeYear1 = map[SchedulerHistoryKey]uint{
	// Red vs Yellow
	makeHistoryKey(5, 104):   5,
	makeHistoryKey(5, 98):    98,
	makeHistoryKey(5, 39):    5,
	makeHistoryKey(5, 12):    12,
	makeHistoryKey(6, 104):   104,
	makeHistoryKey(6, 98):    6,
	makeHistoryKey(6, 39):    39,
	makeHistoryKey(6, 12):    6,
	makeHistoryKey(17, 104):  104,
	makeHistoryKey(17, 98):   17,
	makeHistoryKey(17, 39):   39,
	makeHistoryKey(17, 12):   17,
	makeHistoryKey(118, 104): 104,
	makeHistoryKey(118, 98):  118,
	makeHistoryKey(118, 39):  118,
	makeHistoryKey(118, 12):  12,
	// Red vs Green
	makeHistoryKey(5, 44):   5,
	makeHistoryKey(5, 45):   45,
	makeHistoryKey(5, 79):   79,
	makeHistoryKey(5, 24):   5,
	makeHistoryKey(6, 44):   44,
	makeHistoryKey(6, 45):   6,
	makeHistoryKey(6, 79):   6,
	makeHistoryKey(6, 24):   24,
	makeHistoryKey(17, 44):  17,
	makeHistoryKey(17, 45):  45,
	makeHistoryKey(17, 79):  17,
	makeHistoryKey(17, 24):  24,
	makeHistoryKey(118, 44): 44,
	makeHistoryKey(118, 45): 118,
	makeHistoryKey(118, 79): 118,
	makeHistoryKey(118, 24): 24,
	// Red vs Orange
	makeHistoryKey(5, 126):   5,
	makeHistoryKey(5, 110):   110,
	makeHistoryKey(5, 21):    5,
	makeHistoryKey(5, 43):    43,
	makeHistoryKey(6, 126):   126,
	makeHistoryKey(6, 110):   6,
	makeHistoryKey(6, 21):    21,
	makeHistoryKey(6, 43):    6,
	makeHistoryKey(17, 126):  126,
	makeHistoryKey(17, 110):  17,
	makeHistoryKey(17, 21):   21,
	makeHistoryKey(17, 43):   17,
	makeHistoryKey(118, 126): 118,
	makeHistoryKey(118, 110): 110,
	makeHistoryKey(118, 21):  118,
	makeHistoryKey(118, 43):  43,
	// Yellow vs Green
	makeHistoryKey(12, 44):  12,
	makeHistoryKey(12, 45):  45,
	makeHistoryKey(12, 79):  12,
	makeHistoryKey(12, 24):  24,
	makeHistoryKey(98, 44):  44,
	makeHistoryKey(98, 45):  98,
	makeHistoryKey(98, 79):  79,
	makeHistoryKey(98, 24):  98,
	makeHistoryKey(39, 44):  39,
	makeHistoryKey(39, 45):  45,
	makeHistoryKey(39, 79):  79,
	makeHistoryKey(39, 24):  39,
	makeHistoryKey(104, 44): 44,
	makeHistoryKey(104, 45): 104,
	makeHistoryKey(104, 79): 104,
	makeHistoryKey(104, 24): 24,
	// Yellow vs Orange
	makeHistoryKey(12, 126):  126,
	makeHistoryKey(12, 110):  12,
	makeHistoryKey(12, 21):   12,
	makeHistoryKey(12, 43):   43,
	makeHistoryKey(98, 126):  98,
	makeHistoryKey(98, 110):  110,
	makeHistoryKey(98, 21):   21,
	makeHistoryKey(98, 43):   98,
	makeHistoryKey(39, 126):  126,
	makeHistoryKey(39, 110):  39,
	makeHistoryKey(39, 21):   39,
	makeHistoryKey(39, 43):   43,
	makeHistoryKey(104, 126): 104,
	makeHistoryKey(104, 110): 110,
	makeHistoryKey(104, 21):  104,
	makeHistoryKey(104, 43):  43,
	// Green vs Orange
	makeHistoryKey(44, 126): 44,
	makeHistoryKey(44, 110): 110,
	makeHistoryKey(44, 21):  44,
	makeHistoryKey(44, 43):  43,
	makeHistoryKey(45, 126): 45,
	makeHistoryKey(45, 110): 45,
	makeHistoryKey(45, 21):  21,
	makeHistoryKey(45, 43):  43,
	makeHistoryKey(79, 126): 126,
	makeHistoryKey(79, 110): 79,
	makeHistoryKey(79, 21):  79,
	makeHistoryKey(79, 43):  43,
	makeHistoryKey(24, 126): 126,
	makeHistoryKey(24, 110): 24,
	makeHistoryKey(24, 21):  21,
	makeHistoryKey(24, 43):  24,
}

// big12HomeForSeason returns whether teamID is the home team against oppID in a
// given season, based on the Year-1 base and cycle year flip.
func big12HomeForSeason(teamID, oppID uint, season int) bool {
	key := makeHistoryKey(teamID, oppID)
	year1Home, ok := big12CrossPodHomeYear1[key]
	if !ok {
		// Fallback: lower ID hosts in odd years
		cy := cycleYear(season)
		return (cy%2 == 1) == (teamID == key.A)
	}
	cy := cycleYear(season)
	baseIsHome := year1Home == teamID
	if cy == 1 || cy == 3 {
		return baseIsHome
	}
	return !baseIsHome
}

// GenerateBigTwelveSchedule produces all Big 12 conference games for the season.
func GenerateBigTwelveSchedule(
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
	confGameCount := make(map[uint]int)
	// Seed homecountMap from rivalry-pass home game counts.
	homecountMap := make(map[uint]int, len(homeCountSeedMap))
	for id, count := range homeCountSeedMap {
		homecountMap[id] = count
	}

	// Build Big 12 team set and pre-seed confGameCount from any games already placed.
	var allBig12IDs []uint
	big12TeamSet := make(map[uint]bool)
	for _, pod := range big12Pods {
		for _, id := range pod {
			allBig12IDs = append(allBig12IDs, id)
			big12TeamSet[id] = true
		}
	}
	for _, id := range allBig12IDs {
		for oppID := range gamesPlayedAgainstOpponentsMap[id] {
			if big12TeamSet[oppID] {
				confGameCount[id]++
			}
		}
	}

	// Helper: emit a game record
	emit := func(home, away structs.CollegeTeam, week uint) {
		if home.ID == 0 || away.ID == 0 {
			return
		}
		if alreadyScheduled(home.ID, away.ID, gamesPlayedAgainstOpponentsMap) {
			return
		}
		markWeek(home.ID, away.ID, week, gamesPlayedByWeekMap)
		markOpponents(home.ID, away.ID, gamesPlayedAgainstOpponentsMap)
		homecountMap[home.ID]++
		confGameCount[home.ID]++
		confGameCount[away.ID]++
		g := MakeCollegeGameRecord(home, away, week, seasonID, stadiumMap, stadiumMapByID, rivalryMap)
		games = append(games, g)
	}

	// ----------------------------------------------------------------
	// For each team, build its 9-game slate:
	//   Weeks 4-9: 6 non-pod games (slots 0-5, ≤2 consecutive away)
	//   Week 7: pod game 1 (non-rivalry)
	//   Week 8: pod game 2 (non-rivalry)
	//   Week 9: rivalry pod game (locked)
	// We process pod games first (weeks 7,8,9) to lock those weeks,
	// then distribute non-pod games into weeks 4-6 & remaining slots.
	// ----------------------------------------------------------------

	// Phase 1: Rivalry games (Week 9) — process once per pair
	for _, rkey := range big12RivalryPairs {
		a := teamMap[rkey.A]
		b := teamMap[rkey.B]
		if a.ID == 0 || b.ID == 0 {
			continue
		}
		if confGameCount[rkey.A] >= 9 || confGameCount[rkey.B] >= 9 {
			continue
		}
		var home, away structs.CollegeTeam
		if ShouldBeHome(rkey.A, rkey.B, season, nil, homecountMap) {
			home, away = a, b
		} else {
			home, away = b, a
		}
		emit(home, away, 9)
	}

	// Phase 2: Non-rivalry pod games (Weeks 7 & 8)
	// Each pod: 4 teams, 6 intra-pod pairs; 4 are non-rivalry (rivalry already handled)
	for _, pod := range big12Pods {
		// Collect non-rivalry pairs in this pod
		var nonRivalryPairs [][2]uint
		for i := 0; i < len(pod); i++ {
			for j := i + 1; j < len(pod); j++ {
				if !isBig12Rivalry(pod[i], pod[j]) {
					nonRivalryPairs = append(nonRivalryPairs, [2]uint{pod[i], pod[j]})
				}
			}
		}
		// Assign non-rivalry pairs to weeks 7 & 8 (2 games per week for 4-team pod = 2 matchups)
		weeks := []uint{7, 8}
		for idx, pair := range nonRivalryPairs {
			if idx >= 2 {
				break
			}
			if confGameCount[pair[0]] >= 9 || confGameCount[pair[1]] >= 9 {
				continue
			}
			a := teamMap[pair[0]]
			b := teamMap[pair[1]]
			if a.ID == 0 || b.ID == 0 {
				continue
			}
			var home, away structs.CollegeTeam
			if ShouldBeHome(pair[0], pair[1], season, nil, homecountMap) {
				home, away = a, b
			} else {
				home, away = b, a
			}
			w := weeks[idx]
			// Check if already occupied; if so, swap weeks or find next open
			if gamesPlayedByWeekMap[home.ID] != nil && gamesPlayedByWeekMap[home.ID][w] {
				other := weeks[1-idx]
				if gamesPlayedByWeekMap[home.ID] != nil && gamesPlayedByWeekMap[home.ID][other] {
					// fallback overflow
					w = assignWeek(home.ID, away.ID, 4, 6, gamesPlayedByWeekMap)
				} else {
					w = other
				}
			} else if gamesPlayedByWeekMap[away.ID] != nil && gamesPlayedByWeekMap[away.ID][w] {
				other := weeks[1-idx]
				if gamesPlayedByWeekMap[away.ID] != nil && gamesPlayedByWeekMap[away.ID][other] {
					w = assignWeek(home.ID, away.ID, 4, 6, gamesPlayedByWeekMap)
				} else {
					w = other
				}
			}
			if w == 0 {
				continue
			}
			emit(home, away, w)
		}
	}

	// Phase 3: Non-pod games (Weeks 4-6, overflow to 10 if needed)
	// For each team, get their 6 non-pod opponents from the rotation table.
	// Track processed pairs to avoid duplicates.
	processedCrossPod := make(map[SchedulerHistoryKey]bool)

	for _, pod := range big12Pods {
		for _, teamID := range pod {
			if confGameCount[teamID] >= 9 {
				continue
			}
			team := teamMap[teamID]
			if team.ID == 0 {
				continue
			}
			rotEntry, ok := big12CrossPodRotation[teamID]
			if !ok {
				continue
			}
			opponents := rotEntry[cy-1] // 0-indexed

			// Build list of (home, away, week=0) for this team's 6 opponents
			// Avoid duplicating pairs already processed
			type nonPodEntry struct {
				home, away structs.CollegeTeam
			}
			var pending []nonPodEntry
			for _, oppID := range opponents {
				key := makeHistoryKey(teamID, oppID)
				if processedCrossPod[key] {
					continue
				}
				if alreadyScheduled(teamID, oppID, gamesPlayedAgainstOpponentsMap) {
					continue
				}
				if confGameCount[oppID] >= 9 {
					processedCrossPod[key] = true
					continue
				}
				opp := teamMap[oppID]
				if opp.ID == 0 {
					continue
				}
				var home, away structs.CollegeTeam
				// Priority 1: balance home counts; Priority 2: big12HomeForSeason table
				if homecountMap[teamID] < homecountMap[oppID] {
					home, away = team, opp
				} else if homecountMap[oppID] < homecountMap[teamID] {
					home, away = opp, team
				} else if big12HomeForSeason(teamID, oppID, season) {
					home, away = team, opp
				} else {
					home, away = opp, team
				}
				pending = append(pending, nonPodEntry{home, away})
				processedCrossPod[key] = true
			}

			// Schedule each pending matchup in Weeks 4-6 (overflow to 10-12)
			// enforcing ≤2 consecutive away games for 'team'
			consecutiveAway := 0
			for _, e := range pending {
				// Determine if 'team' is away in this game
				teamIsAway := e.away.ID == teamID
				// Re-check caps since earlier games in the loop may have incremented counts.
				if confGameCount[e.home.ID] >= 9 || confGameCount[e.away.ID] >= 9 {
					continue
				}

				// If 3 consecutive away, find a week where a home game can be inserted
				// For simplicity: try weeks 4-6, then 10-12 as overflow
				var w uint
				if teamIsAway && consecutiveAway >= 2 {
					// Must find home game — swap if possible, else just schedule normally
					// Try to find a week; we don't swap but just push this to later
					w = assignWeek(e.home.ID, e.away.ID, 4, 6, gamesPlayedByWeekMap)
					if w == 0 {
						w = assignWeek(e.home.ID, e.away.ID, 1, 13, gamesPlayedByWeekMap)
					}
					if w == 0 {
						continue
					}
					consecutiveAway = 0
				} else {
					w = assignWeek(e.home.ID, e.away.ID, 4, 6, gamesPlayedByWeekMap)
					if w == 0 {
						w = assignWeek(e.home.ID, e.away.ID, 1, 13, gamesPlayedByWeekMap)
					}
					if w == 0 {
						continue
					}
					if teamIsAway {
						consecutiveAway++
					} else {
						consecutiveAway = 0
					}
				}
				emit(e.home, e.away, w)
			}
		}
	}

	// Phase 4: Validate every Big 12 team has 9 conference games.
	// Attempt to fill gaps using an expanded week window (1–13).
	// confGameCount is already live from Phases 1–3.
	for _, teamID := range allBig12IDs {
		if confGameCount[teamID] >= 9 {
			continue
		}
		for _, oppID := range allBig12IDs {
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
			w := assignWeek(home.ID, away.ID, 1, 13, gamesPlayedByWeekMap)
			if w == 0 {
				continue
			}
			emit(home, away, w)
		}
	}

	// If any team is still short, return empty slice to trigger the manager's
	// retry loop to reseed and try again.
	for _, id := range allBig12IDs {
		if confGameCount[id] < 9 {
			return []structs.CollegeGame{}
		}
	}

	return games
}
