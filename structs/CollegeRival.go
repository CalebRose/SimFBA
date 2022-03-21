package structs

import "github.com/jinzhu/gorm"

type CollegeRival struct {
	gorm.Model
	RivalryName   string
	TeamID        int
	Team          string
	TeamWins      int
	RivalID       int
	RivalTeam     string
	RivalWins     int
	TeamStreak    int
	RivalStreak   int
	CurrentStreak int
	HasTrophy     bool
	LatestVictor  string
}
