package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/CalebRose/SimFBA/managers"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
)

func BootstrapTeamData(w http.ResponseWriter, r *http.Request) {
	data := managers.GetTeamsBootstrap()
	w.Header().Set("Content-Type", "application/json")
	teamData, err := easyjson.Marshal(data)
	if err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	w.Write(teamData)
}

func BootstrapLandingData(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	vars := mux.Vars(r)
	collegeID := vars["collegeID"]
	proID := vars["proID"]
	data := managers.GetLandingBootstrap(collegeID, proID)
	bootstrapData, err := easyjson.Marshal(data)
	if err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	w.Write(bootstrapData)
}

func BootstrapTeamRosterData(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	vars := mux.Vars(r)
	collegeID := vars["collegeID"]
	proID := vars["proID"]
	data := managers.GetTeamRosterBootstrap(collegeID, proID)
	bootstrapData, err := easyjson.Marshal(data)
	if err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	w.Write(bootstrapData)
}

func BootstrapRecruitingData(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	vars := mux.Vars(r)
	collegeID := vars["collegeID"]
	data := managers.GetRecruitingBootstrap(collegeID)
	bootstrapData, err := easyjson.Marshal(data)
	if err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	w.Write(bootstrapData)
}

func BootstrapFreeAgencyData(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	vars := mux.Vars(r)
	proID := vars["proID"]
	data := managers.GetFreeAgencyBootstrap(proID)
	bootstrapData, err := easyjson.Marshal(data)
	if err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	w.Write(bootstrapData)
}

func BootstrapSchedulingData(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	vars := mux.Vars(r)
	username := vars["username"]
	collegeID := vars["collegeID"]
	seasonID := vars["seasonID"]
	data := managers.GetCollegePollsBootstrap(username, collegeID, seasonID)
	bootstrapData, err := easyjson.Marshal(data)
	if err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	w.Write(bootstrapData)
}

func BootstrapDraftData(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	vars := mux.Vars(r)
	proID := vars["proID"]
	data := managers.GetDraftBootstrap(proID)
	bootstrapData, err := easyjson.Marshal(data)
	if err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	w.Write(bootstrapData)
}

func BootstrapPortalData(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	vars := mux.Vars(r)
	collegeID := vars["collegeID"]
	data := managers.GetPortalBootstrap(collegeID)
	bootstrapData, err := easyjson.Marshal(data)
	if err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	w.Write(bootstrapData)
}

func BootstrapGameplanData(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	vars := mux.Vars(r)
	collegeID := vars["collegeID"]
	proID := vars["proID"]
	data := managers.GetTeamRosterBootstrap(collegeID, proID)
	bootstrapData, err := easyjson.Marshal(data)
	if err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	w.Write(bootstrapData)
}

func BootstrapNewsData(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	vars := mux.Vars(r)
	collegeID := vars["collegeID"]
	proID := vars["proID"]
	data := managers.GetNewsBootstrap(collegeID, proID)
	bootstrapData, err := easyjson.Marshal(data)
	if err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	w.Write(bootstrapData)
}

func GetCollegeHistoryProfile(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	data := managers.GetCollegeTeamProfilePageData()
	json.NewEncoder(w).Encode(data)
}
