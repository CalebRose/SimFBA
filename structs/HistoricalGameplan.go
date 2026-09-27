package structs

import "gorm.io/gorm"

type CollegeGameplanRecord struct {
    gorm.Model
    GameID       uint   `gorm:"index"`
    TeamID       uint   `gorm:"index"`
    SeasonID     uint   `gorm:"index"`
    WeekID       uint   `gorm:"index"`
    IsHome       bool
    GameplanJSON string `gorm:"type:text"`
}

type NFLGameplanRecord struct {
    gorm.Model
    GameID       uint   `gorm:"index"`
    TeamID       uint   `gorm:"index"`
    SeasonID     uint   `gorm:"index"`
    WeekID       uint   `gorm:"index"`
    IsHome       bool
    GameplanJSON string `gorm:"type:text"`
}
