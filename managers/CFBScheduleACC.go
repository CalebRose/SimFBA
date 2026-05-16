package managers

import (
	"math/rand/v2"
	"sort"

	"github.com/CalebRose/SimFBA/structs"
)

// ============================================================
// ACC Conference Schedule Generator
// 17 teams | 9 conf games (8 for the designated team)
// Tobacco Road group plays each other every season
// Week 14 conf locks + OOC locks
// ============================================================
//
// ACC team IDs:
// Boston College(14), California(18), Clemson(22), Duke(26),
// Florida State(32), Georgia Tech(37), Louisville(52), Miami FL(57),
// NC State(66), North Carolina(71), Pittsburgh(85), SMU(91),
// Stanford(96), Syracuse(97), Virginia(121), Virginia Tech(122),
// Wake Forest(123)

// accTeamIDs for convenience
var accTeamIDs = []uint{14, 18, 22, 26, 32, 37, 52, 57, 66, 71, 85, 91, 96, 97, 121, 122, 123}

// accWeek14ConfLocks: conference games always in Week 14
var accWeek14ConfLocks = []SchedulerHistoryKey{
	makeHistoryKey(18, 96),   // California vs Stanford
	makeHistoryKey(71, 66),   // North Carolina vs NC State
	makeHistoryKey(97, 14),   // Syracuse vs Boston College
	makeHistoryKey(121, 122), // Virginia vs Virginia Tech
}

// accWeek14OOCTeams: teams busy Week 14 with OOC rivals (do not schedule conf game)
var accWeek14OOCTeams = map[uint]bool{
	32: true, // Florida State  (vs Florida)
	22: true, // Clemson        (vs South Carolina)
	37: true, // Georgia Tech   (vs Georgia)
	52: true, // Louisville     (vs Kentucky)
	85: true, // Pittsburgh     (vs West Virginia)
}

// accWeek14HardBye: teams with neither conf lock nor OOC in Week 14
var accWeek14HardBye = map[uint]bool{
	57:  true, // Miami
	66:  true, // NC State (already plays UNC in W14 conf lock)
	91:  true, // SMU
	123: true, // Wake Forest
}

// accEightGameTeamHistory: tracks which team has been given 8 instead of 9 games.
// In practice, we pick the team least-recently selected. For simulation, we derive
// from game count — the team with the most conf games in recent history gets
// the "rest" of 8 games. We always use the simple heuristic: rotate annually
// through teams sorted by ID, picking the one whose season%len cycle points to them.
func accEightGameTeamForSeason(season int) uint {
	idx := season % len(accTeamIDs)
	return accTeamIDs[idx]
}

// GenerateACCSchedule produces all ACC conference games for the season.
// Returns an empty slice if Phase 4 cannot fill all required games, which signals
// the manager's retry loop to reseed and try again.
func GenerateACCSchedule(
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
	eightGameTeam := accEightGameTeamForSeason(season)

	teamMap := buildTeamMapFromSlice(collegeTeams)

	// Per-team max conference game counts: 9 for all teams, 8 for the designated team.
	maxGameMap := make(map[uint]int, len(accTeamIDs))
	for _, id := range accTeamIDs {
		maxGameMap[id] = 9
	}
	maxGameMap[eightGameTeam] = 8

	// Build conference team set for membership checks.
	accTeamSet := make(map[uint]bool, len(accTeamIDs))
	for _, id := range accTeamIDs {
		accTeamSet[id] = true
	}

	// Pre-seed confGameCount from any conference games already placed before this
	// generator runs (e.g. rivalry games between ACC teammates from the rivalry pass).
	confGameCount := make(map[uint]int)
	for _, teamID := range accTeamIDs {
		for oppID := range gamesPlayedAgainstOpponentsMap[teamID] {
			if accTeamSet[oppID] {
				confGameCount[teamID]++
			}
		}
	}

	// Seed homecountMap from rivalry-pass home game counts so ShouldBeHome sees
	// correct context from the very first conference game assignment.
	homecountMap := make(map[uint]int, len(homeCountSeedMap))
	for id, count := range homeCountSeedMap {
		homecountMap[id] = count
	}

	// recordGame commits a game for a week that has already been secured (either by
	// assignWeek or by an explicit free-slot check). It updates all tracking maps.
	// Callers must NOT call assignWeek and then pass that week to recordGame —
	// assignWeek already marks the week, and recordGame marks it again (idempotent).
	// For Phase 1 locked games, callers verify the week is free before calling.
	recordGame := func(home, away structs.CollegeTeam, week uint) {
		markWeek(home.ID, away.ID, week, gamesPlayedByWeekMap)
		markOpponents(home.ID, away.ID, gamesPlayedAgainstOpponentsMap)
		homecountMap[home.ID]++
		confGameCount[home.ID]++
		confGameCount[away.ID]++
		g := MakeCollegeGameRecord(home, away, week, seasonID, stadiumMap, stadiumMapByID, rivalryMap)
		games = append(games, g)
	}

	// Mark Week 14 unavailable for OOC and hard-bye teams before any phase runs.
	// This must happen before Phase 0 so that assignWeek in Phase 0 never claims
	// week 14 for any team that has an OOC or hard-bye constraint there.
	for teamID := range accWeek14OOCTeams {
		if gamesPlayedByWeekMap[teamID] == nil {
			gamesPlayedByWeekMap[teamID] = make(map[uint]bool)
		}
		gamesPlayedByWeekMap[teamID][14] = true
	}
	for teamID := range accWeek14HardBye {
		if gamesPlayedByWeekMap[teamID] == nil {
			gamesPlayedByWeekMap[teamID] = make(map[uint]bool)
		}
		gamesPlayedByWeekMap[teamID][14] = true
	}

	// Phase 0: Schedule the 8-game team's full conference slate before any other
	// phase runs. This eliminates the cap-check asymmetry that causes ripple failures
	// when the 8-game team's pairs are processed late in the shuffled round-robin.
	// Once this phase completes, all subsequent phases can skip that team entirely.
	//
	// Week 14 conf-locked opponents are skipped here — Phase 1 owns those pairs and
	// must place them at the correct week.
	{
		// Build the set of opponents the 8-game team has locked in Week 14 so we
		// don't consume that pair in the floating window and cause Phase 1 to miss it.
		week14LockedOpponents := make(map[uint]bool)
		for _, key := range accWeek14ConfLocks {
			if key.A == eightGameTeam {
				week14LockedOpponents[key.B] = true
			} else if key.B == eightGameTeam {
				week14LockedOpponents[key.A] = true
			}
		}

		opponents := make([]uint, 0, len(accTeamIDs)-1)
		for _, id := range accTeamIDs {
			if id != eightGameTeam {
				opponents = append(opponents, id)
			}
		}
		rand.Shuffle(len(opponents), func(i, j int) {
			opponents[i], opponents[j] = opponents[j], opponents[i]
		})
		eightObj := teamMap[eightGameTeam]
		for _, oppID := range opponents {
			if confGameCount[eightGameTeam] >= maxGameMap[eightGameTeam] {
				break
			}
			if week14LockedOpponents[oppID] {
				continue // reserved for Phase 1
			}
			if alreadyScheduled(eightGameTeam, oppID, gamesPlayedAgainstOpponentsMap) {
				continue
			}
			oppObj := teamMap[oppID]
			if oppObj.ID == 0 || eightObj.ID == 0 {
				continue
			}
			var home, away structs.CollegeTeam
			if ShouldBeHome(eightGameTeam, oppID, season, lastHomeMap, homecountMap) {
				home, away = eightObj, oppObj
			} else {
				home, away = oppObj, eightObj
			}
			w := assignWeek(home.ID, away.ID, 3, 13, gamesPlayedByWeekMap)
			if w == 0 {
				continue
			}
			recordGame(home, away, w)
		}
	}

	// Phase 1: Week 14 conference locks.
	// Explicitly verify week 14 is free for both teams before recording.
	for _, key := range accWeek14ConfLocks {
		a := teamMap[key.A]
		b := teamMap[key.B]
		if a.ID == 0 || b.ID == 0 {
			continue
		}
		if alreadyScheduled(a.ID, b.ID, gamesPlayedAgainstOpponentsMap) {
			continue
		}
		if confGameCount[a.ID] >= maxGameMap[a.ID] || confGameCount[b.ID] >= maxGameMap[b.ID] {
			continue
		}
		if gamesPlayedByWeekMap[a.ID] != nil && gamesPlayedByWeekMap[a.ID][14] {
			continue
		}
		if gamesPlayedByWeekMap[b.ID] != nil && gamesPlayedByWeekMap[b.ID][14] {
			continue
		}
		var home, away structs.CollegeTeam
		if ShouldBeHome(key.A, key.B, season, lastHomeMap, homecountMap) {
			home, away = a, b
		} else {
			home, away = b, a
		}
		recordGame(home, away, 14)
	}

	// Phase 3: Remaining conference games via shuffled round-robin.
	// Tobacco Road pairs (Duke, UNC, NC State, Wake Forest) are included in the
	// general pool here rather than in a dedicated phase, so homecountMap has
	// proper context when their home/away is decided and Duke's assignments are
	// naturally balanced against non-Tobacco-Road games.
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
		if alreadyScheduled(a, b, gamesPlayedAgainstOpponentsMap) {
			continue
		}
		if confGameCount[a] >= maxGameMap[a] || confGameCount[b] >= maxGameMap[b] {
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
	// attempt to fill any remaining gaps using an expanded week window (1–13).
	// This catches cases where Phase 3's pair ordering left some teams short.
	for _, teamID := range accTeamIDs {
		if confGameCount[teamID] >= maxGameMap[teamID] {
			continue
		}
		for _, oppID := range accTeamIDs {
			if oppID == teamID {
				continue
			}
			if confGameCount[teamID] >= maxGameMap[teamID] {
				break
			}
			if confGameCount[oppID] >= maxGameMap[oppID] {
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

	// Final validation: if any team still hasn't reached its target, return an
	// empty slice to signal the manager's retry loop to reseed and try again.
	for _, teamID := range accTeamIDs {
		if confGameCount[teamID] < maxGameMap[teamID] {
			return []structs.CollegeGame{}
		}
	}

	return games
}
