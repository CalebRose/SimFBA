package managers

import (
	"github.com/CalebRose/SimFBA/structs"
)

// ============================================================
// C-USA Conference Schedule Generator
// 10 teams | 8 conf games | generic partial round-robin
// Locked: Middle Tennessee/Western Kentucky — Week 13
// ============================================================
//
// C-USA team IDs:
// FIU(29), Liberty(48), Middle Tennessee(61), New Mexico State(70),
// Western Kentucky(127), Jacksonville State(132), Sam Houston State(133),
// Kennesaw State(134), Missouri State(138), Delaware(206)

var cusaLockedMatchups = []struct {
	key  SchedulerHistoryKey
	week uint
}{
	{makeHistoryKey(61, 127), 13}, // Middle Tennessee vs Western Kentucky
}

// GenerateCUSASchedule produces all C-USA conference games for the season.
func GenerateCUSASchedule(
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
	for _, l := range cusaLockedMatchups {
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
		1, 13, // overflow
		0, // retrySeed
	)
}
