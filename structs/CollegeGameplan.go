package structs

import "github.com/jinzhu/gorm"

type CollegeGameplan struct {
	gorm.Model
	TeamID int
	BaseGameplan
}
