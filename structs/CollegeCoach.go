package structs

import (
	"github.com/jinzhu/gorm"
)

type CollegeCoach struct {
	gorm.Model
	CoachName                      string
	Age                            int
	TeamID                         uint
	Team                           string
	AlmaMaterID                    uint
	AlmaMater                      string
	FormerPlayerID                 uint
	Prestige                       int    // Level system. Every 6 wins, every playoff win & conference tourney championship nets 1 point which can then be applied towards one of the five odds categories, if they qualify
	PointMin                       int    // Minimum number of points the coach will put towards a player
	PointMax                       int    // Maximum number of points the coach will place on a player
	StarMin                        int    // Minimum star rating they will target on a croot (floor)
	StarMax                        int    // Maximum star rating they will target on a croot (ceiling)
	Odds1                          int    // Modifier towards adding 1 star croots to board
	Odds2                          int    // Modifier towards adding 2 star croots to board
	Odds3                          int    // Modifier towards adding 3 star croots to board
	Odds4                          int    // Modifier towards adding 4 star croots to board
	Odds5                          int    // Modifier towards adding 5 star croots to board
	PositionOne                    string // Preferred positions to recruit. Small 5% modifier at recruiting positions the coach specializes in
	PositionTwo                    string
	PositionThree                  string
	OffensiveScheme                string // Desired offensive scheme the coach wants to run -- will recruit based on the desired scheme
	DefensiveScheme                string // Desired defensive scheme the coach wants to run -- will recruit based on the desired scheme
	TeambuildingPreference         string // Coaches that prefer to recruit vs Coaches that will utilize the transfer portal
	CareerPreference               string // "Prefers to stay at their current job", "Wants to coach at Alma Mater", "Wants a more competitive job", "Average"
	PromiseTendency                string // Coach will either under-promise, over-promise, or be average on promises within transfer portal
	PortalReputation               int    // A value between 1-100 signifying the coach's reputation and behavior in the transfer portal.
	SchoolTenure                   int    // Number of years coach is participating on the team
	CareerTenure                   int    // Number of years the coach is actively coaching in the college sim
	ContractLength                 int    // Number of years of the current contract
	YearsRemaining                 int    // Years left on current contract
	OverallWins                    int
	OverallLosses                  int
	OverallConferenceChampionships int
	BowlWins                       int
	BowlLosses                     int
	PlayoffWins                    int
	PlayoffLosses                  int
	NationalChampionships          int
	IsUser                         bool
	IsActive                       bool
	IsRetired                      bool
	IsFormerPlayer                 bool
}

func (cc *CollegeCoach) SetTeam(TeamID uint) {
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

func (c *CollegeCoach) IncrementOdds(star int) {
	if star == 1 {
		c.Odds1 += 1
	} else if star == 2 {
		c.Odds2 += 1
	} else if star == 3 {
		c.Odds3 += 1
	} else if star == 4 {
		c.Odds4 += 1
	} else if star == 5 {
		c.Odds5 += 1
	}
}
