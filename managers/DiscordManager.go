package managers

import (
	"strconv"

	"github.com/CalebRose/SimFBA/structs"
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
