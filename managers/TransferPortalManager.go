package managers

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"sync"

	"github.com/CalebRose/SimFBA/constants"
	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/models"
	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
	"gorm.io/gorm"
)

var upcomingTeam = "Prefers to play for an up-and-coming team"
var differentState = "Prefers to play in a different state"
var immediateStart = "Prefers to play for a team where he can start immediately"
var closeToHome = "Prefers to be close to home"
var nationalChampionshipContender = "Prefers to play for a national championship contender"
var specificCoach = "Prefers to play for a specific coach"
var legacy = "Legacy"
var richHistory = "Prefers to play for a team with a rich history"

func ProcessTransferIntention(w http.ResponseWriter) {
	db := dbprovider.GetInstance().GetDB()
	// w.Header().Set("Content-Disposition", "attachment;filename=transferStats.csv")
	// w.Header().Set("Transfer-Encoding", "chunked")
	// Initialize writer
	// writer := csv.NewWriter(w)
	ts := GetTimestamp()
	seasonID := strconv.Itoa(ts.CollegeSeasonID)
	allCollegePlayers := GetAllCollegePlayers()
	seasonSnapMap := GetALLCFBPlayerSeasonStatMapBySeason(seasonID)
	fullRosterMap := GetFullTeamRosterWithCrootsMap()
	teamProfileMap := GetTeamProfileMap()
	transferCount := 0
	freshmanCount := 0
	redshirtFreshmanCount := 0
	sophomoreCount := 0
	redshirtSophomoreCount := 0
	juniorCount := 0
	redshirtJuniorCount := 0
	seniorCount := 0
	redshirtSeniorCount := 0
	lowCount := 0
	mediumCount := 0
	highCount := 0
	bigDrop := -25.0
	smallDrop := -10.0
	tinyDrop := -5.0
	tinyGain := 5.0
	smallGain := 10.0
	mediumGain := 15.0
	mediumDrop := -15.0
	bigGain := 25.0

	// HeaderRow := []string{
	// 	"Team", "First Name", "Last Name", "Stars",
	// 	"Archetype", "Position", "Year", "Age", "Redshirt Status",
	// 	"Overall", "Transfer Bias", "Transfer Status", "Transfer Weight", "Dice Roll",
	// 	"Age Mod", "Snap Mod", "star Mod", "DC Comp Mod", "Scheme Mod", "FCS Mod",
	// }

	// err := writer.Write(HeaderRow)
	// if err != nil {
	// 	log.Fatal("Cannot write header row", err)
	// }

	for _, p := range allCollegePlayers {
		// Do not include redshirts and all graduating players
		if p.IsRedshirting || p.TeamID > 194 || p.TeamID == 0 {
			continue
		}
		// Weight will be the initial barrier required for a player to consider transferring.
		// The lower the number gets, the more likely the player will transfer
		transferWeight := 0.0

		// Modifiers on reasons why they would transfer
		snapMod := 0.0
		ageMod := 1.125
		starMod := 0.0
		depthChartCompetitionMod := 0.0
		schemeMod := 0.0
		// closeToHomeMod := 0.0

		// Check Snaps
		seasonStats := seasonSnapMap[p.ID]
		totalSnaps := 0
		for _, s := range seasonStats {
			totalSnaps += s.Snaps
		}
		snapsPerGame := totalSnaps / 12

		if p.Position == "P" || p.Position == "K" {
			if snapsPerGame > 1 {
				snapMod = bigDrop
			} else {
				snapMod = smallGain
			}
		} else if p.Position == "QB" {
			if snapsPerGame > 50 {
				snapMod = bigDrop
			} else if snapsPerGame > 30 {
				snapMod = tinyGain
			} else if snapsPerGame > 15 {
				snapMod = mediumGain
			} else {
				snapMod = bigGain
			}
		} else if p.Position == "OG" || p.Position == "OT" || p.Position == "C" {
			if snapsPerGame > 50 {
				snapMod = bigDrop
			} else if snapsPerGame > 30 {
				snapMod = tinyGain
			} else if snapsPerGame > 15 {
				snapMod = mediumGain
			} else {
				snapMod = bigGain
			}
		} else if p.Position == "SS" || p.Position == "FS" || p.Position == "RB" {
			/// next positions are usually single full time starters, but could also have backups playing ST so more tiers needed
			if snapsPerGame > 35 {
				snapMod = bigDrop
			} else if snapsPerGame > 15 {
				snapMod = tinyGain
			} else if snapsPerGame > 5 {
				snapMod = mediumGain
			} else {
				snapMod = bigGain
			}
		} else if p.Position == "FB" || p.Position == "TE" {
			if snapsPerGame > 30 {
				snapMod = bigDrop
			} else if snapsPerGame > 10 {
				snapMod = tinyDrop
			} else if snapsPerGame > 1 {
				snapMod = mediumGain
			} else {
				snapMod = bigGain
			}
		} else if p.Position == "RB" || p.Position == "DT" || p.Position == "ILB" {
			if snapsPerGame > 35 {
				snapMod = bigDrop
			} else if snapsPerGame > 20 {
				snapMod = smallDrop
			} else if snapsPerGame > 10 {
				snapMod = smallGain
			} else {
				snapMod = bigGain
			}
		} else if p.Position == "CB" || p.Position == "OLB" || p.Position == "WR" || p.Position == "DE" {
			if snapsPerGame > 35 {
				snapMod = bigDrop
			} else if snapsPerGame > 25 {
				snapMod = smallDrop
			} else if snapsPerGame > 15 {
				snapMod = smallGain
			} else if snapsPerGame > 5 {
				snapMod = mediumGain
			} else {
				snapMod = bigGain
			}
		}

		// Check Age
		// The more experienced the player is in the league,
		// the more likely they will transfer.
		/// Have this be a multiplicative factor to odds
		if p.Year == 1 {
			ageMod = .01
		} else if p.Year == 2 && p.IsRedshirt {
			ageMod = .1
		} else if p.Year == 2 && !p.IsRedshirt {
			ageMod = .4
		} else if p.Year == 3 && p.IsRedshirt {
			ageMod = .7
		} else if p.Year == 3 && !p.IsRedshirt {
			ageMod = 1
		} else if p.Year == 4 {
			ageMod = 1.25
		} else if p.Year == 5 {
			ageMod = 1.45
		}

		/// Higher star players are more likely to transfer
		if p.Stars == 0 {
			starMod = 1
		} else if p.Stars == 1 {
			starMod = .66
		} else if p.Stars == 2 {
			starMod = .75
		} else if p.Stars == 3 {
			starMod = util.GenerateFloatFromRange(0.9, 1.1)
		} else if p.Stars == 4 {
			starMod = util.GenerateFloatFromRange(1.11, 1.3)
		} else if p.Stars == 5 {
			starMod = util.GenerateFloatFromRange(1.31, 1.75)
		}

		// Check Team Position Rank
		teamRoster := fullRosterMap[uint(p.TeamID)]
		filteredRosterByPosition := filterRosterByPosition(teamRoster, p.Position)
		youngerPlayerAhead := false
		idFound := false
		for idx, pl := range filteredRosterByPosition {
			if pl.Age < p.Age && !idFound {
				youngerPlayerAhead = true
			}
			if pl.ID == p.ID {
				idFound = true
				// Check the index of the player.
				// If they're at the top of the list, they're considered to be starting caliber.
				if (p.Position == "QB" ||
					p.Position == "P" ||
					p.Position == "K" ||
					p.Position == "FB" ||
					p.Position == "C") && idx > 1 {
					depthChartCompetitionMod += 33
				}

				if (p.Position == "RB" ||
					p.Position == "TE" ||
					p.Position == "FS" ||
					p.Position == "OT" ||
					p.Position == "OG" ||
					p.Position == "DT" ||
					p.Position == "DE" ||
					p.Position == "OLB" ||
					p.Position == "ILB" ||
					p.Position == "SS") && idx > 2 {
					depthChartCompetitionMod += 33
				}

				if (p.Position == "WR" ||
					p.Position == "CB") && idx > 3 {
					depthChartCompetitionMod += 33
				}
			}
		}

		// If there's a modifier applied and there's a younger player ahead on the roster, double the amount on the modifier
		if depthChartCompetitionMod > 0 {
			if youngerPlayerAhead {
				depthChartCompetitionMod += 33
			} else {
				depthChartCompetitionMod = .63 * depthChartCompetitionMod
			}
		}

		// Check for scheme based on Team Recruiting Profile.
		// If it is not a good fit for the player, they will want to transfer
		// Will Need to Lock Scheme Dropdown by halfway through the season or by end of season

		teamID := p.TeamID
		if teamID == 0 {
			teamID = int(p.PreviousTeamID)
		}
		teamIdStr := strconv.Itoa(teamID)
		teamProfile := teamProfileMap[teamIdStr]
		schemeMod = getSchemeMod(teamProfile, p, mediumDrop, mediumGain)

		fcsMod := 1.0
		if p.TeamID > 134 && p.TeamID != 138 && p.TeamID != 206 {
			if p.Year > 2 && p.Overall > 39 {
				fcsMod += (0.1 * float64(p.Year))
			}
			if p.Personality == "Loyal" {
				fcsMod = 0.0
			}
		}

		/// Not playing = 25, low depth chart = 16 or 33, scheme = 10, if you're all 3, that's a ~60% chance of transferring pre- modifiers
		transferWeight = starMod * ageMod * (snapMod + depthChartCompetitionMod + schemeMod) * fcsMod
		diceRoll := util.GenerateIntFromRange(1, 100)

		// NOT INTENDING TO TRANSFER
		transferInt := int(transferWeight)
		if diceRoll > transferInt {
			continue
		}

		if p.Year == 1 {
			fmt.Println("Dice Roll: ", diceRoll)
		}

		// Is Intending to transfer
		p.DeclareTransferIntention(int(transferWeight))
		transferCount++
		if p.Year == 1 && !p.IsRedshirt {
			freshmanCount++
		} else if p.Year == 2 && p.IsRedshirt {
			redshirtFreshmanCount++
		} else if p.Year == 2 && !p.IsRedshirt {
			sophomoreCount++
		} else if p.Year == 3 && p.IsRedshirt {
			redshirtSophomoreCount++
		} else if p.Year == 3 && !p.IsRedshirt {
			juniorCount++
		} else if p.Year == 4 && p.IsRedshirt {
			redshirtJuniorCount++
		} else if p.Year == 4 && !p.IsRedshirt {
			seniorCount++
		} else if p.Year == 5 && p.IsRedshirt {
			redshirtSeniorCount++
		}

		if transferWeight < 30 {
			lowCount++
		} else if transferWeight < 70 {
			mediumCount++
		} else {
			highCount++
		}

		repository.SaveCFBPlayer(p, db)
		if p.Stars > 2 {
			message := "Breaking News! " + strconv.Itoa(p.Stars) + " star " + p.Position + " " + p.FirstName + " " + p.LastName + " has announced their intention to transfer from " + p.TeamAbbr + "!"
			CreateNewsLog("CFB", message, "Transfer Portal", int(p.TeamID), ts)
		}
		notificationMessage := strconv.Itoa(p.Stars) + " star " + p.Position + " " + p.FirstName + " " + p.LastName + " has a " + p.TransferLikeliness + " likeliness of entering the transfer portal. Please navigate to the Roster page to submit a promise."
		CreateNotification("CFB", notificationMessage, "Transfer Intention", uint(p.TeamID))
		// fmt.Println(strconv.Itoa(p.Year)+" YEAR "+p.TeamAbbr+" "+p.Position+" "+p.FirstName+" "+p.LastName+" HAS ANNOUNCED THEIR INTENTION TO TRANSFER | Weight: ", int(transferWeight))
		// // db.Save(&p)
		// csvModel := structs.MapPlayerToCSVModel(p)
		// playerRow := []string{
		// 	p.TeamAbbr, csvModel.FirstName, csvModel.LastName, strconv.Itoa(p.Stars),
		// 	csvModel.Archetype, csvModel.Position,
		// 	csvModel.Year, strconv.Itoa(p.Age), csvModel.RedshirtStatus,
		// 	csvModel.OverallGrade, p.RecruitingBias, p.TransferLikeliness, strconv.Itoa(transferInt), strconv.Itoa(diceRoll),
		// 	fmt.Sprintf("%.3f", ageMod), fmt.Sprintf("%.3f", snapMod), fmt.Sprintf("%.3f", starMod), fmt.Sprintf("%.3f", depthChartCompetitionMod), fmt.Sprintf("%.3f", schemeMod), fmt.Sprintf("%.3f", fcsMod),
		// }

		// err = writer.Write(playerRow)
		// if err != nil {
		// 	log.Fatal("Cannot write player row to CSV", err)
		// }

		// writer.Flush()
		// err = writer.Error()
		// if err != nil {
		// 	log.Fatal("Error while writing to file ::", err)
		// }
	}
	transferPortalMessage := "Breaking News! About " + strconv.Itoa(transferCount) + " players intend to transfer from their current schools. Teams have one week to commit promises to retain players."
	CreateNewsLog("CFB", transferPortalMessage, "Transfer Portal", 0, ts)
	ts.EnactPromisePhase()
	repository.SaveTimestamp(ts, db)
	fmt.Println("Total number of players entering the transfer portal: ", transferCount)
	fmt.Println("Total number of freshmen entering the transfer portal: ", freshmanCount)
	fmt.Println("Total number of redshirt freshmen entering the transfer portal: ", redshirtFreshmanCount)
	fmt.Println("Total number of sophomores entering the transfer portal: ", sophomoreCount)
	fmt.Println("Total number of redshirt sophomores entering the transfer portal: ", redshirtSophomoreCount)
	fmt.Println("Total number of juniors entering the transfer portal: ", juniorCount)
	fmt.Println("Total number of redshirt juniors entering the transfer portal: ", redshirtJuniorCount)
	fmt.Println("Total number of seniors entering the transfer portal: ", seniorCount)
	fmt.Println("Total number of redshirt seniors entering the transfer portal: ", redshirtSeniorCount)
	fmt.Println("Total number of players with low likeliness to enter transfer portal: ", lowCount)
	fmt.Println("Total number of players with medium likeliness to enter transfer portal: ", mediumCount)
	fmt.Println("Total number of players with high likeliness to enter transfer portal: ", highCount)
}

func AICoachPromisePhase() {
	db := dbprovider.GetInstance().GetDB()

	aiTeamProfiles := GetOnlyAITeamRecruitingProfiles()

	coachMap := GetActiveCollegeCoachMap()

	for _, team := range aiTeamProfiles {
		if !team.IsAI || team.ID > 194 || team.ID == 0 || team.IsUserTeam {
			continue
		}
		coach := coachMap[team.ID]
		if coach.ID == 0 {
			continue
		}
		teamID := strconv.Itoa(int(team.ID))
		roster := GetAllCollegePlayersByTeamId(teamID)
		for _, p := range roster {
			if p.TransferStatus > 1 || p.TransferStatus == 0 {
				continue
			}
			collegePlayerID := strconv.Itoa(int(p.ID))
			promise := GetCollegePromiseByCollegePlayerID(collegePlayerID, teamID)
			if promise.ID != 0 {
				continue
			}

			promiseOdds := getBasePromiseOdds(coach.TeambuildingPreference, coach.PromiseTendency)
			diceRoll := util.GenerateIntFromRange(1, 100)

			if diceRoll < promiseOdds {
				// Commit Promise
				promiseLevel := getPromiseLevel(coach.PromiseTendency)
				promiseWeight := "Medium"
				promiseType := ""
				benchmarkStr := ""
				promiseBenchmark := 0

				bias := p.RecruitingBias
				if bias == closeToHome {
					promiseType = "Home State Game"
					benchmarkStr = p.State
				} else if bias == immediateStart && p.Overall > 42 {
					promiseType = "Snaps"
					// REWRITE
					promiseBenchmark = 0
					if promiseLevel == 1 {
						promiseBenchmark += 10
					} else if promiseLevel == -1 {
						promiseBenchmark -= 1
					}

					promiseWeight = getPromiseWeightBySnapsOrWins(p.Position, "Snap Count", promiseBenchmark)
				} else if bias == nationalChampionshipContender || bias == richHistory {
					// Promise based on wins
					promiseBenchmark = 6
					promiseType = "Wins"
					if promiseLevel == 1 {
						promiseBenchmark += 3
					} else if promiseLevel == -1 {
						promiseBenchmark -= 3
					}
					promiseWeight = getPromiseWeightBySnapsOrWins(p.Position, "Wins", promiseBenchmark)
				}

				if promiseType == "" {
					continue
				}

				collegePromise := structs.CollegePromise{
					TeamID:          team.ID,
					CollegePlayerID: p.ID,
					PromiseType:     promiseType,
					PromiseWeight:   promiseWeight,
					Benchmark:       promiseBenchmark,
					BenchmarkStr:    benchmarkStr,
					IsActive:        true,
				}
				repository.CreateCollegePromiseRecord(collegePromise, db)
			}
		}
	}
}

func GetAllCollegePromises() []structs.CollegePromise {
	db := dbprovider.GetInstance().GetDB()

	p := []structs.CollegePromise{}

	err := db.Find(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []structs.CollegePromise{}
		} else {
			log.Fatal(err)
		}
	}
	return p
}

func GetCollegePromiseByID(id string) structs.CollegePromise {
	db := dbprovider.GetInstance().GetDB()

	p := structs.CollegePromise{}

	err := db.Where("id = ?", id).Find(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return structs.CollegePromise{}
		} else {
			log.Fatal(err)
		}
	}
	return p
}

func GetCollegePromisesByTeamID(teamID string) []structs.CollegePromise {
	db := dbprovider.GetInstance().GetDB()

	p := []structs.CollegePromise{}

	err := db.Where("team_id = ?", teamID).Find(&p).Error
	if err != nil {
		return []structs.CollegePromise{}
	}
	return p
}

func GetCollegePromiseByCollegePlayerID(id, teamID string) structs.CollegePromise {
	db := dbprovider.GetInstance().GetDB()

	p := structs.CollegePromise{}

	err := db.Where("college_player_id = ? AND team_id = ?", id, teamID).Find(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return structs.CollegePromise{}
		} else {
			log.Fatal(err)
		}
	}
	return p
}

func CreatePromise(promise structs.CollegePromise) structs.CollegePromise {
	db := dbprovider.GetInstance().GetDB()
	collegePlayerID := strconv.Itoa(int(promise.CollegePlayerID))
	profileID := strconv.Itoa(int(promise.TeamID))

	existingPromise := GetCollegePromiseByCollegePlayerID(collegePlayerID, profileID)
	if existingPromise.ID != 0 && existingPromise.ID > 0 {
		existingPromise.Reactivate(promise.PromiseType, promise.PromiseWeight, promise.Benchmark)
		repository.SaveCollegePromiseRecord(promise, db)
		assignPromiseToProfile(db, collegePlayerID, profileID, existingPromise.ID)
		return existingPromise
	}

	db.Create(&promise)

	assignPromiseToProfile(db, collegePlayerID, profileID, promise.ID)

	return promise
}

func assignPromiseToProfile(db *gorm.DB, collegePlayerID, profileID string, id uint) {
	tpProfile := GetOnlyTransferPortalProfileByPlayerID(collegePlayerID, profileID)
	if tpProfile.ID > 0 {
		tpProfile.AssignPromise(id)
		db.Save(&tpProfile)
	}
}

func UpdatePromise(promise structs.CollegePromise) {
	db := dbprovider.GetInstance().GetDB()
	id := strconv.Itoa(int(promise.ID))
	existingPromise := GetCollegePromiseByID(id)
	existingPromise.UpdatePromise(promise.PromiseType, promise.PromiseWeight, promise.Benchmark)
	db.Save(&existingPromise)
}

func CancelPromise(id string) {
	db := dbprovider.GetInstance().GetDB()
	promise := GetCollegePromiseByID(id)
	promise.Deactivate()
	db.Save(&promise)
}

func EnterTheTransferPortal() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	// Get All Teams
	teams := GetAllCollegeTeams()

	for _, t := range teams {
		teamID := strconv.Itoa(int(t.ID))
		roster := GetAllCollegePlayersByTeamId(teamID)

		for _, p := range roster {
			if p.TransferStatus != 1 {
				continue
			}

			playerID := strconv.Itoa(int(p.ID))

			promise := GetCollegePromiseByCollegePlayerID(playerID, teamID)
			if promise.ID == 0 {
				p.WillTransfer()
				repository.SaveCollegePlayerRecord(p, db)
				continue
			}
			// 1-100
			baseFloor := getTransferFloor(p.TransferLikeliness)
			// 10, 20, 40, 60, 70
			promiseModifier := getPromiseFloor(promise.PromiseWeight)
			difference := baseFloor - promiseModifier
			// In the future, add something like a bias modifier.
			// If the coach promises something
			// That does NOT match the bias of the player, it should not be as impactful.
			// However, this should be implemented after investigating how to make bias more impactful.

			diceRoll := util.GenerateIntFromRange(1, 100)

			// Lets say the difference is 40. 60-20.
			if diceRoll < difference {
				// If the dice roll is within the 40%. They leave.
				// Okay this makes sense.

				p.WillTransfer()

				// Create News Log
				message := "Breaking News! " + p.TeamAbbr + " " + strconv.Itoa(p.Stars) + " Star " + p.Position + " " + p.FirstName + " " + p.LastName + " has officially entered the transfer portal!"
				CreateNewsLog("CFB", message, "Transfer Portal", int(p.PreviousTeamID), ts)

				repository.SaveCFBPlayer(p, db)
				repository.DeleteCollegePromise(promise, db)
				continue
			}

			// Create News Log
			message := "Breaking News! " + p.TeamAbbr + " " + strconv.Itoa(p.Stars) + " Star " + p.Position + " " + p.FirstName + " " + p.LastName + " has withdrawn their name from the transfer portal!"
			CreateNewsLog("CFB", message, "Transfer Portal", int(p.PreviousTeamID), ts)

			promise.MakePromise()
			repository.SaveCollegePromiseRecord(promise, db)
			p.WillStay()
			repository.SaveCFBPlayer(p, db)
		}
	}

	ts.EnactPortalPhase()
	repository.SaveTimestamp(ts, db)
}

func AddTransferPlayerToBoard(transferPortalProfileDto structs.TransferPortalProfile) structs.TransferPortalProfile {
	db := dbprovider.GetInstance().GetDB()

	portalProfile := GetOnlyTransferPortalProfileByPlayerID(strconv.Itoa(int(transferPortalProfileDto.CollegePlayerID)), strconv.Itoa(int(transferPortalProfileDto.ProfileID)))

	// If Recruit Already Exists
	if portalProfile.CollegePlayerID != 0 && portalProfile.ProfileID != 0 {
		portalProfile.Reactivate()
		db.Save(&portalProfile)
		return portalProfile
	}

	newProfileForRecruit := structs.TransferPortalProfile{
		SeasonID:           uint(transferPortalProfileDto.SeasonID),
		CollegePlayerID:    uint(transferPortalProfileDto.CollegePlayerID),
		ProfileID:          uint(transferPortalProfileDto.ProfileID),
		TeamAbbreviation:   transferPortalProfileDto.TeamAbbreviation,
		TotalPoints:        0,
		CurrentWeeksPoints: 0,
		SpendingCount:      0,
		RemovedFromBoard:   false,
	}

	db.Create(&newProfileForRecruit)

	return newProfileForRecruit
}

func RemovePlayerFromTransferPortalBoard(id string) {
	db := dbprovider.GetInstance().GetDB()

	profile := GetOnlyTransferPortalProfileByID(id)

	profile.Deactivate()
	pid := profile.PromiseID.Int64
	profile.RemovePromise()
	repository.SaveTransferPortalProfile(profile, db)
	if pid > 0 && !profile.IsSigned {
		promiseID := strconv.Itoa(int(pid))
		promise := GetCollegePromiseByID(promiseID)
		promise.Deactivate()
		repository.DeleteCollegePromise(promise, db)
	}
}

func AllocatePointsToTransferPlayer(updateTransferPortalBoardDto structs.UpdateTransferPortalBoard) {
	db := dbprovider.GetInstance().GetDB()

	var teamId = strconv.Itoa(updateTransferPortalBoardDto.TeamID)
	var profile = GetOnlyRecruitingProfileByTeamID(teamId)
	var portalProfiles = GetOnlyTransferPortalProfilesByTeamID(teamId)
	var updatedPlayers = updateTransferPortalBoardDto.Players

	currentPoints := 0.0

	for i := 0; i < len(portalProfiles); i++ {
		updatedRecruit := GetPlayerFromTransferPortalList(int(portalProfiles[i].CollegePlayerID), updatedPlayers)

		if portalProfiles[i].CurrentWeeksPoints != updatedRecruit.CurrentWeeksPoints {

			// Allocate Points to Profile
			currentPoints += float64(updatedRecruit.CurrentWeeksPoints)
			profile.AllocateSpentPoints(currentPoints)
			// If total not surpassed, allocate to the recruit and continue
			if profile.SpentPoints <= profile.WeeklyPoints {
				portalProfiles[i].AllocatePoints(updatedRecruit.CurrentWeeksPoints)
				fmt.Println("Saving recruit " + strconv.Itoa(int(portalProfiles[i].CollegePlayerID)))
			} else {
				panic("Error: Allocated more points for Profile " + strconv.Itoa(int(profile.TeamID)) + " than what is allowed.")
			}
			db.Save(&portalProfiles[i])
		} else {
			currentPoints += float64(portalProfiles[i].CurrentWeeksPoints)
			profile.AllocateSpentPoints(currentPoints)
		}
	}

	// Save profile
	repository.SaveRecruitingTeamProfile(profile, db)
}

func AICoachFillBoardsPhase() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID := strconv.Itoa(ts.CollegeSeasonID)
	AITeams := GetOnlyAITeamRecruitingProfiles()
	// Shuffles the list of AI teams so that it's not always iterating from A-Z. Gives the teams at the lower end of the list a chance to recruit other croots
	rand.Shuffle(len(AITeams), func(i, j int) {
		AITeams[i], AITeams[j] = AITeams[j], AITeams[i]
	})
	transferPortalPlayers := GetTransferPortalPlayers()
	teamMap := GetCollegeTeamMap()
	standingsMap := GetCollegeStandingsMap(seasonID)
	profiles := []structs.TransferPortalProfile{}

	for idx, teamProfile := range AITeams {
		if teamProfile.IsUserTeam {
			continue
		}
		fmt.Println("Iterating "+teamProfile.TeamAbbreviation+" on IDX: ", idx)
		team := teamMap[teamProfile.ID]
		teamStandings := standingsMap[uint(teamProfile.TeamID)]
		teamID := strconv.Itoa(int(teamProfile.ID))
		portalProfileMap := getTransferPortalProfileMapByTeamID(teamID)
		if len(portalProfileMap) > 100 {
			continue
		}
		roster := GetAllCollegePlayersByTeamId(teamID)
		rosterSize := len(roster)
		teamCap := 105
		if !teamProfile.IsFBS {
			teamCap = 80
		}
		if rosterSize >= teamCap {
			continue
		}

		majorNeedsMap := getMajorNeedsMap()

		for _, r := range roster {
			if (team.IsFBS && r.Overall > 42 && majorNeedsMap[r.Position]) ||
				(!team.IsFBS && r.Overall > 34 && majorNeedsMap[r.Position]) {
				majorNeedsMap[r.Position] = false
			}
		}
		profileCount := len(portalProfileMap)

		for _, tp := range transferPortalPlayers {
			if profileCount >= 100 {
				break
			}
			isBadFit := IsBadSchemeFit(teamProfile.OffensiveScheme, teamProfile.DefensiveScheme, tp.Archetype, tp.Position)
			if isBadFit || !majorNeedsMap[tp.Position] || tp.PreviousTeamID == team.ID || portalProfileMap[tp.ID].CollegePlayerID == tp.ID || portalProfileMap[tp.ID].ID > 0 {
				continue
			}

			// Put together a player prestige rating to use as a qualifier on which teams will target specific players. Ideally more experienced coaches will be able to target higher rated players
			// playerPrestige := getPlayerPrestigeRating(tp.Stars, tp.Overall)
			// if coach.Prestige < playerPrestige {
			// 	continue
			// }
			bias := tp.RecruitingBias
			biasMod := 0
			postSeasonStatus := teamStandings.PostSeasonStatus
			if bias == richHistory {
				// Get multiple season standings
				teamHistory := GetStandingsHistoryByTeamID(teamID)
				averageWins := getAverageWins(teamHistory)
				biasMod += averageWins
				if teamProfile.AIQuality == "Blue Blood" {
					biasMod += 20
				}
			} else if bias == nationalChampionshipContender {
				if postSeasonStatus == "Bowl Game" {
					biasMod += 10
				} else if postSeasonStatus == "Playoffs" {
					biasMod += 15
				} else if postSeasonStatus == "National Championship Participant" {
					biasMod += 20
				} else if postSeasonStatus == "National Champions" {
					biasMod += 25
				}
			} else if bias == upcomingTeam {
				biasMod += (teamStandings.TotalWins * 2)
				if teamProfile.AIQuality == "Playoff Buster" {
					biasMod += 15
				}
			} else if bias == differentState && tp.State != team.State {
				biasMod += 15
			} else if bias == specificCoach && tp.LegacyID == team.ID {
				biasMod += 25
			} else if bias == legacy && tp.LegacyID == team.ID {
				biasMod += 25
			} else {
				biasMod = 5
			}

			if !teamProfile.IsFBS && teamProfile.ID > 194 {
				biasMod += 20
			}

			diceRoll := util.GenerateIntFromRange(1, 50)
			if diceRoll < biasMod {
				profileCount++
				portalProfile := structs.TransferPortalProfile{
					ProfileID:        teamProfile.ID,
					CollegePlayerID:  tp.ID,
					SeasonID:         uint(ts.CollegeSeasonID),
					TeamAbbreviation: teamProfile.TeamAbbreviation,
				}
				profiles = append(profiles, portalProfile)
			}
		}

	}

	repository.CreateTransferPortalProfileRecordsBatch(db, profiles, 500)
}

func AICoachAllocateAndPromisePhase() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	AITeams := GetOnlyAITeamRecruitingProfiles()
	portalProfiles := GetAllTransferPortalProfiles()
	profileMapByTeamID := MakePortalProfileMapByTeamID(portalProfiles)
	profileMapByPlayerID := MakePortalProfileMapByPlayerID(portalProfiles)
	collegePlayers := GetAllCollegePlayers()
	transferPortalPlayerMap := MakeCollegePlayerMap(collegePlayers)
	collegePlayerMap := MakeCollegePlayerMapByTeamID(collegePlayers, true)
	// regionMap := util.GetRegionMap()
	// Shuffles the list of AI teams so that it's not always iterating from A-Z. Gives the teams at the lower end of the list a chance to recruit other croots
	rand.Shuffle(len(AITeams), func(i, j int) {
		AITeams[i], AITeams[j] = AITeams[j], AITeams[i]
	})

	for _, teamProfile := range AITeams {
		if teamProfile.IsUserTeam {
			continue
		}

		roster := collegePlayerMap[teamProfile.ID]
		teamCap := 105
		if !teamProfile.IsFBS {
			teamCap = 80
		}
		if len(roster) >= teamCap {
			continue
		}

		majorNeedsMap := getMajorNeedsMap()

		for _, r := range roster {
			if r.Overall > 42 && majorNeedsMap[r.Position] {
				majorNeedsMap[r.Position] = false
			}
		}

		teamProfile.ResetSpentPoints()
		points := 0.0

		portalProfiles := profileMapByTeamID[teamProfile.ID]
		for _, profile := range portalProfiles {
			if points >= teamProfile.WeeklyPoints {
				break
			}
			if profile.CurrentWeeksPoints > 0 || profile.RemovedFromBoard {
				continue
			}
			tp := transferPortalPlayerMap[profile.CollegePlayerID]
			// If player has already signed or if the position has been fulfilled
			isBadFit := IsBadSchemeFit(teamProfile.OffensiveScheme, teamProfile.DefensiveScheme, tp.Archetype, tp.Position)
			if isBadFit || tp.TeamID > 0 || tp.TransferStatus == 0 || tp.ID == 0 || !majorNeedsMap[tp.Position] {
				profile.Deactivate()
				repository.SaveTransferPortalProfile(profile, db)
				continue
			}
			pointsRemaining := teamProfile.WeeklyPoints - teamProfile.SpentPoints
			if teamProfile.SpentPoints >= teamProfile.WeeklyPoints || pointsRemaining <= 0 || (pointsRemaining < 1 && pointsRemaining > 0) {
				break
			}

			removePlayerFromBoard := false
			num := 0.0

			profiles := profileMapByPlayerID[profile.CollegePlayerID]
			leadingTeamVal := IsAITeamContendingForPortalPlayer(profiles)
			if profile.CurrentWeeksPoints > 0 && profile.TotalPoints+float64(profile.CurrentWeeksPoints) >= float64(leadingTeamVal)*0.66 {
				// continue, leave everything alone
				points += float64(profile.CurrentWeeksPoints)
				continue
			} else if profile.CurrentWeeksPoints > 0 && profile.TotalPoints+float64(profile.CurrentWeeksPoints) < float64(leadingTeamVal)*0.66 {
				profile.Deactivate()
				repository.SaveTransferPortalProfile(profile, db)
				continue
			}

			maxChance := 2
			if ts.CollegeWeek > 3 {
				maxChance = 4
			}
			chance := util.GenerateIntFromRange(1, maxChance)
			if (chance < 2 && ts.TransferPortalPhase <= 3) || (chance < 4 && ts.TransferPortalPhase > 3) {
				continue
			}

			min := teamProfile.AIMaxThreshold
			max := teamProfile.AIMinThreshold
			if max > 10 {
				max = 10
			}
			if !teamProfile.IsUserTeam && !teamProfile.IsFBS && max > 5 {
				max = 5
			}
			if min >= max {
				min = 4
			}
			num = util.GenerateFloatFromRange(float64(min), float64(max))
			if num > pointsRemaining {
				num = pointsRemaining
			}

			if num+profile.TotalPoints < float64(leadingTeamVal)*0.66 {
				removePlayerFromBoard = true
			}
			if leadingTeamVal < 8 {
				removePlayerFromBoard = false
			}

			if removePlayerFromBoard {
				profile.Deactivate()
				repository.SaveTransferPortalProfile(profile, db)
				continue
			}
			profile.AllocatePoints(int(num))
			points += num

			// Generate Promise based on coach bias
			if profile.PromiseID.Int64 == 0 && !profile.RolledOnPromise {
				promiseOdds := getBasePromiseOdds("", "")
				diceRoll := util.GenerateIntFromRange(1, 100)

				if diceRoll < promiseOdds {
					// Commit Promise
					promiseLevel := getPromiseLevel("")
					promiseWeight := "Medium"
					promiseType := ""
					benchmarkStr := ""
					promiseBenchmark := 0

					bias := tp.RecruitingBias
					// Rewrite
					if bias == closeToHome && (teamProfile.State == tp.State) {
						promiseType = "Home State Game"
						benchmarkStr = tp.State
						promiseWeight = "Low"
					} else if bias == immediateStart && tp.Overall > 40 {
						promiseType = "Snaps"
						// Rewrite
						if promiseLevel == 1 {
							promiseBenchmark += 5
							if promiseBenchmark > tp.Stamina {
								promiseBenchmark = tp.Stamina - 1
							}
						} else if promiseLevel == -1 {
							promiseBenchmark -= 1
						}

						promiseWeight = getPromiseWeightBySnapsOrWins(tp.Position, "Snap Count", promiseBenchmark)
					} else if bias == nationalChampionshipContender || bias == richHistory {
						// Promise based on wins
						promiseBenchmark = 6
						promiseType = "Wins"
						if promiseLevel == 1 {
							promiseBenchmark += 3
						} else if promiseLevel == -1 {
							promiseBenchmark -= 3
						}
						promiseWeight = getPromiseWeightBySnapsOrWins(tp.Position, "Snap Count", promiseBenchmark)
					} else if bias == legacy && tp.LegacyID == uint(teamProfile.TeamID) {
						promiseType = "Legacy"
						promiseWeight = "Medium"
					} else if bias == specificCoach && tp.LegacyID == teamProfile.ID {
						promiseType = "Specific Coach"
						promiseWeight = "Low"
					} else if bias == differentState && teamProfile.State != tp.State {
						promiseType = "Different State"
						promiseWeight = "Low"
					}

					if promiseType != "" {
						collegePromise := structs.CollegePromise{
							TeamID:          uint(teamProfile.TeamID),
							CollegePlayerID: tp.ID,
							PromiseType:     promiseType,
							PromiseWeight:   promiseWeight,
							Benchmark:       promiseBenchmark,
							BenchmarkStr:    benchmarkStr,
							IsActive:        true,
						}
						repository.CreateCollegePromiseRecord(collegePromise, db)
					}
				}

				profile.ToggleRolledOnPromise()
			}
			// Save Profile
			if profile.CurrentWeeksPoints > 0 {
				repository.SaveTransferPortalProfile(profile, db)
			}
		}
		teamProfile.AIAllocateSpentPoints(points)
		repository.SaveRecruitingTeamProfile(teamProfile, db)
	}
}

func SyncTransferPortal() {
	db := dbprovider.GetInstance().GetDB()
	//GetCurrentWeek
	ts := GetTimestamp()
	// Use IsRecruitingLocked to lock the TP when not in use
	teamProfileMap := GetTeamProfileMap()
	transferPortalPlayers := GetTransferPortalPlayers()
	transferPortalProfiles := GetAllTransferPortalProfiles()
	transferPortalProfileMap := MakePortalProfileMapByPlayerID(transferPortalProfiles)
	rosterMap := GetFullTeamRosterWithCrootsMap()
	collegePromises := GetAllCollegePromises()
	collegePromiseMap := MakeCollegePromiseMap(collegePromises)

	if !ts.IsRecruitingLocked {
		ts.ToggleLockRecruiting()
		repository.SaveTimestamp(ts, db)
	}

	for _, portalPlayer := range transferPortalPlayers {
		// Skip over players that have already transferred
		if portalPlayer.TransferStatus != 2 || portalPlayer.TeamID > 0 {
			continue
		}

		portalProfiles := transferPortalProfileMap[portalPlayer.ID]
		if len(portalProfiles) == 0 && ts.TransferPortalRound < constants.FinalPortalRound {
			continue
		}

		// If no one has a profile on them during round 10
		if len(portalProfiles) == 0 && ts.TransferPortalRound == 10 && len(portalPlayer.TransferLikeliness) > 0 {
			roster := rosterMap[portalPlayer.PreviousTeamID]
			tp := teamProfileMap[strconv.Itoa(int(portalPlayer.PreviousTeamID))]
			if (len(roster) > constants.MaxFBSRosterSize && tp.IsFBS) || (len(roster) > constants.MaxFCSRosterSize && !tp.IsFBS) {
				continue
			}
			rosterMap[portalPlayer.PreviousTeamID] = append(rosterMap[portalPlayer.PreviousTeamID], portalPlayer)
			portalPlayer.WillReturn()
			repository.SaveCFBPlayer(portalPlayer, db)
			continue
		}

		totalPointsOnPlayer := 0.0
		eligiblePointThreshold := 0.0
		readyToSign := false
		minSpendingCount := 100
		maxSpendingCount := 0
		signingMinimum := constants.PortalSigningMinimum
		teamCount := 0
		eligibleTeams := []structs.TransferPortalProfile{}

		for i := range portalProfiles {
			if portalProfiles[i].RemovedFromBoard {
				continue
			}
			promise := collegePromiseMap[uint(portalProfiles[i].PromiseID.Int64)]

			multiplier := getMultiplier(promise)
			portalProfiles[i].AddPointsToTotal(multiplier)
		}

		sort.Slice(portalProfiles, func(i, j int) bool {
			return portalProfiles[i].TotalPoints > portalProfiles[j].TotalPoints
		})

		for i := range portalProfiles {
			// roster := rosterMap[portalProfiles[i].ProfileID]
			// tp := teamProfileMap[strconv.Itoa(int(portalProfiles[i].ProfileID))]
			// // if (len(roster) > 105 && tp.IsFBS) || (len(roster) > 80 && !tp.IsFBS) {
			// // 	continue
			// // }
			if eligiblePointThreshold == 0.0 {
				eligiblePointThreshold = portalProfiles[i].TotalPoints * signingMinimum
			}
			if portalProfiles[i].TotalPoints >= eligiblePointThreshold {
				if portalProfiles[i].SpendingCount < minSpendingCount {
					minSpendingCount = portalProfiles[i].SpendingCount
				}
				if portalProfiles[i].SpendingCount > maxSpendingCount {
					maxSpendingCount = portalProfiles[i].SpendingCount
				}
				eligibleTeams = append(eligibleTeams, portalProfiles[i])
				totalPointsOnPlayer += portalProfiles[i].TotalPoints
				teamCount += 1
			}

		}

		if (teamCount >= 1 && minSpendingCount >= 2) || (teamCount > 1 && minSpendingCount > 3) || (ts.TransferPortalRound >= constants.FinalPortalRound) && totalPointsOnPlayer > 0 {
			// threshold met
			readyToSign = true
		}
		var winningTeamID uint = 0
		var odds float64 = 0
		if readyToSign {
			for winningTeamID == 0 && len(eligibleTeams) > 0 {
				percentageOdds := rand.Float64() * (totalPointsOnPlayer)
				currentProbability := 0.0
				for _, profile := range eligibleTeams {
					currentProbability += profile.TotalPoints
					if percentageOdds <= currentProbability {
						// WINNING TEAM
						winningTeamID = profile.ProfileID
						odds = profile.TotalPoints / totalPointsOnPlayer * 100
						break
					}
				}

				if winningTeamID > 0 {
					winningTeamIDSTR := strconv.Itoa(int(winningTeamID))

					teamProfile := teamProfileMap[winningTeamIDSTR]
					currentRoster := rosterMap[teamProfile.ID]
					teamCap := constants.MaxFBSRosterSize
					if !teamProfile.IsFBS {
						teamCap = constants.MaxFCSRosterSize
					}
					if len(currentRoster) < teamCap {
						promise := GetCollegePromiseByCollegePlayerID(strconv.Itoa(int(portalPlayer.ID)), winningTeamIDSTR)
						if promise.ID > 0 {
							promise.MakePromise()
							repository.SaveCollegePromiseRecord(promise, db)
						}
						portalPlayer.SignWithNewTeam(teamProfile.TeamID, teamProfile.TeamAbbreviation)
						message := portalPlayer.FirstName + " " + portalPlayer.LastName + ", " + strconv.Itoa(portalPlayer.Stars) + " star " + portalPlayer.Position + " from " + portalPlayer.PreviousTeam + " has signed with " + portalPlayer.TeamAbbr + " with " + strconv.Itoa(int(odds)) + " percent odds."
						CreateNewsLog("CFB", message, "Transfer Portal", int(winningTeamID), ts)
						fmt.Println("Created new log!")
						// Add player to existing roster map
						rosterMap[teamProfile.ID] = append(rosterMap[teamProfile.ID], portalPlayer)
						for i := range portalProfiles {
							if portalProfiles[i].ID == winningTeamID {
								portalProfiles[i].SignPlayer()
								break
							}
						}

					} else {
						// Filter out profile
						eligibleTeams = FilterOutPortalProfile(eligibleTeams, winningTeamID)
						winningTeamID = 0
						if len(eligibleTeams) == 0 {
							break
						}

						totalPointsOnPlayer = 0
						for _, p := range eligibleTeams {
							totalPointsOnPlayer += p.TotalPoints
						}
					}

				}
			}

		}
		for _, p := range portalProfiles {
			if winningTeamID > 0 && p.ID != winningTeamID {
				p.RemovePromise()
				p.Lock()
			}
			if winningTeamID > 0 || p.SpendingCount > 0 {
				repository.SaveTransferPortalProfile(p, db)
			}
			fmt.Println("Save transfer portal profile from " + portalPlayer.TeamAbbr + " towards " + portalPlayer.FirstName + " " + portalPlayer.LastName)
			if winningTeamID > 0 && p.ProfileID != winningTeamID {
				promise := GetCollegePromiseByCollegePlayerID(strconv.Itoa(int(portalPlayer.ID)), strconv.Itoa(int(p.ProfileID)))
				if promise.ID > 0 {
					repository.DeleteCollegePromise(promise, db)
				}
			}
		}
		// Save Recruit
		if portalPlayer.TeamID > 0 {
			repository.SaveCollegePlayerRecord(portalPlayer, db)
		}
	}

	ts.IncrementTransferPortalRound()
	repository.SaveTimestamp(ts, db)
}

func GetPromisesByTeamID(teamID string) []structs.CollegePromise {
	db := dbprovider.GetInstance().GetDB()

	var promises []structs.CollegePromise

	db.Where("team_id = ?", teamID).Find(&promises)

	return promises
}

func GetOnlyTransferPortalProfilesByTeamID(teamID string) []structs.TransferPortalProfile {
	db := dbprovider.GetInstance().GetDB()

	var profiles []structs.TransferPortalProfile

	db.Where("profile_id = ?", teamID).Find(&profiles)

	return profiles
}

func GetTransferPortalProfilesByPlayerID(playerID string) []structs.TransferPortalProfile {
	db := dbprovider.GetInstance().GetDB()

	var profiles []structs.TransferPortalProfile

	db.Where("college_player_id = ?", playerID).Find(&profiles)

	return profiles
}

func GetTransferPortalProfilesForPage(teamID string) []structs.TransferPortalProfileResponse {
	var response []structs.TransferPortalProfileResponse
	allProfiles := GetAllTransferPortalProfiles()
	teamProfiles := MakePortalProfileMapByTeamID(allProfiles)
	profiles := teamProfiles[uint(util.ConvertStringToInt(teamID))]

	portalProfileMapByPlayerID := MakePortalProfileMapByPlayerID(allProfiles)
	collegePlayers := GetAllCollegePlayers()
	collegePlayerMap := MakeCollegePlayerMap(collegePlayers)

	promises := GetAllCollegePromises()
	promisesTeamMap := MakePromiseMapByTeamID(promises)

	for _, p := range profiles {
		if p.RemovedFromBoard {
			continue
		}
		cp := collegePlayerMap[p.CollegePlayerID]
		cpResponse := structs.TransferPlayerResponse{}
		ovr := util.GetOverallGrade(cp.Overall, cp.Year)
		teamPromises := promisesTeamMap[p.ProfileID]
		promiseMapByPlayerID := MakePromiseMapByPlayerIDByTeam(teamPromises)
		playerPromise := promiseMapByPlayerID[p.CollegePlayerID]
		cpResponse.Map(cp, ovr, portalProfileMapByPlayerID[cp.ID])

		pResponse := structs.TransferPortalProfileResponse{
			ID:                    p.ID,
			SeasonID:              p.SeasonID,
			CollegePlayerID:       p.CollegePlayerID,
			ProfileID:             p.ProfileID,
			PromiseID:             uint(p.PromiseID.Int64),
			TeamAbbreviation:      p.TeamAbbreviation,
			TotalPoints:           p.TotalPoints,
			CurrentWeeksPoints:    p.CurrentWeeksPoints,
			PreviouslySpentPoints: p.PreviouslySpentPoints,
			SpendingCount:         p.SpendingCount,
			RemovedFromBoard:      p.RemovedFromBoard,
			RolledOnPromise:       p.RolledOnPromise,
			LockProfile:           p.LockProfile,
			IsSigned:              p.IsSigned,
			Recruiter:             p.Recruiter,
			CollegePlayer:         cpResponse,
			Promise:               playerPromise,
		}

		response = append(response, pResponse)

	}

	return response
}

func GetTransferPortalProfilesByTeamID(teamID string) []structs.TransferPortalProfile {
	db := dbprovider.GetInstance().GetDB()

	var profiles []structs.TransferPortalProfile

	db.Preload("CollegePlayer.Profiles").Where("profile_id = ?", teamID).Find(&profiles)

	return profiles
}

func GetAllTransferPortalProfiles() []structs.TransferPortalProfile {
	db := dbprovider.GetInstance().GetDB()

	var profiles []structs.TransferPortalProfile

	db.Find(&profiles)

	return profiles
}

func GetOnlyTransferPortalProfileByID(tppID string) structs.TransferPortalProfile {
	db := dbprovider.GetInstance().GetDB()

	var profile structs.TransferPortalProfile

	db.Where("id = ?", tppID).Find(&profile)

	return profile
}

func GetOnlyTransferPortalProfileByPlayerID(playerId, teamId string) structs.TransferPortalProfile {
	db := dbprovider.GetInstance().GetDB()

	var profile structs.TransferPortalProfile

	db.Where("college_player_id = ? and profile_id = ?", playerId, teamId).Find(&profile)

	return profile
}

func GetTransferPortalData(teamID string) structs.TransferPortalResponse {
	var waitgroup sync.WaitGroup
	waitgroup.Add(5)
	profileChan := make(chan structs.RecruitingTeamProfile)
	playersChan := make(chan []structs.TransferPlayerResponse)
	boardChan := make(chan []structs.TransferPortalProfileResponse)
	promiseChan := make(chan []structs.CollegePromise)
	teamsChan := make(chan []structs.CollegeTeam)

	go func() {
		waitgroup.Wait()
		close(profileChan)
		close(playersChan)
		close(boardChan)
		close(promiseChan)
		close(teamsChan)
	}()

	go func() {
		defer waitgroup.Done()
		profile := GetOnlyRecruitingProfileByTeamID(teamID)
		profileChan <- profile
	}()
	go func() {
		defer waitgroup.Done()
		tpPlayers := GetTransferPortalPlayersForPage()
		playersChan <- tpPlayers
	}()
	go func() {
		defer waitgroup.Done()
		tpProfiles := GetTransferPortalProfilesForPage(teamID)
		boardChan <- tpProfiles
	}()
	go func() {
		defer waitgroup.Done()
		cPromises := GetPromisesByTeamID(teamID)
		promiseChan <- cPromises
	}()
	go func() {
		defer waitgroup.Done()
		cTeams := GetAllCollegeTeams()
		teamsChan <- cTeams
	}()

	teamProfile := <-profileChan
	players := <-playersChan
	board := <-boardChan
	promises := <-promiseChan
	teams := <-teamsChan

	return structs.TransferPortalResponse{
		Team:         teamProfile,
		Players:      players,
		TeamBoard:    board,
		TeamPromises: promises,
		TeamList:     teams,
	}
}

func GetTransferPortalPlayersForPage() []structs.TransferPlayerResponse {
	db := dbprovider.GetInstance().GetDB()

	portalProfiles := GetAllTransferPortalProfiles()
	portalProfileMap := MakePortalProfileMapByPlayerID(portalProfiles)

	var players []structs.CollegePlayer

	db.Where("transfer_status = 2").Order("overall DESC").Find(&players)

	playerList := []structs.TransferPlayerResponse{}

	for _, p := range players {
		res := structs.TransferPlayerResponse{}
		ovr := util.GetOverallGrade(p.Overall, p.Year)
		res.Map(p, ovr, portalProfileMap[p.ID])

		playerList = append(playerList, res)
	}

	return playerList
}

func GetNFLPlayerSeasonSnapMap(seasonID string) map[uint]structs.NFLPlayerSeasonSnaps {
	seasonStatMap := make(map[uint]structs.NFLPlayerSeasonSnaps)

	seasonStats := GetNFLSeasonSnapsBySeason(seasonID)
	for _, stat := range seasonStats {
		seasonStatMap[stat.PlayerID] = stat
	}

	return seasonStatMap
}

func GetCollegePlayerSeasonSnapMap(seasonID string) map[uint]structs.CollegePlayerSeasonSnaps {
	seasonStatMap := make(map[uint]structs.CollegePlayerSeasonSnaps)

	seasonStats := GetCollegeSeasonSnapsBySeason(seasonID)
	for _, stat := range seasonStats {
		seasonStatMap[stat.PlayerID] = stat
	}

	return seasonStatMap
}

func GetCollegePlayerSeasonStatsMap(seasonID, gameType string) map[uint]structs.CollegePlayerSeasonStats {
	seasonStatMap := make(map[uint]structs.CollegePlayerSeasonStats)

	seasonStats := GetCollegePlayerSeasonStatsBySeason(seasonID, gameType)
	for _, stat := range seasonStats {
		seasonStatMap[stat.CollegePlayerID] = stat
	}

	return seasonStatMap
}

func GetCollegePlayerStatsMap(seasonID string) map[uint][]structs.CollegePlayerStats {
	seasonStatMap := make(map[uint][]structs.CollegePlayerStats)

	seasonStats := GetCollegeAllPlayerStatsBySeason(seasonID)
	for _, stat := range seasonStats {
		stats := seasonStatMap[uint(stat.CollegePlayerID)]
		if len(stats) == 0 {
			seasonStatMap[uint(stat.CollegePlayerID)] = []structs.CollegePlayerStats{}
		}
		seasonStatMap[uint(stat.CollegePlayerID)] = append(seasonStatMap[uint(stat.CollegePlayerID)], stat)
	}

	return seasonStatMap
}

func GetNFLPlayerStatsMap(seasonID string) map[uint][]structs.NFLPlayerStats {
	seasonStatMap := make(map[uint][]structs.NFLPlayerStats)

	seasonStats := GetAllNFLPlayerStatsBySeason(seasonID)
	for _, stat := range seasonStats {
		stats := seasonStatMap[uint(stat.NFLPlayerID)]
		if len(stats) == 0 {
			seasonStatMap[uint(stat.NFLPlayerID)] = []structs.NFLPlayerStats{}
		}
		seasonStatMap[uint(stat.NFLPlayerID)] = append(seasonStatMap[uint(stat.NFLPlayerID)], stat)
	}

	return seasonStatMap
}

func GetFullTeamRosterWithCrootsMap() map[uint][]structs.CollegePlayer {
	m := &sync.Mutex{}
	var wg sync.WaitGroup
	collegeTeams := GetAllCollegeTeams()
	fullMap := make(map[uint][]structs.CollegePlayer)
	wg.Add(len(collegeTeams))
	semaphore := make(chan struct{}, 10)
	for _, team := range collegeTeams {
		semaphore <- struct{}{}
		go func(t structs.CollegeTeam) {
			defer wg.Done()
			id := strconv.Itoa(int(t.ID))
			collegePlayers := GetAllCollegePlayersByTeamId(id)
			croots := GetSignedRecruitsByTeamProfileID(id)
			fullList := collegePlayers
			for _, croot := range croots {
				p := structs.CollegePlayer{}
				p.MapFromRecruit(croot, t)

				fullList = append(fullList, p)
			}

			m.Lock()
			fullMap[t.ID] = fullList
			m.Unlock()
			<-semaphore
		}(team)
	}

	wg.Wait()
	close(semaphore)
	return fullMap
}

func filterRosterByPosition(roster []structs.CollegePlayer, pos string) []structs.CollegePlayer {
	estimatedSize := len(roster) / 5 // Adjust this based on your data
	filteredList := make([]structs.CollegePlayer, 0, estimatedSize)
	for _, p := range roster {
		if p.Position != pos || (p.Year == 5 || (p.Year == 4 && p.IsRedshirt)) {
			continue
		}
		filteredList = append(filteredList, p)
	}
	sort.Slice(filteredList, func(i, j int) bool {
		return filteredList[i].Overall > filteredList[j].Overall
	})

	return filteredList
}

func getTransferPortalProfileMap(players []structs.CollegePlayer) map[uint][]structs.TransferPortalProfile {
	portalMap := make(map[uint][]structs.TransferPortalProfile)
	var mu sync.Mutex     // to safely update the map
	var wg sync.WaitGroup // to wait for all goroutines to finish
	semaphore := make(chan struct{}, 10)
	for _, p := range players {
		semaphore <- struct{}{}
		wg.Add(1)
		go func(c structs.CollegePlayer) {
			defer wg.Done()
			playerID := strconv.Itoa(int(c.ID))
			portalProfiles := GetTransferPortalProfilesByPlayerID(playerID)
			mu.Lock()
			portalMap[c.ID] = portalProfiles
			mu.Unlock()

			<-semaphore
		}(p)
	}
	wg.Wait()
	close(semaphore)
	return portalMap
}

// GetTransferFloor -- Get the Base Floor to determine if a player will transfer or not based on a promise
func getTransferFloor(likeliness string) int {
	min := 25
	max := 100
	if likeliness == "Low" {
		max = 40
	} else if likeliness == "Medium" {
		min = 45
		max = 70
	} else {
		min = 75
	}

	return util.GenerateIntFromRange(min, max)
}

// getPromiseFloor -- Get the modifier towards the floor value above
func getPromiseFloor(weight string) int {
	if weight == "Very Low" {
		return 10
	}
	if weight == "Low" {
		return 20
	}
	if weight == "Medium" {
		return 40
	}
	if weight == "High" {
		return 60
	}
	return util.GenerateIntFromRange(70, 80)
}

func getPromiseWeightBySnapsOrWins(position, category string, benchmark int) string {
	if benchmark == 0 {
		return "Very Low"
	}
	weight := "Medium"
	if category == "Wins" {
		if benchmark <= 4 {
			weight = "Very Low"
		}
		if benchmark <= 6 {
			weight = "Low"
		}
		if benchmark <= 10 {
			weight = "Medium"
		}
		if benchmark <= 12 {
			weight = "High"
		}
		if benchmark <= 15 {
			weight = "Very High"
		}

	}
	if category == "Snap Count" {
		if position == "P" || position == "K" {
			if benchmark <= 5 {
				weight = "Low"
			}
			if benchmark <= 8 {
				weight = "Medium"
			}
			if benchmark <= 10 {
				weight = "High"
			}
			if benchmark <= 20 {
				weight = "Very High"
			}
		} else {
			if benchmark <= 10 {
				weight = "Very Low"
			}
			if benchmark <= 20 {
				weight = "Low"
			}
			if benchmark <= 30 {
				weight = "Medium"
			}
			if benchmark <= 50 {
				weight = "High"
			}
		}
	}
	return weight
}

func getPlayerPrestigeRating(stars, overall int) int {
	prestige := 1

	starMod := stars / 2
	if starMod <= 0 {
		starMod = 1
	}

	overallMod := overall / 10
	if overallMod <= 0 {
		overallMod = 1
	}

	return prestige + starMod + overallMod
}

func getAverageWins(standings []structs.CollegeStandings) int {
	wins := 0
	for _, s := range standings {
		wins += s.TotalWins
	}

	totalStandings := len(standings)
	if totalStandings > 0 {
		wins = wins / len(standings)
	}

	return wins
}

func getBasePromiseOdds(tbPref, ptTendency string) int {
	promiseOdds := 50
	if tbPref == "Recruiting" {
		promiseOdds += 20
	} else if ptTendency == "Transfer" {
		promiseOdds -= 20
	}

	return promiseOdds
}

func getTransferStatus(weight int) string {
	if weight < 20 {
		return "Low"
	}
	if weight < 40 {
		return "Medium"
	}
	return "High"
}

func getPromiseLevel(pt string) int {
	promiseLevel := 0
	if pt == "Over-Promise" {
		promiseLevel = 1
	} else if pt == "Under-Promise" {
		promiseLevel = -1
	}
	return promiseLevel
}

func getMultiplier(pr structs.CollegePromise) float64 {
	if pr.ID == 0 || !pr.IsActive {
		return 1
	}
	weight := pr.PromiseWeight
	if weight == "Very Low" {
		return 1.05
	} else if weight == "Low" {
		return 1.1
	} else if weight == "Medium" {
		return 1.3
	} else if weight == "High" {
		return 1.5
	}
	// Very High
	return 1.75
}

func GetPlayerFromTransferPortalList(id int, profiles []structs.TransferPortalProfileResponse) structs.TransferPortalProfileResponse {
	var profile structs.TransferPortalProfileResponse

	for i := 0; i < len(profiles); i++ {
		if profiles[i].CollegePlayerID == uint(id) {
			profile = profiles[i]
			break
		}
	}

	return profile
}

func GetTransferScoutingDataByPlayerID(id string) models.ScoutingDataResponse {
	ts := GetTimestamp()

	seasonID := ts.CollegeSeasonID
	seasonIDSTR := strconv.Itoa(int(seasonID))

	draftee := GetCollegePlayerByCollegePlayerId(id)

	seasonStats := GetCollegePlayerSeasonStatsByPlayerIDAndSeason(id, seasonIDSTR, "2")
	teamID := strconv.Itoa(int(draftee.PreviousTeamID))
	collegeStandings := GetCollegeStandingsRecordByTeamID(teamID, seasonIDSTR)

	return models.ScoutingDataResponse{
		DrafteeSeasonStats: seasonStats,
		TeamStandings:      collegeStandings,
	}
}

func getTransferPortalProfileMapByTeamID(id string) map[uint]structs.TransferPortalProfile {
	profiles := GetOnlyTransferPortalProfilesByTeamID(id)

	profileMap := make(map[uint]structs.TransferPortalProfile)

	for _, profile := range profiles {
		profileMap[profile.CollegePlayerID] = profile
	}

	return profileMap
}

func GetTeamProfileMap() map[string]*structs.RecruitingTeamProfile {
	var mu sync.Mutex

	teamRecruitingProfiles := GetRecruitingProfileForRecruitSync()

	teamMap := make(map[string]*structs.RecruitingTeamProfile)
	for i := 0; i < len(teamRecruitingProfiles); i++ {
		mu.Lock()
		teamMap[strconv.Itoa(int(teamRecruitingProfiles[i].ID))] = &teamRecruitingProfiles[i]
		mu.Unlock()
	}

	return teamMap
}

func getSchemeMod(tp *structs.RecruitingTeamProfile, p structs.CollegePlayer, drop, gain float64) float64 {
	schemeMod := 0.0
	if tp.OffensiveScheme == "" {
		fmt.Println("PING!")
	}
	goodFit := IsGoodSchemeFit(tp.OffensiveScheme, tp.DefensiveScheme, p.Archetype, p.Position)
	badFit := IsBadSchemeFit(tp.OffensiveScheme, tp.DefensiveScheme, p.Archetype, p.Position)
	if goodFit {
		schemeMod += drop
	} else if badFit {
		schemeMod += gain
	}

	return schemeMod
}

func IsGoodSchemeFit(offensiveScheme, defensiveScheme, arch, position string) bool {
	archType := arch + " " + position
	offensiveSchemeList := GetFitsByScheme(offensiveScheme, false)
	defensiveSchemeList := GetFitsByScheme(defensiveScheme, false)
	totalFitList := append(offensiveSchemeList, defensiveSchemeList...)

	return CheckPlayerFits(archType, totalFitList)
}

func IsBadSchemeFit(offensiveScheme, defensiveScheme, arch, position string) bool {
	archType := arch + " " + position
	offensiveSchemeList := GetFitsByScheme(offensiveScheme, true)
	defensiveSchemeList := GetFitsByScheme(defensiveScheme, true)
	totalFitList := append(offensiveSchemeList, defensiveSchemeList...)

	return CheckPlayerFits(archType, totalFitList)
}

func getMajorNeedsMap() map[string]bool {
	majorNeedsMap := make(map[string]bool)

	if _, ok := majorNeedsMap["QB"]; !ok {
		majorNeedsMap["QB"] = true
	}

	if _, ok := majorNeedsMap["RB"]; !ok {
		majorNeedsMap["RB"] = true
	}

	if _, ok := majorNeedsMap["FB"]; !ok {
		majorNeedsMap["FB"] = true
	}

	if _, ok := majorNeedsMap["WR"]; !ok {
		majorNeedsMap["WR"] = true
	}

	if _, ok := majorNeedsMap["TE"]; !ok {
		majorNeedsMap["TE"] = true
	}

	if _, ok := majorNeedsMap["OT"]; !ok {
		majorNeedsMap["OT"] = true
	}

	if _, ok := majorNeedsMap["OG"]; !ok {
		majorNeedsMap["OG"] = true
	}

	if _, ok := majorNeedsMap["C"]; !ok {
		majorNeedsMap["C"] = true
	}

	if _, ok := majorNeedsMap["DT"]; !ok {
		majorNeedsMap["DT"] = true
	}

	if _, ok := majorNeedsMap["DE"]; !ok {
		majorNeedsMap["DE"] = true
	}

	if _, ok := majorNeedsMap["OLB"]; !ok {
		majorNeedsMap["OLB"] = true
	}

	if _, ok := majorNeedsMap["ILB"]; !ok {
		majorNeedsMap["ILB"] = true
	}

	if _, ok := majorNeedsMap["CB"]; !ok {
		majorNeedsMap["CB"] = true
	}

	if _, ok := majorNeedsMap["FS"]; !ok {
		majorNeedsMap["FS"] = true
	}

	if _, ok := majorNeedsMap["SS"]; !ok {
		majorNeedsMap["SS"] = true
	}

	if _, ok := majorNeedsMap["P"]; !ok {
		majorNeedsMap["P"] = true
	}

	if _, ok := majorNeedsMap["K"]; !ok {
		majorNeedsMap["K"] = true
	}

	if _, ok := majorNeedsMap["ATH"]; !ok {
		majorNeedsMap["ATH"] = true
	}

	return majorNeedsMap
}
