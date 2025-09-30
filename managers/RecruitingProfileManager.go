package managers

import (
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
	"gorm.io/gorm"
)

func RecalibrateCrootProfiles() {
	// Ensure that affinities are calibrated properly
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	UnsignedRecruits := GetAllUnsignedRecruits()
	stateMatcher := util.GetStateMatcher()
	regionMatcher := util.GetStateRegionMatcher()

	recruitInfos := make(map[uint]structs.RecruitInfo)
	crootMap := make(map[uint]structs.Recruit)
	fmt.Println("Loading recruits...")
	for _, croot := range UnsignedRecruits {
		info := structs.RecruitInfo{
			HasAcademicAffinity:       doesCrootHaveAffinity("Academics", croot),
			HasCloseToHomeAffinity:    doesCrootHaveAffinity("Close to Home", croot),
			HasServiceAffinity:        doesCrootHaveAffinity("Service", croot),
			HasFrontRunnerAffinity:    doesCrootHaveAffinity("Frontrunner", croot),
			HasReligionAffinity:       doesCrootHaveAffinity("Religion", croot),
			HasLargeCrowdAffinity:     doesCrootHaveAffinity("Large Crowds", croot),
			HasSmallSchoolAffinity:    doesCrootHaveAffinity("Small School", croot),
			HasSmallTownAffinity:      doesCrootHaveAffinity("Small Town", croot),
			HasBigCityAffinity:        doesCrootHaveAffinity("Big City", croot),
			HasMediaSpotlightAffinity: doesCrootHaveAffinity("Media Spotlight", croot),
			HasRisingStars:            doesCrootHaveAffinity("Rising Stars", croot),
		}
		recruitInfos[croot.ID] = info
		crootMap[croot.ID] = croot
	}

	teamProfiles := GetAllCollegeTeamsWithRecruitingProfileAndCoach()

	for _, t := range teamProfiles {
		tp := t.RecruitingProfile

		profileID := strconv.Itoa(int(t.ID))
		profiles := GetOnlyRecruitProfilesByTeamProfileID(profileID)

		for _, p := range profiles {
			croot := crootMap[uint(p.RecruitID)]
			closeToHome := util.IsCrootCloseToHome(croot.State, croot.City, t.State, t.TeamAbbr, stateMatcher, regionMatcher)
			oddsObject := getRecruitingOdds(ts, croot, tp, closeToHome, recruitInfos)

			triggerSave := false
			if oddsObject.Af1 && !p.AffinityOneEligible {
				p.ToggleAffinityOne()
				triggerSave = true
			}

			if oddsObject.Af2 && !p.AffinityTwoEligible {
				p.ToggleAffinityTwo()
				triggerSave = true
			}

			crootInfo := recruitInfos[uint(p.RecruitID)]

			for _, affinity := range tp.Affinities {
				if !affinity.IsApplicable {
					continue
				}
				if crootInfo.HasAcademicAffinity && isAffinityApplicable("Academics", affinity) {
					triggerSave = true
					if croot.AffinityOne == "Academics" && !p.AffinityOneEligible {
						p.ToggleAffinityOne()
					} else if croot.AffinityTwo == "Academics" && !p.AffinityTwoEligible {
						p.ToggleAffinityTwo()
					}
				}

				if crootInfo.HasFrontRunnerAffinity && isAffinityApplicable("Frontrunner", affinity) {
					triggerSave = true
					if croot.AffinityOne == "Frontrunner" && !p.AffinityOneEligible {
						p.ToggleAffinityOne()
					} else if croot.AffinityTwo == "Frontrunner" && !p.AffinityTwoEligible {
						p.ToggleAffinityTwo()
					}
				}

				if crootInfo.HasReligionAffinity && isAffinityApplicable("Religion", affinity) {
					triggerSave = true
					if croot.AffinityOne == "Religion" && !p.AffinityOneEligible {
						p.ToggleAffinityOne()
					} else if croot.AffinityTwo == "Religion" && !p.AffinityTwoEligible {
						p.ToggleAffinityTwo()
					}
				}

				if crootInfo.HasServiceAffinity && isAffinityApplicable("Service", affinity) {
					triggerSave = true
					if croot.AffinityOne == "Service" && !p.AffinityOneEligible {
						p.ToggleAffinityOne()
					} else if croot.AffinityTwo == "Service" && !p.AffinityTwoEligible {
						p.ToggleAffinityTwo()
					}
				}

				if crootInfo.HasSmallSchoolAffinity && isAffinityApplicable("Small School", affinity) {
					triggerSave = true
					if croot.AffinityOne == "Small School" && !p.AffinityOneEligible {
						p.ToggleAffinityOne()
					} else if croot.AffinityTwo == "Small School" && !p.AffinityTwoEligible {
						p.ToggleAffinityTwo()
					}
				}

				if crootInfo.HasSmallTownAffinity && isAffinityApplicable("Small Town", affinity) {
					triggerSave = true
					if croot.AffinityOne == "Small Town" && !p.AffinityOneEligible {
						p.ToggleAffinityOne()
					} else if croot.AffinityTwo == "Small Town" && !p.AffinityTwoEligible {
						p.ToggleAffinityTwo()
					}
				}

				if crootInfo.HasBigCityAffinity && isAffinityApplicable("Big City", affinity) {
					triggerSave = true
					if croot.AffinityOne == "Big City" && !p.AffinityOneEligible {
						p.ToggleAffinityOne()
					} else if croot.AffinityTwo == "Big City" && !p.AffinityTwoEligible {
						p.ToggleAffinityTwo()
					}
				}

				if crootInfo.HasMediaSpotlightAffinity && isAffinityApplicable("Media Spotlight", affinity) {
					triggerSave = true
					if croot.AffinityOne == "Media Spotlight" && !p.AffinityOneEligible {
						p.ToggleAffinityOne()
					} else if croot.AffinityTwo == "Media Spotlight" && !p.AffinityTwoEligible {
						p.ToggleAffinityTwo()
					}
				}

				if crootInfo.HasLargeCrowdAffinity && isAffinityApplicable("Large Crowds", affinity) {
					triggerSave = true
					if croot.AffinityOne == "Large Crowds" && !p.AffinityOneEligible {
						p.ToggleAffinityOne()
					} else if croot.AffinityTwo == "Large Crowds" && !p.AffinityTwoEligible {
						p.ToggleAffinityTwo()
					}
				}

				if crootInfo.HasRisingStars && isAffinityApplicable("Rising Stars", affinity) {
					triggerSave = true
					if croot.AffinityOne == "Rising Stars" && !p.AffinityOneEligible {
						p.ToggleAffinityOne()
					} else if croot.AffinityTwo == "Rising Stars" && !p.AffinityTwoEligible {
						p.ToggleAffinityTwo()
					}
				}
			}

			if triggerSave {
				repository.SaveRecruitProfile(p, db)
			}
		}
	}
}

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

func GetRecruitingProfileForTeamBoardByTeamID(TeamID string) structs.SimTeamBoardResponse {
	db := dbprovider.GetInstance().GetDB()

	var profile structs.RecruitingTeamProfile

	err := db.Preload("Affinities").Preload("Recruits.Recruit.RecruitPlayerProfiles", func(db *gorm.DB) *gorm.DB {
		return db.Order("total_points DESC")
	}).Where("id = ?", TeamID).Find(&profile).Error
	if err != nil {
		log.Panicln(err)
	}

	var teamProfileResponse structs.SimTeamBoardResponse
	var crootProfiles []structs.CrootProfile

	// iterate through player recruit profiles --> get recruit with preload to player profiles
	for i := 0; i < len(profile.Recruits); i++ {
		var crootProfile structs.CrootProfile
		var croot structs.Croot

		croot.Map(profile.Recruits[i].Recruit)

		crootProfile.Map(profile.Recruits[i], croot)

		crootProfiles = append(crootProfiles, crootProfile)
	}

	sort.Sort(structs.ByCrootProfileTotal(crootProfiles))

	teamProfileResponse.Map(profile, crootProfiles)

	return teamProfileResponse
}

func GetAllTeamRecruitingProfiles() []structs.RecruitingTeamProfile {
	db := dbprovider.GetInstance().GetDB()

	var profiles []structs.RecruitingTeamProfile

	err := db.Preload("Affinities").Find(&profiles).Error
	if err != nil {
		log.Panicln(err)
	}

	return profiles
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

func GetRecruitingClassByTeamID(TeamID string) structs.SimTeamBoardResponse {
	db := dbprovider.GetInstance().GetDB()

	var profile structs.RecruitingTeamProfile

	err := db.Preload("Recruits.Recruit", func(db *gorm.DB) *gorm.DB {
		return db.Order("overall DESC").Where("team_id = ? AND is_signed = true", TeamID)
	}).Where("id = ?", TeamID).Find(&profile).Error
	if err != nil {
		log.Panicln(err)
	}

	var teamProfileResponse structs.SimTeamBoardResponse
	var crootProfiles []structs.CrootProfile

	// iterate through player recruit profiles --> get recruit with preload to player profiles
	for i := 0; i < len(profile.Recruits); i++ {
		if profile.Recruits[i].IsSigned && profile.Recruits[i].Recruit.College == profile.TeamAbbreviation {
			var crootProfile structs.CrootProfile

			var croot structs.Croot

			croot.Map(profile.Recruits[i].Recruit)

			crootProfile.Map(profile.Recruits[i], croot)

			crootProfiles = append(crootProfiles, crootProfile)
		}

	}

	sort.Sort(structs.ByCrootProfileTotal(crootProfiles))

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
	// The roster is in dire shape, just open up the team needs
	if len(teamRoster) <= 60 {
		needsMap["QB"] = 2
		needsMap["RB"] = 2
		needsMap["FB"] = 1
		needsMap["TE"] = 2
		needsMap["WR"] = 3
		needsMap["OT"] = 3
		needsMap["OG"] = 3
		needsMap["C"] = 2
		needsMap["DE"] = 2
		needsMap["DT"] = 2
		needsMap["OLB"] = 2
		needsMap["ILB"] = 2
		needsMap["CB"] = 2
		needsMap["FS"] = 2
		needsMap["SS"] = 2
		needsMap["P"] = 1
		needsMap["K"] = 1
		needsMap["ATH"] = 3
		return needsMap
	}

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
	repository.SaveRecruitingTeamProfile(recruitingProfile, db)
}

func ResetTeamRecruitingProfiles() {
	db := dbprovider.GetInstance().GetDB()

	teamProfiles := GetAllTeamRecruitingProfiles()

	for _, t := range teamProfiles {
		t.SeasonReset()
		repository.SaveRecruitingTeamProfile(t, db)
	}
}
