package structs

import "github.com/jinzhu/gorm"

type NFLDraftPick struct {
	gorm.Model
	OriginalTeamID uint
	OriginalTeam   string
	PreviousTeamID uint
	PreviousTeam   string
	TeamID         uint
	Team           string
	PickSelection  string
	Notes          string
	PlayerID       uint
	Round          uint
	PickNumber     uint
	SeasonID       uint
	Season         uint
	IsCompensation bool
	IsVoid         bool
	TradeValue     float64
}

func (p *NFLDraftPick) TradePick(id uint, team string) {
	p.PreviousTeamID = p.TeamID
	p.PreviousTeam = p.Team
	p.TeamID = id
	p.Team = team
	if p.PreviousTeamID == p.OriginalTeamID {
		p.Notes = "From " + p.OriginalTeam
	} else {
		p.Notes = "From " + p.PreviousTeam + " via " + p.OriginalTeam
	}
}
