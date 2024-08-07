package managers

import (
	"log"
	"strconv"
	"sync"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/models"
	"github.com/CalebRose/SimFBA/structs"
	"gorm.io/gorm"
)

func GetAllCollegeTeamsForRosterPage() models.RosterPageResponse {

	teams := GetAllCollegeTeams()
	coaches := GetAllCollegeCoaches()

	return models.RosterPageResponse{
		Teams:   teams,
		Coaches: coaches,
	}
}

func GetCollegeTeamMap() map[uint]structs.CollegeTeam {
	teams := GetAllCollegeTeams()
	teamMap := make(map[uint]structs.CollegeTeam)
	for _, t := range teams {
		teamMap[t.ID] = t
	}

	return teamMap
}

func GetAllCollegeTeams() []structs.CollegeTeam {
	db := dbprovider.GetInstance().GetDB()

	var teams []structs.CollegeTeam

	db.Find(&teams)

	return teams
}

func GetAllAvailableCollegeTeams() []structs.CollegeTeam {
	db := dbprovider.GetInstance().GetDB()

	var teams []structs.CollegeTeam

	db.Where("coach IN (?,?)", "", "AI").Find(&teams)

	return teams
}

func GetAllCoachedCollegeTeams() []structs.CollegeTeam {
	db := dbprovider.GetInstance().GetDB()

	var teams []structs.CollegeTeam

	db.Where("coach is not null AND coach NOT IN (?,?)", "", "AI").Find(&teams)

	return teams
}

// GetTeamByTeamID - straightforward
func GetTeamByTeamID(teamId string) structs.CollegeTeam {
	var team structs.CollegeTeam
	db := dbprovider.GetInstance().GetDB()
	err := db.Where("id = ?", teamId).Find(&team).Error
	if err != nil {
		log.Fatal(err)
	}
	return team
}

// GetTeamByTeamID - straightforward
func GetAllNFLTeams() []structs.NFLTeam {
	var teams []structs.NFLTeam
	db := dbprovider.GetInstance().GetDB()
	err := db.Preload("Capsheet").Order("team_name asc").Find(&teams).Error
	if err != nil {
		log.Fatal(err)
	}
	return teams
}

// Get NFL Records for Roster Page
func GetNFLRecordsForRosterPage(teamId string) models.NFLRosterPageResponse {

	team := GetNFLTeamWithCapsheetByTeamID(teamId)

	players := GetNFLPlayersForRosterPage(teamId)

	return models.NFLRosterPageResponse{
		Team:   team,
		Roster: players,
	}
}

// GetTeamByTeamID - straightforward
func GetNFLTeamByTeamID(teamId string) structs.NFLTeam {
	var team structs.NFLTeam
	db := dbprovider.GetInstance().GetDB()
	err := db.Where("id = ?", teamId).Find(&team).Error
	if err != nil {
		log.Fatal(err)
	}
	return team
}

// GetTeamByTeamID - straightforward
func GetNFLTeamWithCapsheetByTeamID(teamId string) structs.NFLTeam {
	var team structs.NFLTeam
	db := dbprovider.GetInstance().GetDB()
	err := db.Preload("Capsheet").Where("id = ?", teamId).Find(&team).Error
	if err != nil {
		log.Fatal(err)
	}
	return team
}

func GetTeamByTeamIDForDiscord(teamId string) structs.CollegeTeam {
	var team structs.CollegeTeam
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	err := db.Preload("TeamStandings", func(db *gorm.DB) *gorm.DB {
		return db.Where("season_id = ?", ts.CollegeSeasonID)
	}).Where("id = ?", teamId).Find(&team).Error
	if err != nil {
		log.Fatal(err)
	}
	return team
}

// GetTeamsByConferenceID
func GetTeamsByConferenceID(conferenceID string) []structs.CollegeTeam {
	var teams []structs.CollegeTeam
	db := dbprovider.GetInstance().GetDB()
	err := db.Where("conference_id = ?", conferenceID).Find(&teams).Error
	if err != nil {
		log.Fatal(err)
	}
	return teams
}

// GetTeamsByConferenceIDWithStandings
func GetTeamsByConferenceIDWithStandings(conferenceID string, seasonID string) []structs.CollegeTeam {
	var teams []structs.CollegeTeam
	db := dbprovider.GetInstance().GetDB()
	err := db.Preload("TeamStandings").
		Where("conference_id = ? AND TeamStandings.season_id = ?", conferenceID, seasonID).
		Find(&teams).Error
	if err != nil {
		log.Fatal(err)
	}
	return teams
}

// GetTeamsByDivisionID
func GetTeamsByDivisionID(conferenceID string) []structs.CollegeTeam {
	var teams []structs.CollegeTeam
	db := dbprovider.GetInstance().GetDB()
	err := db.Where("division_id = ?", conferenceID).Find(&teams).Error
	if err != nil {
		log.Fatal(err)
	}
	return teams
}

func GetTeamsInConference(conference string) []structs.CollegeTeam {
	db := dbprovider.GetInstance().GetDB()

	var teams []structs.CollegeTeam

	err := db.Where("conference = ?", conference).Find(&teams).Error
	if err != nil {
		log.Fatal(err)
	}
	return teams
}

func GetTeamByTeamAbbr(abbr string) structs.CollegeTeam {
	db := dbprovider.GetInstance().GetDB()
	var team structs.CollegeTeam
	err := db.Preload("TeamGameplan").Preload("TeamDepthChart.DepthChartPlayers").Where("team_abbr = ?", abbr).Find(&team).Error
	if err != nil {
		log.Panicln("Could not find team by given abbreviation:"+abbr+"\n", err)
	}

	return team
}

func GetNFLTeamByTeamIDForSim(id string) structs.NFLTeam {
	db := dbprovider.GetInstance().GetDB()

	var team structs.NFLTeam

	err := db.Preload("TeamGameplan").Preload("TeamDepthChart.DepthChartPlayers").Where("id = ?", id).Find(&team).Error
	if err != nil {
		log.Panicln("Could not find team by given id:"+id+"\n", err)
	}

	return team
}

func GetNFLTeamByTeamAbbr(abbr string) structs.NFLTeam {
	db := dbprovider.GetInstance().GetDB()

	var team structs.NFLTeam

	err := db.Preload("TeamGameplan").
		Preload("TeamDepthChart.DepthChartPlayers").Where("team_abbr = ?", abbr).Find(&team).Error
	if err != nil {
		log.Panicln("Could not find team by given abbreviation:"+abbr+"\n", err)
	}

	return team
}

func GetAllCollegeTeamsWithRecruitingProfileAndCoach() []structs.CollegeTeam {
	db := dbprovider.GetInstance().GetDB()

	var teams []structs.CollegeTeam

	db.Preload("CollegeCoach").Preload("RecruitingProfile").Find(&teams)

	return teams
}

func GetAllCollegeTeamsWithCurrentYearStandings() []structs.CollegeTeam {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	var teams []structs.CollegeTeam

	db.Preload("TeamStandings", func(db *gorm.DB) *gorm.DB {
		return db.Where("season_id = ?", strconv.Itoa(ts.CollegeSeasonID))
	}).Find(&teams)

	return teams
}

func GetCollegeConferences() []structs.CollegeConference {
	db := dbprovider.GetInstance().GetDB()

	var conferences []structs.CollegeConference

	db.Preload("Divisions").Find(&conferences)

	return conferences
}

func GetCollegeTeamsByConference(conf string) []structs.CollegeTeam {
	db := dbprovider.GetInstance().GetDB()

	var teams []structs.CollegeTeam

	db.Where("conference = ?", conf).Find(&teams)

	return teams
}

func GetDashboardByTeamID(isCFB bool, teamID string) structs.DashboardResponseData {
	ts := GetTimestamp()
	seasonID := strconv.Itoa(ts.CollegeSeasonID)
	collegeTeam := structs.CollegeTeam{}
	nflTeam := structs.NFLTeam{}
	if isCFB {
		collegeTeam = GetTeamByTeamID(teamID)
	} else {
		nflTeam = GetNFLTeamByTeamID(teamID)
	}
	cStandings := make(chan []structs.CollegeStandings)
	nStandings := make(chan []structs.NFLStandings)
	cGames := make(chan []structs.CollegeGame)
	nGames := make(chan []structs.NFLGame)
	newsChan := make(chan []structs.NewsLog)

	var waitGroup sync.WaitGroup
	waitGroup.Add(5)
	go func() {
		waitGroup.Wait()
		close(cStandings)
		close(nStandings)
		close(cGames)
		close(nGames)
		close(newsChan)
	}()

	go func() {
		defer waitGroup.Done()
		cSt := []structs.CollegeStandings{}
		if isCFB {
			cSt = GetStandingsByConferenceIDAndSeasonID(strconv.Itoa(collegeTeam.ConferenceID), seasonID)
		}
		cStandings <- cSt
	}()

	go func() {
		defer waitGroup.Done()
		nSt := []structs.NFLStandings{}
		if !isCFB {
			nSt = GetNFLStandingsByDivisionIDAndSeasonID(strconv.Itoa(int(nflTeam.DivisionID)), seasonID)
		}
		nStandings <- nSt
	}()

	go func() {
		defer waitGroup.Done()
		cG := []structs.CollegeGame{}
		if isCFB {
			cG = GetCollegeGamesByTeamIdAndSeasonId(teamID, seasonID)
		}
		cGames <- cG
	}()

	go func() {
		defer waitGroup.Done()
		nG := []structs.NFLGame{}
		if !isCFB {
			nG = GetNFLGamesByTeamIdAndSeasonId(teamID, seasonID)
		}
		nGames <- nG
	}()
	go func() {
		defer waitGroup.Done()
		nL := []structs.NewsLog{}
		if isCFB {
			nL = GetCFBRelatedNews(teamID)
		} else {
			nL = GetNFLRelatedNews(teamID)
		}
		newsChan <- nL
	}()

	collegeStandings := <-cStandings
	nflStandings := <-nStandings
	collegeGames := <-cGames
	nflGames := <-nGames
	newsLogs := <-newsChan

	return structs.DashboardResponseData{
		CollegeStandings: collegeStandings,
		NFLStandings:     nflStandings,
		CollegeGames:     collegeGames,
		NFLGames:         nflGames,
		NewsLogs:         newsLogs,
	}
}
