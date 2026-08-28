package controller

import (
	"fmt"
	"strconv"

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

	if ts.RunCron && (ts.IsOffSeason || ts.CollegeSeasonOver) && ts.TransferPortalPhase == 3 && ts.TransferPortalRound <= 10 {
		managers.AICoachFillBoardsPhase()
	}
}

func SyncAIBoardsViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && !ts.IsOffSeason && !ts.CollegeSeasonOver {
		managers.ResetAIBoardsForCompletedTeams()
		managers.AllocatePointsToAIBoards()
	}

	if ts.RunCron && (ts.IsOffSeason || ts.CollegeSeasonOver) && ts.TransferPortalPhase == 3 && ts.TransferPortalRound <= 10 {
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
	if (ts.RunCron && ts.CollegeSeasonOver && ts.TransferPortalPhase == 1) || ts.Phase == 31 {
		managers.ProcessTransferIntention()
	}
	if (ts.RunCron && ts.CollegeSeasonOver && ts.TransferPortalPhase == 2) || ts.Phase == 32 {
		managers.EnterTheTransferPortal()
	} else if ts.RunCron && (ts.CollegeSeasonOver || ts.IsOffSeason) && ts.TransferPortalPhase == 3 && ts.TransferPortalRound <= 10 {
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
	db := dbprovider.GetInstance().GetDB()
	ts := managers.GetTimestamp()

	ts.MoveUpPhase()

	if ts.Phase < 7 {
		return
	}

	if ts.RunCron {
		if !ts.IsOffSeason && !ts.IsNFLOffSeason {
			ts = managers.MoveUpWeek()
		}

		managers.AssignTeamGrades()

		if ts.CollegeWeek >= 2 && !ts.CFBSpringGames && ts.CollegeWeek < 21 {
			seasonID := strconv.Itoa(ts.CollegeSeasonID)
			managers.GenerateCollegeRankings(seasonID)
		}

		// Once National Championship is over and we move up a week.
		if (ts.CollegeSeasonOver && ts.CollegeWeek == 21) || ts.Phase == 31 {
			// Sync Promises
			managers.SyncPromises()
			ts.TransferPortalPhase = 1
			ts.TransferPortalRound = 1
		}

		if (ts.NFLSeasonOver && ts.CollegeSeasonOver && !ts.IsNFLOffSeason && !ts.IsOffSeason && ts.ProgressedCollegePlayers && ts.ProgressedProfessionalPlayers) || ts.Phase > 34 {
			db := dbprovider.GetInstance().GetDB()
			ts.MoveUpSeason()
			repository.SaveTimestamp(ts, db)
			managers.GenerateOffseasonData()
		}
		repository.SaveTimestamp(ts, db)
	}
}

// Sync Phase Tuesday via Cron - Handling separate actions based on the current phase of the season.
func SyncPhaseTuesdayViaCron() {
	ts := managers.GetTimestamp()
	if !ts.RunCron {
		return
	}
	db := dbprovider.GetInstance().GetDB()
	if ts.Phase == 1 {
		// Generate weather for existing CFB Games
		managers.GenerateWeatherForGames()
	}

	if ts.Phase == 6 {
		managers.GenerateOOCSchedule()
		managers.GenerateWeatherForGames()
	}

	// Calculate Player Minimum Values and Average Annual Value (AAV) for players
	if ts.Phase == 19 || ts.NFLWeek == 9 {
		// w := http.ResponseWriter(nil)
		// managers.CalculatePlayerMinimumAndAAVValues(w)
	}

	if (ts.CollegeSeasonOver && !ts.ProgressedCollegePlayers) || (ts.Phase == 31 && !ts.ProgressedCollegePlayers) {
		// Reset progression flags on all teams and players before running
		db.Model(&structs.CollegeTeam{}).Where("id > ?", 0).Updates(map[string]interface{}{"players_progressed": false, "recruits_added": false})
		db.Model(&structs.CollegePlayer{}).Where("id > ?", 0).Update("has_progressed", false)

		managers.CFBProgressionMain()
		ts.ToggleCollegeProgression()
		managers.RecruitingAndTransferPortalCleanUp()
		repository.SaveTimestamp(ts, db)
	}

	if (ts.NFLWeek == 23 || ts.Phase == 33) && ts.NFLSeasonOver && !ts.ProgressedProfessionalPlayers {
		// Reset progression flags on all teams and players before running
		db.Model(&structs.NFLPlayer{}).Where("id > ?", 0).Update("has_progressed", false)
		managers.NFLProgressionMain()
		ts.ToggleProfessionalProgression()
		managers.FreeAgencyCleanUp()
		repository.SaveTimestamp(ts, db)
	}
}

// Sync Phase Wednesday via Cron - Handling separate actions based on the current phase of the season.
func SyncPhaseWednesdayViaCron() {
	ts := managers.GetTimestamp()
	if !ts.RunCron {
		return
	}

	// Run UDFAs
	if ts.Phase == 4 {
		managers.ProcessUDFAs(false)
	}

	if ts.Phase == 5 {
		managers.RunTrainingCamps("")
	}

	// Generate Walkons in the middle of the night.
	if ts.CollegeWeek == 20 || ts.Phase == 30 {
		managers.GenerateWalkOns()
		managers.AssignAllRecruitRanks()
	}
}

// Sync Phase Thursday via Cron - Handling separate actions based on the current phase of the season.
func SyncPhaseThursdayViaCron() {
	ts := managers.GetTimestamp()
	if !ts.RunCron {
		return
	}
}

// Sync Phase Friday via Cron - Handling separate actions based on the current phase of the season.
func SyncPhaseFridayViaCron() {
	ts := managers.GetTimestamp()
	if !ts.RunCron {
		return
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

func StreamCFBThursdayGamesToInterfaceViaCron() {
	managers.StartCFBLiveStreamingCron()
}

func StreamCFBFridayGamesToInterfaceViaCron() {
	managers.StartCFBLiveStreamingCron()
}

func StreamCFBSaturdayGamesToInterfaceViaCron() {
	managers.StartCFBLiveStreamingCron()
}

func StreamNFLThursdayGamesToInterfaceViaCron() {
	managers.StartNFLLiveStreamingCron()
}

func StreamNFLSundayGamesToInterfaceViaCron() {
	managers.StartNFLLiveStreamingCron()
}

func StreamNFLMondayGamesToInterfaceViaCron() {
	managers.StartNFLLiveStreamingCron()
}
