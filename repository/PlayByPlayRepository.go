package repository

import (
	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/structs"
)

func FindCFBPlayByPlaysRecordsByGameID(id string) []structs.CollegePlayByPlay {
	db := dbprovider.GetInstance().GetDB()

	plays := []structs.CollegePlayByPlay{}

	db.Where("game_id = ?", id).Find(&plays)

	return plays
}

func FindNFLPlayByPlaysRecordsByGameID(id string) []structs.NFLPlayByPlay {
	db := dbprovider.GetInstance().GetDB()

	plays := []structs.NFLPlayByPlay{}

	db.Where("game_id = ?", id).Find(&plays)

	return plays
}
