package structs

import "github.com/jinzhu/gorm"

// Player - The generic player for the sim. Will hold foreign key references for Recruit, College, and NFL table
type Player struct {
	gorm.Model
	RecruitID       int
	CollegePlayerID int
	NFLPlayerID     int
}
