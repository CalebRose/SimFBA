package controller

import (
	"net/http"

	"github.com/CalebRose/SimFBA/managers"
)

func ImportCustomCroots(w http.ResponseWriter, r *http.Request) {
	managers.CreateCustomCroots()
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

func ImportCFBGames(w http.ResponseWriter, r *http.Request) {
	managers.ImportCFBGames()
}

func ImportNFLGames(w http.ResponseWriter, r *http.Request) {
	managers.ImportNFLGames()
}

func ImportUDFAs(w http.ResponseWriter, r *http.Request) {
	managers.ImportUDFAs()
}

func GetMissingRecruitingClasses(w http.ResponseWriter, r *http.Request) {
	managers.GetMissingRecruitingClasses()
}

func FixDTOveralls(w http.ResponseWriter, r *http.Request) {
	managers.FixOverallRatingDT()
}
