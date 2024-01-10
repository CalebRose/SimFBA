package managers

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/models"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
	"gorm.io/gorm"
)

func GetAllFreeAgents() []structs.NFLPlayer {
	db := dbprovider.GetInstance().GetDB()

	fas := []structs.NFLPlayer{}

	db.Where("is_free_agent = ?", true).Find(&fas)

	return fas
}

func GetAllWaiverWirePlayers() []structs.NFLPlayer {
	db := dbprovider.GetInstance().GetDB()

	fas := []structs.NFLPlayer{}

	db.Where("is_waived = ?", true).Find(&fas)

	return fas
}

func GetAllAvailableNFLPlayers(TeamID string) models.FreeAgencyResponse {
	var wg sync.WaitGroup
	wg.Add(5)
	var (
		FAs           []models.FreeAgentResponse
		WaiverPlayers []models.WaiverWirePlayerResponse
		Offers        []structs.FreeAgencyOffer
		PracticeSquad []models.FreeAgentResponse
		roster        []structs.NFLPlayer
	)
	go func() {
		defer wg.Done()
		FAs = GetAllFreeAgentsWithOffers()
	}()
	go func() {
		defer wg.Done()
		WaiverPlayers = GetAllWaiverWirePlayersFAPage()
	}()
	go func() {
		defer wg.Done()
		Offers = GetFreeAgentOffersByTeamID(TeamID)
	}()
	go func() {
		defer wg.Done()
		PracticeSquad = GetAllPracticeSquadPlayersForFAPage()

	}()
	go func() {
		defer wg.Done()
		roster = GetNFLPlayersWithContractsByTeamID(TeamID)

	}()
	wg.Wait()

	count := 0

	for _, p := range roster {
		if p.IsPracticeSquad || p.InjuryReserve {
			continue
		}
		count += 1
	}

	return models.FreeAgencyResponse{
		FreeAgents:    FAs,
		WaiverPlayers: WaiverPlayers,
		PracticeSquad: PracticeSquad,
		TeamOffers:    Offers,
		RosterCount:   uint(count),
	}
}

// GetAllFreeAgentsWithOffers -- For Free Agency UI Page.
func GetAllFreeAgentsWithOffers() []models.FreeAgentResponse {
	ts := GetTimestamp()
	db := dbprovider.GetInstance().GetDB()

	fas := []structs.NFLPlayer{}

	seasonID := 0
	if !ts.IsNFLOffSeason {
		seasonID = ts.NFLSeasonID
	} else {
		seasonID = ts.NFLSeasonID - 1
	}
	seasonStr := strconv.Itoa(seasonID)
	db.Preload("Offers", func(db *gorm.DB) *gorm.DB {
		return db.Order("contract_value DESC").Where("is_active = true")
	}).Preload("SeasonStats", func(db *gorm.DB) *gorm.DB {
		return db.Where("season_id = ?", seasonStr)
	}).Order("overall desc").Where("is_free_agent = ? AND overall > ?", true, "48").Find(&fas)

	sort.Slice(fas[:], func(i, j int) bool {
		if fas[i].ShowLetterGrade {
			return true
		}
		if fas[j].ShowLetterGrade {
			return false
		}
		return fas[i].Overall > fas[j].Overall
	})

	faResponseList := make([]models.FreeAgentResponse, len(fas))

	for i, fa := range fas {
		faResponseList[i] = models.FreeAgentResponse{
			ID:              fa.ID,
			PlayerID:        fa.PlayerID,
			TeamID:          fa.TeamID,
			College:         fa.College,
			TeamAbbr:        fa.TeamAbbr,
			BasePlayer:      fa.BasePlayer,
			Experience:      fa.Experience,
			Hometown:        fa.Hometown,
			State:           fa.State,
			IsActive:        fa.IsActive,
			IsWaived:        fa.IsWaived,
			IsPracticeSquad: fa.IsPracticeSquad,
			IsFreeAgent:     fa.IsFreeAgent,
			MinimumValue:    fa.MinimumValue,
			PreviousTeam:    fa.PreviousTeam,
			DraftedTeam:     fa.DraftedTeam,
			ShowLetterGrade: fa.ShowLetterGrade,
			Stats:           fa.Stats,
			SeasonStats:     fa.SeasonStats,
			Offers:          fa.Offers,
		}
	}

	return faResponseList
}

func GetAllWaiverWirePlayersFAPage() []models.WaiverWirePlayerResponse {
	db := dbprovider.GetInstance().GetDB()

	WaivedPlayers := []structs.NFLPlayer{}

	db.Preload("WaiverOffers", func(db *gorm.DB) *gorm.DB {
		return db.Order("waiver_order asc").Where("is_active = true")
	}).Preload("Contract", func(db *gorm.DB) *gorm.DB {
		return db.Where("is_active = true")
	}).Where("is_waived = ?", true).Find(&WaivedPlayers)

	sort.Slice(WaivedPlayers[:], func(i, j int) bool {
		if WaivedPlayers[i].ShowLetterGrade {
			return true
		}
		if WaivedPlayers[j].ShowLetterGrade {
			return false
		}
		return WaivedPlayers[i].Overall > WaivedPlayers[j].Overall
	})

	faResponseList := make([]models.WaiverWirePlayerResponse, len(WaivedPlayers))

	for i, fa := range WaivedPlayers {
		faResponseList[i] = models.WaiverWirePlayerResponse{
			ID:              fa.ID,
			PlayerID:        fa.PlayerID,
			TeamID:          fa.TeamID,
			College:         fa.College,
			TeamAbbr:        fa.TeamAbbr,
			BasePlayer:      fa.BasePlayer,
			Experience:      fa.Experience,
			Hometown:        fa.Hometown,
			State:           fa.State,
			IsActive:        fa.IsActive,
			IsWaived:        fa.IsWaived,
			IsPracticeSquad: fa.IsPracticeSquad,
			IsFreeAgent:     fa.IsFreeAgent,
			MinimumValue:    fa.MinimumValue,
			PreviousTeam:    fa.PreviousTeam,
			DraftedTeam:     fa.DraftedTeam,
			ShowLetterGrade: fa.ShowLetterGrade,
			WaiverOffers:    fa.WaiverOffers,
			Contract:        fa.Contract,
		}
	}

	return faResponseList
}

func GetFreeAgentOffersByTeamID(TeamID string) []structs.FreeAgencyOffer {
	db := dbprovider.GetInstance().GetDB()

	offers := []structs.FreeAgencyOffer{}

	err := db.Where("team_id = ? AND is_active = ?", TeamID, true).Find(&offers).Error
	if err != nil {
		return offers
	}

	return offers
}

func GetFreeAgentOffersByPlayerID(PlayerID string) []structs.FreeAgencyOffer {
	db := dbprovider.GetInstance().GetDB()

	offers := []structs.FreeAgencyOffer{}

	err := db.Where("nfl_player_id = ? AND is_active = ?", PlayerID, true).Find(&offers).Error
	if err != nil {
		return offers
	}

	return offers
}

func GetFreeAgentOfferByOfferID(OfferID string) structs.FreeAgencyOffer {
	db := dbprovider.GetInstance().GetDB()

	offer := structs.FreeAgencyOffer{}

	err := db.Where("id = ?", OfferID).Find(&offer).Error
	if err != nil {
		return offer
	}

	return offer
}

func CreateFAOffer(offer structs.FreeAgencyOfferDTO) structs.FreeAgencyOffer {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	freeAgentOffer := GetFreeAgentOfferByOfferID(strconv.Itoa(int(offer.ID)))
	player := GetNFLPlayerRecord(strconv.Itoa(int(offer.NFLPlayerID)))

	if freeAgentOffer.ID == 0 {
		id := GetLatestFreeAgentOfferInDB(db)
		freeAgentOffer.AssignID(id)
	}

	if ts.IsFreeAgencyLocked {
		return freeAgentOffer
	}

	freeAgentOffer.CalculateOffer(offer)

	// If the owning team is sending an offer to a player
	if player.IsPracticeSquad && player.TeamID == int(offer.TeamID) {
		SignFreeAgent(freeAgentOffer, player, ts)
	} else {
		db.Save(&freeAgentOffer)

		fmt.Println("Creating offer!")
	}

	if player.IsPracticeSquad && player.TeamID != int(offer.TeamID) {
		message := offer.Team + " have placed an offer on " + player.TeamAbbr + " " + player.Position + " " + player.FirstName + " " + player.LastName + " to pick up from the practice squad."
		CreateNewsLog("NFL", message, "Free Agency", player.TeamID, ts)
	}

	return freeAgentOffer
}

func CancelOffer(offer structs.FreeAgencyOfferDTO) {
	db := dbprovider.GetInstance().GetDB()

	OfferID := strconv.Itoa(int(offer.ID))

	freeAgentOffer := GetFreeAgentOfferByOfferID(OfferID)

	freeAgentOffer.CancelOffer()

	db.Save(&freeAgentOffer)
}

func SignFreeAgent(offer structs.FreeAgencyOffer, FreeAgent structs.NFLPlayer, ts structs.Timestamp) {
	db := dbprovider.GetInstance().GetDB()

	NFLTeam := GetNFLTeamByTeamID(strconv.Itoa(int(offer.TeamID)))
	Contract := structs.NFLContract{}
	messageStart := "FA "
	if !FreeAgent.IsPracticeSquad {
		Contract = structs.NFLContract{
			PlayerID:       FreeAgent.PlayerID,
			NFLPlayerID:    FreeAgent.PlayerID,
			TeamID:         NFLTeam.ID,
			Team:           NFLTeam.TeamAbbr,
			OriginalTeamID: NFLTeam.ID,
			OriginalTeam:   NFLTeam.TeamAbbr,
			ContractLength: offer.ContractLength,
			Y1BaseSalary:   offer.Y1BaseSalary,
			Y1Bonus:        offer.Y1Bonus,
			Y2BaseSalary:   offer.Y2BaseSalary,
			Y2Bonus:        offer.Y2Bonus,
			Y3BaseSalary:   offer.Y3BaseSalary,
			Y3Bonus:        offer.Y3Bonus,
			Y4BaseSalary:   offer.Y4BaseSalary,
			Y4Bonus:        offer.Y4Bonus,
			Y5BaseSalary:   offer.Y5BaseSalary,
			Y5Bonus:        offer.Y5Bonus,
			ContractValue:  offer.ContractValue,
			IsActive:       true,
			IsComplete:     false,
			IsExtended:     false,
		}

		db.Create(&Contract)
	} else {
		Contract = GetContractByPlayerID(strconv.Itoa(int(FreeAgent.ID)))
		Contract.MapPracticeSquadOffer(offer)
		db.Save(&Contract)
		messageStart = "PS "
	}
	FreeAgent.SignPlayer(int(NFLTeam.ID), NFLTeam.TeamAbbr)
	db.Save(&FreeAgent)

	// News Log
	message := messageStart + FreeAgent.Position + " " + FreeAgent.FirstName + " " + FreeAgent.LastName + " has signed with the " + NFLTeam.TeamName + " with a contract worth approximately $" + strconv.Itoa(int(Contract.ContractValue)) + " Million Dollars."
	CreateNewsLog("NFL", message, "Free Agency", int(offer.TeamID), ts)
}

func SyncFreeAgencyOffers() {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()
	ts.ToggleFALock()
	db.Save(&ts)
	FreeAgents := GetAllFreeAgents()

	capsheetMap := getCapsheetMap()

	for _, FA := range FreeAgents {
		// If the Free Agent is not available in off-season free agency anymore
		if ts.IsNFLOffSeason && !FA.IsNegotiating && !FA.IsAcceptingOffers {
			continue
		}

		// Check if still accepting offers
		if ts.IsNFLOffSeason && FA.IsAcceptingOffers && ts.FreeAgencyRound < FA.NegotiationRound {
			continue
		}

		if ts.IsNFLOffSeason && FA.IsAcceptingOffers && ts.FreeAgencyRound >= FA.NegotiationRound {
			FA.ToggleIsNegotiating()
			db.Save(&FA)
			continue
		}

		// Check if still negotiation
		if ts.IsNFLOffSeason && FA.IsNegotiating && ts.FreeAgencyRound < FA.SigningRound {
			continue
		}

		// Is Ready to Sign
		Offers := GetFreeAgentOffersByPlayerID(strconv.Itoa(int(FA.ID)))

		// Sort by highest contract value
		sort.Sort(structs.ByContractValue(Offers))

		WinningOffer := structs.FreeAgencyOffer{}

		for _, Offer := range Offers {
			// Calculate to see if team can afford to pay for contract in Y1
			capsheet := capsheetMap[Offer.TeamID]
			if capsheet.ID == 0 {
				// Invalid!!
				continue
			}
			y1CapSpace := ts.Y1Capspace - capsheet.Y1Bonus - capsheet.Y1Salary - capsheet.Y1CapHit
			y1Remaining := y1CapSpace - Offer.Y1BaseSalary - Offer.Y1Bonus
			if y1CapSpace < 0 || y1Remaining < 0 {
				continue
			}
			// Get the Contract with the best value for the FA
			if Offer.IsActive && WinningOffer.ID == 0 {
				WinningOffer = Offer
			}
			if Offer.IsActive {
				Offer.CancelOffer()
			}

			db.Save(&Offer)
		}

		if WinningOffer.ID > 0 {
			SignFreeAgent(WinningOffer, FA, ts)
		} else if ts.IsNFLOffSeason {
			FA.WaitUntilAfterDraft()
			db.Save(&FA)
		}
	}

	WaiverWirePlayers := GetAllWaiverWirePlayers()

	for _, w := range WaiverWirePlayers {
		if len(w.WaiverOffers) == 0 {
			// Deactivate Contract, convert to Free Agent
			w.ConvertWaivedPlayerToFA()
			contract := GetContractByPlayerID(strconv.Itoa(int(w.ID)))
			contract.DeactivateContract()
			db.Delete(&contract)
		} else {
			var winningOffer structs.NFLWaiverOffer
			offers := GetWaiverOffersByPlayerID(strconv.Itoa(int(w.ID)))
			contract := GetContractByPlayerID(strconv.Itoa(int(w.ID)))
			for _, Offer := range offers {
				// Calculate to see if team can afford to pay for contract in Y1
				capsheet := capsheetMap[Offer.TeamID]
				if capsheet.ID == 0 {
					// Invalid!!
					continue
				}
				y1CapSpace := ts.Y1Capspace - capsheet.Y1Bonus - capsheet.Y1Salary - capsheet.Y1CapHit
				y1Remaining := y1CapSpace - contract.Y1BaseSalary - contract.Y1Bonus
				if y1CapSpace < 0 || y1Remaining < 0 {
					continue
				}
				winningOffer = Offer
				break
			}
			if winningOffer.ID == 0 {
				continue
			}
			w.SignPlayer(int(winningOffer.TeamID), winningOffer.Team)
			contract.ReassignTeam(winningOffer.TeamID, winningOffer.Team)
			db.Save(&contract)

			message := w.Position + " " + w.FirstName + " " + w.LastName + " was picked up on the Waiver Wire by " + winningOffer.Team
			CreateNewsLog("NFL", message, "Free Agency", int(winningOffer.TeamID), ts)

			// Recalibrate winning team's remaining offers
			teamOffers := GetWaiverOffersByTeamID(strconv.Itoa(int(winningOffer.TeamID)))
			team := GetNFLTeamByTeamID(strconv.Itoa(int(winningOffer.TeamID)))

			team.AssignWaiverOrder(team.WaiverOrder + 32)
			db.Save(&team)

			for _, o := range teamOffers {
				o.AssignNewWaiverOrder(team.WaiverOrder + 32)
				db.Save(&o)
			}

			// Delete current waiver offers
			for _, o := range offers {
				db.Delete(&o)
			}
		}
		db.Save(&w)
	}

	practiceSquad := GetAllPracticeSquadPlayers()

	for _, p := range practiceSquad {
		Offers := GetFreeAgentOffersByPlayerID(strconv.Itoa(int(p.ID)))
		contract := GetContractByPlayerID(strconv.Itoa(int(p.ID)))
		if len(Offers) == 0 {
			continue
		}
		ownerTeam := p.TeamID
		ownerOffer := structs.FreeAgencyOffer{}

		for _, o := range Offers {
			if int(o.TeamID) == ownerTeam && o.IsActive {
				ownerOffer = o
				break
			}
		}
		if ownerOffer.ID > 0 {
			SignFreeAgent(ownerOffer, p, ts)
			db.Save(&p)
		} else {
			sort.Sort(structs.ByContractValue(Offers))

			WinningOffer := structs.FreeAgencyOffer{}

			for _, Offer := range Offers {
				// Calculate to see if team can afford to pay for contract in Y1
				capsheet := capsheetMap[Offer.TeamID]
				if capsheet.ID == 0 {
					// Invalid!!
					continue
				}
				y1CapSpace := ts.Y1Capspace - capsheet.Y1Bonus - capsheet.Y1Salary - capsheet.Y1CapHit
				y1Remaining := y1CapSpace - contract.Y1BaseSalary - contract.Y1Bonus
				if y1CapSpace < 0 || y1Remaining < 0 {
					continue
				}
				// Get the Contract with the best value for the FA
				if Offer.IsActive && WinningOffer.ID == 0 {
					WinningOffer = Offer
				}
				if Offer.IsActive {
					Offer.CancelOffer()
				}

				db.Save(&Offer)
			}

			if WinningOffer.ID > 0 {
				SignFreeAgent(WinningOffer, p, ts)
			} else if ts.IsNFLOffSeason {
				p.WaitUntilAfterDraft()
				db.Save(&p)
			}
		}
	}

	ts.ToggleFALock()
	ts.ToggleGMActions()
	db.Save(&ts)
}

func SyncExtensionOffers() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID := strconv.Itoa(ts.NFLSeasonID)

	nflTeams := GetAllNFLTeams()

	for _, team := range nflTeams {
		teamID := strconv.Itoa(int(team.ID))
		roster := GetNFLPlayersForRosterPage(teamID)

		for _, player := range roster {
			min := player.MinimumValue
			contract := player.Contract
			if contract.ContractLength == 1 && len(player.Extensions) > 0 {
				for idx, e := range player.Extensions {
					if e.IsRejected || !e.IsActive {
						continue
					}
					minimumValueMultiplier := 1.0
					validation := validateFreeAgencyPref(player, roster, team, seasonID, e.ContractLength, idx)
					// If the offer is valid and meets the player's free agency bias, reduce the minimum value required by 15%
					if validation && player.FreeAgency != "Average" {
						minimumValueMultiplier = 0.85
						// If the offer does not meet the player's free agency bias, increase the minimum value required by 15%
					} else if !validation && player.FreeAgency != "Average" {
						minimumValueMultiplier = 1.15
					}
					percentage := (e.ContractValue / (min * minimumValueMultiplier) * 100)
					odds := getExtensionPercentageOdds(percentage)
					// Run Check on the Extension

					roll := util.GenerateFloatFromRange(1, 100)
					message := ""
					if odds == 0 || roll > odds {
						// Rejects offer
						e.DeclineOffer(ts.NFLWeek)
						player.DeclineOffer(ts.NFLWeek)
						if e.IsRejected || player.Rejections > 2 {
							message = player.Position + " " + player.FirstName + " " + player.LastName + " has rejected an extension offer from " + e.Team + " worth approximately $" + strconv.Itoa(int(e.ContractValue)) + " Million Dollars and will enter Free Agency."
						} else {
							message = player.Position + " " + player.FirstName + " " + player.LastName + " has declined an extension offer from " + e.Team + " with an extension worth approximately $" + strconv.Itoa(int(e.ContractValue)) + " Million Dollars, and is still negotiating."
						}
						CreateNewsLog("NFL", message, "Free Agency", int(e.TeamID), ts)
						db.Save(&player)
					} else {
						e.AcceptOffer()
						message = player.Position + " " + player.FirstName + " " + player.LastName + " has accepted an extension offer from " + e.Team + " worth approximately $" + strconv.Itoa(int(e.ContractValue)) + " Million Dollars."
						CreateNewsLog("NFL", message, "Free Agency", int(e.TeamID), ts)
						db.Save(&team)
					}
					db.Save(&e)
				}
			}
		}
	}
}

func GetLatestFreeAgentOfferInDB(db *gorm.DB) uint {
	var latestOffer structs.FreeAgencyOffer

	err := db.Last(&latestOffer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 1
		}
		log.Fatalln("ERROR! Could not find latest record" + err.Error())
	}

	return latestOffer.ID + 1
}

func GetLatestExtensionOfferInDB(db *gorm.DB) uint {
	var latestOffer structs.NFLExtensionOffer

	err := db.Last(&latestOffer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 1
		}
		log.Fatalln("ERROR! Could not find latest record" + err.Error())
	}

	return latestOffer.ID + 1
}

func GetLatestWaiverOfferInDB(db *gorm.DB) uint {
	var latestOffer structs.NFLWaiverOffer

	err := db.Last(&latestOffer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 1
		}
		log.Fatalln("ERROR! Could not find latest record" + err.Error())
	}

	return latestOffer.ID + 1
}

func SetWaiverOrder() {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	nflTeams := GetAllNFLTeams()

	teamMap := make(map[uint]*structs.NFLTeam)

	for i := 0; i < len(nflTeams); i++ {
		teamMap[nflTeams[i].ID] = &nflTeams[i]
	}

	var nflStandings []structs.NFLStandings

	if ts.IsNFLOffSeason || ts.NFLWeek < 3 {
		nflStandings = GetNFLStandingsBySeasonID(strconv.Itoa(int(ts.NFLSeasonID - 1)))
	} else {
		nflStandings = GetNFLStandingsBySeasonID(strconv.Itoa(int(ts.NFLSeasonID)))
	}

	for idx, ns := range nflStandings {
		rank := len(nflStandings) - idx
		teamMap[ns.TeamID].AssignWaiverOrder(uint(rank))
	}

	for _, t := range nflTeams {
		db.Save(&t)
	}
}

func GetWaiverOfferByOfferID(OfferID string) structs.NFLWaiverOffer {
	db := dbprovider.GetInstance().GetDB()

	offer := structs.NFLWaiverOffer{}

	err := db.Where("id = ?", OfferID).Find(&offer).Error
	if err != nil {
		return offer
	}

	return offer
}

func GetWaiverOffersByPlayerID(playerID string) []structs.NFLWaiverOffer {
	db := dbprovider.GetInstance().GetDB()

	offers := []structs.NFLWaiverOffer{}

	err := db.Where("nfl_player_id = ?", playerID).Find(&offers).Error
	if err != nil {
		return offers
	}

	return offers
}

func GetWaiverOffersByTeamID(teamID string) []structs.NFLWaiverOffer {
	db := dbprovider.GetInstance().GetDB()

	offers := []structs.NFLWaiverOffer{}

	err := db.Where("team_id = ?", teamID).Find(&offers).Error
	if err != nil {
		return offers
	}

	return offers
}

func CreateWaiverOffer(offer structs.NFLWaiverOffDTO) structs.NFLWaiverOffer {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	waiverOffer := GetWaiverOfferByOfferID(strconv.Itoa(int(offer.ID)))

	if waiverOffer.ID == 0 {
		id := GetLatestWaiverOfferInDB(db)
		waiverOffer.AssignID(id)
	}

	if ts.IsFreeAgencyLocked {
		return waiverOffer
	}

	waiverOffer.Map(offer)

	db.Save(&waiverOffer)

	fmt.Println("Creating offer!")

	return waiverOffer
}

func CancelWaiverOffer(offer structs.NFLWaiverOffDTO) {
	db := dbprovider.GetInstance().GetDB()

	OfferID := strconv.Itoa(int(offer.ID))

	waiverOffer := GetWaiverOfferByOfferID(OfferID)

	db.Delete(&waiverOffer)
}

func getExtensionPercentageOdds(percentage float64) float64 {
	if percentage >= 100 {
		return 100
	} else if percentage >= 90 {
		return 75
	} else if percentage >= 80 {
		return 50
	} else if percentage >= 70 {
		return 25
	}
	return 0
}

func validateFreeAgencyPref(playerRecord structs.NFLPlayer, roster []structs.NFLPlayer, team structs.NFLTeam, seasonID string, offerLength int, offerIdx int) bool {
	preference := playerRecord.FreeAgency

	if preference == "Average" {
		return true
	}
	if preference == "Drafted team discount" && playerRecord.DraftedTeamID == team.ID {
		return true
	}
	if preference == "Loyal" && (playerRecord.PreviousTeamID == team.ID || playerRecord.TeamID == int(team.ID)) {
		return true
	}

	if preference == "Hometown Hero" && playerRecord.State == team.State {
		return true
	}
	if preference == "Adversarial" && playerRecord.PreviousTeamID != team.ID && playerRecord.DraftedTeamID != team.ID {
		return true
	}

	if preference == "I'm the starter" {
		depthChart := GetNFLDepthchartByTeamID(strconv.Itoa(int(team.ID)))
		dc := depthChart.DepthChartPlayers
		depthChartByPosition := []structs.NFLDepthChartPosition{}

		for _, dcp := range dc {
			if dcp.Position == playerRecord.Position {
				depthChartByPosition = append(depthChartByPosition, dcp)
			}
		}

		sort.Slice(depthChartByPosition, func(i, j int) bool {
			iNum := util.ConvertStringToInt(depthChartByPosition[i].PositionLevel)
			jNum := util.ConvertStringToInt(depthChartByPosition[j].PositionLevel)
			return iNum > jNum
		})
		for idx, p := range depthChartByPosition {
			if idx > 3 {
				return false
			}
			if playerRecord.Overall >= p.NFLPlayer.Overall {
				return true
			}
		}
	}
	if preference == "Market-driven" && offerLength < 3 {
		return true
	}
	if preference == "Wants Extension" && offerLength > 2 {
		return true
	}
	if preference == "Money motivated" {
		return false
	}
	if preference == "Highest bidder" && offerIdx == 0 {
		return true
	}
	if preference == "Championship seeking" {
		standings := GetNFLStandingsByTeamIDAndSeasonID(strconv.Itoa(int(team.ID)), seasonID)
		if standings.TotalWins > standings.TotalLosses {
			return true
		}
	}

	hateBias := strings.Fields(preference)
	if hateBias[0] == "Hates" {
		check := hateCheck(hateBias[1:], team.TeamName)
		return check
	}

	return false
}

// func checkMarketCity(city string) bool {
// 	return city == "Los Angeles" || city == "New York" || city == "New Jersey" || city == "Chicago" || city == "Philadelphia" || city == "Boston" || city == "Dallas" || city == "San Francisco" || city == "Atlanta" || city == "Houston" || city == "Washington"
// }

func hateCheck(bias []string, teamName string) bool {
	joinedBias := strings.Join(bias, " ")
	return joinedBias != teamName
}

func GetRetiredContracts() []structs.NFLContract {
	db := dbprovider.GetInstance().GetDB()

	contracts := []structs.NFLContract{}

	db.Where("player_retired = ?", true).Find(&contracts)

	return contracts
}

func getCapsheetMap() map[uint]structs.NFLCapsheet {
	capsheetMap := make(map[uint]structs.NFLCapsheet)
	var mu sync.Mutex     // to safely update the map
	var wg sync.WaitGroup // to wait for all goroutines to finish
	semaphore := make(chan struct{}, 10)
	nflTeams := GetAllNFLTeams()

	for _, team := range nflTeams {
		semaphore <- struct{}{}
		wg.Add(1)
		go func(t structs.NFLTeam) {
			defer wg.Done()
			capsheet := GetCapsheetByTeamID(strconv.Itoa(int(t.ID)))
			mu.Lock()
			capsheetMap[t.ID] = capsheet
			mu.Unlock()

			<-semaphore
		}(team)
	}

	wg.Wait()
	close(semaphore)
	return capsheetMap
}
