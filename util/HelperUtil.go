package util

func GetNFLFullTeamName(teamName, mascot string) string {
	return teamName + " " + mascot
}

func GetWeekID(seasonID uint, week uint) uint {
	// WeekID structure is the final two digits of the season year followed by the two digits representing the week.
	// Season 1 == 2021
	// WeekID should look like 2101 for season 1 week 1, 2102 for season 1 week 2, etc.
	return (seasonID+2020-2000)*100 + week
}
