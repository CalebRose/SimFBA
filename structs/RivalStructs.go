package structs

import "github.com/jinzhu/gorm"

type CollegeRival struct {
	gorm.Model
	RivalryName     string
	TrophyName      string
	TeamOneID       uint
	TeamTwoID       uint
	HasTrophy       bool
	TeamOnePriority uint
	TeamTwoPriority uint
	IsAnnualRivalry bool
	ConferenceID    uint
	PreferredWeek   uint8
	Timeslot        string
	IsNeutralSite   bool
	StadiumID       uint
}
