package structs

import (
	"github.com/jinzhu/gorm"
)

type CollegeCoach struct {
	gorm.Model
	CoachName                      string
	TeamID                         int
	OverallWins                    int
	OverallLosses                  int
	OverallConferenceChampionships int
	BowlWins                       int
	BowlLosses                     int
	PlayoffWins                    int
	PlayoffLosses                  int
	NationalChampionships          int
	IsActive                       bool
}

func (cc *CollegeCoach) SetTeam(TeamID int) {
	cc.TeamID = TeamID
	cc.IsActive = true
}

func (cc *CollegeCoach) SetAsInactive() {
	cc.IsActive = false
	cc.TeamID = 0
}

func (cc *CollegeCoach) UpdateCoachRecord(game CollegeGame) {
	isAway := game.AwayTeamCoach == cc.CoachName
	winner := (!isAway && game.HomeTeamWin) || (isAway && game.AwayTeamWin)
	if winner {
		cc.OverallWins += 1
		if game.IsConferenceChampionship {
			cc.OverallConferenceChampionships += 1
		}
		if game.IsBowlGame {
			cc.BowlWins += 1
		}
		if game.IsPlayoffGame {
			cc.PlayoffWins += 1
		}
		if game.IsNationalChampionship {
			cc.NationalChampionships += 1
		}
	} else {
		cc.OverallLosses += 1
		if game.IsBowlGame {
			cc.BowlLosses += 1
		}
		if game.IsPlayoffGame {
			cc.PlayoffLosses += 1
		}
	}
}

func (cc *CollegeCoach) ReduceCoachRecord(game CollegeGame) {
	isAway := game.AwayTeamCoach == cc.CoachName
	winner := (!isAway && game.HomeTeamWin) || (isAway && game.AwayTeamWin)
	if winner {
		cc.OverallWins -= 1
		if game.IsConferenceChampionship {
			cc.OverallConferenceChampionships -= 1
		}
		if game.IsBowlGame {
			cc.BowlWins -= 1
		}
		if game.IsPlayoffGame {
			cc.PlayoffWins -= 1
		}
		if game.IsNationalChampionship {
			cc.NationalChampionships -= 1
		}
	} else {
		cc.OverallLosses -= 1
		if game.IsBowlGame {
			cc.BowlLosses -= 1
		}
		if game.IsPlayoffGame {
			cc.PlayoffLosses -= 1
		}
	}
}
