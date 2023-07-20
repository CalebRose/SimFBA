package managers

import (
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/models"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
	"gorm.io/gorm"
)

// GetAllPlayers - Returns all player reference records
func GetAllPlayers() []structs.Player {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.Player

	db.Find(&players)

	return players
}

func GetAllCollegePlayers() []structs.CollegePlayer {
	db := dbprovider.GetInstance().GetDB()

	var CollegePlayers []structs.CollegePlayer

	db.Find(&CollegePlayers)

	return CollegePlayers
}

func GetAllHistoricCollegePlayers() []structs.HistoricCollegePlayer {
	db := dbprovider.GetInstance().GetDB()

	var CollegePlayers []structs.HistoricCollegePlayer

	db.Find(&CollegePlayers)

	return CollegePlayers
}

func GetAllNFLPlayers() []structs.NFLPlayer {
	db := dbprovider.GetInstance().GetDB()

	var nflPlayers []structs.NFLPlayer

	db.Find(&nflPlayers)

	return nflPlayers
}

func GetAllNFLPlayersWithCurrentSeasonStats(seasonID string) []structs.NFLPlayer {
	db := dbprovider.GetInstance().GetDB()

	var nflPlayers []structs.NFLPlayer

	db.Preload("SeasonStats", func(db *gorm.DB) *gorm.DB {
		return db.Where("season_id = ?", seasonID)
	}).Find(&nflPlayers)

	return nflPlayers
}

func GetAllUnsignedPlayers() []structs.UnsignedPlayer {
	db := dbprovider.GetInstance().GetDB()

	var unsignedPlayers []structs.UnsignedPlayer

	db.Find(&unsignedPlayers)

	return unsignedPlayers
}

func GetAllCollegePlayersByTeamId(TeamID string) []structs.CollegePlayer {
	db := dbprovider.GetInstance().GetDB()

	var CollegePlayers []structs.CollegePlayer

	err := db.Order("overall desc").Where("team_id = ?", TeamID).Where("has_graduated = ?", false).Find(&CollegePlayers).Error
	if err != nil {
		fmt.Println(err.Error())
	}

	return CollegePlayers
}

func GetAllCollegePlayersByTeamIdWithoutRedshirts(TeamID string) []structs.CollegePlayer {
	db := dbprovider.GetInstance().GetDB()

	var CollegePlayers []structs.CollegePlayer

	db.Where("team_id = ?", TeamID).Where("is_redshirting = ?", false).Where("has_graduated = ?", false).Find(&CollegePlayers)

	return CollegePlayers
}

func GetCollegePlayerByCollegePlayerId(CollegePlayerId string) structs.CollegePlayer {
	db := dbprovider.GetInstance().GetDB()

	var CollegePlayer structs.CollegePlayer

	db.Where("id = ?", CollegePlayerId).Find(&CollegePlayer)

	return CollegePlayer
}

func GetCollegePlayerByNameAndTeam(firstName string, lastName string, teamID string) models.CollegePlayerCSV {
	db := dbprovider.GetInstance().GetDB()

	var CollegePlayer structs.CollegePlayer

	db.Where("first_name = ? and last_name = ? and team_id = ?", firstName, lastName, teamID).Find(&CollegePlayer)

	collegePlayerResponse := models.MapPlayerToCSVModel(CollegePlayer)

	return collegePlayerResponse
}

func GetCollegePlayerByNameTeamAndWeek(firstName string, lastName string, teamID string, week string) models.CollegePlayerCSV {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	collegeWeek := GetCollegeWeek(week, ts)

	if collegeWeek.ID == uint(ts.CollegeWeekID) {
		return models.CollegePlayerCSV{}
	} else {
		var CollegePlayer structs.CollegePlayer

		db.Preload("Stats", func(db *gorm.DB) *gorm.DB {
			return db.Where("season_id = ? AND week_id = ?", collegeWeek.SeasonID, collegeWeek.ID)
		}).Where("first_name = ? and last_name = ? and team_id = ?", firstName, lastName, teamID).Find(&CollegePlayer)

		collegePlayerResponse := models.MapPlayerForStats(CollegePlayer)

		return collegePlayerResponse
	}
}

func GetSeasonalCollegePlayerByNameTeam(firstName string, lastName string, teamID string) models.CollegePlayerResponse {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	var CollegePlayer structs.CollegePlayer

	db.Preload("SeasonStats", func(db *gorm.DB) *gorm.DB {
		return db.Where("season_id = ? AND week_id < ?", strconv.Itoa(ts.CollegeSeasonID), strconv.Itoa(ts.CollegeWeekID))
	}).Where("first_name = ? and last_name = ? and team_id = ?", firstName, lastName, teamID).Find(&CollegePlayer)

	collegePlayerResponse := models.CollegePlayerResponse{
		ID:          int(CollegePlayer.ID),
		BasePlayer:  CollegePlayer.BasePlayer,
		TeamID:      CollegePlayer.TeamID,
		TeamAbbr:    CollegePlayer.TeamAbbr,
		City:        CollegePlayer.City,
		State:       CollegePlayer.State,
		Year:        CollegePlayer.Year,
		IsRedshirt:  CollegePlayer.IsRedshirt,
		SeasonStats: CollegePlayer.SeasonStats,
	}

	return collegePlayerResponse
}

func UpdateCollegePlayer(cp structs.CollegePlayer) {
	db := dbprovider.GetInstance().GetDB()
	err := db.Save(&cp).Error
	if err != nil {
		log.Fatal(err)
	}
}

func SetRedshirtStatusForPlayer(playerId string) structs.CollegePlayer {
	player := GetCollegePlayerByCollegePlayerId(playerId)

	player.SetRedshirtingStatus()

	UpdateCollegePlayer(player)

	return player
}

func GetAllNFLDraftees() []structs.NFLDraftee {
	db := dbprovider.GetInstance().GetDB()

	var NFLDraftees []structs.NFLDraftee

	db.Order("overall desc").Find(&NFLDraftees)

	return NFLDraftees
}

func GetNFLDrafteeByPlayerID(PlayerID string) structs.NFLDraftee {
	db := dbprovider.GetInstance().GetDB()

	var player structs.NFLDraftee

	db.Where("id = ?", PlayerID).Find(&player)

	return player
}

func GetAllCollegePlayersWithStatsBySeasonID(cMap map[int]int, cNMap map[int]string, seasonID, weekID, viewType string) []models.CollegePlayerResponse {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	seasonIDVal := util.ConvertStringToInt(seasonID)

	var collegePlayers []structs.CollegePlayer

	// var distinctCollegeStats []structs.CollegePlayerStats
	var distinctCollegeStats []structs.CollegePlayerSeasonStats

	db.Distinct("college_player_id").Where("snaps > 0 AND season_id = ?", seasonID).Find(&distinctCollegeStats)

	distinctCollegePlayerIDs := util.GetCollegePlayerIDsBySeasonStats(distinctCollegeStats)

	if viewType == "SEASON" {
		db.Preload("SeasonStats", func(db *gorm.DB) *gorm.DB {
			return db.Where("season_id = ?", seasonID)
		}).Where("id in ?", distinctCollegePlayerIDs).Find(&collegePlayers)
	} else {
		db.Preload("Stats", func(db *gorm.DB) *gorm.DB {
			return db.Where("season_id = ? AND week_id = ? and snaps > 0", seasonID, weekID)
		}).Where("id in ?", distinctCollegePlayerIDs).Find(&collegePlayers)
	}

	var cpResponse []models.CollegePlayerResponse

	for _, player := range collegePlayers {
		if len(player.Stats) == 0 && viewType == "WEEK" {
			continue
		}
		var stat structs.CollegePlayerStats
		if viewType == "WEEK" {
			stat = player.Stats[0]
		}
		cp := models.CollegePlayerResponse{
			ID:           int(player.ID),
			BasePlayer:   player.BasePlayer,
			ConferenceID: cMap[player.TeamID],
			Conference:   cNMap[player.TeamID],
			TeamID:       player.TeamID,
			TeamAbbr:     player.TeamAbbr,
			City:         player.City,
			State:        player.State,
			Year:         player.Year,
			IsRedshirt:   player.IsRedshirt,
			SeasonStats:  player.SeasonStats,
			Stats:        stat,
		}

		// cp.MapSeasonalStats()

		cpResponse = append(cpResponse, cp)
	}

	// If viewing a past season, get all past season players too
	if seasonIDVal < ts.CollegeSeasonID {
		var historicCollegePlayers []structs.HistoricCollegePlayer

		if viewType == "SEASON" {
			db.Preload("SeasonStats", func(db *gorm.DB) *gorm.DB {
				return db.Where("season_id = ?", seasonID)
			}).Where("id in ?", distinctCollegePlayerIDs).Find(&historicCollegePlayers)
		} else {
			db.Preload("Stats", func(db *gorm.DB) *gorm.DB {
				return db.Where("season_id = ? AND week_id = ?", seasonID, weekID)
			}).Where("id in ?", distinctCollegePlayerIDs).Find(&historicCollegePlayers)
		}

		for _, player := range historicCollegePlayers {
			if len(player.Stats) == 0 && viewType == "WEEK" {
				continue
			}
			var stat structs.CollegePlayerStats
			if viewType == "WEEK" {
				stat = player.Stats[0]
			}
			cp := models.CollegePlayerResponse{
				ID:           int(player.ID),
				BasePlayer:   player.BasePlayer,
				ConferenceID: cMap[player.TeamID],
				Conference:   cNMap[player.TeamID],
				TeamID:       player.TeamID,
				TeamAbbr:     player.TeamAbbr,
				City:         player.City,
				State:        player.State,
				Year:         player.Year,
				IsRedshirt:   player.IsRedshirt,
				SeasonStats:  player.SeasonStats,
				Stats:        stat,
			}

			cpResponse = append(cpResponse, cp)
		}
	}

	return cpResponse
}

func GetAllNFLPlayersWithStatsBySeasonID(cMap, dMap map[int]int, cNMap, dNMap map[int]string, seasonID, weekID, viewType string) []models.NFLPlayerResponse {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	seasonIDVal := util.ConvertStringToInt(seasonID)

	var nflPlayers []structs.NFLPlayer

	// var distinctNFLStats []structs.CollegePlayerStats
	var distinctNFLStats []structs.NFLPlayerSeasonStats

	db.Distinct("nfl_player_id").Where("snaps > 0 AND season_id = ?", seasonID).Find(&distinctNFLStats)

	distinctCollegePlayerIDs := util.GetNFLPlayerIDsBySeasonStats(distinctNFLStats)

	if viewType == "SEASON" {
		db.Preload("SeasonStats", func(db *gorm.DB) *gorm.DB {
			return db.Where("season_id = ?", seasonID)
		}).Where("id in ?", distinctCollegePlayerIDs).Find(&nflPlayers)
	} else {
		db.Preload("Stats", func(db *gorm.DB) *gorm.DB {
			return db.Where("season_id = ? AND week_id = ? and snaps > 0", seasonID, weekID)
		}).Where("id in ?", distinctCollegePlayerIDs).Find(&nflPlayers)
	}

	var cpResponse []models.NFLPlayerResponse

	for _, player := range nflPlayers {
		if len(player.Stats) == 0 && viewType == "WEEK" {
			continue
		}
		var stat structs.NFLPlayerStats
		if viewType == "WEEK" {
			stat = player.Stats[0]
		}
		cp := models.NFLPlayerResponse{
			ID:           int(player.ID),
			BasePlayer:   player.BasePlayer,
			ConferenceID: cMap[player.TeamID],
			Conference:   cNMap[player.TeamID],
			DivisionID:   dMap[player.TeamID],
			Division:     dNMap[player.TeamID],
			TeamID:       player.TeamID,
			TeamAbbr:     player.TeamAbbr,
			State:        player.State,
			Year:         int(player.Experience),
			SeasonStats:  player.SeasonStats,
			Stats:        stat,
		}

		// cp.MapSeasonalStats()

		cpResponse = append(cpResponse, cp)
	}

	// If viewing a past season, get all past season players too
	if seasonIDVal < ts.NFLSeasonID {
		var historicNFLPlayers []structs.NFLRetiredPlayer

		if viewType == "SEASON" {
			db.Preload("SeasonStats", func(db *gorm.DB) *gorm.DB {
				return db.Where("season_id = ?", seasonID)
			}).Where("id in ?", distinctCollegePlayerIDs).Find(&historicNFLPlayers)
		} else {
			db.Preload("Stats", func(db *gorm.DB) *gorm.DB {
				return db.Where("season_id = ? AND week_id = ?", seasonID, weekID)
			}).Where("id in ?", distinctCollegePlayerIDs).Find(&historicNFLPlayers)
		}

		for _, player := range historicNFLPlayers {
			if len(player.Stats) == 0 && viewType == "WEEK" {
				continue
			}
			var stat structs.NFLPlayerStats
			if viewType == "WEEK" {
				stat = player.Stats[0]
			}
			cp := models.NFLPlayerResponse{
				ID:           int(player.ID),
				BasePlayer:   player.BasePlayer,
				ConferenceID: cMap[player.TeamID],
				Conference:   cNMap[player.TeamID],
				TeamID:       player.TeamID,
				TeamAbbr:     player.TeamAbbr,
				State:        player.State,
				Year:         int(player.Experience),
				SeasonStats:  player.SeasonStats,
				Stats:        stat,
			}

			cpResponse = append(cpResponse, cp)
		}
	}

	return cpResponse
}

func GetAllCollegePlayersWithStatsByTeamID(TeamID string, SeasonID string) []structs.CollegePlayer {
	db := dbprovider.GetInstance().GetDB()

	var collegePlayers []structs.CollegePlayer

	db.Preload("Stats", func(db *gorm.DB) *gorm.DB {
		return db.Where("season_id = ? and team_id = ? and snaps > 0", SeasonID, TeamID)
	}).Where("team_id = ?", TeamID).Find(&collegePlayers)

	return collegePlayers
}

func GetAllCollegePlayersWithGameStatsByTeamID(TeamID string, GameID string) []structs.GameResultsPlayer {
	db := dbprovider.GetInstance().GetDB()

	var collegePlayers []structs.CollegePlayer
	var matchRows []structs.GameResultsPlayer

	db.Preload("Stats", func(db *gorm.DB) *gorm.DB {
		return db.Where("game_id = ? and team_id = ? and snaps > 0", GameID, TeamID)
	}).Where("team_id = ?", TeamID).Find(&collegePlayers)

	for _, p := range collegePlayers {
		if len(p.Stats) == 0 {
			continue
		}

		s := p.Stats[0]
		if s.Snaps == 0 {
			continue
		}

		row := structs.GameResultsPlayer{
			FirstName:            p.FirstName,
			LastName:             p.LastName,
			Position:             p.Position,
			Archetype:            p.Archetype,
			Year:                 uint(p.Year),
			League:               "CFB",
			Snaps:                s.Snaps,
			PassingYards:         s.PassingYards,
			PassAttempts:         s.PassAttempts,
			PassCompletions:      s.PassCompletions,
			PassingTDs:           s.PassingTDs,
			Interceptions:        s.Interceptions,
			LongestPass:          s.LongestPass,
			Sacks:                s.Sacks,
			RushAttempts:         s.RushAttempts,
			RushingYards:         s.RushingYards,
			RushingTDs:           s.RushingTDs,
			Fumbles:              s.Fumbles,
			LongestRush:          s.LongestRush,
			Targets:              s.Targets,
			Catches:              s.Catches,
			ReceivingYards:       s.ReceivingYards,
			ReceivingTDs:         s.ReceivingTDs,
			LongestReception:     s.LongestReception,
			SoloTackles:          s.SoloTackles,
			AssistedTackles:      s.AssistedTackles,
			TacklesForLoss:       s.TacklesForLoss,
			SacksMade:            s.SacksMade,
			ForcedFumbles:        s.ForcedFumbles,
			RecoveredFumbles:     s.RecoveredFumbles,
			PassDeflections:      s.PassDeflections,
			InterceptionsCaught:  s.InterceptionsCaught,
			Safeties:             s.Safeties,
			DefensiveTDs:         s.DefensiveTDs,
			FGMade:               s.FGMade,
			FGAttempts:           s.FGAttempts,
			LongestFG:            s.LongestFG,
			ExtraPointsMade:      s.ExtraPointsMade,
			ExtraPointsAttempted: s.ExtraPointsAttempted,
			KickoffTouchbacks:    s.KickoffTouchbacks,
			Punts:                s.Punts,
			PuntTouchbacks:       s.PuntTouchbacks,
			PuntsInside20:        s.PuntsInside20,
			KickReturns:          s.KickReturns,
			KickReturnTDs:        s.KickReturnTDs,
			KickReturnYards:      s.KickReturnYards,
			PuntReturns:          s.PuntReturns,
			PuntReturnTDs:        s.PuntReturnTDs,
			PuntReturnYards:      s.PuntReturnYards,
			STSoloTackles:        s.STSoloTackles,
			STAssistedTackles:    s.STAssistedTackles,
			PuntsBlocked:         s.PuntsBlocked,
			FGBlocked:            s.FGBlocked,
			Pancakes:             s.Pancakes,
			SacksAllowed:         s.SacksAllowed,
			PlayedGame:           s.PlayedGame,
			StartedGame:          s.StartedGame,
			WasInjured:           s.WasInjured,
			WeeksOfRecovery:      s.WeeksOfRecovery,
			InjuryType:           s.InjuryType,
		}

		matchRows = append(matchRows, row)
	}

	historicPlayers := []structs.HistoricCollegePlayer{}
	db.Preload("Stats", func(db *gorm.DB) *gorm.DB {
		return db.Where("game_id = ? and team_id = ? and snaps > 0", GameID, TeamID)
	}).Where("team_id = ?", TeamID).Find(&historicPlayers)

	for _, p := range historicPlayers {
		if len(p.Stats) == 0 {
			continue
		}

		s := p.Stats[0]
		if s.Snaps == 0 {
			continue
		}

		row := structs.GameResultsPlayer{
			FirstName:            p.FirstName,
			LastName:             p.LastName,
			Position:             p.Position,
			Archetype:            p.Archetype,
			Year:                 uint(p.Year),
			League:               "CFB",
			Snaps:                s.Snaps,
			PassingYards:         s.PassingYards,
			PassAttempts:         s.PassAttempts,
			PassCompletions:      s.PassCompletions,
			PassingTDs:           s.PassingTDs,
			Interceptions:        s.Interceptions,
			LongestPass:          s.LongestPass,
			Sacks:                s.Sacks,
			RushAttempts:         s.RushAttempts,
			RushingYards:         s.RushingYards,
			RushingTDs:           s.RushingTDs,
			Fumbles:              s.Fumbles,
			LongestRush:          s.LongestRush,
			Targets:              s.Targets,
			Catches:              s.Catches,
			ReceivingYards:       s.ReceivingYards,
			ReceivingTDs:         s.ReceivingTDs,
			LongestReception:     s.LongestReception,
			SoloTackles:          s.SoloTackles,
			AssistedTackles:      s.AssistedTackles,
			TacklesForLoss:       s.TacklesForLoss,
			SacksMade:            s.SacksMade,
			ForcedFumbles:        s.ForcedFumbles,
			RecoveredFumbles:     s.RecoveredFumbles,
			PassDeflections:      s.PassDeflections,
			InterceptionsCaught:  s.InterceptionsCaught,
			Safeties:             s.Safeties,
			DefensiveTDs:         s.DefensiveTDs,
			FGMade:               s.FGMade,
			FGAttempts:           s.FGAttempts,
			LongestFG:            s.LongestFG,
			ExtraPointsMade:      s.ExtraPointsMade,
			ExtraPointsAttempted: s.ExtraPointsAttempted,
			KickoffTouchbacks:    s.KickoffTouchbacks,
			Punts:                s.Punts,
			PuntTouchbacks:       s.PuntTouchbacks,
			PuntsInside20:        s.PuntsInside20,
			KickReturns:          s.KickReturns,
			KickReturnTDs:        s.KickReturnTDs,
			KickReturnYards:      s.KickReturnYards,
			PuntReturns:          s.PuntReturns,
			PuntReturnTDs:        s.PuntReturnTDs,
			PuntReturnYards:      s.PuntReturnYards,
			STSoloTackles:        s.STSoloTackles,
			STAssistedTackles:    s.STAssistedTackles,
			PuntsBlocked:         s.PuntsBlocked,
			FGBlocked:            s.FGBlocked,
			Pancakes:             s.Pancakes,
			SacksAllowed:         s.SacksAllowed,
			PlayedGame:           s.PlayedGame,
			StartedGame:          s.StartedGame,
			WasInjured:           s.WasInjured,
			WeeksOfRecovery:      s.WeeksOfRecovery,
			InjuryType:           s.InjuryType,
		}

		matchRows = append(matchRows, row)
	}

	return matchRows
}

func GetAllNFLPlayersWithGameStatsByTeamID(TeamID string, GameID string) []structs.GameResultsPlayer {
	db := dbprovider.GetInstance().GetDB()

	var nflPlayers []structs.NFLPlayer
	var matchRows []structs.GameResultsPlayer

	db.Preload("Stats", func(db *gorm.DB) *gorm.DB {
		return db.Where("game_id = ? and team_id = ? and snaps > 0", GameID, TeamID)
	}).Where("team_id = ?", TeamID).Find(&nflPlayers)

	for _, p := range nflPlayers {
		if len(p.Stats) == 0 {
			continue
		}

		s := p.Stats[0]
		if s.Snaps == 0 {
			continue
		}

		row := structs.GameResultsPlayer{
			FirstName:            p.FirstName,
			LastName:             p.LastName,
			Position:             p.Position,
			Archetype:            p.Archetype,
			Year:                 uint(s.Year),
			League:               "NFL",
			Snaps:                s.Snaps,
			PassingYards:         s.PassingYards,
			PassAttempts:         s.PassAttempts,
			PassCompletions:      s.PassCompletions,
			PassingTDs:           s.PassingTDs,
			Interceptions:        s.Interceptions,
			LongestPass:          s.LongestPass,
			Sacks:                s.Sacks,
			RushAttempts:         s.RushAttempts,
			RushingYards:         s.RushingYards,
			RushingTDs:           s.RushingTDs,
			Fumbles:              s.Fumbles,
			LongestRush:          s.LongestRush,
			Targets:              s.Targets,
			Catches:              s.Catches,
			ReceivingYards:       s.ReceivingYards,
			ReceivingTDs:         s.ReceivingTDs,
			LongestReception:     s.LongestReception,
			SoloTackles:          s.SoloTackles,
			AssistedTackles:      s.AssistedTackles,
			TacklesForLoss:       s.TacklesForLoss,
			SacksMade:            s.SacksMade,
			ForcedFumbles:        s.ForcedFumbles,
			RecoveredFumbles:     s.RecoveredFumbles,
			PassDeflections:      s.PassDeflections,
			InterceptionsCaught:  s.InterceptionsCaught,
			Safeties:             s.Safeties,
			DefensiveTDs:         s.DefensiveTDs,
			FGMade:               s.FGMade,
			FGAttempts:           s.FGAttempts,
			LongestFG:            s.LongestFG,
			ExtraPointsMade:      s.ExtraPointsMade,
			ExtraPointsAttempted: s.ExtraPointsAttempted,
			KickoffTouchbacks:    s.KickoffTouchbacks,
			Punts:                s.Punts,
			PuntTouchbacks:       s.PuntTouchbacks,
			PuntsInside20:        s.PuntsInside20,
			KickReturns:          s.KickReturns,
			KickReturnTDs:        s.KickReturnTDs,
			KickReturnYards:      s.KickReturnYards,
			PuntReturns:          s.PuntReturns,
			PuntReturnTDs:        s.PuntReturnTDs,
			PuntReturnYards:      s.PuntReturnYards,
			STSoloTackles:        s.STSoloTackles,
			STAssistedTackles:    s.STAssistedTackles,
			PuntsBlocked:         s.PuntsBlocked,
			FGBlocked:            s.FGBlocked,
			Pancakes:             s.Pancakes,
			SacksAllowed:         s.SacksAllowed,
			PlayedGame:           s.PlayedGame,
			StartedGame:          s.StartedGame,
			WasInjured:           s.WasInjured,
			WeeksOfRecovery:      s.WeeksOfRecovery,
			InjuryType:           s.InjuryType,
		}

		matchRows = append(matchRows, row)
	}

	historicPlayers := []structs.NFLRetiredPlayer{}
	db.Preload("Stats", func(db *gorm.DB) *gorm.DB {
		return db.Where("game_id = ? and team_id = ? and snaps > 0", GameID, TeamID)
	}).Where("team_id = ?", TeamID).Find(&historicPlayers)

	for _, p := range historicPlayers {
		if len(p.Stats) == 0 {
			continue
		}

		s := p.Stats[0]
		if s.Snaps == 0 {
			continue
		}

		row := structs.GameResultsPlayer{
			FirstName:            p.FirstName,
			LastName:             p.LastName,
			Position:             p.Position,
			Archetype:            p.Archetype,
			Year:                 uint(s.Year),
			League:               "NFL",
			Snaps:                s.Snaps,
			PassingYards:         s.PassingYards,
			PassAttempts:         s.PassAttempts,
			PassCompletions:      s.PassCompletions,
			PassingTDs:           s.PassingTDs,
			Interceptions:        s.Interceptions,
			LongestPass:          s.LongestPass,
			Sacks:                s.Sacks,
			RushAttempts:         s.RushAttempts,
			RushingYards:         s.RushingYards,
			RushingTDs:           s.RushingTDs,
			Fumbles:              s.Fumbles,
			LongestRush:          s.LongestRush,
			Targets:              s.Targets,
			Catches:              s.Catches,
			ReceivingYards:       s.ReceivingYards,
			ReceivingTDs:         s.ReceivingTDs,
			LongestReception:     s.LongestReception,
			SoloTackles:          s.SoloTackles,
			AssistedTackles:      s.AssistedTackles,
			TacklesForLoss:       s.TacklesForLoss,
			SacksMade:            s.SacksMade,
			ForcedFumbles:        s.ForcedFumbles,
			RecoveredFumbles:     s.RecoveredFumbles,
			PassDeflections:      s.PassDeflections,
			InterceptionsCaught:  s.InterceptionsCaught,
			Safeties:             s.Safeties,
			DefensiveTDs:         s.DefensiveTDs,
			FGMade:               s.FGMade,
			FGAttempts:           s.FGAttempts,
			LongestFG:            s.LongestFG,
			ExtraPointsMade:      s.ExtraPointsMade,
			ExtraPointsAttempted: s.ExtraPointsAttempted,
			KickoffTouchbacks:    s.KickoffTouchbacks,
			Punts:                s.Punts,
			PuntTouchbacks:       s.PuntTouchbacks,
			PuntsInside20:        s.PuntsInside20,
			KickReturns:          s.KickReturns,
			KickReturnTDs:        s.KickReturnTDs,
			KickReturnYards:      s.KickReturnYards,
			PuntReturns:          s.PuntReturns,
			PuntReturnTDs:        s.PuntReturnTDs,
			PuntReturnYards:      s.PuntReturnYards,
			STSoloTackles:        s.STSoloTackles,
			STAssistedTackles:    s.STAssistedTackles,
			PuntsBlocked:         s.PuntsBlocked,
			FGBlocked:            s.FGBlocked,
			Pancakes:             s.Pancakes,
			SacksAllowed:         s.SacksAllowed,
			PlayedGame:           s.PlayedGame,
			StartedGame:          s.StartedGame,
			WasInjured:           s.WasInjured,
			WeeksOfRecovery:      s.WeeksOfRecovery,
			InjuryType:           s.InjuryType,
		}

		matchRows = append(matchRows, row)
	}

	return matchRows
}

func GetHeismanList() []models.HeismanWatchModel {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	var collegePlayers []structs.CollegePlayer

	var heismanCandidates []models.HeismanWatchModel

	teamWithStandings := GetAllCollegeTeamsWithCurrentYearStandings()

	var teamWeight = make(map[string]float64)

	var homeTeamMapper = make(map[int]string)
	var teamGameMapper = make(map[int][]structs.CollegeGame)

	for _, team := range teamWithStandings {
		homeTeamMapper[int(team.ID)] = team.TeamAbbr

		games := GetCollegeGamesByTeamIdAndSeasonId(strconv.Itoa(int(team.ID)), strconv.Itoa(ts.CollegeSeasonID))

		if len(games) == 0 || len(team.TeamStandings) == 0 {
			continue
		}

		teamGameMapper[int(team.ID)] = games

		currentYearStandings := team.TeamStandings[0]

		var weight float64 = 0 // 1
		if currentYearStandings.TotalLosses+currentYearStandings.TotalWins > 0 {
			// newWeight := (float64(currentYearStandings.TotalWins) / 12) + 1

			// if newWeight > weight {
			// 	weight = newWeight
			// }
			weight = float64(currentYearStandings.TotalWins) / 100
		}

		teamWeight[team.TeamAbbr] = weight
	}

	// db.Preload("Stats", func(db *gorm.DB) *gorm.DB {
	// 	return db.Where("snaps > 0 and season_id = ? and week_id < ?", strconv.Itoa(ts.CollegeSeasonID), strconv.Itoa(ts.CollegeWeekID))
	// }).Where("Stats.snaps > 0").Find(&collegePlayers)

	var distinctCollegeStats []structs.CollegePlayerStats

	db.Distinct("college_player_id").Where("snaps > 0").Find(&distinctCollegeStats)

	distinctCollegePlayerIDs := util.GetCollegePlayerIDs(distinctCollegeStats)

	db.Preload("Stats", func(db *gorm.DB) *gorm.DB {
		return db.Where("snaps > 0 and season_id = ? and week_id < ?", strconv.Itoa(ts.CollegeSeasonID), strconv.Itoa(ts.CollegeWeekID))
	}).Where("id IN ?", distinctCollegePlayerIDs).Find(&collegePlayers)

	for _, cp := range collegePlayers {
		if len(cp.Stats) == 0 {
			continue
		}

		score := util.GetHeismanScore(cp, teamWeight, homeTeamMapper, teamGameMapper[cp.TeamID])

		h := models.HeismanWatchModel{
			FirstName: cp.FirstName,
			LastName:  cp.LastName,
			Position:  cp.Position,
			Archetype: cp.Archetype,
			School:    cp.TeamAbbr,
			Score:     score,
			Games:     len(cp.Stats),
		}

		heismanCandidates = append(heismanCandidates, h)

	}

	sort.Sort(models.ByScore(heismanCandidates))

	return heismanCandidates
}

func GetGlobalPlayerRecord(playerID string) structs.Player {
	db := dbprovider.GetInstance().GetDB()

	var player structs.Player

	db.Where("id = ?", playerID).Find(&player)

	return player
}

func GetOnlyNFLPlayerRecord(playerID string) structs.NFLPlayer {
	db := dbprovider.GetInstance().GetDB()

	var player structs.NFLPlayer

	db.Where("id = ?", playerID).Find(&player)

	return player
}

func GetNFLPlayerRecord(playerID string) structs.NFLPlayer {
	db := dbprovider.GetInstance().GetDB()

	var player structs.NFLPlayer

	db.Preload("Contract").Where("id = ?", playerID).Find(&player)

	return player
}

func GetNFLPlayersForDCPage(TeamID string) []structs.NFLPlayer {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.NFLPlayer

	db.Preload("Contract").Where("team_id = ? AND is_practice_squad = ?", TeamID, false).Find(&players)

	return players
}

func GetTradableNFLPlayersByTeamID(TeamID string) []structs.NFLPlayer {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.NFLPlayer

	db.Preload("Contract").Where("team_id = ? AND is_on_trade_block = ?", TeamID, true).Find(&players)

	return players
}

func GetNFLPlayersWithContractsByTeamID(TeamID string) []structs.NFLPlayer {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.NFLPlayer

	db.Preload("Contract", func(db *gorm.DB) *gorm.DB {
		return db.Where("is_active = true")
	}).Where("team_id = ?", TeamID).Find(&players)

	return players
}

func CutNFLPlayer(playerId string) {
	db := dbprovider.GetInstance().GetDB()

	player := GetOnlyNFLPlayerRecord(playerId)
	contract := GetContractByPlayerID(playerId)
	capsheet := GetCapsheetByTeamID(strconv.Itoa(int(player.TeamID)))
	ts := GetTimestamp()

	if player.Experience < 4 && !ts.IsNFLOffSeason && !player.IsPracticeSquad {
		player.WaivePlayer()
	} else {
		player.ToggleIsFreeAgent()
		contract.DeactivateContract()
	}

	capsheet.CutPlayerFromCapsheet(contract)
	db.Save(&contract)
	db.Save(&player)
	db.Save(&capsheet)
}

func PlaceNFLPlayerOnPracticeSquad(playerId string) {
	db := dbprovider.GetInstance().GetDB()

	player := GetOnlyNFLPlayerRecord(playerId)
	player.ToggleIsPracticeSquad()
	teamID := strconv.Itoa(player.TeamID)

	contract := GetContractByPlayerID(playerId)
	contract.DeactivateContract()
	db.Save(&contract)

	teamCapsheet := GetCapsheetByTeamID(teamID)
	teamCapsheet.SubtractFromCapsheet(contract)
	db.Save(&teamCapsheet)

	db.Save(&player)
}

func GetNFLRosterForSimulation(TeamID string) []structs.NFLPlayer {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.NFLPlayer

	db.Where("team_id = ? AND is_practice_squad = ?", TeamID, false).Find(&players)

	return players
}

func RecoverPlayers() {
	db := dbprovider.GetInstance().GetDB()

	collegePlayers := GetAllCollegePlayers()

	for _, p := range collegePlayers {
		if !p.IsInjured {
			continue
		}
		p.RecoveryCheck()
		db.Save(&p)
	}

	nflPlayers := GetAllNFLPlayers()

	for _, p := range nflPlayers {
		if !p.IsInjured {
			continue
		}

		p.RecoveryCheck()
		db.Save(&p)
	}
}

func CheckNFLRookiesForLetterGrade(seasonID string) {
	db := dbprovider.GetInstance().GetDB()

	nflPlayers := GetAllNFLPlayersWithCurrentSeasonStats(seasonID)

	for _, p := range nflPlayers {
		if !p.ShowLetterGrade {
			continue
		}

		seasonStats := p.SeasonStats
		if seasonStats.Snaps >= 250 {
			p.ShowRealAttributeValue()
			db.Save(&p)
		}
	}
}

func GetAllPracticeSquadPlayers() []structs.NFLPlayer {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.NFLPlayer

	db.Where("is_practice_squad = ?", true).Find(&players)

	return players
}

func GetAllPracticeSquadPlayersForFAPage() []structs.NFLPlayer {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.NFLPlayer

	db.Preload("Offers", func(db *gorm.DB) *gorm.DB {
		return db.Order("contract_value DESC").Where("is_active = true")
	}).Where("is_practice_squad = ?", true).Find(&players)

	return players
}

func GetInjuredCollegePlayers() []structs.CollegePlayer {
	db := dbprovider.GetInstance().GetDB()

	var collegePlayers []structs.CollegePlayer

	db.Order("team_id asc").Where("is_injured = true").Find(&collegePlayers)

	return collegePlayers
}

func GetInjuredNFLPlayers() []structs.NFLPlayer {
	db := dbprovider.GetInstance().GetDB()

	var nflPlayers []structs.NFLPlayer

	db.Order("team_id asc").Where("is_injured = true").Find(&nflPlayers)

	return nflPlayers
}
