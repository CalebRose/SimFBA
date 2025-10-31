package managers

import (
	"fmt"
	"strconv"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/models"
	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
)

func GetCollegeGamesByWeekIdAndSeasonID(WeekID string, SeasonID string) []structs.CollegeGame {
	db := dbprovider.GetInstance().GetDB()

	var games []structs.CollegeGame

	db.Where("week_id = ? AND season_id = ?", WeekID, SeasonID).Find(&games)

	return games
}

func GetCollegeGamesByWeekId(id, timeslot string, springGames bool) []structs.CollegeGame {
	db := dbprovider.GetInstance().GetDB()

	var games []structs.CollegeGame

	db.Order("time_slot asc").Where("week_id = ? AND is_spring_game = ?", id, springGames).Find(&games)

	return games
}

func GetCollegeGamesByTimeslotAndWeekId(id, timeslot string, springGames bool) []structs.CollegeGame {
	db := dbprovider.GetInstance().GetDB()

	var games []structs.CollegeGame

	db.Order("time_slot asc").Where("week_id = ? AND time_slot = ? AND is_spring_game = ?", id, timeslot, springGames).Find(&games)

	return games
}

func GetCollegeGamesByTeamIdAndSeasonId(TeamID string, SeasonID string, isSpringGame bool) []structs.CollegeGame {
	db := dbprovider.GetInstance().GetDB()

	var games []structs.CollegeGame

	db.Order("week_id asc").Where("season_id = ? AND (home_team_id = ? OR away_team_id = ?) AND is_spring_game = ?", SeasonID, TeamID, TeamID, isSpringGame).Find(&games)

	return games
}

func GetCollegeGamesBySeasonID(SeasonID string) []structs.CollegeGame {
	ts := GetTimestamp()
	return repository.FindCollegeGamesRecords(SeasonID, ts.CFBSpringGames)
}

func GetNFLGamesByTeamIdAndSeasonId(TeamID string, SeasonID string) []structs.NFLGame {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	var games []structs.NFLGame

	db.Order("week_id asc").Where("season_id = ? AND (home_team_id = ? OR away_team_id = ?) AND is_preseason_game = ?", SeasonID, TeamID, TeamID, ts.NFLPreseason).Find(&games)

	return games
}

func GetNFLGamesByWeekAndSeasonID(WeekID string, SeasonID string) []structs.NFLGame {
	db := dbprovider.GetInstance().GetDB()

	var games []structs.NFLGame

	db.Order("time_slot asc").Where("week_id = ? and season_id = ?", WeekID, SeasonID).Find(&games)

	return games
}

func GetNFLGamesBySeasonID(SeasonID string) []structs.NFLGame {
	ts := GetTimestamp()
	return repository.FindNFLGamesRecords(SeasonID, ts.NFLPreseason)
}

func GetNFLGamesByWeekId(id string, isPreseason bool) []structs.NFLGame {
	db := dbprovider.GetInstance().GetDB()

	var games []structs.NFLGame

	db.Order("time_slot asc").Where("week_id = ? AND is_preseason_game = ?", id, isPreseason).Find(&games)

	return games
}

func GetNFLGamesByTimeslotAndWeekId(id, timeslot string, isPreseason bool) []structs.NFLGame {
	db := dbprovider.GetInstance().GetDB()

	var games []structs.NFLGame

	db.Order("time_slot asc").Where("week_id = ? AND time_slot = ? AND is_preseason_game = ?", id, timeslot, isPreseason).Find(&games)

	return games
}

func GetTeamScheduleForBot(TeamID string, SeasonID string) []models.CollegeGameResponse {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	var games []structs.CollegeGame

	db.Order("week_id asc").Where("season_id = ? AND (home_team_id = ? OR away_team_id = ?) AND is_spring_game = ?", SeasonID, TeamID, TeamID, ts.CFBSpringGames).Find(&games)

	var gameResponses []models.CollegeGameResponse

	for _, game := range games {
		showGame := game.WeekID < ts.CollegeWeekID
		gameResponse := models.CollegeGameResponse{
			Week:                     game.Week,
			SeasonID:                 game.SeasonID,
			HomeTeamID:               game.HomeTeamID,
			HomeTeam:                 game.HomeTeam,
			HomeTeamCoach:            game.HomeTeamCoach,
			HomeTeamWin:              game.HomeTeamWin,
			HomeTeamScore:            game.HomeTeamScore,
			AwayTeamID:               game.AwayTeamID,
			AwayTeam:                 game.AwayTeam,
			AwayTeamCoach:            game.AwayTeamCoach,
			AwayTeamWin:              game.AwayTeamWin,
			AwayTeamScore:            game.AwayTeamScore,
			Stadium:                  game.Stadium,
			City:                     game.City,
			State:                    game.State,
			IsNeutral:                game.IsNeutral,
			IsConference:             game.IsConference,
			IsDivisional:             game.IsDivisional,
			IsConferenceChampionship: game.IsConferenceChampionship,
			IsBowlGame:               game.IsBowlGame,
			IsPlayoffGame:            game.IsPlayoffGame,
			IsNationalChampionship:   game.IsNationalChampionship,
			GameComplete:             game.GameComplete,
			ShowGame:                 showGame,
			GameTitle:                game.GameTitle,
			TimeSlot:                 game.TimeSlot,
		}

		gameResponses = append(gameResponses, gameResponse)
	}

	return gameResponses
}

func GetCFBScheduleByConference(conf string) []structs.CollegeGame {
	ts := GetTimestamp()

	var games []structs.CollegeGame

	teams := GetCollegeTeamsByConference(conf)
	teamMap := make(map[uint]bool)

	for _, t := range teams {
		teamMap[t.ID] = true
	}

	weekID := strconv.Itoa(int(ts.CollegeWeekID))
	seasonID := strconv.Itoa(int(ts.CollegeSeasonID))
	matches := GetCollegeGamesByWeekIdAndSeasonID(weekID, seasonID)
	for _, m := range matches {
		if teamMap[uint(m.HomeTeamID)] || teamMap[uint(m.AwayTeamID)] {
			games = append(games, m)
		}
	}

	return games
}

func GetCFBCurrentWeekSchedule() []structs.CollegeGame {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	var games []structs.CollegeGame

	db.Order("time_slot asc").
		Where("season_id = ? AND week_id = ?", strconv.Itoa(int(ts.CollegeSeasonID)), strconv.Itoa(int(ts.CollegeWeekID))).
		Find(&games)

	return games
}

func GetNFLCurrentWeekSchedule() []structs.NFLGame {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	var games []structs.NFLGame

	db.Order("time_slot asc").
		Where("season_id = ? AND week_id = ?", strconv.Itoa(int(ts.NFLSeasonID)), strconv.Itoa(int(ts.NFLWeekID))).
		Find(&games)

	return games
}

func GetAllCollegeGames() []structs.CollegeGame {
	db := dbprovider.GetInstance().GetDB()

	var games []structs.CollegeGame

	db.Order("week_id asc").Find(&games)

	return games
}

func GetCollegeGamesByTeamId(TeamID string) []structs.CollegeGame {
	db := dbprovider.GetInstance().GetDB()

	var games []structs.CollegeGame

	db.Order("week_id asc").Where("home_team_id = ? OR away_team_id = ?", TeamID, TeamID).Find(&games)

	return games
}

func GetNFLGamesByTeamId(TeamID string) []structs.NFLGame {
	db := dbprovider.GetInstance().GetDB()

	var games []structs.NFLGame

	db.Order("week_id asc").Where("home_team_id = ? OR away_team_id = ?", TeamID, TeamID).Find(&games)

	return games
}

func GetCollegeGameByAbbreviationsWeekAndSeasonID(HomeTeamAbbr string, WeekID string, SeasonID string) structs.CollegeGame {
	db := dbprovider.GetInstance().GetDB()

	var game structs.CollegeGame

	err := db.Where("season_id = ? AND week_id = ? AND (home_team = ? OR away_team = ?)", SeasonID, WeekID, HomeTeamAbbr, HomeTeamAbbr).Find(&game).Error
	if err != nil {
		fmt.Println("Could not find game!")
	}

	return game
}

// func GetCollegeBoxScoreResultsByGameID(gameID string) {
// 	db := dbprovider.GetInstance().GetDB()

// 	game := GetCollegeGameByGameID()
// 	homeTeamID := strconv.Itoa(int(game.HomeTeamID))
// 	awayTeamID := strconv.Itoa(int(game.AwayTeamID))
// }

func GetCollegeGameByGameID(gameID string) structs.CollegeGame {
	db := dbprovider.GetInstance().GetDB()

	var game structs.CollegeGame

	err := db.Where("id = ?", gameID).Find(&game).Error
	if err != nil {
		fmt.Println("Could not find game!")
	}

	return game
}

func GetNFLGameByAbbreviationsWeekAndSeasonID(HomeTeamAbbr string, WeekID string, SeasonID string) structs.NFLGame {
	db := dbprovider.GetInstance().GetDB()

	var game structs.NFLGame

	err := db.Where("season_id = ? AND week_id = ? AND (home_team = ? OR away_team = ?)", SeasonID, WeekID, HomeTeamAbbr, HomeTeamAbbr).Find(&game).Error
	if err != nil {
		fmt.Println("Could not find game!")
	}

	return game
}

func GetNFLGameByGameID(gameID string) structs.NFLGame {
	db := dbprovider.GetInstance().GetDB()

	var game structs.NFLGame

	err := db.Where("id = ?", gameID).Find(&game).Error
	if err != nil {
		fmt.Println("Could not find game!")
	}

	return game
}

func UpdateTimeslot(dto structs.UpdateTimeslotDTO) {
	db := dbprovider.GetInstance().GetDB()
	gameID := strconv.Itoa(int(dto.GameID))
	// regions := getRegionalWeather()
	// rainForecasts := getRainChart()
	// mixForecasts := getMixChart()
	// snowForecasts := getSnowChart()
	// teamRegions := getRegionsForSchools()
	if dto.League == "CFB" {
		game := GetCollegeGameByGameID(gameID)
		game.UpdateTimeslot(dto.Timeslot)
		// GenerateWeatherForGame(db, game, teamRegions, regions, rainForecasts, mixForecasts, snowForecasts)
		// db.Save(&game)
		// return game, structs.NFLGame{}
	} else {
		game := GetNFLGameByGameID(gameID)
		// homeTeam := GetNFLTeamByTeamID(strconv.Itoa(game.HomeTeamID))
		game.UpdateTimeslot(dto.Timeslot)
		db.Save(&game)
		// GenerateWeatherForNFLGame(db, game, homeTeam.TeamAbbr, teamRegions, regions, rainForecasts, mixForecasts, snowForecasts)
		// return structs.CollegeGame{}, game
	}
}

func RunTheGames() {
	db := dbprovider.GetInstance().GetDB()

	tsChn := make(chan structs.Timestamp)

	go func() {
		timestamp := GetTimestamp()
		tsChn <- timestamp
	}()

	ts := <-tsChn
	close(tsChn)

	ts.ToggleRunGames()
	repository.SaveTimestamp(ts, db)

}

func FixByeWeekLogic() {
	// DB
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID := strconv.Itoa(ts.CollegeSeasonID)
	// Teams
	collegeTeams := GetAllCollegeTeams()

	// Loop through each team and their games
	for _, t := range collegeTeams {
		teamID := strconv.Itoa(int(t.ID))
		games := GetCollegeGamesByTeamIdAndSeasonId(teamID, seasonID, ts.CFBSpringGames)
		prevWeek := 0
		for _, game := range games {
			diff := game.Week - prevWeek
			if diff > 1 {
				game.AssignByeWeek(t.ID)
				repository.SaveCFBGameRecord(game, db)
			}
			prevWeek = game.Week
		}
	}

	nflTeams := GetAllNFLTeams()

	// Loop through each team and their games
	for _, t := range nflTeams {
		teamID := strconv.Itoa(int(t.ID))
		games := GetNFLGamesByTeamIdAndSeasonId(teamID, seasonID)
		prevWeek := 0
		for _, game := range games {
			diff := game.Week - prevWeek
			if diff > 1 {
				game.AssignByeWeek(t.ID)
				repository.SaveNFLGameRecord(game, db)
			}
			prevWeek = game.Week
		}
	}
}

func GetNFLGamesMapSeasonId(SeasonID string) map[uint]structs.NFLGame {
	db := dbprovider.GetInstance().GetDB()

	var games []structs.NFLGame

	db.Order("week_id asc").Where("season_id = ?", SeasonID).Find(&games)

	gameMap := make(map[uint]structs.NFLGame)

	for _, g := range games {
		gameMap[g.ID] = g
	}

	return gameMap
}

func GetCFBGamesMapSeasonId(SeasonID string) map[uint]structs.CollegeGame {
	db := dbprovider.GetInstance().GetDB()

	var games []structs.CollegeGame

	db.Order("week_id asc").Where("season_id = ?", SeasonID).Find(&games)

	gameMap := make(map[uint]structs.CollegeGame)

	for _, g := range games {
		gameMap[g.ID] = g
	}

	return gameMap
}
