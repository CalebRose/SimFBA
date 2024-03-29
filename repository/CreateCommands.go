package repository

import (
	"log"

	"github.com/CalebRose/SimFBA/structs"
	"gorm.io/gorm"
)

func CreateCFBPlayByPlaysInBatch(plays []structs.CollegePlayByPlay, db *gorm.DB) {
	err := db.CreateInBatches(&plays, len(plays)).Error
	if err != nil {
		log.Panicln("Could not save play by plays!")
	}
}

func CreateNFLPlayByPlaysInBatch(plays []structs.NFLPlayByPlay, db *gorm.DB) {
	err := db.CreateInBatches(&plays, len(plays)).Error
	if err != nil {
		log.Panicln("Could not save play by plays!")
	}
}

func CreateCFBTeamStatsInBatch(stats []structs.CollegeTeamStats, db *gorm.DB) {
	err := db.CreateInBatches(&stats, len(stats)).Error
	if err != nil {
		log.Panicln("Could not save team stats!")
	}
}

func CreateCFBPlayerStatsInBatch(stats []structs.CollegePlayerStats, db *gorm.DB) {
	err := db.CreateInBatches(&stats, len(stats)).Error
	if err != nil {
		log.Panicln("Could not save team stats!")
	}
}

func CreateNFLTeamStatsInBatch(stats []structs.NFLTeamStats, db *gorm.DB) {
	err := db.CreateInBatches(&stats, len(stats)).Error
	if err != nil {
		log.Panicln("Could not save team stats!")
	}
}

func CreateNFLPlayerStatsInBatch(stats []structs.NFLPlayerStats, db *gorm.DB) {
	err := db.CreateInBatches(&stats, len(stats)).Error
	if err != nil {
		log.Panicln("Could not save team stats!")
	}
}

func CreateRecruitProfileRecord(cp structs.RecruitPlayerProfile, db *gorm.DB) {
	err := db.Create(&cp).Error
	if err != nil {
		log.Panicln("Could not create recruit profile record!")
	}
}
