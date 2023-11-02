package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimFBA/managers"
	"github.com/gorilla/mux"
)

// Flex: Compare Two Program's history against one another
func CompareTeams(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamOneID := vars["teamOneID"]
	if len(teamOneID) == 0 {
		panic("User did not provide teamID")
	}

	teamTwoID := vars["teamTwoID"]
	if len(teamTwoID) == 0 {
		panic("User did not provide teamID")
	}

	res := managers.CompareTwoTeams(teamOneID, teamTwoID)

	json.NewEncoder(w).Encode(res)
}

func GetTeamByTeamIDForDiscord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	if len(teamID) == 0 {
		panic("User did not provide TeamID")
	}
	team := managers.GetTeamByTeamIDForDiscord(teamID)
	json.NewEncoder(w).Encode(team)
}

func GetCollegePlayerStatsByNameTeamAndWeek(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	firstName := vars["firstName"]
	lastName := vars["lastName"]
	teamID := vars["team"]
	week := vars["week"]

	if len(firstName) == 0 {
		panic("User did not provide a first name")
	}

	player := managers.GetCollegePlayerByNameTeamAndWeek(firstName, lastName, teamID, week)

	json.NewEncoder(w).Encode(player)
}

func GetCurrentSeasonCollegePlayerStatsByNameTeam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	firstName := vars["firstName"]
	lastName := vars["lastName"]
	teamID := vars["team"]

	if len(firstName) == 0 {
		panic("User did not provide a first name")
	}

	player := managers.GetSeasonalCollegePlayerByNameTeam(firstName, lastName, teamID)

	json.NewEncoder(w).Encode(player)
}

func GetWeeklyTeamStatsByTeamAbbrAndWeek(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["team"]
	week := vars["week"]

	if len(teamID) == 0 {
		panic("User did not provide a first name")
	}

	team := managers.GetTeamStatsByWeekAndTeam(teamID, week)

	json.NewEncoder(w).Encode(team)
}

func GetSeasonTeamStatsByTeamAbbrAndSeason(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["team"]
	season := vars["season"]

	if len(teamID) == 0 {
		panic("User did not provide a first name")
	}

	team := managers.GetSeasonalTeamStats(teamID, season)

	json.NewEncoder(w).Encode(team)
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

func GetRecruitingClassByTeamID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]

	if len(teamID) == 0 {
		panic("User did not provide teamID")
	}

	recruitingProfile := managers.GetRecruitingClassByTeamID(teamID)

	json.NewEncoder(w).Encode(recruitingProfile)
}

func GetRecruitByFirstNameAndLastName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	firstName := vars["firstName"]
	lastName := vars["lastName"]

	if len(firstName) == 0 {
		panic("User did not provide a first name")
	}

	recruit := managers.GetCollegeRecruitByName(firstName, lastName)

	json.NewEncoder(w).Encode(recruit)
}

// GetCollegeGamesByTeamIDAndSeasonID
func GetCurrentWeekGamesByLeague(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	league := vars["league"]
	if len(league) == 0 {
		panic("User did not provide teamID")
	}

	if league == "simcfb" {
		collegeGames := managers.GetCFBCurrentWeekSchedule()
		json.NewEncoder(w).Encode(collegeGames)
	} else {
		nflgames := managers.GetNFLCurrentWeekSchedule()
		json.NewEncoder(w).Encode(nflgames)
	}
}

// GetCollegeGamesByConference
func GetCollegeGamesByConference(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	conference := vars["conference"]
	if len(conference) == 0 {
		panic("User did not provide conference name")
	}

	collegeGames := managers.GetCFBScheduleByConference(conference)
	json.NewEncoder(w).Encode(collegeGames)
}
