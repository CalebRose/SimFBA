package repository

import (
	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/structs"
	"gorm.io/gorm"
)

type GamesQuery struct {
	TimeSlot        string
	WeekID          string
	SeasonID        string
	IsSpringGames   string
	IsPreseasonGame string
}

func FindCollegeGamesRecords(clauses GamesQuery) []structs.CollegeGame {
	db := dbprovider.GetInstance().GetDB()

	var games []structs.CollegeGame

	query := db.Model(&games)

	if len(clauses.SeasonID) > 0 {
		query = query.Where("season_id = ?", clauses.SeasonID)
	}

	if len(clauses.IsSpringGames) > 0 {
		isSpringGames := clauses.IsSpringGames == "Y"
		query = query.Where("is_spring_game = ?", isSpringGames)
	}

	if len(clauses.TimeSlot) > 0 {
		query = query.Where("time_slot = ?", clauses.TimeSlot)
	}

	if len(clauses.WeekID) > 0 {
		query = query.Where("week_id = ?", clauses.WeekID)
	}

	if err := query.Order("week_id asc").Find(&games).Error; err != nil {
		return []structs.CollegeGame{}
	}
	AssignCFBStreamChannels(games)
	return games
}

func FindNFLGamesRecords(clauses GamesQuery) []structs.NFLGame {
	db := dbprovider.GetInstance().GetDB()

	var games []structs.NFLGame

	query := db.Model(&games)

	if len(clauses.SeasonID) > 0 {
		query = query.Where("season_id = ?", clauses.SeasonID)
	}

	if len(clauses.WeekID) > 0 {
		query = query.Where("week_id = ?", clauses.WeekID)
	}

	if len(clauses.IsPreseasonGame) > 0 {
		isPreseasonGame := clauses.IsPreseasonGame == "Y"
		query = query.Where("is_preseason_game = ?", isPreseasonGame)
	}

	if err := query.Order("week_id asc").Find(&games).Error; err != nil {
		return []structs.NFLGame{}
	}
	return games
}

func CreateCFBGameRecordsBatch(db *gorm.DB, games []structs.CollegeGame, batchSize int) error {
	total := len(games)
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}

		if err := db.CreateInBatches(games[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateNFLGameRecordsBatch(db *gorm.DB, games []structs.NFLGame, batchSize int) error {
	total := len(games)
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}

		if err := db.CreateInBatches(games[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

// AssignCFBStreamChannels fills the projected StreamChannel on each game. It uses the same
// structs.GetCFBStreamChannel as the stream endpoint so the label and the stream cannot drift.
func AssignCFBStreamChannels(games []structs.CollegeGame) {
	if len(games) == 0 {
		return
	}
	db := dbprovider.GetInstance().GetDB()
	var fbsIDs []uint
	db.Model(&structs.CollegeTeam{}).Where("is_fbs = ?", true).Pluck("id", &fbsIDs)
	isFBS := make(map[uint]bool, len(fbsIDs))
	for _, id := range fbsIDs {
		isFBS[id] = true
	}
	for i := range games {
		g := &games[i]
		g.StreamChannel = structs.GetCFBStreamChannel(*g, isFBS[uint(g.HomeTeamID)], isFBS[uint(g.AwayTeamID)])
	}
}
