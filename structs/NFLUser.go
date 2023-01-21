package structs

import "github.com/jinzhu/gorm"

type NFLUser struct {
	gorm.Model
	Username         string
	TeamID           uint
	TeamAbbreviation string
	IsOwner          bool
	IsManager        bool
	IsHeadCoach      bool
	IsAssistant      bool
	TotalWins        uint
	TotalLosses      uint
	TotalTies        uint
	IsActive         bool
}

func (u *NFLUser) SetTeam(r NFLRequest) {
	u.TeamID = r.NFLTeamID
	u.TeamAbbreviation = r.NFLTeamAbbreviation
	if r.IsOwner {
		u.IsOwner = true
	}
	if r.IsManager {
		u.IsManager = true
	}
	if r.IsCoach {
		u.IsHeadCoach = true
	}
	if r.IsAssistant {
		u.IsAssistant = true
	}
}

func (u *NFLUser) RemoveOwnership() {
	u.IsOwner = false

	if !u.IsHeadCoach && !u.IsManager && !u.IsAssistant {
		u.TeamID = 0
		u.TeamAbbreviation = ""
	}
}

func (u *NFLUser) RemoveManagerPosition() {
	u.IsManager = false

	if !u.IsHeadCoach && !u.IsOwner && !u.IsAssistant {
		u.TeamID = 0
		u.TeamAbbreviation = ""
	}
}

func (u *NFLUser) RemoveCoachPosition() {
	u.IsHeadCoach = false

	if !u.IsManager && !u.IsOwner && !u.IsAssistant {
		u.TeamID = 0
		u.TeamAbbreviation = ""
	}
}

func (u *NFLUser) RemoveAssistantPosition() {
	u.IsHeadCoach = false

	if !u.IsManager && !u.IsOwner && !u.IsHeadCoach {
		u.TeamID = 0
		u.TeamAbbreviation = ""
	}
}
