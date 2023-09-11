package structs

import "github.com/jinzhu/gorm"

type NFLTeam struct {
	gorm.Model
	BaseTeam
	ConferenceID     uint
	Conference       string
	DivisionID       uint
	Division         string
	NFLOwnerID       uint
	NFLOwnerName     string
	NFLCoachID       uint
	NFLCoachName     string
	NFLGMID          uint
	NFLGMName        string
	NFLAssistantID   uint
	NFLAssistantName string
	WaiverOrder      uint
	OffersAccepted   int
	Capsheet         NFLCapsheet          `gorm:"foreignKey:NFLTeamID"`
	Contracts        []NFLContract        `gorm:"foreignKey:TeamID"`
	DraftPicks       []NFLDraftPick       `gorm:"foreignKey:TeamID"`
	TeamStats        []NFLTeamStats       `gorm:"foreignKey:TeamID"`
	TeamSeasonStats  []NFLTeamSeasonStats `gorm:"foreignKey:TeamID"`
	TeamDepthChart   NFLDepthChart        `gorm:"foreignKey:TeamID"`
	TeamGameplan     NFLGameplan          `gorm:"foreignKey:TeamID"`
	Standings        []NFLStandings       `gorm:"foreignKey:TeamID"`

	// Offers           []FreeAgencyOffer `gorm:"foreignKey:TeamID"`
	// NFLCoach     NFLUser        `gorm:"foreignKey:TeamID"`
	// NFLOwner     NFLUser        `gorm:"foreignKey:TeamID"`
	// NFLGM        NFLUser        `gorm:"foreignKey:TeamID"`
}

func (t *NFLTeam) AssignNFLUserToTeam(r NFLRequest, u NFLUser) {
	if r.IsOwner {
		t.NFLOwnerID = u.ID
		t.NFLOwnerName = r.Username
	}
	if r.IsManager {
		t.NFLGMID = u.ID
		t.NFLGMName = r.Username
	}
	if r.IsCoach {
		t.NFLCoachID = u.ID
		t.NFLCoachName = r.Username
	}
	if r.IsAssistant {
		t.NFLAssistantID = u.ID
		t.NFLAssistantName = r.Username
	}
}

func (t *NFLTeam) RemoveNFLUserFromTeam(r NFLRequest, u NFLUser) {
	if r.IsOwner && t.NFLOwnerName == u.Username {
		t.NFLOwnerID = 0
		t.NFLOwnerName = ""
	}
	if r.IsManager && t.NFLGMName == u.Username {
		t.NFLGMID = 0
		t.NFLGMName = ""
	}
	if r.IsCoach && t.NFLCoachName == u.Username {
		t.NFLCoachID = 0
		t.NFLCoachName = ""
	}
	if r.IsAssistant && t.NFLAssistantName == u.Username {
		t.NFLAssistantID = 0
		t.NFLAssistantName = ""
	}
}

func (t *NFLTeam) AssignWaiverOrder(val uint) {
	t.WaiverOrder = val
}

func (t *NFLTeam) ResetSeasonData() {
	t.OffersAccepted = 0
	t.WaiverOrder = 0
}

func (t *NFLTeam) IncrementExtensionOffers() {
	t.OffersAccepted += 1
}
