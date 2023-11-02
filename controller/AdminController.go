package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/CalebRose/SimFBA/managers"
	"github.com/CalebRose/SimFBA/structs"
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

func RegressTimeslot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	timeslot := vars["timeslot"]
	if len(timeslot) == 0 {
		log.Panicln("Missing timeslot!")
	}

	managers.RegressTimeslot(timeslot)

	json.NewEncoder(w).Encode("Timeslot updated")
}

func SyncFreeAgencyRound(w http.ResponseWriter, r *http.Request) {
	managers.MoveUpInOffseasonFreeAgency()
	managers.SyncFreeAgencyOffers()
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

// CreateCollegeRecruit?

// CreateNFLPlayer -- Create NFL Player from template, and then synthetically progress them based on the year of input

// UpdateTeamRecruitingProfile

// ApproveCoachForTeam

// RemoveCoachFromTeam

// UpdateTeam

// RunProgressionsForCollege
func RunProgressionsForCollege(w http.ResponseWriter, r *http.Request) {

}

// GenerateWalkons
func GenerateWalkOns(w http.ResponseWriter, r *http.Request) {
	managers.GenerateWalkOns()
	fmt.Println(w, "Walk ons successfully generated.")
}

// RunProgressionsForNFL

// RunProgressionsForJuco?

func SyncTeamRecruitingRanks(w http.ResponseWriter, r *http.Request) {
	managers.SyncTeamRankings()
	fmt.Println(w, "Team Ranks successfully generated.")
}

func ProgressToNextSeason(w http.ResponseWriter, r *http.Request) {
	managers.ProgressionMain()
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

func RunTheGames(w http.ResponseWriter, r *http.Request) {
	managers.RunTheGames()
	fmt.Println(w, "Games for current week are set to run.")
}

func WeatherGenerator(w http.ResponseWriter, r *http.Request) {
	managers.GenerateWeatherForGames()
	fmt.Println(w, "Congrats, you generated the GODDAM WEATHER!")
}

func FixSmallTownBigCityAIBoards(w http.ResponseWriter, r *http.Request) {
	managers.FixSmallTownBigCityAIBoards()
	fmt.Println(w, "Affinities fixed!")
}
