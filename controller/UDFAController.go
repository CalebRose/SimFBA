package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimFBA/managers"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/gorilla/mux"
)

// GetUDFABoardByTeamID returns the user's specific UDFA bidding board
func GetUDFABoardByTeamID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	if len(teamID) == 0 {
		panic("User did not provide TeamID")
	}

	board := managers.GetUDFABoardByTeamID(teamID)
	json.NewEncoder(w).Encode(board)
}

// AddPlayerToUDFABoard adds a single player to the user's board
func AddPlayerToUDFABoard(w http.ResponseWriter, r *http.Request) {
	var profile structs.NFLUDFAProfile
	err := json.NewDecoder(r.Body).Decode(&profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call the manager to execute the DB logic
	managers.AddPlayerToUDFABoard(profile)

	// Return success to the frontend
	json.NewEncoder(w).Encode(true)
}

// SaveUDFABoard updates the points assigned to all players on the board
func SaveUDFABoard(w http.ResponseWriter, r *http.Request) {
	var board structs.NFLUDFABoard
	err := json.NewDecoder(r.Body).Decode(&board)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call the manager to execute the DB logic
	managers.SaveUDFABoard(board)

	// Return success to the frontend
	json.NewEncoder(w).Encode(true)
}

// RemovePlayerFromUDFABoard deletes a player from the user's board
func RemovePlayerFromUDFABoard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	profileID := vars["profileID"]
	if len(profileID) == 0 {
		panic("User did not provide profileID")
	}

	// Call the manager to execute the DB logic
	managers.RemovePlayerFromUDFABoard(profileID)

	// Return success to the frontend
	json.NewEncoder(w).Encode(true)
}

// ProcessUDFAs triggers the batch signing process (Used by Admin)
func ProcessUDFAs(w http.ResponseWriter, r *http.Request) {
	dryRunParam := r.URL.Query().Get("dryRun")
	isDryRun := dryRunParam == "true"

	// Trigger the batch signing logic
	managers.ProcessUDFAs(isDryRun)

	// Return a clean success message to display in the Admin Panel
	msg := "LIVE UDFA Signings Executed Successfully!"
	if isDryRun {
		msg = "Dry Run Complete. Check backend server console for simulation results."
	}

	json.NewEncoder(w).Encode(msg)
}
