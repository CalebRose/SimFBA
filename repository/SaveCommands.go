package repository

import (
	"log"
	"strconv"

	"github.com/CalebRose/SimFBA/structs"
	"gorm.io/gorm"
)

func SaveTimestamp(ts structs.Timestamp, db *gorm.DB) {
	err := db.Save(&ts).Error
	if err != nil {
		log.Panicln("Could not save timestamp")
	}
}

func SaveCFBGameplanRecord(gameRecord structs.CollegeGameplan, db *gorm.DB) {
	err := db.Save(&gameRecord).Error
	if err != nil {
		log.Panicln("Could not save Gameplan " + strconv.Itoa(int(gameRecord.ID)))
	}
}

func SaveCFBGameRecord(gameRecord structs.CollegeGame, db *gorm.DB) {
	err := db.Save(&gameRecord).Error
	if err != nil {
		log.Panicln("Could not save Game " + strconv.Itoa(int(gameRecord.ID)) + "Between " + gameRecord.HomeTeam + " and " + gameRecord.AwayTeam)
	}
}

func SaveNFLGameRecord(gameRecord structs.NFLGame, db *gorm.DB) {
	err := db.Save(&gameRecord).Error
	if err != nil {
		log.Panicln("Could not save Game " + strconv.Itoa(int(gameRecord.ID)) + "Between " + gameRecord.HomeTeam + " and " + gameRecord.AwayTeam)
	}
}

func SaveCFBPlayer(player structs.CollegePlayer, db *gorm.DB) {
	player.SeasonStats = structs.CollegePlayerSeasonStats{}
	player.Stats = nil
	err := db.Save(&player).Error
	if err != nil {
		log.Panicln("Could not save player record")
	}
}

func SaveNFLPlayer(player structs.NFLPlayer, db *gorm.DB) {
	player.SeasonStats = structs.NFLPlayerSeasonStats{}
	player.Stats = nil
	player.Offers = nil
	player.WaiverOffers = nil
	player.Extensions = nil
	player.Contract = structs.NFLContract{}
	err := db.Save(&player).Error
	if err != nil {
		log.Panicln("Could not save player record")
	}
}

func SaveNFLContract(c structs.NFLContract, db *gorm.DB) {
	err := db.Save(&c).Error
	if err != nil {
		log.Panicln("Could not save contract record")
	}
}

func SaveNFLCapsheet(c structs.NFLCapsheet, db *gorm.DB) {
	err := db.Save(&c).Error
	if err != nil {
		log.Panicln("Could not save capsheet record")
	}
}

func SaveRecruitingTeamProfile(profile structs.RecruitingTeamProfile, db *gorm.DB) {
	err := db.Save(&profile).Error
	if err != nil {
		log.Panicln("Could not save team profile")
	}
}

func SaveRecruitProfile(profile structs.RecruitPlayerProfile, db *gorm.DB) {
	profile.Recruit = structs.Recruit{}
	err := db.Save(&profile).Error
	if err != nil {
		log.Panicln("Could not save team profile")
	}
}

func SaveCollegeTeamRecord(team structs.CollegeTeam, db *gorm.DB) {
	team.CollegeCoach = structs.CollegeCoach{}
	team.RecruitingProfile = structs.RecruitingTeamProfile{}
	team.TeamStats = nil
	team.TeamStandings = nil
	team.TeamRecord = structs.CollegeTeamRecords{}
	team.TeamGameplan = structs.CollegeGameplan{}
	team.TeamDepthChart = structs.CollegeTeamDepthChart{}
	team.TeamSeasonStats = structs.CollegeTeamSeasonStats{}
	err := db.Save(&team).Error
	if err != nil {
		log.Panicln("Could not save team profile")
	}
}

func SaveNFLTeamRecord(team structs.NFLTeam, db *gorm.DB) {
	team.Capsheet = structs.NFLCapsheet{}
	team.Contracts = nil
	team.TeamStats = nil
	team.Standings = nil
	team.DraftPicks = nil
	team.TeamGameplan = structs.NFLGameplan{}
	team.TeamDepthChart = structs.NFLDepthChart{}
	team.TeamSeasonStats = nil
	err := db.Save(&team).Error
	if err != nil {
		log.Panicln("Could not save team profile")
	}
}
