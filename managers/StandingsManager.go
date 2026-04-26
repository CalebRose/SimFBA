package managers

import (
	"log"
	"sort"
	"strconv"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/models"
	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
)

// GetStandingsByConferenceIDAndSeasonID
func GetStandingsByConferenceIDAndSeasonID(conferenceID string, seasonID string) []structs.CollegeStandings {
	var standings []structs.CollegeStandings
	db := dbprovider.GetInstance().GetDB()
	err := db.Where("conference_id = ? AND season_id = ?", conferenceID, seasonID).Order("conference_losses asc").Order("conference_wins desc").
		Order("total_losses asc").Order("total_wins desc").
		Find(&standings).Error
	if err != nil {
		log.Fatal(err)
	}
	return standings
}

func GetNFLStandingsBySeasonID(seasonID string) []structs.NFLStandings {
	var standings []structs.NFLStandings
	db := dbprovider.GetInstance().GetDB()
	err := db.Where("season_id = ?", seasonID).Order("total_losses asc").Order("total_ties asc").Order("total_wins desc").
		Find(&standings).Error
	if err != nil {
		log.Fatal(err)
	}
	return standings
}

func GetNFLStandingsByTeamIDAndSeasonID(teamID string, seasonID string) structs.NFLStandings {
	var standings structs.NFLStandings
	db := dbprovider.GetInstance().GetDB()
	err := db.Where("team_id = ? AND season_id = ?", teamID, seasonID).Order("division_losses asc").Order("division_ties asc").Order("division_wins desc").
		Order("total_losses asc").Order("total_ties asc").Order("total_wins desc").
		Find(&standings).Error
	if err != nil {
		log.Fatal(err)
	}
	return standings
}

func GetNFLStandingsByDivisionIDAndSeasonID(divisionID string, seasonID string) []structs.NFLStandings {
	var standings []structs.NFLStandings
	db := dbprovider.GetInstance().GetDB()
	err := db.Where("division_id = ? AND season_id = ?", divisionID, seasonID).Order("division_losses asc").Order("division_ties asc").Order("division_wins desc").
		Order("total_losses asc").Order("total_ties asc").Order("total_wins desc").
		Find(&standings).Error
	if err != nil {
		log.Fatal(err)
	}
	return standings
}

// GetHistoricalRecordsByTeamID
func GetHistoricalRecordsByTeamID(TeamID string) models.TeamRecordResponse {
	tsChn := make(chan structs.Timestamp)

	go func() {
		ts := GetTimestamp()
		tsChn <- ts
	}()

	timestamp := <-tsChn
	close(tsChn)

	season := strconv.Itoa(timestamp.Season)

	historicGames := GetCollegeGamesByTeamId(TeamID)
	var conferenceChampionships []string
	var divisionTitles []string
	var nationalChampionships []string
	overallWins := 0
	overallLosses := 0
	currentSeasonWins := 0
	currentSeasonLosses := 0
	bowlWins := 0
	bowlLosses := 0

	for _, game := range historicGames {
		if !game.GameComplete || (game.GameComplete && game.SeasonID == timestamp.CollegeSeasonID && game.WeekID == timestamp.CollegeWeekID) || game.IsSpringGame {
			continue
		}
		winningSeason := game.SeasonID + 2020
		winningSeasonStr := strconv.Itoa(winningSeason)
		isAway := strconv.Itoa(game.AwayTeamID) == TeamID

		if (isAway && game.AwayTeamWin) || (!isAway && game.HomeTeamWin) {
			overallWins++

			if game.SeasonID == timestamp.CollegeSeasonID {
				currentSeasonWins++
			}

			if game.IsBowlGame {
				bowlWins++
			}

			if game.IsConferenceChampionship {
				conferenceChampionships = append(conferenceChampionships, winningSeasonStr)
				divisionTitles = append(divisionTitles, winningSeasonStr)
			}

			if game.IsNationalChampionship {
				nationalChampionships = append(nationalChampionships, winningSeasonStr)
			}
		} else {
			overallLosses++

			if game.SeasonID == timestamp.CollegeSeasonID {
				currentSeasonLosses++
			}

			if game.IsBowlGame {
				bowlLosses++
			}

			if game.IsConferenceChampionship {
				divisionTitles = append(divisionTitles, season)
			}
		}
	}

	response := models.TeamRecordResponse{
		OverallWins:             overallWins,
		OverallLosses:           overallLosses,
		CurrentSeasonWins:       currentSeasonWins,
		CurrentSeasonLosses:     currentSeasonLosses,
		BowlWins:                bowlWins,
		BowlLosses:              bowlLosses,
		ConferenceChampionships: conferenceChampionships,
		DivisionTitles:          divisionTitles,
		NationalChampionships:   nationalChampionships,
	}

	return response
}

func GetHistoricalNFLRecordsByTeamID(TeamID string) models.TeamRecordResponse {
	tsChn := make(chan structs.Timestamp)

	go func() {
		ts := GetTimestamp()
		tsChn <- ts
	}()

	timestamp := <-tsChn
	close(tsChn)

	historicGames := GetNFLGamesByTeamId(TeamID)
	var conferenceChampionships []string
	var divisionTitles []string
	var nationalChampionships []string
	overallWins := 0
	overallLosses := 0
	currentSeasonWins := 0
	currentSeasonLosses := 0
	bowlWins := 0
	bowlLosses := 0

	for _, game := range historicGames {
		if !game.GameComplete || (game.GameComplete && game.SeasonID == timestamp.CollegeSeasonID && game.WeekID == timestamp.CollegeWeekID) || game.IsPreseasonGame {
			continue
		}
		gameSeason := game.SeasonID + 2020
		isAway := strconv.Itoa(game.AwayTeamID) == TeamID

		if (isAway && game.AwayTeamWin) || (!isAway && game.HomeTeamWin) {
			overallWins++

			if game.SeasonID == timestamp.CollegeSeasonID {
				currentSeasonWins++
			}

			if game.IsPlayoffGame {
				bowlWins++
			}

			if game.IsConferenceChampionship {
				conferenceChampionships = append(conferenceChampionships, strconv.Itoa(gameSeason))
				divisionTitles = append(divisionTitles, strconv.Itoa(gameSeason))
			}

			if game.IsSuperBowl {
				nationalChampionships = append(nationalChampionships, strconv.Itoa(gameSeason))
			}
		} else {
			overallLosses++

			if game.SeasonID == timestamp.CollegeSeasonID {
				currentSeasonLosses++
			}

			if game.IsPlayoffGame {
				bowlLosses++
			}
		}
	}

	response := models.TeamRecordResponse{
		OverallWins:             overallWins,
		OverallLosses:           overallLosses,
		CurrentSeasonWins:       currentSeasonWins,
		CurrentSeasonLosses:     currentSeasonLosses,
		BowlWins:                bowlWins,
		BowlLosses:              bowlLosses,
		ConferenceChampionships: conferenceChampionships,
		DivisionTitles:          divisionTitles,
		NationalChampionships:   nationalChampionships,
	}

	return response
}

// GetStandingsByConferenceIDAndSeasonID
func GetCFBStandingsByTeamIDAndSeasonID(TeamID string, seasonID string) structs.CollegeStandings {
	var standings structs.CollegeStandings
	db := dbprovider.GetInstance().GetDB()
	err := db.Where("team_id = ? AND season_id = ?", TeamID, seasonID).
		Find(&standings).Error
	if err != nil {
		log.Fatal(err)
	}
	return standings
}

func GetAllCollegeStandingsBySeasonID(seasonID string) []structs.CollegeStandings {
	return repository.FindAllCollegeStandingsRecords(repository.StandingsQuery{
		SeasonID: seasonID,
	})
}

func GetAllNFLStandingsBySeasonID(seasonID string) []structs.NFLStandings {
	return repository.FindAllNFLStandingsRecords(repository.StandingsQuery{
		SeasonID: seasonID,
	})
}

func GetCollegeStandingsRecordByTeamID(id string, seasonID string) structs.CollegeStandings {
	return repository.FindAllCollegeStandingsRecords(repository.StandingsQuery{
		TeamID:   id,
		SeasonID: seasonID,
	})[0]
}

func ResetCollegeStandingsRanks() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.CollegeSeasonID))
	db.Model(&structs.CollegeStandings{}).Where("season_id = ?", seasonID).Updates(structs.CollegeStandings{Rank: 0})
}

func ResetCollegeStandings() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.CollegeSeasonID))
	collegeStandings := repository.FindAllCollegeStandingsRecords(repository.StandingsQuery{SeasonID: seasonID})
	for _, standings := range collegeStandings {
		standings.ResetCFBStandings()

		teamID := strconv.Itoa(standings.TeamID)
		games := GetCollegeGamesByTeamIdAndSeasonId(teamID, seasonID, false)
		for _, game := range games {
			if !game.GameComplete {
				continue
			}
			standings.UpdateCollegeStandings(game)
		}
		repository.SaveCFBStandingsRecord(standings, db)
	}
}

func GetCollegeStandingsMap(seasonID string) map[uint]structs.CollegeStandings {
	standingsMap := make(map[uint]structs.CollegeStandings)

	standings := repository.FindAllCollegeStandingsRecords(repository.StandingsQuery{SeasonID: seasonID})
	for _, stat := range standings {
		standingsMap[uint(stat.TeamID)] = stat
	}

	return standingsMap
}

func GetStandingsHistoryByTeamID(id string) []structs.CollegeStandings {
	return repository.FindAllCollegeStandingsRecords(repository.StandingsQuery{
		TeamID: id,
	})
}

func GenerateNewSeasonStandings() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	teams := GetAllCollegeTeams()
	collegeStandings := []structs.CollegeStandings{}
	nflStandings := []structs.NFLStandings{}

	nflTeams := GetAllNFLTeams()

	for _, t := range teams {
		if !t.IsActive {
			continue
		}
		leagueID := 1
		league := "FBS"
		if !t.IsFBS {
			leagueID = 2
			league = "FCS"
		}

		standings := structs.CollegeStandings{
			TeamID:           int(t.ID),
			TeamName:         t.TeamName,
			SeasonID:         ts.CollegeSeasonID,
			Season:           ts.Season,
			ConferenceID:     t.ConferenceID,
			ConferenceName:   t.Conference,
			PostSeasonStatus: "None",
			IsFBS:            t.IsFBS,
			DivisionID:       t.DivisionID,
			LeagueID:         uint(leagueID),
			LeagueName:       league,
			BaseStandings: structs.BaseStandings{
				Coach:    t.Coach,
				TeamAbbr: t.TeamAbbr,
			},
		}

		collegeStandings = append(collegeStandings, standings)
	}
	repository.CreateCFBStandingsBatch(db, collegeStandings, 100)

	for _, t := range nflTeams {

		standings := structs.NFLStandings{
			TeamID:           t.ID,
			TeamName:         t.TeamName,
			SeasonID:         uint(ts.CollegeSeasonID),
			Season:           uint(ts.Season),
			ConferenceID:     t.ConferenceID,
			ConferenceName:   t.Conference,
			PostSeasonStatus: "None",
			DivisionID:       t.DivisionID,
			BaseStandings: structs.BaseStandings{
				Coach:    t.Coach,
				TeamAbbr: t.TeamAbbr,
			},
		}

		nflStandings = append(nflStandings, standings)
	}
	repository.CreateNFLStandingsBatch(db, nflStandings, 20)
}

func BaseGetCollegeRankings() {
	ts := GetTimestamp()
	seasonID := ts.CollegeSeasonID

	for round := 1; round <= seasonID; round++ {
		seasonIDStr := strconv.Itoa(int(round))
		GenerateCollegeRankings(seasonIDStr)
	}
}

func GenerateCollegeRankings(seasonID string) {
	db := dbprovider.GetInstance().GetDB()
	collegeTeams := GetAllCollegeTeams()
	collegeTeamMap := MakeCollegeTeamMap(collegeTeams)
	collegeStandings := repository.FindAllCollegeStandingsRecords(repository.StandingsQuery{SeasonID: seasonID})
	collegeGames := repository.FindCollegeGamesRecords(seasonID, false)

	preSeasonCFBSeasonStats := GetAllCollegeTeamSeasonStatsBySeason(seasonID, "1")
	regularSeasonCFBSeasonStats := GetAllCollegeTeamSeasonStatsBySeason(seasonID, "2")
	postSeasonCFBSeasonStats := GetAllCollegeTeamSeasonStatsBySeason(seasonID, "3")

	_ = preSeasonCFBSeasonStats
	_ = regularSeasonCFBSeasonStats
	_ = postSeasonCFBSeasonStats

	// Build lookup structures
	standingsMap := MakeCollegeStandingsMapByTeamID(collegeStandings)
	gamesMap := MakeHistoricGamesMapByTeamID(collegeGames)

	// calcWinPct computes win percentage on the fly, since TotalWinPercentage
	// may not be populated in the DB yet.
	calcWinPct := func(wins, losses int) float64 {
		if wins+losses == 0 {
			return 0
		}
		return float64(wins) / float64(wins+losses)
	}

	// -------------------------------------------------------------------------
	// Pass 1: adjusted SOS, preliminary SOR, RPI, ConferenceStrengthAdj.
	//
	// Adjusted SOS: each opponent's win% is computed excluding the game they
	// played against the team being evaluated. This prevents a 12-0 team and
	// a 0-12 team from getting different SOS values just because of their own
	// result — if both faced the same 12 opponents, they should see the same SOS.
	// -------------------------------------------------------------------------
	for idx, standings := range collegeStandings {
		teamID := uint(standings.TeamID)
		teamGames := gamesMap[teamID]

		var adjOppWinPcts []float64
		var oppOppWinPcts []float64
		var adjConfOppWinPcts []float64

		sorNumerator := 0.0
		sorDenominator := 0.0

		for _, game := range teamGames {
			if !game.GameComplete {
				continue
			}

			isAway := game.AwayTeamID == int(teamID)
			var oppID uint
			if isAway {
				oppID = uint(game.HomeTeamID)
			} else {
				oppID = uint(game.AwayTeamID)
			}

			oppStandings, ok := standingsMap[oppID]
			if !ok {
				continue
			}

			isWinner := (!isAway && game.HomeTeamWin) || (isAway && game.AwayTeamWin)

			// Adjusted opponent win%: remove the head-to-head game from O's record.
			// If O beat T (!isWinner), subtract one win; always subtract one game.
			adjOppWins := oppStandings.TotalWins
			if !isWinner {
				adjOppWins--
			}
			adjOppGames := oppStandings.TotalWins + oppStandings.TotalLosses - 1
			adjOppWinPct := 0.0
			if adjOppGames > 0 && adjOppWins >= 0 {
				adjOppWinPct = float64(adjOppWins) / float64(adjOppGames)
			}

			adjOppWinPcts = append(adjOppWinPcts, adjOppWinPct)

			if game.IsConference {
				adjConfOppWinPcts = append(adjConfOppWinPcts, adjOppWinPct)
			}

			// Opponents' opponents win% (RPI third component) — unadjusted
			oppGames := gamesMap[oppID]
			sumOppOpp := 0.0
			countOppOpp := 0
			for _, og := range oppGames {
				if !og.GameComplete {
					continue
				}
				isOppAway := og.AwayTeamID == int(oppID)
				var oppOppID uint
				if isOppAway {
					oppOppID = uint(og.HomeTeamID)
				} else {
					oppOppID = uint(og.AwayTeamID)
				}
				if ooStandings, ok2 := standingsMap[oppOppID]; ok2 {
					sumOppOpp += calcWinPct(ooStandings.TotalWins, ooStandings.TotalLosses)
					countOppOpp++
				}
			}
			if countOppOpp > 0 {
				oppOppWinPcts = append(oppOppWinPcts, sumOppOpp/float64(countOppOpp))
			}

			// Preliminary SOR using adjusted opponent win% (refined in Pass 2)
			sorDenominator++
			if isWinner {
				sorNumerator += 0.5 + 0.5*adjOppWinPct
			}
		}

		// Adjusted SOS
		sos := float32(0)
		if len(adjOppWinPcts) > 0 {
			sum := 0.0
			for _, v := range adjOppWinPcts {
				sum += v
			}
			sos = float32(sum / float64(len(adjOppWinPcts)))
		}

		// Preliminary SOR
		sor := float32(0)
		if sorDenominator > 0 {
			sor = float32(sorNumerator / sorDenominator)
		}

		// RPI: 25% team win% + 50% adjusted avg opp win% (SOS) + 25% avg opp-opp win%
		avgOppOppWinPct := float32(0)
		if len(oppOppWinPcts) > 0 {
			sum := 0.0
			for _, v := range oppOppWinPcts {
				sum += v
			}
			avgOppOppWinPct = float32(sum / float64(len(oppOppWinPcts)))
		}
		teamWinPct := float32(calcWinPct(standings.TotalWins, standings.TotalLosses))
		rpi := 0.25*teamWinPct + 0.50*sos + 0.25*avgOppOppWinPct

		// ConferenceStrengthAdj — fall back to SOS for independent teams (conf IDs 13 and 22)
		// which never have IsConference = true games.
		confStrengthAdj := sos
		if len(adjConfOppWinPcts) > 0 {
			sum := 0.0
			for _, v := range adjConfOppWinPcts {
				sum += v
			}
			confStrengthAdj = float32(sum / float64(len(adjConfOppWinPcts)))
		}

		collegeStandings[idx].SOS = sos
		collegeStandings[idx].SOR = sor
		collegeStandings[idx].RPI = rpi
		collegeStandings[idx].ConferenceStrengthAdj = confStrengthAdj
	}

	// -------------------------------------------------------------------------
	// Interim sort: rank FBS teams by the four quality metrics + win% so that
	// Tier1/Tier2/BadLoss thresholds are fully self-referential and never depend
	// on the poll Rank column.
	// -------------------------------------------------------------------------
	type interimEntry struct {
		teamID uint
		score  float64
		isFBS  bool
	}

	interimEntries := make([]interimEntry, len(collegeStandings))
	for i, s := range collegeStandings {
		team, ok := collegeTeamMap[uint(s.TeamID)]
		isFBS := ok && team.IsFBS
		winPct := calcWinPct(s.TotalWins, s.TotalLosses)
		score := winPct*0.35 +
			float64(s.RPI)*0.30 +
			float64(s.SOR)*0.20 +
			float64(s.SOS)*0.10 +
			float64(s.ConferenceStrengthAdj)*0.05
		interimEntries[i] = interimEntry{teamID: uint(s.TeamID), score: score, isFBS: isFBS}
	}

	sort.SliceStable(interimEntries, func(i, j int) bool {
		if interimEntries[i].isFBS != interimEntries[j].isFBS {
			return interimEntries[i].isFBS
		}
		return interimEntries[i].score > interimEntries[j].score
	})

	// Map teamID → interim rank (FBS only, 1-indexed)
	interimRankMap := make(map[uint]uint16)
	interimRank := uint16(1)
	totalFBSTeams := 0
	for _, e := range interimEntries {
		if !e.isFBS {
			break
		}
		interimRankMap[e.teamID] = interimRank
		interimRank++
		totalFBSTeams++
	}

	// -------------------------------------------------------------------------
	// Pass 2: quality-adjusted SOR, Tier1Wins, Tier2Wins, scaled bad loss penalty.
	//
	// Quality-adjusted SOR: opponent quality is derived from interim rank so that
	// beating a highly-ranked FBS team is worth more than beating an FCS team
	// with an identical record.
	//   FBS opponent at rank R: quality = 1.0 - (R-1)/totalFBSTeams
	//   FCS opponent:           quality = oppWinPct * 0.5 (half-credit)
	//
	// Scaled bad loss penalty (replaces binary count in composite score):
	//   FBS loses to FCS:                   penalty += 0.001 * totalFBSTeams
	//   spread = oppRank - teamRank > 25:   penalty += 0.001 * spread
	// BadLosses (uint16) still stores the count for display purposes.
	// -------------------------------------------------------------------------
	const badLossSpreadThreshold = 25

	// Per-team scaled penalty used in Pass 3 (indexed by collegeStandings position)
	badLossPenalties := make([]float64, len(collegeStandings))

	for idx, standings := range collegeStandings {
		teamID := uint(standings.TeamID)
		teamGames := gamesMap[teamID]
		teamInterimRank, teamIsFBS := interimRankMap[teamID]

		tier1Wins := 0
		tier2Wins := 0
		badLossCount := 0

		adjSorNumerator := 0.0
		adjSorDenominator := 0.0

		for _, game := range teamGames {
			if !game.GameComplete {
				continue
			}

			isAway := game.AwayTeamID == int(teamID)
			var oppID uint
			if isAway {
				oppID = uint(game.HomeTeamID)
			} else {
				oppID = uint(game.AwayTeamID)
			}

			isWinner := (!isAway && game.HomeTeamWin) || (isAway && game.AwayTeamWin)
			oppInterimRank, oppIsFBS := interimRankMap[oppID]
			oppStandings, hasOppStandings := standingsMap[oppID]

			// Quality-adjusted SOR
			adjSorDenominator++
			if isWinner {
				var oppQuality float64
				if oppIsFBS && totalFBSTeams > 0 {
					// FBS: linear scale from 1.0 (rank 1) down to ~0 (last rank)
					oppQuality = 1.0 - float64(oppInterimRank-1)/float64(totalFBSTeams)
				} else if hasOppStandings {
					// FCS: half-credit regardless of record
					oppQuality = calcWinPct(oppStandings.TotalWins, oppStandings.TotalLosses) * 0.5
				}
				adjSorNumerator += oppQuality
			}

			// Tier wins
			if isWinner {
				if oppIsFBS && oppInterimRank <= 25 {
					tier1Wins++
				} else if oppIsFBS && oppInterimRank <= 50 {
					tier2Wins++
				}
			} else {
				// Scaled bad loss penalty
				if teamIsFBS && !oppIsFBS {
					// Loss to FCS: harshest penalty
					badLossPenalties[idx] += 0.001 * float64(totalFBSTeams)
					badLossCount++
				} else if teamIsFBS && oppIsFBS {
					spread := int(oppInterimRank) - int(teamInterimRank)
					if spread > badLossSpreadThreshold {
						badLossPenalties[idx] += 0.001 * float64(spread)
						badLossCount++
					}
				}
			}
		}

		adjSor := float32(0)
		if adjSorDenominator > 0 {
			adjSor = float32(adjSorNumerator / adjSorDenominator)
		}

		collegeStandings[idx].SOR = adjSor
		collegeStandings[idx].Tier1Wins = uint16(tier1Wins)
		collegeStandings[idx].Tier2Wins = uint16(tier2Wins)
		collegeStandings[idx].BadLosses = uint16(badLossCount)
	}

	// -------------------------------------------------------------------------
	// Pass 3: compute full composite ToucanScore and assign final ToucanRank.
	// -------------------------------------------------------------------------
	type teamScore struct {
		idx   int
		score float64
		isFBS bool
	}

	scores := make([]teamScore, len(collegeStandings))
	for idx, s := range collegeStandings {
		team, ok := collegeTeamMap[uint(s.TeamID)]
		isFBS := ok && team.IsFBS

		winPct := calcWinPct(s.TotalWins, s.TotalLosses)

		preseasonBonus := 0.0
		if s.PreseasonRank > 0 && s.PreseasonRank <= 25 {
			preseasonBonus = float64(26-int(s.PreseasonRank)) / float64(25) * 0.02
		}

		regularSeasonBonus := 0.0
		if s.RegularSeasonRank > 0 && s.RegularSeasonRank <= 25 {
			regularSeasonBonus = float64(26-int(s.RegularSeasonRank)) / float64(25) * 0.03
		}

		// Win%: 30%, RPI: 25%, SOR: 20%, SOS: 10%, Conference Strength: 5%,
		// Tier1 wins: +4% each, Tier2 wins: +2% each, preseason: up to 2%,
		// regular season rank: up to 3%, bad losses: scaled penalty from Pass 2
		score := winPct*0.30 +
			float64(s.RPI)*0.25 +
			float64(s.SOR)*0.20 +
			float64(s.SOS)*0.10 +
			float64(s.ConferenceStrengthAdj)*0.05 +
			float64(s.Tier1Wins)*0.04 +
			float64(s.Tier2Wins)*0.02 +
			preseasonBonus +
			regularSeasonBonus -
			badLossPenalties[idx]

		scores[idx] = teamScore{idx: idx, score: score, isFBS: isFBS}
	}

	sort.SliceStable(scores, func(i, j int) bool {
		if scores[i].isFBS != scores[j].isFBS {
			return scores[i].isFBS
		}
		return scores[i].score > scores[j].score
	})

	rank := uint16(1)
	for _, sc := range scores {
		collegeStandings[sc.idx].ToucanRank = rank
		rank++
	}

	// Persist all updated standings
	for idx := range collegeStandings {
		collegeStandings[idx].CalculatePercentages()
		repository.SaveCFBStandingsRecord(collegeStandings[idx], db)
	}
}
