package structs

import "github.com/jinzhu/gorm"

type CollegeGame struct {
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
	IsConference             bool
	IsDivisional             bool
	IsConferenceChampionship bool
	IsBowlGame               bool
	IsPlayoffGame            bool
	IsNationalChampionship   bool
	IsRivalryGame            bool
	GameComplete             bool
	IsSpringGame             bool
	GameTitle                string // For rivalry match-ups, bowl games, championships, and more
	NextGameID               uint
	NextGameHOA              string
}

func (cg *CollegeGame) UpdateScore(HomeScore int, AwayScore int) {
	cg.HomeTeamScore = HomeScore
	cg.AwayTeamScore = AwayScore
	if HomeScore > AwayScore {
		cg.HomeTeamWin = true
	} else {
		cg.AwayTeamWin = true
	}
	cg.GameComplete = true
}

func (cg *CollegeGame) UpdateCoach(TeamID int, Username string) {
	if cg.HomeTeamID == TeamID {
		cg.HomeTeamCoach = Username
	} else if cg.AwayTeamID == TeamID {
		cg.AwayTeamCoach = Username
	}
}

func (cg *CollegeGame) ApplyWeather(precip string, lowTemp float64, highTemp float64, gameTemp float64, cloud string, wind float64, windCategory string, region string) {
	cg.Precip = precip
	cg.LowTemp = lowTemp
	cg.HighTemp = highTemp
	cg.WindSpeed = wind
	cg.WindCategory = windCategory
	cg.Region = region
	cg.GameTemp = gameTemp
	cg.Cloud = cloud
}

func (cg *CollegeGame) UpdateTimeslot(ts string) {
	cg.TimeSlot = ts
}

func (cg *CollegeGame) AddTeam(isHome bool, id int, team, coach string) {
	if isHome {
		cg.HomeTeam = team
		cg.HomeTeamID = id
		cg.HomeTeamCoach = coach
	} else {
		cg.AwayTeam = team
		cg.AwayTeamID = id
		cg.AwayTeamCoach = coach
	}
}

func (cg *CollegeGame) AddLocation(stadiumID int, stadium, city, state string, isDomed bool) {
	cg.StadiumID = uint(stadiumID)
	cg.Stadium = stadium
	cg.City = city
	cg.State = state
	cg.IsDomed = isDomed
}

type WeatherResponse struct {
	LowTemp      float64
	HighTemp     float64
	GameTemp     float64
	Cloud        string
	Precip       string
	WindSpeed    float64
	WindCategory string
}
