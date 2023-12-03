package structs

import "gorm.io/gorm"

type RecruitingVisit struct {
	gorm.Model
	RecruitID    int
	GameID       int
	WeekID       int
	Week         int
	HomeTeamID   int
	HomeTeamAbbr string
	GameWon      bool
	ActiveVisit  bool
}

func (v *RecruitingVisit) ToggleGameWon() {
	v.GameWon = true
}

func (v *RecruitingVisit) ActivateVisit() {
	v.ActiveVisit = true
}

func (v *RecruitingVisit) CancelVisit() {
	v.ActiveVisit = false
}
