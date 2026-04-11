package structs

import "gorm.io/gorm"

type CollegeGameplan struct {
	gorm.Model
	TeamID int
	BaseGameplan
}

type SchemeFits struct {
	GoodFits []string
	BadFits  []string
}
