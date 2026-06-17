package managers

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	time "time"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/firebase"
	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
)

// ─────────────────────────────────────────────
// CFB Game Request
// ─────────────────────────────────────────────

// CreateCFBGameRequest saves a new CFBGameRequest record to the database.
func CreateCFBGameRequest(request structs.CFBGameRequest) {
	db := dbprovider.GetInstance().GetDB()
	receivingTeamID := request.RequestingTeamID
	team := GetTeamByTeamID(strconv.Itoa(int(receivingTeamID)))
	if !isCFBUserTeam(team) {
		request.Accepted()
	}
	repository.CreateCFBGameRequest(request, db)
	if !isCFBUserTeam(team) {
		sendingTeam := GetTeamByTeamID(strconv.Itoa(int(request.SendingTeamID)))
		ctx := context.Background()
		uids := firebase.ResolveUIDsByUsernames(ctx, []string{sendingTeam.Coach})
		firebase.NotifyScheduleEvent(ctx, firebase.ScheduleEventNotificationInput{
			League:         "cfb",
			Domain:         firebase.DomainCFB,
			TeamID:         team.ID,
			RecipientUIDs:  uids,
			Message:        fmt.Sprintf("Out of Conference Game created against AI Team %s for Week %d.", team.TeamName, request.Week),
			SourceEventKey: firebase.BuildSourceEventKey("gamerequest", "cfb", "accept", strconv.Itoa(int(request.ID))),
		})
	} else {
		sendingTeam := GetTeamByTeamID(strconv.Itoa(int(request.SendingTeamID)))
		ctx := context.Background()
		uids := firebase.ResolveUIDsByUsernames(ctx, []string{team.Coach})
		firebase.NotifyScheduleEvent(ctx, firebase.ScheduleEventNotificationInput{
			League:         "cfb",
			Domain:         firebase.DomainCFB,
			TeamID:         team.ID,
			RecipientUIDs:  uids,
			Message:        fmt.Sprintf("Received OOC Game request for Week %d against %s.", request.Week, sendingTeam.TeamName),
			SourceEventKey: firebase.BuildSourceEventKey("gamerequest", "cfb", "accept", strconv.Itoa(int(request.ID))),
		})
	}
}

// AcceptCFBGameRequest marks the request as accepted and notifies the sending
// team's coach if they are a user-managed team.
func AcceptCFBGameRequest(requestID string) {
	db := dbprovider.GetInstance().GetDB()

	request := repository.FindCFBGameRequestRecord(repository.SchedulerQuery{ID: requestID})

	request.Accepted()
	repository.SaveCFBGameRequest(request, db)

	sendingTeam := GetTeamByTeamID(strconv.Itoa(int(request.SendingTeamID)))
	if isCFBUserTeam(sendingTeam) {
		receivingTeam := GetTeamByTeamID(strconv.Itoa(int(request.RequestingTeamID)))
		ctx := context.Background()
		uids := firebase.ResolveUIDsByUsernames(ctx, []string{sendingTeam.Coach})
		firebase.NotifyScheduleEvent(ctx, firebase.ScheduleEventNotificationInput{
			League:         "cfb",
			Domain:         firebase.DomainCFB,
			TeamID:         sendingTeam.ID,
			RecipientUIDs:  uids,
			Message:        fmt.Sprintf("%s has accepted your game request for Week %d.", receivingTeam.TeamName, request.Week),
			SourceEventKey: firebase.BuildSourceEventKey("gamerequest", "cfb", "accept", requestID),
		})
	}
}

// RejectCFBGameRequest deletes the request and notifies the sending team's coach
// if they are a user-managed team.
func RejectCFBGameRequest(requestID string) {
	db := dbprovider.GetInstance().GetDB()

	request := repository.FindCFBGameRequestRecord(repository.SchedulerQuery{ID: requestID})

	sendingTeam := GetTeamByTeamID(strconv.Itoa(int(request.SendingTeamID)))

	repository.DeleteCFBGameRequest(request, db)

	if isCFBUserTeam(sendingTeam) {
		receivingTeam := GetTeamByTeamID(strconv.Itoa(int(request.RequestingTeamID)))
		ctx := context.Background()
		uids := firebase.ResolveUIDsByUsernames(ctx, []string{sendingTeam.Coach})
		firebase.NotifyScheduleEvent(ctx, firebase.ScheduleEventNotificationInput{
			League:         "cfb",
			Domain:         firebase.DomainCFB,
			TeamID:         sendingTeam.ID,
			RecipientUIDs:  uids,
			Message:        fmt.Sprintf("%s has rejected your game request for Week %d.", receivingTeam.TeamName, request.Week),
			SourceEventKey: firebase.BuildSourceEventKey("gamerequest", "cfb", "reject", requestID),
		})
	}
}

// ProcessCFBGameRequest creates a CollegeGame record from the existing game
// request and marks the request as approved.
func ProcessCFBGameRequest(requestID string) {
	db := dbprovider.GetInstance().GetDB()

	request := repository.FindCFBGameRequestRecord(repository.SchedulerQuery{ID: requestID})

	homeTeam := GetTeamByTeamID(strconv.Itoa(int(request.HomeTeamID)))
	awayTeam := GetTeamByTeamID(strconv.Itoa(int(request.AwayTeamID)))

	stadiums := GetAllStadiums()
	stadiumByID := make(map[uint]structs.Stadium, len(stadiums))
	for _, s := range stadiums {
		stadiumByID[s.ID] = s
	}
	arenaID := request.ArenaID
	if arenaID == 0 {
		// If no arena specified, assign the home team's stadium
		arenaID = homeTeam.StadiumID
	}
	stadium := stadiumByID[arenaID]
	stadiumName := stadium.StadiumName
	city := stadium.City
	state := stadium.State
	region := stadium.Region
	if stadium.ID == 0 {
		stadiumName = homeTeam.Stadium
		city = homeTeam.City
		state = homeTeam.State
	}

	timeslot := request.Timeslot
	if timeslot == "" {
		timeslot = util.GetTimeslot(stadium.State, uint(homeTeam.ConferenceID))
	}

	game := structs.CollegeGame{
		HomeTeamID:    int(request.HomeTeamID),
		HomeTeam:      homeTeam.TeamAbbr,
		HomeTeamCoach: homeTeam.Coach,
		AwayTeamCoach: awayTeam.Coach,
		AwayTeamID:    int(request.AwayTeamID),
		AwayTeam:      awayTeam.TeamAbbr,
		Week:          int(request.Week),
		WeekID:        int(request.WeekID),
		SeasonID:      int(request.SeasonID),
		StadiumID:     arenaID,
		Stadium:       stadiumName,
		City:          city,
		State:         state,
		Region:        region,
		TimeSlot:      timeslot,
		IsNeutral:     request.IsNeutralSite,
		IsConference:  false,
		IsDivisional:  false,
		IsSpringGame:  request.IsSpringGame,
	}

	repository.CreateCFBGameRecordsBatch(db, []structs.CollegeGame{game}, 1)

	request.Approved()
	repository.SaveCFBGameRequest(request, db)
}

// VetoCFBGameRequest deletes the request and notifies both the sending and
// receiving teams' coaches if either is a user-managed team.
func VetoCFBGameRequest(requestID string) {
	db := dbprovider.GetInstance().GetDB()

	request := repository.FindCFBGameRequestRecord(repository.SchedulerQuery{ID: requestID})

	sendingTeam := GetTeamByTeamID(strconv.Itoa(int(request.SendingTeamID)))
	receivingTeam := GetTeamByTeamID(strconv.Itoa(int(request.RequestingTeamID)))

	repository.DeleteCFBGameRequest(request, db)

	ctx := context.Background()
	msg := fmt.Sprintf("The game request between %s and %s for Week %d has been vetoed.", sendingTeam.TeamName, receivingTeam.TeamName, request.Week)
	vetoKey := firebase.BuildSourceEventKey("gamerequest", "cfb", "veto", requestID)

	if isCFBUserTeam(sendingTeam) {
		uids := firebase.ResolveUIDsByUsernames(ctx, []string{sendingTeam.Coach})
		firebase.NotifyScheduleEvent(ctx, firebase.ScheduleEventNotificationInput{
			League:         "cfb",
			Domain:         firebase.DomainCFB,
			TeamID:         sendingTeam.ID,
			RecipientUIDs:  uids,
			Message:        msg,
			SourceEventKey: vetoKey + ":sending",
		})
	}
	if isCFBUserTeam(receivingTeam) {
		uids := firebase.ResolveUIDsByUsernames(ctx, []string{receivingTeam.Coach})
		firebase.NotifyScheduleEvent(ctx, firebase.ScheduleEventNotificationInput{
			League:         "cfb",
			Domain:         firebase.DomainCFB,
			TeamID:         receivingTeam.ID,
			RecipientUIDs:  uids,
			Message:        msg,
			SourceEventKey: vetoKey + ":receiving",
		})
	}
}

// ─────────────────────────────────────────────
// NFL Game Request
// ─────────────────────────────────────────────

// CreateNFLGameRequest saves a new NFLGameRequest record to the database.
func CreateNFLGameRequest(request structs.NFLGameRequest) {
	db := dbprovider.GetInstance().GetDB()
	repository.CreateNFLGameRequest(request, db)
}

// AcceptNFLGameRequest marks the request as accepted and notifies the sending
// team's owner if they are a user-managed team.
func AcceptNFLGameRequest(requestID string) {
	db := dbprovider.GetInstance().GetDB()

	request := repository.FindNFLGameRequestRecord(repository.SchedulerQuery{ID: requestID})

	request.Accepted()
	repository.SaveNFLGameRequest(request, db)

	sendingTeam := GetNFLTeamByTeamID(strconv.Itoa(int(request.SendingTeamID)))
	if isNFLUserTeam(sendingTeam) {
		receivingTeam := GetNFLTeamByTeamID(strconv.Itoa(int(request.RequestingTeamID)))
		ctx := context.Background()
		uids := firebase.ResolveUIDsByUsernames(ctx, []string{sendingTeam.NFLOwnerName})
		firebase.NotifyScheduleEvent(ctx, firebase.ScheduleEventNotificationInput{
			League:         "nfl",
			Domain:         firebase.DomainNFL,
			TeamID:         sendingTeam.ID,
			RecipientUIDs:  uids,
			Message:        fmt.Sprintf("%s has accepted your game request for Week %d.", receivingTeam.TeamName, request.Week),
			SourceEventKey: firebase.BuildSourceEventKey("gamerequest", "nfl", "accept", requestID),
		})
	}
}

// RejectNFLGameRequest deletes the request and notifies the sending team's owner
// if they are a user-managed team.
func RejectNFLGameRequest(requestID string) {
	db := dbprovider.GetInstance().GetDB()

	request := repository.FindNFLGameRequestRecord(repository.SchedulerQuery{ID: requestID})

	sendingTeam := GetNFLTeamByTeamID(strconv.Itoa(int(request.SendingTeamID)))

	repository.DeleteNFLGameRequest(request, db)

	if isNFLUserTeam(sendingTeam) {
		receivingTeam := GetNFLTeamByTeamID(strconv.Itoa(int(request.RequestingTeamID)))
		ctx := context.Background()
		uids := firebase.ResolveUIDsByUsernames(ctx, []string{sendingTeam.NFLOwnerName})
		firebase.NotifyScheduleEvent(ctx, firebase.ScheduleEventNotificationInput{
			League:         "nfl",
			Domain:         firebase.DomainNFL,
			TeamID:         sendingTeam.ID,
			RecipientUIDs:  uids,
			Message:        fmt.Sprintf("%s has rejected your game request for Week %d.", receivingTeam.TeamName, request.Week),
			SourceEventKey: firebase.BuildSourceEventKey("gamerequest", "nfl", "reject", requestID),
		})
	}
}

// ProcessNFLGameRequest creates an NFLGame record from the existing game
// request and marks the request as approved.
func ProcessNFLGameRequest(requestID string) {
	db := dbprovider.GetInstance().GetDB()

	request := repository.FindNFLGameRequestRecord(repository.SchedulerQuery{ID: requestID})

	homeTeam := GetNFLTeamByTeamID(strconv.Itoa(int(request.HomeTeamID)))
	awayTeam := GetNFLTeamByTeamID(strconv.Itoa(int(request.AwayTeamID)))

	stadiums := GetAllStadiums()
	stadiumByID := make(map[uint]structs.Stadium, len(stadiums))
	for _, s := range stadiums {
		stadiumByID[s.ID] = s
	}
	stadium := stadiumByID[request.ArenaID]

	isDivisional := homeTeam.DivisionID > 0 && homeTeam.DivisionID == awayTeam.DivisionID

	game := structs.NFLGame{
		HomeTeamID:      int(request.HomeTeamID),
		HomeTeam:        homeTeam.TeamName,
		AwayTeamID:      int(request.AwayTeamID),
		AwayTeam:        awayTeam.TeamName,
		Week:            int(request.Week),
		WeekID:          int(request.WeekID),
		SeasonID:        int(request.SeasonID),
		StadiumID:       request.ArenaID,
		Stadium:         stadium.StadiumName,
		City:            stadium.City,
		State:           stadium.State,
		Region:          stadium.Region,
		TimeSlot:        request.Timeslot,
		IsNeutral:       request.IsNeutralSite,
		IsConference:    homeTeam.ConferenceID == awayTeam.ConferenceID,
		IsDivisional:    isDivisional,
		IsPreseasonGame: request.IsPreseason,
	}

	repository.CreateNFLGameRecordsBatch(db, []structs.NFLGame{game}, 1)

	request.Approved()
	repository.SaveNFLGameRequest(request, db)
}

// VetoNFLGameRequest deletes the request and notifies both the sending and
// receiving teams' owners if either is a user-managed team.
func VetoNFLGameRequest(requestID string) {
	db := dbprovider.GetInstance().GetDB()

	request := repository.FindNFLGameRequestRecord(repository.SchedulerQuery{ID: requestID})

	sendingTeam := GetNFLTeamByTeamID(strconv.Itoa(int(request.SendingTeamID)))
	receivingTeam := GetNFLTeamByTeamID(strconv.Itoa(int(request.RequestingTeamID)))

	repository.DeleteNFLGameRequest(request, db)

	ctx := context.Background()
	msg := fmt.Sprintf("The game request between %s and %s for Week %d has been vetoed.", sendingTeam.TeamName, receivingTeam.TeamName, request.Week)
	vetoKey := firebase.BuildSourceEventKey("gamerequest", "nfl", "veto", requestID)

	if isNFLUserTeam(sendingTeam) {
		uids := firebase.ResolveUIDsByUsernames(ctx, []string{sendingTeam.NFLOwnerName})
		firebase.NotifyScheduleEvent(ctx, firebase.ScheduleEventNotificationInput{
			League:         "nfl",
			Domain:         firebase.DomainNFL,
			TeamID:         sendingTeam.ID,
			RecipientUIDs:  uids,
			Message:        msg,
			SourceEventKey: vetoKey + ":sending",
		})
	}
	if isNFLUserTeam(receivingTeam) {
		uids := firebase.ResolveUIDsByUsernames(ctx, []string{receivingTeam.NFLOwnerName})
		firebase.NotifyScheduleEvent(ctx, firebase.ScheduleEventNotificationInput{
			League:         "nfl",
			Domain:         firebase.DomainNFL,
			TeamID:         receivingTeam.ID,
			RecipientUIDs:  uids,
			Message:        msg,
			SourceEventKey: vetoKey + ":receiving",
		})
	}
}

// ─────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────

func isCFBUserTeam(team structs.CollegeTeam) bool {
	return team.Coach != "" && team.Coach != "AI"
}

func isNFLUserTeam(team structs.NFLTeam) bool {
	return team.NFLOwnerName != "" && team.NFLOwnerName != "AI"
}

func ShuffleTeams(teams []structs.CollegeTeam) []structs.CollegeTeam {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	shuffled := make([]structs.CollegeTeam, len(teams))
	copy(shuffled, teams)
	rng.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled
}

func GenerateOOCSchedule() {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.CollegeSeasonID))
	collegeTeams := GetAllCollegeTeams()
	stadium := GetAllStadiums()
	stadiumMap := MakeStadiumMapByTeamID(stadium, true)
	stadiumMapByID := MakeStadiumMapByID(stadium)

	rivalries := GetAllRivalries()
	rivalryMap := MakeHistoricRivalriesMapByTeamID(rivalries)
	maxAttempts := 5000
	var matchesToUpload []structs.CollegeGame
	var err error

	// Try multiple times to generate a complete schedule
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		matchesToUpload, err = attemptGenerateOOCSchedule(ts, seasonID, collegeTeams, stadiumMap, stadiumMapByID, rivalryMap)
		if err == nil {
			fmt.Printf("Successfully generated OOC schedule on attempt %d with %d matches.\n", attempt, len(matchesToUpload))
			break
		}
		if attempt < maxAttempts {
			fmt.Printf("Attempt %d failed: %v. Retrying...\n", attempt, err)
		} else {
			fmt.Printf("Failed to generate complete schedule after %d attempts. Last error: %v\n", maxAttempts, err)
			fmt.Printf("Uploading partial schedule with %d matches.\n", len(matchesToUpload))
		}
	}

	// Upload matches to database
	repository.CreateCFBGameRecordsBatch(db, matchesToUpload, 100)
}

func attemptGenerateOOCSchedule(ts structs.Timestamp, seasonID string, collegeTeams []structs.CollegeTeam, stadiumMap map[uint]structs.Stadium, stadiumMapByID map[uint]structs.Stadium, rivalryMap map[uint][]structs.CollegeRival) ([]structs.CollegeGame, error) {
	maxNumberOfGamesPerTeam := 12
	oocWeeks := []uint{14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0} // Schedule backwards: fill constrained weeks (7-15) first

	// Initialize tracking maps
	remainingGamesLeftByTeam := make(map[uint]int)
	homeGamesByTeam := make(map[uint]int)
	awayGamesByTeam := make(map[uint]int)
	opponentsFacedMap := make(map[uint]map[uint]bool) // teamID -> opponentTeamID -> bool
	// Count the number of games per week and time slot to ensure we are filling them correctly
	weekTimeSlotCounts := make(map[uint]int)
	numberOfGamesExpectedPerSlot := len(collegeTeams) / 2 // Each game has 2 teams
	gamesScheduledPerTeam := make(map[uint]int)

	teamMap := GetCollegeTeamMap()

	// Keep track of whether FBS teams scheduled an FCS opponent or if an FCS team has scheduled an FBS opponent.
	playedFCSOpponent := make(map[uint]bool)
	playedFBSOpponent := make(map[uint]bool)

	// Create a map to track user vs AI teams
	isUserTeam := make(map[uint]bool)
	for _, team := range collegeTeams {
		remainingGamesLeftByTeam[team.ID] = maxNumberOfGamesPerTeam
		homeGamesByTeam[team.ID] = 0
		awayGamesByTeam[team.ID] = 0
		opponentsFacedMap[team.ID] = make(map[uint]bool)
		// Assume teams with a coach are user teams
		isUserTeam[team.ID] = team.Coach != "" && team.Coach != "AI"
	}

	// Process existing matches to update counters
	allActiveMatches := repository.FindCollegeGamesRecords(repository.GamesQuery{SeasonID: seasonID})
	// Check all active games to ensure no team is scheduled twice in the same timeslot
	// teamSlotSeen: teamID -> "week:timeSlot" -> matchID
	teamSlotSeen := make(map[uint]uint)
	for _, match := range allActiveMatches {
		if match.AwayTeamID == 0 && match.HomeTeamID == 0 {
			continue
		}
		weekTimeSlotCounts[uint(match.Week)]++
		for _, teamID := range []uint{uint(match.HomeTeamID), uint(match.AwayTeamID)} {
			if firstMatchID, exists := teamSlotSeen[teamID]; exists {
				fmt.Printf("DOUBLE-SCHEDULED: TeamID %d appears in Week %d in both MatchID %d and MatchID %d\n",
					teamID, match.Week, firstMatchID, match.ID)
			} else {
				teamSlotSeen[teamID] = match.ID
			}
		}
	}

	for _, match := range allActiveMatches {
		if match.AwayTeamID == 0 && match.HomeTeamID == 0 {
			continue
		}
		// Initialize nested maps if needed, guarding against team IDs not present in collegeTeams
		gamesScheduledPerTeam[uint(match.HomeTeamID)]++
		gamesScheduledPerTeam[uint(match.AwayTeamID)]++

		homeTeam := teamMap[uint(match.HomeTeamID)]
		awayTeam := teamMap[uint(match.AwayTeamID)]

		if homeTeam.IsFBS && !awayTeam.IsFBS {
			playedFCSOpponent[homeTeam.ID] = true
			playedFBSOpponent[awayTeam.ID] = true
		}
		if !homeTeam.IsFBS && awayTeam.IsFBS {
			playedFBSOpponent[homeTeam.ID] = true
			playedFCSOpponent[awayTeam.ID] = true
		}

		// Mark time slots as occupied
		if opponentsFacedMap[uint(match.HomeTeamID)] == nil {
			opponentsFacedMap[uint(match.HomeTeamID)] = make(map[uint]bool)
		}
		if opponentsFacedMap[uint(match.AwayTeamID)] == nil {
			opponentsFacedMap[uint(match.AwayTeamID)] = make(map[uint]bool)
		}
		opponentsFacedMap[uint(match.HomeTeamID)][uint(match.AwayTeamID)] = true
		opponentsFacedMap[uint(match.AwayTeamID)][uint(match.HomeTeamID)] = true

		// Update remaining games and home/away counts
		remainingGamesLeftByTeam[uint(match.HomeTeamID)]--
		remainingGamesLeftByTeam[uint(match.AwayTeamID)]--
		homeGamesByTeam[uint(match.HomeTeamID)]++
		awayGamesByTeam[uint(match.AwayTeamID)]++
	}

	matchesToUpload := []structs.CollegeGame{}

	// Iterate through each week and time slot
	for _, week := range oocWeeks {
		var availableTeams []structs.CollegeTeam
		for _, team := range collegeTeams {
			// Check if team needs games and has the time slot available
			if remainingGamesLeftByTeam[team.ID] > 0 {
				availableTeams = append(availableTeams, team)

			}
		}

		// Shuffle available teams for randomness
		availableTeams = ShuffleTeams(availableTeams)

		// Track which teams have been matched in this slot
		matchedInSlot := make(map[uint]bool)

		// Try to pair teams
		for i := 0; i < len(availableTeams); i++ {
			team := availableTeams[i]

			// Skip if already matched in this slot
			if matchedInSlot[team.ID] {
				continue
			}

			// Find an opponent
			opponentIndex := -1
			for j := i + 1; j < len(availableTeams); j++ {
				opponent := availableTeams[j]

				// Check if valid opponent: different conference, not yet matched, has games remaining, and haven't faced each other
				if !matchedInSlot[opponent.ID] &&
					opponent.ConferenceID != team.ConferenceID &&
					remainingGamesLeftByTeam[opponent.ID] > 0 &&
					!opponentsFacedMap[team.ID][opponent.ID] {
					// Ensure FBS teams haven't already scheduled an FCS opponent and vice versa
					if team.IsFBS && !opponent.IsFBS && playedFCSOpponent[team.ID] {
						continue
					}
					if !team.IsFBS && opponent.IsFBS && playedFBSOpponent[team.ID] {
						continue
					}

					opponentIndex = j
					break
				}
			}

			// If no opponent found, this attempt fails
			if opponentIndex == -1 {
				continue // Try next team
			}

			opponent := availableTeams[opponentIndex]

			// Decide which team is home based on home/away balance
			var homeTeam, awayTeam structs.CollegeTeam

			// Prefer to balance user teams first
			if isUserTeam[team.ID] && !isUserTeam[opponent.ID] {
				// User team: prioritize balance
				if homeGamesByTeam[team.ID] <= awayGamesByTeam[team.ID] {
					homeTeam, awayTeam = team, opponent
				} else {
					homeTeam, awayTeam = opponent, team
				}
			} else if !isUserTeam[team.ID] && isUserTeam[opponent.ID] {
				// Opponent is user team: prioritize their balance
				if homeGamesByTeam[opponent.ID] <= awayGamesByTeam[opponent.ID] {
					homeTeam, awayTeam = opponent, team
				} else {
					homeTeam, awayTeam = team, opponent
				}
			} else {
				// Both user or both AI: balance based on current counts
				teamHomeDeficit := homeGamesByTeam[team.ID] - awayGamesByTeam[team.ID]
				oppHomeDeficit := homeGamesByTeam[opponent.ID] - awayGamesByTeam[opponent.ID]

				if teamHomeDeficit < oppHomeDeficit {
					homeTeam, awayTeam = team, opponent
				} else if teamHomeDeficit > oppHomeDeficit {
					homeTeam, awayTeam = opponent, team
				} else {
					// Equal deficit, randomize
					if rand.Intn(2) == 0 {
						homeTeam, awayTeam = team, opponent
					} else {
						homeTeam, awayTeam = opponent, team
					}
				}
			}

			// Create match
			match := MakeCollegeGameRecord(homeTeam, awayTeam, week, uint(ts.CollegeSeasonID), stadiumMap, stadiumMapByID, rivalryMap)

			matchesToUpload = append(matchesToUpload, match)

			// Update tracking data
			matchedInSlot[team.ID] = true
			matchedInSlot[opponent.ID] = true

			remainingGamesLeftByTeam[homeTeam.ID]--
			remainingGamesLeftByTeam[awayTeam.ID]--

			homeGamesByTeam[homeTeam.ID]++
			awayGamesByTeam[awayTeam.ID]++

			// Mark that these teams have faced each other
			opponentsFacedMap[homeTeam.ID][awayTeam.ID] = true
			opponentsFacedMap[awayTeam.ID][homeTeam.ID] = true

			// Update FBS/FCS opponent tracking
			if homeTeam.IsFBS && !awayTeam.IsFBS {
				playedFCSOpponent[homeTeam.ID] = true
				playedFBSOpponent[awayTeam.ID] = true
			}
			if !homeTeam.IsFBS && awayTeam.IsFBS {
				playedFBSOpponent[homeTeam.ID] = true
				playedFCSOpponent[awayTeam.ID] = true
			}
		}

		// Log teams that were available this slot but could not be paired
		for _, team := range availableTeams {
			if !matchedInSlot[team.ID] && remainingGamesLeftByTeam[team.ID] > 0 {
				fmt.Printf("Week %d: Could not schedule TeamID %d (%s)\n", week, team.ID, team.TeamAbbr)
			}
		}

	}

	// Cleanup pass: schedule any remaining games with relaxed constraints.
	// The only requirement is that the two teams are from different conferences;
	// the opponentsFacedMap check is intentionally dropped to allow rematches.
	fmt.Println("Starting cleanup scheduling pass for remaining unscheduled games...")
	for _, week := range oocWeeks {
		var cleanupAvailable []structs.CollegeTeam
		for _, team := range collegeTeams {
			if remainingGamesLeftByTeam[team.ID] > 0 {
				cleanupAvailable = append(cleanupAvailable, team)

			}
		}

		if len(cleanupAvailable) < 2 {
			continue
		}

		cleanupAvailable = ShuffleTeams(cleanupAvailable)
		cleanupMatchedInSlot := make(map[uint]bool)

		for i := 0; i < len(cleanupAvailable); i++ {
			team := cleanupAvailable[i]
			if cleanupMatchedInSlot[team.ID] {
				continue
			}

			opponentIndex := -1
			for j := i + 1; j < len(cleanupAvailable); j++ {
				opponent := cleanupAvailable[j]
				// Relaxed: only require different conference and no previous matchup
				if !cleanupMatchedInSlot[opponent.ID] &&
					opponent.ConferenceID != team.ConferenceID &&
					remainingGamesLeftByTeam[opponent.ID] > 0 &&
					!opponentsFacedMap[team.ID][opponent.ID] {

					// Ensure FBS teams haven't already scheduled an FCS opponent and vice versa
					if team.IsFBS && !opponent.IsFBS && playedFCSOpponent[team.ID] {
						continue
					}
					if !team.IsFBS && opponent.IsFBS && playedFBSOpponent[team.ID] {
						continue
					}

					opponentIndex = j
					break
				}
			}

			if opponentIndex == -1 {
				continue
			}

			opponent := cleanupAvailable[opponentIndex]

			var homeTeam, awayTeam structs.CollegeTeam
			if isUserTeam[team.ID] && !isUserTeam[opponent.ID] {
				if homeGamesByTeam[team.ID] <= awayGamesByTeam[team.ID] {
					homeTeam, awayTeam = team, opponent
				} else {
					homeTeam, awayTeam = opponent, team
				}
			} else if !isUserTeam[team.ID] && isUserTeam[opponent.ID] {
				if homeGamesByTeam[opponent.ID] <= awayGamesByTeam[opponent.ID] {
					homeTeam, awayTeam = opponent, team
				} else {
					homeTeam, awayTeam = team, opponent
				}
			} else {
				teamHomeDeficit := homeGamesByTeam[team.ID] - awayGamesByTeam[team.ID]
				oppHomeDeficit := homeGamesByTeam[opponent.ID] - awayGamesByTeam[opponent.ID]
				if teamHomeDeficit < oppHomeDeficit {
					homeTeam, awayTeam = team, opponent
				} else if teamHomeDeficit > oppHomeDeficit {
					homeTeam, awayTeam = opponent, team
				} else {
					if rand.Intn(2) == 0 {
						homeTeam, awayTeam = team, opponent
					} else {
						homeTeam, awayTeam = opponent, team
					}
				}
			}

			match := MakeCollegeGameRecord(homeTeam, awayTeam, week, uint(ts.CollegeSeasonID), stadiumMap, stadiumMapByID, rivalryMap)

			matchesToUpload = append(matchesToUpload, match)

			cleanupMatchedInSlot[team.ID] = true
			cleanupMatchedInSlot[opponent.ID] = true

			remainingGamesLeftByTeam[homeTeam.ID]--
			remainingGamesLeftByTeam[awayTeam.ID]--

			homeGamesByTeam[homeTeam.ID]++
			awayGamesByTeam[awayTeam.ID]++

			opponentsFacedMap[homeTeam.ID][awayTeam.ID] = true
			opponentsFacedMap[awayTeam.ID][homeTeam.ID] = true

			// Update FBS/FCS opponent tracking
			if homeTeam.IsFBS && !awayTeam.IsFBS {
				playedFCSOpponent[homeTeam.ID] = true
				playedFBSOpponent[awayTeam.ID] = true
			}
			if !homeTeam.IsFBS && awayTeam.IsFBS {
				playedFBSOpponent[homeTeam.ID] = true
				playedFCSOpponent[awayTeam.ID] = true
			}

			fmt.Printf("Cleanup: Scheduled Week %d: TeamID %d (%s) vs TeamID %d (%s)\n",
				week, homeTeam.ID, homeTeam.TeamAbbr, awayTeam.ID, awayTeam.TeamAbbr)
		}

	}

	// Validate the schedule - check if all teams that need OOC games got enough
	expectedGamesPerTeam := len(oocWeeks) // 4 weeks * 2 slots = 8 games expected per team

	for _, match := range matchesToUpload {
		gamesScheduledPerTeam[uint(match.HomeTeamID)]++
		gamesScheduledPerTeam[uint(match.AwayTeamID)]++
	}

	// Count teams with insufficient games
	insufficientTeams := 0
	for _, team := range collegeTeams {
		if gamesScheduledPerTeam[team.ID] < expectedGamesPerTeam {
			insufficientTeams++
		}
	}

	for _, match := range matchesToUpload {
		weekTimeSlotCounts[uint(match.Week)]++
	}

	// Print out the distribution of games per week and time slot
	for week, count := range weekTimeSlotCounts {
		fmt.Printf("Week %d: %d games scheduled (expected ~%d)\n", week, count, numberOfGamesExpectedPerSlot)
	}

	// If too many teams don't have enough games, consider this attempt failed
	// No, we need a perfect schedule generated.
	tolerance := 0
	if insufficientTeams > tolerance {
		return matchesToUpload, fmt.Errorf("incomplete schedule: %d teams have insufficient games", insufficientTeams)
	}

	return matchesToUpload, nil
}
