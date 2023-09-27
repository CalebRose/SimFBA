package managers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/structs"
)

func GetAllTeamRequests() []structs.CreateRequestDTO {
	db := dbprovider.GetInstance().GetDB()
	var CollegeTeamRequests []structs.CreateRequestDTO
	// var NFLTeamRequests []structs.CreateRequestDTO
	var AllRequests []structs.CreateRequestDTO

	// College Team Requests
	db.Raw("SELECT team_requests.id, team_requests.team_id, college_teams.team_name, college_teams.team_abbr, team_requests.username, college_teams.conference, team_requests.is_approved FROM simfbaah_interface_3.team_requests INNER JOIN simfbaah_interface_3.college_teams on college_teams.id = team_requests.team_id WHERE team_requests.deleted_at is null AND team_requests.is_approved = 0").
		Scan(&CollegeTeamRequests)

	// NFL Team Requests
	// db.Raw("SELECT team_requests.id, team_requests.team_id, nfl_teams.team_name, nfl_teams.team_abbr, team_requests.username, nfl_teams.conference, team_requests.is_approved FROM simfbaah_interface_3.team_requests INNER JOIN simfbaah_interface_3.nfl_teams on nfl_teams.id = team_requests.team_id WHERE team_requests.deleted_at is null AND requests.is_approved = 0").
	// 	Scan(&NFLTeamRequests)

	// Append
	AllRequests = append(AllRequests, CollegeTeamRequests...)
	// AllRequests = append(AllRequests, NFLTeamRequests...)

	return AllRequests
}

func GetAllNFLTeamRequests() []structs.NFLRequest {
	db := dbprovider.GetInstance().GetDB()
	var NFLTeamRequests []structs.NFLRequest

	//NFL Team Requests
	db.Where("is_approved = false").Find(&NFLTeamRequests)

	return NFLTeamRequests
}

func CreateTeamRequest(request structs.TeamRequest) {
	db := dbprovider.GetInstance().GetDB()

	var ExistingTeamRequest structs.TeamRequest
	err := db.Where("username = ? AND team_id = ? AND is_approved = false AND deleted_at is null", request.Username, request.TeamID).Find(&ExistingTeamRequest).Error
	if err != nil {
		// Then there's no existing record, I guess? Which is fine.
		fmt.Println("Creating Team Request for TEAM " + strconv.Itoa(request.TeamID))
	}
	if ExistingTeamRequest.ID != 0 {
		// There is already an existing record.
		panic("There is already an existing request in place for the user. Please be patient while admin approves your formal request. If there is an issue, please reach out to TuscanSota.")
	}

	db.Create(&request)
}

func CreateNFLTeamRequest(request structs.NFLRequest) {
	db := dbprovider.GetInstance().GetDB()

	var existingRequest structs.NFLRequest
	err := db.Where("username = ? AND nfl_team_id = ? AND is_owner = ? AND is_manager = ? AND is_coach = ? AND is_assistant = ? AND is_approved = false AND deleted_at is null", request.Username, request.NFLTeamID, request.IsOwner, request.IsManager, request.IsCoach, request.IsAssistant).Find(&existingRequest).Error
	if err != nil {
		// Then there's no existing record, I guess? Which is fine.
		fmt.Println("Creating Team Request for TEAM " + strconv.Itoa(int(request.NFLTeamID)))
	}
	if existingRequest.ID != 0 {
		// There is already an existing record.
		log.Fatalln("There is already an existing request in place for the user. Please be patient while admin approves your formal request. If there is an issue, please reach out to TuscanSota.")
	}

	db.Create(&request)
}

func ApproveTeamRequest(request structs.TeamRequest) structs.TeamRequest {
	db := dbprovider.GetInstance().GetDB()

	timestamp := GetTimestamp()

	teamId := strconv.Itoa(request.TeamID)
	seasonID := strconv.Itoa(timestamp.CollegeSeasonID)

	// Approve Request
	request.ApproveTeamRequest()

	fmt.Println("Team Approved...")

	db.Save(&request)

	// Assign Team
	fmt.Println("Assigning team...")

	team := GetTeamByTeamID(teamId)

	coach := GetCollegeCoachByCoachName(request.Username)

	coach.SetTeam(request.TeamID)

	team.AssignUserToTeam(coach.CoachName)

	seasonalGames := GetCollegeGamesByTeamIdAndSeasonId(teamId, seasonID)

	for _, game := range seasonalGames {
		if game.Week >= timestamp.CollegeWeek {
			game.UpdateCoach(request.TeamID, coach.CoachName)
			db.Save(&game)
		}

	}

	standings := GetCFBStandingsByTeamIDAndSeasonID(teamId, seasonID)
	standings.SetCoach(coach.CoachName)
	db.Save(&standings)

	recruitingProfile := GetOnlyRecruitingProfileByTeamID(teamId)

	if recruitingProfile.IsAI {
		recruitingProfile.ActivateAI()
		db.Save(&recruitingProfile)
	}

	err := db.Save(&team).Error
	if err != nil {
		log.Fatalln("Could not assign user to team for some reason?")
	}

	db.Save(&coach)

	newsLog := structs.NewsLog{
		TeamID:      0,
		WeekID:      timestamp.CollegeWeekID,
		SeasonID:    timestamp.CollegeSeasonID,
		Week:        timestamp.CollegeWeek,
		MessageType: "CoachJob",
		League:      "CFB",
		Message:     "Breaking News! The " + team.TeamName + " " + team.Mascot + " have hired " + coach.CoachName + " as their new coach for the " + strconv.Itoa(timestamp.Season) + " season!",
	}

	db.Create(&newsLog)

	return request
}

func RejectTeamRequest(request structs.TeamRequest) {
	db := dbprovider.GetInstance().GetDB()

	request.RejectTeamRequest()

	err := db.Delete(&request).Error
	if err != nil {
		log.Fatalln("Could not delete request: " + err.Error())
	}
}

func ApproveNFLTeamRequest(request structs.NFLRequest) structs.NFLRequest {
	db := dbprovider.GetInstance().GetDB()

	timestamp := GetTimestamp()

	// Approve Request
	request.ApproveTeamRequest()

	fmt.Println("Team Approved...")

	db.Save(&request)

	// Assign Team
	fmt.Println("Assigning team...")

	team := GetNFLTeamByTeamID(strconv.Itoa(int(request.NFLTeamID)))

	coach := GetNFLUserByUsername(request.Username)

	coach.SetTeam(request)

	team.AssignNFLUserToTeam(request, coach)

	// seasonalGames := GetCollegeGamesByTeamIdAndSeasonId(strconv.Itoa(request.TeamID), strconv.Itoa(timestamp.CollegeSeasonID))

	// for _, game := range seasonalGames {
	// 	if game.Week >= timestamp.CollegeWeek {
	// 		game.UpdateCoach(int(request.NFLTeamID), request.Username)
	// 		db.Save(&game)
	// 	}
	// }

	db.Save(&team)

	db.Save(&coach)

	newsLog := structs.NewsLog{
		TeamID:      0,
		WeekID:      timestamp.NFLWeekID,
		SeasonID:    timestamp.NFLSeasonID,
		Week:        timestamp.NFLWeek,
		MessageType: "CoachJob",
		League:      "NFL",
		Message:     "Breaking News! The " + team.TeamName + " " + team.Mascot + " have hired " + coach.Username + " to their staff for the " + strconv.Itoa(timestamp.Season) + " season!",
	}

	db.Create(&newsLog)

	return request
}

func RejectNFLTeamRequest(request structs.NFLRequest) {
	db := dbprovider.GetInstance().GetDB()

	request.RejectTeamRequest()

	err := db.Delete(&request).Error
	if err != nil {
		log.Fatalln("Could not delete request: " + err.Error())
	}
}

func RemoveUserFromTeam(teamId string) {
	db := dbprovider.GetInstance().GetDB()

	team := GetTeamByTeamID(teamId)

	coach := GetCollegeCoachByCoachName(team.Coach)

	coach.SetAsInactive()

	team.RemoveUserFromTeam()

	db.Save(&team)

	db.Save(&coach)

	timestamp := GetTimestamp()
	seasonID := strconv.Itoa(int(timestamp.CollegeSeasonID))
	seasonalGames := GetCollegeGamesByTeamIdAndSeasonId(teamId, seasonID)

	for _, game := range seasonalGames {
		if game.Week >= timestamp.CollegeWeek {
			game.UpdateCoach(int(team.ID), "AI")
			db.Save(&game)
		}

	}

	standings := GetCFBStandingsByTeamIDAndSeasonID(teamId, seasonID)
	standings.SetCoach("AI")
	db.Save(&standings)

	recruitingProfile := GetOnlyRecruitingProfileByTeamID(teamId)

	if !recruitingProfile.IsAI {
		recruitingProfile.ActivateAI()
	}

	db.Save(&recruitingProfile)

	newsLog := structs.NewsLog{
		TeamID:      0,
		WeekID:      timestamp.CollegeWeekID,
		SeasonID:    timestamp.CollegeSeasonID,
		Week:        timestamp.CollegeWeek,
		MessageType: "CoachJob",
		League:      "CFB",
		Message:     coach.CoachName + " has decided to step down as the head coach of the " + team.TeamName + " " + team.Mascot + "!",
	}

	db.Create(&newsLog)
}

func RemoveUserFromNFLTeam(request structs.NFLRequest) {
	db := dbprovider.GetInstance().GetDB()

	teamID := strconv.Itoa(int(request.NFLTeamID))

	team := GetNFLTeamByTeamID(teamID)

	user := GetNFLUserByUsername(request.Username)

	message := ""

	if request.IsOwner {
		user.RemoveOwnership()
		message = request.Username + " has decided to step down as Owner of the " + team.TeamName + " " + team.Mascot + "!"
	}

	if request.IsManager {
		user.RemoveManagerPosition()
		message = request.Username + " has decided to step down as Manager of the " + team.TeamName + " " + team.Mascot + "!"
	}

	if request.IsCoach {
		user.RemoveCoachPosition()
		message = request.Username + " has decided to step down as Head Coach of the " + team.TeamName + " " + team.Mascot + "!"
	}

	if request.IsAssistant {
		user.RemoveAssistantPosition()
		message = request.Username + " has decided to step down as an Assistant of the " + team.TeamName + " " + team.Mascot + "!"
	}

	team.RemoveNFLUserFromTeam(request, user)

	db.Save(&team)

	db.Save(&user)

	timestamp := GetTimestamp()

	newsLog := structs.NewsLog{
		TeamID:      0,
		WeekID:      timestamp.NFLWeekID,
		SeasonID:    timestamp.NFLSeasonID,
		Week:        timestamp.NFLWeek,
		MessageType: "CoachJob",
		Message:     message,
		League:      "NFL",
	}

	db.Create(&newsLog)
}
