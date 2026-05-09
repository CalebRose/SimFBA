package managers

import (
	"github.com/CalebRose/SimFBA/structs"
)

// ============================================================
// MWC Conference Schedule Generator
// 10 teams | 8 conf games per season (not full round-robin)
// Locked: UNLV/Nevada Week 14, Wyoming/Hawaii Week 13,
//         Hawaii/San Jose State Week 4
// ============================================================
//
// MWC team IDs:
// Air Force(1), Hawaii(38), Nevada(68), New Mexico(69),
// Northern Illinois(73), San Jose State(90), UNLV(114),
// UTEP(116), Wyoming(130), North Dakota State(143)

var mwcTeamIDs = []uint{1, 38, 68, 69, 73, 90, 114, 116, 130, 143}

// mwcLockedMatchups: pairs with a specific week assignment
var mwcLockedMatchups = []struct {
	key  SchedulerHistoryKey
	week uint
}{
	{makeHistoryKey(114, 68), 14}, // UNLV vs Nevada
	{makeHistoryKey(130, 38), 13}, // Wyoming vs Hawaii
	{makeHistoryKey(38, 90), 4},   // Hawaii vs San Jose State
}

func GenerateMWCSchedule(
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
	lockedSet := make(map[SchedulerHistoryKey]uint)
	for _, l := range mwcLockedMatchups {
		lockedSet[l.key] = l.week
	}

	return generateGenericRoundRobinSchedule(
		collegeTeams,
		stadiumMap,
		stadiumMapByID,
		rivalryMap,
		gamesPlayedAgainstOpponentsMap,
		gamesPlayedByWeekMap,
		playCountMap,
		lastHomeMap,
		homeCountSeedMap,
		ts,
		lockedSet,
		8,     // games per team
		4,     // min week
		13,    // max week
		1, 14, // overflow
		0, // retrySeed
	)
}
