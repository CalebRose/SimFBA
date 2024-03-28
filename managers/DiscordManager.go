package managers

import (
	"strconv"
	"sync"

	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
)

func CompareTwoTeams(t1ID, t2ID string) structs.CFBComparisonModel {

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
		if !game.GameComplete {
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
	if t1CurrentStreak > t2CurrentStreak {
		currentStreak = t1CurrentStreak
	} else {
		currentStreak = t2CurrentStreak
	}

	return structs.CFBComparisonModel{
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
	matches := GetCollegeGamesByTeamIdAndSeasonId(id, seasonId)
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

func GetCFBPlayByPlayStreamData(timeslot, week string, isFBS bool) []structs.StreamResponse {
	ts := GetTimestamp()
	weekNum := util.ConvertStringToInt(week)
	collegeWeek := ts.CollegeWeek
	collegeWeekID := ts.CollegeWeekID
	if collegeWeek == weekNum {
		// Continue
	} else {
		diff := collegeWeek - weekNum
		collegeWeekID = ts.CollegeWeekID - diff
	}
	teamMap := GetTeamProfileMap()
	games := GetCollegeGamesByTimeslotAndWeekId(timeslot, strconv.Itoa(collegeWeekID))

	streams := []structs.StreamResponse{}

	for _, game := range games {
		homeTeam := teamMap[game.HomeTeam]
		awayTeam := teamMap[game.AwayTeam]
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
		wg.Add(5)
		var (
			homeGameplan structs.CollegeGameplan
			awayGameplan structs.CollegeGameplan
			playByPlays  []structs.CollegePlayByPlay
			homePlayers  []structs.GameResultsPlayer
			awayPlayers  []structs.GameResultsPlayer
		)
		homeTeamID := strconv.Itoa(game.HomeTeamID)
		awayTeamID := strconv.Itoa(game.HomeTeamID)

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
			homePlayers = GetAllCollegePlayersWithGameStatsByTeamID(homeTeamID, gameID)
		}()

		go func() {
			defer wg.Done()
			awayPlayers = GetAllCollegePlayersWithGameStatsByTeamID(awayTeamID, gameID)
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
			AwayTeamID:          uint(game.AwayTeamID),
			AwayTeam:            game.AwayTeam,
			AwayTeamCoach:       game.AwayTeam,
			AwayTeamRank:        game.AwayTeamRank,
			HomeOffensiveScheme: homeGameplan.OffensiveScheme,
			HomeDefensiveScheme: homeGameplan.DefensiveScheme,
			AwayOffensiveScheme: awayGameplan.OffensiveScheme,
			AwayDefensiveScheme: awayGameplan.DefensiveScheme,
			Streams:             playbyPlayResponse,
		}

		streams = append(streams, stream)
	}

	return streams
}
