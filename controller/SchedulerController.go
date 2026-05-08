package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimFBA/managers"
)

func TestCFBScheduler(w http.ResponseWriter, r *http.Request) {
	managers.BaseGenerateCFBSchedule(true)
	json.NewEncoder(w).Encode(true)
}
