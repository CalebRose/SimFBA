package structs

import "github.com/jinzhu/gorm"

type CollegeGameplan struct {
	gorm.Model
	TeamID int
	BaseGameplan
}

type CollegeGameplanTEST struct {
	gorm.Model
	TeamID int
	BaseGameplan
}

type SchemeFits struct {
	GoodFits []string
	BadFits  []string
}
