package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CalebRose/SimFBA/managers"
	"github.com/CalebRose/SimFBA/models"
	"github.com/gorilla/mux"
)

func ExportDrafteesToCSV(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")

	managers.ExportDrafteesToCSV(w)
}

func GetDraftPageData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	if len(teamID) == 0 {
		panic("User did not provide TeamID")
	}
	// Get War Room
	// Get Scouting Profiles?
	// Get full list of draftable players

	warRoom := managers.GetNFLWarRoomByTeamID(teamID)
	draftees := managers.GetNFLDrafteesForDraftPage()
	allNFLTeams := managers.GetAllNFLTeams()
	draftPicks := managers.GetAllCurrentSeasonDraftPicksForDraftRoom()

	res := models.NFLDraftPageResponse{
		WarRoom:          warRoom,
		DraftablePlayers: draftees,
		NFLTeams:         allNFLTeams,
		AllDraftPicks:    draftPicks,
	}

	json.NewEncoder(w).Encode(res)
}

func AddPlayerToScoutBoard(w http.ResponseWriter, r *http.Request) {

	var scoutProfileDto models.ScoutingProfileDTO
	err := json.NewDecoder(r.Body).Decode(&scoutProfileDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	scoutingProfile := managers.CreateScoutingProfile(scoutProfileDto)

	json.NewEncoder(w).Encode(scoutingProfile)
}

func ExportDraftedPicks(w http.ResponseWriter, r *http.Request) {
	var draftPickDTO models.ExportDraftPicksDTO
	err := json.NewDecoder(r.Body).Decode(&draftPickDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	saveComplete := managers.ExportDraftedPlayers(draftPickDTO.DraftPicks)

	json.NewEncoder(w).Encode(saveComplete)

	fmt.Fprintf(w, "Exported Players to new tables")
}

func RevealScoutingAttribute(w http.ResponseWriter, r *http.Request) {
	var revealAttributeDTO models.RevealAttributeDTO
	err := json.NewDecoder(r.Body).Decode(&revealAttributeDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	saveComplete := managers.RevealScoutingAttribute(revealAttributeDTO)

	json.NewEncoder(w).Encode(saveComplete)

	fmt.Fprintf(w, "New Scout Profile Created")
}

func RemovePlayerFromScoutBoard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if len(id) == 0 {
		panic("User did not provide scout profile id")
	}

	managers.RemovePlayerFromScoutBoard(id)

	json.NewEncoder(w).Encode("Removed Player From Scout Board")
}

func GetScoutingDataByDraftee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if len(id) == 0 {
		panic("User did not provide scout profile id")
	}

	data := managers.GetScoutingDataByPlayerID(id)

	json.NewEncoder(w).Encode(data)
}

func ToggleDraftTime(w http.ResponseWriter, r *http.Request) {
	managers.ToggleDraftTime()

	json.NewEncoder(w).Encode("Draft Time Changed")
}

func BoomOrBust(w http.ResponseWriter, r *http.Request) {
	managers.BoomOrBust()

	json.NewEncoder(w).Encode("Boom/Bust Complete")
}
