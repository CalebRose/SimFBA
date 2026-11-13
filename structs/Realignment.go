package structs

import "github.com/jinzhu/gorm"

type BaseRealignment struct {
	gorm.Model
	TeamID       uint
	ConferenceID uint
	DivisionID   uint
	SeasonID     uint
	IsProcessed  bool
}

type CFBRealignment struct {
	BaseRealignment
}

type NFLRealignment struct {
	BaseRealignment
}
