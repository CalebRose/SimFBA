package managers

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/models"
	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
)

func CFBProgressionMain() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	SeasonID := strconv.Itoa(ts.CollegeSeasonID)
	// Get All Teams
	snapMap := GetCollegePlayerSeasonSnapMap(SeasonID)
	statMap := GetCollegePlayerStatsMap(SeasonID)

	collegeTeams := GetAllCollegeTeams()
	// Loop
	var graduatingPlayers []models.NFLDraftee

	for _, team := range collegeTeams {
		teamID := strconv.Itoa(int(team.ID))
		roster := GetAllCollegePlayersByTeamId(teamID)
		croots := GetSignedRecruitsByTeamProfileID(teamID)

		if !team.PlayersProgressed {
			for _, player := range roster {
				if player.HasProgressed {
					continue
				}
				// Get Latest Stats
				stats := statMap[player.ID]
				snaps := snapMap[player.ID]

				// Get Average Snaps
				avgSnaps := getAverageSnaps(stats)

				// Run Function to Determine if Player is Declaring Early
				willDeclare := DetermineIfDeclaring(player, avgSnaps)

				// Progress the Player
				player = ProgressCollegePlayer(player, SeasonID, stats, snaps)

				if willDeclare {
					player.GraduatePlayer()
					draftee := models.NFLDraftee{}
					draftee.Map(player)
					// Map New Progression value for NFL
					newProgression := util.GenerateNFLPotential(player.Progression)
					newPotentialGrade := util.GetWeightedPotentialGrade(newProgression)
					draftee.MapProgression(newProgression, newPotentialGrade)

					if draftee.Position == "RB" {
						draftee = BoomBustDraftee(draftee, SeasonID, 31, true)
					}

					draftee.GetLetterGrades()

					/*
						Boom/Bust Function
					*/
					tier := 1
					isBoom := false
					enableBoomBust := false
					boomBustStatus := "None"
					tierRoll := util.GenerateIntFromRange(1, 10)
					diceRoll := util.GenerateIntFromRange(1, 20)

					if tierRoll > 7 && tierRoll < 10 {
						tier = 2
					} else if tierRoll > 9 {
						tier = 3
					}

					// Generate Tier
					if diceRoll == 1 && player.ID != 24984 {
						boomBustStatus = "Bust"
						enableBoomBust = true
						// Bust

						draftee.AssignBoomBustStatus(boomBustStatus)

					} else if diceRoll == 20 || player.ID == 24984 {
						enableBoomBust = true
						// Boom

						boomBustStatus = "Boom"
						isBoom = true
						draftee.AssignBoomBustStatus(boomBustStatus)
					} else {
						tier = 0
					}
					if enableBoomBust {
						for i := 0; i < tier; i++ {
							draftee = BoomBustDraftee(draftee, SeasonID, 51, isBoom)
						}
					}

					// Create Historic Player Record
					hcp := (structs.HistoricCollegePlayer)(player)
					repository.CreateHistoricCFBPlayerRecord(hcp, db)

					message := player.Position + " " + player.FirstName + " " + player.LastName + " has graduated from " + player.TeamAbbr + "!"
					if (player.Year < 5 && player.IsRedshirt) || (player.Year < 4 && !player.IsRedshirt) {
						message = "Breaking News! " + player.Position + " " + player.FirstName + " " + player.LastName + " is declaring early from " + player.TeamAbbr + ", and will be eligible for the SimNFL Draft!"
					}

					CreateNewsLog("CFB", message, "Graduation", player.TeamID, ts)

					// Create Draftee Record
					graduatingPlayers = append(graduatingPlayers, draftee)
					// CollegePlayer record will be deleted, but record will be mapped to a GraduatedCollegePlayer struct, and then saved in that table, along side with NFL Draftees table
					// GraduatedCollegePlayer will be a copy of the collegeplayers table, but only for historical players
					repository.DeleteCollegePlayerRecord(player, db)
					continue
				}
				fmt.Println("Saved " + player.FirstName + " " + player.LastName + "'s record")
				repository.SaveCollegePlayerRecord(player, db)
			}

			team.TogglePlayersProgressed()
		}

		if !team.RecruitsAdded {
			for _, croot := range croots {
				// Convert to College Player Record
				cp := structs.CollegePlayer{}
				cp.MapFromRecruit(croot, team)

				// Add in Boom/Bust
				// Tiering only for FCS teams
				tier := 1
				isBoom := false
				enableBoomBust := false
				tierRoll := util.GenerateIntFromRange(1, 10)
				diceRoll := util.GenerateIntFromRange(1, 20)

				if !team.IsFBS && tierRoll > 7 && tierRoll < 10 {
					tier = 2
				} else if !team.IsFBS && tierRoll == 10 {
					tier = 3
				}

				// Generate Tier
				if diceRoll == 1 {

					enableBoomBust = true
				} else if diceRoll == 20 || (cp.ID == 84719 || cp.ID == 84504) {
					enableBoomBust = true
					isBoom = true
				} else {
					tier = 0
				}
				if enableBoomBust {
					for i := 0; i < tier; i++ {
						cp = BoomBustRecruit(cp, SeasonID, 51, isBoom)
					}
				}

				fmt.Println("Adding " + croot.FirstName + " " + croot.LastName + "to " + team.TeamAbbr)

				// Save College Player Record
				repository.CreateCFBPlayerRecord(cp, db)

				// Delete Recruit Record
				repository.DeleteCollegeRecruitRecord(croot, db)
			}

			team.ToggleRecruitsAdded()
		}
		repository.SaveCFBTeam(team, db)
	}

	// Unsigned Players
	unsignedPlayers := GetAllCollegePlayersByTeamId("0")
	for _, player := range unsignedPlayers {
		player = ProgressCollegePlayer(player, SeasonID, []structs.CollegePlayerStats{}, structs.CollegePlayerSeasonSnaps{})
		if (player.IsRedshirt && player.Year > 5) ||
			(!player.IsRedshirt && player.Year > 4) {
			player.GraduatePlayer()
			draftee := models.NFLDraftee{}
			draftee.Map(player)

			// Map New Progression value for NFL
			newProgression := util.GenerateNFLPotential(player.Progression)
			newPotentialGrade := util.GetWeightedPotentialGrade(newProgression)
			draftee.MapProgression(newProgression, newPotentialGrade)

			if draftee.Position == "RB" {
				draftee = BoomBustDraftee(draftee, SeasonID, 31, true)
			}

			draftee.GetLetterGrades()

			/*
				Boom/Bust Function
			*/
			tier := 1
			isBoom := false
			enableBoomBust := false
			boomBustStatus := "None"
			tierRoll := util.GenerateIntFromRange(1, 10)
			diceRoll := util.GenerateIntFromRange(1, 20)

			if tierRoll > 7 && tierRoll < 10 {
				tier = 2
			} else if tierRoll > 9 {
				tier = 3
			}

			// Generate Tier
			if diceRoll == 1 {
				boomBustStatus = "Bust"
				enableBoomBust = true
				// Bust
				fmt.Println("BUST!")
				draftee.AssignBoomBustStatus(boomBustStatus)

			} else if diceRoll == 20 {
				enableBoomBust = true
				// Boom
				fmt.Println("BOOM!")
				boomBustStatus = "Boom"
				isBoom = true
				draftee.AssignBoomBustStatus(boomBustStatus)
			} else {
				tier = 0
			}
			if enableBoomBust {
				for i := 0; i < tier; i++ {
					draftee = BoomBustDraftee(draftee, SeasonID, 51, isBoom)
				}
			}

			hcp := (structs.HistoricCollegePlayer)(player)

			repository.CreateHistoricCFBPlayerRecord(hcp, db)
			graduatingPlayers = append(graduatingPlayers, draftee)
			// CollegePlayer record will be deleted, but record will be mapped to a GraduatedCollegePlayer struct, and then saved in that table, along side with NFL Draftees table
			// GraduatedCollegePlayer will be a copy of the collegeplayers table, but only for historical players

			repository.DeleteCollegePlayerRecord(player, db)
			continue
		}
		repository.SaveCFBPlayer(player, db)
	}

	// Graduating players
	for _, grad := range graduatingPlayers {
		repository.CreateNFLDrafteeRecord(grad, db)
	}
	// get all unsigned players
	// progress through all unsigned players
	// move all seniors + to graduates table
	// move all unsigned croots to unsigned players table

	unsignedCroots := GetAllUnsignedRecruits()
	for _, croot := range unsignedCroots {
		// Unsigned Players
		up := structs.CollegePlayer{}

		up.MapFromRecruit(croot, structs.CollegeTeam{})
		up.WillTransfer()

		tier := 1
		isBoom := false
		enableBoomBust := false
		tierRoll := util.GenerateIntFromRange(1, 10)
		diceRoll := util.GenerateIntFromRange(1, 20)

		if tierRoll > 7 && tierRoll < 10 {
			tier = 2
		} else if tierRoll == 10 {
			tier = 3
		}

		// Generate Tier
		if diceRoll == 1 {

			enableBoomBust = true
		} else if diceRoll == 20 {

			enableBoomBust = true
			isBoom = true
		} else {
			tier = 0
		}
		if enableBoomBust {
			for i := 0; i < tier; i++ {
				up = BoomBustRecruit(up, SeasonID, 31, isBoom)
			}
		}

		repository.CreateCFBPlayerRecord(up, db)
		repository.DeleteCollegeRecruitRecord(croot, db)
	}
}

func NFLProgressionMain() {
	ProgressNFLPlayers()
	AllocateRetiredContracts()
}

func ProgressNFLPlayers() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	SeasonID := strconv.Itoa(ts.CollegeSeasonID)
	fmt.Println(time.Now().UnixNano())
	snapMap := GetNFLPlayerSeasonSnapMap(SeasonID)
	statMap := GetNFLPlayerStatsMap(SeasonID)
	teams := GetAllNFLTeams()
	freeAgents := GetAllFreeAgents()
	lastTwoStatMap := GetNFLLastTwoSeasonStatMap(ts.NFLSeasonID)

	for _, team := range teams {
		teamID := strconv.Itoa(int(team.ID))
		nflPlayers := GetNFLPlayersRecordsByTeamID(teamID)
		for _, player := range nflPlayers {
			if player.HasProgressed {
				continue
			}

			// Progress the Player
			// Get Latest Stats
			stats := statMap[player.ID]
			// Get Average Snaps
			totalSnaps, avgSnaps := getAverageNFLSnaps(stats)
			snaps := snapMap[player.ID]
			// Run Function to Determine if Player is Declaring Early
			willRetire := DetermineIfRetiring(player, lastTwoStatMap)

			// Progress the Player
			player = ProgressNFLPlayer(player, SeasonID, totalSnaps, avgSnaps, snaps)

			playerID := strconv.Itoa(int(player.ID))
			if !player.IsFreeAgent || player.TeamID > 0 {
				activeContract := GetContractByPlayerID(playerID)
				// Get Contract if Applicable
				activeContract.ProgressContract()
				if (activeContract.IsComplete || activeContract.ContractLength == 0) && !willRetire {
					// Reduce contract by one year
					// If contract length == 0, check for extension contract offer
					// if no offer, set player as free agent
					extensions := GetExtensionOffersByPlayerID(playerID)
					acceptedExtension := structs.NFLExtensionOffer{}
					for _, e := range extensions {
						if !e.IsAccepted {
							db.Delete(&e)
							continue
						}
						acceptedExtension = e
						break
					}
					if acceptedExtension.ID > 0 {
						activeContract.MapExtension(acceptedExtension)
						message := "Breaking News: " + player.Position + " " + player.FirstName + " " + player.LastName + " has official signed his extended offer with " + player.TeamAbbr + " for $" + strconv.Itoa(int(activeContract.ContractValue)) + " Million Dollars!"
						CreateNewsLog("NFL", message, "Free Agency", player.TeamID, ts)
						db.Delete(&acceptedExtension)
					} else {
						// No extension offer
						player.ToggleIsFreeAgent()
					}
				}
				if willRetire {
					activeContract.ToggleRetirement()
				}
				repository.SaveNFLContract(activeContract, db)
			}
			if !willRetire {
				player.ToggleHasProgressed()
				repository.SaveNFLPlayer(player, db)
				continue
			}

			// Retire Player
			message := "Breaking News: " + player.Position + " " + player.FirstName + " " + player.LastName + " has decided to retire from SimNFL. He was drafted by " + player.DraftedTeam + " and last played with " + player.TeamAbbr + " and " + player.PreviousTeam + ". We thank him for his wondrous, extensive career and hope he enjoys his retirement!"
			CreateNewsLog("NFL", message, "Retirement", player.TeamID, ts)
			retiredPlayer := (structs.NFLRetiredPlayer)(player)
			repository.CreateRetireeRecord(retiredPlayer, db)
			repository.DeleteNFLPlayerRecord(player, db)
		}
	}

	for _, player := range freeAgents {
		if player.HasProgressed {
			continue
		}

		// Progress the Player
		// Get Latest Stats
		stats := statMap[player.ID]

		// Get Average Snaps
		totalSnaps, avgSnaps := getAverageNFLSnaps(stats)
		snaps := snapMap[player.ID]
		// Run Function to Determine if Player is Declaring Early
		willRetire := DetermineIfRetiring(player, lastTwoStatMap)

		// Progress the Player
		player = ProgressNFLPlayer(player, SeasonID, totalSnaps, avgSnaps, snaps)

		if !willRetire {
			player.ToggleHasProgressed()
			repository.SaveNFLPlayer(player, db)
			continue
		}

		// Retire Player
		message := "Breaking News: " + player.Position + " " + player.FirstName + " " + player.LastName + " has decided to retire from SimNFL. He was drafted by " + player.DraftedTeam + " and last played with " + player.TeamAbbr + " and " + player.PreviousTeam + ". We thank him for his wondrous, extensive career and hope he enjoys his retirement!"
		CreateNewsLog("NFL", message, "Retirement", player.TeamID, ts)
		retiredPlayer := (structs.NFLRetiredPlayer)(player)
		repository.CreateRetireeRecord(retiredPlayer, db)
		repository.DeleteNFLPlayerRecord(player, db)
	}
}

func ProgressNFLPlayer(np structs.NFLPlayer, SeasonID string, totalSnaps, SnapsPerGame int, snaps structs.NFLPlayerSeasonSnaps) structs.NFLPlayer {
	Agility := 0
	ThrowPower := 0
	ThrowAccuracy := 0
	Speed := 0
	FootballIQ := 0
	Strength := 0
	RunBlock := 0
	PassBlock := 0
	RunDefense := 0
	PassRush := 0
	Carrying := 0
	Tackle := 0
	RouteRunning := 0
	Catching := 0
	PuntPower := 0
	PuntAccuracy := 0
	KickPower := 0
	KickAccuracy := 0
	ManCoverage := 0
	ZoneCoverage := 0

	mostPlayedPosition, mostPlayedSnaps := getMostPlayedPosition(snaps.BasePlayerSeasonSnaps, np.Position)
	totalSnaps -= int(snaps.STSnaps)
	posThreshold := float64(totalSnaps) * 0.8

	if mostPlayedPosition != np.Position && float64(mostPlayedSnaps) > posThreshold {
		// Designate New Position
		newArchetype, archCheck := getNewArchetype(np.Position, np.Archetype, mostPlayedPosition)
		// If Archhetype exists, assign new position. Otherwise, progress by old position
		if archCheck {
			np.DesignateNewPosition(mostPlayedPosition, newArchetype)
		} else {
			mostPlayedPosition = np.Position
		}
	} else {
		mostPlayedPosition = np.Position
	}

	/*
		Future work - converting ATH to new primary position
		If a player is an ATH, and if they have played a full position by snaps somewhere, convert the ATH to a the new position, with an archetype that is a fit to the scheme.
		If the position and archetype has no fit with the position, then set the new primary position & select archetype based on the ATH archetype


	*/

	if mostPlayedPosition == "QB" {
		// Primary Progressions
		Agility = PrimaryProgression(np.Progression, np.Agility, mostPlayedPosition, np.Archetype, SnapsPerGame, "Agility", np.IsPracticeSquad)
		ThrowPower = PrimaryProgression(np.Progression, np.ThrowPower, mostPlayedPosition, np.Archetype, SnapsPerGame, "Throw Power", np.IsPracticeSquad)
		ThrowAccuracy = PrimaryProgression(np.Progression, np.ThrowAccuracy, mostPlayedPosition, np.Archetype, SnapsPerGame, "Throw Accuracy", np.IsPracticeSquad)
		Speed = PrimaryProgression(np.Progression, np.Speed, mostPlayedPosition, np.Archetype, SnapsPerGame, "Speed", np.IsPracticeSquad)
		FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, mostPlayedPosition, np.Archetype, SnapsPerGame, "Football IQ", np.IsPracticeSquad)
		Strength = PrimaryProgression(np.Progression, np.Strength, mostPlayedPosition, np.Archetype, SnapsPerGame, "Strength", np.IsPracticeSquad)

		// Secondary Progressions
		RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
		PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
		RunDefense = SecondaryProgression(np.Progression, np.RunDefense)
		PassRush = SecondaryProgression(np.Progression, np.PassRush)
		Carrying = SecondaryProgression(np.Progression, np.Carrying)
		Tackle = SecondaryProgression(np.Progression, np.Tackle)
		RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
		Catching = SecondaryProgression(np.Progression, np.Catching)
		PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
		PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
		KickPower = SecondaryProgression(np.Progression, np.KickPower)
		KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
		ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
		ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
	} else if mostPlayedPosition == "RB" {
		// Primary
		Agility = PrimaryProgression(np.Progression, np.Agility, mostPlayedPosition, np.Archetype, SnapsPerGame, "Agility", np.IsPracticeSquad)
		Carrying = PrimaryProgression(np.Progression, np.Carrying, mostPlayedPosition, np.Archetype, SnapsPerGame, "Carrying", np.IsPracticeSquad)
		Speed = PrimaryProgression(np.Progression, np.Speed, mostPlayedPosition, np.Archetype, SnapsPerGame, "Speed", np.IsPracticeSquad)
		FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, mostPlayedPosition, np.Archetype, SnapsPerGame, "Football IQ", np.IsPracticeSquad)
		Strength = PrimaryProgression(np.Progression, np.Strength, mostPlayedPosition, np.Archetype, SnapsPerGame, "Strength", np.IsPracticeSquad)
		Catching = PrimaryProgression(np.Progression, np.Catching, mostPlayedPosition, np.Archetype, SnapsPerGame, "Catching", np.IsPracticeSquad)
		PassBlock = PrimaryProgression(np.Progression, np.PassBlock, mostPlayedPosition, np.Archetype, SnapsPerGame, "Pass Blocking", np.IsPracticeSquad)
		// Secondary
		RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
		RunDefense = SecondaryProgression(np.Progression, np.RunDefense)
		PassRush = SecondaryProgression(np.Progression, np.PassRush)
		Tackle = SecondaryProgression(np.Progression, np.Tackle)
		RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
		PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
		PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
		KickPower = SecondaryProgression(np.Progression, np.KickPower)
		KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
		ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
		ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
		ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
		ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
	} else if mostPlayedPosition == "FB" {
		// Primary
		Agility = PrimaryProgression(np.Progression, np.Agility, mostPlayedPosition, np.Archetype, SnapsPerGame, "Agility", np.IsPracticeSquad)
		Carrying = PrimaryProgression(np.Progression, np.Carrying, mostPlayedPosition, np.Archetype, SnapsPerGame, "Carrying", np.IsPracticeSquad)
		Speed = PrimaryProgression(np.Progression, np.Speed, mostPlayedPosition, np.Archetype, SnapsPerGame, "Speed", np.IsPracticeSquad)
		FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, mostPlayedPosition, np.Archetype, SnapsPerGame, "Football IQ", np.IsPracticeSquad)
		Strength = PrimaryProgression(np.Progression, np.Strength, mostPlayedPosition, np.Archetype, SnapsPerGame, "Strength", np.IsPracticeSquad)
		Catching = PrimaryProgression(np.Progression, np.Catching, mostPlayedPosition, np.Archetype, SnapsPerGame, "Catching", np.IsPracticeSquad)
		PassBlock = PrimaryProgression(np.Progression, np.PassBlock, mostPlayedPosition, np.Archetype, SnapsPerGame, "Pass Blocking", np.IsPracticeSquad)
		RunBlock = PrimaryProgression(np.Progression, np.RunBlock, mostPlayedPosition, np.Archetype, SnapsPerGame, "Run Blocking", np.IsPracticeSquad)

		// Secondary
		RunDefense = SecondaryProgression(np.Progression, np.RunDefense)
		PassRush = SecondaryProgression(np.Progression, np.PassRush)
		Tackle = SecondaryProgression(np.Progression, np.Tackle)
		RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
		PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
		PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
		KickPower = SecondaryProgression(np.Progression, np.KickPower)
		KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
		ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
		ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
		ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
		ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)

	} else if mostPlayedPosition == "TE" {
		// Primary
		Agility = PrimaryProgression(np.Progression, np.Agility, mostPlayedPosition, np.Archetype, SnapsPerGame, "Agility", np.IsPracticeSquad)
		Carrying = PrimaryProgression(np.Progression, np.Carrying, mostPlayedPosition, np.Archetype, SnapsPerGame, "Carrying", np.IsPracticeSquad)
		Speed = PrimaryProgression(np.Progression, np.Speed, mostPlayedPosition, np.Archetype, SnapsPerGame, "Speed", np.IsPracticeSquad)
		FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, mostPlayedPosition, np.Archetype, SnapsPerGame, "Football IQ", np.IsPracticeSquad)
		Strength = PrimaryProgression(np.Progression, np.Strength, mostPlayedPosition, np.Archetype, SnapsPerGame, "Strength", np.IsPracticeSquad)
		Catching = PrimaryProgression(np.Progression, np.Catching, mostPlayedPosition, np.Archetype, SnapsPerGame, "Catching", np.IsPracticeSquad)
		RouteRunning = PrimaryProgression(np.Progression, np.RouteRunning, mostPlayedPosition, np.Archetype, SnapsPerGame, "Route Running", np.IsPracticeSquad)
		PassBlock = PrimaryProgression(np.Progression, np.PassBlock, mostPlayedPosition, np.Archetype, SnapsPerGame, "Pass Blocking", np.IsPracticeSquad)
		RunBlock = PrimaryProgression(np.Progression, np.RunBlock, mostPlayedPosition, np.Archetype, SnapsPerGame, "Run Blocking", np.IsPracticeSquad)

		// Secondary
		RunDefense = SecondaryProgression(np.Progression, np.RunDefense)
		PassRush = SecondaryProgression(np.Progression, np.PassRush)
		Tackle = SecondaryProgression(np.Progression, np.Tackle)
		PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
		PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
		KickPower = SecondaryProgression(np.Progression, np.KickPower)
		KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
		ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
		ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
		ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
		ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
	} else if mostPlayedPosition == "WR" {
		// Primary
		Agility = PrimaryProgression(np.Progression, np.Agility, mostPlayedPosition, np.Archetype, SnapsPerGame, "Agility", np.IsPracticeSquad)
		Carrying = PrimaryProgression(np.Progression, np.Carrying, mostPlayedPosition, np.Archetype, SnapsPerGame, "Carrying", np.IsPracticeSquad)
		Speed = PrimaryProgression(np.Progression, np.Speed, mostPlayedPosition, np.Archetype, SnapsPerGame, "Speed", np.IsPracticeSquad)
		FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, mostPlayedPosition, np.Archetype, SnapsPerGame, "Football IQ", np.IsPracticeSquad)
		Strength = PrimaryProgression(np.Progression, np.Strength, mostPlayedPosition, np.Archetype, SnapsPerGame, "Strength", np.IsPracticeSquad)
		Catching = PrimaryProgression(np.Progression, np.Catching, mostPlayedPosition, np.Archetype, SnapsPerGame, "Catching", np.IsPracticeSquad)
		RouteRunning = PrimaryProgression(np.Progression, np.RouteRunning, mostPlayedPosition, np.Archetype, SnapsPerGame, "Route Running", np.IsPracticeSquad)

		// Secondary
		RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
		PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
		RunDefense = SecondaryProgression(np.Progression, np.RunDefense)
		PassRush = SecondaryProgression(np.Progression, np.PassRush)
		Tackle = SecondaryProgression(np.Progression, np.Tackle)
		PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
		PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
		KickPower = SecondaryProgression(np.Progression, np.KickPower)
		KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
		ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
		ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
		ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
		ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
	} else if mostPlayedPosition == "OT" || mostPlayedPosition == "OG" || mostPlayedPosition == "C" {
		// Primary
		Agility = PrimaryProgression(np.Progression, np.Agility, mostPlayedPosition, np.Archetype, SnapsPerGame, "Agility", np.IsPracticeSquad)
		FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, mostPlayedPosition, np.Archetype, SnapsPerGame, "Football IQ", np.IsPracticeSquad)
		Strength = PrimaryProgression(np.Progression, np.Strength, mostPlayedPosition, np.Archetype, SnapsPerGame, "Strength", np.IsPracticeSquad)
		PassBlock = PrimaryProgression(np.Progression, np.PassBlock, mostPlayedPosition, np.Archetype, SnapsPerGame, "Pass Blocking", np.IsPracticeSquad)
		RunBlock = PrimaryProgression(np.Progression, np.RunBlock, mostPlayedPosition, np.Archetype, SnapsPerGame, "Run Blocking", np.IsPracticeSquad)

		// Secondary
		RunDefense = SecondaryProgression(np.Progression, np.RunDefense)
		PassRush = SecondaryProgression(np.Progression, np.PassRush)
		Tackle = SecondaryProgression(np.Progression, np.Tackle)
		RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
		PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
		PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
		KickPower = SecondaryProgression(np.Progression, np.KickPower)
		KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
		ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
		ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
		ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
		ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
		Carrying = SecondaryProgression(np.Progression, np.Carrying)
		Speed = SecondaryProgression(np.Progression, np.Speed)
		Catching = SecondaryProgression(np.Progression, np.Catching)
	} else if mostPlayedPosition == "DT" {
		// Primary
		Agility = PrimaryProgression(np.Progression, np.Agility, mostPlayedPosition, np.Archetype, SnapsPerGame, "Agility", np.IsPracticeSquad)
		FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, mostPlayedPosition, np.Archetype, SnapsPerGame, "Football IQ", np.IsPracticeSquad)
		Strength = PrimaryProgression(np.Progression, np.Strength, mostPlayedPosition, np.Archetype, SnapsPerGame, "Strength", np.IsPracticeSquad)
		PassRush = PrimaryProgression(np.Progression, np.PassRush, mostPlayedPosition, np.Archetype, SnapsPerGame, "Pass Rush", np.IsPracticeSquad)
		RunDefense = PrimaryProgression(np.Progression, np.RunDefense, mostPlayedPosition, np.Archetype, SnapsPerGame, "Run Defense", np.IsPracticeSquad)
		Tackle = PrimaryProgression(np.Progression, np.Tackle, mostPlayedPosition, np.Archetype, SnapsPerGame, "Tackle", np.IsPracticeSquad)

		// Secondary
		PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
		PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
		KickPower = SecondaryProgression(np.Progression, np.KickPower)
		KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
		ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
		ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
		ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
		ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
		Carrying = SecondaryProgression(np.Progression, np.Carrying)
		Catching = SecondaryProgression(np.Progression, np.Catching)
		RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
		Speed = SecondaryProgression(np.Progression, np.Speed)
		RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
		PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
	} else if mostPlayedPosition == "DE" {
		// Primary
		Agility = PrimaryProgression(np.Progression, np.Agility, mostPlayedPosition, np.Archetype, SnapsPerGame, "Agility", np.IsPracticeSquad)
		FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, mostPlayedPosition, np.Archetype, SnapsPerGame, "Football IQ", np.IsPracticeSquad)
		Strength = PrimaryProgression(np.Progression, np.Strength, mostPlayedPosition, np.Archetype, SnapsPerGame, "Strength", np.IsPracticeSquad)
		Speed = PrimaryProgression(np.Progression, np.Speed, mostPlayedPosition, np.Archetype, SnapsPerGame, "Speed", np.IsPracticeSquad)
		PassRush = PrimaryProgression(np.Progression, np.PassRush, mostPlayedPosition, np.Archetype, SnapsPerGame, "Pass Rush", np.IsPracticeSquad)
		RunDefense = PrimaryProgression(np.Progression, np.RunDefense, mostPlayedPosition, np.Archetype, SnapsPerGame, "Run Defense", np.IsPracticeSquad)
		Tackle = PrimaryProgression(np.Progression, np.Tackle, mostPlayedPosition, np.Archetype, SnapsPerGame, "Tackle", np.IsPracticeSquad)

		// Secondary
		PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
		PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
		KickPower = SecondaryProgression(np.Progression, np.KickPower)
		KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
		ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
		ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
		ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
		ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
		Carrying = SecondaryProgression(np.Progression, np.Carrying)
		Catching = SecondaryProgression(np.Progression, np.Catching)
		RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
		RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
		PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
	} else if mostPlayedPosition == "OLB" || mostPlayedPosition == "ILB" {
		// Primary
		Agility = PrimaryProgression(np.Progression, np.Agility, mostPlayedPosition, np.Archetype, SnapsPerGame, "Agility", np.IsPracticeSquad)
		FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, mostPlayedPosition, np.Archetype, SnapsPerGame, "Football IQ", np.IsPracticeSquad)
		Strength = PrimaryProgression(np.Progression, np.Strength, mostPlayedPosition, np.Archetype, SnapsPerGame, "Strength", np.IsPracticeSquad)
		Speed = PrimaryProgression(np.Progression, np.Speed, mostPlayedPosition, np.Archetype, SnapsPerGame, "Speed", np.IsPracticeSquad)
		PassRush = PrimaryProgression(np.Progression, np.PassRush, mostPlayedPosition, np.Archetype, SnapsPerGame, "Pass Rush", np.IsPracticeSquad)
		RunDefense = PrimaryProgression(np.Progression, np.RunDefense, mostPlayedPosition, np.Archetype, SnapsPerGame, "Run Defense", np.IsPracticeSquad)
		Tackle = PrimaryProgression(np.Progression, np.Tackle, mostPlayedPosition, np.Archetype, SnapsPerGame, "Tackle", np.IsPracticeSquad)
		ManCoverage = PrimaryProgression(np.Progression, np.ManCoverage, mostPlayedPosition, np.Archetype, SnapsPerGame, "Man Coverage", np.IsPracticeSquad)
		ZoneCoverage = PrimaryProgression(np.Progression, np.ZoneCoverage, mostPlayedPosition, np.Archetype, SnapsPerGame, "Zone Coverage", np.IsPracticeSquad)

		// Secondary
		PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
		PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
		KickPower = SecondaryProgression(np.Progression, np.KickPower)
		KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
		ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
		ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
		Carrying = SecondaryProgression(np.Progression, np.Carrying)
		Catching = SecondaryProgression(np.Progression, np.Catching)
		RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
		RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
		PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
	} else if mostPlayedPosition == "CB" {
		// Primary
		Agility = PrimaryProgression(np.Progression, np.Agility, mostPlayedPosition, np.Archetype, SnapsPerGame, "Agility", np.IsPracticeSquad)
		FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, mostPlayedPosition, np.Archetype, SnapsPerGame, "Football IQ", np.IsPracticeSquad)
		Strength = PrimaryProgression(np.Progression, np.Strength, mostPlayedPosition, np.Archetype, SnapsPerGame, "Strength", np.IsPracticeSquad)
		Speed = PrimaryProgression(np.Progression, np.Speed, mostPlayedPosition, np.Archetype, SnapsPerGame, "Speed", np.IsPracticeSquad)
		Tackle = PrimaryProgression(np.Progression, np.Tackle, mostPlayedPosition, np.Archetype, SnapsPerGame, "Tackle", np.IsPracticeSquad)
		ManCoverage = PrimaryProgression(np.Progression, np.ManCoverage, mostPlayedPosition, np.Archetype, SnapsPerGame, "Man Coverage", np.IsPracticeSquad)
		ZoneCoverage = PrimaryProgression(np.Progression, np.ZoneCoverage, mostPlayedPosition, np.Archetype, SnapsPerGame, "Zone Coverage", np.IsPracticeSquad)
		Catching = PrimaryProgression(np.Progression, np.Catching, mostPlayedPosition, np.Archetype, SnapsPerGame, "Catching", np.IsPracticeSquad)

		// Secondary
		ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
		ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
		RunDefense = SecondaryProgression(np.Progression, np.RunDefense)
		PassRush = SecondaryProgression(np.Progression, np.PassRush)
		PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
		PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
		KickPower = SecondaryProgression(np.Progression, np.KickPower)
		KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
		Carrying = SecondaryProgression(np.Progression, np.Carrying)
		RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
		PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
		RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
	} else if mostPlayedPosition == "FS" || mostPlayedPosition == "SS" {
		// Primary
		Agility = PrimaryProgression(np.Progression, np.Agility, mostPlayedPosition, np.Archetype, SnapsPerGame, "Agility", np.IsPracticeSquad)
		FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, mostPlayedPosition, np.Archetype, SnapsPerGame, "Football IQ", np.IsPracticeSquad)
		Strength = PrimaryProgression(np.Progression, np.Strength, mostPlayedPosition, np.Archetype, SnapsPerGame, "Strength", np.IsPracticeSquad)
		Speed = PrimaryProgression(np.Progression, np.Speed, mostPlayedPosition, np.Archetype, SnapsPerGame, "Speed", np.IsPracticeSquad)
		RunDefense = PrimaryProgression(np.Progression, np.RunDefense, mostPlayedPosition, np.Archetype, SnapsPerGame, "Run Defense", np.IsPracticeSquad)
		Tackle = PrimaryProgression(np.Progression, np.Tackle, mostPlayedPosition, np.Archetype, SnapsPerGame, "Tackle", np.IsPracticeSquad)
		ManCoverage = PrimaryProgression(np.Progression, np.ManCoverage, mostPlayedPosition, np.Archetype, SnapsPerGame, "Man Coverage", np.IsPracticeSquad)
		ZoneCoverage = PrimaryProgression(np.Progression, np.ZoneCoverage, mostPlayedPosition, np.Archetype, SnapsPerGame, "Zone Coverage", np.IsPracticeSquad)
		Catching = PrimaryProgression(np.Progression, np.Catching, mostPlayedPosition, np.Archetype, SnapsPerGame, "Catching", np.IsPracticeSquad)

		// Secondary
		ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
		ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
		PassRush = SecondaryProgression(np.Progression, np.PassRush)
		PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
		PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
		KickPower = SecondaryProgression(np.Progression, np.KickPower)
		KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
		Carrying = SecondaryProgression(np.Progression, np.Carrying)
		RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
		PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
		RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
	} else if mostPlayedPosition == "K" {
		// Primary
		KickPower = PrimaryProgression(np.Progression, np.KickPower, mostPlayedPosition, np.Archetype, SnapsPerGame, "Kick Power", np.IsPracticeSquad)
		KickAccuracy = PrimaryProgression(np.Progression, np.KickAccuracy, mostPlayedPosition, np.Archetype, SnapsPerGame, "Kick Accuracy", np.IsPracticeSquad)
		// Secondary
		ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
		ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
		PassRush = SecondaryProgression(np.Progression, np.PassRush)
		RunDefense = SecondaryProgression(np.Progression, np.RunDefense)
		PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
		PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
		Carrying = SecondaryProgression(np.Progression, np.Carrying)
		RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
		PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
		Catching = SecondaryProgression(np.Progression, np.Catching)
		RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
		Strength = SecondaryProgression(np.Progression, np.Strength)
		Speed = SecondaryProgression(np.Progression, np.Speed)
		Agility = SecondaryProgression(np.Progression, np.Agility)
		ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
		ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
		FootballIQ = SecondaryProgression(np.Progression, np.FootballIQ)
		Tackle = SecondaryProgression(np.Progression, np.Tackle)
	} else if mostPlayedPosition == "P" {
		// If David Ross
		if np.ID == 24984 {
			ThrowPower = PrimaryProgression(np.Progression, np.ThrowPower, mostPlayedPosition, np.Archetype, SnapsPerGame, "Throw Power", np.IsPracticeSquad)
			ThrowAccuracy = PrimaryProgression(np.Progression, np.ThrowAccuracy, mostPlayedPosition, np.Archetype, SnapsPerGame, "Throw Accuracy", np.IsPracticeSquad)
		} else {
			ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
			ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
		}
		// Primary Progressions
		PuntPower = PrimaryProgression(np.Progression, np.PuntPower, mostPlayedPosition, np.Archetype, SnapsPerGame, "Punt Power", np.IsPracticeSquad)
		PuntAccuracy = PrimaryProgression(np.Progression, np.PuntAccuracy, mostPlayedPosition, np.Archetype, SnapsPerGame, "Punt Accuracy", np.IsPracticeSquad)
		// Secondary Progressions
		RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
		PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
		RunDefense = SecondaryProgression(np.Progression, np.RunDefense)
		PassRush = SecondaryProgression(np.Progression, np.PassRush)
		Carrying = SecondaryProgression(np.Progression, np.Carrying)
		Tackle = SecondaryProgression(np.Progression, np.Tackle)
		RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
		Catching = SecondaryProgression(np.Progression, np.Catching)
		KickPower = SecondaryProgression(np.Progression, np.KickPower)
		KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
		ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
		ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
		Strength = SecondaryProgression(np.Progression, np.Strength)
		Speed = SecondaryProgression(np.Progression, np.Speed)
		Agility = SecondaryProgression(np.Progression, np.Agility)
		FootballIQ = SecondaryProgression(np.Progression, np.FootballIQ)
	} else if mostPlayedPosition == "ATH" {
		// Use Archetype
		archetype := np.Archetype
		if archetype == "Field General" {
			Agility = PrimaryProgression(np.Progression, np.Agility, mostPlayedPosition, np.Archetype, SnapsPerGame, "Agility", np.IsPracticeSquad)
			ThrowPower = PrimaryProgression(np.Progression, np.ThrowPower, mostPlayedPosition, np.Archetype, SnapsPerGame, "Throw Power", np.IsPracticeSquad)
			ThrowAccuracy = PrimaryProgression(np.Progression, np.ThrowAccuracy, mostPlayedPosition, np.Archetype, SnapsPerGame, "Throw Accuracy", np.IsPracticeSquad)
			Speed = PrimaryProgression(np.Progression, np.Speed, mostPlayedPosition, np.Archetype, SnapsPerGame, "Speed", np.IsPracticeSquad)
			FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, mostPlayedPosition, np.Archetype, SnapsPerGame, "Football IQ", np.IsPracticeSquad)
			Strength = PrimaryProgression(np.Progression, np.Strength, mostPlayedPosition, np.Archetype, SnapsPerGame, "Strength", np.IsPracticeSquad)
			RunDefense = PrimaryProgression(np.Progression, np.RunDefense, mostPlayedPosition, np.Archetype, SnapsPerGame, "Run Defense", np.IsPracticeSquad)
			Tackle = PrimaryProgression(np.Progression, np.Tackle, mostPlayedPosition, np.Archetype, SnapsPerGame, "Tackle", np.IsPracticeSquad)
			ManCoverage = PrimaryProgression(np.Progression, np.ManCoverage, mostPlayedPosition, np.Archetype, SnapsPerGame, "Man Coverage", np.IsPracticeSquad)
			ZoneCoverage = PrimaryProgression(np.Progression, np.ZoneCoverage, mostPlayedPosition, np.Archetype, SnapsPerGame, "Zone Coverage", np.IsPracticeSquad)

			// Secondary
			PassRush = SecondaryProgression(np.Progression, np.PassRush)
			RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
			PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
			Carrying = SecondaryProgression(np.Progression, np.Carrying)
			RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
			Catching = SecondaryProgression(np.Progression, np.Catching)
			PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
			PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
			KickPower = SecondaryProgression(np.Progression, np.KickPower)
			KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
		} else if archetype == "Triple-Threat" {
			Agility = PrimaryProgression(np.Progression, np.Agility, mostPlayedPosition, np.Archetype, SnapsPerGame, "Agility", np.IsPracticeSquad)
			ThrowPower = PrimaryProgression(np.Progression, np.ThrowPower, mostPlayedPosition, np.Archetype, SnapsPerGame, "Throw Power", np.IsPracticeSquad)
			ThrowAccuracy = PrimaryProgression(np.Progression, np.ThrowAccuracy, mostPlayedPosition, np.Archetype, SnapsPerGame, "Throw Accuracy", np.IsPracticeSquad)
			Speed = PrimaryProgression(np.Progression, np.Speed, mostPlayedPosition, np.Archetype, SnapsPerGame, "Speed", np.IsPracticeSquad)
			FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, mostPlayedPosition, np.Archetype, SnapsPerGame, "Football IQ", np.IsPracticeSquad)
			KickPower = PrimaryProgression(np.Progression, np.KickPower, mostPlayedPosition, np.Archetype, SnapsPerGame, "Kick Power", np.IsPracticeSquad)
			KickAccuracy = PrimaryProgression(np.Progression, np.KickAccuracy, mostPlayedPosition, np.Archetype, SnapsPerGame, "Kick Accuracy", np.IsPracticeSquad)
			PuntPower = PrimaryProgression(np.Progression, np.PuntPower, mostPlayedPosition, np.Archetype, SnapsPerGame, "Punt Power", np.IsPracticeSquad)
			PuntAccuracy = PrimaryProgression(np.Progression, np.PuntAccuracy, mostPlayedPosition, np.Archetype, SnapsPerGame, "Punt Accuracy", np.IsPracticeSquad)
			RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
			Catching = SecondaryProgression(np.Progression, np.Catching)
			RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
			PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
			RunDefense = SecondaryProgression(np.Progression, np.RunDefense)
			PassRush = SecondaryProgression(np.Progression, np.PassRush)
			Tackle = SecondaryProgression(np.Progression, np.Tackle)
			ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
			ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
			Strength = SecondaryProgression(np.Progression, np.Strength)
		} else if archetype == "Wingback" {
			Agility = PrimaryProgression(np.Progression, np.Agility, mostPlayedPosition, np.Archetype, SnapsPerGame, "Agility", np.IsPracticeSquad)
			Speed = PrimaryProgression(np.Progression, np.Speed, mostPlayedPosition, np.Archetype, SnapsPerGame, "Speed", np.IsPracticeSquad)
			FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, mostPlayedPosition, np.Archetype, SnapsPerGame, "Football IQ", np.IsPracticeSquad)
			Catching = PrimaryProgression(np.Progression, np.Catching, mostPlayedPosition, np.Archetype, SnapsPerGame, "Catching", np.IsPracticeSquad)
			RouteRunning = PrimaryProgression(np.Progression, np.RouteRunning, mostPlayedPosition, np.Archetype, SnapsPerGame, "Route Running", np.IsPracticeSquad)
			Carrying = PrimaryProgression(np.Progression, np.Carrying, mostPlayedPosition, np.Archetype, SnapsPerGame, "Carrying", np.IsPracticeSquad)
			RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
			PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
			RunDefense = SecondaryProgression(np.Progression, np.RunDefense)
			PassRush = SecondaryProgression(np.Progression, np.PassRush)
			Tackle = SecondaryProgression(np.Progression, np.Tackle)
			PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
			PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
			KickPower = SecondaryProgression(np.Progression, np.KickPower)
			KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
			ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
			ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
			Strength = SecondaryProgression(np.Progression, np.Strength)
			ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
			ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
		} else if archetype == "Slotback" {
			Agility = PrimaryProgression(np.Progression, np.Agility, mostPlayedPosition, np.Archetype, SnapsPerGame, "Agility", np.IsPracticeSquad)
			Speed = PrimaryProgression(np.Progression, np.Speed, mostPlayedPosition, np.Archetype, SnapsPerGame, "Speed", np.IsPracticeSquad)
			FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, mostPlayedPosition, np.Archetype, SnapsPerGame, "Football IQ", np.IsPracticeSquad)
			Catching = PrimaryProgression(np.Progression, np.Catching, mostPlayedPosition, np.Archetype, SnapsPerGame, "Catching", np.IsPracticeSquad)
			RouteRunning = PrimaryProgression(np.Progression, np.RouteRunning, mostPlayedPosition, np.Archetype, SnapsPerGame, "Route Running", np.IsPracticeSquad)
			PassBlock = PrimaryProgression(np.Progression, np.PassBlock, mostPlayedPosition, np.Archetype, SnapsPerGame, "Pass Blocking", np.IsPracticeSquad)
			RunBlock = PrimaryProgression(np.Progression, np.RunBlock, mostPlayedPosition, np.Archetype, SnapsPerGame, "Run Blocking", np.IsPracticeSquad)
			Strength = PrimaryProgression(np.Progression, np.Strength, mostPlayedPosition, np.Archetype, SnapsPerGame, "Strength", np.IsPracticeSquad)
			Carrying = SecondaryProgression(np.Progression, np.Carrying)
			RunDefense = SecondaryProgression(np.Progression, np.RunDefense)
			PassRush = SecondaryProgression(np.Progression, np.PassRush)
			Tackle = SecondaryProgression(np.Progression, np.Tackle)
			PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
			PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
			KickPower = SecondaryProgression(np.Progression, np.KickPower)
			KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
			ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
			ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
			ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
			ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
		} else if archetype == "Lineman" {
			Agility = PrimaryProgression(np.Progression, np.Agility, mostPlayedPosition, np.Archetype, SnapsPerGame, "Agility", np.IsPracticeSquad)
			FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, mostPlayedPosition, np.Archetype, SnapsPerGame, "Football IQ", np.IsPracticeSquad)
			Strength = PrimaryProgression(np.Progression, np.Strength, mostPlayedPosition, np.Archetype, SnapsPerGame, "Strength", np.IsPracticeSquad)
			PassBlock = PrimaryProgression(np.Progression, np.PassBlock, mostPlayedPosition, np.Archetype, SnapsPerGame, "Pass Blocking", np.IsPracticeSquad)
			RunBlock = PrimaryProgression(np.Progression, np.RunBlock, mostPlayedPosition, np.Archetype, SnapsPerGame, "Run Blocking", np.IsPracticeSquad)
			PassRush = PrimaryProgression(np.Progression, np.PassRush, mostPlayedPosition, np.Archetype, SnapsPerGame, "Pass Rush", np.IsPracticeSquad)
			RunDefense = PrimaryProgression(np.Progression, np.RunDefense, mostPlayedPosition, np.Archetype, SnapsPerGame, "Run Defense", np.IsPracticeSquad)
			Tackle = PrimaryProgression(np.Progression, np.Tackle, mostPlayedPosition, np.Archetype, SnapsPerGame, "Tackle", np.IsPracticeSquad)
			PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
			PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
			KickPower = SecondaryProgression(np.Progression, np.KickPower)
			KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
			ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
			ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
			ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
			ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
			Carrying = SecondaryProgression(np.Progression, np.Carrying)
			Catching = SecondaryProgression(np.Progression, np.Catching)
			RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
			Speed = SecondaryProgression(np.Progression, np.Speed)
		} else if archetype == "Strongside" {
			Agility = PrimaryProgression(np.Progression, np.Agility, mostPlayedPosition, np.Archetype, SnapsPerGame, "Agility", np.IsPracticeSquad)
			FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, mostPlayedPosition, np.Archetype, SnapsPerGame, "Football IQ", np.IsPracticeSquad)
			Strength = PrimaryProgression(np.Progression, np.Strength, mostPlayedPosition, np.Archetype, SnapsPerGame, "Strength", np.IsPracticeSquad)
			Speed = PrimaryProgression(np.Progression, np.Speed, mostPlayedPosition, np.Archetype, SnapsPerGame, "Speed", np.IsPracticeSquad)
			PassRush = PrimaryProgression(np.Progression, np.PassRush, mostPlayedPosition, np.Archetype, SnapsPerGame, "Pass Rush", np.IsPracticeSquad)
			RunDefense = PrimaryProgression(np.Progression, np.RunDefense, mostPlayedPosition, np.Archetype, SnapsPerGame, "Run Defense", np.IsPracticeSquad)
			Tackle = PrimaryProgression(np.Progression, np.Tackle, mostPlayedPosition, np.Archetype, SnapsPerGame, "Tackle", np.IsPracticeSquad)
			ManCoverage = PrimaryProgression(np.Progression, np.ManCoverage, mostPlayedPosition, np.Archetype, SnapsPerGame, "Man Coverage", np.IsPracticeSquad)
			ZoneCoverage = PrimaryProgression(np.Progression, np.ZoneCoverage, mostPlayedPosition, np.Archetype, SnapsPerGame, "Zone Coverage", np.IsPracticeSquad)
			PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
			PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
			KickPower = SecondaryProgression(np.Progression, np.KickPower)
			KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
			ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
			ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
			Carrying = SecondaryProgression(np.Progression, np.Carrying)
			Catching = SecondaryProgression(np.Progression, np.Catching)
			RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
			RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
			PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
		} else if archetype == "Weakside" {
			Agility = PrimaryProgression(np.Progression, np.Agility, mostPlayedPosition, np.Archetype, SnapsPerGame, "Agility", np.IsPracticeSquad)
			FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, mostPlayedPosition, np.Archetype, SnapsPerGame, "Football IQ", np.IsPracticeSquad)
			Strength = PrimaryProgression(np.Progression, np.Strength, mostPlayedPosition, np.Archetype, SnapsPerGame, "Strength", np.IsPracticeSquad)
			Speed = PrimaryProgression(np.Progression, np.Speed, mostPlayedPosition, np.Archetype, SnapsPerGame, "Speed", np.IsPracticeSquad)
			PassRush = PrimaryProgression(np.Progression, np.PassRush, mostPlayedPosition, np.Archetype, SnapsPerGame, "Pass Rush", np.IsPracticeSquad)
			RunDefense = PrimaryProgression(np.Progression, np.RunDefense, mostPlayedPosition, np.Archetype, SnapsPerGame, "Run Defense", np.IsPracticeSquad)
			Tackle = PrimaryProgression(np.Progression, np.Tackle, mostPlayedPosition, np.Archetype, SnapsPerGame, "Tackle", np.IsPracticeSquad)
			ManCoverage = PrimaryProgression(np.Progression, np.ManCoverage, mostPlayedPosition, np.Archetype, SnapsPerGame, "Man Coverage", np.IsPracticeSquad)
			ZoneCoverage = PrimaryProgression(np.Progression, np.ZoneCoverage, mostPlayedPosition, np.Archetype, SnapsPerGame, "Zone Coverage", np.IsPracticeSquad)
			PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
			PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
			KickPower = SecondaryProgression(np.Progression, np.KickPower)
			KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
			ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
			ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
			Carrying = SecondaryProgression(np.Progression, np.Carrying)
			Catching = SecondaryProgression(np.Progression, np.Catching)
			RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
			RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
			PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
		} else if archetype == "Bandit" {
			Agility = PrimaryProgression(np.Progression, np.Agility, mostPlayedPosition, np.Archetype, SnapsPerGame, "Agility", np.IsPracticeSquad)
			FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, mostPlayedPosition, np.Archetype, SnapsPerGame, "Football IQ", np.IsPracticeSquad)
			Strength = PrimaryProgression(np.Progression, np.Strength, mostPlayedPosition, np.Archetype, SnapsPerGame, "Strength", np.IsPracticeSquad)
			Speed = PrimaryProgression(np.Progression, np.Speed, mostPlayedPosition, np.Archetype, SnapsPerGame, "Speed", np.IsPracticeSquad)
			PassRush = PrimaryProgression(np.Progression, np.PassRush, mostPlayedPosition, np.Archetype, SnapsPerGame, "Pass Rush", np.IsPracticeSquad)
			RunDefense = PrimaryProgression(np.Progression, np.RunDefense, mostPlayedPosition, np.Archetype, SnapsPerGame, "Run Defense", np.IsPracticeSquad)
			Tackle = PrimaryProgression(np.Progression, np.Tackle, mostPlayedPosition, np.Archetype, SnapsPerGame, "Tackle", np.IsPracticeSquad)
			ManCoverage = PrimaryProgression(np.Progression, np.ManCoverage, mostPlayedPosition, np.Archetype, SnapsPerGame, "Man Coverage", np.IsPracticeSquad)
			ZoneCoverage = PrimaryProgression(np.Progression, np.ZoneCoverage, mostPlayedPosition, np.Archetype, SnapsPerGame, "Zone Coverage", np.IsPracticeSquad)
			PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
			PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
			KickPower = SecondaryProgression(np.Progression, np.KickPower)
			KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
			ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
			ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
			Carrying = SecondaryProgression(np.Progression, np.Carrying)
			Catching = SecondaryProgression(np.Progression, np.Catching)
			RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
			RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
			PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
		} else if archetype == "Return Specialist" {
			Agility = PrimaryProgression(np.Progression, np.Agility, mostPlayedPosition, np.Archetype, SnapsPerGame, "Agility", np.IsPracticeSquad)
			Speed = PrimaryProgression(np.Progression, np.Speed, mostPlayedPosition, np.Archetype, SnapsPerGame, "Speed", np.IsPracticeSquad)
			Carrying = PrimaryProgression(np.Progression, np.Carrying, mostPlayedPosition, np.Archetype, SnapsPerGame, "Carrying", np.IsPracticeSquad)
			Catching = PrimaryProgression(np.Progression, np.Catching, mostPlayedPosition, np.Archetype, SnapsPerGame, "Catching", np.IsPracticeSquad)
			RouteRunning = PrimaryProgression(np.Progression, np.RouteRunning, mostPlayedPosition, np.Archetype, SnapsPerGame, "Route Running", np.IsPracticeSquad)
			FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, mostPlayedPosition, np.Archetype, SnapsPerGame, "Football IQ", np.IsPracticeSquad)
			PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
			PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
			KickPower = SecondaryProgression(np.Progression, np.KickPower)
			KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
			ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
			ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
			RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
			PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
			Strength = SecondaryProgression(np.Progression, np.Strength)
		} else if archetype == "Soccer Player" {
			KickPower = PrimaryProgression(np.Progression, np.KickPower, mostPlayedPosition, np.Archetype, SnapsPerGame, "Kick Power", np.IsPracticeSquad)
			KickAccuracy = PrimaryProgression(np.Progression, np.KickAccuracy, mostPlayedPosition, np.Archetype, SnapsPerGame, "Kick Accuracy", np.IsPracticeSquad)
			PuntPower = PrimaryProgression(np.Progression, np.PuntPower, mostPlayedPosition, np.Archetype, SnapsPerGame, "Punt Power", np.IsPracticeSquad)
			PuntAccuracy = PrimaryProgression(np.Progression, np.PuntAccuracy, mostPlayedPosition, np.Archetype, SnapsPerGame, "Punt Accuracy", np.IsPracticeSquad)
			Agility = PrimaryProgression(np.Progression, np.Agility, mostPlayedPosition, np.Archetype, SnapsPerGame, "Agility", np.IsPracticeSquad)
			Speed = PrimaryProgression(np.Progression, np.Speed, mostPlayedPosition, np.Archetype, SnapsPerGame, "Speed", np.IsPracticeSquad)
			FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, mostPlayedPosition, np.Archetype, SnapsPerGame, "Football IQ", np.IsPracticeSquad)
			Catching = PrimaryProgression(np.Progression, np.Catching, mostPlayedPosition, np.Archetype, SnapsPerGame, "Catching", np.IsPracticeSquad)
			RouteRunning = PrimaryProgression(np.Progression, np.RouteRunning, mostPlayedPosition, np.Archetype, SnapsPerGame, "Route Running", np.IsPracticeSquad)
			Strength = SecondaryProgression(np.Progression, np.Strength)
			RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
			PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
			RunDefense = SecondaryProgression(np.Progression, np.RunDefense)
			PassRush = SecondaryProgression(np.Progression, np.PassRush)
			Tackle = SecondaryProgression(np.Progression, np.Tackle)
			ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
			ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
			ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
			ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
			Carrying = SecondaryProgression(np.Progression, np.Carrying)
		}
	}

	ThrowPower = RegressAttribute(ThrowPower, np.Age, int(np.PrimeAge), false)
	ThrowAccuracy = RegressAttribute(ThrowAccuracy, np.Age, int(np.PrimeAge), false)
	RunBlock = RegressAttribute(RunBlock, np.Age, int(np.PrimeAge), false)
	PassBlock = RegressAttribute(PassBlock, np.Age, int(np.PrimeAge), false)
	RunDefense = RegressAttribute(RunDefense, np.Age, int(np.PrimeAge), false)
	PassRush = RegressAttribute(PassRush, np.Age, int(np.PrimeAge), false)
	Carrying = RegressAttribute(Carrying, np.Age, int(np.PrimeAge), false)
	Tackle = RegressAttribute(Tackle, np.Age, int(np.PrimeAge), false)
	RouteRunning = RegressAttribute(RouteRunning, np.Age, int(np.PrimeAge), false)
	Catching = RegressAttribute(Catching, np.Age, int(np.PrimeAge), false)
	KickPower = RegressAttribute(KickPower, np.Age, int(np.PrimeAge), false)
	KickAccuracy = RegressAttribute(KickAccuracy, np.Age, int(np.PrimeAge), false)
	ManCoverage = RegressAttribute(ManCoverage, np.Age, int(np.PrimeAge), false)
	ZoneCoverage = RegressAttribute(ZoneCoverage, np.Age, int(np.PrimeAge), false)
	Strength = RegressAttribute(Strength, np.Age, int(np.PrimeAge), false)
	Speed = RegressAttribute(Speed, np.Age, int(np.PrimeAge), false)
	Agility = RegressAttribute(Agility, np.Age, int(np.PrimeAge), false)
	FootballIQ = RegressAttribute(FootballIQ, np.Age, int(np.PrimeAge), true)
	PuntPower = RegressAttribute(PuntPower, np.Age, int(np.PrimeAge), false)
	PuntAccuracy = RegressAttribute(PuntAccuracy, np.Age, int(np.PrimeAge), false)
	newPotentialGrade := util.GetWeightedPotentialGrade(np.Progression)

	progressions := structs.CollegePlayerProgressions{
		Agility:        Agility,
		Catching:       Catching,
		Carrying:       Carrying,
		Speed:          Speed,
		RouteRunning:   RouteRunning,
		RunBlock:       RunBlock,
		PassBlock:      PassBlock,
		RunDefense:     RunDefense,
		PassRush:       PassRush,
		Strength:       Strength,
		Tackle:         Tackle,
		ThrowPower:     ThrowPower,
		ThrowAccuracy:  ThrowAccuracy,
		PuntAccuracy:   PuntAccuracy,
		PuntPower:      PuntPower,
		KickAccuracy:   KickAccuracy,
		KickPower:      KickPower,
		FootballIQ:     FootballIQ,
		ManCoverage:    ManCoverage,
		ZoneCoverage:   ZoneCoverage,
		PotentialGrade: newPotentialGrade,
	}

	np.Progress(progressions)

	np.GetOverall()

	return np
}

func ProgressCollegePlayer(cp structs.CollegePlayer, SeasonID string, stats []structs.CollegePlayerStats, snaps structs.CollegePlayerSeasonSnaps) structs.CollegePlayer {
	totalSnaps := 0

	for _, stat := range stats {
		totalSnaps += stat.Snaps
	}

	var SnapsPerGame int = 0
	if len(stats) > 0 {
		SnapsPerGame = totalSnaps / 12 // 12
	}

	Agility := 0
	ThrowPower := 0
	ThrowAccuracy := 0
	Speed := 0
	FootballIQ := 0
	Strength := 0
	RunBlock := 0
	PassBlock := 0
	RunDefense := 0
	PassRush := 0
	Carrying := 0
	Tackle := 0
	RouteRunning := 0
	Catching := 0
	PuntPower := 0
	PuntAccuracy := 0
	KickPower := 0
	KickAccuracy := 0
	ManCoverage := 0
	ZoneCoverage := 0

	// Get most played position
	mostPlayedPosition, mostPlayedSnaps := getMostPlayedPosition(snaps.BasePlayerSeasonSnaps, cp.Position)
	totalSnaps -= int(snaps.STSnaps)
	posThreshold := float64(totalSnaps) * 0.8

	if mostPlayedPosition != cp.Position && float64(mostPlayedSnaps) > posThreshold {
		// Designate New Position
		newArchetype, archCheck := getNewArchetype(cp.Position, cp.Archetype, mostPlayedPosition)
		// If Archhetype exists, assign new position. Otherwise, progress by old position
		if archCheck {
			cp.DesignateNewPosition(mostPlayedPosition, newArchetype)
		} else {
			mostPlayedPosition = cp.Position
		}
	} else {
		mostPlayedPosition = cp.Position
	}

	if mostPlayedPosition == "QB" {
		// Primary Progressions
		Agility = PrimaryProgression(cp.Progression, cp.Agility, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Agility", cp.IsRedshirting)
		ThrowPower = PrimaryProgression(cp.Progression, cp.ThrowPower, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Throw Power", cp.IsRedshirting)
		ThrowAccuracy = PrimaryProgression(cp.Progression, cp.ThrowAccuracy, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Throw Accuracy", cp.IsRedshirting)
		Speed = PrimaryProgression(cp.Progression, cp.Speed, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Speed", cp.IsRedshirting)
		FootballIQ = PrimaryProgression(cp.Progression, cp.FootballIQ, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Football IQ", cp.IsRedshirting)
		Strength = PrimaryProgression(cp.Progression, cp.Strength, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Strength", cp.IsRedshirting)

		// Secondary Progressions
		RunBlock = SecondaryProgression(cp.Progression, cp.RunBlock)
		PassBlock = SecondaryProgression(cp.Progression, cp.PassBlock)
		RunDefense = SecondaryProgression(cp.Progression, cp.RunDefense)
		PassRush = SecondaryProgression(cp.Progression, cp.PassRush)
		Carrying = SecondaryProgression(cp.Progression, cp.Carrying)
		Tackle = SecondaryProgression(cp.Progression, cp.Tackle)
		RouteRunning = SecondaryProgression(cp.RouteRunning, cp.RouteRunning)
		Catching = SecondaryProgression(cp.Progression, cp.Catching)
		PuntPower = SecondaryProgression(cp.Progression, cp.PuntPower)
		PuntAccuracy = SecondaryProgression(cp.Progression, cp.PuntAccuracy)
		KickPower = SecondaryProgression(cp.Progression, cp.KickPower)
		KickAccuracy = SecondaryProgression(cp.Progression, cp.KickAccuracy)
		ManCoverage = SecondaryProgression(cp.Progression, cp.ManCoverage)
		ZoneCoverage = SecondaryProgression(cp.Progression, cp.ZoneCoverage)
	} else if mostPlayedPosition == "RB" {
		// Primary
		Agility = PrimaryProgression(cp.Progression, cp.Agility, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Agility", cp.IsRedshirting)
		Carrying = PrimaryProgression(cp.Progression, cp.Carrying, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Carrying", cp.IsRedshirting)
		Speed = PrimaryProgression(cp.Progression, cp.Speed, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Speed", cp.IsRedshirting)
		FootballIQ = PrimaryProgression(cp.Progression, cp.FootballIQ, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Football IQ", cp.IsRedshirting)
		Strength = PrimaryProgression(cp.Progression, cp.Strength, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Strength", cp.IsRedshirting)
		Catching = PrimaryProgression(cp.Progression, cp.Catching, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Catching", cp.IsRedshirting)
		PassBlock = PrimaryProgression(cp.Progression, cp.PassBlock, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Pass Blocking", cp.IsRedshirting)
		// Secondary
		RunBlock = SecondaryProgression(cp.Progression, cp.RunBlock)
		RunDefense = SecondaryProgression(cp.Progression, cp.RunDefense)
		PassRush = SecondaryProgression(cp.Progression, cp.PassRush)
		Tackle = SecondaryProgression(cp.Progression, cp.Tackle)
		RouteRunning = SecondaryProgression(cp.RouteRunning, cp.RouteRunning)
		PuntPower = SecondaryProgression(cp.Progression, cp.PuntPower)
		PuntAccuracy = SecondaryProgression(cp.Progression, cp.PuntAccuracy)
		KickPower = SecondaryProgression(cp.Progression, cp.KickPower)
		KickAccuracy = SecondaryProgression(cp.Progression, cp.KickAccuracy)
		ManCoverage = SecondaryProgression(cp.Progression, cp.ManCoverage)
		ZoneCoverage = SecondaryProgression(cp.Progression, cp.ZoneCoverage)
		ThrowPower = SecondaryProgression(cp.Progression, cp.ThrowPower)
		ThrowAccuracy = SecondaryProgression(cp.Progression, cp.ThrowAccuracy)
	} else if mostPlayedPosition == "FB" {
		// Primary
		Agility = PrimaryProgression(cp.Progression, cp.Agility, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Agility", cp.IsRedshirting)
		Carrying = PrimaryProgression(cp.Progression, cp.Carrying, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Carrying", cp.IsRedshirting)
		Speed = PrimaryProgression(cp.Progression, cp.Speed, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Speed", cp.IsRedshirting)
		FootballIQ = PrimaryProgression(cp.Progression, cp.FootballIQ, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Football IQ", cp.IsRedshirting)
		Strength = PrimaryProgression(cp.Progression, cp.Strength, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Strength", cp.IsRedshirting)
		Catching = PrimaryProgression(cp.Progression, cp.Catching, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Catching", cp.IsRedshirting)
		PassBlock = PrimaryProgression(cp.Progression, cp.PassBlock, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Pass Blocking", cp.IsRedshirting)
		RunBlock = PrimaryProgression(cp.Progression, cp.RunBlock, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Run Blocking", cp.IsRedshirting)

		// Secondary
		RunDefense = SecondaryProgression(cp.Progression, cp.RunDefense)
		PassRush = SecondaryProgression(cp.Progression, cp.PassRush)
		Tackle = SecondaryProgression(cp.Progression, cp.Tackle)
		RouteRunning = SecondaryProgression(cp.RouteRunning, cp.RouteRunning)
		PuntPower = SecondaryProgression(cp.Progression, cp.PuntPower)
		PuntAccuracy = SecondaryProgression(cp.Progression, cp.PuntAccuracy)
		KickPower = SecondaryProgression(cp.Progression, cp.KickPower)
		KickAccuracy = SecondaryProgression(cp.Progression, cp.KickAccuracy)
		ManCoverage = SecondaryProgression(cp.Progression, cp.ManCoverage)
		ZoneCoverage = SecondaryProgression(cp.Progression, cp.ZoneCoverage)
		ThrowPower = SecondaryProgression(cp.Progression, cp.ThrowPower)
		ThrowAccuracy = SecondaryProgression(cp.Progression, cp.ThrowAccuracy)

	} else if mostPlayedPosition == "TE" {
		// Primary
		Agility = PrimaryProgression(cp.Progression, cp.Agility, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Agility", cp.IsRedshirting)
		Carrying = PrimaryProgression(cp.Progression, cp.Carrying, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Carrying", cp.IsRedshirting)
		Speed = PrimaryProgression(cp.Progression, cp.Speed, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Speed", cp.IsRedshirting)
		FootballIQ = PrimaryProgression(cp.Progression, cp.FootballIQ, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Football IQ", cp.IsRedshirting)
		Strength = PrimaryProgression(cp.Progression, cp.Strength, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Strength", cp.IsRedshirting)
		Catching = PrimaryProgression(cp.Progression, cp.Catching, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Catching", cp.IsRedshirting)
		RouteRunning = PrimaryProgression(cp.Progression, cp.RouteRunning, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Route Running", cp.IsRedshirting)
		PassBlock = PrimaryProgression(cp.Progression, cp.PassBlock, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Pass Blocking", cp.IsRedshirting)
		RunBlock = PrimaryProgression(cp.Progression, cp.RunBlock, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Run Blocking", cp.IsRedshirting)

		// Secondary
		RunDefense = SecondaryProgression(cp.Progression, cp.RunDefense)
		PassRush = SecondaryProgression(cp.Progression, cp.PassRush)
		Tackle = SecondaryProgression(cp.Progression, cp.Tackle)
		PuntPower = SecondaryProgression(cp.Progression, cp.PuntPower)
		PuntAccuracy = SecondaryProgression(cp.Progression, cp.PuntAccuracy)
		KickPower = SecondaryProgression(cp.Progression, cp.KickPower)
		KickAccuracy = SecondaryProgression(cp.Progression, cp.KickAccuracy)
		ManCoverage = SecondaryProgression(cp.Progression, cp.ManCoverage)
		ZoneCoverage = SecondaryProgression(cp.Progression, cp.ZoneCoverage)
		ThrowPower = SecondaryProgression(cp.Progression, cp.ThrowPower)
		ThrowAccuracy = SecondaryProgression(cp.Progression, cp.ThrowAccuracy)
	} else if mostPlayedPosition == "WR" {
		// Primary
		Agility = PrimaryProgression(cp.Progression, cp.Agility, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Agility", cp.IsRedshirting)
		Carrying = PrimaryProgression(cp.Progression, cp.Carrying, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Carrying", cp.IsRedshirting)
		Speed = PrimaryProgression(cp.Progression, cp.Speed, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Speed", cp.IsRedshirting)
		FootballIQ = PrimaryProgression(cp.Progression, cp.FootballIQ, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Football IQ", cp.IsRedshirting)
		Strength = PrimaryProgression(cp.Progression, cp.Strength, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Strength", cp.IsRedshirting)
		Catching = PrimaryProgression(cp.Progression, cp.Catching, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Catching", cp.IsRedshirting)
		RouteRunning = PrimaryProgression(cp.Progression, cp.RouteRunning, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Route Running", cp.IsRedshirting)

		// Secondary
		RunBlock = SecondaryProgression(cp.Progression, cp.RunBlock)
		PassBlock = SecondaryProgression(cp.Progression, cp.PassBlock)
		RunDefense = SecondaryProgression(cp.Progression, cp.RunDefense)
		PassRush = SecondaryProgression(cp.Progression, cp.PassRush)
		Tackle = SecondaryProgression(cp.Progression, cp.Tackle)
		PuntPower = SecondaryProgression(cp.Progression, cp.PuntPower)
		PuntAccuracy = SecondaryProgression(cp.Progression, cp.PuntAccuracy)
		KickPower = SecondaryProgression(cp.Progression, cp.KickPower)
		KickAccuracy = SecondaryProgression(cp.Progression, cp.KickAccuracy)
		ManCoverage = SecondaryProgression(cp.Progression, cp.ManCoverage)
		ZoneCoverage = SecondaryProgression(cp.Progression, cp.ZoneCoverage)
		ThrowPower = SecondaryProgression(cp.Progression, cp.ThrowPower)
		ThrowAccuracy = SecondaryProgression(cp.Progression, cp.ThrowAccuracy)
	} else if mostPlayedPosition == "OT" || mostPlayedPosition == "OG" || mostPlayedPosition == "C" {
		// Primary
		Agility = PrimaryProgression(cp.Progression, cp.Agility, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Agility", cp.IsRedshirting)
		FootballIQ = PrimaryProgression(cp.Progression, cp.FootballIQ, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Football IQ", cp.IsRedshirting)
		Strength = PrimaryProgression(cp.Progression, cp.Strength, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Strength", cp.IsRedshirting)
		PassBlock = PrimaryProgression(cp.Progression, cp.PassBlock, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Pass Blocking", cp.IsRedshirting)
		RunBlock = PrimaryProgression(cp.Progression, cp.RunBlock, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Run Blocking", cp.IsRedshirting)

		// Secondary
		RunDefense = SecondaryProgression(cp.Progression, cp.RunDefense)
		PassRush = SecondaryProgression(cp.Progression, cp.PassRush)
		Tackle = SecondaryProgression(cp.Progression, cp.Tackle)
		RouteRunning = SecondaryProgression(cp.RouteRunning, cp.RouteRunning)
		PuntPower = SecondaryProgression(cp.Progression, cp.PuntPower)
		PuntAccuracy = SecondaryProgression(cp.Progression, cp.PuntAccuracy)
		KickPower = SecondaryProgression(cp.Progression, cp.KickPower)
		KickAccuracy = SecondaryProgression(cp.Progression, cp.KickAccuracy)
		ManCoverage = SecondaryProgression(cp.Progression, cp.ManCoverage)
		ZoneCoverage = SecondaryProgression(cp.Progression, cp.ZoneCoverage)
		ThrowPower = SecondaryProgression(cp.Progression, cp.ThrowPower)
		ThrowAccuracy = SecondaryProgression(cp.Progression, cp.ThrowAccuracy)
		Carrying = SecondaryProgression(cp.Progression, cp.Carrying)
		Speed = SecondaryProgression(cp.Progression, cp.Speed)
		Catching = SecondaryProgression(cp.Progression, cp.Catching)
	} else if mostPlayedPosition == "DT" {
		// Primary
		Agility = PrimaryProgression(cp.Progression, cp.Agility, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Agility", cp.IsRedshirting)
		FootballIQ = PrimaryProgression(cp.Progression, cp.FootballIQ, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Football IQ", cp.IsRedshirting)
		Strength = PrimaryProgression(cp.Progression, cp.Strength, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Strength", cp.IsRedshirting)
		PassRush = PrimaryProgression(cp.Progression, cp.PassRush, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Pass Rush", cp.IsRedshirting)
		RunDefense = PrimaryProgression(cp.Progression, cp.RunDefense, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Run Defense", cp.IsRedshirting)
		Tackle = PrimaryProgression(cp.Progression, cp.Tackle, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Tackle", cp.IsRedshirting)

		// Secondary
		PuntPower = SecondaryProgression(cp.Progression, cp.PuntPower)
		PuntAccuracy = SecondaryProgression(cp.Progression, cp.PuntAccuracy)
		KickPower = SecondaryProgression(cp.Progression, cp.KickPower)
		KickAccuracy = SecondaryProgression(cp.Progression, cp.KickAccuracy)
		ManCoverage = SecondaryProgression(cp.Progression, cp.ManCoverage)
		ZoneCoverage = SecondaryProgression(cp.Progression, cp.ZoneCoverage)
		ThrowPower = SecondaryProgression(cp.Progression, cp.ThrowPower)
		ThrowAccuracy = SecondaryProgression(cp.Progression, cp.ThrowAccuracy)
		Carrying = SecondaryProgression(cp.Progression, cp.Carrying)
		Catching = SecondaryProgression(cp.Progression, cp.Catching)
		RouteRunning = SecondaryProgression(cp.RouteRunning, cp.RouteRunning)
		Speed = SecondaryProgression(cp.Progression, cp.Speed)
		RunBlock = SecondaryProgression(cp.Progression, cp.RunBlock)
		PassBlock = SecondaryProgression(cp.Progression, cp.PassBlock)
	} else if mostPlayedPosition == "DE" {
		// Primary
		Agility = PrimaryProgression(cp.Progression, cp.Agility, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Agility", cp.IsRedshirting)
		FootballIQ = PrimaryProgression(cp.Progression, cp.FootballIQ, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Football IQ", cp.IsRedshirting)
		Strength = PrimaryProgression(cp.Progression, cp.Strength, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Strength", cp.IsRedshirting)
		Speed = PrimaryProgression(cp.Progression, cp.Speed, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Speed", cp.IsRedshirting)
		PassRush = PrimaryProgression(cp.Progression, cp.PassRush, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Pass Rush", cp.IsRedshirting)
		RunDefense = PrimaryProgression(cp.Progression, cp.RunDefense, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Run Defense", cp.IsRedshirting)
		Tackle = PrimaryProgression(cp.Progression, cp.Tackle, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Tackle", cp.IsRedshirting)

		// Secondary
		PuntPower = SecondaryProgression(cp.Progression, cp.PuntPower)
		PuntAccuracy = SecondaryProgression(cp.Progression, cp.PuntAccuracy)
		KickPower = SecondaryProgression(cp.Progression, cp.KickPower)
		KickAccuracy = SecondaryProgression(cp.Progression, cp.KickAccuracy)
		ManCoverage = SecondaryProgression(cp.Progression, cp.ManCoverage)
		ZoneCoverage = SecondaryProgression(cp.Progression, cp.ZoneCoverage)
		ThrowPower = SecondaryProgression(cp.Progression, cp.ThrowPower)
		ThrowAccuracy = SecondaryProgression(cp.Progression, cp.ThrowAccuracy)
		Carrying = SecondaryProgression(cp.Progression, cp.Carrying)
		Catching = SecondaryProgression(cp.Progression, cp.Catching)
		RouteRunning = SecondaryProgression(cp.RouteRunning, cp.RouteRunning)
		RunBlock = SecondaryProgression(cp.Progression, cp.RunBlock)
		PassBlock = SecondaryProgression(cp.Progression, cp.PassBlock)
	} else if mostPlayedPosition == "OLB" || mostPlayedPosition == "ILB" {
		// Primary
		Agility = PrimaryProgression(cp.Progression, cp.Agility, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Agility", cp.IsRedshirting)
		FootballIQ = PrimaryProgression(cp.Progression, cp.FootballIQ, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Football IQ", cp.IsRedshirting)
		Strength = PrimaryProgression(cp.Progression, cp.Strength, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Strength", cp.IsRedshirting)
		Speed = PrimaryProgression(cp.Progression, cp.Speed, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Speed", cp.IsRedshirting)
		PassRush = PrimaryProgression(cp.Progression, cp.PassRush, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Pass Rush", cp.IsRedshirting)
		RunDefense = PrimaryProgression(cp.Progression, cp.RunDefense, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Run Defense", cp.IsRedshirting)
		Tackle = PrimaryProgression(cp.Progression, cp.Tackle, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Tackle", cp.IsRedshirting)
		ManCoverage = PrimaryProgression(cp.Progression, cp.ManCoverage, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Man Coverage", cp.IsRedshirting)
		ZoneCoverage = PrimaryProgression(cp.Progression, cp.ZoneCoverage, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Zone Coverage", cp.IsRedshirting)

		// Secondary
		PuntPower = SecondaryProgression(cp.Progression, cp.PuntPower)
		PuntAccuracy = SecondaryProgression(cp.Progression, cp.PuntAccuracy)
		KickPower = SecondaryProgression(cp.Progression, cp.KickPower)
		KickAccuracy = SecondaryProgression(cp.Progression, cp.KickAccuracy)
		ThrowPower = SecondaryProgression(cp.Progression, cp.ThrowPower)
		ThrowAccuracy = SecondaryProgression(cp.Progression, cp.ThrowAccuracy)
		Carrying = SecondaryProgression(cp.Progression, cp.Carrying)
		Catching = SecondaryProgression(cp.Progression, cp.Catching)
		RouteRunning = SecondaryProgression(cp.RouteRunning, cp.RouteRunning)
		RunBlock = SecondaryProgression(cp.Progression, cp.RunBlock)
		PassBlock = SecondaryProgression(cp.Progression, cp.PassBlock)
	} else if mostPlayedPosition == "CB" {
		// Primary
		Agility = PrimaryProgression(cp.Progression, cp.Agility, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Agility", cp.IsRedshirting)
		FootballIQ = PrimaryProgression(cp.Progression, cp.FootballIQ, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Football IQ", cp.IsRedshirting)
		Strength = PrimaryProgression(cp.Progression, cp.Strength, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Strength", cp.IsRedshirting)
		Speed = PrimaryProgression(cp.Progression, cp.Speed, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Speed", cp.IsRedshirting)
		Tackle = PrimaryProgression(cp.Progression, cp.Tackle, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Tackle", cp.IsRedshirting)
		ManCoverage = PrimaryProgression(cp.Progression, cp.ManCoverage, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Man Coverage", cp.IsRedshirting)
		ZoneCoverage = PrimaryProgression(cp.Progression, cp.ZoneCoverage, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Zone Coverage", cp.IsRedshirting)
		Catching = PrimaryProgression(cp.Progression, cp.Catching, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Catching", cp.IsRedshirting)

		// Secondary
		ThrowPower = SecondaryProgression(cp.Progression, cp.ThrowPower)
		ThrowAccuracy = SecondaryProgression(cp.Progression, cp.ThrowAccuracy)
		RunDefense = SecondaryProgression(cp.Progression, cp.RunDefense)
		PassRush = SecondaryProgression(cp.Progression, cp.PassRush)
		PuntPower = SecondaryProgression(cp.Progression, cp.PuntPower)
		PuntAccuracy = SecondaryProgression(cp.Progression, cp.PuntAccuracy)
		KickPower = SecondaryProgression(cp.Progression, cp.KickPower)
		KickAccuracy = SecondaryProgression(cp.Progression, cp.KickAccuracy)
		Carrying = SecondaryProgression(cp.Progression, cp.Carrying)
		RunBlock = SecondaryProgression(cp.Progression, cp.RunBlock)
		PassBlock = SecondaryProgression(cp.Progression, cp.PassBlock)
		RouteRunning = SecondaryProgression(cp.RouteRunning, cp.RouteRunning)
	} else if mostPlayedPosition == "FS" || mostPlayedPosition == "SS" {
		// Primary
		Agility = PrimaryProgression(cp.Progression, cp.Agility, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Agility", cp.IsRedshirting)
		FootballIQ = PrimaryProgression(cp.Progression, cp.FootballIQ, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Football IQ", cp.IsRedshirting)
		Strength = PrimaryProgression(cp.Progression, cp.Strength, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Strength", cp.IsRedshirting)
		Speed = PrimaryProgression(cp.Progression, cp.Speed, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Speed", cp.IsRedshirting)
		RunDefense = PrimaryProgression(cp.Progression, cp.RunDefense, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Run Defense", cp.IsRedshirting)
		Tackle = PrimaryProgression(cp.Progression, cp.Tackle, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Tackle", cp.IsRedshirting)
		ManCoverage = PrimaryProgression(cp.Progression, cp.ManCoverage, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Man Coverage", cp.IsRedshirting)
		ZoneCoverage = PrimaryProgression(cp.Progression, cp.ZoneCoverage, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Zone Coverage", cp.IsRedshirting)
		Catching = PrimaryProgression(cp.Progression, cp.Catching, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Catching", cp.IsRedshirting)

		// Secondary
		ThrowPower = SecondaryProgression(cp.Progression, cp.ThrowPower)
		ThrowAccuracy = SecondaryProgression(cp.Progression, cp.ThrowAccuracy)
		PassRush = SecondaryProgression(cp.Progression, cp.PassRush)
		PuntPower = SecondaryProgression(cp.Progression, cp.PuntPower)
		PuntAccuracy = SecondaryProgression(cp.Progression, cp.PuntAccuracy)
		KickPower = SecondaryProgression(cp.Progression, cp.KickPower)
		KickAccuracy = SecondaryProgression(cp.Progression, cp.KickAccuracy)
		Carrying = SecondaryProgression(cp.Progression, cp.Carrying)
		RunBlock = SecondaryProgression(cp.Progression, cp.RunBlock)
		PassBlock = SecondaryProgression(cp.Progression, cp.PassBlock)
		RouteRunning = SecondaryProgression(cp.RouteRunning, cp.RouteRunning)
	} else if mostPlayedPosition == "K" {
		// Primary
		KickPower = PrimaryProgression(cp.Progression, cp.KickPower, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Kick Power", cp.IsRedshirting)
		KickAccuracy = PrimaryProgression(cp.Progression, cp.KickAccuracy, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Kick Accuracy", cp.IsRedshirting)
		// Secondary
		ThrowPower = SecondaryProgression(cp.Progression, cp.ThrowPower)
		ThrowAccuracy = SecondaryProgression(cp.Progression, cp.ThrowAccuracy)
		PassRush = SecondaryProgression(cp.Progression, cp.PassRush)
		RunDefense = SecondaryProgression(cp.Progression, cp.RunDefense)
		PuntPower = SecondaryProgression(cp.Progression, cp.PuntPower)
		PuntAccuracy = SecondaryProgression(cp.Progression, cp.PuntAccuracy)
		Carrying = SecondaryProgression(cp.Progression, cp.Carrying)
		RunBlock = SecondaryProgression(cp.Progression, cp.RunBlock)
		PassBlock = SecondaryProgression(cp.Progression, cp.PassBlock)
		Catching = SecondaryProgression(cp.Progression, cp.Catching)
		RouteRunning = SecondaryProgression(cp.RouteRunning, cp.RouteRunning)
		Strength = SecondaryProgression(cp.Progression, cp.Strength)
		Speed = SecondaryProgression(cp.Progression, cp.Speed)
		Agility = SecondaryProgression(cp.Progression, cp.Agility)
		ManCoverage = SecondaryProgression(cp.Progression, cp.ManCoverage)
		ZoneCoverage = SecondaryProgression(cp.Progression, cp.ZoneCoverage)
		FootballIQ = SecondaryProgression(cp.Progression, cp.FootballIQ)
		Tackle = SecondaryProgression(cp.Progression, cp.Tackle)
	} else if mostPlayedPosition == "P" {
		// If David Ross
		if cp.ID == 24984 {
			ThrowPower = PrimaryProgression(cp.Progression, cp.ThrowPower, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Throw Power", cp.IsRedshirting)
			FootballIQ = PrimaryProgression(cp.Progression, cp.FootballIQ, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Football IQ", cp.IsRedshirting)
			ThrowAccuracy = PrimaryProgression(cp.Progression, cp.ThrowAccuracy, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Throw Accuracy", cp.IsRedshirting)
		} else {
			ThrowPower = SecondaryProgression(cp.Progression, cp.ThrowPower)
			ThrowAccuracy = SecondaryProgression(cp.Progression, cp.ThrowAccuracy)
			FootballIQ = SecondaryProgression(cp.Progression, cp.FootballIQ)
		}
		// Primary Progressions
		PuntPower = PrimaryProgression(cp.Progression, cp.PuntPower, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Punt Power", cp.IsRedshirting)
		PuntAccuracy = PrimaryProgression(cp.Progression, cp.PuntAccuracy, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Punt Accuracy", cp.IsRedshirting)
		// Secondary Progressions
		RunBlock = SecondaryProgression(cp.Progression, cp.RunBlock)
		PassBlock = SecondaryProgression(cp.Progression, cp.PassBlock)
		RunDefense = SecondaryProgression(cp.Progression, cp.RunDefense)
		PassRush = SecondaryProgression(cp.Progression, cp.PassRush)
		Carrying = SecondaryProgression(cp.Progression, cp.Carrying)
		Tackle = SecondaryProgression(cp.Progression, cp.Tackle)
		RouteRunning = SecondaryProgression(cp.RouteRunning, cp.RouteRunning)
		Catching = SecondaryProgression(cp.Progression, cp.Catching)
		KickPower = SecondaryProgression(cp.Progression, cp.KickPower)
		KickAccuracy = SecondaryProgression(cp.Progression, cp.KickAccuracy)
		ManCoverage = SecondaryProgression(cp.Progression, cp.ManCoverage)
		ZoneCoverage = SecondaryProgression(cp.Progression, cp.ZoneCoverage)
		Strength = SecondaryProgression(cp.Progression, cp.Strength)
		Speed = SecondaryProgression(cp.Progression, cp.Speed)
		Agility = SecondaryProgression(cp.Progression, cp.Agility)
	} else if mostPlayedPosition == "ATH" {
		// Use Archetype
		archetype := cp.Archetype
		if archetype == "Field General" {
			Agility = PrimaryProgression(cp.Progression, cp.Agility, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Agility", cp.IsRedshirting)
			ThrowPower = PrimaryProgression(cp.Progression, cp.ThrowPower, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Throw Power", cp.IsRedshirting)
			ThrowAccuracy = PrimaryProgression(cp.Progression, cp.ThrowAccuracy, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Throw Accuracy", cp.IsRedshirting)
			Speed = PrimaryProgression(cp.Progression, cp.Speed, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Speed", cp.IsRedshirting)
			FootballIQ = PrimaryProgression(cp.Progression, cp.FootballIQ, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Football IQ", cp.IsRedshirting)
			Strength = PrimaryProgression(cp.Progression, cp.Strength, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Strength", cp.IsRedshirting)
			RunDefense = PrimaryProgression(cp.Progression, cp.RunDefense, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Run Defense", cp.IsRedshirting)
			Tackle = PrimaryProgression(cp.Progression, cp.Tackle, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Tackle", cp.IsRedshirting)
			ManCoverage = PrimaryProgression(cp.Progression, cp.ManCoverage, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Man Coverage", cp.IsRedshirting)
			ZoneCoverage = PrimaryProgression(cp.Progression, cp.ZoneCoverage, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Zone Coverage", cp.IsRedshirting)

			// Secondary
			PassRush = SecondaryProgression(cp.Progression, cp.PassRush)
			RunBlock = SecondaryProgression(cp.Progression, cp.RunBlock)
			PassBlock = SecondaryProgression(cp.Progression, cp.PassBlock)
			Carrying = SecondaryProgression(cp.Progression, cp.Carrying)
			RouteRunning = SecondaryProgression(cp.RouteRunning, cp.RouteRunning)
			Catching = SecondaryProgression(cp.Progression, cp.Catching)
			PuntPower = SecondaryProgression(cp.Progression, cp.PuntPower)
			PuntAccuracy = SecondaryProgression(cp.Progression, cp.PuntAccuracy)
			KickPower = SecondaryProgression(cp.Progression, cp.KickPower)
			KickAccuracy = SecondaryProgression(cp.Progression, cp.KickAccuracy)
		} else if archetype == "Triple-Threat" {
			Agility = PrimaryProgression(cp.Progression, cp.Agility, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Agility", cp.IsRedshirting)
			ThrowPower = PrimaryProgression(cp.Progression, cp.ThrowPower, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Throw Power", cp.IsRedshirting)
			ThrowAccuracy = PrimaryProgression(cp.Progression, cp.ThrowAccuracy, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Throw Accuracy", cp.IsRedshirting)
			Speed = PrimaryProgression(cp.Progression, cp.Speed, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Speed", cp.IsRedshirting)
			FootballIQ = PrimaryProgression(cp.Progression, cp.FootballIQ, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Football IQ", cp.IsRedshirting)
			KickPower = PrimaryProgression(cp.Progression, cp.KickPower, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Kick Power", cp.IsRedshirting)
			KickAccuracy = PrimaryProgression(cp.Progression, cp.KickAccuracy, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Kick Accuracy", cp.IsRedshirting)
			PuntPower = PrimaryProgression(cp.Progression, cp.PuntPower, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Punt Power", cp.IsRedshirting)
			PuntAccuracy = PrimaryProgression(cp.Progression, cp.PuntAccuracy, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Punt Accuracy", cp.IsRedshirting)
			RouteRunning = SecondaryProgression(cp.RouteRunning, cp.RouteRunning)
			Carrying = SecondaryProgression(cp.Progression, cp.Carrying)
			Catching = SecondaryProgression(cp.Progression, cp.Catching)
			RunBlock = SecondaryProgression(cp.Progression, cp.RunBlock)
			PassBlock = SecondaryProgression(cp.Progression, cp.PassBlock)
			RunDefense = SecondaryProgression(cp.Progression, cp.RunDefense)
			PassRush = SecondaryProgression(cp.Progression, cp.PassRush)
			Tackle = SecondaryProgression(cp.Progression, cp.Tackle)
			ManCoverage = SecondaryProgression(cp.Progression, cp.ManCoverage)
			ZoneCoverage = SecondaryProgression(cp.Progression, cp.ZoneCoverage)
			Strength = SecondaryProgression(cp.Progression, cp.Strength)
		} else if archetype == "Wingback" {
			Agility = PrimaryProgression(cp.Progression, cp.Agility, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Agility", cp.IsRedshirting)
			Speed = PrimaryProgression(cp.Progression, cp.Speed, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Speed", cp.IsRedshirting)
			FootballIQ = PrimaryProgression(cp.Progression, cp.FootballIQ, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Football IQ", cp.IsRedshirting)
			Catching = PrimaryProgression(cp.Progression, cp.Catching, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Catching", cp.IsRedshirting)
			RouteRunning = PrimaryProgression(cp.Progression, cp.RouteRunning, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Route Running", cp.IsRedshirting)
			Carrying = PrimaryProgression(cp.Progression, cp.Carrying, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Carrying", cp.IsRedshirting)
			RunBlock = SecondaryProgression(cp.Progression, cp.RunBlock)
			PassBlock = SecondaryProgression(cp.Progression, cp.PassBlock)
			RunDefense = SecondaryProgression(cp.Progression, cp.RunDefense)
			PassRush = SecondaryProgression(cp.Progression, cp.PassRush)
			Tackle = SecondaryProgression(cp.Progression, cp.Tackle)
			PuntPower = SecondaryProgression(cp.Progression, cp.PuntPower)
			PuntAccuracy = SecondaryProgression(cp.Progression, cp.PuntAccuracy)
			KickPower = SecondaryProgression(cp.Progression, cp.KickPower)
			KickAccuracy = SecondaryProgression(cp.Progression, cp.KickAccuracy)
			ManCoverage = SecondaryProgression(cp.Progression, cp.ManCoverage)
			ZoneCoverage = SecondaryProgression(cp.Progression, cp.ZoneCoverage)
			Strength = SecondaryProgression(cp.Progression, cp.Strength)
			ThrowPower = SecondaryProgression(cp.Progression, cp.ThrowPower)
			ThrowAccuracy = SecondaryProgression(cp.Progression, cp.ThrowAccuracy)
		} else if archetype == "Slotback" {
			Agility = PrimaryProgression(cp.Progression, cp.Agility, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Agility", cp.IsRedshirting)
			Speed = PrimaryProgression(cp.Progression, cp.Speed, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Speed", cp.IsRedshirting)
			FootballIQ = PrimaryProgression(cp.Progression, cp.FootballIQ, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Football IQ", cp.IsRedshirting)
			Catching = PrimaryProgression(cp.Progression, cp.Catching, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Catching", cp.IsRedshirting)
			RouteRunning = PrimaryProgression(cp.Progression, cp.RouteRunning, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Route Running", cp.IsRedshirting)
			PassBlock = PrimaryProgression(cp.Progression, cp.PassBlock, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Pass Blocking", cp.IsRedshirting)
			RunBlock = PrimaryProgression(cp.Progression, cp.RunBlock, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Run Blocking", cp.IsRedshirting)
			Strength = PrimaryProgression(cp.Progression, cp.Strength, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Strength", cp.IsRedshirting)
			Carrying = SecondaryProgression(cp.Progression, cp.Carrying)
			RunDefense = SecondaryProgression(cp.Progression, cp.RunDefense)
			PassRush = SecondaryProgression(cp.Progression, cp.PassRush)
			Tackle = SecondaryProgression(cp.Progression, cp.Tackle)
			PuntPower = SecondaryProgression(cp.Progression, cp.PuntPower)
			PuntAccuracy = SecondaryProgression(cp.Progression, cp.PuntAccuracy)
			KickPower = SecondaryProgression(cp.Progression, cp.KickPower)
			KickAccuracy = SecondaryProgression(cp.Progression, cp.KickAccuracy)
			ManCoverage = SecondaryProgression(cp.Progression, cp.ManCoverage)
			ZoneCoverage = SecondaryProgression(cp.Progression, cp.ZoneCoverage)
			ThrowPower = SecondaryProgression(cp.Progression, cp.ThrowPower)
			ThrowAccuracy = SecondaryProgression(cp.Progression, cp.ThrowAccuracy)
		} else if archetype == "Lineman" {
			Agility = PrimaryProgression(cp.Progression, cp.Agility, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Agility", cp.IsRedshirting)
			FootballIQ = PrimaryProgression(cp.Progression, cp.FootballIQ, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Football IQ", cp.IsRedshirting)
			Strength = PrimaryProgression(cp.Progression, cp.Strength, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Strength", cp.IsRedshirting)
			PassBlock = PrimaryProgression(cp.Progression, cp.PassBlock, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Pass Blocking", cp.IsRedshirting)
			RunBlock = PrimaryProgression(cp.Progression, cp.RunBlock, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Run Blocking", cp.IsRedshirting)
			PassRush = PrimaryProgression(cp.Progression, cp.PassRush, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Pass Rush", cp.IsRedshirting)
			RunDefense = PrimaryProgression(cp.Progression, cp.RunDefense, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Run Defense", cp.IsRedshirting)
			Tackle = PrimaryProgression(cp.Progression, cp.Tackle, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Tackle", cp.IsRedshirting)
			PuntPower = SecondaryProgression(cp.Progression, cp.PuntPower)
			PuntAccuracy = SecondaryProgression(cp.Progression, cp.PuntAccuracy)
			KickPower = SecondaryProgression(cp.Progression, cp.KickPower)
			KickAccuracy = SecondaryProgression(cp.Progression, cp.KickAccuracy)
			ManCoverage = SecondaryProgression(cp.Progression, cp.ManCoverage)
			ZoneCoverage = SecondaryProgression(cp.Progression, cp.ZoneCoverage)
			ThrowPower = SecondaryProgression(cp.Progression, cp.ThrowPower)
			ThrowAccuracy = SecondaryProgression(cp.Progression, cp.ThrowAccuracy)
			Carrying = SecondaryProgression(cp.Progression, cp.Carrying)
			Catching = SecondaryProgression(cp.Progression, cp.Catching)
			RouteRunning = SecondaryProgression(cp.RouteRunning, cp.RouteRunning)
			Speed = SecondaryProgression(cp.Progression, cp.Speed)
		} else if archetype == "Strongside" {
			Agility = PrimaryProgression(cp.Progression, cp.Agility, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Agility", cp.IsRedshirting)
			FootballIQ = PrimaryProgression(cp.Progression, cp.FootballIQ, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Football IQ", cp.IsRedshirting)
			Strength = PrimaryProgression(cp.Progression, cp.Strength, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Strength", cp.IsRedshirting)
			Speed = PrimaryProgression(cp.Progression, cp.Speed, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Speed", cp.IsRedshirting)
			PassRush = PrimaryProgression(cp.Progression, cp.PassRush, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Pass Rush", cp.IsRedshirting)
			RunDefense = PrimaryProgression(cp.Progression, cp.RunDefense, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Run Defense", cp.IsRedshirting)
			Tackle = PrimaryProgression(cp.Progression, cp.Tackle, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Tackle", cp.IsRedshirting)
			ManCoverage = PrimaryProgression(cp.Progression, cp.ManCoverage, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Man Coverage", cp.IsRedshirting)
			ZoneCoverage = PrimaryProgression(cp.Progression, cp.ZoneCoverage, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Zone Coverage", cp.IsRedshirting)
			PuntPower = SecondaryProgression(cp.Progression, cp.PuntPower)
			PuntAccuracy = SecondaryProgression(cp.Progression, cp.PuntAccuracy)
			KickPower = SecondaryProgression(cp.Progression, cp.KickPower)
			KickAccuracy = SecondaryProgression(cp.Progression, cp.KickAccuracy)
			ThrowPower = SecondaryProgression(cp.Progression, cp.ThrowPower)
			ThrowAccuracy = SecondaryProgression(cp.Progression, cp.ThrowAccuracy)
			Carrying = SecondaryProgression(cp.Progression, cp.Carrying)
			Catching = SecondaryProgression(cp.Progression, cp.Catching)
			RouteRunning = SecondaryProgression(cp.RouteRunning, cp.RouteRunning)
			RunBlock = SecondaryProgression(cp.Progression, cp.RunBlock)
			PassBlock = SecondaryProgression(cp.Progression, cp.PassBlock)
		} else if archetype == "Weakside" {
			Agility = PrimaryProgression(cp.Progression, cp.Agility, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Agility", cp.IsRedshirting)
			FootballIQ = PrimaryProgression(cp.Progression, cp.FootballIQ, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Football IQ", cp.IsRedshirting)
			Strength = PrimaryProgression(cp.Progression, cp.Strength, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Strength", cp.IsRedshirting)
			Speed = PrimaryProgression(cp.Progression, cp.Speed, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Speed", cp.IsRedshirting)
			PassRush = PrimaryProgression(cp.Progression, cp.PassRush, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Pass Rush", cp.IsRedshirting)
			RunDefense = PrimaryProgression(cp.Progression, cp.RunDefense, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Run Defense", cp.IsRedshirting)
			Tackle = PrimaryProgression(cp.Progression, cp.Tackle, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Tackle", cp.IsRedshirting)
			ManCoverage = PrimaryProgression(cp.Progression, cp.ManCoverage, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Man Coverage", cp.IsRedshirting)
			ZoneCoverage = PrimaryProgression(cp.Progression, cp.ZoneCoverage, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Zone Coverage", cp.IsRedshirting)
			PuntPower = SecondaryProgression(cp.Progression, cp.PuntPower)
			PuntAccuracy = SecondaryProgression(cp.Progression, cp.PuntAccuracy)
			KickPower = SecondaryProgression(cp.Progression, cp.KickPower)
			KickAccuracy = SecondaryProgression(cp.Progression, cp.KickAccuracy)
			ThrowPower = SecondaryProgression(cp.Progression, cp.ThrowPower)
			ThrowAccuracy = SecondaryProgression(cp.Progression, cp.ThrowAccuracy)
			Carrying = SecondaryProgression(cp.Progression, cp.Carrying)
			Catching = SecondaryProgression(cp.Progression, cp.Catching)
			RouteRunning = SecondaryProgression(cp.RouteRunning, cp.RouteRunning)
			RunBlock = SecondaryProgression(cp.Progression, cp.RunBlock)
			PassBlock = SecondaryProgression(cp.Progression, cp.PassBlock)
		} else if archetype == "Bandit" {
			Agility = PrimaryProgression(cp.Progression, cp.Agility, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Agility", cp.IsRedshirting)
			FootballIQ = PrimaryProgression(cp.Progression, cp.FootballIQ, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Football IQ", cp.IsRedshirting)
			Strength = PrimaryProgression(cp.Progression, cp.Strength, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Strength", cp.IsRedshirting)
			Speed = PrimaryProgression(cp.Progression, cp.Speed, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Speed", cp.IsRedshirting)
			PassRush = PrimaryProgression(cp.Progression, cp.PassRush, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Pass Rush", cp.IsRedshirting)
			RunDefense = PrimaryProgression(cp.Progression, cp.RunDefense, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Run Defense", cp.IsRedshirting)
			Tackle = PrimaryProgression(cp.Progression, cp.Tackle, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Tackle", cp.IsRedshirting)
			ManCoverage = PrimaryProgression(cp.Progression, cp.ManCoverage, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Man Coverage", cp.IsRedshirting)
			ZoneCoverage = PrimaryProgression(cp.Progression, cp.ZoneCoverage, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Zone Coverage", cp.IsRedshirting)
			PuntPower = SecondaryProgression(cp.Progression, cp.PuntPower)
			PuntAccuracy = SecondaryProgression(cp.Progression, cp.PuntAccuracy)
			KickPower = SecondaryProgression(cp.Progression, cp.KickPower)
			KickAccuracy = SecondaryProgression(cp.Progression, cp.KickAccuracy)
			ThrowPower = SecondaryProgression(cp.Progression, cp.ThrowPower)
			ThrowAccuracy = SecondaryProgression(cp.Progression, cp.ThrowAccuracy)
			Carrying = SecondaryProgression(cp.Progression, cp.Carrying)
			Catching = SecondaryProgression(cp.Progression, cp.Catching)
			RouteRunning = SecondaryProgression(cp.RouteRunning, cp.RouteRunning)
			RunBlock = SecondaryProgression(cp.Progression, cp.RunBlock)
			PassBlock = SecondaryProgression(cp.Progression, cp.PassBlock)
		} else if archetype == "Return Specialist" {
			Agility = PrimaryProgression(cp.Progression, cp.Agility, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Agility", cp.IsRedshirting)
			Speed = PrimaryProgression(cp.Progression, cp.Speed, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Speed", cp.IsRedshirting)
			Carrying = PrimaryProgression(cp.Progression, cp.Carrying, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Carrying", cp.IsRedshirting)
			Catching = PrimaryProgression(cp.Progression, cp.Catching, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Catching", cp.IsRedshirting)
			RouteRunning = PrimaryProgression(cp.Progression, cp.RouteRunning, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Route Running", cp.IsRedshirting)
			FootballIQ = PrimaryProgression(cp.Progression, cp.FootballIQ, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Football IQ", cp.IsRedshirting)
			PuntPower = SecondaryProgression(cp.Progression, cp.PuntPower)
			PuntAccuracy = SecondaryProgression(cp.Progression, cp.PuntAccuracy)
			KickPower = SecondaryProgression(cp.Progression, cp.KickPower)
			KickAccuracy = SecondaryProgression(cp.Progression, cp.KickAccuracy)
			ThrowPower = SecondaryProgression(cp.Progression, cp.ThrowPower)
			ThrowAccuracy = SecondaryProgression(cp.Progression, cp.ThrowAccuracy)
			RunBlock = SecondaryProgression(cp.Progression, cp.RunBlock)
			PassBlock = SecondaryProgression(cp.Progression, cp.PassBlock)
			Strength = SecondaryProgression(cp.Progression, cp.Strength)
			ManCoverage = SecondaryProgression(cp.Progression, cp.ManCoverage)
			ZoneCoverage = SecondaryProgression(cp.Progression, cp.ZoneCoverage)
			Tackle = SecondaryProgression(cp.Progression, cp.Tackle)
			PassRush = SecondaryProgression(cp.Progression, cp.PassRush)
			RunDefense = SecondaryProgression(cp.Progression, cp.RunDefense)
		} else if archetype == "Soccer Player" {
			KickPower = PrimaryProgression(cp.Progression, cp.KickPower, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Kick Power", cp.IsRedshirting)
			KickAccuracy = PrimaryProgression(cp.Progression, cp.KickAccuracy, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Kick Accuracy", cp.IsRedshirting)
			PuntPower = PrimaryProgression(cp.Progression, cp.PuntPower, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Punt Power", cp.IsRedshirting)
			PuntAccuracy = PrimaryProgression(cp.Progression, cp.PuntAccuracy, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Punt Accuracy", cp.IsRedshirting)
			Agility = PrimaryProgression(cp.Progression, cp.Agility, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Agility", cp.IsRedshirting)
			Speed = PrimaryProgression(cp.Progression, cp.Speed, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Speed", cp.IsRedshirting)
			FootballIQ = PrimaryProgression(cp.Progression, cp.FootballIQ, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Football IQ", cp.IsRedshirting)
			Catching = PrimaryProgression(cp.Progression, cp.Catching, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Catching", cp.IsRedshirting)
			RouteRunning = PrimaryProgression(cp.Progression, cp.RouteRunning, mostPlayedPosition, cp.Archetype, SnapsPerGame, "Route Running", cp.IsRedshirting)
			Strength = SecondaryProgression(cp.Progression, cp.Strength)
			RunBlock = SecondaryProgression(cp.Progression, cp.RunBlock)
			PassBlock = SecondaryProgression(cp.Progression, cp.PassBlock)
			RunDefense = SecondaryProgression(cp.Progression, cp.RunDefense)
			PassRush = SecondaryProgression(cp.Progression, cp.PassRush)
			Tackle = SecondaryProgression(cp.Progression, cp.Tackle)
			ManCoverage = SecondaryProgression(cp.Progression, cp.ManCoverage)
			ZoneCoverage = SecondaryProgression(cp.Progression, cp.ZoneCoverage)
			ThrowPower = SecondaryProgression(cp.Progression, cp.ThrowPower)
			ThrowAccuracy = SecondaryProgression(cp.Progression, cp.ThrowAccuracy)
			Carrying = SecondaryProgression(cp.Progression, cp.Carrying)
		}
	}

	progressions := structs.CollegePlayerProgressions{
		Agility:       Agility,
		Catching:      Catching,
		Carrying:      Carrying,
		Speed:         Speed,
		RouteRunning:  RouteRunning,
		RunBlock:      RunBlock,
		PassBlock:     PassBlock,
		RunDefense:    RunDefense,
		PassRush:      PassRush,
		Strength:      Strength,
		Tackle:        Tackle,
		ThrowPower:    ThrowPower,
		ThrowAccuracy: ThrowAccuracy,
		PuntAccuracy:  PuntAccuracy,
		PuntPower:     PuntPower,
		KickAccuracy:  KickAccuracy,
		KickPower:     KickPower,
		FootballIQ:    FootballIQ,
		ManCoverage:   ManCoverage,
		ZoneCoverage:  ZoneCoverage,
	}

	cp.Progress(progressions, false)

	if cp.IsRedshirting {
		cp.SetRedshirtStatus()
	}

	cp.GetOverall()

	return cp
}

func BoomBustDraftee(np models.NFLDraftee, SeasonID string, SnapsPerGame int, isBoom bool) models.NFLDraftee {
	Agility := 0
	ThrowPower := 0
	ThrowAccuracy := 0
	Speed := 0
	FootballIQ := 0
	Strength := 0
	RunBlock := 0
	PassBlock := 0
	RunDefense := 0
	PassRush := 0
	Carrying := 0
	Tackle := 0
	RouteRunning := 0
	Catching := 0
	PuntPower := 0
	PuntAccuracy := 0
	KickPower := 0
	KickAccuracy := 0
	ManCoverage := 0
	ZoneCoverage := 0

	if isBoom {
		if np.Position == "QB" {
			// Primary Progressions
			Agility = PrimaryProgression(np.Progression, np.Agility, np.Position, np.Archetype, SnapsPerGame, "Agility", false)
			ThrowPower = PrimaryProgression(np.Progression, np.ThrowPower, np.Position, np.Archetype, SnapsPerGame, "Throw Power", false)
			ThrowAccuracy = PrimaryProgression(np.Progression, np.ThrowAccuracy, np.Position, np.Archetype, SnapsPerGame, "Throw Accuracy", false)
			Speed = PrimaryProgression(np.Progression, np.Speed, np.Position, np.Archetype, SnapsPerGame, "Speed", false)
			FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, np.Position, np.Archetype, SnapsPerGame, "Football IQ", false)
			Strength = PrimaryProgression(np.Progression, np.Strength, np.Position, np.Archetype, SnapsPerGame, "Strength", false)

			// Secondary Progressions
			RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
			PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
			RunDefense = SecondaryProgression(np.Progression, np.RunDefense)
			PassRush = SecondaryProgression(np.Progression, np.PassRush)
			Carrying = SecondaryProgression(np.Progression, np.Carrying)
			Tackle = SecondaryProgression(np.Progression, np.Tackle)
			RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
			Catching = SecondaryProgression(np.Progression, np.Catching)
			PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
			PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
			KickPower = SecondaryProgression(np.Progression, np.KickPower)
			KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
			ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
			ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
		} else if np.Position == "RB" {
			// Primary
			Agility = PrimaryProgression(np.Progression, np.Agility, np.Position, np.Archetype, SnapsPerGame, "Agility", false)
			Carrying = PrimaryProgression(np.Progression, np.Carrying, np.Position, np.Archetype, SnapsPerGame, "Carrying", false)
			Speed = PrimaryProgression(np.Progression, np.Speed, np.Position, np.Archetype, SnapsPerGame, "Speed", false)
			FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, np.Position, np.Archetype, SnapsPerGame, "Football IQ", false)
			Strength = PrimaryProgression(np.Progression, np.Strength, np.Position, np.Archetype, SnapsPerGame, "Strength", false)
			Catching = PrimaryProgression(np.Progression, np.Catching, np.Position, np.Archetype, SnapsPerGame, "Catching", false)
			PassBlock = PrimaryProgression(np.Progression, np.PassBlock, np.Position, np.Archetype, SnapsPerGame, "Pass Blocking", false)
			// Secondary
			RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
			RunDefense = SecondaryProgression(np.Progression, np.RunDefense)
			PassRush = SecondaryProgression(np.Progression, np.PassRush)
			Tackle = SecondaryProgression(np.Progression, np.Tackle)
			RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
			PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
			PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
			KickPower = SecondaryProgression(np.Progression, np.KickPower)
			KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
			ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
			ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
			ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
			ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
		} else if np.Position == "FB" {
			// Primary
			Agility = PrimaryProgression(np.Progression, np.Agility, np.Position, np.Archetype, SnapsPerGame, "Agility", false)
			Carrying = PrimaryProgression(np.Progression, np.Carrying, np.Position, np.Archetype, SnapsPerGame, "Carrying", false)
			Speed = PrimaryProgression(np.Progression, np.Speed, np.Position, np.Archetype, SnapsPerGame, "Speed", false)
			FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, np.Position, np.Archetype, SnapsPerGame, "Football IQ", false)
			Strength = PrimaryProgression(np.Progression, np.Strength, np.Position, np.Archetype, SnapsPerGame, "Strength", false)
			Catching = PrimaryProgression(np.Progression, np.Catching, np.Position, np.Archetype, SnapsPerGame, "Catching", false)
			PassBlock = PrimaryProgression(np.Progression, np.PassBlock, np.Position, np.Archetype, SnapsPerGame, "Pass Blocking", false)
			RunBlock = PrimaryProgression(np.Progression, np.RunBlock, np.Position, np.Archetype, SnapsPerGame, "Run Blocking", false)

			// Secondary
			RunDefense = SecondaryProgression(np.Progression, np.RunDefense)
			PassRush = SecondaryProgression(np.Progression, np.PassRush)
			Tackle = SecondaryProgression(np.Progression, np.Tackle)
			RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
			PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
			PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
			KickPower = SecondaryProgression(np.Progression, np.KickPower)
			KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
			ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
			ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
			ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
			ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)

		} else if np.Position == "TE" {
			// Primary
			Agility = PrimaryProgression(np.Progression, np.Agility, np.Position, np.Archetype, SnapsPerGame, "Agility", false)
			Carrying = PrimaryProgression(np.Progression, np.Carrying, np.Position, np.Archetype, SnapsPerGame, "Carrying", false)
			Speed = PrimaryProgression(np.Progression, np.Speed, np.Position, np.Archetype, SnapsPerGame, "Speed", false)
			FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, np.Position, np.Archetype, SnapsPerGame, "Football IQ", false)
			Strength = PrimaryProgression(np.Progression, np.Strength, np.Position, np.Archetype, SnapsPerGame, "Strength", false)
			Catching = PrimaryProgression(np.Progression, np.Catching, np.Position, np.Archetype, SnapsPerGame, "Catching", false)
			RouteRunning = PrimaryProgression(np.Progression, np.RouteRunning, np.Position, np.Archetype, SnapsPerGame, "Route Running", false)
			PassBlock = PrimaryProgression(np.Progression, np.PassBlock, np.Position, np.Archetype, SnapsPerGame, "Pass Blocking", false)
			RunBlock = PrimaryProgression(np.Progression, np.RunBlock, np.Position, np.Archetype, SnapsPerGame, "Run Blocking", false)

			// Secondary
			RunDefense = SecondaryProgression(np.Progression, np.RunDefense)
			PassRush = SecondaryProgression(np.Progression, np.PassRush)
			Tackle = SecondaryProgression(np.Progression, np.Tackle)
			PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
			PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
			KickPower = SecondaryProgression(np.Progression, np.KickPower)
			KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
			ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
			ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
			ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
			ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
		} else if np.Position == "WR" {
			// Primary
			Agility = PrimaryProgression(np.Progression, np.Agility, np.Position, np.Archetype, SnapsPerGame, "Agility", false)
			Carrying = PrimaryProgression(np.Progression, np.Carrying, np.Position, np.Archetype, SnapsPerGame, "Carrying", false)
			Speed = PrimaryProgression(np.Progression, np.Speed, np.Position, np.Archetype, SnapsPerGame, "Speed", false)
			FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, np.Position, np.Archetype, SnapsPerGame, "Football IQ", false)
			Strength = PrimaryProgression(np.Progression, np.Strength, np.Position, np.Archetype, SnapsPerGame, "Strength", false)
			Catching = PrimaryProgression(np.Progression, np.Catching, np.Position, np.Archetype, SnapsPerGame, "Catching", false)
			RouteRunning = PrimaryProgression(np.Progression, np.RouteRunning, np.Position, np.Archetype, SnapsPerGame, "Route Running", false)

			// Secondary
			RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
			PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
			RunDefense = SecondaryProgression(np.Progression, np.RunDefense)
			PassRush = SecondaryProgression(np.Progression, np.PassRush)
			Tackle = SecondaryProgression(np.Progression, np.Tackle)
			PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
			PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
			KickPower = SecondaryProgression(np.Progression, np.KickPower)
			KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
			ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
			ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
			ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
			ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
		} else if np.Position == "OT" || np.Position == "OG" || np.Position == "C" {
			// Primary
			Agility = PrimaryProgression(np.Progression, np.Agility, np.Position, np.Archetype, SnapsPerGame, "Agility", false)
			FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, np.Position, np.Archetype, SnapsPerGame, "Football IQ", false)
			Strength = PrimaryProgression(np.Progression, np.Strength, np.Position, np.Archetype, SnapsPerGame, "Strength", false)
			PassBlock = PrimaryProgression(np.Progression, np.PassBlock, np.Position, np.Archetype, SnapsPerGame, "Pass Blocking", false)
			RunBlock = PrimaryProgression(np.Progression, np.RunBlock, np.Position, np.Archetype, SnapsPerGame, "Run Blocking", false)

			// Secondary
			RunDefense = SecondaryProgression(np.Progression, np.RunDefense)
			PassRush = SecondaryProgression(np.Progression, np.PassRush)
			Tackle = SecondaryProgression(np.Progression, np.Tackle)
			RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
			PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
			PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
			KickPower = SecondaryProgression(np.Progression, np.KickPower)
			KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
			ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
			ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
			ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
			ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
			Carrying = SecondaryProgression(np.Progression, np.Carrying)
			Speed = SecondaryProgression(np.Progression, np.Speed)
			Catching = SecondaryProgression(np.Progression, np.Catching)
		} else if np.Position == "DT" {
			// Primary
			Agility = PrimaryProgression(np.Progression, np.Agility, np.Position, np.Archetype, SnapsPerGame, "Agility", false)
			FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, np.Position, np.Archetype, SnapsPerGame, "Football IQ", false)
			Strength = PrimaryProgression(np.Progression, np.Strength, np.Position, np.Archetype, SnapsPerGame, "Strength", false)
			PassRush = PrimaryProgression(np.Progression, np.PassRush, np.Position, np.Archetype, SnapsPerGame, "Pass Rush", false)
			RunDefense = PrimaryProgression(np.Progression, np.RunDefense, np.Position, np.Archetype, SnapsPerGame, "Run Defense", false)
			Tackle = PrimaryProgression(np.Progression, np.Tackle, np.Position, np.Archetype, SnapsPerGame, "Tackle", false)

			// Secondary
			PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
			PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
			KickPower = SecondaryProgression(np.Progression, np.KickPower)
			KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
			ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
			ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
			ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
			ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
			Carrying = SecondaryProgression(np.Progression, np.Carrying)
			Catching = SecondaryProgression(np.Progression, np.Catching)
			RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
			Speed = SecondaryProgression(np.Progression, np.Speed)
			RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
			PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
		} else if np.Position == "DE" {
			// Primary
			Agility = PrimaryProgression(np.Progression, np.Agility, np.Position, np.Archetype, SnapsPerGame, "Agility", false)
			FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, np.Position, np.Archetype, SnapsPerGame, "Football IQ", false)
			Strength = PrimaryProgression(np.Progression, np.Strength, np.Position, np.Archetype, SnapsPerGame, "Strength", false)
			Speed = PrimaryProgression(np.Progression, np.Speed, np.Position, np.Archetype, SnapsPerGame, "Speed", false)
			PassRush = PrimaryProgression(np.Progression, np.PassRush, np.Position, np.Archetype, SnapsPerGame, "Pass Rush", false)
			RunDefense = PrimaryProgression(np.Progression, np.RunDefense, np.Position, np.Archetype, SnapsPerGame, "Run Defense", false)
			Tackle = PrimaryProgression(np.Progression, np.Tackle, np.Position, np.Archetype, SnapsPerGame, "Tackle", false)

			// Secondary
			PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
			PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
			KickPower = SecondaryProgression(np.Progression, np.KickPower)
			KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
			ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
			ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
			ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
			ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
			Carrying = SecondaryProgression(np.Progression, np.Carrying)
			Catching = SecondaryProgression(np.Progression, np.Catching)
			RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
			RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
			PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
		} else if np.Position == "OLB" || np.Position == "ILB" {
			// Primary
			Agility = PrimaryProgression(np.Progression, np.Agility, np.Position, np.Archetype, SnapsPerGame, "Agility", false)
			FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, np.Position, np.Archetype, SnapsPerGame, "Football IQ", false)
			Strength = PrimaryProgression(np.Progression, np.Strength, np.Position, np.Archetype, SnapsPerGame, "Strength", false)
			Speed = PrimaryProgression(np.Progression, np.Speed, np.Position, np.Archetype, SnapsPerGame, "Speed", false)
			PassRush = PrimaryProgression(np.Progression, np.PassRush, np.Position, np.Archetype, SnapsPerGame, "Pass Rush", false)
			RunDefense = PrimaryProgression(np.Progression, np.RunDefense, np.Position, np.Archetype, SnapsPerGame, "Run Defense", false)
			Tackle = PrimaryProgression(np.Progression, np.Tackle, np.Position, np.Archetype, SnapsPerGame, "Tackle", false)
			ManCoverage = PrimaryProgression(np.Progression, np.ManCoverage, np.Position, np.Archetype, SnapsPerGame, "Man Coverage", false)
			ZoneCoverage = PrimaryProgression(np.Progression, np.ZoneCoverage, np.Position, np.Archetype, SnapsPerGame, "Zone Coverage", false)

			// Secondary
			PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
			PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
			KickPower = SecondaryProgression(np.Progression, np.KickPower)
			KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
			ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
			ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
			Carrying = SecondaryProgression(np.Progression, np.Carrying)
			Catching = SecondaryProgression(np.Progression, np.Catching)
			RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
			RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
			PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
		} else if np.Position == "CB" {
			// Primary
			Agility = PrimaryProgression(np.Progression, np.Agility, np.Position, np.Archetype, SnapsPerGame, "Agility", false)
			FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, np.Position, np.Archetype, SnapsPerGame, "Football IQ", false)
			Strength = PrimaryProgression(np.Progression, np.Strength, np.Position, np.Archetype, SnapsPerGame, "Strength", false)
			Speed = PrimaryProgression(np.Progression, np.Speed, np.Position, np.Archetype, SnapsPerGame, "Speed", false)
			Tackle = PrimaryProgression(np.Progression, np.Tackle, np.Position, np.Archetype, SnapsPerGame, "Tackle", false)
			ManCoverage = PrimaryProgression(np.Progression, np.ManCoverage, np.Position, np.Archetype, SnapsPerGame, "Man Coverage", false)
			ZoneCoverage = PrimaryProgression(np.Progression, np.ZoneCoverage, np.Position, np.Archetype, SnapsPerGame, "Zone Coverage", false)
			Catching = PrimaryProgression(np.Progression, np.Catching, np.Position, np.Archetype, SnapsPerGame, "Catching", false)

			// Secondary
			ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
			ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
			RunDefense = SecondaryProgression(np.Progression, np.RunDefense)
			PassRush = SecondaryProgression(np.Progression, np.PassRush)
			PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
			PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
			KickPower = SecondaryProgression(np.Progression, np.KickPower)
			KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
			Carrying = SecondaryProgression(np.Progression, np.Carrying)
			RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
			PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
			RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
		} else if np.Position == "FS" || np.Position == "SS" {
			// Primary
			Agility = PrimaryProgression(np.Progression, np.Agility, np.Position, np.Archetype, SnapsPerGame, "Agility", false)
			FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, np.Position, np.Archetype, SnapsPerGame, "Football IQ", false)
			Strength = PrimaryProgression(np.Progression, np.Strength, np.Position, np.Archetype, SnapsPerGame, "Strength", false)
			Speed = PrimaryProgression(np.Progression, np.Speed, np.Position, np.Archetype, SnapsPerGame, "Speed", false)
			RunDefense = PrimaryProgression(np.Progression, np.RunDefense, np.Position, np.Archetype, SnapsPerGame, "Run Defense", false)
			Tackle = PrimaryProgression(np.Progression, np.Tackle, np.Position, np.Archetype, SnapsPerGame, "Tackle", false)
			ManCoverage = PrimaryProgression(np.Progression, np.ManCoverage, np.Position, np.Archetype, SnapsPerGame, "Man Coverage", false)
			ZoneCoverage = PrimaryProgression(np.Progression, np.ZoneCoverage, np.Position, np.Archetype, SnapsPerGame, "Zone Coverage", false)
			Catching = PrimaryProgression(np.Progression, np.Catching, np.Position, np.Archetype, SnapsPerGame, "Catching", false)

			// Secondary
			ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
			ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
			PassRush = SecondaryProgression(np.Progression, np.PassRush)
			PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
			PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
			KickPower = SecondaryProgression(np.Progression, np.KickPower)
			KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
			Carrying = SecondaryProgression(np.Progression, np.Carrying)
			RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
			PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
			RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
		} else if np.Position == "K" {
			// Primary
			KickPower = PrimaryProgression(np.Progression, np.KickPower, np.Position, np.Archetype, SnapsPerGame, "Kick Power", false)
			KickAccuracy = PrimaryProgression(np.Progression, np.KickAccuracy, np.Position, np.Archetype, SnapsPerGame, "Kick Accuracy", false)
			// Secondary
			ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
			ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
			PassRush = SecondaryProgression(np.Progression, np.PassRush)
			RunDefense = SecondaryProgression(np.Progression, np.RunDefense)
			PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
			PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
			Carrying = SecondaryProgression(np.Progression, np.Carrying)
			RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
			PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
			Catching = SecondaryProgression(np.Progression, np.Catching)
			RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
			Strength = SecondaryProgression(np.Progression, np.Strength)
			Speed = SecondaryProgression(np.Progression, np.Speed)
			Agility = SecondaryProgression(np.Progression, np.Agility)
			ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
			ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
			FootballIQ = SecondaryProgression(np.Progression, np.FootballIQ)
			Tackle = SecondaryProgression(np.Progression, np.Tackle)
		} else if np.Position == "P" {
			// If David Ross
			if np.ID == 24984 {
				ThrowPower = PrimaryProgression(np.Progression, np.ThrowPower, np.Position, np.Archetype, SnapsPerGame, "Throw Power", false)
				ThrowAccuracy = PrimaryProgression(np.Progression, np.ThrowAccuracy, np.Position, np.Archetype, SnapsPerGame, "Throw Accuracy", false)
			} else {
				ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
				ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
			}
			// Primary Progressions
			PuntPower = PrimaryProgression(np.Progression, np.PuntPower, np.Position, np.Archetype, SnapsPerGame, "Punt Power", false)
			PuntAccuracy = PrimaryProgression(np.Progression, np.PuntAccuracy, np.Position, np.Archetype, SnapsPerGame, "Punt Accuracy", false)
			// Secondary Progressions
			RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
			PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
			RunDefense = SecondaryProgression(np.Progression, np.RunDefense)
			PassRush = SecondaryProgression(np.Progression, np.PassRush)
			Carrying = SecondaryProgression(np.Progression, np.Carrying)
			Tackle = SecondaryProgression(np.Progression, np.Tackle)
			RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
			Catching = SecondaryProgression(np.Progression, np.Catching)
			KickPower = SecondaryProgression(np.Progression, np.KickPower)
			KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
			ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
			ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
			Strength = SecondaryProgression(np.Progression, np.Strength)
			Speed = SecondaryProgression(np.Progression, np.Speed)
			Agility = SecondaryProgression(np.Progression, np.Agility)
			FootballIQ = SecondaryProgression(np.Progression, np.FootballIQ)
		}
	} else {
		// Change regression to be attribute - Normalized Random between 1-3
		ThrowPower = util.RegressValue(np.ThrowPower, 1, 3)
		ThrowAccuracy = util.RegressValue(np.ThrowAccuracy, 1, 3)
		RunBlock = util.RegressValue(np.RunBlock, 1, 3)
		PassBlock = util.RegressValue(np.PassBlock, 1, 3)
		RunDefense = util.RegressValue(np.RunDefense, 1, 3)
		PassRush = util.RegressValue(np.PassRush, 1, 3)
		Carrying = util.RegressValue(np.Carrying, 1, 3)
		Tackle = util.RegressValue(np.Tackle, 1, 3)
		RouteRunning = util.RegressValue(np.RouteRunning, 1, 3)
		Catching = util.RegressValue(np.Catching, 1, 3)
		KickPower = util.RegressValue(np.KickPower, 1, 3)
		KickAccuracy = util.RegressValue(np.KickAccuracy, 1, 3)
		ManCoverage = util.RegressValue(np.ManCoverage, 1, 3)
		ZoneCoverage = util.RegressValue(np.ZoneCoverage, 1, 3)
		Strength = util.RegressValue(np.Strength, 1, 3)
		Speed = util.RegressValue(np.Speed, 1, 3)
		Agility = util.RegressValue(np.Agility, 1, 3)
		FootballIQ = util.RegressValue(np.FootballIQ, 1, 3)
		PuntPower = util.RegressValue(np.PuntPower, 1, 3)
		PuntAccuracy = util.RegressValue(np.PuntAccuracy, 1, 3)
	}

	progressions := structs.CollegePlayerProgressions{
		Agility:       Agility,
		Catching:      Catching,
		Carrying:      Carrying,
		Speed:         Speed,
		RouteRunning:  RouteRunning,
		RunBlock:      RunBlock,
		PassBlock:     PassBlock,
		RunDefense:    RunDefense,
		PassRush:      PassRush,
		Strength:      Strength,
		Tackle:        Tackle,
		ThrowPower:    ThrowPower,
		ThrowAccuracy: ThrowAccuracy,
		PuntAccuracy:  PuntAccuracy,
		PuntPower:     PuntPower,
		KickAccuracy:  KickAccuracy,
		KickPower:     KickPower,
		FootballIQ:    FootballIQ,
		ManCoverage:   ManCoverage,
		ZoneCoverage:  ZoneCoverage,
	}

	np.Progress(progressions)

	np.GetOverall()

	return np
}

func BoomBustRecruit(np structs.CollegePlayer, SeasonID string, SnapsPerGame int, isBoom bool) structs.CollegePlayer {
	Agility := 0
	ThrowPower := 0
	ThrowAccuracy := 0
	Speed := 0
	FootballIQ := 0
	Strength := 0
	RunBlock := 0
	PassBlock := 0
	RunDefense := 0
	PassRush := 0
	Carrying := 0
	Tackle := 0
	RouteRunning := 0
	Catching := 0
	PuntPower := 0
	PuntAccuracy := 0
	KickPower := 0
	KickAccuracy := 0
	ManCoverage := 0
	ZoneCoverage := 0

	if isBoom {
		if np.Position == "QB" {
			// Primary Progressions
			Agility = PrimaryProgression(np.Progression, np.Agility, np.Position, np.Archetype, SnapsPerGame, "Agility", false)
			ThrowPower = PrimaryProgression(np.Progression, np.ThrowPower, np.Position, np.Archetype, SnapsPerGame, "Throw Power", false)
			ThrowAccuracy = PrimaryProgression(np.Progression, np.ThrowAccuracy, np.Position, np.Archetype, SnapsPerGame, "Throw Accuracy", false)
			Speed = PrimaryProgression(np.Progression, np.Speed, np.Position, np.Archetype, SnapsPerGame, "Speed", false)
			FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, np.Position, np.Archetype, SnapsPerGame, "Football IQ", false)
			Strength = PrimaryProgression(np.Progression, np.Strength, np.Position, np.Archetype, SnapsPerGame, "Strength", false)

			// Secondary Progressions
			RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
			PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
			RunDefense = SecondaryProgression(np.Progression, np.RunDefense)
			PassRush = SecondaryProgression(np.Progression, np.PassRush)
			Carrying = SecondaryProgression(np.Progression, np.Carrying)
			Tackle = SecondaryProgression(np.Progression, np.Tackle)
			RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
			Catching = SecondaryProgression(np.Progression, np.Catching)
			PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
			PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
			KickPower = SecondaryProgression(np.Progression, np.KickPower)
			KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
			ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
			ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
		} else if np.Position == "RB" {
			// Primary
			Agility = PrimaryProgression(np.Progression, np.Agility, np.Position, np.Archetype, SnapsPerGame, "Agility", false)
			Carrying = PrimaryProgression(np.Progression, np.Carrying, np.Position, np.Archetype, SnapsPerGame, "Carrying", false)
			Speed = PrimaryProgression(np.Progression, np.Speed, np.Position, np.Archetype, SnapsPerGame, "Speed", false)
			FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, np.Position, np.Archetype, SnapsPerGame, "Football IQ", false)
			Strength = PrimaryProgression(np.Progression, np.Strength, np.Position, np.Archetype, SnapsPerGame, "Strength", false)
			Catching = PrimaryProgression(np.Progression, np.Catching, np.Position, np.Archetype, SnapsPerGame, "Catching", false)
			PassBlock = PrimaryProgression(np.Progression, np.PassBlock, np.Position, np.Archetype, SnapsPerGame, "Pass Blocking", false)
			// Secondary
			RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
			RunDefense = SecondaryProgression(np.Progression, np.RunDefense)
			PassRush = SecondaryProgression(np.Progression, np.PassRush)
			Tackle = SecondaryProgression(np.Progression, np.Tackle)
			RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
			PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
			PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
			KickPower = SecondaryProgression(np.Progression, np.KickPower)
			KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
			ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
			ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
			ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
			ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
		} else if np.Position == "FB" {
			// Primary
			Agility = PrimaryProgression(np.Progression, np.Agility, np.Position, np.Archetype, SnapsPerGame, "Agility", false)
			Carrying = PrimaryProgression(np.Progression, np.Carrying, np.Position, np.Archetype, SnapsPerGame, "Carrying", false)
			Speed = PrimaryProgression(np.Progression, np.Speed, np.Position, np.Archetype, SnapsPerGame, "Speed", false)
			FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, np.Position, np.Archetype, SnapsPerGame, "Football IQ", false)
			Strength = PrimaryProgression(np.Progression, np.Strength, np.Position, np.Archetype, SnapsPerGame, "Strength", false)
			Catching = PrimaryProgression(np.Progression, np.Catching, np.Position, np.Archetype, SnapsPerGame, "Catching", false)
			PassBlock = PrimaryProgression(np.Progression, np.PassBlock, np.Position, np.Archetype, SnapsPerGame, "Pass Blocking", false)
			RunBlock = PrimaryProgression(np.Progression, np.RunBlock, np.Position, np.Archetype, SnapsPerGame, "Run Blocking", false)

			// Secondary
			RunDefense = SecondaryProgression(np.Progression, np.RunDefense)
			PassRush = SecondaryProgression(np.Progression, np.PassRush)
			Tackle = SecondaryProgression(np.Progression, np.Tackle)
			RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
			PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
			PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
			KickPower = SecondaryProgression(np.Progression, np.KickPower)
			KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
			ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
			ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
			ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
			ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)

		} else if np.Position == "TE" {
			// Primary
			Agility = PrimaryProgression(np.Progression, np.Agility, np.Position, np.Archetype, SnapsPerGame, "Agility", false)
			Carrying = PrimaryProgression(np.Progression, np.Carrying, np.Position, np.Archetype, SnapsPerGame, "Carrying", false)
			Speed = PrimaryProgression(np.Progression, np.Speed, np.Position, np.Archetype, SnapsPerGame, "Speed", false)
			FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, np.Position, np.Archetype, SnapsPerGame, "Football IQ", false)
			Strength = PrimaryProgression(np.Progression, np.Strength, np.Position, np.Archetype, SnapsPerGame, "Strength", false)
			Catching = PrimaryProgression(np.Progression, np.Catching, np.Position, np.Archetype, SnapsPerGame, "Catching", false)
			RouteRunning = PrimaryProgression(np.Progression, np.RouteRunning, np.Position, np.Archetype, SnapsPerGame, "Route Running", false)
			PassBlock = PrimaryProgression(np.Progression, np.PassBlock, np.Position, np.Archetype, SnapsPerGame, "Pass Blocking", false)
			RunBlock = PrimaryProgression(np.Progression, np.RunBlock, np.Position, np.Archetype, SnapsPerGame, "Run Blocking", false)

			// Secondary
			RunDefense = SecondaryProgression(np.Progression, np.RunDefense)
			PassRush = SecondaryProgression(np.Progression, np.PassRush)
			Tackle = SecondaryProgression(np.Progression, np.Tackle)
			PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
			PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
			KickPower = SecondaryProgression(np.Progression, np.KickPower)
			KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
			ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
			ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
			ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
			ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
		} else if np.Position == "WR" {
			// Primary
			Agility = PrimaryProgression(np.Progression, np.Agility, np.Position, np.Archetype, SnapsPerGame, "Agility", false)
			Carrying = PrimaryProgression(np.Progression, np.Carrying, np.Position, np.Archetype, SnapsPerGame, "Carrying", false)
			Speed = PrimaryProgression(np.Progression, np.Speed, np.Position, np.Archetype, SnapsPerGame, "Speed", false)
			FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, np.Position, np.Archetype, SnapsPerGame, "Football IQ", false)
			Strength = PrimaryProgression(np.Progression, np.Strength, np.Position, np.Archetype, SnapsPerGame, "Strength", false)
			Catching = PrimaryProgression(np.Progression, np.Catching, np.Position, np.Archetype, SnapsPerGame, "Catching", false)
			RouteRunning = PrimaryProgression(np.Progression, np.RouteRunning, np.Position, np.Archetype, SnapsPerGame, "Route Running", false)

			// Secondary
			RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
			PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
			RunDefense = SecondaryProgression(np.Progression, np.RunDefense)
			PassRush = SecondaryProgression(np.Progression, np.PassRush)
			Tackle = SecondaryProgression(np.Progression, np.Tackle)
			PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
			PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
			KickPower = SecondaryProgression(np.Progression, np.KickPower)
			KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
			ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
			ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
			ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
			ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
		} else if np.Position == "OT" || np.Position == "OG" || np.Position == "C" {
			// Primary
			Agility = PrimaryProgression(np.Progression, np.Agility, np.Position, np.Archetype, SnapsPerGame, "Agility", false)
			FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, np.Position, np.Archetype, SnapsPerGame, "Football IQ", false)
			Strength = PrimaryProgression(np.Progression, np.Strength, np.Position, np.Archetype, SnapsPerGame, "Strength", false)
			PassBlock = PrimaryProgression(np.Progression, np.PassBlock, np.Position, np.Archetype, SnapsPerGame, "Pass Blocking", false)
			RunBlock = PrimaryProgression(np.Progression, np.RunBlock, np.Position, np.Archetype, SnapsPerGame, "Run Blocking", false)

			// Secondary
			RunDefense = SecondaryProgression(np.Progression, np.RunDefense)
			PassRush = SecondaryProgression(np.Progression, np.PassRush)
			Tackle = SecondaryProgression(np.Progression, np.Tackle)
			RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
			PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
			PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
			KickPower = SecondaryProgression(np.Progression, np.KickPower)
			KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
			ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
			ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
			ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
			ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
			Carrying = SecondaryProgression(np.Progression, np.Carrying)
			Speed = SecondaryProgression(np.Progression, np.Speed)
			Catching = SecondaryProgression(np.Progression, np.Catching)
		} else if np.Position == "DT" {
			// Primary
			Agility = PrimaryProgression(np.Progression, np.Agility, np.Position, np.Archetype, SnapsPerGame, "Agility", false)
			FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, np.Position, np.Archetype, SnapsPerGame, "Football IQ", false)
			Strength = PrimaryProgression(np.Progression, np.Strength, np.Position, np.Archetype, SnapsPerGame, "Strength", false)
			PassRush = PrimaryProgression(np.Progression, np.PassRush, np.Position, np.Archetype, SnapsPerGame, "Pass Rush", false)
			RunDefense = PrimaryProgression(np.Progression, np.RunDefense, np.Position, np.Archetype, SnapsPerGame, "Run Defense", false)
			Tackle = PrimaryProgression(np.Progression, np.Tackle, np.Position, np.Archetype, SnapsPerGame, "Tackle", false)

			// Secondary
			PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
			PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
			KickPower = SecondaryProgression(np.Progression, np.KickPower)
			KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
			ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
			ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
			ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
			ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
			Carrying = SecondaryProgression(np.Progression, np.Carrying)
			Catching = SecondaryProgression(np.Progression, np.Catching)
			RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
			Speed = SecondaryProgression(np.Progression, np.Speed)
			RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
			PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
		} else if np.Position == "DE" {
			// Primary
			Agility = PrimaryProgression(np.Progression, np.Agility, np.Position, np.Archetype, SnapsPerGame, "Agility", false)
			FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, np.Position, np.Archetype, SnapsPerGame, "Football IQ", false)
			Strength = PrimaryProgression(np.Progression, np.Strength, np.Position, np.Archetype, SnapsPerGame, "Strength", false)
			Speed = PrimaryProgression(np.Progression, np.Speed, np.Position, np.Archetype, SnapsPerGame, "Speed", false)
			PassRush = PrimaryProgression(np.Progression, np.PassRush, np.Position, np.Archetype, SnapsPerGame, "Pass Rush", false)
			RunDefense = PrimaryProgression(np.Progression, np.RunDefense, np.Position, np.Archetype, SnapsPerGame, "Run Defense", false)
			Tackle = PrimaryProgression(np.Progression, np.Tackle, np.Position, np.Archetype, SnapsPerGame, "Tackle", false)

			// Secondary
			PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
			PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
			KickPower = SecondaryProgression(np.Progression, np.KickPower)
			KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
			ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
			ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
			ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
			ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
			Carrying = SecondaryProgression(np.Progression, np.Carrying)
			Catching = SecondaryProgression(np.Progression, np.Catching)
			RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
			RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
			PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
		} else if np.Position == "OLB" || np.Position == "ILB" {
			// Primary
			Agility = PrimaryProgression(np.Progression, np.Agility, np.Position, np.Archetype, SnapsPerGame, "Agility", false)
			FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, np.Position, np.Archetype, SnapsPerGame, "Football IQ", false)
			Strength = PrimaryProgression(np.Progression, np.Strength, np.Position, np.Archetype, SnapsPerGame, "Strength", false)
			Speed = PrimaryProgression(np.Progression, np.Speed, np.Position, np.Archetype, SnapsPerGame, "Speed", false)
			PassRush = PrimaryProgression(np.Progression, np.PassRush, np.Position, np.Archetype, SnapsPerGame, "Pass Rush", false)
			RunDefense = PrimaryProgression(np.Progression, np.RunDefense, np.Position, np.Archetype, SnapsPerGame, "Run Defense", false)
			Tackle = PrimaryProgression(np.Progression, np.Tackle, np.Position, np.Archetype, SnapsPerGame, "Tackle", false)
			ManCoverage = PrimaryProgression(np.Progression, np.ManCoverage, np.Position, np.Archetype, SnapsPerGame, "Man Coverage", false)
			ZoneCoverage = PrimaryProgression(np.Progression, np.ZoneCoverage, np.Position, np.Archetype, SnapsPerGame, "Zone Coverage", false)

			// Secondary
			PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
			PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
			KickPower = SecondaryProgression(np.Progression, np.KickPower)
			KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
			ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
			ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
			Carrying = SecondaryProgression(np.Progression, np.Carrying)
			Catching = SecondaryProgression(np.Progression, np.Catching)
			RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
			RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
			PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
		} else if np.Position == "CB" {
			// Primary
			Agility = PrimaryProgression(np.Progression, np.Agility, np.Position, np.Archetype, SnapsPerGame, "Agility", false)
			FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, np.Position, np.Archetype, SnapsPerGame, "Football IQ", false)
			Strength = PrimaryProgression(np.Progression, np.Strength, np.Position, np.Archetype, SnapsPerGame, "Strength", false)
			Speed = PrimaryProgression(np.Progression, np.Speed, np.Position, np.Archetype, SnapsPerGame, "Speed", false)
			Tackle = PrimaryProgression(np.Progression, np.Tackle, np.Position, np.Archetype, SnapsPerGame, "Tackle", false)
			ManCoverage = PrimaryProgression(np.Progression, np.ManCoverage, np.Position, np.Archetype, SnapsPerGame, "Man Coverage", false)
			ZoneCoverage = PrimaryProgression(np.Progression, np.ZoneCoverage, np.Position, np.Archetype, SnapsPerGame, "Zone Coverage", false)
			Catching = PrimaryProgression(np.Progression, np.Catching, np.Position, np.Archetype, SnapsPerGame, "Catching", false)

			// Secondary
			ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
			ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
			RunDefense = SecondaryProgression(np.Progression, np.RunDefense)
			PassRush = SecondaryProgression(np.Progression, np.PassRush)
			PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
			PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
			KickPower = SecondaryProgression(np.Progression, np.KickPower)
			KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
			Carrying = SecondaryProgression(np.Progression, np.Carrying)
			RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
			PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
			RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
		} else if np.Position == "FS" || np.Position == "SS" {
			// Primary
			Agility = PrimaryProgression(np.Progression, np.Agility, np.Position, np.Archetype, SnapsPerGame, "Agility", false)
			FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, np.Position, np.Archetype, SnapsPerGame, "Football IQ", false)
			Strength = PrimaryProgression(np.Progression, np.Strength, np.Position, np.Archetype, SnapsPerGame, "Strength", false)
			Speed = PrimaryProgression(np.Progression, np.Speed, np.Position, np.Archetype, SnapsPerGame, "Speed", false)
			RunDefense = PrimaryProgression(np.Progression, np.RunDefense, np.Position, np.Archetype, SnapsPerGame, "Run Defense", false)
			Tackle = PrimaryProgression(np.Progression, np.Tackle, np.Position, np.Archetype, SnapsPerGame, "Tackle", false)
			ManCoverage = PrimaryProgression(np.Progression, np.ManCoverage, np.Position, np.Archetype, SnapsPerGame, "Man Coverage", false)
			ZoneCoverage = PrimaryProgression(np.Progression, np.ZoneCoverage, np.Position, np.Archetype, SnapsPerGame, "Zone Coverage", false)
			Catching = PrimaryProgression(np.Progression, np.Catching, np.Position, np.Archetype, SnapsPerGame, "Catching", false)

			// Secondary
			ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
			ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
			PassRush = SecondaryProgression(np.Progression, np.PassRush)
			PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
			PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
			KickPower = SecondaryProgression(np.Progression, np.KickPower)
			KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
			Carrying = SecondaryProgression(np.Progression, np.Carrying)
			RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
			PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
			RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
		} else if np.Position == "K" {
			// Primary
			KickPower = PrimaryProgression(np.Progression, np.KickPower, np.Position, np.Archetype, SnapsPerGame, "Kick Power", false)
			KickAccuracy = PrimaryProgression(np.Progression, np.KickAccuracy, np.Position, np.Archetype, SnapsPerGame, "Kick Accuracy", false)
			// Secondary
			ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
			ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
			PassRush = SecondaryProgression(np.Progression, np.PassRush)
			RunDefense = SecondaryProgression(np.Progression, np.RunDefense)
			PuntPower = SecondaryProgression(np.Progression, np.PuntPower)
			PuntAccuracy = SecondaryProgression(np.Progression, np.PuntAccuracy)
			Carrying = SecondaryProgression(np.Progression, np.Carrying)
			RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
			PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
			Catching = SecondaryProgression(np.Progression, np.Catching)
			RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
			Strength = SecondaryProgression(np.Progression, np.Strength)
			Speed = SecondaryProgression(np.Progression, np.Speed)
			Agility = SecondaryProgression(np.Progression, np.Agility)
			ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
			ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
			FootballIQ = SecondaryProgression(np.Progression, np.FootballIQ)
			Tackle = SecondaryProgression(np.Progression, np.Tackle)
		} else if np.Position == "P" {
			// If David Ross
			if np.ID == 24984 {
				ThrowPower = PrimaryProgression(np.Progression, np.ThrowPower, np.Position, np.Archetype, SnapsPerGame, "Throw Power", false)
				ThrowAccuracy = PrimaryProgression(np.Progression, np.ThrowAccuracy, np.Position, np.Archetype, SnapsPerGame, "Throw Accuracy", false)
			} else {
				ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
				ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
			}
			// Primary Progressions
			PuntPower = PrimaryProgression(np.Progression, np.PuntPower, np.Position, np.Archetype, SnapsPerGame, "Punt Power", false)
			PuntAccuracy = PrimaryProgression(np.Progression, np.PuntAccuracy, np.Position, np.Archetype, SnapsPerGame, "Punt Accuracy", false)
			// Secondary Progressions
			RunBlock = SecondaryProgression(np.Progression, np.RunBlock)
			PassBlock = SecondaryProgression(np.Progression, np.PassBlock)
			RunDefense = SecondaryProgression(np.Progression, np.RunDefense)
			PassRush = SecondaryProgression(np.Progression, np.PassRush)
			Carrying = SecondaryProgression(np.Progression, np.Carrying)
			Tackle = SecondaryProgression(np.Progression, np.Tackle)
			RouteRunning = SecondaryProgression(np.RouteRunning, np.RouteRunning)
			Catching = SecondaryProgression(np.Progression, np.Catching)
			KickPower = SecondaryProgression(np.Progression, np.KickPower)
			KickAccuracy = SecondaryProgression(np.Progression, np.KickAccuracy)
			ManCoverage = SecondaryProgression(np.Progression, np.ManCoverage)
			ZoneCoverage = SecondaryProgression(np.Progression, np.ZoneCoverage)
			Strength = SecondaryProgression(np.Progression, np.Strength)
			Speed = SecondaryProgression(np.Progression, np.Speed)
			Agility = SecondaryProgression(np.Progression, np.Agility)
			FootballIQ = SecondaryProgression(np.Progression, np.FootballIQ)
		} else if np.Position == "ATH" {
			ThrowPower = np.ThrowPower + util.GenerateNormalizedIntFromRange(1, 3)
			ThrowAccuracy = np.ThrowAccuracy + util.GenerateNormalizedIntFromRange(1, 3)
			PuntPower = np.PuntPower + util.GenerateNormalizedIntFromRange(1, 3)
			PuntAccuracy = np.PuntAccuracy + util.GenerateNormalizedIntFromRange(1, 3)
			RunBlock = np.RunBlock + util.GenerateNormalizedIntFromRange(1, 3)
			PassBlock = np.PassBlock + util.GenerateNormalizedIntFromRange(1, 3)
			RunDefense = np.RunDefense + util.GenerateNormalizedIntFromRange(1, 3)
			PassRush = np.PassRush + util.GenerateNormalizedIntFromRange(1, 3)
			Carrying = np.Carrying + util.GenerateNormalizedIntFromRange(1, 3)
			Tackle = np.Tackle + util.GenerateNormalizedIntFromRange(1, 3)
			RouteRunning = np.RouteRunning + util.GenerateNormalizedIntFromRange(1, 3)
			Catching = np.Catching + util.GenerateNormalizedIntFromRange(1, 3)
			KickPower = np.KickPower + util.GenerateNormalizedIntFromRange(1, 3)
			KickAccuracy = np.KickAccuracy + util.GenerateNormalizedIntFromRange(1, 3)
			ManCoverage = np.ManCoverage + util.GenerateNormalizedIntFromRange(1, 3)
			ZoneCoverage = np.ZoneCoverage + util.GenerateNormalizedIntFromRange(1, 3)
			Strength = np.Strength + util.GenerateNormalizedIntFromRange(1, 3)
			Speed = np.Speed + util.GenerateNormalizedIntFromRange(1, 3)
			Agility = np.Agility + util.GenerateNormalizedIntFromRange(1, 3)
			FootballIQ = np.FootballIQ + util.GenerateNormalizedIntFromRange(1, 3)
		}
	} else {
		// Change regression to be attribute - Normalized Random between 1-3
		ThrowPower = util.RegressValue(np.ThrowPower, 1, 2)
		ThrowAccuracy = util.RegressValue(np.ThrowAccuracy, 1, 2)
		RunBlock = util.RegressValue(np.RunBlock, 1, 2)
		PassBlock = util.RegressValue(np.PassBlock, 1, 2)
		RunDefense = util.RegressValue(np.RunDefense, 1, 2)
		PassRush = util.RegressValue(np.PassRush, 1, 2)
		Carrying = util.RegressValue(np.Carrying, 1, 2)
		Tackle = util.RegressValue(np.Tackle, 1, 2)
		RouteRunning = util.RegressValue(np.RouteRunning, 1, 2)
		Catching = util.RegressValue(np.Catching, 1, 2)
		KickPower = util.RegressValue(np.KickPower, 1, 2)
		KickAccuracy = util.RegressValue(np.KickAccuracy, 1, 2)
		ManCoverage = util.RegressValue(np.ManCoverage, 1, 2)
		ZoneCoverage = util.RegressValue(np.ZoneCoverage, 1, 2)
		Strength = util.RegressValue(np.Strength, 1, 2)
		Speed = util.RegressValue(np.Speed, 1, 2)
		Agility = util.RegressValue(np.Agility, 1, 2)
		FootballIQ = util.RegressValue(np.FootballIQ, 1, 2)
		PuntPower = util.RegressValue(np.PuntPower, 1, 2)
		PuntAccuracy = util.RegressValue(np.PuntAccuracy, 1, 2)
	}

	progressions := structs.CollegePlayerProgressions{
		Agility:       Agility,
		Catching:      Catching,
		Carrying:      Carrying,
		Speed:         Speed,
		RouteRunning:  RouteRunning,
		RunBlock:      RunBlock,
		PassBlock:     PassBlock,
		RunDefense:    RunDefense,
		PassRush:      PassRush,
		Strength:      Strength,
		Tackle:        Tackle,
		ThrowPower:    ThrowPower,
		ThrowAccuracy: ThrowAccuracy,
		PuntAccuracy:  PuntAccuracy,
		PuntPower:     PuntPower,
		KickAccuracy:  KickAccuracy,
		KickPower:     KickPower,
		FootballIQ:    FootballIQ,
		ManCoverage:   ManCoverage,
		ZoneCoverage:  ZoneCoverage,
	}

	np.Progress(progressions, true)

	np.GetOverall()

	return np
}

func PrimaryProgression(progression int, input int, position string, archetype string, spg int, attribute string, isRedshirting bool) int {
	if input == 0 {
		return 1
	}

	modifier := GetModifiers(position, archetype, spg, attribute)

	var progress float64 = 0

	// Is Redshirting or Practice Squad Check
	if !isRedshirting {
		progress = ((1 - math.Pow((float64(input)/99.0), 15)) * math.Log10(float64(input)) * (0.6 + modifier)) * (1 + (float64(progression) / 70))
	} else {
		rsMod := util.GenerateFloatFromRange(1.7, 2.3)
		progress = ((1 - math.Pow((float64(input)/99), 15)) * math.Log10(float64(input)) * rsMod * (1 + (float64(progression / 60))))
	}

	if progress+float64(input) > 99 {
		progress = 99
	} else {
		progress = progress + float64(input)
	}

	return int(math.Round(progress))
}

func SecondaryProgression(progression int, input int) int {
	num := rand.Intn(99)

	if num < progression && input < 99 {
		newInput := input + 1
		return newInput
	} else {
		return input
	}
}

func GetModifiers(position string, archetype string, spg int, attrib string) float64 {
	var snapMod float64 = 0.0
	if spg > 50 {
		snapMod = rand.Float64()*(1.25-1) + 1
	} else if spg > 30 {
		snapMod = rand.Float64()*(1.1-0.9) + 0.9
	} else if spg > 20 {
		snapMod = rand.Float64()*(1.0-0.75) + 0.75
	} else if spg > 10 {
		snapMod = rand.Float64()*(.9-0.6) + 0.6
	} else {
		snapMod = rand.Float64()*(0.75-0.5) + 0.5
	}
	archMod := GetArchetypeMod(position, archetype, attrib)
	return snapMod + archMod
}

func GetArchetypeMod(pos string, arch string, attribute string) float64 {
	if pos == "QB" {
		if arch == "Pocket" {
			if attribute == "Throw Power" || attribute == "Throw Accuracy" {
				return 0.05
			}
		} else if arch == "Balanced" {
			if attribute == "Agility" {
				return 0.05
			} else if attribute == "Throw Power" || attribute == "Throw Accuracy" {
				return 0.025
			}

		} else if arch == "Field General" {
			if attribute == "Football IQ" || attribute == "Throw Accuracy" {
				return 0.05
			}
		} else if arch == "Scrambler" {
			if attribute == "Agility" || attribute == "Speed" {
				return 0.05
			}
			if attribute == "Throw Power" || attribute == "Throw Accuracy" {
				return 0.0175
			}
		}
		return 0
	} else if pos == "RB" {
		if arch == "Speed" {
			if attribute == "Speed" || attribute == "Carrying" {
				return 0.05
			}
			return 0.0
		} else if arch == "Power" {
			if attribute == "Carrying" || attribute == "Strength" {
				return 0.05
			}
			return 0.0
		} else if arch == "Balanced" {
			if attribute == "Agility" || attribute == "Pass Block" {
				return 0.05
			}
			return 0.0
		} else if arch == "Receiving" {
			if attribute == "Agility" || attribute == "Catching" {
				return 0.05
			}
			return 0.0
		}
		return 0
	} else if pos == "FB" {
		if arch == "Blocking" && (attribute == "Pass Blocking" || attribute == "Run Blocking") {
			return 0.05
		} else if arch == "Rushing" && (attribute == "Carrying" || attribute == "Strength") {
			return 0.05
		} else if arch == "Balanced" && (attribute == "Agility" || attribute == "Speed") {
			return 0.05
		} else if arch == "Receiving" && (attribute == "Agility" || attribute == "Catching") {
			return 0.05
		}
		return 0
	} else if pos == "WR" {
		if arch == "Possession" && (attribute == "Carrying" || attribute == "Catching") {
			return 0.05
		} else if arch == "Route Runner" && (attribute == "Agility" || attribute == "Route Running") {
			return 0.05
		} else if arch == "Speed" && attribute == "Speed" {
			return 0.1
		} else if arch == "Red Zone Threat" && (attribute == "Strength" || attribute == "Catching") {
			return 0.05
		}
		return 0
	} else if pos == "TE" {
		if arch == "Blocking" && (attribute == "Pass Blocking" || attribute == "Run Blocking") {
			return 0.05
		} else if arch == "Vertical Threat" && (attribute == "Catching" || attribute == "Speed") {
			return 0.05
		} else if arch == "Receiving" && (attribute == "Agility" || attribute == "Catching") {
			return 0.05
		}
		return 0
	} else if pos == "OT" {
		if arch == "Pass Blocking" && (attribute == "Pass Blocking" || attribute == "Strength") {
			return 0.05
		} else if arch == "Run Blocking" && (attribute == "Run Blocking" || attribute == "Strength") {
			return 0.05
		} else if arch == "Balanced" {
			if attribute == "Run Blocking" || attribute == "Pass Blocking" {
				return 0.025
			} else if attribute == "Strength" {
				return 0.05
			}
		}
	} else if pos == "OG" {
		if arch == "Pass Blocking" && (attribute == "Pass Blocking" || attribute == "Strength") {
			return 0.05
		} else if arch == "Run Blocking" && (attribute == "Run Blocking" || attribute == "Strength") {
			return 0.05
		} else if arch == "Balanced" {
			if attribute == "Run Blocking" || attribute == "Pass Blocking" {
				return 0.025
			} else if attribute == "Strength" {
				return 0.05
			}
		}
	} else if pos == "C" {
		if arch == "Pass Blocking" && (attribute == "Pass Blocking" || attribute == "Strength") {
			return 0.05
		} else if arch == "Run Blocking" && (attribute == "Run Blocking" || attribute == "Strength") {
			return 0.05
		} else if arch == "Balanced" {
			if attribute == "Run Blocking" || attribute == "Pass Blocking" {
				return 0.025
			} else if attribute == "Strength" {
				return 0.05
			}
		} else if arch == "Line Captain" {
			if attribute == "Run Blocking" || attribute == "Pass Blocking" {
				return 0.025
			} else if attribute == "Football IQ" {
				return 0.05
			}
		}
	} else if pos == "DT" {
		if arch == "Pass Rusher" && (attribute == "Pass Rush" || attribute == "Strength") {
			return 0.05
		} else if arch == "Nose Tackle" && (attribute == "Run Defense" || attribute == "Strength") {
			return 0.05
		} else if arch == "Balanced" && (attribute == "Tackle" || attribute == "Strength") {
			return 0.05
		}
	} else if pos == "DE" {
		if arch == "Speed Rusher" && attribute == "Pass Rush" {
			return 0.1
		} else if arch == "Run Stopper" && attribute == "Run Defense" {
			return 0.1
		} else if arch == "Balanced" && (attribute == "Pass Rush" || attribute == "Run Defense") {
			return 0.05
		}
	} else if pos == "ILB" {
		if arch == "Field General" && (attribute == "Football IQ" || attribute == "Tackle") {
			return 0.05
		} else if arch == "Run Stopper" && (attribute == "Run Defense" || attribute == "Strength") {
			return 0.05
		} else if arch == "Coverage" && (attribute == "Zone Coverage" || attribute == "Man Coverage") {
			return 0.05
		} else if arch == "Speed" && (attribute == "Speed" || attribute == "Agility") {
			return 0.05
		}
	} else if pos == "OLB" {
		if arch == "Pass Rush" && attribute == "Pass Rush" {
			return 0.1
		} else if arch == "Run Stopper" && (attribute == "Run Defense" || attribute == "Strength") {
			return 0.05
		} else if arch == "Coverage" && (attribute == "Zone Coverage" || attribute == "Man Coverage") {
			return 0.05
		} else if arch == "Speed" && (attribute == "Speed" || attribute == "Agility") {
			return 0.05
		}
	} else if pos == "CB" {
		if arch == "Man Coverage" && (attribute == "Man Coverage" || attribute == "Speed") {
			return 0.05
		} else if arch == "Zone Coverage" && (attribute == "Zone Coverage" || attribute == "Football IQ") {
			return 0.05
		} else if arch == "Ball Hawk" {
			if attribute == "Man Coverage" || attribute == "Zone Coverage" {
				return 0.025
			} else if attribute == "Catching" {
				return 0.05
			}
		}
	} else if pos == "FS" || pos == "SS" {
		if arch == "Man Coverage" && (attribute == "Man Coverage" || attribute == "Speed") {
			return 0.05
		} else if arch == "Zone Coverage" && (attribute == "Zone Coverage" || attribute == "Football IQ") {
			return 0.05
		} else if arch == "Ball Hawk" {
			if attribute == "Man Coverage" || attribute == "Zone Coverage" {
				return 0.025
			} else if attribute == "Catching" {
				return 0.05
			}
		} else if arch == "Run Stopper" && (attribute == "Run Defense" || attribute == "Tackle") {
			return 0.05
		}
	} else if pos == "K" {
		if arch == "Power" && attribute == "Kick Power" {
			return 0.1
		} else if arch == "Accuracy" && attribute == "Kick Accuracy" {
			return 0.1
		} else if arch == "Balanced" && (attribute == "Kick Power" || attribute == "Kick Accuracy") {
			return 0.05
		}
	} else if pos == "P" {
		if arch == "Power" && attribute == "Punt Power" {
			return 0.1
		} else if arch == "Accuracy" && attribute == "Punt Accuracy" {
			return 0.1
		} else if arch == "Balanced" && (attribute == "Punt Power" || attribute == "Punt Accuracy") {
			return 0.05
		}
	} else if pos == "ATH" {
		return 0.0075
	}
	return 0
}

func DetermineIfRetiring(player structs.NFLPlayer, statMap map[uint][]structs.NFLPlayerSeasonStats) bool {
	if player.IsFreeAgent && player.Experience > 1 {
		lastTwoSeasonStats := statMap[player.ID]
		totalSnaps := 0
		for _, stat := range lastTwoSeasonStats {
			totalSnaps += stat.Snaps
		}
		return totalSnaps == 0
	}

	if player.Age < int(player.PrimeAge) {
		return false
	}

	/*
		Thoughts - we could implement historic injuries into this somewhere, although we are impacting prime age upon injuries.
	*/
	benchmark := 0
	age := player.Age
	primeAge := player.PrimeAge
	retirementAge := primeAge + uint(util.GenerateIntFromRange(3, 5))
	if age > int(retirementAge) {
		benchmark += 50
	}
	if age > int(primeAge) && player.Overall < 56 {
		benchmark += (10*age - int(primeAge))
	}
	diceRoll := util.GenerateIntFromRange(1, 100)
	// If the roll is less than the benchmark, player will retire. Otherwise, they are staying.
	return diceRoll < benchmark
}

func DetermineIfDeclaring(player structs.CollegePlayer, avgSnaps int) bool {
	// Redshirt senior or just a senior
	if (player.IsRedshirt && player.Year == 5) || (!player.IsRedshirt && player.Year == 4 && !player.IsRedshirting) {
		return true
	}
	ovr := player.Overall

	// YEAR 3 and NOT redshirt == Junior
	// Year 4 AND redshirt == Redshirt Junior
	// Year 3 AND redshirt == Redshirt Sophomore
	// Year 2 AND redshirt == Redshirt Freshmen
	// All players that are freshmen, redshirt freshmen, sophomore, redshirt sophomores, and non-redshirt juniors
	isRedshirtJunior := player.Year == 4 && player.IsRedshirt
	if !isRedshirtJunior {
		return false
	}

	if ovr < 55 || player.IsRedshirting {
		return false
	}

	snapMod := 0
	if avgSnaps > 50 {
		snapMod = 16
	} else if avgSnaps > 30 {
		snapMod = 12
	} else if avgSnaps > 20 {
		snapMod = 8
	} else if avgSnaps > 10 {
		snapMod = 4
	}

	// Dice Roll
	odds := util.GenerateIntFromRange(1, 100) - snapMod
	if ovr > 54 && odds <= 25 {
		return true
	} else if ovr > 56 && odds <= 30 {
		return true
	} else if ovr > 57 && odds <= 35 {
		return true
	} else if ovr > 58 && odds <= 45 {
		return true
	} else if ovr > 59 && odds <= 70 {
		return true
	} else if ovr > 60 && odds <= 75 {
		return true
	} else if ovr > 61 && odds <= 80 {
		return true
	} else if ovr > 62 && odds <= 95 {
		return true
	}
	return false
}

func getAverageNFLSnaps(stats []structs.NFLPlayerStats) (int, int) {
	totalSnaps := 0

	for _, stat := range stats {
		totalSnaps += stat.Snaps
	}

	var SnapsPerGame int = 0
	if len(stats) > 0 {
		SnapsPerGame = totalSnaps / 17
	}

	return totalSnaps, SnapsPerGame
}

func getAverageSnaps(stats []structs.CollegePlayerStats) int {
	totalSnaps := 0

	for _, stat := range stats {
		totalSnaps += stat.Snaps
	}

	var SnapsPerGame int = 0
	if len(stats) > 0 {
		SnapsPerGame = totalSnaps / 12 // 12
	}

	return SnapsPerGame
}

func getMostPlayedPosition(snaps structs.BasePlayerSeasonSnaps, pos string) (string, uint16) {
	position := ""
	mostSnaps := uint16(0)

	if snaps.QBSnaps > mostSnaps {
		mostSnaps = snaps.QBSnaps
		position = "QB"
	}

	if snaps.RBSnaps > mostSnaps {
		mostSnaps = snaps.RBSnaps
		position = "RB"
	}

	if snaps.FBSnaps > mostSnaps {
		mostSnaps = snaps.FBSnaps
		position = "FB"
	}

	if snaps.WRSnaps > mostSnaps {
		mostSnaps = snaps.WRSnaps
		position = "WR"
	}

	if snaps.TESnaps > mostSnaps {
		mostSnaps = snaps.TESnaps
		position = "TE"
	}

	if snaps.OTSnaps > mostSnaps {
		mostSnaps = snaps.OTSnaps
		position = "OT"
	}

	if snaps.OGSnaps > mostSnaps {
		mostSnaps = snaps.OGSnaps
		position = "OG"
	}

	if snaps.CSnaps > mostSnaps {
		mostSnaps = snaps.CSnaps
		position = "C"
	}

	if snaps.DTSnaps > mostSnaps {
		mostSnaps = snaps.DTSnaps
		position = "DT"
	}

	if snaps.DESnaps > mostSnaps {
		mostSnaps = snaps.DESnaps
		position = "DE"
	}

	if snaps.OLBSnaps > mostSnaps {
		mostSnaps = snaps.OLBSnaps
		position = "OLB"
	}

	if snaps.ILBSnaps > mostSnaps {
		mostSnaps = snaps.ILBSnaps
		position = "ILB"
	}

	if snaps.CBSnaps > mostSnaps {
		mostSnaps = snaps.CBSnaps
		position = "CB"
	}

	if snaps.FSSnaps > mostSnaps {
		mostSnaps = snaps.FSSnaps
		position = "FS"
	}

	if snaps.SSSnaps > mostSnaps {
		mostSnaps = snaps.SSSnaps
		position = "SS"
	}

	if snaps.KSnaps > mostSnaps {
		mostSnaps = snaps.KSnaps
		position = "K"
	}

	if snaps.PSnaps > mostSnaps {
		mostSnaps = snaps.PSnaps
		position = "P"
	}

	if snaps.STSnaps > mostSnaps {
		return pos, 0
	}

	return position, mostSnaps
}

func getNewArchetype(pos, arch, newPos string) (string, bool) {
	newArchetypeMap := util.GetNewArchetypeMap()
	posMap, posExists := newArchetypeMap[pos]
	if !posExists {
		return "", false
	}
	oriArchMap, archExists := posMap[arch]
	if !archExists {
		return "", false
	}
	newArchtype, newPosExists := oriArchMap[newPos]
	if !newPosExists {
		return "", false
	}

	return newArchtype, true
}

func RegressAttribute(attr, age, primeAge int, isIQ bool) int {
	min := 0.0
	max := 0.0
	ageRequirement := primeAge
	mod := float64(primeAge) - 5.5
	if isIQ {
		ageRequirement = primeAge + 1
		mod = float64(primeAge)
	}
	if age < ageRequirement {
		max = 0.02
	} else {
		percent := float64(age) - mod
		min = percent / 100.0
		max = (percent + 5) / 100.0
	}

	regression := util.GenerateFloatFromRange(min, max)
	newAttrValue := float64(attr) * (1 - regression)

	return int(newAttrValue)
}
