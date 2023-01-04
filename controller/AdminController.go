package controller

import (
	"encoding/json"
	"fmt"
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

func SyncMissingRES(w http.ResponseWriter, r *http.Request) {
	managers.SyncAllMissingEfficiencies()
}

func GetNewsLogs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	weekID := vars["weekID"]
	seasonID := vars["seasonID"]

	newsLogs := managers.GetNewsLogs(weekID, seasonID)

	json.NewEncoder(w).Encode(newsLogs)
}

func GetAllNewsLogsForASeason(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	seasonID := vars["seasonID"]

	newsLogs := managers.GetAllNewsLogs(seasonID)

	json.NewEncoder(w).Encode(newsLogs)
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

func FillAIBoards(w http.ResponseWriter, r *http.Request) {
	managers.FillAIRecruitingBoards()
	fmt.Println(w, "Team Ranks successfully generated.")
}
