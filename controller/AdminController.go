package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/managers"
	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/ws"
	"github.com/gorilla/mux"
)

// GetTimeStamp
func GetCurrentTimestamp(w http.ResponseWriter, r *http.Request) {
	timestamp := managers.GetTimestamp()
	json.NewEncoder(w).Encode(timestamp)
}

// SyncWeek?
func SyncTimestamp(w http.ResponseWriter, r *http.Request) {
	var updateTimestampDto structs.UpdateTimestampDto
	err := json.NewDecoder(r.Body).Decode(&updateTimestampDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newTimestamp := managers.UpdateTimestamp(updateTimestampDto)
	json.NewEncoder(w).Encode(newTimestamp)
}

func SyncRecruiting(w http.ResponseWriter, r *http.Request) {
	ts := managers.GetTimestamp()
	managers.SyncRecruiting(ts)
	json.NewEncoder(w).Encode("Sync Complete")
}

func SyncWeek(w http.ResponseWriter, r *http.Request) {
	newTimestamp := managers.MoveUpWeek()
	json.NewEncoder(w).Encode(newTimestamp)
}

func SyncTimeslot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	timeslot := vars["timeslot"]
	if len(timeslot) == 0 {
		log.Panicln("Missing timeslot!")
	}
	managers.SyncTimeslot(timeslot)
	json.NewEncoder(w).Encode("Timeslot updated")
}

func SyncFreeAgencyRound(w http.ResponseWriter, r *http.Request) {
	managers.SyncFreeAgencyOffers()
	managers.MoveUpInOffseasonFreeAgency()
	json.NewEncoder(w).Encode("Moved to next free agency round")
}

func SyncMissingRES(w http.ResponseWriter, r *http.Request) {
	managers.SyncAllMissingEfficiencies()
}

func GetWeeksInSeason(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	seasonID := vars["seasonID"]
	weekID := vars["weekID"]
	weeks := managers.GetWeeksInASeason(seasonID, weekID)
	json.NewEncoder(w).Encode(weeks)
}

func GenerateWalkOns(w http.ResponseWriter, r *http.Request) {
	managers.GenerateWalkOns()
	fmt.Println(w, "Walk ons successfully generated.")
}

func SyncTeamRecruitingRanks(w http.ResponseWriter, r *http.Request) {
	managers.SyncTeamRankings()
	fmt.Println(w, "Team Ranks successfully generated.")
}

func ProgressToNextSeason(w http.ResponseWriter, r *http.Request) {
	managers.CFBProgressionMain()
	fmt.Println(w, "Team Ranks successfully generated.")
}

func ProgressNFL(w http.ResponseWriter, r *http.Request) {
	managers.NFLProgressionMain()
	fmt.Println(w, "Progressions Complete.")
}

func FillAIBoards(w http.ResponseWriter, r *http.Request) {
	managers.FillAIRecruitingBoards()
	fmt.Println(w, "Team Ranks successfully generated.")
}

func SyncAIBoards(w http.ResponseWriter, r *http.Request) {
	managers.ResetAIBoardsForCompletedTeams()
	managers.AllocatePointsToAIBoards()
	fmt.Println(w, "Team Ranks successfully generated.")
}

func AllocateAIRedshirts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	seasonID := vars["seasonID"]
	managers.AllocateAIRedshirts(seasonID)
	http.ResponseWriter(w).Write([]byte("Redshirts allocated to AI teams."))
}

func RunTheGames(w http.ResponseWriter, r *http.Request) {
	managers.FixBrokenGameplans()
	managers.CheckForSchemePenalties()
	managers.RunTheGames()
	ts := managers.GetTimestamp()
	json.NewEncoder(w).Encode(ts)
}

func UpdateTimestampAdmin(w http.ResponseWriter, r *http.Request) {
	UpdateTimestampDirect(w, r)
}

func GetCFBPostseasonGames(w http.ResponseWriter, r *http.Request) {
	games := managers.GetCFBPostseasonGamesForCurrentSeason()
	json.NewEncoder(w).Encode(games)
}

func DeleteCFBPostseasonGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameID"]
	if err := managers.DeleteCFBPostseasonGame(gameID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]bool{"deleted": true})
}

func SaveCFBConferenceChampionship(w http.ResponseWriter, r *http.Request) {
	var dto structs.CFBConferenceChampionshipDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	game, err := managers.SaveCFBConferenceChampionship(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(game)
}

func UploadCFBBowlGames(w http.ResponseWriter, r *http.Request) {
	var req structs.CFBBowlUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	games, err := managers.UploadCFBBowlGames(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(games)
}

func GenerateFBSPlayoffSchedule(w http.ResponseWriter, r *http.Request) {
	var req structs.CFBPlayoffScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	games, err := managers.GenerateFBSPlayoffSchedule(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(games)
}

func GenerateFCSPlayoffSchedule(w http.ResponseWriter, r *http.Request) {
	var req structs.CFBPlayoffScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	games, err := managers.GenerateFCSPlayoffSchedule(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(games)
}

func GetNFLPostseasonGames(w http.ResponseWriter, r *http.Request) {
	games := managers.GetNFLPostseasonGamesForCurrentSeason()
	json.NewEncoder(w).Encode(games)
}

func DeleteNFLPostseasonGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameID"]
	if err := managers.DeleteNFLPostseasonGame(gameID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]bool{"deleted": true})
}

func GenerateNFLPlayoffSchedule(w http.ResponseWriter, r *http.Request) {
	var req structs.NFLPlayoffScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	games, err := managers.GenerateNFLPlayoffSchedule(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(games)
}

func WeatherGenerator(w http.ResponseWriter, r *http.Request) {
	managers.GenerateWeatherForGames()
	fmt.Println(w, "Congrats, you generated the GODDAM WEATHER!")
}

func FixSmallTownBigCityAIBoards(w http.ResponseWriter, r *http.Request) {
	managers.FixSmallTownBigCityAIBoards()
	fmt.Println(w, "Affinities fixed!")
}

func ImportMissingDraftees(w http.ResponseWriter, r *http.Request) {
	managers.MigrateHistoricPlayersToNFLDraftees()
	fmt.Println(w, "Affinities fixed!")
}

func ImportMissingStats(w http.ResponseWriter, r *http.Request) {
	managers.MigrateCFBPlayerStatsFromPreviousSeason()
	fmt.Println(w, "NFL Stats fixed!")
}

func GenerateNewGameplans(w http.ResponseWriter, r *http.Request) {
	managers.MigrateCFBGameplansAndDepthChartsForRemainingFCSTeams()
	fmt.Println(w, "NFL Stats fixed!")
}

func SyncRES(w http.ResponseWriter, r *http.Request) {
	ts := managers.GetTimestamp()
	managers.SyncRecruitingEfficiency(ts)
}

func SyncToNextSeason(w http.ResponseWriter, r *http.Request) {
	db := dbprovider.GetInstance().GetDB()
	ts := managers.GetTimestamp()
	ts.MoveUpSeason()
	repository.SaveTimestamp(ts, db)
	managers.GenerateOffseasonData()
}

func NFLAttributePercentilesReport(w http.ResponseWriter, r *http.Request) {
	managers.NFLAttributePercentilesReport()
}

func CFBAttributePercentilesReport(w http.ResponseWriter, r *http.Request) {
	managers.CFBAttributePercentilesReport()
}

func UpdateTimestampDirect(w http.ResponseWriter, r *http.Request) {
	var dto structs.UpdateTimestampDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ts, err := managers.UpdateTimestampDirect(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ws.BroadcastTSUpdate(ts)
	json.NewEncoder(w).Encode(ts)
}

func AddCompensatoryDraftPick(w http.ResponseWriter, r *http.Request) {
	var dto structs.CompensatoryDraftPickDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pick, err := managers.AddCompensatoryDraftPick(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(pick)
}
