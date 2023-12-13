package managers

import (
	"log"
	"sort"
	"strconv"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/models"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
)

func ToggleDraftTime() {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	ts.ToggleDraftTime()

	db.Save(&ts)
}

// Gets all Current Season and Beyond Draft Picks
func GetDraftPicksByTeamID(TeamID string) []structs.NFLDraftPick {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	seasonID := strconv.Itoa(int(ts.NFLSeasonID))
	var picks []structs.NFLDraftPick

	db.Where("team_id = ? AND season_id >= ?", TeamID, seasonID).Find(&picks)

	return picks
}

// Gets all Current Season and Beyond Draft Picks
func GetDraftPickByDraftPickID(DraftPickID string) structs.NFLDraftPick {
	db := dbprovider.GetInstance().GetDB()

	var pick structs.NFLDraftPick

	db.Where("id = ?", DraftPickID).Find(&pick)

	return pick
}

func GetAllCurrentSeasonDraftPicks() []structs.NFLDraftPick {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	draftPicks := []structs.NFLDraftPick{}

	db.Order("draft_number asc").Where("season_id = ?", strconv.Itoa(int(ts.NFLSeasonID))).Find(&draftPicks)

	return draftPicks
}

func GetAllCurrentSeasonDraftPicksForDraftRoom() [7][]structs.NFLDraftPick {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	draftPicks := []structs.NFLDraftPick{}

	db.Order("draft_number asc").Where("season_id = ?", strconv.Itoa(int(ts.NFLSeasonID))).Find(&draftPicks)

	draftList := [7][]structs.NFLDraftPick{}
	for _, pick := range draftPicks {
		roundIdx := int(pick.DraftRound) - 1
		if roundIdx >= 0 && roundIdx < len(draftList) {
			draftList[roundIdx] = append(draftList[roundIdx], pick)
		} else {
			log.Panicln("Invalid round to insert pick!")
		}

	}

	return draftList
}

func GetOnlyNFLWarRoomByTeamID(TeamID string) models.NFLWarRoom {
	db := dbprovider.GetInstance().GetDB()

	warRoom := models.NFLWarRoom{}

	err := db.
		Where("team_id = ?", TeamID).Find(&warRoom).Error
	if err != nil {
		return warRoom
	}

	return warRoom
}

func GetNFLWarRoomByTeamID(TeamID string) models.NFLWarRoom {
	db := dbprovider.GetInstance().GetDB()

	warRoom := models.NFLWarRoom{}
	ts := GetTimestamp()

	err := db.Preload("DraftPicks", "season_id = ?", strconv.Itoa(int(ts.NFLSeasonID))).
		Preload("ScoutProfiles.Draftee").
		Preload("ScoutProfiles", "removed_from_board = ?", false).
		Where("team_id = ?", TeamID).Find(&warRoom).Error
	if err != nil {
		return warRoom
	}

	return warRoom
}

func GetNFLDrafteesForDraftPage() []models.NFLDraftee {
	db := dbprovider.GetInstance().GetDB()
	draftees := []models.NFLDraftee{}

	db.Find(&draftees)

	sort.Slice(draftees, func(i, j int) bool {
		iVal := util.GetNumericalSortValueByLetterGrade(draftees[i].OverallGrade)
		jVal := util.GetNumericalSortValueByLetterGrade(draftees[j].OverallGrade)
		return iVal < jVal
	})

	return draftees
}

func GetScoutProfileByScoutProfileID(profileID string) models.ScoutingProfile {
	db := dbprovider.GetInstance().GetDB()

	var scoutProfile models.ScoutingProfile

	err := db.Where("id = ?", profileID).Find(&scoutProfile).Error
	if err != nil {
		return models.ScoutingProfile{}
	}

	return scoutProfile
}

func GetOnlyScoutProfileByPlayerIDandTeamID(playerID, teamID string) models.ScoutingProfile {
	db := dbprovider.GetInstance().GetDB()

	var scoutProfile models.ScoutingProfile

	err := db.Where("player_id = ? AND team_id = ?", playerID, teamID).Error
	if err != nil {
		return models.ScoutingProfile{}
	}

	return scoutProfile
}

func CreateScoutingProfile(dto models.ScoutingProfileDTO) models.ScoutingProfile {
	db := dbprovider.GetInstance().GetDB()

	scoutProfile := GetOnlyScoutProfileByPlayerIDandTeamID(strconv.Itoa(int(dto.PlayerID)), strconv.Itoa(int(dto.TeamID)))

	// If Recruit Already Exists
	if scoutProfile.PlayerID > 0 && scoutProfile.TeamID > 0 {
		scoutProfile.ReplaceOnBoard()
		db.Save(&scoutProfile)
		return scoutProfile
	}

	newScoutingProfile := models.ScoutingProfile{
		PlayerID:         dto.PlayerID,
		TeamID:           dto.TeamID,
		ShowCount:        0,
		RemovedFromBoard: false,
	}

	db.Create(&newScoutingProfile)

	return newScoutingProfile
}

func RemovePlayerFromScoutBoard(id string) {
	db := dbprovider.GetInstance().GetDB()

	scoutProfile := GetScoutProfileByScoutProfileID(id)

	scoutProfile.RemoveFromBoard()

	db.Save(&scoutProfile)
}

func GetScoutingDataByPlayerID(id string) models.ScoutingDataResponse {
	ts := GetTimestamp()

	lastSeasonID := ts.NFLSeasonID - 1
	lastSeasonIDSTR := strconv.Itoa(int(lastSeasonID))

	draftee := GetHistoricCollegePlayerByID(id)

	seasonStats := GetCollegeSeasonStatsByPlayerAndSeason(id, lastSeasonIDSTR)
	teamID := strconv.Itoa(int(draftee.TeamID))
	collegeStandings := GetCollegeStandingsRecordByTeamID(teamID, lastSeasonIDSTR)

	return models.ScoutingDataResponse{
		DrafteeSeasonStats: seasonStats,
		TeamStandings:      collegeStandings,
	}
}

func RevealScoutingAttribute(dto models.RevealAttributeDTO) bool {
	db := dbprovider.GetInstance().GetDB()

	scoutProfile := GetScoutProfileByScoutProfileID(strconv.Itoa(int(dto.ScoutProfileID)))

	if scoutProfile.ID == 0 {
		return false
	}

	scoutProfile.RevealAttribute(dto.Attribute)

	warRoom := GetOnlyNFLWarRoomByTeamID(strconv.Itoa(int(dto.TeamID)))

	if warRoom.ID == 0 || warRoom.SpentPoints >= 100 || warRoom.SpentPoints+dto.Points > 100 {
		return false
	}

	warRoom.AddToSpentPoints(dto.Points)

	err := db.Save(&scoutProfile).Error
	if err != nil {
		return false
	}
	err = db.Save(&warRoom).Error
	return err == nil
}

func ExportDraftedPlayers(picks []structs.NFLDraftPick) bool {
	db := dbprovider.GetInstance().GetDB()

	for _, pick := range picks {
		playerId := strconv.Itoa(int(pick.SelectedPlayerID))
		draftee := GetNFLDrafteeByID(playerId)

		draftee.AssignDraftedTeam(pick.DraftNumber, pick.ID, pick.TeamID, pick.Team)

		NFLPlayer := structs.NFLPlayer{
			BasePlayer:    draftee.BasePlayer, // Assuming BasePlayer fields are common
			PlayerID:      draftee.PlayerID,
			TeamID:        int(pick.TeamID),
			TeamAbbr:      pick.Team,
			CollegeID:     draftee.CollegeID,
			College:       draftee.College,
			DraftPickID:   pick.ID,
			DraftedTeamID: pick.TeamID,
			DraftedTeam:   pick.Team,
			PrimeAge:      draftee.PrimeAge,
		}

		NFLPlayer.SetID(pick.SelectedPlayerID)

		year1Salary := util.GetDrafteeSalary(pick.DraftNumber, 1, pick.DraftRound, true)
		year2Salary := util.GetDrafteeSalary(pick.DraftNumber, 2, pick.DraftRound, true)
		year3Salary := util.GetDrafteeSalary(pick.DraftNumber, 3, pick.DraftRound, true)
		year4Salary := util.GetDrafteeSalary(pick.DraftNumber, 4, pick.DraftRound, true)
		year1Bonus := util.GetDrafteeSalary(pick.DraftNumber, 1, pick.DraftRound, false)
		year2Bonus := util.GetDrafteeSalary(pick.DraftNumber, 2, pick.DraftRound, false)
		year3Bonus := util.GetDrafteeSalary(pick.DraftNumber, 3, pick.DraftRound, false)
		year4Bonus := util.GetDrafteeSalary(pick.DraftNumber, 4, pick.DraftRound, false)
		yearsRemaining := 4
		contract := structs.NFLContract{
			PlayerID:       NFLPlayer.PlayerID,
			TeamID:         uint(NFLPlayer.TeamID),
			Team:           NFLPlayer.TeamAbbr,
			OriginalTeamID: uint(NFLPlayer.TeamID),
			OriginalTeam:   NFLPlayer.TeamAbbr,
			ContractLength: yearsRemaining,
			ContractType:   "Rookie",
			Y1BaseSalary:   year1Salary,
			Y2BaseSalary:   year2Salary,
			Y3BaseSalary:   year3Salary,
			Y4BaseSalary:   year4Salary,
			Y1Bonus:        year1Bonus,
			Y2Bonus:        year2Bonus,
			Y3Bonus:        year3Bonus,
			Y4Bonus:        year4Bonus,

			IsActive: true,
		}

		db.Create(&contract)

		db.Create(&NFLPlayer)
		db.Save(&draftee)
	}

	draftablePlayers := GetAllNFLDraftees()

	for _, draftee := range draftablePlayers {
		if draftee.DraftPickID > 0 {
			continue
		}

		nflPlayer := structs.NFLPlayer{
			BasePlayer:        draftee.BasePlayer, // Assuming BasePlayer fields are common
			PlayerID:          draftee.PlayerID,
			TeamID:            0,
			TeamAbbr:          "FA",
			CollegeID:         draftee.CollegeID,
			College:           draftee.College,
			DraftPickID:       0,
			DraftedTeamID:     0,
			DraftedTeam:       "UDFA",
			PrimeAge:          uint(draftee.PrimeAge),
			IsNegotiating:     false,
			IsAcceptingOffers: true,
			IsFreeAgent:       true,
			MinimumValue:      0.7,
		}

		nflPlayer.SetID(draftee.ID)

		NegotiationRound := 0
		if draftee.Overall < 80 {
			NegotiationRound = util.GenerateIntFromRange(2, 4)
		} else {
			NegotiationRound = util.GenerateIntFromRange(3, 6)
		}

		SigningRound := NegotiationRound + util.GenerateIntFromRange(2, 5)
		if SigningRound > 10 {
			SigningRound = 10
		}
		nflPlayer.AssignFAPreferences(uint(NegotiationRound), uint(SigningRound))

		db.Create(&nflPlayer)
	}

	return true
}

func BoomOrBust() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	SeasonID := strconv.Itoa(ts.NFLSeasonID)
	draftees := GetAllNFLDraftees()

	for _, player := range draftees {
		if player.BoomOrBust && player.BoomOrBustStatus == "Bust" {
			player = BoomBustDraftee(player, SeasonID, 51, false)
		} else {
			continue
		}
		// diceRoll := util.GenerateIntFromRange(1, 20)
		// if diceRoll == 1 {
		// 	// Bust
		// 	fmt.Println("BUST!")
		// 	player.AssignBoomBustStatus("Bust")
		// 	player = BoomBustDraftee(player, SeasonID, 51, false)
		// } else if diceRoll == 20 {
		// 	// Boom
		// 	fmt.Println("BOOM!")
		// 	player.AssignBoomBustStatus("Boom")
		// 	player = BoomBustDraftee(player, SeasonID, 51, true)
		// } else {
		// 	continue
		// }
		db.Save(&player)
	}
}
