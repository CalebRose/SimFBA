package managers

import (
	"log"
	"strconv"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/models"
	"github.com/CalebRose/SimFBA/repository"
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
		if !game.GameComplete || (game.GameComplete && game.SeasonID == timestamp.CollegeSeasonID && game.WeekID == timestamp.CollegeWeekID) || game.IsSpringGame {
			continue
		}
		winningSeason := game.SeasonID + 2020
		winningSeasonStr := strconv.Itoa(winningSeason)
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
				conferenceChampionships = append(conferenceChampionships, winningSeasonStr)
				divisionTitles = append(divisionTitles, winningSeasonStr)
			}

			if game.IsNationalChampionship {
				nationalChampionships = append(nationalChampionships, winningSeasonStr)
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

func GetHistoricalNFLRecordsByTeamID(TeamID string) models.TeamRecordResponse {
	tsChn := make(chan structs.Timestamp)

	go func() {
		ts := GetTimestamp()
		tsChn <- ts
	}()

	timestamp := <-tsChn
	close(tsChn)

	historicGames := GetNFLGamesByTeamId(TeamID)
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
		if !game.GameComplete || (game.GameComplete && game.SeasonID == timestamp.CollegeSeasonID && game.WeekID == timestamp.CollegeWeekID) || game.IsPreseasonGame {
			continue
		}
		gameSeason := game.SeasonID + 2020
		isAway := strconv.Itoa(game.AwayTeamID) == TeamID

		if (isAway && game.AwayTeamWin) || (!isAway && game.HomeTeamWin) {
			overallWins++

			if game.SeasonID == timestamp.CollegeSeasonID {
				currentSeasonWins++
			}

			if game.IsPlayoffGame {
				bowlWins++
			}

			if game.IsConferenceChampionship {
				conferenceChampionships = append(conferenceChampionships, strconv.Itoa(gameSeason))
				divisionTitles = append(divisionTitles, strconv.Itoa(gameSeason))
			}

			if game.IsSuperBowl {
				nationalChampionships = append(nationalChampionships, strconv.Itoa(gameSeason))
			}
		} else {
			overallLosses++

			if game.SeasonID == timestamp.CollegeSeasonID {
				currentSeasonLosses++
			}

			if game.IsPlayoffGame {
				bowlLosses++
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

func GetAllCollegeStandingsBySeasonID(seasonID string) []structs.CollegeStandings {
	return repository.FindAllCollegeStandingsRecords(repository.StandingsQuery{
		SeasonID: seasonID,
	})
}

func GetAllNFLStandingsBySeasonID(seasonID string) []structs.NFLStandings {
	return repository.FindAllNFLStandingsRecords(repository.StandingsQuery{
		SeasonID: seasonID,
	})
}

func GetCollegeStandingsRecordByTeamID(id string, seasonID string) structs.CollegeStandings {
	return repository.FindAllCollegeStandingsRecords(repository.StandingsQuery{
		TeamID:   id,
		SeasonID: seasonID,
	})[0]
}

func ResetCollegeStandingsRanks() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.CollegeSeasonID))
	db.Model(&structs.CollegeStandings{}).Where("season_id = ?", seasonID).Updates(structs.CollegeStandings{Rank: 0})
}

func ResetCollegeStandings() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.CollegeSeasonID))
	collegeStandings := repository.FindAllCollegeStandingsRecords(repository.StandingsQuery{SeasonID: seasonID})
	for _, standings := range collegeStandings {
		standings.ResetCFBStandings()

		teamID := strconv.Itoa(standings.TeamID)
		games := GetCollegeGamesByTeamIdAndSeasonId(teamID, seasonID, false)
		for _, game := range games {
			if game.GameComplete {
				continue
			}
			standings.UpdateCollegeStandings(game)
		}
		repository.SaveCFBStandingsRecord(standings, db)
	}
}

func GetCollegeStandingsMap(seasonID string) map[uint]structs.CollegeStandings {
	standingsMap := make(map[uint]structs.CollegeStandings)

	standings := repository.FindAllCollegeStandingsRecords(repository.StandingsQuery{SeasonID: seasonID})
	for _, stat := range standings {
		standingsMap[uint(stat.TeamID)] = stat
	}

	return standingsMap
}

func GetStandingsHistoryByTeamID(id string) []structs.CollegeStandings {
	return repository.FindAllCollegeStandingsRecords(repository.StandingsQuery{
		TeamID: id,
	})
}

func GenerateNewSeasonStandings() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	teams := GetAllCollegeTeams()
	collegeStandings := []structs.CollegeStandings{}
	nflStandings := []structs.NFLStandings{}

	nflTeams := GetAllNFLTeams()

	for _, t := range teams {
		if !t.IsActive {
			continue
		}
		leagueID := 1
		league := "FBS"
		if !t.IsFBS {
			leagueID = 2
			league = "FCS"
		}

		standings := structs.CollegeStandings{
			TeamID:           int(t.ID),
			TeamName:         t.TeamName,
			SeasonID:         ts.CollegeSeasonID,
			Season:           ts.Season,
			ConferenceID:     t.ConferenceID,
			ConferenceName:   t.Conference,
			PostSeasonStatus: "None",
			IsFBS:            t.IsFBS,
			DivisionID:       t.DivisionID,
			LeagueID:         uint(leagueID),
			LeagueName:       league,
			BaseStandings: structs.BaseStandings{
				Coach:    t.Coach,
				TeamAbbr: t.TeamAbbr,
			},
		}

		collegeStandings = append(collegeStandings, standings)
	}
	repository.CreateCFBStandingsBatch(db, collegeStandings, 100)

	for _, t := range nflTeams {

		standings := structs.NFLStandings{
			TeamID:           t.ID,
			TeamName:         t.TeamName,
			SeasonID:         uint(ts.CollegeSeasonID),
			Season:           uint(ts.Season),
			ConferenceID:     t.ConferenceID,
			ConferenceName:   t.Conference,
			PostSeasonStatus: "None",
			DivisionID:       t.DivisionID,
			BaseStandings: structs.BaseStandings{
				Coach:    t.Coach,
				TeamAbbr: t.TeamAbbr,
			},
		}

		nflStandings = append(nflStandings, standings)
	}
	repository.CreateNFLStandingsBatch(db, nflStandings, 20)
}
