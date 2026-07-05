package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimFBA/managers"
	"github.com/gorilla/mux"
)

// GetBulkPlayByPlay returns the massive array of plays to feed the frontend spoofing loop
func GetBulkPlayByPlay(w http.ResponseWriter, r *http.Request) {
	isCollege := r.URL.Query().Get("isCollege") == "true"
	season := r.URL.Query().Get("season")
	week := r.URL.Query().Get("week")
	timeslot := r.URL.Query().Get("timeslot")

	response := managers.GetBulkPlayByPlayData(isCollege, season, week, timeslot)
	json.NewEncoder(w).Encode(response)
}

// RunAdminGames manually triggers the game engine via POST from the Control Room
func RunAdminGames(w http.ResponseWriter, r *http.Request) {
	// managers.RunGames()
	json.NewEncoder(w).Encode("Live Broadcast Engine Started!")
}

func TestCFBCronJob(w http.ResponseWriter, r *http.Request) {
	managers.StartCFBLiveStreamingCron()
}

func TestNFLCronJob(w http.ResponseWriter, r *http.Request) {
	managers.StartNFLLiveStreamingCron()
}

func GetFBGameQueue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	league := vars["league"]
	weekID := r.URL.Query().Get("weekID")
	seasonID := r.URL.Query().Get("seasonID")
	gameDay := r.URL.Query().Get("gameDay")
	isPreseason := r.URL.Query().Get("isPreseason") == "true"
	queue := managers.BuildStreamGameQueue(league, weekID, seasonID, gameDay, isPreseason)
	json.NewEncoder(w).Encode(queue)
}
