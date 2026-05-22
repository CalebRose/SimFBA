package structs

import "github.com/jinzhu/gorm"

// NFLTagData stores the calculated tag amounts for a position group.
// One row per position group (e.g., "QB", "EDGE", "DL", "LB", "WR", etc.).
// Updated each offseason by CalculateTagValues.
type NFLTagData struct {
	gorm.Model
	Position   string
	Franchise  float64
	Transition float64
	Playtime   float64
	Basic      float64
}
