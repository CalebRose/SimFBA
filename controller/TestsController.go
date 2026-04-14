package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/CalebRose/SimFBA/managers"
	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/gorilla/mux"
)

func TestCFBProgressionAlgorithm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")

	managers.CFBProgressionExport(w)
}

func TestNFLProgressionAlgorithm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")

	managers.NFLProgressionExport(w)
}

func UpdateCollegeAIDepthChartsTEST(w http.ResponseWriter, r *http.Request) {
	// managers.SetAIGameplan()
	managers.UpdateCollegeAIDepthChartsTEST()
	json.NewEncoder(w).Encode("Updated all CFB Depth Charts")
}

func MassCFBUpdateGameplansTEST(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	off := vars["off"]
	if len(off) == 0 {
		panic("User did not provide a teamID")
	}
	def := vars["def"]
	if len(def) == 0 {
		panic("User did not provide a teamID")
	}
	managers.MassUpdateGameplanSchemesTEST(off, def)
	json.NewEncoder(w).Encode("Updated all CFB Depth Charts For Testing")
}

func UpdateCFBIndividualGameplanTEST(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	if len(teamID) == 0 {
		panic("User did not provide a teamID")
	}
	off := vars["off"]
	if len(off) == 0 {
		panic("User did not provide a teamID")
	}
	def := vars["def"]
	if len(def) == 0 {
		panic("User did not provide a teamID")
	}
	managers.UpdateIndividualGameplanSchemeTEST(teamID, off, def)
	json.NewEncoder(w).Encode("Updated all CFB Depth Charts For Testing")
}

// GetCFBHomeAndAwayTeamTestData
func GetCFBHomeAndAwayTeamTestData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameID"]

	gameChan := make(chan structs.CollegeGame)
	go func() {
		g := managers.GetCollegeGameByGameID(gameID)
		gameChan <- g
	}()
	game := <-gameChan
	close(gameChan)

	responseModel := structs.SimGameDataResponseTEST{
		GameID:   int(game.ID),
		WeekID:   game.WeekID,
		SeasonID: game.SeasonID,
	}
	var waitgroup sync.WaitGroup
	waitgroup.Add(6)

	var rosterGroup sync.WaitGroup
	rosterGroup.Add(3)

	homeTeamChan := make(chan structs.CollegeTeam)
	awayTeamChan := make(chan structs.CollegeTeam)
	hDepthChartChan := make(chan structs.CollegeTeamDepthChartTEST)
	aDepthChartChan := make(chan structs.CollegeTeamDepthChartTEST)
	hGameplanChan := make(chan structs.CollegeGameplanTEST)
	aGameplanChan := make(chan structs.CollegeGameplanTEST)
	homeRosterChan := make(chan []structs.CollegePlayer)
	awayRosterChan := make(chan []structs.CollegePlayer)
	stadiumChan := make(chan structs.Stadium)
	htID := strconv.Itoa(game.HomeTeamID)
	atID := strconv.Itoa(game.AwayTeamID)

	go func() {
		waitgroup.Wait()
		close(hDepthChartChan)
		close(aDepthChartChan)
		close(hGameplanChan)
		close(aGameplanChan)
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
		ht := managers.GetTeamByTeamID(htID)
		homeTeamChan <- ht
	}()

	go func() {
		defer waitgroup.Done()
		at := managers.GetTeamByTeamID(atID)
		awayTeamChan <- at
	}()

	go func() {
		defer waitgroup.Done()
		hg := repository.GetGameplanTESTByTeamID(htID)
		hGameplanChan <- hg
	}()

	go func() {
		defer waitgroup.Done()
		hdc := repository.GetDCTESTByTeamID(htID)
		hDepthChartChan <- hdc
	}()

	go func() {
		defer waitgroup.Done()
		ag := repository.GetGameplanTESTByTeamID(atID)
		aGameplanChan <- ag
	}()

	go func() {
		defer waitgroup.Done()
		adc := repository.GetDCTESTByTeamID(atID)
		aDepthChartChan <- adc
	}()

	homeTeam := <-homeTeamChan
	awayTeam := <-awayTeamChan
	homeTeamDC := <-hDepthChartChan
	awayTeamDC := <-aDepthChartChan
	hGP := <-hGameplanChan
	aGP := <-aGameplanChan
	stadiumID := strconv.Itoa(int(game.StadiumID))

	go func() {
		defer rosterGroup.Done()
		hr := managers.GetAllCollegePlayersByTeamIdWithoutRedshirts(htID)
		homeRosterChan <- hr
	}()

	go func() {
		defer rosterGroup.Done()
		ar := managers.GetAllCollegePlayersByTeamIdWithoutRedshirts(atID)
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

	var homeTeamResponse structs.SimTeamDataResponseTEST
	var homeDCResponse structs.SimTeamDepthChartResponseTEST
	var homeDCList []structs.SimDepthChartPosResponseTEST

	var awayTeamResponse structs.SimTeamDataResponseTEST
	var awayDCResponse structs.SimTeamDepthChartResponseTEST
	var awayDCList []structs.SimDepthChartPosResponseTEST

	for _, dcp := range homeTeamDC.DepthChartPlayers {
		var simDCPR structs.SimDepthChartPosResponseTEST
		simDCPR.Map(dcp)
		homeDCList = append(homeDCList, simDCPR)
	}

	for _, dcp := range awayTeamDC.DepthChartPlayers {
		var simDCPR structs.SimDepthChartPosResponseTEST
		simDCPR.Map(dcp)
		awayDCList = append(awayDCList, simDCPR)
	}

	homeDCResponse.Map(homeTeamDC, homeDCList)
	awayDCResponse.Map(awayTeamDC, awayDCList)

	homeTeamResponse.Map(homeTeam, hGP, homeDCResponse, game.HomePreviousBye)
	awayTeamResponse.Map(awayTeam, aGP, awayDCResponse, game.AwayPreviousBye)

	responseModel.AssignHomeTeam(homeTeamResponse, homeTeamRoster)
	responseModel.AssignAwayTeam(awayTeamResponse, awayTeamRoster)

	responseModel.AssignWeather(game.GameTemp, game.Cloud, game.Precip, game.WindCategory, game.WindSpeed)
	responseModel.AssignStadium(stadium)
	responseModel.AssignPostSeasonStatus(game.IsBowlGame || game.IsConferenceChampionship || game.IsNationalChampionship || game.IsPlayoffGame)

	json.NewEncoder(w).Encode(responseModel)
}

// GetNFLHomeAndAwayTeamTestData
func GetNFLHomeAndAwayTeamTestData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameID"]

	gameChan := make(chan structs.NFLGame)
	go func() {
		g := managers.GetNFLGameByGameID(gameID)
		gameChan <- g
	}()
	game := <-gameChan
	close(gameChan)

	responseModel := structs.SimNFLGameDataResponseTEST{
		GameID:   int(game.ID),
		WeekID:   game.WeekID,
		SeasonID: game.SeasonID,
	}
	var waitgroup sync.WaitGroup
	waitgroup.Add(6)

	var rosterGroup sync.WaitGroup
	rosterGroup.Add(3)

	homeTeamChan := make(chan structs.NFLTeam)
	awayTeamChan := make(chan structs.NFLTeam)
	hDepthChartChan := make(chan structs.NFLDepthChartTEST)
	aDepthChartChan := make(chan structs.NFLDepthChartTEST)
	hGameplanChan := make(chan structs.NFLGameplanTEST)
	aGameplanChan := make(chan structs.NFLGameplanTEST)
	homeRosterChan := make(chan []structs.NFLPlayer)
	awayRosterChan := make(chan []structs.NFLPlayer)
	stadiumChan := make(chan structs.Stadium)
	htID := strconv.Itoa(game.HomeTeamID)
	atID := strconv.Itoa(game.AwayTeamID)

	go func() {
		waitgroup.Wait()
		close(hDepthChartChan)
		close(aDepthChartChan)
		close(hGameplanChan)
		close(aGameplanChan)
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
		ht := managers.GetNFLTeamByTeamIDForSim(htID)
		homeTeamChan <- ht
	}()

	go func() {
		defer waitgroup.Done()
		at := managers.GetNFLTeamByTeamIDForSim(atID)
		awayTeamChan <- at
	}()

	go func() {
		defer waitgroup.Done()
		hg := repository.GetNFLGameplanTESTByTeamID(htID)
		hGameplanChan <- hg
	}()

	go func() {
		defer waitgroup.Done()
		hdc := repository.GetNFLDCTESTByTeamID(htID)
		hDepthChartChan <- hdc
	}()

	go func() {
		defer waitgroup.Done()
		ag := repository.GetNFLGameplanTESTByTeamID(atID)
		aGameplanChan <- ag
	}()

	go func() {
		defer waitgroup.Done()
		adc := repository.GetNFLDCTESTByTeamID(atID)
		aDepthChartChan <- adc
	}()

	homeTeam := <-homeTeamChan
	awayTeam := <-awayTeamChan
	homeTeamDC := <-hDepthChartChan
	awayTeamDC := <-aDepthChartChan
	hGP := <-hGameplanChan
	aGP := <-aGameplanChan
	stadiumID := strconv.Itoa(int(game.StadiumID))

	go func() {
		defer rosterGroup.Done()
		hr := managers.GetNFLRosterForSimulation(htID)
		homeRosterChan <- hr
	}()

	go func() {
		defer rosterGroup.Done()
		ar := managers.GetNFLRosterForSimulation(atID)
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

	var homeTeamResponse structs.SimNFLTeamDataResponseTEST
	var homeDCResponse structs.SimNFLTeamDepthChartResponseTEST
	var homeDCList []structs.SimNFLDepthChartPosResponseTEST

	var awayTeamResponse structs.SimNFLTeamDataResponseTEST
	var awayDCResponse structs.SimNFLTeamDepthChartResponseTEST
	var awayDCList []structs.SimNFLDepthChartPosResponseTEST

	for _, dcp := range homeTeamDC.DepthChartPlayers {
		var simDCPR structs.SimNFLDepthChartPosResponseTEST
		simDCPR.Map(dcp)
		homeDCList = append(homeDCList, simDCPR)
	}

	for _, dcp := range awayTeamDC.DepthChartPlayers {
		var simDCPR structs.SimNFLDepthChartPosResponseTEST
		simDCPR.Map(dcp)
		awayDCList = append(awayDCList, simDCPR)
	}

	homeDCResponse.Map(homeTeamDC, homeDCList)
	awayDCResponse.Map(awayTeamDC, awayDCList)

	homeTeamResponse.Map(homeTeam, hGP, homeDCResponse, game.HomePreviousBye)
	awayTeamResponse.Map(awayTeam, aGP, awayDCResponse, game.AwayPreviousBye)

	responseModel.AssignHomeTeam(homeTeamResponse, homeTeamRoster)
	responseModel.AssignAwayTeam(awayTeamResponse, awayTeamRoster)

	responseModel.AssignWeather(game.GameTemp, game.Cloud, game.Precip, game.WindCategory, game.WindSpeed)
	responseModel.AssignStadium(stadium)
	responseModel.AssignPostSeasonStatus(game.IsConferenceChampionship || game.IsSuperBowl || game.IsPlayoffGame)

	json.NewEncoder(w).Encode(responseModel)
}

func SetUpNFLTestDataStructs(w http.ResponseWriter, r *http.Request) {

	managers.SetupNFLTestDataStructs()
	json.NewEncoder(w).Encode("Set Up NFL Test Data Structs")
}

func UpdateNFLAIDepthChartsTEST(w http.ResponseWriter, r *http.Request) {
	// managers.SetAIGameplan()
	managers.UpdateNFLAIDepthChartsTEST()
	json.NewEncoder(w).Encode("Updated all NFL Depth Charts")
}

func MassNFLUpdateGameplansTEST(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	off := vars["off"]
	if len(off) == 0 {
		panic("User did not provide a teamID")
	}
	def := vars["def"]
	if len(def) == 0 {
		panic("User did not provide a teamID")
	}
	managers.MassUpdateNFLGameplanSchemesTEST(off, def)
	json.NewEncoder(w).Encode("Updated all NFL Depth Charts For Testing")
}

func UpdateNFLIndividualGameplanTEST(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	if len(teamID) == 0 {
		panic("User did not provide a teamID")
	}
	off := vars["off"]
	if len(off) == 0 {
		panic("User did not provide a teamID")
	}
	def := vars["def"]
	if len(def) == 0 {
		panic("User did not provide a teamID")
	}
	managers.UpdateIndividualNFLGameplanSchemeTEST(teamID, off, def)
	json.NewEncoder(w).Encode("Updated all NFL Depth Charts For Testing")
}
