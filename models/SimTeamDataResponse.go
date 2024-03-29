package models

import "github.com/CalebRose/SimFBA/structs"

type SimTeamDataResponse struct {
	TeamName        string
	Mascot          string
	Coach           string
	City            string
	State           string
	Stadium         string
	StadiumCapacity int
	ColorOne        string
	ColorTwo        string
	TeamGameplan    structs.CollegeGameplan
	TeamDepthChart  SimTeamDepthChartResponse
	PreviousByeWeek bool
}

type SimTeamDepthChartResponse struct {
	ID                uint
	TeamID            int
	DepthChartPlayers []SimDepthChartPosResponse
}

type SimDepthChartPosResponse struct {
	PlayerID      int
	Position      string
	PositionLevel string
}

func (stdr *SimTeamDataResponse) Map(team structs.CollegeTeam, dcr SimTeamDepthChartResponse, isPrev bool) {
	stdr.TeamName = team.TeamName
	stdr.Mascot = team.Mascot
	stdr.City = team.City
	stdr.State = team.State
	stdr.Stadium = team.Stadium
	stdr.ColorOne = team.ColorOne
	stdr.ColorTwo = team.ColorTwo
	stdr.StadiumCapacity = team.StadiumCapacity
	stdr.TeamGameplan = team.TeamGameplan
	stdr.TeamDepthChart = dcr
	stdr.PreviousByeWeek = isPrev
}

func (stdcr *SimTeamDepthChartResponse) Map(dc structs.CollegeTeamDepthChart, dcp []SimDepthChartPosResponse) {
	stdcr.ID = dc.ID
	stdcr.TeamID = dc.TeamID
	stdcr.DepthChartPlayers = dcp
}

func (sdcpr *SimDepthChartPosResponse) Map(dc structs.CollegeDepthChartPosition) {
	sdcpr.PlayerID = dc.PlayerID
	sdcpr.Position = dc.Position
	sdcpr.PositionLevel = dc.PositionLevel
}

type NFLSimTeamDataResponse struct {
	TeamName        string
	Mascot          string
	Coach           string
	City            string
	State           string
	Stadium         string
	StadiumCapacity int
	ColorOne        string
	ColorTwo        string
	TeamGameplan    structs.NFLGameplan
	TeamDepthChart  NFLSimTeamDepthChartResponse
}

type NFLSimTeamDepthChartResponse struct {
	ID                uint
	TeamID            int
	DepthChartPlayers []NFLSimDepthChartPosResponse
}

type NFLSimDepthChartPosResponse struct {
	PlayerID      int
	Position      string
	PositionLevel string
}

func (stdr *NFLSimTeamDataResponse) Map(team structs.NFLTeam, dcr NFLSimTeamDepthChartResponse) {
	stdr.TeamName = team.TeamName
	stdr.Mascot = team.Mascot
	stdr.City = team.City
	stdr.State = team.State
	stdr.Stadium = team.Stadium
	stdr.ColorOne = team.ColorOne
	stdr.ColorTwo = team.ColorTwo
	stdr.StadiumCapacity = team.StadiumCapacity
	stdr.TeamGameplan = team.TeamGameplan
	stdr.TeamDepthChart = dcr
}

func (stdcr *NFLSimTeamDepthChartResponse) Map(dc structs.NFLDepthChart, dcp []NFLSimDepthChartPosResponse) {
	stdcr.ID = dc.ID
	stdcr.TeamID = dc.TeamID
	stdcr.DepthChartPlayers = dcp
}

func (sdcpr *NFLSimDepthChartPosResponse) Map(dc structs.NFLDepthChartPosition) {
	sdcpr.PlayerID = int(dc.PlayerID)
	sdcpr.Position = dc.Position
	sdcpr.PositionLevel = dc.PositionLevel
}

/*
TEST PURPOSES ONLY
*/

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
	TeamGameplan    structs.CollegeGameplanTEST
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

func (stdr *SimTeamDataResponseTEST) Map(team structs.CollegeTeam, gp structs.CollegeGameplanTEST, dcr SimTeamDepthChartResponseTEST, isPrev bool) {
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

func (stdcr *SimTeamDepthChartResponseTEST) Map(dc structs.CollegeTeamDepthChartTEST, dcp []SimDepthChartPosResponseTEST) {
	stdcr.ID = dc.ID
	stdcr.TeamID = dc.TeamID
	stdcr.DepthChartPlayers = dcp
}

func (sdcpr *SimDepthChartPosResponseTEST) Map(dc structs.CollegeDepthChartPositionTEST) {
	sdcpr.PlayerID = dc.PlayerID
	sdcpr.Position = dc.Position
	sdcpr.PositionLevel = dc.PositionLevel
}
