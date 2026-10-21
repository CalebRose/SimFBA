package repository

import (
	"log"

	"github.com/CalebRose/SimFBA/structs"
	"gorm.io/gorm"
)

func GetAllCollegeDepthChartPositions(db *gorm.DB) []structs.CollegeDepthChartPosition {
	var positions []structs.CollegeDepthChartPosition
	err := db.Find(&positions).Error
	if err != nil {
		log.Printf("Error fetching college depth chart positions: %v\n", err)
		return []structs.CollegeDepthChartPosition{}
	}
	return positions
}

func GetAllNFLDepthChartPositions(db *gorm.DB) []structs.NFLDepthChartPosition {
	var positions []structs.NFLDepthChartPosition
	err := db.Find(&positions).Error
	if err != nil {
		log.Printf("Error fetching NFL depth chart positions: %v\n", err)
		return []structs.NFLDepthChartPosition{}
	}
	return positions
}

func GetCollegeDepthChartSnapshots(seasonID uint, db *gorm.DB) []structs.CollegeDepthChartSeasonSnapshot {
	var records []structs.CollegeDepthChartSeasonSnapshot
	err := db.Where("season_id = ?", seasonID).Find(&records).Error
	if err != nil {
		log.Printf("Error fetching college depth chart snapshots for season %d: %v\n", seasonID, err)
		return []structs.CollegeDepthChartSeasonSnapshot{}
	}
	return records
}

func GetNFLDepthChartSnapshots(seasonID uint, db *gorm.DB) []structs.NFLDepthChartSeasonSnapshot {
	var records []structs.NFLDepthChartSeasonSnapshot
	err := db.Where("season_id = ?", seasonID).Find(&records).Error
	if err != nil {
		log.Printf("Error fetching NFL depth chart snapshots for season %d: %v\n", seasonID, err)
		return []structs.NFLDepthChartSeasonSnapshot{}
	}
	return records
}

func SaveCollegeDepthChartSnapshots(seasonID uint, records []structs.CollegeDepthChartSeasonSnapshot, db *gorm.DB) error {
	err := db.Where("season_id = ?", seasonID).Delete(&structs.CollegeDepthChartSeasonSnapshot{}).Error
	if err != nil {
		log.Printf("Error clearing previous college depth chart snapshots for season %d: %v\n", seasonID, err)
		return err
	}

	if len(records) > 0 {
		err = db.CreateInBatches(&records, 500).Error
		if err != nil {
			log.Printf("Error batch inserting college depth chart snapshots: %v\n", err)
			return err
		}
	}
	return nil
}

func SaveNFLDepthChartSnapshots(seasonID uint, records []structs.NFLDepthChartSeasonSnapshot, db *gorm.DB) error {
	err := db.Where("season_id = ?", seasonID).Delete(&structs.NFLDepthChartSeasonSnapshot{}).Error
	if err != nil {
		log.Printf("Error clearing previous NFL depth chart snapshots for season %d: %v\n", seasonID, err)
		return err
	}

	if len(records) > 0 {
		err = db.CreateInBatches(&records, 500).Error
		if err != nil {
			log.Printf("Error batch inserting NFL depth chart snapshots: %v\n", err)
			return err
		}
	}
	return nil
}
