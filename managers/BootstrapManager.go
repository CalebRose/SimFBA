package managers

import (
	"log"
	"sort"
	"strconv"
	"sync"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
)

type BootstrapData struct {
	CollegeTeam          structs.CollegeTeam
	AllCollegeTeams      []structs.CollegeTeam
	CollegeRosterMap     map[uint][]structs.CollegePlayer
	TopCFBPassers        []structs.CollegePlayer
	TopCFBRushers        []structs.CollegePlayer
	TopCFBReceivers      []structs.CollegePlayer
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
	FaceData             map[uint]structs.FaceDataResponse
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
	PracticeSquadPlayers []structs.NFLPlayer
	TopNFLPassers        []structs.NFLPlayer
	TopNFLRushers        []structs.NFLPlayer
	TopNFLReceivers      []structs.NFLPlayer
	ProInjuryReport      []structs.NFLPlayer
}

type BootstrapDataThree struct {
	Recruits             []structs.Croot
	CollegeDepthChartMap map[uint]structs.CollegeTeamDepthChart
	FreeAgentOffers      []structs.FreeAgencyOffer
	WaiverWireOffers     []structs.NFLWaiverOffer
	ProNews              []structs.NewsLog
	NFLDepthChartMap     map[uint]structs.NFLDepthChart
	ContractMap          map[uint]structs.NFLContract
	ExtensionMap         map[uint]structs.NFLExtensionOffer
}

func GetFirstBootstrapData(collegeID, proID string) BootstrapData {
	var wg sync.WaitGroup
	var mu sync.Mutex

	// College Data
	var (
		collegeTeam           structs.CollegeTeam
		collegePlayers        []structs.CollegePlayer
		allCollegeTeams       []structs.CollegeTeam
		collegePlayerMap      map[uint][]structs.CollegePlayer
		portalPlayers         []structs.CollegePlayer
		injuredCollegePlayers []structs.CollegePlayer
		collegeNotifications  []structs.Notification
		collegeGameplan       structs.CollegeGameplan
		collegeDepthChart     structs.CollegeTeamDepthChart
		topPassers            []structs.CollegePlayer
		topRushers            []structs.CollegePlayer
		topReceivers          []structs.CollegePlayer
		faceDataMap           map[uint]structs.FaceDataResponse
	)

	// Professional Data
	var (
		proTeam          structs.NFLTeam
		allProTeams      []structs.NFLTeam
		proNotifications []structs.Notification
		proGameplan      structs.NFLGameplan
		proDepthChart    structs.NFLDepthChart
	)

	ts := GetTimestamp()

	_, gtStr := ts.GetCFBCurrentGameType()
	seasonID := strconv.Itoa(int(ts.CollegeSeasonID))

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

	if len(collegeID) > 0 {
		wg.Add(5)
		go func() {
			defer wg.Done()
			collegeTeam = GetTeamByTeamID(collegeID)
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
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		faceDataMap = GetAllFaces()
	}()

	wg.Wait()
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
		TopCFBPassers:        topPassers,
		TopCFBRushers:        topRushers,
		TopCFBReceivers:      topReceivers,
		FaceData:             faceDataMap,
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
	)
	ts := GetTimestamp()
	log.Println("Timestamp:", ts)
	_, gtStr := ts.GetNFLCurrentGameType()
	seasonID := strconv.Itoa(int(ts.NFLSeasonID))
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
		log.Println("Initiated all College data queries.")

	}
	if len(proID) > 0 {
		nflTeamID := util.ConvertStringToInt(proID)
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
		wg.Wait()
		log.Println("Completed all football data queries.")
	}
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
		proNews          []structs.NewsLog
		proDepthChartMap map[uint]structs.NFLDepthChart
		contractMap      map[uint]structs.NFLContract
		extensionMap     map[uint]structs.NFLExtensionOffer
		freeAgentoffers  []structs.FreeAgencyOffer
		waiverOffers     []structs.NFLWaiverOffer
	)

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
	}

	if len(proID) > 0 {
		wg.Add(6)

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

		wg.Wait()
	}
	return BootstrapDataThree{
		CollegeDepthChartMap: collegeDepthChartMap,
		Recruits:             recruits,
		FreeAgentOffers:      freeAgentoffers,
		WaiverWireOffers:     waiverOffers,
		ProNews:              proNews,
		NFLDepthChartMap:     proDepthChartMap,
		ContractMap:          contractMap,
		ExtensionMap:         extensionMap,
	}
}

func getCFBOrderedListByStatType(statType string, teamID uint, CollegeStats []structs.CollegePlayerSeasonStats, collegePlayerMap map[uint]structs.CollegePlayer) []structs.CollegePlayer {
	orderedStats := CollegeStats
	resultList := []structs.CollegePlayer{}
	if statType == "PASSING" {
		sort.Slice(orderedStats[:], func(i, j int) bool {
			return orderedStats[i].PassingTDs > orderedStats[j].PassingTDs
		})
	} else if statType == "RUSHING" {
		sort.Slice(orderedStats[:], func(i, j int) bool {
			return orderedStats[i].RushingYards > orderedStats[j].RushingYards
		})
	} else if statType == "RECEIVING" {
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
	if statType == "PASSING" {
		sort.Slice(orderedStats[:], func(i, j int) bool {
			return orderedStats[i].PassingTDs > orderedStats[j].PassingTDs
		})
	} else if statType == "RUSHING" {
		sort.Slice(orderedStats[:], func(i, j int) bool {
			return orderedStats[i].RushingYards > orderedStats[j].RushingYards
		})
	} else if statType == "RECEIVING" {
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
}

func GetCollegeTeamProfilePageData(teamID string) CollegeTeamProfileData {
	ts := GetTimestamp()
	// Get Career Stats
	standings := repository.FindAllCollegeStandingsRecords(repository.StandingsQuery{
		TeamID:   teamID,
		SeasonID: "",
	})

	careerStatsList := []structs.CollegePlayerSeasonStats{}
	seasonStats := GetCollegePlayerSeasonStatsByTeamID(teamID)

	collegePlayers := GetAllCollegePlayersByTeamId(teamID)
	historicPlayers := GetAllHistoricCollegePlayers()

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

	collegePlayerMap := MakeCollegePlayerMap(collegePlayers)
	statsMap := make(map[uint][]structs.CollegePlayerSeasonStats)

	for _, stat := range seasonStats {
		if len(statsMap[stat.CollegePlayerID]) == 0 {
			statsMap[stat.CollegePlayerID] = []structs.CollegePlayerSeasonStats{stat}
		} else {
			statsMap[stat.CollegePlayerID] = append(statsMap[stat.CollegePlayerID], stat)
		}
	}

	for _, player := range collegePlayers {
		stats := statsMap[player.ID]
		if len(stats) == 0 {
			continue
		}
		careerStats := structs.CollegePlayerSeasonStats{CollegePlayerID: player.ID, SeasonID: uint(ts.CollegeSeasonID)}
		careerStats.MapSeasonStats(stats)
		careerStatsList = append(careerStatsList, careerStats)
	}
	collegeTeamMap := GetCollegeTeamMap()
	rivals := GetRivalriesByTeamID(teamID)
	games := GetCollegeGamesByTeamId(teamID)

	rivalryModels := []structs.FlexComparisonModel{}

	for _, rivalry := range rivals {
		t1 := strconv.Itoa(int(rivalry.TeamOneID))
		team1ID := 0
		team2ID := 0
		team1 := structs.CollegeTeam{}
		team2 := structs.CollegeTeam{}
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
		if t1 == teamID {
			team1 = collegeTeamMap[rivalry.TeamOneID]
			team2 = collegeTeamMap[rivalry.TeamTwoID]
			team1ID = int(team1.ID)
			team2ID = int(team2.ID)
		} else {
			team1 = collegeTeamMap[rivalry.TeamTwoID]
			team2 = collegeTeamMap[rivalry.TeamOneID]
			team1ID = int(team1.ID)
			team2ID = int(team2.ID)
		}

		if t1CurrentStreak > 0 && t1CurrentStreak > t1Streak {
			t1Streak = t1CurrentStreak
		}
		if t2CurrentStreak > 0 && t2CurrentStreak > t2Streak {
			t2Streak = t2CurrentStreak
		}

		for _, game := range games {
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

	return CollegeTeamProfileData{
		CareerStats:      careerStatsList,
		CollegeStandings: standings,
		PlayerMap:        collegePlayerMap,
		Rivalries:        rivalryModels,
	}
}

func GetRivalriesByTeamID(teamID string) []structs.CollegeRival {
	db := dbprovider.GetInstance().GetDB()

	rivals := []structs.CollegeRival{}

	db.Where("team_one_id = ? OR team_two_id = ?", teamID, teamID).Find(&rivals)

	return rivals
}
