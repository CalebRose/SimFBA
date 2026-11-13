package managers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
)

type confChampVenueDefault struct {
	GameTitle string
	IsNeutral bool
	StadiumID uint
	Stadium   string
	City      string
	State     string
	IsDomed   bool
	TimeSlot  string
}

var defaultConfChampVenues = map[uint]confChampVenueDefault{
	3:  {GameTitle: "ACC Championship Game", IsNeutral: true, StadiumID: 158, Stadium: "Bank of America Stadium", City: "Charlotte", State: "NC", IsDomed: false, TimeSlot: "Saturday Evening"},
	4:  {GameTitle: "Big Ten Championship Game", IsNeutral: true, StadiumID: 154, Stadium: "Lucas Oil Stadium", City: "Indianapolis", State: "IN", IsDomed: true, TimeSlot: "Saturday Evening"},
	5:  {GameTitle: "Big 12 Championship Game", IsNeutral: true, StadiumID: 99, Stadium: "AT&T Stadium", City: "Arlington", State: "TX", IsDomed: true, TimeSlot: "Saturday Morning"},
	6:  {GameTitle: "Pac-12 Championship Game", IsNeutral: true, StadiumID: 155, Stadium: "Allegiant Stadium", City: "Las Vegas", State: "NV", IsDomed: true, TimeSlot: "Friday Night"},
	7:  {GameTitle: "SEC Championship Game", IsNeutral: true, StadiumID: 25, Stadium: "Mercedes-Benz Stadium", City: "Atlanta", State: "GA", IsDomed: true, TimeSlot: "Saturday Afternoon"},
	8:  {GameTitle: "AAC Championship Game", IsNeutral: false, TimeSlot: "Saturday Afternoon"},
	9:  {GameTitle: "C-USA Championship Game", IsNeutral: false, TimeSlot: "Friday Night"},
	10: {GameTitle: "MAC Championship Game", IsNeutral: true, StadiumID: 147, Stadium: "Ford Field", City: "Detroit", State: "MI", IsDomed: true, TimeSlot: "Saturday Morning"},
	11: {GameTitle: "MWC Championship Game", IsNeutral: false, TimeSlot: "Saturday Afternoon"},
	12: {GameTitle: "Sun Belt Championship Game", IsNeutral: false, TimeSlot: "Saturday Afternoon"},
}

func GetCFBPostseasonGamesForCurrentSeason() []structs.CollegeGame {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	var games []structs.CollegeGame
	db.Order("week_id asc, id asc").
		Where("season_id = ? AND (is_conference_championship = ? OR is_bowl_game = ? OR is_playoff_game = ? OR is_national_championship = ?)",
			ts.CollegeSeasonID, true, true, true, true).
		Find(&games)

	repository.AssignCFBStreamChannels(games)
	return games
}

func DeleteCFBPostseasonGame(gameID string) error {
	db := dbprovider.GetInstance().GetDB()
	game := GetCollegeGameByGameID(gameID)
	if game.ID == 0 {
		return errors.New("game not found")
	}
	if game.GameComplete {
		return errors.New("cannot delete a completed game")
	}
	return db.Delete(&game).Error
}

func SaveCFBConferenceChampionship(dto structs.CFBConferenceChampionshipDTO) (structs.CollegeGame, error) {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	if dto.HomeTeamID == 0 || dto.AwayTeamID == 0 {
		return structs.CollegeGame{}, errors.New("both Team 1 and Team 2 must be selected")
	}
	if dto.HomeTeamID == dto.AwayTeamID {
		return structs.CollegeGame{}, errors.New("Team 1 and Team 2 cannot be the same team")
	}
	if dto.ConferenceID == 0 || dto.ConferenceID == 13 || dto.ConferenceID == 22 {
		return structs.CollegeGame{}, errors.New("invalid FBS conference selected")
	}

	homeTeam := GetTeamByTeamID(strconv.Itoa(int(dto.HomeTeamID)))
	awayTeam := GetTeamByTeamID(strconv.Itoa(int(dto.AwayTeamID)))
	if homeTeam.ID == 0 || awayTeam.ID == 0 {
		return structs.CollegeGame{}, errors.New("could not find selected teams")
	}
	if !homeTeam.IsFBS || !awayTeam.IsFBS {
		return structs.CollegeGame{}, errors.New("only FBS conference teams are eligible for FBS Conference Championships")
	}
	if uint(homeTeam.ConferenceID) != dto.ConferenceID || uint(awayTeam.ConferenceID) != dto.ConferenceID {
		return structs.CollegeGame{}, errors.New("both selected teams must belong to the selected conference")
	}

	week := dto.Week
	if week <= 0 {
		week = 15
	}
	seasonID := ts.CollegeSeasonID
	weekID := int(util.GetWeekID(uint(seasonID), uint(week)))

	standings := GetAllCollegeStandingsBySeasonID(strconv.Itoa(seasonID))
	standingsMap := MakeCollegeStandingsMapByTeamID(standings)
	homeRank := uint(standingsMap[homeTeam.ID].Rank)
	awayRank := uint(standingsMap[awayTeam.ID].Rank)

	venueDefault, hasDefault := defaultConfChampVenues[dto.ConferenceID]
	gameTitle := strings.TrimSpace(dto.GameTitle)
	if gameTitle == "" {
		if hasDefault && venueDefault.GameTitle != "" {
			gameTitle = venueDefault.GameTitle
		} else {
			confName := homeTeam.Conference
			if dto.Conference != "" {
				confName = dto.Conference
			}
			gameTitle = fmt.Sprintf("%s Championship Game", confName)
		}
	}

	timeSlot := strings.TrimSpace(dto.TimeSlot)
	if timeSlot == "" {
		if hasDefault && venueDefault.TimeSlot != "" {
			timeSlot = venueDefault.TimeSlot
		} else {
			timeSlot = "Saturday Afternoon"
		}
	}

	isNeutral := dto.IsNeutral
	stadiumID := dto.StadiumID
	stadium := strings.TrimSpace(dto.Stadium)
	city := strings.TrimSpace(dto.City)
	state := strings.TrimSpace(dto.State)
	isDomed := dto.IsDomed

	if !isNeutral {
		stadiumID = homeTeam.StadiumID
		stadium = homeTeam.Stadium
		city = homeTeam.City
		state = homeTeam.State
		stadiumRec := GetStadiumByStadiumID(strconv.Itoa(int(stadiumID)))
		isDomed = stadiumRec.IsDomed
	} else if stadium == "" {
		if hasDefault && venueDefault.IsNeutral {
			stadiumID = venueDefault.StadiumID
			stadium = venueDefault.Stadium
			city = venueDefault.City
			state = venueDefault.State
			isDomed = venueDefault.IsDomed
		} else if stadiumID > 0 {
			stadiumRec := GetStadiumByStadiumID(strconv.Itoa(int(stadiumID)))
			stadium = stadiumRec.StadiumName
			city = stadiumRec.City
			state = stadiumRec.State
			isDomed = stadiumRec.IsDomed
		} else {
			stadiumID = homeTeam.StadiumID
			stadium = homeTeam.Stadium
			city = homeTeam.City
			state = homeTeam.State
			stadiumRec := GetStadiumByStadiumID(strconv.Itoa(int(stadiumID)))
			isDomed = stadiumRec.IsDomed
		}
	}

	var existing structs.CollegeGame
	if dto.GameID > 0 {
		db.Where("id = ?", dto.GameID).First(&existing)
	}
	if existing.ID == 0 {
		db.Where("season_id = ? AND is_conference_championship = ? AND conference_id = ?", seasonID, true, dto.ConferenceID).First(&existing)
	}
	if existing.ID == 0 {
		db.Where("season_id = ? AND is_conference_championship = ? AND game_title = ?", seasonID, true, gameTitle).First(&existing)
	}

	if existing.ID > 0 {
		if existing.GameComplete {
			return structs.CollegeGame{}, errors.New("this conference championship game has already been completed")
		}
		existing.Week = week
		existing.WeekID = weekID
		existing.HomeTeamID = int(homeTeam.ID)
		existing.HomeTeam = homeTeam.TeamAbbr
		existing.HomeTeamCoach = homeTeam.Coach
		existing.HomeTeamRank = homeRank
		existing.AwayTeamID = int(awayTeam.ID)
		existing.AwayTeam = awayTeam.TeamAbbr
		existing.AwayTeamCoach = awayTeam.Coach
		existing.AwayTeamRank = awayRank
		existing.TimeSlot = timeSlot
		existing.StadiumID = stadiumID
		existing.Stadium = stadium
		existing.City = city
		existing.State = state
		existing.IsDomed = isDomed
		existing.IsNeutral = isNeutral
		existing.IsNightGame = strings.Contains(timeSlot, "Evening") || strings.Contains(timeSlot, "Night")
		existing.IsConference = true
		existing.IsConferenceChampionship = true
		existing.ConferenceID = dto.ConferenceID
		existing.GameTitle = gameTitle
		repository.SaveCFBGameRecord(existing, db)
		GenerateWeatherForGames()
		games := []structs.CollegeGame{existing}
		repository.AssignCFBStreamChannels(games)
		return games[0], nil
	}

	newGame := structs.CollegeGame{
		SeasonID:                 seasonID,
		WeekID:                   weekID,
		Week:                     week,
		HomeTeamID:               int(homeTeam.ID),
		HomeTeam:                 homeTeam.TeamAbbr,
		HomeTeamCoach:            homeTeam.Coach,
		HomeTeamRank:             homeRank,
		AwayTeamID:               int(awayTeam.ID),
		AwayTeam:                 awayTeam.TeamAbbr,
		AwayTeamCoach:            awayTeam.Coach,
		AwayTeamRank:             awayRank,
		TimeSlot:                 timeSlot,
		StadiumID:                stadiumID,
		Stadium:                  stadium,
		City:                     city,
		State:                    state,
		IsDomed:                  isDomed,
		IsNeutral:                isNeutral,
		IsNightGame:              strings.Contains(timeSlot, "Evening") || strings.Contains(timeSlot, "Night"),
		IsConference:             true,
		IsDivisional:             false,
		IsConferenceChampionship: true,
		IsBowlGame:               false,
		IsPlayoffGame:            false,
		IsNationalChampionship:   false,
		ConferenceID:             dto.ConferenceID,
		GameTitle:                gameTitle,
	}

	if err := db.Create(&newGame).Error; err != nil {
		return structs.CollegeGame{}, err
	}

	GenerateWeatherForGames()
	games := []structs.CollegeGame{newGame}
	repository.AssignCFBStreamChannels(games)
	return games[0], nil
}

func UploadCFBBowlGames(req structs.CFBBowlUploadRequest) ([]structs.CollegeGame, error) {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	if ts.CollegeWeek <= 14 && !req.OverrideWeekCheck {
		return nil, fmt.Errorf("bowl games can only be added to the schedule after Week 14 (current week is Week %d)", ts.CollegeWeek)
	}

	if len(req.Fixtures) == 0 {
		return nil, errors.New("no bowl game fixtures provided in upload")
	}

	allTeams := GetAllCollegeTeams()
	teamByID := make(map[uint]structs.CollegeTeam)
	teamByLookup := make(map[string]structs.CollegeTeam)
	for _, t := range allTeams {
		teamByID[t.ID] = t
		teamByLookup[strings.ToUpper(strings.TrimSpace(t.TeamAbbr))] = t
		teamByLookup[strings.ToUpper(strings.TrimSpace(t.TeamName))] = t
		fullName := strings.ToUpper(strings.TrimSpace(t.TeamName + " " + t.Mascot))
		teamByLookup[fullName] = t
	}

	allStadiums := GetAllStadiums()
	stadiumByName := make(map[string]structs.Stadium)
	for _, s := range allStadiums {
		stadiumByName[strings.ToUpper(strings.TrimSpace(s.StadiumName))] = s
	}

	seasonID := ts.CollegeSeasonID
	standings := GetAllCollegeStandingsBySeasonID(strconv.Itoa(seasonID))
	standingsMap := MakeCollegeStandingsMapByTeamID(standings)

	var savedGames []structs.CollegeGame

	for _, f := range req.Fixtures {
		var ht structs.CollegeTeam
		var at structs.CollegeTeam

		if f.HomeTeamID > 0 {
			ht = teamByID[f.HomeTeamID]
		}
		if ht.ID == 0 && strings.TrimSpace(f.HomeTeam) != "" {
			ht = teamByLookup[strings.ToUpper(strings.TrimSpace(f.HomeTeam))]
		}

		if f.AwayTeamID > 0 {
			at = teamByID[f.AwayTeamID]
		}
		if at.ID == 0 && strings.TrimSpace(f.AwayTeam) != "" {
			at = teamByLookup[strings.ToUpper(strings.TrimSpace(f.AwayTeam))]
		}

		if ht.ID == 0 && at.ID == 0 {
			continue
		}
		if ht.ID == 0 || at.ID == 0 {
			return nil, fmt.Errorf("could not resolve teams for bowl '%s' (Home: '%s', Away: '%s')", f.GameTitle, f.HomeTeam, f.AwayTeam)
		}
		if ht.ID == at.ID {
			return nil, fmt.Errorf("Home Team and Away Team cannot be the same in '%s'", f.GameTitle)
		}

		week := f.Week
		if week < 15 {
			week = 16
		}
		weekID := int(util.GetWeekID(uint(seasonID), uint(week)))

		timeSlot := strings.TrimSpace(f.TimeSlot)
		if timeSlot == "" {
			timeSlot = "Saturday Afternoon"
		}

		gameTitle := strings.TrimSpace(f.GameTitle)
		if gameTitle == "" {
			gameTitle = "Bowl Game"
		}

		stadiumID := f.StadiumID
		stadium := strings.TrimSpace(f.Stadium)
		city := strings.TrimSpace(f.City)
		state := strings.TrimSpace(f.State)
		isDomed := f.IsDomed

		if stadium != "" && stadiumID == 0 {
			if stRec, ok := stadiumByName[strings.ToUpper(stadium)]; ok {
				stadiumID = stRec.ID
				if city == "" {
					city = stRec.City
				}
				if state == "" {
					state = stRec.State
				}
				isDomed = isDomed || stRec.IsDomed
			}
		}
		if !f.IsNeutral && stadium == "" {
			stadiumID = ht.StadiumID
			stadium = ht.Stadium
			city = ht.City
			state = ht.State
		}

		homeRank := uint(standingsMap[ht.ID].Rank)
		awayRank := uint(standingsMap[at.ID].Rank)
		isConf := ht.ConferenceID == at.ConferenceID
		isDiv := isConf && ht.DivisionID == at.DivisionID && ht.DivisionID > 0
		confID := uint(0)
		if isConf {
			confID = uint(ht.ConferenceID)
		}

		var existing structs.CollegeGame
		if f.GameID > 0 {
			db.Where("id = ?", f.GameID).First(&existing)
		}
		if existing.ID == 0 {
			db.Where("season_id = ? AND is_bowl_game = ? AND is_playoff_game = ? AND game_title = ?", seasonID, true, false, gameTitle).First(&existing)
		}

		if existing.ID > 0 && !existing.GameComplete {
			existing.Week = week
			existing.WeekID = weekID
			existing.HomeTeamID = int(ht.ID)
			existing.HomeTeam = ht.TeamAbbr
			existing.HomeTeamCoach = ht.Coach
			existing.HomeTeamRank = homeRank
			existing.AwayTeamID = int(at.ID)
			existing.AwayTeam = at.TeamAbbr
			existing.AwayTeamCoach = at.Coach
			existing.AwayTeamRank = awayRank
			existing.TimeSlot = timeSlot
			existing.StadiumID = stadiumID
			existing.Stadium = stadium
			existing.City = city
			existing.State = state
			existing.IsDomed = isDomed
			existing.IsNeutral = f.IsNeutral
			existing.IsNightGame = strings.Contains(timeSlot, "Evening") || strings.Contains(timeSlot, "Night")
			existing.IsConference = isConf
			existing.IsDivisional = isDiv
			existing.ConferenceID = confID
			existing.IsBowlGame = true
			existing.GameTitle = gameTitle
			repository.SaveCFBGameRecord(existing, db)
			savedGames = append(savedGames, existing)
			continue
		}

		newGame := structs.CollegeGame{
			SeasonID:       seasonID,
			WeekID:         weekID,
			Week:           week,
			HomeTeamID:     int(ht.ID),
			HomeTeam:       ht.TeamAbbr,
			HomeTeamCoach:  ht.Coach,
			HomeTeamRank:   homeRank,
			AwayTeamID:     int(at.ID),
			AwayTeam:       at.TeamAbbr,
			AwayTeamCoach:  at.Coach,
			AwayTeamRank:   awayRank,
			TimeSlot:       timeSlot,
			StadiumID:      stadiumID,
			Stadium:        stadium,
			City:           city,
			State:          state,
			IsDomed:        isDomed,
			IsNeutral:      f.IsNeutral,
			IsNightGame:    strings.Contains(timeSlot, "Evening") || strings.Contains(timeSlot, "Night"),
			IsConference:   isConf,
			IsDivisional:   isDiv,
			ConferenceID:   confID,
			IsBowlGame:     true,
			IsPlayoffGame:  false,
			GameTitle:      gameTitle,
		}
		if err := db.Create(&newGame).Error; err != nil {
			return nil, err
		}
		savedGames = append(savedGames, newGame)
	}

	if len(savedGames) == 0 {
		return nil, errors.New("no valid bowl matchups with both Home Team and Away Team were found in the uploaded file")
	}

	GenerateWeatherForGames()
	repository.AssignCFBStreamChannels(savedGames)
	return savedGames, nil
}

func resolveSeedsToTeams(seeds []structs.CFBPlayoffSeedDTO, expectedCount int, requireFBS bool) (map[uint]structs.CollegeTeam, error) {
	if len(seeds) != expectedCount {
		return nil, fmt.Errorf("expected exactly %d ranked teams, got %d", expectedCount, len(seeds))
	}

	allTeams := GetAllCollegeTeams()
	teamByID := make(map[uint]structs.CollegeTeam)
	teamByLookup := make(map[string]structs.CollegeTeam)
	for _, t := range allTeams {
		teamByID[t.ID] = t
		teamByLookup[strings.ToUpper(strings.TrimSpace(t.TeamAbbr))] = t
		teamByLookup[strings.ToUpper(strings.TrimSpace(t.TeamName))] = t
		fullName := strings.ToUpper(strings.TrimSpace(t.TeamName + " " + t.Mascot))
		teamByLookup[fullName] = t
	}

	seedMap := make(map[uint]structs.CollegeTeam)
	usedTeamIDs := make(map[uint]bool)

	for _, s := range seeds {
		if s.Seed < 1 || int(s.Seed) > expectedCount {
			return nil, fmt.Errorf("invalid seed number %d (must be between 1 and %d)", s.Seed, expectedCount)
		}
		if _, exists := seedMap[s.Seed]; exists {
			return nil, fmt.Errorf("duplicate seed number #%d", s.Seed)
		}

		var team structs.CollegeTeam
		if s.TeamID > 0 {
			team = teamByID[s.TeamID]
		}
		if team.ID == 0 && strings.TrimSpace(s.TeamAbbr) != "" {
			team = teamByLookup[strings.ToUpper(strings.TrimSpace(s.TeamAbbr))]
		}
		if team.ID == 0 {
			return nil, fmt.Errorf("could not find team for seed #%d", s.Seed)
		}
		if requireFBS && !team.IsFBS {
			return nil, fmt.Errorf("team %s (seed #%d) is not an FBS team", team.TeamName, s.Seed)
		}
		if !requireFBS && team.IsFBS {
			return nil, fmt.Errorf("team %s (seed #%d) is an FBS team, not an FCS team", team.TeamName, s.Seed)
		}
		if usedTeamIDs[team.ID] {
			return nil, fmt.Errorf("team %s is selected more than once", team.TeamName)
		}
		usedTeamIDs[team.ID] = true
		seedMap[s.Seed] = team
	}

	return seedMap, nil
}

func GenerateFBSPlayoffSchedule(req structs.CFBPlayoffScheduleRequest) ([]structs.CollegeGame, error) {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	seedMap, err := resolveSeedsToTeams(req.Seeds, 12, true)
	if err != nil {
		return nil, err
	}

	seasonID := ts.CollegeSeasonID
	seasonYear := ts.Season
	if seasonYear == 0 {
		seasonYear = 2020 + seasonID
	}

	// Remove any existing unplayed FBS playoff games for this season so re-uploading replaces cleanly
	db.Where("season_id = ? AND is_playoff_game = ? AND is_bowl_game = ? AND game_complete = ?", seasonID, true, true, false).
		Delete(&structs.CollegeGame{})

	allStadiums := GetAllStadiums()
	stadiumMap := make(map[uint]structs.Stadium)
	for _, s := range allStadiums {
		stadiumMap[s.ID] = s
	}

	w17ID := int(util.GetWeekID(uint(seasonID), 17))
	w18ID := int(util.GetWeekID(uint(seasonID), 18))
	w19ID := int(util.GetWeekID(uint(seasonID), 19))
	w20ID := int(util.GetWeekID(uint(seasonID), 20))

	// 1. Week 20: National Championship Game
	natty := structs.CollegeGame{
		SeasonID:               seasonID,
		WeekID:                 w20ID,
		Week:                   20,
		TimeSlot:               "Saturday Evening",
		StadiumID:              25,
		Stadium:                "Mercedes-Benz Stadium",
		City:                   "Atlanta",
		State:                  "GA",
		IsDomed:                true,
		IsNeutral:              true,
		IsNightGame:            true,
		IsBowlGame:             true,
		IsPlayoffGame:          true,
		IsNationalChampionship: true,
		GameTitle:              fmt.Sprintf("%d College Football Playoff National Championship Game", seasonYear),
	}
	if err := db.Create(&natty).Error; err != nil {
		return nil, err
	}

	// 2. Week 19: Semifinals (Fiesta Bowl & Peach Bowl)
	semi1 := structs.CollegeGame{
		SeasonID:      seasonID,
		WeekID:        w19ID,
		Week:          19,
		TimeSlot:      "Saturday Afternoon",
		StadiumID:     294,
		Stadium:       "State Farm Stadium",
		City:          "Glendale",
		State:         "AZ",
		IsDomed:       true,
		IsNeutral:     true,
		IsNightGame:   false,
		IsBowlGame:    true,
		IsPlayoffGame: true,
		GameTitle:     fmt.Sprintf("%d Fiesta Bowl - CFP Semifinals", seasonYear),
		NextGameID:    natty.ID,
		NextGameHOA:   "H",
	}
	if err := db.Create(&semi1).Error; err != nil {
		return nil, err
	}

	semi2 := structs.CollegeGame{
		SeasonID:      seasonID,
		WeekID:        w19ID,
		Week:          19,
		TimeSlot:      "Saturday Evening",
		StadiumID:     25,
		Stadium:       "Mercedes-Benz Stadium",
		City:          "Atlanta",
		State:         "GA",
		IsDomed:       true,
		IsNeutral:     true,
		IsNightGame:   true,
		IsBowlGame:    true,
		IsPlayoffGame: true,
		GameTitle:     fmt.Sprintf("%d Peach Bowl - CFP Semifinals", seasonYear),
		NextGameID:    natty.ID,
		NextGameHOA:   "A",
	}
	if err := db.Create(&semi2).Error; err != nil {
		return nil, err
	}

	// 3. Week 18: Quarterfinals (Seeds 1-4 have First Round byes and are Home Teams)
	s1 := seedMap[1]
	s2 := seedMap[2]
	s3 := seedMap[3]
	s4 := seedMap[4]

	// QF1: Rose Bowl (#1 vs Winner of #8/#9) -> Semi 1 Home
	qf1 := structs.CollegeGame{
		SeasonID:        seasonID,
		WeekID:          w18ID,
		Week:            18,
		HomeTeamID:      int(s1.ID),
		HomeTeam:        s1.TeamAbbr,
		HomeTeamCoach:   s1.Coach,
		HomeTeamRank:    1,
		HomePreviousBye: true,
		TimeSlot:        "Saturday Evening",
		StadiumID:       156,
		Stadium:         "Rose Bowl",
		City:            "Pasadena",
		State:           "CA",
		IsDomed:         false,
		IsNeutral:       true,
		IsNightGame:     true,
		IsBowlGame:      true,
		IsPlayoffGame:   true,
		GameTitle:       fmt.Sprintf("%d Rose Bowl - CFP Quarterfinals", seasonYear),
		NextGameID:      semi1.ID,
		NextGameHOA:     "H",
	}
	if err := db.Create(&qf1).Error; err != nil {
		return nil, err
	}

	// QF2: Cotton Bowl (#4 vs Winner of #5/#12) -> Semi 1 Away
	qf2 := structs.CollegeGame{
		SeasonID:        seasonID,
		WeekID:          w18ID,
		Week:            18,
		HomeTeamID:      int(s4.ID),
		HomeTeam:        s4.TeamAbbr,
		HomeTeamCoach:   s4.Coach,
		HomeTeamRank:    4,
		HomePreviousBye: true,
		TimeSlot:        "Friday Night",
		StadiumID:       99,
		Stadium:         "AT&T Stadium",
		City:            "Arlington",
		State:           "TX",
		IsDomed:         true,
		IsNeutral:       true,
		IsNightGame:     true,
		IsBowlGame:      true,
		IsPlayoffGame:   true,
		GameTitle:       fmt.Sprintf("%d Cotton Bowl - CFP Quarterfinals", seasonYear),
		NextGameID:      semi1.ID,
		NextGameHOA:     "A",
	}
	if err := db.Create(&qf2).Error; err != nil {
		return nil, err
	}

	// QF3: Sugar Bowl (#2 vs Winner of #7/#10) -> Semi 2 Home
	qf3 := structs.CollegeGame{
		SeasonID:        seasonID,
		WeekID:          w18ID,
		Week:            18,
		HomeTeamID:      int(s2.ID),
		HomeTeam:        s2.TeamAbbr,
		HomeTeamCoach:   s2.Coach,
		HomeTeamRank:    2,
		HomePreviousBye: true,
		TimeSlot:        "Saturday Night",
		StadiumID:       100,
		Stadium:         "Caesars Superdome",
		City:            "New Orleans",
		State:           "LA",
		IsDomed:         true,
		IsNeutral:       true,
		IsNightGame:     true,
		IsBowlGame:      true,
		IsPlayoffGame:   true,
		GameTitle:       fmt.Sprintf("%d Sugar Bowl - CFP Quarterfinals", seasonYear),
		NextGameID:      semi2.ID,
		NextGameHOA:     "H",
	}
	if err := db.Create(&qf3).Error; err != nil {
		return nil, err
	}

	// QF4: Orange Bowl (#3 vs Winner of #6/#11) -> Semi 2 Away
	qf4 := structs.CollegeGame{
		SeasonID:        seasonID,
		WeekID:          w18ID,
		Week:            18,
		HomeTeamID:      int(s3.ID),
		HomeTeam:        s3.TeamAbbr,
		HomeTeamCoach:   s3.Coach,
		HomeTeamRank:    3,
		HomePreviousBye: true,
		TimeSlot:        "Saturday Afternoon",
		StadiumID:       55,
		Stadium:         "Hard Rock Stadium",
		City:            "Miami Gardens",
		State:           "FL",
		IsDomed:         false,
		IsNeutral:       true,
		IsNightGame:     false,
		IsBowlGame:      true,
		IsPlayoffGame:   true,
		GameTitle:       fmt.Sprintf("%d Orange Bowl - CFP Quarterfinals", seasonYear),
		NextGameID:      semi2.ID,
		NextGameHOA:     "A",
	}
	if err := db.Create(&qf4).Error; err != nil {
		return nil, err
	}

	// 4. Week 17: First Round (Seeds 5-8 host Seeds 9-12 on campus)
	type r1Def struct {
		homeSeed uint
		awaySeed uint
		timeSlot string
		nextID   uint
	}
	r1Defs := []r1Def{
		{homeSeed: 8, awaySeed: 9, timeSlot: "Friday Night", nextID: qf1.ID},
		{homeSeed: 5, awaySeed: 12, timeSlot: "Saturday Morning", nextID: qf2.ID},
		{homeSeed: 6, awaySeed: 11, timeSlot: "Saturday Afternoon", nextID: qf4.ID},
		{homeSeed: 7, awaySeed: 10, timeSlot: "Saturday Evening", nextID: qf3.ID},
	}

	var r1Games []structs.CollegeGame
	for _, d := range r1Defs {
		ht := seedMap[d.homeSeed]
		at := seedMap[d.awaySeed]
		st := stadiumMap[ht.StadiumID]
		isConf := ht.ConferenceID == at.ConferenceID
		confID := uint(0)
		if isConf {
			confID = uint(ht.ConferenceID)
		}
		g := structs.CollegeGame{
			SeasonID:      seasonID,
			WeekID:        w17ID,
			Week:          17,
			HomeTeamID:    int(ht.ID),
			HomeTeam:      ht.TeamAbbr,
			HomeTeamCoach: ht.Coach,
			HomeTeamRank:  d.homeSeed,
			AwayTeamID:    int(at.ID),
			AwayTeam:      at.TeamAbbr,
			AwayTeamCoach: at.Coach,
			AwayTeamRank:  d.awaySeed,
			TimeSlot:      d.timeSlot,
			StadiumID:     ht.StadiumID,
			Stadium:       ht.Stadium,
			City:          ht.City,
			State:         ht.State,
			IsDomed:       st.IsDomed,
			IsNeutral:     false,
			IsNightGame:   strings.Contains(d.timeSlot, "Evening") || strings.Contains(d.timeSlot, "Night"),
			IsConference:  isConf,
			ConferenceID:  confID,
			IsBowlGame:    true,
			IsPlayoffGame: true,
			GameTitle:     fmt.Sprintf("%d FBS Playoffs Round 1", seasonYear),
			NextGameID:    d.nextID,
			NextGameHOA:   "A",
		}
		if err := db.Create(&g).Error; err != nil {
			return nil, err
		}
		r1Games = append(r1Games, g)
	}

	GenerateWeatherForGames()

	allCreated := append([]structs.CollegeGame{}, r1Games...)
	allCreated = append(allCreated, qf1, qf2, qf3, qf4, semi1, semi2, natty)
	repository.AssignCFBStreamChannels(allCreated)
	return allCreated, nil
}

func GenerateFCSPlayoffSchedule(req structs.CFBPlayoffScheduleRequest) ([]structs.CollegeGame, error) {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	seedMap, err := resolveSeedsToTeams(req.Seeds, 24, false)
	if err != nil {
		return nil, err
	}

	seasonID := ts.CollegeSeasonID

	// Remove any existing unplayed FCS playoff games for this season
	db.Where("season_id = ? AND is_playoff_game = ? AND is_bowl_game = ? AND game_complete = ?", seasonID, true, false, false).
		Delete(&structs.CollegeGame{})

	allStadiums := GetAllStadiums()
	stadiumMap := make(map[uint]structs.Stadium)
	for _, s := range allStadiums {
		stadiumMap[s.ID] = s
	}

	w16ID := int(util.GetWeekID(uint(seasonID), 16))
	w17ID := int(util.GetWeekID(uint(seasonID), 17))
	w18ID := int(util.GetWeekID(uint(seasonID), 18))
	w19ID := int(util.GetWeekID(uint(seasonID), 19))
	w20ID := int(util.GetWeekID(uint(seasonID), 20))

	// 1. Week 20: FCS National Championship Game
	fcsNatty := structs.CollegeGame{
		SeasonID:               seasonID,
		WeekID:                 w20ID,
		Week:                   20,
		TimeSlot:               "Saturday Afternoon",
		StadiumID:              288,
		Stadium:                "Toyota Stadium",
		City:                   "Frisco",
		State:                  "TX",
		IsDomed:                false,
		IsNeutral:              true,
		IsNightGame:            false,
		IsBowlGame:             false,
		IsPlayoffGame:          true,
		IsNationalChampionship: true,
		GameTitle:              "FCS National Championship Game",
	}
	if err := db.Create(&fcsNatty).Error; err != nil {
		return nil, err
	}

	// 2. Week 19: FCS Semifinals (2 games)
	sf1 := structs.CollegeGame{
		SeasonID:      seasonID,
		WeekID:        w19ID,
		Week:          19,
		TimeSlot:      "Saturday Afternoon",
		IsNeutral:     false,
		IsPlayoffGame: true,
		GameTitle:     "FCS Playoffs Semifinals",
		NextGameID:    fcsNatty.ID,
		NextGameHOA:   "H",
	}
	if err := db.Create(&sf1).Error; err != nil {
		return nil, err
	}

	sf2 := structs.CollegeGame{
		SeasonID:      seasonID,
		WeekID:        w19ID,
		Week:          19,
		TimeSlot:      "Saturday Evening",
		IsNeutral:     false,
		IsNightGame:   true,
		IsPlayoffGame: true,
		GameTitle:     "FCS Playoffs Semifinals",
		NextGameID:    fcsNatty.ID,
		NextGameHOA:   "A",
	}
	if err := db.Create(&sf2).Error; err != nil {
		return nil, err
	}

	// 3. Week 18: FCS Quarterfinals (4 games)
	type fcsQFDef struct {
		timeSlot string
		nextID   uint
		nextHOA  string
	}
	qfDefs := []fcsQFDef{
		{timeSlot: "Saturday Morning", nextID: sf1.ID, nextHOA: "H"},
		{timeSlot: "Saturday Afternoon", nextID: sf1.ID, nextHOA: "A"},
		{timeSlot: "Saturday Evening", nextID: sf2.ID, nextHOA: "H"},
		{timeSlot: "Saturday Night", nextID: sf2.ID, nextHOA: "A"},
	}
	var qfGames []structs.CollegeGame
	for _, d := range qfDefs {
		g := structs.CollegeGame{
			SeasonID:      seasonID,
			WeekID:        w18ID,
			Week:          18,
			TimeSlot:      d.timeSlot,
			IsNeutral:     false,
			IsNightGame:   strings.Contains(d.timeSlot, "Evening") || strings.Contains(d.timeSlot, "Night"),
			IsPlayoffGame: true,
			GameTitle:     "FCS Playoffs Quarterfinals",
			NextGameID:    d.nextID,
			NextGameHOA:   d.nextHOA,
		}
		if err := db.Create(&g).Error; err != nil {
			return nil, err
		}
		qfGames = append(qfGames, g)
	}

	// 4. Week 17: FCS Second Round (8 games, Seeds 1-8 have First Round byes and host)
	// Bracket pods:
	// Pod 1 (QF 0): #1 (H) vs #16/#17 winner, #8 (A) vs #9/#24 winner
	// Pod 2 (QF 1): #4 (H) vs #13/#20 winner, #5 (A) vs #12/#21 winner
	// Pod 3 (QF 2): #2 (H) vs #15/#18 winner, #7 (A) vs #10/#23 winner
	// Pod 4 (QF 3): #3 (H) vs #14/#19 winner, #6 (A) vs #11/#22 winner
	type fcsSRDef struct {
		homeSeed uint
		r1Home   uint
		r1Away   uint
		timeSlot string
		qfIndex  int
		nextHOA  string
	}
	srDefs := []fcsSRDef{
		{homeSeed: 1, r1Home: 16, r1Away: 17, timeSlot: "Saturday Morning", qfIndex: 0, nextHOA: "H"},
		{homeSeed: 8, r1Home: 9, r1Away: 24, timeSlot: "Saturday Morning", qfIndex: 0, nextHOA: "A"},
		{homeSeed: 4, r1Home: 13, r1Away: 20, timeSlot: "Saturday Afternoon", qfIndex: 1, nextHOA: "H"},
		{homeSeed: 5, r1Home: 12, r1Away: 21, timeSlot: "Saturday Afternoon", qfIndex: 1, nextHOA: "A"},
		{homeSeed: 2, r1Home: 15, r1Away: 18, timeSlot: "Saturday Evening", qfIndex: 2, nextHOA: "H"},
		{homeSeed: 7, r1Home: 10, r1Away: 23, timeSlot: "Saturday Evening", qfIndex: 2, nextHOA: "A"},
		{homeSeed: 3, r1Home: 14, r1Away: 19, timeSlot: "Saturday Night", qfIndex: 3, nextHOA: "H"},
		{homeSeed: 6, r1Home: 11, r1Away: 22, timeSlot: "Saturday Night", qfIndex: 3, nextHOA: "A"},
	}

	var srGames []structs.CollegeGame
	var r1Games []structs.CollegeGame

	for _, d := range srDefs {
		byeTeam := seedMap[d.homeSeed]
		byeSt := stadiumMap[byeTeam.StadiumID]

		srGame := structs.CollegeGame{
			SeasonID:        seasonID,
			WeekID:          w17ID,
			Week:            17,
			HomeTeamID:      int(byeTeam.ID),
			HomeTeam:        byeTeam.TeamAbbr,
			HomeTeamCoach:   byeTeam.Coach,
			HomeTeamRank:    d.homeSeed,
			HomePreviousBye: true,
			TimeSlot:        d.timeSlot,
			StadiumID:       byeTeam.StadiumID,
			Stadium:         byeTeam.Stadium,
			City:            byeTeam.City,
			State:           byeTeam.State,
			IsDomed:         byeSt.IsDomed,
			IsNeutral:       false,
			IsNightGame:     strings.Contains(d.timeSlot, "Evening") || strings.Contains(d.timeSlot, "Night"),
			IsPlayoffGame:   true,
			GameTitle:       "FCS Playoffs Second Round",
			NextGameID:      qfGames[d.qfIndex].ID,
			NextGameHOA:     d.nextHOA,
		}
		if err := db.Create(&srGame).Error; err != nil {
			return nil, err
		}
		srGames = append(srGames, srGame)

		// Corresponding Week 16 First Round game feeding into this Second Round game as Away ("A")
		ht := seedMap[d.r1Home]
		at := seedMap[d.r1Away]
		r1St := stadiumMap[ht.StadiumID]
		isConf := ht.ConferenceID == at.ConferenceID
		confID := uint(0)
		if isConf {
			confID = uint(ht.ConferenceID)
		}

		r1Game := structs.CollegeGame{
			SeasonID:      seasonID,
			WeekID:        w16ID,
			Week:          16,
			HomeTeamID:    int(ht.ID),
			HomeTeam:      ht.TeamAbbr,
			HomeTeamCoach: ht.Coach,
			HomeTeamRank:  d.r1Home,
			AwayTeamID:    int(at.ID),
			AwayTeam:      at.TeamAbbr,
			AwayTeamCoach: at.Coach,
			AwayTeamRank:  d.r1Away,
			TimeSlot:      d.timeSlot,
			StadiumID:     ht.StadiumID,
			Stadium:       ht.Stadium,
			City:          ht.City,
			State:         ht.State,
			IsDomed:       r1St.IsDomed,
			IsNeutral:     false,
			IsNightGame:   strings.Contains(d.timeSlot, "Evening") || strings.Contains(d.timeSlot, "Night"),
			IsConference:  isConf,
			ConferenceID:  confID,
			IsPlayoffGame: true,
			GameTitle:     "FCS Playoffs First Round",
			NextGameID:    srGame.ID,
			NextGameHOA:   "A",
		}
		if err := db.Create(&r1Game).Error; err != nil {
			return nil, err
		}
		r1Games = append(r1Games, r1Game)
	}

	GenerateWeatherForGames()

	allCreated := append([]structs.CollegeGame{}, r1Games...)
	allCreated = append(allCreated, srGames...)
	allCreated = append(allCreated, qfGames...)
	allCreated = append(allCreated, sf1, sf2, fcsNatty)
	repository.AssignCFBStreamChannels(allCreated)
	return allCreated, nil
}
