package managers

import (
	"strconv"

	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
)

// StreamGameQueueItem is the wire shape simsn-live consumes to build its
// own PendingGame queue without touching SimHockey's database directly.
type StreamGameQueueItem struct {
	GameID       uint   `json:"gameID"`
	HomeTeamID   uint   `json:"homeTeamID"`
	AwayTeamID   uint   `json:"awayTeamID"`
	HomeTeam     string `json:"homeTeam"`
	AwayTeam     string `json:"awayTeam"`
	IsUserGame   bool   `json:"isUserGame"`
	HomeTeamRank int    `json:"homeTeamRank"`
	AwayTeamRank int    `json:"awayTeamRank"`
	Arena        string `json:"arena"`
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
	TotalSeconds int    `json:"totalSeconds"`
	TotalPlays   int    `json:"totalPlays"`
}

// BuildStreamGameQueue returns the queued, unrevealed games for a league
// (chl|phl), ordered user-games-first, with PbP totals attached.
func BuildStreamGameQueue(league, weekID, seasonID, gameDay string, isPreseason bool) []StreamGameQueueItem {
	if league == "cfb" {
		return buildCFBStreamQueue(weekID, seasonID, gameDay, isPreseason)
	}
	return buildNFLStreamQueue(weekID, seasonID, gameDay, isPreseason)
}

func buildCFBStreamQueue(weekID, seasonID, gameDay string, isPreseason bool) []StreamGameQueueItem {
	var userGames, aiGames []StreamGameQueueItem

	isSpringGames := "N"
	if isPreseason {
		isSpringGames = "Y"
	}
	gamesFetch := repository.FindCollegeGamesRecords(repository.GamesQuery{WeekID: weekID, SeasonID: seasonID, IsSpringGames: isSpringGames})
	games := []structs.CollegeGame{}
	for _, game := range gamesFetch {
		if game.TimeSlot == "Thursday Night" && game.TimeSlot == gameDay {
			games = append(games, game)
		}
		if game.TimeSlot == "Friday Night" && game.TimeSlot == gameDay {
			games = append(games, game)
		}
		if gameDay == "Saturday Morning" && game.TimeSlot != "Thursday Night" && game.TimeSlot != "Friday Night" {
			games = append(games, game)
		}
	}
	cfbPlayByPlayMap := make(map[uint][]structs.CollegePlayByPlay)
	nflPlayByPlayMap := make(map[uint][]structs.NFLPlayByPlay)
	gameIDs := make([]string, len(games))
	for i, g := range games {
		gameIDs[i] = strconv.Itoa(int(g.ID))
	}
	playByPlays := GetCFBPlayByPlaysByGameIDs(gameIDs)
	for _, play := range playByPlays {
		cfbPlayByPlayMap[uint(play.GameID)] = append(cfbPlayByPlayMap[uint(play.GameID)], play)
	}

	teamMap := GetCollegeTeamMap()
	for _, g := range games {
		if !g.GameComplete || g.IsRevealed {
			continue
		}
		homeTeam := teamMap[uint(g.HomeTeamID)]
		awayTeam := teamMap[uint(g.AwayTeamID)]
		item := StreamGameQueueItem{
			GameID:       g.ID,
			HomeTeamID:   uint(g.HomeTeamID),
			AwayTeamID:   uint(g.AwayTeamID),
			HomeTeam:     homeTeam.TeamAbbr,
			AwayTeam:     awayTeam.TeamAbbr,
			IsUserGame:   homeTeam.Coach != "AI" || awayTeam.Coach != "AI",
			HomeTeamRank: int(g.HomeTeamRank),
			AwayTeamRank: int(g.AwayTeamRank),
			Arena:        g.Stadium,
			City:         g.City,
			State:        g.State,
			Country:      "",
			TotalSeconds: loadTotalSeconds(g.ID, cfbPlayByPlayMap, nflPlayByPlayMap, true),
			TotalPlays:   loadTotalPlays(g.ID, cfbPlayByPlayMap, nflPlayByPlayMap, true),
		}
		if item.IsUserGame {
			userGames = append(userGames, item)
		} else {
			aiGames = append(aiGames, item)
		}
	}
	return append(userGames, aiGames...)
}

func buildNFLStreamQueue(weekID, seasonID, gameDay string, isPreseason bool) []StreamGameQueueItem {
	var userGames, aiGames []StreamGameQueueItem

	gamesFetch := repository.FindNFLGamesRecords(repository.GamesQuery{WeekID: weekID, SeasonID: seasonID, TimeSlot: gameDay, IsPreseasonGame: func() string {
		if isPreseason {
			return "Y"
		} else {
			return "N"
		}
	}()})
	games := []structs.NFLGame{}
	for _, game := range gamesFetch {
		if game.TimeSlot == "Thursday Night Football" && game.TimeSlot == gameDay {
			games = append(games, game)
		}
		if game.TimeSlot == "Monday Night Football" && game.TimeSlot == gameDay {
			games = append(games, game)
		}
		if gameDay == "Sunday Noon" && game.TimeSlot != "Thursday Night Football" && game.TimeSlot != "Monday Night Football" {
			games = append(games, game)
		}
	}
	cfbPlayByPlayMap := make(map[uint][]structs.CollegePlayByPlay)
	nflPlayByPlayMap := make(map[uint][]structs.NFLPlayByPlay)
	gameIDs := make([]string, len(games))
	for i, g := range games {
		gameIDs[i] = strconv.Itoa(int(g.ID))
	}
	playByPlays := GetNFLPlayByPlaysByGameIDs(gameIDs)
	for _, play := range playByPlays {
		nflPlayByPlayMap[uint(play.GameID)] = append(nflPlayByPlayMap[uint(play.GameID)], play)
	}

	nflTeams := GetAllNFLTeams()
	teamMap := MakeNFLTeamMap(nflTeams)
	for _, g := range games {
		if !g.GameComplete || g.IsRevealed {
			continue
		}
		homeTeam := teamMap[uint(g.HomeTeamID)]
		awayTeam := teamMap[uint(g.AwayTeamID)]
		isUser := homeTeam.NFLOwnerName != "" || awayTeam.NFLOwnerName != "" ||
			homeTeam.NFLGMName != "" || awayTeam.NFLGMName != ""
		item := StreamGameQueueItem{
			GameID:       g.ID,
			HomeTeamID:   uint(g.HomeTeamID),
			AwayTeamID:   uint(g.AwayTeamID),
			HomeTeam:     homeTeam.TeamAbbr,
			AwayTeam:     awayTeam.TeamAbbr,
			IsUserGame:   isUser,
			HomeTeamRank: 0,
			AwayTeamRank: 0,
			Arena:        g.Stadium,
			City:         g.City,
			State:        g.State,
			Country:      "",
			TotalSeconds: loadTotalSeconds(g.ID, cfbPlayByPlayMap, nflPlayByPlayMap, false),
			TotalPlays:   loadTotalPlays(g.ID, cfbPlayByPlayMap, nflPlayByPlayMap, false),
		}
		if item.IsUserGame {
			userGames = append(userGames, item)
		} else {
			aiGames = append(aiGames, item)
		}
	}
	return append(userGames, aiGames...)
}
