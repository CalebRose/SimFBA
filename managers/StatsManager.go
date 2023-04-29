package managers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/models"
	"github.com/CalebRose/SimFBA/structs"
	"gorm.io/gorm"
)

func GetCollegePlayerStatsByGame(PlayerID string, GameID string) structs.CollegePlayerStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats structs.CollegePlayerStats

	db.Where("college_player_id = ? and game_id = ?", PlayerID, GameID).Find(&playerStats)

	return playerStats
}

func GetCareerCollegePlayerStatsByPlayerID(PlayerID string) []structs.CollegePlayerStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.CollegePlayerStats

	db.Where("college_player_id = ?", PlayerID).Find(&playerStats)

	return playerStats
}

func GetCollegePlayerStatsByPlayerIDAndSeason(PlayerID string, SeasonID string) []structs.CollegePlayerStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.CollegePlayerStats

	db.Where("college_player_id = ? and season_id = ?", PlayerID, SeasonID).Find(&playerStats)

	return playerStats
}

func GetAllCollegePlayerStatsByGame(GameID string, TeamID string) []structs.CollegePlayerStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.CollegePlayerStats

	db.Where("game_id = ? and season_id = ?", GameID, TeamID).Find(&playerStats)

	return playerStats
}

func GetAllNFLPlayerStatsByGame(GameID string, TeamID string) []structs.NFLPlayerStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.NFLPlayerStats

	db.Where("game_id = ? and season_id = ?", GameID, TeamID).Find(&playerStats)

	return playerStats
}

func GetAllPlayerStatsByWeek(WeekID string) []structs.CollegePlayerStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.CollegePlayerStats

	db.Where("week_id = ?", WeekID).Find(&playerStats)

	return playerStats
}

func GetTeamStatsByWeekAndTeam(TeamID string, Week string) structs.CollegeTeam {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	collegeWeek := GetCollegeWeek(Week, ts)
	var collegeTeam structs.CollegeTeam

	if collegeWeek.ID == uint(ts.CollegeWeekID) {
		return structs.CollegeTeam{}
	} else {
		err := db.Preload("TeamStats", func(db *gorm.DB) *gorm.DB {
			return db.Where("season_id = ? AND week_id = ?", collegeWeek.SeasonID, collegeWeek.ID)
		}).Where("id = ?", TeamID).Find(&collegeTeam).Error
		if err != nil {
			fmt.Println("Could not find college team and stats from week")
		}

	}
	return collegeTeam
}

// TEAM STATS
func GetSeasonalTeamStats(TeamID string, SeasonID string) models.CollegeTeamResponse {
	db := dbprovider.GetInstance().GetDB()

	var collegeTeam structs.CollegeTeam

	err := db.Preload("TeamStats", func(db *gorm.DB) *gorm.DB {
		return db.Where("season_id = ?", SeasonID)
	}).Where("id = ?", TeamID).Find(&collegeTeam).Error
	if err != nil {
		fmt.Println("Could not find college team and stats from week")
	}

	ct := models.CollegeTeamResponse{
		ID:           int(collegeTeam.ID),
		BaseTeam:     collegeTeam.BaseTeam,
		ConferenceID: collegeTeam.ConferenceID,
		Conference:   collegeTeam.Conference,
		DivisionID:   collegeTeam.DivisionID,
		Division:     collegeTeam.Division,
		TeamStats:    collegeTeam.TeamStats,
	}

	ct.MapSeasonalStats()

	return ct
}

func GetCollegeTeamSeasonStatsBySeason(TeamID string, SeasonID string) structs.CollegeTeamSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var teamStats structs.CollegeTeamSeasonStats

	db.Where("team_id = ?, season_id = ?", TeamID, SeasonID).Find(&teamStats)

	return teamStats
}

func GetCollegeSeasonStatsBySeason(SeasonID string) []structs.CollegeTeamSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var teamStats []structs.CollegeTeamSeasonStats

	db.Where("season_id = ?", SeasonID).Find(&teamStats)

	return teamStats
}

func GetCollegePlayerSeasonStatsBySeason(SeasonID string) []structs.CollegePlayerSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.CollegePlayerSeasonStats

	db.Where("season_id = ?", SeasonID).Find(&playerStats)

	return playerStats
}

func GetNFLTeamSeasonStatsBySeason(SeasonID string) []structs.NFLTeamSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var teamStats []structs.NFLTeamSeasonStats

	db.Where("season_id = ?", SeasonID).Find(&teamStats)

	return teamStats
}

func GetNFLPlayerSeasonStatsBySeason(SeasonID string) []structs.NFLPlayerSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.NFLPlayerSeasonStats

	db.Where("season_id = ?", SeasonID).Find(&playerStats)

	return playerStats
}

func GetCollegeTeamStatsByGame(TeamID string, GameID string) structs.CollegeTeamStats {
	db := dbprovider.GetInstance().GetDB()

	var teamStats structs.CollegeTeamStats

	db.Where("team_id = ?, game_id = ?", TeamID, GameID).Find(&teamStats)

	return teamStats
}

func GetNFLTeamStatsByGame(TeamID string, GameID string) structs.NFLTeamStats {
	db := dbprovider.GetInstance().GetDB()

	var teamStats structs.NFLTeamStats

	db.Where("team_id = ?, game_id = ?", TeamID, GameID).Find(&teamStats)

	return teamStats
}

func GetCollegeTeamStatsByWeek(WeekID string) []structs.CollegeTeamStats {
	db := dbprovider.GetInstance().GetDB()

	var teamStats []structs.CollegeTeamStats

	db.Where("week_id = ?", WeekID).Find(&teamStats)

	return teamStats
}

func GetCollegeTeamStatsBySeason(SeasonID string) []structs.CollegeTeamStats {
	db := dbprovider.GetInstance().GetDB()

	var teamStats []structs.CollegeTeamStats

	db.Where("season_id = ?", SeasonID).Find(&teamStats)

	return teamStats
}

func GetHistoricalTeamStats(TeamID string, SeasonID string) []structs.CollegeTeamStats {
	db := dbprovider.GetInstance().GetDB()

	var teamStats []structs.CollegeTeamStats

	db.Where("team_id = ? AND season_id = ?", TeamID, SeasonID).Find(&teamStats)

	return teamStats
}

func ExportCFBStatisticsFromSim(gameStats []structs.GameStatDTO) {
	db := dbprovider.GetInstance().GetDB()
	fmt.Println("START")

	tsChn := make(chan structs.Timestamp)

	go func() {
		ts := GetTimestamp()
		tsChn <- ts
	}()

	timestamp := <-tsChn
	close(tsChn)

	var teamStats []structs.CollegeTeamStats

	for _, gameDataDTO := range gameStats {

		record := make(chan structs.CollegeGame)

		go func() {
			asynchronousGame := GetCollegeGameByAbbreviationsWeekAndSeasonID(gameDataDTO.HomeTeam.Abbreviation, strconv.Itoa(timestamp.CollegeWeekID), strconv.Itoa(timestamp.CollegeSeasonID))
			record <- asynchronousGame
		}()

		gameRecord := <-record
		close(record)
		var playerStats []structs.CollegePlayerStats

		// Team Stats Export
		homeTeamChn := make(chan structs.CollegeTeam)

		go func() {
			homeTeam := GetTeamByTeamAbbr(gameDataDTO.HomeTeam.Abbreviation)
			homeTeamChn <- homeTeam
		}()

		ht := <-homeTeamChn
		close(homeTeamChn)

		homeTeam := structs.CollegeTeamStats{
			TeamID:        int(ht.ID),
			GameID:        int(gameRecord.ID),
			WeekID:        gameRecord.WeekID,
			SeasonID:      gameRecord.SeasonID,
			OpposingTeam:  gameDataDTO.AwayTeam.Abbreviation,
			BaseTeamStats: gameDataDTO.HomeTeam.MapToBaseTeamStatsObject(),
		}

		teamStats = append(teamStats, homeTeam)

		// Away Team
		awayTeamChn := make(chan structs.CollegeTeam)

		go func() {
			awayTeam := GetTeamByTeamAbbr(gameDataDTO.AwayTeam.Abbreviation)
			awayTeamChn <- awayTeam
		}()

		at := <-awayTeamChn
		close(awayTeamChn)

		awayTeam := structs.CollegeTeamStats{
			TeamID:        int(at.ID),
			GameID:        int(gameRecord.ID),
			WeekID:        gameRecord.WeekID,
			SeasonID:      gameRecord.SeasonID,
			OpposingTeam:  gameDataDTO.HomeTeam.Abbreviation,
			BaseTeamStats: gameDataDTO.AwayTeam.MapToBaseTeamStatsObject(),
		}

		teamStats = append(teamStats, awayTeam)

		// Player Stat Export
		// HomePlayers
		for _, player := range gameDataDTO.HomePlayers {
			if player.IsInjured && player.WeeksOfRecovery > 0 {
				playerRecord := GetCollegePlayerByCollegePlayerId(strconv.Itoa(player.PlayerID))
				playerRecord.SetIsInjured(player.IsInjured, player.InjuryType, player.WeeksOfRecovery)
				db.Save(&playerRecord)
			}
			collegePlayerStats := structs.CollegePlayerStats{
				CollegePlayerID: player.GetPlayerID(),
				TeamID:          homeTeam.TeamID,
				GameID:          homeTeam.GameID,
				WeekID:          gameRecord.WeekID,
				SeasonID:        gameRecord.SeasonID,
				OpposingTeam:    gameDataDTO.AwayTeam.Abbreviation,
				BasePlayerStats: player.MapTobasePlayerStatsObject(),
			}
			playerStats = append(playerStats, collegePlayerStats)
		}

		// AwayPlayers
		for _, player := range gameDataDTO.AwayPlayers {
			if player.IsInjured && player.WeeksOfRecovery > 0 {
				playerRecord := GetCollegePlayerByCollegePlayerId(strconv.Itoa(player.PlayerID))
				playerRecord.SetIsInjured(player.IsInjured, player.InjuryType, player.WeeksOfRecovery)
				db.Save(&playerRecord)
			}
			collegePlayerStats := structs.CollegePlayerStats{
				CollegePlayerID: player.GetPlayerID(),
				TeamID:          awayTeam.TeamID,
				GameID:          awayTeam.GameID,
				WeekID:          gameRecord.WeekID,
				SeasonID:        gameRecord.SeasonID,
				OpposingTeam:    gameDataDTO.HomeTeam.Abbreviation,
				BasePlayerStats: player.MapTobasePlayerStatsObject(),
			}
			playerStats = append(playerStats, collegePlayerStats)
		}

		// Update Game
		gameRecord.UpdateScore(gameDataDTO.HomeScore, gameDataDTO.AwayScore)

		err := db.Save(&gameRecord).Error
		if err != nil {
			log.Panicln("Could not save Game " + strconv.Itoa(int(gameRecord.ID)) + "Between " + gameRecord.HomeTeam + " and " + gameRecord.AwayTeam)
		}

		err = db.CreateInBatches(&playerStats, len(playerStats)).Error
		if err != nil {
			log.Panicln("Could not save player stats from week " + strconv.Itoa(timestamp.CollegeWeek))
		}

		fmt.Println("Finished Game " + strconv.Itoa(int(gameRecord.ID)) + " Between " + gameRecord.HomeTeam + " and " + gameRecord.AwayTeam)
	}

	err := db.CreateInBatches(&teamStats, len(teamStats)).Error
	if err != nil {
		log.Panicln("Could not save team stats!")
	}
}

func GetAllCollegeTeamsWithStatsBySeasonID(seasonID string) []models.CollegeTeamResponse {
	db := dbprovider.GetInstance().GetDB()

	var teams []structs.CollegeTeam

	db.Preload("TeamSeasonStats", func(db *gorm.DB) *gorm.DB {
		return db.Where("season_id = ?", seasonID)
	}).Find(&teams)

	var ctResponse []models.CollegeTeamResponse

	for _, team := range teams {
		ct := models.CollegeTeamResponse{
			ID:           int(team.ID),
			BaseTeam:     team.BaseTeam,
			ConferenceID: team.ConferenceID,
			Conference:   team.Conference,
			DivisionID:   team.DivisionID,
			Division:     team.Division,
			SeasonStats:  team.TeamSeasonStats,
		}

		ctResponse = append(ctResponse, ct)
	}

	return ctResponse
}

func MapAllStatsToSeason() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	teams := GetAllCollegeTeams()

	for _, team := range teams {
		teamStats := GetHistoricalTeamStats(strconv.Itoa(int(team.ID)), strconv.Itoa(ts.CollegeSeasonID))

		seasonStats := structs.CollegeTeamSeasonStats{
			TeamID:   team.ID,
			SeasonID: uint(ts.CollegeSeasonID),
		}

		seasonStats.MapStats(teamStats)

		db.Save(&seasonStats)
		fmt.Println("Saved Season Stats for " + team.TeamName)
	}

	players := GetAllCollegePlayers()

	for _, player := range players {
		playerStats := GetCollegePlayerStatsByPlayerIDAndSeason(strconv.Itoa(int(player.ID)), strconv.Itoa(ts.CollegeSeasonID))

		seasonStats := structs.CollegePlayerSeasonStats{
			CollegePlayerID: player.ID,
			TeamID:          uint(player.TeamID),
			SeasonID:        uint(ts.CollegeSeasonID),
		}

		seasonStats.MapStats(playerStats)

		db.Save(&seasonStats)
		fmt.Println("Saved Season Stats for " + player.FirstName + " " + player.LastName + " " + player.Position)
	}
}

func ExportNFLStatisticsFromSim(gameStats []structs.GameStatDTO) {
	db := dbprovider.GetInstance().GetDB()
	fmt.Println("START")

	tsChn := make(chan structs.Timestamp)

	go func() {
		ts := GetTimestamp()
		tsChn <- ts
	}()

	timestamp := <-tsChn
	close(tsChn)

	var teamStats []structs.NFLTeamStats

	for _, gameDataDTO := range gameStats {

		record := make(chan structs.NFLGame)

		go func() {
			asynchronousGame := GetNFLGameByAbbreviationsWeekAndSeasonID(gameDataDTO.HomeTeam.Abbreviation, strconv.Itoa(timestamp.NFLWeekID), strconv.Itoa(timestamp.NFLSeasonID))
			record <- asynchronousGame
		}()

		gameRecord := <-record
		close(record)
		var playerStats []structs.NFLPlayerStats

		// Team Stats Export
		homeTeamChn := make(chan structs.NFLTeam)

		go func() {
			homeTeam := GetNFLTeamByTeamAbbr(gameDataDTO.HomeTeam.Abbreviation)
			homeTeamChn <- homeTeam
		}()

		ht := <-homeTeamChn
		close(homeTeamChn)

		homeTeam := structs.NFLTeamStats{
			TeamID:        ht.ID,
			GameID:        gameRecord.ID,
			WeekID:        uint(gameRecord.WeekID),
			SeasonID:      uint(gameRecord.SeasonID),
			OpposingTeam:  gameDataDTO.AwayTeam.Abbreviation,
			BaseTeamStats: gameDataDTO.HomeTeam.MapToBaseTeamStatsObject(),
		}

		teamStats = append(teamStats, homeTeam)

		// Away Team
		awayTeamChn := make(chan structs.NFLTeam)

		go func() {
			awayTeam := GetNFLTeamByTeamAbbr(gameDataDTO.AwayTeam.Abbreviation)
			awayTeamChn <- awayTeam
		}()

		at := <-awayTeamChn
		close(awayTeamChn)

		awayTeam := structs.NFLTeamStats{
			TeamID:        at.ID,
			GameID:        gameRecord.ID,
			WeekID:        uint(gameRecord.WeekID),
			SeasonID:      uint(gameRecord.SeasonID),
			OpposingTeam:  gameDataDTO.HomeTeam.Abbreviation,
			BaseTeamStats: gameDataDTO.AwayTeam.MapToBaseTeamStatsObject(),
		}

		teamStats = append(teamStats, awayTeam)

		// Player Stat Export
		// HomePlayers
		for _, player := range gameDataDTO.HomePlayers {
			if player.IsInjured && player.WeeksOfRecovery > 0 {
				playerRecord := GetNFLPlayerRecord(strconv.Itoa(player.PlayerID))
				playerRecord.SetIsInjured(player.IsInjured, player.InjuryType, player.WeeksOfRecovery)
				db.Save(&playerRecord)
			}
			nflPlayerStats := structs.NFLPlayerStats{
				NFLPlayerID:     player.GetPlayerID(),
				TeamID:          int(homeTeam.TeamID),
				GameID:          int(homeTeam.GameID),
				WeekID:          gameRecord.WeekID,
				SeasonID:        gameRecord.SeasonID,
				OpposingTeam:    gameDataDTO.AwayTeam.Abbreviation,
				BasePlayerStats: player.MapTobasePlayerStatsObject(),
			}
			playerStats = append(playerStats, nflPlayerStats)
		}

		// AwayPlayers
		for _, player := range gameDataDTO.AwayPlayers {
			if player.IsInjured && player.WeeksOfRecovery > 0 {
				playerRecord := GetNFLPlayerRecord(strconv.Itoa(player.PlayerID))
				playerRecord.SetIsInjured(player.IsInjured, player.InjuryType, player.WeeksOfRecovery)
				db.Save(&playerRecord)
			}

			nflPlayerStats := structs.NFLPlayerStats{
				NFLPlayerID:     player.GetPlayerID(),
				TeamID:          int(awayTeam.TeamID),
				GameID:          int(awayTeam.GameID),
				WeekID:          gameRecord.WeekID,
				SeasonID:        gameRecord.SeasonID,
				OpposingTeam:    gameDataDTO.HomeTeam.Abbreviation,
				BasePlayerStats: player.MapTobasePlayerStatsObject(),
			}

			playerStats = append(playerStats, nflPlayerStats)
		}

		// Update Game
		gameRecord.UpdateScore(gameDataDTO.HomeScore, gameDataDTO.AwayScore)

		err := db.Save(&gameRecord).Error
		if err != nil {
			log.Panicln("Could not save Game " + strconv.Itoa(int(gameRecord.ID)) + "Between " + gameRecord.HomeTeam + " and " + gameRecord.AwayTeam)
		}

		err = db.CreateInBatches(&playerStats, len(playerStats)).Error
		if err != nil {
			log.Panicln("Could not save player stats from week " + strconv.Itoa(timestamp.CollegeWeek))
		}

		fmt.Println("Finished Game " + strconv.Itoa(int(gameRecord.ID)) + " Between " + gameRecord.HomeTeam + " and " + gameRecord.AwayTeam)
	}

	err := db.CreateInBatches(&teamStats, len(teamStats)).Error
	if err != nil {
		log.Panicln("Could not save team stats!")
	}
}
