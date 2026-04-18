package repository

import (
	"log"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/structs"
	"gorm.io/gorm"
)

type TransferPortalQuery struct {
	ID               string
	TeamID           string
	ProfileID        string
	CollegePlayerID  string
	RemovedFromBoard string
	OrderByPoints    bool
	IsActive         string
}

func FindTransferPortalProfileRecord(clauses TransferPortalQuery) structs.TransferPortalProfile {
	db := dbprovider.GetInstance().GetDB()

	var profiles structs.TransferPortalProfile

	query := db.Model(&profiles)

	if len(clauses.ID) > 0 {
		query = query.Where("id = ?", clauses.ID)
	}

	if len(clauses.ProfileID) > 0 {
		query = query.Where("profile_id = ?", clauses.ProfileID)
	}

	if len(clauses.CollegePlayerID) > 0 {
		query = query.Where("college_player_id = ?", clauses.CollegePlayerID)
	}

	if len(clauses.RemovedFromBoard) > 0 {
		isRemoved := clauses.RemovedFromBoard == "Y"
		query = query.Where("removed_from_board = ?", isRemoved)
	}

	if clauses.OrderByPoints {
		query = query.Order("total_points desc")
	}

	if err := query.Find(&profiles).Error; err != nil {
		return structs.TransferPortalProfile{}
	}

	return profiles
}

func FindTransferPortalProfileRecords(clauses TransferPortalQuery) []structs.TransferPortalProfile {
	db := dbprovider.GetInstance().GetDB()

	var profiles []structs.TransferPortalProfile

	query := db.Model(&profiles)

	if len(clauses.ID) > 0 {
		query = query.Where("id = ?", clauses.ID)
	}

	if len(clauses.ProfileID) > 0 {
		query = query.Where("profile_id = ?", clauses.ProfileID)
	}

	if len(clauses.CollegePlayerID) > 0 {
		query = query.Where("college_player_id = ?", clauses.CollegePlayerID)
	}

	if len(clauses.RemovedFromBoard) > 0 {
		isRemoved := clauses.RemovedFromBoard == "Y"
		query = query.Where("removed_from_board = ?", isRemoved)
	}

	if clauses.OrderByPoints {
		query = query.Order("total_points desc")
	}

	if err := query.Find(&profiles).Error; err != nil {
		return []structs.TransferPortalProfile{}
	}

	return profiles
}

func FindCollegePromiseRecord(clauses TransferPortalQuery) structs.CollegePromise {
	db := dbprovider.GetInstance().GetDB()

	p := structs.CollegePromise{}

	query := db.Model(&p)

	if len(clauses.CollegePlayerID) > 0 {
		query = query.Where("college_player_id = ?", clauses.CollegePlayerID)
	}

	if len(clauses.TeamID) > 0 {
		query = query.Where("team_id = ?", clauses.TeamID)
	}

	if err := query.Find(&p).Error; err != nil {
		return structs.CollegePromise{}
	}
	return p
}

func FindCollegePromiseRecords(clauses TransferPortalQuery) []structs.CollegePromise {
	db := dbprovider.GetInstance().GetDB()

	p := []structs.CollegePromise{}

	query := db.Model(&p)

	if len(clauses.CollegePlayerID) > 0 {
		query = query.Where("college_player_id = ?", clauses.ID)
	}

	if len(clauses.TeamID) > 0 {
		query = query.Where("team_id = ?", clauses.TeamID)
	}

	if len(clauses.IsActive) > 0 {
		isActive := clauses.IsActive == "Y"
		query = query.Where("is_active = ?", isActive)
	}

	if err := query.Find(&p).Error; err != nil {
		return []structs.CollegePromise{}
	}
	return p
}

// Saves
func SaveTransferPortalProfile(profile structs.TransferPortalProfile, db *gorm.DB) {

	err := db.Save(&profile).Error
	if err != nil {
		log.Panicln("Could not save player record")
	}
}
