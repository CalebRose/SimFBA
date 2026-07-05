package managers

import "github.com/CalebRose/SimFBA/repository"

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
	games := repository.FindCollegeGamesRecords(repository.GamesQuery{WeekID: weekID, SeasonID: seasonID, TimeSlot: gameDay, IsSpringGames: isSpringGames})
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
			TotalSeconds: loadTotalSeconds(g.ID, true),
			TotalPlays:   loadTotalPlays(g.ID, true),
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

	games := repository.FindNFLGamesRecords(repository.GamesQuery{WeekID: weekID, SeasonID: seasonID, TimeSlot: gameDay, IsPreseasonGame: func() string {
		if isPreseason {
			return "Y"
		} else {
			return "N"
		}
	}()})
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
			TotalSeconds: loadTotalSeconds(g.ID, false),
			TotalPlays:   loadTotalPlays(g.ID, false),
		}
		if item.IsUserGame {
			userGames = append(userGames, item)
		} else {
			aiGames = append(aiGames, item)
		}
	}
	return append(userGames, aiGames...)
}
