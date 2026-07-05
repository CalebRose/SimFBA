package repository

import (
	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/structs"
)

type PlayerQuery struct {
	ID             string
	TeamID         string
	CollegeID      string
	PlayerIDs      []string
	TransferStatus string
	LeagueID       string
	IsInjured      string
	IsFreeAgent    string
	OverallDesc    bool
}

// College Players
func FindAllCollegePlayers(clauses PlayerQuery) []structs.CollegePlayer {
	db := dbprovider.GetInstance().GetDB()

	var CollegePlayers []structs.CollegePlayer

	query := db.Model(&CollegePlayers)

	if len(clauses.TeamID) > 0 {
		query = query.Where("team_id = ?", clauses.TeamID)
	}

	if len(clauses.PlayerIDs) > 0 {
		query = query.Where("id in (?)", clauses.PlayerIDs)
	}

	if len(clauses.TransferStatus) > 0 {
		query = query.Where("transfer_status = ?", clauses.TransferStatus)
	}

	if len(clauses.LeagueID) > 0 {
		query = query.Where("league_id = ?", clauses.LeagueID)
	}

	if len(clauses.IsInjured) > 0 {
		query = query.Where("is_injured = ?", true)
	}

	if clauses.OverallDesc {
		query = query.Order("overall desc")
	}

	if err := query.Find(&CollegePlayers).Error; err != nil {
		return []structs.CollegePlayer{}
	}

	return CollegePlayers
}
