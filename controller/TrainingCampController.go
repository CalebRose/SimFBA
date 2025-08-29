package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimFBA/managers"
)

// GetAllCollegeTeamsForRosterPage
func RunTrainingCamps(w http.ResponseWriter, r *http.Request) {
	managers.RunTrainingCamps()

	json.NewEncoder(w).Encode("Training Camp Complete")
}
