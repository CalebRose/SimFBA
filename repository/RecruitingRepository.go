package repository

import (
	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/structs"
)

type RecruitingClauses struct {
	RecruitID        string
	WeekID           string
	ProfileID        string
	RecruitProfileID string
	IncludeRecruit   bool
	OrderByPoints    bool
	RemoveFromBoard  bool
	CaughtCheating   bool
}

func FindRecruitPlayerProfileRecords(profileID, recruitID string, includeRecruit, orderByPoints, removeFromBoard bool) []structs.RecruitPlayerProfile {
	db := dbprovider.GetInstance().GetDB()

	var croots []structs.RecruitPlayerProfile

	query := db.Model(&croots)

	if includeRecruit {
		query = query.Preload("Recruit")
	}

	if len(profileID) > 0 {
		query = query.Where("profile_id = ?", profileID)
	}

	if len(recruitID) > 0 {
		query = query.Where("recruit_id = ?", recruitID)
	}

	if removeFromBoard {
		query = query.Where("removed_from_board = ?", false)
	}

	if orderByPoints {
		query = query.Order("total_points desc")
	}

	if err := query.Find(&croots).Error; err != nil {
		return []structs.RecruitPlayerProfile{}
	}

	return croots
}

func FindRecruitPointAllocationRecords(clauses RecruitingClauses) []structs.RecruitPointAllocation {
	db := dbprovider.GetInstance().GetDB()

	var croots []structs.RecruitPointAllocation

	query := db.Model(&croots)

	if clauses.IncludeRecruit {
		query = query.Preload("Recruit")
	}

	if len(clauses.WeekID) > 0 {
		query = query.Where("week_id >= ?", clauses.WeekID)
	}

	if len(clauses.ProfileID) > 0 {
		query = query.Where("profile_id = ?", clauses.ProfileID)
	}

	if clauses.CaughtCheating {
		query = query.Where("caught_cheating = ?", true)
	}

	if clauses.RecruitProfileID != "" {
		query = query.Where("recruit_profile_id = ?", clauses.RecruitProfileID)
	}

	if len(clauses.RecruitID) > 0 {
		query = query.Where("recruit_id = ?", clauses.RecruitID)
	}

	if clauses.RemoveFromBoard {
		query = query.Where("removed_from_board = ?", false)
	}

	if clauses.OrderByPoints {
		query = query.Order("total_points desc")
	}

	if err := query.Find(&croots).Error; err != nil {
		return []structs.RecruitPointAllocation{}
	}

	return croots
}
