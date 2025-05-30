package structs

import "github.com/jinzhu/gorm"

type NFLRetiredPlayer struct {
	gorm.Model
	BasePlayer
	PlayerID          int
	TeamID            int
	CollegeID         uint
	College           string
	TeamAbbr          string
	Experience        uint
	HighSchool        string
	Hometown          string
	State             string
	IsActive          bool
	IsPracticeSquad   bool
	IsFreeAgent       bool
	IsWaived          bool
	IsOnTradeBlock    bool
	IsAcceptingOffers bool
	IsNegotiating     bool
	NegotiationRound  uint
	SigningRound      uint
	MinimumValue      float64
	AAV               float64
	DraftedTeamID     uint
	DraftedTeam       string
	DraftedRound      uint
	DraftPickID       uint
	DraftedPick       uint
	ShowLetterGrade   bool
	HasProgressed     bool
	Rejections        int
	ProBowls          uint8
	TagType           uint8                // 0 == Basic, 1 == Franchise, 2 == Transition, 3 == Playtime
	Stats             []NFLPlayerStats     `gorm:"foreignKey:NFLPlayerID"`
	SeasonStats       NFLPlayerSeasonStats `gorm:"foreignKey:NFLPlayerID"`
	Contract          NFLContract          `gorm:"foreignKey:NFLPlayerID"`
	Offers            []FreeAgencyOffer    `gorm:"foreignKey:NFLPlayerID"`
	WaiverOffers      []NFLWaiverOffer     `gorm:"foreignKey:NFLPlayerID"`
	Extensions        []NFLExtensionOffer  `gorm:"foreignKey:NFLPlayerID"`
}
