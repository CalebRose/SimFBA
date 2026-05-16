package managers

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/CalebRose/SimFBA/models"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
)

func CFBProgressionExport(w http.ResponseWriter) {
	w.Header().Set("Content-Disposition", "attachment;filename=2025_progression_sample_six.csv")
	w.Header().Set("Transfer-Encoding", "chunked")
	// Initialize writer
	writer := csv.NewWriter(w)
	HeaderRow := []string{
		"Team", "Player ID", "First Name", "Last Name", "Position",
		"Archetype", "Year", "Age", "Stars",
		"High School", "City", "State", "Height",
		"Weight", "Overall", "Speed",
		"Football IQ", "Agility", "Carrying",
		"Catching", "Route Running", "Zone Coverage", "Man Coverage",
		"Strength", "Tackle", "Pass Block", "Run Block",
		"Pass Rush", "Run Defense", "Throw Power", "Throw Accuracy",
		"Kick Power", "Kick Accuracy", "Punt Power", "Punt Accuracy",
		"Stamina", "Injury", "Potential Grade", "Redshirt Status",
		"BoomBustStatus", "Tier", "DraftStatus", "College",
	}
	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}
	ts := GetTimestamp()
	SeasonID := strconv.Itoa(ts.CollegeSeasonID)
	// Get All Teams
	snapMap := GetCollegePlayerSeasonSnapMap(SeasonID)
	statMap := GetCollegePlayerStatsMap(SeasonID)
	collegeTeams := GetAllCollegeTeams()
	csvRows := [][]string{}
	// Loop
	for _, team := range collegeTeams {
		teamID := strconv.Itoa(int(team.ID))
		roster := GetAllCollegePlayersByTeamId(teamID)
		croots := GetSignedRecruitsByTeamProfileID(teamID)
		fmt.Println("Progressing " + team.TeamAbbr + "...")

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
					newProgression := util.GenerateNFLPotential(int(player.Progression))
					newPotentialGrade := util.GetWeightedPotentialGrade(int8(newProgression))
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
					switch diceRoll {
					case 1:
						boomBustStatus = "Bust"
						enableBoomBust = true
						// Bust
						fmt.Println("BUST!")
						draftee.AssignBoomBustStatus(boomBustStatus)

					case 20:
						enableBoomBust = true
						// Boom
						fmt.Println("BOOM!")
						boomBustStatus = "Boom"
						isBoom = true
						draftee.AssignBoomBustStatus(boomBustStatus)
					default:
						tier = 0
					}
					if enableBoomBust {
						for i := 0; i < tier; i++ {
							draftee = BoomBustDraftee(draftee, SeasonID, 51, isBoom)
						}
					}
					idStr := strconv.Itoa(int(draftee.ID))
					csvModel := models.MapNFLDrafteeToModel(draftee)
					playerRow := []string{
						"", idStr, csvModel.FirstName, csvModel.LastName, csvModel.Position,
						csvModel.Archetype, "0", strconv.Itoa(int(draftee.Age)), strconv.Itoa(int(draftee.Stars)),
						draftee.HighSchool, draftee.City, draftee.State, strconv.Itoa(int(draftee.Height)),
						strconv.Itoa(int(draftee.Weight)), strconv.Itoa(int(draftee.Overall)), strconv.Itoa(int(draftee.Speed)),
						strconv.Itoa(int(draftee.FootballIQ)), strconv.Itoa(int(draftee.Agility)), strconv.Itoa(int(draftee.Carrying)),
						strconv.Itoa(int(draftee.Catching)), strconv.Itoa(int(draftee.RouteRunning)), strconv.Itoa(int(draftee.ZoneCoverage)), strconv.Itoa(int(draftee.ManCoverage)),
						strconv.Itoa(int(draftee.Strength)), strconv.Itoa(int(draftee.Tackle)), strconv.Itoa(int(draftee.PassBlock)), strconv.Itoa(int(draftee.RunBlock)),
						strconv.Itoa(int(draftee.PassRush)), strconv.Itoa(int(draftee.RunDefense)), strconv.Itoa(int(draftee.ThrowPower)), strconv.Itoa(int(draftee.ThrowAccuracy)),
						strconv.Itoa(int(draftee.KickPower)), strconv.Itoa(int(draftee.KickAccuracy)), strconv.Itoa(int(draftee.PuntPower)), strconv.Itoa(int(draftee.PuntAccuracy)),
						strconv.Itoa(int(draftee.Stamina)), strconv.Itoa(int(draftee.Injury)), csvModel.PotentialGrade, "None",
						boomBustStatus, strconv.Itoa(tier), "Draftee", csvModel.College,
					}
					csvRows = append(csvRows, playerRow)
					continue
				}
				csvModel := structs.MapPlayerToCSVModel(player)
				idStr := strconv.Itoa(int(player.ID))
				playerRow := []string{
					team.TeamName, idStr, csvModel.FirstName, csvModel.LastName, csvModel.Position,
					csvModel.Archetype, csvModel.Year, strconv.Itoa(int(player.Age)), strconv.Itoa(int(player.Stars)),
					player.HighSchool, player.City, player.State, strconv.Itoa(int(player.Height)),
					strconv.Itoa(int(player.Weight)), strconv.Itoa(int(player.Overall)), strconv.Itoa(int(player.Speed)),
					strconv.Itoa(int(player.FootballIQ)), strconv.Itoa(int(player.Agility)), strconv.Itoa(int(player.Carrying)),
					strconv.Itoa(int(player.Catching)), strconv.Itoa(int(player.RouteRunning)), strconv.Itoa(int(player.ZoneCoverage)), strconv.Itoa(int(player.ManCoverage)),
					strconv.Itoa(int(player.Strength)), strconv.Itoa(int(player.Tackle)), strconv.Itoa(int(player.PassBlock)), strconv.Itoa(int(player.RunBlock)),
					strconv.Itoa(int(player.PassRush)), strconv.Itoa(int(player.RunDefense)), strconv.Itoa(int(player.ThrowPower)), strconv.Itoa(int(player.ThrowAccuracy)),
					strconv.Itoa(int(player.KickPower)), strconv.Itoa(int(player.KickAccuracy)), strconv.Itoa(int(player.PuntPower)), strconv.Itoa(int(player.PuntAccuracy)),
					strconv.Itoa(int(player.Stamina)), strconv.Itoa(int(player.Injury)), csvModel.PotentialGrade, csvModel.RedshirtStatus,
					"None", "", "Collegiate", "",
				}
				csvRows = append(csvRows, playerRow)
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
				boomBustStatus := "None"
				// Generate Tier
				switch diceRoll {
				case 1:
					boomBustStatus = "Bust"
					enableBoomBust = true
				case 20:
					boomBustStatus = "Boom"
					enableBoomBust = true
					isBoom = true
				default:
					tier = 0
				}
				if enableBoomBust {
					for i := 0; i < tier; i++ {
						cp = BoomBustRecruit(cp, SeasonID, 51, isBoom)
					}
				}

				fmt.Println("Adding " + croot.FirstName + " " + croot.LastName + "to " + team.TeamAbbr)

				csvModel := structs.MapPlayerToCSVModel(cp)
				idStr := strconv.Itoa(int(cp.ID))
				playerRow := []string{
					team.TeamName, idStr, csvModel.FirstName, csvModel.LastName, csvModel.Position,
					csvModel.Archetype, csvModel.Year, strconv.Itoa(int(cp.Age)), strconv.Itoa(int(cp.Stars)),
					cp.HighSchool, cp.City, cp.State, strconv.Itoa(int(cp.Height)),
					strconv.Itoa(int(cp.Weight)), strconv.Itoa(int(cp.Overall)), strconv.Itoa(int(cp.Speed)),
					strconv.Itoa(int(cp.FootballIQ)), strconv.Itoa(int(cp.Agility)), strconv.Itoa(int(cp.Carrying)),
					strconv.Itoa(int(cp.Catching)), strconv.Itoa(int(cp.RouteRunning)), strconv.Itoa(int(cp.ZoneCoverage)), strconv.Itoa(int(cp.ManCoverage)),
					strconv.Itoa(int(cp.Strength)), strconv.Itoa(int(cp.Tackle)), strconv.Itoa(int(cp.PassBlock)), strconv.Itoa(int(cp.RunBlock)),
					strconv.Itoa(int(cp.PassRush)), strconv.Itoa(int(cp.RunDefense)), strconv.Itoa(int(cp.ThrowPower)), strconv.Itoa(int(cp.ThrowAccuracy)),
					strconv.Itoa(int(cp.KickPower)), strconv.Itoa(int(cp.KickAccuracy)), strconv.Itoa(int(cp.PuntPower)), strconv.Itoa(int(cp.PuntAccuracy)),
					strconv.Itoa(int(cp.Stamina)), strconv.Itoa(int(cp.Injury)), csvModel.PotentialGrade, csvModel.RedshirtStatus,
					boomBustStatus, strconv.Itoa(tier), "Collegiate", "",
				}
				csvRows = append(csvRows, playerRow)
			}
			team.ToggleRecruitsAdded()
		}
	}
	fmt.Println("Exporting all data to csv...")
	for _, row := range csvRows {
		err = writer.Write(row)
		if err != nil {
			log.Fatal("Cannot write player row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
}

func NFLProgressionExport(w http.ResponseWriter) {
	w.Header().Set("Content-Disposition", "attachment;filename=2025_nfl_progression_sample_nine.csv")
	w.Header().Set("Transfer-Encoding", "chunked")
	// Initialize writer
	writer := csv.NewWriter(w)
	HeaderRow := []string{
		"Team", "Player ID", "First Name", "Last Name", "Position",
		"Archetype", "Year", "Age", "Prime Age", "Stars",
		"High School", "City", "State", "Height",
		"Weight", "Overall", "Speed",
		"Football IQ", "Agility", "Carrying",
		"Catching", "Route Running", "Zone Coverage", "Man Coverage",
		"Strength", "Tackle", "Pass Block", "Run Block",
		"Pass Rush", "Run Defense", "Throw Power", "Throw Accuracy",
		"Kick Power", "Kick Accuracy", "Punt Power", "Punt Accuracy",
		"Stamina", "Injury", "Potential Grade", "Redshirt Status",
		"RetiringStatus", "College",
	}
	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}
	ts := GetTimestamp()
	SeasonID := strconv.Itoa(ts.CollegeSeasonID)
	// Get All Teams
	snapMap := GetNFLPlayerSeasonSnapMap(SeasonID)
	statMap := GetNFLPlayerStatsMap(SeasonID)
	teams := GetAllNFLTeams()
	freeAgents := GetAllFreeAgents()
	lastTwoStatMap := GetNFLLastTwoSeasonStatMap(ts.NFLSeasonID)

	csvRows := [][]string{}
	// Loop
	for _, team := range teams {
		teamID := strconv.Itoa(int(team.ID))
		nflPlayers := GetNFLPlayersRecordsByTeamID(teamID)
		fmt.Println("Progressing " + team.TeamAbbr + "...")

		for _, player := range nflPlayers {
			if player.HasProgressed {
				continue
			}

			// Get Latest Stats
			stats := statMap[player.ID]
			snaps := snapMap[player.ID]
			totalSnaps, avgSnaps := getAverageNFLSnaps(stats)

			// Run Function to Determine if Player is retiring
			willRetire := DetermineIfRetiring(player, lastTwoStatMap)

			// Progress the Player
			player = ProgressNFLPlayer(player, SeasonID, totalSnaps, avgSnaps, snaps)
			retireStatus := ""
			if willRetire {
				retireStatus = "Retiring"
			}

			csvModel := structs.MapNFLPlayerToCSVModel(player)
			idStr := strconv.Itoa(int(player.ID))
			playerRow := []string{
				team.TeamName, idStr, csvModel.FirstName, csvModel.LastName, csvModel.Position,
				csvModel.Archetype, csvModel.Year, strconv.Itoa(int(player.Age)), strconv.Itoa(int(player.PrimeAge)), strconv.Itoa(int(player.Stars)),
				player.HighSchool, "", player.State, strconv.Itoa(int(player.Height)),
				strconv.Itoa(int(player.Weight)), strconv.Itoa(int(player.Overall)), strconv.Itoa(int(player.Speed)),
				strconv.Itoa(int(player.FootballIQ)), strconv.Itoa(int(player.Agility)), strconv.Itoa(int(player.Carrying)),
				strconv.Itoa(int(player.Catching)), strconv.Itoa(int(player.RouteRunning)), strconv.Itoa(int(player.ZoneCoverage)), strconv.Itoa(int(player.ManCoverage)),
				strconv.Itoa(int(player.Strength)), strconv.Itoa(int(player.Tackle)), strconv.Itoa(int(player.PassBlock)), strconv.Itoa(int(player.RunBlock)),
				strconv.Itoa(int(player.PassRush)), strconv.Itoa(int(player.RunDefense)), strconv.Itoa(int(player.ThrowPower)), strconv.Itoa(int(player.ThrowAccuracy)),
				strconv.Itoa(int(player.KickPower)), strconv.Itoa(int(player.KickAccuracy)), strconv.Itoa(int(player.PuntPower)), strconv.Itoa(int(player.PuntAccuracy)),
				strconv.Itoa(int(player.Stamina)), strconv.Itoa(int(player.Injury)), csvModel.PotentialGrade, "",
				retireStatus, player.College,
			}
			csvRows = append(csvRows, playerRow)
		}

	}

	for _, player := range freeAgents {
		if player.HasProgressed {
			continue
		}
		// Get Latest Stats
		stats := statMap[player.ID]
		snaps := snapMap[player.ID]
		totalSnaps, avgSnaps := getAverageNFLSnaps(stats)

		// Run Function to Determine if Player is retiring
		willRetire := DetermineIfRetiring(player, lastTwoStatMap)
		// Progress the Player
		player = ProgressNFLPlayer(player, SeasonID, totalSnaps, avgSnaps, snaps)
		retireStatus := ""
		if willRetire {
			retireStatus = "Retiring"
		}
		csvModel := structs.MapNFLPlayerToCSVModel(player)
		idStr := strconv.Itoa(int(player.ID))
		playerRow := []string{
			"FA", idStr, csvModel.FirstName, csvModel.LastName, csvModel.Position,
			csvModel.Archetype, csvModel.Year, strconv.Itoa(int(player.Age)), strconv.Itoa(int(player.PrimeAge)), strconv.Itoa(int(player.Stars)),
			player.HighSchool, "", player.State, strconv.Itoa(int(player.Height)),
			strconv.Itoa(int(player.Weight)), strconv.Itoa(int(player.Overall)), strconv.Itoa(int(player.Speed)),
			strconv.Itoa(int(player.FootballIQ)), strconv.Itoa(int(player.Agility)), strconv.Itoa(int(player.Carrying)),
			strconv.Itoa(int(player.Catching)), strconv.Itoa(int(player.RouteRunning)), strconv.Itoa(int(player.ZoneCoverage)), strconv.Itoa(int(player.ManCoverage)),
			strconv.Itoa(int(player.Strength)), strconv.Itoa(int(player.Tackle)), strconv.Itoa(int(player.PassBlock)), strconv.Itoa(int(player.RunBlock)),
			strconv.Itoa(int(player.PassRush)), strconv.Itoa(int(player.RunDefense)), strconv.Itoa(int(player.ThrowPower)), strconv.Itoa(int(player.ThrowAccuracy)),
			strconv.Itoa(int(player.KickPower)), strconv.Itoa(int(player.KickAccuracy)), strconv.Itoa(int(player.PuntPower)), strconv.Itoa(int(player.PuntAccuracy)),
			strconv.Itoa(int(player.Stamina)), strconv.Itoa(int(player.Injury)), csvModel.PotentialGrade, "",
			retireStatus, player.College,
		}
		csvRows = append(csvRows, playerRow)
	}
	fmt.Println("Exporting all data to csv...")
	for _, row := range csvRows {
		err = writer.Write(row)
		if err != nil {
			log.Fatal("Cannot write player row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
}
