package structs

import (
	"strings"

	"github.com/jinzhu/gorm"
)

type CollegeStandings struct {
	gorm.Model
	TeamID                int
	TeamName              string
	SeasonID              int
	Season                int
	LeagueID              uint
	LeagueName            string
	ConferenceID          int
	ConferenceName        string
	DivisionID            int
	PostSeasonStatus      string
	IsFBS                 bool
	Rank                  uint
	IsConferenceChampion  bool
	ToucanRank            uint16
	PreseasonRank         uint16
	RegularSeasonRank     uint16
	SOS                   float32
	SOR                   float32
	RPI                   float32
	ConferenceStrengthAdj float32
	Tier1Wins             uint16
	Tier2Wins             uint16
	BadLosses             uint16
	BaseStandings
}

func (ns *CollegeStandings) CalculatePercentages() {
	totalGames := ns.TotalWins + ns.TotalLosses
	totalConfGames := ns.ConferenceWins + ns.ConferenceLosses
	if totalGames > 0 {
		ns.TotalWinPercentage = float32(ns.TotalWins) / float32(totalGames)
	} else {
		ns.TotalWinPercentage = 0
	}
	if totalConfGames > 0 {
		ns.ConfWinPercentage = float32(ns.ConferenceWins) / float32(totalConfGames)
	} else {
		ns.ConfWinPercentage = 0
	}
}

func (cs *CollegeStandings) UpdateCollegeStandings(game CollegeGame) {
	isAway := cs.TeamID == game.AwayTeamID
	winner := (!isAway && game.HomeTeamWin) || (isAway && game.AwayTeamWin)
	if winner && strings.Contains(game.GameTitle, "Championship") && game.IsConferenceChampionship {
		cs.IsConferenceChampion = true
	} else if winner && strings.Contains(game.GameTitle, "Championship") && game.IsNationalChampionship {
		cs.PostSeasonStatus = "National Champion"
	} else if winner && strings.Contains(game.GameTitle, "Round 1") && game.IsPlayoffGame {
		cs.PostSeasonStatus = "Playoffs"
	} else if winner && strings.Contains(game.GameTitle, "Bowl") && game.IsPlayoffGame && game.Week == 18 {
		cs.PostSeasonStatus = "Quarterfinals"
	} else if winner && strings.Contains(game.GameTitle, "Bowl") && game.IsPlayoffGame && game.Week == 19 {
		cs.PostSeasonStatus = "Semifinals"
	}

	if winner {
		cs.TotalWins += 1
		if isAway {
			cs.AwayWins += 1
		} else {
			cs.HomeWins += 1
		}
		if game.IsConference {
			cs.ConferenceWins += 1
		}
		cs.Streak += 1
	} else {
		cs.TotalLosses += 1
		cs.Streak = 0
		if game.IsConference {
			cs.ConferenceLosses += 1
		}
	}
	if isAway {
		cs.PointsFor += game.AwayTeamScore
		cs.PointsAgainst += game.HomeTeamScore
	} else {
		cs.PointsFor += game.HomeTeamScore
		cs.PointsAgainst += game.AwayTeamScore
	}

	cs.CalculatePercentages()
}

func (cs *CollegeStandings) SubtractCollegeStandings(game CollegeGame) {
	isAway := cs.TeamID == game.AwayTeamID
	winner := (!isAway && game.HomeTeamWin) || (isAway && game.AwayTeamWin)
	if winner {
		cs.TotalWins -= 1
		if isAway {
			cs.AwayWins -= 1
		} else {
			cs.HomeWins -= 1
		}
		if game.IsConference {
			cs.ConferenceWins -= 1
		}
		cs.Streak -= 1
	} else {
		cs.TotalLosses -= 1
		cs.Streak = 0
		if game.IsConference {
			cs.ConferenceLosses -= 1
		}
	}
	if isAway {
		cs.PointsFor -= game.AwayTeamScore
		cs.PointsAgainst -= game.HomeTeamScore
	} else {
		cs.PointsFor -= game.HomeTeamScore
		cs.PointsAgainst -= game.AwayTeamScore
	}
	cs.CalculatePercentages()
}

func (cs *CollegeStandings) ResetCFBStandings() {
	cs.TotalLosses = 0
	cs.TotalWins = 0
	cs.ConferenceLosses = 0
	cs.ConferenceWins = 0
	cs.PostSeasonStatus = ""
	cs.Streak = 0
	cs.PointsFor = 0
	cs.PointsAgainst = 0
	cs.HomeWins = 0
	cs.AwayWins = 0
	cs.RankedWins = 0
	cs.RankedLosses = 0
	cs.CalculatePercentages()
}

func (cs *CollegeStandings) SetCoach(coach string) {
	cs.Coach = coach
}

func (cs *CollegeStandings) AssignRank(rank int) {
	cs.Rank = uint(rank)
}

func (cs *CollegeStandings) MaskGames(wins, losses, confWins, confLosses int) {
	cs.TotalWins = wins
	cs.TotalLosses = losses
	cs.ConferenceWins = confWins
	cs.ConferenceLosses = confLosses
}
