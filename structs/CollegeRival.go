package structs

import "github.com/jinzhu/gorm"

type CollegeRival struct {
	gorm.Model
	RivalryName    string
	TeamOneID      int
	TeamOne        string
	TeamOneWins    int
	TeamTwoID      int
	TeamTwo        string
	TeamTwoWins    int
	DesignatedWeek int
	TeamOneStreak  int
	TeamTwoStreak  int
	CurrentStreak  int
	HasTrophy      bool
	LatestVictor   string
}
