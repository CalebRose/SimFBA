package managers

import (
	"encoding/csv"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
)

func ExportAllRostersToCSV(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment;filename=2023_Rosters.csv")
	w.Header().Set("Transfer-Encoding", "chunked")
	// Initialize writer
	writer := csv.NewWriter(w)

	HeaderRow := []string{
		"Team", "Player ID", "First Name", "Last Name", "Position",
		"Archetype", "Position Two", "Archetype Two", "Year", "Age", "Stars",
		"High School", "City", "State", "Height",
		"Weight", "Overall", "Speed",
		"Football IQ", "Agility", "Carrying",
		"Catching", "Route Running", "Zone Coverage", "Man Coverage",
		"Strength", "Tackle", "Pass Block", "Run Block",
		"Pass Rush", "Run Defense", "Throw Power", "Throw Accuracy",
		"Kick Power", "Kick Accuracy", "Punt Power", "Punt Accuracy",
		"Stamina", "Injury", "Potential Grade", "Redshirt Status",
	}

	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	players := GetAllCollegePlayers()

	for _, player := range players {
		csvModel := structs.MapPlayerToCSVModel(player)
		idStr := strconv.Itoa(int(player.ID))
		playerRow := []string{
			player.TeamAbbr, idStr, csvModel.FirstName, csvModel.LastName, csvModel.Position,
			csvModel.Archetype, csvModel.PositionTwo, csvModel.ArchetypeTwo, csvModel.Year, strconv.Itoa(player.Age), strconv.Itoa(player.Stars),
			player.HighSchool, player.City, player.State, strconv.Itoa(player.Height),
			strconv.Itoa(player.Weight), csvModel.OverallGrade, csvModel.SpeedGrade,
			csvModel.FootballIQGrade, csvModel.AgilityGrade, csvModel.CarryingGrade,
			csvModel.CatchingGrade, csvModel.RouteRunningGrade, csvModel.ZoneCoverageGrade, csvModel.ManCoverageGrade,
			csvModel.StrengthGrade, csvModel.TackleGrade, csvModel.PassBlockGrade, csvModel.RunBlockGrade,
			csvModel.PassRushGrade, csvModel.RunDefenseGrade, csvModel.ThrowPowerGrade, csvModel.ThrowAccuracyGrade,
			csvModel.KickPowerGrade, csvModel.KickAccuracyGrade, csvModel.PuntPowerGrade, csvModel.PuntAccuracyGrade,
			csvModel.StaminaGrade, csvModel.InjuryGrade, csvModel.PotentialGrade, csvModel.RedshirtStatus,
		}

		err = writer.Write(playerRow)
		if err != nil {
			log.Fatal("Cannot write player row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
}

func ExportTeamToCSV(TeamID string, w http.ResponseWriter) {
	// Get Team Data
	team := GetTeamByTeamID(TeamID)
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment;filename="+team.TeamName+".csv")
	w.Header().Set("Transfer-Encoding", "chunked")
	// Initialize writer
	writer := csv.NewWriter(w)

	// Get Players
	players := GetAllCollegePlayersByTeamId(TeamID)

	HeaderRow := []string{
		"Team", "Player ID", "First Name", "Last Name", "Position",
		"Archetype", "Position Two", "Archetype Two", "Year", "Age", "Stars",
		"High School", "City", "State", "Height",
		"Weight", "Overall", "Speed",
		"Football IQ", "Agility", "Carrying",
		"Catching", "Route Running", "Zone Coverage", "Man Coverage",
		"Strength", "Tackle", "Pass Block", "Run Block",
		"Pass Rush", "Run Defense", "Throw Power", "Throw Accuracy",
		"Kick Power", "Kick Accuracy", "Punt Power", "Punt Accuracy",
		"Stamina", "Injury", "Potential Grade", "Redshirt Status",
	}

	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	for _, player := range players {
		csvModel := structs.MapPlayerToCSVModel(player)
		idStr := strconv.Itoa(int(player.ID))
		playerRow := []string{
			team.TeamName, idStr, csvModel.FirstName, csvModel.LastName, csvModel.Position,
			csvModel.Archetype, csvModel.PositionTwo, csvModel.ArchetypeTwo, csvModel.Year, strconv.Itoa(player.Age), strconv.Itoa(player.Stars),
			player.HighSchool, player.City, player.State, strconv.Itoa(player.Height),
			strconv.Itoa(player.Weight), csvModel.OverallGrade, csvModel.SpeedGrade,
			csvModel.FootballIQGrade, csvModel.AgilityGrade, csvModel.CarryingGrade,
			csvModel.CatchingGrade, csvModel.RouteRunningGrade, csvModel.ZoneCoverageGrade, csvModel.ManCoverageGrade,
			csvModel.StrengthGrade, csvModel.TackleGrade, csvModel.PassBlockGrade, csvModel.RunBlockGrade,
			csvModel.PassRushGrade, csvModel.RunDefenseGrade, csvModel.ThrowPowerGrade, csvModel.ThrowAccuracyGrade,
			csvModel.KickPowerGrade, csvModel.KickAccuracyGrade, csvModel.PuntPowerGrade, csvModel.PuntAccuracyGrade,
			csvModel.StaminaGrade, csvModel.InjuryGrade, csvModel.PotentialGrade, csvModel.RedshirtStatus,
		}

		err = writer.Write(playerRow)
		if err != nil {
			log.Fatal("Cannot write player row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
}

func ExportNFLTeamToCSV(TeamID string, w http.ResponseWriter) {
	// Get Team Data
	team := GetNFLTeamByTeamID(TeamID)
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment;filename="+team.TeamName+".csv")
	w.Header().Set("Transfer-Encoding", "chunked")
	// Initialize writer
	writer := csv.NewWriter(w)

	// Get Players
	players := GetNFLPlayersWithContractsByTeamID(TeamID)

	HeaderRow := []string{
		"Team", "ID", "First Name", "Last Name", "Position",
		"Archetype", "Position Two", "Archetype Two", "Year", "Age",
		"High School", "Hometown", "State", "Height",
		"Weight", "Overall", "Speed",
		"Football IQ", "Agility", "Carrying",
		"Catching", "Route Running", "Zone Coverage", "Man Coverage",
		"Strength", "Tackle", "Pass Block", "Run Block",
		"Pass Rush", "Run Defense", "Throw Power", "Throw Accuracy",
		"Kick Power", "Kick Accuracy", "Punt Power", "Punt Accuracy",
		"Stamina", "Injury", "Potential Grade",
	}

	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	for _, player := range players {
		csvModel := structs.MapNFLPlayerToCSVModel(player)
		playerRow := []string{
			team.TeamName, strconv.Itoa(int(player.ID)), csvModel.FirstName, csvModel.LastName, csvModel.Position,
			csvModel.Archetype, csvModel.PositionTwo, csvModel.ArchetypeTwo, csvModel.Year, strconv.Itoa(player.Age),
			player.HighSchool, player.Hometown, player.State, strconv.Itoa(player.Height),
			strconv.Itoa(player.Weight), csvModel.OverallGrade, csvModel.SpeedGrade,
			csvModel.FootballIQGrade, csvModel.AgilityGrade, csvModel.CarryingGrade,
			csvModel.CatchingGrade, csvModel.RouteRunningGrade, csvModel.ZoneCoverageGrade, csvModel.ManCoverageGrade,
			csvModel.StrengthGrade, csvModel.TackleGrade, csvModel.PassBlockGrade, csvModel.RunBlockGrade,
			csvModel.PassRushGrade, csvModel.RunDefenseGrade, csvModel.ThrowPowerGrade, csvModel.ThrowAccuracyGrade,
			csvModel.KickPowerGrade, csvModel.KickAccuracyGrade, csvModel.PuntPowerGrade, csvModel.PuntAccuracyGrade,
			csvModel.StaminaGrade, csvModel.InjuryGrade, csvModel.PotentialGrade,
		}

		err = writer.Write(playerRow)
		if err != nil {
			log.Fatal("Cannot write player row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
}

func ExportAllNFLTeamsToCSV(w http.ResponseWriter) {
	// Get Team Data
	nflPlayers := GetAllNFLPlayers()
	nflPlayerMap := MakeNFLPlayerMapByTeamID(nflPlayers, true)
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment;filename=all_nfl_players_on_roster.csv")
	w.Header().Set("Transfer-Encoding", "chunked")
	// Initialize writer
	writer := csv.NewWriter(w)

	// Get Players
	nflTeams := GetAllNFLTeams()

	HeaderRow := []string{
		"Team", "ID", "First Name", "Last Name", "Position",
		"Archetype", "Position Two", "Archetype Two", "Year", "Age",
		"High School", "Hometown", "State", "Height",
		"Weight", "Overall", "Speed",
		"Football IQ", "Agility", "Carrying",
		"Catching", "Route Running", "Zone Coverage", "Man Coverage",
		"Strength", "Tackle", "Pass Block", "Run Block",
		"Pass Rush", "Run Defense", "Throw Power", "Throw Accuracy",
		"Kick Power", "Kick Accuracy", "Punt Power", "Punt Accuracy",
		"Stamina", "Injury", "Potential Grade",
	}

	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	for _, team := range nflTeams {
		players := nflPlayerMap[team.ID]

		for _, player := range players {
			csvModel := structs.MapNFLPlayerToCSVModel(player)
			playerRow := []string{
				team.TeamName, strconv.Itoa(int(player.ID)), csvModel.FirstName, csvModel.LastName, csvModel.Position,
				csvModel.Archetype, csvModel.PositionTwo, csvModel.ArchetypeTwo, csvModel.Year, strconv.Itoa(player.Age),
				player.HighSchool, player.Hometown, player.State, strconv.Itoa(player.Height),
				strconv.Itoa(player.Weight), csvModel.OverallGrade, csvModel.SpeedGrade,
				csvModel.FootballIQGrade, csvModel.AgilityGrade, csvModel.CarryingGrade,
				csvModel.CatchingGrade, csvModel.RouteRunningGrade, csvModel.ZoneCoverageGrade, csvModel.ManCoverageGrade,
				csvModel.StrengthGrade, csvModel.TackleGrade, csvModel.PassBlockGrade, csvModel.RunBlockGrade,
				csvModel.PassRushGrade, csvModel.RunDefenseGrade, csvModel.ThrowPowerGrade, csvModel.ThrowAccuracyGrade,
				csvModel.KickPowerGrade, csvModel.KickAccuracyGrade, csvModel.PuntPowerGrade, csvModel.PuntAccuracyGrade,
				csvModel.StaminaGrade, csvModel.InjuryGrade, csvModel.PotentialGrade,
			}

			err = writer.Write(playerRow)
			if err != nil {
				log.Fatal("Cannot write player row to CSV", err)
			}

			writer.Flush()
			err = writer.Error()
			if err != nil {
				log.Fatal("Error while writing to file ::", err)
			}
		}
	}

}

func ExportCrootsToCSV(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment;filename=2022SimNFLDraftClass.csv")
	w.Header().Set("Transfer-Encoding", "chunked")
	// Initialize writer
	writer := csv.NewWriter(w)

	// Get NFL Draft Class
	croots := GetAllRecruits()

	HeaderRow := []string{
		"RecruitID", "First Name", "Last Name", "Position",
		"Archetype", "Stars", "College",
		"High School", "City", "State", "Height",
		"Weight", "Overall", "Potential Grade", "Affinity One", "Affinity Two", "Personality",
		"Recruiting Bias", "Academic Bias", "Work Ethic", "LeadingTeams",
	}

	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	for _, croot := range croots {
		var leadingAbbr []string

		for _, lt := range croot.LeadingTeams {
			if lt.Odds > 0 {
				leadingAbbr = append(leadingAbbr, lt.TeamAbbr)
			}
		}

		crootRow := []string{
			strconv.Itoa(int(croot.ID)), croot.FirstName, croot.LastName, croot.Position,
			croot.Archetype, strconv.Itoa(croot.Stars), croot.College,
			croot.HighSchool, croot.City, croot.State, strconv.Itoa(croot.Height),
			strconv.Itoa(croot.Weight), croot.OverallGrade, croot.PotentialGrade,
			croot.AffinityOne, croot.AffinityTwo, croot.Personality, croot.RecruitingBias,
			croot.AcademicBias, croot.WorkEthic, strings.Join(leadingAbbr, ", "),
		}

		err = writer.Write(crootRow)
		if err != nil {
			log.Fatal("Cannot write croot row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
}

func ExportDrafteesToCSV(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment;filename=2024SimNFLDraftClass.csv")
	w.Header().Set("Transfer-Encoding", "chunked")
	// Initialize writer
	writer := csv.NewWriter(w)

	// Get NFL Draft Class
	draftees := GetAllNFLDraftees()

	HeaderRow := []string{
		"PlayerID", "First Name", "Last Name", "Position",
		"Archetype", "Position Two", "Archetype Two", "Age", "Stars", "College",
		"High School", "City", "State", "Height",
		"Weight", "Overall", "Speed",
		"Football IQ", "Agility", "Carrying",
		"Catching", "Route Running", "Zone Coverage", "Man Coverage",
		"Strength", "Tackle", "Pass Block", "Run Block",
		"Pass Rush", "Run Defense", "Throw Power", "Throw Accuracy",
		"Kick Power", "Kick Accuracy", "Punt Power", "Punt Accuracy",
		"Stamina", "Injury", "Potential Grade",
	}

	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	for _, player := range draftees {
		playerRow := []string{
			strconv.Itoa(int(player.ID)), player.FirstName, player.LastName, player.Position,
			player.Archetype, player.PositionTwo, player.ArchetypeTwo, strconv.Itoa(player.Age), strconv.Itoa(player.Stars), player.College,
			player.HighSchool, player.City, player.State, strconv.Itoa(player.Height),
			strconv.Itoa(player.Weight), player.OverallGrade, player.SpeedGrade,
			player.FootballIQGrade, player.AgilityGrade, player.CarryingGrade,
			player.CatchingGrade, player.RouteRunningGrade, player.ZoneCoverageGrade, player.ManCoverageGrade,
			player.StrengthGrade, player.TackleGrade, player.PassBlockGrade, player.RunBlockGrade,
			player.PassRushGrade, player.RunDefenseGrade, player.ThrowPowerGrade, player.ThrowAccuracyGrade,
			player.KickPowerGrade, player.KickAccuracyGrade, player.PuntPowerGrade, player.PuntAccuracyGrade,
			player.StaminaGrade, player.InjuryGrade, player.PotentialGrade,
		}

		err = writer.Write(playerRow)
		if err != nil {
			log.Fatal("Cannot write player row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
}

func ExportPlayerStatsToCSV(cp []structs.CollegePlayerResponse, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment;filename=special_season_stats.csv")
	w.Header().Set("Transfer-Encoding", "chunked")
	// Initialize writer
	writer := csv.NewWriter(w)

	// Get Players

	HeaderRow := []string{
		"ID", "First Name", "Last Name", "Position", "Position Two",
		"Archetype", "Archetype Two", "Year", "Is Redshirt?", "Team", "Conference", "Age", "Stars",
		"Passing Yards", "Pass Attempts", "Pass Completions", "Passing Avg",
		"Passing TDs", "Interceptions", "Longest Pass", "QB Sacks",
		"QB Rating", "Rush Attempts", "Rushing Yards", "Rushing Avg",
		"Rushing TDs", "Fumbles", "Longest Rush", "Targets",
		"Catches", "Receiving Yards", "Receiving Avg", "Receiving TDs",
		"Longest Reception", "Solo Tackles", "Assisted Tackles", "Tackles For Loss",
		"Defensive Sacks", "Forced Fumbles", "Pass Deflections", "Interceptions Caught",
		"Safeties", "Defensive TDs", "FG Made", "FG Attempts",
		"Longest FG", "Extra Points Made", "Extra Point Attempts", "Kickoff TBs",
		"Punts", "Punt Touchbacks", "Punts Inside 20", "Kick Returns",
		"Kick Return TDs", "Kick Return Yards", "Punt Returns", "Punt Return TDs",
		"Punt Return Yards", "ST Solo Tackles", "ST Assisted Tackles", "Punts Blocked",
		"FG Blocked", "Snaps", "Pancakes", "Ready for Mars?",
	}

	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	for _, p := range cp {
		Year, RedshirtStatus := util.GetYearAndRedshirtStatus(p.Year, p.IsRedshirt)
		seasonStats := p.SeasonStats

		pr := []string{strconv.Itoa(int(p.ID)), p.FirstName, p.LastName, p.Position, p.PositionTwo,
			p.Archetype, p.ArchetypeTwo, Year, RedshirtStatus, p.TeamAbbr, p.Conference, strconv.Itoa(p.Age), strconv.Itoa(p.Stars),
			strconv.Itoa(seasonStats.PassingYards), strconv.Itoa(seasonStats.PassAttempts), strconv.Itoa(seasonStats.PassCompletions), strconv.Itoa(int(seasonStats.PassingAvg)),
			strconv.Itoa(seasonStats.PassingTDs), strconv.Itoa(seasonStats.Interceptions), strconv.Itoa(seasonStats.LongestPass), strconv.Itoa(seasonStats.Sacks),
			strconv.Itoa(int(seasonStats.QBRating)), strconv.Itoa(seasonStats.RushAttempts), strconv.Itoa(seasonStats.RushingYards), strconv.Itoa(int(seasonStats.RushingAvg)),
			strconv.Itoa(seasonStats.RushingTDs), strconv.Itoa(seasonStats.Fumbles), strconv.Itoa(seasonStats.LongestRush), strconv.Itoa(seasonStats.Targets),
			strconv.Itoa(seasonStats.Catches), strconv.Itoa(seasonStats.ReceivingYards), strconv.Itoa(int(seasonStats.ReceivingAvg)), strconv.Itoa(seasonStats.ReceivingTDs),
			strconv.Itoa(seasonStats.LongestReception), strconv.Itoa(int(seasonStats.SoloTackles)), strconv.Itoa(int(seasonStats.AssistedTackles)), strconv.Itoa(int(seasonStats.TacklesForLoss)),
			strconv.Itoa(int(seasonStats.SacksMade)), strconv.Itoa(seasonStats.ForcedFumbles), strconv.Itoa(seasonStats.PassDeflections), strconv.Itoa(seasonStats.InterceptionsCaught),
			strconv.Itoa(seasonStats.Safeties), strconv.Itoa(seasonStats.DefensiveTDs), strconv.Itoa(seasonStats.FGMade), strconv.Itoa(seasonStats.FGAttempts),
			strconv.Itoa(seasonStats.LongestFG), strconv.Itoa(seasonStats.ExtraPointsMade), strconv.Itoa(seasonStats.ExtraPointsAttempted), strconv.Itoa(seasonStats.KickoffTouchbacks),
			strconv.Itoa(seasonStats.Punts), strconv.Itoa(seasonStats.PuntTouchbacks), strconv.Itoa(seasonStats.PuntsInside20), strconv.Itoa(seasonStats.KickReturns),
			strconv.Itoa(seasonStats.KickReturnTDs), strconv.Itoa(seasonStats.KickReturnYards), strconv.Itoa(seasonStats.PuntReturns), strconv.Itoa(seasonStats.PuntReturnTDs),
			strconv.Itoa(seasonStats.PuntReturnYards), strconv.Itoa(int(seasonStats.STSoloTackles)), strconv.Itoa(int(seasonStats.STAssistedTackles)), strconv.Itoa(seasonStats.PuntsBlocked),
			strconv.Itoa(seasonStats.FGBlocked), strconv.Itoa(seasonStats.Snaps), strconv.Itoa(seasonStats.Pancakes), "No.",
		}
		err = writer.Write(pr)
		if err != nil {
			log.Fatal("Cannot write player row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
}

func ExportTransferPlayersToCSV(transfers []structs.CollegePlayer, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment;filename=transferStats.csv")
	w.Header().Set("Transfer-Encoding", "chunked")
	// Initialize writer
	writer := csv.NewWriter(w)
	HeaderRow := []string{
		"Team", "ID", "First Name", "Last Name", "Stars",
		"Archetype", "Position", "Year", "Age", "Redshirt Status",
		"Overall", "Transfer Weight", "Dice Roll",
	}

	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	for _, player := range transfers {
		csvModel := structs.MapPlayerToCSVModel(player)
		playerRow := []string{
			player.TeamAbbr, strconv.Itoa(int(player.ID)), csvModel.FirstName, csvModel.LastName, strconv.Itoa(player.Stars),
			csvModel.Archetype, csvModel.Position,
			csvModel.Year, strconv.Itoa(player.Age), csvModel.RedshirtStatus,
			csvModel.OverallGrade,
		}

		err = writer.Write(playerRow)
		if err != nil {
			log.Fatal("Cannot write player row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
}

func ExportNFLPlayByPlayToCSV(gameID string, w http.ResponseWriter) {
	game := GetNFLGameByGameID(gameID)
	htID := strconv.Itoa(game.HomeTeamID)
	atID := strconv.Itoa(game.AwayTeamID)
	homePlayerStats := GetAllNFLPlayerStatsByGame(gameID, htID)
	awayPlayerStats := GetAllNFLPlayerStatsByGame(gameID, atID)
	homePlayers := GetAllNFLPlayersWithGameStatsByTeamID(gameID, homePlayerStats)
	awayPlayers := GetAllNFLPlayersWithGameStatsByTeamID(gameID, awayPlayerStats)
	participantMap := getGameParticipantMap(homePlayers, awayPlayers)

	playByPlays := GetNFLPlayByPlaysByGameID(gameID)
	// Generate the Play By Play Response
	playbyPlayResponseList := GenerateNFLPlayByPlayResponse(playByPlays, participantMap, false, game.HomeTeam, game.AwayTeam)

	// Begin Writing
	fileName := gameID + "_" + game.HomeTeam + "_vs_" + game.AwayTeam + "_play_by_play"
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment;filename="+fileName+".csv")
	w.Header().Set("Transfer-Encoding", "chunked")
	// Initialize writer
	writer := csv.NewWriter(w)
	HeaderRow := []string{
		"Play", game.HomeTeam + " Score", game.AwayTeam + " Score", "Quarter", "Time Remaining",
		"Possession", "Down", "Distance", "Line of Scrimmage", "Type of Play",
		"Offensive Formation", "Offensive Play", "Offensive PoA", "Defensive Formation",
		"Defensive Tendency", "# of Blitzers", "LB Coverage", "CB Coverage", "S Coverage",
		"QB Player ID", "Ballcarrier ID", "Tackler1 ID", "Tackler2 ID", "Yards Gained",
		"Result",
		"QB Id", "Back1 Id", "Back2 Id", "Back3 Id", "Slot1 Id", "Slot2 Id", "Le Id", "Re Id", "Lt Id", "LG id", "C id", "Rg Id", "Rt Id",
		"Rde id", "Rdt Id", "Nt Id", "Ldt Id", "Lde Id", "Rolb Id", "Rilb Id", "Mlb Id", "Lilb Id", "Lolb Id", "Rcb Id", "DB1 Id", "DB2 Id", "DB3 ID", "Fs ID", "SS Id", "Lcb Id",
		"Blitzer1 Id", "Blitzer2 Id", "Blitzer3 Id",
	}
	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	for _, play := range playbyPlayResponseList {
		num := strconv.Itoa(int(play.PlayNumber))
		hcs := strconv.Itoa(int(play.HomeTeamScore))
		acs := strconv.Itoa(int(play.AwayTeamScore))
		qt := strconv.Itoa(int(play.Quarter))
		tr := play.TimeRemaining
		down := strconv.Itoa(int(play.Down))
		dist := strconv.Itoa(int(play.Distance))
		qbID := strconv.Itoa(int(play.QBPlayerID))
		bcID := strconv.Itoa(int(play.BallCarrierID))
		t1ID := strconv.Itoa(int(play.Tackler1ID))
		t2ID := strconv.Itoa(int(play.Tackler2ID))
		yards := strconv.Itoa(int(play.ResultYards))
		blitzNumber := strconv.Itoa(int(play.BlitzNumber))

		row := []string{
			num, hcs, acs, qt, tr, play.Possession, down, dist, play.LineOfScrimmage,
			play.PlayType, play.OffensiveFormation, play.PlayName, play.PointOfAttack, play.DefensiveFormation,
			play.DefensiveTendency, blitzNumber, play.LBCoverage, play.CBCoverage, play.SCoverage,
			qbID, bcID, t1ID, t2ID, yards,
			play.Result,
			strconv.Itoa(int(play.Qb)), strconv.Itoa(int(play.Back1)), strconv.Itoa(int(play.Back2)), strconv.Itoa(int(play.Back3)),
			strconv.Itoa(int(play.Slot1)), strconv.Itoa(int(play.Slot2)), strconv.Itoa(int(play.Le)), strconv.Itoa(int(play.Re)),
			strconv.Itoa(int(play.Lt)), strconv.Itoa(int(play.Lg)), strconv.Itoa(int(play.C)), strconv.Itoa(int(play.Rg)), strconv.Itoa(int(play.Rt)),
			strconv.Itoa(int(play.Rde)), strconv.Itoa(int(play.Rdt)), strconv.Itoa(int(play.Nt)), strconv.Itoa(int(play.Ldt)), strconv.Itoa(int(play.Lde)),
			strconv.Itoa(int(play.Rolb)), strconv.Itoa(int(play.Rilb)), strconv.Itoa(int(play.Mlb)), strconv.Itoa(int(play.Lilb)), strconv.Itoa(int(play.Lolb)),
			strconv.Itoa(int(play.Rcb)),
			strconv.Itoa(int(play.Extradb1)), strconv.Itoa(int(play.Extradb2)), strconv.Itoa(int(play.Extradb3)),
			strconv.Itoa(int(play.Fs)), strconv.Itoa(int(play.Ss)), strconv.Itoa(int(play.Fcb)),
			strconv.Itoa(int(play.Blitzer1)), strconv.Itoa(int(play.Blitzer2)), strconv.Itoa(int(play.Blitzer3)),
		}

		err = writer.Write(row)
		if err != nil {
			log.Fatal("Cannot write player row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
}

func ExportCFBPlayByPlayToCSV(gameID string, w http.ResponseWriter) {
	game := GetCollegeGameByGameID(gameID)
	htID := strconv.Itoa(game.HomeTeamID)
	atID := strconv.Itoa(game.AwayTeamID)

	homeStats := GetAllCollegePlayerStatsByGame(gameID, htID)
	awayStats := GetAllCollegePlayerStatsByGame(gameID, atID)

	homePlayers := GetAllCollegePlayersWithGameStatsByTeamID(gameID, homeStats)
	awayPlayers := GetAllCollegePlayersWithGameStatsByTeamID(gameID, awayStats)
	participantMap := getGameParticipantMap(homePlayers, awayPlayers)

	playByPlays := GetCFBPlayByPlaysByGameID(gameID)
	// Generate the Play By Play Response
	playbyPlayResponseList := GenerateCFBPlayByPlayResponse(playByPlays, participantMap, false, game.HomeTeam, game.AwayTeam)

	// Begin Writing
	fileName := gameID + "_" + game.HomeTeam + "_vs_" + game.AwayTeam + "_play_by_play"
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment;filename="+fileName+".csv")
	w.Header().Set("Transfer-Encoding", "chunked")
	// Initialize writer
	writer := csv.NewWriter(w)
	HeaderRow := []string{
		"Play", game.HomeTeam + " Score", game.AwayTeam + " Score", "Quarter", "Time Remaining",
		"Possession", "Down", "Distance", "Line of Scrimmage", "Type of Play",
		"Offensive Formation", "Offensive Play", "Offensive PoA", "Defensive Formation",
		"Defensive Tendency", "# of Blitzers", "LB Coverage", "CB Coverage", "S Coverage",
		"QB Player ID", "Ballcarrier ID", "Tackler1 ID", "Tackler2 ID", "Yards Gained",
		"Result",
		"QB Id", "Back1 Id", "Back2 Id", "Back3 Id", "Slot1 Id", "Slot2 Id", "Le Id", "Re Id", "Lt Id", "LG id", "C id", "Rg Id", "Rt Id",
		"Rde id", "Rdt Id", "Nt Id", "Ldt Id", "Lde Id", "Rolb Id", "Rilb Id", "Mlb Id", "Lilb Id", "Lolb Id", "Rcb Id", "DB1 Id", "DB2 Id", "DB3 ID", "Fs ID", "SS Id", "Lcb Id",
		"Blitzer1 Id", "Blitzer2 Id", "Blitzer3 Id",
	}
	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	for _, play := range playbyPlayResponseList {
		num := strconv.Itoa(int(play.PlayNumber))
		hcs := strconv.Itoa(int(play.HomeTeamScore))
		acs := strconv.Itoa(int(play.AwayTeamScore))
		qt := strconv.Itoa(int(play.Quarter))
		tr := play.TimeRemaining
		down := strconv.Itoa(int(play.Down))
		dist := strconv.Itoa(int(play.Distance))
		qbID := strconv.Itoa(int(play.QBPlayerID))
		bcID := strconv.Itoa(int(play.BallCarrierID))
		t1ID := strconv.Itoa(int(play.Tackler1ID))
		t2ID := strconv.Itoa(int(play.Tackler2ID))
		yards := strconv.Itoa(int(play.ResultYards))
		blitzNumber := strconv.Itoa(int(play.BlitzNumber))

		row := []string{
			num, hcs, acs, qt, tr, play.Possession, down, dist, play.LineOfScrimmage,
			play.PlayType, play.OffensiveFormation, play.PlayName, play.PointOfAttack, play.DefensiveFormation,
			play.DefensiveTendency, blitzNumber, play.LBCoverage, play.CBCoverage, play.SCoverage,
			qbID, bcID, t1ID, t2ID, yards,
			play.Result,
			strconv.Itoa(int(play.Qb)), strconv.Itoa(int(play.Back1)), strconv.Itoa(int(play.Back2)), strconv.Itoa(int(play.Back3)),
			strconv.Itoa(int(play.Slot1)), strconv.Itoa(int(play.Slot2)), strconv.Itoa(int(play.Le)), strconv.Itoa(int(play.Re)),
			strconv.Itoa(int(play.Lt)), strconv.Itoa(int(play.Lg)), strconv.Itoa(int(play.C)), strconv.Itoa(int(play.Rg)), strconv.Itoa(int(play.Rt)),
			strconv.Itoa(int(play.Rde)), strconv.Itoa(int(play.Rdt)), strconv.Itoa(int(play.Nt)), strconv.Itoa(int(play.Ldt)), strconv.Itoa(int(play.Lde)),
			strconv.Itoa(int(play.Rolb)), strconv.Itoa(int(play.Rilb)), strconv.Itoa(int(play.Mlb)), strconv.Itoa(int(play.Lilb)), strconv.Itoa(int(play.Lolb)),
			strconv.Itoa(int(play.Rcb)),
			strconv.Itoa(int(play.Extradb1)), strconv.Itoa(int(play.Extradb2)), strconv.Itoa(int(play.Extradb3)),
			strconv.Itoa(int(play.Fs)), strconv.Itoa(int(play.Ss)), strconv.Itoa(int(play.Lcb)),
			strconv.Itoa(int(play.Blitzer1)), strconv.Itoa(int(play.Blitzer2)), strconv.Itoa(int(play.Blitzer3)),
		}

		err = writer.Write(row)
		if err != nil {
			log.Fatal("Cannot write player row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
}

func ExportCFBGameResults(w http.ResponseWriter, seasonID, weekID, nflWeekID, timeslot string) {
	fileName := "slippery_jim_secret_results_list.csv"
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment;"+fileName)
	w.Header().Set("Transfer-Encoding", "chunked")
	writer := csv.NewWriter(w)
	ts := GetTimestamp()
	isExactWeek := weekID == strconv.Itoa(ts.CollegeWeekID) && seasonID == strconv.Itoa(ts.CollegeSeasonID)

	// Get All needed data
	matchChn := make(chan []structs.CollegeGame)
	nflMatchChn := make(chan []structs.NFLGame)

	go func() {
		matches := GetCollegeGamesByWeekId(weekID, timeslot, ts.CFBSpringGames)
		matchChn <- matches
	}()

	go func() {
		nbamatches := GetNFLGamesByWeekId(nflWeekID, ts.NFLPreseason)
		nflMatchChn <- nbamatches
	}()

	collegeGames := <-matchChn
	close(matchChn)
	nflGames := <-nflMatchChn
	close(nflMatchChn)

	HeaderRow := []string{
		"League", "Week", "Home Team", "Home Score",
		"Away Team", "Away Score", "Home Coach", "Home Rank", "Away Coach", "Away Rank", "Game Title",
		"Neutral Site", "Conference", "Timeslot", "Stadium", "City", "State",
	}

	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	for _, m := range collegeGames {
		if !m.GameComplete {
			continue
		}
		gameTime := m.TimeSlot
		gameNotRan := (gameTime == "Thursday Night" && !ts.ThursdayGames) ||
			(gameTime == "Friday Night" && !ts.FridayGames) ||
			(gameTime == "Saturday Morning" && !ts.SaturdayMorning) ||
			(gameTime == "Saturday Afternoon" && !ts.SaturdayNoon) ||
			(gameTime == "Saturday Evening" && !ts.SaturdayEvening) ||
			(gameTime == "Saturday Night" && !ts.SaturdayNight)
		if isExactWeek && gameNotRan {
			m.HideScore()
		}
		neutralStr := "N"
		if m.IsNeutral {
			neutralStr = "Y"
		}
		confStr := "N"
		if m.IsConference {
			confStr = "Y"
		}

		row := []string{
			"CFB", strconv.Itoa(int(m.Week)), m.HomeTeam, strconv.Itoa(int(m.HomeTeamScore)),
			m.AwayTeam, strconv.Itoa(int(m.AwayTeamScore)), m.HomeTeamCoach,
			strconv.Itoa(int(m.HomeTeamRank)), m.AwayTeamCoach, strconv.Itoa(int(m.AwayTeamRank)), m.GameTitle,
			neutralStr, confStr, m.TimeSlot, m.Stadium, m.City, m.State,
		}
		err = writer.Write(row)
		if err != nil {
			log.Fatal("Cannot write croot row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
	for _, m := range nflGames {
		if !m.GameComplete {
			continue
		}
		gameTime := m.TimeSlot
		gameNotRan := (gameTime == "Thursday Night Football" && !ts.NFLThursday) ||
			(gameTime == "Sunday Noon" && !ts.NFLSundayNoon) ||
			(gameTime == "Sunday Afternoon" && !ts.NFLSundayAfternoon) ||
			(gameTime == "Sunday Night Football" && !ts.NFLSundayEvening) ||
			(gameTime == "Monday Night Football" && !ts.NFLMondayEvening)
		if isExactWeek && gameNotRan {
			m.HideScore()
		}
		neutralStr := "N"
		if m.IsNeutral {
			neutralStr = "Y"
		}
		confStr := "N"
		if m.IsConference {
			confStr = "Y"
		}
		divStr := "N"
		if m.IsDivisional {
			divStr = "Y"
		}

		row := []string{
			"NFL", strconv.Itoa(int(m.Week)), m.GameTitle, m.HomeTeam, m.AwayTeamCoach,
			"N/A", strconv.Itoa(int(m.HomeTeamScore)),
			m.AwayTeam, m.AwayTeamCoach, "N/A", strconv.Itoa(int(m.AwayTeamScore)),
			neutralStr, confStr, divStr, m.TimeSlot, m.Stadium, m.City, m.State,
		}
		err = writer.Write(row)
		if err != nil {
			log.Fatal("Cannot write croot row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
}

func ExportCollegePlayerStatsToCSV(cp []structs.CollegePlayerResponse, viewType string, w http.ResponseWriter) {
	ts := GetTimestamp()
	seasonStr := strconv.Itoa(ts.Season)
	weekStr := ""
	if viewType != "SEASON" {
		weekStr = "WEEK_" + strconv.Itoa(ts.CollegeWeek) + "_"
	}
	fileName := "toucans_secret_" + seasonStr + "_" + weekStr + "stats"
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment;filename="+fileName+".csv")
	w.Header().Set("Transfer-Encoding", "chunked")
	// Initialize writer
	writer := csv.NewWriter(w)

	// Get Players

	HeaderRow := []string{
		"First Name", "Last Name", "Position",
		"Archetype", "Year", "Is Redshirt?", "Team", "Conference", "Age", "Stars",
		"Passing Yards", "Pass Attempts", "Pass Completions", "Passing Avg",
		"Passing TDs", "Interceptions", "Longest Pass", "QB Sacks",
		"QB Rating", "Rush Attempts", "Rushing Yards", "Rushing Avg",
		"Rushing TDs", "Fumbles", "Longest Rush", "Targets",
		"Catches", "Receiving Yards", "Receiving Avg", "Receiving TDs",
		"Longest Reception", "Solo Tackles", "Assisted Tackles", "Tackles For Loss",
		"Defensive Sacks", "Forced Fumbles", "Pass Deflections", "Interceptions Caught",
		"Safeties", "Defensive TDs", "FG Made", "FG Attempts",
		"Longest FG", "Extra Points Made", "Extra Point Attempts", "Kickoff TBs",
		"Punts", "Punt Touchbacks", "Punts Inside 20", "Kick Returns",
		"Kick Return TDs", "Kick Return Yards", "Punt Returns", "Punt Return TDs",
		"Punt Return Yards", "ST Solo Tackles", "ST Assisted Tackles", "Punts Blocked",
		"FG Blocked", "Snaps", "Pancakes", "Likely to transfer to Guam?",
	}

	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	for _, p := range cp {
		Year, RedshirtStatus := util.GetYearAndRedshirtStatus(p.Year, p.IsRedshirt)
		seasonStats := p.SeasonStats
		if viewType != "SEASON" {
			seasonStats = structs.CollegePlayerSeasonStats{}
			seasonStats.MapStats([]structs.CollegePlayerStats{p.Stats})
		}

		answer := "No."
		diceRoll := util.GenerateIntFromRange(1, 1000)
		if diceRoll == 1000 {
			answer = "Yes."
		}

		pr := []string{p.FirstName, p.LastName, p.Position,
			p.Archetype, Year, RedshirtStatus, p.TeamAbbr, p.Conference, strconv.Itoa(p.Age), strconv.Itoa(p.Stars),
			strconv.Itoa(seasonStats.PassingYards), strconv.Itoa(seasonStats.PassAttempts), strconv.Itoa(seasonStats.PassCompletions), strconv.Itoa(int(seasonStats.PassingAvg)),
			strconv.Itoa(seasonStats.PassingTDs), strconv.Itoa(seasonStats.Interceptions), strconv.Itoa(seasonStats.LongestPass), strconv.Itoa(seasonStats.Sacks),
			strconv.Itoa(int(seasonStats.QBRating)), strconv.Itoa(seasonStats.RushAttempts), strconv.Itoa(seasonStats.RushingYards), strconv.Itoa(int(seasonStats.RushingAvg)),
			strconv.Itoa(seasonStats.RushingTDs), strconv.Itoa(seasonStats.Fumbles), strconv.Itoa(seasonStats.LongestRush), strconv.Itoa(seasonStats.Targets),
			strconv.Itoa(seasonStats.Catches), strconv.Itoa(seasonStats.ReceivingYards), strconv.Itoa(int(seasonStats.ReceivingAvg)), strconv.Itoa(seasonStats.ReceivingTDs),
			strconv.Itoa(seasonStats.LongestReception), strconv.Itoa(int(seasonStats.SoloTackles)), strconv.Itoa(int(seasonStats.AssistedTackles)), strconv.Itoa(int(seasonStats.TacklesForLoss)),
			strconv.Itoa(int(seasonStats.SacksMade)), strconv.Itoa(seasonStats.ForcedFumbles), strconv.Itoa(seasonStats.PassDeflections), strconv.Itoa(seasonStats.InterceptionsCaught),
			strconv.Itoa(seasonStats.Safeties), strconv.Itoa(seasonStats.DefensiveTDs), strconv.Itoa(seasonStats.FGMade), strconv.Itoa(seasonStats.FGAttempts),
			strconv.Itoa(seasonStats.LongestFG), strconv.Itoa(seasonStats.ExtraPointsMade), strconv.Itoa(seasonStats.ExtraPointsAttempted), strconv.Itoa(seasonStats.KickoffTouchbacks),
			strconv.Itoa(seasonStats.Punts), strconv.Itoa(seasonStats.PuntTouchbacks), strconv.Itoa(seasonStats.PuntsInside20), strconv.Itoa(seasonStats.KickReturns),
			strconv.Itoa(seasonStats.KickReturnTDs), strconv.Itoa(seasonStats.KickReturnYards), strconv.Itoa(seasonStats.PuntReturns), strconv.Itoa(seasonStats.PuntReturnTDs),
			strconv.Itoa(seasonStats.PuntReturnYards), strconv.Itoa(int(seasonStats.STSoloTackles)), strconv.Itoa(int(seasonStats.STAssistedTackles)), strconv.Itoa(seasonStats.PuntsBlocked),
			strconv.Itoa(seasonStats.FGBlocked), strconv.Itoa(seasonStats.Snaps), strconv.Itoa(seasonStats.Pancakes), answer,
		}
		err = writer.Write(pr)
		if err != nil {
			log.Fatal("Cannot write player row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
}

func ExportNFLPlayerStatsToCSV(cp []structs.NFLPlayerResponse, viewType string, w http.ResponseWriter) {
	ts := GetTimestamp()
	seasonStr := strconv.Itoa(ts.Season)
	weekStr := ""
	if viewType != "SEASON" {
		weekStr = "WEEK_" + strconv.Itoa(ts.NFLWeek) + "_"
	}
	fileName := "toucans_other_secret_" + seasonStr + "_" + weekStr + "stats"
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment;filename="+fileName+".csv")
	w.Header().Set("Transfer-Encoding", "chunked")
	// Initialize writer
	writer := csv.NewWriter(w)

	// Get Players

	HeaderRow := []string{
		"First Name", "Last Name", "Position",
		"Archetype", "Position Two", "Archetype Two", "Year", "Team", "Conference", "Division", "Age", "Stars",
		"Passing Yards", "Pass Attempts", "Pass Completions", "Passing Avg",
		"Passing TDs", "Interceptions", "Longest Pass", "QB Sacks",
		"QB Rating", "Rush Attempts", "Rushing Yards", "Rushing Avg",
		"Rushing TDs", "Fumbles", "Longest Rush", "Targets",
		"Catches", "Receiving Yards", "Receiving Avg", "Receiving TDs",
		"Longest Reception", "Solo Tackles", "Assisted Tackles", "Tackles For Loss",
		"Defensive Sacks", "Forced Fumbles", "Pass Deflections", "Interceptions Caught",
		"Safeties", "Defensive TDs", "FG Made", "FG Attempts",
		"Longest FG", "Extra Points Made", "Extra Point Attempts", "Kickoff TBs",
		"Punts", "Punt Touchbacks", "Punts Inside 20", "Kick Returns",
		"Kick Return TDs", "Kick Return Yards", "Punt Returns", "Punt Return TDs",
		"Punt Return Yards", "ST Solo Tackles", "ST Assisted Tackles", "Punts Blocked",
		"FG Blocked", "Snaps", "Pancakes", "Likelihood to be traded to WAS?",
	}

	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	for _, p := range cp {
		Year := strconv.Itoa(p.Year)
		seasonStats := p.SeasonStats
		if viewType != "SEASON" {
			seasonStats = structs.NFLPlayerSeasonStats{}
			seasonStats.MapStats([]structs.NFLPlayerStats{p.Stats}, ts)
		}

		pr := []string{p.FirstName, p.LastName, p.Position,
			p.Archetype, p.PositionTwo, p.ArchetypeTwo, Year, p.TeamAbbr, p.Conference, p.Division, strconv.Itoa(p.Age), strconv.Itoa(p.Stars),
			strconv.Itoa(seasonStats.PassingYards), strconv.Itoa(seasonStats.PassAttempts), strconv.Itoa(seasonStats.PassCompletions), strconv.Itoa(int(seasonStats.PassingAvg)),
			strconv.Itoa(seasonStats.PassingTDs), strconv.Itoa(seasonStats.Interceptions), strconv.Itoa(seasonStats.LongestPass), strconv.Itoa(seasonStats.Sacks),
			strconv.Itoa(int(seasonStats.QBRating)), strconv.Itoa(seasonStats.RushAttempts), strconv.Itoa(seasonStats.RushingYards), strconv.Itoa(int(seasonStats.RushingAvg)),
			strconv.Itoa(seasonStats.RushingTDs), strconv.Itoa(seasonStats.Fumbles), strconv.Itoa(seasonStats.LongestRush), strconv.Itoa(seasonStats.Targets),
			strconv.Itoa(seasonStats.Catches), strconv.Itoa(seasonStats.ReceivingYards), strconv.Itoa(int(seasonStats.ReceivingAvg)), strconv.Itoa(seasonStats.ReceivingTDs),
			strconv.Itoa(seasonStats.LongestReception), strconv.Itoa(int(seasonStats.SoloTackles)), strconv.Itoa(int(seasonStats.AssistedTackles)), strconv.Itoa(int(seasonStats.TacklesForLoss)),
			strconv.Itoa(int(seasonStats.SacksMade)), strconv.Itoa(seasonStats.ForcedFumbles), strconv.Itoa(seasonStats.PassDeflections), strconv.Itoa(seasonStats.InterceptionsCaught),
			strconv.Itoa(seasonStats.Safeties), strconv.Itoa(seasonStats.DefensiveTDs), strconv.Itoa(seasonStats.FGMade), strconv.Itoa(seasonStats.FGAttempts),
			strconv.Itoa(seasonStats.LongestFG), strconv.Itoa(seasonStats.ExtraPointsMade), strconv.Itoa(seasonStats.ExtraPointsAttempted), strconv.Itoa(seasonStats.KickoffTouchbacks),
			strconv.Itoa(seasonStats.Punts), strconv.Itoa(seasonStats.PuntTouchbacks), strconv.Itoa(seasonStats.PuntsInside20), strconv.Itoa(seasonStats.KickReturns),
			strconv.Itoa(seasonStats.KickReturnTDs), strconv.Itoa(seasonStats.KickReturnYards), strconv.Itoa(seasonStats.PuntReturns), strconv.Itoa(seasonStats.PuntReturnTDs),
			strconv.Itoa(seasonStats.PuntReturnYards), strconv.Itoa(int(seasonStats.STSoloTackles)), strconv.Itoa(int(seasonStats.STAssistedTackles)), strconv.Itoa(seasonStats.PuntsBlocked),
			strconv.Itoa(seasonStats.FGBlocked), strconv.Itoa(seasonStats.Snaps), strconv.Itoa(seasonStats.Pancakes), "No.",
		}
		err = writer.Write(pr)
		if err != nil {
			log.Fatal("Cannot write player row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
}

func ExportTransferPortalToCSV(w http.ResponseWriter) {
	// Get Team Data
	w.Header().Set("Content-Disposition", "attachment;filename=Official_Portal_List.csv")
	w.Header().Set("Transfer-Encoding", "chunked")
	// Initialize writer
	writer := csv.NewWriter(w)

	// Get Players
	players := GetTransferPortalPlayersForPage()

	HeaderRow := []string{
		"Previous Team", "Player ID", "First Name", "Last Name", "Position",
		"Archetype", "Position Two", "Archetype Two", "Year", "Age", "Stars",
		"State", "Height",
		"Weight", "Overall", "Speed",
		"Football IQ", "Agility", "Carrying",
		"Catching", "Route Running", "Zone Coverage", "Man Coverage",
		"Strength", "Tackle", "Pass Block", "Run Block",
		"Pass Rush", "Run Defense", "Throw Power", "Throw Accuracy",
		"Kick Power", "Kick Accuracy", "Punt Power", "Punt Accuracy",
		"Stamina", "Injury", "Potential Grade", "Redshirt Status",
	}

	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	for _, player := range players {
		csvModel := structs.MapPortalPlayerToCSVModel(player)
		idStr := strconv.Itoa(int(player.PlayerID))
		playerRow := []string{
			csvModel.Team, idStr, csvModel.FirstName, csvModel.LastName, csvModel.Position,
			csvModel.Archetype, csvModel.PositionTwo, csvModel.ArchetypeTwo, csvModel.Year, strconv.Itoa(player.Age), strconv.Itoa(player.Stars),
			player.State, strconv.Itoa(player.Height),
			strconv.Itoa(player.Weight), player.OverallGrade, csvModel.SpeedGrade,
			csvModel.FootballIQGrade, csvModel.AgilityGrade, csvModel.CarryingGrade,
			csvModel.CatchingGrade, csvModel.RouteRunningGrade, csvModel.ZoneCoverageGrade, csvModel.ManCoverageGrade,
			csvModel.StrengthGrade, csvModel.TackleGrade, csvModel.PassBlockGrade, csvModel.RunBlockGrade,
			csvModel.PassRushGrade, csvModel.RunDefenseGrade, csvModel.ThrowPowerGrade, csvModel.ThrowAccuracyGrade,
			csvModel.KickPowerGrade, csvModel.KickAccuracyGrade, csvModel.PuntPowerGrade, csvModel.PuntAccuracyGrade,
			csvModel.StaminaGrade, csvModel.InjuryGrade, csvModel.PotentialGrade, csvModel.RedshirtStatus,
		}

		err = writer.Write(playerRow)
		if err != nil {
			log.Fatal("Cannot write player row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
}

func ExportCFBSpringPlayByPlayToCSV(w http.ResponseWriter) {
	games := repository.FindCollegeGamesRecords("6", true)
	collegePlayers := GetAllCollegePlayers()
	participantMap := make(map[uint]structs.GameResultsPlayer)
	for _, p := range collegePlayers {
		participantMap[p.ID] = structs.GameResultsPlayer{
			ID:        p.ID,
			FirstName: p.FirstName,
			LastName:  p.LastName,
			Position:  p.Position,
			Archetype: p.Archetype,
			Year:      uint(p.Year),
			League:    "CFB",
		}
	}

	fileName := "cfb_spring_game_play_by_plays_all"
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment;filename="+fileName+".csv")
	w.Header().Set("Transfer-Encoding", "chunked")
	// Initialize writer
	writer := csv.NewWriter(w)
	HeaderRow := []string{
		"Game ID", "Play", "Home Team", "Home Team Score", "Away Team", "Away Team Score", "Quarter", "Time Remaining",
		"Possession", "Down", "Distance", "Line of Scrimmage", "Type of Play",
		"Offensive Formation", "Offensive Play", "Offensive PoA", "Defensive Formation",
		"Defensive Tendency", "# of Blitzers", "LB Coverage", "CB Coverage", "S Coverage",
		"QB Player ID", "Ballcarrier ID", "Tackler1 ID", "Tackler2 ID", "Yards Gained",
		"Result",
		"QB Id", "Back1 Id", "Back2 Id", "Back3 Id", "Slot1 Id", "Slot2 Id", "Le Id", "Re Id", "Lt Id", "LG id", "C id", "Rg Id", "Rt Id",
		"Rde id", "Rdt Id", "Nt Id", "Ldt Id", "Lde Id", "Rolb Id", "Rilb Id", "Mlb Id", "Lilb Id", "Lolb Id", "Rcb Id", "DB1 Id", "DB2 Id", "DB3 ID", "Fs ID", "SS Id", "Lcb Id",
		"Blitzer1 Id", "Blitzer2 Id", "Blitzer3 Id",
	}
	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	for _, game := range games {
		if !game.IsSpringGame {
			continue
		}
		gameID := strconv.Itoa(int(game.ID))

		playByPlays := GetCFBPlayByPlaysByGameID(gameID)
		// Generate the Play By Play Response
		playbyPlayResponseList := GenerateCFBPlayByPlayResponse(playByPlays, participantMap, false, game.HomeTeam, game.AwayTeam)

		for _, play := range playbyPlayResponseList {
			num := strconv.Itoa(int(play.PlayNumber))
			hcs := strconv.Itoa(int(play.HomeTeamScore))
			acs := strconv.Itoa(int(play.AwayTeamScore))
			qt := strconv.Itoa(int(play.Quarter))
			tr := play.TimeRemaining
			down := strconv.Itoa(int(play.Down))
			dist := strconv.Itoa(int(play.Distance))
			qbID := strconv.Itoa(int(play.QBPlayerID))
			bcID := strconv.Itoa(int(play.BallCarrierID))
			t1ID := strconv.Itoa(int(play.Tackler1ID))
			t2ID := strconv.Itoa(int(play.Tackler2ID))
			yards := strconv.Itoa(int(play.ResultYards))
			blitzNumber := strconv.Itoa(int(play.BlitzNumber))

			row := []string{
				gameID, num, game.HomeTeam, hcs, game.AwayTeam, acs, qt, tr, play.Possession, down, dist, play.LineOfScrimmage,
				play.PlayType, play.OffensiveFormation, play.PlayName, play.PointOfAttack, play.DefensiveFormation,
				play.DefensiveTendency, blitzNumber, play.LBCoverage, play.CBCoverage, play.SCoverage,
				qbID, bcID, t1ID, t2ID, yards,
				play.Result,
				strconv.Itoa(int(play.Qb)), strconv.Itoa(int(play.Back1)), strconv.Itoa(int(play.Back2)), strconv.Itoa(int(play.Back3)),
				strconv.Itoa(int(play.Slot1)), strconv.Itoa(int(play.Slot2)), strconv.Itoa(int(play.Le)), strconv.Itoa(int(play.Re)),
				strconv.Itoa(int(play.Lt)), strconv.Itoa(int(play.Lg)), strconv.Itoa(int(play.C)), strconv.Itoa(int(play.Rg)), strconv.Itoa(int(play.Rt)),
				strconv.Itoa(int(play.Rde)), strconv.Itoa(int(play.Rdt)), strconv.Itoa(int(play.Nt)), strconv.Itoa(int(play.Ldt)), strconv.Itoa(int(play.Lde)),
				strconv.Itoa(int(play.Rolb)), strconv.Itoa(int(play.Rilb)), strconv.Itoa(int(play.Mlb)), strconv.Itoa(int(play.Lilb)), strconv.Itoa(int(play.Lolb)),
				strconv.Itoa(int(play.Rcb)),
				strconv.Itoa(int(play.Extradb1)), strconv.Itoa(int(play.Extradb2)), strconv.Itoa(int(play.Extradb3)),
				strconv.Itoa(int(play.Fs)), strconv.Itoa(int(play.Ss)), strconv.Itoa(int(play.Lcb)),
				strconv.Itoa(int(play.Blitzer1)), strconv.Itoa(int(play.Blitzer2)), strconv.Itoa(int(play.Blitzer3)),
			}

			err = writer.Write(row)
			if err != nil {
				log.Fatal("Cannot write player row to CSV", err)
			}

			writer.Flush()
			err = writer.Error()
			if err != nil {
				log.Fatal("Error while writing to file ::", err)
			}
		}
	}
}

func ExportNFLPreseasonPlayByPlayToCSV(w http.ResponseWriter) {
	games := repository.FindNFLGamesRecords("6", true)
	nflPlayers := GetAllNFLPlayers()
	participantMap := make(map[uint]structs.GameResultsPlayer)
	for _, p := range nflPlayers {
		participantMap[p.ID] = structs.GameResultsPlayer{
			ID:        p.ID,
			FirstName: p.FirstName,
			LastName:  p.LastName,
			Position:  p.Position,
			Archetype: p.Archetype,
			Year:      uint(p.Experience),
			League:    "NFL",
		}
	}

	fileName := "nfl_preseason_play_by_plays_all"
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment;filename="+fileName+".csv")
	w.Header().Set("Transfer-Encoding", "chunked")
	// Initialize writer
	writer := csv.NewWriter(w)
	HeaderRow := []string{
		"Game ID", "Play", "Home Team", "Home Team Score", "Away Team", "Away Team Score", "Quarter", "Time Remaining",
		"Possession", "Down", "Distance", "Line of Scrimmage", "Type of Play",
		"Offensive Formation", "Offensive Play", "Offensive PoA", "Defensive Formation",
		"Defensive Tendency", "# of Blitzers", "LB Coverage", "CB Coverage", "S Coverage",
		"QB Player ID", "Ballcarrier ID", "Tackler1 ID", "Tackler2 ID", "Yards Gained",
		"Result",
		"QB Id", "Back1 Id", "Back2 Id", "Back3 Id", "Slot1 Id", "Slot2 Id", "Le Id", "Re Id", "Lt Id", "LG id", "C id", "Rg Id", "Rt Id",
		"Rde id", "Rdt Id", "Nt Id", "Ldt Id", "Lde Id", "Rolb Id", "Rilb Id", "Mlb Id", "Lilb Id", "Lolb Id", "Rcb Id", "DB1 Id", "DB2 Id", "DB3 ID", "Fs ID", "SS Id", "Lcb Id",
		"Blitzer1 Id", "Blitzer2 Id", "Blitzer3 Id",
	}
	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	for _, game := range games {
		if !game.IsPreseasonGame {
			continue
		}
		gameID := strconv.Itoa(int(game.ID))

		playByPlays := GetNFLPlayByPlaysByGameID(gameID)
		// Generate the Play By Play Response
		playbyPlayResponseList := GenerateNFLPlayByPlayResponse(playByPlays, participantMap, false, game.HomeTeam, game.AwayTeam)

		for _, play := range playbyPlayResponseList {
			num := strconv.Itoa(int(play.PlayNumber))
			hcs := strconv.Itoa(int(play.HomeTeamScore))
			acs := strconv.Itoa(int(play.AwayTeamScore))
			qt := strconv.Itoa(int(play.Quarter))
			tr := play.TimeRemaining
			down := strconv.Itoa(int(play.Down))
			dist := strconv.Itoa(int(play.Distance))
			qbID := strconv.Itoa(int(play.QBPlayerID))
			bcID := strconv.Itoa(int(play.BallCarrierID))
			t1ID := strconv.Itoa(int(play.Tackler1ID))
			t2ID := strconv.Itoa(int(play.Tackler2ID))
			yards := strconv.Itoa(int(play.ResultYards))
			blitzNumber := strconv.Itoa(int(play.BlitzNumber))

			row := []string{
				gameID, num, game.HomeTeam, hcs, game.AwayTeam, acs, qt, tr, play.Possession, down, dist, play.LineOfScrimmage,
				play.PlayType, play.OffensiveFormation, play.PlayName, play.PointOfAttack, play.DefensiveFormation,
				play.DefensiveTendency, blitzNumber, play.LBCoverage, play.CBCoverage, play.SCoverage,
				qbID, bcID, t1ID, t2ID, yards,
				play.Result,
				strconv.Itoa(int(play.Qb)), strconv.Itoa(int(play.Back1)), strconv.Itoa(int(play.Back2)), strconv.Itoa(int(play.Back3)),
				strconv.Itoa(int(play.Slot1)), strconv.Itoa(int(play.Slot2)), strconv.Itoa(int(play.Le)), strconv.Itoa(int(play.Re)),
				strconv.Itoa(int(play.Lt)), strconv.Itoa(int(play.Lg)), strconv.Itoa(int(play.C)), strconv.Itoa(int(play.Rg)), strconv.Itoa(int(play.Rt)),
				strconv.Itoa(int(play.Rde)), strconv.Itoa(int(play.Rdt)), strconv.Itoa(int(play.Nt)), strconv.Itoa(int(play.Ldt)), strconv.Itoa(int(play.Lde)),
				strconv.Itoa(int(play.Rolb)), strconv.Itoa(int(play.Rilb)), strconv.Itoa(int(play.Mlb)), strconv.Itoa(int(play.Lilb)), strconv.Itoa(int(play.Lolb)),
				strconv.Itoa(int(play.Rcb)),
				strconv.Itoa(int(play.Extradb1)), strconv.Itoa(int(play.Extradb2)), strconv.Itoa(int(play.Extradb3)),
				strconv.Itoa(int(play.Fs)), strconv.Itoa(int(play.Ss)), strconv.Itoa(int(play.Lcb)),
				strconv.Itoa(int(play.Blitzer1)), strconv.Itoa(int(play.Blitzer2)), strconv.Itoa(int(play.Blitzer3)),
			}

			err = writer.Write(row)
			if err != nil {
				log.Fatal("Cannot write player row to CSV", err)
			}

			writer.Flush()
			err = writer.Error()
			if err != nil {
				log.Fatal("Error while writing to file ::", err)
			}
		}
	}
}
