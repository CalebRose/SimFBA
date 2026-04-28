package managers

import (
	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
)

func FixPreferencesForAllRecruitsAndCollegePlayers() {
	db := dbprovider.GetInstance().GetDB()

	recruits := GetAllRecruitRecords()
	collegePlayers := GetAllCollegePlayers()

	for _, recruit := range recruits {

		// Preferences
		program := util.GenerateNormalizedIntFromRange(1, 9)
		profDevelopment := util.GenerateNormalizedIntFromRange(1, 9)
		traditions := util.GenerateNormalizedIntFromRange(1, 9)
		facilities := util.GenerateNormalizedIntFromRange(1, 9)
		atmosphere := util.GenerateNormalizedIntFromRange(1, 9)
		academics := util.GenerateNormalizedIntFromRange(1, 9)
		conferencePrestige := util.GenerateNormalizedIntFromRange(1, 9)
		coachPref := util.GenerateNormalizedIntFromRange(1, 9)
		seasonMomentumPref := util.GenerateNormalizedIntFromRange(1, 9)
		campusLife := util.GenerateNormalizedIntFromRange(1, 9)
		religionPref := util.GenerateNormalizedIntFromRange(1, 9)
		serviceAcademyPref := util.GenerateNormalizedIntFromRange(1, 9)
		smallTownPref := util.GenerateNormalizedIntFromRange(1, 9)
		bigCityPref := util.GenerateNormalizedIntFromRange(1, 9)
		mediaSpotlightPref := util.GenerateNormalizedIntFromRange(1, 9)

		playerPref := structs.PlayerPreferences{
			ProgramPref:        uint8(program),
			ProfDevPref:        uint8(profDevelopment),
			TraditionsPref:     uint8(traditions),
			FacilitiesPref:     uint8(facilities),
			AtmospherePref:     uint8(atmosphere),
			AcademicsPref:      uint8(academics),
			ConferencePref:     uint8(conferencePrestige),
			CoachPref:          uint8(coachPref),
			SeasonMomentumPref: uint8(seasonMomentumPref),
			CampusLifePref:     uint8(campusLife),
			ReligionPref:       uint8(religionPref),
			ServiceAcademyPref: uint8(serviceAcademyPref),
			SmallTownPref:      uint8(smallTownPref),
			BigCityPref:        uint8(bigCityPref),
			MediaSpotlightPref: uint8(mediaSpotlightPref),
		}

		recruit.AssignPreferences(playerPref)
		repository.SaveRecruitRecord(recruit, db)
	}

	for _, player := range collegePlayers {

		// Preferences
		program := util.GenerateNormalizedIntFromRange(1, 9)
		profDevelopment := util.GenerateNormalizedIntFromRange(1, 9)
		traditions := util.GenerateNormalizedIntFromRange(1, 9)
		facilities := util.GenerateNormalizedIntFromRange(1, 9)
		atmosphere := util.GenerateNormalizedIntFromRange(1, 9)
		academics := util.GenerateNormalizedIntFromRange(1, 9)
		conferencePrestige := util.GenerateNormalizedIntFromRange(1, 9)
		coachPref := util.GenerateNormalizedIntFromRange(1, 9)
		seasonMomentumPref := util.GenerateNormalizedIntFromRange(1, 9)
		campusLife := util.GenerateNormalizedIntFromRange(1, 9)
		religionPref := util.GenerateNormalizedIntFromRange(1, 9)
		serviceAcademyPref := util.GenerateNormalizedIntFromRange(1, 9)
		smallTownPref := util.GenerateNormalizedIntFromRange(1, 9)
		bigCityPref := util.GenerateNormalizedIntFromRange(1, 9)
		mediaSpotlightPref := util.GenerateNormalizedIntFromRange(1, 9)

		playerPref := structs.PlayerPreferences{
			ProgramPref:        uint8(program),
			ProfDevPref:        uint8(profDevelopment),
			TraditionsPref:     uint8(traditions),
			FacilitiesPref:     uint8(facilities),
			AtmospherePref:     uint8(atmosphere),
			AcademicsPref:      uint8(academics),
			ConferencePref:     uint8(conferencePrestige),
			CoachPref:          uint8(coachPref),
			SeasonMomentumPref: uint8(seasonMomentumPref),
			CampusLifePref:     uint8(campusLife),
			ReligionPref:       uint8(religionPref),
			ServiceAcademyPref: uint8(serviceAcademyPref),
			SmallTownPref:      uint8(smallTownPref),
			BigCityPref:        uint8(bigCityPref),
			MediaSpotlightPref: uint8(mediaSpotlightPref),
		}

		player.AssignPreferences(playerPref)
		repository.SaveCollegePlayerRecord(player, db)
	}
}
