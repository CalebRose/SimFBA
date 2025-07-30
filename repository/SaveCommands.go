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

func SaveNFLGameplanRecord(gameRecord structs.NFLGameplan, db *gorm.DB) {
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

func SaveCFBTeam(team structs.CollegeTeam, db *gorm.DB) {
	team.TeamSeasonStats = structs.CollegeTeamSeasonStats{}
	team.CollegeCoach = structs.CollegeCoach{}
	team.RecruitingProfile = structs.RecruitingTeamProfile{}
	team.TeamRecord = structs.CollegeTeamRecords{}
	team.TeamGameplan = structs.CollegeGameplan{}
	team.TeamDepthChart = structs.CollegeTeamDepthChart{}
	team.TeamStandings = nil
	team.TeamStats = nil
	err := db.Save(&team).Error
	if err != nil {
		log.Panicln("Could not save player record")
	}
}

func SaveNFLTeam(team structs.NFLTeam, db *gorm.DB) {
	team.TeamSeasonStats = nil
	team.Capsheet = structs.NFLCapsheet{}
	team.Contracts = nil
	team.DraftPicks = nil
	team.TeamGameplan = structs.NFLGameplan{}
	team.TeamDepthChart = structs.NFLDepthChart{}
	team.Standings = nil
	team.TeamStats = nil
	err := db.Save(&team).Error
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
	profile.Affinities = nil
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

func SaveRecruitRecord(croot structs.Recruit, db *gorm.DB) {
	croot.RecruitPlayerProfiles = nil
	err := db.Save(&croot).Error
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

func SaveTransferPortalProfile(profile structs.TransferPortalProfile, db *gorm.DB) {

	err := db.Save(&profile).Error
	if err != nil {
		log.Panicln("Could not save player record")
	}
}

func SaveCollegePlayerRecord(player structs.CollegePlayer, db *gorm.DB) {
	player.Stats = nil
	player.SeasonStats = structs.CollegePlayerSeasonStats{}
	err := db.Save(&player).Error
	if err != nil {
		log.Panicln("Could not save player record")
	}
}

func SaveNFLDepthChartPosition(dcp structs.NFLDepthChartPosition, db *gorm.DB) {
	dcp.NFLPlayer = structs.NFLPlayer{}
	err := db.Save(&dcp).Error
	if err != nil {
		log.Panicln("Could not save player record")
	}
}

func SaveCFBSeasonSnaps(snap structs.CollegePlayerSeasonSnaps, db *gorm.DB) {
	err := db.Save(&snap).Error
	if err != nil {
		log.Panicln("Could not save cfb season snaps record!")
	}
}

func SaveNFLSeasonSnaps(snap structs.NFLPlayerSeasonSnaps, db *gorm.DB) {
	err := db.Save(&snap).Error
	if err != nil {
		log.Panicln("Could not save nfl season snaps record!")
	}
}

func SaveNotification(noti structs.Notification, db *gorm.DB) {
	err := db.Save(&noti).Error
	if err != nil {
		log.Panicln("Could not save notification record!")
	}
}

func SaveCollegePromiseRecord(promise structs.CollegePromise, db *gorm.DB) {
	// Save College Player Record
	err := db.Save(&promise).Error
	if err != nil {
		log.Panicln("Could not save new college recruit record")
	}
}

func SaveCFBTeamStats(teamStat structs.CollegeTeamStats, db *gorm.DB) {
	// Save NFL Team Stats Record
	err := db.Save(&teamStat).Error
	if err != nil {
		log.Panicln("Could not save nfl team seasons stats record")
	}
}

func SaveCFBTeamSeasonStats(seasonStats structs.CollegeTeamSeasonStats, db *gorm.DB) {
	// Save NFL Team Season Stats Record
	err := db.Save(&seasonStats).Error
	if err != nil {
		log.Panicln("Could not save nfl team seasons stats record")
	}
}

func SaveNFLTeamStats(teamStat structs.NFLTeamStats, db *gorm.DB) {
	// Save NFL Team Stats Record
	err := db.Save(&teamStat).Error
	if err != nil {
		log.Panicln("Could not save nfl team seasons stats record")
	}
}

func SaveNFLTeamSeasonStats(seasonStats structs.NFLTeamSeasonStats, db *gorm.DB) {
	// Save NFL Team Season Stats Record
	err := db.Save(&seasonStats).Error
	if err != nil {
		log.Panicln("Could not save nfl team seasons stats record")
	}
}

func SaveNFLPlayerSeasonStats(seasonStats structs.NFLPlayerSeasonStats, db *gorm.DB) {
	// Save NFL Player Season Stats Record
	err := db.Save(&seasonStats).Error
	if err != nil {
		log.Panicln("Could not save nfl player seasons stats record")
	}
}

func SaveCollegePlayerSeasonStats(seasonStats structs.CollegePlayerSeasonStats, db *gorm.DB) {
	// Save CFB Player Season Stats Record
	err := db.Save(&seasonStats).Error
	if err != nil {
		log.Panicln("Could not save cfb player seasons stats record")
	}
}

func SaveCFBStandingsRecord(record structs.CollegeStandings, db *gorm.DB) {
	err := db.Save(&record).Error
	if err != nil {
		log.Panicln("Could not save Standings Record for " + record.TeamName)
	}
}

func SaveNFLStandingsRecord(record structs.NFLStandings, db *gorm.DB) {
	err := db.Save(&record).Error
	if err != nil {
		log.Panicln("Could not save Standings Record for " + record.TeamName)
	}
}
