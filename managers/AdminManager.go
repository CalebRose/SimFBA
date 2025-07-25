package managers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
)

// Timestamp Funcs
// GetTimestamp -- Get the Timestamp
func GetTimestamp() structs.Timestamp {
	db := dbprovider.GetInstance().GetDB()

	var timestamp structs.Timestamp

	err := db.First(&timestamp).Error
	if err != nil {
		log.Printf("Error querying for timestamp: %v", err)
	}

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
	ts := GetTimestamp()
	if ts.CollegeWeek < 21 || !ts.IsOffSeason {
		ResetCollegeStandingsRanks()
	}

	// Sync to Next Week
	UpdateGameplanPenalties()
	RecoverPlayers()
	CheckNFLRookiesForLetterGrade(strconv.Itoa(int(ts.NFLSeasonID)))
	ts.SyncToNextWeek()

	if ts.CollegeWeek < 21 && !ts.CollegeSeasonOver && !ts.IsOffSeason && !ts.CFBSpringGames {
		SyncCollegePollSubmissionForCurrentWeek(uint(ts.CollegeWeek), uint(ts.CollegeWeekID), uint(ts.CollegeSeasonID))
		ts.TogglePollRan()
	}
	if ts.NFLWeek > 15 {
		SyncExtensionOffers()
		AllocateCapsheets()
	}
	repository.SaveTimestamp(ts, db)

	return ts
}

func SyncTimeslot(timeslot string) {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	// Update timeslot
	ts.ToggleTimeSlot(timeslot)
	cgt, cfbgt := ts.GetCFBCurrentGameType()
	ngt, nflgt := ts.GetNFLCurrentGameType()

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
		gameIDs := []string{}
		games := GetCollegeGamesByTimeslotAndWeekId(strconv.Itoa(ts.CollegeWeekID), timeslot, ts.CFBSpringGames)
		seasonStats := GetCollegeSeasonStatsBySeason(strconv.Itoa(ts.CollegeSeasonID), cfbgt)
		seasonStatsMap := make(map[int]*structs.CollegeTeamSeasonStats)
		for _, s := range seasonStats {
			seasonStatsMap[int(s.TeamID)] = &s
		}

		for _, game := range games {
			// Get team stats
			gameID := strconv.Itoa(int(game.ID))
			gameIDs = append(gameIDs, gameID)

			homeTeamID := game.HomeTeamID
			awayTeamID := game.AwayTeamID

			// homeTeamSeasonStats := seasonStatsMap[homeTeamID]
			// awayTeamSeasonStats := seasonStatsMap[awayTeamID]
			homeTeamSeasonStats := GetCollegeTeamSeasonStatsBySeason(strconv.Itoa(homeTeamID), strconv.Itoa(int(ts.CollegeSeasonID)), cfbgt)
			awayTeamSeasonStats := GetCollegeTeamSeasonStatsBySeason(strconv.Itoa(awayTeamID), strconv.Itoa(int(ts.CollegeSeasonID)), cfbgt)
			if homeTeamSeasonStats.ID == 0 {
				homeTeamSeasonStats = structs.CollegeTeamSeasonStats{
					TeamID:   uint(homeTeamID),
					SeasonID: uint(game.SeasonID),
					Year:     ts.Season,
					BaseTeamStats: structs.BaseTeamStats{
						GameType: uint8(cgt),
					},
				}
			}

			if awayTeamSeasonStats.ID == 0 {
				awayTeamSeasonStats = structs.CollegeTeamSeasonStats{
					TeamID:   uint(awayTeamID),
					SeasonID: uint(game.SeasonID),
					Year:     ts.Season,
					BaseTeamStats: structs.BaseTeamStats{
						GameType: uint8(cgt),
					},
				}
			}

			homeTeamStats := GetCollegeTeamStatsByGame(strconv.Itoa(homeTeamID), gameID)
			awayTeamStats := GetCollegeTeamStatsByGame(strconv.Itoa(awayTeamID), gameID)

			homeTeamSeasonStats.MapStats([]structs.CollegeTeamStats{homeTeamStats}, ts.CollegeSeasonID)
			awayTeamSeasonStats.MapStats([]structs.CollegeTeamStats{awayTeamStats}, ts.CollegeSeasonID)
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
					repository.SaveCFBPlayer(playerRecord, db)
				}
				// playerSeasonStat := playerSeasonStatsMap[h.CollegePlayerID]
				playerSeasonStat := GetCollegeSeasonStatsByPlayerAndSeason(strconv.Itoa(h.CollegePlayerID), strconv.Itoa(int(ts.CollegeSeasonID)), cfbgt)
				if playerSeasonStat.ID == 0 {
					playerSeasonStat = structs.CollegePlayerSeasonStats{
						CollegePlayerID: uint(h.CollegePlayerID),
						SeasonID:        uint(ts.CollegeSeasonID),
						Year:            uint(h.Year),
						BasePlayerStats: structs.BasePlayerStats{
							GameType: uint8(cgt),
						},
					}
				}

				playerSeasonStat.MapStats([]structs.CollegePlayerStats{h})
				repository.SaveCollegePlayerSeasonStats(playerSeasonStat, db)
			}

			for _, a := range awayPlayerStats {
				if a.Snaps == 0 {
					continue
				}
				if a.WasInjured {
					playerRecord := GetCollegePlayerByCollegePlayerId(strconv.Itoa(a.CollegePlayerID))
					playerRecord.SetIsInjured(a.WasInjured, a.InjuryType, a.WeeksOfRecovery)
					repository.SaveCFBPlayer(playerRecord, db)
				}
				playerSeasonStat := GetCollegeSeasonStatsByPlayerAndSeason(strconv.Itoa(a.CollegePlayerID), strconv.Itoa(int(ts.CollegeSeasonID)), cfbgt)
				if playerSeasonStat.ID == 0 {
					playerSeasonStat = structs.CollegePlayerSeasonStats{
						CollegePlayerID: uint(a.CollegePlayerID),
						SeasonID:        uint(ts.CollegeSeasonID),
						Year:            uint(a.Year),
						BasePlayerStats: structs.BasePlayerStats{
							GameType: uint8(cgt),
						},
					}
				}
				playerSeasonStat.MapStats([]structs.CollegePlayerStats{a})
				repository.SaveCollegePlayerSeasonStats(playerSeasonStat, db)
			}

			playerSnaps := GetAllCollegePlayerSnapsByGame(gameID)
			for _, snap := range playerSnaps {
				playerID := strconv.Itoa(int(snap.PlayerID))
				seasonID := strconv.Itoa(int(snap.SeasonID))
				seasonSnaps := GetCollegeSeasonSnapsByPlayerAndSeason(playerID, seasonID)
				if seasonSnaps.ID == 0 {
					seasonSnaps = structs.CollegePlayerSeasonSnaps{
						BasePlayerSeasonSnaps: structs.BasePlayerSeasonSnaps{
							PlayerID: snap.PlayerID,
							SeasonID: snap.SeasonID,
						},
					}
				}
				seasonSnaps.AddToSeason(snap.BasePlayerGameSnaps)
				repository.SaveCFBSeasonSnaps(seasonSnaps, db)
			}

			// Update Standings
			homeTeamStandings := GetCFBStandingsByTeamIDAndSeasonID(strconv.Itoa(homeTeamID), strconv.Itoa(ts.CollegeSeasonID))
			awayTeamStandings := GetCFBStandingsByTeamIDAndSeasonID(strconv.Itoa(awayTeamID), strconv.Itoa(ts.CollegeSeasonID))

			homeTeamStandings.UpdateCollegeStandings(game)
			awayTeamStandings.UpdateCollegeStandings(game)

			if game.HomeTeamCoach != "AI" && !ts.CFBSpringGames {
				homeCoach := GetCollegeCoachByCoachName(game.HomeTeamCoach)
				homeCoach.UpdateCoachRecord(game)

				err := db.Save(&homeCoach).Error
				if err != nil {
					log.Panicln("Could not save coach record for team " + strconv.Itoa(homeTeamID))
				}
			}

			if game.AwayTeamCoach != "AI" && !ts.CFBSpringGames {
				awayCoach := GetCollegeCoachByCoachName(game.AwayTeamCoach)
				awayCoach.UpdateCoachRecord(game)
				err := db.Save(&awayCoach).Error
				if err != nil {
					log.Panicln("Could not save coach record for team " + strconv.Itoa(awayTeamID))
				}
			}

			// Save
			if !ts.CFBSpringGames {
				if game.NextGameID > 0 {
					nextGameID := strconv.Itoa(int(game.NextGameID))
					winningTeamID := 0
					winningTeam := ""
					winningCoach := ""
					winningRank := 0
					if game.HomeTeamWin {
						winningTeamID = game.HomeTeamID
						winningTeam = game.HomeTeam
						winningCoach = game.HomeTeamCoach
						winningRank = int(game.HomeTeamRank)
					} else {
						winningTeamID = game.AwayTeamID
						winningTeam = game.AwayTeam
						winningCoach = game.AwayTeamCoach
						winningRank = int(game.AwayTeamRank)
					}

					nextGame := GetCollegeGameByGameID(nextGameID)
					nextGame.AddTeam(game.NextGameHOA == "H", winningTeamID, winningRank, winningTeam, winningCoach)
					if !nextGame.IsNeutral && !game.IsNationalChampionship {
						stadiumID := 0
						stadium := ""
						city := ""
						state := ""
						isDomed := false

						if game.HomeTeamWin {
							homeTeam := GetTeamByTeamID(strconv.Itoa(homeTeamID))
							stadiumRecord := GetStadiumByStadiumID(strconv.Itoa(int(homeTeam.StadiumID)))
							stadium = homeTeam.Stadium
							city = homeTeam.City
							state = homeTeam.State
							stadiumID = int(homeTeam.StadiumID)
							isDomed = stadiumRecord.IsDomed
						} else {
							awayTeam := GetTeamByTeamID(strconv.Itoa(awayTeamID))
							stadiumRecord := GetStadiumByStadiumID(strconv.Itoa(int(awayTeam.StadiumID)))
							stadium = awayTeam.Stadium
							city = awayTeam.City
							state = awayTeam.State
							stadiumID = int(awayTeam.StadiumID)
							isDomed = stadiumRecord.IsDomed
						}
						if game.NextGameHOA == "H" && !nextGame.IsNationalChampionship && !nextGame.IsNeutral {
							nextGame.AddLocation(stadiumID, stadium, city, state, isDomed)
						}
					}

					// Updating matchup for playoff game!
					repository.SaveCFBGameRecord(nextGame, db)
				}
				repository.SaveCFBStandingsRecord(homeTeamStandings, db)
				repository.SaveCFBStandingsRecord(awayTeamStandings, db)
			}
			repository.SaveCFBTeamSeasonStats(homeTeamSeasonStats, db)
			repository.SaveCFBTeamSeasonStats(awayTeamSeasonStats, db)
			if (game.HomeTeamWin && ((game.HomeTeamRank == 0 && game.AwayTeamRank > 0) || (game.HomeTeamRank > 0 && game.AwayTeamRank > 0 && game.HomeTeamRank > game.AwayTeamRank))) ||
				(game.AwayTeamWin && ((game.AwayTeamRank == 0 && game.HomeTeamRank > 0) || (game.AwayTeamRank > 0 && game.HomeTeamRank > 0 && game.AwayTeamRank > game.HomeTeamRank))) {
				// NEWS LOG
				messageStart := "UPSET: "
				messageType := "Upset Alert"
				winningTeam := ""
				losingTeam := ""
				ws := 0
				ls := 0
				message := ""
				htr := ""
				atr := ""
				wtr := ""
				ltr := ""
				if game.HomeTeamRank > 0 {
					htr = "(" + strconv.Itoa(int(game.HomeTeamRank)) + ") "
				}
				if game.AwayTeamRank > 0 {
					atr = "(" + strconv.Itoa(int(game.AwayTeamRank)) + ") "
				}

				if game.HomeTeamWin {
					winningTeam = game.HomeTeam
					losingTeam = game.AwayTeam
					ws = game.HomeTeamScore
					ls = game.AwayTeamScore
					wtr = htr
					ltr = atr
				} else {
					winningTeam = game.AwayTeam
					losingTeam = game.HomeTeam
					ws = game.AwayTeamScore
					ls = game.HomeTeamScore
					wtr = atr
					ltr = htr
				}
				winningVerb := util.GetWinningVerb(ws, ls)

				message = messageStart + wtr + winningTeam + winningVerb + ltr + losingTeam + " " + strconv.Itoa(ws) + "-" + strconv.Itoa(ls) + " at " + game.Stadium + " in " + game.City + ", " + game.State

				CreateNewsLog("CFB", message, messageType, 0, ts)
			}
		}

		// Reveal Individual Stats
		db.Model(&structs.CollegePlayerStats{}).Where("game_id in (?)", gameIDs).Update("reveal_results", true)
		db.Model(&structs.CollegeTeamStats{}).Where("game_id in (?)", gameIDs).Update("reveal_results", true)
	} else {
		// Get Games
		games := GetNFLGamesByTimeslotAndWeekId(strconv.Itoa(ts.NFLWeekID), timeslot, ts.NFLPreseason)

		// seasonStatsMap := make(map[int]structs.NFLTeamSeasonStats)
		// for _, s := range seasonStats {
		// 	seasonStatsMap[int(s.TeamID)] = s
		// }
		gameIDs := []string{}

		for _, game := range games {
			// Get team stats
			gameID := strconv.Itoa(int(game.ID))
			gameIDs = append(gameIDs, gameID)
			homeTeamID := game.HomeTeamID
			awayTeamID := game.AwayTeamID

			// homeTeamSeasonStats := seasonStatsMap[homeTeamID]
			// awayTeamSeasonStats := seasonStatsMap[awayTeamID]
			homeTeamSeasonStats := GetNFLTeamSeasonStatsByTeamANDSeason(strconv.Itoa(homeTeamID), strconv.Itoa(int(ts.NFLSeasonID)), nflgt)
			awayTeamSeasonStats := GetNFLTeamSeasonStatsByTeamANDSeason(strconv.Itoa(awayTeamID), strconv.Itoa(int(ts.NFLSeasonID)), nflgt)

			homeTeamStats := GetNFLTeamStatsByGame(strconv.Itoa(homeTeamID), gameID)
			awayTeamStats := GetNFLTeamStatsByGame(strconv.Itoa(awayTeamID), gameID)

			homeTeamSeasonStats.MapStats([]structs.NFLTeamStats{homeTeamStats}, ts.Season, ts.CollegeSeasonID)
			awayTeamSeasonStats.MapStats([]structs.NFLTeamStats{awayTeamStats}, ts.Season, ts.CollegeSeasonID)
			homeTeamSeasonStats.AddGameType(uint8(ngt))
			awayTeamSeasonStats.AddGameType(uint8(ngt))
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
					repository.SaveNFLPlayer(playerRecord, db)
				}
				// playerSeasonStat := playerSeasonStatsMap[h.NFLPlayerID]
				seasonStats := GetNFLSeasonStatsByPlayerAndSeason(strconv.Itoa(h.NFLPlayerID), strconv.Itoa(int(ts.NFLSeasonID)), nflgt)
				if seasonStats.ID == 0 {
					seasonStats = structs.NFLPlayerSeasonStats{
						NFLPlayerID: uint(h.NFLPlayerID),
						SeasonID:    uint(ts.NFLSeasonID),
						Year:        uint(ts.Season),
						BasePlayerStats: structs.BasePlayerStats{
							GameType: uint8(ngt),
						},
					}
				}

				seasonStats.MapStats([]structs.NFLPlayerStats{h}, ts)
				repository.SaveNFLPlayerSeasonStats(seasonStats, db)
			}

			for _, a := range awayPlayerStats {
				if a.Snaps == 0 {
					continue
				}
				if a.WasInjured {
					playerRecord := GetNFLPlayerRecord(strconv.Itoa(a.NFLPlayerID))
					playerRecord.SetIsInjured(a.WasInjured, a.InjuryType, a.WeeksOfRecovery)
					repository.SaveNFLPlayer(playerRecord, db)
				}
				// playerSeasonStat := playerSeasonStatsMap[a.NFLPlayerID]
				seasonStats := GetNFLSeasonStatsByPlayerAndSeason(strconv.Itoa(a.NFLPlayerID), strconv.Itoa(int(ts.NFLSeasonID)), nflgt)
				if seasonStats.ID == 0 {
					seasonStats = structs.NFLPlayerSeasonStats{
						NFLPlayerID: uint(a.NFLPlayerID),
						SeasonID:    uint(ts.NFLSeasonID),
						Year:        uint(ts.Season),
						BasePlayerStats: structs.BasePlayerStats{
							GameType: uint8(ngt),
						},
					}
				}

				seasonStats.MapStats([]structs.NFLPlayerStats{a}, ts)
				repository.SaveNFLPlayerSeasonStats(seasonStats, db)
			}

			playerSnaps := GetAllNFLPlayerSnapsByGame(gameID)
			for _, snap := range playerSnaps {
				playerID := strconv.Itoa(int(snap.PlayerID))
				seasonID := strconv.Itoa(int(snap.SeasonID))
				seasonSnaps := GetNFLSeasonSnapsByPlayerAndSeason(playerID, seasonID)
				if seasonSnaps.ID == 0 {
					seasonSnaps = structs.NFLPlayerSeasonSnaps{
						BasePlayerSeasonSnaps: structs.BasePlayerSeasonSnaps{
							PlayerID: snap.PlayerID,
							SeasonID: snap.SeasonID,
						},
					}
				}
				seasonSnaps.AddToSeason(snap.BasePlayerGameSnaps)
				repository.SaveNFLSeasonSnaps(seasonSnaps, db)
			}

			// Update Standings
			homeTeamStandings := GetNFLStandingsByTeamIDAndSeasonID(strconv.Itoa(homeTeamID), strconv.Itoa(ts.NFLSeasonID))
			awayTeamStandings := GetNFLStandingsByTeamIDAndSeasonID(strconv.Itoa(awayTeamID), strconv.Itoa(ts.NFLSeasonID))

			homeTeamStandings.UpdateNFLStandings(game)
			awayTeamStandings.UpdateNFLStandings(game)

			if game.HomeTeamCoach != "AI" && !ts.NFLPreseason {
				homeCoach := GetNFLUserByUsername(game.HomeTeamCoach)
				homeCoach.UpdateCoachRecord(game)

				err := db.Save(&homeCoach).Error
				if err != nil {
					log.Panicln("Could not save coach record for team " + strconv.Itoa(homeTeamID))
				}
			}

			if game.AwayTeamCoach != "AI" && !ts.NFLPreseason {
				awayCoach := GetNFLUserByUsername(game.AwayTeamCoach)
				awayCoach.UpdateCoachRecord(game)
				err := db.Save(&awayCoach).Error
				if err != nil {
					log.Panicln("Could not save coach record for team " + strconv.Itoa(awayTeamID))
				}
			}

			// Save
			if !ts.NFLPreseason {
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
					if !nextGame.IsNeutral && !game.IsSuperBowl && game.NextGameHOA == "H" {
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
					repository.SaveNFLGameRecord(nextGame, db)
				}

				repository.SaveNFLTeamSeasonStats(homeTeamSeasonStats, db)
				repository.SaveNFLTeamSeasonStats(awayTeamSeasonStats, db)
				repository.SaveNFLStandingsRecord(homeTeamStandings, db)
				repository.SaveNFLStandingsRecord(awayTeamStandings, db)
			}
		}

		db.Model(&structs.NFLPlayerStats{}).Where("game_id in (?)", gameIDs).Update("reveal_results", true)
		db.Model(&structs.NFLTeamStats{}).Where("game_id in (?)", gameIDs).Update("reveal_results", true)

	}

	repository.SaveTimestamp(ts, db)
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

func GenerateOffseasonData() {
	// Fulfill Promises

	// Run First Phase of Portal

	// Create Standings Records for Leagues
	GenerateNewSeasonStandings()

	// Reset Recruiting Profiles
	ResetTeamRecruitingProfiles()

	// Generate Next Croot Class
	GenerateCroots()
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
