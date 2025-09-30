package managers

import (
	"log"
	"sort"
	"strconv"
	"sync"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/models"
	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
)

type BootstrapData struct {
	CollegeTeam             structs.CollegeTeam
	AllCollegeTeams         []structs.CollegeTeam
	CollegeRosterMap        map[uint][]structs.CollegePlayer
	HistoricCollegePlayers  []structs.HistoricCollegePlayer
	TopCFBPassers           []structs.CollegePlayer
	TopCFBRushers           []structs.CollegePlayer
	TopCFBReceivers         []structs.CollegePlayer
	PortalPlayers           []structs.CollegePlayer
	CollegeInjuryReport     []structs.CollegePlayer
	CollegeNotifications    []structs.Notification
	CollegeGameplan         structs.CollegeGameplan
	CollegeDepthChart       structs.CollegeTeamDepthChart
	ProTeam                 structs.NFLTeam
	AllProTeams             []structs.NFLTeam
	ProNotifications        []structs.Notification
	NFLGameplan             structs.NFLGameplan
	NFLDepthChart           structs.NFLDepthChart
	Recruits                []structs.Croot
	RecruitProfiles         []structs.RecruitPlayerProfile
	CollegeDepthChartMap    map[uint]structs.CollegeTeamDepthChart
	FreeAgentOffers         []structs.FreeAgencyOffer
	WaiverWireOffers        []structs.NFLWaiverOffer
	FreeAgents              []structs.NFLPlayer
	WaiverPlayers           []structs.NFLPlayer
	NFLDraftees             []models.NFLDraftee
	ProNews                 []structs.NewsLog
	NFLDepthChartMap        map[uint]structs.NFLDepthChart
	ContractMap             map[uint]structs.NFLContract
	ExtensionMap            map[uint]structs.NFLExtensionOffer
	FaceData                map[uint]structs.FaceDataResponse
	CollegeNews             []structs.NewsLog
	AllCollegeGames         []structs.CollegeGame
	TeamProfileMap          map[string]*structs.RecruitingTeamProfile
	CollegeStandings        []structs.CollegeStandings
	ProStandings            []structs.NFLStandings
	AllProGames             []structs.NFLGame
	CapsheetMap             map[uint]structs.NFLCapsheet
	ProRosterMap            map[uint][]structs.NFLPlayer
	RetiredPlayers          []structs.NFLRetiredPlayer
	PracticeSquadPlayers    []structs.NFLPlayer
	TopNFLPassers           []structs.NFLPlayer
	TopNFLRushers           []structs.NFLPlayer
	TopNFLReceivers         []structs.NFLPlayer
	ProInjuryReport         []structs.NFLPlayer
	NFLWarRoomMap           map[uint]models.NFLWarRoom      // BY TEAM
	DraftScoutingProfileMap map[uint]models.ScoutingProfile // By TEAM
	TransferPortalProfiles  []structs.TransferPortalProfile
	CollegePromises         []structs.CollegePromise
	PollSubmission          structs.CollegePollSubmission
	OfficialPolls           []structs.CollegePollOfficial
}

type BootstrapDataTwo struct {
	CollegeNews          []structs.NewsLog
	AllCollegeGames      []structs.CollegeGame
	TeamProfileMap       map[string]*structs.RecruitingTeamProfile
	CollegeStandings     []structs.CollegeStandings
	ProStandings         []structs.NFLStandings
	AllProGames          []structs.NFLGame
	CapsheetMap          map[uint]structs.NFLCapsheet
	ProRosterMap         map[uint][]structs.NFLPlayer
	RetiredPlayers       []structs.NFLRetiredPlayer
	PracticeSquadPlayers []structs.NFLPlayer
	TopNFLPassers        []structs.NFLPlayer
	TopNFLRushers        []structs.NFLPlayer
	TopNFLReceivers      []structs.NFLPlayer
	ProInjuryReport      []structs.NFLPlayer
}

type BootstrapDataThree struct {
	Recruits             []structs.Croot
	RecruitProfiles      []structs.RecruitPlayerProfile
	CollegeDepthChartMap map[uint]structs.CollegeTeamDepthChart
	FreeAgentOffers      []structs.FreeAgencyOffer
	WaiverWireOffers     []structs.NFLWaiverOffer
	FreeAgents           []structs.NFLPlayer
	WaiverPlayers        []structs.NFLPlayer
	NFLDraftees          []models.NFLDraftee
	ProNews              []structs.NewsLog
	NFLDepthChartMap     map[uint]structs.NFLDepthChart
	ContractMap          map[uint]structs.NFLContract
	ExtensionMap         map[uint]structs.NFLExtensionOffer
	FaceData             map[uint]structs.FaceDataResponse
}

func GetTeamsBootstrap() BootstrapData {
	var wg sync.WaitGroup

	var (
		allCollegeTeams []structs.CollegeTeam
		allProTeams     []structs.NFLTeam
	)

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

	return BootstrapData{
		AllCollegeTeams: allCollegeTeams,
		AllProTeams:     allProTeams,
	}
}

func GetFirstBootstrapData(collegeID, proID string) BootstrapData {
	var wg sync.WaitGroup
	var mu sync.Mutex

	// College Data
	var (
		collegeTeam            structs.CollegeTeam
		collegePlayers         []structs.CollegePlayer
		collegePlayerMap       map[uint][]structs.CollegePlayer
		historicCollegePlayers []structs.HistoricCollegePlayer
		portalPlayers          []structs.CollegePlayer
		injuredCollegePlayers  []structs.CollegePlayer
		collegeNotifications   []structs.Notification
		collegeGameplan        structs.CollegeGameplan
		collegeDepthChart      structs.CollegeTeamDepthChart
		topPassers             []structs.CollegePlayer
		topRushers             []structs.CollegePlayer
		topReceivers           []structs.CollegePlayer
	)

	// Professional Data
	var (
		proTeam          structs.NFLTeam
		proNotifications []structs.Notification
		proGameplan      structs.NFLGameplan
		proDepthChart    structs.NFLDepthChart
	)

	ts := GetTimestamp()

	_, gtStr := ts.GetCFBCurrentGameType()
	seasonID := strconv.Itoa(int(ts.CollegeSeasonID))

	// Start concurrent queries

	if len(collegeID) > 0 && collegeID != "0" {
		wg.Add(6)
		go func() {
			defer wg.Done()
			mu.Lock()
			collegeTeam = GetTeamByTeamID(collegeID)
			collegeTeam.UpdateLatestInstance()
			repository.SaveCFBTeam(collegeTeam, dbprovider.GetInstance().GetDB())
			mu.Unlock()
		}()
		go func() {
			defer wg.Done()
			collegePlayers = GetAllCollegePlayers()
			cfbStats := GetCollegePlayerSeasonStatsBySeason(seasonID, gtStr)

			mu.Lock()
			collegePlayerMap = MakeCollegePlayerMapByTeamID(collegePlayers, true)
			fullCollegePlayerMap := MakeCollegePlayerMap(collegePlayers)
			topPassers = getCFBOrderedListByStatType("PASSING", collegeTeam.ID, cfbStats, fullCollegePlayerMap)
			topRushers = getCFBOrderedListByStatType("RUSHING", collegeTeam.ID, cfbStats, fullCollegePlayerMap)
			topReceivers = getCFBOrderedListByStatType("RECEIVING", collegeTeam.ID, cfbStats, fullCollegePlayerMap)
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
		go func() {
			defer wg.Done()
			historicCollegePlayers = GetAllHistoricCollegePlayers()
		}()
	}
	if len(proID) > 0 && proID != "0" {
		wg.Add(4)
		go func() {
			defer wg.Done()
			mu.Lock()
			proTeam = GetNFLTeamByTeamID(proID)
			proTeam.UpdateLatestInstance()
			repository.SaveNFLTeam(proTeam, dbprovider.GetInstance().GetDB())
			mu.Unlock()
		}()

		go func() {
			defer wg.Done()
			proNotifications = GetNotificationByTeamIDAndLeague("NFL", proID)
		}()
		go func() {
			defer wg.Done()
			proGameplan = GetNFLGameplanByTeamID(proID)
		}()
		go func() {
			defer wg.Done()
			proDepthChart = GetNFLDepthchartByTeamID(proID)
		}()
	}

	wg.Wait()
	return BootstrapData{
		CollegeTeam:            collegeTeam,
		CollegeRosterMap:       collegePlayerMap,
		CollegeInjuryReport:    injuredCollegePlayers,
		CollegeNotifications:   collegeNotifications,
		CollegeGameplan:        collegeGameplan,
		CollegeDepthChart:      collegeDepthChart,
		PortalPlayers:          portalPlayers,
		ProTeam:                proTeam,
		ProNotifications:       proNotifications,
		NFLGameplan:            proGameplan,
		NFLDepthChart:          proDepthChart,
		TopCFBPassers:          topPassers,
		TopCFBRushers:          topRushers,
		TopCFBReceivers:        topReceivers,
		HistoricCollegePlayers: historicCollegePlayers,
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
		proStandings         []structs.NFLStandings
		proRosterMap         map[uint][]structs.NFLPlayer
		practiceSquadPlayers []structs.NFLPlayer
		capsheetMap          map[uint]structs.NFLCapsheet
		injuredProPlayers    []structs.NFLPlayer
		proGames             []structs.NFLGame
		topPassers           []structs.NFLPlayer
		topRushers           []structs.NFLPlayer
		topReceivers         []structs.NFLPlayer
		retiredPlayers       []structs.NFLRetiredPlayer
	)
	ts := GetTimestamp()
	log.Println("Timestamp:", ts)
	_, gtStr := ts.GetNFLCurrentGameType()
	seasonID := strconv.Itoa(int(ts.NFLSeasonID))
	// Start concurrent queries
	if len(collegeID) > 0 && collegeID != "0" {
		wg.Add(4)
		go func() {
			defer wg.Done()
			log.Println("Fetching College News Logs...")
			collegeNews = GetAllCFBNewsLogs()
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
		log.Println("Initiated all College data queries.")
	}
	if len(proID) > 0 && proID != "0" {
		nflTeamID := util.ConvertStringToInt(proID)
		wg.Add(5)
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
			log.Println("Fetching Capsheet Map...")
			retiredPlayers = GetAllRetiredPlayers()
			log.Println("Fetched Capsheet Map, count:", len(capsheetMap))
		}()
		go func() {
			defer wg.Done()
			log.Println("Fetching NFL Players for roster mapping...")
			proPlayers := GetAllNFLPlayers()
			nflStats := GetNFLPlayerSeasonStatsBySeason(seasonID, gtStr)

			mu.Lock()
			nflPlayerMap := MakeNFLPlayerMap(proPlayers)
			proRosterMap = MakeNFLPlayerMapByTeamID(proPlayers, true)
			injuredProPlayers = MakeProInjuryList(proPlayers)
			practiceSquadPlayers = MakePracticeSquadList(proPlayers)
			topPassers = getNFLOrderedListByStatType("PASSING", uint(nflTeamID), nflStats, nflPlayerMap)
			topRushers = getNFLOrderedListByStatType("RUSHING", uint(nflTeamID), nflStats, nflPlayerMap)
			topReceivers = getNFLOrderedListByStatType("RECEIVING", uint(nflTeamID), nflStats, nflPlayerMap)
			mu.Unlock()
			log.Println("Fetched NFL Players, roster count:", len(proRosterMap), "injured count:", len(injuredProPlayers))
		}()

		log.Println("Initiated all Pro data queries.")
	}
	wg.Wait()
	log.Println("Completed all football data queries.")
	return BootstrapDataTwo{
		CollegeStandings:     collegeStandings,
		CollegeNews:          collegeNews,
		AllCollegeGames:      collegeGames,
		TeamProfileMap:       teamProfileMap,
		ProStandings:         proStandings,
		ProRosterMap:         proRosterMap,
		PracticeSquadPlayers: practiceSquadPlayers,
		CapsheetMap:          capsheetMap,
		ProInjuryReport:      injuredProPlayers,
		AllProGames:          proGames,
		TopNFLPassers:        topPassers,
		TopNFLRushers:        topRushers,
		TopNFLReceivers:      topReceivers,
		RetiredPlayers:       retiredPlayers,
	}
}

func GetThirdBootstrapData(collegeID, proID string) BootstrapDataThree {
	var wg sync.WaitGroup
	var mu sync.Mutex

	// College Data
	var (
		recruits             []structs.Croot
		recruitProfiles      []structs.RecruitPlayerProfile
		collegeDepthChartMap map[uint]structs.CollegeTeamDepthChart
		faceDataMap          map[uint]structs.FaceDataResponse
	)

	// Professional Data
	var (
		proNews          []structs.NewsLog
		proDepthChartMap map[uint]structs.NFLDepthChart
		contractMap      map[uint]structs.NFLContract
		extensionMap     map[uint]structs.NFLExtensionOffer
		freeAgents       []structs.NFLPlayer
		waiverPlayers    []structs.NFLPlayer
		freeAgentoffers  []structs.FreeAgencyOffer
		waiverOffers     []structs.NFLWaiverOffer
		nflDraftees      []models.NFLDraftee
	)

	if len(collegeID) > 0 && collegeID != "0" {
		wg.Add(3)
		go func() {
			defer wg.Done()
			recruits = GetAllRecruits()
		}()

		go func() {
			defer wg.Done()
			collegeDCs := GetAllCollegeDepthcharts()
			collegeDepthChartMap = MakeCollegeDepthChartMap(collegeDCs)
		}()
		go func() {
			defer wg.Done()
			recruitProfiles = repository.FindRecruitPlayerProfileRecords(collegeID, "", false, false, true)
		}()
	}

	if len(proID) > 0 && proID != "0" {
		wg.Add(9)

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
			freeAgentoffers = repository.FindAllFreeAgentOffers(repository.FreeAgencyQuery{IsActive: true})
		}()

		go func() {
			defer wg.Done()
			waiverOffers = repository.FindAllWaiverOffers(repository.FreeAgencyQuery{IsActive: true})
		}()

		go func() {
			defer wg.Done()
			contractMap = GetContractMap()
		}()

		go func() {
			defer wg.Done()
			extensionMap = GetExtensionMap()
		}()
		go func() {
			defer wg.Done()
			freeAgents = GetAllFreeAgents()
		}()

		go func() {
			defer wg.Done()
			waiverPlayers = GetAllWaiverWirePlayers()
		}()

		go func() {
			defer wg.Done()
			nflDraftees = GetAllNFLDraftees()
		}()

	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		faceDataMap = GetAllFaces()
	}()

	wg.Wait()

	return BootstrapDataThree{
		CollegeDepthChartMap: collegeDepthChartMap,
		Recruits:             recruits,
		RecruitProfiles:      recruitProfiles,
		FreeAgentOffers:      freeAgentoffers,
		WaiverWireOffers:     waiverOffers,
		ProNews:              proNews,
		NFLDepthChartMap:     proDepthChartMap,
		ContractMap:          contractMap,
		ExtensionMap:         extensionMap,
		FreeAgents:           freeAgents,
		WaiverPlayers:        waiverPlayers,
		FaceData:             faceDataMap,
		NFLDraftees:          nflDraftees,
	}
}

func GetTeamRosterBootstrap(collegeID, nflID string) BootstrapData {
	var wg sync.WaitGroup
	var (
		contractMap     map[uint]structs.NFLContract
		extensionMap    map[uint]structs.NFLExtensionOffer
		collegePromises []structs.CollegePromise
	)

	if len(collegeID) > 0 && collegeID != "0" {
		wg.Add(1)
		go func() {
			promises := GetCollegePromisesByTeamID(collegeID)
			collegePromises = promises
		}()
	}

	if len(nflID) > 0 && nflID != "0" {
		wg.Add(2)
		go func() {
			defer wg.Done()
			contractMap = GetContractMap()
		}()

		go func() {
			defer wg.Done()
			extensionMap = GetExtensionMap()
		}()
	}
	return BootstrapData{
		ContractMap:     contractMap,
		ExtensionMap:    extensionMap,
		CollegePromises: collegePromises,
	}
}

func GetRecruitingBootstrap(collegeID string) BootstrapData {
	var wg sync.WaitGroup
	var (
		recruits        []structs.Croot
		recruitProfiles []structs.RecruitPlayerProfile
		teamProfileMap  map[string]*structs.RecruitingTeamProfile
	)

	if len(collegeID) > 0 && collegeID != "0" {
		wg.Add(3)
		go func() {
			defer wg.Done()
			recruits = GetAllRecruits()
		}()

		go func() {
			defer wg.Done()
			teamProfileMap = GetTeamProfileMap()
		}()
		go func() {
			defer wg.Done()
			recruitProfiles = repository.FindRecruitPlayerProfileRecords(collegeID, "", false, false, true)
		}()
	}

	wg.Wait()

	return BootstrapData{
		Recruits:        recruits,
		RecruitProfiles: recruitProfiles,
		TeamProfileMap:  teamProfileMap,
	}
}

func GetFreeAgencyBootstrap(proID string) BootstrapData {
	var wg sync.WaitGroup
	var (
		freeAgents      []structs.NFLPlayer
		waiverPlayers   []structs.NFLPlayer
		freeAgentoffers []structs.FreeAgencyOffer
		waiverOffers    []structs.NFLWaiverOffer
	)

	if len(proID) > 0 && proID != "0" {
		wg.Add(4)

		go func() {
			defer wg.Done()
			freeAgentoffers = repository.FindAllFreeAgentOffers(repository.FreeAgencyQuery{IsActive: true})
		}()

		go func() {
			defer wg.Done()
			waiverOffers = repository.FindAllWaiverOffers(repository.FreeAgencyQuery{IsActive: true})
		}()

		go func() {
			defer wg.Done()
			freeAgents = GetAllFreeAgents()
		}()

		go func() {
			defer wg.Done()
			waiverPlayers = GetAllWaiverWirePlayers()
		}()

	}

	wg.Wait()

	return BootstrapData{
		FreeAgentOffers:  freeAgentoffers,
		WaiverWireOffers: waiverOffers,
		FreeAgents:       freeAgents,
		WaiverPlayers:    waiverPlayers,
	}
}

func GetSchedulingBootstrap(username, collegeID, proID, seasonID string) BootstrapData {
	var wg sync.WaitGroup
	var (
		collegeStandings []structs.CollegeStandings
		collegeGames     []structs.CollegeGame
		proStandings     []structs.NFLStandings
		proGames         []structs.NFLGame
		pollSubmission   structs.CollegePollSubmission
		officialPolls    []structs.CollegePollOfficial
	)
	if len(collegeID) > 0 && collegeID != "0" {
		wg.Add(4)
		go func() {
			defer wg.Done()
			log.Println("Fetching College Games for seasonID:", seasonID)
			collegeGames = GetCollegeGamesBySeasonID(seasonID)
			log.Println("Fetched College Games, count:", len(collegeGames))
		}()
		go func() {
			defer wg.Done()
			log.Println("Fetching College Standings for seasonID:", seasonID)
			collegeStandings = GetAllCollegeStandingsBySeasonID(seasonID)
			log.Println("Fetched College Standings, count:", len(collegeStandings))
		}()
		go func() {
			defer wg.Done()
			officialPolls = GetOfficialPollBySeasonID(seasonID)
		}()
		go func() {
			defer wg.Done()
			pollSubmission = GetPollSubmissionByUsernameWeekAndSeason(username)
		}()
		log.Println("Initiated all College data queries.")
	}

	if len(proID) > 0 && proID != "0" {
		wg.Add(2)
		go func() {
			defer wg.Done()
			log.Println("Fetching NFL Standings for seasonID:", seasonID)
			proStandings = GetAllNFLStandingsBySeasonID(seasonID)
			log.Println("Fetched NFL Standings, count:", len(proStandings))
		}()
		go func() {
			defer wg.Done()
			log.Println("Fetching NFL Games for seasonID:", seasonID)
			proGames = GetNFLGamesBySeasonID(seasonID)
			log.Println("Fetched NFL Games, count:", len(proGames))
		}()

		log.Println("Initiated all Pro data queries.")
	}

	wg.Wait()
	return BootstrapData{
		CollegeStandings: collegeStandings,
		AllCollegeGames:  collegeGames,
		ProStandings:     proStandings,
		AllProGames:      proGames,
		PollSubmission:   pollSubmission,
		OfficialPolls:    officialPolls,
	}
}

func GetDraftBootstrap(proID string) BootstrapData {
	var wg sync.WaitGroup

	var (
		nflDraftees        []models.NFLDraftee
		warRoomMap         map[uint]models.NFLWarRoom      // BY TEAM
		scoutingProfileMap map[uint]models.ScoutingProfile // By TEAM
	)

	if len(proID) > 0 && proID != "0" {
		wg.Add(3)
		go func() {
			defer wg.Done()
			nflDraftees = GetAllNFLDraftees()
		}()

		go func() {
			defer wg.Done()
			nflWarRooms := GetNFLWarRooms()
			warRoomMap = MakeNFLWarRoomMap(nflWarRooms)
		}()

		go func() {
			defer wg.Done()
			scoutingProfiles := GetAllScoutingProfiles()
			scoutingProfileMap = MakeScoutingProfileMapByTeam(scoutingProfiles)

		}()

		log.Println("Initiated all Pro data queries.")
	}
	wg.Wait()

	return BootstrapData{
		NFLDraftees:             nflDraftees,
		NFLWarRoomMap:           warRoomMap,
		DraftScoutingProfileMap: scoutingProfileMap,
	}
}

func GetTransferPortalBootstrap(collegeID string) BootstrapData {
	// On assumption that initial bootstrap still returns an entire college player map including transfers
	var wg sync.WaitGroup
	var (
		teamProfileMap         map[string]*structs.RecruitingTeamProfile // Get Just in Case because this page also uses this data
		transferPortalProfiles []structs.TransferPortalProfile
		collegePromises        []structs.CollegePromise
	)

	if len(collegeID) > 0 && collegeID != "0" {
		wg.Add(3)
		go func() {
			defer wg.Done()
			transferPortalProfiles = GetTransferPortalProfilesByTeamID(collegeID)
		}()

		go func() {
			defer wg.Done()
			teamProfileMap = GetTeamProfileMap()
		}()

		go func() {
			promises := GetCollegePromisesByTeamID(collegeID)
			collegePromises = promises
		}()

	}

	wg.Wait()

	return BootstrapData{
		TransferPortalProfiles: transferPortalProfiles,
		TeamProfileMap:         teamProfileMap,
		CollegePromises:        collegePromises,
	}
}

func GetGameplanBootstrap(collegeID, proID string) BootstrapData {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var (
		collegeGameplan      structs.CollegeGameplan
		collegeDepthChart    structs.CollegeTeamDepthChart
		proGameplan          structs.NFLGameplan
		proDepthChart        structs.NFLDepthChart
		collegeDepthChartMap map[uint]structs.CollegeTeamDepthChart
		proDepthChartMap     map[uint]structs.NFLDepthChart
	)

	if len(collegeID) > 0 && collegeID != "0" {
		wg.Add(3)
		go func() {
			defer wg.Done()
			collegeGameplan = GetGameplanByTeamID(collegeID)
		}()
		go func() {
			defer wg.Done()
			collegeDepthChart = GetDepthchartByTeamID(collegeID)
		}()

		go func() {
			defer wg.Done()
			collegeDCs := GetAllCollegeDepthcharts()
			collegeDepthChartMap = MakeCollegeDepthChartMap(collegeDCs)
		}()
	}

	if len(proID) > 0 && proID != "0" {
		wg.Add(3)
		go func() {
			defer wg.Done()
			proGameplan = GetNFLGameplanByTeamID(proID)
		}()
		go func() {
			defer wg.Done()
			proDepthChart = GetNFLDepthchartByTeamID(proID)
		}()
		go func() {
			defer wg.Done()
			dcs := GetAllNFLDepthcharts()
			mu.Lock()
			proDepthChartMap = MakeNFLDepthChartMap(dcs)
			mu.Unlock()
		}()

	}

	wg.Wait()
	return BootstrapData{
		CollegeGameplan:      collegeGameplan,
		CollegeDepthChart:    collegeDepthChart,
		CollegeDepthChartMap: collegeDepthChartMap,
		NFLGameplan:          proGameplan,
		NFLDepthChart:        proDepthChart,
		NFLDepthChartMap:     proDepthChartMap,
	}
}

func GetNewsBootstrap(collegeID, proID string) BootstrapData {
	var wg sync.WaitGroup

	var (
		collegeNews []structs.NewsLog
		proNews     []structs.NewsLog
	)

	if len(collegeID) > 0 && collegeID != "0" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Println("Fetching College News Logs...")
			collegeNews = GetAllCFBNewsLogs()
			log.Println("Fetched College News Logs, count:", len(collegeNews))
		}()
		log.Println("Initiated all College data queries.")
	}

	if len(proID) > 0 && proID != "0" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			proNews = GetAllNFLNewsLogs()
		}()

	}

	wg.Wait()

	return BootstrapData{
		CollegeNews: collegeNews,
		ProNews:     proNews,
	}
}

func getCFBOrderedListByStatType(statType string, teamID uint, CollegeStats []structs.CollegePlayerSeasonStats, collegePlayerMap map[uint]structs.CollegePlayer) []structs.CollegePlayer {
	orderedStats := CollegeStats
	resultList := []structs.CollegePlayer{}
	switch statType {
	case "PASSING":
		sort.Slice(orderedStats[:], func(i, j int) bool {
			return orderedStats[i].PassingTDs > orderedStats[j].PassingTDs
		})
	case "RUSHING":
		sort.Slice(orderedStats[:], func(i, j int) bool {
			return orderedStats[i].RushingYards > orderedStats[j].RushingYards
		})
	case "RECEIVING":
		sort.Slice(orderedStats[:], func(i, j int) bool {
			return orderedStats[i].ReceivingYards > orderedStats[j].ReceivingYards
		})
	}

	teamLeaderInTopStats := false
	for idx, stat := range orderedStats {
		if idx > 4 {
			break
		}
		player := collegePlayerMap[stat.CollegePlayerID]
		if stat.TeamID == teamID {
			teamLeaderInTopStats = true
		}
		player.AddSeasonStats(stat)
		resultList = append(resultList, player)
	}

	if !teamLeaderInTopStats {
		for _, stat := range orderedStats {
			if stat.TeamID == teamID {
				player := collegePlayerMap[stat.CollegePlayerID]
				player.AddSeasonStats(stat)
				resultList = append(resultList, player)
				break
			}
		}
	}
	return resultList
}

func getNFLOrderedListByStatType(statType string, teamID uint, CollegeStats []structs.NFLPlayerSeasonStats, proPlayerMap map[uint]structs.NFLPlayer) []structs.NFLPlayer {
	orderedStats := CollegeStats
	resultList := []structs.NFLPlayer{}
	switch statType {
	case "PASSING":
		sort.Slice(orderedStats[:], func(i, j int) bool {
			return orderedStats[i].PassingTDs > orderedStats[j].PassingTDs
		})
	case "RUSHING":
		sort.Slice(orderedStats[:], func(i, j int) bool {
			return orderedStats[i].RushingYards > orderedStats[j].RushingYards
		})
	case "RECEIVING":
		sort.Slice(orderedStats[:], func(i, j int) bool {
			return orderedStats[i].ReceivingYards > orderedStats[j].ReceivingYards
		})
	}

	teamLeaderInTopStats := false
	for idx, stat := range orderedStats {
		if idx > 4 {
			break
		}
		player := proPlayerMap[stat.NFLPlayerID]
		if stat.TeamID == teamID {
			teamLeaderInTopStats = true
		}
		player.AddSeasonStats(stat)
		resultList = append(resultList, player)
	}

	if !teamLeaderInTopStats {
		for _, stat := range orderedStats {
			if stat.TeamID == teamID {
				player := proPlayerMap[stat.NFLPlayerID]
				player.AddSeasonStats(stat)
				resultList = append(resultList, player)
				break
			}
		}
	}
	return resultList
}

type CollegeTeamProfileData struct {
	CareerStats      []structs.CollegePlayerSeasonStats
	CollegeStandings []structs.CollegeStandings
	Rivalries        []structs.FlexComparisonModel
	PlayerMap        map[uint]structs.CollegePlayer
	CollegeGames     []structs.CollegeGame
}

func GetCollegeTeamProfilePageData() map[uint]CollegeTeamProfileData {
	var wg sync.WaitGroup
	var mu sync.Mutex

	// College Data
	var (
		standingsMap     map[uint][]structs.CollegeStandings
		statsMap         map[uint][]structs.CollegePlayerSeasonStats
		collegePlayerMap map[uint][]structs.CollegePlayer
		teams            []structs.CollegeTeam
		teamMap          map[uint]structs.CollegeTeam
		rivalryMap       map[uint][]structs.CollegeRival
		gameMap          map[uint][]structs.CollegeGame
	)

	gamesByPair := make(map[uint]map[uint][]structs.CollegeGame)

	ts := GetTimestamp()
	// Get Career Stats

	wg.Add(6)
	go func() {
		defer wg.Done()
		standings := repository.FindAllCollegeStandingsRecords(repository.StandingsQuery{
			SeasonID: "",
		})
		mu.Lock()
		standingsMap = MakeHistoricCollegeStandingsMapByTeamID(standings)
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		seasonStats := GetCollegePlayerSeasonStatsByTeamID("")
		mu.Lock()
		statsMap = MakeHistoricCollegeSeasonStatsMapByTeamID(seasonStats)
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		collegePlayers := GetAllCollegePlayers()
		historicPlayers := GetAllHistoricCollegePlayers()
		mu.Lock()
		for _, player := range historicPlayers {
			collegePlayerResponse := structs.CollegePlayer{
				Model:      player.Model,
				BasePlayer: player.BasePlayer,
				TeamID:     player.TeamID,
				TeamAbbr:   player.TeamAbbr,
				City:       player.City,
				State:      player.State,
				Year:       player.Year,
				IsRedshirt: player.IsRedshirt,
			}
			collegePlayers = append(collegePlayers, collegePlayerResponse)
		}
		collegePlayerMap = MakeCollegePlayerMapByTeamID(collegePlayers, false)
		mu.Unlock()
	}()
	go func() {
		defer wg.Done()
		teams = GetAllCollegeTeams()
		mu.Lock()
		teamMap = GetCollegeTeamMap()
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		rivals := GetAllRivalries()
		mu.Lock()
		rivalryMap = MakeHistoricRivalriesMapByTeamID(rivals)
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		games := GetAllCollegeGames()
		mu.Lock()
		for _, g := range games {
			home, away := uint(g.HomeTeamID), uint(g.AwayTeamID)
			if gamesByPair[home] == nil {
				gamesByPair[home] = make(map[uint][]structs.CollegeGame)
			}
			gamesByPair[home][away] = append(gamesByPair[home][away], g)

			if gamesByPair[away] == nil {
				gamesByPair[away] = make(map[uint][]structs.CollegeGame)
			}
			gamesByPair[away][home] = append(gamesByPair[away][home], g)
		}
		gameMap = make(map[uint][]structs.CollegeGame)
		for _, g := range games {
			home, away := uint(g.HomeTeamID), uint(g.AwayTeamID)
			if len(gameMap[home]) > 0 {
				gameMap[home] = append(gameMap[home], g)
			} else {
				gameMap[home] = []structs.CollegeGame{g}
			}
			if len(gameMap[away]) > 0 {
				gameMap[away] = append(gameMap[away], g)
			} else {
				gameMap[away] = []structs.CollegeGame{g}
			}
		}
		mu.Unlock()
	}()

	wg.Wait()

	result := make(map[uint]CollegeTeamProfileData, len(teams))

	for _, team := range teams {
		collegePlayers := collegePlayerMap[team.ID]
		careerStatsList := make([]structs.CollegePlayerSeasonStats, 0, len(collegePlayers))
		playerMap := MakeCollegePlayerMap(collegePlayers)
		teamCareerStats := statsMap[team.ID]

		rosterStatsMap := make(map[uint][]structs.CollegePlayerSeasonStats)
		for _, s := range teamCareerStats {
			if len(rosterStatsMap[s.CollegePlayerID]) > 0 {
				rosterStatsMap[s.CollegePlayerID] = append(rosterStatsMap[s.CollegePlayerID], s)
			} else {
				rosterStatsMap[s.CollegePlayerID] = []structs.CollegePlayerSeasonStats{s}
			}
		}

		for _, player := range collegePlayers {
			stats := rosterStatsMap[player.ID]
			if len(stats) == 0 {
				continue
			}
			careerStats := structs.CollegePlayerSeasonStats{CollegePlayerID: player.ID, SeasonID: uint(ts.CollegeSeasonID)}
			careerStats.MapSeasonStats(stats)
			careerStatsList = append(careerStatsList, careerStats)
		}
		rivals := rivalryMap[team.ID]

		rivalryModels := []structs.FlexComparisonModel{}

		for _, rival := range rivals {
			var opp uint
			if rival.TeamOneID == team.ID {
				opp = rival.TeamTwoID
			} else {
				opp = rival.TeamOneID
			}
			team1, team2 := teamMap[rival.TeamOneID], teamMap[rival.TeamTwoID]
			team1ID := team1.ID
			team2ID := team2.ID
			t1Wins := 0
			t1Losses := 0
			t1Streak := 0
			t1CurrentStreak := 0
			t1LargestMarginSeason := 0
			t1LargestMarginDiff := 0
			t1LargestMarginScore := ""
			t2Wins := 0
			t2Losses := 0
			t2Streak := 0
			t2CurrentStreak := 0
			latestWin := ""
			t2LargestMarginSeason := 0
			t2LargestMarginDiff := 0
			t2LargestMarginScore := ""

			if t1CurrentStreak > 0 && t1CurrentStreak > t1Streak {
				t1Streak = t1CurrentStreak
			}
			if t2CurrentStreak > 0 && t2CurrentStreak > t2Streak {
				t2Streak = t2CurrentStreak
			}
			head2head := gamesByPair[team.ID][opp]

			for _, game := range head2head {
				if !game.GameComplete ||
					(game.Week == ts.CollegeWeek &&
						((game.TimeSlot == "Thursday Night" && !ts.ThursdayGames) ||
							(game.TimeSlot == "Friday Night" && !ts.FridayGames) ||
							(game.TimeSlot == "Saturday Morning" && !ts.SaturdayMorning) ||
							(game.TimeSlot == "Saturday Afternoon" && !ts.SaturdayNoon) ||
							(game.TimeSlot == "Saturday Evening" && !ts.SaturdayEvening) ||
							(game.TimeSlot == "Saturday Night" && !ts.SaturdayNight))) {
					continue
				}
				doComparison := (game.HomeTeamID == int(team1ID) && game.AwayTeamID == int(team2ID)) ||
					(game.HomeTeamID == int(team2ID) && game.AwayTeamID == int(team1ID))

				if !doComparison {
					continue
				}
				homeTeamTeamOne := game.HomeTeamID == int(team1ID)
				if homeTeamTeamOne {
					if game.HomeTeamWin {
						t1Wins += 1
						t1CurrentStreak += 1
						latestWin = game.HomeTeam
						diff := game.HomeTeamScore - game.AwayTeamScore
						if diff > t1LargestMarginDiff {
							t1LargestMarginDiff = diff
							t1LargestMarginSeason = game.SeasonID + 2020
							t1LargestMarginScore = "" + strconv.Itoa(game.HomeTeamScore) + "-" + strconv.Itoa(game.AwayTeamScore)
						}
					} else {
						t1Streak = t1CurrentStreak
						t1CurrentStreak = 0
						t1Losses += 1
					}
				} else {
					if game.HomeTeamWin {
						t2Wins += 1
						t2CurrentStreak += 1
						latestWin = game.HomeTeam
						diff := game.HomeTeamScore - game.AwayTeamScore
						if diff > t2LargestMarginDiff {
							t2LargestMarginDiff = diff
							t2LargestMarginSeason = game.SeasonID + 2020
							t2LargestMarginScore = "" + strconv.Itoa(game.HomeTeamScore) + "-" + strconv.Itoa(game.AwayTeamScore)
						}
					} else {
						t2Streak = t2CurrentStreak
						t2CurrentStreak = 0
						t2Losses += 1
					}
				}

				awayTeamTeamOne := game.AwayTeamID == int(team1ID)
				if awayTeamTeamOne {
					if game.AwayTeamWin {
						t1Wins += 1
						t1CurrentStreak += 1
						latestWin = game.AwayTeam
						diff := game.AwayTeamScore - game.HomeTeamScore
						if diff > t1LargestMarginDiff {
							t1LargestMarginDiff = diff
							t1LargestMarginSeason = game.SeasonID + 2020
							t1LargestMarginScore = "" + strconv.Itoa(game.AwayTeamScore) + "-" + strconv.Itoa(game.HomeTeamScore)
						}
					} else {
						t1Streak = t1CurrentStreak
						t1CurrentStreak = 0
						t1Losses += 1
					}
				} else {
					if game.AwayTeamWin {
						t2Wins += 1
						t2CurrentStreak += 1
						latestWin = game.AwayTeam
						diff := game.AwayTeamScore - game.HomeTeamScore
						if diff > t2LargestMarginDiff {
							t2LargestMarginDiff = diff
							t2LargestMarginSeason = game.SeasonID + 2020
							t2LargestMarginScore = "" + strconv.Itoa(game.AwayTeamScore) + "-" + strconv.Itoa(game.HomeTeamScore)
						}
					} else {
						t2Streak = t2CurrentStreak
						t2CurrentStreak = 0
						t2Losses += 1
					}
				}
			}

			currentStreak := max(t1CurrentStreak, t2CurrentStreak)

			rivalryModel := structs.FlexComparisonModel{
				TeamOneID:      uint(team1ID),
				TeamOne:        team1.TeamAbbr,
				TeamOneWins:    uint(t1Wins),
				TeamOneLosses:  uint(t1Losses),
				TeamOneStreak:  uint(t1Streak),
				TeamOneMSeason: t1LargestMarginSeason,
				TeamOneMScore:  t1LargestMarginScore,
				TeamTwoID:      uint(team2ID),
				TeamTwo:        team2.TeamAbbr,
				TeamTwoWins:    uint(t2Wins),
				TeamTwoLosses:  uint(t2Losses),
				TeamTwoStreak:  uint(t2Streak),
				TeamTwoMSeason: t2LargestMarginSeason,
				TeamTwoMScore:  t2LargestMarginScore,
				CurrentStreak:  uint(currentStreak),
				LatestWin:      latestWin,
			}

			rivalryModels = append(rivalryModels, rivalryModel)
		}
		data := CollegeTeamProfileData{
			CareerStats:      careerStatsList,
			CollegeStandings: standingsMap[team.ID],
			PlayerMap:        playerMap,
			Rivalries:        rivalryModels,
			CollegeGames:     gameMap[team.ID],
		}
		result[team.ID] = data
	}

	return result
}

func GetRivalriesByTeamID(teamID string) []structs.CollegeRival {
	db := dbprovider.GetInstance().GetDB()

	rivals := []structs.CollegeRival{}

	db.Where("team_one_id = ? OR team_two_id = ?", teamID, teamID).Find(&rivals)

	return rivals
}

func GetAllRivalries() []structs.CollegeRival {
	db := dbprovider.GetInstance().GetDB()

	rivals := []structs.CollegeRival{}

	db.Find(&rivals)

	return rivals
}
