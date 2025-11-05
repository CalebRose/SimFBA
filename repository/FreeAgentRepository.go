package repository

import (
	"log"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/structs"
	"gorm.io/gorm"
)

type FreeAgencyQuery struct {
	PlayerID string
	IsActive bool
	TeamID   string
	OfferID  string
}

func FindAllFreeAgentOffers(clauses FreeAgencyQuery) []structs.FreeAgencyOffer {
	db := dbprovider.GetInstance().GetDB()

	offers := []structs.FreeAgencyOffer{}
	query := db.Model(&offers)

	if len(clauses.TeamID) > 0 {
		query = query.Where("team_id = ?", clauses.TeamID)
	}

	if len(clauses.PlayerID) > 0 {
		query = query.Where("player_id = ?", clauses.PlayerID)
	}

	if len(clauses.OfferID) > 0 {
		query = query.Where("id = ?", clauses.OfferID)
	}

	if clauses.IsActive {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Find(&offers).Error; err != nil {
		return []structs.FreeAgencyOffer{}
	}

	return offers
}

func FindAllWaiverOffers(clauses FreeAgencyQuery) []structs.NFLWaiverOffer {
	db := dbprovider.GetInstance().GetDB()

	offers := []structs.NFLWaiverOffer{}
	query := db.Model(&offers)

	if len(clauses.TeamID) > 0 {
		query = query.Where("team_id = ?", clauses.TeamID)
	}

	if len(clauses.PlayerID) > 0 {
		query = query.Where("player_id = ?", clauses.PlayerID)
	}

	if len(clauses.OfferID) > 0 {
		query = query.Where("id = ?", clauses.OfferID)
	}

	if clauses.IsActive {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Find(&offers).Error; err != nil {
		return []structs.NFLWaiverOffer{}
	}

	return offers
}

func SaveFreeAgencyOfferRecord(contract structs.FreeAgencyOffer, db *gorm.DB) {
	err := db.Save(&contract).Error
	if err != nil {
		log.Panicln("Could not save offer record")
	}
}

func MassDeleteFreeAgencyOffers(offerIDs []string, db *gorm.DB) {
	err := db.Unscoped().Where("id IN ?", offerIDs).Delete(&structs.FreeAgencyOffer{}).Error
	if err != nil {
		log.Panicln("Could not delete free agency offers")
	}
}
