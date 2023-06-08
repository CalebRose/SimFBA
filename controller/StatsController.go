package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/CalebRose/SimFBA/managers"
	"github.com/CalebRose/SimFBA/models"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/gorilla/mux"
)

func ExportCFBStatisticsFromSim(w http.ResponseWriter, r *http.Request) {
	// Create DTO for College Recruit
	var exportStatsDTO structs.ExportStatsDTO
	fmt.Println("PING!")
	err := json.NewDecoder(r.Body).Decode(&exportStatsDTO)
	if err != nil {
		fmt.Println("CANNOT DECODE BODY!")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Send DTO to Manager Class
	if len(exportStatsDTO.CFBGameStatDTOs) > 0 {
		managers.ExportCFBStatisticsFromSim(exportStatsDTO.CFBGameStatDTOs)
	}
	if len(exportStatsDTO.NFLGameStatDTOs) > 0 {
		managers.ExportNFLStatisticsFromSim(exportStatsDTO.NFLGameStatDTOs)
	}

	fmt.Println(w, "Game Data Exported")
}

func ExportPlayerStatsToCSV(w http.ResponseWriter, r *http.Request) {

	ts := managers.GetTimestamp()

	teamsChan := make(chan []structs.CollegeTeam)

	go func() {
		ct := managers.GetAllCollegeTeams()
		teamsChan <- ct
	}()

	collegeTeams := <-teamsChan
	close(teamsChan)

	var conferenceMap = make(map[int]int)
	var conferenceNameMap = make(map[int]string)

	for _, team := range collegeTeams {
		conferenceMap[int(team.ID)] = team.ConferenceID
		conferenceNameMap[int(team.ID)] = team.Conference
	}

	playersChan := make(chan []models.CollegePlayerResponse)
	go func() {
		cp := managers.GetAllCollegePlayersWithStatsBySeasonID(conferenceMap, conferenceNameMap, strconv.Itoa(ts.CollegeSeasonID), "", "SEASON")
		playersChan <- cp
	}()

	collegePlayers := <-playersChan
	close(playersChan)

	managers.ExportPlayerStatsToCSV(collegePlayers, w)
}

func GetStatsPageContentForSeason(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	seasonID := vars["seasonID"]
	viewType := vars["viewType"]
	weekID := vars["weekID"]

	if len(viewType) == 0 {
		panic("User did not provide view type")
	}

	if len(seasonID) == 0 {
		panic("User did not provide TeamID")
	}

	teamsChan := make(chan []models.CollegeTeamResponse)

	go func() {
		ct := managers.GetAllCollegeTeamsWithStatsBySeasonID(seasonID, weekID, viewType)
		teamsChan <- ct
	}()

	collegeTeams := <-teamsChan
	close(teamsChan)

	var conferenceMap = make(map[int]int)
	var conferenceNameMap = make(map[int]string)

	for _, team := range collegeTeams {
		conferenceMap[int(team.ID)] = team.ConferenceID
		conferenceNameMap[int(team.ID)] = team.Conference
	}

	playersChan := make(chan []models.CollegePlayerResponse)
	go func() {
		cp := managers.GetAllCollegePlayersWithStatsBySeasonID(conferenceMap, conferenceNameMap, seasonID, weekID, viewType)
		playersChan <- cp
	}()

	collegePlayers := <-playersChan
	close(playersChan)

	confChan := make(chan []structs.CollegeConference)
	go func() {
		cf := managers.GetCollegeConferences()
		confChan <- cf
	}()

	collegeConferences := <-confChan
	close(confChan)

	response := models.SimCFBStatsResponse{
		CollegePlayers:     collegePlayers,
		CollegeTeams:       collegeTeams,
		CollegeConferences: collegeConferences,
	}

	json.NewEncoder(w).Encode(response)
}

func GetNFLStatsPageContent(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	vars := mux.Vars(r)
	seasonID := vars["seasonID"]
	viewType := vars["viewType"]
	weekID := vars["weekID"]

	if len(viewType) == 0 {
		panic("User did not provide view type")
	}

	if len(seasonID) == 0 {
		panic("User did not provide TeamID")
	}

	teamsChan := make(chan []models.NFLTeamResponse)

	go func() {
		ct := managers.GetAllNFLTeamsWithStatsBySeasonID(seasonID, weekID, viewType)
		teamsChan <- ct
	}()

	nflTeams := <-teamsChan
	close(teamsChan)

	var conferenceMap = make(map[int]int)
	var conferenceNameMap = make(map[int]string)
	var divisionMap = make(map[int]int)
	var divisionNameMap = make(map[int]string)
	for _, team := range nflTeams {
		conferenceMap[int(team.ID)] = team.ConferenceID
		conferenceNameMap[int(team.ID)] = team.Conference
		divisionMap[int(team.ID)] = team.DivisionID
		divisionNameMap[int(team.ID)] = team.Division
	}

	playersChan := make(chan []models.NFLPlayerResponse)
	go func() {
		cp := managers.GetAllNFLPlayersWithStatsBySeasonID(conferenceMap, divisionMap, conferenceNameMap, divisionNameMap, seasonID, weekID, viewType)
		playersChan <- cp
	}()

	nflPlayers := <-playersChan
	close(playersChan)

	response := models.SimNFLStatsResponse{
		NFLPlayers: nflPlayers,
		NFLTeams:   nflTeams,
	}

	json.NewEncoder(w).Encode(response)
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

func MapAllStatsToSeason(w http.ResponseWriter, r *http.Request) {
	managers.MapAllStatsToSeason()
}

func GetCollegeGameResultsByGameID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameID"]
	if len(gameID) == 0 {
		panic("User did not provide a first name")
	}

	player := managers.GetCFBGameResultsByGameID(gameID)

	json.NewEncoder(w).Encode(player)
}

func GetNFLGameResultsByGameID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameID"]
	if len(gameID) == 0 {
		panic("User did not provide a first name")
	}

	player := managers.GetNFLGameResultsByGameID(gameID)

	json.NewEncoder(w).Encode(player)
}
