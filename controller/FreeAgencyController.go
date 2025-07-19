package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimFBA/managers"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/gorilla/mux"
)

// FreeAgencyAvailablePlayers - Get All Available NFL Players for Free Agency Page
func FreeAgencyAvailablePlayers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamId := vars["teamID"]

	var players = managers.GetAllAvailableNFLPlayers(teamId)

	json.NewEncoder(w).Encode(players)
}

// FreeAgencyAvailablePlayers - Get All Available NFL Players for Free Agency Page
func CreateFreeAgencyOffer(w http.ResponseWriter, r *http.Request) {
	var freeAgencyOfferDTO structs.FreeAgencyOfferDTO
	err := json.NewDecoder(r.Body).Decode(&freeAgencyOfferDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var offer = managers.CreateFAOffer(freeAgencyOfferDTO)

	json.NewEncoder(w).Encode(offer)
}

// FreeAgencyAvailablePlayers - Get All Available NFL Players for Free Agency Page
func CancelFreeAgencyOffer(w http.ResponseWriter, r *http.Request) {
	var freeAgencyOfferDTO structs.FreeAgencyOfferDTO
	err := json.NewDecoder(r.Body).Decode(&freeAgencyOfferDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.CancelOffer(freeAgencyOfferDTO)

	json.NewEncoder(w).Encode(true)
}

func SetWaiverOrderForNFLTeams(w http.ResponseWriter, r *http.Request) {
	managers.SetWaiverOrder()
	json.NewEncoder(w).Encode("Waiver Order Set")
}

// CreateWaiverWireOffer
func CreateWaiverWireOffer(w http.ResponseWriter, r *http.Request) {
	var waiverWireOfferDTO structs.NFLWaiverOffDTO
	err := json.NewDecoder(r.Body).Decode(&waiverWireOfferDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var offer = managers.CreateWaiverOffer(waiverWireOfferDTO)

	json.NewEncoder(w).Encode(offer)
}

// CancelWaiverWireOffer
func CancelWaiverWireOffer(w http.ResponseWriter, r *http.Request) {
	var waiverWireOfferDTO structs.NFLWaiverOffDTO
	err := json.NewDecoder(r.Body).Decode(&waiverWireOfferDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.CancelWaiverOffer(waiverWireOfferDTO)

	json.NewEncoder(w).Encode(true)
}

func TagPlayer(w http.ResponseWriter, r *http.Request) {
	var tagDTO structs.NFLTagDTO
	err := json.NewDecoder(r.Body).Decode(&tagDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.TagPlayer(tagDTO)

	json.NewEncoder(w).Encode(true)
}

func SyncExtensions(w http.ResponseWriter, r *http.Request) {
	managers.SyncExtensionOffers()
	json.NewEncoder(w).Encode(true)
}
