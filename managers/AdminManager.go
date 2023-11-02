package managers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/structs"
)

// Timestamp Funcs
// GetTimestamp -- Get the Timestamp
func GetTimestamp() structs.Timestamp {
	db := dbprovider.GetInstance().GetDB()

	var timestamp structs.Timestamp

	db.First(&timestamp)

	return timestamp
}

func GetCollegeWeek(weekID string, ts structs.Timestamp) structs.CollegeWeek {
	db := dbprovider.GetInstance().GetDB()

	var week structs.CollegeWeek

	db.Where("week = ? AND season_id = ?", weekID, ts.CollegeSeasonID).Find(&week)

	return week
}

func MoveUpWeek() structs.Timestamp {
	db := dbprovider.GetInstance().GetDB()
	timestamp := GetTimestamp()

	// Sync to Next Week
	UpdateGameplanPenalties()
	RecoverPlayers()
	CheckNFLRookiesForLetterGrade(strconv.Itoa(int(timestamp.NFLSeasonID)))
	timestamp.SyncToNextWeek()
	if timestamp.NFLWeek > 15 {
		SyncExtensionOffers()
	}
	if timestamp.NFLWeek > 22 {
		timestamp.MoveUpSeason()
		// Run Progressions
		//
	}
	db.Save(&timestamp)

	return timestamp
}

func RegressTimeslot(timeslot string) {
	db := dbprovider.GetInstance().GetDB()
	timestamp := GetTimestamp()

	// Update timeslot
	// timestamp.ToggleTimeSlot(timeslot)

	isCFB := false
	if timeslot == "Thursday Night" ||
		timeslot == "Friday Night" ||
		timeslot == "Saturday Morning" ||
		timeslot == "Saturday Afternoon" ||
		timeslot == "Saturday Evening" ||
		timeslot == "Saturday Night" {
		isCFB = true
	}

	if isCFB {
		// Get Games
		games := GetCollegeGamesByTimeslotAndWeekId(strconv.Itoa(timestamp.CollegeWeekID), timeslot)
		seasonStats := GetCollegeSeasonStatsBySeason(strconv.Itoa(timestamp.CollegeSeasonID))
		seasonStatsMap := make(map[int]*structs.CollegeTeamSeasonStats)
		for _, s := range seasonStats {
			seasonStatsMap[int(s.TeamID)] = &s
		}

		for _, game := range games {
			if game.ID == 2048 || game.ID == 2049 {
				continue
			}
			// Get team stats
			gameID := strconv.Itoa(int(game.ID))
			homeTeamID := game.HomeTeamID
			awayTeamID := game.AwayTeamID

			// homeTeamSeasonStats := seasonStatsMap[homeTeamID]
			// awayTeamSeasonStats := seasonStatsMap[awayTeamID]
			homeTeamSeasonStats := GetCollegeTeamSeasonStatsBySeason(strconv.Itoa(homeTeamID), strconv.Itoa(int(timestamp.CollegeSeasonID)))
			awayTeamSeasonStats := GetCollegeTeamSeasonStatsBySeason(strconv.Itoa(awayTeamID), strconv.Itoa(int(timestamp.CollegeSeasonID)))
			if homeTeamSeasonStats.ID == 0 {
				homeTeamSeasonStats = structs.CollegeTeamSeasonStats{
					TeamID:   uint(homeTeamID),
					SeasonID: uint(game.SeasonID),
					Year:     timestamp.Season,
				}
			}

			if awayTeamSeasonStats.ID == 0 {
				awayTeamSeasonStats = structs.CollegeTeamSeasonStats{
					TeamID:   uint(awayTeamID),
					SeasonID: uint(game.SeasonID),
					Year:     timestamp.Season,
				}
			}

			homeTeamStats := GetCollegeTeamStatsByGame(strconv.Itoa(homeTeamID), gameID)
			awayTeamStats := GetCollegeTeamStatsByGame(strconv.Itoa(awayTeamID), gameID)

			homeTeamSeasonStats.ReduceStats([]structs.CollegeTeamStats{homeTeamStats})
			awayTeamSeasonStats.ReduceStats([]structs.CollegeTeamStats{awayTeamStats})
			// Get Player Stats
			homePlayerStats := GetAllCollegePlayerStatsByGame(gameID, strconv.Itoa(homeTeamID))
			awayPlayerStats := GetAllCollegePlayerStatsByGame(gameID, strconv.Itoa(awayTeamID))

			for _, h := range homePlayerStats {
				if h.Snaps == 0 {
					continue
				}
				// playerSeasonStat := playerSeasonStatsMap[h.CollegePlayerID]
				playerSeasonStat := GetCollegeSeasonStatsByPlayerAndSeason(strconv.Itoa(h.CollegePlayerID), strconv.Itoa(int(timestamp.CollegeSeasonID)))
				if playerSeasonStat.ID == 0 {
					playerSeasonStat = structs.CollegePlayerSeasonStats{
						CollegePlayerID: uint(h.CollegePlayerID),
						SeasonID:        uint(timestamp.CollegeSeasonID),
						TeamID:          uint(h.TeamID),
						Year:            uint(h.Year),
					}
				}

				playerSeasonStat.ReduceStats([]structs.CollegePlayerStats{h})

				if !timestamp.CFBSpringGames {
					db.Save(&playerSeasonStat)
				}
			}

			for _, a := range awayPlayerStats {
				if a.Snaps == 0 {
					continue
				}
				// playerSeasonStat := playerSeasonStatsMap[a.CollegePlayerID]
				playerSeasonStat := GetCollegeSeasonStatsByPlayerAndSeason(strconv.Itoa(a.CollegePlayerID), strconv.Itoa(int(timestamp.CollegeSeasonID)))
				if playerSeasonStat.ID == 0 {
					playerSeasonStat = structs.CollegePlayerSeasonStats{
						CollegePlayerID: uint(a.CollegePlayerID),
						SeasonID:        uint(timestamp.CollegeSeasonID),
						TeamID:          uint(a.TeamID),
						Year:            uint(a.Year),
					}
				}
				playerSeasonStat.ReduceStats([]structs.CollegePlayerStats{a})

				if !timestamp.CFBSpringGames {
					db.Save(&playerSeasonStat)
				}
			}

			// Update Standings
			homeTeamStandings := GetCFBStandingsByTeamIDAndSeasonID(strconv.Itoa(homeTeamID), strconv.Itoa(timestamp.CollegeSeasonID))
			awayTeamStandings := GetCFBStandingsByTeamIDAndSeasonID(strconv.Itoa(awayTeamID), strconv.Itoa(timestamp.CollegeSeasonID))

			homeTeamStandings.SubtractCollegeStandings(game)
			awayTeamStandings.SubtractCollegeStandings(game)

			if game.HomeTeamCoach != "AI" && !timestamp.CFBSpringGames {
				homeCoach := GetCollegeCoachByCoachName(game.HomeTeamCoach)
				homeCoach.UpdateCoachRecord(game)

				err := db.Save(&homeCoach).Error
				if err != nil {
					log.Panicln("Could not save coach record for team " + strconv.Itoa(homeTeamID))
				}
			}

			if game.AwayTeamCoach != "AI" && !timestamp.CFBSpringGames {
				awayCoach := GetCollegeCoachByCoachName(game.AwayTeamCoach)
				awayCoach.UpdateCoachRecord(game)
				err := db.Save(&awayCoach).Error
				if err != nil {
					log.Panicln("Could not save coach record for team " + strconv.Itoa(awayTeamID))
				}
			}

			// Save
			if !timestamp.CFBSpringGames {
				db.Save(&homeTeamSeasonStats)
				db.Save(&awayTeamSeasonStats)
				db.Save(&homeTeamStandings)
				db.Save(&awayTeamStandings)
			}
		}
	} else {
		// Get Games
		games := GetNFLGamesByTimeslotAndWeekId(strconv.Itoa(timestamp.NFLWeekID), timeslot)

		for _, game := range games {
			// Get team stats
			gameID := strconv.Itoa(int(game.ID))
			homeTeamID := game.HomeTeamID
			awayTeamID := game.AwayTeamID

			// homeTeamSeasonStats := seasonStatsMap[homeTeamID]
			// awayTeamSeasonStats := seasonStatsMap[awayTeamID]
			homeTeamSeasonStats := GetNFLTeamSeasonStatsByTeamANDSeason(strconv.Itoa(homeTeamID), strconv.Itoa(int(timestamp.NFLSeasonID)))
			awayTeamSeasonStats := GetNFLTeamSeasonStatsByTeamANDSeason(strconv.Itoa(awayTeamID), strconv.Itoa(int(timestamp.NFLSeasonID)))

			homeTeamStats := GetNFLTeamStatsByGame(strconv.Itoa(homeTeamID), gameID)
			awayTeamStats := GetNFLTeamStatsByGame(strconv.Itoa(awayTeamID), gameID)

			homeTeamSeasonStats.SubtractStats([]structs.NFLTeamStats{homeTeamStats})
			awayTeamSeasonStats.SubtractStats([]structs.NFLTeamStats{awayTeamStats})
			// Get Player Stats
			homePlayerStats := GetAllNFLPlayerStatsByGame(gameID, strconv.Itoa(homeTeamID))
			awayPlayerStats := GetAllNFLPlayerStatsByGame(gameID, strconv.Itoa(awayTeamID))

			for _, h := range homePlayerStats {
				if h.Snaps == 0 {
					continue
				}
				// playerSeasonStat := playerSeasonStatsMap[h.NFLPlayerID]
				seasonStats := GetNFLSeasonStatsByPlayerAndSeason(strconv.Itoa(h.NFLPlayerID), strconv.Itoa(int(timestamp.NFLSeasonID)))
				if seasonStats.ID == 0 {
					seasonStats = structs.NFLPlayerSeasonStats{
						NFLPlayerID: uint(h.NFLPlayerID),
						SeasonID:    uint(timestamp.NFLSeasonID),
						TeamID:      uint(h.TeamID),
						Team:        h.Team,
						Year:        uint(timestamp.Season),
					}
				}

				seasonStats.SubtractStats([]structs.NFLPlayerStats{h}, timestamp)

				db.Save(&seasonStats)
			}

			for _, a := range awayPlayerStats {
				if a.Snaps == 0 {
					continue
				}

				// playerSeasonStat := playerSeasonStatsMap[a.NFLPlayerID]
				seasonStats := GetNFLSeasonStatsByPlayerAndSeason(strconv.Itoa(a.NFLPlayerID), strconv.Itoa(int(timestamp.NFLSeasonID)))
				if seasonStats.ID == 0 {
					seasonStats = structs.NFLPlayerSeasonStats{
						NFLPlayerID: uint(a.NFLPlayerID),
						SeasonID:    uint(timestamp.NFLSeasonID),
						TeamID:      uint(a.TeamID),
						Team:        a.Team,
						Year:        uint(timestamp.Season),
					}
				}

				seasonStats.SubtractStats([]structs.NFLPlayerStats{a}, timestamp)

				db.Save(&seasonStats)
			}

			// Update Standings
			homeTeamStandings := GetNFLStandingsByTeamIDAndSeasonID(strconv.Itoa(homeTeamID), strconv.Itoa(timestamp.NFLSeasonID))
			awayTeamStandings := GetNFLStandingsByTeamIDAndSeasonID(strconv.Itoa(awayTeamID), strconv.Itoa(timestamp.NFLSeasonID))

			homeTeamStandings.ReduceNFLStandings(game)
			awayTeamStandings.ReduceNFLStandings(game)

			if game.HomeTeamCoach != "AI" && !timestamp.NFLPreseason {
				homeCoach := GetNFLUserByUsername(game.HomeTeamCoach)
				homeCoach.ReduceCoachRecord(game)

				err := db.Save(&homeCoach).Error
				if err != nil {
					log.Panicln("Could not save coach record for team " + strconv.Itoa(homeTeamID))
				}
			}

			if game.AwayTeamCoach != "AI" && !timestamp.NFLPreseason {
				awayCoach := GetNFLUserByUsername(game.AwayTeamCoach)
				awayCoach.ReduceCoachRecord(game)
				err := db.Save(&awayCoach).Error
				if err != nil {
					log.Panicln("Could not save coach record for team " + strconv.Itoa(awayTeamID))
				}
			}

			// Save
			if !timestamp.NFLPreseason {
				db.Save(&homeTeamSeasonStats)
				db.Save(&awayTeamSeasonStats)
				db.Save(&homeTeamStandings)
				db.Save(&awayTeamStandings)
			}

		}
	}

	db.Save(&timestamp)
}

func SyncTimeslot(timeslot string) {
	db := dbprovider.GetInstance().GetDB()
	timestamp := GetTimestamp()

	// Update timeslot
	timestamp.ToggleTimeSlot(timeslot)

	isCFB := false
	if timeslot == "Thursday Night" ||
		timeslot == "Friday Night" ||
		timeslot == "Saturday Morning" ||
		timeslot == "Saturday Afternoon" ||
		timeslot == "Saturday Evening" ||
		timeslot == "Saturday Night" {
		isCFB = true
	}

	if isCFB {
		// Get Games
		games := GetCollegeGamesByTimeslotAndWeekId(strconv.Itoa(timestamp.CollegeWeekID), timeslot)
		seasonStats := GetCollegeSeasonStatsBySeason(strconv.Itoa(timestamp.CollegeSeasonID))
		seasonStatsMap := make(map[int]*structs.CollegeTeamSeasonStats)
		for _, s := range seasonStats {
			seasonStatsMap[int(s.TeamID)] = &s
		}

		for _, game := range games {
			// Get team stats
			gameID := strconv.Itoa(int(game.ID))
			homeTeamID := game.HomeTeamID
			awayTeamID := game.AwayTeamID

			// homeTeamSeasonStats := seasonStatsMap[homeTeamID]
			// awayTeamSeasonStats := seasonStatsMap[awayTeamID]
			homeTeamSeasonStats := GetCollegeTeamSeasonStatsBySeason(strconv.Itoa(homeTeamID), strconv.Itoa(int(timestamp.CollegeSeasonID)))
			awayTeamSeasonStats := GetCollegeTeamSeasonStatsBySeason(strconv.Itoa(awayTeamID), strconv.Itoa(int(timestamp.CollegeSeasonID)))
			if homeTeamSeasonStats.ID == 0 {
				homeTeamSeasonStats = structs.CollegeTeamSeasonStats{
					TeamID:   uint(homeTeamID),
					SeasonID: uint(game.SeasonID),
					Year:     timestamp.Season,
				}
			}

			if awayTeamSeasonStats.ID == 0 {
				awayTeamSeasonStats = structs.CollegeTeamSeasonStats{
					TeamID:   uint(awayTeamID),
					SeasonID: uint(game.SeasonID),
					Year:     timestamp.Season,
				}
			}

			homeTeamStats := GetCollegeTeamStatsByGame(strconv.Itoa(homeTeamID), gameID)
			awayTeamStats := GetCollegeTeamStatsByGame(strconv.Itoa(awayTeamID), gameID)

			homeTeamSeasonStats.MapStats([]structs.CollegeTeamStats{homeTeamStats})
			awayTeamSeasonStats.MapStats([]structs.CollegeTeamStats{awayTeamStats})
			// Get Player Stats
			homePlayerStats := GetAllCollegePlayerStatsByGame(gameID, strconv.Itoa(homeTeamID))
			awayPlayerStats := GetAllCollegePlayerStatsByGame(gameID, strconv.Itoa(awayTeamID))

			for _, h := range homePlayerStats {
				if h.Snaps == 0 {
					continue
				}
				if h.WasInjured {
					playerRecord := GetCollegePlayerByCollegePlayerId(strconv.Itoa(h.CollegePlayerID))
					playerRecord.SetIsInjured(h.WasInjured, h.InjuryType, h.WeeksOfRecovery)
					db.Save(&playerRecord)
				}
				// playerSeasonStat := playerSeasonStatsMap[h.CollegePlayerID]
				playerSeasonStat := GetCollegeSeasonStatsByPlayerAndSeason(strconv.Itoa(h.CollegePlayerID), strconv.Itoa(int(timestamp.CollegeSeasonID)))
				if playerSeasonStat.ID == 0 {
					playerSeasonStat = structs.CollegePlayerSeasonStats{
						CollegePlayerID: uint(h.CollegePlayerID),
						SeasonID:        uint(timestamp.CollegeSeasonID),
						TeamID:          uint(h.TeamID),
						Year:            uint(h.Year),
					}
				}

				playerSeasonStat.MapStats([]structs.CollegePlayerStats{h})

				if !timestamp.CFBSpringGames {
					db.Save(&playerSeasonStat)
				}
			}

			for _, a := range awayPlayerStats {
				if a.Snaps == 0 {
					continue
				}
				if a.WasInjured {
					playerRecord := GetCollegePlayerByCollegePlayerId(strconv.Itoa(a.CollegePlayerID))
					playerRecord.SetIsInjured(a.WasInjured, a.InjuryType, a.WeeksOfRecovery)
					db.Save(&playerRecord)
				}
				// playerSeasonStat := playerSeasonStatsMap[a.CollegePlayerID]
				playerSeasonStat := GetCollegeSeasonStatsByPlayerAndSeason(strconv.Itoa(a.CollegePlayerID), strconv.Itoa(int(timestamp.CollegeSeasonID)))
				if playerSeasonStat.ID == 0 {
					playerSeasonStat = structs.CollegePlayerSeasonStats{
						CollegePlayerID: uint(a.CollegePlayerID),
						SeasonID:        uint(timestamp.CollegeSeasonID),
						TeamID:          uint(a.TeamID),
						Year:            uint(a.Year),
					}
				}
				playerSeasonStat.MapStats([]structs.CollegePlayerStats{a})

				if !timestamp.CFBSpringGames {
					db.Save(&playerSeasonStat)
				}
			}

			// Update Standings
			homeTeamStandings := GetCFBStandingsByTeamIDAndSeasonID(strconv.Itoa(homeTeamID), strconv.Itoa(timestamp.CollegeSeasonID))
			awayTeamStandings := GetCFBStandingsByTeamIDAndSeasonID(strconv.Itoa(awayTeamID), strconv.Itoa(timestamp.CollegeSeasonID))

			homeTeamStandings.UpdateCollegeStandings(game)
			awayTeamStandings.UpdateCollegeStandings(game)

			if game.HomeTeamCoach != "AI" && !timestamp.CFBSpringGames {
				homeCoach := GetCollegeCoachByCoachName(game.HomeTeamCoach)
				homeCoach.UpdateCoachRecord(game)

				err := db.Save(&homeCoach).Error
				if err != nil {
					log.Panicln("Could not save coach record for team " + strconv.Itoa(homeTeamID))
				}
			}

			if game.AwayTeamCoach != "AI" && !timestamp.CFBSpringGames {
				awayCoach := GetCollegeCoachByCoachName(game.AwayTeamCoach)
				awayCoach.UpdateCoachRecord(game)
				err := db.Save(&awayCoach).Error
				if err != nil {
					log.Panicln("Could not save coach record for team " + strconv.Itoa(awayTeamID))
				}
			}

			// Save
			if !timestamp.CFBSpringGames {
				if game.NextGameID > 0 {
					nextGameID := strconv.Itoa(int(game.NextGameID))
					winningTeamID := 0
					winningTeam := ""
					winningCoach := ""
					if game.HomeTeamWin {
						winningTeamID = game.HomeTeamID
						winningTeam = game.HomeTeam
						winningCoach = game.HomeTeamCoach
					} else {
						winningTeamID = game.AwayTeamID
						winningTeam = game.AwayTeam
						winningCoach = game.AwayTeamCoach
					}

					nextGame := GetCollegeGameByGameID(nextGameID)
					nextGame.AddTeam(game.NextGameHOA == "H", winningTeamID, winningTeam, winningCoach)
					if !nextGame.IsNeutral && !game.IsNationalChampionship {
						stadiumID := 0
						stadium := ""
						city := ""
						state := ""
						isDomed := false

						if game.HomeTeamWin {
							stadium = game.Stadium
							city = game.City
							state = game.State
							stadiumID = int(game.StadiumID)
						} else {
							awayTeam := GetTeamByTeamID(strconv.Itoa(awayTeamID))
							stadiumRecord := GetStadiumByStadiumID(strconv.Itoa(int(awayTeam.StadiumID)))
							stadium = awayTeam.Stadium
							city = awayTeam.City
							state = awayTeam.State
							stadiumID = int(awayTeam.StadiumID)
							isDomed = stadiumRecord.IsDomed
						}
						nextGame.AddLocation(stadiumID, stadium, city, state, isDomed)
					}

					// Updating matchup for playoff game!
					db.Save(&nextGame)
				}

				db.Save(&homeTeamSeasonStats)
				db.Save(&awayTeamSeasonStats)
				db.Save(&homeTeamStandings)
				db.Save(&awayTeamStandings)
			}
		}
	} else {
		// Get Games
		games := GetNFLGamesByTimeslotAndWeekId(strconv.Itoa(timestamp.NFLWeekID), timeslot)

		// seasonStatsMap := make(map[int]structs.NFLTeamSeasonStats)
		// for _, s := range seasonStats {
		// 	seasonStatsMap[int(s.TeamID)] = s
		// }

		for _, game := range games {
			// Get team stats
			gameID := strconv.Itoa(int(game.ID))
			homeTeamID := game.HomeTeamID
			awayTeamID := game.AwayTeamID

			// homeTeamSeasonStats := seasonStatsMap[homeTeamID]
			// awayTeamSeasonStats := seasonStatsMap[awayTeamID]
			homeTeamSeasonStats := GetNFLTeamSeasonStatsByTeamANDSeason(strconv.Itoa(homeTeamID), strconv.Itoa(int(timestamp.NFLSeasonID)))
			awayTeamSeasonStats := GetNFLTeamSeasonStatsByTeamANDSeason(strconv.Itoa(awayTeamID), strconv.Itoa(int(timestamp.NFLSeasonID)))

			homeTeamStats := GetNFLTeamStatsByGame(strconv.Itoa(homeTeamID), gameID)
			awayTeamStats := GetNFLTeamStatsByGame(strconv.Itoa(awayTeamID), gameID)

			homeTeamSeasonStats.MapStats([]structs.NFLTeamStats{homeTeamStats})
			awayTeamSeasonStats.MapStats([]structs.NFLTeamStats{awayTeamStats})
			// Get Player Stats
			homePlayerStats := GetAllNFLPlayerStatsByGame(gameID, strconv.Itoa(homeTeamID))
			awayPlayerStats := GetAllNFLPlayerStatsByGame(gameID, strconv.Itoa(awayTeamID))

			for _, h := range homePlayerStats {
				if h.Snaps == 0 {
					continue
				}
				if h.WasInjured {
					playerRecord := GetNFLPlayerRecord(strconv.Itoa(h.NFLPlayerID))
					playerRecord.SetIsInjured(h.WasInjured, h.InjuryType, h.WeeksOfRecovery)
					db.Save(&playerRecord)
				}
				// playerSeasonStat := playerSeasonStatsMap[h.NFLPlayerID]
				seasonStats := GetNFLSeasonStatsByPlayerAndSeason(strconv.Itoa(h.NFLPlayerID), strconv.Itoa(int(timestamp.NFLSeasonID)))
				if seasonStats.ID == 0 {
					seasonStats = structs.NFLPlayerSeasonStats{
						NFLPlayerID: uint(h.NFLPlayerID),
						SeasonID:    uint(timestamp.NFLSeasonID),
						TeamID:      uint(h.TeamID),
						Team:        h.Team,
						Year:        uint(timestamp.Season),
					}
				}

				seasonStats.MapStats([]structs.NFLPlayerStats{h}, timestamp)

				db.Save(&seasonStats)
			}

			for _, a := range awayPlayerStats {
				if a.Snaps == 0 {
					continue
				}
				if a.WasInjured {
					playerRecord := GetNFLPlayerRecord(strconv.Itoa(a.NFLPlayerID))
					playerRecord.SetIsInjured(a.WasInjured, a.InjuryType, a.WeeksOfRecovery)
					db.Save(&playerRecord)
				}
				// playerSeasonStat := playerSeasonStatsMap[a.NFLPlayerID]
				seasonStats := GetNFLSeasonStatsByPlayerAndSeason(strconv.Itoa(a.NFLPlayerID), strconv.Itoa(int(timestamp.NFLSeasonID)))
				if seasonStats.ID == 0 {
					seasonStats = structs.NFLPlayerSeasonStats{
						NFLPlayerID: uint(a.NFLPlayerID),
						SeasonID:    uint(timestamp.NFLSeasonID),
						TeamID:      uint(a.TeamID),
						Team:        a.Team,
						Year:        uint(timestamp.Season),
					}
				}

				seasonStats.MapStats([]structs.NFLPlayerStats{a}, timestamp)

				db.Save(&seasonStats)
			}

			// Update Standings
			homeTeamStandings := GetNFLStandingsByTeamIDAndSeasonID(strconv.Itoa(homeTeamID), strconv.Itoa(timestamp.NFLSeasonID))
			awayTeamStandings := GetNFLStandingsByTeamIDAndSeasonID(strconv.Itoa(awayTeamID), strconv.Itoa(timestamp.NFLSeasonID))

			homeTeamStandings.UpdateNFLStandings(game)
			awayTeamStandings.UpdateNFLStandings(game)

			if game.HomeTeamCoach != "AI" && !timestamp.NFLPreseason {
				homeCoach := GetNFLUserByUsername(game.HomeTeamCoach)
				homeCoach.UpdateCoachRecord(game)

				err := db.Save(&homeCoach).Error
				if err != nil {
					log.Panicln("Could not save coach record for team " + strconv.Itoa(homeTeamID))
				}
			}

			if game.AwayTeamCoach != "AI" && !timestamp.NFLPreseason {
				awayCoach := GetNFLUserByUsername(game.AwayTeamCoach)
				awayCoach.UpdateCoachRecord(game)
				err := db.Save(&awayCoach).Error
				if err != nil {
					log.Panicln("Could not save coach record for team " + strconv.Itoa(awayTeamID))
				}
			}

			// Save
			if !timestamp.NFLPreseason {
				if game.NextGameID > 0 {
					nextGameID := strconv.Itoa(int(game.NextGameID))
					winningTeamID := 0
					winningTeam := ""
					winningCoach := ""
					if game.HomeTeamWin {
						winningTeamID = game.HomeTeamID
						winningTeam = game.HomeTeam
						winningCoach = game.HomeTeamCoach
					} else {
						winningTeamID = game.AwayTeamID
						winningTeam = game.AwayTeam
						winningCoach = game.AwayTeamCoach
					}

					nextGame := GetNFLGameByGameID(nextGameID)
					nextGame.AddTeam(game.NextGameHOA == "H", winningTeamID, winningTeam, winningCoach)
					if !nextGame.IsNeutral && !game.IsSuperBowl {
						stadiumID := 0
						stadium := ""
						city := ""
						state := ""
						isDomed := false
						if game.HomeTeamWin {
							stadium = game.Stadium
							city = game.City
							state = game.State
							stadiumID = int(game.StadiumID)
							isDomed = game.IsDomed
						} else {
							awayTeam := GetNFLTeamByTeamID(strconv.Itoa(awayTeamID))
							stadiumRecord := GetStadiumByStadiumID(strconv.Itoa(int(awayTeam.StadiumID)))
							stadium = awayTeam.Stadium
							city = awayTeam.City
							state = awayTeam.State
							stadiumID = int(awayTeam.StadiumID)
							isDomed = stadiumRecord.IsDomed
						}
						nextGame.AddLocation(stadiumID, stadium, city, state, isDomed)
					}

					// Updating matchup for playoff game!
					db.Save(&nextGame)
				}

				db.Save(&homeTeamSeasonStats)
				db.Save(&awayTeamSeasonStats)
				db.Save(&homeTeamStandings)
				db.Save(&awayTeamStandings)
			}

		}
	}

	db.Save(&timestamp)
}

// UpdateTimestamp - Update the timestamp
func UpdateTimestamp(updateTimestampDto structs.UpdateTimestampDto) structs.Timestamp {
	db := dbprovider.GetInstance().GetDB()

	timestamp := GetTimestamp()

	if updateTimestampDto.MoveUpCollegeWeek {
		// Update Standings based on current week's games

		// Sync to Next Week
		// UpdateStandings(timestamp)
		UpdateGameplanPenalties()
		timestamp.SyncToNextWeek()
	}
	// else if updateTimestampDto.ThursdayGames && !timestamp.ThursdayGames {
	// 	timestamp.ToggleThursdayGames()
	// } else if updateTimestampDto.FridayGames && !timestamp.FridayGames {
	// 	timestamp.ToggleFridayGames()
	// } else if updateTimestampDto.SaturdayMorning && !timestamp.SaturdayMorning {
	// 	timestamp.ToggleSaturdayMorningGames()
	// } else if updateTimestampDto.SaturdayNoon && !timestamp.SaturdayNoon {
	// 	timestamp.ToggleSaturdayNoonGames()
	// } else if updateTimestampDto.SaturdayEvening && !timestamp.SaturdayEvening {
	// 	timestamp.ToggleSaturdayEveningGames()
	// } else if updateTimestampDto.SaturdayNight && !timestamp.SaturdayNight {
	// 	timestamp.ToggleSaturdayNightGames()
	// }

	if updateTimestampDto.ToggleRecruitingLock {
		timestamp.ToggleLockRecruiting()
	}

	// if updateTimestampDto.RESSynced && !timestamp.RecruitingEfficiencySynced {
	// 	timestamp.ToggleRES()
	// 	SyncRecruitingEfficiency(timestamp)
	// }

	if updateTimestampDto.RecruitingSynced && !timestamp.RecruitingSynced && timestamp.IsRecruitingLocked {
		SyncRecruiting(timestamp)
		timestamp.ToggleRecruiting()
	}

	err := db.Save(&timestamp).Error
	if err != nil {
		fmt.Println(err.Error())
		log.Fatalf("Could not save timestamp")
	}

	return timestamp
}

// Week Funcs
func CreateCollegeWeek() {

}

// Season Funcs
func CreateCollegeSeason() {

}

// Season Funcs
func MoveUpInOffseasonFreeAgency() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	if ts.IsNFLOffSeason {
		ts.MoveUpFreeAgencyRound()
	}
	db.Save(&ts)
}

func GetNewsLogs(weekID string, seasonID string) []structs.NewsLog {
	db := dbprovider.GetInstance().GetDB()

	var logs []structs.NewsLog

	err := db.Where("week_id = ? AND season_id = ?", weekID, seasonID).Find(&logs).Error
	if err != nil {
		fmt.Println(err)
	}

	return logs
}

func GetAllNewsLogs() []structs.NewsLog {
	db := dbprovider.GetInstance().GetDB()

	var logs []structs.NewsLog

	err := db.Where("league = ?", "CFB").Find(&logs).Error
	if err != nil {
		fmt.Println(err)
	}

	return logs
}

func GetAllNFLNewsLogs() []structs.NewsLog {
	db := dbprovider.GetInstance().GetDB()

	var logs []structs.NewsLog

	err := db.Where("league = ?", "NFL").Find(&logs).Error
	if err != nil {
		fmt.Println(err)
	}

	return logs
}

func GetWeeksInASeason(seasonID string, weekID string) []structs.CollegeWeek {
	db := dbprovider.GetInstance().GetDB()

	var weeks []structs.CollegeWeek

	err := db.Where("season_id = ?", seasonID).Find(&weeks).Error
	if err != nil {
		fmt.Println(err)
	}

	return weeks
}
