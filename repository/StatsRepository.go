package repository

import (
	"log"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/structs"
)

type StatsQuery struct {
	SeasonID string
	GameType string
	WeekID   string
	TeamID   string
	PlayerID string
}

func FindCollegePlayerSeasonStatsRecords(clauses StatsQuery) []structs.CollegePlayerSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.CollegePlayerSeasonStats

	query := db.Model(&playerStats)
	if len(clauses.SeasonID) > 0 {
		query = query.Where("season_id = ?", clauses.SeasonID)
	}
	if len(clauses.GameType) > 0 {
		query = query.Where("game_type = ?", clauses.GameType)
	}
	if len(clauses.TeamID) > 0 {
		query = query.Where("team_id = ?", clauses.TeamID)
	}
	if len(clauses.PlayerID) > 0 {
		query = query.Where("player_id = ?", clauses.PlayerID)
	}

	if err := query.Order("passing_yards desc").Find(&playerStats).Error; err != nil {
		// handle the error, for example, log it
		log.Println("Error finding college player season stats records:", err)
	}

	return playerStats
}

func FindProPlayerSeasonStatsRecords(clauses StatsQuery) []structs.NFLPlayerSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.NFLPlayerSeasonStats

	query := db.Model(&playerStats)
	if len(clauses.SeasonID) > 0 {
		query = query.Where("season_id = ?", clauses.SeasonID)
	}
	if len(clauses.GameType) > 0 {
		query = query.Where("game_type = ?", clauses.GameType)
	}
	if len(clauses.TeamID) > 0 {
		query = query.Where("team_id = ?", clauses.TeamID)
	}
	if len(clauses.PlayerID) > 0 {
		query = query.Where("player_id = ?", clauses.PlayerID)
	}

	if err := query.Order("passing_yards desc").Find(&playerStats).Error; err != nil {
		log.Println("Error finding pro player season stats records:", err)
	}

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

	query.Order("passing_yards desc").Find(&playerStats)

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

	query.Order("passing_yards desc").Find(&playerStats)

	return playerStats
}

func FindCollegeTeamSeasonStatsRecords(SeasonID, gameType string) []structs.CollegeTeamSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var teamStats []structs.CollegeTeamSeasonStats

	db.Order("passing_yards desc").Where("season_id = ? AND game_type = ?", SeasonID, gameType).Find(&teamStats)

	return teamStats
}

func FindProTeamSeasonStatsRecords(clauses StatsQuery) []structs.NFLTeamSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var teamStats []structs.NFLTeamSeasonStats

	query := db.Model(&teamStats)
	if len(clauses.SeasonID) > 0 {
		query = query.Where("season_id = ?", clauses.SeasonID)
	}
	if len(clauses.GameType) > 0 {
		query = query.Where("game_type = ?", clauses.GameType)
	}
	if len(clauses.TeamID) > 0 {
		query = query.Where("team_id = ?", clauses.TeamID)
	}
	if len(clauses.PlayerID) > 0 {
		query = query.Where("player_id = ?", clauses.PlayerID)
	}

	if err := query.Order("passing_yards desc").Find(&teamStats).Error; err != nil {
		// handle the error, for example, log it
		log.Println("Error finding pro team season stats records:", err)
	}

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

	query.Order("passing_yards desc").Find(&teamStats)

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

	query.Order("passing_yards desc").Find(&teamStats)

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
