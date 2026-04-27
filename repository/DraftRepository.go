package repository

import (
	"github.com/CalebRose/SimFBA/structs"
	"gorm.io/gorm"
)

func CreateNFLDraftPickRecord(pick structs.NFLDraftPick, db *gorm.DB) error {
	if pick.ID == 0 {
		err := db.Create(&pick).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func SaveNFLDraftPickRecord(pick structs.NFLDraftPick, db *gorm.DB) error {
	if pick.ID == 0 {
		return nil
	} else {
		err := db.Save(&pick).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateNFLDraftPickBatch(db *gorm.DB, picks []structs.NFLDraftPick, batchSize int) error {
	total := len(picks)
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}

		if err := db.CreateInBatches(picks[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}
