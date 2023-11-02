package controller

import "github.com/CalebRose/SimFBA/managers"

func FillAIBoardsViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && !ts.IsOffSeason {
		managers.FillAIRecruitingBoards()
	}
}

func SyncAIBoardsViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && !ts.IsOffSeason {
		managers.ResetAIBoardsForCompletedTeams()
		managers.AllocatePointsToAIBoards()
	}
}

func SyncRecruitingViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron {
		managers.SyncRecruiting(ts)
	}

	managers.MoveUpInOffseasonFreeAgency()
	managers.SyncFreeAgencyOffers()
}

func SyncFreeAgencyOffersViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron {
		managers.MoveUpInOffseasonFreeAgency()
		managers.SyncFreeAgencyOffers()
	}
}

func SyncToNextWeekViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron {
		if !ts.IsOffSeason {
			// managers.ResetCollegeStandingsRanks()
		}
		managers.MoveUpWeek()
		if !ts.IsOffSeason {
			// managers.SyncCollegePollSubmissionForCurrentWeek()
		}
	}
}

func ShowGamesViaCron() {
	ts := managers.GetTimestamp()
	timeslot := ""
	if !ts.ThursdayGames {
		timeslot = "Thursday Night"
	} else if !ts.NFLThursday {
		timeslot = "Thursday Night Football"
	} else if !ts.SaturdayMorning {
		timeslot = "Saturday Noon"
	} else if !ts.SaturdayNoon {
		timeslot = "Saturday Afternoon"
	} else if !ts.SaturdayEvening {
		timeslot = "Saturday Evening"
	} else if !ts.SaturdayNight {
		timeslot = "Saturday Night"
	} else if !ts.NFLSundayNoon {
		timeslot = "Sunday Noon"
	} else if !ts.NFLSundayAfternoon {
		timeslot = "Sunday Afternoon"
	} else if !ts.NFLSundayEvening {
		timeslot = "Sunday Night Football"
	} else if !ts.NFLMondayEvening {
		timeslot = "Monday Night Football"
	}
	if ts.RunCron {
		managers.SyncTimeslot(timeslot)
	}
}
