package repository

import (
	"fmt"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/structs"
	"gorm.io/gorm"
)

func GetGameplanTESTByTeamID(teamID string) structs.CollegeGameplanTEST {
	db := dbprovider.GetInstance().GetDB()

	var gamePlan structs.CollegeGameplanTEST

	err := db.Where("id = ?", teamID).Find(&gamePlan).Error
	if err != nil {
		fmt.Println(err)
		return structs.CollegeGameplanTEST{}
	}

	return gamePlan
}

func GetDCTESTByTeamID(teamID string) structs.CollegeTeamDepthChartTEST {
	db := dbprovider.GetInstance().GetDB()

	var dc structs.CollegeTeamDepthChartTEST

	err := db.Preload("DepthChartPlayers").Where("id = ?", teamID).Find(&dc).Error
	if err != nil {
		fmt.Println(err)
		return structs.CollegeTeamDepthChartTEST{}
	}

	return dc
}

func GetDepthChartPositionPlayersByDepthchartIDTEST(depthChartID string) []structs.CollegeDepthChartPositionTEST {
	db := dbprovider.GetInstance().GetDB()

	var positionPlayers []structs.CollegeDepthChartPositionTEST

	err := db.Where("depth_chart_id = ?", depthChartID).Find(&positionPlayers).Error
	if err != nil {
		fmt.Println(err)
		panic("Depth Chart does not exist for this ID")
	}

	return positionPlayers
}

func CreateCollegeGameplansTESTRecordsBatch(db *gorm.DB, gps []structs.CollegeGameplanTEST, batchSize int) error {
	total := len(gps)
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}

		if err := db.CreateInBatches(gps[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func GetNFLGameplanTESTByTeamID(teamID string) structs.NFLGameplanTEST {
	db := dbprovider.GetInstance().GetDB()

	var gamePlan structs.NFLGameplanTEST

	err := db.Where("id = ?", teamID).Find(&gamePlan).Error
	if err != nil {
		fmt.Println(err)
		return structs.NFLGameplanTEST{}
	}

	return gamePlan
}

func GetNFLDCTESTByTeamID(teamID string) structs.NFLDepthChartTEST {
	db := dbprovider.GetInstance().GetDB()

	var dc structs.NFLDepthChartTEST

	err := db.Preload("DepthChartPlayers").Where("id = ?", teamID).Find(&dc).Error
	if err != nil {
		fmt.Println(err)
		return structs.NFLDepthChartTEST{}
	}

	return dc
}

func GetNFLDepthChartPositionPlayersByDepthchartIDTEST(depthChartID string) []structs.NFLDepthChartPositionTEST {
	db := dbprovider.GetInstance().GetDB()

	var positionPlayers []structs.NFLDepthChartPositionTEST

	err := db.Where("depth_chart_id = ?", depthChartID).Find(&positionPlayers).Error
	if err != nil {
		fmt.Println(err)
		panic("Depth Chart does not exist for this ID")
	}

	return positionPlayers
}

/* Batch Create Functions */

func CreateNFLGameplansTESTRecordsBatch(db *gorm.DB, gps []structs.NFLGameplanTEST, batchSize int) error {
	total := len(gps)
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}
		if err := db.CreateInBatches(gps[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateNFLDepthChartTESTBatch(db *gorm.DB, positions []structs.NFLDepthChartTEST, batchSize int) error {
	total := len(positions)
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}
		if err := db.CreateInBatches(positions[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateNFLDepthChartPositionsTESTBatch(db *gorm.DB, positions []structs.NFLDepthChartPositionTEST, batchSize int) error {
	total := len(positions)
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}
		if err := db.CreateInBatches(positions[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}
