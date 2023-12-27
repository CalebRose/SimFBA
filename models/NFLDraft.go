package models

import (
	"github.com/CalebRose/SimFBA/structs"
	"github.com/jinzhu/gorm"
)

type NFLDraftPageResponse struct {
	WarRoom          NFLWarRoom
	DraftablePlayers []NFLDraftee
	NFLTeams         []structs.NFLTeam
	AllDraftPicks    [7][]structs.NFLDraftPick
	CollegeTeams     []structs.CollegeTeam
}

type NFLWarRoom struct {
	gorm.Model
	TeamID         uint
	Team           string
	ScoutingPoints uint
	SpentPoints    uint
	DraftPicks     []structs.NFLDraftPick `gorm:"foreignKey:TeamID"`
	ScoutProfiles  []ScoutingProfile      `gorm:"foreignKey:TeamID"`
}

func (w *NFLWarRoom) ResetSpentPoints() {
	w.SpentPoints = 0
}

func (w *NFLWarRoom) AddToSpentPoints(points uint) {
	w.SpentPoints += points
}

type ScoutingProfile struct {
	gorm.Model
	PlayerID         uint
	TeamID           uint
	ShowAttribute1   bool
	ShowAttribute2   bool
	ShowAttribute3   bool
	ShowAttribute4   bool
	ShowAttribute5   bool
	ShowAttribute6   bool
	ShowAttribute7   bool
	ShowAttribute8   bool
	ShowPotential    bool
	RemovedFromBoard bool
	ShowCount        uint
	Draftee          NFLDraftee `gorm:"foreignKey:PlayerID;references:PlayerID"`
}

func (sp *ScoutingProfile) RevealAttribute(attr string) {
	if attr == "1" {
		sp.ShowAttribute1 = true
	} else if attr == "2" {
		sp.ShowAttribute2 = true
	} else if attr == "3" {
		sp.ShowAttribute3 = true
	} else if attr == "4" {
		sp.ShowAttribute4 = true
	} else if attr == "5" {
		sp.ShowAttribute5 = true
	} else if attr == "6" {
		sp.ShowAttribute6 = true
	} else if attr == "7" {
		sp.ShowAttribute7 = true
	} else if attr == "8" {
		sp.ShowAttribute8 = true
	} else if attr == "Potential" {
		sp.ShowPotential = true
	}
	sp.ShowCount++
}

func (sp *ScoutingProfile) RemoveFromBoard() {
	sp.RemovedFromBoard = true
}

func (sp *ScoutingProfile) ReplaceOnBoard() {
	sp.RemovedFromBoard = false
}

type ScoutingProfileDTO struct {
	PlayerID uint
	TeamID   uint
}

type RevealAttributeDTO struct {
	ScoutProfileID uint
	Attribute      string
	Points         uint
	TeamID         uint
}

type ScoutingDataResponse struct {
	DrafteeSeasonStats structs.CollegePlayerSeasonStats
	TeamStandings      structs.CollegeStandings
}

type ExportDraftPicksDTO struct {
	DraftPicks []structs.NFLDraftPick
}
