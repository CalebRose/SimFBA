package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/structs"
)

// Collusion Button
func CollusionButton(w http.ResponseWriter, r *http.Request) {
	db := dbprovider.GetInstance().GetDB()
	var collusionButton structs.CollusionDto
	err := json.NewDecoder(r.Body).Decode(&collusionButton)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newsLog := structs.NewsLog{
		WeekID:      collusionButton.WeekID,
		SeasonID:    collusionButton.SeasonID,
		MessageType: "Collusion",
		Message:     collusionButton.Message,
	}

	db.Create(&newsLog)
}