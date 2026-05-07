package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimFBA/managers"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/gorilla/mux"
)

// GetUDFABoardByTeamID - Returns the bidding board for a specific team
func GetUDFABoardByTeamID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]

	if len(teamID) == 0 {
		http.Error(w, "No Team ID provided", http.StatusBadRequest)
		return
	}

	board := managers.GetUDFABoardByTeamID(teamID)
	json.NewEncoder(w).Encode(board)
}

// AddPlayerToUDFABoard - Adds a player to the team's bidding list
func AddPlayerToUDFABoard(w http.ResponseWriter, r *http.Request) {
	var profile structs.NFLUDFAProfile
	err := json.NewDecoder(r.Body).Decode(&profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedProfile := managers.AddPlayerToUDFABoard(profile)
	json.NewEncoder(w).Encode(updatedProfile)
}

// SaveUDFABoard - Saves the point allocations (1-20)
func SaveUDFABoard(w http.ResponseWriter, r *http.Request) {
	var board structs.NFLUDFABoard
	err := json.NewDecoder(r.Body).Decode(&board)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = managers.SaveUDFABoard(board)
	if err != nil {
		http.Error(w, "Could not save UDFA board: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(true)
}

// RemovePlayerFromUDFABoard - Deletes a bid
func RemovePlayerFromUDFABoard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	profileID := vars["profileID"]

	managers.RemovePlayerFromUDFABoard(profileID)
	json.NewEncoder(w).Encode(true)
}