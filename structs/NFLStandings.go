package structs

import "github.com/jinzhu/gorm"

type NFLStandings struct {
	gorm.Model
	TeamID           uint
	TeamName         string
	Mascot           string
	SeasonID         uint
	Season           uint
	LeagueID         uint
	LeagueName       string
	ConferenceID     uint
	ConferenceName   string
	TotalTies        uint
	ConferenceTies   uint
	DivisionID       uint
	DivisionName     string
	DivisionWins     uint
	DivisionLosses   uint
	DivisionTies     uint
	PostSeasonStatus string
	BaseStandings
}

func (ns *NFLStandings) UpdateNFLStandings(game NFLGame) {
	isAway := ns.TeamID == uint(game.AwayTeamID)
	winner := (!isAway && game.HomeTeamWin) || (isAway && game.AwayTeamWin)
	tie := game.HomeTeamScore == game.AwayTeamScore
	if tie {
		ns.TotalTies += 1
		if game.IsConference {
			ns.ConferenceTies += 1
		}
		if game.IsDivisional {
			ns.DivisionTies += 1
		}
		return
	}
	if winner {
		ns.TotalWins += 1
		if isAway {
			ns.AwayWins += 1
		} else {
			ns.HomeWins += 1
		}
		if game.IsConference {
			ns.ConferenceWins += 1
		}
		if game.IsDivisional {
			ns.DivisionWins += 1
		}
		if game.IsSuperBowl {
			ns.PostSeasonStatus = "Super Bowl Winner"
		}
		ns.Streak += 1
	} else {
		ns.TotalLosses += 1
		ns.Streak = 0
		if game.IsConference {
			ns.ConferenceLosses += 1
		}
		if game.IsDivisional {
			ns.DivisionLosses += 1
		}
		if game.IsSuperBowl {
			if ns.ConferenceID == 1 {
				ns.PostSeasonStatus = "NFC Champion"
			} else {
				ns.PostSeasonStatus = "AFC Champion"
			}
		}
	}
	if isAway {
		ns.PointsFor += game.AwayTeamScore
		ns.PointsAgainst += game.HomeTeamScore
	} else {
		ns.PointsFor += game.HomeTeamScore
		ns.PointsAgainst += game.AwayTeamScore
	}
}

func (ns *NFLStandings) ReduceNFLStandings(game NFLGame) {
	isAway := ns.TeamID == uint(game.AwayTeamID)
	winner := (!isAway && game.HomeTeamWin) || (isAway && game.AwayTeamWin)
	tie := game.HomeTeamScore == game.AwayTeamScore
	if tie {
		ns.TotalTies -= 1
		if game.IsConference {
			ns.ConferenceTies -= 1
		}
		if game.IsDivisional {
			ns.DivisionTies -= 1
		}
		return
	}
	if winner {
		ns.TotalWins -= 1
		if isAway {
			ns.AwayWins -= 1
		} else {
			ns.HomeWins -= 1
		}
		if game.IsConference {
			ns.ConferenceWins -= 1
		}
		if game.IsDivisional {
			ns.DivisionWins -= 1
		}
		if game.IsSuperBowl {
			ns.PostSeasonStatus = "Super Bowl Winner"
		}
		ns.Streak -= 1
	} else {
		ns.TotalLosses -= 1
		ns.Streak = 0
		if game.IsConference {
			ns.ConferenceLosses -= 1
		}
		if game.IsDivisional {
			ns.DivisionLosses -= 1
		}
		if game.IsSuperBowl {
			if ns.ConferenceID == 1 {
				ns.PostSeasonStatus = "NFC Champion"
			} else {
				ns.PostSeasonStatus = "AFC Champion"
			}
		}
	}
	if isAway {
		ns.PointsFor -= game.AwayTeamScore
		ns.PointsAgainst -= game.HomeTeamScore
	} else {
		ns.PointsFor -= game.HomeTeamScore
		ns.PointsAgainst -= game.AwayTeamScore
	}
}

func (ns *NFLStandings) ResetNFLStandings() {
	ns.TotalTies = 0
	ns.TotalLosses = 0
	ns.TotalWins = 0
	ns.DivisionLosses = 0
	ns.DivisionTies = 0
	ns.DivisionWins = 0
	ns.ConferenceLosses = 0
	ns.ConferenceWins = 0
	ns.ConferenceTies = 0
	ns.PostSeasonStatus = ""
	ns.Streak = 0
	ns.PointsFor = 0
	ns.PointsAgainst = 0
	ns.HomeWins = 0
	ns.AwayWins = 0
}
