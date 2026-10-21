package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimFBA/managers"
	"github.com/gorilla/mux"
)

func GetHistoricalCollegeGameplansBySeasonAndWeek(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	seasonID := vars["seasonID"]
	weekID := vars["weekID"]
	if len(seasonID) == 0 || len(weekID) == 0 {
		panic("User did not provide both a seasonID and a weekID")
	}
	gameplans := managers.GetHistoricalCollegeGameplansBySeasonAndWeek(seasonID, weekID)
	json.NewEncoder(w).Encode(gameplans)
}

func GetHistoricalNFLGameplansBySeasonAndWeek(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	seasonID := vars["seasonID"]
	weekID := vars["weekID"]
	if len(seasonID) == 0 || len(weekID) == 0 {
		panic("User did not provide both a seasonID and a weekID")
	}
	gameplans := managers.GetHistoricalNFLGameplansBySeasonAndWeek(seasonID, weekID)
	json.NewEncoder(w).Encode(gameplans)
}

func BackfillGameplans(w http.ResponseWriter, r *http.Request) {
	managers.BackfillCurrentGameplans()
	json.NewEncoder(w).Encode("Backfill Complete")
}
