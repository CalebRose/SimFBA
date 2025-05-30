package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CalebRose/SimFBA/managers"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/gorilla/mux"
)

// GetAllCollegeTeamsForRosterPage
func GetAllCollegeTeamsForRosterPage(w http.ResponseWriter, r *http.Request) {
	collegeTeams := managers.GetAllCollegeTeamsForRosterPage()

	json.NewEncoder(w).Encode(collegeTeams)
}

// GetAllCollegeTeams
func GetAllCollegeTeams(w http.ResponseWriter, r *http.Request) {
	collegeTeams := managers.GetAllCollegeTeams()

	json.NewEncoder(w).Encode(collegeTeams)
}

// GetAllNFLTeams
func GetAllNFLTeams(w http.ResponseWriter, r *http.Request) {
	nflTeams := managers.GetAllNFLTeams()

	json.NewEncoder(w).Encode(nflTeams)
}

// GetAllActiveCollegeTeams
func GetAllActiveCollegeTeams(w http.ResponseWriter, r *http.Request) {
	collegeTeams := managers.GetAllCoachedCollegeTeams()

	json.NewEncoder(w).Encode(collegeTeams)
}

// GetAllAvailableCollegeTeams
func GetAllAvailableCollegeTeams(w http.ResponseWriter, r *http.Request) {
	collegeTeams := managers.GetAllAvailableCollegeTeams()

	json.NewEncoder(w).Encode(collegeTeams)
}

// GetAllAvailableNFLTeams

// GetAllCoachedNFLTeams

// GetTeamByTeamID
func GetTeamByTeamID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	if len(teamID) == 0 {
		panic("User did not provide TeamID")
	}
	team := managers.GetTeamByTeamID(teamID)
	json.NewEncoder(w).Encode(team)
}

func GetNFLRecordsForRosterPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	if len(teamID) == 0 {
		panic("User did not provide TeamID")
	}
	team := managers.GetNFLRecordsForRosterPage(teamID)
	json.NewEncoder(w).Encode(team)
}

func GetNFLTeamByTeamID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	if len(teamID) == 0 {
		panic("User did not provide TeamID")
	}
	team := managers.GetNFLTeamWithCapsheetByTeamID(teamID)
	json.NewEncoder(w).Encode(team)
}

// GetTeamsByConferenceID
func GetTeamsByConferenceID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	conferenceID := vars["conferenceID"]
	if len(conferenceID) == 0 {
		panic("User did not provide conferenceID")
	}
	team := managers.GetTeamByTeamID(conferenceID)
	json.NewEncoder(w).Encode(team)
}

// GetTeamsByDivisionID
func GetTeamsByDivisionID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	divisionID := vars["divisionID"]
	if len(divisionID) == 0 {
		panic("User did not provide divisionID")
	}
	team := managers.GetTeamByTeamID(divisionID)
	json.NewEncoder(w).Encode(team)
}

// GetTeamsByDivisionID
func GetRecruitingClassSizeForTeams(w http.ResponseWriter, r *http.Request) {
	managers.GetRecruitingClassSizeForTeams()
	json.NewEncoder(w).Encode("Sync for Class Size complete")
}

func GetCFBDashboardByTeamID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	if len(teamID) == 0 {
		panic("User did not provide teamID")
	}
	// Schedule, Standings, News Logs, top players by stats
	dashboard := managers.GetDashboardByTeamID(true, teamID)
	json.NewEncoder(w).Encode(dashboard)
}

func GetNFLDashboardByTeamID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	if len(teamID) == 0 {
		panic("User did not provide teamID")
	}
	// Schedule, Standings, News Logs, top players by stats
	dashboard := managers.GetDashboardByTeamID(false, teamID)
	json.NewEncoder(w).Encode(dashboard)
}

func AssignCFBTeamGrades(w http.ResponseWriter, r *http.Request) {
	managers.AssignTeamGrades()
}

func UpdateCFBJersey(w http.ResponseWriter, r *http.Request) {
	// Create DTO for College Recruit
	var jerseyDTO structs.JerseyDTO
	err := json.NewDecoder(r.Body).Decode(&jerseyDTO)
	if err != nil {
		fmt.Println("CANNOT DECODE BODY!")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.UpdateCFBJersey(jerseyDTO)

	fmt.Println(w, "Game Jersey Updated")
}

func UpdateNFLJersey(w http.ResponseWriter, r *http.Request) {
	// Create DTO for College Recruit
	var jerseyDTO structs.JerseyDTO
	err := json.NewDecoder(r.Body).Decode(&jerseyDTO)
	if err != nil {
		fmt.Println("CANNOT DECODE BODY!")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.UpdateNFLJersey(jerseyDTO)

	fmt.Println(w, "Game Jersey Updated")
}
