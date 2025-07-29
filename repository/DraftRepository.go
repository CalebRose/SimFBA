package repository

import (
	"github.com/CalebRose/SimFBA/structs"
	"gorm.io/gorm"
)

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
