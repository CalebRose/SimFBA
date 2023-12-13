package managers

import (
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/models"
	"github.com/CalebRose/SimFBA/structs"
	"gorm.io/gorm"
)

// GetRecruitingProfileByTeamID
func GetOnlyRecruitingProfileByTeamID(TeamID string) structs.RecruitingTeamProfile {
	db := dbprovider.GetInstance().GetDB()

	var profile structs.RecruitingTeamProfile

	err := db.Where("id = ?", TeamID).Find(&profile).Error
	if err != nil {
		log.Fatal(err)
	}

	return profile
}

// GetRecruitingProfileForDashboardByTeamID -- Dashboard
func GetRecruitingProfileForDashboardByTeamID(TeamID string) structs.RecruitingTeamProfile {
	db := dbprovider.GetInstance().GetDB()

	var profile structs.RecruitingTeamProfile

	err := db.Preload("Affinities").Preload("Recruits.Recruit.RecruitPlayerProfiles", func(db *gorm.DB) *gorm.DB {
		return db.Order("total_points DESC").Where("total_points > 0")
	}).Where("id = ?", TeamID).Find(&profile).Error
	if err != nil {
		log.Panicln(err)
	}

	return profile
}

func GetRecruitingProfileForTeamBoardByTeamID(TeamID string) models.SimTeamBoardResponse {
	db := dbprovider.GetInstance().GetDB()

	var profile structs.RecruitingTeamProfile

	err := db.Preload("Affinities").Preload("Recruits.Recruit.RecruitPlayerProfiles", func(db *gorm.DB) *gorm.DB {
		return db.Order("total_points DESC").Where("total_points > 0")
	}).Where("id = ?", TeamID).Find(&profile).Error
	if err != nil {
		log.Panicln(err)
	}

	var teamProfileResponse models.SimTeamBoardResponse
	var crootProfiles []models.CrootProfile

	// iterate through player recruit profiles --> get recruit with preload to player profiles
	for i := 0; i < len(profile.Recruits); i++ {
		var crootProfile models.CrootProfile
		var croot models.Croot

		croot.Map(profile.Recruits[i].Recruit)

		crootProfile.Map(profile.Recruits[i], croot)

		crootProfiles = append(crootProfiles, crootProfile)
	}

	sort.Sort(models.ByCrootProfileTotal(crootProfiles))

	teamProfileResponse.Map(profile, crootProfiles)

	return teamProfileResponse
}

func GetOnlyAITeamRecruitingProfiles() []structs.RecruitingTeamProfile {
	db := dbprovider.GetInstance().GetDB()

	var profiles []structs.RecruitingTeamProfile

	err := db.Preload("Affinities").Where("is_ai = ?", true).Find(&profiles).Error
	if err != nil {
		log.Panicln(err)
	}

	return profiles
}

func GetRecruitingClassByTeamID(TeamID string) models.SimTeamBoardResponse {
	db := dbprovider.GetInstance().GetDB()

	var profile structs.RecruitingTeamProfile

	err := db.Preload("Recruits.Recruit", func(db *gorm.DB) *gorm.DB {
		return db.Order("overall DESC").Where("team_id = ? AND is_signed = true", TeamID)
	}).Where("id = ?", TeamID).Find(&profile).Error
	if err != nil {
		log.Panicln(err)
	}

	var teamProfileResponse models.SimTeamBoardResponse
	var crootProfiles []models.CrootProfile

	// iterate through player recruit profiles --> get recruit with preload to player profiles
	for i := 0; i < len(profile.Recruits); i++ {
		if profile.Recruits[i].IsSigned && profile.Recruits[i].Recruit.College == profile.TeamAbbreviation {
			var crootProfile models.CrootProfile

			var croot models.Croot

			croot.Map(profile.Recruits[i].Recruit)

			crootProfile.Map(profile.Recruits[i], croot)

			crootProfiles = append(crootProfiles, crootProfile)
		}

	}

	sort.Sort(models.ByCrootProfileTotal(crootProfiles))

	teamProfileResponse.Map(profile, crootProfiles)

	return teamProfileResponse
}

func GetRecruitingProfileForRecruitSync() []structs.RecruitingTeamProfile {
	db := dbprovider.GetInstance().GetDB()

	var profiles []structs.RecruitingTeamProfile

	err := db.Find(&profiles).Error
	if err != nil {
		log.Panicln(err)
	}

	return profiles
}

func GetRecruitingNeeds(TeamID string) map[string]int {
	needsMap := make(map[string]int)

	teamRoster := GetAllCollegePlayersByTeamId(TeamID)

	for _, player := range teamRoster {
		if player.IsRedshirting {
			continue
		}
		if (player.Year == 4 && !player.IsRedshirt) || (player.Year == 5 && player.IsRedshirt) {
			needsMap[player.Position] += 1
		}
	}

	return needsMap
}

func AddRecruitToBoard(RecruitDTO structs.CreateRecruitProfileDto) structs.RecruitPlayerProfile {
	//
	db := dbprovider.GetInstance().GetDB()

	recruitProfile := GetRecruitProfileByPlayerId(strconv.Itoa(RecruitDTO.RecruitID), strconv.Itoa(RecruitDTO.ProfileID))

	if recruitProfile.RecruitID != 0 && recruitProfile.ProfileID != 0 {
		// Replace Recruit Onto Board
		recruitProfile.ToggleRemoveFromBoard()
		db.Save(&recruitProfile)
		return recruitProfile
	}

	recruitingProfile := structs.RecruitPlayerProfile{
		SeasonID:                  RecruitDTO.SeasonID,
		RecruitID:                 RecruitDTO.RecruitID,
		ProfileID:                 RecruitDTO.ProfileID,
		TotalPoints:               0,
		CurrentWeeksPoints:        0,
		SpendingCount:             0,
		Scholarship:               false,
		ScholarshipRevoked:        false,
		RecruitingEfficiencyScore: RecruitDTO.RES,
		AffinityOneEligible:       RecruitDTO.AffinityOneEligible,
		AffinityTwoEligible:       RecruitDTO.AffinityTwoEligible,
		TeamAbbreviation:          RecruitDTO.Team,
		RemovedFromBoard:          false,
		IsSigned:                  false,
		Recruiter:                 RecruitDTO.Recruiter,
	}

	// Save
	db.Create(&recruitingProfile)

	return recruitingProfile
}

func RemoveRecruitFromBoard(updateRecruitPointsDto structs.UpdateRecruitPointsDto) structs.RecruitPlayerProfile {
	db := dbprovider.GetInstance().GetDB()

	recruitingPointsProfile := GetRecruitProfileByPlayerId(
		strconv.Itoa(updateRecruitPointsDto.RecruitID),
		strconv.Itoa(updateRecruitPointsDto.ProfileID))

	if recruitingPointsProfile.RemovedFromBoard {
		log.Panicln("Recruit has already been removed from Team Recruiting Board.")
	}

	recruitingPointsProfile.ToggleRemoveFromBoard()
	db.Save(&recruitingPointsProfile)

	return recruitingPointsProfile
}

func UpdateRecruitingProfile(updateRecruitingBoardDto structs.UpdateRecruitingBoardDTO) structs.RecruitingTeamProfile {
	db := dbprovider.GetInstance().GetDB()

	var teamID = strconv.Itoa(updateRecruitingBoardDto.TeamID)

	var teamProfile = GetOnlyRecruitingProfileByTeamID(teamID)

	var recruitProfiles = GetOnlyRecruitProfilesByTeamProfileID(teamID)

	var updatedRecruits = updateRecruitingBoardDto.Recruits

	var currentPoints float64 = 0

	for i := 0; i < len(recruitProfiles); i++ {
		updatedRecruit := GetRecruitFromRecruitsList(recruitProfiles[i].RecruitID, updatedRecruits)

		if recruitProfiles[i].CurrentWeeksPoints != updatedRecruit.CurrentWeeksPoints {
			// Allocate Points to Profile
			currentPoints += updatedRecruit.CurrentWeeksPoints
			teamProfile.AllocateSpentPoints(currentPoints)
			// If total not surpassed, allocate to the recruit and continue
			if teamProfile.SpentPoints <= teamProfile.WeeklyPoints {
				recruitProfiles[i].AllocateCurrentWeekPoints(updatedRecruit.CurrentWeeksPoints)
				recruitProfiles[i].AssignRES(teamProfile.RecruitingEfficiencyScore)
				fmt.Println("Saving recruit " + strconv.Itoa(recruitProfiles[i].RecruitID))
			} else {
				panic("Error: Allocated more points for Profile " + strconv.Itoa(teamProfile.TeamID) + " than what is allowed.")
			}
			// Save Recruit Profile
			db.Save(&recruitProfiles[i])
		} else {
			currentPoints += recruitProfiles[i].CurrentWeeksPoints
			teamProfile.AllocateSpentPoints(currentPoints)
		}
	}

	// Save team recruiting profile
	db.Save(&teamProfile)

	return teamProfile
}

func GetRecruitingClassSizeForTeams() {
	db := dbprovider.GetInstance().GetDB()
	profiles := GetRecruitingProfileForRecruitSync()

	for _, team := range profiles {
		count := 0

		players := GetAllCollegePlayersByTeamId(strconv.Itoa(int(team.ID)))

		rosterSize := len(players)

		for _, player := range players {
			if (player.Year == 4 && !player.IsRedshirt) || (player.Year == 5 && player.IsRedshirt) && player.Stars > 0 {
				count++
			}

			// if player.Year == 4 && player.IsRedshirt && player.Overall > 55 && player.Stars > 0 {
			// 	count++
			// }
		}

		rosterMinusGrads := rosterSize - count

		if rosterMinusGrads+25 > 105 {
			count = 105 - rosterMinusGrads
		} else if rosterMinusGrads+25 < 85 {
			count = 85 - rosterMinusGrads
		} else {
			count = 25
		}

		team.SetRecruitingClassSize(count)

		db.Save(&team)
	}
}

// SaveAIBehavior -- Toggle whether a Team will use AI recruiting or not
func SaveAIBehavior(profile structs.RecruitingTeamProfile) {
	db := dbprovider.GetInstance().GetDB()
	TeamID := strconv.Itoa(int(profile.TeamID))
	recruitingProfile := GetOnlyRecruitingProfileByTeamID(TeamID)
	recruitingProfile.UpdateAIBehavior(profile.IsAI, profile.AIAutoOfferscholarships, profile.AIStarMax, profile.AIStarMin, profile.AIMinThreshold, profile.AIMaxThreshold, profile.OffensiveScheme, profile.DefensiveScheme)
	db.Save(&recruitingProfile)
}
