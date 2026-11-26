package controller

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"

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
	ip := getClientIP(r)
	log.Printf("Request IP address: %s", ip)
	collegeID := vars["collegeID"]
	proID := vars["proID"]
	data := managers.GetLandingBootstrap(collegeID, proID, ip)
	bootstrapData, err := easyjson.Marshal(data)
	if err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	w.Write(bootstrapData)
}

func getClientIP(r *http.Request) string {
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		candidate := strings.TrimSpace(strings.SplitN(forwardedFor, ",", 2)[0])
		if ip := parseIPAddress(candidate); ip != "" {
			return ip
		}
		log.Printf("Invalid client IP in X-Forwarded-For: %q", candidate)
	}

	if ip := parseIPAddress(r.RemoteAddr); ip != "" {
		return ip
	}
	log.Printf("Failed to get client IP from RemoteAddr %q", r.RemoteAddr)
	return ""
}

func parseIPAddress(value string) string {
	if net.ParseIP(value) != nil {
		return value
	}
	host, _, err := net.SplitHostPort(value)
	if err == nil && net.ParseIP(host) != nil {
		return host
	}
	return ""
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
	proID := vars["proID"]
	data := managers.GetSchedulePageBootstrap(username, collegeID, proID, seasonID)
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
	data := managers.GetGameplanBootstrap(collegeID, proID)
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

func BootstrapStatsData(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	vars := mux.Vars(r)
	collegeID := vars["collegeID"]
	proID := vars["proID"]
	data := managers.GetStatsBootstrap(collegeID, proID)
	bootstrapData, err := easyjson.Marshal(data)
	if err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	w.Write(bootstrapData)
}

func BootstrapPlayerData(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	vars := mux.Vars(r)
	collegeID := vars["collegeID"]
	proID := vars["proID"]
	data := managers.GetPlayerBootstrap(collegeID, proID)
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
