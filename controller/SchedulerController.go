package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimFBA/managers"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/gorilla/mux"
)

func TestCFBScheduler(w http.ResponseWriter, r *http.Request) {
	managers.BaseGenerateCFBSchedule(true)
	json.NewEncoder(w).Encode(true)
}

// CreateCFBGameRequest accepts a CFBGameRequest body and persists it.
func CreateCFBGameRequest(w http.ResponseWriter, r *http.Request) {
	var request structs.CFBGameRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	managers.CreateCFBGameRequest(request)
	json.NewEncoder(w).Encode(true)
}

// AcceptCFBGameRequest marks the request as accepted.
func AcceptCFBGameRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := vars["requestID"]
	if len(requestID) == 0 {
		panic("User did not provide a requestID")
	}
	managers.AcceptCFBGameRequest(requestID)
	json.NewEncoder(w).Encode(true)
}

// RejectCFBGameRequest deletes the request.
func RejectCFBGameRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := vars["requestID"]
	if len(requestID) == 0 {
		panic("User did not provide a requestID")
	}
	managers.RejectCFBGameRequest(requestID)
	json.NewEncoder(w).Encode(true)
}

// ProcessCFBGameRequest converts an accepted request into a CollegeGame record.
func ProcessCFBGameRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := vars["requestID"]
	if len(requestID) == 0 {
		panic("User did not provide a requestID")
	}
	managers.ProcessCFBGameRequest(requestID)
	json.NewEncoder(w).Encode(true)
}

// VetoCFBGameRequest deletes the request via admin veto.
func VetoCFBGameRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := vars["requestID"]
	if len(requestID) == 0 {
		panic("User did not provide a requestID")
	}
	managers.VetoCFBGameRequest(requestID)
	json.NewEncoder(w).Encode(true)
}

// CreateNFLGameRequest accepts an NFLGameRequest body and persists it.
func CreateNFLGameRequest(w http.ResponseWriter, r *http.Request) {
	var request structs.NFLGameRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	managers.CreateNFLGameRequest(request)
	json.NewEncoder(w).Encode(true)
}

// AcceptNFLGameRequest marks the NFL request as accepted.
func AcceptNFLGameRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := vars["requestID"]
	if len(requestID) == 0 {
		panic("User did not provide a requestID")
	}
	managers.AcceptNFLGameRequest(requestID)
	json.NewEncoder(w).Encode(true)
}

// RejectNFLGameRequest deletes the NFL request.
func RejectNFLGameRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := vars["requestID"]
	if len(requestID) == 0 {
		panic("User did not provide a requestID")
	}
	managers.RejectNFLGameRequest(requestID)
	json.NewEncoder(w).Encode(true)
}

// ProcessNFLGameRequest converts an accepted NFL request into an NFLGame record.
func ProcessNFLGameRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := vars["requestID"]
	if len(requestID) == 0 {
		panic("User did not provide a requestID")
	}
	managers.ProcessNFLGameRequest(requestID)
	json.NewEncoder(w).Encode(true)
}

// VetoNFLGameRequest deletes the NFL request via admin veto.
func VetoNFLGameRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := vars["requestID"]
	if len(requestID) == 0 {
		panic("User did not provide a requestID")
	}
	managers.VetoNFLGameRequest(requestID)
	json.NewEncoder(w).Encode(true)
}

// SwapCFBHomeAndAwayTeams accepts a gameID and swaps the home and away teams for that game.
func SwapCFBHomeAndAwayTeams(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameID"]
	if len(gameID) == 0 {
		panic("User did not provide a gameID")
	}
	managers.SwapCFBGameHomeAndAwayTeams(gameID)
	json.NewEncoder(w).Encode(true)
}
