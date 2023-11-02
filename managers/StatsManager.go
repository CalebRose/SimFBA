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

func GetNFLPlayerStatsByPlayerIDAndSeason(PlayerID string, SeasonID string) []structs.NFLPlayerStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.NFLPlayerStats

	db.Where("nfl_player_id = ? and season_id = ?", PlayerID, SeasonID).Find(&playerStats)

	return playerStats
}

func GetAllCollegePlayerStatsByGame(GameID string, TeamID string) []structs.CollegePlayerStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.CollegePlayerStats

	db.Where("game_id = ? and team_id = ?", GameID, TeamID).Find(&playerStats)

	return playerStats
}

func GetAllNFLPlayerStatsByGame(GameID string, TeamID string) []structs.NFLPlayerStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.NFLPlayerStats

	db.Where("game_id = ? and team_id = ?", GameID, TeamID).Find(&playerStats)

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

	err := db.Preload("TeamSeasonStats", func(db *gorm.DB) *gorm.DB {
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
		SeasonStats:  collegeTeam.TeamSeasonStats,
	}

	return ct
}

func GetCollegeTeamSeasonStatsBySeason(TeamID string, SeasonID string) structs.CollegeTeamSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var teamStats structs.CollegeTeamSeasonStats

	err := db.Where("team_id = ? AND season_id = ?", TeamID, SeasonID).Find(&teamStats).Error
	if err != nil {
		return structs.CollegeTeamSeasonStats{}
	}

	return teamStats
}

func GetNFLTeamSeasonStatsByTeamANDSeason(TeamID string, SeasonID string) structs.NFLTeamSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var teamStats structs.NFLTeamSeasonStats

	err := db.Where("team_id = ? AND season_id = ?", TeamID, SeasonID).Find(&teamStats).Error
	if err != nil {
		return structs.NFLTeamSeasonStats{}
	}

	return teamStats
}

func GetCollegeSeasonStatsByPlayerAndSeason(PlayerID, SeasonID string) structs.CollegePlayerSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats structs.CollegePlayerSeasonStats

	err := db.Where("college_player_id = ? AND season_id = ?", PlayerID, SeasonID).Find(&playerStats).Error
	if err != nil {
		return structs.CollegePlayerSeasonStats{}
	}

	return playerStats
}

func GetNFLSeasonStatsByPlayerAndSeason(PlayerID, SeasonID string) structs.NFLPlayerSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats structs.NFLPlayerSeasonStats

	err := db.Where("nfl_player_id = ? AND season_id = ?", PlayerID, SeasonID).Find(&playerStats).Error
	if err != nil {
		return structs.NFLPlayerSeasonStats{}
	}

	return playerStats
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

func GetAllNFLPlayerSeasonStatsByPlayerID(playerID string) []structs.NFLPlayerSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.NFLPlayerSeasonStats

	db.Where("nfl_player_id = ?", playerID).Find(&playerStats)

	return playerStats
}

func GetCollegeTeamStatsByGame(TeamID string, GameID string) structs.CollegeTeamStats {
	db := dbprovider.GetInstance().GetDB()

	var teamStats structs.CollegeTeamStats

	db.Where("team_id = ? AND game_id = ?", TeamID, GameID).Find(&teamStats)

	return teamStats
}

func GetNFLTeamStatsByGame(TeamID string, GameID string) structs.NFLTeamStats {
	db := dbprovider.GetInstance().GetDB()

	var teamStats structs.NFLTeamStats

	db.Where("team_id = ? AND game_id = ?", TeamID, GameID).Find(&teamStats)

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

func GetNFLHistoricalTeamStats(TeamID string, SeasonID string) []structs.NFLTeamStats {
	db := dbprovider.GetInstance().GetDB()

	var teamStats []structs.NFLTeamStats

	db.Where("team_id = ? AND season_id = ?", TeamID, SeasonID).Find(&teamStats)

	return teamStats
}

func GetAllCollegeTeamsWithStatsBySeasonID(seasonID, weekID, viewType string) []models.CollegeTeamResponse {
	db := dbprovider.GetInstance().GetDB()

	var teams []structs.CollegeTeam

	if viewType == "SEASON" {
		db.Preload("TeamSeasonStats", func(db *gorm.DB) *gorm.DB {
			return db.Where("season_id = ?", seasonID)
		}).Find(&teams)
	} else {
		db.Preload("TeamStats", func(db *gorm.DB) *gorm.DB {
			return db.Where("season_id = ? AND week_id = ?", seasonID, weekID)
		}).Find(&teams)
	}

	var ctResponse []models.CollegeTeamResponse

	for _, team := range teams {
		if len(team.TeamStats) == 0 && viewType == "WEEK" {
			continue
		}
		var teamstat structs.CollegeTeamStats
		if viewType == "WEEK" {
			teamstat = team.TeamStats[0]
		}
		ct := models.CollegeTeamResponse{
			ID:           int(team.ID),
			BaseTeam:     team.BaseTeam,
			ConferenceID: team.ConferenceID,
			Conference:   team.Conference,
			DivisionID:   team.DivisionID,
			Division:     team.Division,
			SeasonStats:  team.TeamSeasonStats,
			Stats:        teamstat,
		}

		ctResponse = append(ctResponse, ct)
	}

	return ctResponse
}

func GetAllNFLTeamsWithStatsBySeasonID(seasonID, weekID, viewType string) []models.NFLTeamResponse {
	db := dbprovider.GetInstance().GetDB()

	var teams []structs.NFLTeam

	if viewType == "SEASON" {
		db.Preload("TeamSeasonStats", func(db *gorm.DB) *gorm.DB {
			return db.Where("season_id = ?", seasonID)
		}).Find(&teams)
	} else {
		db.Preload("TeamStats", func(db *gorm.DB) *gorm.DB {
			return db.Where("season_id = ? AND week_id = ?", seasonID, weekID)
		}).Find(&teams)
	}

	var ctResponse []models.NFLTeamResponse

	for _, team := range teams {
		if len(team.TeamStats) == 0 && viewType == "WEEK" {
			continue
		}
		var teamstat structs.NFLTeamStats
		var seasonstat structs.NFLTeamSeasonStats
		if viewType == "WEEK" && len(team.TeamStats) > 0 {
			teamstat = team.TeamStats[0]
		} else if len(team.TeamSeasonStats) > 0 {
			seasonstat = team.TeamSeasonStats[0]
		}
		ct := models.NFLTeamResponse{
			ID:           int(team.ID),
			BaseTeam:     team.BaseTeam,
			ConferenceID: int(team.ConferenceID),
			Conference:   team.Conference,
			DivisionID:   int(team.DivisionID),
			Division:     team.Division,
			SeasonStats:  seasonstat,
			Stats:        teamstat,
		}

		ctResponse = append(ctResponse, ct)
	}

	return ctResponse
}

func ResetCFBSeasonalStats() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.CollegeSeasonID))
	teams := GetAllCollegeTeams()

	for _, team := range teams {
		teamID := strconv.Itoa(int(team.ID))
		teamStats := GetHistoricalTeamStats(teamID, seasonID)
		seasonStats := GetCollegeTeamSeasonStatsBySeason(teamID, seasonID)
		seasonStats.ResetStats()
		seasonStats.MapStats(teamStats)
		db.Save(&seasonStats)
		fmt.Println("Reset Season Stats for " + team.TeamName)
	}

	players := GetAllCollegePlayers()

	for _, player := range players {
		playerID := strconv.Itoa(int(player.ID))
		playerStats := GetCollegePlayerStatsByPlayerIDAndSeason(playerID, seasonID)
		if len(playerStats) > 0 {
			seasonStats := GetCollegeSeasonStatsByPlayerAndSeason(playerID, seasonID)
			seasonStats.ResetStats()
			seasonStats.MapStats(playerStats)
			db.Save(&seasonStats)
		}
		fmt.Println("Reset Season Stats for " + player.FirstName + " " + player.LastName + " " + player.Position)
	}

	standings := GetAllConferenceStandingsBySeasonID(seasonID)

	for _, standing := range standings {
		if standing.ID == 0 {
			continue
		}
		standing.ResetCFBStandings()

		games := GetCollegeGamesByTeamIdAndSeasonId(strconv.Itoa(int(standing.TeamID)), seasonID)

		for _, game := range games {
			if !game.GameComplete {
				continue
			}
			standing.UpdateCollegeStandings(game)
		}

		db.Save(&standing)
		fmt.Println("Updated Standings for " + standing.TeamName)
	}
}

func ResetNFLSeasonalStats() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.NFLSeasonID))
	teams := GetAllNFLTeams()

	for _, team := range teams {
		teamID := strconv.Itoa(int(team.ID))
		teamStats := GetNFLHistoricalTeamStats(teamID, seasonID)
		seasonStats := GetNFLTeamSeasonStatsByTeamANDSeason(teamID, seasonID)
		seasonStats.ResetStats()
		seasonStats.MapStats(teamStats)
		db.Save(&seasonStats)
		fmt.Println("Reset Season Stats for " + team.TeamName)
	}

	players := GetAllNFLPlayers()

	for _, player := range players {
		playerID := strconv.Itoa(int(player.ID))
		playerStats := GetNFLPlayerStatsByPlayerIDAndSeason(playerID, seasonID)
		if len(playerStats) > 0 {
			seasonStats := GetNFLSeasonStatsByPlayerAndSeason(playerID, seasonID)
			seasonStats.ResetStats()
			seasonStats.MapStats(playerStats, ts)
			db.Save(&seasonStats)
		}
		fmt.Println("Reset Season Stats for " + player.FirstName + " " + player.LastName + " " + player.Position)
	}

	standings := GetAllNFLStandingsBySeasonID(seasonID)

	for _, standing := range standings {
		if standing.ID == 0 {
			continue
		}
		standing.ResetNFLStandings()

		games := GetNFLGamesByTeamIdAndSeasonId(strconv.Itoa(int(standing.TeamID)), seasonID)

		for _, game := range games {
			if !game.GameComplete {
				continue
			}
			standing.UpdateNFLStandings(game)
		}

		db.Save(&standing)
		fmt.Println("Updated Standings for " + standing.TeamName)
	}
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
		gameID := strconv.Itoa(int(gameDataDTO.GameID))
		record := make(chan structs.CollegeGame)

		go func() {
			asynchronousGame := GetCollegeGameByGameID(gameID)
			record <- asynchronousGame
		}()

		gameRecord := <-record
		close(record)
		var playerStats []structs.CollegePlayerStats

		homeTeamID := strconv.Itoa(int(gameRecord.HomeTeamID))
		awayTeamID := strconv.Itoa(int(gameRecord.AwayTeamID))

		// Team Stats Export
		homeTeamChn := make(chan structs.CollegeTeam)

		go func() {
			homeTeam := GetTeamByTeamID(homeTeamID)
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
			awayTeam := GetTeamByTeamID(awayTeamID)
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
			collegePlayerStats := structs.CollegePlayerStats{
				CollegePlayerID: player.GetPlayerID(),
				TeamID:          homeTeam.TeamID,
				GameID:          homeTeam.GameID,
				WeekID:          gameRecord.WeekID,
				SeasonID:        gameRecord.SeasonID,
				OpposingTeam:    at.TeamAbbr,
				BasePlayerStats: player.MapTobasePlayerStatsObject(),
				Year:            player.Year,
			}
			playerStats = append(playerStats, collegePlayerStats)
		}

		// AwayPlayers
		for _, player := range gameDataDTO.AwayPlayers {
			collegePlayerStats := structs.CollegePlayerStats{
				CollegePlayerID: player.GetPlayerID(),
				TeamID:          awayTeam.TeamID,
				GameID:          awayTeam.GameID,
				WeekID:          gameRecord.WeekID,
				SeasonID:        gameRecord.SeasonID,
				OpposingTeam:    ht.TeamAbbr,
				BasePlayerStats: player.MapTobasePlayerStatsObject(),
				Year:            player.Year,
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
		gameID := strconv.Itoa(int(gameDataDTO.GameID))
		record := make(chan structs.NFLGame)

		go func() {
			asynchronousGame := GetNFLGameByGameID(gameID)
			record <- asynchronousGame
		}()

		gameRecord := <-record
		close(record)

		homeTeamID := strconv.Itoa(int(gameRecord.HomeTeamID))
		awayTeamID := strconv.Itoa(int(gameRecord.AwayTeamID))

		var playerStats []structs.NFLPlayerStats

		// Team Stats Export
		homeTeamChn := make(chan structs.NFLTeam)

		go func() {
			homeTeam := GetNFLTeamByTeamID(homeTeamID)
			homeTeamChn <- homeTeam
		}()

		ht := <-homeTeamChn
		close(homeTeamChn)

		// Away Team
		awayTeamChn := make(chan structs.NFLTeam)

		go func() {
			awayTeam := GetNFLTeamByTeamID(awayTeamID)
			awayTeamChn <- awayTeam
		}()

		at := <-awayTeamChn
		close(awayTeamChn)

		homeTeam := structs.NFLTeamStats{
			TeamID:        ht.ID,
			GameID:        gameRecord.ID,
			WeekID:        uint(gameRecord.WeekID),
			SeasonID:      uint(gameRecord.SeasonID),
			OpposingTeam:  at.TeamAbbr,
			BaseTeamStats: gameDataDTO.HomeTeam.MapToBaseTeamStatsObject(),
		}

		teamStats = append(teamStats, homeTeam)
		awayTeam := structs.NFLTeamStats{
			TeamID:        at.ID,
			GameID:        gameRecord.ID,
			WeekID:        uint(gameRecord.WeekID),
			SeasonID:      uint(gameRecord.SeasonID),
			OpposingTeam:  ht.TeamAbbr,
			BaseTeamStats: gameDataDTO.AwayTeam.MapToBaseTeamStatsObject(),
		}

		teamStats = append(teamStats, awayTeam)

		// Player Stat Export
		// HomePlayers
		for _, player := range gameDataDTO.HomePlayers {
			nflPlayerStats := structs.NFLPlayerStats{
				NFLPlayerID:     player.GetPlayerID(),
				TeamID:          int(homeTeam.TeamID),
				GameID:          int(homeTeam.GameID),
				WeekID:          gameRecord.WeekID,
				SeasonID:        gameRecord.SeasonID,
				OpposingTeam:    gameDataDTO.AwayTeam.Abbreviation,
				BasePlayerStats: player.MapTobasePlayerStatsObject(),
				Year:            player.Year,
			}
			playerStats = append(playerStats, nflPlayerStats)
		}

		// AwayPlayers
		for _, player := range gameDataDTO.AwayPlayers {
			nflPlayerStats := structs.NFLPlayerStats{
				NFLPlayerID:     player.GetPlayerID(),
				TeamID:          int(awayTeam.TeamID),
				GameID:          int(awayTeam.GameID),
				WeekID:          gameRecord.WeekID,
				SeasonID:        gameRecord.SeasonID,
				OpposingTeam:    gameDataDTO.HomeTeam.Abbreviation,
				BasePlayerStats: player.MapTobasePlayerStatsObject(),
				Year:            player.Year,
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

func GetCFBGameResultsByGameID(gameID string) structs.GameResultsResponse {
	game := GetCollegeGameByGameID(gameID)

	homePlayers := GetAllCollegePlayersWithGameStatsByTeamID(strconv.Itoa(game.HomeTeamID), gameID)
	awayPlayers := GetAllCollegePlayersWithGameStatsByTeamID(strconv.Itoa(game.AwayTeamID), gameID)

	return structs.GameResultsResponse{
		HomePlayers: homePlayers,
		AwayPlayers: awayPlayers,
	}
}

func GetNFLGameResultsByGameID(gameID string) structs.GameResultsResponse {
	game := GetNFLGameByGameID(gameID)

	homePlayers := GetAllNFLPlayersWithGameStatsByTeamID(strconv.Itoa(game.HomeTeamID), gameID)
	awayPlayers := GetAllNFLPlayersWithGameStatsByTeamID(strconv.Itoa(game.AwayTeamID), gameID)

	return structs.GameResultsResponse{
		HomePlayers: homePlayers,
		AwayPlayers: awayPlayers,
	}
}
