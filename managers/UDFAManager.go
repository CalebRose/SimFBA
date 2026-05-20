package managers

import (
	"fmt"
	"strconv"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
)

// ------------------------------------------------------------------------
// USER BOARD MANAGEMENT (The functions your Controller was looking for)
// ------------------------------------------------------------------------

func GetUDFABoardByTeamID(teamID string) structs.NFLUDFABoard {
	db := dbprovider.GetInstance().GetDB()
	var board structs.NFLUDFABoard

	// Find the board and load the players currently on it
	db.Preload("Profiles").Where("team_id = ?", teamID).Find(&board)

	// If the team doesn't have a board yet, create an empty one
	if board.ID == 0 {
		var team structs.NFLTeam
		db.Where("id = ?", teamID).Find(&team)

		board = structs.NFLUDFABoard{
			TeamID:   team.ID,
			TeamAbbr: team.TeamAbbr,
		}
		db.Create(&board)
	}

	return board
}

func AddPlayerToUDFABoard(dto structs.NFLUDFAProfile) {
	db := dbprovider.GetInstance().GetDB()

	// Get the user's board
	board := GetUDFABoardByTeamID(strconv.Itoa(int(dto.TeamID)))

	// Check if this player is already on the board to prevent duplicates
	var existing structs.NFLUDFAProfile
	db.Where("nfl_udfa_board_id = ? AND player_id = ?", board.ID, dto.PlayerID).Find(&existing)

	if existing.ID == 0 {
		dto.NFLUDFABoardID = board.ID
		db.Create(&dto)
	}
}

func SaveUDFABoard(dto structs.NFLUDFABoard) {
	db := dbprovider.GetInstance().GetDB()

	// Loop through the submitted board and update the points for each player
	for _, profile := range dto.Profiles {
		var existing structs.NFLUDFAProfile
		db.Where("id = ?", profile.ID).Find(&existing)

		if existing.ID > 0 {
			existing.Points = profile.Points
			db.Save(&existing)
		}
	}
}

func RemovePlayerFromUDFABoard(profileID string) {
	db := dbprovider.GetInstance().GetDB()
	db.Where("id = ?", profileID).Delete(&structs.NFLUDFAProfile{})
}

// ------------------------------------------------------------------------
// ADMIN BATCH PROCESSING (The logic to actually sign the players)
// ------------------------------------------------------------------------

func ProcessUDFAs(isDryRun bool) {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	// 1. Get all Undrafted Players
	var undraftedPlayers []structs.NFLPlayer
	db.Where("experience = 1 AND is_free_agent = true").Find(&undraftedPlayers)

	// 2. Get all Bids
	var allBids []structs.NFLUDFAProfile
	db.Find(&allBids)

	// Group bids by PlayerID
	bidsByPlayer := make(map[uint][]structs.NFLUDFAProfile)
	for _, bid := range allBids {
		bidsByPlayer[bid.PlayerID] = append(bidsByPlayer[bid.PlayerID], bid)
	}

	// NEW: Create a map to group signings by Team for the Forum Post
	teamSignings := make(map[string][]string)

	for _, player := range undraftedPlayers {
		bids := bidsByPlayer[player.ID]
		if len(bids) == 0 {
			continue // No one bid on this player
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
			randIdx := util.GenerateIntFromRange(0, len(tiedBids)-1)
			winningBid = tiedBids[randIdx]
		}

		// Execute
		if !isDryRun && winningBid.Points > 0 {
			SignUDFA(player, winningBid)

			// NEW: Record the signing for the forum post
			playerString := fmt.Sprintf("%s %s %s", player.Position, player.FirstName, player.LastName)
			teamSignings[winningBid.TeamAbbr] = append(teamSignings[winningBid.TeamAbbr], playerString)

		} else if isDryRun {
			fmt.Printf("DRY RUN: %s %s would sign with %s for %d points\n", player.FirstName, player.LastName, winningBid.TeamAbbr, winningBid.Points)
		}
	}

	// NEW: Generate the Forum Post if it was a live run and players were signed!
	if !isDryRun && len(teamSignings) > 0 {
		// Get the current NFL Season ID from the Timestamp
		// ts := GetTimestamp()

		var forumSignings []string

		for teamAbbr, players := range teamSignings {
			// Add the Team abbreviation as a bolded paragraph line
			forumSignings = append(forumSignings, fmt.Sprintf("**%s**", teamAbbr))

			// Add each player as a bullet point below the team
			for _, p := range players {
				forumSignings = append(forumSignings, fmt.Sprintf("• %s", p))
			}
		}

		// Fire off the automated post to Firebase!
		CreateNFLUDFASyncForumThread(ts.Season, forumSignings)
	}
}

func SignUDFA(draftee structs.NFLPlayer, bid structs.NFLUDFAProfile) {
	db := dbprovider.GetInstance().GetDB()
	draftee.SignPlayer(int(bid.TeamID), bid.TeamAbbr)

	repository.SaveNFLPlayer(draftee, db)

	// Create 3-year contract: 0.5M Salary, 0 Bonus
	contract := structs.NFLContract{
		NFLPlayerID:    int(draftee.ID),
		TeamID:         bid.TeamID,
		Team:           bid.TeamAbbr,
		ContractLength: 3,
		Y1BaseSalary:   0.5,
		Y2BaseSalary:   0.5,
		Y3BaseSalary:   0.5,
		IsActive:       true,
	}
	contract.CalculateContract() // Updates the AAV and total values
	db.Create(&contract)
}
