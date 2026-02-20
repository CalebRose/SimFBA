package structs

import "github.com/jinzhu/gorm"

type CollegeGame struct {
	gorm.Model
	WeekID                   int
	Week                     int
	SeasonID                 int
	HomeTeamRank             uint
	HomeTeamID               int
	HomeTeam                 string
	HomeTeamCoach            string
	HomeTeamWin              bool
	AwayTeamRank             uint
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
	HomePreviousBye          bool
	AwayPreviousBye          bool
	ConferenceID             uint
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

func (cg *CollegeGame) AddTeam(isHome bool, id, rank int, team, coach string) {
	if isHome {
		cg.HomeTeam = team
		cg.HomeTeamID = id
		cg.HomeTeamRank = uint(rank)
		cg.HomeTeamCoach = coach
	} else {
		cg.AwayTeam = team
		cg.AwayTeamID = id
		cg.AwayTeamRank = uint(rank)
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
func (cg *CollegeGame) AssignRank(id, rank uint) {
	isHome := id == uint(cg.HomeTeamID)
	if isHome {
		cg.HomeTeamRank = rank
	} else {
		cg.AwayTeamRank = rank
	}
}

func (cg *CollegeGame) AssignByeWeek(id uint) {
	isHome := id == uint(cg.HomeTeamID)
	if isHome {
		cg.HomePreviousBye = true
	} else {
		cg.AwayPreviousBye = true
	}
}

func (cg *CollegeGame) HideScore() {
	cg.HomeTeamScore = 0
	cg.AwayTeamScore = 0
	cg.HomeTeamWin = false
	cg.AwayTeamWin = false
	cg.GameComplete = false
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

type GameRequest struct {
	gorm.Model
	HomeTeamID       uint
	AwayTeamID       uint
	SendingTeamID    uint
	RequestingTeamID uint
	IsAccepted       bool
	IsApproved       bool
	ArenaID          uint
	Arena            string
	IsNeutralSite    bool
	SeasonID         uint
	WeekID           uint
	Week             uint
}

func (g *GameRequest) Accepted() {
	g.IsAccepted = true
}

func (g *GameRequest) Approved() {
	g.IsApproved = true
}

func (g *GameRequest) UpdateRequest(arenaID uint, arena string, isNeutralSite bool, seasonID uint, weekID uint, week uint) {
	g.ArenaID = arenaID
	g.Arena = arena
	g.IsNeutralSite = isNeutralSite
	g.SeasonID = seasonID
	g.WeekID = weekID
	g.Week = week
}
