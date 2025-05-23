package structs

import "github.com/jinzhu/gorm"

type CollegeStandings struct {
	gorm.Model
	TeamID           int
	TeamName         string
	SeasonID         int
	Season           int
	LeagueID         uint
	LeagueName       string
	ConferenceID     int
	ConferenceName   string
	DivisionID       int
	PostSeasonStatus string
	IsFBS            bool
	Rank             uint
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
		ns.ConfWinPercentage = float32(ns.TotalWins) / float32(totalConfGames)
	} else {
		ns.ConfWinPercentage = 0
	}
}

func (cs *CollegeStandings) UpdateCollegeStandings(game CollegeGame) {
	isAway := cs.TeamID == game.AwayTeamID
	winner := (!isAway && game.HomeTeamWin) || (isAway && game.AwayTeamWin)
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
