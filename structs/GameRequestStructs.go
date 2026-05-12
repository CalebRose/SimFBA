package structs

import "github.com/jinzhu/gorm"

type CFBGameRequest struct {
	GameRequest
	IsSpringGame bool
}

type NFLGameRequest struct {
	GameRequest
	IsPreseason bool
}

type GameRequest struct {
	gorm.Model
	HomeTeamID       uint
	AwayTeamID       uint
	SendingTeamID    uint
	RequestingTeamID uint
	IsAccepted       bool
	IsApproved       bool
	ArenaID          uint
	Arena            string
	IsNeutralSite    bool
	SeasonID         uint
	WeekID           uint
	Week             uint
	Timeslot         string
}

func (g *GameRequest) Accepted() {
	g.IsAccepted = true
}

func (g *GameRequest) Approved() {
	g.IsApproved = true
}

func (g *GameRequest) UpdateRequest(arenaID uint, arena string, isNeutralSite bool, seasonID uint, weekID uint, week uint) {
	g.ArenaID = arenaID
	g.Arena = arena
	g.IsNeutralSite = isNeutralSite
	g.SeasonID = seasonID
	g.WeekID = weekID
	g.Week = week
}
