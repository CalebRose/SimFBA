package managers

import (
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/models"
	"github.com/CalebRose/SimFBA/repository"
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

func GetAllRetiredPlayers() []structs.NFLRetiredPlayer {
	db := dbprovider.GetInstance().GetDB()

	var nflPlayers []structs.NFLRetiredPlayer

	db.Find(&nflPlayers)

	return nflPlayers
}

func GetAllNFLPlayersWithCurrentSeasonStats(seasonID, gameType string) []structs.NFLPlayer {
	db := dbprovider.GetInstance().GetDB()

	var nflPlayers []structs.NFLPlayer

	db.Preload("SeasonStats", func(db *gorm.DB) *gorm.DB {
		return db.Where("season_id = ? AND game_type = ?", seasonID, gameType)
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

func GetCollegePlayerViaDiscord(id, gameType string) structs.DiscordPlayerResponse {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	seasonID := strconv.Itoa(ts.CollegeSeasonID)
	var collegePlayer structs.CollegePlayer

	db.Preload("SeasonStats", "season_id = ? AND game_type = ?", seasonID, gameType).Where("id = ?", id).Find(&collegePlayer)

	collegePlayerResponse := structs.MapPlayerToCSVModel(collegePlayer)

	return structs.DiscordPlayerResponse{
		Player:       collegePlayerResponse,
		CollegeStats: collegePlayer.SeasonStats,
	}
}

func GetCollegePlayerByNameViaDiscord(firstName, lastName, teamID, gameType string) structs.DiscordPlayerResponse {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	seasonID := strconv.Itoa(ts.CollegeSeasonID)
	var collegePlayer structs.CollegePlayer

	db.Preload("SeasonStats", "season_id = ? AND game_type = ?", seasonID, gameType).Where("first_name = ? AND last_name = ? and team_id = ?", firstName, lastName, teamID).Find(&collegePlayer)

	collegePlayerResponse := structs.MapPlayerToCSVModel(collegePlayer)

	return structs.DiscordPlayerResponse{
		Player:       collegePlayerResponse,
		CollegeStats: collegePlayer.SeasonStats,
	}
}

func GetNFLPlayerViaDiscord(id, gameType string) structs.DiscordPlayerResponse {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	seasonID := strconv.Itoa(ts.NFLSeasonID)
	var nflPlayer structs.NFLPlayer

	db.Preload("SeasonStats", "season_id = ? AND game_type = ?", seasonID, gameType).Where("id = ?", id).Find(&nflPlayer)

	nflPlayerResponse := structs.MapNFLPlayerToCSVModel(nflPlayer)

	return structs.DiscordPlayerResponse{
		Player:   nflPlayerResponse,
		NFLStats: nflPlayer.SeasonStats,
	}
}

func GetCareerNFLPlayerViaDiscord(id string) structs.DiscordPlayerResponse {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	var nflPlayer structs.NFLPlayer

	db.Where("id = ?", id).Find(&nflPlayer)

	indStats := []structs.NFLPlayerStats{}

	db.Where("nfl_player_id = ?", strconv.Itoa(int(nflPlayer.ID))).Find(&indStats)

	seasonStats := structs.NFLPlayerSeasonStats{}
	seasonStats.MapStats(indStats, ts)

	nflPlayerResponse := structs.MapNFLPlayerToCSVModel(nflPlayer)

	return structs.DiscordPlayerResponse{
		Player:   nflPlayerResponse,
		NFLStats: seasonStats,
	}
}

func GetNFLPlayerByNameViaDiscord(firstName, lastName, teamID, gameType string) structs.DiscordPlayerResponse {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	seasonID := strconv.Itoa(ts.NFLSeasonID)
	var nflPlayer structs.NFLPlayer

	db.Preload("SeasonStats", "season_id = ? AND game_type = ?", seasonID, gameType).Where("first_name = ? AND last_name = ? and team_id = ?", firstName, lastName, teamID).Find(&nflPlayer)

	nflPlayerResponse := structs.MapNFLPlayerToCSVModel(nflPlayer)

	return structs.DiscordPlayerResponse{
		Player:   nflPlayerResponse,
		NFLStats: nflPlayer.SeasonStats,
	}
}

func GetCollegePlayerByIdAndWeek(id, week string) structs.CollegePlayerCSV {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	collegeWeek := GetCollegeWeek(week, ts)

	if collegeWeek.ID == uint(ts.CollegeWeekID) {
		return structs.CollegePlayerCSV{}
	} else {
		var CollegePlayer structs.CollegePlayer

		db.Preload("Stats", func(db *gorm.DB) *gorm.DB {
			return db.Where("season_id = ? AND week_id = ?", collegeWeek.SeasonID, collegeWeek.ID)
		}).Where("id = ?", id).Find(&CollegePlayer)

		collegePlayerResponse := structs.MapPlayerForStats(CollegePlayer)

		return collegePlayerResponse
	}
}

func GetCareerCollegePlayerByNameTeam(id string) structs.CollegePlayerResponse {
	db := dbprovider.GetInstance().GetDB()

	var CollegePlayer structs.CollegePlayer

	db.Where("id = ?", id).Find(&CollegePlayer)

	indStats := []structs.CollegePlayerStats{}

	db.Where("college_player_id = ?", strconv.Itoa(int(CollegePlayer.ID))).Find(&indStats)

	seasonStats := structs.CollegePlayerSeasonStats{}
	seasonStats.MapStats(indStats)

	collegePlayerResponse := structs.CollegePlayerResponse{
		ID:          int(CollegePlayer.ID),
		BasePlayer:  CollegePlayer.BasePlayer,
		TeamID:      CollegePlayer.TeamID,
		TeamAbbr:    CollegePlayer.TeamAbbr,
		City:        CollegePlayer.City,
		State:       CollegePlayer.State,
		Year:        CollegePlayer.Year,
		IsRedshirt:  CollegePlayer.IsRedshirt,
		SeasonStats: seasonStats,
	}

	return collegePlayerResponse
}

func GetSeasonalCollegePlayerByNameTeam(id string) structs.CollegePlayerResponse {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	var CollegePlayer structs.CollegePlayer

	db.Preload("SeasonStats", func(db *gorm.DB) *gorm.DB {
		return db.Where("season_id = ?", strconv.Itoa(ts.CollegeSeasonID))
	}).Where("id = ?", id).Find(&CollegePlayer)

	collegePlayerResponse := structs.CollegePlayerResponse{
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

func GetAllNFLDraftees() []models.NFLDraftee {
	db := dbprovider.GetInstance().GetDB()

	var NFLDraftees []models.NFLDraftee

	db.Find(&NFLDraftees)

	sort.Slice(NFLDraftees, func(i, j int) bool {
		iVal := util.GetNumericalSortValueByLetterGrade(NFLDraftees[i].OverallGrade)
		jVal := util.GetNumericalSortValueByLetterGrade(NFLDraftees[j].OverallGrade)
		return iVal < jVal
	})

	return NFLDraftees
}

func GetNFLDrafteeByPlayerID(PlayerID string) models.NFLDraftee {
	db := dbprovider.GetInstance().GetDB()

	var player models.NFLDraftee

	db.Where("id = ?", PlayerID).Find(&player)

	return player
}

func GetAllCollegePlayersWithStatsBySeasonID(cMap map[int]int, cNMap map[int]string, seasonID, weekID, viewType, gameType string) []structs.CollegePlayerResponse {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	seasonIDVal := util.ConvertStringToInt(seasonID)

	var collegePlayers []structs.CollegePlayer

	// var distinctCollegeStats []structs.CollegePlayerStats
	var distinctCollegeStats []structs.CollegePlayerSeasonStats

	db.Distinct("college_player_id").Where("snaps > 0 AND season_id = ? AND game_type = ?", seasonID, gameType).Find(&distinctCollegeStats)

	distinctCollegePlayerIDs := GetCollegePlayerIDsBySeasonStats(distinctCollegeStats)

	if viewType == "SEASON" {
		db.Preload("SeasonStats", func(db *gorm.DB) *gorm.DB {
			return db.Where("season_id = ? AND game_type = ?", seasonID, gameType)
		}).Where("id in ?", distinctCollegePlayerIDs).Find(&collegePlayers)
	} else {
		db.Preload("Stats", func(db *gorm.DB) *gorm.DB {
			return db.Where("season_id = ? AND week_id = ? and snaps > 0 AND reveal_results = ? AND game_type = ?", seasonID, weekID, true, gameType)
		}).Where("id in ?", distinctCollegePlayerIDs).Find(&collegePlayers)
	}

	var cpResponse []structs.CollegePlayerResponse

	for _, player := range collegePlayers {
		if len(player.Stats) == 0 && viewType == "WEEK" {
			continue
		}
		var stat structs.CollegePlayerStats
		if viewType == "WEEK" {
			stat = player.Stats[0]
		}
		cp := structs.CollegePlayerResponse{
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
	if seasonIDVal <= ts.CollegeSeasonID {
		var historicCollegePlayers []structs.HistoricCollegePlayer

		if viewType == "SEASON" {
			db.Preload("SeasonStats", func(db *gorm.DB) *gorm.DB {
				return db.Where("season_id = ? AND game_type = ?", seasonID, gameType)
			}).Where("id in ?", distinctCollegePlayerIDs).Find(&historicCollegePlayers)
		} else {
			db.Preload("Stats", func(db *gorm.DB) *gorm.DB {
				return db.Where("season_id = ? AND week_id = ? AND game_type = ?", seasonID, weekID, gameType)
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
			cp := structs.CollegePlayerResponse{
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

func GetAllNFLPlayersWithStatsBySeasonID(cMap, dMap map[int]int, cNMap, dNMap map[int]string, seasonID, weekID, viewType, gameType string) []structs.NFLPlayerResponse {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	seasonIDVal := util.ConvertStringToInt(seasonID)

	var nflPlayers []structs.NFLPlayer

	// var distinctNFLStats []structs.CollegePlayerStats
	var distinctNFLStats []structs.NFLPlayerSeasonStats

	db.Distinct("nfl_player_id").Where("snaps > 0 AND season_id = ?", seasonID).Find(&distinctNFLStats)

	distinctCollegePlayerIDs := GetNFLPlayerIDsBySeasonStats(distinctNFLStats)

	if viewType == "SEASON" {
		db.Preload("SeasonStats", func(db *gorm.DB) *gorm.DB {
			return db.Where("season_id = ? AND game_type = ?", seasonID, gameType)
		}).Where("id in ?", distinctCollegePlayerIDs).Find(&nflPlayers)
	} else {
		db.Preload("Stats", func(db *gorm.DB) *gorm.DB {
			return db.Where("season_id = ? AND week_id = ? and snaps > 0 AND reveal_results = ? AND game_type = ?", seasonID, weekID, true, gameType)
		}).Where("id in ?", distinctCollegePlayerIDs).Find(&nflPlayers)
	}

	var cpResponse []structs.NFLPlayerResponse

	for _, player := range nflPlayers {
		if len(player.Stats) == 0 && viewType == "WEEK" {
			continue
		}
		var stat structs.NFLPlayerStats
		if viewType == "WEEK" {
			stat = player.Stats[0]
		}
		cp := structs.NFLPlayerResponse{
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
				return db.Where("season_id = ? AND game_type = ?", seasonID, gameType)
			}).Where("id in ?", distinctCollegePlayerIDs).Find(&historicNFLPlayers)
		} else {
			db.Preload("Stats", func(db *gorm.DB) *gorm.DB {
				return db.Where("season_id = ? AND week_id = ? AND game_type = ?", seasonID, weekID, gameType)
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
			cp := structs.NFLPlayerResponse{
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

func GetAllCollegePlayersWithSeasonStatsByTeamID(TeamID, SeasonID, gameType string) []structs.CollegePlayerResponse {
	db := dbprovider.GetInstance().GetDB()

	var collegePlayers []structs.CollegePlayer

	var responseList []structs.CollegePlayerResponse

	db.Preload("SeasonStats", func(db *gorm.DB) *gorm.DB {
		return db.Where("season_id = ? and team_id = ? and snaps > 0", SeasonID, TeamID)
	}).Where("team_id = ?", TeamID).Find(&collegePlayers)

	for _, p := range collegePlayers {
		res := structs.CollegePlayerResponse{
			ID:          p.PlayerID,
			BasePlayer:  p.BasePlayer,
			TeamID:      p.TeamID,
			TeamAbbr:    p.TeamAbbr,
			City:        p.City,
			State:       p.State,
			Year:        p.Year,
			IsRedshirt:  p.IsRedshirt,
			SeasonStats: p.SeasonStats,
		}
		responseList = append(responseList, res)
	}

	return responseList
}

func GetAllCollegePlayersWithStatsByTeamID(TeamID string, SeasonID string) []structs.CollegePlayer {
	db := dbprovider.GetInstance().GetDB()

	var collegePlayers []structs.CollegePlayer

	db.Preload("Stats", func(db *gorm.DB) *gorm.DB {
		return db.Where("season_id = ? and team_id = ? and snaps > 0", SeasonID, TeamID)
	}).Where("team_id = ?", TeamID).Find(&collegePlayers)

	return collegePlayers
}

func GetAllCollegePlayersWithGameStatsByTeamID(GameID string, stats []structs.CollegePlayerStats) []structs.GameResultsPlayer {
	db := dbprovider.GetInstance().GetDB()
	ids := []string{}
	statMap := make(map[uint]structs.CollegePlayerStats)
	for _, s := range stats {
		playerID := strconv.Itoa(s.CollegePlayerID)
		ids = append(ids, playerID)
		statMap[uint(s.CollegePlayerID)] = s
	}

	var collegePlayers []structs.CollegePlayer
	var matchRows []structs.GameResultsPlayer

	db.Where("id in (?)", ids).Find(&collegePlayers)

	for _, p := range collegePlayers {
		s := statMap[p.ID]
		if s.ID == 0 || s.Snaps == 0 {
			continue
		}

		row := structs.GameResultsPlayer{
			ID:                   p.ID,
			FirstName:            p.FirstName,
			LastName:             p.LastName,
			Position:             p.Position,
			Archetype:            p.Archetype,
			Year:                 uint(p.Year),
			TeamAbbr:             p.TeamAbbr,
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
	db.Where("id in (?)", ids).Find(&historicPlayers)

	for _, p := range historicPlayers {
		s := statMap[p.ID]
		if s.ID == 0 || s.Snaps == 0 {
			continue
		}

		row := structs.GameResultsPlayer{
			ID:                   p.ID,
			FirstName:            p.FirstName,
			LastName:             p.LastName,
			Position:             p.Position,
			Archetype:            p.Archetype,
			TeamAbbr:             p.TeamAbbr,
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

func GetAllNFLPlayersWithGameStatsByTeamID(GameID string, stats []structs.NFLPlayerStats) []structs.GameResultsPlayer {
	db := dbprovider.GetInstance().GetDB()

	ids := []string{}
	statMap := make(map[uint]structs.NFLPlayerStats)
	for _, s := range stats {
		playerID := strconv.Itoa(s.NFLPlayerID)
		ids = append(ids, playerID)
		statMap[uint(s.NFLPlayerID)] = s
	}

	var nflPlayers []structs.NFLPlayer
	var matchRows []structs.GameResultsPlayer

	db.Where("id in (?)", ids).Find(&nflPlayers)

	for _, p := range nflPlayers {
		s := statMap[p.ID]
		if s.ID == 0 || s.Snaps == 0 {
			continue
		}

		row := structs.GameResultsPlayer{
			ID:                   p.ID,
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
	db.Where("id in (?)", ids).Find(&historicPlayers)

	for _, p := range historicPlayers {
		s := statMap[p.ID]
		if s.ID == 0 || s.Snaps == 0 {
			continue
		}
		row := structs.GameResultsPlayer{
			ID:                   p.ID,
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

func GetNFLPlayersForRosterPage(TeamID string) []structs.NFLPlayer {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.NFLPlayer

	db.Preload("Contract", func(db *gorm.DB) *gorm.DB {
		return db.Where("is_active = true")
	}).Preload("SeasonStats").Preload("Extensions").Where("team_id = ?", TeamID).Find(&players)

	return players
}

func GetNFLPlayersRecordsByTeamID(TeamID string) []structs.NFLPlayer {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.NFLPlayer

	db.Where("team_id = ?", TeamID).Find(&players)

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

func CutCFBPlayer(playerId string) {
	db := dbprovider.GetInstance().GetDB()

	player := GetCollegePlayerByCollegePlayerId(playerId)
	player.WillTransfer()
	ts := GetTimestamp()
	if ts.IsOffSeason || ts.CollegeWeek <= 1 || ts.CollegeWeek >= 21 || ts.TransferPortalPhase == 3 {
		previousTeamID := strconv.Itoa(int(player.PreviousTeamID))
		deduction := 0
		promiseDeduction := 0
		if player.Stars > 2 {
			deduction = player.Stars / 2
		}
		collegePromise := GetCollegePromiseByCollegePlayerID(strconv.Itoa(int(player.ID)), previousTeamID)
		if collegePromise.IsActive && collegePromise.PromiseMade {
			weight := collegePromise.PromiseWeight
			switch weight {
			case "Vew Low":
				promiseDeduction = 3
			case "Low":
				promiseDeduction = 8
			case "Medium":
				promiseDeduction = 13
			case "High":
				promiseDeduction = 23
			case "Very High":
				promiseDeduction = 28
			}
		}

		points := (-1 * deduction) - promiseDeduction
		teamProfile := GetOnlyRecruitingProfileByTeamID(previousTeamID)
		teamProfile.IncrementClassSize()
		if player.Stars > 0 {
			teamProfile.AdjustPortalReputation(points)
			repository.SaveRecruitingTeamProfile(teamProfile, db)
		}
	}
	repository.SaveCFBPlayer(player, db)
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
		contract.CutContract()
	}

	capsheet.CutPlayerFromCapsheet(contract)
	repository.SaveNFLContract(contract, db)
	repository.SaveNFLPlayer(player, db)
	repository.SaveNFLCapsheet(capsheet, db)
}

func PlaceNFLPlayerOnPracticeSquad(playerId string) {
	db := dbprovider.GetInstance().GetDB()

	player := GetOnlyNFLPlayerRecord(playerId)
	player.ToggleIsPracticeSquad()

	if !player.IsPracticeSquad {
		Offers := GetFreeAgentOffersByPlayerID(strconv.Itoa(int(player.ID)))
		for _, o := range Offers {
			db.Delete(&o)
		}
	}

	db.Save(&player)
}

func PlaceNFLPlayerOnInjuryReserve(playerId string) {
	db := dbprovider.GetInstance().GetDB()

	player := GetOnlyNFLPlayerRecord(playerId)
	player.ToggleInjuryReserve()

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

	nflPlayers := GetAllNFLPlayersWithCurrentSeasonStats(seasonID, "2")

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

func GetAllPracticeSquadPlayersForFAPage() []models.FreeAgentResponse {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.NFLPlayer

	db.Preload("Offers", func(db *gorm.DB) *gorm.DB {
		return db.Order("contract_value DESC").Where("is_active = true")
	}).Where("is_practice_squad = ?", true).Find(&players)

	faResponseList := make([]models.FreeAgentResponse, len(players))

	for i, fa := range players {
		faResponseList[i] = models.FreeAgentResponse{
			ID:                fa.ID,
			PlayerID:          fa.PlayerID,
			TeamID:            fa.TeamID,
			College:           fa.College,
			TeamAbbr:          fa.TeamAbbr,
			FirstName:         fa.FirstName,
			LastName:          fa.LastName,
			Position:          fa.Position,
			Archetype:         fa.Archetype,
			Age:               fa.Age,
			Overall:           fa.Overall,
			PotentialGrade:    fa.PotentialGrade,
			FreeAgency:        fa.FreeAgency,
			Personality:       fa.Personality,
			RecruitingBias:    fa.RecruitingBias,
			WorkEthic:         fa.WorkEthic,
			AcademicBias:      fa.AcademicBias,
			PreviousTeam:      fa.PreviousTeam,
			PreviousTeamID:    fa.PreviousTeamID,
			Shotgun:           fa.Shotgun,
			Experience:        fa.Experience,
			Hometown:          fa.Hometown,
			State:             fa.State,
			IsActive:          fa.IsActive,
			IsWaived:          fa.IsWaived,
			IsPracticeSquad:   fa.IsPracticeSquad,
			IsFreeAgent:       fa.IsFreeAgent,
			IsAcceptingOffers: fa.IsAcceptingOffers,
			IsNegotiating:     fa.IsNegotiating,
			MinimumValue:      fa.MinimumValue,
			DraftedTeam:       fa.DraftedTeam,
			ShowLetterGrade:   fa.ShowLetterGrade,
			SeasonStats:       fa.SeasonStats,
			Offers:            fa.Offers,
		}
	}

	return faResponseList
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

func GetExtensionOffersByPlayerID(playerID string) []structs.NFLExtensionOffer {
	db := dbprovider.GetInstance().GetDB()

	offers := []structs.NFLExtensionOffer{}

	err := db.Where("nfl_player_id = ?", playerID).Find(&offers).Error
	if err != nil {
		return offers
	}

	return offers
}

func GetExtensionOfferByOfferID(OfferID string) structs.NFLExtensionOffer {
	db := dbprovider.GetInstance().GetDB()

	offer := structs.NFLExtensionOffer{}

	err := db.Where("id = ?", OfferID).Find(&offer).Error
	if err != nil {
		return offer
	}

	return offer
}

func CreateExtensionOffer(offer structs.FreeAgencyOfferDTO) structs.NFLExtensionOffer {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	extensionOffer := GetExtensionOfferByOfferID(strconv.Itoa(int(offer.ID)))
	player := GetNFLPlayerRecord(strconv.Itoa(int(offer.NFLPlayerID)))

	extensionOffer.CalculateOffer(offer)

	// If the owning team is sending an offer to a player
	if extensionOffer.ID == 0 {
		id := GetLatestExtensionOfferInDB(db)
		extensionOffer.AssignID(id)
		db.Create(&extensionOffer)
		fmt.Println("Creating Extension Offer!")

		message := offer.Team + " have offered a " + strconv.Itoa(offer.ContractLength) + " year contract extension for " + player.Position + " " + player.FirstName + " " + player.LastName + "."
		CreateNewsLog("NFL", message, "Free Agency", player.TeamID, ts)
	} else {
		fmt.Println("Updating Extension Offer!")
		db.Save(&extensionOffer)
	}

	return extensionOffer
}

func CancelExtensionOffer(offer structs.FreeAgencyOfferDTO) {
	db := dbprovider.GetInstance().GetDB()

	OfferID := strconv.Itoa(int(offer.ID))

	freeAgentOffer := GetExtensionOfferByOfferID(OfferID)

	freeAgentOffer.CancelOffer()

	db.Save(&freeAgentOffer)
}

func GetRetiredSimNFLPlayers() []structs.NFLRetiredPlayer {
	db := dbprovider.GetInstance().GetDB()

	players := []structs.NFLRetiredPlayer{}

	db.Find(&players)

	return players
}

func GetHistoricCollegePlayerByID(id string) structs.HistoricCollegePlayer {
	db := dbprovider.GetInstance().GetDB()

	var player structs.HistoricCollegePlayer

	err := db.Where("id = ?", id).Find(&player).Error
	if err != nil {
		fmt.Println("Could not find player in historics DB")
	}

	return player
}

func GetNFLDrafteeByID(id string) models.NFLDraftee {
	db := dbprovider.GetInstance().GetDB()

	var player models.NFLDraftee

	err := db.Where("id = ?", id).Find(&player).Error
	if err != nil {
		fmt.Println("Could not find player in historics DB")
	}

	return player
}

func GetTransferPortalPlayers() []structs.CollegePlayer {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.CollegePlayer

	db.Where("transfer_status > 0").Find(&players)

	return players
}

func GetCollegePlayerMap() map[uint]structs.CollegePlayer {

	portalMap := make(map[uint]structs.CollegePlayer)

	players := GetAllCollegePlayers()

	for _, p := range players {
		portalMap[p.ID] = p
	}

	return portalMap
}

func GetAllNFLPlayersWithSeasonStatsByTeamID(TeamID, SeasonID, gameType string) []structs.NFLPlayerResponse {
	db := dbprovider.GetInstance().GetDB()

	var nflPlayers []structs.NFLPlayer

	var responseList []structs.NFLPlayerResponse

	db.Preload("SeasonStats", func(db *gorm.DB) *gorm.DB {
		return db.Where("season_id = ? and team_id = ? and snaps > 0 AND game_type = ?", SeasonID, TeamID, gameType)
	}).Where("team_id = ?", TeamID).Find(&nflPlayers)

	for _, p := range nflPlayers {
		res := structs.NFLPlayerResponse{
			ID:          p.PlayerID,
			BasePlayer:  p.BasePlayer,
			TeamID:      p.TeamID,
			TeamAbbr:    p.TeamAbbr,
			State:       p.State,
			Year:        int(p.Experience),
			SeasonStats: p.SeasonStats,
		}
		responseList = append(responseList, res)
	}

	return responseList
}

func GetAllNFLPlayersMap() map[uint]structs.NFLPlayer {
	nflPlayers := GetAllNFLPlayers()

	return MapNFLPlayers(nflPlayers)
}

func GetHistoricCollegePlayersByTeamID(teamID string) []structs.HistoricCollegePlayer {
	db := dbprovider.GetInstance().GetDB()

	var CollegePlayers []structs.HistoricCollegePlayer

	db.Where("team_id = ?").Find(&CollegePlayers)

	return CollegePlayers
}
