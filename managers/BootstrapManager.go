package managers

import (
	"log"
	"strconv"
	"sync"

	"github.com/CalebRose/SimFBA/models"
	"github.com/CalebRose/SimFBA/structs"
)

type BootstrapData struct {
	CollegeTeam          structs.CollegeTeam
	AllCollegeTeams      []structs.CollegeTeam
	CollegeRosterMap     map[uint][]structs.CollegePlayer
	PortalPlayers        []structs.CollegePlayer
	CollegeInjuryReport  []structs.CollegePlayer
	CollegeNotifications []structs.Notification
	CollegeGameplan      structs.CollegeGameplan
	CollegeDepthChart    structs.CollegeTeamDepthChart
	ProTeam              structs.NFLTeam
	AllProTeams          []structs.NFLTeam
	ProNotifications     []structs.Notification
	NFLGameplan          structs.NFLGameplan
	NFLDepthChart        structs.NFLDepthChart
}

type BootstrapDataTwo struct {
	CollegeNews      []structs.NewsLog
	AllCollegeGames  []structs.CollegeGame
	TeamProfileMap   map[string]*structs.RecruitingTeamProfile
	CollegeStandings []structs.CollegeStandings
	ProStandings     []structs.NFLStandings
	AllProGames      []structs.NFLGame
	CapsheetMap      map[uint]structs.NFLCapsheet
	ProRosterMap     map[uint][]structs.NFLPlayer
	ProInjuryReport  []structs.NFLPlayer
}

type BootstrapDataThree struct {
	Recruits             []structs.Croot
	CollegeDepthChartMap map[uint]structs.CollegeTeamDepthChart
	FreeAgency           models.FreeAgencyResponse
	ProNews              []structs.NewsLog
	NFLDepthChartMap     map[uint]structs.NFLDepthChart
}

func GetFirstBootstrapData(collegeID, proID string) BootstrapData {
	var wg sync.WaitGroup
	var mu sync.Mutex

	// College Data
	var (
		collegeTeam           structs.CollegeTeam
		allCollegeTeams       []structs.CollegeTeam
		collegePlayerMap      map[uint][]structs.CollegePlayer
		portalPlayers         []structs.CollegePlayer
		injuredCollegePlayers []structs.CollegePlayer
		collegeNotifications  []structs.Notification
		collegeGameplan       structs.CollegeGameplan
		collegeDepthChart     structs.CollegeTeamDepthChart
	)

	// Professional Data
	var (
		proTeam          structs.NFLTeam
		allProTeams      []structs.NFLTeam
		proNotifications []structs.Notification
		proGameplan      structs.NFLGameplan
		proDepthChart    structs.NFLDepthChart
	)

	// Start concurrent queries
	wg.Add(2)
	go func() {
		defer wg.Done()
		allCollegeTeams = GetAllCollegeTeams()
	}()
	go func() {
		defer wg.Done()
		allProTeams = GetAllNFLTeams()
	}()

	wg.Wait()

	if len(collegeID) > 0 {
		wg.Add(4)
		go func() {
			defer wg.Done()
			collegeTeam = GetTeamByTeamID(collegeID)
		}()

		go func() {
			defer wg.Done()
			collegePlayers := GetAllCollegePlayers()
			mu.Lock()
			collegePlayerMap = MakeCollegePlayerMapByTeamID(collegePlayers, true)
			injuredCollegePlayers = MakeCollegeInjuryList(collegePlayers)
			portalPlayers = MakeCollegePortalList(collegePlayers)
			mu.Unlock()
		}()
		go func() {
			defer wg.Done()
			collegeNotifications = GetNotificationByTeamIDAndLeague("CFB", collegeID)
		}()
		go func() {
			defer wg.Done()
			collegeGameplan = GetGameplanByTeamID(collegeID)
		}()
		go func() {
			defer wg.Done()
			collegeDepthChart = GetDepthchartByTeamID(collegeID)

		}()

		wg.Wait()
	}
	if len(proID) > 0 {
		wg.Add(3)
		go func() {
			defer wg.Done()
			proTeam = GetNFLTeamByTeamID(proID)
		}()

		go func() {
			defer wg.Done()
			proNotifications = GetNotificationByTeamIDAndLeague("NFL", proID)
		}()
		go func() {
			defer wg.Done()
			proGameplan = GetNFLGameplanByTeamID(proID)
		}()
		wg.Wait()
	}
	return BootstrapData{
		CollegeTeam:          collegeTeam,
		AllCollegeTeams:      allCollegeTeams,
		CollegeRosterMap:     collegePlayerMap,
		CollegeInjuryReport:  injuredCollegePlayers,
		CollegeNotifications: collegeNotifications,
		CollegeGameplan:      collegeGameplan,
		CollegeDepthChart:    collegeDepthChart,
		PortalPlayers:        portalPlayers,
		ProTeam:              proTeam,
		AllProTeams:          allProTeams,
		ProNotifications:     proNotifications,
		NFLGameplan:          proGameplan,
		NFLDepthChart:        proDepthChart,
	}
}

func GetSecondBootstrapData(collegeID, proID string) BootstrapDataTwo {
	log.Println("GetSecondBootstrapData called with collegeID:", collegeID, "and proID:", proID)

	var wg sync.WaitGroup
	var mu sync.Mutex
	// College Data
	var (
		collegeStandings []structs.CollegeStandings
		teamProfileMap   map[string]*structs.RecruitingTeamProfile
		collegeNews      []structs.NewsLog
		collegeGames     []structs.CollegeGame
	)

	// Professional Data
	var (
		proStandings      []structs.NFLStandings
		proRosterMap      map[uint][]structs.NFLPlayer
		capsheetMap       map[uint]structs.NFLCapsheet
		injuredProPlayers []structs.NFLPlayer
		proGames          []structs.NFLGame
	)
	ts := GetTimestamp()
	log.Println("Timestamp:", ts)

	// Start concurrent queries
	if len(collegeID) > 0 {
		wg.Add(4)
		go func() {
			defer wg.Done()
			log.Println("Fetching College News Logs...")
			collegeNews = GetAllNewsLogs()
			log.Println("Fetched College News Logs, count:", len(collegeNews))
		}()
		go func() {
			defer wg.Done()
			log.Println("Fetching College Games for seasonID:", ts.CollegeSeasonID)
			collegeGames = GetCollegeGamesBySeasonID(strconv.Itoa(int(ts.CollegeSeasonID)))
			log.Println("Fetched College Games, count:", len(collegeGames))
		}()
		go func() {
			defer wg.Done()
			log.Println("Fetching Team Profile Map...")
			teamProfileMap = GetTeamProfileMap()
			log.Println("Fetched Team Profile Map, count:", len(teamProfileMap))
		}()
		go func() {
			defer wg.Done()
			log.Println("Fetching College Standings for seasonID:", ts.CollegeSeasonID)
			collegeStandings = GetAllCollegeStandingsBySeasonID(strconv.Itoa(int(ts.CollegeSeasonID)))
			log.Println("Fetched College Standings, count:", len(collegeStandings))
		}()
		wg.Wait()
		log.Println("Completed all College data queries.")

	}
	if len(proID) > 0 {
		wg.Add(4)
		go func() {
			defer wg.Done()
			log.Println("Fetching NFL Standings for seasonID:", ts.NFLSeasonID)
			proStandings = GetAllNFLStandingsBySeasonID(strconv.Itoa(int(ts.NFLSeasonID)))
			log.Println("Fetched NFL Standings, count:", len(proStandings))
		}()
		go func() {
			defer wg.Done()
			log.Println("Fetching NFL Games for seasonID:", ts.NFLSeasonID)
			proGames = GetNFLGamesBySeasonID(strconv.Itoa(int(ts.NFLSeasonID)))
			log.Println("Fetched NFL Games, count:", len(proGames))
		}()
		go func() {
			defer wg.Done()
			log.Println("Fetching Capsheet Map...")
			capsheetMap = getCapsheetMap()
			log.Println("Fetched Capsheet Map, count:", len(capsheetMap))
		}()
		go func() {
			defer wg.Done()
			log.Println("Fetching NFL Players for roster mapping...")
			proPlayers := GetAllNFLPlayers()
			mu.Lock()
			proRosterMap = MakeNFLPlayerMapByTeamID(proPlayers, true)
			injuredProPlayers = MakeProInjuryList(proPlayers)
			mu.Unlock()
			log.Println("Fetched NFL Players, roster count:", len(proRosterMap), "injured count:", len(injuredProPlayers))
		}()
		wg.Wait()
		log.Println("Completed all Pro data queries.")
	}
	return BootstrapDataTwo{
		CollegeStandings: collegeStandings,
		CollegeNews:      collegeNews,
		AllCollegeGames:  collegeGames,
		TeamProfileMap:   teamProfileMap,
		ProStandings:     proStandings,
		ProRosterMap:     proRosterMap,
		CapsheetMap:      capsheetMap,
		ProInjuryReport:  injuredProPlayers,
		AllProGames:      proGames,
	}
}

func GetThirdBootstrapData(collegeID, proID string) BootstrapDataThree {
	var wg sync.WaitGroup
	var mu sync.Mutex
	// College Data
	var (
		recruits             []structs.Croot
		collegeDepthChartMap map[uint]structs.CollegeTeamDepthChart
	)

	// Professional Data
	var (
		freeAgency       models.FreeAgencyResponse
		proNews          []structs.NewsLog
		proDepthChartMap map[uint]structs.NFLDepthChart
	)

	freeAgencyCh := make(chan models.FreeAgencyResponse, 1)

	if len(collegeID) > 0 {
		wg.Add(2)
		go func() {
			defer wg.Done()
			recruits = GetAllRecruits()
		}()

		go func() {
			defer wg.Done()
			collegeDCs := GetAllCollegeDepthcharts()
			collegeDepthChartMap = MakeCollegeDepthChartMap(collegeDCs)
		}()

		wg.Wait()
	}
	if len(proID) > 0 {
		wg.Add(3)

		go func() {
			defer wg.Done()
			dcs := GetAllNFLDepthcharts()
			mu.Lock()
			proDepthChartMap = MakeNFLDepthChartMap(dcs)
			mu.Unlock()
		}()
		go func() {
			defer wg.Done()
			proNews = GetAllNFLNewsLogs()
		}()

		go func() {
			defer wg.Done()
			GetAllAvailableNFLPlayersViaChan(proID, freeAgencyCh)
		}()

		freeAgency = <-freeAgencyCh
		wg.Wait()
	}
	return BootstrapDataThree{
		CollegeDepthChartMap: collegeDepthChartMap,
		Recruits:             recruits,
		FreeAgency:           freeAgency,
		ProNews:              proNews,
		NFLDepthChartMap:     proDepthChartMap,
	}
}
