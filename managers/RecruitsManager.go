package managers

import (
	"errors"
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

func GetAllRecruits() []structs.Croot {
	db := dbprovider.GetInstance().GetDB()

	var recruits []structs.Recruit

	db.Preload("RecruitPlayerProfiles", func(db *gorm.DB) *gorm.DB {
		return db.Order("total_points DESC")
	}).Find(&recruits)

	var croots []structs.Croot
	for _, recruit := range recruits {
		var croot structs.Croot
		croot.Map(recruit)

		croots = append(croots, croot)
	}

	sort.Sort(structs.ByCrootRank(croots))

	return croots
}

func GetAllRecruitRecords() []structs.Recruit {
	db := dbprovider.GetInstance().GetDB()

	var croots []structs.Recruit

	db.Find(&croots)

	return croots
}

func GetAllUnsignedRecruits() []structs.Recruit {
	db := dbprovider.GetInstance().GetDB()

	var croots []structs.Recruit

	db.Where("is_signed = ?", false).Find(&croots)

	return croots
}

func GetCollegeRecruitByRecruitID(recruitID string) structs.Recruit {
	db := dbprovider.GetInstance().GetDB()

	var recruit structs.Recruit

	err := db.Where("id = ?", recruitID).Find(&recruit).Error
	if err != nil {
		log.Fatalln(err)
	}

	return recruit
}

func GetCollegeRecruitViaDiscord(id string) structs.Croot {
	db := dbprovider.GetInstance().GetDB()

	var recruit structs.Recruit

	err := db.Preload("RecruitPlayerProfiles").Where("id = ?", id).Find(&recruit).Error
	if err != nil {
		log.Fatalln(err)
	}

	var croot structs.Croot

	croot.Map(recruit)

	return croot
}

func GetCollegeRecruitByRecruitIDForTeamBoard(recruitID string) structs.Recruit {
	db := dbprovider.GetInstance().GetDB()

	var recruit structs.Recruit

	err := db.Preload("RecruitPlayerProfiles", func(db *gorm.DB) *gorm.DB {
		return db.Order("total_points DESC").Where("total_points > ?", "0")
	}).Where("id = ?", recruitID).Find(&recruit).Error
	if err != nil {
		log.Fatalln(err)
	}

	return recruit
}

func GetRecruitsByTeamProfileID(ProfileID string) []structs.RecruitPlayerProfile {
	db := dbprovider.GetInstance().GetDB()

	var croots []structs.RecruitPlayerProfile

	err := db.Preload("Recruit").Where("profile_id = ?", ProfileID).Find(&croots).Error
	if err != nil {
		log.Fatal(err)
	}

	return croots
}

func GetRecruitsForAIPointSync(ProfileID string) []structs.RecruitPlayerProfile {
	db := dbprovider.GetInstance().GetDB()

	var croots []structs.RecruitPlayerProfile

	err := db.Preload("Recruit", func(db *gorm.DB) *gorm.DB {
		return db.Order("stars DESC")
	}).Where("profile_id = ? AND removed_from_board = ?", ProfileID, false).Order("total_points DESC").Find(&croots).Error
	if err != nil {
		log.Fatal(err)
	}

	return croots
}

func GetOnlyRecruitProfilesByTeamProfileID(ProfileID string) []structs.RecruitPlayerProfile {
	db := dbprovider.GetInstance().GetDB()

	var croots []structs.RecruitPlayerProfile

	err := db.Where("profile_id = ?", ProfileID).Find(&croots).Error

	if err != nil {
		log.Fatal(err)
	}

	return croots
}

func GetSignedRecruitsByTeamProfileID(ProfileID string) []structs.Recruit {
	db := dbprovider.GetInstance().GetDB()

	var croots []structs.Recruit

	err := db.Order("overall DESC").Where("team_id = ? AND is_signed = ?", ProfileID, true).Find(&croots).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []structs.Recruit{}
		} else {
			log.Fatal(err)
		}
	}

	return croots
}

func GetRecruitProfileByPlayerId(recruitID string, profileID string) structs.RecruitPlayerProfile {
	db := dbprovider.GetInstance().GetDB()

	var croot structs.RecruitPlayerProfile
	err := db.Where("recruit_id = ? and profile_id = ?", recruitID, profileID).Find(&croot).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return structs.RecruitPlayerProfile{}
		} else {
			log.Fatal(err)
		}
	}
	return croot
}

func GetRecruitPlayerProfilesByRecruitId(recruitID string) []structs.RecruitPlayerProfile {
	db := dbprovider.GetInstance().GetDB()

	var croots []structs.RecruitPlayerProfile
	err := db.Where("recruit_id = ?", recruitID).Order("total_points desc").Find(&croots).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []structs.RecruitPlayerProfile{}
		} else {
			log.Fatal(err)
		}
	}
	return croots
}

func CreateRecruitingProfileForRecruit(recruitPointsDto structs.CreateRecruitProfileDto) structs.RecruitPlayerProfile {
	db := dbprovider.GetInstance().GetDB()

	recruitEntry := GetRecruitProfileByPlayerId(strconv.Itoa(recruitPointsDto.RecruitID),
		strconv.Itoa(recruitPointsDto.ProfileID))

	teamRecruitingProfile := GetOnlyRecruitingProfileByTeamID(strconv.Itoa(int(recruitPointsDto.ProfileID)))

	if recruitEntry.RecruitID != 0 && recruitEntry.ProfileID != 0 {
		// Replace Recruit
		recruitEntry.ToggleRemoveFromBoard()
		repository.SaveRecruitProfile(recruitEntry, db)
		return recruitEntry
	}
	stateMatcher := util.GetStateMatcher()
	regionMatcher := util.GetStateRegionMatcher()

	modifier := CalculateModifierTowardsRecruit(recruitPointsDto.PlayerRecruit, teamRecruitingProfile, stateMatcher, regionMatcher)

	createRecruitEntry := structs.RecruitPlayerProfile{
		SeasonID:            recruitPointsDto.SeasonID,
		RecruitID:           recruitPointsDto.RecruitID,
		ProfileID:           recruitPointsDto.ProfileID,
		TotalPoints:         0,
		CurrentWeeksPoints:  0,
		SpendingCount:       0,
		Scholarship:         false,
		ScholarshipRevoked:  false,
		AffinityOneEligible: recruitPointsDto.AffinityOneEligible,
		AffinityTwoEligible: recruitPointsDto.AffinityTwoEligible,
		TeamAbbreviation:    recruitPointsDto.Team,
		RemovedFromBoard:    false,
		IsSigned:            false,
		PreferenceModifier:  modifier,
	}

	// Create
	repository.CreateRecruitProfileRecord(createRecruitEntry, db)
	return createRecruitEntry
}

func SendScholarshipToRecruit(updateRecruitPointsDto structs.UpdateRecruitPointsDto) (structs.RecruitPlayerProfile, structs.RecruitingTeamProfile) {
	db := dbprovider.GetInstance().GetDB()

	recruitingProfile := GetOnlyRecruitingProfileByTeamID(strconv.Itoa(updateRecruitPointsDto.ProfileID))

	if recruitingProfile.ScholarshipsAvailable == 0 && (updateRecruitPointsDto.RewardScholarship || updateRecruitPointsDto.RevokeScholarship) {
		log.Panicln("\nTeamId: " + strconv.Itoa(updateRecruitPointsDto.ProfileID) + " does not have any availabe scholarships")
	}

	crootProfile := GetRecruitProfileByPlayerId(
		strconv.Itoa(updateRecruitPointsDto.RecruitID),
		strconv.Itoa(updateRecruitPointsDto.ProfileID),
	)

	crootProfile.ToggleScholarship()
	if !crootProfile.ScholarshipRevoked {
		recruitingProfile.SubtractScholarshipsAvailable()
	} else {
		recruitingProfile.ReallocateScholarship()
	}

	repository.SaveRecruitProfile(crootProfile, db)
	repository.SaveRecruitingTeamProfile(recruitingProfile, db)

	return crootProfile, recruitingProfile
}

func RevokeScholarshipFromRecruit(updateRecruitPointsDto structs.UpdateRecruitPointsDto) (structs.RecruitPlayerProfile, structs.RecruitingTeamProfile) {
	db := dbprovider.GetInstance().GetDB()

	recruitingProfile := GetOnlyRecruitingProfileByTeamID(strconv.Itoa(updateRecruitPointsDto.ProfileID))

	recruitingPointsProfile := GetRecruitProfileByPlayerId(
		strconv.Itoa(updateRecruitPointsDto.RecruitID),
		strconv.Itoa(updateRecruitPointsDto.ProfileID),
	)

	if !recruitingPointsProfile.Scholarship {
		fmt.Printf("%s", "\nCannot revoke an inexistant scholarship from Recruit "+strconv.Itoa(recruitingPointsProfile.RecruitID))
		return recruitingPointsProfile, recruitingProfile
	}

	// recruitingPointsProfile.ToggleScholarship()
	recruitingProfile.ReallocateScholarship()

	repository.SaveRecruitProfile(recruitingPointsProfile, db)
	repository.SaveRecruitingTeamProfile(recruitingProfile, db)

	return recruitingPointsProfile, recruitingProfile
}

func GetRecruitFromRecruitsList(id int, recruits []structs.RecruitPlayerProfile) structs.RecruitPlayerProfile {
	var recruit structs.RecruitPlayerProfile

	for i := 0; i < len(recruits); i++ {
		if recruits[i].RecruitID == id {
			recruit = recruits[i]
			break
		}
	}

	return recruit
}

func CreateCollegeRecruit(createRecruitDTO structs.CreateRecruitDTO) {
	db := dbprovider.GetInstance().GetDB()

	var lastPlayerRecord structs.Player

	err := db.Last(&lastPlayerRecord).Error
	if err != nil {
		log.Fatalln("Could not grab last player record from players table...")
	}

	newID := lastPlayerRecord.ID + 1

	collegeRecruit := &structs.Recruit{}
	collegeRecruit.Map(createRecruitDTO, newID)

	// No Player Record exists, so we shall make one.

	db.Create(&collegeRecruit)

	playerRecord := structs.Player{
		RecruitID: int(collegeRecruit.ID),
	}
	// Create Player Record
	db.Create(&playerRecord)
	// Assign PlayerID to Recruit
	collegeRecruit.AssignPlayerID(int(playerRecord.ID))
	// Save Recruit
	db.Save(&collegeRecruit)
}

func CalculateModifierTowardsRecruit(player structs.Recruit, team structs.RecruitingTeamProfile, stateMatcher map[string][]string, regionMatcher map[string]map[string]string) float32 {
	prefs := player.PlayerPreferences
	closeToHome := util.IsCrootCloseToHome(player.State, player.City, team.State, team.TeamAbbreviation, stateMatcher, regionMatcher)
	programMod := calculateMultiplier(uint(team.ProgramPrestige), uint(prefs.ProgramPref), false, false)
	professionalDevMod := calculateMultiplier(uint(team.ProfessionalPrestige), uint(prefs.ProfDevPref), false, false)
	traditionsMod := calculateMultiplier(uint(team.Traditions), uint(prefs.TraditionsPref), false, false)
	facilitiesMod := calculateMultiplier(uint(team.Facilities), uint(prefs.FacilitiesPref), false, false)
	atmosphereMod := calculateMultiplier(uint(team.Atmosphere), uint(prefs.AtmospherePref), false, false)
	academicsMod := calculateMultiplier(uint(team.Academics), uint(prefs.AcademicsPref), false, false)
	conferenceMod := calculateMultiplier(uint(team.ConferencePrestige), uint(prefs.ConferencePref), false, false)
	coachMod := calculateMultiplier(uint(team.CoachRating), uint(prefs.CoachPref), false, false)
	seasonMod := calculateMultiplier(uint(team.SeasonMomentum), uint(prefs.SeasonMomentumPref), false, false)
	collegeLifeMod := calculateMultiplier(uint(team.CampusLife), uint(prefs.CampusLifePref), false, false)
	mediaSpotlightMod := calculateMultiplier(uint(team.MediaSpotlight), uint(prefs.MediaSpotlightPref), false, false)
	religionMod := calculateMultiplier(1, uint(prefs.ReligionPref), true, team.ReligionAffinity)
	serviceMod := calculateMultiplier(1, uint(prefs.ServiceAcademyPref), true, team.ServiceAffinity)
	smallTownMod := calculateMultiplier(1, uint(prefs.SmallTownPref), true, team.SmallTownAffinity)
	bigCityMod := calculateMultiplier(1, uint(prefs.BigCityPref), true, team.BigCityAffinity)

	dynamicModSum := programMod + professionalDevMod + traditionsMod + atmosphereMod + conferenceMod + coachMod + seasonMod + mediaSpotlightMod
	staticMod := facilitiesMod + academicsMod + religionMod + serviceMod + smallTownMod + bigCityMod + collegeLifeMod

	closeToHomeMod := float32(0.0)
	if closeToHome {
		closeToHomeMod = 0.2
	}
	// Weighted average of dynamic and static modifiers, with dynamic modifiers having a higher weight
	// since they are more likely to change over time and thus be more influential in a
	// recruit's decision making process
	largeMod := (dynamicModSum*1.7 + staticMod*0.3) / 15
	return largeMod + closeToHomeMod
}

func calculateBaseModifier(attr int, isBool, booleanAttr bool) float32 {
	attrVal := attr
	if isBool {
		if booleanAttr {
			attrVal = 10
		} else {
			attrVal = 5
		}
	}
	return 1 + float32(attrVal-5)/5
}

func calculateAdjustmentFactor(teamAttr, playerPref int) float32 {
	return 1 + float32((teamAttr-playerPref)/10)
}

func calculateMultiplier(teamAttr uint, playerPref uint, isBool, booleanAttr bool) float32 {
	baseMod := calculateBaseModifier(int(teamAttr), isBool, booleanAttr)
	adjFactor := calculateAdjustmentFactor(int(teamAttr), int(playerPref))
	return baseMod * adjFactor
}
