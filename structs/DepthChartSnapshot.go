package structs

type CollegeDepthChartSeasonSnapshot struct {
	ID               uint `gorm:"primaryKey"`
	SeasonID         uint `gorm:"index:idx_cfb_snap_season"`
	DepthChartID     int  `gorm:"index:idx_cfb_snap_dc"`
	PlayerID         int  `gorm:"column:player_id"`
	Position         string
	PositionLevel    string
	FirstName        string
	LastName         string
	OriginalPosition string
}

func (CollegeDepthChartSeasonSnapshot) TableName() string {
	return "college_depth_chart_season_snapshot"
}

type NFLDepthChartSeasonSnapshot struct {
	ID               uint `gorm:"primaryKey"`
	SeasonID         uint `gorm:"index:idx_nfl_snap_season"`
	DepthChartID     uint `gorm:"index:idx_nfl_snap_dc"`
	PlayerID         uint `gorm:"column:player_id"`
	Position         string
	PositionLevel    string
	FirstName        string
	LastName         string
	OriginalPosition string
}

func (NFLDepthChartSeasonSnapshot) TableName() string {
	return "nfl_depth_chart_season_snapshot"
}
