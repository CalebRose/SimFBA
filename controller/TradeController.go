package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CalebRose/SimFBA/managers"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/gorilla/mux"
)

// Get Trade Block Data for Trade Block Page
func GetNFLTradeBlockDataByTeamID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	if len(teamID) == 0 {
		panic("User did not provide TeamID")
	}

	response := managers.GetTradeBlockDataByTeamID(teamID)

	json.NewEncoder(w).Encode(response)
}

// Get Trade Block Data for Trade Block Page
func GetAllAcceptedTrades(w http.ResponseWriter, r *http.Request) {
	response := managers.GetAcceptedTradeProposals()

	json.NewEncoder(w).Encode(response)
}

// Get Trade Block Data for Trade Block Page
func GetAllRejectedTrades(w http.ResponseWriter, r *http.Request) {
	response := managers.GetRejectedTradeProposals()

	json.NewEncoder(w).Encode(response)
}

// Place player on NFL Trade block
func PlaceNFLPlayerOnTradeBlock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playerID := vars["playerID"]
	if len(playerID) == 0 {
		panic("User did not provide playerID")
	}

	managers.PlaceNFLPlayerOnTradeBlock(playerID)

	json.NewEncoder(w).Encode("Player " + playerID + " placed on trade block.")
}

// Update Trade Preferences
func UpdateTradePreferences(w http.ResponseWriter, r *http.Request) {

	var tradePreferenceDTO structs.NFLTradePreferencesDTO
	err := json.NewDecoder(r.Body).Decode(&tradePreferenceDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.UpdateTradePreferences(tradePreferenceDTO)

	fmt.Fprintf(w, "Trade Preferences Updated")
}

// Create NFL Trade Proposal
func CreateNFLTradeProposal(w http.ResponseWriter, r *http.Request) {

	var tradeProposalDTO structs.NFLTradeProposalDTO
	err := json.NewDecoder(r.Body).Decode(&tradeProposalDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.CreateTradeProposal(tradeProposalDTO)

}

// Accept Trade Offer
func AcceptTradeOffer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proposalID := vars["proposalID"]
	if len(proposalID) == 0 {
		panic("User did not provide a proposalID")
	}

	managers.AcceptTradeProposal(proposalID)
}

// Reject Trade Offer
func RejectTradeOffer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proposalID := vars["proposalID"]
	if len(proposalID) == 0 {
		panic("User did not provide a proposalID")
	}

	managers.RejectTradeProposal(proposalID)
}

// Cancels Trade Offer
func CancelTradeOffer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proposalID := vars["proposalID"]
	if len(proposalID) == 0 {
		panic("User did not provide a proposalID")
	}

	managers.CancelTradeProposal(proposalID)
}

// SyncAcceptedTrade -- Admin approve a trade
func SyncAcceptedTrade(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proposalID := vars["proposalID"]
	if len(proposalID) == 0 {
		panic("User did not provide a proposalID")
	}

	managers.SyncAcceptedTrade(proposalID)
}

// SyncTradeFromDraftPage -- Admin approve a trade
func SyncTradeFromDraftPage(w http.ResponseWriter, r *http.Request) {
	var tradeProposalDTO structs.NFLTradeProposal
	err := json.NewDecoder(r.Body).Decode(&tradeProposalDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	managers.SyncTradeFromDraftPage(tradeProposalDTO)
}

// SyncAcceptedTrade -- Admin approve a trade
func VetoAcceptedTrade(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proposalID := vars["proposalID"]
	if len(proposalID) == 0 {
		panic("User did not provide a proposalID")
	}

	managers.VetoTrade(proposalID)
}

// CleanUpRejectedTrades -- Remove all rejected trades from the DB
func CleanUpRejectedTrades(w http.ResponseWriter, r *http.Request) {
	managers.RemoveRejectedTrades()
}
