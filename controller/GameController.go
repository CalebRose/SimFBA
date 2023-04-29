package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimFBA/managers"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/gorilla/mux"
)

// GetCollegeGamesByWeekId
func GetCollegeGamesByWeekId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	weekID := vars["weekID"]
	seasonID := vars["seasonID"]

	if len(weekID) == 0 {
		panic("User did not provide weekID")
	}

	collegeGames := managers.GetCollegeGamesByWeekIdAndSeasonID(weekID, seasonID)

	json.NewEncoder(w).Encode(collegeGames)
}

// GetCollegeGamesByTimeslotAndWeekID
func GetCollegeGamesByTimeslotWeekId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	weekID := vars["weekID"]
	timeSlot := vars["timeSlot"]

	if len(weekID) == 0 {
		panic("User did not provide weekID")
	}

	if len(timeSlot) == 0 {
		panic("No time slot selected!")
	}

	collegeGames := managers.GetCollegeGamesByTimeslotAndWeekId(weekID, timeSlot)

	json.NewEncoder(w).Encode(collegeGames)
}

// GetCollegeGamesByTeamIDAndSeasonID
func GetCollegeGamesByTeamIDAndSeasonID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]

	seasonID := vars["seasonID"]

	if len(teamID) == 0 {
		panic("User did not provide teamID")
	}

	collegeGames := managers.GetTeamScheduleForBot(teamID, seasonID)

	json.NewEncoder(w).Encode(collegeGames)
}

func GetCollegeGamesBySeasonID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	seasonID := vars["seasonID"]

	if len(seasonID) == 0 {
		panic("User did not provide teamID")
	}

	collegeGames := managers.GetCollegeGamesBySeasonID(seasonID)

	json.NewEncoder(w).Encode(collegeGames)
}

// GetNFLGamesByWeekID

// GetNFLGamesByTimeslotAndWeekID

// GetNFLGamesByTeamIDAndSeasonID
func GetNFLGamesByTeamIDAndSeasonID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]

	seasonID := vars["seasonID"]

	if len(teamID) == 0 {
		panic("User did not provide teamID")
	}

	nflGames := managers.GetNFLGamesByTeamIdAndSeasonId(teamID, seasonID)

	json.NewEncoder(w).Encode(nflGames)
}

func GetNFLGamesBySeasonID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	seasonID := vars["seasonID"]

	if len(seasonID) == 0 {
		panic("User did not provide teamID")
	}

	nflGames := managers.GetNFLGamesBySeasonID(seasonID)

	json.NewEncoder(w).Encode(nflGames)
}

func UpdateTimeslot(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	w.Header().Set("Content-Type", "application/json")
	var timeslotDTO structs.UpdateTimeslotDTO
	err := json.NewDecoder(r.Body).Decode(&timeslotDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	g1, g2 := managers.UpdateTimeslot(timeslotDTO)
	if g1.ID > 0 {
		json.NewEncoder(w).Encode(g1)
	} else {
		json.NewEncoder(w).Encode(g2)
	}
}

func GetWeatherForecast(w http.ResponseWriter, r *http.Request) {
	res := managers.GetCurrentWeekWeather()
	json.NewEncoder(w).Encode(res)
}
