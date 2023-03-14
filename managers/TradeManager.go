package managers

import (
	"errors"
	"log"
	"strconv"
	"sync"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
	"gorm.io/gorm"
)

func GetTradeBlockDataByTeamID(TeamID string) structs.NFLTradeBlockResponse {
	var waitgroup sync.WaitGroup
	waitgroup.Add(5)
	nflTeamChan := make(chan structs.NFLTeam)
	playersChan := make(chan []structs.NFLPlayer)
	picksChan := make(chan []structs.NFLDraftPick)
	proposalsChan := make(chan structs.NFLTeamProposals)
	preferencesChan := make(chan structs.NFLTradePreferences)

	go func() {
		waitgroup.Wait()
		close(nflTeamChan)
		close(playersChan)
		close(picksChan)
		close(proposalsChan)
		close(preferencesChan)
	}()

	go func() {
		defer waitgroup.Done()
		team := GetNFLTeamWithCapsheetByTeamID(TeamID)
		nflTeamChan <- team
	}()

	go func() {
		defer waitgroup.Done()
		players := GetTradableNFLPlayersByTeamID(TeamID)
		playersChan <- players
	}()

	go func() {
		defer waitgroup.Done()
		picks := GetDraftPicksByTeamID(TeamID)
		picksChan <- picks
	}()

	go func() {
		defer waitgroup.Done()
		proposals := GetTradeProposalsByNFLID(TeamID)
		proposalsChan <- proposals
	}()

	go func() {
		defer waitgroup.Done()
		pref := GetTradePreferencesByTeamID(TeamID)
		preferencesChan <- pref
	}()

	nflTeam := <-nflTeamChan
	tradablePlayers := <-playersChan
	draftPicks := <-picksChan
	teamProposals := <-proposalsChan
	tradePreferences := <-preferencesChan

	// close(nflTeamChan)
	// close(playersChan)
	// close(picksChan)
	// close(proposalsChan)

	return structs.NFLTradeBlockResponse{
		Team:                   nflTeam,
		TradablePlayers:        tradablePlayers,
		DraftPicks:             draftPicks,
		SentTradeProposals:     teamProposals.SentTradeProposals,
		ReceivedTradeProposals: teamProposals.ReceivedTradeProposals,
		TradePreferences:       tradePreferences,
	}
}

func GetOnlyTradeProposalByProposalID(proposalID string) structs.NFLTradeProposal {
	db := dbprovider.GetInstance().GetDB()

	proposal := structs.NFLTradeProposal{}

	db.Preload("NFLTeamTradeOptions").Where("id = ?", proposalID).Find(&proposal)

	return proposal
}

func GetTradePreferencesByTeamID(TeamID string) structs.NFLTradePreferences {
	db := dbprovider.GetInstance().GetDB()

	preferences := structs.NFLTradePreferences{}

	db.Where("id = ?", TeamID).Find(&preferences)

	return preferences
}

func UpdateTradePreferences(pref structs.NFLTradePreferencesDTO) {
	db := dbprovider.GetInstance().GetDB()

	preferences := GetTradePreferencesByTeamID(strconv.Itoa(int(pref.NFLTeamID)))

	preferences.UpdatePreferences(pref)

	db.Save(&preferences)
}

func GetAcceptedTradeProposals() []structs.NFLTradeProposal {
	db := dbprovider.GetInstance().GetDB()

	proposals := []structs.NFLTradeProposal{}

	db.Preload("NFLTeamTradeOptions").Preload("RecepientTeamTradeOptions").Where("is_trade_accepted = ?", true).Find(&proposals)

	return proposals
}

func GetRejectedTradeProposals() []structs.NFLTradeProposal {
	db := dbprovider.GetInstance().GetDB()

	proposals := []structs.NFLTradeProposal{}

	db.Preload("NFLTeamTradeOptions").Where("is_trade_rejected = ?", true).Find(&proposals)

	return proposals
}

// GetTradeProposalsByNFLID -- Returns all trade proposals that were either sent or received
func GetTradeProposalsByNFLID(TeamID string) structs.NFLTeamProposals {
	db := dbprovider.GetInstance().GetDB()

	proposals := []structs.NFLTradeProposal{}

	db.Preload("NFLTeamTradeOptions").Where("nfl_team_id = ? OR recepient_team_id = ?", TeamID, TeamID).Find(&proposals)

	SentProposals := []structs.NFLTradeProposalDTO{}
	ReceivedProposals := []structs.NFLTradeProposalDTO{}

	id := uint(util.ConvertStringToInt(TeamID))

	for _, proposal := range proposals {
		sentOptions := []structs.NFLTradeOptionObj{}
		receivedOptions := []structs.NFLTradeOptionObj{}
		for _, option := range proposal.NFLTeamTradeOptions {
			opt := structs.NFLTradeOptionObj{
				ID:               option.Model.ID,
				TradeProposalID:  option.TradeProposalID,
				NFLTeamID:        option.NFLTeamID,
				SalaryPercentage: option.SalaryPercentage,
				OptionType:       option.OptionType,
			}
			if option.NFLPlayerID > 0 {
				player := GetNFLPlayerRecord(strconv.Itoa(int(option.NFLPlayerID)))
				opt.AssignPlayer(player)
			} else if option.NFLDraftPickID > 0 {
				draftPick := GetDraftPickByDraftPickID(strconv.Itoa((int(option.NFLDraftPickID))))
				opt.AssignPick(draftPick)
			}
			if option.NFLTeamID == proposal.NFLTeamID {
				sentOptions = append(sentOptions, opt)
			} else {
				receivedOptions = append(receivedOptions, opt)
			}
		}

		proposalResponse := structs.NFLTradeProposalDTO{
			ID:                        proposal.Model.ID,
			NFLTeamID:                 proposal.NFLTeamID,
			NFLTeam:                   proposal.NFLTeam,
			RecepientTeamID:           proposal.RecepientTeamID,
			RecepientTeam:             proposal.RecepientTeam,
			IsTradeAccepted:           proposal.IsTradeAccepted,
			IsTradeRejected:           proposal.IsTradeRejected,
			NFLTeamTradeOptions:       sentOptions,
			RecepientTeamTradeOptions: receivedOptions,
		}

		if proposal.NFLTeamID == id {
			SentProposals = append(SentProposals, proposalResponse)
		} else {
			ReceivedProposals = append(ReceivedProposals, proposalResponse)
		}
	}
	return structs.NFLTeamProposals{
		SentTradeProposals:     SentProposals,
		ReceivedTradeProposals: ReceivedProposals,
	}
}

func PlaceNFLPlayerOnTradeBlock(playerID string) {
	db := dbprovider.GetInstance().GetDB()

	player := GetNFLPlayerRecord(playerID)

	player.ToggleTradeBlock()

	db.Save(&player)
}

func CreateTradeProposal(TradeProposal structs.NFLTradeProposalDTO) {
	db := dbprovider.GetInstance().GetDB()
	latestID := GetLatestProposalInDB(db)

	// Create Trade Proposal Object
	proposal := structs.NFLTradeProposal{
		NFLTeamID:       TradeProposal.NFLTeamID,
		NFLTeam:         TradeProposal.NFLTeam,
		RecepientTeamID: TradeProposal.RecepientTeamID,
		RecepientTeam:   TradeProposal.RecepientTeam,
		IsTradeAccepted: false,
		IsTradeRejected: false,
	}
	proposal.AssignID(latestID)

	db.Create(&proposal)

	// Create Trade Options
	SentTradeOptions := TradeProposal.NFLTeamTradeOptions
	ReceivedTradeOptions := TradeProposal.RecepientTeamTradeOptions

	for _, sentOption := range SentTradeOptions {
		tradeOption := structs.NFLTradeOption{
			TradeProposalID:  latestID,
			NFLTeamID:        TradeProposal.NFLTeamID,
			NFLPlayerID:      sentOption.NFLPlayerID,
			NFLDraftPickID:   sentOption.NFLDraftPickID,
			SalaryPercentage: sentOption.SalaryPercentage,
			OptionType:       sentOption.OptionType,
		}
		db.Create(&tradeOption)
	}

	for _, recepientOption := range ReceivedTradeOptions {
		tradeOption := structs.NFLTradeOption{
			TradeProposalID:  latestID,
			NFLTeamID:        TradeProposal.RecepientTeamID,
			NFLPlayerID:      recepientOption.NFLPlayerID,
			NFLDraftPickID:   recepientOption.NFLDraftPickID,
			SalaryPercentage: recepientOption.SalaryPercentage,
			OptionType:       recepientOption.OptionType,
		}
		db.Create(&tradeOption)
	}
}

func AcceptTradeProposal(proposalID string) {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	proposal := GetOnlyTradeProposalByProposalID(proposalID)

	proposal.AcceptTrade()

	// Create News Log
	newsLogMessage := proposal.RecepientTeam + " has accepted a trade offer from " + proposal.NFLTeam + " for trade the following players:\n\n"

	for _, options := range proposal.NFLTeamTradeOptions {
		if options.NFLPlayerID > 0 {
			playerRecord := GetNFLPlayerRecord(strconv.Itoa(int(options.NFLPlayerID)))
			newsLogMessage += playerRecord.Position + " " + playerRecord.FirstName + " " + playerRecord.LastName + " to " + proposal.RecepientTeam + "\n"
		} else if options.NFLDraftPickID > 0 {
			draftPick := GetDraftPickByDraftPickID(strconv.Itoa(int(options.NFLDraftPickID)))
			pickRound := strconv.Itoa(int(draftPick.Round))
			roundAbbreviation := util.GetRoundAbbreviation(pickRound)
			season := strconv.Itoa(int(draftPick.Season))
			newsLogMessage += season + " " + roundAbbreviation + " pick to " + proposal.RecepientTeam + "\n"
		}
	}
	newsLogMessage += "\n"

	for _, options := range proposal.RecepientTeamTradeOptions {
		if options.NFLPlayerID > 0 {
			playerRecord := GetNFLPlayerRecord(strconv.Itoa(int(options.NFLPlayerID)))
			newsLogMessage += playerRecord.Position + " " + playerRecord.FirstName + " " + playerRecord.LastName + " to " + proposal.NFLTeam + "\n"
		} else if options.NFLDraftPickID > 0 {
			draftPick := GetDraftPickByDraftPickID(strconv.Itoa(int(options.NFLDraftPickID)))
			pickRound := strconv.Itoa(int(draftPick.Round))
			roundAbbreviation := util.GetRoundAbbreviation(pickRound)
			season := strconv.Itoa(int(draftPick.Season))
			newsLogMessage += season + " " + roundAbbreviation + " pick to " + proposal.NFLTeam + "\n"
		}
	}

	newsLog := structs.NewsLog{
		WeekID:      ts.NFLWeekID,
		Week:        ts.NFLWeek,
		SeasonID:    ts.NFLSeasonID,
		League:      "NFL",
		MessageType: "Trade",
		Message:     newsLogMessage,
	}

	db.Create(&newsLog)
	db.Save(&proposal)
}

func RejectTradeProposal(proposalID string) {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	proposal := GetOnlyTradeProposalByProposalID(proposalID)

	proposal.RejectTrade()
	newsLog := structs.NewsLog{
		WeekID:      ts.NFLWeekID,
		Week:        ts.NFLWeek,
		SeasonID:    ts.NFLSeasonID,
		League:      "NFL",
		MessageType: "Trade",
		Message:     proposal.RecepientTeam + " has rejected a trade from " + proposal.NFLTeam,
	}

	db.Create(&newsLog)
	db.Save(&proposal)
}

func GetLatestProposalInDB(db *gorm.DB) uint {
	var latestProposal structs.NFLTradeProposal

	err := db.Last(&latestProposal).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 1
		}
		log.Fatalln("ERROR! Could not find latest record" + err.Error())
	}

	return latestProposal.ID + 1
}

func RemoveRejectedTrades() {
	db := dbprovider.GetInstance().GetDB()

	rejectedProposals := GetRejectedTradeProposals()

	for _, proposal := range rejectedProposals {
		sentOptions := proposal.NFLTeamTradeOptions
		deleteOptions(db, sentOptions)

		// Delete Proposal
		db.Delete(&proposal)
	}
}

func SyncAcceptedTrade(proposalID string) {
	db := dbprovider.GetInstance().GetDB()

	proposal := GetOnlyTradeProposalByProposalID(proposalID)
	SentOptions := proposal.NFLTeamTradeOptions

	syncAcceptedOptions(db, SentOptions, proposal.NFLTeamID, proposal.NFLTeam, proposal.RecepientTeamID, proposal.RecepientTeam)

	proposal.ToggleSyncStatus()

	db.Save(&proposal)
}

func syncAcceptedOptions(db *gorm.DB, options []structs.NFLTradeOption, senderID uint, senderTeam string, recepientID uint, recepientTeam string) {
	for _, option := range options {
		if option.NFLPlayerID > 0 {
			playerRecord := GetNFLPlayerRecord(strconv.Itoa(int(option.NFLPlayerID)))
			contract := playerRecord.Contract
			if playerRecord.TeamID == int(senderID) {
				playerRecord.TradePlayer(recepientID, recepientTeam)
				contract.TradePlayer(recepientID, recepientTeam)
			} else {
				playerRecord.TradePlayer(senderID, senderTeam)
				contract.TradePlayer(senderID, senderTeam)
			}

			// Contract
			percentage := option.SalaryPercentage

			// Subtract Contract from Senders's Capsheet
			sendersPercentage := 100 - percentage
			SendersCapsheet := GetCapsheetByTeamID(strconv.Itoa(int(senderID)))
			SendersCapsheet.SubtractFromCapsheet(contract)
			SendersCapsheet.NegotiateSalaryDifference(contract.Y1BaseSalary, float64(contract.Y1BaseSalary*sendersPercentage))

			db.Save(&SendersCapsheet)

			// Add to Recepient Capsheet
			recepientCapsheet := GetCapsheetByTeamID(strconv.Itoa(int(recepientID)))
			recepientCapsheet.AddContractViaTrade(contract, float64(percentage*contract.Y1BaseSalary))
			db.Save(&recepientCapsheet)

		} else if option.NFLDraftPickID > 0 {
			draftPick := GetDraftPickByDraftPickID(strconv.Itoa(int(option.NFLDraftPickID)))
			if draftPick.TeamID == senderID {
				draftPick.TradePick(recepientID, recepientTeam)
			} else {
				draftPick.TradePick(senderID, senderTeam)
			}

			db.Save(&draftPick)
		}

		db.Delete(&option)
	}
}

func VetoTrade(proposalID string) {
	db := dbprovider.GetInstance().GetDB()

	proposal := GetOnlyTradeProposalByProposalID(proposalID)
	SentOptions := proposal.NFLTeamTradeOptions

	deleteOptions(db, SentOptions)

	db.Delete(&proposal)
}

func deleteOptions(db *gorm.DB, options []structs.NFLTradeOption) {
	// Delete Recepient Trade Options
	for _, option := range options {
		// Deletes the option
		db.Delete(&option)
	}
}
