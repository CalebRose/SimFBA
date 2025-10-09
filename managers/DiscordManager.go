package managers

import (
	"strconv"
	"sync"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
)

func CompareTwoCFBTeams(t1ID, t2ID string) structs.FlexComparisonModel {
	ts := GetTimestamp()
	teamOneChan := make(chan structs.CollegeTeam)
	teamTwoChan := make(chan structs.CollegeTeam)

	go func() {
		t1 := GetTeamByTeamID(t1ID)
		teamOneChan <- t1
	}()

	teamOne := <-teamOneChan
	close(teamOneChan)

	go func() {
		t2 := GetTeamByTeamID(t2ID)
		teamTwoChan <- t2
	}()

	teamTwo := <-teamTwoChan
	close(teamTwoChan)

	allTeamOneGames := GetCollegeGamesByTeamId(t1ID)

	t1Wins := 0
	t1Losses := 0
	t1Streak := 0
	t1CurrentStreak := 0
	t1LargestMarginSeason := 0
	t1LargestMarginDiff := 0
	t1LargestMarginScore := ""
	t2Wins := 0
	t2Losses := 0
	t2Streak := 0
	t2CurrentStreak := 0
	latestWin := ""
	t2LargestMarginSeason := 0
	t2LargestMarginDiff := 0
	t2LargestMarginScore := ""

	for _, game := range allTeamOneGames {
		if !game.GameComplete ||
			(game.Week == ts.CollegeWeek &&
				((game.TimeSlot == "Thursday Night" && !ts.ThursdayGames) ||
					(game.TimeSlot == "Friday Night" && !ts.FridayGames) ||
					(game.TimeSlot == "Saturday Morning" && !ts.SaturdayMorning) ||
					(game.TimeSlot == "Saturday Afternoon" && !ts.SaturdayNoon) ||
					(game.TimeSlot == "Saturday Evening" && !ts.SaturdayEvening) ||
					(game.TimeSlot == "Saturday Night" && !ts.SaturdayNight))) {
			continue
		}
		doComparison := (game.HomeTeamID == int(teamOne.ID) && game.AwayTeamID == int(teamTwo.ID)) ||
			(game.HomeTeamID == int(teamTwo.ID) && game.AwayTeamID == int(teamOne.ID))

		if !doComparison {
			continue
		}
		homeTeamTeamOne := game.HomeTeamID == int(teamOne.ID)
		if homeTeamTeamOne {
			if game.HomeTeamWin {
				t1Wins += 1
				t1CurrentStreak += 1
				latestWin = game.HomeTeam
				diff := game.HomeTeamScore - game.AwayTeamScore
				if diff > t1LargestMarginDiff {
					t1LargestMarginDiff = diff
					t1LargestMarginSeason = game.SeasonID + 2020
					t1LargestMarginScore = "" + strconv.Itoa(game.HomeTeamScore) + "-" + strconv.Itoa(game.AwayTeamScore)
				}
			} else {
				t1Streak = t1CurrentStreak
				t1CurrentStreak = 0
				t1Losses += 1
			}
		} else {
			if game.HomeTeamWin {
				t2Wins += 1
				t2CurrentStreak += 1
				latestWin = game.HomeTeam
				diff := game.HomeTeamScore - game.AwayTeamScore
				if diff > t2LargestMarginDiff {
					t2LargestMarginDiff = diff
					t2LargestMarginSeason = game.SeasonID + 2020
					t2LargestMarginScore = "" + strconv.Itoa(game.HomeTeamScore) + "-" + strconv.Itoa(game.AwayTeamScore)
				}
			} else {
				t2Streak = t2CurrentStreak
				t2CurrentStreak = 0
				t2Losses += 1
			}
		}

		awayTeamTeamOne := game.AwayTeamID == int(teamOne.ID)
		if awayTeamTeamOne {
			if game.AwayTeamWin {
				t1Wins += 1
				t1CurrentStreak += 1
				latestWin = game.AwayTeam
				diff := game.AwayTeamScore - game.HomeTeamScore
				if diff > t1LargestMarginDiff {
					t1LargestMarginDiff = diff
					t1LargestMarginSeason = game.SeasonID + 2020
					t1LargestMarginScore = "" + strconv.Itoa(game.AwayTeamScore) + "-" + strconv.Itoa(game.HomeTeamScore)
				}
			} else {
				t1Streak = t1CurrentStreak
				t1CurrentStreak = 0
				t1Losses += 1
			}
		} else {
			if game.AwayTeamWin {
				t2Wins += 1
				t2CurrentStreak += 1
				latestWin = game.AwayTeam
				diff := game.AwayTeamScore - game.HomeTeamScore
				if diff > t2LargestMarginDiff {
					t2LargestMarginDiff = diff
					t2LargestMarginSeason = game.SeasonID + 2020
					t2LargestMarginScore = "" + strconv.Itoa(game.AwayTeamScore) + "-" + strconv.Itoa(game.HomeTeamScore)
				}
			} else {
				t2Streak = t2CurrentStreak
				t2CurrentStreak = 0
				t2Losses += 1
			}
		}
	}

	if t1CurrentStreak > 0 && t1CurrentStreak > t1Streak {
		t1Streak = t1CurrentStreak
	}
	if t2CurrentStreak > 0 && t2CurrentStreak > t2Streak {
		t2Streak = t2CurrentStreak
	}

	currentStreak := 0
	currentStreak = max(t1CurrentStreak, t2CurrentStreak)

	return structs.FlexComparisonModel{
		TeamOneID:      teamOne.ID,
		TeamOne:        teamOne.TeamAbbr,
		TeamOneWins:    uint(t1Wins),
		TeamOneLosses:  uint(t1Losses),
		TeamOneStreak:  uint(t1Streak),
		TeamOneMSeason: t1LargestMarginSeason,
		TeamOneMScore:  t1LargestMarginScore,
		TeamTwoID:      teamTwo.ID,
		TeamTwo:        teamTwo.TeamAbbr,
		TeamTwoWins:    uint(t2Wins),
		TeamTwoLosses:  uint(t2Losses),
		TeamTwoStreak:  uint(t2Streak),
		TeamTwoMSeason: t2LargestMarginSeason,
		TeamTwoMScore:  t2LargestMarginScore,
		CurrentStreak:  uint(currentStreak),
		LatestWin:      latestWin,
	}
}

func CompareTwoNFLTeams(t1ID, t2ID string) structs.FlexComparisonModel {
	ts := GetTimestamp()
	teamOneChan := make(chan structs.NFLTeam)
	teamTwoChan := make(chan structs.NFLTeam)

	go func() {
		t1 := GetNFLTeamByTeamID(t1ID)
		teamOneChan <- t1
	}()

	teamOne := <-teamOneChan
	close(teamOneChan)

	go func() {
		t2 := GetNFLTeamByTeamID(t2ID)
		teamTwoChan <- t2
	}()

	teamTwo := <-teamTwoChan
	close(teamTwoChan)

	allTeamOneGames := GetNFLGamesByTeamId(t1ID)

	t1Wins := 0
	t1Losses := 0
	t1Streak := 0
	t1CurrentStreak := 0
	t1LargestMarginSeason := 0
	t1LargestMarginDiff := 0
	t1LargestMarginScore := ""
	t2Wins := 0
	t2Losses := 0
	t2Streak := 0
	t2CurrentStreak := 0
	latestWin := ""
	t2LargestMarginSeason := 0
	t2LargestMarginDiff := 0
	t2LargestMarginScore := ""

	for _, game := range allTeamOneGames {
		if !game.GameComplete ||
			(game.Week == ts.NFLWeek &&
				((game.TimeSlot == "Sunday Noon" && !ts.NFLSundayNoon) ||
					(game.TimeSlot == "Sunday Afternoon" && !ts.NFLSundayAfternoon) ||
					(game.TimeSlot == "Sunday Night Football" && !ts.NFLSundayEvening) ||
					(game.TimeSlot == "Monday Night Football" && !ts.NFLMondayEvening) ||
					(game.TimeSlot == "Thursday Night Football" && !ts.NFLThursday))) {
			continue
		}
		doComparison := (game.HomeTeamID == int(teamOne.ID) && game.AwayTeamID == int(teamTwo.ID)) ||
			(game.HomeTeamID == int(teamTwo.ID) && game.AwayTeamID == int(teamOne.ID))

		if !doComparison {
			continue
		}
		homeTeamTeamOne := game.HomeTeamID == int(teamOne.ID)
		if homeTeamTeamOne {
			if game.HomeTeamWin {
				t1Wins += 1
				t1CurrentStreak += 1
				latestWin = game.HomeTeam
				diff := game.HomeTeamScore - game.AwayTeamScore
				if diff > t1LargestMarginDiff {
					t1LargestMarginDiff = diff
					t1LargestMarginSeason = game.SeasonID + 2020
					t1LargestMarginScore = "" + strconv.Itoa(game.HomeTeamScore) + "-" + strconv.Itoa(game.AwayTeamScore)
				}
			} else {
				t1Streak = t1CurrentStreak
				t1CurrentStreak = 0
				t1Losses += 1
			}
		} else {
			if game.HomeTeamWin {
				t2Wins += 1
				t2CurrentStreak += 1
				latestWin = game.HomeTeam
				diff := game.HomeTeamScore - game.AwayTeamScore
				if diff > t2LargestMarginDiff {
					t2LargestMarginDiff = diff
					t2LargestMarginSeason = game.SeasonID + 2020
					t2LargestMarginScore = "" + strconv.Itoa(game.HomeTeamScore) + "-" + strconv.Itoa(game.AwayTeamScore)
				}
			} else {
				t2Streak = t2CurrentStreak
				t2CurrentStreak = 0
				t2Losses += 1
			}
		}

		awayTeamTeamOne := game.AwayTeamID == int(teamOne.ID)
		if awayTeamTeamOne {
			if game.AwayTeamWin {
				t1Wins += 1
				t1CurrentStreak += 1
				latestWin = game.AwayTeam
				diff := game.AwayTeamScore - game.HomeTeamScore
				if diff > t1LargestMarginDiff {
					t1LargestMarginDiff = diff
					t1LargestMarginSeason = game.SeasonID + 2020
					t1LargestMarginScore = "" + strconv.Itoa(game.AwayTeamScore) + "-" + strconv.Itoa(game.HomeTeamScore)
				}
			} else {
				t1Streak = t1CurrentStreak
				t1CurrentStreak = 0
				t1Losses += 1
			}
		} else {
			if game.AwayTeamWin {
				t2Wins += 1
				t2CurrentStreak += 1
				latestWin = game.AwayTeam
				diff := game.AwayTeamScore - game.HomeTeamScore
				if diff > t2LargestMarginDiff {
					t2LargestMarginDiff = diff
					t2LargestMarginSeason = game.SeasonID + 2020
					t2LargestMarginScore = "" + strconv.Itoa(game.AwayTeamScore) + "-" + strconv.Itoa(game.HomeTeamScore)
				}
			} else {
				t2Streak = t2CurrentStreak
				t2CurrentStreak = 0
				t2Losses += 1
			}
		}
	}

	if t1CurrentStreak > 0 && t1CurrentStreak > t1Streak {
		t1Streak = t1CurrentStreak
	}
	if t2CurrentStreak > 0 && t2CurrentStreak > t2Streak {
		t2Streak = t2CurrentStreak
	}

	currentStreak := 0
	currentStreak = max(t1CurrentStreak, t2CurrentStreak)

	return structs.FlexComparisonModel{
		TeamOneID:      teamOne.ID,
		TeamOne:        teamOne.TeamAbbr,
		TeamOneWins:    uint(t1Wins),
		TeamOneLosses:  uint(t1Losses),
		TeamOneStreak:  uint(t1Streak),
		TeamOneMSeason: t1LargestMarginSeason,
		TeamOneMScore:  t1LargestMarginScore,
		TeamTwoID:      teamTwo.ID,
		TeamTwo:        teamTwo.TeamAbbr,
		TeamTwoWins:    uint(t2Wins),
		TeamTwoLosses:  uint(t2Losses),
		TeamTwoStreak:  uint(t2Streak),
		TeamTwoMSeason: t2LargestMarginSeason,
		TeamTwoMScore:  t2LargestMarginScore,
		CurrentStreak:  uint(currentStreak),
		LatestWin:      latestWin,
	}
}

func GetCFBTeamDataForDiscord(id string) structs.CollegeTeamResponseData {
	ts := GetTimestamp()
	seasonId := strconv.Itoa(ts.CollegeSeasonID)

	team := GetTeamByTeamID(id)
	standings := GetCollegeStandingsRecordByTeamID(id, seasonId)
	matches := GetCollegeGamesByTeamIdAndSeasonId(id, seasonId, ts.CFBSpringGames)
	wins := 0
	losses := 0
	confWins := 0
	confLosses := 0
	matchList := []structs.CollegeGame{}

	for _, m := range matches {
		if m.Week > ts.CollegeWeek {
			break
		}
		gameNotRan := (m.TimeSlot == "Thursday Night" && !ts.ThursdayGames) ||
			(m.TimeSlot == "Friday Night" && !ts.FridayGames) ||
			(m.TimeSlot == "Saturday Morning" && !ts.SaturdayMorning) ||
			(m.TimeSlot == "Saturday Afternoon" && !ts.SaturdayNoon) ||
			(m.TimeSlot == "Saturday Evening" && !ts.SaturdayEvening) ||
			(m.TimeSlot == "Saturday Night" && !ts.SaturdayNight)

		earlierWeek := m.Week < ts.CollegeWeek

		if ((strconv.Itoa(int(m.HomeTeamID)) == id && m.HomeTeamWin) ||
			(strconv.Itoa(int(m.AwayTeamID)) == id && m.AwayTeamWin)) && (earlierWeek || !gameNotRan) {
			wins += 1
			if m.IsConference {
				confWins += 1
			}
		} else if ((strconv.Itoa(int(m.HomeTeamID)) == id && m.AwayTeamWin) ||
			(strconv.Itoa(int(m.AwayTeamID)) == id && m.HomeTeamWin)) && (earlierWeek || !gameNotRan) {
			losses += 1
			if m.IsConference {
				confLosses += 1
			}
		}
		if gameNotRan {
			m.HideScore()
		}
		if m.Week == ts.CollegeWeek {
			matchList = append(matchList, m)
		}
	}

	standings.MaskGames(wins, losses, confWins, confLosses)

	return structs.CollegeTeamResponseData{
		TeamData:        team,
		TeamStandings:   standings,
		UpcomingMatches: matchList,
	}
}

func GetCFBPlayByPlayStreamData(timeslot string, isFBS bool) []structs.StreamResponse {
	ts := GetTimestamp()
	collegeWeekID := ts.CollegeWeekID
	teamMap := GetCollegeTeamMap()
	games := GetCollegeGamesByTimeslotAndWeekId(strconv.Itoa(collegeWeekID), timeslot, ts.CFBSpringGames)

	streams := []structs.StreamResponse{}

	for _, game := range games {
		if !game.GameComplete || game.ID == 4117 {
			continue
		}
		homeTeam := teamMap[uint(game.HomeTeamID)]
		awayTeam := teamMap[uint(game.AwayTeamID)]
		// If it's a full FCS match up and we're not streaming FCS, skip
		if !homeTeam.IsFBS && !awayTeam.IsFBS && isFBS {
			continue
		}

		// If it's a partial FCS and we're streaming FCS, skip. It was streamed in FBS
		if ((!homeTeam.IsFBS && awayTeam.IsFBS) || (homeTeam.IsFBS && !awayTeam.IsFBS)) && !isFBS {
			continue
		}

		// If it's an FBS match and we're streaming FCS, skip
		if homeTeam.IsFBS && awayTeam.IsFBS && !isFBS {
			continue
		}
		gameID := strconv.Itoa(int(game.ID))
		var wg sync.WaitGroup
		var (
			homeGameplan structs.CollegeGameplan
			awayGameplan structs.CollegeGameplan
			playByPlays  []structs.CollegePlayByPlay
			homePlayers  []structs.GameResultsPlayer
			awayPlayers  []structs.GameResultsPlayer
			homeStats    []structs.CollegePlayerStats
			awayStats    []structs.CollegePlayerStats
		)
		homeTeamID := strconv.Itoa(game.HomeTeamID)
		awayTeamID := strconv.Itoa(game.AwayTeamID)
		wg.Add(2)

		go func() {
			defer wg.Done()
			homeStats = GetAllCollegePlayerStatsByGame(gameID, homeTeamID)
		}()

		go func() {
			defer wg.Done()
			awayStats = GetAllCollegePlayerStatsByGame(gameID, awayTeamID)
		}()

		wg.Wait()

		wg.Add(5)

		go func() {
			defer wg.Done()
			homeGameplan = GetGameplanByTeamID(homeTeamID)
		}()

		go func() {
			defer wg.Done()
			awayGameplan = GetGameplanByTeamID(awayTeamID)
		}()

		go func() {
			defer wg.Done()
			playByPlays = GetCFBPlayByPlaysByGameID(gameID)
		}()

		go func() {
			defer wg.Done()
			homePlayers = GetAllCollegePlayersWithGameStatsByTeamID(gameID, homeStats)
		}()

		go func() {
			defer wg.Done()
			awayPlayers = GetAllCollegePlayersWithGameStatsByTeamID(gameID, awayStats)
		}()

		wg.Wait()

		participantMap := getGameParticipantMap(homePlayers, awayPlayers)
		playbyPlayResponse := GenerateCFBPlayByPlayResponse(playByPlays, participantMap, true, game.HomeTeam, game.AwayTeam)

		stream := structs.StreamResponse{
			GameID:              game.ID,
			HomeTeamID:          uint(game.HomeTeamID),
			HomeTeam:            game.HomeTeam,
			HomeTeamCoach:       game.HomeTeamCoach,
			HomeTeamRank:        game.HomeTeamRank,
			HomeLabel:           homeTeam.TeamName + " " + homeTeam.Mascot,
			HomeTeamDiscordID:   homeTeam.DiscordID,
			AwayTeamID:          uint(game.AwayTeamID),
			AwayTeam:            game.AwayTeam,
			AwayTeamCoach:       game.AwayTeamCoach,
			AwayTeamRank:        game.AwayTeamRank,
			AwayTeamDiscordID:   awayTeam.DiscordID,
			AwayLabel:           awayTeam.TeamName + " " + awayTeam.Mascot,
			HomeOffensiveScheme: homeGameplan.OffensiveScheme,
			HomeDefensiveScheme: homeGameplan.DefensiveScheme,
			AwayOffensiveScheme: awayGameplan.OffensiveScheme,
			AwayDefensiveScheme: awayGameplan.DefensiveScheme,
			Streams:             playbyPlayResponse,
			GameCloud:           game.Cloud,
			GameWind:            game.WindCategory,
			GamePrecip:          game.Precip,
			GameWindSpeed:       game.WindSpeed,
			GameTemp:            game.GameTemp,
			Stadium:             game.Stadium,
			City:                game.City,
			State:               game.State,
		}

		streams = append(streams, stream)
	}

	return streams
}

func GetNFLPlayByPlayStreamData(timeslot string) []structs.StreamResponse {
	ts := GetTimestamp()
	nflWeekID := ts.NFLWeekID
	games := GetNFLGamesByTimeslotAndWeekId(strconv.Itoa(nflWeekID), timeslot, ts.NFLPreseason)

	streams := []structs.StreamResponse{}

	for _, game := range games {
		if !game.GameComplete {
			continue
		}

		gameID := strconv.Itoa(int(game.ID))
		var wg sync.WaitGroup
		var (
			homeGameplan structs.NFLGameplan
			awayGameplan structs.NFLGameplan
			playByPlays  []structs.NFLPlayByPlay
			homePlayers  []structs.GameResultsPlayer
			awayPlayers  []structs.GameResultsPlayer
			homeStats    []structs.NFLPlayerStats
			awayStats    []structs.NFLPlayerStats
		)
		homeTeamID := strconv.Itoa(game.HomeTeamID)
		awayTeamID := strconv.Itoa(game.AwayTeamID)

		wg.Add(2)

		go func() {
			defer wg.Done()
			homeStats = GetAllNFLPlayerStatsByGame(gameID, homeTeamID)
		}()

		go func() {
			defer wg.Done()
			awayStats = GetAllNFLPlayerStatsByGame(gameID, awayTeamID)
		}()

		wg.Wait()
		wg.Add(5)

		go func() {
			defer wg.Done()
			homeGameplan = GetNFLGameplanByTeamID(homeTeamID)
		}()

		go func() {
			defer wg.Done()
			awayGameplan = GetNFLGameplanByTeamID(awayTeamID)
		}()

		go func() {
			defer wg.Done()
			playByPlays = GetNFLPlayByPlaysByGameID(gameID)
		}()

		go func() {
			defer wg.Done()
			homePlayers = GetAllNFLPlayersWithGameStatsByTeamID(gameID, homeStats)
		}()

		go func() {
			defer wg.Done()
			awayPlayers = GetAllNFLPlayersWithGameStatsByTeamID(gameID, awayStats)
		}()

		wg.Wait()

		participantMap := getGameParticipantMap(homePlayers, awayPlayers)
		playbyPlayResponse := GenerateNFLPlayByPlayResponse(playByPlays, participantMap, true, game.HomeTeam, game.AwayTeam)

		stream := structs.StreamResponse{
			GameID:              game.ID,
			HomeTeamID:          uint(game.HomeTeamID),
			HomeTeam:            game.HomeTeam,
			HomeTeamCoach:       game.HomeTeamCoach,
			HomeLabel:           game.HomeTeam,
			AwayTeamID:          uint(game.AwayTeamID),
			AwayTeam:            game.AwayTeam,
			AwayTeamCoach:       game.AwayTeamCoach,
			AwayLabel:           game.AwayTeam,
			HomeOffensiveScheme: homeGameplan.OffensiveScheme,
			HomeDefensiveScheme: homeGameplan.DefensiveScheme,
			AwayOffensiveScheme: awayGameplan.OffensiveScheme,
			AwayDefensiveScheme: awayGameplan.DefensiveScheme,
			Streams:             playbyPlayResponse,
			GameCloud:           game.Cloud,
			GameWind:            game.WindCategory,
			GamePrecip:          game.Precip,
			GameWindSpeed:       game.WindSpeed,
			GameTemp:            game.GameTemp,
			Stadium:             game.Stadium,
			City:                game.City,
			State:               game.State,
		}

		streams = append(streams, stream)
	}

	return streams
}

func AssignDiscordIDToCollegeTeam(tID, dID string) {
	db := dbprovider.GetInstance().GetDB()

	team := GetTeamByTeamID(tID)

	team.AssignDiscordID(dID)

	repository.SaveCollegeTeamRecord(team, db)
}

func AssignDiscordIDToNFLTeam(tID, dID string) {
	db := dbprovider.GetInstance().GetDB()

	team := GetNFLTeamByTeamID(tID)

	team.AssignDiscordID(dID)

	repository.SaveNFLTeamRecord(team, db)
}
