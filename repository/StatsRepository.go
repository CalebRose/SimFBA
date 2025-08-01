package repository

import (
	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/structs"
)

func FindCollegePlayerSeasonStatsRecords(SeasonID, gameType string) []structs.CollegePlayerSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.CollegePlayerSeasonStats

	db.Where("season_id = ? AND game_type = ?", SeasonID, gameType).Find(&playerStats)

	return playerStats
}

func FindProPlayerSeasonStatsRecords(SeasonID, gameType string) []structs.NFLPlayerSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.NFLPlayerSeasonStats

	db.Where("season_id = ? AND game_type = ?", SeasonID, gameType).Find(&playerStats)

	return playerStats
}

func FindCollegePlayerGameStatsRecords(SeasonID, WeekID, GameType, GameID string) []structs.CollegePlayerStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.CollegePlayerStats

	query := db.Model(&playerStats)
	if len(SeasonID) > 0 {
		query = query.Where("season_id = ?", SeasonID)
	}

	if len(WeekID) > 0 {
		query = query.Where("week_id = ?", WeekID)
	}

	if len(GameType) > 0 {
		query = query.Where("game_type = ?", GameType)
	}

	if len(GameID) > 0 {
		query = query.Where("game_id = ?", GameID)
	}

	query.Find(&playerStats)

	return playerStats
}

func FindProPlayerGameStatsRecords(SeasonID, WeekID, GameType, GameID string) []structs.NFLPlayerStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.NFLPlayerStats
	query := db.Model(&playerStats)
	if len(SeasonID) > 0 {
		query = query.Where("season_id = ?", SeasonID)
	}

	if len(WeekID) > 0 {
		query = query.Where("week_id = ?", WeekID)
	}

	if len(GameType) > 0 {
		query = query.Where("game_type = ?", GameType)
	}

	if len(GameID) > 0 {
		query = query.Where("game_id = ?", GameID)
	}

	query.Find(&playerStats)

	return playerStats
}

func FindCollegeTeamSeasonStatsRecords(SeasonID, gameType string) []structs.CollegeTeamSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var teamStats []structs.CollegeTeamSeasonStats

	db.Where("season_id = ? AND game_type = ?", SeasonID, gameType).Find(&teamStats)

	return teamStats
}

func FindProTeamSeasonStatsRecords(SeasonID, gameType string) []structs.NFLTeamSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var teamStats []structs.NFLTeamSeasonStats

	db.Where("season_id = ? AND game_type = ?", SeasonID, gameType).Find(&teamStats)

	return teamStats
}

func FindCollegeTeamGameStatsRecords(SeasonID, WeekID, gameType string) []structs.CollegeTeamStats {
	db := dbprovider.GetInstance().GetDB()

	var teamStats []structs.CollegeTeamStats
	query := db.Model(&teamStats)
	if len(SeasonID) > 0 {
		query = query.Where("season_id = ?", SeasonID)
	}

	if len(WeekID) > 0 {
		query = query.Where("week_id = ?", WeekID)
	}

	if len(gameType) > 0 {
		query = query.Where("game_type = ?", gameType)
	}

	query.Find(&teamStats)

	return teamStats
}

func FindProTeamGameStatsRecords(SeasonID, WeekID, gameType string) []structs.NFLTeamStats {
	db := dbprovider.GetInstance().GetDB()

	var teamStats []structs.NFLTeamStats

	query := db.Model(&teamStats)
	if len(SeasonID) > 0 {
		query = query.Where("season_id = ?", SeasonID)
	}

	if len(WeekID) > 0 {
		query = query.Where("week_id = ?", WeekID)
	}

	if len(gameType) > 0 {
		query = query.Where("game_type = ?", gameType)
	}

	query.Find(&teamStats)

	return teamStats
}

func FindCollegeTeamStatsRecordByGame(gameID, teamID string) structs.CollegeTeamStats {
	db := dbprovider.GetInstance().GetDB()

	var teamStats structs.CollegeTeamStats

	db.Where("game_id = ? AND team_id = ?", gameID, teamID).Find(&teamStats)

	return teamStats
}

func FindProTeamStatsRecordByGame(gameID, teamID string) structs.NFLTeamStats {
	db := dbprovider.GetInstance().GetDB()

	var teamStats structs.NFLTeamStats

	db.Where("game_id = ? AND team_id = ?", gameID, teamID).Find(&teamStats)

	return teamStats
}

func FindCollegePlayerStatsRecordByGame(gameID string) []structs.CollegePlayerStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.CollegePlayerStats

	db.Where("game_id = ?", gameID).Find(&playerStats)

	return playerStats
}

func FindProPlayerStatsRecordByGame(gameID string) []structs.NFLPlayerStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.NFLPlayerStats

	db.Where("game_id = ?", gameID).Find(&playerStats)

	return playerStats
}
