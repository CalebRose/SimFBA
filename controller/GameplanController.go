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

// GetGameplanByTeamID
func GetTeamGameplanByTeamID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	if len(teamID) == 0 {
		panic("User did not provide a teamID")
	}

	gamePlan := managers.GetGameplanByTeamID(teamID)

	json.NewEncoder(w).Encode(gamePlan)
}

// GetDepthChartByTeamID
func GetTeamDepthchartByTeamID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	if len(teamID) == 0 {
		panic("User did not provide a teamID")
	}

	depthchart := managers.GetDepthchartByTeamID(teamID)

	json.NewEncoder(w).Encode(depthchart)
}

// GetDepthChartByTeamID
func GetNFLDepthChart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	if len(teamID) == 0 {
		panic("User did not provide a teamID")
	}

	depthchart := managers.GetNFLDepthchartByTeamID(teamID)

	json.NewEncoder(w).Encode(depthchart)
}

// GetDepthChartPositionsByDepthChartID
func GetDepthChartPositionsByDepthChartID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	depthChartID := vars["depthChartID"]
	if len(depthChartID) == 0 {
		panic("User did not provide a depthChartID")
	}

	positionPlayers := managers.GetDepthChartPositionPlayersByDepthchartID(depthChartID)

	json.NewEncoder(w).Encode(positionPlayers)
}

// UpdateGameplan
func UpdateGameplan(w http.ResponseWriter, r *http.Request) {
	var updateGameplanDto structs.UpdateGameplanDTO
	err := json.NewDecoder(r.Body).Decode(&updateGameplanDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.UpdateGameplan(updateGameplanDto)

	fmt.Println("Updated Team Gameplan")
}

// UpdateDepthChart
func UpdateDepthChart(w http.ResponseWriter, r *http.Request) {
	var updateDepthChartDto structs.UpdateDepthChartDTO
	err := json.NewDecoder(r.Body).Decode(&updateDepthChartDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.UpdateDepthChart(updateDepthChartDto)

	fmt.Println("Updated Depth Chart for Team " + strconv.Itoa(updateDepthChartDto.DepthChartID))
}

// NFL //
// GetGameplanByTeamID
func GetNFLGameplanByTeamID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	if len(teamID) == 0 {
		panic("User did not provide a teamID")
	}

	gamePlan := managers.GetNFLGameplanByTeamID(teamID)

	json.NewEncoder(w).Encode(gamePlan)
}

// UpdateGameplan
func UpdateNFLGameplan(w http.ResponseWriter, r *http.Request) {
	var updateGameplanDto structs.UpdateGameplanDTO
	err := json.NewDecoder(r.Body).Decode(&updateGameplanDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.UpdateNFLGameplan(updateGameplanDto)

	fmt.Println("Updated Team Gameplan")
}

// UpdateDepthChart
func UpdateNFLDepthChart(w http.ResponseWriter, r *http.Request) {
	var updateDepthChartDto structs.UpdateNFLDepthChartDTO
	err := json.NewDecoder(r.Body).Decode(&updateDepthChartDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.UpdateNFLDepthChart(updateDepthChartDto)

	fmt.Println("Updated Depth Chart for Team " + strconv.Itoa(updateDepthChartDto.DepthChartID))
}

// UpdateCollegeAIDepthCharts
func UpdateCollegeAIDepthCharts(w http.ResponseWriter, r *http.Request) {
	managers.UpdateCollegeAIDepthCharts()
	json.NewEncoder(w).Encode("Updated all CFB Depth Charts")
}

// UpdateCollegeAIDepthCharts
func UpdateNFLAIDepthCharts(w http.ResponseWriter, r *http.Request) {
	managers.UpdateNFLAIDepthCharts()
	json.NewEncoder(w).Encode("Updated all NFL Depth Charts")
}
