package managers

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/models"
	"github.com/CalebRose/SimFBA/structs"
)

func ProcessUDFAs(isDryRun bool) {
	db := dbprovider.GetInstance().GetDB()

	// 1. Get all Undrafted Players
	var undraftedPlayers []models.NFLDraftee
	db.Where("draft_pick_id = 0").Find(&undraftedPlayers)

	// 2. Get all Bids
	var allBids []structs.NFLUDFAProfile
	db.Find(&allBids)

	// Group bids by PlayerID
	bidsByPlayer := make(map[uint][]structs.NFLUDFAProfile)
	for _, bid := range allBids {
		bidsByPlayer[bid.PlayerID] = append(bidsByPlayer[bid.PlayerID], bid)
	}

	for _, player := range undraftedPlayers {
		bids := bidsByPlayer[player.ID]
		if len(bids) == 0 {
			continue
		}

		// Find winning bid (Highest points)
		var winningBid structs.NFLUDFAProfile
		maxPoints := 0
		var tiedBids []structs.NFLUDFAProfile

		for _, bid := range bids {
			if bid.Points > maxPoints {
				maxPoints = bid.Points
				winningBid = bid
				tiedBids = []structs.NFLUDFAProfile{bid}
			} else if bid.Points == maxPoints {
				tiedBids = append(tiedBids, bid)
			}
		}

		// Tie-breaker: Random Roll
		if len(tiedBids) > 1 {
			rand.Seed(time.Now().UnixNano())
			winningBid = tiedBids[rand.Intn(len(tiedBids))]
		}

		if !isDryRun && winningBid.Points > 0 {
			SignUDFA(player, winningBid)
		} else if isDryRun {
			fmt.Printf("DRY RUN: %s %s would sign with %s for %d points\n", player.FirstName, player.LastName, winningBid.TeamAbbr, winningBid.Points)
		}
	}
}

func SignUDFA(draftee models.NFLDraftee, bid structs.NFLUDFAProfile) {
	db := dbprovider.GetInstance().GetDB()

	// Convert Draftee to NFLPlayer
	nflPlayer := structs.NFLPlayer{
		BasePlayer: draftee.BasePlayer,
		TeamID:     int(bid.TeamID),
		TeamAbbr:   bid.TeamAbbr,
		Experience: 1,
		IsActive:   true,
	}
	nflPlayer.ID = draftee.ID
	db.Create(&nflPlayer)

	// Create 3-year contract: 0.5 Salary, 0 Bonus
	contract := structs.NFLContract{
		NFLPlayerID:    int(nflPlayer.ID),
		TeamID:         bid.TeamID,
		Team:           bid.TeamAbbr,
		ContractLength: 3,
		Y1BaseSalary:   0.5,
		Y2BaseSalary:   0.5,
		Y3BaseSalary:   0.5,
		IsActive:       true,
	}
	db.Create(&contract)

	// Delete from draftee table
	db.Delete(&draftee)
}