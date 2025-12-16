package structs

import "github.com/jinzhu/gorm"

type RecruitPointAllocation struct {
	gorm.Model
	RecruitID          uint
	TeamProfileID      uint
	RecruitProfileID   uint
	WeekID             uint
	Points             float32
	RESAffectedPoints  float32
	RES                float32
	AffinityOneApplied bool
	AffinityTwoApplied bool
	CaughtCheating     bool
}

func (rpa *RecruitPointAllocation) UpdatePointsSpent(points float64, res float64) {
	rpa.Points = float32(points)
	rpa.RESAffectedPoints = float32(res)
}

func (rpa *RecruitPointAllocation) ApplyAffinityOne() {
	rpa.AffinityOneApplied = true
}

func (rpa *RecruitPointAllocation) ApplyAffinityTwo() {
	rpa.AffinityTwoApplied = true
}

func (rpa *RecruitPointAllocation) ApplyCaughtCheating() {
	rpa.CaughtCheating = true
}
