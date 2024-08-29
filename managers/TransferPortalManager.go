package managers

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"sync"

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
	// db := dbprovider.GetInstance().GetDB()
	w.Header().Set("Content-Disposition", "attachment;filename=transferStats.csv")
	w.Header().Set("Transfer-Encoding", "chunked")
	// Initialize writer
	writer := csv.NewWriter(w)
	ts := GetTimestamp()
	seasonID := strconv.Itoa(ts.CollegeSeasonID)
	allCollegePlayers := GetAllCollegePlayers()
	seasonSnapMap := GetCollegePlayerSeasonSnapMap(seasonID)
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
	bigGain := 25.0

	HeaderRow := []string{
		"Team", "First Name", "Last Name", "Stars",
		"Archetype", "Position", "Year", "Age", "Redshirt Status",
		"Overall", "Transfer Bias", "Transfer Status", "Transfer Weight", "Dice Roll",
	}

	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	for _, p := range allCollegePlayers {
		// Do not include redshirts and all graduating players
		if p.IsRedshirting || p.Year == 5 || (p.Year == 4 && !p.IsRedshirt) {
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
		totalSnaps := seasonStats.GetTotalSnaps()
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
			ageMod = .15
		} else if p.Year == 2 {
			ageMod = .4
		} else if p.Year == 3 {
			ageMod = .75
		} else if p.Year == 4 {
			ageMod = 1
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
					p.Position == "SS") && idx > 2 {
					depthChartCompetitionMod += 33
				}

				if (p.Position == "WR" ||
					p.Position == "OT" ||
					p.Position == "OG" ||
					p.Position == "DE" ||
					p.Position == "DT" ||
					p.Position == "OLB" ||
					p.Position == "ILB" ||
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
		schemeMod = getSchemeMod(teamProfile, p, smallDrop, smallGain)

		/// Not playing = 25, low depth chart = 16 or 33, scheme = 10, if you're all 3, that's a ~60% chance of transferring pre- modifiers
		transferWeight = starMod * ageMod * (snapMod + depthChartCompetitionMod + schemeMod)
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

		fmt.Println(strconv.Itoa(p.Year)+" YEAR "+p.TeamAbbr+" "+p.Position+" "+p.FirstName+" "+p.LastName+" HAS ANNOUNCED THEIR INTENTION TO TRANSFER | Weight: ", int(transferWeight))
		// db.Save(&p)
		csvModel := structs.MapPlayerToCSVModel(p)
		playerRow := []string{
			p.TeamAbbr, csvModel.FirstName, csvModel.LastName, strconv.Itoa(p.Stars),
			csvModel.Archetype, csvModel.Position,
			csvModel.Year, strconv.Itoa(p.Age), csvModel.RedshirtStatus,
			csvModel.OverallGrade, p.RecruitingBias, p.TransferLikeliness, strconv.Itoa(transferInt), strconv.Itoa(diceRoll),
		}

		err = writer.Write(playerRow)
		if err != nil {
			log.Fatal("Cannot write player row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
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
		if !team.IsAI {
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
				} else if bias == immediateStart && p.Overall > 48 {
					promiseType = "Snaps"
					// REWRITE
					promiseBenchmark = 0
					if promiseLevel == 1 {
						promiseBenchmark += 5
						if promiseBenchmark > p.Stamina {
							promiseBenchmark = p.Stamina - 1
						}
					} else if promiseLevel == -1 {
						promiseBenchmark -= 1
					}

					promiseWeight = getPromiseWeightByMinutesOrWins(promiseBenchmark)
				} else if bias == nationalChampionshipContender || bias == richHistory {
					// Promise based on wins
					promiseBenchmark = 20
					promiseType = "Wins"
					if promiseLevel == 1 {
						promiseBenchmark += 5
					} else if promiseLevel == -1 {
						promiseBenchmark -= 5
					}
					promiseWeight = getPromiseWeightByMinutesOrWins(promiseBenchmark)
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

				db.Create(&collegePromise)
			}
		}
	}
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
	id := strconv.Itoa(int(promise.ID))
	collegePlayerID := strconv.Itoa(int(promise.CollegePlayerID))
	profileID := strconv.Itoa(int(promise.TeamID))

	existingPromise := GetCollegePromiseByID(id)
	if existingPromise.ID != 0 && existingPromise.ID > 0 {
		existingPromise.Reactivate(promise.PromiseType, promise.PromiseWeight, promise.Benchmark)
		db.Save(&existingPromise)
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
				db.Save(&p)
				continue
			}
			// 1-100
			baseFloor := getTransferFloor(p.TransferLikeliness)
			// 10, 20, 40, 60, 70
			promiseModifier := getPromiseFloor(promise.PromiseWeight)
			difference := baseFloor - promiseModifier

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
				db.Delete(&promise)
				continue
			}

			// Create News Log
			message := "Breaking News! " + p.TeamAbbr + " " + strconv.Itoa(p.Stars) + " Star " + p.Position + " " + p.FirstName + " " + p.LastName + " has withdrawn their name from the transfer portal!"
			CreateNewsLog("CFB", message, "Transfer Portal", int(p.PreviousTeamID), ts)

			promise.MakePromise()
			db.Save(&promise)
			p.WillStay()
			repository.SaveCFBPlayer(p, db)
			db.Save(&p)
		}
	}

	ts.EnactPortalPhase()
	db.Save(&ts)
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

	if profile.PromiseID.Int64 > 0 {
		promiseID := strconv.Itoa(int(profile.PromiseID.Int64))
		promise := GetCollegePromiseByID(promiseID)
		repository.DeleteCollegePromise(promise, db)
	}

	db.Save(&profile)
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
	db.Save(&profile)
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
	coachMap := GetActiveCollegeCoachMap()
	teamMap := GetCollegeTeamMap()
	standingsMap := GetCollegeStandingsMap(seasonID)
	for _, teamProfile := range AITeams {
		if !teamProfile.IsAI {
			continue
		}
		team := teamMap[teamProfile.ID]
		teamStandings := standingsMap[uint(teamProfile.TeamID)]
		teamID := strconv.Itoa(int(teamProfile.ID))
		coach := coachMap[teamProfile.ID]
		portalProfileMap := getTransferPortalProfileMapByTeamID(teamID)
		roster := GetAllCollegePlayersByTeamId(teamID)
		rosterSize := len(roster)
		teamCap := 105
		if !teamProfile.IsFBS {
			teamCap = 80
		}
		if rosterSize >= teamCap {
			continue
		}
		teamNeedsMap := GetRecruitingNeeds(teamID)

		for _, r := range roster {
			if r.Year == 1 && !r.IsRedshirt {
				teamNeedsMap[r.Position] -= 1
			}
		}

		for _, tp := range transferPortalPlayers {
			isBadFit := IsBadSchemeFit(teamProfile.OffensiveScheme, teamProfile.DefensiveScheme, tp.Archetype, tp.Position)
			if isBadFit || teamNeedsMap[tp.Position] <= 0 || tp.PreviousTeamID == team.ID || portalProfileMap[tp.ID].CollegePlayerID == tp.ID || portalProfileMap[tp.ID].ID > 0 {
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
			} else if bias == specificCoach && tp.LegacyID == coach.ID {
				biasMod += 25
			} else if bias == legacy && tp.LegacyID == team.ID {
				biasMod += 25
			}

			diceRoll := util.GenerateIntFromRange(1, 50)
			if diceRoll < biasMod {
				portalProfile := structs.TransferPortalProfile{
					ProfileID:        teamProfile.ID,
					CollegePlayerID:  tp.ID,
					SeasonID:         uint(ts.CollegeSeasonID),
					TeamAbbreviation: teamProfile.TeamAbbreviation,
				}

				db.Create(&portalProfile)
			}
		}
	}
}

func AICoachAllocateAndPromisePhase() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	AITeams := GetOnlyAITeamRecruitingProfiles()
	transferPortalPlayerMap := GetCollegePlayerMap()
	coachMap := GetActiveCollegeCoachMap()
	// regionMap := util.GetRegionMap()
	// Shuffles the list of AI teams so that it's not always iterating from A-Z. Gives the teams at the lower end of the list a chance to recruit other croots
	rand.Shuffle(len(AITeams), func(i, j int) {
		AITeams[i], AITeams[j] = AITeams[j], AITeams[i]
	})

	for _, teamProfile := range AITeams {
		if teamProfile.SpentPoints >= teamProfile.WeeklyPoints {
			continue
		}

		teamID := strconv.Itoa(int(teamProfile.ID))
		roster := GetAllCollegePlayersByTeamId(teamID)
		teamCap := 105
		if !teamProfile.IsFBS {
			teamCap = 80
		}
		if len(roster) >= teamCap {
			continue
		}

		teamNeedsMap := GetRecruitingNeeds(teamID)

		for _, r := range roster {
			if r.Year == 1 && !r.IsRedshirt {
				teamNeedsMap[r.Position] -= 1
			}
		}

		portalProfiles := GetTransferPortalProfilesByTeamID(teamID)
		for _, profile := range portalProfiles {
			if profile.CurrentWeeksPoints > 0 || profile.RemovedFromBoard {
				continue
			}
			tp := transferPortalPlayerMap[profile.CollegePlayerID]
			// If player has already signed or if the position has been fulfilled
			isBadFit := IsBadSchemeFit(teamProfile.OffensiveScheme, teamProfile.DefensiveScheme, tp.Archetype, tp.Position)
			if isBadFit || tp.TeamID > 0 || tp.TransferStatus == 0 || tp.ID == 0 || teamNeedsMap[tp.Position] <= 0 {
				profile.Deactivate()
				repository.SaveTransferPortalProfile(profile, db)
				continue
			}
			playerID := strconv.Itoa(int(profile.CollegePlayerID))
			pointsRemaining := teamProfile.WeeklyPoints - teamProfile.SpentPoints
			if teamProfile.SpentPoints >= teamProfile.WeeklyPoints || pointsRemaining <= 0 || (pointsRemaining < 1 && pointsRemaining > 0) {
				break
			}

			removePlayerFromBoard := false
			num := 0.0

			profiles := GetTransferPortalProfilesByPlayerID(playerID)
			leadingTeamVal := IsAITeamContendingForPortalPlayer(profiles)
			if profile.CurrentWeeksPoints > 0 && profile.TotalPoints+float64(profile.CurrentWeeksPoints) >= float64(leadingTeamVal)*0.66 {
				// continue, leave everything alone
				continue
			} else if profile.CurrentWeeksPoints > 0 && profile.TotalPoints+float64(profile.CurrentWeeksPoints) < float64(leadingTeamVal)*0.66 {
				profile.Deactivate()
				db.Save(&profile)
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
			coach := coachMap[uint(teamProfile.TeamID)]

			min := coach.PointMin
			max := coach.PointMax
			if max > 10 {
				max = 10
			}
			num = util.GenerateFloatFromRange(float64(min), float64(max))
			if num > pointsRemaining {
				num = pointsRemaining
			}

			if num+profile.TotalPoints < float64(leadingTeamVal)*0.66 {
				removePlayerFromBoard = true
			}
			if leadingTeamVal < 14 {
				removePlayerFromBoard = false
			}

			if removePlayerFromBoard {
				profile.Deactivate()
				db.Save(&profile)
				continue
			}
			profile.AllocatePoints(int(num))

			// Generate Promise based on coach bias
			if profile.PromiseID.Int64 == 0 && !profile.RolledOnPromise {
				promiseOdds := getBasePromiseOdds(coach.TeambuildingPreference, coach.PromiseTendency)
				diceRoll := util.GenerateIntFromRange(1, 100)

				if diceRoll < promiseOdds {
					// Commit Promise
					promiseLevel := getPromiseLevel(coach.PromiseTendency)
					promiseWeight := "Medium"
					promiseType := ""
					benchmarkStr := ""
					promiseBenchmark := 0

					bias := tp.RecruitingBias
					// Rewrite
					if bias == closeToHome && (teamProfile.State == tp.State) {
						promiseType = "Home State Game"
						benchmarkStr = tp.State
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

						promiseWeight = getPromiseWeightByMinutesOrWins(promiseBenchmark)
					} else if bias == nationalChampionshipContender || bias == richHistory {
						// Promise based on wins
						promiseBenchmark = 5
						promiseType = "Wins"
						if promiseLevel == 1 {
							promiseBenchmark += 5
						} else if promiseLevel == -1 {
							promiseBenchmark -= 5
						}
						promiseWeight = getPromiseWeightByMinutesOrWins(promiseBenchmark)
					} else if bias == legacy && tp.LegacyID == uint(teamProfile.TeamID) {
						promiseType = "Legacy"
					} else if bias == specificCoach && tp.LegacyID == coach.ID {
						promiseType = "Specific Coach"
					} else if bias == differentState && teamProfile.State != tp.State {
						promiseType = "Different State"
						promiseWeight = "Low"
					}

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

				profile.ToggleRolledOnPromise()
			}
			// Save Profile
			repository.SaveTransferPortalProfile(profile, db)
		}
	}
}

func SyncTransferPortal() {
	db := dbprovider.GetInstance().GetDB()
	//GetCurrentWeek
	ts := GetTimestamp()
	// Use IsRecruitingLocked to lock the TP when not in use
	teamProfileMap := GetTeamProfileMap()
	transferPortalPlayers := GetTransferPortalPlayers()
	transferPortalProfileMap := getTransferPortalProfileMap(transferPortalPlayers)
	rosterMap := GetFullTeamRosterWithCrootsMap()

	if !ts.IsRecruitingLocked {
		ts.ToggleLockRecruiting()
		db.Save(&ts)
	}

	for _, portalPlayer := range transferPortalPlayers {
		// Skip over players that have already transferred
		if portalPlayer.TransferStatus != 2 || portalPlayer.TeamID > 0 {
			continue
		}

		portalProfiles := transferPortalProfileMap[portalPlayer.ID]
		if len(portalProfiles) == 0 && ts.TransferPortalRound < 10 {
			continue
		}

		// If no one has a profile on them during round 10
		if len(portalProfiles) == 0 && ts.TransferPortalRound == 10 && len(portalPlayer.TransferLikeliness) > 0 {
			portalPlayer.WillReturn()
			db.Save(&portalPlayer)
			continue
		}

		totalPointsOnPlayer := 0.0
		eligiblePointThreshold := 0.0
		readyToSign := false
		minSpendingCount := 100
		maxSpendingCount := 0
		signingMinimum := 0.66
		teamCount := 0
		eligibleTeams := []structs.TransferPortalProfile{}

		for i := range portalProfiles {
			promiseID := strconv.Itoa(int(portalProfiles[i].PromiseID.Int64))

			promise := GetCollegePromiseByID(promiseID)

			multiplier := getMultiplier(promise)
			portalProfiles[i].AddPointsToTotal(multiplier)
		}

		sort.Slice(portalProfiles, func(i, j int) bool {
			return portalProfiles[i].TotalPoints > portalProfiles[j].TotalPoints
		})

		for i := range portalProfiles {
			roster := rosterMap[portalProfiles[i].ProfileID]
			if len(roster) >= 13 {
				continue
			}
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

		if (teamCount == 1 && minSpendingCount >= 2) || (teamCount > 1 && minSpendingCount > 3 || ts.TransferPortalRound == 10) {
			// threshold met
			readyToSign = true
		}

		if readyToSign {
			var winningTeamID uint = 0
			var odds float64 = 0

			for winningTeamID == 0 {
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
					promise := GetCollegePromiseByCollegePlayerID(strconv.Itoa(int(portalPlayer.ID)), winningTeamIDSTR)

					promise.MakePromise()
					db.Save(&promise)

					teamProfile := teamProfileMap[winningTeamIDSTR]
					currentRoster := rosterMap[teamProfile.ID]
					teamCap := 105
					if !teamProfile.IsFBS {
						teamCap = 80
					}
					if len(currentRoster) < teamCap {
						portalPlayer.SignWithNewTeam(teamProfile.TeamID, teamProfile.TeamAbbreviation)
						message := portalPlayer.FirstName + " " + portalPlayer.LastName + ", " + strconv.Itoa(portalPlayer.Stars) + " star " + portalPlayer.Position + " from " + portalPlayer.PreviousTeam + " has signed with " + portalPlayer.TeamAbbr + " with " + strconv.Itoa(int(odds)) + " percent odds."
						CreateNewsLog("CFB", message, "Transfer Portal", int(winningTeamID), ts)
						fmt.Println("Created new log!")
						// Add player to existing roster map
						rosterMap[teamProfile.ID] = append(rosterMap[teamProfile.ID], portalPlayer)
						for i := range portalProfiles {
							if portalProfiles[i].ID == winningTeamID {
								portalProfiles[i].SignPlayer()
							} else {
								promise := GetCollegePromiseByCollegePlayerID(strconv.Itoa(int(portalPlayer.ID)), strconv.Itoa(int(portalProfiles[i].ProfileID)))
								if promise.ID > 0 {
									repository.DeleteCollegePromise(promise, db)
								}
							}
							portalProfiles[i].Lock()
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
			repository.SaveTransferPortalProfile(p, db)
			fmt.Println("Save transfer portal profile from " + portalPlayer.TeamAbbr + " towards " + portalPlayer.FirstName + " " + portalPlayer.LastName)
		}
		// Save Recruit
		repository.SaveCollegePlayerRecord(portalPlayer, db)
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
	db := dbprovider.GetInstance().GetDB()

	var profiles []structs.TransferPortalProfile
	var response []structs.TransferPortalProfileResponse
	err := db.Preload("CollegePlayer.Profiles").Preload("Promise").Where("profile_id = ? AND removed_from_board = ?", teamID, false).Find(&profiles).Error
	if err != nil {
		log.Fatalln("Error!: ", err)
	}

	for _, p := range profiles {
		if p.RemovedFromBoard {
			continue
		}
		cp := p.CollegePlayer
		cpResponse := structs.TransferPlayerResponse{}
		ovr := util.GetOverallGrade(cp.Overall, cp.Year)
		cpResponse.Map(cp, ovr)

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
			Promise:               p.Promise,
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

	var players []structs.CollegePlayer

	db.Preload("Profiles").Where("transfer_status = 2").Order("overall DESC").Find(&players)

	playerList := []structs.TransferPlayerResponse{}

	for _, p := range players {
		res := structs.TransferPlayerResponse{}
		ovr := util.GetOverallGrade(p.Overall, p.Year)
		res.Map(p, ovr)

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

func GetCollegePlayerSeasonStatMap(seasonID string) map[uint]structs.CollegePlayerSeasonStats {
	seasonStatMap := make(map[uint]structs.CollegePlayerSeasonStats)

	seasonStats := GetCollegePlayerSeasonStatsBySeason(seasonID)
	for _, stat := range seasonStats {
		seasonStatMap[stat.CollegePlayerID] = stat
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
	min := 15
	max := 100
	if likeliness == "Low" {
		max = 33
	} else if likeliness == "Medium" {
		min = 34
		max = 60
	} else {
		min = 61
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

func getPromiseWeightByMinutesOrWins(benchmark int) string {
	weight := "Medium"
	if benchmark <= 40 {
		weight = "Very High"
	}
	if benchmark <= 25 {
		weight = "High"
	}
	if benchmark <= 20 {
		weight = "Medium"
	}
	if benchmark <= 10 {
		weight = "Low"
	}
	if benchmark <= 5 {
		weight = "Very Low"
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
	if pr.ID == 0 {
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

	seasonStats := GetCollegePlayerSeasonStatsByPlayerIDAndSeason(id, seasonIDSTR)
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
	teamRecruitingProfiles := GetRecruitingProfileForRecruitSync()

	teamMap := make(map[string]*structs.RecruitingTeamProfile)
	for i := 0; i < len(teamRecruitingProfiles); i++ {
		teamMap[strconv.Itoa(int(teamRecruitingProfiles[i].ID))] = &teamRecruitingProfiles[i]
	}

	return teamMap
}

func getSchemeMod(tp *structs.RecruitingTeamProfile, p structs.CollegePlayer, drop, gain float64) float64 {
	schemeMod := 0.0

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
