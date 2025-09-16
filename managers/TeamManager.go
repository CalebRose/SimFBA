package managers

import (
	"log"
	"math"
	"strconv"
	"strings"
	"sync"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/models"
	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
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

	db.Preload("CollegeCoach").Preload("RecruitingProfile.Affinities").Find(&teams)

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
	_, cfbGT := ts.GetCFBCurrentGameType()
	_, nflGT := ts.GetNFLCurrentGameType()
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
	cfbPlayerChan := make(chan []structs.CollegePlayerResponse)
	nflPlayerChan := make(chan []structs.NFLPlayerResponse)
	cfbTeamStatsChan := make(chan structs.CollegeTeamSeasonStats)
	nflTeamStatsChan := make(chan structs.NFLTeamSeasonStats)
	pollChan := make(chan structs.CollegePollOfficial)

	var waitGroup sync.WaitGroup
	waitGroup.Add(10)
	go func() {
		waitGroup.Wait()
		close(cStandings)
		close(nStandings)
		close(cGames)
		close(nGames)
		close(newsChan)
		close(cfbPlayerChan)
		close(nflPlayerChan)
		close(cfbTeamStatsChan)
		close(nflTeamStatsChan)
		close(pollChan)
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
			cG = GetCollegeGamesByTeamIdAndSeasonId(teamID, seasonID, ts.CFBSpringGames)
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

	go func() {
		defer waitGroup.Done()
		players := []structs.CollegePlayerResponse{}
		if isCFB {
			seasonKey := ts.CollegeSeasonID
			if ts.IsOffSeason {
				seasonKey -= 1
			}
			players = GetAllCollegePlayersWithSeasonStatsByTeamID(teamID, strconv.Itoa(seasonKey), cfbGT)
		}
		cfbPlayerChan <- players
	}()

	go func() {
		defer waitGroup.Done()
		players := []structs.NFLPlayerResponse{}
		if !isCFB {
			seasonKey := ts.NFLSeasonID
			if ts.IsNFLOffSeason {
				seasonKey -= 1
			}
			players = GetAllNFLPlayersWithSeasonStatsByTeamID(teamID, strconv.Itoa(seasonKey), nflGT)
		}
		nflPlayerChan <- players
	}()

	go func() {
		defer waitGroup.Done()
		stats := structs.CollegeTeamSeasonStats{}
		if isCFB {
			seasonKey := ts.CollegeSeasonID
			if ts.IsOffSeason {
				seasonKey -= 1
			}
			stats = GetCollegeTeamSeasonStatsBySeason(teamID, strconv.Itoa(seasonKey), cfbGT)
		}
		cfbTeamStatsChan <- stats
	}()

	go func() {
		defer waitGroup.Done()
		stats := structs.NFLTeamSeasonStats{}
		if !isCFB {
			seasonKey := ts.NFLSeasonID
			if ts.IsNFLOffSeason {
				seasonKey -= 1
			}
			stats = GetNFLTeamSeasonStatsByTeamANDSeason(teamID, strconv.Itoa(seasonKey), nflGT)
		}
		nflTeamStatsChan <- stats
	}()

	go func() {
		defer waitGroup.Done()
		poll := structs.CollegePollOfficial{}
		if isCFB {
			seasonKey := ts.NFLSeasonID
			if ts.IsOffSeason {
				seasonKey -= 1
			}
			polls := GetOfficialPollBySeasonID(strconv.Itoa(seasonKey))
			if len(polls) > 0 {
				poll = polls[len(polls)-1]
			}
		}
		pollChan <- poll
	}()

	collegeStandings := <-cStandings
	nflStandings := <-nStandings
	collegeGames := <-cGames
	nflGames := <-nGames
	newsLogs := <-newsChan
	collegePlayers := <-cfbPlayerChan
	nflPlayers := <-nflPlayerChan
	cfbTeamStats := <-cfbTeamStatsChan
	nflTeamStats := <-nflTeamStatsChan
	collegePoll := <-pollChan

	return structs.DashboardResponseData{
		CollegeStandings: collegeStandings,
		NFLStandings:     nflStandings,
		CollegeGames:     collegeGames,
		NFLGames:         nflGames,
		NewsLogs:         newsLogs,
		TopCFBPlayers:    collegePlayers,
		TopNFLPlayers:    nflPlayers,
		CFBTeamStats:     cfbTeamStats,
		NFLTeamStats:     nflTeamStats,
		TopTenPoll:       collegePoll,
	}
}

func GetCollegePlayer(depthChartPlayers structs.CollegeTeamDepthChart, position string, level int) structs.CollegePlayer {
	for _, player := range depthChartPlayers.DepthChartPlayers {
		if player.Position == position && util.ConvertStringToInt(player.PositionLevel) == level {
			return player.CollegePlayer
		}
	}

	return structs.CollegePlayer{}
}

func GetNFLPlayer(depthChartPlayers structs.NFLDepthChart, position string, level int) structs.NFLPlayer {
	for _, player := range depthChartPlayers.DepthChartPlayers {
		if player.Position == position && util.ConvertStringToInt(player.PositionLevel) == level {
			return player.NFLPlayer
		}
	}

	return structs.NFLPlayer{}
}

func GetKickReturnOverall(speed int, agility int) float64 {
	grade := (float64(speed) * 0.75) + (float64(agility) * 0.25)
	return grade
}

func GetPuntReturnOverall(speed int, agility int) float64 {
	grade := (float64(speed) * 0.25) + (float64(agility) * 0.75)
	return grade
}

// Returns the CFB team's numerical value for their entire offense
func OffenseGradeCFB(depthChartPlayers structs.CollegeTeamDepthChart, gameplan structs.CollegeGameplan) float64 {
	// Get overall values for all relevant positions
	qb1 := GetCollegePlayer(depthChartPlayers, "QB", 1)
	if qb1.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: QB1 was not found!!! Player ID: " + strconv.Itoa(int(qb1.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	rb1 := GetCollegePlayer(depthChartPlayers, "RB", 1)
	if rb1.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: RB1 was not found!!! Player ID: " + strconv.Itoa(int(rb1.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	rb2 := GetCollegePlayer(depthChartPlayers, "RB", 2)
	if rb2.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: RB2 was not found!!! Player ID: " + strconv.Itoa(int(rb2.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	rb3 := GetCollegePlayer(depthChartPlayers, "RB", 3)
	if rb3.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: RB3 was not found!!! Player ID: " + strconv.Itoa(int(rb3.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	fb1 := GetCollegePlayer(depthChartPlayers, "FB", 1)
	if fb1.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: FB1 was not found!!! Player ID: " + strconv.Itoa(int(fb1.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	te1 := GetCollegePlayer(depthChartPlayers, "TE", 1)
	if te1.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: TE1 was not found!!! Player ID: " + strconv.Itoa(int(te1.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	te2 := GetCollegePlayer(depthChartPlayers, "TE", 2)
	if te2.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: TE2 was not found!!! Player ID: " + strconv.Itoa(int(te2.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	wr1 := GetCollegePlayer(depthChartPlayers, "WR", 1)
	if wr1.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: WR1 was not found!!! Player ID: " + strconv.Itoa(int(wr1.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	wr2 := GetCollegePlayer(depthChartPlayers, "WR", 2)
	if wr2.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: WR2 was not found!!! Player ID: " + strconv.Itoa(int(wr2.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	wr3 := GetCollegePlayer(depthChartPlayers, "WR", 3)
	if wr3.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: WR3 was not found!!! Player ID: " + strconv.Itoa(int(wr3.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	wr4 := GetCollegePlayer(depthChartPlayers, "WR", 4)
	if wr4.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: WR4 was not found!!! Player ID: " + strconv.Itoa(int(wr4.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	wr5 := GetCollegePlayer(depthChartPlayers, "WR", 5)
	if wr5.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: WR5 was not found!!! Player ID: " + strconv.Itoa(int(wr5.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	lt1 := GetCollegePlayer(depthChartPlayers, "LT", 1)
	if lt1.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: LT1 was not found!!! Player ID: " + strconv.Itoa(int(lt1.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	lg1 := GetCollegePlayer(depthChartPlayers, "LG", 1)
	if lg1.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: LG1 was not found!!! Player ID: " + strconv.Itoa(int(lg1.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	c1 := GetCollegePlayer(depthChartPlayers, "C", 1)
	if c1.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: C1 was not found!!! Player ID: " + strconv.Itoa(int(c1.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	rg1 := GetCollegePlayer(depthChartPlayers, "RG", 1)
	if rg1.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: RG1 was not found!!! Player ID: " + strconv.Itoa(int(rg1.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	rt1 := GetCollegePlayer(depthChartPlayers, "RT", 1)
	if rt1.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: RT1 was not found!!! Player ID: " + strconv.Itoa(int(rt1.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}

	// If the player is a scheme fit, give them a bonus, if they are a bad fit, give them a malus
	// ENSURE TO TRANSLATE THEIR POSITION PROPERLY FOR THE SCHEME FIT!!!
	qb1Overall := ApplySchemeModifiers(float64(qb1.Overall), "QB", qb1.Archetype, gameplan.OffensiveScheme)
	rb1Overall := ApplySchemeModifiers(float64(rb1.Overall), "RB", rb1.Archetype, gameplan.OffensiveScheme)
	rb2Overall := ApplySchemeModifiers(float64(rb2.Overall), "RB", rb2.Archetype, gameplan.OffensiveScheme)
	rb3Overall := ApplySchemeModifiers(float64(rb3.Overall), "RB", rb3.Archetype, gameplan.OffensiveScheme)
	fb1Overall := ApplySchemeModifiers(float64(fb1.Overall), "FB", fb1.Archetype, gameplan.OffensiveScheme)
	te1Overall := ApplySchemeModifiers(float64(te1.Overall), "TE", te1.Archetype, gameplan.OffensiveScheme)
	te2Overall := ApplySchemeModifiers(float64(te2.Overall), "TE", te2.Archetype, gameplan.OffensiveScheme)
	wr1Overall := ApplySchemeModifiers(float64(wr1.Overall), "WR", wr1.Archetype, gameplan.OffensiveScheme)
	wr2Overall := ApplySchemeModifiers(float64(wr2.Overall), "WR", wr2.Archetype, gameplan.OffensiveScheme)
	wr3Overall := ApplySchemeModifiers(float64(wr3.Overall), "WR", wr3.Archetype, gameplan.OffensiveScheme)
	wr4Overall := ApplySchemeModifiers(float64(wr4.Overall), "WR", wr4.Archetype, gameplan.OffensiveScheme)
	wr5Overall := ApplySchemeModifiers(float64(wr5.Overall), "WR", wr5.Archetype, gameplan.OffensiveScheme)
	lt1Overall := ApplySchemeModifiers(float64(lt1.Overall), "OT", lt1.Archetype, gameplan.OffensiveScheme)
	lg1Overall := ApplySchemeModifiers(float64(lg1.Overall), "OG", lg1.Archetype, gameplan.OffensiveScheme)
	c1Overall := ApplySchemeModifiers(float64(c1.Overall), "C", c1.Archetype, gameplan.OffensiveScheme)
	rg1Overall := ApplySchemeModifiers(float64(rg1.Overall), "RG", rg1.Archetype, gameplan.OffensiveScheme)
	rt1Overall := ApplySchemeModifiers(float64(rt1.Overall), "RT", rt1.Archetype, gameplan.OffensiveScheme)

	// Depending on scheme, weight them
	qb1Overall = qb1Overall * GetOffensePositionGradeWeight("QB1", gameplan.OffensiveScheme)
	rb1Overall = rb1Overall * GetOffensePositionGradeWeight("RB1", gameplan.OffensiveScheme)
	rb2Overall = rb2Overall * GetOffensePositionGradeWeight("RB2", gameplan.OffensiveScheme)
	rb3Overall = rb3Overall * GetOffensePositionGradeWeight("RB3", gameplan.OffensiveScheme)
	fb1Overall = fb1Overall * GetOffensePositionGradeWeight("FB1", gameplan.OffensiveScheme)
	te1Overall = te1Overall * GetOffensePositionGradeWeight("TE1", gameplan.OffensiveScheme)
	te2Overall = te2Overall * GetOffensePositionGradeWeight("TE2", gameplan.OffensiveScheme)
	wr1Overall = wr1Overall * GetOffensePositionGradeWeight("WR1", gameplan.OffensiveScheme)
	wr2Overall = wr2Overall * GetOffensePositionGradeWeight("WR2", gameplan.OffensiveScheme)
	wr3Overall = wr3Overall * GetOffensePositionGradeWeight("WR3", gameplan.OffensiveScheme)
	wr4Overall = wr4Overall * GetOffensePositionGradeWeight("WR4", gameplan.OffensiveScheme)
	wr5Overall = wr5Overall * GetOffensePositionGradeWeight("WR5", gameplan.OffensiveScheme)
	lt1Overall = lt1Overall * GetOffensePositionGradeWeight("LT1", gameplan.OffensiveScheme)
	lg1Overall = lg1Overall * GetOffensePositionGradeWeight("LG1", gameplan.OffensiveScheme)
	c1Overall = c1Overall * GetOffensePositionGradeWeight("C1", gameplan.OffensiveScheme)
	rg1Overall = rg1Overall * GetOffensePositionGradeWeight("RG1", gameplan.OffensiveScheme)
	rt1Overall = rt1Overall * GetOffensePositionGradeWeight("RT1", gameplan.OffensiveScheme)

	// Sum them all up
	grade := qb1Overall + rb1Overall + rb2Overall + rb3Overall + fb1Overall + te1Overall + te2Overall + wr1Overall + wr2Overall + wr3Overall + wr4Overall + wr5Overall + lt1Overall + lg1Overall + c1Overall + rg1Overall + rt1Overall
	// Divide by 11.5 (offense weight normalization value)
	grade = grade / 11.5
	// return the resulting value
	return grade
}

// Returns the CFB team's numerical value for their entire offense
func DefenseGradeCFB(depthChartPlayers structs.CollegeTeamDepthChart, gameplan structs.CollegeGameplan) float64 {
	// Get overall values for all relevant positions
	le1 := GetCollegePlayer(depthChartPlayers, "LE", 1)
	if le1.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: LE1 was not found!!! Player ID: " + strconv.Itoa(int(le1.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	dt1 := GetCollegePlayer(depthChartPlayers, "DT", 1)
	if dt1.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: DT1 was not found!!! Player ID: " + strconv.Itoa(int(dt1.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	dt2 := GetCollegePlayer(depthChartPlayers, "DT", 2)
	if dt2.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: DT2 was not found!!! Player ID: " + strconv.Itoa(int(dt2.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	re1 := GetCollegePlayer(depthChartPlayers, "RE", 1)
	if re1.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: RE1 was not found!!! Player ID: " + strconv.Itoa(int(re1.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	lolb1 := GetCollegePlayer(depthChartPlayers, "LOLB", 1)
	if lolb1.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: LOLB1 was not found!!! Player ID: " + strconv.Itoa(int(lolb1.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	mlb1 := GetCollegePlayer(depthChartPlayers, "MLB", 1)
	if mlb1.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: MLB1 was not found!!! Player ID: " + strconv.Itoa(int(mlb1.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	mlb2 := GetCollegePlayer(depthChartPlayers, "MLB", 2)
	if mlb2.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: MLB2 was not found!!! Player ID: " + strconv.Itoa(int(mlb2.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	rolb1 := GetCollegePlayer(depthChartPlayers, "ROLB", 1)
	if rolb1.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: ROLB1 was not found!!! Player ID: " + strconv.Itoa(int(rolb1.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	cb1 := GetCollegePlayer(depthChartPlayers, "CB", 1)
	if cb1.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: CB1 was not found!!! Player ID: " + strconv.Itoa(int(cb1.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	cb2 := GetCollegePlayer(depthChartPlayers, "CB", 2)
	if cb2.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: CB2 was not found!!! Player ID: " + strconv.Itoa(int(cb2.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	cb3 := GetCollegePlayer(depthChartPlayers, "CB", 3)
	if cb3.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: CB3 was not found!!! Player ID: " + strconv.Itoa(int(cb3.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	cb4 := GetCollegePlayer(depthChartPlayers, "CB", 4)
	if cb4.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: CB4 was not found!!! Player ID: " + strconv.Itoa(int(cb4.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	fs1 := GetCollegePlayer(depthChartPlayers, "FS", 1)
	if fs1.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: FS1 was not found!!! Player ID: " + strconv.Itoa(int(fs1.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	ss1 := GetCollegePlayer(depthChartPlayers, "SS", 1)
	if ss1.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: SS1 was not found!!! Player ID: " + strconv.Itoa(int(ss1.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	ss2 := GetCollegePlayer(depthChartPlayers, "SS", 2)
	if ss2.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: SS2 was not found!!! Player ID: " + strconv.Itoa(int(ss2.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}

	// If the player is a scheme fit, give them a bonus, if they are a bad fit, give them a malus
	// ENSURE TO TRANSLATE THEIR POSITION PROPERLY FOR THE SCHEME FIT!!!
	le1Overall := ApplySchemeModifiers(float64(le1.Overall), "DE", le1.Archetype, gameplan.DefensiveScheme)
	dt1Overall := ApplySchemeModifiers(float64(dt1.Overall), "DT", dt1.Archetype, gameplan.DefensiveScheme)
	dt2Overall := ApplySchemeModifiers(float64(dt2.Overall), "DT", dt2.Archetype, gameplan.DefensiveScheme)
	re1Overall := ApplySchemeModifiers(float64(re1.Overall), "DE", re1.Archetype, gameplan.DefensiveScheme)
	lolb1Overall := ApplySchemeModifiers(float64(lolb1.Overall), "OLB", lolb1.Archetype, gameplan.DefensiveScheme)
	mlb1Overall := ApplySchemeModifiers(float64(mlb1.Overall), "ILB", mlb1.Archetype, gameplan.DefensiveScheme)
	mlb2Overall := ApplySchemeModifiers(float64(mlb2.Overall), "ILB", mlb2.Archetype, gameplan.DefensiveScheme)
	rolb1Overall := ApplySchemeModifiers(float64(rolb1.Overall), "OLB", rolb1.Archetype, gameplan.DefensiveScheme)
	cb1Overall := ApplySchemeModifiers(float64(cb1.Overall), "CB", cb1.Archetype, gameplan.DefensiveScheme)
	cb2Overall := ApplySchemeModifiers(float64(cb2.Overall), "CB", cb2.Archetype, gameplan.DefensiveScheme)
	cb3Overall := ApplySchemeModifiers(float64(cb3.Overall), "CB", cb3.Archetype, gameplan.DefensiveScheme)
	cb4Overall := ApplySchemeModifiers(float64(cb4.Overall), "CB", cb4.Archetype, gameplan.DefensiveScheme)
	fs1Overall := ApplySchemeModifiers(float64(fs1.Overall), "FS", fs1.Archetype, gameplan.DefensiveScheme)
	ss1Overall := ApplySchemeModifiers(float64(ss1.Overall), "SS", ss1.Archetype, gameplan.DefensiveScheme)
	ss2Overall := ApplySchemeModifiers(float64(ss2.Overall), "SS", ss2.Archetype, gameplan.DefensiveScheme)

	// Depending on scheme, weight them
	le1Overall = le1Overall * GetDefensePositionGradeWeight("LE1", gameplan.DefensiveScheme)
	dt1Overall = dt1Overall * GetDefensePositionGradeWeight("DT1", gameplan.DefensiveScheme)
	dt2Overall = dt2Overall * GetDefensePositionGradeWeight("DT2", gameplan.DefensiveScheme)
	re1Overall = re1Overall * GetDefensePositionGradeWeight("RE1", gameplan.DefensiveScheme)
	lolb1Overall = lolb1Overall * GetDefensePositionGradeWeight("LOLB1", gameplan.DefensiveScheme)
	mlb1Overall = mlb1Overall * GetDefensePositionGradeWeight("MLB1", gameplan.DefensiveScheme)
	mlb2Overall = mlb2Overall * GetDefensePositionGradeWeight("MLB2", gameplan.DefensiveScheme)
	rolb1Overall = rolb1Overall * GetDefensePositionGradeWeight("ROLB1", gameplan.DefensiveScheme)
	cb1Overall = cb1Overall * GetDefensePositionGradeWeight("CB1", gameplan.DefensiveScheme)
	cb2Overall = cb2Overall * GetDefensePositionGradeWeight("CB2", gameplan.DefensiveScheme)
	cb3Overall = cb3Overall * GetDefensePositionGradeWeight("CB3", gameplan.DefensiveScheme)
	cb4Overall = cb4Overall * GetDefensePositionGradeWeight("CB4", gameplan.DefensiveScheme)
	fs1Overall = fs1Overall * GetDefensePositionGradeWeight("FS1", gameplan.DefensiveScheme)
	ss1Overall = ss1Overall * GetDefensePositionGradeWeight("SS1", gameplan.DefensiveScheme)
	ss2Overall = ss2Overall * GetDefensePositionGradeWeight("SS2", gameplan.DefensiveScheme)

	// Sum them all up
	grade := le1Overall + dt1Overall + dt2Overall + re1Overall + lolb1Overall + mlb1Overall + mlb2Overall + rolb1Overall + cb1Overall + cb2Overall + cb3Overall + cb4Overall + fs1Overall + ss1Overall + ss2Overall
	// Divide by 11 (defense weight normalization value)
	grade = grade / 11.0
	// return the resulting value
	return grade
}

// Returns the CFB team's numerical value for their entire offense
func STGradeCFB(depthChartPlayers structs.CollegeTeamDepthChart) float64 {
	// Get overall values for all relevant positions
	k1 := GetCollegePlayer(depthChartPlayers, "K", 1)
	if k1.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: K1 was not found!!! Player ID: " + strconv.Itoa(int(k1.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	p1 := GetCollegePlayer(depthChartPlayers, "P", 1)
	if p1.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: p1 was not found!!! Player ID: " + strconv.Itoa(int(p1.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	fg1 := GetCollegePlayer(depthChartPlayers, "FG", 1)
	if fg1.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: FG1 was not found!!! Player ID: " + strconv.Itoa(int(fg1.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	kr1 := GetCollegePlayer(depthChartPlayers, "KR", 1)
	if kr1.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: KR1 was not found!!! Player ID: " + strconv.Itoa(int(kr1.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	kr2 := GetCollegePlayer(depthChartPlayers, "KR", 2)
	if kr2.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: KR2 was not found!!! Player ID: " + strconv.Itoa(int(kr2.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	pr1 := GetCollegePlayer(depthChartPlayers, "PR", 1)
	if pr1.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: PR1 was not found!!! Player ID: " + strconv.Itoa(int(pr1.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	pr2 := GetCollegePlayer(depthChartPlayers, "PR", 2)
	if pr1.PlayerID == 0 {
		log.Println("ERROR DURING COLLEGE TEAM GRADING: PR2 was not found!!! Player ID: " + strconv.Itoa(int(pr2.ID)) + " COLLEGE TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}

	// Weight them by position
	k1Overall := float64(k1.Overall)
	p1Overall := float64(p1.Overall)
	fg1Overall := float64(fg1.Overall)
	kr1Overall := GetKickReturnOverall(kr1.Speed, kr1.Agility) * 0.5
	kr2Overall := GetKickReturnOverall(kr2.Speed, kr2.Agility) * 0.5
	pr1Overall := GetPuntReturnOverall(pr1.Speed, pr1.Agility) * 0.5
	pr2Overall := GetPuntReturnOverall(pr2.Speed, pr2.Agility) * 0.5

	// Sum them all up
	grade := k1Overall + p1Overall + fg1Overall + kr1Overall + kr2Overall + pr1Overall + pr2Overall
	// Divide by 5 (Special Teams weight normalization value)
	grade = grade / 5.0
	// return the resulting value
	return grade
}

// Returns the CFB team's numerical value for their entire offense
func OffenseGradeNFL(depthChartPlayers structs.NFLDepthChart, gameplan structs.NFLGameplan) float64 {
	// Get overall values for all relevant positions
	qb1 := GetNFLPlayer(depthChartPlayers, "QB", 1)
	if qb1.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: QB1 was not found!!! Player ID: " + strconv.Itoa(int(qb1.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	rb1 := GetNFLPlayer(depthChartPlayers, "RB", 1)
	if rb1.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: RB1 was not found!!! Player ID: " + strconv.Itoa(int(rb1.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	rb2 := GetNFLPlayer(depthChartPlayers, "RB", 2)
	if rb2.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: RB2 was not found!!! Player ID: " + strconv.Itoa(int(rb2.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	rb3 := GetNFLPlayer(depthChartPlayers, "RB", 3)
	if rb3.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: RB3 was not found!!! Player ID: " + strconv.Itoa(int(rb3.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	fb1 := GetNFLPlayer(depthChartPlayers, "FB", 1)
	if fb1.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: FB1 was not found!!! Player ID: " + strconv.Itoa(int(fb1.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	te1 := GetNFLPlayer(depthChartPlayers, "TE", 1)
	if te1.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: TE1 was not found!!! Player ID: " + strconv.Itoa(int(te1.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	te2 := GetNFLPlayer(depthChartPlayers, "TE", 2)
	if te2.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: TE2 was not found!!! Player ID: " + strconv.Itoa(int(te2.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	wr1 := GetNFLPlayer(depthChartPlayers, "WR", 1)
	if wr1.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: WR1 was not found!!! Player ID: " + strconv.Itoa(int(wr1.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	wr2 := GetNFLPlayer(depthChartPlayers, "WR", 2)
	if wr2.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: WR2 was not found!!! Player ID: " + strconv.Itoa(int(wr2.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	wr3 := GetNFLPlayer(depthChartPlayers, "WR", 3)
	if wr3.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: WR3 was not found!!! Player ID: " + strconv.Itoa(int(wr3.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	wr4 := GetNFLPlayer(depthChartPlayers, "WR", 4)
	if wr4.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: WR4 was not found!!! Player ID: " + strconv.Itoa(int(wr4.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	wr5 := GetNFLPlayer(depthChartPlayers, "WR", 5)
	if wr5.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: WR5 was not found!!! Player ID: " + strconv.Itoa(int(wr5.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	lt1 := GetNFLPlayer(depthChartPlayers, "LT", 1)
	if lt1.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: LT1 was not found!!! Player ID: " + strconv.Itoa(int(lt1.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	lg1 := GetNFLPlayer(depthChartPlayers, "LG", 1)
	if lg1.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: LG1 was not found!!! Player ID: " + strconv.Itoa(int(lg1.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	c1 := GetNFLPlayer(depthChartPlayers, "C", 1)
	if c1.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: C1 was not found!!! Player ID: " + strconv.Itoa(int(c1.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	rg1 := GetNFLPlayer(depthChartPlayers, "RG", 1)
	if rg1.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: RG1 was not found!!! Player ID: " + strconv.Itoa(int(rg1.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	rt1 := GetNFLPlayer(depthChartPlayers, "RT", 1)
	if rt1.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: RT1 was not found!!! Player ID: " + strconv.Itoa(int(rt1.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}

	// If the player is a scheme fit, give them a bonus, if they are a bad fit, give them a malus
	// ENSURE TO TRANSLATE THEIR POSITION PROPERLY FOR THE SCHEME FIT!!!
	qb1Overall := ApplySchemeModifiers(float64(qb1.Overall), "QB", qb1.Archetype, gameplan.OffensiveScheme)
	rb1Overall := ApplySchemeModifiers(float64(rb1.Overall), "RB", rb1.Archetype, gameplan.OffensiveScheme)
	rb2Overall := ApplySchemeModifiers(float64(rb2.Overall), "RB", rb2.Archetype, gameplan.OffensiveScheme)
	rb3Overall := ApplySchemeModifiers(float64(rb3.Overall), "RB", rb3.Archetype, gameplan.OffensiveScheme)
	fb1Overall := ApplySchemeModifiers(float64(fb1.Overall), "FB", fb1.Archetype, gameplan.OffensiveScheme)
	te1Overall := ApplySchemeModifiers(float64(te1.Overall), "TE", te1.Archetype, gameplan.OffensiveScheme)
	te2Overall := ApplySchemeModifiers(float64(te2.Overall), "TE", te2.Archetype, gameplan.OffensiveScheme)
	wr1Overall := ApplySchemeModifiers(float64(wr1.Overall), "WR", wr1.Archetype, gameplan.OffensiveScheme)
	wr2Overall := ApplySchemeModifiers(float64(wr2.Overall), "WR", wr2.Archetype, gameplan.OffensiveScheme)
	wr3Overall := ApplySchemeModifiers(float64(wr3.Overall), "WR", wr3.Archetype, gameplan.OffensiveScheme)
	wr4Overall := ApplySchemeModifiers(float64(wr4.Overall), "WR", wr4.Archetype, gameplan.OffensiveScheme)
	wr5Overall := ApplySchemeModifiers(float64(wr5.Overall), "WR", wr5.Archetype, gameplan.OffensiveScheme)
	lt1Overall := ApplySchemeModifiers(float64(lt1.Overall), "OT", lt1.Archetype, gameplan.OffensiveScheme)
	lg1Overall := ApplySchemeModifiers(float64(lg1.Overall), "OG", lg1.Archetype, gameplan.OffensiveScheme)
	c1Overall := ApplySchemeModifiers(float64(c1.Overall), "C", c1.Archetype, gameplan.OffensiveScheme)
	rg1Overall := ApplySchemeModifiers(float64(rg1.Overall), "RG", rg1.Archetype, gameplan.OffensiveScheme)
	rt1Overall := ApplySchemeModifiers(float64(rt1.Overall), "RT", rt1.Archetype, gameplan.OffensiveScheme)

	// Depending on scheme, weight them
	qb1Overall = qb1Overall * GetOffensePositionGradeWeight("QB1", gameplan.OffensiveScheme)
	rb1Overall = rb1Overall * GetOffensePositionGradeWeight("RB1", gameplan.OffensiveScheme)
	rb2Overall = rb2Overall * GetOffensePositionGradeWeight("RB2", gameplan.OffensiveScheme)
	rb3Overall = rb3Overall * GetOffensePositionGradeWeight("RB3", gameplan.OffensiveScheme)
	fb1Overall = fb1Overall * GetOffensePositionGradeWeight("FB1", gameplan.OffensiveScheme)
	te1Overall = te1Overall * GetOffensePositionGradeWeight("TE1", gameplan.OffensiveScheme)
	te2Overall = te2Overall * GetOffensePositionGradeWeight("TE2", gameplan.OffensiveScheme)
	wr1Overall = wr1Overall * GetOffensePositionGradeWeight("WR1", gameplan.OffensiveScheme)
	wr2Overall = wr2Overall * GetOffensePositionGradeWeight("WR2", gameplan.OffensiveScheme)
	wr3Overall = wr3Overall * GetOffensePositionGradeWeight("WR3", gameplan.OffensiveScheme)
	wr4Overall = wr4Overall * GetOffensePositionGradeWeight("WR4", gameplan.OffensiveScheme)
	wr5Overall = wr5Overall * GetOffensePositionGradeWeight("WR5", gameplan.OffensiveScheme)
	lt1Overall = lt1Overall * GetOffensePositionGradeWeight("LT1", gameplan.OffensiveScheme)
	lg1Overall = lg1Overall * GetOffensePositionGradeWeight("LG1", gameplan.OffensiveScheme)
	c1Overall = c1Overall * GetOffensePositionGradeWeight("C1", gameplan.OffensiveScheme)
	rg1Overall = rg1Overall * GetOffensePositionGradeWeight("RG1", gameplan.OffensiveScheme)
	rt1Overall = rt1Overall * GetOffensePositionGradeWeight("RT1", gameplan.OffensiveScheme)

	// Sum them all up
	grade := qb1Overall + rb1Overall + rb2Overall + rb3Overall + fb1Overall + te1Overall + te2Overall + wr1Overall + wr2Overall + wr3Overall + wr4Overall + wr5Overall + lt1Overall + lg1Overall + c1Overall + rg1Overall + rt1Overall
	// Divide by 11.5 (offense weight normalization value)
	grade = grade / 11.5
	// return the resulting value
	return grade
}

// Returns the CFB team's numerical value for their entire offense
func DefenseGradeNFL(depthChartPlayers structs.NFLDepthChart, gameplan structs.NFLGameplan) float64 {
	// Get overall values for all relevant positions
	le1 := GetNFLPlayer(depthChartPlayers, "LE", 1)
	if le1.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: LE1 was not found!!! Player ID: " + strconv.Itoa(int(le1.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	dt1 := GetNFLPlayer(depthChartPlayers, "DT", 1)
	if dt1.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: DT1 was not found!!! Player ID: " + strconv.Itoa(int(dt1.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	dt2 := GetNFLPlayer(depthChartPlayers, "DT", 2)
	if dt2.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: DT2 was not found!!! Player ID: " + strconv.Itoa(int(dt2.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	re1 := GetNFLPlayer(depthChartPlayers, "RE", 1)
	if re1.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: RE1 was not found!!! Player ID: " + strconv.Itoa(int(re1.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	lolb1 := GetNFLPlayer(depthChartPlayers, "LOLB", 1)
	if lolb1.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: LOLB1 was not found!!! Player ID: " + strconv.Itoa(int(lolb1.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	mlb1 := GetNFLPlayer(depthChartPlayers, "MLB", 1)
	if mlb1.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: MLB1 was not found!!! Player ID: " + strconv.Itoa(int(mlb1.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	mlb2 := GetNFLPlayer(depthChartPlayers, "MLB", 2)
	if mlb2.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: MLB2 was not found!!! Player ID: " + strconv.Itoa(int(mlb2.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	rolb1 := GetNFLPlayer(depthChartPlayers, "ROLB", 1)
	if rolb1.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: ROLB1 was not found!!! Player ID: " + strconv.Itoa(int(rolb1.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	cb1 := GetNFLPlayer(depthChartPlayers, "CB", 1)
	if cb1.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: CB1 was not found!!! Player ID: " + strconv.Itoa(int(cb1.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	cb2 := GetNFLPlayer(depthChartPlayers, "CB", 2)
	if cb2.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: CB2 was not found!!! Player ID: " + strconv.Itoa(int(cb2.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	cb3 := GetNFLPlayer(depthChartPlayers, "CB", 3)
	if cb3.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: CB3 was not found!!! Player ID: " + strconv.Itoa(int(cb3.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	cb4 := GetNFLPlayer(depthChartPlayers, "CB", 4)
	if cb4.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: CB4 was not found!!! Player ID: " + strconv.Itoa(int(cb4.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	fs1 := GetNFLPlayer(depthChartPlayers, "FS", 1)
	if fs1.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: FS1 was not found!!! Player ID: " + strconv.Itoa(int(fs1.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	ss1 := GetNFLPlayer(depthChartPlayers, "SS", 1)
	if ss1.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: SS1 was not found!!! Player ID: " + strconv.Itoa(int(ss1.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	ss2 := GetNFLPlayer(depthChartPlayers, "SS", 2)
	if ss2.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: SS2 was not found!!! Player ID: " + strconv.Itoa(int(ss2.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}

	// If the player is a scheme fit, give them a bonus, if they are a bad fit, give them a malus
	// ENSURE TO TRANSLATE THEIR POSITION PROPERLY FOR THE SCHEME FIT!!!
	le1Overall := ApplySchemeModifiers(float64(le1.Overall), "DE", le1.Archetype, gameplan.DefensiveScheme)
	dt1Overall := ApplySchemeModifiers(float64(dt1.Overall), "DT", dt1.Archetype, gameplan.DefensiveScheme)
	dt2Overall := ApplySchemeModifiers(float64(dt2.Overall), "DT", dt2.Archetype, gameplan.DefensiveScheme)
	re1Overall := ApplySchemeModifiers(float64(re1.Overall), "DE", re1.Archetype, gameplan.DefensiveScheme)
	lolb1Overall := ApplySchemeModifiers(float64(lolb1.Overall), "OLB", lolb1.Archetype, gameplan.DefensiveScheme)
	mlb1Overall := ApplySchemeModifiers(float64(mlb1.Overall), "ILB", mlb1.Archetype, gameplan.DefensiveScheme)
	mlb2Overall := ApplySchemeModifiers(float64(mlb2.Overall), "ILB", mlb2.Archetype, gameplan.DefensiveScheme)
	rolb1Overall := ApplySchemeModifiers(float64(rolb1.Overall), "OLB", rolb1.Archetype, gameplan.DefensiveScheme)
	cb1Overall := ApplySchemeModifiers(float64(cb1.Overall), "CB", cb1.Archetype, gameplan.DefensiveScheme)
	cb2Overall := ApplySchemeModifiers(float64(cb2.Overall), "CB", cb2.Archetype, gameplan.DefensiveScheme)
	cb3Overall := ApplySchemeModifiers(float64(cb3.Overall), "CB", cb3.Archetype, gameplan.DefensiveScheme)
	cb4Overall := ApplySchemeModifiers(float64(cb4.Overall), "CB", cb4.Archetype, gameplan.DefensiveScheme)
	fs1Overall := ApplySchemeModifiers(float64(fs1.Overall), "FS", fs1.Archetype, gameplan.DefensiveScheme)
	ss1Overall := ApplySchemeModifiers(float64(ss1.Overall), "SS", ss1.Archetype, gameplan.DefensiveScheme)
	ss2Overall := ApplySchemeModifiers(float64(ss2.Overall), "SS", ss2.Archetype, gameplan.DefensiveScheme)

	// Depending on scheme, weight them
	le1Overall = le1Overall * GetDefensePositionGradeWeight("LE1", gameplan.DefensiveScheme)
	dt1Overall = dt1Overall * GetDefensePositionGradeWeight("DT1", gameplan.DefensiveScheme)
	dt2Overall = dt2Overall * GetDefensePositionGradeWeight("DT2", gameplan.DefensiveScheme)
	re1Overall = re1Overall * GetDefensePositionGradeWeight("RE1", gameplan.DefensiveScheme)
	lolb1Overall = lolb1Overall * GetDefensePositionGradeWeight("LOLB1", gameplan.DefensiveScheme)
	mlb1Overall = mlb1Overall * GetDefensePositionGradeWeight("MLB1", gameplan.DefensiveScheme)
	mlb2Overall = mlb2Overall * GetDefensePositionGradeWeight("MLB2", gameplan.DefensiveScheme)
	rolb1Overall = rolb1Overall * GetDefensePositionGradeWeight("ROLB1", gameplan.DefensiveScheme)
	cb1Overall = cb1Overall * GetDefensePositionGradeWeight("CB1", gameplan.DefensiveScheme)
	cb2Overall = cb2Overall * GetDefensePositionGradeWeight("CB2", gameplan.DefensiveScheme)
	cb3Overall = cb3Overall * GetDefensePositionGradeWeight("CB3", gameplan.DefensiveScheme)
	cb4Overall = cb4Overall * GetDefensePositionGradeWeight("CB4", gameplan.DefensiveScheme)
	fs1Overall = fs1Overall * GetDefensePositionGradeWeight("FS1", gameplan.DefensiveScheme)
	ss1Overall = ss1Overall * GetDefensePositionGradeWeight("SS1", gameplan.DefensiveScheme)
	ss2Overall = ss2Overall * GetDefensePositionGradeWeight("SS2", gameplan.DefensiveScheme)

	// Sum them all up
	grade := le1Overall + dt1Overall + dt2Overall + re1Overall + lolb1Overall + mlb1Overall + mlb2Overall + rolb1Overall + cb1Overall + cb2Overall + cb3Overall + cb4Overall + fs1Overall + ss1Overall + ss2Overall
	// Divide by 11 (defense weight normalization value)
	grade = grade / 11.0
	// return the resulting value
	return grade
}

// Returns the CFB team's numerical value for their entire offense
func STGradeNFL(depthChartPlayers structs.NFLDepthChart) float64 {
	// Get overall values for all relevant positions
	k1 := GetNFLPlayer(depthChartPlayers, "K", 1)
	if k1.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: K1 was not found!!! Player ID: " + strconv.Itoa(int(k1.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	p1 := GetNFLPlayer(depthChartPlayers, "P", 1)
	if p1.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: p1 was not found!!! Player ID: " + strconv.Itoa(int(p1.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	fg1 := GetNFLPlayer(depthChartPlayers, "FG", 1)
	if fg1.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: FG1 was not found!!! Player ID: " + strconv.Itoa(int(fg1.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	kr1 := GetNFLPlayer(depthChartPlayers, "KR", 1)
	if kr1.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: KR1 was not found!!! Player ID: " + strconv.Itoa(int(kr1.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	kr2 := GetNFLPlayer(depthChartPlayers, "KR", 2)
	if kr2.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: KR2 was not found!!! Player ID: " + strconv.Itoa(int(kr2.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	pr1 := GetNFLPlayer(depthChartPlayers, "PR", 1)
	if pr1.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: PR1 was not found!!! Player ID: " + strconv.Itoa(int(pr1.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}
	pr2 := GetNFLPlayer(depthChartPlayers, "PR", 2)
	if pr1.PlayerID == 0 {
		log.Println("ERROR DURING NFL TEAM GRADING: PR2 was not found!!! Player ID: " + strconv.Itoa(int(pr2.ID)) + " NFL TEAM ID: " + strconv.Itoa(int(depthChartPlayers.TeamID)))
	}

	// Weight them by position
	k1Overall := float64(k1.Overall)
	p1Overall := float64(p1.Overall)
	fg1Overall := float64(fg1.Overall)
	kr1Overall := GetKickReturnOverall(kr1.Speed, kr1.Agility) * 0.5
	kr2Overall := GetKickReturnOverall(kr2.Speed, kr2.Agility) * 0.5
	pr1Overall := GetPuntReturnOverall(pr1.Speed, pr1.Agility) * 0.5
	pr2Overall := GetPuntReturnOverall(pr2.Speed, pr2.Agility) * 0.5

	// Sum them all up
	grade := k1Overall + p1Overall + fg1Overall + kr1Overall + kr2Overall + pr1Overall + pr2Overall
	// Divide by 5 (Special Teams weight normalization value)
	grade = grade / 5.0
	// return the resulting value
	return grade
}

// League agnostic
func OverallGrade(offense float64, defense float64, specialTeams float64) float64 {
	var overallGrade float64 = offense * 0.45
	overallGrade = overallGrade + (defense * 0.45)
	overallGrade = overallGrade + (specialTeams * 0.1)
	return overallGrade
}

func ApplySchemeModifiers(overall float64, position string, archetype string, scheme string) float64 {
	if IsSchemeFit(position, archetype, scheme) {
		return (overall * 1.05)
	} else if IsBadFit(position, archetype, scheme) {
		return (overall * 0.95)
	} else {
		return overall
	}
}

func IsSchemeFit(position string, archetype string, scheme string) bool {
	scheme = strings.ToLower(scheme)
	archetype = strings.ToLower(archetype)
	switch position {
	case "QB":
		switch archetype {
		case "scrambler":
			if strings.Contains(scheme, "option") || strings.Contains(scheme, "flexbone") {
				return true
			} else {
				return false
			}
		case "balanced":
			if strings.Contains(scheme, "wing-t") || strings.Contains(scheme, "pistol") || strings.Contains(scheme, "wishbone") {
				return true
			} else {
				return false
			}
		case "pocket":
			if strings.Contains(scheme, "raid") || strings.Contains(scheme, "vert") || strings.Contains(scheme, "pistol") {
				return true
			} else {
				return false
			}
		case "field":
			if strings.Contains(scheme, "west") || strings.Contains(scheme, "shoot") || strings.Contains(scheme, "wish") {
				return true
			} else {
				return false
			}
		default:
			return false
		}
	case "RB":
		switch archetype {
		case "balanced":
			if strings.Contains(scheme, "pistol") || strings.Contains(scheme, "wing-t") || strings.Contains(scheme, "wish") {
				return true
			} else {
				return false
			}
		case "receiving":
			if strings.Contains(scheme, "vert") || strings.Contains(scheme, "shoot") || strings.Contains(scheme, "raid") {
				return true
			} else {
				return false
			}
		case "speed":
			if strings.Contains(scheme, "spread") || strings.Contains(scheme, "shoot") || strings.Contains(scheme, "flex") {
				return true
			} else {
				return false
			}
		case "power":
			if strings.Contains(scheme, "power") || strings.Contains(scheme, "i option") || strings.Contains(scheme, "double") {
				return true
			} else {
				return false
			}
		default:
			return false
		}
	case "FB":
		switch archetype {
		case "balanced":
			if strings.Contains(scheme, "west") || strings.Contains(scheme, "wing-t") || strings.Contains(scheme, "flex") {
				return true
			} else {
				return false
			}
		case "receiving":
			if strings.Contains(scheme, "west") || strings.Contains(scheme, "spread") || strings.Contains(scheme, "vert") {
				return true
			} else {
				return false
			}
		case "rushing":
			if strings.Contains(scheme, "i option") || strings.Contains(scheme, "pistol") || strings.Contains(scheme, "double") {
				return true
			} else {
				return false
			}
		case "blocking":
			if strings.Contains(scheme, "power") || strings.Contains(scheme, "double") {
				return true
			} else {
				return false
			}
		default:
			return false
		}
	case "WR":
		switch archetype {
		case "route runner":
			if strings.Contains(scheme, "west") || strings.Contains(scheme, "vert") || strings.Contains(scheme, "pistol") || strings.Contains(scheme, "spread") {
				return true
			} else {
				return false
			}
		case "red zone threat":
			if strings.Contains(scheme, "power") || strings.Contains(scheme, "double") || strings.Contains(scheme, "bone") {
				return true
			} else {
				return false
			}
		case "possession":
			if strings.Contains(scheme, "option") || strings.Contains(scheme, "pistol") || strings.Contains(scheme, "west") {
				return true
			} else {
				return false
			}
		case "possesion":
			if strings.Contains(scheme, "option") || strings.Contains(scheme, "pistol") || strings.Contains(scheme, "west") {
				return true
			} else {
				return false
			}
		case "speed":
			if strings.Contains(scheme, "vert") || strings.Contains(scheme, "shoot") || strings.Contains(scheme, "raid") || strings.Contains(scheme, "wing-t") {
				return true
			} else {
				return false
			}
		default:
			return false
		}
	case "TE":
		switch archetype {
		case "vertical threat":
			if strings.Contains(scheme, "vert") || strings.Contains(scheme, "raid") || strings.Contains(scheme, "pistol") {
				return true
			} else {
				return false
			}
		case "receiving":
			if strings.Contains(scheme, "vert") || strings.Contains(scheme, "west") || strings.Contains(scheme, "raid") {
				return true
			} else {
				return false
			}
		case "blocking":
			if strings.Contains(scheme, "power") || strings.Contains(scheme, "double") || strings.Contains(scheme, "i option") {
				return true
			} else {
				return false
			}
		default:
			return false
		}
	case "OT":
		switch archetype {
		case "run blocking":
			if strings.Contains(scheme, "power") || strings.Contains(scheme, "double") {
				return true
			} else {
				return false
			}
		case "pass blocking":
			if strings.Contains(scheme, "vert") || strings.Contains(scheme, "raid") {
				return true
			} else {
				return false
			}
		default:
			return false
		}
	case "OG":
		switch archetype {
		case "run blocking":
			if strings.Contains(scheme, "power") || strings.Contains(scheme, "double") {
				return true
			} else {
				return false
			}
		case "pass blocking":
			if strings.Contains(scheme, "vert") || strings.Contains(scheme, "raid") {
				return true
			} else {
				return false
			}
		default:
			return false
		}
	case "C":
		switch archetype {
		case "run blocking":
			if strings.Contains(scheme, "power") || strings.Contains(scheme, "double") {
				return true
			} else {
				return false
			}
		case "pass blocking":
			if strings.Contains(scheme, "vert") || strings.Contains(scheme, "raid") {
				return true
			} else {
				return false
			}
		case "line captain":
			if strings.Contains(scheme, "west") || strings.Contains(scheme, "shoot") {
				return true
			} else {
				return false
			}
		default:
			return false
		}
	case "DE":
		switch archetype {
		case "speed rusher":
			if strings.Contains(scheme, "4") || strings.Contains(scheme, "speed") {
				return true
			} else {
				return false
			}
		case "run stopper":
			if strings.Contains(scheme, "2") || strings.Contains(scheme, "multiple") {
				return true
			} else {
				return false
			}
		default:
			return false
		}
	case "DT":
		switch archetype {
		case "pass rusher":
			if strings.Contains(scheme, "4") || strings.Contains(scheme, "speed") {
				return true
			} else {
				return false
			}
		case "nose tackle":
			if strings.Contains(scheme, "2") || strings.Contains(scheme, "3") {
				return true
			} else {
				return false
			}
		default:
			return false
		}
	case "OLB":
		switch archetype {
		case "pass rush":
			if strings.Contains(scheme, "2") || strings.Contains(scheme, "3") {
				return true
			} else {
				return false
			}
		case "speed":
			if strings.Contains(scheme, "speed") || strings.Contains(scheme, "multiple") {
				return true
			} else {
				return false
			}
		case "coverage":
			if strings.Contains(scheme, "4") || strings.Contains(scheme, "speed") {
				return true
			} else {
				return false
			}
		case "run stopper":
			if strings.Contains(scheme, "old") || strings.Contains(scheme, "2") {
				return true
			} else {
				return false
			}
		default:
			return false
		}
	case "ILB":
		switch archetype {
		case "field general":
			if strings.Contains(scheme, "old") || strings.Contains(scheme, "multiple") {
				return true
			} else {
				return false
			}
		case "speed":
			if strings.Contains(scheme, "speed") || strings.Contains(scheme, "multiple") {
				return true
			} else {
				return false
			}
		case "coverage":
			if strings.Contains(scheme, "4") || strings.Contains(scheme, "3") {
				return true
			} else {
				return false
			}
		case "run stopper":
			if strings.Contains(scheme, "old") || strings.Contains(scheme, "2") {
				return true
			} else {
				return false
			}
		default:
			return false
		}
	default:
		return false
	}
}

func IsBadFit(position string, archetype string, scheme string) bool {
	scheme = strings.ToLower(scheme)
	archetype = strings.ToLower(archetype)
	switch position {
	case "QB":
		switch archetype {
		case "scrambler":
			if strings.Contains(scheme, "raid") || strings.Contains(scheme, "vert") {
				return true
			} else {
				return false
			}
		case "balanced":
			if strings.Contains(scheme, "vert") || strings.Contains(scheme, "raid") {
				return true
			} else {
				return false
			}
		case "pocket":
			if strings.Contains(scheme, "i option") || strings.Contains(scheme, "flex") || strings.Contains(scheme, "double") {
				return true
			} else {
				return false
			}
		case "field":
			if strings.Contains(scheme, "vert") || strings.Contains(scheme, "raid") {
				return true
			} else {
				return false
			}
		default:
			return false
		}
	case "RB":
		switch archetype {
		case "balanced":
			if strings.Contains(scheme, "shoot") || strings.Contains(scheme, "spread") || strings.Contains(scheme, "flex") {
				return true
			} else {
				return false
			}
		case "receiving":
			if strings.Contains(scheme, "power") || strings.Contains(scheme, "i option") || strings.Contains(scheme, "double") {
				return true
			} else {
				return false
			}
		case "speed":
			if strings.Contains(scheme, "power") || strings.Contains(scheme, "i option") || strings.Contains(scheme, "double") {
				return true
			} else {
				return false
			}
		case "power":
			if strings.Contains(scheme, "vert") || strings.Contains(scheme, "shoot") || strings.Contains(scheme, "raid") {
				return true
			} else {
				return false
			}
		default:
			return false
		}
	case "FB":
		switch archetype {
		case "balanced":
			if strings.Contains(scheme, "pistol") || strings.Contains(scheme, "spread") || strings.Contains(scheme, "wish") {
				return true
			} else {
				return false
			}
		case "receiving":
			if strings.Contains(scheme, "power") || strings.Contains(scheme, "i option") || strings.Contains(scheme, "double") {
				return true
			} else {
				return false
			}
		case "rushing":
			if strings.Contains(scheme, "vert") || strings.Contains(scheme, "shoot") || strings.Contains(scheme, "raid") {
				return true
			} else {
				return false
			}
		case "blocking":
			if strings.Contains(scheme, "vert") || strings.Contains(scheme, "west") || strings.Contains(scheme, "shoot") || strings.Contains(scheme, "raid") {
				return true
			} else {
				return false
			}
		default:
			return false
		}
	case "WR":
		switch archetype {
		case "route runner":
			if strings.Contains(scheme, "wish") {
				return true
			} else {
				return false
			}
		case "red zone threat":
			if strings.Contains(scheme, "vert") || strings.Contains(scheme, "west") {
				return true
			} else {
				return false
			}
		case "possession":
			if strings.Contains(scheme, "shoot") || strings.Contains(scheme, "flex") {
				return true
			} else {
				return false
			}
		case "possesion":
			if strings.Contains(scheme, "shoot") || strings.Contains(scheme, "flex") {
				return true
			} else {
				return false
			}
		case "speed":
			if strings.Contains(scheme, "flex") {
				return true
			} else {
				return false
			}
		default:
			return false
		}
	case "TE":
		switch archetype {
		case "vertical threat":
			if strings.Contains(scheme, "power") || strings.Contains(scheme, "i option") || strings.Contains(scheme, "double") {
				return true
			} else {
				return false
			}
		case "receiving":
			if strings.Contains(scheme, "power") || strings.Contains(scheme, "i option") || strings.Contains(scheme, "double") {
				return true
			} else {
				return false
			}
		case "blocking":
			if strings.Contains(scheme, "vert") || strings.Contains(scheme, "shoot") || strings.Contains(scheme, "raid") {
				return true
			} else {
				return false
			}
		default:
			return false
		}
	case "OT":
		switch archetype {
		case "run blocking":
			if strings.Contains(scheme, "vert") || strings.Contains(scheme, "raid") {
				return true
			} else {
				return false
			}
		case "pass blocking":
			if strings.Contains(scheme, "power") || strings.Contains(scheme, "double") {
				return true
			} else {
				return false
			}
		default:
			return false
		}
	case "OG":
		switch archetype {
		case "run blocking":
			if strings.Contains(scheme, "vert") || strings.Contains(scheme, "raid") {
				return true
			} else {
				return false
			}
		case "pass blocking":
			if strings.Contains(scheme, "power") || strings.Contains(scheme, "double") {
				return true
			} else {
				return false
			}
		default:
			return false
		}
	case "C":
		switch archetype {
		case "run blocking":
			if strings.Contains(scheme, "vert") || strings.Contains(scheme, "raid") {
				return true
			} else {
				return false
			}
		case "pass blocking":
			if strings.Contains(scheme, "power") || strings.Contains(scheme, "double") {
				return true
			} else {
				return false
			}
		case "line captain":
			if strings.Contains(scheme, "pistol") || strings.Contains(scheme, "wish") {
				return true
			} else {
				return false
			}
		default:
			return false
		}
	case "DE":
		switch archetype {
		case "speed rusher":
			if strings.Contains(scheme, "2") || strings.Contains(scheme, "multiple") {
				return true
			} else {
				return false
			}
		case "run stopper":
			if strings.Contains(scheme, "4") || strings.Contains(scheme, "speed") {
				return true
			} else {
				return false
			}
		default:
			return false
		}
	case "DT":
		switch archetype {
		case "pass rusher":
			if strings.Contains(scheme, "2") || strings.Contains(scheme, "multiple") {
				return true
			} else {
				return false
			}
		case "nose tackle":
			if strings.Contains(scheme, "4") || strings.Contains(scheme, "speed") {
				return true
			} else {
				return false
			}
		default:
			return false
		}
	case "OLB":
		switch archetype {
		case "pass rush":
			if strings.Contains(scheme, "speed") {
				return true
			} else {
				return false
			}
		case "speed":
			if strings.Contains(scheme, "2") || strings.Contains(scheme, "3") {
				return true
			} else {
				return false
			}
		case "coverage":
			if strings.Contains(scheme, "old") || strings.Contains(scheme, "multiple") {
				return true
			} else {
				return false
			}
		case "run stopper":
			if strings.Contains(scheme, "4") || strings.Contains(scheme, "3") {
				return true
			} else {
				return false
			}
		default:
			return false
		}
	case "ILB":
		switch archetype {
		case "field general":
			if strings.Contains(scheme, "3") || strings.Contains(scheme, "speed") {
				return true
			} else {
				return false
			}
		case "speed":
			if strings.Contains(scheme, "2") || strings.Contains(scheme, "3") {
				return true
			} else {
				return false
			}
		case "coverage":
			if strings.Contains(scheme, "old") || strings.Contains(scheme, "multiple") {
				return true
			} else {
				return false
			}
		case "run stopper":
			if strings.Contains(scheme, "4") || strings.Contains(scheme, "3") {
				return true
			} else {
				return false
			}
		default:
			return false
		}
	default:
		return false
	}
}

func GetOffensePositionGradeWeight(position string, scheme string) float64 {
	scheme = strings.ToLower(scheme)
	switch position {
	case "QB1":
		if strings.Contains(scheme, "shoot") || strings.Contains(scheme, "raid") || strings.Contains(scheme, "pistol") || strings.Contains(scheme, "spread") {
			return 1.5
		} else if strings.Contains(scheme, "wing") || strings.Contains(scheme, "bone") {
			return 0.5
		} else {
			return 1.0
		}
	case "RB1":
		if strings.Contains(scheme, "shoot") || strings.Contains(scheme, "raid") || strings.Contains(scheme, "double") {
			return 0.8
		} else {
			return 1.0
		}
	case "RB2":
		if strings.Contains(scheme, "wing-t") || strings.Contains(scheme, "bone") {
			return 1.0
		} else if strings.Contains(scheme, "double") {
			return 0.8
		} else if strings.Contains(scheme, "spread") {
			return 0.7
		} else if strings.Contains(scheme, "raid") {
			return 0.6
		} else if strings.Contains(scheme, "shoot") {
			return 0.4
		} else {
			return 0.5
		}
	case "RB3":
		if strings.Contains(scheme, "wing-t") || strings.Contains(scheme, "wishbone") {
			return 0.1
		} else if strings.Contains(scheme, "flexbone") {
			return 0.5
		} else if strings.Contains(scheme, "double") {
			return 0.3
		} else {
			return 0.0
		}
	case "FB1":
		if strings.Contains(scheme, "wing") || strings.Contains(scheme, "bone") {
			return 1.0
		} else if strings.Contains(scheme, "power") {
			return 0.6
		} else if strings.Contains(scheme, "vert") || strings.Contains(scheme, "spread") {
			return 0.2
		} else if strings.Contains(scheme, "shoot") || strings.Contains(scheme, "raid") {
			return 0.0
		} else {
			return 0.4
		}
	case "TE1":
		if strings.Contains(scheme, "vert") || strings.Contains(scheme, "west") || strings.Contains(scheme, "pistol") || strings.Contains(scheme, "wishbone") {
			return 0.8
		} else if strings.Contains(scheme, "spread") || strings.Contains(scheme, "double") || strings.Contains(scheme, "flex") {
			return 0.6
		} else if strings.Contains(scheme, "shoot") || strings.Contains(scheme, "raid") {
			return 0.4
		} else {
			return 1.0
		}
	case "TE2":
		if strings.Contains(scheme, "power") || strings.Contains(scheme, "pistol") || strings.Contains(scheme, "double") || strings.Contains(scheme, "bone") {
			return 0.4
		} else if strings.Contains(scheme, "shoot") || strings.Contains(scheme, "west") {
			return 0.0
		} else {
			return 0.2
		}
	case "WR1":
		if strings.Contains(scheme, "wing") || strings.Contains(scheme, "wishbone") {
			return 0.8
		} else if strings.Contains(scheme, "flexbone") {
			return 0.6
		} else {
			return 1.0
		}
	case "WR2":
		if strings.Contains(scheme, "power") || strings.Contains(scheme, "pistol") {
			return 0.8
		} else if strings.Contains(scheme, "wing") || strings.Contains(scheme, "bone") {
			return 0.4
		} else {
			return 1.0
		}
	case "WR3":
		if strings.Contains(scheme, "raid") {
			return 0.8
		} else if strings.Contains(scheme, "vert") || strings.Contains(scheme, "west") || strings.Contains(scheme, "spread") {
			return 0.6
		} else if strings.Contains(scheme, "power") || strings.Contains(scheme, "double") {
			return 0.2
		} else if strings.Contains(scheme, "raid") {
			return 1.0
		} else if strings.Contains(scheme, "bone") || strings.Contains(scheme, "wing-t") {
			return 0.0
		} else {
			return 0.4
		}
	case "WR4":
		if strings.Contains(scheme, "shoot") {
			return 0.6
		} else if strings.Contains(scheme, "vert") || strings.Contains(scheme, "west") || strings.Contains(scheme, "pistol") || strings.Contains(scheme, "spread") || strings.Contains(scheme, "double") {
			return 0.2
		} else if strings.Contains(scheme, "raid") {
			return 0.4
		} else {
			return 0.0
		}
	case "WR5":
		if strings.Contains(scheme, "shoot") || strings.Contains(scheme, "raid") {
			return 0.3
		} else {
			return 0.0
		}
	case "LT1":
		if strings.Contains(scheme, "shoot") || strings.Contains(scheme, "raid") || strings.Contains(scheme, "pistol") || strings.Contains(scheme, "spread") {
			return 0.9
		} else if strings.Contains(scheme, "wing") || strings.Contains(scheme, "bone") {
			return 1.1
		} else {
			return 1.0
		}
	case "LG1":
		if strings.Contains(scheme, "shoot") || strings.Contains(scheme, "raid") || strings.Contains(scheme, "pistol") || strings.Contains(scheme, "spread") {
			return 0.9
		} else if strings.Contains(scheme, "wing") || strings.Contains(scheme, "bone") {
			return 1.1
		} else {
			return 1.0
		}
	case "C1":
		if strings.Contains(scheme, "shoot") || strings.Contains(scheme, "raid") || strings.Contains(scheme, "pistol") || strings.Contains(scheme, "spread") {
			return 0.9
		} else if strings.Contains(scheme, "wing") || strings.Contains(scheme, "bone") {
			return 1.1
		} else {
			return 1.0
		}
	case "RG1":
		if strings.Contains(scheme, "shoot") || strings.Contains(scheme, "raid") || strings.Contains(scheme, "pistol") || strings.Contains(scheme, "spread") {
			return 0.9
		} else if strings.Contains(scheme, "wing") || strings.Contains(scheme, "bone") {
			return 1.1
		} else {
			return 1.0
		}
	case "RT1":
		if strings.Contains(scheme, "shoot") || strings.Contains(scheme, "raid") || strings.Contains(scheme, "pistol") || strings.Contains(scheme, "spread") {
			return 0.9
		} else if strings.Contains(scheme, "wing") || strings.Contains(scheme, "bone") {
			return 1.1
		} else {
			return 1.0
		}
	default:
		return 0.0
	}
}

func GetDefensePositionGradeWeight(position string, scheme string) float64 {
	scheme = strings.ToLower(scheme)
	switch position {
	case "LE1":
		return 1.0
	case "DT1":
		if strings.Contains(scheme, "2") {
			return 0.6
		} else {
			return 1.0
		}
	case "DT2":
		if strings.Contains(scheme, "2") || strings.Contains(scheme, "3") {
			return 0.0
		} else if strings.Contains(scheme, "multiple") {
			return 0.6
		} else {
			return 1.0
		}
	case "RE1":
		return 1.0
	case "LOLB1":
		if strings.Contains(scheme, "4") {
			return 0.0
		} else if strings.Contains(scheme, "old") || strings.Contains(scheme, "2") {
			return 0.8
		} else {
			return 0.6
		}
	case "MLB1":
		return 1.0
	case "MLB2":
		if strings.Contains(scheme, "2") {
			return 1.0
		} else if strings.Contains(scheme, "old") {
			return 0.6
		} else if strings.Contains(scheme, "multiple") {
			return 0.4
		} else {
			return 0.0
		}
	case "ROLB1":
		if strings.Contains(scheme, "speed") {
			return 0.8
		} else if strings.Contains(scheme, "4") {
			return 0.6
		} else {
			return 1.0
		}
	case "CB1":
		return 1.0
	case "CB2":
		return 1.0
	case "CB3":
		if strings.Contains(scheme, "3") {
			return 1.0
		} else if strings.Contains(scheme, "4") {
			return 0.6
		} else if strings.Contains(scheme, "old") {
			return 0.2
		} else {
			return 0.4
		}
	case "CB4":
		if strings.Contains(scheme, "old") || strings.Contains(scheme, "multiple") {
			return 1.0
		} else {
			return 0.2
		}
	case "FS1":
		return 1.0
	case "SS1":
		if strings.Contains(scheme, "old") {
			return 0.4
		} else {
			return 1.0
		}
	case "SS2":
		if strings.Contains(scheme, "4") {
			return 0.6
		} else if strings.Contains(scheme, "3") {
			return 0.2
		} else {
			return 0.0
		}
	default:
		return 0.0
	}
}

func TeamLetterGrade(value float64, mean float64, stdDev float64) string {
	// Assign a letter grade for each grade value by comparing the value to the following:
	if value > (mean + (2 * stdDev)) {
		// A+: 2.0+ std dev above the mean
		return "A+"
	} else if value > (mean + (1.75 * stdDev)) {
		// A: between 1.75-2.0 std dev above the mean
		return "A"
	} else if value > (mean + (1.5 * stdDev)) {
		// A-: between 1.5-1.75 std dev above the mean
		return "A-"
	} else if value > (mean + (1.25 * stdDev)) {
		// B+: between 1.25-1.5 std dev above the mean
		return "B+"
	} else if value > (mean + (1.0 * stdDev)) {
		// B: between 1.0-1.25 std dev above the mean
		return "B"
	} else if value > (mean + (0.75 * stdDev)) {
		// B-: between .75-1.0 std dev above the mean
		return "B-"
	} else if value > (mean + (0.5 * stdDev)) {
		// C+: between .5-.75 std dev above the mean
		return "C+"
	} else if value > (mean - (0.5 * stdDev)) {
		// C: between +/- .5 std dev from mean
		return "C"
	} else if value > (mean - (0.75 * stdDev)) {
		// C-: between .5-.75 std dev below the mean
		return "C-"
	} else if value > (mean - (1.0 * stdDev)) {
		// D+: between .75-1.0 std dev below the mean
		return "D+"
	} else if value > (mean - (1.5 * stdDev)) {
		// D: between 1.0-1.5 std dev below the mean
		return "D"
	} else if value > (mean - (2.0 * stdDev)) {
		// D-: between 1.5-2.0 std dev below the mean
		return "D-"
	} else {
		// F: 2.0+ std dev below the mean
		return "F"
	}
}

// This function should be called weekly, once 2.0 is released.
func AssignTeamGrades() {
	db := dbprovider.GetInstance().GetDB()

	// College
	collegeTeams := GetAllCollegeTeams()
	collegeDepthChartMap := GetDepthChartMap()
	collegeGameplanMap := GetCollegeGameplanMap()
	collegeTeamGrades := make(map[uint]*structs.TeamGrade)

	for _, t := range collegeTeams {
		depthChart := collegeDepthChartMap[t.ID]
		gameplan := collegeGameplanMap[t.ID]
		offenseGrade := OffenseGradeCFB(depthChart, gameplan)
		defenseGrade := DefenseGradeCFB(depthChart, gameplan)
		STGrade := STGradeCFB(depthChart)

		collegeTeamGrades[t.ID] = &structs.TeamGrade{
			OffenseGradeNumber:      offenseGrade,
			DefenseGradeNumber:      defenseGrade,
			SpecialTeamsGradeNumber: STGrade,
			OverallGradeNumber:      OverallGrade(offenseGrade, defenseGrade, STGrade),
			OffenseGradeLetter:      "",
			DefenseGradeLetter:      "",
			SpecialTeamsGradeLetter: "",
			OverallGradeLetter:      "",
		}
	}

	// Determine the mean and std dev for the data set contained in the collegeTeamGrades map
	offenseMean := 0.0
	defenseMean := 0.0
	stMean := 0.0
	overallMean := 0.0
	offenseVar := 0.0
	defenseVar := 0.0
	stVar := 0.0
	overallVar := 0.0
	dataLength := len(collegeTeamGrades)

	// Mean for all values
	for _, element := range collegeTeamGrades {
		offenseMean += element.OffenseGradeNumber
		defenseMean += element.DefenseGradeNumber
		stMean += element.SpecialTeamsGradeNumber
		overallMean += element.OverallGradeNumber
	}

	offenseMean = (offenseMean / float64(dataLength))
	defenseMean = (defenseMean / float64(dataLength))
	stMean = (stMean / float64(dataLength))
	overallMean = (overallMean / float64(dataLength))

	// Variance for all values
	for _, element := range collegeTeamGrades {
		offenseDiff := element.OffenseGradeNumber - offenseMean
		defenseDiff := element.DefenseGradeNumber - defenseMean
		stDiff := element.SpecialTeamsGradeNumber - stMean
		overallDiff := element.OffenseGradeNumber - overallMean

		offenseVar += (offenseDiff * offenseDiff)
		defenseVar += (defenseDiff * defenseDiff)
		stVar += (stDiff * stDiff)
		overallVar += (overallDiff * overallDiff)
	}

	offenseVar = (offenseVar / float64(dataLength-1))
	defenseVar = (defenseVar / float64(dataLength-1))
	stVar = (stVar / float64(dataLength-1))
	overallVar = (overallVar / float64(dataLength-1))

	// Std Dev for all values
	offenseStdDev := math.Sqrt(offenseVar)
	defenseStdDev := math.Sqrt(defenseVar)
	stStdDev := math.Sqrt(stVar)
	overallStdDev := math.Sqrt(overallVar)

	// Iterate back through the map and set the letter grades based on the number grades' relationship to the mean and std dev of the entire data set for that value
	for i := uint(1); i < uint(len(collegeTeams)); i++ {
		_, ok := collegeTeamGrades[i]
		if ok {
			collegeTeamGrades[i].SetOffenseGradeLetter(TeamLetterGrade(collegeTeamGrades[i].OffenseGradeNumber, offenseMean, offenseStdDev))
			collegeTeamGrades[i].SetDefenseGradeLetter(TeamLetterGrade(collegeTeamGrades[i].DefenseGradeNumber, defenseMean, defenseStdDev))
			collegeTeamGrades[i].SetSpecialTeamsGradeLetter(TeamLetterGrade(collegeTeamGrades[i].SpecialTeamsGradeNumber, stMean, stStdDev))
			collegeTeamGrades[i].SetOverallGradeLetter(TeamLetterGrade(collegeTeamGrades[i].OverallGradeNumber, overallMean, overallStdDev))
		}
	}

	// Assign those letter grades to that team's grade properties
	for _, collegeTeam := range collegeTeams {
		if collegeTeam.ID > 194 {
			break
		}
		collegeTeam.AssignTeamGrades(collegeTeamGrades[collegeTeam.ID].OverallGradeLetter, collegeTeamGrades[collegeTeam.ID].OffenseGradeLetter,
			collegeTeamGrades[collegeTeam.ID].DefenseGradeLetter, collegeTeamGrades[collegeTeam.ID].SpecialTeamsGradeLetter)

		repository.SaveCFBTeam(collegeTeam, db)
	}

	// NFL
	nflTeams := GetAllNFLTeams()
	nflDepthChartList := GetAllNFLDepthcharts()
	nflGameplanList := GetAllNFLGameplans()
	nflDepthChartMap := make(map[uint]structs.NFLDepthChart)
	nflGameplanMap := make(map[uint]structs.NFLGameplan)

	for _, t := range nflDepthChartList {
		nflDepthChartMap[t.ID] = t
	}

	for _, t := range nflGameplanList {
		nflGameplanMap[t.ID] = t
	}

	nflTeamGrades := make(map[uint]*structs.TeamGrade)

	// CHANGE ALL REFERENCES TO COLLEGE TO NFL FROM HERE ON
	for _, t := range nflTeams {
		depthChart := nflDepthChartMap[t.ID]
		gameplan := nflGameplanMap[t.ID]
		offenseGrade := OffenseGradeNFL(depthChart, gameplan)
		defenseGrade := DefenseGradeNFL(depthChart, gameplan)
		STGrade := STGradeNFL(depthChart)

		nflTeamGrades[t.ID] = &structs.TeamGrade{
			OffenseGradeNumber:      offenseGrade,
			DefenseGradeNumber:      defenseGrade,
			SpecialTeamsGradeNumber: STGrade,
			OverallGradeNumber:      OverallGrade(offenseGrade, defenseGrade, STGrade),
			OffenseGradeLetter:      "",
			DefenseGradeLetter:      "",
			SpecialTeamsGradeLetter: "",
			OverallGradeLetter:      "",
		}
	}

	// Determine the mean and std dev for the data set contained in the nflTeamGrades map
	offenseMean = 0.0
	defenseMean = 0.0
	stMean = 0.0
	overallMean = 0.0
	offenseVar = 0.0
	defenseVar = 0.0
	stVar = 0.0
	overallVar = 0.0
	dataLength = len(nflTeamGrades)

	// Mean for all values
	for _, element := range nflTeamGrades {
		offenseMean += element.OffenseGradeNumber
		defenseMean += element.DefenseGradeNumber
		stMean += element.SpecialTeamsGradeNumber
		overallMean += element.OverallGradeNumber
	}

	offenseMean = (offenseMean / float64(dataLength))
	defenseMean = (defenseMean / float64(dataLength))
	stMean = (stMean / float64(dataLength))
	overallMean = (overallMean / float64(dataLength))

	// Variance for all values
	for _, element := range nflTeamGrades {
		offenseDiff := element.OffenseGradeNumber - offenseMean
		defenseDiff := element.DefenseGradeNumber - defenseMean
		stDiff := element.SpecialTeamsGradeNumber - stMean
		overallDiff := element.OffenseGradeNumber - overallMean

		offenseVar += (offenseDiff * offenseDiff)
		defenseVar += (defenseDiff * defenseDiff)
		stVar += (stDiff * stDiff)
		overallVar += (overallDiff * overallDiff)
	}

	offenseVar = (offenseVar / float64(dataLength-1))
	defenseVar = (defenseVar / float64(dataLength-1))
	stVar = (stVar / float64(dataLength-1))
	overallVar = (overallVar / float64(dataLength-1))

	// Std Dev for all values
	offenseStdDev = math.Sqrt(offenseVar)
	defenseStdDev = math.Sqrt(defenseVar)
	stStdDev = math.Sqrt(stVar)
	overallStdDev = math.Sqrt(overallVar)

	// Iterate back through the map and set the letter grades based on the number grades' relationship to the mean and std dev of the entire data set for that value
	for i := uint(1); i < uint(len(nflTeams)+1); i++ {
		_, ok := nflTeamGrades[i]
		if ok {
			nflTeamGrades[i].SetOffenseGradeLetter(TeamLetterGrade(nflTeamGrades[i].OffenseGradeNumber, offenseMean, offenseStdDev))
			nflTeamGrades[i].SetDefenseGradeLetter(TeamLetterGrade(nflTeamGrades[i].DefenseGradeNumber, defenseMean, defenseStdDev))
			nflTeamGrades[i].SetSpecialTeamsGradeLetter(TeamLetterGrade(nflTeamGrades[i].SpecialTeamsGradeNumber, stMean, stStdDev))
			nflTeamGrades[i].SetOverallGradeLetter(TeamLetterGrade(nflTeamGrades[i].OverallGradeNumber, overallMean, overallStdDev))
		}
	}

	// Assign those letter grades to that team's grade properties
	for _, nflTeam := range nflTeams {
		nflTeam.AssignTeamGrades(nflTeamGrades[nflTeam.ID].OverallGradeLetter, nflTeamGrades[nflTeam.ID].OffenseGradeLetter,
			nflTeamGrades[nflTeam.ID].DefenseGradeLetter, nflTeamGrades[nflTeam.ID].SpecialTeamsGradeLetter)

		repository.SaveNFLTeam(nflTeam, db)
	}
}

func UpdateCFBJersey(jersey structs.JerseyDTO) {
	db := dbprovider.GetInstance().GetDB()

	team := GetTeamByTeamID(jersey.TeamID)

	team.AdjustJerseyType(jersey.JerseyType)

	repository.SaveCollegeTeamRecord(team, db)
}

func UpdateNFLJersey(jersey structs.JerseyDTO) {
	db := dbprovider.GetInstance().GetDB()

	team := GetNFLTeamByTeamID(jersey.TeamID)

	team.AdjustJerseyType(jersey.JerseyType)

	repository.SaveNFLTeamRecord(team, db)
}
