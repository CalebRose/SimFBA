package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/CalebRose/SimFBA/managers"
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

	// Turn off Run Games Boolean
	managers.RunTheGames()

	fmt.Println(w, "Game Data Exported")
}

func ExportPlayerStatsToCSV(w http.ResponseWriter, r *http.Request) {

	ts := managers.GetTimestamp()
	_, gt := ts.GetCFBCurrentGameType()

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

	playersChan := make(chan []structs.CollegePlayerResponse)
	go func() {
		cp := managers.GetAllCollegePlayersWithStatsBySeasonID(conferenceMap, conferenceNameMap, strconv.Itoa(ts.CollegeSeasonID), "", "SEASON", gt)
		playersChan <- cp
	}()

	collegePlayers := <-playersChan
	close(playersChan)

	managers.ExportPlayerStatsToCSV(collegePlayers, w)
}

func ExportNFLPlayByPlayToCSV(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameID"]
	if len(gameID) == 0 {
		log.Panicln("Missing game ID!")
	}
	managers.ExportNFLPlayByPlayToCSV(gameID, w)
}

func ExportPlayByPlayToCSV(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameID"]
	if len(gameID) == 0 {
		log.Panicln("Missing game ID!")
	}
	managers.ExportCFBPlayByPlayToCSV(gameID, w)
}

func ExportAllSpringGamesPbPToCSV(w http.ResponseWriter, r *http.Request) {
	managers.ExportCFBSpringPlayByPlayToCSV(w)
}

func ExportAllPreseasonGamesPbPToCSV(w http.ResponseWriter, r *http.Request) {
	managers.ExportNFLPreseasonPlayByPlayToCSV(w)
}

func GetInjuryReport(w http.ResponseWriter, r *http.Request) {

	// GetInjuredCollegePlayers
	collegePlayers := managers.GetInjuredCollegePlayers()

	// GetInjuredNFLPlayers
	nflPlayers := managers.GetInjuredNFLPlayers()

	response := structs.InjuryReportResponse{
		CollegePlayers: collegePlayers,
		NFLPlayers:     nflPlayers,
	}

	json.NewEncoder(w).Encode(response)
}

func GetStatsPageContentForSeason(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	seasonID := vars["seasonID"]
	viewType := vars["viewType"]
	weekID := vars["weekID"]
	gameType := vars["gameType"]

	if len(viewType) == 0 {
		panic("User did not provide view type")
	}

	if len(seasonID) == 0 {
		panic("User did not provide TeamID")
	}

	teamsChan := make(chan []structs.CollegeTeamResponse)

	go func() {
		ct := managers.GetAllCollegeTeamsWithStatsBySeasonID(seasonID, weekID, viewType, gameType)
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

	playersChan := make(chan []structs.CollegePlayerResponse)
	go func() {
		cp := managers.GetAllCollegePlayersWithStatsBySeasonID(conferenceMap, conferenceNameMap, seasonID, weekID, viewType, gameType)
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

	response := structs.SimCFBStatsResponse{
		CollegePlayers:     collegePlayers,
		CollegeTeams:       collegeTeams,
		CollegeConferences: collegeConferences,
	}

	json.NewEncoder(w).Encode(response)
}

func ExportStatsPageContentForSeason(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	seasonID := vars["seasonID"]
	viewType := vars["viewType"]
	weekID := vars["weekID"]
	gameType := vars["gameType"]

	if len(viewType) == 0 {
		panic("User did not provide view type")
	}

	if len(seasonID) == 0 {
		panic("User did not provide TeamID")
	}

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

	playersChan := make(chan []structs.CollegePlayerResponse)
	go func() {
		cp := managers.GetAllCollegePlayersWithStatsBySeasonID(conferenceMap, conferenceNameMap, seasonID, weekID, viewType, gameType)
		playersChan <- cp
	}()

	collegePlayers := <-playersChan
	close(playersChan)

	managers.ExportCollegePlayerStatsToCSV(collegePlayers, viewType, w)
}

func GetNFLStatsPageContent(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	vars := mux.Vars(r)
	seasonID := vars["seasonID"]
	viewType := vars["viewType"]
	weekID := vars["weekID"]
	gameType := vars["gameType"]
	if len(viewType) == 0 {
		panic("User did not provide view type")
	}

	if len(seasonID) == 0 {
		panic("User did not provide TeamID")
	}

	teamsChan := make(chan []structs.NFLTeamResponse)

	go func() {
		ct := managers.GetAllNFLTeamsWithStatsBySeasonID(seasonID, weekID, viewType, gameType)
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

	playersChan := make(chan []structs.NFLPlayerResponse)
	go func() {
		cp := managers.GetAllNFLPlayersWithStatsBySeasonID(conferenceMap, divisionMap, conferenceNameMap, divisionNameMap, seasonID, weekID, viewType, gameType)
		playersChan <- cp
	}()

	nflPlayers := <-playersChan
	close(playersChan)

	response := structs.SimNFLStatsResponse{
		NFLPlayers: nflPlayers,
		NFLTeams:   nflTeams,
	}

	json.NewEncoder(w).Encode(response)
}

func ExportNFLStatsPageContent(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	vars := mux.Vars(r)
	seasonID := vars["seasonID"]
	viewType := vars["viewType"]
	weekID := vars["weekID"]
	gameType := vars["gameType"]

	if len(viewType) == 0 {
		panic("User did not provide view type")
	}

	if len(seasonID) == 0 {
		panic("User did not provide TeamID")
	}

	teamsChan := make(chan []structs.NFLTeamResponse)

	go func() {
		ct := managers.GetAllNFLTeamsWithStatsBySeasonID(seasonID, weekID, viewType, gameType)
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

	playersChan := make(chan []structs.NFLPlayerResponse)
	go func() {
		cp := managers.GetAllNFLPlayersWithStatsBySeasonID(conferenceMap, divisionMap, conferenceNameMap, divisionNameMap, seasonID, weekID, viewType, gameType)
		playersChan <- cp
	}()

	nflPlayers := <-playersChan
	close(playersChan)

	managers.ExportNFLPlayerStatsToCSV(nflPlayers, viewType, w)
}

func ResetCFBSeasonalStats(w http.ResponseWriter, r *http.Request) {
	managers.MigrateCFBPlayerStatsFromPreviousSeason()
	managers.MigrateCFBTeamSnapsFromPreviousSeason()
}

func ResetNFLSeasonalStats(w http.ResponseWriter, r *http.Request) {
	managers.MigrateNFLPlayerStatsFromPreviousSeason()
	managers.MigrateNFLTeamSnapsFromPreviousSeason()
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

func GetCFBSeasonStatsRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playerID := vars["playerID"]
	seasonID := vars["seasonID"]
	gameType := vars["gameType"]
	if len(playerID) == 0 {
		panic("User did not provide a first name")
	}

	player := managers.GetCollegePlayerSeasonStatsByPlayerIDAndSeason(playerID, seasonID, gameType)

	json.NewEncoder(w).Encode(player)
}

func GetCFBStatsPageContent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	seasonID := vars["seasonID"]
	viewType := vars["viewType"]
	weekID := vars["weekID"]
	gameType := vars["gameType"]

	if len(viewType) == 0 {
		panic("User did not provide view type")
	}

	if len(seasonID) == 0 {
		panic("User did not provide TeamID")
	}

	response := managers.SearchCollegeStats(seasonID, weekID, viewType, gameType)
	json.NewEncoder(w).Encode(response)
}

func GetProStatsPageContent(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	vars := mux.Vars(r)
	seasonID := vars["seasonID"]
	viewType := vars["viewType"]
	weekID := vars["weekID"]
	gameType := vars["gameType"]

	response := managers.SearchProStats(seasonID, weekID, viewType, gameType)

	json.NewEncoder(w).Encode(response)
}
