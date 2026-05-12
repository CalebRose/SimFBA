package managers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/firebase"
	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
)

// ─────────────────────────────────────────────
// CFB Game Request
// ─────────────────────────────────────────────

// CreateCFBGameRequest saves a new CFBGameRequest record to the database.
func CreateCFBGameRequest(request structs.CFBGameRequest) {
	db := dbprovider.GetInstance().GetDB()
	repository.CreateCFBGameRequest(request, db)
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
	stadium := stadiumByID[request.ArenaID]

	isDivisional := homeTeam.DivisionID > 0 && homeTeam.DivisionID == awayTeam.DivisionID

	game := structs.CollegeGame{
		HomeTeamID:   int(request.HomeTeamID),
		HomeTeam:     homeTeam.TeamName,
		AwayTeamID:   int(request.AwayTeamID),
		AwayTeam:     awayTeam.TeamName,
		Week:         int(request.Week),
		WeekID:       int(request.WeekID),
		SeasonID:     int(request.SeasonID),
		StadiumID:    request.ArenaID,
		Stadium:      stadium.StadiumName,
		City:         stadium.City,
		State:        stadium.State,
		Region:       stadium.Region,
		TimeSlot:     request.Timeslot,
		IsNeutral:    request.IsNeutralSite,
		IsConference: homeTeam.ConferenceID == awayTeam.ConferenceID,
		IsDivisional: isDivisional,
		IsSpringGame: request.IsSpringGame,
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
