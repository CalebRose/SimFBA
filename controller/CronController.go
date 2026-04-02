package controller

import (
	"fmt"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/managers"
	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
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
		if ts.CollegeWeek == 20 {
			managers.GenerateWalkOns()
		}
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

func RunCFBProgressionsViaCron() {
	db := dbprovider.GetInstance().GetDB()
	ts := managers.GetTimestamp()
	if ts.CollegeWeek < 21 {
		return
	}
	if ts.CollegeSeasonOver && !ts.ProgressedCollegePlayers {
		// Reset progression flags on all teams and players before running
		db.Model(&structs.CollegeTeam{}).Where("id > ?", 0).Updates(map[string]interface{}{"players_progressed": false, "recruits_added": false})
		db.Model(&structs.CollegePlayer{}).Where("id > ?", 0).Update("has_progressed", false)

		managers.CFBProgressionMain()
		ts.ToggleCollegeProgression()
		repository.SaveTimestamp(ts, db)
	}
}

func RunNFLProgressionsViaCron() {
	db := dbprovider.GetInstance().GetDB()
	ts := managers.GetTimestamp()
	if ts.NFLWeek < 23 {
		return
	}

	if ts.NFLSeasonOver && !ts.ProgressedProfessionalPlayers {
		// Reset progression flags on all teams and players before running
		db.Model(&structs.NFLPlayer{}).Where("id > ?", 0).Update("has_progressed", false)
		managers.NFLProgressionMain()
		ts.ToggleProfessionalProgression()
		repository.SaveTimestamp(ts, db)
	}
}

func SyncToNextWeekViaCron() {
	ts := managers.GetTimestamp()
	if !ts.RunGames {
		return
	}
	if ts.RunCron {
		if !ts.IsOffSeason && !ts.IsNFLOffSeason {
			ts = managers.MoveUpWeek()
		}
		managers.AssignTeamGrades()

		// Once National Championship is over and we move up a week.
		if ts.CollegeSeasonOver && ts.CollegeWeek == 21 {
			// Sync Promises
			managers.SyncPromises()
			ts.TransferPortalPhase = 1
			ts.TransferPortalRound = 1
			db := dbprovider.GetInstance().GetDB()
			repository.SaveTimestamp(ts, db)
		}

		if ts.NFLSeasonOver && ts.CollegeSeasonOver && !ts.IsNFLOffSeason && !ts.IsOffSeason && ts.ProgressedCollegePlayers && ts.ProgressedProfessionalPlayers {
			db := dbprovider.GetInstance().GetDB()
			ts.MoveUpSeason()
			repository.SaveTimestamp(ts, db)
			managers.GenerateOffseasonData()
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

	if !ts.RunGames {
		return
	}

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
	if !ts.RunGames {
		return
	}
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
	if !ts.RunGames {
		return
	}
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
	if !ts.RunGames {
		return
	}
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
	if !ts.RunGames {
		return
	}
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
	if !ts.RunGames {
		return
	}
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
	if !ts.RunGames {
		return
	}
	timeslot := ""
	if !ts.SaturdayNight {
		timeslot = "Saturday Night"
	}
	if ts.RunCron && (!ts.IsOffSeason || !ts.IsNFLOffSeason || !ts.CollegeSeasonOver || !ts.NFLSeasonOver) {
		managers.SyncTimeslot(timeslot)
	}
}

func ShowNFLSunNoonViaCron() {
	ts := managers.GetTimestamp()
	if !ts.RunGames {
		return
	}
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
	if !ts.RunGames {
		return
	}
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
	if !ts.RunGames {
		return
	}
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
	if !ts.RunGames {
		return
	}
	timeslot := ""
	if !ts.NFLMondayEvening {
		timeslot = "Monday Night Football"
	}
	if ts.RunCron && (!ts.IsOffSeason || !ts.IsNFLOffSeason || !ts.CollegeSeasonOver || !ts.NFLSeasonOver) {
		managers.SyncTimeslot(timeslot)
	}
}
