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

func ImportTradePreferences(w http.ResponseWriter, r *http.Request) {
	managers.ImportTradePreferences()
}

func GetMissingRecruitingClasses(w http.ResponseWriter, r *http.Request) {
	managers.GetMissingRecruitingClasses()
}
