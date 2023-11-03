package structs

import "github.com/jinzhu/gorm"

type NFLGameplan struct {
	gorm.Model
	TeamID uint
	BaseGameplan
}
