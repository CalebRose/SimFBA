package managers

import (
	"github.com/CalebRose/SimFBA/structs"
)

// ============================================================
// Small / FCS Conference Schedule Generators
// Handles conference IDs: 14 (MVFC), 15 (Patriot), 16 (Big Sky),
// 17 (Ivy), 18 (SoCon), 19 (SWAC), 20 (Big South-OVC), 21 (CAA),
// 23 (MEAC), 24 (NEC), 25 (Pioneer), 26 (Southland), 27 (UAC)
// ============================================================

// Per-conference game counts
var smallConfGameCount = map[int]int{
	14: 8, // MVFC
	15: 7, // Patriot
	16: 8, // Big Sky
	17: 7, // Ivy
	18: 9, // SoCon
	19: 8, // SWAC
	20: 7, // Big South-OVC
	21: 8, // CAA
	23: 5, // MEAC
	24: 8, // NEC
	25: 8, // Pioneer
	26: 8, // Southland
	27: 7, // UAC
}

// smallConfLockedMatchups defines intra-conference rivalry weeks per conference.
// key = conferenceID, value = list of locked matchups with their weeks.
var smallConfLockedMatchups = map[int][]struct {
	key  SchedulerHistoryKey
	week uint
}{
	// Mac (13)
	10: {
		{makeHistoryKey(2, 46), 14},   // Akron vs Kent State
		{makeHistoryKey(15, 105), 14}, // Bowling Green vs Toledo
		{makeHistoryKey(19, 128), 14}, // CMU vs WMU
		{makeHistoryKey(16, 113), 14}, // Buffalo vs UMass
		{makeHistoryKey(58, 76), 14},  // Miami OH vs Ohio
		{makeHistoryKey(11, 164), 13}, // BG vs Sacramento State
		{makeHistoryKey(19, 28), 4},   // CMU vs EMU
		{makeHistoryKey(16, 2), 4},    // Buffalo vs Akron
		{makeHistoryKey(11, 58), 4},   // Ball State vs Miami (OH)
		{makeHistoryKey(113, 46), 4},  // UMass vs Kent State
		{makeHistoryKey(15, 76), 4},   // Bowling Green vs Ohio
		{makeHistoryKey(128, 28), 8},  // WMU vs EMU
		{makeHistoryKey(11, 76), 8},   // Ball State vs Ohio
	},
	// Big Sky (16)
	16: {
		// UC Davis(154) vs Cal Poly(155) — Week 12
		// Eastern Washington(156) vs Montana(159) — Week 12
		// Eastern Washington(156) vs Portland State(163) — Week 13
		// Idaho(157) vs Idaho State(158) — Week 13
		// Idaho(157) vs Montana(159) — Week 8
		// Montana(159) vs Montana State(160) — Week 13
		// Northern Arizona(161) vs Southern Utah(260) — Week 12
		// Weber State(165) vs Southern Utah(260) — Week 5
		// Southern Utah(260) vs Utah Tech(263) — Week 8
		{makeHistoryKey(154, 155), 12},
		{makeHistoryKey(156, 159), 12},
		{makeHistoryKey(156, 163), 13},
		{makeHistoryKey(157, 158), 13},
		{makeHistoryKey(157, 159), 8},
		{makeHistoryKey(159, 160), 13},
		{makeHistoryKey(161, 260), 12},
		{makeHistoryKey(165, 260), 5},
		{makeHistoryKey(260, 263), 8},
	},
	// MVFC (14)
	14: {
		// South Dakota State(144) vs South Dakota(145) — Week 11
		// North Dakota(142) vs South Dakota(145) — Week 8
		{makeHistoryKey(144, 145), 11},
		{makeHistoryKey(142, 145), 8},
	},
	// CAA (21)
	21: {
		// Albany(135) vs Stony Brook(197) — Week 12
		// Maine(158) vs New Hampshire(184) — Week 13
		{makeHistoryKey(135, 197), 12},
		{makeHistoryKey(158, 184), 13},
	},
	// SoCon (18)
	18: {
		// VMI(202) vs The Citadel(198) — Week 13
		// The Citadel(198) vs Furman(153) — Week 10
		{makeHistoryKey(202, 198), 13},
		{makeHistoryKey(198, 153), 10},
	},
	// Big South-OVC (20)
	20: {
		// Tennessee State(199) vs UT Martin(201) — Week 10
		// Western Illinois(203) vs Eastern Illinois(151) — Week 12
		{makeHistoryKey(199, 201), 10},
		{makeHistoryKey(203, 151), 12},
	},
	// Ivy League (17)
	17: {
		// Harvard(154) vs Yale(205) — Week 14
		// Penn(183) vs Princeton(186) — Week 14
		{makeHistoryKey(154, 205), 14},
		{makeHistoryKey(183, 186), 14},
	},
	// Southland (26)
	26: {
		// Stephen F. Austin(174) vs Northwestern State(163) — Week 12
		// Lamar(157) vs McNeese(172) — Week 12
		// Northwestern State(163) vs McNeese(172) — Week 4
		// Northwestern State(163) vs Nicholls(185) — Week 8
		// Northwestern State(163) vs Southeastern Louisiana(192) — Week 8
		{makeHistoryKey(174, 163), 12},
		{makeHistoryKey(157, 172), 12},
		{makeHistoryKey(163, 172), 4},
		{makeHistoryKey(163, 185), 8},
		{makeHistoryKey(163, 192), 8},
	},
}

// swacEast and swacWest divisions
// SWAC East: Bethune-Cookman(184), Alabama A&M(188), Florida A&M(190),
//
//	Alabama State(191), Mississippi Valley State(193), Jackson State(194)
//
// SWAC West: Grambling State(183), Prairie View A&M(185), Alcorn State(186),
//
//	Southern(187), Arkansas-Pine Bluff(189), Texas Southern(192)
var swacEast = []uint{184, 188, 190, 191, 193, 194} // BCookman, AlaA&M, FAMU, AlaSt, MVST, JSU
var swacWest = []uint{183, 185, 186, 187, 189, 192} // Grambling, PVA&M, Alcorn, Southern, UAPB, TxSo

// GenerateSWACSchedule handles the SWAC 12-team divisional schedule (8 conf games).
func GenerateSWACSchedule(
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

	// SWAC locked intra-conference matchups
	swacLocked := []struct {
		key  SchedulerHistoryKey
		week uint
	}{
		{makeHistoryKey(190, 194), 1},  // FAMU(190) vs Jackson State(194) — Week 1
		{makeHistoryKey(185, 192), 1},  // Prairie View(185) vs Texas Southern(192) — Week 1
		{makeHistoryKey(183, 185), 4},  // Grambling(183) vs Prairie View(185) — Week 4
		{makeHistoryKey(188, 191), 10}, // Alabama A&M(188) vs Alabama State(191) — Week 10
		{makeHistoryKey(184, 190), 13}, // Bethune-Cookman(184) vs FAMU(190) — Week 13
		{makeHistoryKey(183, 187), 13}, // Grambling(183) vs Southern(187) — Week 13
		{makeHistoryKey(193, 194), 13}, // MVST(193) vs Jackson State(194) — Week 13
	}
	lockedSet := make(map[SchedulerHistoryKey]uint)
	for _, l := range swacLocked {
		lockedSet[l.key] = l.week
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
		if alreadyScheduled(home.ID, away.ID, gamesPlayedAgainstOpponentsMap) {
			continue
		}
		// Skip if either team already has a game this week (OOC conflict)
		if gamesPlayedByWeekMap[home.ID] != nil && gamesPlayedByWeekMap[home.ID][l.week] {
			continue
		}
		if gamesPlayedByWeekMap[away.ID] != nil && gamesPlayedByWeekMap[away.ID][l.week] {
			continue
		}
		markWeek(home.ID, away.ID, l.week, gamesPlayedByWeekMap)
		markOpponents(home.ID, away.ID, gamesPlayedAgainstOpponentsMap)
		homecountMap[home.ID]++
		g := CreateCollegeGameRecord(home, away, l.week, seasonID, stadiumMap, stadiumMapByID, rivalryMap)
		games = append(games, g)
	}

	emitGame := func(home, away structs.CollegeTeam, week uint) {
		if home.ID == 0 || away.ID == 0 {
			return
		}
		if alreadyScheduled(home.ID, away.ID, gamesPlayedAgainstOpponentsMap) {
			return
		}
		markWeek(home.ID, away.ID, week, gamesPlayedByWeekMap)
		markOpponents(home.ID, away.ID, gamesPlayedAgainstOpponentsMap)
		homecountMap[home.ID]++
		g := CreateCollegeGameRecord(home, away, week, seasonID, stadiumMap, stadiumMapByID, rivalryMap)
		games = append(games, g)
	}

	// Schedule 5 intra-divisional games per team
	confGameCount := make(map[uint]int)
	scheduleDivGames := func(div []uint, maxGames int) {
		for i := 0; i < len(div); i++ {
			for j := i + 1; j < len(div); j++ {
				a := div[i]
				b := div[j]
				key := makeHistoryKey(a, b)
				if _, locked := lockedSet[key]; locked {
					continue
				}
				if alreadyScheduled(a, b, gamesPlayedAgainstOpponentsMap) {
					continue
				}
				if confGameCount[a] >= maxGames || confGameCount[b] >= maxGames {
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
				w := assignWeek(home.ID, away.ID, 4, 12, gamesPlayedByWeekMap)
				if w == 0 {
					continue
				}
				emitGame(home, away, w)
				confGameCount[home.ID]++
				confGameCount[away.ID]++
			}
		}
	}
	scheduleDivGames(swacEast, 5)
	scheduleDivGames(swacWest, 5)

	// 3 cross-divisional games per team, rotating
	cy := cycleYear(season)
	crossProcessed := make(map[SchedulerHistoryKey]bool)
	for eIdx, eTeam := range swacEast {
		for k := 0; k < 3; k++ {
			wIdx := (eIdx + (cy-1)*3 + k) % len(swacWest)
			wTeam := swacWest[wIdx]
			key := makeHistoryKey(eTeam, wTeam)
			if crossProcessed[key] || alreadyScheduled(eTeam, wTeam, gamesPlayedAgainstOpponentsMap) {
				crossProcessed[key] = true
				continue
			}
			if confGameCount[eTeam] >= 8 || confGameCount[wTeam] >= 8 {
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
			w := assignWeek(home.ID, away.ID, 4, 12, gamesPlayedByWeekMap)
			if w == 0 {
				crossProcessed[key] = true
				continue
			}
			crossProcessed[key] = true
			emitGame(home, away, w)
			confGameCount[home.ID]++
			confGameCount[away.ID]++
		}
	}

	return games
}

// GenerateSmallConferenceSchedule dispatches to the correct generator for FCS
// conferences. Falls back to the generic round-robin helper for most conferences.
func GenerateSmallConferenceSchedule(
	conferenceID int,
	retrySeed uint,
	collegeTeams []structs.CollegeTeam,
	stadiumMap map[uint]structs.Stadium,
	stadiumMapByID map[uint]structs.Stadium,
	rivalryMap map[uint][]structs.CollegeRival,
	gamesPlayedAgainstOpponentsMap map[uint]map[uint]bool,
	gamesPlayedByWeekMap map[uint]map[uint]bool,
	ts structs.Timestamp,
) []structs.CollegeGame {
	if conferenceID == 19 {
		return GenerateSWACSchedule(collegeTeams, stadiumMap, stadiumMapByID, rivalryMap, gamesPlayedAgainstOpponentsMap, gamesPlayedByWeekMap, ts)
	}

	gamesPerTeam, ok := smallConfGameCount[conferenceID]
	if !ok {
		gamesPerTeam = 8
	}

	// Build locked set from table
	lockedSet := make(map[SchedulerHistoryKey]uint)
	if locks, found := smallConfLockedMatchups[conferenceID]; found {
		for _, l := range locks {
			lockedSet[l.key] = l.week
		}
	}

	return generateGenericRoundRobinSchedule(
		collegeTeams,
		stadiumMap,
		stadiumMapByID,
		rivalryMap,
		gamesPlayedAgainstOpponentsMap,
		gamesPlayedByWeekMap,
		nil, // playCountMap
		nil, // lastHomeMap
		ts,
		lockedSet,
		gamesPerTeam,
		4,
		14,
		1, 14, // overflow: weeks 1–3 and week 14
		retrySeed,
	)
}
