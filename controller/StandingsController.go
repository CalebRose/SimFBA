package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimFBA/managers"
	"github.com/gorilla/mux"
)

// GetCollegeStandingsByConferenceIDAndSeasonID
func GetCollegeStandingsByConferenceIDAndSeasonID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	conferenceID := vars["conferenceID"]
	seasonID := vars["seasonID"]
	if len(conferenceID) == 0 || len(seasonID) == 0 {
		panic("User did not provide enough information")
	}
	standings := managers.GetStandingsByConferenceIDAndSeasonID(conferenceID, seasonID)
	json.NewEncoder(w).Encode(standings)
}

func GetNFLStandingsByDivisionIDAndSeasonID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	divisionID := vars["divisionID"]
	seasonID := vars["seasonID"]
	if len(divisionID) == 0 || len(seasonID) == 0 {
		panic("User did not provide enough information")
	}
	standings := managers.GetNFLStandingsByDivisionIDAndSeasonID(divisionID, seasonID)
	json.NewEncoder(w).Encode(standings)
}

// GetCollegeStandingsByDivisionIDAndSeasonID

// GetHistoricalRecordsByTeamID
func GetHistoricalRecordsByTeamID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	if len(teamID) == 0 {
		panic("User did not provide enough information")
	}
	standings := managers.GetHistoricalRecordsByTeamID(teamID)
	json.NewEncoder(w).Encode(standings)
}

// GetAllCollegeStandings
func GetAllCollegeStandings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	seasonID := vars["seasonID"]
	if len(seasonID) == 0 {
		panic("User did not provide enough information")
	}
	standings := managers.GetAllCollegeStandingsBySeasonID(seasonID)
	json.NewEncoder(w).Encode(standings)
}

// Gets all NFL Standings in a Season
func GetAllNFLStandings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	seasonID := vars["seasonID"]
	if len(seasonID) == 0 {
		panic("User did not provide enough information")
	}
	standings := managers.GetAllNFLStandingsBySeasonID(seasonID)
	json.NewEncoder(w).Encode(standings)
}

func CreateCollegeStandings(w http.ResponseWriter, r *http.Request) {
	managers.GenerateNewSeasonStandings()
	json.NewEncoder(w).Encode("Standings Generated")
}

func GenerateCollegeRankingsHandler(w http.ResponseWriter, r *http.Request) {
	managers.BaseGetCollegeRankings()
	json.NewEncoder(w).Encode("College Rankings Generated")
}
