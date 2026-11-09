package repository

import (
    "github.com/CalebRose/SimFBA/structs"
    "github.com/CalebRose/SimFBA/dbprovider"
)

func GetHistoricalCollegeGameplansBySeasonAndWeek(seasonID string, weekID string) []structs.CollegeGameplanRecord {
    db := dbprovider.GetInstance().GetDB()
    var records []structs.CollegeGameplanRecord
    db.Where("season_id = ? AND week_id = ?", seasonID, weekID).Find(&records)
    return records
}

func GetHistoricalNFLGameplansBySeasonAndWeek(seasonID string, weekID string) []structs.NFLGameplanRecord {
    db := dbprovider.GetInstance().GetDB()
    var records []structs.NFLGameplanRecord
    db.Where("season_id = ? AND week_id = ?", seasonID, weekID).Find(&records)
    return records
}
