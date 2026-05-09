package repository

import (
	"log"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/structs"
	"gorm.io/gorm"
)

type SchedulerQuery struct {
	ID             string
	TeamID         string
	IsSpringGames  bool
	IsNFLPreseason bool
	IsAccepted     bool
	SeasonID       string
	WeekID         string
}

func FindCFBGameRequestRecord(scheduleParams SchedulerQuery) structs.CFBGameRequest {
	db := dbprovider.GetInstance().GetDB()
	var request structs.CFBGameRequest

	query := db.Model(&request)

	if len(scheduleParams.ID) > 0 {
		query = query.Where("id = ?", scheduleParams.ID)
	}

	if len(scheduleParams.TeamID) > 0 {
		query = query.Where("home_team_id = ? OR away_team_id = ?", scheduleParams.TeamID, scheduleParams.TeamID)
	}
	if scheduleParams.IsSpringGames {
		query = query.Where("is_spring_game = ?", true)
	}
	if err := query.Find(&request).Error; err != nil {
		return structs.CFBGameRequest{}
	}
	return request
}

func FindCFBGameRequestRecords(scheduleParams SchedulerQuery) []structs.CFBGameRequest {
	db := dbprovider.GetInstance().GetDB()
	var requests []structs.CFBGameRequest

	query := db.Model(&requests)

	if len(scheduleParams.TeamID) > 0 {
		query = query.Where("home_team_id = ? OR away_team_id = ?", scheduleParams.TeamID, scheduleParams.TeamID)
	}
	if scheduleParams.IsAccepted {
		query = query.Where("is_accepted = ?", true)
	}
	if len(scheduleParams.SeasonID) > 0 {
		query = query.Where("season_id = ?", scheduleParams.SeasonID)
	}
	if len(scheduleParams.WeekID) > 0 {
		query = query.Where("week_id = ?", scheduleParams.WeekID)
	}
	if scheduleParams.IsSpringGames {
		query = query.Where("is_spring_game = ?", true)
	}

	if err := query.Find(&requests).Error; err != nil {
		return []structs.CFBGameRequest{}
	}
	return requests
}

func FindNFLGameRequestRecord(scheduleParams SchedulerQuery) structs.NFLGameRequest {
	db := dbprovider.GetInstance().GetDB()
	var request structs.NFLGameRequest

	query := db.Model(&request)

	if len(scheduleParams.ID) > 0 {
		query = query.Where("id = ?", scheduleParams.ID)
	}

	if len(scheduleParams.TeamID) > 0 {
		query = query.Where("home_team_id = ? OR away_team_id = ?", scheduleParams.TeamID, scheduleParams.TeamID)
	}
	if scheduleParams.IsNFLPreseason {
		query = query.Where("is_nfl_preseason = ?", true)
	}
	if err := query.Find(&request).Error; err != nil {
		return structs.NFLGameRequest{}
	}
	return request
}

func FindNFLGameRequestRecords(scheduleParams SchedulerQuery) []structs.NFLGameRequest {
	db := dbprovider.GetInstance().GetDB()
	var requests []structs.NFLGameRequest

	query := db.Model(&requests)

	if len(scheduleParams.TeamID) > 0 {
		query = query.Where("home_team_id = ? OR away_team_id = ?", scheduleParams.TeamID, scheduleParams.TeamID)
	}

	// Return accepted but not yet approved requests for processing
	if scheduleParams.IsAccepted {
		query = query.Where("is_accepted = ? AND is_approved = ?", true, false)
	}
	if len(scheduleParams.SeasonID) > 0 {
		query = query.Where("season_id = ?", scheduleParams.SeasonID)
	}
	if len(scheduleParams.WeekID) > 0 {
		query = query.Where("week_id = ?", scheduleParams.WeekID)
	}
	if scheduleParams.IsNFLPreseason {
		query = query.Where("is_nfl_preseason = ?", true)
	}

	if err := query.Find(&requests).Error; err != nil {
		return []structs.NFLGameRequest{}
	}
	return requests
}

func CreateCFBGameRequest(request structs.CFBGameRequest, db *gorm.DB) {
	err := db.Create(&request).Error
	if err != nil {
		log.Panicln("Could not create CFB game request record!")
	}
}

func SaveCFBGameRequest(request structs.CFBGameRequest, db *gorm.DB) error {
	err := db.Save(&request).Error
	if err != nil {
		return err
	}
	return nil
}

func CreateNFLGameRequest(request structs.NFLGameRequest, db *gorm.DB) {
	err := db.Create(&request).Error
	if err != nil {
		log.Panicln("Could not create NFL game request record!")
	}
}

func SaveNFLGameRequest(request structs.NFLGameRequest, db *gorm.DB) error {
	err := db.Save(&request).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteCFBGameRequest(request structs.CFBGameRequest, db *gorm.DB) error {
	err := db.Delete(&request).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteNFLGameRequest(request structs.NFLGameRequest, db *gorm.DB) error {
	err := db.Delete(&request).Error
	if err != nil {
		return err
	}
	return nil
}
