package structs

import "github.com/jinzhu/gorm"

type NFLUDFABoard struct {
	gorm.Model
	TeamID     uint
	TeamAbbr   string
	Profiles   []NFLUDFAProfile `gorm:"foreignKey:NFLUDFABoardID"`
}

type NFLUDFAProfile struct {
	gorm.Model
	NFLUDFABoardID uint
	PlayerID       uint
	PlayerName     string
	Position       string
	TeamID         uint
	TeamAbbr       string
	Points         int // The 1-20 points allocated
	IsSigned       bool
}