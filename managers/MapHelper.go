package managers

import "github.com/CalebRose/SimFBA/structs"

func MapNFLPlayers(nflPlayers []structs.NFLPlayer) map[uint]structs.NFLPlayer {
	playerMap := make(map[uint]structs.NFLPlayer)

	for _, p := range nflPlayers {
		playerMap[p.ID] = p
	}

	return playerMap
}

func MakeCollegePlayerMap(players []structs.CollegePlayer) map[uint]structs.CollegePlayer {
	playerMap := make(map[uint]structs.CollegePlayer)

	for _, p := range players {
		playerMap[p.ID] = p
	}

	return playerMap
}

func MakeCollegePlayerMapByTeamID(players []structs.CollegePlayer, excludeUnsigned bool) map[uint][]structs.CollegePlayer {
	playerMap := make(map[uint][]structs.CollegePlayer)

	for _, p := range players {
		if p.TeamID == 0 && excludeUnsigned {
			continue
		}
		if len(playerMap[uint(p.TeamID)]) > 0 {
			playerMap[uint(p.TeamID)] = append(playerMap[uint(p.TeamID)], p)
		} else {
			playerMap[uint(p.TeamID)] = []structs.CollegePlayer{p}
		}
	}

	return playerMap
}

func MakeNFLPlayerMap(players []structs.NFLPlayer) map[uint]structs.NFLPlayer {
	playerMap := make(map[uint]structs.NFLPlayer)

	for _, p := range players {
		playerMap[p.ID] = p
	}

	return playerMap
}

func MakeNFLPlayerMapByTeamID(players []structs.NFLPlayer, excludeFAs bool) map[uint][]structs.NFLPlayer {
	playerMap := make(map[uint][]structs.NFLPlayer)

	for _, p := range players {
		if p.TeamID == 0 && excludeFAs {
			continue
		}
		if len(playerMap[uint(p.TeamID)]) > 0 {
			playerMap[uint(p.TeamID)] = append(playerMap[uint(p.TeamID)], p)
		} else {
			playerMap[uint(p.TeamID)] = []structs.NFLPlayer{p}
		}
	}

	return playerMap
}

func MakeCollegeDepthChartMap(dcs []structs.CollegeTeamDepthChart) map[uint]structs.CollegeTeamDepthChart {
	dcMap := make(map[uint]structs.CollegeTeamDepthChart)

	for _, dc := range dcs {
		dcMap[uint(dc.TeamID)] = dc
	}

	return dcMap
}

func MakeNFLDepthChartMap(dcs []structs.NFLDepthChart) map[uint]structs.NFLDepthChart {
	dcMap := make(map[uint]structs.NFLDepthChart)

	for _, dc := range dcs {
		dcMap[uint(dc.TeamID)] = dc
	}

	return dcMap
}

func MakeContractMap(contracts []structs.NFLContract) map[uint]structs.NFLContract {
	contractMap := make(map[uint]structs.NFLContract)

	for _, c := range contracts {
		contractMap[uint(c.NFLPlayerID)] = c
	}

	return contractMap
}

func MakeExtensionMap(extensions []structs.NFLExtensionOffer) map[uint]structs.NFLExtensionOffer {
	contractMap := make(map[uint]structs.NFLExtensionOffer)

	for _, c := range extensions {
		contractMap[uint(c.NFLPlayerID)] = c
	}

	return contractMap
}

func MakeHistoricCollegeStandingsMapByTeamID(standings []structs.CollegeStandings) map[uint][]structs.CollegeStandings {
	standingsMap := make(map[uint][]structs.CollegeStandings)

	for _, p := range standings {
		if p.TeamID == 0 {
			continue
		}
		if len(standingsMap[uint(p.TeamID)]) > 0 {
			standingsMap[uint(p.TeamID)] = append(standingsMap[uint(p.TeamID)], p)
		} else {
			standingsMap[uint(p.TeamID)] = []structs.CollegeStandings{p}
		}
	}

	return standingsMap
}

func MakeHistoricCollegeSeasonStatsMapByTeamID(stats []structs.CollegePlayerSeasonStats) map[uint][]structs.CollegePlayerSeasonStats {
	statsMap := make(map[uint][]structs.CollegePlayerSeasonStats)

	for _, p := range stats {
		if p.TeamID == 0 {
			continue
		}
		if len(statsMap[uint(p.TeamID)]) > 0 {
			statsMap[uint(p.TeamID)] = append(statsMap[uint(p.TeamID)], p)
		} else {
			statsMap[uint(p.TeamID)] = []structs.CollegePlayerSeasonStats{p}
		}
	}

	return statsMap
}

/*
Where("team_one_id = ? OR team_two_id = ?", teamID, teamID)
*/
func MakeHistoricRivalriesMapByTeamID(rivals []structs.CollegeRival) map[uint][]structs.CollegeRival {
	statsMap := make(map[uint][]structs.CollegeRival)

	for _, r := range rivals {
		if r.TeamOneID == 0 || r.TeamTwoID == 0 {
			continue
		}
		if len(statsMap[uint(r.TeamOneID)]) > 0 {
			statsMap[uint(r.TeamOneID)] = append(statsMap[uint(r.TeamOneID)], r)
		} else {
			statsMap[uint(r.TeamOneID)] = []structs.CollegeRival{r}
		}
		if len(statsMap[uint(r.TeamTwoID)]) > 0 {
			statsMap[uint(r.TeamTwoID)] = append(statsMap[uint(r.TeamTwoID)], r)
		} else {
			statsMap[uint(r.TeamTwoID)] = []structs.CollegeRival{r}
		}
	}

	return statsMap
}

func MakeHistoricGamesMapByTeamID(games []structs.CollegeGame) map[uint][]structs.CollegeGame {
	gamesMap := make(map[uint][]structs.CollegeGame)

	for _, r := range games {
		if r.HomeTeamID == 0 || r.AwayTeamID == 0 {
			continue
		}
		if len(gamesMap[uint(r.HomeTeamID)]) > 0 {
			gamesMap[uint(r.HomeTeamID)] = append(gamesMap[uint(r.HomeTeamID)], r)
		} else {
			gamesMap[uint(r.HomeTeamID)] = []structs.CollegeGame{r}
		}
		if len(gamesMap[uint(r.AwayTeamID)]) > 0 {
			gamesMap[uint(r.AwayTeamID)] = append(gamesMap[uint(r.AwayTeamID)], r)
		} else {
			gamesMap[uint(r.AwayTeamID)] = []structs.CollegeGame{r}
		}
	}

	return gamesMap
}
