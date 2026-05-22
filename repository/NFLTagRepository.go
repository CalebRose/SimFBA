package repository

import (
	"log"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/structs"
	"gorm.io/gorm"
)

func FindAllNFLTagData() []structs.NFLTagData {
	var tagData []structs.NFLTagData
	db := dbprovider.GetInstance().GetDB()
	if err := db.Find(&tagData).Error; err != nil {
		log.Printf("Error loading NFLTagData: %v", err)
	}
	return tagData
}

func FindNFLTagDataByPosition(position string) structs.NFLTagData {
	var tagData structs.NFLTagData
	db := dbprovider.GetInstance().GetDB()
	db.Where("position = ?", position).First(&tagData)
	return tagData
}

func SaveNFLTagData(tagData structs.NFLTagData, db *gorm.DB) {
	if err := db.Save(&tagData).Error; err != nil {
		log.Printf("Error saving NFLTagData for position %s: %v", tagData.Position, err)
	}
}

func FindAllNFLPlayerSeasonSnaps() []structs.NFLPlayerSeasonSnaps {
	var snaps []structs.NFLPlayerSeasonSnaps
	db := dbprovider.GetInstance().GetDB()
	if err := db.Find(&snaps).Error; err != nil {
		log.Printf("Error loading NFLPlayerSeasonSnaps: %v", err)
	}
	return snaps
}
