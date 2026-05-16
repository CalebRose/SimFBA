package managers

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
	"github.com/jinzhu/gorm"
)

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
	var draftPicksToUpload []structs.NFLDraftPick

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
	playerPath := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\2026\\2026_Free_Agency_Expected_Values.csv"

	nflCSV := util.ReadCSV(playerPath)

	nflPlayerMap := GetAllNFLPlayersMap()

	for idx, row := range nflCSV {
		if idx < 1 {
			continue
		}

		playerID := row[0]
		id := util.ConvertStringToInt(playerID)
		valueStr := strings.TrimSpace(row[3])
		aavStr := strings.TrimSpace(row[4])
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

func ImportCFBGames(isSpringGames bool) {
	db := dbprovider.GetInstance().GetDB()

	path := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\2026\\2026_fbs_playoffs.csv"

	gamesCSV := util.ReadCSV(path)
	ts := GetTimestamp()
	teamMap := make(map[string]structs.CollegeTeam)
	collegeGames := GetCollegeGamesBySeasonID(strconv.Itoa(int(ts.CollegeSeasonID)))
	collegeGameMap := MakeCollegeGameMapByID(collegeGames)

	var games []structs.CollegeGame

	allCollegeTeams := GetAllCollegeTeams()

	for _, t := range allCollegeTeams {
		teamMap[t.TeamAbbr] = t
	}

	for idx, row := range gamesCSV {
		if idx == 0 {
			continue
		}

		gameID := util.ConvertStringToInt(row[0])
		existingGame := collegeGameMap[uint(gameID)]
		if existingGame.ID > 0 {
			continue
		}
		season := ts.Season
		seasonID := season - 2020
		week := util.ConvertStringToInt(row[2])
		weekID := util.GetWeekID(uint(seasonID), uint(week))
		homeTeamAbbr := row[3]
		awayTeamAbbr := row[4]
		ht := teamMap[homeTeamAbbr]
		at := teamMap[awayTeamAbbr]
		homeTeamID := ht.ID
		awayTeamID := at.ID
		homeTeamCoach := ht.Coach
		awayTeamCoach := at.Coach
		timeSlot := row[16]
		if week < 17 {
			timeSlot = SelectTimeslotForGameByConferenceID(uint(ht.ConferenceID))
		}
		// Need to implement Stadium ID
		stadium := row[17]
		city := row[18]
		state := row[19]
		isDomed := util.ConvertStringToBool(row[20])
		// Need to check for if a game is in a domed stadium or not
		isConferenceGame := !isSpringGames && ht.ConferenceID == at.ConferenceID
		isDivisionGame := !isSpringGames && isConferenceGame && ht.DivisionID == at.DivisionID && ht.DivisionID > 0
		conferenceID := 0
		if isConferenceGame {
			conferenceID = ht.ConferenceID
		}
		isNeutralSite := util.ConvertStringToBool(row[7])
		if !isNeutralSite {
			stadium = ht.Stadium
			city = ht.City
			state = ht.State
		}
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
			WeekID:                   int(weekID),
			Week:                     week,
			HomeTeamID:               int(homeTeamID),
			AwayTeamID:               int(awayTeamID),
			HomeTeam:                 homeTeamAbbr,
			AwayTeam:                 awayTeamAbbr,
			HomeTeamCoach:            homeTeamCoach,
			AwayTeamCoach:            awayTeamCoach,
			IsConferenceChampionship: isConferenceChampionship,
			IsSpringGame:             isSpringGames,
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

		games = append(games, game)
	}
	repository.CreateCFBGameRecordsBatch(db, games, 250)

	GenerateWeatherForGames()
}

func ImportNFLGames() {
	db := dbprovider.GetInstance().GetDB()

	path := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\2026\\2026_nfl_postseason_games.csv"

	gamesCSV := util.ReadCSV(path)

	teamMap := make(map[string]structs.NFLTeam)

	allNFLTeams := GetAllNFLTeams()

	var games []structs.NFLGame
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
		weekID := util.GetWeekID(uint(seasonID), uint(week))
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
		timeSlot := row[12]
		// Need to implement Stadium ID
		stadium := ht.Stadium
		city := ht.City
		state := ht.State
		// Need to check for if a game is in a domed stadium or not
		isConferenceGame := ht.ConferenceID == at.ConferenceID
		isDivisionGame := ht.DivisionID == at.DivisionID && ht.DivisionID > 0
		isNeutralSite := util.ConvertStringToBool(row[5])
		isPreseasonGame := util.ConvertStringToBool(row[6])
		// isConferenceChampionship := util.ConvertStringToBool(row[7])
		isPlayoffGame := util.ConvertStringToBool(row[8])
		if week > 18 {
			isConferenceGame = false
			isDivisionGame = false
			isPlayoffGame = true
			isPreseasonGame = false
		}
		isNationalChampionship := util.ConvertStringToBool(row[9])
		gameTitle := row[17]
		nextGame := util.ConvertStringToInt(row[18])
		nextGameHOA := row[19]

		game := structs.NFLGame{
			Model:           gorm.Model{ID: uint(gameID)},
			SeasonID:        seasonID,
			WeekID:          int(weekID),
			Week:            week,
			HomeTeamID:      int(homeTeamID),
			AwayTeamID:      int(awayTeamID),
			HomeTeam:        homeTeamName,
			AwayTeam:        awayTeamName,
			HomeTeamCoach:   homeTeamCoach,
			AwayTeamCoach:   awayTeamCoach,
			IsPreseasonGame: isPreseasonGame,
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

		games = append(games, game)
	}

	repository.CreateNFLGameRecordsBatch(db, games, 250)
	GenerateWeatherForGames()
}

func ImportCFBTeams() {
	db := dbprovider.GetInstance().GetDB()

	teamPath := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\2027\\teams.csv"
	stadiumPath := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\2027\\stadia.csv"
	// profilePath := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\2027\\profiles.csv"

	teamCSV := util.ReadCSV(teamPath)
	stadiumCSV := util.ReadCSV(stadiumPath)
	// profileCSV := util.ReadCSV(profilePath)

	for idx, row := range teamCSV {
		if idx == 0 {
			continue
		}

		stadiumRecord := stadiumCSV[idx]
		// profileRecord := profileCSV[idx]

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

		// aiBehavior := profileRecord[10]
		// aiQuality := profileRecord[11]
		// min := util.ConvertStringToInt(profileRecord[12])
		// max := util.ConvertStringToInt(profileRecord[13])
		// off := profileRecord[17]
		// def := profileRecord[18]

		// teamProfile := structs.RecruitingTeamProfile{
		// 	Model: gorm.Model{
		// 		ID: uint(teamID),
		// 	},
		// 	TeamID:                    teamID,
		// 	Team:                      teamName,
		// 	TeamAbbreviation:          abbr,
		// 	State:                     state,
		// 	ScholarshipsAvailable:     40,
		// 	WeeklyPoints:              100,
		// 	SpentPoints:               0,
		// 	TotalCommitments:          0,
		// 	RecruitClassSize:          20,
		// 	PortalReputation:          100,
		// 	BaseEfficiencyScore:       0.6,
		// 	RecruitingEfficiencyScore: 0.8,
		// 	IsFBS:                     false,
		// 	IsUserTeam:                false,
		// 	IsAI:                      true,
		// 	AIBehavior:                aiBehavior,
		// 	AIQuality:                 aiQuality,
		// 	AIMinThreshold:            min,
		// 	AIMaxThreshold:            max,
		// 	AIStarMin:                 1,
		// 	AIStarMax:                 2,
		// 	OffensiveScheme:           off,
		// 	DefensiveScheme:           def,
		// }

		db.Create(&team)
		db.Create(&stadium)
		// db.Create(&teamProfile)
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

	path := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\2027\\rivalries.csv"

	rivalsCSV := util.ReadCSV(path)

	teamMap := make(map[string]structs.CollegeTeam)

	allCollegeTeams := GetAllCollegeTeams()

	for _, t := range allCollegeTeams {
		teamMap[t.TeamName] = t
	}

	for idx, row := range rivalsCSV {
		if idx == 0 {
			continue
		}

		id := util.ConvertStringToInt(row[0])
		rival1 := row[4]
		rival2 := row[6]
		if len(rival1) == 0 && len(rival2) == 0 {
			break
		}
		rivalryName := row[1]
		trophyName := row[2]
		priority1 := util.ConvertStringToInt(row[8])
		priority2 := util.ConvertStringToInt(row[9])

		team1, ok := teamMap[rival1]
		if !ok {
			fmt.Println("FIX!!!")
		}

		team2, ok2 := teamMap[rival2]
		if !ok2 {
			fmt.Println("FIX!!!")
		}

		isAnnualRivalry := util.ConvertStringToBool(row[10])
		conferenceID := 0
		if team1.ConferenceID > 0 && team1.ConferenceID == team2.ConferenceID {
			conferenceID = int(team1.ConferenceID)
		}

		preferredWeek := util.ConvertStringToInt(row[12])
		preferredTimeSlot := row[13]
		isNeutral := util.ConvertStringToBool(row[14])
		stadiumID := util.ConvertStringToInt(row[15])

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
			IsAnnualRivalry: isAnnualRivalry,
			ConferenceID:    uint(conferenceID),
			PreferredWeek:   uint8(preferredWeek),
			Timeslot:        preferredTimeSlot,
			IsNeutralSite:   isNeutral,
			StadiumID:       uint(stadiumID),
		}

		db.Create(&rivalry)
	}
}

func MigrateRetiredAndNFLPlayersToHistoricCFBTable() {
	db := dbprovider.GetInstance().GetDB()

	var historicPlayersToUpload []structs.HistoricCollegePlayer

	historicPlayers := GetAllHistoricCollegePlayers()
	var collegePlayers []structs.CollegePlayer
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

	var positionsUpload []structs.CollegeDepthChartPosition
	var positionsUploadTEST []structs.CollegeDepthChartPositionTEST

	var nflPositionsUpload []structs.NFLDepthChartPosition

	collegeTeams := GetAllCollegeTeams()
	nflTeams := GetAllNFLTeams()

	for _, team := range collegeTeams {
		newPositions := []structs.CollegeDepthChartPosition{
			{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "RB",
				PositionLevel: "4",
			},
			{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "LT",
				PositionLevel: "3",
			},
			{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "LG",
				PositionLevel: "3",
			},
			{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "C",
				PositionLevel: "3",
			},
			{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "RG",
				PositionLevel: "3",
			},
			{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "RT",
				PositionLevel: "3",
			},
			{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "LE",
				PositionLevel: "3",
			},
			{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "RE",
				PositionLevel: "3",
			},
			{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "LOLB",
				PositionLevel: "3",
			},
			{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "ROLB",
				PositionLevel: "3",
			},
			{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "FS",
				PositionLevel: "3",
			},
			{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "SS",
				PositionLevel: "3",
			},
			{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "KR",
				PositionLevel: "3",
			},
			{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "PR",
				PositionLevel: "3",
			},
		}

		newPositionsTEST := []structs.CollegeDepthChartPositionTEST{
			{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "RB",
				PositionLevel: "4",
			},
			{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "LT",
				PositionLevel: "3",
			},
			{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "LG",
				PositionLevel: "3",
			},
			{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "C",
				PositionLevel: "3",
			},
			{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "RG",
				PositionLevel: "3",
			},
			{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "RT",
				PositionLevel: "3",
			},
			{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "LE",
				PositionLevel: "3",
			},
			{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "RE",
				PositionLevel: "3",
			},
			{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "LOLB",
				PositionLevel: "3",
			},
			{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "ROLB",
				PositionLevel: "3",
			},
			{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "FS",
				PositionLevel: "3",
			},
			{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "SS",
				PositionLevel: "3",
			},
			{
				DepthChartID:  int(team.ID),
				PlayerID:      0,
				Position:      "KR",
				PositionLevel: "3",
			},
			{
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
			{
				DepthChartID:  team.ID,
				PlayerID:      0,
				Position:      "RB",
				PositionLevel: "4",
			},
			{
				DepthChartID:  team.ID,
				PlayerID:      0,
				Position:      "LT",
				PositionLevel: "3",
			},
			{
				DepthChartID:  team.ID,
				PlayerID:      0,
				Position:      "LG",
				PositionLevel: "3",
			},
			{
				DepthChartID:  team.ID,
				PlayerID:      0,
				Position:      "C",
				PositionLevel: "3",
			},
			{
				DepthChartID:  team.ID,
				PlayerID:      0,
				Position:      "RG",
				PositionLevel: "3",
			},
			{
				DepthChartID:  team.ID,
				PlayerID:      0,
				Position:      "RT",
				PositionLevel: "3",
			},
			{
				DepthChartID:  team.ID,
				PlayerID:      0,
				Position:      "LE",
				PositionLevel: "3",
			},
			{
				DepthChartID:  team.ID,
				PlayerID:      0,
				Position:      "RE",
				PositionLevel: "3",
			},
			{
				DepthChartID:  team.ID,
				PlayerID:      0,
				Position:      "LOLB",
				PositionLevel: "3",
			},
			{
				DepthChartID:  team.ID,
				PlayerID:      0,
				Position:      "ROLB",
				PositionLevel: "3",
			},
			{
				DepthChartID:  team.ID,
				PlayerID:      0,
				Position:      "FS",
				PositionLevel: "3",
			},
			{
				DepthChartID:  team.ID,
				PlayerID:      0,
				Position:      "SS",
				PositionLevel: "3",
			},
			{
				DepthChartID:  team.ID,
				PlayerID:      0,
				Position:      "KR",
				PositionLevel: "3",
			},
			{
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

func FixRecruitingProfiles() {
	db := dbprovider.GetInstance().GetDB()

	recruitProfiles := repository.FindRecruitPlayerProfileRecords("", "", false, false, false)
	recruitPointAllocations := repository.FindRecruitPointAllocationRecords(repository.RecruitingClauses{
		WeekID: "2600",
	})

	recruitPointAllocationsMap := make(map[uint][]structs.RecruitPointAllocation)
	for _, rpa := range recruitPointAllocations {
		if _, ok := recruitPointAllocationsMap[uint(rpa.RecruitProfileID)]; !ok {
			recruitPointAllocationsMap[uint(rpa.RecruitProfileID)] = []structs.RecruitPointAllocation{}
		}
		recruitPointAllocationsMap[uint(rpa.RecruitProfileID)] = append(recruitPointAllocationsMap[uint(rpa.RecruitProfileID)], rpa)
	}

	for _, rp := range recruitProfiles {
		if !rp.CaughtCheating {
			continue
		}

		pointAllocations, ok := recruitPointAllocationsMap[rp.ID]
		if !ok {
			continue
		}

		totalPoints := 0.0
		previousWeekPoints := 0.0
		for _, pa := range pointAllocations {
			totalPoints += float64(pa.RESAffectedPoints)
			previousWeekPoints = float64(pa.RESAffectedPoints)
		}
		streak := len(pointAllocations)

		rp.FixPoints(totalPoints, previousWeekPoints, streak)

		repository.SaveRecruitProfile(rp, db)
	}
}
