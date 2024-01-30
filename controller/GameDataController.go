package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/CalebRose/SimFBA/managers"
	"github.com/CalebRose/SimFBA/models"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/gorilla/mux"
)

// GetHomeAndAwayTeamData
func GetHomeAndAwayTeamData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameID"]

	gameChan := make(chan structs.CollegeGame)
	go func() {
		g := managers.GetCollegeGameByGameID(gameID)
		gameChan <- g
	}()
	game := <-gameChan
	close(gameChan)

	var responseModel models.SimGameDataResponse
	var waitgroup sync.WaitGroup
	waitgroup.Add(2)

	var rosterGroup sync.WaitGroup
	rosterGroup.Add(3)

	homeTeamChan := make(chan structs.CollegeTeam)
	awayTeamChan := make(chan structs.CollegeTeam)
	homeRosterChan := make(chan []structs.CollegePlayer)
	awayRosterChan := make(chan []structs.CollegePlayer)
	stadiumChan := make(chan structs.Stadium)

	go func() {
		waitgroup.Wait()
		close(homeTeamChan)
		close(awayTeamChan)
	}()

	go func() {
		rosterGroup.Wait()
		close(homeRosterChan)
		close(awayRosterChan)
		close(stadiumChan)
	}()

	go func() {
		defer waitgroup.Done()
		ht := managers.GetTeamByTeamAbbr(game.HomeTeam)
		homeTeamChan <- ht
	}()

	go func() {
		defer waitgroup.Done()
		at := managers.GetTeamByTeamAbbr(game.AwayTeam)
		awayTeamChan <- at
	}()

	homeTeam := <-homeTeamChan
	awayTeam := <-awayTeamChan
	homeTeamID := strconv.Itoa(int(homeTeam.ID))
	awayTeamID := strconv.Itoa(int(awayTeam.ID))
	stadiumID := strconv.Itoa(int(game.StadiumID))

	go func() {
		defer rosterGroup.Done()
		hr := managers.GetAllCollegePlayersByTeamIdWithoutRedshirts(homeTeamID)
		homeRosterChan <- hr
	}()

	go func() {
		defer rosterGroup.Done()
		ar := managers.GetAllCollegePlayersByTeamIdWithoutRedshirts(awayTeamID)
		awayRosterChan <- ar
	}()

	go func() {
		defer rosterGroup.Done()
		st := managers.GetStadiumByStadiumID(stadiumID)
		stadiumChan <- st
	}()

	homeTeamRoster := <-homeRosterChan
	awayTeamRoster := <-awayRosterChan
	stadium := <-stadiumChan

	var homeTeamResponse models.SimTeamDataResponse
	var homeDCResponse models.SimTeamDepthChartResponse
	var homeDCList []models.SimDepthChartPosResponse

	var awayTeamResponse models.SimTeamDataResponse
	var awayDCResponse models.SimTeamDepthChartResponse
	var awayDCList []models.SimDepthChartPosResponse

	hdc := homeTeam.TeamDepthChart

	for _, dcp := range hdc.DepthChartPlayers {
		var simDCPR models.SimDepthChartPosResponse
		simDCPR.Map(dcp)
		homeDCList = append(homeDCList, simDCPR)
	}

	adc := awayTeam.TeamDepthChart
	for _, dcp := range adc.DepthChartPlayers {
		var simDCPR models.SimDepthChartPosResponse
		simDCPR.Map(dcp)
		awayDCList = append(awayDCList, simDCPR)
	}

	homeDCResponse.Map(hdc, homeDCList)
	awayDCResponse.Map(adc, awayDCList)

	homeTeamResponse.Map(homeTeam, homeDCResponse)
	awayTeamResponse.Map(awayTeam, awayDCResponse)

	responseModel.AssignHomeTeam(homeTeamResponse, homeTeamRoster)
	responseModel.AssignAwayTeam(awayTeamResponse, awayTeamRoster)

	responseModel.AssignWeather(game.GameTemp, game.Cloud, game.Precip, game.WindCategory, game.WindSpeed)
	responseModel.AssignStadium(stadium)
	responseModel.AssignPostSeasonStatus(game.IsBowlGame || game.IsConferenceChampionship || game.IsNationalChampionship || game.IsPlayoffGame)

	json.NewEncoder(w).Encode(responseModel)
}

func GetNFLHomeAndAwayTeamData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameID"]

	gameChan := make(chan structs.NFLGame)
	go func() {
		g := managers.GetNFLGameByGameID(gameID)
		gameChan <- g
	}()
	game := <-gameChan
	close(gameChan)

	var responseModel models.NFLSimGameDataResponse

	var waitgroup sync.WaitGroup
	waitgroup.Add(2)
	var rosterGroup sync.WaitGroup
	rosterGroup.Add(3)

	homeTeamChan := make(chan structs.NFLTeam)
	awayTeamChan := make(chan structs.NFLTeam)
	homeRosterChan := make(chan []structs.NFLPlayer)
	awayRosterChan := make(chan []structs.NFLPlayer)
	stadiumChan := make(chan structs.Stadium)

	go func() {
		waitgroup.Wait()
		close(homeTeamChan)
		close(awayTeamChan)
	}()

	go func() {
		rosterGroup.Wait()
		close(homeRosterChan)
		close(awayRosterChan)
		close(stadiumChan)
	}()

	go func() {
		defer waitgroup.Done()
		ht := managers.GetNFLTeamByTeamIDForSim(strconv.Itoa(game.HomeTeamID))
		homeTeamChan <- ht
	}()

	go func() {
		defer waitgroup.Done()
		at := managers.GetNFLTeamByTeamIDForSim(strconv.Itoa(game.AwayTeamID))
		awayTeamChan <- at
	}()

	homeTeam := <-homeTeamChan
	awayTeam := <-awayTeamChan
	homeTeamID := strconv.Itoa(int(homeTeam.ID))
	awayTeamID := strconv.Itoa(int(awayTeam.ID))
	stadiumID := strconv.Itoa(int(game.StadiumID))

	go func() {
		defer rosterGroup.Done()
		hr := managers.GetNFLRosterForSimulation(homeTeamID)
		homeRosterChan <- hr
	}()

	go func() {
		defer rosterGroup.Done()
		ar := managers.GetNFLRosterForSimulation(awayTeamID)
		awayRosterChan <- ar
	}()

	go func() {
		defer rosterGroup.Done()
		st := managers.GetStadiumByStadiumID(stadiumID)
		stadiumChan <- st
	}()

	homeTeamRoster := <-homeRosterChan
	awayTeamRoster := <-awayRosterChan
	stadium := <-stadiumChan

	var homeTeamResponse models.NFLSimTeamDataResponse
	var homeDCResponse models.NFLSimTeamDepthChartResponse
	var homeDCList []models.NFLSimDepthChartPosResponse
	var awayTeamResponse models.NFLSimTeamDataResponse
	var awayDCResponse models.NFLSimTeamDepthChartResponse
	var awayDCList []models.NFLSimDepthChartPosResponse

	hdc := homeTeam.TeamDepthChart
	for _, dcp := range hdc.DepthChartPlayers {
		var simDCPR models.NFLSimDepthChartPosResponse
		simDCPR.Map(dcp)
		homeDCList = append(homeDCList, simDCPR)
	}
	adc := awayTeam.TeamDepthChart
	for _, dcp := range adc.DepthChartPlayers {
		var simDCPR models.NFLSimDepthChartPosResponse
		simDCPR.Map(dcp)
		awayDCList = append(awayDCList, simDCPR)
	}

	homeDCResponse.Map(hdc, homeDCList)
	awayDCResponse.Map(adc, awayDCList)
	homeTeamResponse.Map(homeTeam, homeDCResponse)
	awayTeamResponse.Map(awayTeam, awayDCResponse)
	responseModel.AssignHomeTeam(homeTeamResponse, homeTeamRoster)
	responseModel.AssignAwayTeam(awayTeamResponse, awayTeamRoster)
	responseModel.AssignWeather(game.GameTemp, game.Cloud, game.Precip, game.WindCategory, game.WindSpeed)
	responseModel.AssignStadium(stadium)
	responseModel.AssignPostSeasonStatus(game.IsConferenceChampionship || game.IsSuperBowl || game.IsPlayoffGame)
	json.NewEncoder(w).Encode(responseModel)
}
