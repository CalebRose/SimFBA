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

func GetTimeslot(state string, conferenceID uint) string {
	stateKey := GetStateKey(state)
	// Implement logic to determine timeslot based on stateKey
	switch stateKey {
	case "TX", "OK", "AR", "LA", "KS", "IA", "MO", "NE", "SD", "ND", "MN", "WI", "IL", "NM":
		if (conferenceID > 7 && conferenceID < 13) || (conferenceID > 13) {
			return PickFromStringList([]string{"Saturday Afternoon", "Saturday Evening", "Friday Night", "Friday Night", "Saturday Night"})
		}
		return PickFromStringList([]string{"Saturday Afternoon", "Saturday Evening", "Saturday Afternoon", "Friday Night", "Saturday Night"})

	case "MS", "AL", "GA", "FL", "SC", "NC", "VA", "TN", "KY", "IN", "OH", "MI", "PA", "WV", "MD":
		if (conferenceID > 7 && conferenceID < 13) || (conferenceID > 13) {
			return PickFromStringList([]string{"Saturday Afternoon", "Saturday Morning", "Thursday Night", "Thursday Night"})
		}
		return PickFromStringList([]string{"Saturday Afternoon", "Saturday Evening", "Saturday Morning", "Thursday Night"})

	case "NY", "NJ", "CT", "RI", "MA", "VT", "NH", "ME":
		if (conferenceID > 7 && conferenceID < 13) || (conferenceID > 13) {
			return PickFromStringList([]string{"Saturday Afternoon", "Saturday Morning", "Thursday Night", "Thursday Night"})
		}
		return PickFromStringList([]string{"Saturday Afternoon", "Saturday Morning", "Saturday Morning", "Thursday Night"})

	case "WA", "OR", "CA", "NV", "ID", "UT", "AZ", "CO", "MT", "WY", "GM", "AS", "AK", "HI":
		if (conferenceID > 7 && conferenceID < 13) || (conferenceID > 13) {
			return PickFromStringList([]string{"Saturday Evening", "Saturday Evening", "Friday Night", "Friday Night", "Thursday Night", "Saturday Night"})
		}
		return PickFromStringList([]string{"Saturday Afternoon", "Saturday Evening", "Saturday Evening", "Saturday Night", "Saturday Night", "Saturday Night", "Friday Night"})
	}
	return "Thursday Night"
}
