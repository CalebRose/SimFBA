package managers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/structs"
)

func GetAllCollegeGameplans() []structs.CollegeGameplan {
	db := dbprovider.GetInstance().GetDB()

	gameplans := []structs.CollegeGameplan{}

	db.Find(&gameplans)

	return gameplans
}

func GetAllNFLGameplans() []structs.NFLGameplan {
	db := dbprovider.GetInstance().GetDB()

	gameplans := []structs.NFLGameplan{}

	db.Find(&gameplans)

	return gameplans
}

func UpdateGameplanPenalties() {
	db := dbprovider.GetInstance().GetDB()

	collegeGPs := GetAllCollegeGameplans()

	for _, gp := range collegeGPs {
		if gp.HasSchemePenalty {
			gp.LowerPenalty()
			db.Save(&gp)
		}
	}

	nflGPs := GetAllNFLGameplans()

	for _, gp := range nflGPs {
		if gp.HasSchemePenalty {
			gp.LowerPenalty()
			db.Save(&gp)
		}
	}
}

func GetGameplanByTeamID(teamID string) structs.CollegeGameplan {
	db := dbprovider.GetInstance().GetDB()

	var gamePlan structs.CollegeGameplan

	err := db.Where("id = ?", teamID).Find(&gamePlan).Error
	if err != nil {
		fmt.Println(err)
		log.Fatalln("Gameplan does not exist for team.")
	}
	return gamePlan
}

func GetNFLGameplanByTeamID(teamID string) structs.NFLGameplan {
	db := dbprovider.GetInstance().GetDB()

	var gamePlan structs.NFLGameplan

	err := db.Where("id = ?", teamID).Find(&gamePlan).Error
	if err != nil {
		fmt.Println(err)
		log.Fatalln("Gameplan does not exist for team.")
	}
	return gamePlan
}

func GetGameplanByGameplanID(gameplanID string) structs.CollegeGameplan {
	db := dbprovider.GetInstance().GetDB()

	var gamePlan structs.CollegeGameplan

	err := db.Where("id = ?", gameplanID).Find(&gamePlan).Error
	if err != nil {
		fmt.Println(err)
		log.Fatalln("Gameplan does not exist for team.")
	}
	return gamePlan
}

func GetDepthchartByTeamID(teamID string) structs.CollegeTeamDepthChart {
	db := dbprovider.GetInstance().GetDB()

	var depthChart structs.CollegeTeamDepthChart

	// Preload Depth Chart Positions
	err := db.Preload("DepthChartPlayers.CollegePlayer").Where("team_id = ?", teamID).Find(&depthChart).Error
	if err != nil {
		fmt.Println(err)
		log.Fatalln("Depthchart does not exist for team.")
	}
	return depthChart
}

func GetNFLDepthchartByTeamID(teamID string) structs.NFLDepthChart {
	db := dbprovider.GetInstance().GetDB()

	var depthChart structs.NFLDepthChart

	// Preload Depth Chart Positions
	err := db.Preload("DepthChartPlayers.NFLPlayer").Where("team_id = ?", teamID).Find(&depthChart).Error
	if err != nil {
		fmt.Println(err)
		log.Fatalln("Depthchart does not exist for team.")
	}
	return depthChart
}

func GetDepthChartPositionPlayersByDepthchartID(depthChartID string) []structs.CollegeDepthChartPosition {
	db := dbprovider.GetInstance().GetDB()

	var positionPlayers []structs.CollegeDepthChartPosition

	err := db.Where("depth_chart_id = ?", depthChartID).Find(&positionPlayers).Error
	if err != nil {
		fmt.Println(err)
		panic("Depth Chart does not exist for this ID")
	}

	return positionPlayers
}

func GetNFLDepthChartPositionsByDepthchartID(depthChartID string) []structs.NFLDepthChartPosition {
	db := dbprovider.GetInstance().GetDB()

	var positionPlayers []structs.NFLDepthChartPosition

	err := db.Where("depth_chart_id = ?", depthChartID).Find(&positionPlayers).Error
	if err != nil {
		fmt.Println(err)
		panic("Depth Chart does not exist for this ID")
	}

	return positionPlayers
}

func UpdateGameplan(updateGameplanDto structs.UpdateGameplanDTO) {
	db := dbprovider.GetInstance().GetDB()

	gameplanID := updateGameplanDto.GameplanID

	currentGameplan := GetGameplanByGameplanID(gameplanID)

	ts := GetTimestamp()

	schemePenalty := false

	if currentGameplan.OffensiveScheme != updateGameplanDto.UpdatedGameplan.OffensiveScheme {

		if ts.CollegeWeek != 0 {
			currentGameplan.ApplySchemePenalty(true)
		}
		schemePenalty = true
	}

	if currentGameplan.DefensiveScheme != updateGameplanDto.UpdatedGameplan.DefensiveScheme {

		if ts.CollegeWeek != 0 {
			currentGameplan.ApplySchemePenalty(false)
		}
		schemePenalty = true
	}

	if schemePenalty {

		newsLog := structs.NewsLog{
			WeekID:      ts.CollegeWeekID,
			Week:        ts.CollegeWeek,
			SeasonID:    ts.CollegeSeasonID,
			MessageType: "Gameplan",
			League:      "CFB",
			Message:     "Coach " + updateGameplanDto.Username + " has updated " + updateGameplanDto.TeamName + "'s offensive scheme from " + currentGameplan.OffensiveScheme + " to " + updateGameplanDto.UpdatedGameplan.OffensiveScheme,
		}

		db.Create(&newsLog)
	}

	currentGameplan.UpdateGameplan(updateGameplanDto.UpdatedGameplan)

	db.Save(&currentGameplan)
}

func UpdateNFLGameplan(updateGameplanDto structs.UpdateGameplanDTO) {
	db := dbprovider.GetInstance().GetDB()

	gameplanID := updateGameplanDto.GameplanID

	currentGameplan := GetNFLGameplanByTeamID(gameplanID)
	UpdatedGameplan := updateGameplanDto.UpdatedNFLGameplan

	schemeChange := false
	ts := GetTimestamp()
	if currentGameplan.OffensiveScheme != UpdatedGameplan.OffensiveScheme && !ts.IsNFLOffSeason {

		if ts.NFLWeek != 0 {
			currentGameplan.ApplySchemePenalty(true)
		}

		schemeChange = true

	}

	if currentGameplan.DefensiveScheme != UpdatedGameplan.DefensiveScheme && !ts.IsNFLOffSeason {

		if ts.NFLWeek != 0 {
			currentGameplan.ApplySchemePenalty(false)
		}
		schemeChange = true
	}

	if schemeChange {

		newsLog := structs.NewsLog{
			WeekID:      ts.NFLWeekID,
			Week:        ts.NFLWeek,
			SeasonID:    ts.NFLSeasonID,
			League:      "NFL",
			MessageType: "Gameplan",
			Message:     "Coach " + updateGameplanDto.Username + " has updated " + updateGameplanDto.TeamName + "'s offensive scheme from " + currentGameplan.OffensiveScheme + " to " + updateGameplanDto.UpdatedGameplan.OffensiveScheme,
		}

		db.Create(&newsLog)
	}

	currentGameplan.UpdateGameplan(UpdatedGameplan)

	db.Save(&currentGameplan)
}

func UpdateDepthChart(updateDepthchartDTO structs.UpdateDepthChartDTO) {

	depthChartID := strconv.Itoa(updateDepthchartDTO.DepthChartID)
	depthChartPlayers := GetDepthChartPositionPlayersByDepthchartID(depthChartID)

	updatedPlayers := updateDepthchartDTO.UpdatedPlayerPositions
	updateCounter := 0

	fmt.Println(len(depthChartPlayers))
	fmt.Println(len(updatedPlayers))
	db := dbprovider.GetInstance().GetDB()

	for i := 0; i < len(depthChartPlayers); i++ {
		player := depthChartPlayers[i]

		updatedPlayer := GetPlayerFromDClist(player.ID, updatedPlayers)

		if player.ID == updatedPlayer.ID &&
			player.PlayerID == updatedPlayer.PlayerID &&
			player.OriginalPosition == updatedPlayer.OriginalPosition {
			continue
		}

		player.UpdateDepthChartPosition(updatedPlayer)

		updateCounter++

		if updateCounter == len(updatedPlayers) {
			break
		}
		db.Save(&player)
	}
}

func UpdateNFLDepthChart(updateDepthchartDTO structs.UpdateNFLDepthChartDTO) {

	depthChartID := strconv.Itoa(updateDepthchartDTO.DepthChartID)
	depthChartPlayers := GetNFLDepthChartPositionsByDepthchartID(depthChartID)

	updatedPlayers := updateDepthchartDTO.UpdatedPlayerPositions
	updateCounter := 0

	db := dbprovider.GetInstance().GetDB()

	for i := 0; i < len(depthChartPlayers); i++ {
		player := depthChartPlayers[i]

		updatedPlayer := GetPlayerFromNFLDClist(player.ID, updatedPlayers)

		if player.ID == updatedPlayer.ID &&
			uint(player.PlayerID) == updatedPlayer.PlayerID &&
			player.OriginalPosition == updatedPlayer.OriginalPosition {
			continue
		}

		player.UpdateDepthChartPosition(updatedPlayer)

		updateCounter++

		if updateCounter == len(updatedPlayers) {
			break
		}
		db.Save(&player)
	}
}

func GetPlayerFromDClist(id uint, updatedPlayers []structs.CollegeDepthChartPosition) structs.CollegeDepthChartPosition {
	var player structs.CollegeDepthChartPosition

	for i := 0; i < len(updatedPlayers); i++ {
		if updatedPlayers[i].ID == id {
			player = updatedPlayers[i]
			break
		}
	}

	return player
}

func GetPlayerFromNFLDClist(id uint, updatedPlayers []structs.NFLDepthChartPosition) structs.NFLDepthChartPosition {
	var player structs.NFLDepthChartPosition

	for i := 0; i < len(updatedPlayers); i++ {
		if updatedPlayers[i].ID == id {
			player = updatedPlayers[i]
			break
		}
	}

	return player
}
