package util

import (
	"strconv"

	"github.com/CalebRose/SimFBA/structs"
)

func GetHeismanScore(cp structs.CollegePlayer, weightMap map[string]float64, homeTeamMapper map[int]string, games []structs.CollegeGame) float64 {
	var score float64 = 0

	homeTeam := homeTeamMapper[cp.TeamID]

	stats := cp.Stats

	var totalMod float64 = 0

	for idx, stat := range stats {
		if idx > 12 {
			continue
		}
		var statScore float64 = 0

		opposingTeam := stat.OpposingTeam
		statScore += (float64(stat.PassingYards) * 0.069)
		statScore += (float64(stat.PassingTDs) * 4)
		statScore -= (float64(stat.Interceptions) * 2.25)
		statScore -= (float64(stat.Sacks) * 2.25)
		if cp.Position == "RB" || cp.Position == "FB" {
			statScore += (float64(stat.RushingYards) * 0.1)
			statScore += (float64(stat.RushingTDs) * 6)
		} else {
			statScore += (float64(stat.RushingYards) * 0.0775)
			statScore += (float64(stat.RushingTDs) * 4.75)
		}

		drops := stat.Targets - stat.Catches
		if cp.Position == "WR" || cp.Position == "TE" {
			statScore += (float64(stat.Catches) * 0.525)
			statScore += (float64(stat.ReceivingYards) * 0.1125)
			statScore += (float64(stat.ReceivingTDs) * 6)
			statScore -= float64(drops) * 0.75
		} else {
			statScore += (float64(stat.Catches) * 0.25)
			statScore += (float64(stat.ReceivingYards) * 0.05)
			statScore += (float64(stat.ReceivingTDs) * 4)
			statScore -= float64(drops) * 0.75
		}
		statScore -= (float64(stat.Fumbles) * 6)

		statScore += (float64(stat.SoloTackles) * 1)
		statScore += (float64(stat.STSoloTackles) * 1)
		statScore += (float64(stat.AssistedTackles) * .9)
		statScore += (float64(stat.STAssistedTackles) * .9)

		statScore += (float64(stat.TacklesForLoss) * 6.25)
		statScore += (float64(stat.SacksMade) * 7.125)
		statScore += (float64(stat.PassDeflections) * 6.25)

		statScore += (float64(stat.ForcedFumbles) * 8)
		statScore += (float64(stat.RecoveredFumbles) * 6)
		statScore += (float64(stat.InterceptionsCaught) * 15)
		statScore += (float64(stat.PuntsBlocked) * 10)

		statScore += (float64(stat.Safeties) * 10)
		statScore += (float64(stat.DefensiveTDs) * 20)
		statScore += (float64(stat.KickReturnTDs) * 14)
		statScore += (float64(stat.PuntReturnTDs) * 14)
		statScore += (float64(stat.FGBlocked) * 12)
		statScore += (float64(stat.FGMade) * 0.3)
		statScore += (float64(stat.ExtraPointsMade) * 0.1)

		game := GetCollegeGameByGameID(games, stat.GameID)

		opposingTeamWeight := weightMap[opposingTeam]

		if (game.HomeTeamWin && cp.TeamID != game.HomeTeamID) || (game.AwayTeamWin && cp.TeamID != game.AwayTeamID) {
			opposingTeamWeight *= -.4125
		}

		totalMod += opposingTeamWeight

		// statScore *= weightMap[opposingTeam]

		// statScore *= weightMap[homeTeam]

		score += statScore
	}

	score = (score / float64(len(stats))) * (float64(len(stats)) / 12)

	score = score * (1 + totalMod + weightMap[homeTeam])

	return score
}

func GetCollegeGameByGameID(games []structs.CollegeGame, gameID int) structs.CollegeGame {
	for _, game := range games {
		if int(game.ID) == gameID {
			return game
		}
	}

	return structs.CollegeGame{}
}

func GetCollegePlayerIDsBySeasonStats(cps []structs.CollegePlayerSeasonStats) []string {
	var list []string

	for _, cp := range cps {
		list = append(list, strconv.Itoa(int(cp.CollegePlayerID)))
	}

	return list
}

func GetNFLPlayerIDsBySeasonStats(cps []structs.NFLPlayerSeasonStats) []string {
	var list []string

	for _, cp := range cps {
		list = append(list, strconv.Itoa(int(cp.NFLPlayerID)))
	}

	return list
}

func GetCollegePlayerIDs(cps []structs.CollegePlayerStats) []string {
	var list []string

	for _, cp := range cps {
		list = append(list, strconv.Itoa(int(cp.CollegePlayerID)))
	}

	return list
}
