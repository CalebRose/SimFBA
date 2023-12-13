package managers

import (
	"log"
	"strconv"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/models"
	"github.com/CalebRose/SimFBA/structs"
)

// GetStandingsByConferenceIDAndSeasonID
func GetStandingsByConferenceIDAndSeasonID(conferenceID string, seasonID string) []structs.CollegeStandings {
	var standings []structs.CollegeStandings
	db := dbprovider.GetInstance().GetDB()
	err := db.Where("conference_id = ? AND season_id = ?", conferenceID, seasonID).Order("conference_losses asc").Order("conference_wins desc").
		Order("total_losses asc").Order("total_wins desc").
		Find(&standings).Error
	if err != nil {
		log.Fatal(err)
	}
	return standings
}

func GetNFLStandingsBySeasonID(seasonID string) []structs.NFLStandings {
	var standings []structs.NFLStandings
	db := dbprovider.GetInstance().GetDB()
	err := db.Where("season_id = ?", seasonID).Order("total_losses asc").Order("total_ties asc").Order("total_wins desc").
		Find(&standings).Error
	if err != nil {
		log.Fatal(err)
	}
	return standings
}

func GetNFLStandingsByTeamIDAndSeasonID(teamID string, seasonID string) structs.NFLStandings {
	var standings structs.NFLStandings
	db := dbprovider.GetInstance().GetDB()
	err := db.Where("team_id = ? AND season_id = ?", teamID, seasonID).Order("division_losses asc").Order("division_ties asc").Order("division_wins desc").
		Order("total_losses asc").Order("total_ties asc").Order("total_wins desc").
		Find(&standings).Error
	if err != nil {
		log.Fatal(err)
	}
	return standings
}

func GetNFLStandingsByDivisionIDAndSeasonID(divisionID string, seasonID string) []structs.NFLStandings {
	var standings []structs.NFLStandings
	db := dbprovider.GetInstance().GetDB()
	err := db.Where("division_id = ? AND season_id = ?", divisionID, seasonID).Order("division_losses asc").Order("division_ties asc").Order("division_wins desc").
		Order("total_losses asc").Order("total_ties asc").Order("total_wins desc").
		Find(&standings).Error
	if err != nil {
		log.Fatal(err)
	}
	return standings
}

// GetHistoricalRecordsByTeamID
func GetHistoricalRecordsByTeamID(TeamID string) models.TeamRecordResponse {
	tsChn := make(chan structs.Timestamp)

	go func() {
		ts := GetTimestamp()
		tsChn <- ts
	}()

	timestamp := <-tsChn
	close(tsChn)

	season := strconv.Itoa(timestamp.Season)

	historicGames := GetCollegeGamesByTeamId(TeamID)
	var conferenceChampionships []string
	var divisionTitles []string
	var nationalChampionships []string
	overallWins := 0
	overallLosses := 0
	currentSeasonWins := 0
	currentSeasonLosses := 0
	bowlWins := 0
	bowlLosses := 0

	for _, game := range historicGames {
		if !game.GameComplete || (game.GameComplete && game.SeasonID == timestamp.CollegeSeasonID && game.WeekID == timestamp.CollegeWeekID) {
			continue
		}

		isAway := strconv.Itoa(game.AwayTeamID) == TeamID

		if (isAway && game.AwayTeamWin) || (!isAway && game.HomeTeamWin) {
			overallWins++

			if game.SeasonID == timestamp.CollegeSeasonID {
				currentSeasonWins++
			}

			if game.IsBowlGame {
				bowlWins++
			}

			if game.IsConferenceChampionship {
				conferenceChampionships = append(conferenceChampionships, season)
				divisionTitles = append(divisionTitles, season)
			}

			if game.IsNationalChampionship {
				nationalChampionships = append(nationalChampionships, season)
			}
		} else {
			overallLosses++

			if game.SeasonID == timestamp.CollegeSeasonID {
				currentSeasonLosses++
			}

			if game.IsBowlGame {
				bowlLosses++
			}

			if game.IsConferenceChampionship {
				divisionTitles = append(divisionTitles, season)
			}
		}
	}

	response := models.TeamRecordResponse{
		OverallWins:             overallWins,
		OverallLosses:           overallLosses,
		CurrentSeasonWins:       currentSeasonWins,
		CurrentSeasonLosses:     currentSeasonLosses,
		BowlWins:                bowlWins,
		BowlLosses:              bowlLosses,
		ConferenceChampionships: conferenceChampionships,
		DivisionTitles:          divisionTitles,
		NationalChampionships:   nationalChampionships,
	}

	return response
}

// GetStandingsByConferenceIDAndSeasonID
func GetCFBStandingsByTeamIDAndSeasonID(TeamID string, seasonID string) structs.CollegeStandings {
	var standings structs.CollegeStandings
	db := dbprovider.GetInstance().GetDB()
	err := db.Where("team_id = ? AND season_id = ?", TeamID, seasonID).
		Find(&standings).Error
	if err != nil {
		log.Fatal(err)
	}
	return standings
}

func GetAllConferenceStandingsBySeasonID(seasonID string) []structs.CollegeStandings {
	db := dbprovider.GetInstance().GetDB()

	var standings []structs.CollegeStandings

	db.Where("season_id = ?", seasonID).Order("conference_id asc").Order("conference_losses asc").Order("conference_wins desc").
		Order("total_losses asc").Order("total_wins desc").Find(&standings)

	return standings
}

func GetAllNFLStandingsBySeasonID(seasonID string) []structs.NFLStandings {
	db := dbprovider.GetInstance().GetDB()

	var standings []structs.NFLStandings

	db.Where("season_id = ?", seasonID).Order("conference_id asc").Order("conference_losses asc").Order("conference_wins desc").
		Order("total_losses asc").Order("total_wins desc").Find(&standings)

	return standings
}

func GetCollegeStandingsRecordByTeamID(id string, seasonID string) structs.CollegeStandings {
	db := dbprovider.GetInstance().GetDB()

	var standing structs.CollegeStandings

	err := db.Where("team_id = ? AND season_id = ?", id, seasonID).Find(&standing).Error
	if err != nil {
		return structs.CollegeStandings{}
	}

	return standing
}
