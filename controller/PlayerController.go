package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CalebRose/SimFBA/managers"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/gorilla/mux"
)

// AllPlayers - Get All Players Record references in table
func AllPlayers(w http.ResponseWriter, r *http.Request) {
	var players = managers.GetAllPlayers()

	json.NewEncoder(w).Encode(players)
}

// AllCollegePlayers

// AllCollegePlayersByTeamID
func AllCollegePlayersByTeamID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamId := vars["teamID"]

	if len(teamId) == 0 {
		panic("User did not provide TeamID")
	}

	players := managers.GetAllCollegePlayersByTeamId(teamId)
	promises := managers.GetCollegePromisesByTeamID(teamId)

	response := structs.CFBRosterPageResponse{
		Players:  players,
		Promises: promises,
	}

	json.NewEncoder(w).Encode(response)
}

// AllCollegePlayersByTeamIDWithoutRedshirts
func AllNFLPlayersByTeamIDForDC(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamId := vars["teamID"]

	if len(teamId) == 0 {
		panic("User did not provide TeamID")
	}

	players := managers.GetNFLPlayersForDCPage(teamId)

	json.NewEncoder(w).Encode(players)
}

// AllCollegePlayersByTeamIDWithoutRedshirts
func AllCollegePlayersByTeamIDWithoutRedshirts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamId := vars["teamID"]

	if len(teamId) == 0 {
		panic("User did not provide TeamID")
	}

	players := managers.GetAllCollegePlayersByTeamIdWithoutRedshirts(teamId)

	json.NewEncoder(w).Encode(players)
}

// GetCollegePlayerByID
func GetCollegePlayerByCollegePlayerId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collegePlayerID := vars["playerID"]

	if len(collegePlayerID) == 0 {
		panic("User did not provide TeamID")
	}

	players := managers.GetCollegePlayerByCollegePlayerId(collegePlayerID)

	json.NewEncoder(w).Encode(players)
}

func GetHeismanList(w http.ResponseWriter, r *http.Request) {

	heismanList := managers.GetAllPostSeasonAwardsLists()

	json.NewEncoder(w).Encode(heismanList)
}

// UpdateCollegePlayer
func UpdateCollegePlayer(w http.ResponseWriter, r *http.Request) {

	fmt.Println(w, "College Player successfully updated.")
}

// ToggleRedshirtStatusForPlayer
func ToggleRedshirtStatusForPlayer(w http.ResponseWriter, r *http.Request) {
	// Update DTO for College Player
	vars := mux.Vars(r)
	PlayerID := vars["PlayerID"]

	managers.SetRedshirtStatusForPlayer(PlayerID)

	fmt.Println(w, "College Player successfully redshirted.")
}

// ToggleRedshirtStatusForPlayer
func CutCFBPlayerFromRoster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	PlayerID := vars["PlayerID"]

	managers.CutCFBPlayer(PlayerID)

	fmt.Println(w, "NFL Player Cut from Roster")
}

// CutNFLPlayer
func CutNFLPlayerFromRoster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	PlayerID := vars["PlayerID"]

	managers.CutNFLPlayer(PlayerID)

	fmt.Println(w, "NFL Player Cut from Roster")
}

func ExportRosterToCSV(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")

	vars := mux.Vars(r)
	teamId := vars["teamID"]

	if len(teamId) == 0 {
		panic("User did not provide TeamID")
	}

	managers.ExportTeamToCSV(teamId, w)
}

func ExportNFLRosterToCSV(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")

	vars := mux.Vars(r)
	teamId := vars["teamID"]

	if len(teamId) == 0 {
		panic("User did not provide TeamID")
	}

	managers.ExportNFLTeamToCSV(teamId, w)
	// ?
}

func ExportAllCFBRostersToCSV(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")

	managers.ExportAllRostersToCSV(w)
	// ?
}

func ExportAllNFLRostersToCSV(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")

	managers.ExportAllNFLTeamsToCSV(w)
}

// Place player on NFL Trade block
func PlaceNFLPlayerOnPracticeSquad(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playerID := vars["PlayerID"]
	if len(playerID) == 0 {
		panic("User did not provide playerID")
	}

	managers.PlaceNFLPlayerOnPracticeSquad(playerID)

	json.NewEncoder(w).Encode("Player " + playerID + " placed on trade block.")
}

// Place player on NFL Trade block
func PlaceNFLPlayerOnInjuryReserve(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playerID := vars["PlayerID"]
	if len(playerID) == 0 {
		panic("User did not provide playerID")
	}

	managers.PlaceNFLPlayerOnInjuryReserve(playerID)

	json.NewEncoder(w).Encode("Player " + playerID + " placed on trade block.")
}

// CreateExtensionOffer - Extend Offer to NFL player to extend contract with existing team
func CreateExtensionOffer(w http.ResponseWriter, r *http.Request) {
	var freeAgencyOfferDTO structs.FreeAgencyOfferDTO
	err := json.NewDecoder(r.Body).Decode(&freeAgencyOfferDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var offer = managers.CreateExtensionOffer(freeAgencyOfferDTO)

	json.NewEncoder(w).Encode(offer)
}

// CancelExtensionOffer - Cancel an extension offer with an NFL player
func CancelExtensionOffer(w http.ResponseWriter, r *http.Request) {
	var freeAgencyOfferDTO structs.FreeAgencyOfferDTO
	err := json.NewDecoder(r.Body).Decode(&freeAgencyOfferDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.CancelExtensionOffer(freeAgencyOfferDTO)

	json.NewEncoder(w).Encode(true)
}

func ExportPortalPlayersToCSV(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")
	managers.ExportTransferPortalToCSV(w)
}
