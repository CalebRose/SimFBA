package managers

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
	"github.com/jinzhu/gorm"
)

func ImportRecruitAICSV() {
	db := dbprovider.GetInstance().GetDB()
	completedCrootPath := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\FCS_Croot_Weekly_Signings.csv"
	aiPoolPath := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\2022_Croot_Class_AI.csv"
	crootMap := make(map[string][]string)
	f, err := os.Open(completedCrootPath)
	if err != nil {
		log.Fatal("Unable to read input file "+completedCrootPath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	croots, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+completedCrootPath, err)
	}

	for idx, record := range croots {
		if idx == 0 {
			continue
		}
		// Add recruit to map
		crootMap[record[0]] = record
		id := util.ConvertStringToInt(record[0])
		teamID := util.ConvertStringToInt(record[18])
		points := util.ConvertStringToInt(record[21])
		if points <= 0 {
			points = 1
		}

		if teamID > 0 {
			recruitProfile := structs.RecruitPlayerProfile{
				RecruitID:        id,
				IsSigned:         true,
				Scholarship:      false,
				TotalPoints:      float64(points),
				ProfileID:        teamID,
				TeamAbbreviation: record[19],
				SeasonID:         2,
			}

			db.Create(&recruitProfile)
			recruit := GetCollegeRecruitByRecruitID(record[0])
			// Since this croot is not in the DB yet, they will be added later
			if recruit.ID == 0 {
				continue
			}
			recruit.AssignCollege(recruitProfile.TeamAbbreviation)
			recruit.UpdateTeamID(teamID)
			recruit.UpdateSigningStatus()

			db.Save(&recruit)
		}
	}

	fmt.Println("NOW GET THE BIG ONE")

	poolFile, err := os.Open(aiPoolPath)
	if err != nil {
		log.Fatal("Unable to read input file "+aiPoolPath, err)
	}
	defer f.Close()

	csvReader = csv.NewReader(poolFile)
	croots, err = csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+aiPoolPath, err)
	}

	for idx, croot := range croots {
		if idx == 0 {
			continue
		}
		recruit := GetCollegeRecruitByRecruitID(croot[0])
		if recruit.ID == 0 {
			// Create the record
			idStr := croot[0]
			id := util.ConvertStringToInt(croot[0])
			// If for some reason a recruit was not included in the map, then they don't have a team. Skip over for now.

			var mapRecord []string
			if len(crootMap[idStr]) > 0 {
				// Retrieve for the Team Info
				mapRecord = crootMap[idStr]
			}

			teamIDStr := ""
			abbr := ""
			if len(mapRecord) > 0 {
				teamIDStr = mapRecord[18]
				abbr = mapRecord[20]
			}

			teamID := util.ConvertStringToInt(teamIDStr)
			isSigned := false
			if teamID > 0 {
				isSigned = true
			}

			// Attributes

			base := structs.BasePlayer{
				FirstName:      croot[1],
				LastName:       croot[2],
				Stars:          util.ConvertStringToInt(croot[3]),
				Position:       croot[4],
				Archetype:      croot[5],
				Overall:        util.ConvertStringToInt(croot[6]),
				Height:         util.ConvertStringToInt(croot[7]),
				Weight:         util.ConvertStringToInt(croot[8]),
				Carrying:       util.ConvertStringToInt(croot[12]),
				Agility:        util.ConvertStringToInt(croot[13]),
				Catching:       util.ConvertStringToInt(croot[14]),
				ZoneCoverage:   util.ConvertStringToInt(croot[15]),
				ManCoverage:    util.ConvertStringToInt(croot[16]),
				FootballIQ:     util.ConvertStringToInt(croot[17]),
				KickAccuracy:   util.ConvertStringToInt(croot[18]),
				KickPower:      util.ConvertStringToInt(croot[19]),
				PassBlock:      util.ConvertStringToInt(croot[20]),
				PassRush:       util.ConvertStringToInt(croot[21]),
				PuntAccuracy:   util.ConvertStringToInt(croot[22]),
				PuntPower:      util.ConvertStringToInt(croot[23]),
				RouteRunning:   util.ConvertStringToInt(croot[24]),
				RunBlock:       util.ConvertStringToInt(croot[25]),
				RunDefense:     util.ConvertStringToInt(croot[26]),
				Speed:          util.ConvertStringToInt(croot[27]),
				Strength:       util.ConvertStringToInt(croot[28]),
				Tackle:         util.ConvertStringToInt(croot[29]),
				ThrowPower:     util.ConvertStringToInt(croot[30]),
				ThrowAccuracy:  util.ConvertStringToInt(croot[31]),
				Injury:         util.ConvertStringToInt(croot[32]),
				Stamina:        util.ConvertStringToInt(croot[33]),
				Discipline:     util.ConvertStringToInt(croot[34]),
				AcademicBias:   croot[35],
				FreeAgency:     croot[36],
				Personality:    croot[37],
				RecruitingBias: croot[38],
				WorkEthic:      croot[39],
				Progression:    util.ConvertStringToInt(croot[40]),
				PotentialGrade: croot[41],
				Age:            18,
			}

			r := structs.Recruit{
				BasePlayer: base,
				PlayerID:   id,
				TeamID:     teamID,
				HighSchool: croot[9],
				City:       croot[10],
				State:      croot[11],
				IsSigned:   isSigned,
				College:    abbr,
			}

			r.AssignID(id)

			db.Create(&r)
		} else {
			// This recruit is already in the DB
			continue
		}
	}
}

func GetLeftoverRecruits() []structs.UnsignedPlayer {
	db := dbprovider.GetInstance().GetDB()
	var unsignedPlayers []structs.UnsignedPlayer

	db.Where("year = 1").Find(&unsignedPlayers)

	return unsignedPlayers
}

func ImportNFLPlayersCSV() {
	db := dbprovider.GetInstance().GetDB()
	playerPath := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\NFL_Progressed.csv"

	nflCSV := util.ReadCSV(playerPath)

	for idx, row := range nflCSV {
		if idx < 2 {
			continue
		}

		playerID := util.ConvertStringToInt(row[0])
		academic := row[35]
		fa := row[36]
		personality := row[37]
		recruit := row[38]
		we := row[39]
		progression := util.ConvertStringToInt(row[40])

		gp := GetGlobalPlayerRecord(row[0])
		if gp.ID == 0 {
			player := structs.Player{
				NFLPlayerID: playerID,
			}
			player.AssignID(uint(playerID))

			db.Create(&player)
		}

		NFLPlayerRecord := GetNFLPlayerRecord(row[0])
		if NFLPlayerRecord.ID == 0 {
			log.Fatalln("Something is wrong, this player was not uploaded.")
		}
		NFLPlayerRecord.AssignMissingValues(progression, academic, fa, personality, recruit, we)

		db.Save(&NFLPlayerRecord)
	}
}

// Imports 2028-2075 Draft Picks
func ImportNFLDraftPicks() {
	db := dbprovider.GetInstance().GetDB()
	nflTeams := GetAllNFLTeams()
	teamMap := make(map[uint]structs.NFLTeam)
	seasonID := 8
	season := 2028
	draftPicksToUpload := []structs.NFLDraftPick{}

	for _, team := range nflTeams {
		teamMap[team.ID] = team
	}

	for season < 2076 {
		draftNumber := 1
		for draftRound := 1; draftRound < 8; draftRound++ {
			for _, team := range nflTeams {
				tradeValue := util.GetDraftPickValue(draftRound, draftNumber)
				draftPick := structs.NFLDraftPick{
					SeasonID:       uint(seasonID),
					Season:         uint(season),
					DraftRound:     uint(draftRound),
					DraftNumber:    uint(draftNumber),
					Team:           util.GetNFLFullTeamName(team.TeamName, team.Mascot),
					TeamID:         team.ID,
					OriginalTeamID: team.ID,
					OriginalTeam:   util.GetNFLFullTeamName(team.TeamName, team.Mascot),
					Notes:          "",
					DraftValue:     tradeValue,
				}

				draftPicksToUpload = append(draftPicksToUpload, draftPick)
				draftNumber++
			}
		}
		seasonID++
		season++
	}
	repository.CreateNFLDraftPickBatch(db, draftPicksToUpload, 400)
}

func ImportMinimumFAValues() {
	db := dbprovider.GetInstance().GetDB()
	playerPath := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\2025_Free_Agency_Expected_Values.csv"

	nflCSV := util.ReadCSV(playerPath)

	nflPlayerMap := GetAllNFLPlayersMap()

	for idx, row := range nflCSV {
		if idx < 1 {
			continue
		}

		playerID := row[0]
		id := util.ConvertStringToInt(playerID)
		valueStr := strings.TrimSpace(row[7])
		aavStr := strings.TrimSpace(row[8])
		value := util.ConvertStringToFloat(valueStr)
		aav := util.ConvertStringToFloat(aavStr)

		NFLPlayerRecord := nflPlayerMap[uint(id)]
		if NFLPlayerRecord.ID == 0 {
			continue
		}

		NFLPlayerRecord.AssignMinimumValue(value, aav)

		repository.SaveNFLPlayer(NFLPlayerRecord, db)
	}
}

func ImportWorkEthic() {
	fmt.Println(time.Now().UnixNano())
	db := dbprovider.GetInstance().GetDB()

	nflPlayers := GetAllNFLPlayers()

	for _, p := range nflPlayers {
		WorkEthic := util.GetWorkEthic()
		if p.ID == 10 {
			FreeAgency := "Highly Unlikely to Play for the Miami Dolphins."
			Personality := "Worships Himself"
			p.AssignPersonality(Personality)
			p.AssignFreeAgency(FreeAgency)
		}

		p.AssignWorkEthic(WorkEthic)

		db.Save(&p)
	}
}

func ImportFAPreferences() {
	fmt.Println(time.Now().UnixNano())
	db := dbprovider.GetInstance().GetDB()

	nflPlayers := GetAllNFLPlayers()

	for _, p := range nflPlayers {
		NegotiationRound := 0
		if p.Overall > 70 {
			NegotiationRound = util.GenerateIntFromRange(2, 4)
		} else {
			NegotiationRound = util.GenerateIntFromRange(3, 6)
		}

		SigningRound := NegotiationRound + util.GenerateIntFromRange(2, 4)
		if SigningRound > 10 {
			SigningRound = 10
		}

		p.AssignFAPreferences(uint(NegotiationRound), uint(SigningRound))

		db.Save(&p)
	}
}

func RetireAndFreeAgentPlayers() {
	db := dbprovider.GetInstance().GetDB()

	nflPlayers := GetAllNFLPlayers()

	for _, record := range nflPlayers {

		if !record.IsActive {
			retiredPlayerRecord := (structs.NFLRetiredPlayer)(record)

			db.Create(&retiredPlayerRecord)
			db.Delete(&record)
			continue
		}

		if record.TeamID == 0 {
			record.ToggleIsFreeAgent()
			db.Save(&record)
		}
	}
}

func ImportTradePreferences() {
	db := dbprovider.GetInstance().GetDB()

	nflTeams := GetAllNFLTeams()

	for _, t := range nflTeams {

		pref := structs.NFLTradePreferences{
			NFLTeamID: t.ID,
		}

		db.Create(&pref)
	}
}

func Import2023DraftedPlayers() {
	db := dbprovider.GetInstance().GetDB()

	path := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\2023DraftList.csv"

	nflCSV := util.ReadCSV(path)

	nflTeams := GetAllNFLTeams()
	teamMap := make(map[string]uint)

	for _, team := range nflTeams {
		teamMap[team.TeamAbbr] = team.ID
	}

	for idx, draftee := range nflCSV {
		if idx == 0 {
			continue
		}

		team := draftee[0]
		teamID := teamMap[team]
		playerID := draftee[3]
		round := util.ConvertStringToInt(draftee[1])
		pickNumber := util.ConvertStringToInt(draftee[2])

		draftRecord := GetNFLDrafteeByPlayerID(playerID)

		nflPlayerRecord := structs.NFLPlayer{
			Model: gorm.Model{
				ID: draftRecord.ID,
			},
			BasePlayer:        draftRecord.BasePlayer,
			PlayerID:          int(draftRecord.ID),
			TeamID:            int(teamID),
			College:           draftRecord.College,
			TeamAbbr:          team,
			Experience:        1,
			HighSchool:        draftRecord.HighSchool,
			Hometown:          draftRecord.City,
			State:             draftRecord.State,
			IsActive:          true,
			IsPracticeSquad:   false,
			IsFreeAgent:       false,
			IsWaived:          false,
			IsOnTradeBlock:    false,
			IsAcceptingOffers: false,
			IsNegotiating:     false,
			DraftedTeamID:     teamID,
			DraftedTeam:       team,
			DraftedRound:      uint(round),
			DraftedPick:       uint(pickNumber),
			ShowLetterGrade:   true,
		}

		baseSalaryByYear := getBaseSalaryByYear(round, pickNumber)
		bonusByYear := getBonusByYear(round, pickNumber)

		contract := structs.NFLContract{
			PlayerID:       int(draftRecord.ID),
			NFLPlayerID:    int(draftRecord.ID),
			TeamID:         teamID,
			Team:           team,
			OriginalTeamID: teamID,
			OriginalTeam:   team,
			ContractLength: 4,
			Y1BaseSalary:   baseSalaryByYear,
			Y2BaseSalary:   baseSalaryByYear,
			Y3BaseSalary:   baseSalaryByYear,
			Y4BaseSalary:   baseSalaryByYear,
			Y1Bonus:        bonusByYear,
			Y2Bonus:        bonusByYear,
			Y3Bonus:        bonusByYear,
			Y4Bonus:        bonusByYear,
			ContractType:   "Rookie",
			IsActive:       true,
		}

		contract.CalculateContract()

		db.Create(&contract)
		db.Create(&nflPlayerRecord)
		db.Delete(&draftRecord)
	}
}

func UpdateDraftPicks() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	path := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\2026\\2026_draft_picks_upload.csv"

	draftPickCSV := util.ReadCSV(path)

	draftPicks := GetAllCurrentSeasonDraftPicks()
	pickMap := make(map[uint]structs.NFLDraftPick)
	var latestID uint = 1135 // Latest ID from Draft Pick table in DB

	for _, pick := range draftPicks {
		pickMap[pick.ID] = pick
	}

	for idx, row := range draftPickCSV {
		if idx == 0 {
			continue
		}

		draftPickID := util.ConvertStringToInt(row[0])
		draftRound := util.ConvertStringToInt(row[3])
		overallNumber := util.ConvertStringToInt(row[4])
		teamID := util.ConvertStringToInt(row[6])
		team := row[7]
		draftValue := util.ConvertStringToFloat(row[12])
		isCompensation := util.ConvertStringToBool(row[17])
		if isCompensation {
			draftPickID = int(latestID)
			latestID += 1
		}
		isVoid := util.ConvertStringToBool(row[18])
		draftPick := structs.NFLDraftPick{
			SeasonID: uint(ts.NFLSeasonID),
			Season:   uint(ts.Season),
		}
		if !isCompensation {
			draftPick = pickMap[uint(draftPickID)]
		}
		draftPick.MapValuesToDraftPick(uint(draftPickID), uint(draftRound), uint(overallNumber), uint(teamID), team, draftValue, isCompensation, isVoid)

		db.Save(&draftPick)
	}
}

func ImportUDFAs() {
	db := dbprovider.GetInstance().GetDB()

	UDFAs := GetAllNFLDraftees()

	for idx, draftee := range UDFAs {
		if idx == 0 {
			continue
		}

		team := "FA"
		teamID := 0

		nflPlayerRecord := structs.NFLPlayer{
			Model: gorm.Model{
				ID: draftee.ID,
			},
			BasePlayer:        draftee.BasePlayer,
			PlayerID:          int(draftee.ID),
			TeamID:            int(teamID),
			College:           draftee.College,
			TeamAbbr:          team,
			Experience:        1,
			HighSchool:        draftee.HighSchool,
			Hometown:          draftee.City,
			State:             draftee.State,
			IsActive:          true,
			IsPracticeSquad:   false,
			IsFreeAgent:       true,
			IsWaived:          false,
			IsOnTradeBlock:    false,
			IsAcceptingOffers: true,
			IsNegotiating:     false,
			ShowLetterGrade:   true,
		}

		db.Create(&nflPlayerRecord)
		db.Delete(&draftee)
	}
}

func ImportCFBGames() {
	db := dbprovider.GetInstance().GetDB()

	path := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\2026\\2026_cfb_games_preseason.csv"

	gamesCSV := util.ReadCSV(path)

	ts := GetTimestamp()

	teamMap := make(map[string]structs.CollegeTeam)

	allCollegeTeams := GetAllCollegeTeams()

	for _, t := range allCollegeTeams {
		teamMap[t.TeamAbbr] = t
	}

	for idx, row := range gamesCSV {
		if idx == 0 {
			continue
		}

		gameID := util.ConvertStringToInt(row[0])
		season := util.ConvertStringToInt(row[1])
		seasonID := season - 2020
		week := util.ConvertStringToInt(row[2])
		weekID := week + 88 // Week 43 is week 0 of the 2023 Season
		homeTeamAbbr := row[3]
		awayTeamAbbr := row[4]
		ht := teamMap[homeTeamAbbr]
		at := teamMap[awayTeamAbbr]
		homeTeamID := ht.ID
		awayTeamID := at.ID
		homeTeamCoach := ht.Coach
		awayTeamCoach := at.Coach
		timeSlot := row[16]
		// Need to implement Stadium ID
		stadium := row[17]
		city := row[18]
		state := row[19]
		isDomed := util.ConvertStringToBool(row[20])
		// Need to check for if a game is in a domed stadium or not
		isConferenceGame := ht.ConferenceID == at.ConferenceID
		isDivisionGame := isConferenceGame && ht.DivisionID == at.DivisionID && ht.DivisionID > 0
		conferenceID := 0
		if isConferenceGame {
			conferenceID = ht.ConferenceID
		}
		isNeutralSite := util.ConvertStringToBool(row[7])
		isConferenceChampionship := util.ConvertStringToBool(row[8])
		isBowlGame := util.ConvertStringToBool(row[9])
		isPlayoffGame := util.ConvertStringToBool(row[10])
		isNationalChampionship := util.ConvertStringToBool(row[11])
		isHomePrevBye := util.ConvertStringToBool(row[12])
		isAwayPrevBye := util.ConvertStringToBool(row[13])
		gameTitle := row[21]
		nextGame := util.ConvertStringToInt(row[23])
		nextGameHOA := row[24]

		game := structs.CollegeGame{
			Model:                    gorm.Model{ID: uint(gameID)},
			SeasonID:                 seasonID,
			WeekID:                   weekID,
			Week:                     week,
			HomeTeamID:               int(homeTeamID),
			AwayTeamID:               int(awayTeamID),
			HomeTeam:                 homeTeamAbbr,
			AwayTeam:                 awayTeamAbbr,
			HomeTeamCoach:            homeTeamCoach,
			AwayTeamCoach:            awayTeamCoach,
			IsConferenceChampionship: isConferenceChampionship,
			IsSpringGame:             ts.CFBSpringGames,
			IsBowlGame:               isBowlGame,
			IsNeutral:                isNeutralSite,
			IsPlayoffGame:            isPlayoffGame,
			IsNationalChampionship:   isNationalChampionship,
			IsConference:             isConferenceGame,
			IsDivisional:             isDivisionGame,
			TimeSlot:                 timeSlot,
			Stadium:                  stadium,
			City:                     city,
			State:                    state,
			IsDomed:                  isDomed,
			GameTitle:                gameTitle,
			NextGameID:               uint(nextGame),
			NextGameHOA:              nextGameHOA,
			HomePreviousBye:          isHomePrevBye,
			AwayPreviousBye:          isAwayPrevBye,
			ConferenceID:             uint(conferenceID),
		}

		db.Create(&game)
	}
	GenerateWeatherForGames()
}

func ImportNFLGames() {
	db := dbprovider.GetInstance().GetDB()

	path := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\2025\\2025_nfl_postseason_games.csv"

	gamesCSV := util.ReadCSV(path)

	ts := GetTimestamp()

	teamMap := make(map[string]structs.NFLTeam)

	allNFLTeams := GetAllNFLTeams()

	for _, t := range allNFLTeams {
		teamMap[t.TeamAbbr] = t
	}

	for idx, row := range gamesCSV {
		if idx == 0 {
			continue
		}

		gameID := util.ConvertStringToInt(row[0])
		season := util.ConvertStringToInt(row[1])
		seasonID := season - 2020
		week := util.ConvertStringToInt(row[2])
		weekID := week + 23 // Week 43 is week 0 of the 2024 Season
		if gameID > 381 {
			weekID = week + 51
		}
		homeTeamAbbr := row[3]
		awayTeamAbbr := row[4]
		ht := teamMap[homeTeamAbbr]
		at := teamMap[awayTeamAbbr]
		homeTeamName := ht.TeamName + " " + ht.Mascot
		awayTeamName := at.TeamName + " " + at.Mascot
		homeTeamID := ht.ID
		awayTeamID := at.ID
		homeTeamCoach := ht.NFLCoachName
		if len(homeTeamCoach) == 0 {
			homeTeamCoach = ht.NFLOwnerName
		}
		if len(homeTeamCoach) == 0 {
			homeTeamCoach = "AI"
		}
		awayTeamCoach := at.NFLCoachName
		if len(awayTeamCoach) == 0 {
			awayTeamCoach = at.NFLOwnerName
		}
		if len(awayTeamCoach) == 0 {
			awayTeamCoach = "AI"
		}
		timeSlot := row[18]
		// Need to implement Stadium ID
		stadium := ht.Stadium
		city := ht.City
		state := ht.State
		// Need to check for if a game is in a domed stadium or not
		isConferenceGame := util.ConvertStringToBool(row[9])
		isDivisionGame := util.ConvertStringToBool(row[10])
		isNeutralSite := util.ConvertStringToBool(row[11])
		// isPreseasonGame := util.ConvertStringToBool(row[12])
		// isConferenceChampionship := util.ConvertStringToBool(row[13])
		isPlayoffGame := util.ConvertStringToBool(row[14])
		isNationalChampionship := util.ConvertStringToBool(row[15])
		gameTitle := row[23]
		nextGame := util.ConvertStringToInt(row[24])
		nextGameHOA := row[25]

		game := structs.NFLGame{
			Model:           gorm.Model{ID: uint(gameID)},
			SeasonID:        seasonID,
			WeekID:          weekID,
			Week:            week,
			HomeTeamID:      int(homeTeamID),
			AwayTeamID:      int(awayTeamID),
			HomeTeam:        homeTeamName,
			AwayTeam:        awayTeamName,
			HomeTeamCoach:   homeTeamCoach,
			AwayTeamCoach:   awayTeamCoach,
			IsPreseasonGame: ts.NFLPreseason,
			IsNeutral:       isNeutralSite,
			IsConference:    isConferenceGame,
			IsDivisional:    isDivisionGame,
			TimeSlot:        timeSlot,
			Stadium:         stadium,
			City:            city,
			State:           state,
			GameTitle:       gameTitle,
			NextGameID:      uint(nextGame),
			NextGameHOA:     nextGameHOA,
			IsPlayoffGame:   isPlayoffGame,
			IsSuperBowl:     isNationalChampionship,
		}

		db.Create(&game)
	}

	GenerateWeatherForGames()
}

func ImportCFBTeams() {
	db := dbprovider.GetInstance().GetDB()

	teamPath := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\2024\\teams.csv"
	stadiumPath := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\2024\\stadia.csv"
	profilePath := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\2024\\profiles.csv"

	teamCSV := util.ReadCSV(teamPath)
	stadiumCSV := util.ReadCSV(stadiumPath)
	profileCSV := util.ReadCSV(profilePath)

	for idx, row := range teamCSV {
		if idx == 0 {
			continue
		}

		stadiumRecord := stadiumCSV[idx]
		profileRecord := profileCSV[idx]

		teamID := util.ConvertStringToInt(row[0])
		stadiumID := util.ConvertStringToInt(stadiumRecord[0])
		stadiumName := stadiumRecord[1]
		capacity := util.ConvertStringToInt(stadiumRecord[8])
		recordAtt := util.ConvertStringToInt(stadiumRecord[9])
		teamName := row[1]
		mascot := row[2]
		abbr := row[3]
		city := row[4]
		state := row[5]
		country := "USA"
		conferenceID := util.ConvertStringToInt(row[7])
		conference := row[8]
		firstSeason := 2027
		isFBS := false
		isActive := false

		stadium := structs.Stadium{
			Model: gorm.Model{
				ID: uint(stadiumID),
			},
			StadiumName:      stadiumName,
			TeamID:           uint(teamID),
			TeamAbbr:         abbr,
			City:             city,
			State:            state,
			Country:          country,
			Capacity:         uint(capacity),
			RecordAttendance: uint(recordAtt),
			FirstSeason:      uint(firstSeason),
			LeagueID:         2,
			LeagueName:       "FCS",
		}

		team := structs.CollegeTeam{
			Model: gorm.Model{
				ID: uint(teamID),
			},
			BaseTeam: structs.BaseTeam{
				TeamName:         teamName,
				Mascot:           mascot,
				TeamAbbr:         abbr,
				City:             city,
				State:            state,
				Country:          country,
				StadiumID:        uint(stadiumID),
				Stadium:          stadiumName,
				RecordAttendance: recordAtt,
				Enrollment:       0,
				FirstPlayed:      firstSeason,
			},
			ConferenceID: conferenceID,
			Conference:   conference,
			IsFBS:        isFBS,
			IsActive:     isActive,
		}

		aiBehavior := profileRecord[10]
		aiQuality := profileRecord[11]
		min := util.ConvertStringToInt(profileRecord[12])
		max := util.ConvertStringToInt(profileRecord[13])
		off := profileRecord[17]
		def := profileRecord[18]

		teamProfile := structs.RecruitingTeamProfile{
			Model: gorm.Model{
				ID: uint(teamID),
			},
			TeamID:                    teamID,
			Team:                      teamName,
			TeamAbbreviation:          abbr,
			State:                     state,
			ScholarshipsAvailable:     40,
			WeeklyPoints:              100,
			SpentPoints:               0,
			TotalCommitments:          0,
			RecruitClassSize:          20,
			PortalReputation:          100,
			BaseEfficiencyScore:       0.6,
			RecruitingEfficiencyScore: 0.8,
			IsFBS:                     false,
			IsUserTeam:                false,
			IsAI:                      true,
			AIBehavior:                aiBehavior,
			AIQuality:                 aiQuality,
			AIMinThreshold:            min,
			AIMaxThreshold:            max,
			AIStarMin:                 1,
			AIStarMax:                 2,
			OffensiveScheme:           off,
			DefensiveScheme:           def,
		}

		db.Create(&team)
		db.Create(&stadium)
		db.Create(&teamProfile)
	}
}

func ImportSeasonStandings() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID := strconv.Itoa(ts.CollegeSeasonID)
	collegeTeams := GetAllCollegeTeams()
	for _, team := range collegeTeams {
		teamID := strconv.Itoa(int(team.ID))
		standings := GetCollegeStandingsRecordByTeamID(teamID, seasonID)
		if standings.ID > 0 {
			continue
		}
		league := ""
		leagueID := 0
		if team.IsFBS {
			league = "FBS"
			leagueID = 1
		} else {
			league = "FCS"
			leagueID = 2
		}

		newStandings := structs.CollegeStandings{
			TeamID:         int(team.ID),
			TeamName:       team.TeamName,
			IsFBS:          team.IsFBS,
			LeagueID:       uint(leagueID),
			LeagueName:     league,
			SeasonID:       ts.CollegeSeasonID,
			Season:         ts.Season,
			ConferenceID:   team.ConferenceID,
			ConferenceName: team.Conference,
			DivisionID:     team.DivisionID,
			BaseStandings: structs.BaseStandings{
				TeamAbbr: team.TeamAbbr,
				Coach:    team.Coach,
			},
		}
		db.Create(&newStandings)
	}
}

func ImplementPrimeAge() {
	db := dbprovider.GetInstance().GetDB()

	// Recruits
	croots := GetAllRecruitRecords()
	// College Players
	collegePlayers := GetAllCollegePlayers()
	// NFL Players
	nflPlayers := GetAllNFLPlayers()

	for _, c := range croots {
		if c.PrimeAge > 0 {
			continue
		}
		primeAge := util.GetPrimeAge(c.Position, c.Archetype)
		c.AssignPrimeAge(uint(primeAge))
		repository.SaveRecruitRecord(c, db)
	}

	for _, cp := range collegePlayers {
		if cp.PrimeAge > 0 {
			continue
		}
		primeAge := util.GetPrimeAge(cp.Position, cp.Archetype)
		cp.AssignPrimeAge(uint(primeAge))
		repository.SaveCFBPlayer(cp, db)
	}

	for _, nflP := range nflPlayers {
		if nflP.PrimeAge > 0 {
			continue
		}
		primeAge := util.GetPrimeAge(nflP.Position, nflP.Archetype)
		nflP.AssignPrimeAge(uint(primeAge))
		repository.SaveNFLPlayer(nflP, db)
	}
}

func getBaseSalaryByYear(round int, pick int) float64 {
	if round == 1 {
		if pick == 1 {
			return 3.25
		}
		if pick < 6 {
			return 2.75
		}
		if pick < 11 {
			return 2.25
		}
		if pick < 17 {
			return 1.875
		}
		if pick < 25 {
			return 1.5
		}
		return 1.25
	}
	if round == 2 {
		return 1
	}
	if round == 3 {
		return 0.75
	}
	if round == 4 {
		return 0.9
	}
	if round == 5 {
		return 0.75
	}
	if round == 6 {
		return 0.9
	}
	return 0.8
}

func getBonusByYear(round int, pick int) float64 {
	if round == 1 {
		if pick == 1 {
			return 3.25
		}
		if pick < 6 {
			return 2.75
		}
		if pick < 11 {
			return 2.25
		}
		if pick < 17 {
			return 1.875
		}
		if pick < 25 {
			return 1.5
		}
		return 1.25
	}
	if round == 2 {
		return 1
	}
	if round == 3 {
		return 0.75
	}
	if round == 4 {
		return 0.3
	}
	if round == 5 {
		return 0.25
	}
	return 0
}

func getShotgunVal() int {
	shotgunNum := util.GenerateIntFromRange(1, 100)
	if shotgunNum < 61 {
		return 0
	} else if shotgunNum < 86 {
		return 1
	}
	return -1
}

func getClutchValue() int {
	clutchNum := util.GenerateIntFromRange(1, 1000)
	if clutchNum < 145 {
		return -1
	} else if clutchNum < 845 {
		return 0
	} else if clutchNum < 990 {
		return 1
	}
	return 2
}

func MigrateCFBGameplansAndDepthChartsForRemainingFCSTeams() {
	db := dbprovider.GetInstance().GetDB()

	teamPath := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\dc_positions_migration.csv"

	dcPositionsCSV := util.ReadCSV(teamPath)
	gameplansList := []structs.CollegeGameplan{}
	testgameplansList := []structs.CollegeGameplanTEST{}
	dcList := []structs.CollegeTeamDepthChart{}
	testDCList := []structs.CollegeTeamDepthChartTEST{}
	dcPList := []structs.CollegeDepthChartPosition{}
	testDCPList := []structs.CollegeDepthChartPositionTEST{}
	for i := 195; i < 265; i++ {
		gp := structs.CollegeGameplan{
			TeamID: i,
			Model: gorm.Model{
				ID: uint(i),
			},
			BaseGameplan: structs.BaseGameplan{
				OffensiveScheme: "Pistol",
				DefensiveScheme: "Multiple",
			},
		}
		gpt := structs.CollegeGameplanTEST{
			TeamID: i,
			Model: gorm.Model{
				ID: uint(i),
			},
			BaseGameplan: structs.BaseGameplan{
				OffensiveScheme: "Pistol",
				DefensiveScheme: "Multiple",
			},
		}

		gameplansList = append(gameplansList, gp)
		testgameplansList = append(testgameplansList, gpt)

		dc := structs.CollegeTeamDepthChart{
			TeamID: i,
			Model: gorm.Model{
				ID: uint(i),
			},
		}
		dct := structs.CollegeTeamDepthChartTEST{
			TeamID: i,
			Model: gorm.Model{
				ID: uint(i),
			},
		}

		dcList = append(dcList, dc)
		testDCList = append(testDCList, dct)
		for idx, row := range dcPositionsCSV {
			if idx == 0 {
				continue
			}
			positionlevel := row[3]
			originalPosition := row[2]

			dcp := structs.CollegeDepthChartPosition{
				DepthChartID:     i,
				PositionLevel:    positionlevel,
				Position:         originalPosition,
				OriginalPosition: originalPosition,
			}

			dcpt := structs.CollegeDepthChartPositionTEST{
				DepthChartID:     i,
				PositionLevel:    positionlevel,
				Position:         originalPosition,
				OriginalPosition: originalPosition,
			}

			dcPList = append(dcPList, dcp)
			testDCPList = append(testDCPList, dcpt)
		}
	}
	repository.CreateCollegeGameplansRecordsBatch(db, gameplansList, 50)
	repository.CreateCollegeGameplansTESTRecordsBatch(db, testgameplansList, 50)
	repository.CreateCollegeDCRecordsBatch(db, dcList, 50)
	repository.CreateCollegeDCTESTRecordsBatch(db, testDCList, 50)
	repository.CreateCollegeDCPRecordsBatch(db, dcPList, 200)
	repository.CreateCollegeDCPTESTRecordsBatch(db, testDCPList, 200)
}

func FixCollegeDTs() {
	db := dbprovider.GetInstance().GetDB()

	players := GetAllNFLPlayers()

	for _, p := range players {
		if p.Position != "DT" {
			continue
		}

		p.GetOverall()

		repository.SaveNFLPlayer(p, db)
	}
}

func ImportCFBRivals() {
	db := dbprovider.GetInstance().GetDB()

	path := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\2025\\rivalries.csv"

	rivalsCSV := util.ReadCSV(path)

	teamMap := make(map[string]structs.CollegeTeam)

	allCollegeTeams := GetAllCollegeTeams()

	for _, t := range allCollegeTeams {
		teamMap[t.TeamName] = t
	}

	for idx, row := range rivalsCSV {
		if idx < 314 {
			continue
		}

		id := util.ConvertStringToInt(row[0])
		rival1 := row[2]
		rival2 := row[4]
		if len(rival1) == 0 && len(rival2) == 0 {
			break
		}
		rivalryName := row[5]
		trophyName := row[6]
		priority1 := util.ConvertStringToInt(row[7])
		priority2 := util.ConvertStringToInt(row[8])

		team1, ok := teamMap[rival1]
		if !ok {
			fmt.Println("FIX!!!")
		}

		team2, ok2 := teamMap[rival2]
		if !ok2 {
			fmt.Println("FIX!!!")
		}

		rivalry := structs.CollegeRival{
			Model: gorm.Model{
				ID: uint(id),
			},
			RivalryName:     rivalryName,
			TrophyName:      trophyName,
			TeamOneID:       team1.ID,
			TeamTwoID:       team2.ID,
			HasTrophy:       len(trophyName) > 0,
			TeamOnePriority: uint(priority1),
			TeamTwoPriority: uint(priority2),
		}

		db.Create(&rivalry)
	}
}

func MigrateRetiredAndNFLPlayersToHistoricCFBTable() {
	db := dbprovider.GetInstance().GetDB()

	historicPlayersToUpload := []structs.HistoricCollegePlayer{}

	historicPlayers := GetAllHistoricCollegePlayers()
	collegePlayers := []structs.CollegePlayer{}
	for _, player := range historicPlayers {
		collegePlayerResponse := structs.CollegePlayer{
			Model:      player.Model,
			BasePlayer: player.BasePlayer,
			TeamID:     player.TeamID,
			TeamAbbr:   player.TeamAbbr,
			City:       player.City,
			State:      player.State,
			Year:       player.Year,
			IsRedshirt: player.IsRedshirt,
		}
		collegePlayers = append(collegePlayers, collegePlayerResponse)
	}
	collegePlayerMap := MakeCollegePlayerMap(collegePlayers)

	nflPlayers := GetAllNFLPlayers()
	retiredPlayers := GetAllRetiredPlayers()

	for _, p := range nflPlayers {
		historicRecord := collegePlayerMap[p.ID]

		if historicRecord.ID > 0 {
			continue
		}

		playerRecord := structs.HistoricCollegePlayer{
			Model:        p.Model,
			BasePlayer:   p.BasePlayer,
			TeamID:       int(p.CollegeID),
			TeamAbbr:     p.College,
			HasGraduated: true,
		}

		historicPlayersToUpload = append(historicPlayersToUpload, playerRecord)
	}

	for _, p := range retiredPlayers {
		historicRecord := collegePlayerMap[p.ID]

		if historicRecord.ID > 0 {
			continue
		}

		playerRecord := structs.HistoricCollegePlayer{
			Model:        p.Model,
			BasePlayer:   p.BasePlayer,
			TeamID:       int(p.CollegeID),
			TeamAbbr:     p.College,
			HasGraduated: true,
		}

		historicPlayersToUpload = append(historicPlayersToUpload, playerRecord)
	}

	repository.CreateHistoricCFBRecordsBatch(db, historicPlayersToUpload, 200)

}

func ImportCFB2021PlayerStats() {
	db := dbprovider.GetInstance().GetDB()

	path := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\2021\\2021_cfb_player_stats.csv"

	statsCSV := util.ReadCSV(path)

	collegePlayers := GetAllCollegePlayers()
	historicPlayers := GetAllHistoricCollegePlayers()
	for _, player := range historicPlayers {
		collegePlayerResponse := structs.CollegePlayer{
			Model:      player.Model,
			BasePlayer: player.BasePlayer,
			TeamID:     player.TeamID,
			TeamAbbr:   player.TeamAbbr,
			City:       player.City,
			State:      player.State,
			Year:       player.Year,
			IsRedshirt: player.IsRedshirt,
		}
		collegePlayers = append(collegePlayers, collegePlayerResponse)
	}
	collegePlayerMap := make(map[string]structs.CollegePlayer)

	for _, p := range collegePlayers {
		key := p.Position + p.Archetype + p.FirstName + p.LastName
		collegePlayerMap[key] = p
	}

	seasonStats := []structs.CollegePlayerSeasonStats{}

	for idx, row := range statsCSV {
		if idx == 0 {
			continue
		}

		pos := row[51]
		arch := row[52]
		fn := row[53]
		ln := row[54]
		key := pos + arch + fn + ln
		player := collegePlayerMap[key]
		id := player.ID
		if player.ID == 0 {
			continue
		}
		if player.ID == 10 {
			fmt.Println("MATT HOWARD??")
		}
		passAttempts := util.ConvertStringToInt(row[1])
		passCompletions := util.ConvertStringToInt(row[2])
		passYards := util.ConvertStringToInt(row[3])
		passTDs := util.ConvertStringToInt(row[4])
		interceptionsthrown := util.ConvertStringToInt(row[5])
		longestPass := util.ConvertStringToInt(row[6])
		sacksTaken := util.ConvertStringToInt(row[7])
		rushAttempts := util.ConvertStringToInt(row[8])
		rushYards := util.ConvertStringToInt(row[9])
		rushTDs := util.ConvertStringToInt(row[10])
		fumbles := util.ConvertStringToInt(row[11])
		longestRush := util.ConvertStringToInt(row[12])
		targets := util.ConvertStringToInt(row[13])
		catches := util.ConvertStringToInt(row[14])
		receivingYards := util.ConvertStringToInt(row[15])
		receivingTDs := util.ConvertStringToInt(row[16])
		longestCatch := util.ConvertStringToInt(row[17])
		soloTackles := util.ConvertStringToInt(row[18])
		assTackles := util.ConvertStringToInt(row[19])
		tacklesForLoss := util.ConvertStringToInt(row[20])
		sacksMade := util.ConvertStringToInt(row[21])
		forcedFumbles := util.ConvertStringToInt(row[22])
		fumblesRecovered := util.ConvertStringToInt(row[23])
		passDeflections := util.ConvertStringToInt(row[24])
		intsCaught := util.ConvertStringToInt(row[25])
		safeties := util.ConvertStringToInt(row[26])
		defensiveTDs := util.ConvertStringToInt(row[27])
		fgMade := util.ConvertStringToInt(row[28])
		fgAttempts := util.ConvertStringToInt(row[29])
		longestFG := util.ConvertStringToInt(row[30])
		xpMade := util.ConvertStringToInt(row[31])
		xpAttempts := util.ConvertStringToInt(row[32])
		kickoffTBs := util.ConvertStringToInt(row[33])
		punts := util.ConvertStringToInt(row[34])
		puntTBs := util.ConvertStringToInt(row[35])
		puntsInside20 := util.ConvertStringToInt(row[36])
		kickReturns := util.ConvertStringToInt(row[37])
		kickReturnYards := util.ConvertStringToInt(row[38])
		kickReturnTDs := util.ConvertStringToInt(row[39])
		puntReturns := util.ConvertStringToInt(row[40])
		puntReturnYards := util.ConvertStringToInt(row[41])
		puntReturnTDs := util.ConvertStringToInt(row[43])
		stSoloTackles := util.ConvertStringToInt(row[44])
		stassTackles := util.ConvertStringToInt(row[45])
		puntsblocked := util.ConvertStringToInt(row[46])
		fgBlocked := util.ConvertStringToInt(row[47])
		snaps := util.ConvertStringToInt(row[48])
		gamesPlayed := util.ConvertStringToInt(row[49])
		yearStr := row[50]
		year := 1
		if yearStr == "Sr." {
			year = 4
		} else if yearStr == "Jr." {
			year = 3
		} else if yearStr == "So." {
			year = 2
		}

		seasonStat := structs.CollegePlayerSeasonStats{
			CollegePlayerID: uint(id),
			SeasonID:        1,
			Year:            uint(year),
			IsRedshirt:      false,
			GamesPlayed:     gamesPlayed,
			Tackles:         float64(soloTackles) + float64(stSoloTackles) + (float64(assTackles) * 0.5) + (float64(stassTackles) * 0.5),
			RushingAvg:      (float64(rushYards) / float64(rushAttempts)),
			PassingAvg:      float64(passYards) / float64(passCompletions),
			ReceivingAvg:    float64(receivingYards) / float64(catches),
			Completion:      float64(passCompletions) / float64(passAttempts),
			BasePlayerStats: structs.BasePlayerStats{
				PassingYards:         passYards,
				PassAttempts:         passAttempts,
				PassCompletions:      passCompletions,
				PassingTDs:           passTDs,
				Interceptions:        interceptionsthrown,
				LongestPass:          longestPass,
				Sacks:                sacksTaken,
				RushAttempts:         rushAttempts,
				RushingTDs:           rushTDs,
				Fumbles:              fumbles,
				LongestRush:          longestRush,
				Targets:              targets,
				Catches:              catches,
				ReceivingYards:       receivingYards,
				ReceivingTDs:         receivingTDs,
				LongestReception:     longestCatch,
				SoloTackles:          float64(soloTackles),
				AssistedTackles:      float64(assTackles),
				TacklesForLoss:       float64(tacklesForLoss),
				SacksMade:            float64(sacksMade),
				ForcedFumbles:        forcedFumbles,
				RecoveredFumbles:     fumblesRecovered,
				PassDeflections:      passDeflections,
				InterceptionsCaught:  intsCaught,
				Safeties:             safeties,
				DefensiveTDs:         defensiveTDs,
				FGMade:               fgMade,
				FGAttempts:           fgAttempts,
				LongestFG:            longestFG,
				ExtraPointsMade:      xpMade,
				ExtraPointsAttempted: xpAttempts,
				KickoffTouchbacks:    kickoffTBs,
				Punts:                punts,
				GrossPuntDistance:    0,
				NetPuntDistance:      0,
				PuntTouchbacks:       puntTBs,
				PuntsInside20:        puntsInside20,
				KickReturns:          kickReturns,
				KickReturnTDs:        kickReturnTDs,
				KickReturnYards:      kickReturnYards,
				STSoloTackles:        float64(stSoloTackles),
				STAssistedTackles:    float64(stassTackles),
				PuntsBlocked:         puntsblocked,
				FGBlocked:            fgBlocked,
				Snaps:                snaps,
				GameType:             2,
				PuntReturns:          puntReturns,
				PuntReturnYards:      puntReturnYards,
				PuntReturnTDs:        puntReturnTDs,
				TeamID:               uint(player.TeamID),
				Team:                 player.TeamAbbr,
			},
		}

		seasonStats = append(seasonStats, seasonStat)
	}

	repository.CreateCFBPlayerSeasonStatsRecordsBatch(db, seasonStats, 200)
}

func FixATHProgressions() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	players := GetAllCollegePlayers()
	collegePlayerMap := MakeCollegePlayerMap(players)
	snapMap := GetCollegePlayerSeasonSnapMap(strconv.Itoa(int(ts.CollegeSeasonID)))
	statMap := GetCollegePlayerStatsMap(strconv.Itoa(int(ts.CollegeSeasonID)))

	path := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\2026\\2026_ATH_Fix.csv"

	playerCSV := util.ReadCSV(path)
	// Open CSV

	for idx, row := range playerCSV {
		if idx == 0 {
			continue
		}
		playerID := util.ConvertStringToInt(row[0])
		p := collegePlayerMap[uint(playerID)]
		isBroken := p.Position == "ATH"
		if !isBroken {
			continue
		}
		/*
			CSV contains letter grades only. Fuck.
			I will need to generate new values from these letter grades
			So look through how we get letter grades, and then reverse-engineer a numeric value from that
			And THEN progress that value
		*/
		footballIQ := util.ConvertStringToInt(row[6])
		speed := util.ConvertStringToInt(row[7])
		agility := util.ConvertStringToInt(row[9])
		carrying := util.ConvertStringToInt(row[8])
		catching := util.ConvertStringToInt(row[10])
		routeRunning := util.ConvertStringToInt(row[11])
		zoneCoverage := util.ConvertStringToInt(row[12])
		manCoverage := util.ConvertStringToInt(row[13])
		strength := util.ConvertStringToInt(row[14])
		tackle := util.ConvertStringToInt(row[15])
		passBlock := util.ConvertStringToInt(row[16])
		runBlock := util.ConvertStringToInt(row[17])
		passRush := util.ConvertStringToInt(row[18])
		runDefense := util.ConvertStringToInt(row[19])
		throwPower := util.ConvertStringToInt(row[20])
		throwAccuracy := util.ConvertStringToInt(row[21])
		kickAccuracy := util.ConvertStringToInt(row[22])
		kickPower := util.ConvertStringToInt(row[23])
		puntAccuracy := util.ConvertStringToInt(row[24])
		puntPower := util.ConvertStringToInt(row[25])
		isRedshirting := util.ConvertStringToBool(row[27])
		// Apply to player record
		p.ApplyFixedATHAttributes(footballIQ, speed, agility, carrying, catching, routeRunning,
			zoneCoverage, manCoverage, strength, tackle, passBlock, runBlock, passRush, runDefense,
			throwPower, throwAccuracy, kickAccuracy, kickPower, puntAccuracy, puntPower)

		p.RevertRedshirting(isRedshirting)

		// Then progress
		stats := statMap[p.ID]
		snaps := snapMap[p.ID]
		p = ProgressCollegePlayer(p, strconv.Itoa(int(ts.CollegeSeasonID)), stats, snaps)
		// Revert year back by 1
		p.RevertYearProgression()

		// Save player record
		repository.SaveCFBPlayer(p, db)
	}
}

func ImportNewDepthChartPositionRecords() {
	db := dbprovider.GetInstance().GetDB()

	positionsUpload := []structs.CollegeDepthChartPosition{}
	positionsUploadTEST := []structs.CollegeDepthChartPositionTEST{}

	nflPositionsUpload := []structs.NFLDepthChartPosition{}

	collegeTeams := GetAllCollegeTeams()
	nflTeams := GetAllNFLTeams()

	for _, team := range collegeTeams {
		newPositions := []structs.CollegeDepthChartPosition{
			structs.CollegeDepthChartPosition{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "RB",
				PositionLevel: "4",
			},
			structs.CollegeDepthChartPosition{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "LT",
				PositionLevel: "3",
			},
			structs.CollegeDepthChartPosition{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "LG",
				PositionLevel: "3",
			},
			structs.CollegeDepthChartPosition{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "C",
				PositionLevel: "3",
			},
			structs.CollegeDepthChartPosition{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "RG",
				PositionLevel: "3",
			},
			structs.CollegeDepthChartPosition{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "RT",
				PositionLevel: "3",
			},
			structs.CollegeDepthChartPosition{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "LE",
				PositionLevel: "3",
			},
			structs.CollegeDepthChartPosition{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "RE",
				PositionLevel: "3",
			},
			structs.CollegeDepthChartPosition{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "LOLB",
				PositionLevel: "3",
			},
			structs.CollegeDepthChartPosition{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "ROLB",
				PositionLevel: "3",
			},
			structs.CollegeDepthChartPosition{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "FS",
				PositionLevel: "3",
			},
			structs.CollegeDepthChartPosition{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "SS",
				PositionLevel: "3",
			},
			structs.CollegeDepthChartPosition{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "KR",
				PositionLevel: "3",
			},
			structs.CollegeDepthChartPosition{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "PR",
				PositionLevel: "3",
			},
		}

		newPositionsTEST := []structs.CollegeDepthChartPositionTEST{
			structs.CollegeDepthChartPositionTEST{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "RB",
				PositionLevel: "4",
			},
			structs.CollegeDepthChartPositionTEST{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "LT",
				PositionLevel: "3",
			},
			structs.CollegeDepthChartPositionTEST{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "LG",
				PositionLevel: "3",
			},
			structs.CollegeDepthChartPositionTEST{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "C",
				PositionLevel: "3",
			},
			structs.CollegeDepthChartPositionTEST{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "RG",
				PositionLevel: "3",
			},
			structs.CollegeDepthChartPositionTEST{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "RT",
				PositionLevel: "3",
			},
			structs.CollegeDepthChartPositionTEST{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "LE",
				PositionLevel: "3",
			},
			structs.CollegeDepthChartPositionTEST{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "RE",
				PositionLevel: "3",
			},
			structs.CollegeDepthChartPositionTEST{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "LOLB",
				PositionLevel: "3",
			},
			structs.CollegeDepthChartPositionTEST{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "ROLB",
				PositionLevel: "3",
			},
			structs.CollegeDepthChartPositionTEST{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "FS",
				PositionLevel: "3",
			},
			structs.CollegeDepthChartPositionTEST{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "SS",
				PositionLevel: "3",
			},
			structs.CollegeDepthChartPositionTEST{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "KR",
				PositionLevel: "3",
			},
			structs.CollegeDepthChartPositionTEST{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "PR",
				PositionLevel: "3",
			},
		}

		positionsUpload = append(positionsUpload, newPositions...)
		positionsUploadTEST = append(positionsUploadTEST, newPositionsTEST...)
	}

	for _, team := range nflTeams {
		newPositions := []structs.NFLDepthChartPosition{
			structs.NFLDepthChartPosition{
				DepthChartID:  team.ID,
				PlayerID:      0,
				Position:      "RB",
				PositionLevel: "4",
			},
			structs.NFLDepthChartPosition{
				DepthChartID:  team.ID,
				PlayerID:      0,
				Position:      "LT",
				PositionLevel: "3",
			},
			structs.NFLDepthChartPosition{
				DepthChartID:  team.ID,
				PlayerID:      0,
				Position:      "LG",
				PositionLevel: "3",
			},
			structs.NFLDepthChartPosition{
				DepthChartID:  team.ID,
				PlayerID:      0,
				Position:      "C",
				PositionLevel: "3",
			},
			structs.NFLDepthChartPosition{
				DepthChartID:  team.ID,
				PlayerID:      0,
				Position:      "RG",
				PositionLevel: "3",
			},
			structs.NFLDepthChartPosition{
				DepthChartID:  team.ID,
				PlayerID:      0,
				Position:      "RT",
				PositionLevel: "3",
			},
			structs.NFLDepthChartPosition{
				DepthChartID:  team.ID,
				PlayerID:      0,
				Position:      "LE",
				PositionLevel: "3",
			},
			structs.NFLDepthChartPosition{
				DepthChartID:  team.ID,
				PlayerID:      0,
				Position:      "RE",
				PositionLevel: "3",
			},
			structs.NFLDepthChartPosition{
				DepthChartID:  team.ID,
				PlayerID:      0,
				Position:      "LOLB",
				PositionLevel: "3",
			},
			structs.NFLDepthChartPosition{
				DepthChartID:  team.ID,
				PlayerID:      0,
				Position:      "ROLB",
				PositionLevel: "3",
			},
			structs.NFLDepthChartPosition{
				DepthChartID:  team.ID,
				PlayerID:      0,
				Position:      "FS",
				PositionLevel: "3",
			},
			structs.NFLDepthChartPosition{
				DepthChartID:  team.ID,
				PlayerID:      0,
				Position:      "SS",
				PositionLevel: "3",
			},
			structs.NFLDepthChartPosition{
				DepthChartID:  team.ID,
				PlayerID:      0,
				Position:      "KR",
				PositionLevel: "3",
			},
			structs.NFLDepthChartPosition{
				DepthChartID:  team.ID,
				PlayerID:      0,
				Position:      "PR",
				PositionLevel: "3",
			},
		}

		nflPositionsUpload = append(nflPositionsUpload, newPositions...)
	}

	repository.CreateCFBDepthChartPositionRecordsBatch(db, positionsUpload, 200)
	repository.CreateCFBDepthChartPositionTESTRecordsBatch(db, positionsUploadTEST, 200)
	repository.CreateNFLDepthChartPositionRecordsBatch(db, nflPositionsUpload, 200)
}

func FixSecondaryPositions() {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()
	seasonID := ts.CollegeSeasonID - 1
	collegePlayers := GetAllCollegePlayers()
	nflDraftees := GetAllNFLDraftees()
	nflPlayers := GetAllNFLPlayers()

	snapMap := GetCollegePlayerSeasonSnapMap(strconv.Itoa(seasonID))
	nflSnapMap := GetNFLPlayerSeasonSnapMap(strconv.Itoa(seasonID))

	for _, cp := range collegePlayers {
		if len(cp.PositionTwo) > 0 {
			continue
		}
		snaps := snapMap[cp.ID]
		hasSecondaryPosition := false
		totalSnaps := snaps.GetTotalSnaps()
		if snaps.ID == 0 || totalSnaps == 0 {
			continue
		}
		totalSnaps -= int(snaps.STSnaps)
		mostPlayedPosition, mostPlayedSnaps := getMostPlayedPosition(snaps.BasePlayerSeasonSnaps, cp.Position)
		posThreshold := float64(totalSnaps) * 0.8

		if mostPlayedSnaps == 0 {
			continue
		}

		if mostPlayedPosition != cp.Position && float64(mostPlayedSnaps) > posThreshold {
			// Designate New Position
			newArchetype, archCheck := getNewArchetype(cp.Position, cp.Archetype, mostPlayedPosition)
			// If Archhetype exists, assign new position. Otherwise, progress by old position
			hasSecondaryPosition = true
			if archCheck {
				cp.DesignateNewPosition(mostPlayedPosition, newArchetype)
			}
		}
		if !hasSecondaryPosition {
			continue
		}
		repository.SaveCFBPlayer(cp, db)
	}

	for _, p := range nflDraftees {
		if len(p.PositionTwo) > 0 {
			continue
		}
		snaps := snapMap[p.ID]
		hasSecondaryPosition := false
		totalSnaps := snaps.GetTotalSnaps()
		if snaps.ID == 0 || totalSnaps == 0 {
			continue
		}
		totalSnaps -= int(snaps.STSnaps)
		mostPlayedPosition, mostPlayedSnaps := getMostPlayedPosition(snaps.BasePlayerSeasonSnaps, p.Position)
		posThreshold := float64(totalSnaps) * 0.8

		if mostPlayedSnaps == 0 {
			continue
		}

		if mostPlayedPosition != p.Position && float64(mostPlayedSnaps) > posThreshold {
			// Designate New Position
			newArchetype, archCheck := getNewArchetype(p.Position, p.Archetype, mostPlayedPosition)
			// If Archhetype exists, assign new position. Otherwise, progress by old position
			hasSecondaryPosition = true
			if archCheck {
				p.DesignateNewPosition(mostPlayedPosition, newArchetype)
			}
		}
		if !hasSecondaryPosition {
			continue
		}
		repository.SaveNFLDrafteeRecord(p, db)
	}

	for _, p := range nflPlayers {
		if len(p.PositionTwo) > 0 {
			continue
		}
		snaps := nflSnapMap[p.ID]
		hasSecondaryPosition := false
		totalSnaps := snaps.GetTotalSnaps()
		if snaps.ID == 0 || totalSnaps == 0 {
			continue
		}
		totalSnaps -= int(snaps.STSnaps)
		mostPlayedPosition, mostPlayedSnaps := getMostPlayedPosition(snaps.BasePlayerSeasonSnaps, p.Position)
		posThreshold := float64(totalSnaps) * 0.8

		if mostPlayedSnaps == 0 {
			continue
		}

		if mostPlayedPosition != p.Position && float64(mostPlayedSnaps) > posThreshold {
			// Designate New Position
			newArchetype, archCheck := getNewArchetype(p.Position, p.Archetype, mostPlayedPosition)
			// If Archhetype exists, assign new position. Otherwise, progress by old position
			hasSecondaryPosition = true
			if archCheck {
				p.DesignateNewPosition(mostPlayedPosition, newArchetype)
			}
		}
		if !hasSecondaryPosition {
			continue
		}
		repository.SaveNFLPlayer(p, db)
	}
}
