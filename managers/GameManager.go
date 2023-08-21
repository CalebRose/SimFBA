package managers

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/models"
	"github.com/CalebRose/SimFBA/structs"
)

func GetCollegeGamesByWeekIdAndSeasonID(WeekID string, SeasonID string) []structs.CollegeGame {
	db := dbprovider.GetInstance().GetDB()

	var games []structs.CollegeGame

	db.Where("week_id = ? AND season_id = ?", WeekID, SeasonID).Find(&games)

	return games
}

func GetCollegeGamesByTimeslotAndWeekId(id string, timeslot string) []structs.CollegeGame {
	db := dbprovider.GetInstance().GetDB()

	var games []structs.CollegeGame

	db.Order("time_slot asc").Where("week_id = ? AND time_slot = ?", id, timeslot).Find(&games)

	return games
}

func GetCollegeGamesByTeamIdAndSeasonId(TeamID string, SeasonID string) []structs.CollegeGame {
	db := dbprovider.GetInstance().GetDB()

	var games []structs.CollegeGame

	db.Order("week_id asc").Where("season_id = ? AND (home_team_id = ? OR away_team_id = ?)", SeasonID, TeamID, TeamID).Find(&games)

	return games
}

func GetCollegeGamesBySeasonID(SeasonID string) []structs.CollegeGame {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	var games []structs.CollegeGame

	db.Order("week_id asc").Where("season_id = ? and is_spring_game = ?", SeasonID, ts.CFBSpringGames).Find(&games)

	return games
}

func GetNFLGamesByTeamIdAndSeasonId(TeamID string, SeasonID string) []structs.NFLGame {
	db := dbprovider.GetInstance().GetDB()

	var games []structs.NFLGame

	db.Order("week_id asc").Where("season_id = ? AND (home_team_id = ? OR away_team_id = ?)", SeasonID, TeamID, TeamID).Find(&games)

	return games
}

func GetNFLGamesByWeekAndSeasonID(WeekID string, SeasonID string) []structs.NFLGame {
	db := dbprovider.GetInstance().GetDB()

	var games []structs.NFLGame

	db.Order("time_slot asc").Where("week_id = ? and season_id = ?", WeekID, SeasonID).Find(&games)

	return games
}

func GetNFLGamesBySeasonID(SeasonID string) []structs.NFLGame {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	var games []structs.NFLGame

	db.Order("week_id asc").Order("time_slot asc").Where("season_id = ? and is_preseason_game = ?", SeasonID, ts.NFLPreseason).Find(&games)

	return games
}

func GetNFLGamesByTimeslotAndWeekId(id string, timeslot string) []structs.NFLGame {
	db := dbprovider.GetInstance().GetDB()

	var games []structs.NFLGame

	db.Order("time_slot asc").Where("week_id = ? AND time_slot = ?", id, timeslot).Find(&games)

	return games
}

func GetTeamScheduleForBot(TeamID string, SeasonID string) []models.CollegeGameResponse {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	var games []structs.CollegeGame

	db.Order("week_id asc").Where("season_id = ? AND (home_team_id = ? OR away_team_id = ?)", SeasonID, TeamID, TeamID).Find(&games)

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

func GetCollegeGamesByTeamId(TeamID string) []structs.CollegeGame {
	db := dbprovider.GetInstance().GetDB()

	var games []structs.CollegeGame

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
	rand.Seed(time.Now().UnixNano())
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

	ts := GetTimestamp()

	ts.ToggleRunGames()

	db.Save(&ts)
}
