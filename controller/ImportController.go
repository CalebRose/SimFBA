package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimFBA/managers"
)

func ImportCustomCroots(w http.ResponseWriter, r *http.Request) {
	managers.CreateCustomCroots()
}

func ImportTeamGrades(w http.ResponseWriter, r *http.Request) {
	managers.AssignTeamGrades()
}

func ImportNFLDraftPicks(w http.ResponseWriter, r *http.Request) {
	managers.ImportNFLDraftPicks()
}

func ImportRecruitAICSV(w http.ResponseWriter, r *http.Request) {
	managers.ImportRecruitAICSV()
}

func ImportNFLRecords(w http.ResponseWriter, r *http.Request) {
	managers.RetireAndFreeAgentPlayers()
}

func ImportWorkEthic(w http.ResponseWriter, r *http.Request) {
	managers.ImportWorkEthic()
}

func ImportFAPreferences(w http.ResponseWriter, r *http.Request) {
	managers.ImportFAPreferences()
}

func ImportSimNFLMinimumValues(w http.ResponseWriter, r *http.Request) {
	managers.ImportMinimumFAValues()
}

func ImportTradePreferences(w http.ResponseWriter, r *http.Request) {
	managers.ImportTradePreferences()
}

func Import2023DraftedPlayers(w http.ResponseWriter, r *http.Request) {
	managers.Import2023DraftedPlayers()
}

func ImportCFBStandings(w http.ResponseWriter, r *http.Request) {
	managers.ImportSeasonStandings()
}

func ImportCFBGames(w http.ResponseWriter, r *http.Request) {
	managers.ImportCFBGames()
}

func ImportNFLGames(w http.ResponseWriter, r *http.Request) {
	managers.ImportNFLGames()
}

func ImportCFBTeams(w http.ResponseWriter, r *http.Request) {
	managers.ImportCFBTeams()
}

func ImportUDFAs(w http.ResponseWriter, r *http.Request) {
	managers.ImportUDFAs()
}

func UpdateDraftPicksForDraft(w http.ResponseWriter, r *http.Request) {
	managers.UpdateDraftPicks()
}

func ImplementPrimeAge(w http.ResponseWriter, r *http.Request) {
	managers.ImplementPrimeAge()
}

func GenerateCoachesForAITeams(w http.ResponseWriter, r *http.Request) {
	managers.GenerateCoachesForAITeams()
}

func FixCollegeDTOVRs(w http.ResponseWriter, r *http.Request) {
	managers.FixCollegeDTs()
	json.NewEncoder(w).Encode("Players Fixed")
}

func FixSpendingCount(w http.ResponseWriter, r *http.Request) {
	managers.FixSpendingCount()
	json.NewEncoder(w).Encode("Players Fixed")
}

func ImportCFBRivals(w http.ResponseWriter, r *http.Request) {
	managers.ImportCFBRivals()
}

func Import2021CFBStats(w http.ResponseWriter, r *http.Request) {
	// managers.MigrateRetiredAndNFLPlayersToHistoricCFBTable()
	managers.ImportCFB2021PlayerStats()
}

func FixATHProgressions(w http.ResponseWriter, r *http.Request) {
	managers.FixATHProgressions()
}

func FixSecondaryPositions(w http.ResponseWriter, r *http.Request) {
	managers.FixSecondaryPositions()
}

func ImportAdditionalDCPositions(w http.ResponseWriter, r *http.Request) {
	// managers.MigrateRetiredAndNFLPlayersToHistoricCFBTable()
	managers.ImportNewDepthChartPositionRecords()
}
