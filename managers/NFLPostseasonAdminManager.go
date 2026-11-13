package managers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
)

func GetNFLPostseasonGamesForCurrentSeason() []structs.NFLGame {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	var games []structs.NFLGame
	db.Order("week asc, id asc").
		Where("season_id = ? AND (is_playoff_game = ? OR is_super_bowl = ? OR is_conference_championship = ? OR week >= 19)",
			ts.NFLSeasonID, true, true, true, 19).
		Find(&games)

	return games
}

func DeleteNFLPostseasonGame(gameID string) error {
	db := dbprovider.GetInstance().GetDB()
	game := GetNFLGameByGameID(gameID)
	if game.ID == 0 {
		return errors.New("game not found")
	}
	if game.GameComplete {
		return errors.New("cannot delete a completed game")
	}
	return db.Delete(&game).Error
}

func GenerateNFLPlayoffSchedule(req structs.NFLPlayoffScheduleRequest) ([]structs.NFLGame, error) {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID := int(ts.NFLSeasonID)
	seasonYear := ts.Season

	if len(req.Seeds) != 14 {
		return nil, fmt.Errorf("expected 14 total playoff seeds (7 AFC, 7 NFC), got %d", len(req.Seeds))
	}

	allTeams := GetAllNFLTeams()
	teamMap := make(map[uint]structs.NFLTeam)
	for _, t := range allTeams {
		teamMap[t.ID] = t
	}

	afcSeedMap := make(map[uint]structs.NFLTeam)
	nfcSeedMap := make(map[uint]structs.NFLTeam)
	seenTeams := make(map[uint]bool)

	for _, s := range req.Seeds {
		if s.Seed < 1 || s.Seed > 7 {
			return nil, fmt.Errorf("invalid seed number %d: must be 1 to 7", s.Seed)
		}
		if s.TeamID == 0 {
			return nil, fmt.Errorf("seed %d (%s) has no team selected", s.Seed, s.Conference)
		}
		if seenTeams[s.TeamID] {
			return nil, fmt.Errorf("team %d is assigned to multiple playoff seeds", s.TeamID)
		}
		seenTeams[s.TeamID] = true

		team, exists := teamMap[s.TeamID]
		if !exists {
			return nil, fmt.Errorf("team ID %d not found in NFL teams", s.TeamID)
		}

		conf := strings.ToUpper(strings.TrimSpace(s.Conference))
		if conf != "AFC" && conf != "NFC" {
			// Infer from team conference if unspecified or full name
			teamConf := strings.ToUpper(team.Conference)
			if strings.Contains(teamConf, "AFC") || strings.Contains(teamConf, "AMERICAN") || team.ConferenceID == 1 {
				conf = "AFC"
			} else {
				conf = "NFC"
			}
		}

		if conf == "AFC" {
			if _, exists := afcSeedMap[s.Seed]; exists {
				return nil, fmt.Errorf("duplicate AFC seed #%d", s.Seed)
			}
			afcSeedMap[s.Seed] = team
		} else {
			if _, exists := nfcSeedMap[s.Seed]; exists {
				return nil, fmt.Errorf("duplicate NFC seed #%d", s.Seed)
			}
			nfcSeedMap[s.Seed] = team
		}
	}

	for seed := uint(1); seed <= 7; seed++ {
		if _, ok := afcSeedMap[seed]; !ok {
			return nil, fmt.Errorf("missing AFC seed #%d", seed)
		}
		if _, ok := nfcSeedMap[seed]; !ok {
			return nil, fmt.Errorf("missing NFC seed #%d", seed)
		}
	}

	// Check existing postseason games for this season
	var existingGames []structs.NFLGame
	db.Where("season_id = ? AND (is_playoff_game = ? OR is_super_bowl = ? OR is_conference_championship = ? OR week >= 19)",
		seasonID, true, true, true, 19).Find(&existingGames)
	for _, g := range existingGames {
		if g.GameComplete {
			return nil, errors.New("cannot generate new playoff schedule: playoff games have already been played this season")
		}
	}
	if len(existingGames) > 0 {
		db.Delete(&existingGames)
	}

	allStadiums := GetAllStadiums()
	stadiumMap := make(map[uint]structs.Stadium)
	for _, s := range allStadiums {
		stadiumMap[s.ID] = s
	}

	w19ID := int(util.GetWeekID(uint(seasonID), 19))
	w20ID := int(util.GetWeekID(uint(seasonID), 20))
	w21ID := int(util.GetWeekID(uint(seasonID), 21))
	w22ID := int(util.GetWeekID(uint(seasonID), 22))

	// 1. Week 22: Super Bowl (Neutral site, Allegiant Stadium in Las Vegas)
	superBowl := structs.NFLGame{
		SeasonID:      seasonID,
		WeekID:        w22ID,
		Week:          22,
		TimeSlot:      "Sunday Evening",
		StadiumID:     155,
		Stadium:       "Allegiant Stadium",
		City:          "Las Vegas",
		State:         "NV",
		IsDomed:       true,
		IsNeutral:     true,
		IsNightGame:   true,
		IsPlayoffGame: true,
		IsSuperBowl:   true,
		GameTitle:     fmt.Sprintf("%d Super Bowl", seasonYear),
		HomeTeam:      "AFC Champion",
		AwayTeam:      "NFC Champion",
	}
	if err := db.Create(&superBowl).Error; err != nil {
		return nil, err
	}

	// 2. Week 21: Conference Championship Games
	afcChamp := structs.NFLGame{
		SeasonID:                 seasonID,
		WeekID:                   w21ID,
		Week:                     21,
		TimeSlot:                 "Sunday Afternoon",
		IsPlayoffGame:            true,
		IsConferenceChampionship: true,
		IsConference:             true,
		GameTitle:                fmt.Sprintf("%d AFC Championship Game", seasonYear),
		NextGameID:               superBowl.ID,
		NextGameHOA:              "H",
		HomeTeam:                 "TBD",
		AwayTeam:                 "TBD",
	}
	if err := db.Create(&afcChamp).Error; err != nil {
		return nil, err
	}

	nfcChamp := structs.NFLGame{
		SeasonID:                 seasonID,
		WeekID:                   w21ID,
		Week:                     21,
		TimeSlot:                 "Sunday Evening",
		IsPlayoffGame:            true,
		IsConferenceChampionship: true,
		IsConference:             true,
		GameTitle:                fmt.Sprintf("%d NFC Championship Game", seasonYear),
		NextGameID:               superBowl.ID,
		NextGameHOA:              "A",
		HomeTeam:                 "TBD",
		AwayTeam:                 "TBD",
	}
	if err := db.Create(&nfcChamp).Error; err != nil {
		return nil, err
	}

	// 3. Week 20: Divisional Round
	// Seed 1 teams have first round byes and host in Week 20
	afc1 := afcSeedMap[1]
	nfc1 := nfcSeedMap[1]
	afc1Stadium := stadiumMap[afc1.StadiumID]
	nfc1Stadium := stadiumMap[nfc1.StadiumID]

	// AFC Div 1: #1 Seed (Home) vs Winner of #4 vs #5 (Away) -> advances to AFC Champ Game Home ("H")
	afcDiv1 := structs.NFLGame{
		SeasonID:      seasonID,
		WeekID:        w20ID,
		Week:          20,
		TimeSlot:      "Saturday Afternoon",
		StadiumID:     afc1.StadiumID,
		Stadium:       afc1.Stadium,
		City:          afc1.City,
		State:         afc1.State,
		IsDomed:       afc1Stadium.IsDomed,
		HomeTeamID:    int(afc1.ID),
		HomeTeam:      afc1.TeamName + " " + afc1.Mascot,
		HomeTeamCoach: afc1.NFLCoachName,
		AwayTeam:      "Winner AFC 4/5",
		IsPlayoffGame: true,
		IsDivisional:  true,
		GameTitle:     fmt.Sprintf("%d AFC Divisional Round - #1 %s vs TBD", seasonYear, afc1.TeamAbbr),
		NextGameID:    afcChamp.ID,
		NextGameHOA:   "H",
	}
	if err := db.Create(&afcDiv1).Error; err != nil {
		return nil, err
	}

	// AFC Div 2: Winner of #2/#7 (Home) vs Winner of #3/#6 (Away) -> advances to AFC Champ Game Away ("A")
	afcDiv2 := structs.NFLGame{
		SeasonID:      seasonID,
		WeekID:        w20ID,
		Week:          20,
		TimeSlot:      "Sunday Afternoon",
		HomeTeam:      "Winner AFC 2/7",
		AwayTeam:      "Winner AFC 3/6",
		IsPlayoffGame: true,
		IsDivisional:  true,
		GameTitle:     fmt.Sprintf("%d AFC Divisional Round - TBD vs TBD", seasonYear),
		NextGameID:    afcChamp.ID,
		NextGameHOA:   "A",
	}
	if err := db.Create(&afcDiv2).Error; err != nil {
		return nil, err
	}

	// NFC Div 1: #1 Seed (Home) vs Winner of #4 vs #5 (Away) -> advances to NFC Champ Game Home ("H")
	nfcDiv1 := structs.NFLGame{
		SeasonID:      seasonID,
		WeekID:        w20ID,
		Week:          20,
		TimeSlot:      "Saturday Evening",
		StadiumID:     nfc1.StadiumID,
		Stadium:       nfc1.Stadium,
		City:          nfc1.City,
		State:         nfc1.State,
		IsDomed:       nfc1Stadium.IsDomed,
		HomeTeamID:    int(nfc1.ID),
		HomeTeam:      nfc1.TeamName + " " + nfc1.Mascot,
		HomeTeamCoach: nfc1.NFLCoachName,
		AwayTeam:      "Winner NFC 4/5",
		IsPlayoffGame: true,
		IsDivisional:  true,
		GameTitle:     fmt.Sprintf("%d NFC Divisional Round - #1 %s vs TBD", seasonYear, nfc1.TeamAbbr),
		NextGameID:    nfcChamp.ID,
		NextGameHOA:   "H",
	}
	if err := db.Create(&nfcDiv1).Error; err != nil {
		return nil, err
	}

	// NFC Div 2: Winner of #2/#7 (Home) vs Winner of #3/#6 (Away) -> advances to NFC Champ Game Away ("A")
	nfcDiv2 := structs.NFLGame{
		SeasonID:      seasonID,
		WeekID:        w20ID,
		Week:          20,
		TimeSlot:      "Sunday Evening",
		HomeTeam:      "Winner NFC 2/7",
		AwayTeam:      "Winner NFC 3/6",
		IsPlayoffGame: true,
		IsDivisional:  true,
		GameTitle:     fmt.Sprintf("%d NFC Divisional Round - TBD vs TBD", seasonYear),
		NextGameID:    nfcChamp.ID,
		NextGameHOA:   "A",
	}
	if err := db.Create(&nfcDiv2).Error; err != nil {
		return nil, err
	}

	// 4. Week 19: Wild Card Weekend (6 Games)
	buildWCGame := func(homeTeam, awayTeam structs.NFLTeam, title, timeslot, nextHOA string, nextGameID uint) structs.NFLGame {
		st := stadiumMap[homeTeam.StadiumID]
		return structs.NFLGame{
			SeasonID:      seasonID,
			WeekID:        w19ID,
			Week:          19,
			TimeSlot:      timeslot,
			StadiumID:     homeTeam.StadiumID,
			Stadium:       homeTeam.Stadium,
			City:          homeTeam.City,
			State:         homeTeam.State,
			IsDomed:       st.IsDomed,
			HomeTeamID:    int(homeTeam.ID),
			HomeTeam:      homeTeam.TeamName + " " + homeTeam.Mascot,
			HomeTeamCoach: homeTeam.NFLCoachName,
			AwayTeamID:    int(awayTeam.ID),
			AwayTeam:      awayTeam.TeamName + " " + awayTeam.Mascot,
			AwayTeamCoach: awayTeam.NFLCoachName,
			IsPlayoffGame: true,
			GameTitle:     title,
			NextGameID:    nextGameID,
			NextGameHOA:   nextHOA,
		}
	}

	wcGames := []structs.NFLGame{
		// AFC: 2v7 -> Div2 Home ("H")
		buildWCGame(
			afcSeedMap[2], afcSeedMap[7],
			fmt.Sprintf("%d AFC Wild Card - #7 %s at #2 %s", seasonYear, afcSeedMap[7].TeamAbbr, afcSeedMap[2].TeamAbbr),
			"Saturday Afternoon", "H", afcDiv2.ID,
		),
		// AFC: 3v6 -> Div2 Away ("A")
		buildWCGame(
			afcSeedMap[3], afcSeedMap[6],
			fmt.Sprintf("%d AFC Wild Card - #6 %s at #3 %s", seasonYear, afcSeedMap[6].TeamAbbr, afcSeedMap[3].TeamAbbr),
			"Saturday Evening", "A", afcDiv2.ID,
		),
		// AFC: 4v5 -> Div1 Away ("A")
		buildWCGame(
			afcSeedMap[4], afcSeedMap[5],
			fmt.Sprintf("%d AFC Wild Card - #5 %s at #4 %s", seasonYear, afcSeedMap[5].TeamAbbr, afcSeedMap[4].TeamAbbr),
			"Sunday Afternoon", "A", afcDiv1.ID,
		),
		// NFC: 2v7 -> Div2 Home ("H")
		buildWCGame(
			nfcSeedMap[2], nfcSeedMap[7],
			fmt.Sprintf("%d NFC Wild Card - #7 %s at #2 %s", seasonYear, nfcSeedMap[7].TeamAbbr, nfcSeedMap[2].TeamAbbr),
			"Saturday Evening", "H", nfcDiv2.ID,
		),
		// NFC: 3v6 -> Div2 Away ("A")
		buildWCGame(
			nfcSeedMap[3], nfcSeedMap[6],
			fmt.Sprintf("%d NFC Wild Card - #6 %s at #3 %s", seasonYear, nfcSeedMap[6].TeamAbbr, nfcSeedMap[3].TeamAbbr),
			"Sunday Afternoon", "A", nfcDiv2.ID,
		),
		// NFC: 4v5 -> Div1 Away ("A")
		buildWCGame(
			nfcSeedMap[4], nfcSeedMap[5],
			fmt.Sprintf("%d NFC Wild Card - #5 %s at #4 %s", seasonYear, nfcSeedMap[5].TeamAbbr, nfcSeedMap[4].TeamAbbr),
			"Sunday Evening", "A", nfcDiv1.ID,
		),
	}

	for i := range wcGames {
		if err := db.Create(&wcGames[i]).Error; err != nil {
			return nil, err
		}
	}

	// Fetch and return all generated games in schedule order
	allCreated := GetNFLPostseasonGamesForCurrentSeason()
	return allCreated, nil
}
