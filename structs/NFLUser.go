package structs

import "github.com/jinzhu/gorm"

type NFLUser struct {
	gorm.Model
	Username                string
	TeamID                  uint
	TeamAbbreviation        string
	IsOwner                 bool
	IsManager               bool
	IsHeadCoach             bool
	IsAssistant             bool
	TotalWins               uint
	TotalLosses             uint
	TotalTies               uint
	ConferenceChampionships uint
	PlayoffWins             uint
	PlayoffLosses           uint
	SuperBowls              uint
	SuperBowlLosses         uint
	IsActive                bool
}

func (u *NFLUser) SetTeam(r NFLRequest) {
	u.TeamID = r.NFLTeamID
	u.TeamAbbreviation = r.NFLTeamAbbreviation
	if r.IsOwner {
		u.IsOwner = true
	}
	if r.IsManager {
		u.IsManager = true
	}
	if r.IsCoach {
		u.IsHeadCoach = true
	}
	if r.IsAssistant {
		u.IsAssistant = true
	}
}

func (u *NFLUser) RemoveOwnership() {
	u.IsOwner = false

	if !u.IsHeadCoach && !u.IsManager && !u.IsAssistant {
		u.TeamID = 0
		u.TeamAbbreviation = ""
	}
}

func (u *NFLUser) RemoveManagerPosition() {
	u.IsManager = false

	if !u.IsHeadCoach && !u.IsOwner && !u.IsAssistant {
		u.TeamID = 0
		u.TeamAbbreviation = ""
	}
}

func (u *NFLUser) RemoveCoachPosition() {
	u.IsHeadCoach = false

	if !u.IsManager && !u.IsOwner && !u.IsAssistant {
		u.TeamID = 0
		u.TeamAbbreviation = ""
	}
}

func (u *NFLUser) RemoveAssistantPosition() {
	u.IsHeadCoach = false

	if !u.IsManager && !u.IsOwner && !u.IsHeadCoach {
		u.TeamID = 0
		u.TeamAbbreviation = ""
	}
}

func (nu *NFLUser) UpdateCoachRecord(game NFLGame) {
	isAway := game.AwayTeamCoach == nu.Username
	winner := (!isAway && game.HomeTeamWin) || (isAway && game.AwayTeamWin)
	if winner {
		nu.TotalWins += 1
		if game.IsConferenceChampionship {
			nu.ConferenceChampionships += 1
		}
		if game.IsPlayoffGame {
			nu.PlayoffWins += 1
		}
		if game.IsSuperBowl {
			nu.SuperBowls += 1
		}
	} else {
		nu.TotalLosses += 1
		if game.IsSuperBowl {
			nu.SuperBowlLosses += 1
		}
		if game.IsPlayoffGame {
			nu.PlayoffLosses += 1
		}
	}
}

func (nu *NFLUser) ReduceCoachRecord(game NFLGame) {
	isAway := game.AwayTeamCoach == nu.Username
	winner := (!isAway && game.HomeTeamWin) || (isAway && game.AwayTeamWin)
	if winner {
		nu.TotalWins -= 1
		if game.IsConferenceChampionship {
			nu.ConferenceChampionships -= 1
		}
		if game.IsPlayoffGame {
			nu.PlayoffWins -= 1
		}
		if game.IsSuperBowl {
			nu.SuperBowls -= 1
		}
	} else {
		nu.TotalLosses -= 1
		if game.IsSuperBowl {
			nu.SuperBowlLosses -= 1
		}
		if game.IsPlayoffGame {
			nu.PlayoffLosses -= 1
		}
	}
}
