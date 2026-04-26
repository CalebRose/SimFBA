package managers

import (
	"math"
	"sort"
	"strconv"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/models"
	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
)

func PostSeasonStatusCleanUp() {
	db := dbprovider.GetInstance().GetDB()

	collegeGames := GetAllCollegeGames()
	collegeTeams := GetAllCollegeTeams()

	// Sort by seasonID asc
	sort.Slice(collegeGames, func(i, j int) bool {
		return collegeGames[i].SeasonID < collegeGames[j].SeasonID
	})

	// seasonIDs := []uint{1, 2, 3, 4, 5, 6}
	seasonIDs := []uint{6}

	for _, seasonID := range seasonIDs {
		seasonIDStr := strconv.Itoa(int(seasonID))
		collegeStandings := repository.FindAllCollegeStandingsRecords(repository.StandingsQuery{SeasonID: seasonIDStr})
		collegeStandingsMap := make(map[uint]structs.CollegeStandings)
		postSeasonStatusMap := make(map[uint]string)
		for _, standing := range collegeStandings {
			collegeStandingsMap[uint(standing.TeamID)] = standing
			postSeasonStatusMap[uint(standing.TeamID)] = "None"
		}

		for _, game := range collegeGames {
			if game.SeasonID != int(seasonID) {
				continue
			}
			if !game.IsBowlGame && !game.IsPlayoffGame && !game.IsNationalChampionship {
				continue
			}
			if game.IsBowlGame {
				postSeasonStatusMap[uint(game.HomeTeamID)] = "Bowl Game"
				postSeasonStatusMap[uint(game.AwayTeamID)] = "Bowl Game"
			}
			if game.IsPlayoffGame && game.Week < 18 {
				postSeasonStatusMap[uint(game.HomeTeamID)] = "Playoffs"
				postSeasonStatusMap[uint(game.AwayTeamID)] = "Playoffs"
			}
			if game.IsPlayoffGame && game.Week >= 18 {
				postSeasonStatusMap[uint(game.HomeTeamID)] = "CFP Semifinals"
				postSeasonStatusMap[uint(game.AwayTeamID)] = "CFP Semifinals"
			}
			if game.IsNationalChampionship {
				if game.HomeTeamWin {
					postSeasonStatusMap[uint(game.HomeTeamID)] = "National Champion"
					postSeasonStatusMap[uint(game.AwayTeamID)] = "National Runner-Up"
				} else {
					postSeasonStatusMap[uint(game.HomeTeamID)] = "National Runner-Up"
					postSeasonStatusMap[uint(game.AwayTeamID)] = "National Champion"
				}
			}
		}

		for _, team := range collegeTeams {
			standing, ok := collegeStandingsMap[uint(team.ID)]
			if !ok {
				continue
			}
			status, ok := postSeasonStatusMap[uint(team.ID)]
			if !ok {
				status = "None"
			}
			standing.PostSeasonStatus = status
			repository.SaveCFBStandingsRecord(standing, db)

		}
	}
}

func UpdateTeamProfileAffinities() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	teamRecruitingProfiles := GetAllTeamRecruitingProfiles()
	teamProfileMap := MakeRecruitTeamProfileMapByTeamID(teamRecruitingProfiles)
	collegeTeams := GetAllCollegeTeams()
	collegeGames := GetAllCollegeGames()
	collegeStandings := repository.FindAllCollegeStandingsRecords(repository.StandingsQuery{})
	collegeGamesByCoachMap := make(map[string][]structs.CollegeGame)
	collegeGamesByTeamIDMap := make(map[uint][]structs.CollegeGame)
	affinitiesMap := make(map[uint]structs.TeamProfileAffinities)

	for _, game := range collegeGames {
		if !game.GameComplete || game.IsSpringGame {
			continue
		}
		if game.HomeTeamCoach != "" && game.HomeTeamCoach != "AI" {
			collegeGamesByCoachMap[game.HomeTeamCoach] = append(collegeGamesByCoachMap[game.HomeTeamCoach], game)
		}
		if game.AwayTeamCoach != "" && game.AwayTeamCoach != "AI" {
			collegeGamesByCoachMap[game.AwayTeamCoach] = append(collegeGamesByCoachMap[game.AwayTeamCoach], game)
		}
		if game.HomeTeamID > 0 && game.AwayTeamID > 0 {
			collegeGamesByTeamIDMap[uint(game.HomeTeamID)] = append(collegeGamesByTeamIDMap[uint(game.HomeTeamID)], game)
			collegeGamesByTeamIDMap[uint(game.AwayTeamID)] = append(collegeGamesByTeamIDMap[uint(game.AwayTeamID)], game)
		}

	}

	// Record by previous season since this will run in the season migration
	seasonID := ts.CollegeSeasonID - 1
	baseProgramDevSeasonID := seasonID - 4

	// Map all historical standings by team
	collegeStandingsMap := make(map[uint][]structs.CollegeStandings)
	for _, standing := range collegeStandings {
		collegeStandingsMap[uint(standing.TeamID)] = append(collegeStandingsMap[uint(standing.TeamID)], standing)
	}

	// Iterate by team
	// get all standings by team.ID
	// iterate by game
	for _, team := range collegeTeams {
		totalWins := 0
		totalLosses := 0
		totalCoachWins := 0
		totalCoachLosses := 0
		seasonMomentumWins := 0
		seasonMomentumLosses := 0
		homeWins := 0
		homeLosses := 0
		programPrestige := 5
		professionalPrestige := 5
		facilities := 5
		academics := 5
		campusLife := 5
		confPrestige := 5
		mediaSpotlight := 5
		teamStandings := collegeStandingsMap[uint(team.ID)]
		collegeGamesByTeam := collegeGamesByTeamIDMap[uint(team.ID)]
		collegeGamesByCoach := collegeGamesByCoachMap[team.Coach]

		// Iterate and track home wins and losses
		for _, game := range collegeGamesByTeam {
			if game.HomeTeamID == int(team.ID) {
				if !game.IsNeutral && game.HomeTeamScore > game.AwayTeamScore {
					homeWins++
				} else {
					homeLosses++
				}
				if game.HomeTeamScore > game.AwayTeamScore {
					totalWins++
				} else {
					totalLosses++
				}
			} else if game.AwayTeamID == int(team.ID) {
				if game.AwayTeamScore > game.HomeTeamScore {
					totalWins++
				} else {
					totalLosses++
				}
			}
		}

		tradPercentage := float64(totalWins) / float64(totalWins+totalLosses)
		atmospherePct := float64(homeWins) / float64(homeWins+homeLosses)

		// Iterate by historic standings for program development & seasonMomentum
		for _, standing := range teamStandings {
			// Season momentum
			if standing.SeasonID == ts.CollegeSeasonID {
				seasonMomentumWins = standing.TotalWins
				seasonMomentumLosses = standing.TotalLosses
			}
			if standing.SeasonID < baseProgramDevSeasonID {
				continue
			}

			if standing.TotalWins == 0 {
				programPrestige -= 2
			}
			record := float64(standing.TotalWins) / float64(standing.TotalWins+standing.TotalLosses)
			if record < 0.25 {
				programPrestige -= 2
			} else if record < 0.5 {
				programPrestige -= 1
			} else if record >= 0.75 {
				programPrestige += 1
			}
			if standing.TotalLosses == 0 {
				programPrestige += 2
			}

			postSeason := standing.PostSeasonStatus
			postSeasonFlag := 0
			if postSeason == "Bowl Game" {
				postSeasonFlag = 1
			}
			if postSeason == "Playoffs" || postSeason == "National Runner-Up" || postSeason == "National Champion" || postSeason == "CFP Semifinals" {
				postSeasonFlag = 1
			}
			if postSeason == "National Runner-Up" || postSeason == "National Champion" || postSeason == "CFP Semifinals" {
				postSeasonFlag = 2
			}
			if postSeason == "National Champion" {
				postSeasonFlag = 3
			}
			programPrestige += postSeasonFlag

			if programPrestige < 1 {
				programPrestige = 1
			}
			if programPrestige > 10 {
				programPrestige = 10
			}
		}

		seasonMomentumPct := float64(seasonMomentumWins) / float64(seasonMomentumWins+seasonMomentumLosses)

		// Iterate through games by Coach
		if team.Coach != "" && team.Coach != "AI" {
			for _, game := range collegeGamesByCoach {
				if game.HomeTeamCoach == team.Coach {
					if game.HomeTeamScore > game.AwayTeamScore {
						totalCoachWins += 1
					} else {
						totalCoachLosses += 1
					}
				}
				if game.AwayTeamCoach == team.Coach {
					if game.AwayTeamScore > game.HomeTeamScore {
						totalCoachWins += 1
					} else {
						totalCoachLosses += 1
					}
				}
			}
		}

		coachPct := float64(0.0)
		if (totalCoachWins + totalCoachLosses) > 0 {
			coachPct = float64(totalCoachWins) / float64(totalCoachWins+totalCoachLosses)
		}

		newCoachRating := uint8(coachPct * 10)
		if team.Coach == "" || team.Coach == "AI" || len(collegeGamesByCoach) == 0 {
			// Set default to 5
			newCoachRating = 5
		}

		seasonMomentum := uint8(seasonMomentumPct * 10)
		if seasonMomentum < 1 {
			seasonMomentum = 1
		} else if seasonMomentum > 10 {
			seasonMomentum = 10
		}

		affinitiesMap[uint(team.ID)] = structs.TeamProfileAffinities{
			ProgramPrestige:      uint8(programPrestige),
			ProfessionalPrestige: uint8(professionalPrestige),
			Traditions:           uint8(tradPercentage * 10),
			Facilities:           uint8(facilities),
			Atmosphere:           uint8(atmospherePct * 10),
			Academics:            uint8(academics),
			CampusLife:           uint8(campusLife),
			ConferencePrestige:   uint8(confPrestige),
			CoachRating:          newCoachRating,
			SeasonMomentum:       seasonMomentum,
			MediaSpotlight:       uint8(mediaSpotlight),
		}
	}

	// -------------------------------------------------------------------------
	// Conference Prestige: recency-weighted 4-year average playoff metrics
	// combined with current-year bowl game rates. All metrics are computed as
	// per-team averages so that larger conferences are not advantaged over
	// smaller ones.
	//
	// Score components (all rates in [0, 1]):
	//   Playoff appearance rate (4-yr weighted)   × 3  → up to +3
	//   Playoff advancement rate (4-yr weighted)  × 2  → up to +2
	//     (advancement = won at least one playoff game, i.e. CFP Semis or better)
	//   Current-year bowl appearance rate          × 2  → up to +2
	//   Current-year bowl win rate (of bowl teams) × 2  → up to +2
	// Base score = 1; total max = 10.
	// Recency weights: current season = 4, one year ago = 3, two = 2, three = 1.
	// Weights are normalised by the sum of weights actually used, so conferences
	// that have existed for fewer than 4 seasons are not penalised unfairly.
	// -------------------------------------------------------------------------

	// Build the set of teams that won at least one postseason game this season.
	bowlWinsSet := make(map[uint]bool)
	for _, game := range collegeGames {
		if !game.GameComplete {
			continue
		}
		if game.SeasonID != int(seasonID) {
			continue
		}
		if !game.IsBowlGame && !game.IsPlayoffGame && !game.IsNationalChampionship {
			continue
		}
		if game.HomeTeamWin {
			bowlWinsSet[uint(game.HomeTeamID)] = true
		} else {
			bowlWinsSet[uint(game.AwayTeamID)] = true
		}
	}

	// Organise all loaded standings by (conferenceID, seasonID) for O(1) lookup.
	type confSeasonKey struct {
		confID   uint
		seasonID uint
	}
	confSeasonStandingsMap := make(map[confSeasonKey][]structs.CollegeStandings)
	for _, standing := range collegeStandings {
		key := confSeasonKey{
			confID:   uint(standing.ConferenceID),
			seasonID: uint(standing.SeasonID),
		}
		confSeasonStandingsMap[key] = append(confSeasonStandingsMap[key], standing)
	}

	// Collect unique conference IDs from active teams.
	confIDSet := make(map[uint]bool)
	for _, team := range collegeTeams {
		confIDSet[uint(team.ConferenceID)] = true
	}

	const playoffWindow = 4

	conferencePrestigeMap := make(map[uint]int)
	confRawScores := make(map[uint]float64)
	for confID := range confIDSet {
		weightedPlayoffAppRate := 0.0
		weightedPlayoffWinRate := 0.0
		usedWeight := 0.0

		for offset := 0; offset < playoffWindow; offset++ {
			targetSeason := int(seasonID) - offset
			if targetSeason < 1 {
				continue
			}
			// Newest season (offset=0) gets weight 4; oldest (offset=3) gets weight 1.
			weight := float64(playoffWindow - offset)

			key := confSeasonKey{confID: confID, seasonID: uint(targetSeason)}
			seasonStandings := confSeasonStandingsMap[key]
			if len(seasonStandings) == 0 {
				continue
			}
			total := float64(len(seasonStandings))

			playoffAppearances := 0
			playoffWins := 0
			for _, s := range seasonStandings {
				ps := s.PostSeasonStatus
				if ps == "Playoffs" || ps == "CFP Semifinals" ||
					ps == "National Runner-Up" || ps == "National Champion" {
					playoffAppearances++
				}
				// Advancement = won at least one playoff game.
				if ps == "CFP Semifinals" || ps == "National Runner-Up" || ps == "National Champion" {
					playoffWins++
				}
			}

			weightedPlayoffAppRate += weight * float64(playoffAppearances) / total
			weightedPlayoffWinRate += weight * float64(playoffWins) / total
			usedWeight += weight
		}

		if usedWeight > 0 {
			weightedPlayoffAppRate /= usedWeight
			weightedPlayoffWinRate /= usedWeight
		}

		// Current-year bowl metrics.
		currentKey := confSeasonKey{confID: confID, seasonID: uint(seasonID)}
		currentStandings := confSeasonStandingsMap[currentKey]
		currentTotal := float64(len(currentStandings))

		bowlAppearances := 0
		bowlWins := 0
		for _, s := range currentStandings {
			ps := s.PostSeasonStatus
			if ps != "" && ps != "None" {
				bowlAppearances++
			}
			if bowlWinsSet[uint(s.TeamID)] {
				bowlWins++
			}
		}

		// Bowl appearance rate = fraction of conference teams that made any bowl.
		// Bowl win rate = fraction of bowl participants that won their game.
		bowlAppearanceRate := 0.0
		bowlWinRate := 0.0
		if currentTotal > 0 {
			bowlAppearanceRate = float64(bowlAppearances) / currentTotal
		}
		if bowlAppearances > 0 {
			bowlWinRate = float64(bowlWins) / float64(bowlAppearances)
		}

		rawScore := 1.0 +
			weightedPlayoffAppRate*3.0 +
			weightedPlayoffWinRate*2.0 +
			bowlAppearanceRate*2.0 +
			bowlWinRate*2.0

		confRawScores[confID] = rawScore
	}

	// -------------------------------------------------------------------------
	// Rank-distribute prestige so the full 1–10 scale is always used.
	// Independents (conf IDs 13 and 22) are fixed at 5 (the median).
	// FBS conferences occupy the top of the ranking; FCS conferences occupy
	// the bottom. Within each tier, conferences are sorted by raw score
	// descending. The combined ordered list is then linearly interpolated
	// so that rank 1 → prestige 10 and rank 25 → prestige 1.
	// -------------------------------------------------------------------------
	independentConfIDs := map[uint]bool{13: true, 22: true}
	fcsConfIDs := map[uint]bool{
		14: true, 15: true, 16: true, 17: true, 18: true, 19: true,
		20: true, 21: true, 23: true, 24: true, 25: true, 26: true, 27: true,
	}

	fbsConfs := make([]uint, 0)
	fcsConfs := make([]uint, 0)
	for confID := range confIDSet {
		if independentConfIDs[confID] {
			continue
		}
		if fcsConfIDs[confID] {
			fcsConfs = append(fcsConfs, confID)
		} else {
			fbsConfs = append(fbsConfs, confID)
		}
	}
	sort.Slice(fbsConfs, func(i, j int) bool {
		return confRawScores[fbsConfs[i]] > confRawScores[fbsConfs[j]]
	})
	sort.Slice(fcsConfs, func(i, j int) bool {
		return confRawScores[fcsConfs[i]] > confRawScores[fcsConfs[j]]
	})

	// FBS conferences occupy the top prestige slots; FCS follow.
	nonIndepConfs := append(fbsConfs, fcsConfs...)

	n := len(nonIndepConfs)
	for rank, confID := range nonIndepConfs {
		prestige := 10
		if n > 1 {
			// Linear interpolation: rank 0 → 10, rank n-1 → 1
			prestige = int(math.Round(10.0 - float64(rank)*9.0/float64(n-1)))
		}
		conferencePrestigeMap[confID] = prestige
	}

	// Independents always sit at the median.
	conferencePrestigeMap[13] = 5
	conferencePrestigeMap[22] = 5

	for _, team := range collegeTeams {
		teamProfile := teamProfileMap[team.ID]
		affinities := affinitiesMap[team.ID]
		conferencePrestige := conferencePrestigeMap[uint(team.ConferenceID)]
		if conferencePrestige < 1 {
			conferencePrestige = 1
		} else if conferencePrestige > 10 {
			conferencePrestige = 10
		}

		if team.ConferenceID > 13 && conferencePrestige < 1 {
			conferencePrestige = 5
		}
		affinities.ConferencePrestige = uint8(conferencePrestige)
		teamProfile.UpdateTeamProfileAffinities(affinities)
		repository.SaveRecruitingTeamProfile(teamProfile, db)
	}
}

func RecruitingAndTransferPortalCleanUp() {
	db := dbprovider.GetInstance().GetDB()
	db.Model(&models.NFLWarRoom{}).Where("id > ?", 0).Update("spent_points", 0)

	// Clear Transfer Portal Profiles Table
	db.Delete(&structs.TransferPortalProfile{}, "id > ?", 0)

	// Clear Recruiting Profiles Table
	db.Delete(&structs.RecruitPlayerProfile{}, "id > ?", 0)

	// Clear Transfer Profiles Table
	db.Delete(&structs.TransferPortalProfile{}, "id > ?", 0)

	// Clear NFL Scouting Boards
	db.Delete(&models.ScoutingProfile{}, "id > ?", 0)
}

func FreeAgencyCleanUp() {
	db := dbprovider.GetInstance().GetDB()
	db.Delete(&structs.FreeAgencyOffer{}, "id > ?", 0)
	db.Delete(&structs.NFLExtensionOffer{}, "id > ?", 0)
}
