package structs

import "github.com/jinzhu/gorm"

type NFLGame struct {
	gorm.Model
	WeekID                   int
	Week                     int
	SeasonID                 int
	HomeTeamID               int
	HomeTeam                 string
	HomeTeamCoach            string
	HomeTeamWin              bool
	AwayTeamID               int
	AwayTeam                 string
	AwayTeamCoach            string
	AwayTeamWin              bool
	MVP                      string
	HomeTeamScore            int
	AwayTeamScore            int
	TimeSlot                 string
	StadiumID                uint
	Stadium                  string
	City                     string
	State                    string
	Region                   string
	LowTemp                  float64
	HighTemp                 float64
	GameTemp                 float64
	Cloud                    string
	Precip                   string
	WindSpeed                float64
	WindCategory             string
	IsNeutral                bool
	IsDomed                  bool
	IsNightGame              bool
	IsPlayoffGame            bool
	IsRivalryGame            bool
	IsConference             bool
	IsDivisional             bool
	IsConferenceChampionship bool
	IsSuperBowl              bool
	GameComplete             bool
	IsPreseasonGame          bool
	GameTitle                string // For rivalry match-ups, bowl games, championships, and more
}

func (ng *NFLGame) UpdateScore(HomeScore int, AwayScore int) {
	ng.HomeTeamScore = HomeScore
	ng.AwayTeamScore = AwayScore
	if HomeScore > AwayScore {
		ng.HomeTeamWin = true
	} else {
		ng.AwayTeamWin = true
	}
	ng.GameComplete = true
}

func (ng *NFLGame) UpdateCoach(TeamID int, Username string) {
	if ng.HomeTeamID == TeamID {
		ng.HomeTeamCoach = Username
	} else if ng.AwayTeamID == TeamID {
		ng.AwayTeamCoach = Username
	}
}

func (ng *NFLGame) ApplyWeather(precip string, lowTemp float64, highTemp float64, gameTemp float64, cloud string, wind float64, windCategory string, region string) {
	ng.Precip = precip
	ng.LowTemp = lowTemp
	ng.HighTemp = highTemp
	ng.WindSpeed = wind
	ng.WindCategory = windCategory
	ng.Region = region
	ng.GameTemp = gameTemp
	ng.Cloud = cloud
}

func (ng *NFLGame) UpdateTimeslot(ts string) {
	ng.TimeSlot = ts
}
