package managers

import (
	"log"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
)

// GetCollegeDepthChartSnapshots fetches snapshot positions from DB for a given season
func GetCollegeDepthChartSnapshots(seasonID uint) []structs.CollegeDepthChartSeasonSnapshot {
	db := dbprovider.GetInstance().GetDB()
	return repository.GetCollegeDepthChartSnapshots(seasonID, db)
}

// GetNFLDepthChartSnapshots fetches snapshot positions from DB for a given season
func GetNFLDepthChartSnapshots(seasonID uint) []structs.NFLDepthChartSeasonSnapshot {
	db := dbprovider.GetInstance().GetDB()
	return repository.GetNFLDepthChartSnapshots(seasonID, db)
}

// CaptureCollegeDepthChartSnapshot snapshots the live college_depth_chart_positions into college_depth_chart_season_snapshot
// adding the seasonID flag.
func CaptureCollegeDepthChartSnapshot(seasonID uint) ([]structs.CollegeDepthChartSeasonSnapshot, error) {
	db := dbprovider.GetInstance().GetDB()

	positions := repository.GetAllCollegeDepthChartPositions(db)
	records := make([]structs.CollegeDepthChartSeasonSnapshot, len(positions))

	for i, pos := range positions {
		records[i] = structs.CollegeDepthChartSeasonSnapshot{
			SeasonID:         seasonID,
			DepthChartID:     pos.DepthChartID,
			PlayerID:         pos.PlayerID,
			Position:         pos.Position,
			PositionLevel:    pos.PositionLevel,
			FirstName:        pos.FirstName,
			LastName:         pos.LastName,
			OriginalPosition: pos.OriginalPosition,
		}
	}

	err := repository.SaveCollegeDepthChartSnapshots(seasonID, records, db)
	if err != nil {
		return nil, err
	}

	log.Printf("[Snapshot] Successfully snapshotted %d college depth chart positions for Season %d\n", len(records), seasonID)
	return records, nil
}

// CaptureNFLDepthChartSnapshot snapshots the live nfl_depth_chart_positions into nfl_depth_chart_season_snapshot
// adding the seasonID flag.
func CaptureNFLDepthChartSnapshot(seasonID uint) ([]structs.NFLDepthChartSeasonSnapshot, error) {
	db := dbprovider.GetInstance().GetDB()

	positions := repository.GetAllNFLDepthChartPositions(db)
	records := make([]structs.NFLDepthChartSeasonSnapshot, len(positions))

	for i, pos := range positions {
		records[i] = structs.NFLDepthChartSeasonSnapshot{
			SeasonID:         seasonID,
			DepthChartID:     pos.DepthChartID,
			PlayerID:         pos.PlayerID,
			Position:         pos.Position,
			PositionLevel:    pos.PositionLevel,
			FirstName:        pos.FirstName,
			LastName:         pos.LastName,
			OriginalPosition: pos.OriginalPosition,
		}
	}

	err := repository.SaveNFLDepthChartSnapshots(seasonID, records, db)
	if err != nil {
		return nil, err
	}

	log.Printf("[Snapshot] Successfully snapshotted %d NFL depth chart positions for Season %d\n", len(records), seasonID)
	return records, nil
}

// SaveCollegeDepthChartSnapshots saves a batch of college snapshot positions
func SaveCollegeDepthChartSnapshots(seasonID uint, records []structs.CollegeDepthChartSeasonSnapshot) error {
	db := dbprovider.GetInstance().GetDB()
	return repository.SaveCollegeDepthChartSnapshots(seasonID, records, db)
}

// SaveNFLDepthChartSnapshots saves a batch of NFL snapshot positions
func SaveNFLDepthChartSnapshots(seasonID uint, records []structs.NFLDepthChartSeasonSnapshot) error {
	db := dbprovider.GetInstance().GetDB()
	return repository.SaveNFLDepthChartSnapshots(seasonID, records, db)
}
