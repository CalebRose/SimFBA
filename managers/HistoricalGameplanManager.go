package managers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
)

func GetHistoricalCollegeGameplansBySeasonAndWeek(seasonID string, weekID string) []structs.CollegeGameplanRecord {
	return repository.GetHistoricalCollegeGameplansBySeasonAndWeek(seasonID, weekID)
}

func GetHistoricalNFLGameplansBySeasonAndWeek(seasonID string, weekID string) []structs.NFLGameplanRecord {
	return repository.GetHistoricalNFLGameplansBySeasonAndWeek(seasonID, weekID)
}

func BackfillCurrentGameplans() {
	fmt.Println("Starting manual backfill of current gameplans for Week 8...")
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	// Ensure tables exist before inserting
	db.AutoMigrate(&structs.CollegeGameplanRecord{})
	db.AutoMigrate(&structs.NFLGameplanRecord{})

	collegeSeasonID := strconv.Itoa(ts.CollegeSeasonID)
	collegeWeekID := strconv.Itoa(ts.CollegeWeekID)

	cfbGames := GetCollegeGamesByWeekIdAndSeasonID(collegeWeekID, collegeSeasonID, false)
	var cfbRecords []structs.CollegeGameplanRecord

	for _, g := range cfbGames {
		homeGP := GetGameplanByTeamID(strconv.Itoa(g.HomeTeamID))
		awayGP := GetGameplanByTeamID(strconv.Itoa(g.AwayTeamID))

		homeJson, _ := json.Marshal(homeGP)
		awayJson, _ := json.Marshal(awayGP)

		cfbRecords = append(cfbRecords, structs.CollegeGameplanRecord{
			GameID:       uint(g.ID),
			TeamID:       uint(g.HomeTeamID),
			SeasonID:     uint(ts.CollegeSeasonID),
			WeekID:       uint(ts.CollegeWeekID),
			IsHome:       true,
			GameplanJSON: string(homeJson),
		})
		cfbRecords = append(cfbRecords, structs.CollegeGameplanRecord{
			GameID:       uint(g.ID),
			TeamID:       uint(g.AwayTeamID),
			SeasonID:     uint(ts.CollegeSeasonID),
			WeekID:       uint(ts.CollegeWeekID),
			IsHome:       false,
			GameplanJSON: string(awayJson),
		})
	}
	if len(cfbRecords) > 0 {
		db.Create(&cfbRecords)
		fmt.Printf("Created %d CFB Gameplan Records for Season %d, Week %d\n", len(cfbRecords), ts.CollegeSeasonID, ts.CollegeWeekID)
	}

	nflSeasonID := strconv.Itoa(ts.NFLSeasonID)
	nflWeekID := strconv.Itoa(ts.NFLWeekID)
	nflGames := GetNFLGamesByWeekAndSeasonID(nflWeekID, nflSeasonID)
	var nflRecords []structs.NFLGameplanRecord

	for _, g := range nflGames {
		homeGP := GetNFLGameplanByTeamID(strconv.Itoa(g.HomeTeamID))
		awayGP := GetNFLGameplanByTeamID(strconv.Itoa(g.AwayTeamID))

		homeJson, _ := json.Marshal(homeGP)
		awayJson, _ := json.Marshal(awayGP)

		nflRecords = append(nflRecords, structs.NFLGameplanRecord{
			GameID:       uint(g.ID),
			TeamID:       uint(g.HomeTeamID),
			SeasonID:     uint(ts.NFLSeasonID),
			WeekID:       uint(ts.NFLWeekID),
			IsHome:       true,
			GameplanJSON: string(homeJson),
		})
		nflRecords = append(nflRecords, structs.NFLGameplanRecord{
			GameID:       uint(g.ID),
			TeamID:       uint(g.AwayTeamID),
			SeasonID:     uint(ts.NFLSeasonID),
			WeekID:       uint(ts.NFLWeekID),
			IsHome:       false,
			GameplanJSON: string(awayJson),
		})
	}
	if len(nflRecords) > 0 {
		db.Create(&nflRecords)
		fmt.Printf("Created %d NFL Gameplan Records for Season %d, Week %d\n", len(nflRecords), ts.NFLSeasonID, ts.NFLWeekID)
	}
	fmt.Println("Backfill Complete.")
}
