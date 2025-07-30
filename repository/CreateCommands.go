package repository

import (
	"log"

	"github.com/CalebRose/SimFBA/models"
	"github.com/CalebRose/SimFBA/structs"
	"gorm.io/gorm"
)

func CreateCFBPlayByPlaysInBatch(plays []structs.CollegePlayByPlay, db *gorm.DB) {
	err := db.CreateInBatches(&plays, len(plays)).Error
	if err != nil {
		log.Panicln("Could not save play by plays!")
	}
}

func CreateNFLPlayByPlaysInBatch(plays []structs.NFLPlayByPlay, db *gorm.DB) {
	err := db.CreateInBatches(&plays, len(plays)).Error
	if err != nil {
		log.Panicln("Could not save play by plays!")
	}
}

func CreateCFBTeamStatsInBatch(stats []structs.CollegeTeamStats, db *gorm.DB) {
	err := db.CreateInBatches(&stats, len(stats)).Error
	if err != nil {
		log.Panicln("Could not save team stats!")
	}
}

func CreateCFBPlayerStatsInBatch(stats []structs.CollegePlayerStats, db *gorm.DB) {
	err := db.CreateInBatches(&stats, len(stats)).Error
	if err != nil {
		log.Panicln("Could not save team stats!")
	}
}

func CreateNFLTeamStatsInBatch(stats []structs.NFLTeamStats, db *gorm.DB) {
	err := db.CreateInBatches(&stats, len(stats)).Error
	if err != nil {
		log.Panicln("Could not save team stats!")
	}
}

func CreateNFLPlayerStatsInBatch(stats []structs.NFLPlayerStats, db *gorm.DB) {
	err := db.CreateInBatches(&stats, len(stats)).Error
	if err != nil {
		log.Panicln("Could not save team stats!")
	}
}

func CreateRecruitProfileRecord(cp structs.RecruitPlayerProfile, db *gorm.DB) {
	err := db.Create(&cp).Error
	if err != nil {
		log.Panicln("Could not create recruit profile record!")
	}
}

func CreateCFBSnapsInBatch(snaps []structs.CollegePlayerGameSnaps, db *gorm.DB) {
	err := db.CreateInBatches(&snaps, len(snaps)).Error
	if err != nil {
		log.Panicln("Could not create college snaps in batch!")
	}
}

func CreateNFLSnapsInBatch(snaps []structs.NFLPlayerGameSnaps, db *gorm.DB) {
	err := db.CreateInBatches(&snaps, len(snaps)).Error
	if err != nil {
		log.Panicln("Could not create nfl snaps in batch!")
	}
}

func CreateCFBSeasonSnaps(snap structs.CollegePlayerSeasonSnaps, db *gorm.DB) {
	err := db.Create(&snap).Error
	if err != nil {
		log.Panicln("Could not create cfb season snaps record!")
	}
}

func CreateNFLSeasonSnaps(snap structs.NFLPlayerSeasonSnaps, db *gorm.DB) {
	err := db.Create(&snap).Error
	if err != nil {
		log.Panicln("Could not create nfl season snaps record!")
	}
}

func CreateCollegePromiseRecord(promise structs.CollegePromise, db *gorm.DB) {
	// Save College Player Record
	err := db.Create(&promise).Error
	if err != nil {
		log.Panicln("Could not save new college recruit record")
	}
}

func CreateNotification(noti structs.Notification, db *gorm.DB) {
	err := db.Create(&noti).Error
	if err != nil {
		log.Panicln("Could not create notification record!")
	}
}

func CreateCFBPlayerRecord(player structs.CollegePlayer, db *gorm.DB) {
	err := db.Create(&player).Error
	if err != nil {
		log.Panicln("Could not create cfb season snaps record!")
	}
}

func CreateHistoricCFBPlayerRecord(player structs.HistoricCollegePlayer, db *gorm.DB) {
	err := db.Create(&player).Error
	if err != nil {
		log.Panicln("Could not create cfb season snaps record!")
	}
}

func CreateNFLDrafteeRecord(player models.NFLDraftee, db *gorm.DB) {
	err := db.Create(&player).Error
	if err != nil {
		log.Panicln("Could not create cfb season snaps record!")
	}
}

func CreateRetireeRecord(player structs.NFLRetiredPlayer, db *gorm.DB) {
	err := db.Create(&player).Error
	if err != nil {
		log.Panicln("Could not create cfb season snaps record!")
	}
}

func CreateNFLDrafteesInBatches(db *gorm.DB, draftees []models.NFLDraftee, batchSize int) error {
	// Create the records in batches with the specified batch size
	if err := db.CreateInBatches(draftees, batchSize).Error; err != nil {
		return err
	}

	return nil
}

func CreateCFBRecruitProfileRecordsBatch(db *gorm.DB, profiles []structs.RecruitPlayerProfile, batchSize int) error {
	total := len(profiles)
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}

		if err := db.CreateInBatches(profiles[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateCFBRecruitRecordsBatch(db *gorm.DB, croots []structs.Recruit, batchSize int) error {
	total := len(croots)
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}

		if err := db.CreateInBatches(croots[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateGlobalPlayerRecordsBatch(db *gorm.DB, croots []structs.Player, batchSize int) error {
	total := len(croots)
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}

		if err := db.CreateInBatches(croots[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateNFLDrafteesSafely(db *gorm.DB, draftees []models.NFLDraftee, batchSize int) error {
	total := len(draftees)
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}

		if err := db.CreateInBatches(draftees[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateTransferPortalProfileRecord(profile structs.TransferPortalProfile, db *gorm.DB) {
	err := db.Create(&profile).Error
	if err != nil {
		log.Panicln("Could not create cfb season snaps record!")
	}
}

func CreateTransferPortalProfileRecordsBatch(db *gorm.DB, profiles []structs.TransferPortalProfile, batchSize int) error {
	total := len(profiles)
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}

		if err := db.CreateInBatches(profiles[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

// CreateNFLTeamSeasonStats
func CreateNFLTeamSeasonStats(seasonStats structs.NFLTeamSeasonStats, db *gorm.DB) {
	// Save NFL Team Season Stats Record
	err := db.Create(&seasonStats).Error
	if err != nil {
		log.Panicln("Could not save nfl team seasons stats record")
	}
}

// CreateCFBTeamSeasonStats
func CreateCFBTeamSeasonStats(seasonStats structs.CollegeTeamSeasonStats, db *gorm.DB) {
	// Save NFL Team Season Stats Record
	err := db.Create(&seasonStats).Error
	if err != nil {
		log.Panicln("Could not save nfl team seasons stats record")
	}
}

// CreateNFLTeamSeasonStats
func CreateNFLPlayerSeasonStats(seasonStats structs.NFLPlayerSeasonStats, db *gorm.DB) {
	// Save NFL Player Season Stats Record
	err := db.Create(&seasonStats).Error
	if err != nil {
		log.Panicln("Could not save nfl Player seasons stats record")
	}
}

// CreateCFBTeamSeasonStats
func CreateCFBPlayerSeasonStats(seasonStats structs.CollegePlayerSeasonStats, db *gorm.DB) {
	// Save NFL Player Season Stats Record
	err := db.Create(&seasonStats).Error
	if err != nil {
		log.Panicln("Could not save nfl Player seasons stats record")
	}
}

func CreateCollegeGameplansRecordsBatch(db *gorm.DB, gps []structs.CollegeGameplan, batchSize int) error {
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

func CreateCollegeDCRecordsBatch(db *gorm.DB, gps []structs.CollegeTeamDepthChart, batchSize int) error {
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

func CreateCollegeDCTESTRecordsBatch(db *gorm.DB, gps []structs.CollegeTeamDepthChartTEST, batchSize int) error {
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

func CreateCollegeDCPRecordsBatch(db *gorm.DB, gps []structs.CollegeDepthChartPosition, batchSize int) error {
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

func CreateCollegeDCPTESTRecordsBatch(db *gorm.DB, gps []structs.CollegeDepthChartPositionTEST, batchSize int) error {
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

func CreateNFLPlayerRecordsBatch(db *gorm.DB, gps []structs.NFLPlayer, batchSize int) error {
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

func CreateNFLContractRecordsBatch(db *gorm.DB, gps []structs.NFLContract, batchSize int) error {
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

func CreateHistoricCFBRecordsBatch(db *gorm.DB, gps []structs.HistoricCollegePlayer, batchSize int) error {
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

func CreateCFBPlayerSeasonStatsRecordsBatch(db *gorm.DB, gps []structs.CollegePlayerSeasonStats, batchSize int) error {
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

func CreatePreDraftEventResultsBatch(db *gorm.DB, events []models.EventResults, batchSize int) error {
	total := len(events)
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}
		if err := db.CreateInBatches(events[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateCFBStandingsBatch(db *gorm.DB, standings []structs.CollegeStandings, batchSize int) error {
	total := len(standings)
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}
		if err := db.CreateInBatches(standings[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateNFLStandingsBatch(db *gorm.DB, standings []structs.NFLStandings, batchSize int) error {
	total := len(standings)
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}

		if err := db.CreateInBatches(standings[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateCFBRecruitPointAllocationRecordsBatch(db *gorm.DB, profiles []structs.RecruitPointAllocation, batchSize int) error {
	total := len(profiles)
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}

		if err := db.CreateInBatches(profiles[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateCFBDepthChartPositionRecordsBatch(db *gorm.DB, profiles []structs.CollegeDepthChartPosition, batchSize int) error {
	total := len(profiles)
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}

		if err := db.CreateInBatches(profiles[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateCFBDepthChartPositionTESTRecordsBatch(db *gorm.DB, profiles []structs.CollegeDepthChartPositionTEST, batchSize int) error {
	total := len(profiles)
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}

		if err := db.CreateInBatches(profiles[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateNFLDepthChartPositionRecordsBatch(db *gorm.DB, profiles []structs.NFLDepthChartPosition, batchSize int) error {
	total := len(profiles)
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}

		if err := db.CreateInBatches(profiles[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}
