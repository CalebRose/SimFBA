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
