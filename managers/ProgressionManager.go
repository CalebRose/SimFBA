package managers

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/models"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
)

func ProgressionMain() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	SeasonID := strconv.Itoa(ts.CollegeSeasonID)
	fmt.Println(time.Now().UnixNano())
	rand.Seed(time.Now().UnixNano())
	// Get All Teams

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
				stats := GetCollegePlayerStatsByPlayerIDAndSeason(strconv.Itoa(int(player.ID)), SeasonID)

				// Get Average Snaps
				avgSnaps := getAverageSnaps(stats)

				// Run Function to Determine if Player is Declaring Early
				willDeclare := DetermineIfDeclaring(player, avgSnaps)

				// Progress the Player
				player = ProgressCollegePlayer(player, SeasonID, stats)

				if willDeclare {
					player.GraduatePlayer()
					draftee := models.NFLDraftee{}
					draftee.Map(player)
					// Map New Progression value for NFL
					newProgression := util.GeneratePotential()
					newPotentialGrade := util.GetWeightedPotentialGrade(newProgression)
					draftee.MapProgression(newProgression, newPotentialGrade)

					// Create Historic Player Record
					hcp := (structs.HistoricCollegePlayer)(player)

					err := db.Create(&hcp).Error
					if err != nil {
						log.Panicln("Could not save historic player record!")
					}

					message := player.Position + " " + player.FirstName + " " + player.LastName + " has graduated from " + player.TeamAbbr + "!"
					if (player.Year < 5 && player.IsRedshirt) || (player.Year < 4 && !player.IsRedshirt) {
						message = "Breaking News! " + player.Position + " " + player.FirstName + " " + player.LastName + " is declaring early from " + player.TeamAbbr + ", and will be eligible for the SimNFL Draft!"
					}
					CreateNewsLog("CFB", message, "Graduation", player.TeamID, ts)

					// Create Draftee Record
					err = db.Create(&draftee).Error
					if err != nil {
						log.Panicln("Could not save graduating players")
					}
					graduatingPlayers = append(graduatingPlayers, draftee)
					// CollegePlayer record will be deleted, but record will be mapped to a GraduatedCollegePlayer struct, and then saved in that table, along side with NFL Draftees table
					// GraduatedCollegePlayer will be a copy of the collegeplayers table, but only for historical players

					err = db.Delete(&player).Error
					if err != nil {
						log.Panicln("Could not delete old college player record.")
					}
					continue
				}
				fmt.Println("Saved " + player.FirstName + " " + player.LastName + "'s record")
				err := db.Save(&player).Error
				if err != nil {
					log.Panicln("Could not save player record")
				}

			}

			team.TogglePlayersProgressed()
		}

		if !team.RecruitsAdded {
			for _, croot := range croots {
				// Convert to College Player Record
				cp := structs.CollegePlayer{}
				cp.MapFromRecruit(croot, team)

				fmt.Println("Adding " + croot.FirstName + " " + croot.LastName + "to " + team.TeamAbbr)

				// Save College Player Record
				err := db.Create(&cp).Error
				if err != nil {
					log.Panicln("Could not save new college player record")
				}

				// Delete Recruit Record
				err = db.Delete(&croot).Error
				if err != nil {
					log.Panicln("Could not save recruit record")
				}
			}

			team.ToggleRecruitsAdded()
		}

		db.Save(&team)

	}

	// Unsigned Players
	unsignedPlayers := GetAllUnsignedPlayers()
	for _, player := range unsignedPlayers {
		player = ProgressUnsignedPlayer(player, SeasonID)
		if (player.IsRedshirt && player.Year > 5) ||
			(!player.IsRedshirt && player.Year > 4) {
			player.GraduatePlayer()
			draftee := models.NFLDraftee{}
			draftee.MapUnsignedPlayer(player)
			hcp := structs.HistoricCollegePlayer{}
			hcp.MapUnsignedPlayer(player)

			err := db.Create(&hcp).Error
			if err != nil {
				log.Panicln("Could not save historic player record!")
			}
			graduatingPlayers = append(graduatingPlayers, draftee)
			// CollegePlayer record will be deleted, but record will be mapped to a GraduatedCollegePlayer struct, and then saved in that table, along side with NFL Draftees table
			// GraduatedCollegePlayer will be a copy of the collegeplayers table, but only for historical players

			err = db.Delete(&player).Error
			if err != nil {
				log.Panicln("Could not delete old college player record.")
			}
			continue
		}
		err := db.Save(&player).Error
		if err != nil {
			log.Panicln("Could not save player record")
		}
	}

	// Graduating players
	for _, grad := range graduatingPlayers {
		err := db.Create(&grad).Error
		if err != nil {
			log.Panicln("Could not save graduating players")
		}
	}
	// get all unsigned players
	// progress through all unsigned players
	// move all seniors + to graduates table
	// move all unsigned croots to unsigned players table

	unsignedCroots := GetAllUnsignedRecruits()
	for _, croot := range unsignedCroots {
		// Unsigned Players
		up := structs.UnsignedPlayer{}

		up.MapFromRecruit(croot)

		err := db.Create(&up).Error
		if err != nil {
			log.Panic("Unsigned player could not be created!")
		}
		err = db.Delete(&croot).Error
		if err != nil {
			log.Panic("Recruit could not be deleted!")
		}
	}
}

func NFLProgressionMain() {
	ProgressNFLPlayers()
	AllocateRetiredContracts()
}

func ProgressNFLPlayers() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	// SeasonID := strconv.Itoa(ts.NFLSeasonID - 1)
	SeasonID := "2"
	fmt.Println(time.Now().UnixNano())
	rand.Seed(time.Now().UnixNano())

	teams := GetAllNFLTeams()
	// nflPlayers := GetAllNFLPlayers()
	freeAgents := GetAllFreeAgents()
	// waivedPlayers := GetAllWaiverWirePlayers()

	for _, team := range teams {
		teamID := strconv.Itoa(int(team.ID))
		nflPlayers := GetNFLPlayersRecordsByTeamID(teamID)
		for _, player := range nflPlayers {
			if player.HasProgressed {
				continue
			}

			if player.ID == 387 || player.ID == 67924 ||
				player.ID == 55373 || player.ID == 41718 || player.ID == 57100 || player.ID == 60489 || player.ID == 45264 ||
				player.ID == 30539 || player.ID == 5332 || player.ID == 60148 || player.ID == 45417 || player.ID == 51262 || player.ID == 65418 ||
				player.ID == 5257 || player.ID == 62249 {
				fmt.Println("STOP HERE")
			}

			// Progress the Player
			// Get Latest Stats
			stats := GetNFLPlayerStatsByPlayerIDAndSeason(strconv.Itoa(int(player.ID)), SeasonID)

			// Get Average Snaps
			avgSnaps := getAverageNFLSnaps(stats)

			// Run Function to Determine if Player is Declaring Early
			willRetire := DetermineIfRetiring(player)

			// Progress the Player
			player = ProgressNFLPlayer(player, SeasonID, avgSnaps)

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
				db.Save(&activeContract)
			}
			if !willRetire {
				player.ToggleHasProgressed()
				db.Save(&player)
				continue
			}

			// Retire Player
			message := "Breaking News: " + player.Position + " " + player.FirstName + " " + player.LastName + " has decided to retire from SimNFL. He was drafted by " + player.DraftedTeam + " and last played with " + player.TeamAbbr + " and " + player.PreviousTeam + ". We thank him for his wondrous, extensive career and hope he enjoys his retirement!"
			CreateNewsLog("NFL", message, "Retirement", player.TeamID, ts)
			retiredPlayer := (structs.NFLRetiredPlayer)(player)
			db.Create(&retiredPlayer)
			db.Delete(&player)
		}
	}

	for _, player := range freeAgents {
		if player.HasProgressed {
			continue
		}

		// Progress the Player
		// Get Latest Stats
		stats := GetNFLPlayerStatsByPlayerIDAndSeason(strconv.Itoa(int(player.ID)), SeasonID)

		// Get Average Snaps
		avgSnaps := getAverageNFLSnaps(stats)

		// Run Function to Determine if Player is Declaring Early
		willRetire := DetermineIfRetiring(player)

		// Progress the Player
		player = ProgressNFLPlayer(player, SeasonID, avgSnaps)

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
			db.Save(&activeContract)
		}
		if !willRetire {
			player.ToggleHasProgressed()
			db.Save(&player)
			continue
		}

		// Retire Player
		message := "Breaking News: " + player.Position + " " + player.FirstName + " " + player.LastName + " has decided to retire from SimNFL. He was drafted by " + player.DraftedTeam + " and last played with " + player.TeamAbbr + " and " + player.PreviousTeam + ". We thank him for his wondrous, extensive career and hope he enjoys his retirement!"
		CreateNewsLog("NFL", message, "Retirement", player.TeamID, ts)
		retiredPlayer := (structs.NFLRetiredPlayer)(player)
		db.Create(&retiredPlayer)
		db.Delete(&player)
	}
}

func ProgressNFLPlayer(np structs.NFLPlayer, SeasonID string, SnapsPerGame int) structs.NFLPlayer {
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

	if np.Position == "QB" {
		// Primary Progressions
		Agility = PrimaryProgression(np.Progression, np.Agility, np.Position, np.Archetype, SnapsPerGame, "Agility", np.IsPracticeSquad)
		ThrowPower = PrimaryProgression(np.Progression, np.ThrowPower, np.Position, np.Archetype, SnapsPerGame, "Throw Power", np.IsPracticeSquad)
		ThrowAccuracy = PrimaryProgression(np.Progression, np.ThrowAccuracy, np.Position, np.Archetype, SnapsPerGame, "Throw Accuracy", np.IsPracticeSquad)
		Speed = PrimaryProgression(np.Progression, np.Speed, np.Position, np.Archetype, SnapsPerGame, "Speed", np.IsPracticeSquad)
		FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, np.Position, np.Archetype, SnapsPerGame, "Football IQ", np.IsPracticeSquad)
		Strength = PrimaryProgression(np.Progression, np.Strength, np.Position, np.Archetype, SnapsPerGame, "Strength", np.IsPracticeSquad)

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
		Agility = PrimaryProgression(np.Progression, np.Agility, np.Position, np.Archetype, SnapsPerGame, "Agility", np.IsPracticeSquad)
		Carrying = PrimaryProgression(np.Progression, np.Carrying, np.Position, np.Archetype, SnapsPerGame, "Carrying", np.IsPracticeSquad)
		Speed = PrimaryProgression(np.Progression, np.Speed, np.Position, np.Archetype, SnapsPerGame, "Speed", np.IsPracticeSquad)
		FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, np.Position, np.Archetype, SnapsPerGame, "Football IQ", np.IsPracticeSquad)
		Strength = PrimaryProgression(np.Progression, np.Strength, np.Position, np.Archetype, SnapsPerGame, "Strength", np.IsPracticeSquad)
		Catching = PrimaryProgression(np.Progression, np.Catching, np.Position, np.Archetype, SnapsPerGame, "Catching", np.IsPracticeSquad)
		PassBlock = PrimaryProgression(np.Progression, np.PassBlock, np.Position, np.Archetype, SnapsPerGame, "Pass Blocking", np.IsPracticeSquad)
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
		Agility = PrimaryProgression(np.Progression, np.Agility, np.Position, np.Archetype, SnapsPerGame, "Agility", np.IsPracticeSquad)
		Carrying = PrimaryProgression(np.Progression, np.Carrying, np.Position, np.Archetype, SnapsPerGame, "Carrying", np.IsPracticeSquad)
		Speed = PrimaryProgression(np.Progression, np.Speed, np.Position, np.Archetype, SnapsPerGame, "Speed", np.IsPracticeSquad)
		FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, np.Position, np.Archetype, SnapsPerGame, "Football IQ", np.IsPracticeSquad)
		Strength = PrimaryProgression(np.Progression, np.Strength, np.Position, np.Archetype, SnapsPerGame, "Strength", np.IsPracticeSquad)
		Catching = PrimaryProgression(np.Progression, np.Catching, np.Position, np.Archetype, SnapsPerGame, "Catching", np.IsPracticeSquad)
		PassBlock = PrimaryProgression(np.Progression, np.PassBlock, np.Position, np.Archetype, SnapsPerGame, "Pass Blocking", np.IsPracticeSquad)
		RunBlock = PrimaryProgression(np.Progression, np.RunBlock, np.Position, np.Archetype, SnapsPerGame, "Run Blocking", np.IsPracticeSquad)

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
		Agility = PrimaryProgression(np.Progression, np.Agility, np.Position, np.Archetype, SnapsPerGame, "Agility", np.IsPracticeSquad)
		Carrying = PrimaryProgression(np.Progression, np.Carrying, np.Position, np.Archetype, SnapsPerGame, "Carrying", np.IsPracticeSquad)
		Speed = PrimaryProgression(np.Progression, np.Speed, np.Position, np.Archetype, SnapsPerGame, "Speed", np.IsPracticeSquad)
		FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, np.Position, np.Archetype, SnapsPerGame, "Football IQ", np.IsPracticeSquad)
		Strength = PrimaryProgression(np.Progression, np.Strength, np.Position, np.Archetype, SnapsPerGame, "Strength", np.IsPracticeSquad)
		Catching = PrimaryProgression(np.Progression, np.Catching, np.Position, np.Archetype, SnapsPerGame, "Catching", np.IsPracticeSquad)
		RouteRunning = PrimaryProgression(np.Progression, np.RouteRunning, np.Position, np.Archetype, SnapsPerGame, "Route Running", np.IsPracticeSquad)
		PassBlock = PrimaryProgression(np.Progression, np.PassBlock, np.Position, np.Archetype, SnapsPerGame, "Pass Blocking", np.IsPracticeSquad)
		RunBlock = PrimaryProgression(np.Progression, np.RunBlock, np.Position, np.Archetype, SnapsPerGame, "Run Blocking", np.IsPracticeSquad)

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
		Agility = PrimaryProgression(np.Progression, np.Agility, np.Position, np.Archetype, SnapsPerGame, "Agility", np.IsPracticeSquad)
		Carrying = PrimaryProgression(np.Progression, np.Carrying, np.Position, np.Archetype, SnapsPerGame, "Carrying", np.IsPracticeSquad)
		Speed = PrimaryProgression(np.Progression, np.Speed, np.Position, np.Archetype, SnapsPerGame, "Speed", np.IsPracticeSquad)
		FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, np.Position, np.Archetype, SnapsPerGame, "Football IQ", np.IsPracticeSquad)
		Strength = PrimaryProgression(np.Progression, np.Strength, np.Position, np.Archetype, SnapsPerGame, "Strength", np.IsPracticeSquad)
		Catching = PrimaryProgression(np.Progression, np.Catching, np.Position, np.Archetype, SnapsPerGame, "Catching", np.IsPracticeSquad)
		RouteRunning = PrimaryProgression(np.Progression, np.RouteRunning, np.Position, np.Archetype, SnapsPerGame, "Route Running", np.IsPracticeSquad)

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
		Agility = PrimaryProgression(np.Progression, np.Agility, np.Position, np.Archetype, SnapsPerGame, "Agility", np.IsPracticeSquad)
		FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, np.Position, np.Archetype, SnapsPerGame, "Football IQ", np.IsPracticeSquad)
		Strength = PrimaryProgression(np.Progression, np.Strength, np.Position, np.Archetype, SnapsPerGame, "Strength", np.IsPracticeSquad)
		PassBlock = PrimaryProgression(np.Progression, np.PassBlock, np.Position, np.Archetype, SnapsPerGame, "Pass Blocking", np.IsPracticeSquad)
		RunBlock = PrimaryProgression(np.Progression, np.RunBlock, np.Position, np.Archetype, SnapsPerGame, "Run Blocking", np.IsPracticeSquad)

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
		Agility = PrimaryProgression(np.Progression, np.Agility, np.Position, np.Archetype, SnapsPerGame, "Agility", np.IsPracticeSquad)
		FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, np.Position, np.Archetype, SnapsPerGame, "Football IQ", np.IsPracticeSquad)
		Strength = PrimaryProgression(np.Progression, np.Strength, np.Position, np.Archetype, SnapsPerGame, "Strength", np.IsPracticeSquad)
		PassRush = PrimaryProgression(np.Progression, np.PassRush, np.Position, np.Archetype, SnapsPerGame, "Pass Rush", np.IsPracticeSquad)
		RunDefense = PrimaryProgression(np.Progression, np.RunDefense, np.Position, np.Archetype, SnapsPerGame, "Run Defense", np.IsPracticeSquad)
		Tackle = PrimaryProgression(np.Progression, np.Tackle, np.Position, np.Archetype, SnapsPerGame, "Tackle", np.IsPracticeSquad)

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
		Agility = PrimaryProgression(np.Progression, np.Agility, np.Position, np.Archetype, SnapsPerGame, "Agility", np.IsPracticeSquad)
		FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, np.Position, np.Archetype, SnapsPerGame, "Football IQ", np.IsPracticeSquad)
		Strength = PrimaryProgression(np.Progression, np.Strength, np.Position, np.Archetype, SnapsPerGame, "Strength", np.IsPracticeSquad)
		Speed = PrimaryProgression(np.Progression, np.Speed, np.Position, np.Archetype, SnapsPerGame, "Speed", np.IsPracticeSquad)
		PassRush = PrimaryProgression(np.Progression, np.PassRush, np.Position, np.Archetype, SnapsPerGame, "Pass Rush", np.IsPracticeSquad)
		RunDefense = PrimaryProgression(np.Progression, np.RunDefense, np.Position, np.Archetype, SnapsPerGame, "Run Defense", np.IsPracticeSquad)
		Tackle = PrimaryProgression(np.Progression, np.Tackle, np.Position, np.Archetype, SnapsPerGame, "Tackle", np.IsPracticeSquad)

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
		Agility = PrimaryProgression(np.Progression, np.Agility, np.Position, np.Archetype, SnapsPerGame, "Agility", np.IsPracticeSquad)
		FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, np.Position, np.Archetype, SnapsPerGame, "Football IQ", np.IsPracticeSquad)
		Strength = PrimaryProgression(np.Progression, np.Strength, np.Position, np.Archetype, SnapsPerGame, "Strength", np.IsPracticeSquad)
		Speed = PrimaryProgression(np.Progression, np.Speed, np.Position, np.Archetype, SnapsPerGame, "Speed", np.IsPracticeSquad)
		PassRush = PrimaryProgression(np.Progression, np.PassRush, np.Position, np.Archetype, SnapsPerGame, "Pass Rush", np.IsPracticeSquad)
		RunDefense = PrimaryProgression(np.Progression, np.RunDefense, np.Position, np.Archetype, SnapsPerGame, "Run Defense", np.IsPracticeSquad)
		Tackle = PrimaryProgression(np.Progression, np.Tackle, np.Position, np.Archetype, SnapsPerGame, "Tackle", np.IsPracticeSquad)
		ManCoverage = PrimaryProgression(np.Progression, np.ManCoverage, np.Position, np.Archetype, SnapsPerGame, "Man Coverage", np.IsPracticeSquad)
		ZoneCoverage = PrimaryProgression(np.Progression, np.ZoneCoverage, np.Position, np.Archetype, SnapsPerGame, "Zone Coverage", np.IsPracticeSquad)

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
		Agility = PrimaryProgression(np.Progression, np.Agility, np.Position, np.Archetype, SnapsPerGame, "Agility", np.IsPracticeSquad)
		FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, np.Position, np.Archetype, SnapsPerGame, "Football IQ", np.IsPracticeSquad)
		Strength = PrimaryProgression(np.Progression, np.Strength, np.Position, np.Archetype, SnapsPerGame, "Strength", np.IsPracticeSquad)
		Speed = PrimaryProgression(np.Progression, np.Speed, np.Position, np.Archetype, SnapsPerGame, "Speed", np.IsPracticeSquad)
		Tackle = PrimaryProgression(np.Progression, np.Tackle, np.Position, np.Archetype, SnapsPerGame, "Tackle", np.IsPracticeSquad)
		ManCoverage = PrimaryProgression(np.Progression, np.ManCoverage, np.Position, np.Archetype, SnapsPerGame, "Man Coverage", np.IsPracticeSquad)
		ZoneCoverage = PrimaryProgression(np.Progression, np.ZoneCoverage, np.Position, np.Archetype, SnapsPerGame, "Zone Coverage", np.IsPracticeSquad)
		Catching = PrimaryProgression(np.Progression, np.Catching, np.Position, np.Archetype, SnapsPerGame, "Catching", np.IsPracticeSquad)

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
		Agility = PrimaryProgression(np.Progression, np.Agility, np.Position, np.Archetype, SnapsPerGame, "Agility", np.IsPracticeSquad)
		FootballIQ = PrimaryProgression(np.Progression, np.FootballIQ, np.Position, np.Archetype, SnapsPerGame, "Football IQ", np.IsPracticeSquad)
		Strength = PrimaryProgression(np.Progression, np.Strength, np.Position, np.Archetype, SnapsPerGame, "Strength", np.IsPracticeSquad)
		Speed = PrimaryProgression(np.Progression, np.Speed, np.Position, np.Archetype, SnapsPerGame, "Speed", np.IsPracticeSquad)
		RunDefense = PrimaryProgression(np.Progression, np.RunDefense, np.Position, np.Archetype, SnapsPerGame, "Run Defense", np.IsPracticeSquad)
		Tackle = PrimaryProgression(np.Progression, np.Tackle, np.Position, np.Archetype, SnapsPerGame, "Tackle", np.IsPracticeSquad)
		ManCoverage = PrimaryProgression(np.Progression, np.ManCoverage, np.Position, np.Archetype, SnapsPerGame, "Man Coverage", np.IsPracticeSquad)
		ZoneCoverage = PrimaryProgression(np.Progression, np.ZoneCoverage, np.Position, np.Archetype, SnapsPerGame, "Zone Coverage", np.IsPracticeSquad)
		Catching = PrimaryProgression(np.Progression, np.Catching, np.Position, np.Archetype, SnapsPerGame, "Catching", np.IsPracticeSquad)

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
		KickPower = PrimaryProgression(np.Progression, np.KickPower, np.Position, np.Archetype, SnapsPerGame, "Kick Power", np.IsPracticeSquad)
		KickAccuracy = PrimaryProgression(np.Progression, np.KickAccuracy, np.Position, np.Archetype, SnapsPerGame, "Kick Accuracy", np.IsPracticeSquad)
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
			ThrowPower = PrimaryProgression(np.Progression, np.ThrowPower, np.Position, np.Archetype, SnapsPerGame, "Throw Power", np.IsPracticeSquad)
			ThrowAccuracy = PrimaryProgression(np.Progression, np.ThrowAccuracy, np.Position, np.Archetype, SnapsPerGame, "Throw Accuracy", np.IsPracticeSquad)
		} else {
			ThrowPower = SecondaryProgression(np.Progression, np.ThrowPower)
			ThrowAccuracy = SecondaryProgression(np.Progression, np.ThrowAccuracy)
		}
		// Primary Progressions
		PuntPower = PrimaryProgression(np.Progression, np.PuntPower, np.Position, np.Archetype, SnapsPerGame, "Punt Power", np.IsPracticeSquad)
		PuntAccuracy = PrimaryProgression(np.Progression, np.PuntAccuracy, np.Position, np.Archetype, SnapsPerGame, "Punt Accuracy", np.IsPracticeSquad)
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

	ThrowPower = RegressAttribute(ThrowPower, np.Age, false)
	ThrowAccuracy = RegressAttribute(ThrowAccuracy, np.Age, false)
	RunBlock = RegressAttribute(RunBlock, np.Age, false)
	PassBlock = RegressAttribute(PassBlock, np.Age, false)
	RunDefense = RegressAttribute(RunDefense, np.Age, false)
	PassRush = RegressAttribute(PassRush, np.Age, false)
	Carrying = RegressAttribute(Carrying, np.Age, false)
	Tackle = RegressAttribute(Tackle, np.Age, false)
	RouteRunning = RegressAttribute(RouteRunning, np.Age, false)
	Catching = RegressAttribute(Catching, np.Age, false)
	KickPower = RegressAttribute(KickPower, np.Age, false)
	KickAccuracy = RegressAttribute(KickAccuracy, np.Age, false)
	ManCoverage = RegressAttribute(ManCoverage, np.Age, false)
	ZoneCoverage = RegressAttribute(ZoneCoverage, np.Age, false)
	Strength = RegressAttribute(Strength, np.Age, false)
	Speed = RegressAttribute(Speed, np.Age, false)
	Agility = RegressAttribute(Agility, np.Age, false)
	FootballIQ = RegressAttribute(FootballIQ, np.Age, true)
	PuntPower = RegressAttribute(PuntPower, np.Age, false)
	PuntAccuracy = RegressAttribute(PuntAccuracy, np.Age, false)
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

func ProgressCollegePlayer(cp structs.CollegePlayer, SeasonID string, stats []structs.CollegePlayerStats) structs.CollegePlayer {
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

	if cp.Position == "QB" {
		// Primary Progressions
		Agility = PrimaryProgression(cp.Progression, cp.Agility, cp.Position, cp.Archetype, SnapsPerGame, "Agility", cp.IsRedshirting)
		ThrowPower = PrimaryProgression(cp.Progression, cp.ThrowPower, cp.Position, cp.Archetype, SnapsPerGame, "Throw Power", cp.IsRedshirting)
		ThrowAccuracy = PrimaryProgression(cp.Progression, cp.ThrowAccuracy, cp.Position, cp.Archetype, SnapsPerGame, "Throw Accuracy", cp.IsRedshirting)
		Speed = PrimaryProgression(cp.Progression, cp.Speed, cp.Position, cp.Archetype, SnapsPerGame, "Speed", cp.IsRedshirting)
		FootballIQ = PrimaryProgression(cp.Progression, cp.FootballIQ, cp.Position, cp.Archetype, SnapsPerGame, "Football IQ", cp.IsRedshirting)
		Strength = PrimaryProgression(cp.Progression, cp.Strength, cp.Position, cp.Archetype, SnapsPerGame, "Strength", cp.IsRedshirting)

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
	} else if cp.Position == "RB" {
		// Primary
		Agility = PrimaryProgression(cp.Progression, cp.Agility, cp.Position, cp.Archetype, SnapsPerGame, "Agility", cp.IsRedshirting)
		Carrying = PrimaryProgression(cp.Progression, cp.Carrying, cp.Position, cp.Archetype, SnapsPerGame, "Carrying", cp.IsRedshirting)
		Speed = PrimaryProgression(cp.Progression, cp.Speed, cp.Position, cp.Archetype, SnapsPerGame, "Speed", cp.IsRedshirting)
		FootballIQ = PrimaryProgression(cp.Progression, cp.FootballIQ, cp.Position, cp.Archetype, SnapsPerGame, "Football IQ", cp.IsRedshirting)
		Strength = PrimaryProgression(cp.Progression, cp.Strength, cp.Position, cp.Archetype, SnapsPerGame, "Strength", cp.IsRedshirting)
		Catching = PrimaryProgression(cp.Progression, cp.Catching, cp.Position, cp.Archetype, SnapsPerGame, "Catching", cp.IsRedshirting)
		PassBlock = PrimaryProgression(cp.Progression, cp.PassBlock, cp.Position, cp.Archetype, SnapsPerGame, "Pass Blocking", cp.IsRedshirting)
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
	} else if cp.Position == "FB" {
		// Primary
		Agility = PrimaryProgression(cp.Progression, cp.Agility, cp.Position, cp.Archetype, SnapsPerGame, "Agility", cp.IsRedshirting)
		Carrying = PrimaryProgression(cp.Progression, cp.Carrying, cp.Position, cp.Archetype, SnapsPerGame, "Carrying", cp.IsRedshirting)
		Speed = PrimaryProgression(cp.Progression, cp.Speed, cp.Position, cp.Archetype, SnapsPerGame, "Speed", cp.IsRedshirting)
		FootballIQ = PrimaryProgression(cp.Progression, cp.FootballIQ, cp.Position, cp.Archetype, SnapsPerGame, "Football IQ", cp.IsRedshirting)
		Strength = PrimaryProgression(cp.Progression, cp.Strength, cp.Position, cp.Archetype, SnapsPerGame, "Strength", cp.IsRedshirting)
		Catching = PrimaryProgression(cp.Progression, cp.Catching, cp.Position, cp.Archetype, SnapsPerGame, "Catching", cp.IsRedshirting)
		PassBlock = PrimaryProgression(cp.Progression, cp.PassBlock, cp.Position, cp.Archetype, SnapsPerGame, "Pass Blocking", cp.IsRedshirting)
		RunBlock = PrimaryProgression(cp.Progression, cp.RunBlock, cp.Position, cp.Archetype, SnapsPerGame, "Run Blocking", cp.IsRedshirting)

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

	} else if cp.Position == "TE" {
		// Primary
		Agility = PrimaryProgression(cp.Progression, cp.Agility, cp.Position, cp.Archetype, SnapsPerGame, "Agility", cp.IsRedshirting)
		Carrying = PrimaryProgression(cp.Progression, cp.Carrying, cp.Position, cp.Archetype, SnapsPerGame, "Carrying", cp.IsRedshirting)
		Speed = PrimaryProgression(cp.Progression, cp.Speed, cp.Position, cp.Archetype, SnapsPerGame, "Speed", cp.IsRedshirting)
		FootballIQ = PrimaryProgression(cp.Progression, cp.FootballIQ, cp.Position, cp.Archetype, SnapsPerGame, "Football IQ", cp.IsRedshirting)
		Strength = PrimaryProgression(cp.Progression, cp.Strength, cp.Position, cp.Archetype, SnapsPerGame, "Strength", cp.IsRedshirting)
		Catching = PrimaryProgression(cp.Progression, cp.Catching, cp.Position, cp.Archetype, SnapsPerGame, "Catching", cp.IsRedshirting)
		RouteRunning = PrimaryProgression(cp.Progression, cp.RouteRunning, cp.Position, cp.Archetype, SnapsPerGame, "Route Running", cp.IsRedshirting)
		PassBlock = PrimaryProgression(cp.Progression, cp.PassBlock, cp.Position, cp.Archetype, SnapsPerGame, "Pass Blocking", cp.IsRedshirting)
		RunBlock = PrimaryProgression(cp.Progression, cp.RunBlock, cp.Position, cp.Archetype, SnapsPerGame, "Run Blocking", cp.IsRedshirting)

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
	} else if cp.Position == "WR" {
		// Primary
		Agility = PrimaryProgression(cp.Progression, cp.Agility, cp.Position, cp.Archetype, SnapsPerGame, "Agility", cp.IsRedshirting)
		Carrying = PrimaryProgression(cp.Progression, cp.Carrying, cp.Position, cp.Archetype, SnapsPerGame, "Carrying", cp.IsRedshirting)
		Speed = PrimaryProgression(cp.Progression, cp.Speed, cp.Position, cp.Archetype, SnapsPerGame, "Speed", cp.IsRedshirting)
		FootballIQ = PrimaryProgression(cp.Progression, cp.FootballIQ, cp.Position, cp.Archetype, SnapsPerGame, "Football IQ", cp.IsRedshirting)
		Strength = PrimaryProgression(cp.Progression, cp.Strength, cp.Position, cp.Archetype, SnapsPerGame, "Strength", cp.IsRedshirting)
		Catching = PrimaryProgression(cp.Progression, cp.Catching, cp.Position, cp.Archetype, SnapsPerGame, "Catching", cp.IsRedshirting)
		RouteRunning = PrimaryProgression(cp.Progression, cp.RouteRunning, cp.Position, cp.Archetype, SnapsPerGame, "Route Running", cp.IsRedshirting)

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
	} else if cp.Position == "OT" || cp.Position == "OG" || cp.Position == "C" {
		// Primary
		Agility = PrimaryProgression(cp.Progression, cp.Agility, cp.Position, cp.Archetype, SnapsPerGame, "Agility", cp.IsRedshirting)
		FootballIQ = PrimaryProgression(cp.Progression, cp.FootballIQ, cp.Position, cp.Archetype, SnapsPerGame, "Football IQ", cp.IsRedshirting)
		Strength = PrimaryProgression(cp.Progression, cp.Strength, cp.Position, cp.Archetype, SnapsPerGame, "Strength", cp.IsRedshirting)
		PassBlock = PrimaryProgression(cp.Progression, cp.PassBlock, cp.Position, cp.Archetype, SnapsPerGame, "Pass Blocking", cp.IsRedshirting)
		RunBlock = PrimaryProgression(cp.Progression, cp.RunBlock, cp.Position, cp.Archetype, SnapsPerGame, "Run Blocking", cp.IsRedshirting)

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
	} else if cp.Position == "DT" {
		// Primary
		Agility = PrimaryProgression(cp.Progression, cp.Agility, cp.Position, cp.Archetype, SnapsPerGame, "Agility", cp.IsRedshirting)
		FootballIQ = PrimaryProgression(cp.Progression, cp.FootballIQ, cp.Position, cp.Archetype, SnapsPerGame, "Football IQ", cp.IsRedshirting)
		Strength = PrimaryProgression(cp.Progression, cp.Strength, cp.Position, cp.Archetype, SnapsPerGame, "Strength", cp.IsRedshirting)
		PassRush = PrimaryProgression(cp.Progression, cp.PassRush, cp.Position, cp.Archetype, SnapsPerGame, "Pass Rush", cp.IsRedshirting)
		RunDefense = PrimaryProgression(cp.Progression, cp.RunDefense, cp.Position, cp.Archetype, SnapsPerGame, "Run Defense", cp.IsRedshirting)
		Tackle = PrimaryProgression(cp.Progression, cp.Tackle, cp.Position, cp.Archetype, SnapsPerGame, "Tackle", cp.IsRedshirting)

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
	} else if cp.Position == "DE" {
		// Primary
		Agility = PrimaryProgression(cp.Progression, cp.Agility, cp.Position, cp.Archetype, SnapsPerGame, "Agility", cp.IsRedshirting)
		FootballIQ = PrimaryProgression(cp.Progression, cp.FootballIQ, cp.Position, cp.Archetype, SnapsPerGame, "Football IQ", cp.IsRedshirting)
		Strength = PrimaryProgression(cp.Progression, cp.Strength, cp.Position, cp.Archetype, SnapsPerGame, "Strength", cp.IsRedshirting)
		Speed = PrimaryProgression(cp.Progression, cp.Speed, cp.Position, cp.Archetype, SnapsPerGame, "Speed", cp.IsRedshirting)
		PassRush = PrimaryProgression(cp.Progression, cp.PassRush, cp.Position, cp.Archetype, SnapsPerGame, "Pass Rush", cp.IsRedshirting)
		RunDefense = PrimaryProgression(cp.Progression, cp.RunDefense, cp.Position, cp.Archetype, SnapsPerGame, "Run Defense", cp.IsRedshirting)
		Tackle = PrimaryProgression(cp.Progression, cp.Tackle, cp.Position, cp.Archetype, SnapsPerGame, "Tackle", cp.IsRedshirting)

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
	} else if cp.Position == "OLB" || cp.Position == "ILB" {
		// Primary
		Agility = PrimaryProgression(cp.Progression, cp.Agility, cp.Position, cp.Archetype, SnapsPerGame, "Agility", cp.IsRedshirting)
		FootballIQ = PrimaryProgression(cp.Progression, cp.FootballIQ, cp.Position, cp.Archetype, SnapsPerGame, "Football IQ", cp.IsRedshirting)
		Strength = PrimaryProgression(cp.Progression, cp.Strength, cp.Position, cp.Archetype, SnapsPerGame, "Strength", cp.IsRedshirting)
		Speed = PrimaryProgression(cp.Progression, cp.Speed, cp.Position, cp.Archetype, SnapsPerGame, "Speed", cp.IsRedshirting)
		PassRush = PrimaryProgression(cp.Progression, cp.PassRush, cp.Position, cp.Archetype, SnapsPerGame, "Pass Rush", cp.IsRedshirting)
		RunDefense = PrimaryProgression(cp.Progression, cp.RunDefense, cp.Position, cp.Archetype, SnapsPerGame, "Run Defense", cp.IsRedshirting)
		Tackle = PrimaryProgression(cp.Progression, cp.Tackle, cp.Position, cp.Archetype, SnapsPerGame, "Tackle", cp.IsRedshirting)
		ManCoverage = PrimaryProgression(cp.Progression, cp.ManCoverage, cp.Position, cp.Archetype, SnapsPerGame, "Man Coverage", cp.IsRedshirting)
		ZoneCoverage = PrimaryProgression(cp.Progression, cp.ZoneCoverage, cp.Position, cp.Archetype, SnapsPerGame, "Zone Coverage", cp.IsRedshirting)

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
	} else if cp.Position == "CB" {
		// Primary
		Agility = PrimaryProgression(cp.Progression, cp.Agility, cp.Position, cp.Archetype, SnapsPerGame, "Agility", cp.IsRedshirting)
		FootballIQ = PrimaryProgression(cp.Progression, cp.FootballIQ, cp.Position, cp.Archetype, SnapsPerGame, "Football IQ", cp.IsRedshirting)
		Strength = PrimaryProgression(cp.Progression, cp.Strength, cp.Position, cp.Archetype, SnapsPerGame, "Strength", cp.IsRedshirting)
		Speed = PrimaryProgression(cp.Progression, cp.Speed, cp.Position, cp.Archetype, SnapsPerGame, "Speed", cp.IsRedshirting)
		Tackle = PrimaryProgression(cp.Progression, cp.Tackle, cp.Position, cp.Archetype, SnapsPerGame, "Tackle", cp.IsRedshirting)
		ManCoverage = PrimaryProgression(cp.Progression, cp.ManCoverage, cp.Position, cp.Archetype, SnapsPerGame, "Man Coverage", cp.IsRedshirting)
		ZoneCoverage = PrimaryProgression(cp.Progression, cp.ZoneCoverage, cp.Position, cp.Archetype, SnapsPerGame, "Zone Coverage", cp.IsRedshirting)
		Catching = PrimaryProgression(cp.Progression, cp.Catching, cp.Position, cp.Archetype, SnapsPerGame, "Catching", cp.IsRedshirting)

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
	} else if cp.Position == "FS" || cp.Position == "SS" {
		// Primary
		Agility = PrimaryProgression(cp.Progression, cp.Agility, cp.Position, cp.Archetype, SnapsPerGame, "Agility", cp.IsRedshirting)
		FootballIQ = PrimaryProgression(cp.Progression, cp.FootballIQ, cp.Position, cp.Archetype, SnapsPerGame, "Football IQ", cp.IsRedshirting)
		Strength = PrimaryProgression(cp.Progression, cp.Strength, cp.Position, cp.Archetype, SnapsPerGame, "Strength", cp.IsRedshirting)
		Speed = PrimaryProgression(cp.Progression, cp.Speed, cp.Position, cp.Archetype, SnapsPerGame, "Speed", cp.IsRedshirting)
		RunDefense = PrimaryProgression(cp.Progression, cp.RunDefense, cp.Position, cp.Archetype, SnapsPerGame, "Run Defense", cp.IsRedshirting)
		Tackle = PrimaryProgression(cp.Progression, cp.Tackle, cp.Position, cp.Archetype, SnapsPerGame, "Tackle", cp.IsRedshirting)
		ManCoverage = PrimaryProgression(cp.Progression, cp.ManCoverage, cp.Position, cp.Archetype, SnapsPerGame, "Man Coverage", cp.IsRedshirting)
		ZoneCoverage = PrimaryProgression(cp.Progression, cp.ZoneCoverage, cp.Position, cp.Archetype, SnapsPerGame, "Zone Coverage", cp.IsRedshirting)
		Catching = PrimaryProgression(cp.Progression, cp.Catching, cp.Position, cp.Archetype, SnapsPerGame, "Catching", cp.IsRedshirting)

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
	} else if cp.Position == "K" {
		// Primary
		KickPower = PrimaryProgression(cp.Progression, cp.KickPower, cp.Position, cp.Archetype, SnapsPerGame, "Kick Power", cp.IsRedshirting)
		KickAccuracy = PrimaryProgression(cp.Progression, cp.KickAccuracy, cp.Position, cp.Archetype, SnapsPerGame, "Kick Accuracy", cp.IsRedshirting)
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
	} else if cp.Position == "P" {
		// If David Ross
		if cp.ID == 24984 {
			ThrowPower = PrimaryProgression(cp.Progression, cp.ThrowPower, cp.Position, cp.Archetype, SnapsPerGame, "Throw Power", cp.IsRedshirting)
			ThrowAccuracy = PrimaryProgression(cp.Progression, cp.ThrowAccuracy, cp.Position, cp.Archetype, SnapsPerGame, "Throw Accuracy", cp.IsRedshirting)
		} else {
			ThrowPower = SecondaryProgression(cp.Progression, cp.ThrowPower)
			ThrowAccuracy = SecondaryProgression(cp.Progression, cp.ThrowAccuracy)
		}
		// Primary Progressions
		PuntPower = PrimaryProgression(cp.Progression, cp.PuntPower, cp.Position, cp.Archetype, SnapsPerGame, "Punt Power", cp.IsRedshirting)
		PuntAccuracy = PrimaryProgression(cp.Progression, cp.PuntAccuracy, cp.Position, cp.Archetype, SnapsPerGame, "Punt Accuracy", cp.IsRedshirting)
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
		FootballIQ = SecondaryProgression(cp.Progression, cp.FootballIQ)
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

	cp.Progress(progressions)

	if cp.IsRedshirting {
		cp.SetRedshirtStatus()
	}

	cp.GetOverall()

	return cp
}

func ProgressUnsignedPlayer(up structs.UnsignedPlayer, SeasonID string) structs.UnsignedPlayer {
	var SnapsPerGame int = 0
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

	if up.Position == "QB" {
		// Primary Progressions
		Agility = PrimaryProgression(up.Progression, up.Agility, up.Position, up.Archetype, SnapsPerGame, "Agility", true)
		ThrowPower = PrimaryProgression(up.Progression, up.ThrowPower, up.Position, up.Archetype, SnapsPerGame, "Throw Power", true)
		ThrowAccuracy = PrimaryProgression(up.Progression, up.ThrowAccuracy, up.Position, up.Archetype, SnapsPerGame, "Throw Accuracy", true)
		Speed = PrimaryProgression(up.Progression, up.Speed, up.Position, up.Archetype, SnapsPerGame, "Speed", true)
		FootballIQ = PrimaryProgression(up.Progression, up.FootballIQ, up.Position, up.Archetype, SnapsPerGame, "Football IQ", true)
		Strength = PrimaryProgression(up.Progression, up.Strength, up.Position, up.Archetype, SnapsPerGame, "Strength", true)

		// Secondary Progressions
		RunBlock = SecondaryProgression(up.Progression, up.RunBlock)
		PassBlock = SecondaryProgression(up.Progression, up.PassBlock)
		RunDefense = SecondaryProgression(up.Progression, up.RunDefense)
		PassRush = SecondaryProgression(up.Progression, up.PassRush)
		Carrying = SecondaryProgression(up.Progression, up.Carrying)
		Tackle = SecondaryProgression(up.Progression, up.Tackle)
		RouteRunning = SecondaryProgression(up.RouteRunning, up.RouteRunning)
		Catching = SecondaryProgression(up.Progression, up.Catching)
		PuntPower = SecondaryProgression(up.Progression, up.PuntPower)
		PuntAccuracy = SecondaryProgression(up.Progression, up.PuntAccuracy)
		KickPower = SecondaryProgression(up.Progression, up.KickPower)
		KickAccuracy = SecondaryProgression(up.Progression, up.KickAccuracy)
		ManCoverage = SecondaryProgression(up.Progression, up.ManCoverage)
		ZoneCoverage = SecondaryProgression(up.Progression, up.ZoneCoverage)
	} else if up.Position == "RB" {
		// Primary
		Agility = PrimaryProgression(up.Progression, up.Agility, up.Position, up.Archetype, SnapsPerGame, "Agility", true)
		Carrying = PrimaryProgression(up.Progression, up.Carrying, up.Position, up.Archetype, SnapsPerGame, "Carrying", true)
		Speed = PrimaryProgression(up.Progression, up.Speed, up.Position, up.Archetype, SnapsPerGame, "Speed", true)
		FootballIQ = PrimaryProgression(up.Progression, up.FootballIQ, up.Position, up.Archetype, SnapsPerGame, "Football IQ", true)
		Strength = PrimaryProgression(up.Progression, up.Strength, up.Position, up.Archetype, SnapsPerGame, "Strength", true)
		Catching = PrimaryProgression(up.Progression, up.Catching, up.Position, up.Archetype, SnapsPerGame, "Catching", true)
		PassBlock = PrimaryProgression(up.Progression, up.PassBlock, up.Position, up.Archetype, SnapsPerGame, "Pass Blocking", true)
		// Secondary
		RunBlock = SecondaryProgression(up.Progression, up.RunBlock)
		RunDefense = SecondaryProgression(up.Progression, up.RunDefense)
		PassRush = SecondaryProgression(up.Progression, up.PassRush)
		Tackle = SecondaryProgression(up.Progression, up.Tackle)
		RouteRunning = SecondaryProgression(up.RouteRunning, up.RouteRunning)
		PuntPower = SecondaryProgression(up.Progression, up.PuntPower)
		PuntAccuracy = SecondaryProgression(up.Progression, up.PuntAccuracy)
		KickPower = SecondaryProgression(up.Progression, up.KickPower)
		KickAccuracy = SecondaryProgression(up.Progression, up.KickAccuracy)
		ManCoverage = SecondaryProgression(up.Progression, up.ManCoverage)
		ZoneCoverage = SecondaryProgression(up.Progression, up.ZoneCoverage)
		ThrowPower = SecondaryProgression(up.Progression, up.ThrowPower)
		ThrowAccuracy = SecondaryProgression(up.Progression, up.ThrowAccuracy)
	} else if up.Position == "FB" {
		// Primary
		Agility = PrimaryProgression(up.Progression, up.Agility, up.Position, up.Archetype, SnapsPerGame, "Agility", true)
		Carrying = PrimaryProgression(up.Progression, up.Carrying, up.Position, up.Archetype, SnapsPerGame, "Carrying", true)
		Speed = PrimaryProgression(up.Progression, up.Speed, up.Position, up.Archetype, SnapsPerGame, "Speed", true)
		FootballIQ = PrimaryProgression(up.Progression, up.FootballIQ, up.Position, up.Archetype, SnapsPerGame, "Football IQ", true)
		Strength = PrimaryProgression(up.Progression, up.Strength, up.Position, up.Archetype, SnapsPerGame, "Strength", true)
		Catching = PrimaryProgression(up.Progression, up.Catching, up.Position, up.Archetype, SnapsPerGame, "Catching", true)
		PassBlock = PrimaryProgression(up.Progression, up.PassBlock, up.Position, up.Archetype, SnapsPerGame, "Pass Blocking", true)
		RunBlock = PrimaryProgression(up.Progression, up.RunBlock, up.Position, up.Archetype, SnapsPerGame, "Run Blocking", true)

		// Secondary
		RunDefense = SecondaryProgression(up.Progression, up.RunDefense)
		PassRush = SecondaryProgression(up.Progression, up.PassRush)
		Tackle = SecondaryProgression(up.Progression, up.Tackle)
		RouteRunning = SecondaryProgression(up.RouteRunning, up.RouteRunning)
		PuntPower = SecondaryProgression(up.Progression, up.PuntPower)
		PuntAccuracy = SecondaryProgression(up.Progression, up.PuntAccuracy)
		KickPower = SecondaryProgression(up.Progression, up.KickPower)
		KickAccuracy = SecondaryProgression(up.Progression, up.KickAccuracy)
		ManCoverage = SecondaryProgression(up.Progression, up.ManCoverage)
		ZoneCoverage = SecondaryProgression(up.Progression, up.ZoneCoverage)
		ThrowPower = SecondaryProgression(up.Progression, up.ThrowPower)
		ThrowAccuracy = SecondaryProgression(up.Progression, up.ThrowAccuracy)

	} else if up.Position == "TE" {
		// Primary
		Agility = PrimaryProgression(up.Progression, up.Agility, up.Position, up.Archetype, SnapsPerGame, "Agility", true)
		Carrying = PrimaryProgression(up.Progression, up.Carrying, up.Position, up.Archetype, SnapsPerGame, "Carrying", true)
		Speed = PrimaryProgression(up.Progression, up.Speed, up.Position, up.Archetype, SnapsPerGame, "Speed", true)
		FootballIQ = PrimaryProgression(up.Progression, up.FootballIQ, up.Position, up.Archetype, SnapsPerGame, "Football IQ", true)
		Strength = PrimaryProgression(up.Progression, up.Strength, up.Position, up.Archetype, SnapsPerGame, "Strength", true)
		Catching = PrimaryProgression(up.Progression, up.Catching, up.Position, up.Archetype, SnapsPerGame, "Catching", true)
		RouteRunning = PrimaryProgression(up.Progression, up.RouteRunning, up.Position, up.Archetype, SnapsPerGame, "Route Running", true)
		PassBlock = PrimaryProgression(up.Progression, up.PassBlock, up.Position, up.Archetype, SnapsPerGame, "Pass Blocking", true)
		RunBlock = PrimaryProgression(up.Progression, up.RunBlock, up.Position, up.Archetype, SnapsPerGame, "Run Blocking", true)

		// Secondary
		RunDefense = SecondaryProgression(up.Progression, up.RunDefense)
		PassRush = SecondaryProgression(up.Progression, up.PassRush)
		Tackle = SecondaryProgression(up.Progression, up.Tackle)
		PuntPower = SecondaryProgression(up.Progression, up.PuntPower)
		PuntAccuracy = SecondaryProgression(up.Progression, up.PuntAccuracy)
		KickPower = SecondaryProgression(up.Progression, up.KickPower)
		KickAccuracy = SecondaryProgression(up.Progression, up.KickAccuracy)
		ManCoverage = SecondaryProgression(up.Progression, up.ManCoverage)
		ZoneCoverage = SecondaryProgression(up.Progression, up.ZoneCoverage)
		ThrowPower = SecondaryProgression(up.Progression, up.ThrowPower)
		ThrowAccuracy = SecondaryProgression(up.Progression, up.ThrowAccuracy)
	} else if up.Position == "WR" {
		// Primary
		Agility = PrimaryProgression(up.Progression, up.Agility, up.Position, up.Archetype, SnapsPerGame, "Agility", true)
		Carrying = PrimaryProgression(up.Progression, up.Carrying, up.Position, up.Archetype, SnapsPerGame, "Carrying", true)
		Speed = PrimaryProgression(up.Progression, up.Speed, up.Position, up.Archetype, SnapsPerGame, "Speed", true)
		FootballIQ = PrimaryProgression(up.Progression, up.FootballIQ, up.Position, up.Archetype, SnapsPerGame, "Football IQ", true)
		Strength = PrimaryProgression(up.Progression, up.Strength, up.Position, up.Archetype, SnapsPerGame, "Strength", true)
		Catching = PrimaryProgression(up.Progression, up.Catching, up.Position, up.Archetype, SnapsPerGame, "Catching", true)
		RouteRunning = PrimaryProgression(up.Progression, up.RouteRunning, up.Position, up.Archetype, SnapsPerGame, "Route Running", true)

		// Secondary
		RunBlock = SecondaryProgression(up.Progression, up.RunBlock)
		PassBlock = SecondaryProgression(up.Progression, up.PassBlock)
		RunDefense = SecondaryProgression(up.Progression, up.RunDefense)
		PassRush = SecondaryProgression(up.Progression, up.PassRush)
		Tackle = SecondaryProgression(up.Progression, up.Tackle)
		PuntPower = SecondaryProgression(up.Progression, up.PuntPower)
		PuntAccuracy = SecondaryProgression(up.Progression, up.PuntAccuracy)
		KickPower = SecondaryProgression(up.Progression, up.KickPower)
		KickAccuracy = SecondaryProgression(up.Progression, up.KickAccuracy)
		ManCoverage = SecondaryProgression(up.Progression, up.ManCoverage)
		ZoneCoverage = SecondaryProgression(up.Progression, up.ZoneCoverage)
		ThrowPower = SecondaryProgression(up.Progression, up.ThrowPower)
		ThrowAccuracy = SecondaryProgression(up.Progression, up.ThrowAccuracy)
	} else if up.Position == "OT" || up.Position == "OG" || up.Position == "C" {
		// Primary
		Agility = PrimaryProgression(up.Progression, up.Agility, up.Position, up.Archetype, SnapsPerGame, "Agility", true)
		FootballIQ = PrimaryProgression(up.Progression, up.FootballIQ, up.Position, up.Archetype, SnapsPerGame, "Football IQ", true)
		Strength = PrimaryProgression(up.Progression, up.Strength, up.Position, up.Archetype, SnapsPerGame, "Strength", true)
		PassBlock = PrimaryProgression(up.Progression, up.PassBlock, up.Position, up.Archetype, SnapsPerGame, "Pass Blocking", true)
		RunBlock = PrimaryProgression(up.Progression, up.RunBlock, up.Position, up.Archetype, SnapsPerGame, "Run Blocking", true)

		// Secondary
		RunDefense = SecondaryProgression(up.Progression, up.RunDefense)
		PassRush = SecondaryProgression(up.Progression, up.PassRush)
		Tackle = SecondaryProgression(up.Progression, up.Tackle)
		RouteRunning = SecondaryProgression(up.RouteRunning, up.RouteRunning)
		PuntPower = SecondaryProgression(up.Progression, up.PuntPower)
		PuntAccuracy = SecondaryProgression(up.Progression, up.PuntAccuracy)
		KickPower = SecondaryProgression(up.Progression, up.KickPower)
		KickAccuracy = SecondaryProgression(up.Progression, up.KickAccuracy)
		ManCoverage = SecondaryProgression(up.Progression, up.ManCoverage)
		ZoneCoverage = SecondaryProgression(up.Progression, up.ZoneCoverage)
		ThrowPower = SecondaryProgression(up.Progression, up.ThrowPower)
		ThrowAccuracy = SecondaryProgression(up.Progression, up.ThrowAccuracy)
		Carrying = SecondaryProgression(up.Progression, up.Carrying)
		Speed = SecondaryProgression(up.Progression, up.Speed)
		Catching = SecondaryProgression(up.Progression, up.Catching)
	} else if up.Position == "DT" {
		// Primary
		Agility = PrimaryProgression(up.Progression, up.Agility, up.Position, up.Archetype, SnapsPerGame, "Agility", true)
		FootballIQ = PrimaryProgression(up.Progression, up.FootballIQ, up.Position, up.Archetype, SnapsPerGame, "Football IQ", true)
		Strength = PrimaryProgression(up.Progression, up.Strength, up.Position, up.Archetype, SnapsPerGame, "Strength", true)
		PassRush = PrimaryProgression(up.Progression, up.PassRush, up.Position, up.Archetype, SnapsPerGame, "Pass Rush", true)
		RunDefense = PrimaryProgression(up.Progression, up.RunDefense, up.Position, up.Archetype, SnapsPerGame, "Run Defense", true)
		Tackle = PrimaryProgression(up.Progression, up.Tackle, up.Position, up.Archetype, SnapsPerGame, "Tackle", true)

		// Secondary
		PuntPower = SecondaryProgression(up.Progression, up.PuntPower)
		PuntAccuracy = SecondaryProgression(up.Progression, up.PuntAccuracy)
		KickPower = SecondaryProgression(up.Progression, up.KickPower)
		KickAccuracy = SecondaryProgression(up.Progression, up.KickAccuracy)
		ManCoverage = SecondaryProgression(up.Progression, up.ManCoverage)
		ZoneCoverage = SecondaryProgression(up.Progression, up.ZoneCoverage)
		ThrowPower = SecondaryProgression(up.Progression, up.ThrowPower)
		ThrowAccuracy = SecondaryProgression(up.Progression, up.ThrowAccuracy)
		Carrying = SecondaryProgression(up.Progression, up.Carrying)
		Catching = SecondaryProgression(up.Progression, up.Catching)
		RouteRunning = SecondaryProgression(up.RouteRunning, up.RouteRunning)
		Speed = SecondaryProgression(up.Progression, up.Speed)
		RunBlock = SecondaryProgression(up.Progression, up.RunBlock)
		PassBlock = SecondaryProgression(up.Progression, up.PassBlock)
	} else if up.Position == "DE" {
		// Primary
		Agility = PrimaryProgression(up.Progression, up.Agility, up.Position, up.Archetype, SnapsPerGame, "Agility", true)
		FootballIQ = PrimaryProgression(up.Progression, up.FootballIQ, up.Position, up.Archetype, SnapsPerGame, "Football IQ", true)
		Strength = PrimaryProgression(up.Progression, up.Strength, up.Position, up.Archetype, SnapsPerGame, "Strength", true)
		Speed = PrimaryProgression(up.Progression, up.Speed, up.Position, up.Archetype, SnapsPerGame, "Speed", true)
		PassRush = PrimaryProgression(up.Progression, up.PassRush, up.Position, up.Archetype, SnapsPerGame, "Pass Rush", true)
		RunDefense = PrimaryProgression(up.Progression, up.RunDefense, up.Position, up.Archetype, SnapsPerGame, "Run Defense", true)
		Tackle = PrimaryProgression(up.Progression, up.Tackle, up.Position, up.Archetype, SnapsPerGame, "Tackle", true)

		// Secondary
		PuntPower = SecondaryProgression(up.Progression, up.PuntPower)
		PuntAccuracy = SecondaryProgression(up.Progression, up.PuntAccuracy)
		KickPower = SecondaryProgression(up.Progression, up.KickPower)
		KickAccuracy = SecondaryProgression(up.Progression, up.KickAccuracy)
		ManCoverage = SecondaryProgression(up.Progression, up.ManCoverage)
		ZoneCoverage = SecondaryProgression(up.Progression, up.ZoneCoverage)
		ThrowPower = SecondaryProgression(up.Progression, up.ThrowPower)
		ThrowAccuracy = SecondaryProgression(up.Progression, up.ThrowAccuracy)
		Carrying = SecondaryProgression(up.Progression, up.Carrying)
		Catching = SecondaryProgression(up.Progression, up.Catching)
		RouteRunning = SecondaryProgression(up.RouteRunning, up.RouteRunning)
		RunBlock = SecondaryProgression(up.Progression, up.RunBlock)
		PassBlock = SecondaryProgression(up.Progression, up.PassBlock)
	} else if up.Position == "OLB" || up.Position == "ILB" {
		// Primary
		Agility = PrimaryProgression(up.Progression, up.Agility, up.Position, up.Archetype, SnapsPerGame, "Agility", true)
		FootballIQ = PrimaryProgression(up.Progression, up.FootballIQ, up.Position, up.Archetype, SnapsPerGame, "Football IQ", true)
		Strength = PrimaryProgression(up.Progression, up.Strength, up.Position, up.Archetype, SnapsPerGame, "Strength", true)
		Speed = PrimaryProgression(up.Progression, up.Speed, up.Position, up.Archetype, SnapsPerGame, "Speed", true)
		PassRush = PrimaryProgression(up.Progression, up.PassRush, up.Position, up.Archetype, SnapsPerGame, "Pass Rush", true)
		RunDefense = PrimaryProgression(up.Progression, up.RunDefense, up.Position, up.Archetype, SnapsPerGame, "Run Defense", true)
		Tackle = PrimaryProgression(up.Progression, up.Tackle, up.Position, up.Archetype, SnapsPerGame, "Tackle", true)
		ManCoverage = PrimaryProgression(up.Progression, up.ManCoverage, up.Position, up.Archetype, SnapsPerGame, "Man Coverage", true)
		ZoneCoverage = PrimaryProgression(up.Progression, up.ZoneCoverage, up.Position, up.Archetype, SnapsPerGame, "Zone Coverage", true)

		// Secondary
		PuntPower = SecondaryProgression(up.Progression, up.PuntPower)
		PuntAccuracy = SecondaryProgression(up.Progression, up.PuntAccuracy)
		KickPower = SecondaryProgression(up.Progression, up.KickPower)
		KickAccuracy = SecondaryProgression(up.Progression, up.KickAccuracy)
		ThrowPower = SecondaryProgression(up.Progression, up.ThrowPower)
		ThrowAccuracy = SecondaryProgression(up.Progression, up.ThrowAccuracy)
		Carrying = SecondaryProgression(up.Progression, up.Carrying)
		Catching = SecondaryProgression(up.Progression, up.Catching)
		RouteRunning = SecondaryProgression(up.RouteRunning, up.RouteRunning)
		RunBlock = SecondaryProgression(up.Progression, up.RunBlock)
		PassBlock = SecondaryProgression(up.Progression, up.PassBlock)
	} else if up.Position == "CB" {
		// Primary
		Agility = PrimaryProgression(up.Progression, up.Agility, up.Position, up.Archetype, SnapsPerGame, "Agility", true)
		FootballIQ = PrimaryProgression(up.Progression, up.FootballIQ, up.Position, up.Archetype, SnapsPerGame, "Football IQ", true)
		Strength = PrimaryProgression(up.Progression, up.Strength, up.Position, up.Archetype, SnapsPerGame, "Strength", true)
		Speed = PrimaryProgression(up.Progression, up.Speed, up.Position, up.Archetype, SnapsPerGame, "Speed", true)
		Tackle = PrimaryProgression(up.Progression, up.Tackle, up.Position, up.Archetype, SnapsPerGame, "Tackle", true)
		ManCoverage = PrimaryProgression(up.Progression, up.ManCoverage, up.Position, up.Archetype, SnapsPerGame, "Man Coverage", true)
		ZoneCoverage = PrimaryProgression(up.Progression, up.ZoneCoverage, up.Position, up.Archetype, SnapsPerGame, "Zone Coverage", true)
		Catching = PrimaryProgression(up.Progression, up.Catching, up.Position, up.Archetype, SnapsPerGame, "Catching", true)

		// Secondary
		ThrowPower = SecondaryProgression(up.Progression, up.ThrowPower)
		ThrowAccuracy = SecondaryProgression(up.Progression, up.ThrowAccuracy)
		RunDefense = SecondaryProgression(up.Progression, up.RunDefense)
		PassRush = SecondaryProgression(up.Progression, up.PassRush)
		PuntPower = SecondaryProgression(up.Progression, up.PuntPower)
		PuntAccuracy = SecondaryProgression(up.Progression, up.PuntAccuracy)
		KickPower = SecondaryProgression(up.Progression, up.KickPower)
		KickAccuracy = SecondaryProgression(up.Progression, up.KickAccuracy)
		Carrying = SecondaryProgression(up.Progression, up.Carrying)
		RunBlock = SecondaryProgression(up.Progression, up.RunBlock)
		PassBlock = SecondaryProgression(up.Progression, up.PassBlock)
		RouteRunning = SecondaryProgression(up.RouteRunning, up.RouteRunning)
	} else if up.Position == "FS" || up.Position == "SS" {
		// Primary
		Agility = PrimaryProgression(up.Progression, up.Agility, up.Position, up.Archetype, SnapsPerGame, "Agility", true)
		FootballIQ = PrimaryProgression(up.Progression, up.FootballIQ, up.Position, up.Archetype, SnapsPerGame, "Football IQ", true)
		Strength = PrimaryProgression(up.Progression, up.Strength, up.Position, up.Archetype, SnapsPerGame, "Strength", true)
		Speed = PrimaryProgression(up.Progression, up.Speed, up.Position, up.Archetype, SnapsPerGame, "Speed", true)
		RunDefense = PrimaryProgression(up.Progression, up.RunDefense, up.Position, up.Archetype, SnapsPerGame, "Run Defense", true)
		Tackle = PrimaryProgression(up.Progression, up.Tackle, up.Position, up.Archetype, SnapsPerGame, "Tackle", true)
		ManCoverage = PrimaryProgression(up.Progression, up.ManCoverage, up.Position, up.Archetype, SnapsPerGame, "Man Coverage", true)
		ZoneCoverage = PrimaryProgression(up.Progression, up.ZoneCoverage, up.Position, up.Archetype, SnapsPerGame, "Zone Coverage", true)
		Catching = PrimaryProgression(up.Progression, up.Catching, up.Position, up.Archetype, SnapsPerGame, "Catching", true)

		// Secondary
		ThrowPower = SecondaryProgression(up.Progression, up.ThrowPower)
		ThrowAccuracy = SecondaryProgression(up.Progression, up.ThrowAccuracy)
		PassRush = SecondaryProgression(up.Progression, up.PassRush)
		PuntPower = SecondaryProgression(up.Progression, up.PuntPower)
		PuntAccuracy = SecondaryProgression(up.Progression, up.PuntAccuracy)
		KickPower = SecondaryProgression(up.Progression, up.KickPower)
		KickAccuracy = SecondaryProgression(up.Progression, up.KickAccuracy)
		Carrying = SecondaryProgression(up.Progression, up.Carrying)
		RunBlock = SecondaryProgression(up.Progression, up.RunBlock)
		PassBlock = SecondaryProgression(up.Progression, up.PassBlock)
		RouteRunning = SecondaryProgression(up.RouteRunning, up.RouteRunning)
	} else if up.Position == "K" {
		// Primary
		KickPower = PrimaryProgression(up.Progression, up.KickPower, up.Position, up.Archetype, SnapsPerGame, "Kick Power", true)
		KickAccuracy = PrimaryProgression(up.Progression, up.KickAccuracy, up.Position, up.Archetype, SnapsPerGame, "Kick Accuracy", true)
		// Secondary
		ThrowPower = SecondaryProgression(up.Progression, up.ThrowPower)
		ThrowAccuracy = SecondaryProgression(up.Progression, up.ThrowAccuracy)
		PassRush = SecondaryProgression(up.Progression, up.PassRush)
		RunDefense = SecondaryProgression(up.Progression, up.RunDefense)
		PuntPower = SecondaryProgression(up.Progression, up.PuntPower)
		PuntAccuracy = SecondaryProgression(up.Progression, up.PuntAccuracy)
		Carrying = SecondaryProgression(up.Progression, up.Carrying)
		Catching = SecondaryProgression(up.Progression, up.Catching)
		RunBlock = SecondaryProgression(up.Progression, up.RunBlock)
		PassBlock = SecondaryProgression(up.Progression, up.PassBlock)
		RouteRunning = SecondaryProgression(up.RouteRunning, up.RouteRunning)
		Strength = SecondaryProgression(up.Progression, up.Strength)
		Speed = SecondaryProgression(up.Progression, up.Speed)
		Agility = SecondaryProgression(up.Progression, up.Agility)
		ManCoverage = SecondaryProgression(up.Progression, up.ManCoverage)
		ZoneCoverage = SecondaryProgression(up.Progression, up.ZoneCoverage)
		FootballIQ = SecondaryProgression(up.Progression, up.FootballIQ)
		Tackle = SecondaryProgression(up.Progression, up.Tackle)
	} else if up.Position == "P" {
		// Primary Progressions
		PuntPower = PrimaryProgression(up.Progression, up.PuntPower, up.Position, up.Archetype, SnapsPerGame, "Punt Power", true)
		PuntAccuracy = PrimaryProgression(up.Progression, up.PuntAccuracy, up.Position, up.Archetype, SnapsPerGame, "Punt Accuracy", true)
		// Secondary Progressions
		ThrowPower = SecondaryProgression(up.Progression, up.ThrowPower)
		ThrowAccuracy = SecondaryProgression(up.Progression, up.ThrowAccuracy)
		RunBlock = SecondaryProgression(up.Progression, up.RunBlock)
		PassBlock = SecondaryProgression(up.Progression, up.PassBlock)
		RunDefense = SecondaryProgression(up.Progression, up.RunDefense)
		PassRush = SecondaryProgression(up.Progression, up.PassRush)
		Carrying = SecondaryProgression(up.Progression, up.Carrying)
		Tackle = SecondaryProgression(up.Progression, up.Tackle)
		RouteRunning = SecondaryProgression(up.RouteRunning, up.RouteRunning)
		Catching = SecondaryProgression(up.Progression, up.Catching)
		KickPower = SecondaryProgression(up.Progression, up.KickPower)
		KickAccuracy = SecondaryProgression(up.Progression, up.KickAccuracy)
		ManCoverage = SecondaryProgression(up.Progression, up.ManCoverage)
		ZoneCoverage = SecondaryProgression(up.Progression, up.ZoneCoverage)
		Strength = SecondaryProgression(up.Progression, up.Strength)
		Speed = SecondaryProgression(up.Progression, up.Speed)
		Agility = SecondaryProgression(up.Progression, up.Agility)
		FootballIQ = SecondaryProgression(up.Progression, up.FootballIQ)
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

	up.Progress(progressions)

	up.GetOverall()

	return up
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
		ThrowPower = RegressAttribute(np.ThrowPower, np.Age, false)
		ThrowAccuracy = RegressAttribute(np.ThrowAccuracy, np.Age, false)
		RunBlock = RegressAttribute(np.RunBlock, np.Age, false)
		PassBlock = RegressAttribute(np.PassBlock, np.Age, false)
		RunDefense = RegressAttribute(np.RunDefense, np.Age, false)
		PassRush = RegressAttribute(np.PassRush, np.Age, false)
		Carrying = RegressAttribute(np.Carrying, np.Age, false)
		Tackle = RegressAttribute(np.Tackle, np.Age, false)
		RouteRunning = RegressAttribute(np.RouteRunning, np.Age, false)
		Catching = RegressAttribute(np.Catching, np.Age, false)
		KickPower = RegressAttribute(np.KickPower, np.Age, false)
		KickAccuracy = RegressAttribute(np.KickAccuracy, np.Age, false)
		ManCoverage = RegressAttribute(np.ManCoverage, np.Age, false)
		ZoneCoverage = RegressAttribute(np.ZoneCoverage, np.Age, false)
		Strength = RegressAttribute(np.Strength, np.Age, false)
		Speed = RegressAttribute(np.Speed, np.Age, false)
		Agility = RegressAttribute(np.Agility, np.Age, false)
		FootballIQ = RegressAttribute(np.FootballIQ, np.Age, true)
		PuntPower = RegressAttribute(np.PuntPower, np.Age, false)
		PuntAccuracy = RegressAttribute(np.PuntAccuracy, np.Age, false)
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
		progress = ((1 - math.Pow((float64(input)/99), 15)) * math.Log10(float64(input)) * 1.6 * (1 + (float64(progression / 60))))
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

func DetermineIfRetiring(player structs.NFLPlayer) bool {
	if player.Age < 28 {
		return false
	}

	drafteeOverall := 56

	odds := getRetirementOdds(player.Age)
	diceRoll := util.GenerateFloatFromRange(1, 2) - 1

	if player.Age > 34 && player.Overall > drafteeOverall && diceRoll < odds {
		newRoll := util.GenerateFloatFromRange(1, 2) - 1
		return newRoll < 0.1
	}

	return diceRoll < odds && player.Overall < drafteeOverall
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

func getAverageNFLSnaps(stats []structs.NFLPlayerStats) int {
	totalSnaps := 0

	for _, stat := range stats {
		totalSnaps += stat.Snaps
	}

	var SnapsPerGame int = 0
	if len(stats) > 0 {
		SnapsPerGame = totalSnaps / 17
	}

	return SnapsPerGame
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

func getRetirementOdds(age int) float64 {
	if age < 28 {
		return 0
	} else if age == 28 {
		return 0.1
	} else if age == 29 {
		return 0.18
	} else if age == 30 {
		return 0.21
	} else if age == 31 {
		return 0.25
	} else if age == 32 {
		return 0.28
	} else if age == 33 {
		return 0.32
	} else if age == 34 {
		return 0.36
	} else if age == 35 {
		return 0.39
	} else if age == 36 {
		return 0.43
	} else if age == 37 {
		return 0.46
	}
	return 0.5
}

func RegressAttribute(attr, age int, isIQ bool) int {
	min := 0.0
	max := 0.0
	ageRequirement := 28
	mod := 24.0
	if isIQ {
		ageRequirement = 29
		mod = 28.0
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
