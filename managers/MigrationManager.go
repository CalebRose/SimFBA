package managers

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/models"
	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
	"gorm.io/gorm"
)

func MigrateHistoricPlayersToNFLDraftees() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	SeasonID := strconv.Itoa(ts.CollegeSeasonID)
	historicPlayers := []structs.HistoricCollegePlayer{}
	draftees := []models.NFLDraftee{}
	targetDate, err := time.Parse("2006-01-02", "2024-09-10")
	if err != nil {
		log.Panic(err)
	}
	err = db.Where("created_at > ?", targetDate).Find(&historicPlayers).Error
	if err != nil {
		log.Panic(err)
	}

	for _, p := range historicPlayers {
		if p.CreatedAt.Before(targetDate) {
			continue
		}

		grad := (structs.CollegePlayer)(p)

		draftee := models.NFLDraftee{}
		draftee.Map(grad)
		// Map New Progression value for NFL
		newProgression := util.GenerateNFLPotential(grad.Progression)
		newPotentialGrade := util.GetWeightedPotentialGrade(newProgression)
		draftee.MapProgression(newProgression, newPotentialGrade)

		if draftee.Position == "RB" {
			draftee = BoomBustDraftee(draftee, SeasonID, 31, true)
		}

		draftee.GetLetterGrades()

		/*
			Boom/Bust Function
		*/
		tier := 1
		isBoom := false
		enableBoomBust := false
		boomBustStatus := "None"
		tierRoll := util.GenerateIntFromRange(1, 10)
		diceRoll := util.GenerateIntFromRange(1, 20)

		if tierRoll > 7 && tierRoll < 10 {
			tier = 2
		} else if tierRoll > 9 {
			tier = 3
		}

		// Generate Tier
		if diceRoll == 1 {
			boomBustStatus = "Bust"
			enableBoomBust = true
			// Bust
			fmt.Println(draftee.FirstName + " " + draftee.LastName + " has BUSTED!")
			draftee.AssignBoomBustStatus(boomBustStatus)

		} else if diceRoll == 20 {
			enableBoomBust = true
			// Boom
			fmt.Println(draftee.FirstName + " " + draftee.LastName + " has BOOMED!")
			boomBustStatus = "Boom"
			isBoom = true
			draftee.AssignBoomBustStatus(boomBustStatus)
		} else {
			tier = 0
		}
		if enableBoomBust {
			for i := 0; i < tier; i++ {
				draftee = BoomBustDraftee(draftee, SeasonID, 51, isBoom)
			}
		}

		draftees = append(draftees, draftee)
	}
	repository.CreateNFLDrafteesSafely(db, draftees, 500)
}

func MigrateNFLTeamSnapsFromPreviousSeason() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID2024 := ts.CollegeSeasonID - 1
	season2024 := ts.Season - 1
	seasonID := strconv.Itoa(seasonID2024)

	nflTeams := GetAllNFLTeams()

	for _, t := range nflTeams {
		teamID := strconv.Itoa(int(t.ID))
		seasonStats := GetALLNFLTeamSeasonStatsByTeamANDSeason(teamID, seasonID)
		regularSeasonStats := structs.NFLTeamSeasonStats{}
		preSeasonStats := structs.NFLTeamSeasonStats{}
		postSeasonStats := structs.NFLTeamSeasonStats{}
		for _, s := range seasonStats {
			if s.GameType == 1 {
				preSeasonStats = s
			}
			if s.GameType == 2 {
				regularSeasonStats = s
			}
			if s.GameType == 3 {
				postSeasonStats = s
			}
		}
		games := GetNFLGamesByTeamIdAndSeasonId(teamID, seasonID)
		teamStats := GetNFLTeamStatsByTeamAndSeason(teamID, seasonID)
		gameMap := make(map[uint]structs.NFLTeamStats)
		regularSeasonStats.ResetStats()

		for _, s := range teamStats {
			gameMap[s.GameID] = s
		}

		// Apply Snaps
		for _, g := range games {
			offensiveSnaps := 0
			defensiveSnaps := 0
			specialTeamSnaps := 0
			stat := gameMap[g.ID]
			gameID := strconv.Itoa(int(g.ID))
			isHomeTeam := g.HomeTeamID == int(t.ID)
			playByPlays := GetNFLPlayByPlaysByGameID(gameID)

			for _, p := range playByPlays {
				playTypeID := p.PlayTypeID
				teamHasBall := (p.HomeHasBall && isHomeTeam) || (!p.HomeHasBall && !isHomeTeam)
				if playTypeID > 1 {
					specialTeamSnaps++
				} else if teamHasBall {
					offensiveSnaps++
				} else {
					defensiveSnaps++
				}
			}

			if g.IsPreseasonGame {
				stat.AddGameType(1)
				preSeasonStats.MapStats([]structs.NFLTeamStats{stat}, season2024, seasonID2024)
				preSeasonStats.AddTeamSnaps(uint16(offensiveSnaps), uint16(defensiveSnaps), uint16(specialTeamSnaps))
			} else if g.IsPlayoffGame || g.IsSuperBowl {
				stat.AddGameType(3)
				postSeasonStats.MapStats([]structs.NFLTeamStats{stat}, season2024, seasonID2024)
				postSeasonStats.AddTeamSnaps(uint16(offensiveSnaps), uint16(defensiveSnaps), uint16(specialTeamSnaps))
			} else {
				stat.AddGameType(2)
				regularSeasonStats.MapStats([]structs.NFLTeamStats{stat}, season2024, seasonID2024)
				regularSeasonStats.AddTeamSnaps(uint16(offensiveSnaps), uint16(defensiveSnaps), uint16(specialTeamSnaps))
			}
			// Add Team Stats to individual record
			stat.AddTeamSnaps(uint16(offensiveSnaps), uint16(defensiveSnaps), uint16(specialTeamSnaps))
			repository.SaveNFLTeamStats(stat, db)
		}
		repository.CreateNFLTeamSeasonStats(preSeasonStats, db)
		repository.CreateNFLTeamSeasonStats(postSeasonStats, db)
		repository.SaveNFLTeamSeasonStats(regularSeasonStats, db)
	}
}

func MigrateCFBTeamSnapsFromPreviousSeason() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID2026 := ts.CollegeSeasonID
	seasonID := strconv.Itoa(seasonID2026)

	collegeTeams := GetAllCollegeTeams()

	for _, t := range collegeTeams {
		teamID := strconv.Itoa(int(t.ID))
		seasonStats := GetALLCollegeTeamSeasonStatsBySeasonANDTeam(teamID, seasonID)
		if len(seasonStats) == 0 {
			continue
		}
		regularSeasonStats := structs.CollegeTeamSeasonStats{}
		preSeasonStats := structs.CollegeTeamSeasonStats{}
		postSeasonStats := structs.CollegeTeamSeasonStats{}
		for _, s := range seasonStats {
			if s.GameType == 1 {
				preSeasonStats = s
			}
			if s.GameType == 2 {
				regularSeasonStats = s
			}
			if s.GameType == 3 {
				postSeasonStats = s
			}
		}
		games := GetCollegeGamesByTeamIdAndSeasonId(teamID, seasonID, ts.CFBSpringGames)
		teamStats := GetCollegeTeamStatsBySeasonANDTeam(teamID, seasonID)
		gameMap := make(map[uint]structs.CollegeTeamStats)
		regularSeasonStats.ResetStats()

		for _, s := range teamStats {
			gameMap[uint(s.GameID)] = s
		}

		// Apply Snaps
		for _, g := range games {
			offensiveSnaps := 0
			defensiveSnaps := 0
			specialTeamSnaps := 0
			stat := gameMap[g.ID]
			gameID := strconv.Itoa(int(g.ID))
			isHomeTeam := g.HomeTeamID == int(t.ID)
			playByPlays := GetCFBPlayByPlaysByGameID(gameID)

			for _, p := range playByPlays {
				playTypeID := p.PlayTypeID
				teamHasBall := (p.HomeHasBall && isHomeTeam) || (!p.HomeHasBall && !isHomeTeam)
				if playTypeID > 1 {
					specialTeamSnaps++
				} else if teamHasBall {
					offensiveSnaps++
				} else {
					defensiveSnaps++
				}
			}

			if g.IsSpringGame {
				stat.AddGameType(1)
				preSeasonStats.MapStats([]structs.CollegeTeamStats{stat}, seasonID2026)
				preSeasonStats.AddTeamSnaps(uint16(offensiveSnaps), uint16(defensiveSnaps), uint16(specialTeamSnaps))
			} else if g.IsPlayoffGame || g.IsNationalChampionship || g.IsConferenceChampionship || g.IsBowlGame {
				stat.AddGameType(3)
				postSeasonStats.MapStats([]structs.CollegeTeamStats{stat}, seasonID2026)
				postSeasonStats.AddTeamSnaps(uint16(offensiveSnaps), uint16(defensiveSnaps), uint16(specialTeamSnaps))
			} else {
				stat.AddGameType(2)
				regularSeasonStats.MapStats([]structs.CollegeTeamStats{stat}, seasonID2026)
				regularSeasonStats.AddTeamSnaps(uint16(offensiveSnaps), uint16(defensiveSnaps), uint16(specialTeamSnaps))
			}
			// Add Team Stats to individual record
			stat.AddTeamSnaps(uint16(offensiveSnaps), uint16(defensiveSnaps), uint16(specialTeamSnaps))
			repository.SaveCFBTeamStats(stat, db)
		}
		repository.CreateCFBTeamSeasonStats(preSeasonStats, db)
		repository.CreateCFBTeamSeasonStats(postSeasonStats, db)
		repository.SaveCFBTeamSeasonStats(regularSeasonStats, db)
	}
}

func MigrateNFLPlayerStatsFromPreviousSeason() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID2024 := ts.CollegeSeasonID
	seasonID := strconv.Itoa(seasonID2024)

	nflPlayerSeasonStatMap := GetALLNFLPlayerSeasonStatMapBySeason(seasonID)
	nflPlayerStatMap := GetNFLPlayerIndividualStatMapBySeason(seasonID)
	GetNFLPlayerIndividualStatMapBySeason(seasonID)
	freeAgents := GetAllFreeAgents()
	retiredPlayers := GetRetiredSimNFLPlayers()
	nflTeams := GetAllNFLTeams()
	gameMap := GetNFLGamesMapSeasonId(seasonID)
	for _, t := range nflTeams {
		teamID := strconv.Itoa(int(t.ID))
		roster := GetNFLPlayersRecordsByTeamID(teamID)
		for _, p := range roster {
			// Apply Snaps
			migrateNFLPlayerSeasonStats(nflPlayerSeasonStatMap, p, nflPlayerStatMap, gameMap, ts, db)
		}
	}

	for _, p := range freeAgents {
		migrateNFLPlayerSeasonStats(nflPlayerSeasonStatMap, p, nflPlayerStatMap, gameMap, ts, db)
	}

	for _, p := range retiredPlayers {
		// Map to NFL player struct
		nflp := structs.NFLPlayer(p)
		migrateNFLPlayerSeasonStats(nflPlayerSeasonStatMap, nflp, nflPlayerStatMap, gameMap, ts, db)
	}
}

func migrateNFLPlayerSeasonStats(nflPlayerSeasonStatMap map[uint][]structs.NFLPlayerSeasonStats, p structs.NFLPlayer, nflPlayerStatMap map[uint][]structs.NFLPlayerStats, gameMap map[uint]structs.NFLGame, ts structs.Timestamp, db *gorm.DB) {
	seasonStats := nflPlayerSeasonStatMap[p.ID]
	playerStats := nflPlayerStatMap[p.ID]
	if len(seasonStats) == 0 || len(playerStats) == 0 {
		return
	}
	regularSeasonStats := structs.NFLPlayerSeasonStats{}
	preSeasonStats := structs.NFLPlayerSeasonStats{}
	postSeasonStats := structs.NFLPlayerSeasonStats{}
	for _, s := range seasonStats {
		if s.GameType == 1 {
			preSeasonStats = s
		}
		if s.GameType == 2 {
			regularSeasonStats = s
		}
		if s.GameType == 3 {
			postSeasonStats = s
		}
	}
	regularSeasonStats.ResetStats()
	preSeasonStats.ResetStats()
	postSeasonStats.ResetStats()

	for _, stat := range playerStats {
		g := gameMap[uint(stat.GameID)]
		if g.IsPreseasonGame {
			preSeasonStats.MapStats([]structs.NFLPlayerStats{stat}, ts)
		} else if g.IsPlayoffGame || g.IsSuperBowl || g.Week > 18 {
			postSeasonStats.MapStats([]structs.NFLPlayerStats{stat}, ts)
		} else {
			regularSeasonStats.MapStats([]structs.NFLPlayerStats{stat}, ts)
		}
	}
	if preSeasonStats.SeasonID > 0 {
		repository.SaveNFLPlayerSeasonStats(preSeasonStats, db)
	}

	if postSeasonStats.SeasonID > 0 {
		repository.SaveNFLPlayerSeasonStats(postSeasonStats, db)
	}

	if regularSeasonStats.SeasonID > 0 {
		repository.SaveNFLPlayerSeasonStats(regularSeasonStats, db)
	}
}

func MigrateCFBPlayerStatsFromPreviousSeason() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID2026 := ts.CollegeSeasonID
	seasonID := strconv.Itoa(seasonID2026)
	cfbPlayerSeasonStatMap := GetALLCFBPlayerSeasonStatMapBySeason(seasonID)
	cfbPlayerStatMap := GetCFBPlayerIndividualStatMapBySeason(seasonID)
	unsignedPlayers := GetAllUnsignedPlayers()
	historicCollegePlayers := GetAllHistoricCollegePlayers()
	gameMap := GetCFBGamesMapSeasonId(seasonID)

	collegeTeams := GetAllCollegeTeams()

	for _, t := range collegeTeams {
		fmt.Println("Iterating through " + t.TeamName)
		teamID := strconv.Itoa(int(t.ID))
		roster := GetAllCollegePlayersByTeamId(teamID)
		for _, p := range roster {
			migrateCFBPlayerSeasonStats(cfbPlayerSeasonStatMap, p, cfbPlayerStatMap, gameMap, db)
		}
	}

	fmt.Println("Iterating through unsigned players...")
	for _, p := range unsignedPlayers {
		player := structs.CollegePlayer(p)
		migrateCFBPlayerSeasonStats(cfbPlayerSeasonStatMap, player, cfbPlayerStatMap, gameMap, db)
	}

	fmt.Println("Iterating through historic graduated players...")
	for _, p := range historicCollegePlayers {
		player := structs.CollegePlayer(p)
		migrateCFBPlayerSeasonStats(cfbPlayerSeasonStatMap, player, cfbPlayerStatMap, gameMap, db)
	}
}

func migrateCFBPlayerSeasonStats(cfbPlayerSeasonStatMap map[uint][]structs.CollegePlayerSeasonStats, p structs.CollegePlayer, cfbPlayerStatMap map[uint][]structs.CollegePlayerStats, gameMap map[uint]structs.CollegeGame, db *gorm.DB) {
	seasonStats := cfbPlayerSeasonStatMap[p.ID]
	playerStats := cfbPlayerStatMap[p.ID]
	if len(seasonStats) == 0 || len(playerStats) == 0 {
		return
	}
	regularSeasonStats := structs.CollegePlayerSeasonStats{}
	preSeasonStats := structs.CollegePlayerSeasonStats{}
	postSeasonStats := structs.CollegePlayerSeasonStats{}
	for _, s := range seasonStats {
		if s.GameType == 1 {
			preSeasonStats = s
		}
		if s.GameType == 2 {
			regularSeasonStats = s
		}
		if s.GameType == 3 {
			postSeasonStats = s
		}
	}

	regularSeasonStats.ResetStats()
	preSeasonStats.ResetStats()
	postSeasonStats.ResetStats()

	for _, stat := range playerStats {
		g := gameMap[uint(stat.GameID)]
		if g.IsSpringGame {
			preSeasonStats.MapStats([]structs.CollegePlayerStats{stat})
		} else if g.IsPlayoffGame || g.IsNationalChampionship || g.IsBowlGame {
			postSeasonStats.MapStats([]structs.CollegePlayerStats{stat})
		} else {
			regularSeasonStats.MapStats([]structs.CollegePlayerStats{stat})
		}
	}
	if preSeasonStats.GamesPlayed > 0 {
		repository.SaveCollegePlayerSeasonStats(preSeasonStats, db)
	}
	if postSeasonStats.GamesPlayed > 0 {
		repository.SaveCollegePlayerSeasonStats(postSeasonStats, db)
	}
	if regularSeasonStats.ID > 0 {
		repository.SaveCollegePlayerSeasonStats(regularSeasonStats, db)
	}
}

func FixSpendingCount() {
	db := dbprovider.GetInstance().GetDB()

	recruits := GetAllRecruitRecords()
	unsignedIDs := []string{}
	for _, r := range recruits {
		if r.IsSigned {
			continue
		}
		id := strconv.Itoa(int(r.ID))
		unsignedIDs = append(unsignedIDs, id)
	}

	allProfiles := []structs.RecruitPlayerProfile{}

	db.Where("recruit_id in (?)", unsignedIDs).Find(&allProfiles)

	profileIDs := []string{}
	profileMap := make(map[uint][]structs.RecruitPlayerProfile)

	for _, p := range allProfiles {
		if p.RemovedFromBoard {
			continue
		}
		id := strconv.Itoa(int(p.ProfileID))
		profileIDs = append(profileIDs, id)
		if len(profileMap[uint(p.RecruitID)]) == 0 {
			profileMap[uint(p.RecruitID)] = []structs.RecruitPlayerProfile{p}
		} else {
			profileMap[uint(p.RecruitID)] = append(profileMap[uint(p.RecruitID)], p)
		}
	}

	allocations := []structs.RecruitPointAllocation{}

	db.Where("team_profile_id in (?)", profileIDs).Find(&allocations)

	allocationMap := make(map[uint][]structs.RecruitPointAllocation)

	for _, a := range allocations {
		if len(allocationMap[uint(a.RecruitProfileID)]) == 0 {
			allocationMap[uint(a.RecruitProfileID)] = []structs.RecruitPointAllocation{a}
		} else {
			allocationMap[uint(a.RecruitProfileID)] = append(allocationMap[uint(a.RecruitProfileID)], a)
		}
	}

	for _, r := range recruits {
		if r.IsSigned {
			continue
		}

		profiles := profileMap[r.ID]

		for _, p := range profiles {
			spendingCount := 0
			profileAllocations := allocationMap[uint(p.ID)]
			weekPointMap := make(map[int]float64)

			for _, a := range profileAllocations {
				weekPointMap[a.WeekID] += a.Points
			}

			weeks := []int{}
			for week := range weekPointMap {
				weeks = append(weeks, week)
			}
			sort.Ints(weeks)
			prevWeek := -1
			for _, week := range weeks {
				points := weekPointMap[week]

				if prevWeek != -1 && week-prevWeek > 1 {
					// Missed one or more weeks – reset the streak
					spendingCount = 0
				}

				if points > 0 {
					spendingCount++
				}

				prevWeek = week
			}

			p.UpdateSpendingCount(spendingCount)
			repository.SaveRecruitProfile(p, db)
		}
	}
}
