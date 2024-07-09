package managers

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
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

// Imports 2023-25 Draft Picks
func ImportNFLDraftPicks() {
	db := dbprovider.GetInstance().GetDB()
	draftpickPath := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\2026_draft_picks_official.csv"

	nflCSV := util.ReadCSV(draftpickPath)

	nflTeams := GetAllNFLTeams()
	teamMap := make(map[uint]structs.NFLTeam)

	for _, team := range nflTeams {
		teamMap[team.ID] = team
	}

	for idx, row := range nflCSV {
		if idx < 1 {
			continue
		}

		season := util.ConvertStringToInt(row[2])
		season_id := util.ConvertStringToInt(row[1])
		round := util.ConvertStringToInt(row[3])
		pick := util.ConvertStringToInt(row[4])
		teamID := util.ConvertStringToInt(row[6])
		team := teamMap[uint(teamID)]
		tradeValue := util.ConvertStringToFloat(row[12])

		draftPick := structs.NFLDraftPick{
			SeasonID:       uint(season_id),
			Season:         uint(season),
			DraftRound:     uint(round),
			DraftNumber:    uint(pick),
			Team:           team.TeamName + " " + team.Mascot,
			TeamID:         team.ID,
			OriginalTeamID: team.ID,
			OriginalTeam:   team.TeamName + " " + team.Mascot,
			Notes:          "",
			DraftValue:     tradeValue,
		}

		db.Create(&draftPick)
	}
}

func ImportMinimumFAValues() {
	db := dbprovider.GetInstance().GetDB()
	playerPath := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\2024_Free_Agency_Expected_Values.csv"

	nflCSV := util.ReadCSV(playerPath)

	for idx, row := range nflCSV {
		if idx < 1 {
			continue
		}

		playerID := row[0]
		value := util.ConvertStringToFloat(row[5])

		NFLPlayerRecord := GetNFLPlayerRecord(playerID)
		if NFLPlayerRecord.ID == 0 {
			log.Fatalln("Something is wrong, this player was not uploaded.")
		}

		NFLPlayerRecord.AssignMinimumValue(value)

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

	path := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\2024_draft_picks_upload.csv"

	draftPickCSV := util.ReadCSV(path)

	draftPicks := GetAllCurrentSeasonDraftPicks()
	pickMap := make(map[uint]structs.NFLDraftPick)
	var latestID uint = 673 // Latest ID from Draft Pick table in DB

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
			SeasonID: 4,
			Season:   2024,
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

	path := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\2024\\2024_games_official.csv"

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
		weekID := week + 65 // Week 43 is week 0 of the 2023 Season
		homeTeamAbbr := row[3]
		awayTeamAbbr := row[4]
		ht := teamMap[homeTeamAbbr]
		at := teamMap[awayTeamAbbr]
		homeTeamID := ht.ID
		awayTeamID := at.ID
		homeTeamCoach := ht.Coach
		awayTeamCoach := at.Coach
		timeSlot := row[20]
		// Need to implement Stadium ID
		stadium := row[21]
		city := row[22]
		state := row[23]
		isDomed := util.ConvertStringToBool(row[24])
		// Need to check for if a game is in a domed stadium or not
		isConferenceGame := util.ConvertStringToBool(row[9])
		isDivisionGame := util.ConvertStringToBool(row[10])
		isNeutralSite := util.ConvertStringToBool(row[11])
		isConferenceChampionship := util.ConvertStringToBool(row[12])
		isBowlGame := util.ConvertStringToBool(row[13])
		isPlayoffGame := util.ConvertStringToBool(row[14])
		isNationalChampionship := util.ConvertStringToBool(row[15])
		isHomePrevBye := util.ConvertStringToBool(row[16])
		isAwayPrevBye := util.ConvertStringToBool(row[17])
		gameTitle := row[25]
		nextGame := util.ConvertStringToInt(row[27])
		nextGameHOA := row[28]

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
		}

		db.Create(&game)
	}
}

func ImportNFLGames() {
	db := dbprovider.GetInstance().GetDB()

	path := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\2024\\2024_NFL_Games.csv"

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
			weekID = week + 26
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
		stadium := row[19]
		city := row[20]
		state := row[21]
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
	croots := GetAllUnsignedRecruits()
	// College Players
	collegePlayers := GetAllCollegePlayers()
	// Unsigned Players
	unsignedPlayers := GetAllUnsignedPlayers()
	// NFL Players
	nflPlayers := GetAllNFLPlayers()

	for _, c := range croots {
		primeAge := util.GetPrimeAge(c.Position, c.Archetype)
		c.AssignPrimeAge(uint(primeAge))
		db.Save(&c)
	}

	for _, cp := range collegePlayers {
		primeAge := util.GetPrimeAge(cp.Position, cp.Archetype)
		cp.AssignPrimeAge(uint(primeAge))
		db.Save(&cp)
	}

	for _, up := range unsignedPlayers {
		primeAge := util.GetPrimeAge(up.Position, up.Archetype)
		up.AssignPrimeAge(uint(primeAge))
		db.Save(&up)
	}

	for _, nflP := range nflPlayers {
		primeAge := util.GetPrimeAge(nflP.Position, nflP.Archetype)
		nflP.AssignPrimeAge(uint(primeAge))
		db.Save(&nflP)
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
