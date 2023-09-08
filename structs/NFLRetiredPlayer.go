package structs

import "github.com/jinzhu/gorm"

type NFLRetiredPlayer struct {
	gorm.Model
	BasePlayer
	PlayerID          int
	TeamID            int
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
	PreviousTeamID    uint
	PreviousTeam      string
	DraftedTeamID     uint
	DraftedTeam       string
	DraftedRound      uint
	DraftedPick       uint
	ShowLetterGrade   bool
	Stats             []NFLPlayerStats     `gorm:"foreignKey:NFLPlayerID"`
	SeasonStats       NFLPlayerSeasonStats `gorm:"foreignKey:NFLPlayerID"`
	Contract          NFLContract          `gorm:"foreignKey:NFLPlayerID"`
	Offers            []FreeAgencyOffer    `gorm:"foreignKey:NFLPlayerID"`
	WaiverOffers      []NFLWaiverOffer     `gorm:"foreignKey:NFLPlayerID"`
	Extensions        []NFLExtensionOffer  `gorm:"foreignKey:NFLPlayerID"`
}
