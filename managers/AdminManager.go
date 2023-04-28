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
	if timestamp.RecruitingSynced {
		// Sync to Next Week
		UpdateGameplanPenalties()
		RecoverPlayers()
		CheckNFLRookiesForLetterGrade(strconv.Itoa(int(timestamp.NFLSeasonID)))
		timestamp.SyncToNextWeek()
		db.Save(&timestamp)
	}

	return timestamp
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
		playerSeasonStats := GetCollegePlayerSeasonStatsBySeason(strconv.Itoa(timestamp.CollegeSeasonID))
		playerSeasonStatsMap := make(map[int]*structs.CollegePlayerSeasonStats)
		for _, s := range seasonStats {
			seasonStatsMap[int(s.TeamID)] = &s
		}

		for _, p := range playerSeasonStats {
			playerSeasonStatsMap[int(p.CollegePlayerID)] = &p
		}

		for _, game := range games {
			// Get team stats
			gameID := strconv.Itoa(int(game.ID))
			homeTeamID := game.HomeTeamID
			awayTeamID := game.AwayTeamID

			homeTeamSeasonStats := seasonStatsMap[homeTeamID]
			awayTeamSeasonStats := seasonStatsMap[awayTeamID]

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
				playerSeasonStat := playerSeasonStatsMap[h.CollegePlayerID]
				playerSeasonStat.MapStats([]structs.CollegePlayerStats{h})

				db.Save(&playerSeasonStat)
			}

			for _, a := range awayPlayerStats {
				if a.Snaps == 0 {
					continue
				}
				playerSeasonStat := playerSeasonStatsMap[a.CollegePlayerID]
				playerSeasonStat.MapStats([]structs.CollegePlayerStats{a})

				db.Save(&playerSeasonStat)
			}

			// Update Standings
			homeTeamStandings := GetCFBStandingsByTeamIDAndSeasonID(strconv.Itoa(homeTeamID), strconv.Itoa(timestamp.CollegeSeasonID))
			awayTeamStandings := GetCFBStandingsByTeamIDAndSeasonID(strconv.Itoa(awayTeamID), strconv.Itoa(timestamp.CollegeSeasonID))

			homeTeamStandings.UpdateCollegeStandings(game)
			awayTeamStandings.UpdateCollegeStandings(game)

			if game.HomeTeamCoach != "AI" {
				homeCoach := GetCollegeCoachByCoachName(game.HomeTeamCoach)
				homeCoach.UpdateCoachRecord(game)

				err := db.Save(&homeCoach).Error
				if err != nil {
					log.Panicln("Could not save coach record for team " + strconv.Itoa(homeTeamID))
				}
			}

			if game.AwayTeamCoach != "AI" {
				awayCoach := GetCollegeCoachByCoachName(game.AwayTeamCoach)
				awayCoach.UpdateCoachRecord(game)
				err := db.Save(&awayCoach).Error
				if err != nil {
					log.Panicln("Could not save coach record for team " + strconv.Itoa(awayTeamID))
				}
			}

			// Save
			db.Save(&homeTeamSeasonStats)
			db.Save(&awayTeamSeasonStats)
			db.Save(&homeTeamStandings)
			db.Save(&awayTeamStandings)
		}
	} else {
		// Get Games
		games := GetNFLGamesByTimeslotAndWeekId(strconv.Itoa(timestamp.NFLWeekID), timeslot)

		seasonStats := GetNFLTeamSeasonStatsBySeason(strconv.Itoa(timestamp.NFLSeasonID))
		seasonStatsMap := make(map[int]*structs.NFLTeamSeasonStats)
		playerSeasonStats := GetNFLPlayerSeasonStatsBySeason(strconv.Itoa(timestamp.NFLSeasonID))
		playerSeasonStatsMap := make(map[int]*structs.NFLPlayerSeasonStats)
		for _, s := range seasonStats {
			seasonStatsMap[int(s.TeamID)] = &s
		}

		for _, p := range playerSeasonStats {
			playerSeasonStatsMap[int(p.NFLPlayerID)] = &p
		}

		for _, game := range games {
			// Get team stats
			gameID := strconv.Itoa(int(game.ID))
			homeTeamID := game.HomeTeamID
			awayTeamID := game.AwayTeamID

			homeTeamSeasonStats := seasonStatsMap[homeTeamID]
			awayTeamSeasonStats := seasonStatsMap[awayTeamID]

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
				playerSeasonStat := playerSeasonStatsMap[h.NFLPlayerID]
				playerSeasonStat.MapStats([]structs.NFLPlayerStats{h})

				db.Save(&playerSeasonStat)
			}

			for _, a := range awayPlayerStats {
				if a.Snaps == 0 {
					continue
				}
				playerSeasonStat := playerSeasonStatsMap[a.NFLPlayerID]
				playerSeasonStat.MapStats([]structs.NFLPlayerStats{a})

				db.Save(&playerSeasonStat)
			}

			// Update Standings
			homeTeamStandings := GetNFLStandingsByTeamIDAndSeasonID(strconv.Itoa(homeTeamID), strconv.Itoa(timestamp.CollegeSeasonID))
			awayTeamStandings := GetNFLStandingsByTeamIDAndSeasonID(strconv.Itoa(awayTeamID), strconv.Itoa(timestamp.CollegeSeasonID))

			homeTeamStandings.UpdateNFLStandings(game)
			awayTeamStandings.UpdateNFLStandings(game)

			if game.HomeTeamCoach != "AI" {
				homeCoach := GetNFLUserByUsername(game.HomeTeamCoach)
				homeCoach.UpdateCoachRecord(game)

				err := db.Save(&homeCoach).Error
				if err != nil {
					log.Panicln("Could not save coach record for team " + strconv.Itoa(homeTeamID))
				}
			}

			if game.AwayTeamCoach != "AI" {
				awayCoach := GetNFLUserByUsername(game.AwayTeamCoach)
				awayCoach.UpdateCoachRecord(game)
				err := db.Save(&awayCoach).Error
				if err != nil {
					log.Panicln("Could not save coach record for team " + strconv.Itoa(awayTeamID))
				}
			}

			// Save
			db.Save(&homeTeamSeasonStats)
			db.Save(&awayTeamSeasonStats)
			db.Save(&homeTeamStandings)
			db.Save(&awayTeamStandings)
		}
	}
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
