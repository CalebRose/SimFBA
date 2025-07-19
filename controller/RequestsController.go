package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimFBA/managers"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/gorilla/mux"
)

func GetFBARequests(w http.ResponseWriter, r *http.Request) {
	requests := managers.GetAllFBARequests()
	json.NewEncoder(w).Encode(requests)
}

func GetTeamRequests(w http.ResponseWriter, r *http.Request) {
	requests := managers.GetAllTeamRequests()

	json.NewEncoder(w).Encode(requests)
}

func GetNFLTeamRequests(w http.ResponseWriter, r *http.Request) {
	requests := managers.GetAllNFLTeamRequests()

	json.NewEncoder(w).Encode(requests)
}

func CreateTeamRequest(w http.ResponseWriter, r *http.Request) {
	var request structs.TeamRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.CreateTeamRequest(request)

	json.NewEncoder(w).Encode(request)
}

func ApproveTeamRequest(w http.ResponseWriter, r *http.Request) {
	var request structs.TeamRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil || request.ID == 0 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.ApproveTeamRequest(request)
	json.NewEncoder(w).Encode(request)

}

func RejectTeamRequest(w http.ResponseWriter, r *http.Request) {
	var request structs.TeamRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.RejectTeamRequest(request)
	json.NewEncoder(w).Encode(request)

}

func RemoveUserFromTeam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	if len(teamID) == 0 {
		panic("User did not provide TeamID")
	}

	managers.RemoveUserFromTeam(teamID)
	json.NewEncoder(w).Encode(true)

}

func ViewCFBTeamUponRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	if len(teamID) == 0 {
		panic("User did not provide TeamID")
	}

	team := managers.GetCFBTeamForAvailableTeamsPage(teamID)

	json.NewEncoder(w).Encode(team)
}

func ViewNFLTeamUponRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	if len(teamID) == 0 {
		panic("User did not provide TeamID")
	}

	team := managers.GetNFLTeamForAvailableTeamsPage(teamID)

	json.NewEncoder(w).Encode(team)
}

func CreateNFLTeamRequest(w http.ResponseWriter, r *http.Request) {
	var request structs.NFLRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.CreateNFLTeamRequest(request)
	json.NewEncoder(w).Encode(request)

}

func ApproveNFLTeamRequest(w http.ResponseWriter, r *http.Request) {
	var request structs.NFLRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil || request.ID == 0 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.ApproveNFLTeamRequest(request)
	json.NewEncoder(w).Encode(request)

}

func RejectNFLTeamRequest(w http.ResponseWriter, r *http.Request) {
	var request structs.NFLRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.RejectNFLTeamRequest(request)
	json.NewEncoder(w).Encode(request)

}

func RemoveNFLUserFromNFLTeam(w http.ResponseWriter, r *http.Request) {
	var request structs.NFLRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.RemoveUserFromNFLTeam(request)

	json.NewEncoder(w).Encode(true)
}
