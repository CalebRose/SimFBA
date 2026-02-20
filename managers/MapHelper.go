package managers

import (
	"github.com/CalebRose/SimFBA/models"
	"github.com/CalebRose/SimFBA/structs"
)

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

func MakeCollegeStandingsMapByTeamID(standings []structs.CollegeStandings) map[uint]structs.CollegeStandings {
	standingsMap := make(map[uint]structs.CollegeStandings)

	for _, p := range standings {
		if p.TeamID == 0 {
			continue
		}
		standingsMap[uint(p.TeamID)] = p
	}

	return standingsMap
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

func MakeFreeAgencyOfferMapByPlayer(offers []structs.FreeAgencyOffer) map[uint][]structs.FreeAgencyOffer {
	playerMap := make(map[uint][]structs.FreeAgencyOffer)

	for _, p := range offers {
		if p.TeamID == 0 || !p.IsActive {
			continue
		}
		if len(playerMap[uint(p.NFLPlayerID)]) > 0 {
			playerMap[uint(p.NFLPlayerID)] = append(playerMap[uint(p.NFLPlayerID)], p)
		} else {
			playerMap[uint(p.NFLPlayerID)] = []structs.FreeAgencyOffer{p}
		}
	}

	return playerMap
}

func MakePortalProfileMapByPlayerID(profiles []structs.TransferPortalProfile) map[uint][]structs.TransferPortalProfile {
	playerMap := make(map[uint][]structs.TransferPortalProfile)

	for _, p := range profiles {
		if p.ProfileID == 0 {
			continue
		}
		if len(playerMap[uint(p.CollegePlayerID)]) > 0 {
			playerMap[uint(p.CollegePlayerID)] = append(playerMap[uint(p.CollegePlayerID)], p)
		} else {
			playerMap[uint(p.CollegePlayerID)] = []structs.TransferPortalProfile{p}
		}
	}

	return playerMap
}

func MakePortalProfileMapByTeamID(profiles []structs.TransferPortalProfile) map[uint][]structs.TransferPortalProfile {
	playerMap := make(map[uint][]structs.TransferPortalProfile)

	for _, p := range profiles {
		if p.ProfileID == 0 {
			continue
		}
		if len(playerMap[uint(p.ProfileID)]) > 0 {
			playerMap[uint(p.ProfileID)] = append(playerMap[uint(p.ProfileID)], p)
		} else {
			playerMap[uint(p.ProfileID)] = []structs.TransferPortalProfile{p}
		}
	}

	return playerMap
}

func MakeCollegePromiseMap(promises []structs.CollegePromise) map[uint]structs.CollegePromise {
	playerMap := make(map[uint]structs.CollegePromise)

	for _, p := range promises {
		playerMap[p.ID] = p
	}

	return playerMap
}

func MakePromiseMapByTeamID(profiles []structs.CollegePromise) map[uint][]structs.CollegePromise {
	playerMap := make(map[uint][]structs.CollegePromise)

	for _, p := range profiles {
		if p.TeamID == 0 || !p.IsActive {
			continue
		}
		if len(playerMap[uint(p.TeamID)]) > 0 {
			playerMap[uint(p.TeamID)] = append(playerMap[uint(p.TeamID)], p)
		} else {
			playerMap[uint(p.TeamID)] = []structs.CollegePromise{p}
		}
	}

	return playerMap
}

func MakeNFLWarRoomMap(warRooms []models.NFLWarRoom) map[uint]models.NFLWarRoom {
	warRoomMap := make(map[uint]models.NFLWarRoom)

	for _, t := range warRooms {
		warRoomMap[t.TeamID] = t
	}

	return warRoomMap
}

func MakeScoutingProfileMapByTeam(profiles []models.ScoutingProfile) map[uint][]models.ScoutingProfile {
	profileMap := make(map[uint][]models.ScoutingProfile)

	for _, t := range profiles {
		if len(profileMap[t.TeamID]) > 0 {
			profileMap[t.TeamID] = append(profileMap[t.TeamID], t)
		} else {
			profileMap[t.TeamID] = []models.ScoutingProfile{t}
		}
	}

	return profileMap
}

func MakePromiseMapByPlayerIDByTeam(promises []structs.CollegePromise) map[uint]structs.CollegePromise {
	playerMap := make(map[uint]structs.CollegePromise)

	for _, p := range promises {
		playerMap[p.CollegePlayerID] = p
	}

	return playerMap
}

func MakeCollegeGameMapByID(collegeGames []structs.CollegeGame) map[uint]structs.CollegeGame {
	gamesMap := make(map[uint]structs.CollegeGame)

	for _, c := range collegeGames {
		gamesMap[uint(c.ID)] = c
	}

	return gamesMap
}

func MakeCollegeGameMapByTeamID(collegeGames []structs.CollegeGame) map[uint][]structs.CollegeGame {
	gamesMap := make(map[uint][]structs.CollegeGame)

	for _, c := range collegeGames {
		if c.IsSpringGame || c.IsPlayoffGame || c.IsBowlGame {
			continue
		}
		if len(gamesMap[uint(c.HomeTeamID)]) > 0 {
			gamesMap[uint(c.HomeTeamID)] = append(gamesMap[uint(c.HomeTeamID)], c)
		} else {
			gamesMap[uint(c.HomeTeamID)] = []structs.CollegeGame{c}
		}
		if len(gamesMap[uint(c.AwayTeamID)]) > 0 {
			gamesMap[uint(c.AwayTeamID)] = append(gamesMap[uint(c.AwayTeamID)], c)
		} else {
			gamesMap[uint(c.AwayTeamID)] = []structs.CollegeGame{c}
		}
	}
	return gamesMap
}

func MakeCollegeGameplanMap(gameplans []structs.CollegeGameplan) map[uint]structs.CollegeGameplan {
	gamesMap := make(map[uint]structs.CollegeGameplan)

	for _, c := range gameplans {
		gamesMap[uint(c.ID)] = c
	}

	return gamesMap
}

func MakeNFLGameplanMap(gameplans []structs.NFLGameplan) map[uint]structs.NFLGameplan {
	gamesMap := make(map[uint]structs.NFLGameplan)

	for _, c := range gameplans {
		gamesMap[uint(c.ID)] = c
	}

	return gamesMap
}

func MakeNFLTradePreferencesMap(tradePreferences []structs.NFLTradePreferences) map[uint]structs.NFLTradePreferences {
	preferencesMap := make(map[uint]structs.NFLTradePreferences)

	for _, c := range tradePreferences {
		preferencesMap[uint(c.NFLTeamID)] = c
	}

	return preferencesMap
}

func MakeRecruitProfileMapByRecruitID(profiles []structs.RecruitPlayerProfile) map[uint][]structs.RecruitPlayerProfile {
	profileMap := make(map[uint][]structs.RecruitPlayerProfile)

	for _, rp := range profiles {
		if len(profileMap[uint(rp.RecruitID)]) > 0 {
			profileMap[uint(rp.RecruitID)] = append(profileMap[uint(rp.RecruitID)], rp)
		} else {
			profileMap[uint(rp.RecruitID)] = []structs.RecruitPlayerProfile{rp}
		}
	}

	return profileMap
}

func MakeTradeProposalMap(proposals []structs.NFLTradeProposal) map[uint][]structs.NFLTradeProposal {
	proposalMap := make(map[uint][]structs.NFLTradeProposal)

	for _, proposal := range proposals {
		if len(proposalMap[proposal.NFLTeamID]) > 0 {
			proposalMap[proposal.NFLTeamID] = append(proposalMap[uint(proposal.NFLTeamID)], proposal)
		} else {
			proposalMap[proposal.NFLTeamID] = []structs.NFLTradeProposal{proposal}
		}
	}

	return proposalMap
}

func MakeRecruitTeamProfileMapByTeamID(profiles []structs.RecruitingTeamProfile) map[uint]structs.RecruitingTeamProfile {
	profileMap := make(map[uint]structs.RecruitingTeamProfile)

	for _, rp := range profiles {
		profileMap[rp.ID] = rp
	}

	return profileMap
}
