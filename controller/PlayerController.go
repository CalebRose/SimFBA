package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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

	json.NewEncoder(w).Encode(players)
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

// GetCollegePlayerByNameAndTeam
func GetCollegePlayerByNameAndTeam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	firstName := vars["firstName"]
	lastName := vars["lastName"]
	teamID := vars["teamID"]

	if len(firstName) == 0 {
		panic("User did not provide a first name")
	}

	player := managers.GetCollegePlayerByNameAndTeam(firstName, lastName, teamID)

	json.NewEncoder(w).Encode(player)
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

	heismanList := managers.GetHeismanList()

	json.NewEncoder(w).Encode(heismanList)
}

// UpdateCollegePlayer
func UpdateCollegePlayer(w http.ResponseWriter, r *http.Request) {

	fmt.Println(w, "College Player successfully updated.")
}

// ToggleRedshirtStatusForPlayer
func ToggleRedshirtStatusForPlayer(w http.ResponseWriter, r *http.Request) {
	// Update DTO for College Player
	var redshirtDTO structs.RedshirtDTO
	err := json.NewDecoder(r.Body).Decode(&redshirtDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.SetRedshirtStatusForPlayer(strconv.Itoa(redshirtDTO.PlayerID))

	fmt.Println(w, "College Player successfully redshirted.")
}

// ToggleRedshirtStatusForPlayer
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

	// ?
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

func ExportAllRostersToCSV(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")

	managers.ExportAllRostersToCSV(w)
	// ?
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
