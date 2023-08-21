package managers

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
	"gorm.io/gorm"
)

func GetRecruitingModifiers() structs.AdminRecruitModifier {
	db := dbprovider.GetInstance().GetDB()

	var recruitModifiers structs.AdminRecruitModifier

	db.First(&recruitModifiers)

	return recruitModifiers
}

func SyncRecruiting(timestamp structs.Timestamp) {
	db := dbprovider.GetInstance().GetDB()
	fmt.Println(time.Now().UnixNano())
	rand.Seed(time.Now().UnixNano())
	//GetCurrentWeek

	if timestamp.RecruitingSynced {
		log.Fatalln("Recruiting already ran for this week. Please wait until next week to sync recruiting again.")
	}

	eligibleThresholdPercentage := 0.66
	pointLimit := 20

	recruitProfilePointsMap := util.GetTeamPointsMap()

	recruitModifiers := GetRecruitingModifiers()

	teamRecruitingProfiles := GetRecruitingProfileForRecruitSync()

	teamMap := make(map[string]*structs.RecruitingTeamProfile)

	for i := 0; i < len(teamRecruitingProfiles); i++ {
		teamMap[strconv.Itoa(int(teamRecruitingProfiles[i].ID))] = &teamRecruitingProfiles[i]
	}

	var recruitProfiles []structs.RecruitPlayerProfile

	// Get every recruit
	recruits := GetAllUnsignedRecruits()

	// Iterate through every recruit
	for _, recruit := range recruits {
		recruitProfiles = GetRecruitPlayerProfilesByRecruitId(strconv.Itoa(int(recruit.ID)))

		if len(recruitProfiles) == 0 {
			fmt.Println("Skipping over " + recruit.FirstName + " " + recruit.LastName + " because no one is recruiting them.")
			continue
		}

		var recruitProfilesWithScholarship []structs.RecruitPlayerProfile

		eligibleTeams := 0

		var totalPointsOnRecruit float64 = 0

		var eligiblePointThreshold float64 = 0

		var signThreshold float64

		pointsPlaced := false
		spendingCountAdjusted := false

		allocatePointsToRecruit(recruit, &recruitProfiles, float64(pointLimit), &spendingCountAdjusted, &pointsPlaced, timestamp, &recruitProfilePointsMap, db)

		if !pointsPlaced && !spendingCountAdjusted {
			fmt.Println("Skipping over " + recruit.FirstName + " " + recruit.LastName)
			continue
		}

		sort.Sort(structs.ByPoints(recruitProfiles))

		for i := 0; i < len(recruitProfiles) && pointsPlaced; i++ {
			recruitTeamProfile := teamMap[strconv.Itoa(recruitProfiles[i].ProfileID)]
			if recruitTeamProfile.TotalCommitments >= recruitTeamProfile.RecruitClassSize {
				continue
			}
			if eligiblePointThreshold == 0 && recruitProfiles[i].Scholarship {
				eligiblePointThreshold = float64(recruitProfiles[i].TotalPoints) * eligibleThresholdPercentage
			}

			if recruitProfiles[i].Scholarship && recruitProfiles[i].TotalPoints >= eligiblePointThreshold {
				totalPointsOnRecruit += recruitProfiles[i].TotalPoints
				eligibleTeams += 1
				recruitProfilesWithScholarship = append(recruitProfilesWithScholarship, recruitProfiles[i])
			}
		}

		// Change?
		// Assign point totals
		// If there are any modifiers
		// Evaluate
		firstMod := float64(recruitModifiers.ModifierOne - timestamp.CollegeWeek)
		secondMod := float64(eligibleTeams) / float64(recruit.RecruitingModifier)
		thirdMod := math.Log10(float64(recruitModifiers.WeeksOfRecruiting - timestamp.CollegeWeek))
		signThreshold = firstMod * secondMod * thirdMod
		recruit.ApplyRecruitingStatus(totalPointsOnRecruit, signThreshold)

		// Change logic to withold teams without available scholarships
		passedTheSigningThreshold := float64(totalPointsOnRecruit) > signThreshold && eligibleTeams > 0 && pointsPlaced
		if passedTheSigningThreshold {
			winningTeamID := 0
			var odds float64 = 0

			for winningTeamID == 0 && len(recruitProfilesWithScholarship) > 0 {
				percentageOdds := rand.Float64() * (totalPointsOnRecruit)
				var currentProbability float64 = 0

				for i := 0; i < len(recruitProfilesWithScholarship); i++ {
					// If a team has no available scholarships or if a team has 25 commitments, continue
					currentProbability += recruitProfilesWithScholarship[i].TotalPoints
					if float64(percentageOdds) <= currentProbability {
						// WINNING TEAM
						winningTeamID = recruitProfilesWithScholarship[i].ProfileID
						odds = float64(recruitProfilesWithScholarship[i].TotalPoints) / float64(totalPointsOnRecruit) * 100
						break
					}
				}

				if winningTeamID > 0 && len(recruitProfilesWithScholarship) > 0 {
					recruitTeamProfile := teamMap[strconv.Itoa(winningTeamID)]
					if recruitTeamProfile.TotalCommitments < recruitTeamProfile.RecruitClassSize {
						recruitTeamProfile.IncreaseCommitCount()
						teamAbbreviation := recruitTeamProfile.TeamAbbreviation
						recruit.AssignCollege(teamAbbreviation)

						newsLog := structs.NewsLog{
							TeamID:      winningTeamID,
							WeekID:      timestamp.CollegeWeekID + 1,
							Week:        timestamp.CollegeWeek,
							SeasonID:    timestamp.CollegeSeasonID,
							MessageType: "Commitment",
							League:      "CFB",
							Message:     recruit.FirstName + " " + recruit.LastName + ", " + strconv.Itoa(recruit.Stars) + " star " + recruit.Position + " from " + recruit.City + ", " + recruit.State + " has signed with " + recruit.College + " with " + strconv.Itoa(int(odds)) + " percent odds.",
						}

						db.Create(&newsLog)
						fmt.Println("Created new log!")

						db.Save(&recruitTeamProfile)
						fmt.Println("Saved " + recruitTeamProfile.TeamAbbreviation + " profile.")

						for i := 0; i < len(recruitProfiles); i++ {
							if recruitProfiles[i].ProfileID == winningTeamID {
								recruitProfiles[i].SignPlayer()
							} else {
								recruitProfiles[i].LockPlayer()
								if recruitProfiles[i].Scholarship {
									tp := teamMap[strconv.Itoa(recruitProfiles[i].ProfileID)]
									tp.ReallocateScholarship()
									err := db.Save(&tp).Error
									if err != nil {
										fmt.Println(err.Error())
										log.Fatalf("Could not sync recruiting profile.")
									}

									fmt.Println("Reallocated Scholarship to " + tp.TeamAbbreviation)
								}
							}
						}

					} else {
						recruitProfilesWithScholarship = util.FilterOutRecruitingProfile(recruitProfilesWithScholarship, winningTeamID)
						winningTeamID = 0
						if len(recruitProfilesWithScholarship) == 0 {
							break
						}
						totalPointsOnRecruit = 0
						for _, rp := range recruitProfilesWithScholarship {
							totalPointsOnRecruit += rp.TotalPoints
						}
					}
				}
			}
			recruit.UpdateTeamID(int(winningTeamID))

		}

		// Save Player Files towards Recruit
		for _, rp := range recruitProfiles {
			// Save Team Profile
			err := db.Save(&rp).Error
			if err != nil {
				fmt.Println(err.Error())
				log.Fatalf("Could not sync recruiting profile.")
			}

			fmt.Println("Save recruit profile from " + rp.TeamAbbreviation + " towards " + recruit.FirstName + " " + recruit.LastName)
		}
		db.Save(&recruit)
	}

	updateTeamRankings(teamRecruitingProfiles, teamMap, recruitProfilePointsMap, db)

}

func SyncRecruitingEfficiency(timestamp structs.Timestamp) {
	db := dbprovider.GetInstance().GetDB()

	// Get All Team Recruiting Profiles
	teams := GetAllCollegeTeamsWithRecruitingProfileAndCoach()

	// Iterate through all profiles

	// var teamProfilesToSave []structs.RecruitingTeamProfile

	for _, team := range teams {
		// Get all games by team within a season

		if team.Coach == "" || team.Coach == "AI" {
			continue
		}

		teamProfile := team.RecruitingProfile

		coach := team.CollegeCoach

		// What about previous season?
		currentSeasonGames := GetCollegeGamesByTeamIdAndSeasonId(
			strconv.Itoa(teamProfile.TeamID), strconv.Itoa(timestamp.CollegeSeasonID))

		currentSeasonWins, currentSeasonLosses := util.GetWinsAndLossesForCollegeGames(currentSeasonGames, teamProfile.TeamID, false)
		currentConferenceWins, currentConferenceLosses := util.GetWinsAndLossesForCollegeGames(currentSeasonGames, teamProfile.TeamID, true)

		previousSeasonGames := GetCollegeGamesByTeamIdAndSeasonId(
			strconv.Itoa(teamProfile.TeamID), strconv.Itoa(timestamp.CollegeSeasonID-1))

		previousSeasonWins, previousSeasonLosses := util.GetWinsAndLossesForCollegeGames(previousSeasonGames, teamProfile.TeamID, false)
		previousConferenceWins, previousConferenceLosses := util.GetWinsAndLossesForCollegeGames(previousSeasonGames, teamProfile.TeamID, true)
		// Do calculation for current season losses

		// Current Season Win Percentage
		var cswp float64 = 1
		var csweight float64 = .125
		var ccwp float64 = 1
		var ccweight float64 = .125
		var pswp float64 = 1
		var psweight float64 = .05
		var pcwp float64 = 1
		var pcweight float64 = .05
		var coachwp float64 = 1
		var coachweight float64 = .1

		var postSeasonVal float64
		var conferenceChampionshipVal float64
		var postseasonweight float64 = 0.025

		if timestamp.CollegeWeek < 15 {
			postSeasonVal = util.GetPostSeasonWeight(previousSeasonGames, teamProfile.TeamID)
			conferenceChampionshipVal = util.GetConferenceChampionshipWeight(previousSeasonGames, teamProfile.TeamID)
		} else {
			postSeasonVal = util.GetPostSeasonWeight(currentSeasonGames, teamProfile.TeamID)
			conferenceChampionshipVal = util.GetConferenceChampionshipWeight(currentSeasonGames, teamProfile.TeamID)
		}

		if currentSeasonWins+currentSeasonLosses > 0 {
			cswp = float64(currentSeasonWins) / float64(currentSeasonWins+currentSeasonLosses)
			if team.ConferenceID != 13 && currentConferenceWins+currentConferenceLosses > 0 {
				ccwp = float64(currentConferenceWins) / float64(currentConferenceWins+currentConferenceLosses)
			}
		}

		// Previous Season Win Percentage
		if previousSeasonWins+previousSeasonLosses > 0 {
			pswp = float64(previousSeasonWins) / float64(previousSeasonWins+previousSeasonLosses)
			if team.ConferenceID != 13 {
				pcwp = float64(previousConferenceWins) / float64(previousConferenceWins+previousConferenceLosses)
			}
		}

		if coach.OverallWins+coach.OverallLosses > 0 {
			coachwp = float64(coach.OverallWins) / float64(coach.OverallWins+coach.OverallLosses)
		}

		res := teamProfile.BaseEfficiencyScore

		cswsum := cswp * csweight                             // Current Season Win Percentage
		ccwsum := ccwp * ccweight                             // Current Conference Wins
		pswsum := pswp * psweight                             // Prev Season Wins
		pcwsum := pcwp * pcweight                             // Prev Conference Wins
		coachsum := coachwp * coachweight                     // Coach Overall Wins
		postseasonsum := postSeasonVal * postseasonweight     // PostSeason Check
		ccsum := conferenceChampionshipVal * postseasonweight // ConferenceChampionship Check

		totalSum := cswsum + ccwsum + pswsum + pcwsum + coachsum + postseasonsum + ccsum
		// totalWeight := csweight + ccweight + psweight + pcweight + coachweight + (postseasonweight * 2)

		// RES Calculation
		// Base of .8

		teamProfile.AssignRES(res + (totalSum * 0.4))

		err := db.Save(&teamProfile).Error
		if err != nil {
			fmt.Println(err.Error())
			log.Fatalf("Could not sync all team profiles.")
		}
		fmt.Println("Saved RES for Team: " + team.TeamAbbr)

		// teamProfilesToSave = append(teamProfilesToSave, teamProfile)
	}
	// Save the Recruiting Profiles
	// err = db.Save(&recruitProfilesToSave).Error
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	log.Fatalf("Could not sync res to all recruits")
	// }
}

func SyncAllMissingEfficiencies() {
	db := dbprovider.GetInstance().GetDB()
	recruitingProfiles := GetRecruitingProfileForRecruitSync()

	for _, rp := range recruitingProfiles {
		playerProfiles := GetOnlyRecruitProfilesByTeamProfileID(strconv.Itoa(int(rp.ID)))

		for _, pp := range playerProfiles {
			if pp.RecruitingEfficiencyScore != 0 {
				continue
			}
			pp.AssignRES(rp.RecruitingEfficiencyScore)

			err := db.Save(&pp).Error
			if err != nil {
				log.Panicln("COULD NOT SAVE ALL PROFILES")
			}
			fmt.Println("Saved profile from " + rp.Team + " towards " + strconv.Itoa(int(pp.ID)))
		}
	}
}

func SyncTeamRankings() {
	db := dbprovider.GetInstance().GetDB()
	// Update rank system for all teams
	teamRecruitingProfiles := GetRecruitingProfileForRecruitSync()

	var maxESPNScore float64 = 0
	var minESPNScore float64 = 100000
	var maxRivalsScore float64 = 0
	var minRivalsScore float64 = 100000
	var max247Score float64 = 0
	var min247Score float64 = 100000

	for i := 0; i < len(teamRecruitingProfiles); i++ {

		signedRecruits := GetSignedRecruitsByTeamProfileID(strconv.Itoa(teamRecruitingProfiles[i].TeamID))

		teamRecruitingProfiles[i].UpdateTotalSignedRecruits(len(signedRecruits))

		team247Rank := Get247TeamRanking(teamRecruitingProfiles[i], signedRecruits)
		teamESPNRank := GetESPNTeamRanking(teamRecruitingProfiles[i], signedRecruits)
		teamRivalsRank := GetRivalsTeamRanking(teamRecruitingProfiles[i], signedRecruits)

		teamRecruitingProfiles[i].Assign247Rank(team247Rank)
		teamRecruitingProfiles[i].AssignESPNRank(teamESPNRank)
		teamRecruitingProfiles[i].AssignRivalsRank(teamRivalsRank)
		if teamESPNRank > maxESPNScore {
			maxESPNScore = teamESPNRank
		}
		if teamESPNRank < minESPNScore {
			minESPNScore = teamESPNRank
		}
		if teamRivalsRank > maxRivalsScore {
			maxRivalsScore = teamRivalsRank
		}
		if teamRivalsRank < minRivalsScore {
			minRivalsScore = teamRivalsRank
		}
		if team247Rank > max247Score {
			max247Score = team247Rank
		}
		if team247Rank < min247Score {
			min247Score = team247Rank
		}
	}

	espnDivisor := (maxESPNScore - minESPNScore)
	divisor247 := (max247Score - min247Score)
	rivalsDivisor := (maxRivalsScore - minRivalsScore)

	for _, rp := range teamRecruitingProfiles {

		var avg float64 = 0
		if espnDivisor > 0 && divisor247 > 0 && rivalsDivisor > 0 {
			distributionESPN := (rp.ESPNScore - minESPNScore) / espnDivisor
			distribution247 := (rp.Rank247Score - min247Score) / divisor247
			distributionRivals := (rp.RivalsScore - minRivalsScore) / rivalsDivisor

			avg = (distributionESPN + distribution247 + distributionRivals)

			rp.AssignCompositeRank(avg)
		}

		// Save TEAM Recruiting Profile
		err := db.Save(&rp).Error
		if err != nil {
			fmt.Println(err.Error())
			log.Fatalf("Could not save timestamp")
		}
		fmt.Println("Saved Rank Scores for Team " + rp.TeamAbbreviation)
	}
}

func FixSmallTownBigCityAIBoards() {
	db := dbprovider.GetInstance().GetDB()

	teams := GetAllCollegeTeams()

	for _, t := range teams {

		teamID := strconv.Itoa(int(t.ID))

		profile := structs.RecruitingTeamProfile{}

		err := db.Preload("Affinities").Where("id = ?", teamID).Find(&profile).Error
		if err != nil {
			log.Panicln(err)
		}

		playerProfiles := GetRecruitingProfileForTeamBoardByTeamID(teamID)

		croots := playerProfiles.Recruits
		smallTownApplicable := true
		bigCityApplicable := true
		for _, croot := range croots {
			r := croot.Recruit

			individualProfile := GetRecruitProfileByPlayerId(strconv.Itoa(int(r.ID)), teamID)
			fixApplied := false

			for _, affinity := range profile.Affinities {
				if affinity.AffinityName != "Small Town" && affinity.AffinityName != "Big City" {
					continue
				}

				if affinity.AffinityName == "Small Town" && !affinity.IsApplicable {
					smallTownApplicable = false
					continue
				}
				if affinity.AffinityName == "Big City" && !affinity.IsApplicable {
					bigCityApplicable = false
					continue
				}
				if r.AffinityOne == "Small Town" && isAffinityApplicable("Small Town", affinity) && !individualProfile.AffinityOneEligible {
					// Fix Affinity One in recruiting player profile
					individualProfile.ToggleAffinityOne()
					fixApplied = true
				}

				if r.AffinityTwo == "Small Town" && isAffinityApplicable("Small Town", affinity) && !individualProfile.AffinityTwoEligible {
					// Fix Affinity One in recruiting player profile
					individualProfile.ToggleAffinityTwo()
					fixApplied = true
				}

				if r.AffinityOne == "Big City" && isAffinityApplicable("Big City", affinity) && !individualProfile.AffinityOneEligible {
					// Fix Affinity One in recruiting player profile
					individualProfile.ToggleAffinityOne()
					fixApplied = true
				}

				if r.AffinityTwo == "Big City" && isAffinityApplicable("Big City", affinity) && !individualProfile.AffinityTwoEligible {
					// Fix Affinity One in recruiting player profile
					individualProfile.ToggleAffinityTwo()
					fixApplied = true
				}
			}

			if !smallTownApplicable && !bigCityApplicable {
				break
			}
			if fixApplied {
				db.Save(&individualProfile)
			}
		}
	}
}

func FillAIRecruitingBoards() {
	db := dbprovider.GetInstance().GetDB()
	fmt.Println(time.Now().UnixNano())
	rand.Seed(time.Now().UnixNano())
	ts := GetTimestamp()

	AITeams := GetOnlyAITeamRecruitingProfiles()
	UnsignedRecruits := GetAllUnsignedRecruits()
	stateMatcher := util.GetStateMatcher()
	regionMatcher := util.GetStateRegionMatcher()

	boardCount := 75

	for _, team := range AITeams {
		count := 0
		if !team.IsAI || team.TotalCommitments >= team.RecruitClassSize {
			continue
		}

		existingBoard := GetOnlyRecruitProfilesByTeamProfileID(strconv.Itoa(int(team.ID)))
		teamNeeds := GetRecruitingNeeds(strconv.Itoa(int(team.ID)))
		// Get Current Count of the existing board
		for _, r := range existingBoard {
			if r.RemovedFromBoard {
				continue
			}

			if r.IsSigned {
				teamNeeds[r.Recruit.Position] -= 1
			}

			count++
		}

		for k := range teamNeeds {
			if teamNeeds[k] > 0 {
				teamNeeds[k] *= 4
			}
		}

		for _, croot := range UnsignedRecruits {
			if count >= boardCount {
				break
			}

			if teamNeeds[croot.Position] < 1 {
				continue
			}

			passOnRecruit := false

			if croot.Stars == 5 && !isBlueBlood(team.AIBehavior) {
				passOnRecruit = true
			}

			if croot.Stars > 3 && !team.IsFBS {
				passOnRecruit = true
			}

			// Conditions in which the team should not recruit this particular recruit
			if passOnRecruit {
				continue
			}

			// Check and see if the croot already exists on the player's board
			crootProfile := GetRecruitProfileByPlayerId(strconv.Itoa(int(croot.ID)), strconv.Itoa(int(team.ID)))
			if uint(crootProfile.ProfileID) == team.ID || crootProfile.ID > 0 || crootProfile.RemovedFromBoard || crootProfile.IsLocked {
				fmt.Println(croot.FirstName + " " + croot.LastName + " is already on " + team.TeamAbbreviation + "'s board.")
				continue
			}

			crootProfiles := GetRecruitPlayerProfilesByRecruitId(strconv.Itoa(int(croot.ID)))
			affinityMod := 0
			teamCount := 0

			for _, crootProfile := range crootProfiles {
				if crootProfile.RemovedFromBoard {
					continue
				}
				teamCount++
			}

			leadingVal := util.IsAITeamContendingForCroot(crootProfiles)
			if leadingVal > 15 {
				continue
			}

			odds := 5

			if ts.CollegeWeek > 5 {
				odds = 10
			}

			if croot.State == team.State {
				odds = 25
			}

			closeToHome := util.IsCrootCloseToHome(croot.State, croot.City, team.State, team.TeamAbbreviation, stateMatcher, regionMatcher)
			// In Region
			if closeToHome && croot.State != team.State {
				odds = 15
			}

			affinityOneApplicable := false
			affinityTwoApplicable := false

			if team.AIBehavior == "G5" {
				if croot.Stars > 3 {
					odds -= 15
				} else {
					odds += 5
				}
			}

			if team.AIBehavior == "Doormat" {
				if croot.Stars > 2 {
					odds -= 20
				} else {
					odds += 10
				}
			}

			for _, affinity := range team.Affinities {
				if (doesCrootHaveAffinity("Close to Home", croot)) && closeToHome {
					if team.IsFBS {
						odds += 33
					} else {
						odds += 20
					}

					if croot.AffinityOne == "Close to Home" {
						affinityOneApplicable = true
						affinityMod += 3
					}

					if croot.AffinityTwo == "Close to Home" {
						affinityTwoApplicable = true
						affinityMod += 3
					}
				}

				if doesCrootHaveAffinity("Academics", croot) && isAffinityApplicable("Academics", affinity) {
					if team.IsFBS {
						odds += 33
					} else {
						odds += 17
					}

					if croot.AffinityOne == "Academics" {
						affinityOneApplicable = true
						affinityMod += 3
					}

					if croot.AffinityTwo == "Academics" {
						affinityTwoApplicable = true
						affinityMod += 3
					}
				}

				if doesCrootHaveAffinity("Frontrunner", croot) && isAffinityApplicable("Frontrunner", affinity) {
					if team.IsFBS {
						odds += 33
					} else {
						odds += 17
					}

					if isBlueBlood(team.AIBehavior) || team.AIBehavior == "Playoff Buster" {
						odds += 5
					}

					if croot.AffinityOne == "Frontrunner" {
						affinityOneApplicable = true
						affinityMod += 3
					}

					if croot.AffinityTwo == "Frontrunner" {
						affinityTwoApplicable = true
						affinityMod += 3
					}
				}

				if doesCrootHaveAffinity("Religion", croot) && isAffinityApplicable("Religion", affinity) {
					if team.IsFBS {
						odds += 33
					} else {
						odds += 17
					}

					if croot.AffinityOne == "Religion" {
						affinityOneApplicable = true
						affinityMod += 3
					}

					if croot.AffinityTwo == "Religion" {
						affinityTwoApplicable = true
						affinityMod += 3
					}
				}

				if doesCrootHaveAffinity("Service", croot) && isAffinityApplicable("Service", affinity) {
					if team.IsFBS {
						odds += 33
					} else {
						odds += 17
					}

					if croot.AffinityOne == "Service" {
						affinityOneApplicable = true
						affinityMod += 3
					}

					if croot.AffinityTwo == "Service" {
						affinityTwoApplicable = true
						affinityMod += 3
					}
				}

				if doesCrootHaveAffinity("Small School", croot) && isAffinityApplicable("Small School", affinity) {
					if team.IsFBS {
						odds += 33
					} else {
						odds += 17
					}

					if croot.AffinityOne == "Small School" {
						affinityOneApplicable = true
						affinityMod += 3
					}

					if croot.AffinityTwo == "Small School" {
						affinityTwoApplicable = true
						affinityMod += 3
					}
				}

				if doesCrootHaveAffinity("Small Town", croot) && isAffinityApplicable("Small Town", affinity) {
					if team.IsFBS {
						odds += 33
					} else {
						odds += 17
					}

					if croot.AffinityOne == "Small Town" {
						affinityOneApplicable = true
						affinityMod += 3
					}

					if croot.AffinityTwo == "Small Town" {
						affinityTwoApplicable = true
						affinityMod += 3
					}
				}

				if doesCrootHaveAffinity("Big City", croot) && isAffinityApplicable("Big City", affinity) {
					if team.IsFBS {
						odds += 33
					} else {
						odds += 17
					}

					if croot.AffinityOne == "Big City" {
						affinityOneApplicable = true
						affinityMod += 3
					}

					if croot.AffinityTwo == "Big City" {
						affinityTwoApplicable = true
						affinityMod += 3
					}
				}
			}

			chance := util.GenerateIntFromRange(1, 100)

			willAddToBoard := isHighlyContestedCroot(affinityMod, teamCount, ts.CollegeWeek)

			if chance <= odds && willAddToBoard {
				playerProfile := structs.RecruitPlayerProfile{
					RecruitID:                 int(croot.ID),
					ProfileID:                 int(team.ID),
					SeasonID:                  ts.CollegeSeasonID,
					TotalPoints:               0,
					CurrentWeeksPoints:        0,
					SpendingCount:             0,
					Scholarship:               false,
					ScholarshipRevoked:        false,
					TeamAbbreviation:          team.TeamAbbreviation,
					AffinityOneEligible:       affinityOneApplicable,
					AffinityTwoEligible:       affinityTwoApplicable,
					RecruitingEfficiencyScore: 1,
					IsSigned:                  false,
					IsLocked:                  false,
				}

				err := db.Create(&playerProfile).Error
				if err != nil {
					log.Fatalln("Could not add " + croot.FirstName + " " + croot.LastName + " to " + team.TeamAbbreviation + "'s Recruiting Board.")
				}

				teamNeeds[croot.Position] -= 1
				count++
			}
		}
	}
}

func AllocatePointsToAIBoards() {
	db := dbprovider.GetInstance().GetDB()
	fmt.Println(time.Now().UnixNano())
	rand.Seed(time.Now().UnixNano())
	ts := GetTimestamp()

	AITeams := GetOnlyAITeamRecruitingProfiles()

	for _, team := range AITeams {
		if team.SpentPoints >= team.WeeklyPoints || team.TotalCommitments >= team.RecruitClassSize {
			continue
		}

		teamID := strconv.Itoa(int(team.ID))

		teamRecruits := GetRecruitsForAIPointSync(teamID)

		teamNeedsMap := GetRecruitingNeeds(teamID)

		// Safety check to make sure teams aren't recruiting too many in one position
		for _, croot := range teamRecruits {
			if croot.IsSigned && croot.TeamAbbreviation == team.TeamAbbreviation && ts.CollegeWeek < 17 {
				teamNeedsMap[croot.Recruit.Position] -= 1
			}
		}

		for _, croot := range teamRecruits {
			pointsRemaining := team.WeeklyPoints - team.SpentPoints
			if team.SpentPoints >= team.WeeklyPoints || pointsRemaining <= 0 || (pointsRemaining < 1 && pointsRemaining > 0) {
				break
			}

			if croot.IsSigned || croot.CurrentWeeksPoints > 0 || croot.ScholarshipRevoked {
				continue
			}

			removeCrootFromBoard := false
			var num float64 = 0
			recruitID := strconv.Itoa(int(croot.RecruitID))

			if (croot.IsLocked && croot.TeamAbbreviation != croot.Recruit.College) || teamNeedsMap[croot.Recruit.Position] <= 0 {
				removeCrootFromBoard = true
			}

			if !removeCrootFromBoard {
				profiles := GetRecruitPlayerProfilesByRecruitId(recruitID)

				if croot.PreviousWeekPoints > 0 {
					leadingTeamVal := util.IsAITeamContendingForCroot(profiles)

					if croot.PreviousWeekPoints+croot.TotalPoints >= leadingTeamVal*0.66 || leadingTeamVal < 15 {
						num = croot.PreviousWeekPoints
						if num > pointsRemaining {
							num = pointsRemaining
						}
					} else {
						removeCrootFromBoard = true
					}
				} else {
					maxChance := 2
					if ts.CollegeWeek > 3 {
						maxChance = 4
					}
					chance := util.GenerateIntFromRange(1, maxChance)
					if (chance < 2 && ts.CollegeWeek <= 3) || (chance < 4 && ts.CollegeWeek > 3) {
						continue
					}

					min := 5
					max := 15

					if team.AIBehavior == "Blue Blood" || team.AIBehavior == "Playoff Buster" {
						min = 8
					} else if team.AIBehavior == "Doormat" {
						max = 7
						min = 1
					} else if team.AIBehavior == "G5" {
						min = 1
						max = 10
					}

					num = float64(util.GenerateIntFromRange(min, max))
					if num > pointsRemaining {
						num = pointsRemaining
					}

					leadingTeamVal := util.IsAITeamContendingForCroot(profiles)

					if float64(num)+croot.TotalPoints < leadingTeamVal*0.66 {
						removeCrootFromBoard = true
					}
					if leadingTeamVal < 15 {
						removeCrootFromBoard = false
					}
				}
			}

			if removeCrootFromBoard || (team.ScholarshipsAvailable == 0 && !croot.Scholarship) {
				if croot.Scholarship {
					croot.ToggleScholarship(false, true)
					team.ReallocateScholarship()
				}
				croot.ToggleRemoveFromBoard()
				fmt.Println("Because " + croot.Recruit.FirstName + " " + croot.Recruit.LastName + " is heavily considering other teams, they are being removed from " + team.TeamAbbreviation + "'s Recruiting Board.")
				db.Save(&croot)
				continue
			}

			if ts.CollegeWeek == 20 {
				num = 2
			}

			croot.AllocateCurrentWeekPoints(num)
			if !croot.Scholarship && team.ScholarshipsAvailable > 0 {
				croot.ToggleScholarship(true, false)
				team.SubtractScholarshipsAvailable()
			}

			team.AIAllocateSpentPoints(num)
			db.Save(&croot)
			fmt.Println(team.TeamAbbreviation + " allocating " + strconv.Itoa(int(num)) + " points to " + croot.Recruit.FirstName + " " + croot.Recruit.LastName)

		}
		// Save Team Profile after iterating through recruits
		fmt.Println("Saved " + team.TeamAbbreviation + " Recruiting Board!")
		db.Save(&team)
	}
}

func ResetAIBoardsForCompletedTeams() {
	db := dbprovider.GetInstance().GetDB()

	AITeams := GetOnlyAITeamRecruitingProfiles()

	for _, team := range AITeams {
		// If a team already has the maximum allowed for their recruiting class, take all Recruit Profiles for that team where the recruit hasn't signed, and reset their total points.
		// This is so that these unsigned recruits can be recruited for and will allow the AI to put points onto those recruits.

		if team.TotalCommitments >= team.RecruitClassSize {
			teamRecruits := GetRecruitsByTeamProfileID(strconv.Itoa(int(team.ID)))

			for _, croot := range teamRecruits {
				if croot.IsSigned || croot.IsLocked || croot.TotalPoints == 0 {
					continue
				}
				croot.ResetTotalPoints()
				db.Save(&croot)
			}
			team.ResetSpentPoints()
			db.Save(&team)
		}
	}
}

func isBlueBlood(behavior string) bool {
	return behavior == "Blue Blood"
}

func isAffinityApplicable(affinity string, af structs.ProfileAffinity) bool {
	return af.AffinityName == affinity && af.IsApplicable
}

func doesCrootHaveAffinity(affinity string, croot structs.Recruit) bool {
	return croot.AffinityOne == affinity || croot.AffinityTwo == affinity
}

func isHighlyContestedCroot(mod int, teams int, CollegeWeek int) bool {
	if CollegeWeek == 20 && teams > 1 {
		return false
	}
	chance := util.GenerateIntFromRange(1, 5)
	chance += mod

	return chance > teams
}

func allocatePointsToRecruit(recruit structs.Recruit, recruitProfiles *[]structs.RecruitPlayerProfile, pointLimit float64, spendingCountAdjusted *bool, pointsPlaced *bool, timestamp structs.Timestamp, recruitProfilePointsMap *map[string]float64, db *gorm.DB) {
	// numWorkers := 3
	numWorkers := runtime.NumCPU()
	if numWorkers > 3 {
		numWorkers = 3
	}
	jobs := make(chan int, len(*recruitProfiles))
	results := make(chan error, len(*recruitProfiles))

	// This starts up numWorkers number of workers, initially blocked because there are no jobs yet.
	for w := 1; w <= numWorkers; w++ {
		go func(jobs <-chan int, results chan<- error, w int) {
			for i := range jobs {
				err := processRecruitProfile(i, recruit, recruitProfiles, pointLimit, spendingCountAdjusted, pointsPlaced, timestamp, recruitProfilePointsMap, db)
				results <- err
			}
		}(jobs, results, w)
	}

	// Here we send len(*recruitProfiles) jobs and then close the channel.
	for i := 0; i < len(*recruitProfiles); i++ {
		jobs <- i
	}
	close(jobs)

	// Finally, we collect all the results.
	// This ensures the function doesn't return until we've processed all recruit profiles.
	for i := 0; i < len(*recruitProfiles); i++ {
		err := <-results
		if err != nil {
			fmt.Println(err)
			log.Fatalf("Could not process recruit profile: %v", err)
		}
	}
}

func processRecruitProfile(i int, recruit structs.Recruit, recruitProfiles *[]structs.RecruitPlayerProfile, pointLimit float64, spendingCountAdjusted *bool, pointsPlaced *bool, timestamp structs.Timestamp, recruitProfilePointsMap *map[string]float64, db *gorm.DB) error {
	m := &sync.Mutex{}
	affinityBonus := 0.1

	if (*recruitProfiles)[i].CurrentWeeksPoints == 0 {
		if (*recruitProfiles)[i].SpendingCount > 0 {
			(*recruitProfiles)[i].ResetSpendingCount()
			*spendingCountAdjusted = true
			fmt.Println("Resetting spending count for " + recruit.FirstName + " " + recruit.LastName + " for " + (*recruitProfiles)[i].TeamAbbreviation)
		}
		return nil
	} else {
		*pointsPlaced = true
	}

	rpa := structs.RecruitPointAllocation{
		RecruitID:        (*recruitProfiles)[i].RecruitID,
		TeamProfileID:    (*recruitProfiles)[i].ProfileID,
		RecruitProfileID: int((*recruitProfiles)[i].ID),
		WeekID:           timestamp.CollegeWeekID,
	}

	var curr float64 = 0

	var res float64 = 1 // recruitProfiles[i].RecruitingEfficiencyScore

	if (*recruitProfiles)[i].AffinityOneEligible {
		res += affinityBonus
		rpa.ApplyAffinityOne()
	}
	if (*recruitProfiles)[i].AffinityTwoEligible {
		res += affinityBonus
		rpa.ApplyAffinityTwo()
	}

	curr = float64((*recruitProfiles)[i].CurrentWeeksPoints) * res

	if (*recruitProfiles)[i].SpendingCount > 0 {
		streakFormula := affinityBonus * float64((*recruitProfiles)[i].SpendingCount)
		curr *= (1 + streakFormula)
	}

	if (*recruitProfiles)[i].CurrentWeeksPoints < 0 || (*recruitProfiles)[i].CurrentWeeksPoints > float64(pointLimit) {
		curr = 0
		rpa.ApplyCaughtCheating()
	}

	rpa.UpdatePointsSpent((*recruitProfiles)[i].CurrentWeeksPoints, curr)
	(*recruitProfiles)[i].AddCurrentWeekPointsToTotal(curr)
	m.Lock()
	(*recruitProfilePointsMap)[(*recruitProfiles)[i].TeamAbbreviation] += (*recruitProfiles)[i].CurrentWeeksPoints
	m.Unlock()

	// Add RPA to point allocations list
	err := db.Create(&rpa).Error
	if err != nil {
		return fmt.Errorf("could not save point allocation: %v", err)
	}
	return nil
}

func updateTeamRankings(teamRecruitingProfiles []structs.RecruitingTeamProfile, teamMap map[string]*structs.RecruitingTeamProfile, recruitProfilePointsMap map[string]float64, db *gorm.DB) {
	// Update rank system for all teams
	var maxESPNScore float64 = 0
	var minESPNScore float64 = 100000
	var maxRivalsScore float64 = 0
	var minRivalsScore float64 = 100000
	var max247Score float64 = 0
	var min247Score float64 = 100000

	for i := 0; i < len(teamRecruitingProfiles); i++ {

		signedRecruits := GetSignedRecruitsByTeamProfileID(strconv.Itoa(teamRecruitingProfiles[i].TeamID))

		teamRecruitingProfiles[i].UpdateTotalSignedRecruits(len(signedRecruits))

		team247Rank := Get247TeamRanking(teamRecruitingProfiles[i], signedRecruits)
		teamESPNRank := GetESPNTeamRanking(teamRecruitingProfiles[i], signedRecruits)
		teamRivalsRank := GetRivalsTeamRanking(teamRecruitingProfiles[i], signedRecruits)
		if teamESPNRank > maxESPNScore {
			maxESPNScore = teamESPNRank
		}
		if teamESPNRank < minESPNScore {
			minESPNScore = teamESPNRank
		}
		if teamRivalsRank > maxRivalsScore {
			maxRivalsScore = teamRivalsRank
		}
		if teamRivalsRank < minRivalsScore {
			minRivalsScore = teamRivalsRank
		}
		if team247Rank > max247Score {
			max247Score = team247Rank
		}
		if team247Rank < min247Score {
			min247Score = team247Rank
		}

		teamRecruitingProfiles[i].Assign247Rank(team247Rank)
		teamRecruitingProfiles[i].AssignESPNRank(teamESPNRank)
		teamRecruitingProfiles[i].AssignRivalsRank(teamRivalsRank)
	}

	espnDivisor := (maxESPNScore - minESPNScore)
	divisor247 := (max247Score - min247Score)
	rivalsDivisor := (maxRivalsScore - minRivalsScore)

	for _, rp := range teamRecruitingProfiles {
		if recruitProfilePointsMap[rp.TeamAbbreviation] > rp.WeeklyPoints {
			rp.ApplyCaughtCheating()
		}

		var avg float64 = 0
		if espnDivisor > 0 && divisor247 > 0 && rivalsDivisor > 0 {
			distributionESPN := (rp.ESPNScore - minESPNScore) / espnDivisor
			distribution247 := (rp.Rank247Score - min247Score) / divisor247
			distributionRivals := (rp.RivalsScore - minRivalsScore) / rivalsDivisor

			avg = (distributionESPN + distribution247 + distributionRivals)

			rp.AssignCompositeRank(avg)
		}
		rp.ResetSpentPoints()

		// Save TEAM Recruiting Profile
		err := db.Save(&rp).Error
		if err != nil {
			fmt.Println(err.Error())
			log.Fatalf("Could not save timestamp")
		}
		fmt.Println("Saved Rank Scores for Team " + rp.TeamAbbreviation)
	}
}
