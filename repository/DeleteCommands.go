package repository

import (
	"log"

	"github.com/CalebRose/SimFBA/structs"
	"gorm.io/gorm"
)

func DeleteCollegePlayerRecord(player structs.CollegePlayer, db *gorm.DB) {
	err := db.Delete(&player).Error
	if err != nil {
		log.Panicln("Could not delete old college player record.")
	}
}

func DeleteCollegeRecruitRecord(player structs.Recruit, db *gorm.DB) {
	err := db.Delete(&player).Error
	if err != nil {
		log.Panicln("Could not delete old college recruit record.")
	}
}

func DeleteTransferPortalProfile(profile structs.TransferPortalProfile, db *gorm.DB) {
	profile.CollegePlayer = structs.CollegePlayer{}
	profile.Promise = structs.CollegePromise{}
	err := db.Delete(&profile).Error
	if err != nil {
		log.Panicln("Could not delete old college player record.")
	}
}

func DeleteCollegePromise(promise structs.CollegePromise, db *gorm.DB) {
	err := db.Delete(&promise).Error
	if err != nil {
		log.Panicln("Could not delete old college player record.")
	}
}

func DeleteNotificationRecord(noti structs.Notification, db *gorm.DB) {
	err := db.Delete(&noti).Error
	if err != nil {
		log.Panicln("Could not delete old notification record.")
	}
}
