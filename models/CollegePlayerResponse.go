package models

import "github.com/CalebRose/SimFBA/structs"

type CollegePlayerResponse struct {
	ID int
	structs.BasePlayer
	TeamID       int
	TeamAbbr     string
	City         string
	State        string
	Year         int
	IsRedshirt   bool
	ConferenceID int
	Conference   string
	Stats        structs.CollegePlayerStats
	SeasonStats  structs.CollegePlayerSeasonStats
}
