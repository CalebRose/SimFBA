package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/CalebRose/SimFBA/managers"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/gorilla/mux"
)

// GetDepthChartSnapshots handles GET /api/depthcharts/snapshots/{league}/season/{seasonID}
func GetDepthChartSnapshots(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	league := strings.ToLower(vars["league"])
	seasonIDStr := vars["seasonID"]

	if len(seasonIDStr) == 0 {
		http.Error(w, "seasonID parameter is required", http.StatusBadRequest)
		return
	}

	seasonID, err := strconv.ParseUint(seasonIDStr, 10, 32)
	if err != nil {
		http.Error(w, "invalid seasonID", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if league == "nfl" || league == "pro" {
		snapshots := managers.GetNFLDepthChartSnapshots(uint(seasonID))
		json.NewEncoder(w).Encode(snapshots)
	} else {
		snapshots := managers.GetCollegeDepthChartSnapshots(uint(seasonID))
		json.NewEncoder(w).Encode(snapshots)
	}
}

// CaptureDepthChartSnapshot handles POST /api/depthcharts/snapshots/{league}/season/{seasonID}
// It snapshots the live depth chart positions table into the snapshot table with the season column.
func CaptureDepthChartSnapshot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	league := strings.ToLower(vars["league"])
	seasonIDStr := vars["seasonID"]

	if len(seasonIDStr) == 0 {
		http.Error(w, "seasonID parameter is required", http.StatusBadRequest)
		return
	}

	seasonID, err := strconv.ParseUint(seasonIDStr, 10, 32)
	if err != nil {
		http.Error(w, "invalid seasonID", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if league == "nfl" || league == "pro" {
		records, err := managers.CaptureNFLDepthChartSnapshot(uint(seasonID))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(records)
	} else {
		records, err := managers.CaptureCollegeDepthChartSnapshot(uint(seasonID))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(records)
	}
}

// UploadDepthChartSnapshot handles PUT /api/depthcharts/snapshots/{league}/season/{seasonID}
func UploadDepthChartSnapshot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	league := strings.ToLower(vars["league"])
	seasonIDStr := vars["seasonID"]

	if len(seasonIDStr) == 0 {
		http.Error(w, "seasonID parameter is required", http.StatusBadRequest)
		return
	}

	seasonID, err := strconv.ParseUint(seasonIDStr, 10, 32)
	if err != nil {
		http.Error(w, "invalid seasonID", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if league == "nfl" || league == "pro" {
		var records []structs.NFLDepthChartSeasonSnapshot
		err := json.NewDecoder(r.Body).Decode(&records)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for i := range records {
			records[i].SeasonID = uint(seasonID)
		}
		err = managers.SaveNFLDepthChartSnapshots(uint(seasonID), records)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(records)
	} else {
		var records []structs.CollegeDepthChartSeasonSnapshot
		err := json.NewDecoder(r.Body).Decode(&records)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for i := range records {
			records[i].SeasonID = uint(seasonID)
		}
		err = managers.SaveCollegeDepthChartSnapshots(uint(seasonID), records)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(records)
	}
}
