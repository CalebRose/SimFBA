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
	if ts.IsCollegeRegularSeason() && !ts.CFBSpringGames && ts.CollegeWeek > 0 && ts.CollegeWeek < 21 {
		managers.SyncRecruiting(ts)
	}
	if ts.RunTransferIntentions() {
		managers.ProcessTransferIntention()
	}
	if ts.ReadyToActivatePortal() {
		managers.EnterTheTransferPortal()
	} else if ts.RunCron && (ts.CollegeSeasonOver || ts.IsOffSeason) && ts.TransferPortalPhase == 3 && ts.TransferPortalRound <= 10 {
		managers.SyncTransferPortal()
	}
}

func SyncFreeAgencyViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron {
		managers.SyncFreeAgencyOffers()
		// if ts.FreeAgencyRound >= 1 && ts.FreeAgencyRound < 25 {
		// 	managers.MoveUpInOffseasonFreeAgency()
		// }
		managers.AllocateCapsheets()
	}
}

func SyncToNextWeekViaCron() {
	db := dbprovider.GetInstance().GetDB()
	ts := managers.GetTimestamp()

	ts.MoveUpPhase()
	repository.SaveTimestamp(ts, db)
	if ts.Phase < 7 {
		return
	}

	if ts.RunCron {
		if !ts.IsOffSeason && !ts.IsNFLOffSeason {
			ts = managers.MoveUpWeek()
		}

		if ts.RunRedshirtsAlgorithm() {
			seasonID := strconv.Itoa(ts.CollegeSeasonID)
			managers.AllocateAIRedshirts(seasonID)
		}

		managers.AssignTeamGrades()

		if ts.RunGenerateCollegeRankings() {
			seasonID := strconv.Itoa(ts.CollegeSeasonID)
			managers.GenerateCollegeRankings(seasonID)
		}

		// Once National Championship is over and we move up a week.
		if ts.RunSyncPromises() {
			// Sync Promises
			managers.SyncPromises()
			ts.TransferPortalPhase = 1
			ts.TransferPortalRound = 1
		}

		if ts.ReadyForSeasonMigration() {
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
	if ts.RunGenerateWeatherFunction() {
		// Generate weather for existing CFB Games
		managers.GenerateWeatherForGames()
	}

	if ts.RunGenerateOOCSchedule() {
		managers.GenerateOOCSchedule()
		managers.GenerateWeatherForGames()
	}

	// Calculate Player Minimum Values and Average Annual Value (AAV) for players
	if ts.RunExtensionAndTagValues() {
		managers.CalculateTagValues()
		managers.CalculatePlayerMinimumAndAAVValues()
	}

	if ts.RunCollegeProgressions() {
		// Reset progression flags on all teams and players before running
		db.Model(&structs.CollegeTeam{}).Where("id > ?", 0).Updates(map[string]interface{}{"players_progressed": false, "recruits_added": false})
		db.Model(&structs.CollegePlayer{}).Where("id > ?", 0).Update("has_progressed", false)

		managers.CFBProgressionMain()
		ts.ToggleCollegeProgression()
		managers.RecruitingAndTransferPortalCleanUp()
		repository.SaveTimestamp(ts, db)
	}

	if ts.RunProfessionalProgressions() {
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
	if ts.RunUDFAs() {
		managers.ProcessUDFAs(false)
	}

	if ts.RunTrainingCamps() {
		managers.RunTrainingCamps("")
	}

	// Generate Walkons in the middle of the night.
	if ts.RunGenerateWalkonAlgorithm() {
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
	if ts.IsCollegeRegularSeason() {
		managers.DetermineAIGameplan()
	}
}

func RunAIGameplanViaCron() {
	ts := managers.GetTimestamp()
	if ts.IsCollegeRegularSeason() {
		managers.SetAIGameplan()
	}
}

func RunTheGamesViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron {
		if ts.RunGames && (ts.IsCollegeRegularSeason() || ts.IsNFLRegularSeason()) {
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
	if ts.RunCron && (ts.IsCollegeRegularSeason()) {
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
	if ts.RunCron && (ts.IsNFLRegularSeason()) {
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
	if ts.RunCron && (ts.IsCollegeRegularSeason()) {
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
	if ts.RunCron && (ts.IsCollegeRegularSeason()) {
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
	if ts.RunCron && (ts.IsCollegeRegularSeason()) {
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
	if ts.RunCron && (ts.IsCollegeRegularSeason()) {
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
	if ts.RunCron && (ts.IsCollegeRegularSeason()) {
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
	if ts.RunCron && (ts.IsNFLRegularSeason()) {
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
	if ts.RunCron && (ts.IsNFLRegularSeason()) {
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
	if ts.RunCron && (ts.IsNFLRegularSeason()) {
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
	if ts.RunCron && (ts.IsNFLRegularSeason()) {
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
