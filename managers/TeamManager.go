package managers

import (
	"log"
	"strconv"
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

// Returns the CFB team's numerical value for their entire offense
func OffenseGradeCFB(depthChartPlayers structs.CollegeTeamDepthChart, gameplan structs.CollegeGameplan) float32 {
	// Get overall values for all relevant positions
	// If the player is a scheme fit, give them a bonus, if they are a bad fit, give them a malus
	// Depending on scheme, weight them
	// Sum them all up
	// Divide by 11.5 (offense weight normalization value)
	// return the resulting value
	return 0.0
}

// Returns the CFB team's numerical value for their entire offense
func DefenseGradeCFB(depthChartPlayers structs.CollegeTeamDepthChart, gameplan structs.CollegeGameplan) float32 {
	// Get overall values for all relevant positions
	// If the player is a scheme fit, give them a bonus, if they are a bad fit, give them a malus
	// Depending on scheme, weight them
	// Sum them all up
	// Divide by 11 (defense weight normalization value)
	// return the resulting value
	return 0.0
}

// Returns the CFB team's numerical value for their entire offense
func STGradeCFB(depthChartPlayers structs.CollegeTeamDepthChart) float32 {
	// Get overall values for all relevant positions
	// Weight them by position
	// Sum them all up
	// Divide by 5 (Special Teams weight normalization value)
	// return the resulting value
	return 0.0
}

// Returns the CFB team's numerical value for their entire offense
func OffenseGradeNFL(depthChartPlayers structs.NFLDepthChart, gameplan structs.NFLGameplan) float32 {
	// Get overall values for all relevant positions
	// If the player is a scheme fit, give them a bonus, if they are a bad fit, give them a malus
	// Depending on scheme, weight them
	// Sum them all up
	// Divide by 11.5 (offense weight normalization value)
	// return the resulting value
	return 0.0
}

// Returns the CFB team's numerical value for their entire offense
func DefenseGradeNFL(depthChartPlayers structs.NFLDepthChart, gameplan structs.NFLGameplan) float32 {
	// Get overall values for all relevant positions
	// If the player is a scheme fit, give them a bonus, if they are a bad fit, give them a malus
	// Depending on scheme, weight them
	// Sum them all up
	// Divide by 11 (defense weight normalization value)
	// return the resulting value
	return 0.0
}

// Returns the CFB team's numerical value for their entire offense
func STGradeNFL(depthChartPlayers structs.NFLDepthChart) float32 {
	// Get overall values for all relevant positions
	// Weight them by position
	// Sum them all up
	// Divide by 5 (Special Teams weight normalization value)
	// return the resulting value
	return 0.0
}

// League agnostic
func OverallGrade(offense float32, defense float32, specialTeams float32) float32 {
	var overallGrade float32 = offense * 0.45
	overallGrade = overallGrade + (defense * 0.45)
	overallGrade = overallGrade + (specialTeams * 0.1)
	return overallGrade
}

// This function should be called weekly, once 2.0 is released.
// Change to league agnostic?
func AssignTeamGrades() {
	db := dbprovider.GetInstance().GetDB()

	// College
	collegeTeams := GetAllCollegeTeams()
	collegeDepthChartMap := GetDepthChartMap()
	collegeGameplanMap := GetCollegeGameplanMap()
	collegeTeamGrades := make(map[uint]structs.TeamGrade)

	for _, t := range collegeTeams {
		if !t.IsActive {
			continue
		}
		depthChart := collegeDepthChartMap[t.ID]
		gameplan := collegeGameplanMap[t.ID]
		offenseGrade := OffenseGradeCFB(depthChart, gameplan)
		defenseGrade := DefenseGradeCFB(depthChart, gameplan)
		STGrade := STGradeCFB(depthChart)

		collegeTeamGrades[t.ID] = structs.TeamGrade{
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

	// Iterate back through the map and set the letter grades based on the number grades' relationship to the mean and std dev of the entire data set for that value

	// Assign those letter grades to that team's grade properties

	// NFL
	nflTeams := GetAllCollegeTeams()
	nflDepthChartMap := GetDepthChartMap()
	nflGameplanMap := GetCollegeGameplanMap()

	// Include reference to current gameplan
	// Move the grade calculation to 4 separate functions:
	// OffenseGrade
	// DefenseGrade
	// SpecialTeamsGrade
	// OverallGrade

	// Follow the following algorithm:
	// Acquire every team's numerical value for all 4 grades
	// Determine the mean and standard deviation of each grade across all teams
	// Assign a letter grade for each grade value by comparing the value to the following:
	// A+: 2.0+ std dev above the mean
	// A: between 1.75-2.0 std dev above the mean
	// A-: between 1.5-1.75 std dev above the mean
	// B+: between 1.25-1.5 std dev above the mean
	// B: between 1.0-1.25 std dev above the mean
	// B-: between .75-1.0 std dev above the mean
	// C+: between .5-.75 std dev above the mean
	// C: between +/- .5 std dev from mean
	// C-: between .5-.75 std dev below the mean
	// D+: between .75-1.0 std dev below the mean
	// D: between 1.0-1.5 std dev below the mean
	// D-: between 1.5-2.0 std dev below the mean
	// F: 2.0+ std dev below the mean

	for _, t := range collegeTeams {
		if !t.IsActive {
			continue
		}
		depthChart := depthChartMap[t.ID]
		players := depthChart.DepthChartPlayers
		overallScore := 0.0
		offensiveScore := 0.0
		defensiveScore := 0.0
		offenseCount := 0
		defenseCount := 0
		totalCount := 0
		for _, p := range players {
			if (p.Position == "QB" || p.Position == "RB" || p.Position == "FB" || p.Position == "TE" || p.Position == "C" || p.Position == "MLB" || p.Position == "FS" || p.Position == "SS" || p.Position == "K" || p.Position == "P") &&
				p.PositionLevel != "1" {
				continue
			}
			if (p.Position == "WR" || p.Position == "OT" || p.Position == "OG" || p.Position == "DT" || p.Position == "DE" || p.Position == "OLB" || p.Position == "CB") &&
				(p.PositionLevel != "1" && p.PositionLevel != "2") {
				continue
			}
			if p.Position == "KR" || p.Position == "PR" || p.Position == "FG" || p.Position == "STU" {
				continue
			}
			if p.Position == "QB" || p.Position == "RB" || p.Position == "FB" || p.Position == "TE" || p.Position == "C" || p.Position == "WR" || p.Position == "OT" || p.Position == "OG" {
				offensiveScore += float64(p.CollegePlayer.Overall)
				offenseCount += 1
			}
			if p.Position == "DT" || p.Position == "DE" || p.Position == "OLB" || p.Position == "CB" || p.Position == "MLB" || p.Position == "FS" || p.Position == "SS" {
				defensiveScore += float64(p.CollegePlayer.Overall)
				defenseCount += 1
			}
			totalCount += 1
			overallScore += float64(p.CollegePlayer.Overall)
		}

		ovrAvg := overallScore / float64(totalCount)
		offAvg := offensiveScore / float64(offenseCount)
		defAvg := defensiveScore / float64(defenseCount)

		ovrGrade := util.GetOverallGrade(int(ovrAvg), 4)
		offGrade := util.GetOverallGrade(int(offAvg), 4)
		defGrade := util.GetOverallGrade(int(defAvg), 4)

		t.AssignTeamGrades(ovrGrade, offGrade, defGrade)

		repository.SaveCFBTeam(t, db)
	}
}
