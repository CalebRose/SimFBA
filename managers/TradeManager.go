package managers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/CalebRose/SimFBA/dbprovider"
	fbsvc "github.com/CalebRose/SimFBA/firebase"
	"github.com/CalebRose/SimFBA/repository"
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

func GetTradePreferencesMap() map[uint]structs.NFLTradePreferences {
	db := dbprovider.GetInstance().GetDB()

	preferences := []structs.NFLTradePreferences{}

	db.Find(&preferences)

	return MakeNFLTradePreferencesMap(preferences)
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

func GetAcceptedTradeProposals() []structs.NFLTradeProposalDTO {
	proposals := repository.FindAllAcceptedTradeProposals()

	acceptedProposals := []structs.NFLTradeProposalDTO{}

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

		acceptedProposals = append(acceptedProposals, proposalResponse)
	}

	return acceptedProposals
}

func GetRejectedTradeProposals() []structs.NFLTradeProposal {
	db := dbprovider.GetInstance().GetDB()

	proposals := []structs.NFLTradeProposal{}

	db.Preload("NFLTeamTradeOptions").Where("is_trade_rejected = ?", true).Find(&proposals)

	return proposals
}

func GetTradeProposalsMap() map[uint][]structs.NFLTradeProposal {
	proposals := repository.FindAllTradeProposalsRecords(repository.TradeClauses{PreloadTradeOptions: true})
	return MakeTradeProposalMap(proposals)
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
		if proposal.IsTradeAccepted || proposal.IsTradeRejected {
			continue
		}
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
		if options.NFLTeamID == proposal.NFLTeamID {
			if options.NFLPlayerID > 0 {
				playerRecord := GetNFLPlayerRecord(strconv.Itoa(int(options.NFLPlayerID)))
				ovrGrade := util.GetNFLOverallGrade(int(playerRecord.Overall))
				ovr := playerRecord.Overall
				if playerRecord.Experience > 1 {
					newsLogMessage += playerRecord.Position + " " + strconv.Itoa(int(ovr)) + " " + playerRecord.FirstName + " " + playerRecord.LastName + " to " + proposal.RecepientTeam + "\n"
				} else {
					newsLogMessage += playerRecord.Position + " " + ovrGrade + " " + playerRecord.FirstName + " " + playerRecord.LastName + " to " + proposal.RecepientTeam + "\n"
				}
			} else if options.NFLDraftPickID > 0 {
				draftPick := GetDraftPickByDraftPickID(strconv.Itoa(int(options.NFLDraftPickID)))
				pickRound := strconv.Itoa(int(draftPick.DraftRound))
				roundAbbreviation := util.GetRoundAbbreviation(pickRound)
				season := strconv.Itoa(int(draftPick.Season))
				newsLogMessage += season + " " + roundAbbreviation + " pick to " + proposal.RecepientTeam + "\n"
			}
		} else {
			if options.NFLPlayerID > 0 {
				playerRecord := GetNFLPlayerRecord(strconv.Itoa(int(options.NFLPlayerID)))
				newsLogMessage += playerRecord.Position + " " + playerRecord.FirstName + " " + playerRecord.LastName + " to " + proposal.NFLTeam + "\n"
			} else if options.NFLDraftPickID > 0 {
				draftPick := GetDraftPickByDraftPickID(strconv.Itoa(int(options.NFLDraftPickID)))
				pickRound := strconv.Itoa(int(draftPick.DraftRound))
				roundAbbreviation := util.GetRoundAbbreviation(pickRound)
				season := strconv.Itoa(int(draftPick.Season))
				newsLogMessage += season + " " + roundAbbreviation + " pick to " + proposal.NFLTeam + "\n"
			}
		}
	}
	newsLogMessage += "\n"

	newsLog := structs.NewsLog{
		TeamID:      0,
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
		TeamID:      0,
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

func CancelTradeProposal(proposalID string) {
	db := dbprovider.GetInstance().GetDB()

	proposal := GetOnlyTradeProposalByProposalID(proposalID)
	options := proposal.NFLTeamTradeOptions

	for _, option := range options {
		db.Delete(&option)
	}

	db.Delete(&proposal)
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
	if proposal.ID == 0 || proposal.IsSynced || len(proposal.NFLTeamTradeOptions) == 0 {
		log.Printf("TradeManager: trade proposal %s is missing, already synced, or has no assets", proposalID)
		return
	}
	SentOptions := proposal.NFLTeamTradeOptions

	details, complete := syncAcceptedOptions(db, SentOptions, proposal.NFLTeamID, proposal.NFLTeam, proposal.RecepientTeamID, proposal.RecepientTeam)

	proposal.ToggleSyncStatus()

	if err := db.Save(&proposal).Error; err != nil {
		log.Printf("TradeManager: failed to save trade proposal %s: %v", proposalID, err)
		return
	}
	if !complete {
		log.Printf("TradeManager: trade proposal %s may be partially processed; report not created", proposalID)
		return
	}
	reportProcessedTrade(proposal, details, false)
}

func SyncTradeFromDraftPage(proposal structs.NFLTradeProposal) bool {
	if len(proposal.NFLTeamTradeOptions)+len(proposal.RecepientTeamTradeOptions) == 0 {
		log.Printf("TradeManager: draft trade has no assets")
		return false
	}
	db := dbprovider.GetInstance().GetDB()
	sent, sentComplete := syncAcceptedOptions(db, proposal.NFLTeamTradeOptions, proposal.NFLTeamID, proposal.NFLTeam, proposal.RecepientTeamID, proposal.RecepientTeam)
	received, receivedComplete := syncAcceptedOptions(db, proposal.RecepientTeamTradeOptions, proposal.RecepientTeamID, proposal.RecepientTeam, proposal.NFLTeamID, proposal.NFLTeam)
	if sentComplete && receivedComplete {
		reportProcessedTrade(proposal, append(sent, received...), true)
	} else {
		log.Printf("TradeManager: draft trade may be partially processed; report not created")
	}
	return false
}

func reportProcessedTrade(proposal structs.NFLTradeProposal, details []string, draft bool) {
	ts := GetTimestamp()
	if ts.ID == 0 {
		log.Printf("TradeManager: cannot report trade %d: timestamp unavailable", proposal.ID)
		return
	}
	db := dbprovider.GetInstance().GetDB()
	var sender, recipient structs.NFLTeam
	if err := db.First(&sender, proposal.NFLTeamID).Error; err != nil {
		log.Printf("TradeManager: cannot report trade %d: sender unavailable: %v", proposal.ID, err)
		return
	}
	if err := db.First(&recipient, proposal.RecepientTeamID).Error; err != nil {
		log.Printf("TradeManager: cannot report trade %d: recipient unavailable: %v", proposal.ID, err)
		return
	}
	proposal.NFLTeam = sender.TeamAbbr
	proposal.RecepientTeam = recipient.TeamAbbr
	id, message, err := tradeReportIdentityAndMessage(proposal, details, draft)
	if err != nil {
		log.Printf("TradeManager: cannot report trade %d: %v", proposal.ID, err)
		return
	}
	report, err := fbsvc.TradeTransactionReport(proposal.NFLTeamID, message, ts)
	if err != nil {
		log.Printf("TradeManager: cannot report trade %d: %v", proposal.ID, err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := fbsvc.CreateAdminReport(ctx, id, report); err != nil {
		log.Printf("TradeManager: failed to create report for trade %d: %v", proposal.ID, err)
	}
}

func tradeReportIdentityAndMessage(proposal structs.NFLTradeProposal, details []string, draft bool) (string, string, error) {
	if proposal.NFLTeamID == 0 || proposal.RecepientTeamID == 0 || proposal.NFLTeamID == proposal.RecepientTeamID {
		return "", "", fmt.Errorf("invalid trade teams %d and %d", proposal.NFLTeamID, proposal.RecepientTeamID)
	}
	id := fmt.Sprintf("trade-proposal-nfl-%d", proposal.ID)
	label := fmt.Sprintf("Trade proposal #%d", proposal.ID)
	if draft {
		label = "Draft trade"
		if proposal.ID != 0 {
			id = fmt.Sprintf("draft-trade-nfl-%d", proposal.ID)
			label = fmt.Sprintf("Draft trade #%d", proposal.ID)
		} else {
			var randomID [16]byte
			if _, err := rand.Read(randomID[:]); err != nil {
				return "", "", fmt.Errorf("generate draft trade report ID: %w", err)
			}
			id = "draft-trade-nfl-" + hex.EncodeToString(randomID[:])
		}
	} else if proposal.ID == 0 {
		return "", "", fmt.Errorf("trade proposal ID is missing")
	}
	message := fmt.Sprintf("%s processed between %s (team %d) and %s (team %d).\n%s",
		label, proposal.NFLTeam, proposal.NFLTeamID, proposal.RecepientTeam, proposal.RecepientTeamID, strings.Join(details, "\n"))
	return id, message, nil
}

func syncAcceptedOptions(db *gorm.DB, options []structs.NFLTradeOption, senderID uint, senderTeam string, recepientID uint, recepientTeam string) ([]string, bool) {
	sendingTeam := GetNFLTeamByTeamID(strconv.Itoa(int(senderID)))
	receivingTeam := GetNFLTeamByTeamID(strconv.Itoa(int(recepientID)))
	sendersCapsheet := GetCapsheetByTeamID(strconv.Itoa(int(senderID)))
	recepientCapsheet := GetCapsheetByTeamID(strconv.Itoa(int(recepientID)))
	if senderID == 0 || recepientID == 0 || senderID == recepientID || sendingTeam.ID == 0 || receivingTeam.ID == 0 || sendersCapsheet.ID == 0 || recepientCapsheet.ID == 0 {
		log.Printf("TradeManager: trade teams %d and %d must exist and have capsheets", senderID, recepientID)
		return nil, false
	}
	salaryMinimum := 0.5
	details := make([]string, 0, len(options))
	complete := true
	for _, option := range options {
		if (option.NFLPlayerID == 0) == (option.NFLDraftPickID == 0) {
			log.Printf("TradeManager: trade option %d must identify exactly one player or draft pick", option.ID)
			complete = false
			continue
		}
		// Contract
		percentage := option.SalaryPercentage
		if option.NFLPlayerID > 0 {
			playerRecord := GetNFLPlayerRecord(strconv.Itoa(int(option.NFLPlayerID)))
			contract := playerRecord.Contract
			if playerRecord.ID == 0 || contract.ID == 0 || (playerRecord.TeamID != int(senderID) && playerRecord.TeamID != int(recepientID)) {
				log.Printf("TradeManager: trade player %d has a missing player, contract, or valid team", option.NFLPlayerID)
				complete = false
				continue
			}
			if playerRecord.TeamID == int(senderID) {
				details = append(details, tradePlayerDetail(playerRecord, sendingTeam.TeamAbbr, receivingTeam.TeamAbbr, percentage))
				sendersPercentage := percentage * 0.01
				receiversPercentage := (100 - percentage) * 0.01
				sendingTeamPay := float64(contract.Y1BaseSalary * sendersPercentage)
				receivingTeamPay := float64(contract.Y1BaseSalary * receiversPercentage)
				// If a team is eating the Y1 Salary for a player
				if sendersPercentage == 1 && contract.Y1Bonus == 0 {
					sendingTeamPay -= salaryMinimum
					receivingTeamPay += salaryMinimum
				} else if contract.Y1BaseSalary == 0 && contract.Y1Bonus > 0 {
					sendingTeamPay -= salaryMinimum
					receivingTeamPay += salaryMinimum
				}
				sendersCapsheet.SubtractFromCapsheetViaTrade(contract)
				sendersCapsheet.NegotiateSalaryDifference(contract.Y1BaseSalary, sendingTeamPay)
				recepientCapsheet.AddContractViaTrade(contract, receivingTeamPay)
				playerRecord.TradePlayer(recepientID, receivingTeam.TeamAbbr)
				contract.TradePlayer(recepientID, receivingTeam.TeamAbbr, receiversPercentage)
			} else {
				details = append(details, tradePlayerDetail(playerRecord, receivingTeam.TeamAbbr, sendingTeam.TeamAbbr, percentage))
				receiversPercentage := percentage * 0.01
				sendersPercentage := (100 - percentage) * 0.01
				sendingTeamPay := float64(contract.Y1BaseSalary * sendersPercentage)
				receivingTeamPay := float64(contract.Y1BaseSalary * receiversPercentage)
				if sendersPercentage == 1 && contract.Y1Bonus == 0 {
					receivingTeamPay -= salaryMinimum
					sendingTeamPay += salaryMinimum
				} else if contract.Y1BaseSalary == 0 && contract.Y1Bonus > 0 {
					receivingTeamPay -= salaryMinimum
					sendingTeamPay += salaryMinimum
				}
				recepientCapsheet.SubtractFromCapsheetViaTrade(contract)
				recepientCapsheet.NegotiateSalaryDifference(contract.Y1BaseSalary, receivingTeamPay)
				sendersCapsheet.AddContractViaTrade(contract, sendingTeamPay)
				playerRecord.TradePlayer(senderID, sendingTeam.TeamAbbr)
				contract.TradePlayer(senderID, sendingTeam.TeamAbbr, sendersPercentage)
			}

			if playerRecord.ID > 0 {
				repository.SaveNFLPlayerRecord(playerRecord, db)
			}
			if contract.ID > 0 {
				repository.SaveNFLContract(contract, db)
			}

		} else if option.NFLDraftPickID > 0 {
			draftPick := GetDraftPickByDraftPickID(strconv.Itoa(int(option.NFLDraftPickID)))
			if draftPick.ID == 0 || (draftPick.TeamID != senderID && draftPick.TeamID != recepientID) {
				log.Printf("TradeManager: trade draft pick %d is missing or has an invalid team", option.NFLDraftPickID)
				complete = false
				continue
			}
			if draftPick.TeamID == senderID {
				details = append(details, tradePickDetail(draftPick, sendingTeam.TeamAbbr, receivingTeam.TeamAbbr))
				draftPick.TradePick(recepientID, recepientTeam)
			} else {
				details = append(details, tradePickDetail(draftPick, receivingTeam.TeamAbbr, sendingTeam.TeamAbbr))
				draftPick.TradePick(senderID, senderTeam)
			}
			if draftPick.ID > 0 {
				if err := db.Save(&draftPick).Error; err != nil {
					log.Printf("TradeManager: failed to save draft pick %d: %v", draftPick.ID, err)
					complete = false
					continue
				}
			}

		}

		if option.ID > 0 {
			if err := db.Where("id = ?", option.ID).Delete(&option).Error; err != nil {
				log.Printf("TradeManager: failed to delete trade option %d: %v", option.ID, err)
				complete = false
			}
		}
	}
	repository.SaveNFLCapsheet(sendersCapsheet, db)
	repository.SaveNFLCapsheet(recepientCapsheet, db)
	return details, complete
}

func tradePlayerDetail(player structs.NFLPlayer, from, to string, percentage float64) string {
	return fmt.Sprintf("%s %s (player %d): %s to %s (salary retention input %.2f%%)", player.FirstName, player.LastName, player.ID, from, to, percentage)
}

func tradePickDetail(pick structs.NFLDraftPick, from, to string) string {
	return fmt.Sprintf("%d round %d pick (pick %d, originally %s): %s to %s", pick.Season, pick.DraftRound, pick.ID, pick.OriginalTeam, from, to)
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
