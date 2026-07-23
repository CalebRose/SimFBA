package managers

import (
	"fmt"
	"strconv"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
)

type LiveGameHubDTO struct {
	GameID                uint   `json:"GameID"`
	HomeTeam              string `json:"HomeTeam"`
	AwayTeam              string `json:"AwayTeam"`
	HomeTeamID            uint   `json:"HomeTeamID"`
	AwayTeamID            uint   `json:"AwayTeamID"`
	HomeTeamScore         uint   `json:"HomeTeamScore"`
	AwayTeamScore         uint   `json:"AwayTeamScore"`
	HomeTeamShootoutScore uint   `json:"HomeTeamShootoutScore"`
	AwayTeamShootoutScore uint   `json:"AwayTeamShootoutScore"`
	Period                uint8  `json:"Period"`
	TimeOnClock           uint16 `json:"TimeOnClock"`
	Zone                  uint8  `json:"Zone"`
	GameComplete          bool   `json:"GameComplete"`
	IsShootout            bool   `json:"IsShootout"`
}

type GameDetailsDTO struct {
	Feeds     []PbPDTO        `json:"Feeds"`
	HomeStats TeamBoxScoreDTO `json:"HomeStats"`
	AwayStats TeamBoxScoreDTO `json:"AwayStats"`
}

type PbPDTO struct {
	Period      uint8  `json:"Period"`
	TimeOnClock uint16 `json:"TimeOnClock"`
	PlayText    string `json:"PlayText"`
	Zone        uint8  `json:"Zone"`
	HomeScore   uint8  `json:"HomeScore"`
	AwayScore   uint8  `json:"AwayScore"`
	HomeSOScore uint8  `json:"HomeSOScore"`
	AwaySOScore uint8  `json:"AwaySOScore"`
}

type TeamBoxScoreDTO struct {
	Forwards  []PlayerBoxScoreDTO `json:"Forwards"`
	Defenders []PlayerBoxScoreDTO `json:"Defenders"`
}

type PlayerBoxScoreDTO struct {
	Name      string `json:"Name"`
	Goals     uint8  `json:"Goals"`
	Assists   uint8  `json:"Assists"`
	PlusMinus int8   `json:"PlusMinus"`
}

type BulkSpoofDataDTO struct {
	Plays   map[uint][]structs.PlayByPlayResponse `json:"Plays"`
	Rosters map[uint]GameRosterDTO                `json:"Rosters"`
}

type GameRosterDTO struct {
	HomeStats TeamBoxScoreDTO `json:"HomeStats"`
	AwayStats TeamBoxScoreDTO `json:"AwayStats"`
}

func GetBulkPlayByPlayData(isCollege bool, reqSeason string, reqWeek string, reqTimeslot string) BulkSpoofDataDTO {
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.CollegeSeasonID))
	weekID := strconv.Itoa(int(ts.CollegeWeekID))

	response := BulkSpoofDataDTO{
		Plays:   make(map[uint][]structs.PlayByPlayResponse),
		Rosters: make(map[uint]GameRosterDTO),
	}
	_, gameType := ts.GetCFBCurrentGameType()

	db := dbprovider.GetInstance().GetDB()

	if isCollege {
		isPreseason := "N"
		if ts.CFBSpringGames {
			isPreseason = "Y"
		}
		clauses := repository.GamesQuery{SeasonID: seasonID, WeekID: weekID, IsSpringGames: isPreseason}
		games := repository.FindCollegeGamesRecords(clauses)
		collegePlayers := repository.FindAllCollegePlayers(repository.PlayerQuery{})
		collegePlayerMap := MakeCollegePlayerMap(collegePlayers)
		collegePlayerMapByTeamID := MakeCollegePlayerMapByTeamID(collegePlayers, true)
		collegePlayerStats := repository.FindCollegePlayerGameStatsRecords(seasonID, weekID, gameType, "")
		collegePlayerStatsMap := MakeCollegePlayerStatsMapByTeamID(collegePlayerStats)

		for _, g := range games {
			// if reqTimeslot != "" && reqTimeslot != "undefined" && g.TimeSlot != reqTimeslot {
			// 	continue
			// }
			response.Plays[g.ID] = []structs.PlayByPlayResponse{}
		}

		var allPbPs []structs.CollegePlayByPlay
		gameIDs := make([]uint, 0, len(response.Plays))
		for id := range response.Plays {
			gameIDs = append(gameIDs, id)
		}
		db.Where("game_id IN ?", gameIDs).Find(&allPbPs)

		for _, g := range games {
			// if reqTimeslot != "" && reqTimeslot != "undefined" && g.TimeSlot != reqTimeslot {
			// 	continue
			// }
			response.Plays[g.ID] = []structs.PlayByPlayResponse{}

			// Build Roster for this game
			roster := GameRosterDTO{
				HomeStats: TeamBoxScoreDTO{Forwards: []PlayerBoxScoreDTO{}, Defenders: []PlayerBoxScoreDTO{}},
				AwayStats: TeamBoxScoreDTO{Forwards: []PlayerBoxScoreDTO{}, Defenders: []PlayerBoxScoreDTO{}},
			}

			htID := strconv.Itoa(g.HomeTeamID)
			atID := strconv.Itoa(g.AwayTeamID)
			homeStats := collegePlayerStatsMap[uint(g.HomeTeamID)]
			awayStats := collegePlayerStatsMap[uint(g.AwayTeamID)]
			homeRoster := collegePlayerMapByTeamID[uint(g.HomeTeamID)]
			awayRoster := collegePlayerMapByTeamID[uint(g.AwayTeamID)]
			homePlayers := MakeGameResultsPlayerListFromCFB(homeStats, homeRoster)
			awayPlayers := MakeGameResultsPlayerListFromCFB(awayStats, awayRoster)

			playerStats := []structs.CollegePlayerStats{}
			playerStats = append(playerStats, homeStats...)
			playerStats = append(playerStats, awayStats...)
			for _, s := range playerStats {
				if s.Snaps <= 0 {
					continue
				}
				pInfo := collegePlayerMap[uint(s.CollegePlayerID)]
				nameStr := fmt.Sprintf("%s. %s", string(pInfo.FirstName[0]), pInfo.LastName)
				isHome := s.TeamID == uint(g.HomeTeamID)

				ps := PlayerBoxScoreDTO{Name: nameStr, Goals: 0, Assists: 0, PlusMinus: 0}
				if pInfo.Position == "D" {
					if isHome {
						roster.HomeStats.Defenders = append(roster.HomeStats.Defenders, ps)
					} else {
						roster.AwayStats.Defenders = append(roster.AwayStats.Defenders, ps)
					}
				} else {
					if isHome {
						roster.HomeStats.Forwards = append(roster.HomeStats.Forwards, ps)
					} else {
						roster.AwayStats.Forwards = append(roster.AwayStats.Forwards, ps)
					}
				}

			}
			response.Rosters[g.ID] = roster
			participantMap := getGameParticipantMap(homePlayers, awayPlayers)

			gamePbps := []structs.CollegePlayByPlay{}
			for _, p := range allPbPs {
				if p.GameID != g.ID {
					continue
				}
				gamePbps = append(gamePbps, p)

			}
			response.Plays[uint(g.ID)] = append(response.Plays[uint(g.ID)], GenerateCFBPlayByPlayResponse(gamePbps, participantMap, true, htID, atID)...)
		}
	} else {
		// PRO LOGIC
		isPreseason := "N"
		if ts.NFLPreseason {
			isPreseason = "Y"
		}
		clauses := repository.GamesQuery{SeasonID: seasonID, WeekID: weekID, IsPreseasonGame: isPreseason}
		games := repository.FindNFLGamesRecords(clauses)
		nflPlayers := GetAllNFLPlayers()
		proPlayerMap := MakeNFLPlayerMap(nflPlayers)
		nflPlayerMapByTeamID := MakeNFLPlayerMapByTeamID(nflPlayers, true)
		nflPlayerStats := repository.FindProPlayerGameStatsRecords(seasonID, weekID, gameType, "")
		nflPlayerStatsMap := MakeNFLPlayerStatsMapByTeamID(nflPlayerStats)

		for _, g := range games {
			// if reqTimeslot != "" && reqTimeslot != "undefined" && g.TimeSlot != reqTimeslot {
			// 	continue
			// }
			response.Plays[g.ID] = []structs.PlayByPlayResponse{}
		}

		var allPbPs []structs.NFLPlayByPlay
		gameIDs := make([]uint, 0, len(response.Plays))
		for id := range response.Plays {
			gameIDs = append(gameIDs, id)
		}
		db.Where("game_id IN ?", gameIDs).Find(&allPbPs)

		for _, g := range games {
			// if reqTimeslot != "" && reqTimeslot != "undefined" && g.TimeSlot != reqTimeslot {
			// 	continue
			// }
			response.Plays[g.ID] = []structs.PlayByPlayResponse{}

			roster := GameRosterDTO{
				HomeStats: TeamBoxScoreDTO{Forwards: []PlayerBoxScoreDTO{}, Defenders: []PlayerBoxScoreDTO{}},
				AwayStats: TeamBoxScoreDTO{Forwards: []PlayerBoxScoreDTO{}, Defenders: []PlayerBoxScoreDTO{}},
			}
			htID := strconv.Itoa(g.HomeTeamID)
			atID := strconv.Itoa(g.AwayTeamID)
			homeStats := nflPlayerStatsMap[uint(g.HomeTeamID)]
			awayStats := nflPlayerStatsMap[uint(g.AwayTeamID)]
			homeRoster := nflPlayerMapByTeamID[uint(g.HomeTeamID)]
			awayRoster := nflPlayerMapByTeamID[uint(g.AwayTeamID)]
			homePlayers := MakeGameResultsPlayerListFromNFL(homeStats, homeRoster)
			awayPlayers := MakeGameResultsPlayerListFromNFL(awayStats, awayRoster)

			playerStats := []structs.NFLPlayerStats{}
			playerStats = append(playerStats, homeStats...)
			playerStats = append(playerStats, awayStats...)
			for _, s := range playerStats {
				if s.Snaps <= 0 {
					continue
				}
				pInfo := proPlayerMap[uint(s.NFLPlayerID)]
				nameStr := fmt.Sprintf("%s. %s", string(pInfo.FirstName[0]), pInfo.LastName)
				isHome := s.TeamID == uint(g.HomeTeamID)

				ps := PlayerBoxScoreDTO{Name: nameStr, Goals: 0, Assists: 0, PlusMinus: 0}
				if pInfo.Position == "D" {
					if isHome {
						roster.HomeStats.Defenders = append(roster.HomeStats.Defenders, ps)
					} else {
						roster.AwayStats.Defenders = append(roster.AwayStats.Defenders, ps)
					}
				} else {
					if isHome {
						roster.HomeStats.Forwards = append(roster.HomeStats.Forwards, ps)
					} else {
						roster.AwayStats.Forwards = append(roster.AwayStats.Forwards, ps)
					}
				}

			}
			participantMap := getGameParticipantMap(homePlayers, awayPlayers)
			response.Rosters[g.ID] = roster
			gamePbps := []structs.NFLPlayByPlay{}
			for _, p := range allPbPs {
				if p.GameID != g.ID {
					continue
				}
				gamePbps = append(gamePbps, p)

			}
			response.Plays[uint(g.ID)] = append(response.Plays[uint(g.ID)], GenerateNFLPlayByPlayResponse(gamePbps, participantMap, true, htID, atID)...)
		}
	}
	return response
}
