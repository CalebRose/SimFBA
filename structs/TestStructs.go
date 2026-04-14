package structs

import "gorm.io/gorm"

type CollegeGameplanTEST struct {
	gorm.Model
	TeamID int
	BaseGameplan
}

type CollegeTeamDepthChartTEST struct {
	gorm.Model
	TeamID            int
	DepthChartPlayers []CollegeDepthChartPositionTEST `gorm:"foreignKey:DepthChartID"`
}

type CollegeDepthChartPositionTEST struct {
	gorm.Model
	DepthChartID     int
	PlayerID         int           `gorm:"column:player_id"` // 123 -- CollegePlayerID
	Position         string        // "QB"
	PositionLevel    string        // "1"
	FirstName        string        // "David"
	LastName         string        // "Ross"
	OriginalPosition string        // The Original Position of the Player. Will only be used for STU position
	CollegePlayer    CollegePlayer `gorm:"foreignKey:PlayerID;references:PlayerID"`
}

// Update DepthChartPosition -- Updates the Player taking the position
func (dcp *CollegeDepthChartPositionTEST) UpdateDepthChartPosition(dto CollegeDepthChartPositionTEST) {
	if dcp.ID != dto.ID || dcp.DepthChartID != dto.DepthChartID {
		return
	}
	dcp.PlayerID = dto.PlayerID
	dcp.FirstName = dto.FirstName
	dcp.LastName = dto.LastName
	dcp.OriginalPosition = dto.OriginalPosition
}

// Update DepthChartPosition -- Updates the Player taking the position
func (dcp *CollegeDepthChartPositionTEST) AssignID(id uint) {
	dcp.ID = id
}

type NFLGameplanTEST struct {
	gorm.Model
	TeamID uint
	BaseGameplan
}

type NFLDepthChartTEST struct {
	gorm.Model
	TeamID            int
	DepthChartPlayers []NFLDepthChartPositionTEST `gorm:"foreignKey:DepthChartID"`
}

type NFLDepthChartPositionTEST struct {
	gorm.Model
	DepthChartID     uint
	PlayerID         uint      `gorm:"column:player_id"` // 123 -- CollegePlayerID
	Position         string    // "QB"
	PositionLevel    string    // "1"
	FirstName        string    // "David"
	LastName         string    // "Ross"
	OriginalPosition string    // The Original Position of the Player. Will only be used for STU position
	NFLPlayer        NFLPlayer `gorm:"foreignKey:PlayerID;references:PlayerID"`
}

// Update DepthChartPosition -- Updates the Player taking the position
func (dcp *NFLDepthChartPositionTEST) UpdateDepthChartPosition(dto NFLDepthChartPositionTEST) {
	if dcp.ID != dto.ID || dcp.DepthChartID != dto.DepthChartID {
		return
	}
	dcp.PlayerID = dto.PlayerID
	dcp.FirstName = dto.FirstName
	dcp.LastName = dto.LastName
	dcp.OriginalPosition = dto.OriginalPosition
}

// Update DepthChartPosition -- Updates the Player taking the position
func (dcp *NFLDepthChartPositionTEST) AssignID(id uint) {
	dcp.ID = id
}

// College Test Structs for Engine Testing
type SimTeamDataResponseTEST struct {
	TeamName        string
	Mascot          string
	Coach           string
	City            string
	State           string
	Stadium         string
	StadiumCapacity int
	ColorOne        string
	ColorTwo        string
	TeamGameplan    CollegeGameplanTEST
	TeamDepthChart  SimTeamDepthChartResponseTEST
	PreviousByeWeek bool
}

type SimTeamDepthChartResponseTEST struct {
	ID                uint
	TeamID            int
	DepthChartPlayers []SimDepthChartPosResponseTEST
}

type SimDepthChartPosResponseTEST struct {
	PlayerID      int
	Position      string
	PositionLevel string
}

func (stdr *SimTeamDataResponseTEST) Map(team CollegeTeam, gp CollegeGameplanTEST, dcr SimTeamDepthChartResponseTEST, isPrev bool) {
	stdr.TeamName = team.TeamName
	stdr.Mascot = team.Mascot
	stdr.City = team.City
	stdr.State = team.State
	stdr.Stadium = team.Stadium
	stdr.ColorOne = team.ColorOne
	stdr.ColorTwo = team.ColorTwo
	stdr.StadiumCapacity = team.StadiumCapacity
	stdr.TeamGameplan = gp
	stdr.TeamDepthChart = dcr
	stdr.PreviousByeWeek = isPrev
}

func (stdcr *SimTeamDepthChartResponseTEST) Map(dc CollegeTeamDepthChartTEST, dcp []SimDepthChartPosResponseTEST) {
	stdcr.ID = dc.ID
	stdcr.TeamID = dc.TeamID
	stdcr.DepthChartPlayers = dcp
}

func (sdcpr *SimDepthChartPosResponseTEST) Map(dc CollegeDepthChartPositionTEST) {
	sdcpr.PlayerID = dc.PlayerID
	sdcpr.Position = dc.Position
	sdcpr.PositionLevel = dc.PositionLevel
}

type SimGameDataResponseTEST struct {
	HomeTeam       SimTeamDataResponseTEST
	HomeTeamRoster []CollegePlayer
	AwayTeam       SimTeamDataResponseTEST
	AwayTeamRoster []CollegePlayer
	Stadium        Stadium
	GameID         int
	WeekID         int
	SeasonID       int
	GameTemp       float64
	Cloud          string
	Precip         string
	WindSpeed      float64
	WindCategory   string
	IsPostSeason   bool
}

func (sgdr *SimGameDataResponseTEST) AssignHomeTeam(team SimTeamDataResponseTEST, roster []CollegePlayer) {
	sgdr.HomeTeam = team
	sgdr.HomeTeamRoster = roster
}

func (sgdr *SimGameDataResponseTEST) AssignAwayTeam(team SimTeamDataResponseTEST, roster []CollegePlayer) {
	sgdr.AwayTeam = team
	sgdr.AwayTeamRoster = roster
}

func (sgdr *SimGameDataResponseTEST) AssignWeather(temp float64, cloud string, precip string, wind string, windspeed float64) {
	sgdr.GameTemp = temp
	sgdr.Cloud = cloud
	sgdr.Precip = precip
	sgdr.WindSpeed = windspeed
	sgdr.WindCategory = wind
}

func (sgdr *SimGameDataResponseTEST) AssignStadium(s Stadium) {
	sgdr.Stadium = s
}

func (sgdr *SimGameDataResponseTEST) AssignPostSeasonStatus(isPostSeason bool) {
	sgdr.IsPostSeason = isPostSeason
}

// NFL Data structs
type SimNFLTeamDataResponseTEST struct {
	TeamName        string
	Mascot          string
	Coach           string
	City            string
	State           string
	Stadium         string
	StadiumCapacity int
	ColorOne        string
	ColorTwo        string
	TeamGameplan    NFLGameplanTEST
	TeamDepthChart  SimNFLTeamDepthChartResponseTEST
	PreviousByeWeek bool
}

type SimNFLTeamDepthChartResponseTEST struct {
	ID                uint
	TeamID            int
	DepthChartPlayers []SimNFLDepthChartPosResponseTEST
}

type SimNFLDepthChartPosResponseTEST struct {
	PlayerID      uint
	Position      string
	PositionLevel string
}

func (stdr *SimNFLTeamDataResponseTEST) Map(team NFLTeam, gp NFLGameplanTEST, dcr SimNFLTeamDepthChartResponseTEST, isPrev bool) {
	stdr.TeamName = team.TeamName
	stdr.Mascot = team.Mascot
	stdr.City = team.City
	stdr.State = team.State
	stdr.Stadium = team.Stadium
	stdr.ColorOne = team.ColorOne
	stdr.ColorTwo = team.ColorTwo
	stdr.StadiumCapacity = team.StadiumCapacity
	stdr.TeamGameplan = gp
	stdr.TeamDepthChart = dcr
	stdr.PreviousByeWeek = isPrev
}

func (stdcr *SimNFLTeamDepthChartResponseTEST) Map(dc NFLDepthChartTEST, dcp []SimNFLDepthChartPosResponseTEST) {
	stdcr.ID = dc.ID
	stdcr.TeamID = dc.TeamID
	stdcr.DepthChartPlayers = dcp
}

func (sdcpr *SimNFLDepthChartPosResponseTEST) Map(dc NFLDepthChartPositionTEST) {
	sdcpr.PlayerID = dc.PlayerID
	sdcpr.Position = dc.Position
	sdcpr.PositionLevel = dc.PositionLevel
}

type SimNFLGameDataResponseTEST struct {
	HomeTeam       SimNFLTeamDataResponseTEST
	HomeTeamRoster []NFLPlayer
	AwayTeam       SimNFLTeamDataResponseTEST
	AwayTeamRoster []NFLPlayer
	Stadium        Stadium
	GameID         int
	WeekID         int
	SeasonID       int
	GameTemp       float64
	Cloud          string
	Precip         string
	WindSpeed      float64
	WindCategory   string
	IsPostSeason   bool
}

func (sgdr *SimNFLGameDataResponseTEST) AssignHomeTeam(team SimNFLTeamDataResponseTEST, roster []NFLPlayer) {
	sgdr.HomeTeam = team
	sgdr.HomeTeamRoster = roster
}

func (sgdr *SimNFLGameDataResponseTEST) AssignAwayTeam(team SimNFLTeamDataResponseTEST, roster []NFLPlayer) {
	sgdr.AwayTeam = team
	sgdr.AwayTeamRoster = roster
}

func (sgdr *SimNFLGameDataResponseTEST) AssignWeather(temp float64, cloud string, precip string, wind string, windspeed float64) {
	sgdr.GameTemp = temp
	sgdr.Cloud = cloud
	sgdr.Precip = precip
	sgdr.WindSpeed = windspeed
	sgdr.WindCategory = wind
}

func (sgdr *SimNFLGameDataResponseTEST) AssignStadium(s Stadium) {
	sgdr.Stadium = s
}

func (sgdr *SimNFLGameDataResponseTEST) AssignPostSeasonStatus(isPostSeason bool) {
	sgdr.IsPostSeason = isPostSeason
}
