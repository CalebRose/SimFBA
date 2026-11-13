package repository

import (
	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/structs"
	"gorm.io/gorm"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func FindCFBRealignments() []structs.CFBRealignment {
	db := dbprovider.GetInstance().GetDB()
	var realignments []structs.CFBRealignment
	db.Where("is_processed = ?", false).Find(&realignments)
	return realignments
}

func FindNFLRealignments() []structs.NFLRealignment {
	db := dbprovider.GetInstance().GetDB()
	var realignments []structs.NFLRealignment
	db.Where("is_processed = ?", false).Find(&realignments)
	return realignments
}

func SaveCFBRealignmentRecord(realignment structs.CFBRealignment, db *gorm.DB) error {
	if err := db.Save(&realignment).Error; err != nil {
		return err
	}
	return nil
}

func SaveNFLRealignmentRecord(realignment structs.NFLRealignment, db *gorm.DB) error {
	if err := db.Save(&realignment).Error; err != nil {
		return err
	}
	return nil
}

func CreateCFBRealignmentBatch(realignments []structs.CFBRealignment, db *gorm.DB, batchSize int) error {
	total := len(realignments)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(realignments[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateNFLRealignmentBatch(realignments []structs.NFLRealignment, db *gorm.DB, batchSize int) error {
	total := len(realignments)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(realignments[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}
