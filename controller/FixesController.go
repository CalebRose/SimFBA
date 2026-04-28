package controller

import (
	"net/http"

	"github.com/CalebRose/SimFBA/managers"
)

func FixCollegeStandings(w http.ResponseWriter, r *http.Request) {
	managers.ResetCollegeStandings()
}

func FixRecruitPoints(w http.ResponseWriter, r *http.Request) {
	managers.FixRecruitingProfiles()
}

func FixOffensiveFormationNames(w http.ResponseWriter, r *http.Request) {
	managers.FixOffensiveFormationNames()
}

func FixNFLDrafteePotentialGrades(w http.ResponseWriter, r *http.Request) {
	managers.FixNFLDrafteePotentialGrades()
}

func FixPlayerPreferences(w http.ResponseWriter, r *http.Request) {
	managers.FixPreferencesForAllRecruitsAndCollegePlayers()
}
