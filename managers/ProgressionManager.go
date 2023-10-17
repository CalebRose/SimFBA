package managers

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/CalebRose/SimFBA/dbprovider"
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
	for _, team := range collegeTeams {
		var graduatingPlayers []structs.NFLDraftee
		teamID := strconv.Itoa(int(team.ID))
		roster := GetAllCollegePlayersByTeamId(teamID)
		croots := GetSignedRecruitsByTeamProfileID(teamID)

		if !team.PlayersProgressed {
			for _, player := range roster {
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
					draftee := structs.NFLDraftee{}
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

					// Add Graduating Player to List
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

			// Graduating players
			for _, grad := range graduatingPlayers {
				err := db.Create(&grad).Error
				if err != nil {
					log.Panicln("Could not save graduating players")
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
	var graduatingPlayers []structs.NFLDraftee

	for _, player := range unsignedPlayers {
		player = ProgressUnsignedPlayer(player, SeasonID)
		if (player.IsRedshirt && player.Year > 5) ||
			(!player.IsRedshirt && player.Year > 4) {
			player.GraduatePlayer()
			draftee := structs.NFLDraftee{}
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

	// nflPlayers := GetAllNFLPlayers()

	// for _, player := range nflPlayers {
	// 	// Progress the Player

	// 	// Get Contract if Applicable
	// 	// Reduce contract by one year
	// 	// If contract length == 0, check for extension contract offer
	// 	// if no offer, set player as free agent

	// 	// Get Stats from latest two seasons
	// 	// Do an age check
	// 	// If player is ready to retire or has not seen the field in two seasons, set player to retire

	// }

	// For all retired players
	// Create retiredplayer record
	// Delete NFL Player Record
	// Retire Contracts?

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

func PrimaryProgression(progression int, input int, position string, archetype string, spg int, attribute string, isRedshirting bool) int {
	if input == 0 {
		return 1
	}

	modifier := GetModifiers(position, archetype, spg, attribute)

	var progress float64 = 0

	if !isRedshirting {
		progress = ((1 - math.Pow((float64(input)/99.0), 15)) * math.Log10(float64(input)) * (0.6 + modifier)) * (1 + (float64(progression) / 70))
	} else {
		progress = ((1 - math.Pow((float64(input)/99), 15)) * math.Log10(float64(input)) * 1.225 * (1 + (float64(progression / 60))))
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

func DetermineIfDeclaring(player structs.CollegePlayer, avgSnaps int) bool {
	// Redshirt senior or just a senior
	if (player.IsRedshirt && player.Year == 5) || (!player.IsRedshirt && player.Year == 4 && !player.IsRedshirting) {
		return true
	}
	ovr := player.Overall
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
	if ovr > 55 && odds <= 25 {
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
	} else if ovr > 63 {
		return true
	}
	return false
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
