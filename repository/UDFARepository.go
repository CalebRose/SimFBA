package repository

import (
	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/structs"
)

// GetUDFABoardByTeamID fetches the board and all associated bid profiles
func GetUDFABoardByTeamID(teamID string) structs.NFLUDFABoard {
	db := dbprovider.GetInstance().GetDB()
	var board structs.NFLUDFABoard
	db.Preload("Profiles").Where("team_id = ?", teamID).First(&board)
	return board
}

// SaveUDFAProfile saves or updates a specific bid
func SaveUDFAProfile(profile structs.NFLUDFAProfile) error {
	db := dbprovider.GetInstance().GetDB()
	return db.Save(&profile).Error
}

// DeleteUDFAProfile removes a bid from the board
func DeleteUDFAProfile(profileID string) error {
	db := dbprovider.GetInstance().GetDB()
	return db.Where("id = ?", profileID).Delete(&structs.NFLUDFAProfile{}).Error
}