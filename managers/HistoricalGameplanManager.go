package managers

import (
    "github.com/CalebRose/SimFBA/repository"
    "github.com/CalebRose/SimFBA/structs"
)

func GetHistoricalCollegeGameplansBySeasonAndWeek(seasonID string, weekID string) []structs.CollegeGameplanRecord {
    return repository.GetHistoricalCollegeGameplansBySeasonAndWeek(seasonID, weekID)
}

func GetHistoricalNFLGameplansBySeasonAndWeek(seasonID string, weekID string) []structs.NFLGameplanRecord {
    return repository.GetHistoricalNFLGameplansBySeasonAndWeek(seasonID, weekID)
}
