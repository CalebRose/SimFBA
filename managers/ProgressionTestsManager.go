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
						csvModel.Archetype, "0", strconv.Itoa(draftee.Age), strconv.Itoa(draftee.Stars),
						draftee.HighSchool, draftee.City, draftee.State, strconv.Itoa(draftee.Height),
						strconv.Itoa(draftee.Weight), strconv.Itoa(draftee.Overall), strconv.Itoa(draftee.Speed),
						strconv.Itoa(draftee.FootballIQ), strconv.Itoa(draftee.Agility), strconv.Itoa(draftee.Carrying),
						strconv.Itoa(draftee.Catching), strconv.Itoa(draftee.RouteRunning), strconv.Itoa(draftee.ZoneCoverage), strconv.Itoa(draftee.ManCoverage),
						strconv.Itoa(draftee.Strength), strconv.Itoa(draftee.Tackle), strconv.Itoa(draftee.PassBlock), strconv.Itoa(draftee.RunBlock),
						strconv.Itoa(draftee.PassRush), strconv.Itoa(draftee.RunDefense), strconv.Itoa(draftee.ThrowPower), strconv.Itoa(draftee.ThrowAccuracy),
						strconv.Itoa(draftee.KickPower), strconv.Itoa(draftee.KickAccuracy), strconv.Itoa(draftee.PuntPower), strconv.Itoa(draftee.PuntAccuracy),
						strconv.Itoa(draftee.Stamina), strconv.Itoa(draftee.Injury), csvModel.PotentialGrade, "None",
						boomBustStatus, strconv.Itoa(tier), "Draftee", csvModel.College,
					}
					csvRows = append(csvRows, playerRow)
					continue
				}
				csvModel := structs.MapPlayerToCSVModel(player)
				idStr := strconv.Itoa(int(player.ID))
				playerRow := []string{
					team.TeamName, idStr, csvModel.FirstName, csvModel.LastName, csvModel.Position,
					csvModel.Archetype, csvModel.Year, strconv.Itoa(player.Age), strconv.Itoa(player.Stars),
					player.HighSchool, player.City, player.State, strconv.Itoa(player.Height),
					strconv.Itoa(player.Weight), strconv.Itoa(player.Overall), strconv.Itoa(player.Speed),
					strconv.Itoa(player.FootballIQ), strconv.Itoa(player.Agility), strconv.Itoa(player.Carrying),
					strconv.Itoa(player.Catching), strconv.Itoa(player.RouteRunning), strconv.Itoa(player.ZoneCoverage), strconv.Itoa(player.ManCoverage),
					strconv.Itoa(player.Strength), strconv.Itoa(player.Tackle), strconv.Itoa(player.PassBlock), strconv.Itoa(player.RunBlock),
					strconv.Itoa(player.PassRush), strconv.Itoa(player.RunDefense), strconv.Itoa(player.ThrowPower), strconv.Itoa(player.ThrowAccuracy),
					strconv.Itoa(player.KickPower), strconv.Itoa(player.KickAccuracy), strconv.Itoa(player.PuntPower), strconv.Itoa(player.PuntAccuracy),
					strconv.Itoa(player.Stamina), strconv.Itoa(player.Injury), csvModel.PotentialGrade, csvModel.RedshirtStatus,
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
					csvModel.Archetype, csvModel.Year, strconv.Itoa(cp.Age), strconv.Itoa(cp.Stars),
					cp.HighSchool, cp.City, cp.State, strconv.Itoa(cp.Height),
					strconv.Itoa(cp.Weight), strconv.Itoa(cp.Overall), strconv.Itoa(cp.Speed),
					strconv.Itoa(cp.FootballIQ), strconv.Itoa(cp.Agility), strconv.Itoa(cp.Carrying),
					strconv.Itoa(cp.Catching), strconv.Itoa(cp.RouteRunning), strconv.Itoa(cp.ZoneCoverage), strconv.Itoa(cp.ManCoverage),
					strconv.Itoa(cp.Strength), strconv.Itoa(cp.Tackle), strconv.Itoa(cp.PassBlock), strconv.Itoa(cp.RunBlock),
					strconv.Itoa(cp.PassRush), strconv.Itoa(cp.RunDefense), strconv.Itoa(cp.ThrowPower), strconv.Itoa(cp.ThrowAccuracy),
					strconv.Itoa(cp.KickPower), strconv.Itoa(cp.KickAccuracy), strconv.Itoa(cp.PuntPower), strconv.Itoa(cp.PuntAccuracy),
					strconv.Itoa(cp.Stamina), strconv.Itoa(cp.Injury), csvModel.PotentialGrade, csvModel.RedshirtStatus,
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
				csvModel.Archetype, csvModel.Year, strconv.Itoa(player.Age), strconv.Itoa(int(player.PrimeAge)), strconv.Itoa(player.Stars),
				player.HighSchool, "", player.State, strconv.Itoa(player.Height),
				strconv.Itoa(player.Weight), strconv.Itoa(player.Overall), strconv.Itoa(player.Speed),
				strconv.Itoa(player.FootballIQ), strconv.Itoa(player.Agility), strconv.Itoa(player.Carrying),
				strconv.Itoa(player.Catching), strconv.Itoa(player.RouteRunning), strconv.Itoa(player.ZoneCoverage), strconv.Itoa(player.ManCoverage),
				strconv.Itoa(player.Strength), strconv.Itoa(player.Tackle), strconv.Itoa(player.PassBlock), strconv.Itoa(player.RunBlock),
				strconv.Itoa(player.PassRush), strconv.Itoa(player.RunDefense), strconv.Itoa(player.ThrowPower), strconv.Itoa(player.ThrowAccuracy),
				strconv.Itoa(player.KickPower), strconv.Itoa(player.KickAccuracy), strconv.Itoa(player.PuntPower), strconv.Itoa(player.PuntAccuracy),
				strconv.Itoa(player.Stamina), strconv.Itoa(player.Injury), csvModel.PotentialGrade, "",
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
			csvModel.Archetype, csvModel.Year, strconv.Itoa(player.Age), strconv.Itoa(int(player.PrimeAge)), strconv.Itoa(player.Stars),
			player.HighSchool, "", player.State, strconv.Itoa(player.Height),
			strconv.Itoa(player.Weight), strconv.Itoa(player.Overall), strconv.Itoa(player.Speed),
			strconv.Itoa(player.FootballIQ), strconv.Itoa(player.Agility), strconv.Itoa(player.Carrying),
			strconv.Itoa(player.Catching), strconv.Itoa(player.RouteRunning), strconv.Itoa(player.ZoneCoverage), strconv.Itoa(player.ManCoverage),
			strconv.Itoa(player.Strength), strconv.Itoa(player.Tackle), strconv.Itoa(player.PassBlock), strconv.Itoa(player.RunBlock),
			strconv.Itoa(player.PassRush), strconv.Itoa(player.RunDefense), strconv.Itoa(player.ThrowPower), strconv.Itoa(player.ThrowAccuracy),
			strconv.Itoa(player.KickPower), strconv.Itoa(player.KickAccuracy), strconv.Itoa(player.PuntPower), strconv.Itoa(player.PuntAccuracy),
			strconv.Itoa(player.Stamina), strconv.Itoa(player.Injury), csvModel.PotentialGrade, "",
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
