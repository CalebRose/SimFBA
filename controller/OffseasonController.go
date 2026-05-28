package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimFBA/managers"
)

func FixPostseasonStatus(w http.ResponseWriter, r *http.Request) {
	managers.PostSeasonStatusCleanUp()
	json.NewEncoder(w).Encode("Done!")
}

func UpdateTeamProfileAffinities(w http.ResponseWriter, r *http.Request) {
	managers.UpdateTeamProfileAffinities()
	json.NewEncoder(w).Encode("Done!")
}

func CalculateExtensionValues(w http.ResponseWriter, r *http.Request) {
	managers.CalculatePlayerMinimumAndAAVValues()
	json.NewEncoder(w).Encode("Done!")
}

func CalculateTagValues(w http.ResponseWriter, r *http.Request) {
	managers.CalculateTagValues()
	json.NewEncoder(w).Encode("Done!")
}

func ResetNFLMinimumValues(w http.ResponseWriter, r *http.Request) {
	managers.ResetNFLPlayerMinimumValues()
	json.NewEncoder(w).Encode("Done!")
}

func SeedTagDataFromJSON(w http.ResponseWriter, r *http.Request) {
	managers.SeedNFLTagDataFromJSON()
	json.NewEncoder(w).Encode("Done!")
}
