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

func FixExistingModifiersForRecruits() {
	db := dbprovider.GetInstance().GetDB()
	stateMatcher := util.GetStateMatcher()
	regionMatcher := util.GetStateRegionMatcher()
	recruits := GetAllRecruitRecords()
	recruitMap := MakeCollegeRecruitMapByID(recruits)
	recruitProfiles := repository.FindRecruitPlayerProfileRecords("", "", false, false, false)
	teamProfiles := GetAllTeamRecruitingProfiles()
	teamProfileMap := MakeRecruitTeamProfileMapByTeamID(teamProfiles)

	for _, profile := range recruitProfiles {
		recruit := recruitMap[uint(profile.RecruitID)]
		teamProfile := teamProfileMap[uint(profile.ProfileID)]

		modifier := CalculateModifierTowardsRecruit(recruit, teamProfile, stateMatcher, regionMatcher)

		profile.PreferenceModifier = modifier
		repository.SaveRecruitProfile(profile, db)
	}
}

func FixPlayerWeights() {
	db := dbprovider.GetInstance().GetDB()
	collegePlayers := GetAllCollegePlayers()
	recruits := GetAllRecruitRecords()
	nflPlayers := GetAllNFLPlayers()
	recruitMap := MakeCollegeRecruitMapByID(recruits)
	collegePlayerMap := MakeCollegePlayerMap(collegePlayers)
	nflPlayerMap := MakeNFLPlayerMap(nflPlayers)

	collegePlayerPath := "data/2027/weight_fix/cfb_players_backup.csv"
	collegePlayersCSV := util.ReadCSV(collegePlayerPath)
	recruitPath := "data/2027/weight_fix/2027Croots_with_weight.csv"
	recruitsCSV := util.ReadCSV(recruitPath)
	nflPlayerPath := "data/2027/weight_fix/simnfl_players.csv"
	nflPlayersCSV := util.ReadCSV(nflPlayerPath)
	draftedPlayerPath := "data/2027/weight_fix/2027_nfl_draftees.csv"
	draftedPlayersCSV := util.ReadCSV(draftedPlayerPath)

	for idx, row := range collegePlayersCSV {
		if idx == 0 {
			continue
		}
		weight := util.ConvertStringToInt(row[4])
		id := util.ConvertStringToInt(row[0])
		player := collegePlayerMap[uint(id)]
		player.Weight = int16(weight)
		repository.SaveCollegePlayerRecord(player, db)
	}

	for idx, row := range recruitsCSV {
		if idx == 0 {
			continue
		}
		weight := util.ConvertStringToInt(row[3])
		id := util.ConvertStringToInt(row[0])
		player := recruitMap[uint(id)]
		player.Weight = int16(weight)
		repository.SaveRecruitRecord(player, db)
	}

	for idx, row := range nflPlayersCSV {
		if idx == 0 {
			continue
		}
		weight := util.ConvertStringToInt(row[4])
		id := util.ConvertStringToInt(row[0])
		player := nflPlayerMap[uint(id)]
		player.Weight = int16(weight)
		repository.SaveNFLPlayerRecord(player, db)
	}

	for idx, row := range draftedPlayersCSV {
		if idx == 0 {
			continue
		}
		weight := util.ConvertStringToInt(row[4])
		id := util.ConvertStringToInt(row[0])
		player := nflPlayerMap[uint(id)]
		player.Weight = int16(weight)
		repository.SaveNFLPlayerRecord(player, db)
	}
}
