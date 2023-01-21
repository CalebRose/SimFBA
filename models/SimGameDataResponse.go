package models

import "github.com/CalebRose/SimFBA/structs"

type SimGameDataResponse struct {
	HomeTeam       SimTeamDataResponse
	HomeTeamRoster []structs.CollegePlayer
	AwayTeam       SimTeamDataResponse
	AwayTeamRoster []structs.CollegePlayer
	GameID         int
	WeekID         int
	SeasonID       int
	GameTemp       float64
	Cloud          string
	Precip         string
	WindSpeed      float64
	WindCategory   string
}

func (sgdr *SimGameDataResponse) AssignHomeTeam(team SimTeamDataResponse, roster []structs.CollegePlayer) {
	sgdr.HomeTeam = team
	sgdr.HomeTeamRoster = roster
}

func (sgdr *SimGameDataResponse) AssignAwayTeam(team SimTeamDataResponse, roster []structs.CollegePlayer) {
	sgdr.AwayTeam = team
	sgdr.AwayTeamRoster = roster
}

func (sgdr *SimGameDataResponse) AssignWeather(temp float64, cloud string, precip string, wind string, windspeed float64) {
	sgdr.GameTemp = temp
	sgdr.Cloud = cloud
	sgdr.Precip = precip
	sgdr.WindSpeed = windspeed
	sgdr.WindCategory = wind
}
