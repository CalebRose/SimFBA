package managers

import (
	"fmt"
	"sort"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
)

func BaseGenerateCFBSchedule(testTheData bool) {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	collegeTeams := GetAllCollegeTeams()
	collegeTeamMap := MakeCollegeTeamMap(collegeTeams)
	collegeTeamsByConference := MakeCollegeTeamsByConference(collegeTeams)
	allHistoricGames := GetAllCollegeGames()
	playCountMap, _, lastHomeMap := BuildScheduleHistoryMaps(allHistoricGames)

	stadium := GetAllStadiums()
	stadiumMap := MakeStadiumMapByTeamID(stadium, true)
	stadiumMapByID := MakeStadiumMapByID(stadium)

	rivalries := GetAllRivalries()
	rivalryMap := MakeHistoricRivalriesMapByTeamID(rivalries)

	gamesPlayedAgainstOpponentsMap := make(map[uint]map[uint]bool) // teamID -> opponentID -> bool
	gamesPlayedByWeekMap := make(map[uint]map[uint]bool)           // teamID -> week -> bool
	// homeCountSeedMap tracks how many times each team was the home team in the
	// OOC rivalry pass. This is forwarded to conference generators so they can
	// seed their homecountMap with pre-existing home game counts and avoid
	// accidentally over-scheduling home games for teams that are already "home-heavy".
	homeCountSeedMap := make(map[uint]int)

	collegeGamesUpload := []structs.CollegeGame{}

	// Iterate through rivalry games first to ensure that annual games will be generated
	for _, rivalry := range rivalries {
		if !rivalry.IsAnnualRivalry {
			continue
		}
		// Skip annual rivalries that have no preferred week — the conference
		// scheduler will handle them as floating locked pairs.
		if rivalry.PreferredWeek == 0 {
			continue
		}
		homeTeam := collegeTeamMap[rivalry.TeamOneID]
		awayTeam := collegeTeamMap[rivalry.TeamTwoID]
		if homeTeam.ID == 0 || awayTeam.ID == 0 {
			continue
		}
		if ts.Season%2 == 0 {
			homeTeam, awayTeam = awayTeam, homeTeam
		}
		preferredWeek := rivalry.PreferredWeek
		// Notre Dame vs USC
		if rivalry.ID == 233 {
			// If it's an odd year, USC and Notre Dame play in week 14 at USC
			if ts.Season%2 != 0 {
				preferredWeek = 14
			}
		}
		// Stanford vs USC
		if rivalry.ID == 274 {
			// If it's an odd year, Stanford and Notre Dame play in week 4 at Notre Dame
			if ts.Season%2 != 0 {
				preferredWeek = 4
			}
		}
		rivalryGame := MakeCollegeGameRecord(homeTeam, awayTeam, uint(preferredWeek), uint(ts.CollegeSeasonID), stadiumMap, stadiumMapByID, rivalryMap)

		collegeGamesUpload = append(collegeGamesUpload, rivalryGame)
		// Update nested map data structures
		if gamesPlayedAgainstOpponentsMap[homeTeam.ID] == nil {
			gamesPlayedAgainstOpponentsMap[homeTeam.ID] = make(map[uint]bool)
		}
		if gamesPlayedAgainstOpponentsMap[awayTeam.ID] == nil {
			gamesPlayedAgainstOpponentsMap[awayTeam.ID] = make(map[uint]bool)
		}
		gamesPlayedAgainstOpponentsMap[homeTeam.ID][awayTeam.ID] = true
		gamesPlayedAgainstOpponentsMap[awayTeam.ID][homeTeam.ID] = true
		if gamesPlayedByWeekMap[homeTeam.ID] == nil {
			gamesPlayedByWeekMap[homeTeam.ID] = make(map[uint]bool)
		}
		if gamesPlayedByWeekMap[awayTeam.ID] == nil {
			gamesPlayedByWeekMap[awayTeam.ID] = make(map[uint]bool)
		}
		// Use the week NUMBER (not WeekID) so conference schedulers can see these slots as occupied.
		gamesPlayedByWeekMap[homeTeam.ID][uint(rivalryGame.Week)] = true
		gamesPlayedByWeekMap[awayTeam.ID][uint(rivalryGame.Week)] = true
		// Track home team for seeding conference generators' homecountMap.
		homeCountSeedMap[homeTeam.ID]++
	}

	// Generate Conference Schedules

	conferenceGames := BaseGenerateCFBConferenceSchedule(collegeTeamMap, collegeTeamsByConference, stadiumMap, stadiumMapByID, rivalryMap, gamesPlayedAgainstOpponentsMap, gamesPlayedByWeekMap, playCountMap, lastHomeMap, homeCountSeedMap, ts)
	collegeGamesUpload = append(collegeGamesUpload, conferenceGames...)

	// Sort by week ID for easier reading in CSV output
	sort.Slice(collegeGamesUpload, func(i, j int) bool {
		// Sort by week ID and then timeslot and then home team ID for consistent ordering
		if collegeGamesUpload[i].WeekID != collegeGamesUpload[j].WeekID {
			return collegeGamesUpload[i].WeekID < collegeGamesUpload[j].WeekID
		}
		if collegeGamesUpload[i].TimeSlot != collegeGamesUpload[j].TimeSlot {
			return collegeGamesUpload[i].TimeSlot < collegeGamesUpload[j].TimeSlot
		}
		return collegeGamesUpload[i].HomeTeamID < collegeGamesUpload[j].HomeTeamID
	})

	// Conduct Upload
	if testTheData {
		exportPath := fmt.Sprintf("schedule_preview_season_%d.csv", ts.Season)
		if err := ExportScheduleToCSV(collegeGamesUpload, exportPath); err != nil {
			fmt.Printf("CSV export failed: %v\n", err)
		} else {
			fmt.Printf("Schedule preview exported to %s (%d total games)\n", exportPath, len(collegeGamesUpload))
		}
	} else {
		repository.CreateCFBGameRecordsBatch(db, collegeGamesUpload, 200)
	}
}

// BaseGenerateCFBConferenceSchedule generates conference games for every conference,
// processing one conference at a time (partition approach). After each conference's
// primary generation, a validation pass checks that every team has the expected
// number of conference games and attempts to schedule any missing matchups.
func BaseGenerateCFBConferenceSchedule(collegeTeamMap map[uint]structs.CollegeTeam, collegeTeamsByConference map[int][]structs.CollegeTeam, stadiumMap map[uint]structs.Stadium, stadiumMapByID map[uint]structs.Stadium, rivalryMap map[uint][]structs.CollegeRival, gamesPlayedAgainstOpponentsMap map[uint]map[uint]bool, gamesPlayedByWeekMap map[uint]map[uint]bool, playCountMap map[SchedulerHistoryKey]int, lastHomeMap map[uint]map[uint]bool, homeCountSeedMap map[uint]int, ts structs.Timestamp) []structs.CollegeGame {
	allConferenceGames := []structs.CollegeGame{}
	// IDs 22 and 13 are excluded because these are independent conferences
	// 3=ACC, 4=Big Ten, 5=Big 12, 6=Pac-12, 7=SEC, 8=AAC, 9=C-USA, 10=MAC,
	// 11=MWC, 12=Sun Belt, 14=MVFC, 15=Patriot, 16=Big Sky, 17=Ivy,
	// 18=SoCon, 19=SWAC, 20=Big South-OVC, 21=CAA, 23=MEAC, 24=NEC,
	// 25=Pioneer, 26=Southland, 27=UAC
	conferenceIDs := []int{3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 14, 15, 16, 17, 18, 19, 20, 21, 23, 24, 25, 26, 27}

	var allMissing []MissingGame

	// smallConferenceIDs is the set of FCS/small-conference IDs that use the
	// generic round-robin generator and are eligible for multi-seed retry.
	smallConferenceIDs := map[int]bool{
		3:  true,  // ACC — generator returns empty slice on failure to trigger reseed
		4:  true,  // Big Ten — generator returns empty slice on failure to trigger reseed
		5:  false, // Big 12 — generator returns empty slice on failure to trigger reseed
		7:  true,  // SEC — generator returns empty slice on failure to trigger reseed
		8:  true,  // AAC — generator returns empty slice on failure to trigger reseed
		9:  true,  // CUSA — generator returns empty slice on failure to trigger reseed
		10: true,
		11: true, 14: true, 15: true, 16: true, 17: true, 18: true, 19: true,
		20: true, 21: true, 23: true, 24: true, 25: true, 26: true, 27: true,
	}
	const maxSmallConfRetries = 10000

	for _, conferenceID := range conferenceIDs {
		teams := collegeTeamsByConference[conferenceID]
		if len(teams) == 0 {
			continue
		}

		expected, ok := ConferenceExpectedGames[conferenceID]
		if !ok {
			expected = 8
		}

		if smallConferenceIDs[conferenceID] {
			// --- Small conference: try up to maxSmallConfRetries seeds ---
			// Each seed shuffles the pair-processing order differently, changing
			// which weeks are claimed first and potentially resolving edge cases
			// where certain team pairs cannot find a shared free week.
			var bestGames []structs.CollegeGame
			var bestMissing []MissingGame
			var bestOppMap map[uint]map[uint]bool
			var bestWeekMap map[uint]map[uint]bool
			bestMissingCount := 999999

			for seed := uint(0); seed < maxSmallConfRetries; seed++ {
				oppClone := deepCopyBoolMap(gamesPlayedAgainstOpponentsMap)
				weekClone := deepCopyBoolMap(gamesPlayedByWeekMap)

				cGames := dispatchConferenceGenerator(conferenceID, seed, teams, stadiumMap, stadiumMapByID, rivalryMap, oppClone, weekClone, playCountMap, lastHomeMap, homeCountSeedMap, ts)
				if len(cGames) == 0 {
					// Generator returned empty slice — signals a failed attempt (e.g. ACC 8-game
					// asymmetry unresolvable, or Big Ten week-exhaustion). Decrement seed so the
					// outer loop increment lands on the same seed value, effectively retrying
					// with a fresh rand.Shuffle ordering without burning a retry slot.
					seed--
					continue
				}
				rGames, missing := ValidateAndRescueConference(
					conferenceID, expected, teams,
					stadiumMap, stadiumMapByID, rivalryMap,
					oppClone, weekClone,
					uint(ts.CollegeSeasonID), ts.Season,
				)
				total := append(cGames, rGames...)
				if len(missing) < bestMissingCount {
					bestMissingCount = len(missing)
					bestGames = total
					bestMissing = missing
					bestOppMap = oppClone
					bestWeekMap = weekClone
				}
				if bestMissingCount == 0 {
					break
				}
			}

			// Commit the winning attempt's maps back into the shared maps.
			for k := range gamesPlayedAgainstOpponentsMap {
				delete(gamesPlayedAgainstOpponentsMap, k)
			}
			for k, v := range bestOppMap {
				gamesPlayedAgainstOpponentsMap[k] = v
			}
			for k := range gamesPlayedByWeekMap {
				delete(gamesPlayedByWeekMap, k)
			}
			for k, v := range bestWeekMap {
				gamesPlayedByWeekMap[k] = v
			}

			allConferenceGames = append(allConferenceGames, bestGames...)
			logConferenceGameCountIssues(conferenceID, expected, teams, gamesPlayedAgainstOpponentsMap)
			if len(bestMissing) > 0 {
				fmt.Printf("SCHEDULE RESCUE: conf %d — %d unresolvable game(s) after %d attempts\n", conferenceID, len(bestMissing), maxSmallConfRetries)
				allMissing = append(allMissing, bestMissing...)
			}
		} else {
			// --- Large conference: single pass ---
			confGames := dispatchConferenceGenerator(conferenceID, 0, teams, stadiumMap, stadiumMapByID, rivalryMap, gamesPlayedAgainstOpponentsMap, gamesPlayedByWeekMap, playCountMap, lastHomeMap, homeCountSeedMap, ts)
			allConferenceGames = append(allConferenceGames, confGames...)

			rescueGames, missing := ValidateAndRescueConference(
				conferenceID, expected, teams,
				stadiumMap, stadiumMapByID, rivalryMap,
				gamesPlayedAgainstOpponentsMap, gamesPlayedByWeekMap,
				uint(ts.CollegeSeasonID), ts.Season,
			)
			if len(rescueGames) > 0 {
				fmt.Printf("SCHEDULE RESCUE: conf %d — added %d game(s) in rescue pass\n", conferenceID, len(rescueGames))
				allConferenceGames = append(allConferenceGames, rescueGames...)
			}
			allMissing = append(allMissing, missing...)
			logConferenceGameCountIssues(conferenceID, expected, teams, gamesPlayedAgainstOpponentsMap)
		}
	}

	if len(allMissing) > 0 {
		fmt.Printf("SCHEDULE VALIDATION: %d unresolvable matchup(s) after rescue:\n", len(allMissing))
		for _, m := range allMissing {
			fmt.Printf("  Conf %d: team %d vs team %d — %s\n", m.ConferenceID, m.TeamAID, m.TeamBID, m.Reason)
		}
	}

	// Rebalance home/away within all conference games so that every team ends
	// up with |home − away| ≤ 2.  Neutral-site and rivalry games are left
	// untouched.  Swapped pairs are rebuilt with the new home team's stadium
	// and timeslot.
	// Rerun a few more times to achieve better balance, as each pass may only fix some of the imbalances.
	// Commenting out for now as it can cause issues with certain conferences (e.g. Big Ten) where the conference generator is already struggling to find any valid schedule, and the rebalance pass can break some of the delicate balance achieved by the conference generator's retry loop.
	// for i := 0; i < 25; i++ {
	// 	allConferenceGames = rebalanceConferenceHomeAway(allConferenceGames, collegeTeamMap, stadiumMap, stadiumMapByID, rivalryMap)
	// }

	return allConferenceGames
}

// dispatchConferenceGenerator routes a conference ID to its dedicated scheduler.
// retrySeed is forwarded to small-conference generators to vary pair ordering on retry.
func dispatchConferenceGenerator(conferenceID int, retrySeed uint, teams []structs.CollegeTeam, stadiumMap map[uint]structs.Stadium, stadiumMapByID map[uint]structs.Stadium, rivalryMap map[uint][]structs.CollegeRival, gamesPlayedAgainstOpponentsMap map[uint]map[uint]bool, gamesPlayedByWeekMap map[uint]map[uint]bool, playCountMap map[SchedulerHistoryKey]int, lastHomeMap map[uint]map[uint]bool, homeCountSeedMap map[uint]int, ts structs.Timestamp) []structs.CollegeGame {
	switch conferenceID {
	case 3:
		return GenerateACCSchedule(teams, stadiumMap, stadiumMapByID, rivalryMap, gamesPlayedAgainstOpponentsMap, gamesPlayedByWeekMap, playCountMap, lastHomeMap, homeCountSeedMap, ts)
	case 4:
		return GenerateBigTenSchedule(teams, stadiumMap, stadiumMapByID, rivalryMap, gamesPlayedAgainstOpponentsMap, gamesPlayedByWeekMap, playCountMap, lastHomeMap, homeCountSeedMap, ts)
	case 5:
		return GenerateBigTwelveSchedule(teams, stadiumMap, stadiumMapByID, rivalryMap, gamesPlayedAgainstOpponentsMap, gamesPlayedByWeekMap, playCountMap, lastHomeMap, homeCountSeedMap, ts)
	case 6:
		return GeneratePacTwelveSchedule(teams, stadiumMap, stadiumMapByID, rivalryMap, gamesPlayedAgainstOpponentsMap, gamesPlayedByWeekMap, playCountMap, lastHomeMap, homeCountSeedMap, ts)
	case 7:
		return GenerateSECSchedule(teams, stadiumMap, stadiumMapByID, rivalryMap, gamesPlayedAgainstOpponentsMap, gamesPlayedByWeekMap, playCountMap, lastHomeMap, homeCountSeedMap, ts)
	case 8:
		return GenerateAACSchedule(teams, stadiumMap, stadiumMapByID, rivalryMap, gamesPlayedAgainstOpponentsMap, gamesPlayedByWeekMap, playCountMap, lastHomeMap, homeCountSeedMap, ts)
	case 9:
		return GenerateCUSASchedule(teams, stadiumMap, stadiumMapByID, rivalryMap, gamesPlayedAgainstOpponentsMap, gamesPlayedByWeekMap, playCountMap, lastHomeMap, homeCountSeedMap, ts)
	// case 10:
	// 	return GenerateMACSchedule(teams, stadiumMap, stadiumMapByID, rivalryMap, gamesPlayedAgainstOpponentsMap, gamesPlayedByWeekMap, playCountMap, lastHomeMap, homeCountSeedMap, ts)
	case 11:
		return GenerateMWCSchedule(teams, stadiumMap, stadiumMapByID, rivalryMap, gamesPlayedAgainstOpponentsMap, gamesPlayedByWeekMap, playCountMap, lastHomeMap, homeCountSeedMap, ts)
	case 12:
		return GenerateSunBeltSchedule(teams, stadiumMap, stadiumMapByID, rivalryMap, gamesPlayedAgainstOpponentsMap, gamesPlayedByWeekMap, playCountMap, lastHomeMap, homeCountSeedMap, ts)
	default:
		return GenerateSmallConferenceSchedule(conferenceID, retrySeed, teams, stadiumMap, stadiumMapByID, rivalryMap, gamesPlayedAgainstOpponentsMap, gamesPlayedByWeekMap, playCountMap, lastHomeMap, homeCountSeedMap, ts)
	}
}

// logConferenceGameCountIssues prints a per-team warning for any team in the
// conference that ends up SHORT or OVER the expected conference game count.
// Call this after the primary generator and rescue pass have both been committed.
func logConferenceGameCountIssues(conferenceID int, expected int, teams []structs.CollegeTeam, gamesPlayedAgainstOpponentsMap map[uint]map[uint]bool) {
	confSet := make(map[uint]bool, len(teams))
	for _, t := range teams {
		confSet[t.ID] = true
	}
	hasIssue := false
	for _, t := range teams {
		count := 0
		for oppID := range gamesPlayedAgainstOpponentsMap[t.ID] {
			if confSet[oppID] {
				count++
			}
		}
		if count < expected {
			fmt.Printf("  CONF %d SHORT: %s (id=%d) %d/%d conf games\n", conferenceID, t.TeamName, t.ID, count, expected)
			hasIssue = true
		} else if count > expected {
			fmt.Printf("  CONF %d OVER:  %s (id=%d) %d/%d conf games\n", conferenceID, t.TeamName, t.ID, count, expected)
			hasIssue = true
		}
	}
	if hasIssue {
		fmt.Printf("  CONF %d: expected %d conf games per team\n", conferenceID, expected)
	}
}

// Helper functions
func MakeCollegeGameRecord(homeTeam structs.CollegeTeam, awayTeam structs.CollegeTeam, week uint, seasonID uint, stadiumMap map[uint]structs.Stadium, stadiumMapByID map[uint]structs.Stadium, rivalryMap map[uint][]structs.CollegeRival) structs.CollegeGame {
	homeStadium := stadiumMap[homeTeam.ID]
	rivalry := rivalryMap[homeTeam.ID]
	weekPlayed := week
	isRivalry := false
	isNeutral := false
	stadiumID := homeStadium.ID
	timeSlot := ""
	for _, r := range rivalry {
		if (r.TeamOneID == homeTeam.ID && r.TeamTwoID == awayTeam.ID) || (r.TeamOneID == awayTeam.ID && r.TeamTwoID == homeTeam.ID) {
			isRivalry = true
			if r.PreferredWeek > 0 {
				weekPlayed = uint(r.PreferredWeek)
			}
			if len(r.Timeslot) > 0 {
				timeSlot = r.Timeslot
			}
			if r.IsNeutralSite {
				isNeutral = true
				stadiumID = r.StadiumID
				homeStadium = stadiumMapByID[r.StadiumID]
			}
			break
		}
	}

	if timeSlot == "" {
		timeSlot = util.GetTimeslot(homeTeam.State, uint(homeTeam.ConferenceID))
	}

	isDivision := false
	if homeTeam.DivisionID > 0 && homeTeam.DivisionID == awayTeam.DivisionID {
		isDivision = true
	}

	weekID := util.GetWeekID(seasonID, weekPlayed)

	return structs.CollegeGame{
		HomeTeamID:    int(homeTeam.ID),
		HomeTeam:      homeTeam.TeamAbbr,
		AwayTeamID:    int(awayTeam.ID),
		AwayTeam:      awayTeam.TeamAbbr,
		Week:          int(weekPlayed),
		SeasonID:      int(seasonID),
		WeekID:        int(weekID),
		StadiumID:     stadiumID,
		Stadium:       homeStadium.StadiumName,
		City:          homeStadium.City,
		State:         homeStadium.State,
		IsRivalryGame: isRivalry,
		IsConference:  homeTeam.ConferenceID == awayTeam.ConferenceID,
		IsDivisional:  isDivision,
		TimeSlot:      timeSlot,
		Region:        homeStadium.Region,
		IsNeutral:     isNeutral,
	}
}
