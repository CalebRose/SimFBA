package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/CalebRose/SimFBA/controller"
	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"github.com/nelkinda/health-go"
	"github.com/nelkinda/health-go/checks/sendgrid"
)

func InitialMigration() {
	initiate := dbprovider.GetInstance().InitDatabase()
	if !initiate {
		log.Println("Initiate pool failure... Ending application")
		os.Exit(1)
	}
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	// Health Controls
	HealthCheck := health.New(
		health.Health{
			Version:   "1",
			ReleaseID: "0.0.7-SNAPSHOT",
		},
		sendgrid.Health(),
	)
	myRouter.HandleFunc("/health", HealthCheck.Handler).Methods("GET")

	// Admin Controls
	myRouter.HandleFunc("/simfba/get/timestamp/", controller.GetCurrentTimestamp).Methods("GET")
	myRouter.HandleFunc("/simfba/sync/timestamp/", controller.SyncTimestamp).Methods("POST")
	myRouter.HandleFunc("/simfba/sync/week/", controller.SyncWeek).Methods("GET")
	myRouter.HandleFunc("/simfba/sync/timeslot/{timeslot}", controller.SyncTimeslot).Methods("GET")
	myRouter.HandleFunc("/simfba/regress/timeslot/{timeslot}", controller.RegressTimeslot).Methods("GET")
	myRouter.HandleFunc("/simfba/sync/freeagency/round", controller.SyncFreeAgencyRound).Methods("GET")
	myRouter.HandleFunc("/simfba/sync/recruiting/", controller.SyncRecruiting).Methods("GET")
	// myRouter.HandleFunc("/simfba/sync/missing/", controller.SyncMissingRES).Methods("GET")
	// myRouter.HandleFunc("/simfba/sync/weather/", controller.WeatherGenerator).Methods("GET")
	myRouter.HandleFunc("/simfba/current/weather/forecast/", controller.GetWeatherForecast).Methods("GET")
	myRouter.HandleFunc("/simfba/future/weather/forecast/", controller.GetFutureWeatherForecast).Methods("GET")
	myRouter.HandleFunc("/news/{weekID}/{seasonID}/", controller.GetNewsLogs).Methods("GET")
	myRouter.HandleFunc("/season/{seasonID}/weeks/{weekID}", controller.GetWeeksInSeason).Methods("GET")
	myRouter.HandleFunc("/admin/teams/croot/sync", controller.SyncTeamRecruitingRanks).Methods("GET")
	myRouter.HandleFunc("/admin/recruiting/class/size", controller.GetRecruitingClassSizeForTeams).Methods("GET")
	myRouter.HandleFunc("/admin/ai/fill/boards", controller.FillAIBoards).Methods("GET")
	myRouter.HandleFunc("/admin/ai/sync/boards", controller.SyncAIBoards).Methods("GET")
	// myRouter.HandleFunc("/admin/fix/affinities", controller.FixSmallTownBigCityAIBoards).Methods("GET")
	myRouter.HandleFunc("/admin/run/the/games/", controller.RunTheGames).Methods("GET")
	// myRouter.HandleFunc("/admin/overall/progressions/next/season", controller.ProgressToNextSeason).Methods("GET")
	myRouter.HandleFunc("/admin/trades/accept/sync/{proposalID}", controller.SyncAcceptedTrade).Methods("GET")
	myRouter.HandleFunc("/admin/trades/veto/sync/{proposalID}", controller.VetoAcceptedTrade).Methods("GET")
	myRouter.HandleFunc("/admin/trades/cleanup", controller.CleanUpRejectedTrades).Methods("GET")

	// Capsheet Controls
	myRouter.HandleFunc("/nfl/capsheet/generate", controller.GenerateCapsheets).Methods("GET")
	myRouter.HandleFunc("/nfl/contracts/get/value", controller.CalculateContracts).Methods("GET")

	// Draft Controls
	myRouter.HandleFunc("/nfl/draft/draftees/export/{season}", controller.ExportDrafteesToCSV).Methods("GET")

	// Free Agency Controls
	myRouter.HandleFunc("/nfl/freeagency/create/offer", controller.CreateFreeAgencyOffer).Methods("POST")
	myRouter.HandleFunc("/nfl/freeagency/cancel/offer", controller.CancelFreeAgencyOffer).Methods("POST")
	myRouter.HandleFunc("/nfl/waiverwire/create/offer", controller.CreateWaiverWireOffer).Methods("POST")
	myRouter.HandleFunc("/nfl/waiverwire/cancel/offer", controller.CancelWaiverWireOffer).Methods("POST")
	myRouter.HandleFunc("/nfl/freeagency/waiver/order/set", controller.SetWaiverOrderForNFLTeams).Methods("GET")

	// Game Controls
	myRouter.HandleFunc("/games/update/time/", controller.UpdateTimeslot).Methods("POST", "OPTIONS")
	myRouter.HandleFunc("/games/college/week/{weekID}/", controller.GetCollegeGamesByTimeslotWeekId).Methods("GET")
	myRouter.HandleFunc("/games/college/timeslot/{timeSlot}/{weekID}", controller.GetCollegeGamesByTimeslotWeekId).Methods("GET")
	myRouter.HandleFunc("/games/college/team/{teamID}/{seasonID}", controller.GetCollegeGamesByTeamIDAndSeasonID).Methods("GET")
	myRouter.HandleFunc("/games/college/season/{seasonID}", controller.GetCollegeGamesBySeasonID).Methods("GET")
	myRouter.HandleFunc("/games/nfl/team/{teamID}/{seasonID}", controller.GetNFLGamesByTeamIDAndSeasonID).Methods("GET")
	myRouter.HandleFunc("/games/nfl/season/{seasonID}", controller.GetNFLGamesBySeasonID).Methods("GET")
	myRouter.HandleFunc("/games/result/cfb/{gameID}", controller.GetCollegeGameResultsByGameID).Methods("GET")
	myRouter.HandleFunc("/games/result/nfl/{gameID}", controller.GetNFLGameResultsByGameID).Methods("GET")

	// Gameplan Controls
	myRouter.HandleFunc("/gameplan/college/team/{teamID}/", controller.GetTeamGameplanByTeamID).Methods("GET")
	myRouter.HandleFunc("/gameplan/college/updategameplan", controller.UpdateGameplan).Methods("POST")
	myRouter.HandleFunc("/gameplan/college/depthchart/{teamID}/", controller.GetTeamDepthchartByTeamID).Methods("GET")
	myRouter.HandleFunc("/gameplan/college/depthchart/user/check/", controller.CheckAllUserDepthChartsForInjuredPlayers).Methods("GET")
	myRouter.HandleFunc("/gameplan/college/depthchart/ai/update/", controller.UpdateCollegeAIDepthCharts).Methods("GET")
	myRouter.HandleFunc("/gameplan/college/depthchart/positions/{depthChartID}/", controller.GetDepthChartPositionsByDepthChartID).Methods("GET")
	myRouter.HandleFunc("/gameplan/college/updatedepthchart", controller.UpdateDepthChart).Methods("PUT")
	myRouter.HandleFunc("/gameplan/nfl/team/{teamID}/", controller.GetNFLGameplanByTeamID).Methods("GET")
	myRouter.HandleFunc("/gameplan/nfl/updategameplan", controller.UpdateNFLGameplan).Methods("POST")
	myRouter.HandleFunc("/gameplan/nfl/depthchart/{teamID}/", controller.GetNFLDepthChart).Methods("GET")
	myRouter.HandleFunc("/gameplan/nfl/updatedepthchart", controller.UpdateNFLDepthChart).Methods("POST")
	myRouter.HandleFunc("/gameplan/nfl/depthchart/ai/update/", controller.UpdateNFLAIDepthCharts).Methods("GET")

	// Generation Controls
	// myRouter.HandleFunc("/admin/generate/walkons", controller.GenerateWalkOns).Methods("GET")

	// Import Controls
	// myRouter.HandleFunc("/admin/import/recruit/ai", controller.ImportRecruitAICSV).Methods("GET")
	// myRouter.HandleFunc("/admin/import/nfl/draft", controller.Import2023DraftedPlayers).Methods("GET")
	// myRouter.HandleFunc("/admin/import/cfb/games", controller.ImportCFBGames).Methods("GET")
	// myRouter.HandleFunc("/admin/import/nfl/games", controller.ImportNFLGames).Methods("GET")
	// myRouter.HandleFunc("/admin/import/nfl/udfas", controller.ImportUDFAs).Methods("GET")
	// myRouter.HandleFunc("/admin/import/missing/recruits", controller.GetMissingRecruitingClasses).Methods("GET")
	// myRouter.HandleFunc("/admin/import/preferences", controller.ImportTradePreferences).Methods("GET")
	// myRouter.HandleFunc("/import/custom/croots", controller.ImportCustomCroots).Methods("GET")
	// myRouter.HandleFunc("/import/simnfl/updated/values", controller.ImportSimNFLMinimumValues).Methods("GET")

	// News Controls
	myRouter.HandleFunc("/cfb/news/all/", controller.GetAllNewsLogsForASeason).Methods("GET")
	myRouter.HandleFunc("/nfl/news/all/", controller.GetAllNFLNewsBySeason).Methods("GET")
	myRouter.HandleFunc("/news/feed/{league}/{teamID}/", controller.GetNewsFeed).Methods("GET")

	// Player Controls
	myRouter.HandleFunc("/players/all/", controller.AllPlayers).Methods("GET")
	myRouter.HandleFunc("/collegeplayers/heisman/", controller.GetHeismanList).Methods("GET")
	myRouter.HandleFunc("/collegeplayers/team/{teamID}/", controller.AllCollegePlayersByTeamID).Methods("GET")
	myRouter.HandleFunc("/collegeplayers/team/nors/{teamID}/", controller.AllCollegePlayersByTeamIDWithoutRedshirts).Methods("GET")
	myRouter.HandleFunc("/collegeplayers/team/export/{teamID}/", controller.ExportRosterToCSV).Methods("GET")
	myRouter.HandleFunc("/collegeplayers/assign/redshirt/", controller.ToggleRedshirtStatusForPlayer).Methods("POST")
	myRouter.HandleFunc("/nflplayers/team/{teamID}/", controller.AllNFLPlayersByTeamIDForDC).Methods("GET")
	myRouter.HandleFunc("/nflplayers/freeagency/available/{teamID}", controller.FreeAgencyAvailablePlayers).Methods("GET")
	myRouter.HandleFunc("/nflplayers/team/export/{teamID}/", controller.ExportNFLRosterToCSV).Methods("GET")
	myRouter.HandleFunc("/nflplayers/cut/player/{PlayerID}/", controller.CutNFLPlayerFromRoster).Methods("GET")
	myRouter.HandleFunc("/nflplayers/place/player/squad/{PlayerID}/", controller.PlaceNFLPlayerOnPracticeSquad).Methods("GET")
	myRouter.HandleFunc("/nflplayers/injury/reserve/player/{PlayerID}/", controller.PlaceNFLPlayerOnInjuryReserve).Methods("GET")
	// myRouter.HandleFunc("/collegeplayers/teams/export/", controller.ExportAllRostersToCSV).Methods("GET") // DO NOT USE

	// Rankings Controls
	myRouter.HandleFunc("/rankings/assign/all/croots/", controller.AssignAllRecruitRanks).Methods("GET")

	// Recruiting Controls
	myRouter.HandleFunc("/recruiting/overview/dashboard/{teamID}", controller.GetRecruitingProfileForDashboardByTeamID).Methods("GET")
	myRouter.HandleFunc("/recruiting/profile/team/{teamID}/", controller.GetRecruitingProfileForTeamBoardByTeamID).Methods("GET")
	myRouter.HandleFunc("/recruiting/profile/all/", controller.GetAllRecruitingProfiles).Methods("GET")
	myRouter.HandleFunc("/recruiting/profile/only/{teamID}/", controller.GetOnlyRecruitingProfileByTeamID).Methods("GET")
	myRouter.HandleFunc("/recruiting/toggle/ai/{teamID}/", controller.ToggleAIBehavior).Methods("GET")
	myRouter.HandleFunc("/recruiting/addrecruit/", controller.CreateRecruitPlayerProfile).Methods("POST")
	// myRouter.HandleFunc("/recruiting/allocaterecruitpoints/", controller.AllocateRecruitingPointsForRecruit).Methods("PUT")
	myRouter.HandleFunc("/recruiting/toggleScholarship/", controller.SendScholarshipToRecruit).Methods("POST")
	// myRouter.HandleFunc("/recruiting/revokescholarship/", controller.RevokeScholarshipFromRecruit).Methods("PUT")
	myRouter.HandleFunc("/recruiting/removecrootfromboard/", controller.RemoveRecruitFromBoard).Methods("PUT")
	myRouter.HandleFunc("/recruiting/savecrootboard/", controller.SaveRecruitingBoard).Methods("POST")

	// ReCroot Controls
	myRouter.HandleFunc("/recruits/all/", controller.AllRecruits).Methods("GET")
	myRouter.HandleFunc("/recruits/export/all/", controller.ExportCroots).Methods("GET")
	// myRouter.HandleFunc("/recruits/juco/all/", controller.AllJUCOCollegeRecruits).Methods("GET")
	myRouter.HandleFunc("/recruits/recruit/{recruitID}/", controller.GetCollegeRecruitByRecruitID).Methods("GET")
	myRouter.HandleFunc("/recruits/profile/recruits/{recruitProfileID}/", controller.GetRecruitsByTeamProfileID).Methods("GET")
	myRouter.HandleFunc("/recruits/recruit/create", controller.CreateCollegeRecruit).Methods("POST")
	// myRouter.HandleFunc("/recruits/recruit/update/", controller.UpdateCollegeRecruit).Methods("PUT")

	// Requests Controls
	myRouter.HandleFunc("/requests/all/", controller.GetTeamRequests).Methods("GET")
	myRouter.HandleFunc("/requests/create/", controller.CreateTeamRequest).Methods("POST")
	myRouter.HandleFunc("/requests/approve/", controller.ApproveTeamRequest).Methods("PUT")
	myRouter.HandleFunc("/requests/reject/", controller.RejectTeamRequest).Methods("DELETE")
	myRouter.HandleFunc("/requests/remove/{teamID}", controller.RemoveUserFromTeam).Methods("PUT")
	myRouter.HandleFunc("/nfl/requests/all/", controller.GetNFLTeamRequests).Methods("GET")
	myRouter.HandleFunc("/nfl/requests/create/", controller.CreateNFLTeamRequest).Methods("POST")
	myRouter.HandleFunc("/nfl/requests/approve/", controller.ApproveNFLTeamRequest).Methods("POST")
	myRouter.HandleFunc("/nfl/requests/reject/", controller.RejectNFLTeamRequest).Methods("DELETE")
	myRouter.HandleFunc("/nfl/requests/remove/{teamID}", controller.RemoveNFLUserFromNFLTeam).Methods("POST")

	// Standings Controls
	myRouter.HandleFunc("/standings/cfb/season/{seasonID}/", controller.GetAllCollegeStandings).Methods("GET")
	myRouter.HandleFunc("/standings/cfb/{conferenceID}/{seasonID}/", controller.GetCollegeStandingsByConferenceIDAndSeasonID).Methods("GET")
	myRouter.HandleFunc("/standings/nfl/season/{seasonID}/", controller.GetAllNFLStandings).Methods("GET")
	myRouter.HandleFunc("/standings/nfl/{divisionID}/{seasonID}/", controller.GetNFLStandingsByDivisionIDAndSeasonID).Methods("GET")
	myRouter.HandleFunc("/standings/cfb/history/team/{teamID}/", controller.GetHistoricalRecordsByTeamID).Methods("GET")

	// Stats Controls
	myRouter.HandleFunc("/statistics/export/cfb/", controller.ExportCFBStatisticsFromSim).Methods("POST")
	// myRouter.HandleFunc("/statistics/export/nfl/", controller.ExportNFLStatisticsFromSim).Methods("POST")
	myRouter.HandleFunc("/statistics/export/players/", controller.ExportPlayerStatsToCSV).Methods("GET")
	myRouter.HandleFunc("/statistics/injured/players/", controller.GetInjuryReport).Methods("GET")
	myRouter.HandleFunc("/statistics/interface/cfb/{seasonID}/{weekID}/{viewType}", controller.GetStatsPageContentForSeason).Methods("GET")
	myRouter.HandleFunc("/statistics/interface/nfl/{seasonID}/{weekID}/{viewType}", controller.GetNFLStatsPageContent).Methods("GET")
	myRouter.HandleFunc("/statistics/reset/season/", controller.ResetCFBSeasonalStats).Methods("GET")

	// Team Controls
	myRouter.HandleFunc("/teams/college/all/", controller.GetAllCollegeTeams).Methods("GET")
	myRouter.HandleFunc("/teams/nfl/all/", controller.GetAllNFLTeams).Methods("GET")
	myRouter.HandleFunc("/teams/nfl/roster/{teamID}/", controller.GetNFLRecordsForRosterPage).Methods("GET")
	myRouter.HandleFunc("/teams/college/active/", controller.GetAllActiveCollegeTeams).Methods("GET")
	myRouter.HandleFunc("/teams/college/available/", controller.GetAllAvailableCollegeTeams).Methods("GET")
	myRouter.HandleFunc("/teams/college/team/{teamID}/", controller.GetTeamByTeamID).Methods("GET")
	myRouter.HandleFunc("/teams/nfl/team/{teamID}/", controller.GetNFLTeamByTeamID).Methods("GET")
	myRouter.HandleFunc("/teams/college/conference/{conferenceID}/", controller.GetTeamsByConferenceID).Methods("GET")
	myRouter.HandleFunc("/teams/college/division/{divisionID}/", controller.GetTeamsByDivisionID).Methods("GET")
	myRouter.HandleFunc("/teams/college/sim/{gameID}/", controller.GetHomeAndAwayTeamData).Methods("GET")
	myRouter.HandleFunc("/teams/nfl/sim/{gameID}/", controller.GetNFLHomeAndAwayTeamData).Methods("GET")

	// Trade Controls
	myRouter.HandleFunc("/trades/nfl/all/accepted", controller.GetAllAcceptedTrades).Methods("GET")
	myRouter.HandleFunc("/trades/nfl/all/rejected", controller.GetAllRejectedTrades).Methods("GET")
	myRouter.HandleFunc("/trades/nfl/block/{teamID}", controller.GetNFLTradeBlockDataByTeamID).Methods("GET")
	myRouter.HandleFunc("/trades/nfl/place/block/{playerID}", controller.PlaceNFLPlayerOnTradeBlock).Methods("GET")
	myRouter.HandleFunc("/trades/nfl/preferences/update", controller.UpdateTradePreferences).Methods("POST")
	myRouter.HandleFunc("/trades/nfl/create/proposal", controller.CreateNFLTradeProposal).Methods("POST")
	myRouter.HandleFunc("/trades/nfl/proposal/accept/{proposalID}", controller.AcceptTradeOffer).Methods("GET")
	myRouter.HandleFunc("/trades/nfl/proposal/reject/{proposalID}", controller.RejectTradeOffer).Methods("GET")
	myRouter.HandleFunc("/trades/nfl/proposal/cancel/{proposalID}", controller.CancelTradeOffer).Methods("GET")

	// Discord Controls
	myRouter.HandleFunc("/teams/ds/college/team/{teamID}/", controller.GetTeamByTeamIDForDiscord).Methods("GET")
	myRouter.HandleFunc("/players/ds/college/player/{firstName}/{lastName}/{team}/{week}/", controller.GetCollegePlayerStatsByNameTeamAndWeek).Methods("GET")
	myRouter.HandleFunc("/players/ds/college/season/player/{firstName}/{lastName}/{team}/", controller.GetCurrentSeasonCollegePlayerStatsByNameTeam).Methods("GET")
	myRouter.HandleFunc("/teams/ds/college/week/team/{week}/{team}/", controller.GetWeeklyTeamStatsByTeamAbbrAndWeek).Methods("GET")
	myRouter.HandleFunc("/teams/ds/college/season/team/{season}/{team}/", controller.GetSeasonTeamStatsByTeamAbbrAndSeason).Methods("GET")
	myRouter.HandleFunc("/players/{firstName}/{lastName}/{teamID}", controller.GetCollegePlayerByNameAndTeam).Methods("GET")
	myRouter.HandleFunc("/croots/ds/class/{teamID}/", controller.GetRecruitingClassByTeamID).Methods("GET")
	myRouter.HandleFunc("/croots/ds/croot/{firstName}/{lastName}", controller.GetRecruitByFirstNameAndLastName).Methods("GET")
	myRouter.HandleFunc("/schedule/ds/current/week/{league}/", controller.GetCurrentWeekGamesByLeague).Methods("GET")
	// Easter Controls
	myRouter.HandleFunc("/easter/egg/collude/", controller.CollusionButton).Methods("POST")

	// Handle Controls
	// handler := cors.AllowAll().Handler(myRouter)
	loadEnvs()
	origins := os.Getenv("ORIGIN_ALLOWED")
	originsOk := handlers.AllowedOrigins([]string{origins})
	headersOk := handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "Accept", "Access-Control-Request-Method", "Access-Control-Request-Headers"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	log.Fatal(http.ListenAndServe(":5001", handlers.CORS(originsOk, methodsOk, headersOk)(myRouter)))
}

func loadEnvs() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("CANNOT LOAD ENV VARIABLES")
	}
}

func main() {
	InitialMigration()
	fmt.Println("Football Server Initialized.")

	handleRequests()
	fmt.Println("Hello There!")
}
