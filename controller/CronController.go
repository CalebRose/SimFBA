package controller

import (
	"fmt"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/managers"
	"github.com/CalebRose/SimFBA/repository"
)

func CronTest() {
	fmt.Println("PING!")
}

func FillAIBoardsViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && !ts.IsOffSeason && !ts.CollegeSeasonOver {
		managers.FillAIRecruitingBoards()
	}

	if ts.RunCron && ts.IsOffSeason && ts.CollegeSeasonOver && ts.TransferPortalPhase == 3 {
		managers.AICoachFillBoardsPhase()
	}
}

func SyncAIBoardsViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && !ts.IsOffSeason && !ts.CollegeSeasonOver {
		managers.ResetAIBoardsForCompletedTeams()
		managers.AllocatePointsToAIBoards()
	}

	if ts.RunCron && ts.IsOffSeason && ts.CollegeSeasonOver && ts.TransferPortalPhase == 3 {
		managers.AICoachAllocateAndPromisePhase()
	}
}

func RunRESViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && !ts.IsOffSeason && !ts.CollegeSeasonOver {
		managers.SyncRecruitingEfficiency(ts)
	}
}

func SyncRecruitingViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && !ts.IsOffSeason && !ts.CollegeSeasonOver && !ts.CFBSpringGames && ts.CollegeWeek > 0 && ts.CollegeWeek < 21 {
		managers.SyncRecruiting(ts)
	}
	if ts.RunCron && ts.IsOffSeason && ts.TransferPortalPhase == 2 {
		managers.EnterTheTransferPortal()
	} else if ts.RunCron && ts.IsOffSeason && ts.TransferPortalPhase == 3 {
		managers.SyncTransferPortal()
	}
}

func SyncFreeAgencyViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron {
		managers.SyncFreeAgencyOffers()
		if ts.FreeAgencyRound >= 1 && ts.FreeAgencyRound < 25 {
			managers.MoveUpInOffseasonFreeAgency()
		}
		managers.AllocateCapsheets()
	}
}

func SyncToNextWeekViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron {
		if !ts.IsOffSeason && !ts.IsNFLOffSeason {
			ts = managers.MoveUpWeek()
		}
		managers.AssignTeamGrades()
		if ts.NFLSeasonOver && ts.CollegeSeasonOver && !ts.IsNFLOffSeason && !ts.IsOffSeason {
			if !ts.ProgressedCollegePlayers && !ts.ProgressedProfessionalPlayers {
				// Progress College

				ts.ToggleCollegeProgression()
				// Progress NFL

				ts.ToggleProfessionalProgression()

				db := dbprovider.GetInstance().GetDB()
				ts.MoveUpSeason()
				repository.SaveTimestamp(ts, db)
				managers.GenerateOffseasonData()
			}
		}
	}
}

func RunAISchemeAndDCViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && !ts.IsOffSeason && !ts.CollegeSeasonOver {
		managers.DetermineAIGameplan()
	}
}

func RunAIGameplanViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && !ts.IsOffSeason && !ts.CollegeSeasonOver {
		managers.SetAIGameplan()
	}
}

func RunTheGamesViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron {
		if !ts.IsOffSeason && ts.RunGames && ts.NFLWeek < 22 && ts.CollegeWeek < 21 {
			managers.FixBrokenGameplans()
			managers.CheckForSchemePenalties()
			managers.RunTheGames()
		}
	}
}

func ShowCFBThursdayViaCron() {
	ts := managers.GetTimestamp()
	timeslot := ""
	if !ts.ThursdayGames {
		timeslot = "Thursday Night"
	}
	if ts.RunCron && (!ts.IsOffSeason || !ts.IsNFLOffSeason || !ts.CollegeSeasonOver || !ts.NFLSeasonOver) {
		managers.SyncTimeslot(timeslot)
	}
}

func ShowNFLThursdayViaCron() {
	ts := managers.GetTimestamp()
	timeslot := ""
	if !ts.NFLThursday {
		timeslot = "Thursday Night Football"
	}
	if ts.RunCron && (!ts.IsOffSeason || !ts.IsNFLOffSeason || !ts.CollegeSeasonOver || !ts.NFLSeasonOver) {
		managers.SyncTimeslot(timeslot)
	}
}

func ShowCFBFridayViaCron() {
	ts := managers.GetTimestamp()
	timeslot := ""
	if !ts.FridayGames {
		timeslot = "Friday Night"
	}
	if ts.RunCron && (!ts.IsOffSeason || !ts.IsNFLOffSeason || !ts.CollegeSeasonOver || !ts.NFLSeasonOver) {
		managers.SyncTimeslot(timeslot)
	}
}

func ShowCFBSatMornViaCron() {
	ts := managers.GetTimestamp()
	timeslot := ""
	if !ts.SaturdayMorning {
		timeslot = "Saturday Morning"
	}
	if ts.RunCron && (!ts.IsOffSeason || !ts.IsNFLOffSeason || !ts.CollegeSeasonOver || !ts.NFLSeasonOver) {
		managers.SyncTimeslot(timeslot)
	}
}

func ShowCFBSatAftViaCron() {
	ts := managers.GetTimestamp()
	timeslot := ""
	if !ts.SaturdayNoon {
		timeslot = "Saturday Afternoon"
	}
	if ts.RunCron && (!ts.IsOffSeason || !ts.IsNFLOffSeason || !ts.CollegeSeasonOver || !ts.NFLSeasonOver) {
		managers.SyncTimeslot(timeslot)
	}
}

func ShowCFBSatEveViaCron() {
	ts := managers.GetTimestamp()
	timeslot := ""
	if !ts.SaturdayEvening {
		timeslot = "Saturday Evening"
	}
	if ts.RunCron && (!ts.IsOffSeason || !ts.IsNFLOffSeason || !ts.CollegeSeasonOver || !ts.NFLSeasonOver) {
		managers.SyncTimeslot(timeslot)
	}
}

func ShowCFBSatNitViaCron() {
	ts := managers.GetTimestamp()
	timeslot := ""
	if !ts.ThursdayGames {
		timeslot = "Thursday Night"
	} else if !ts.NFLThursday {
		timeslot = "Thursday Night Football"
	} else if !ts.SaturdayMorning {
		timeslot = "Saturday Morning"
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
	if ts.RunCron && (!ts.IsOffSeason || !ts.IsNFLOffSeason || !ts.CollegeSeasonOver || !ts.NFLSeasonOver) {
		managers.SyncTimeslot(timeslot)
	}
}

func ShowNFLSunNoonViaCron() {
	ts := managers.GetTimestamp()
	timeslot := ""
	if !ts.NFLSundayNoon {
		timeslot = "Sunday Noon"
	}
	if ts.RunCron && (!ts.IsOffSeason || !ts.IsNFLOffSeason || !ts.CollegeSeasonOver || !ts.NFLSeasonOver) {
		managers.SyncTimeslot(timeslot)
	}
}

func ShowNFLSunAftViaCron() {
	ts := managers.GetTimestamp()
	timeslot := ""
	if !ts.NFLSundayAfternoon {
		timeslot = "Sunday Afternoon"
	}
	if ts.RunCron && (!ts.IsOffSeason || !ts.IsNFLOffSeason || !ts.CollegeSeasonOver || !ts.NFLSeasonOver) {
		managers.SyncTimeslot(timeslot)
	}
}

func ShowNFLSunNitViaCron() {
	ts := managers.GetTimestamp()
	timeslot := ""
	if !ts.NFLSundayEvening {
		timeslot = "Sunday Night Football"
	}
	if ts.RunCron && (!ts.IsOffSeason || !ts.IsNFLOffSeason || !ts.CollegeSeasonOver || !ts.NFLSeasonOver) {
		managers.SyncTimeslot(timeslot)
	}
}

func ShowNFLMonNitViaCron() {
	ts := managers.GetTimestamp()
	timeslot := ""
	if !ts.NFLMondayEvening {
		timeslot = "Monday Night Football"
	}
	if ts.RunCron && (!ts.IsOffSeason || !ts.IsNFLOffSeason || !ts.CollegeSeasonOver || !ts.NFLSeasonOver) {
		managers.SyncTimeslot(timeslot)
	}
}
