package managers

import (
	"encoding/csv"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/CalebRose/SimFBA/models"
	"github.com/CalebRose/SimFBA/util"
)

func ExportAllRostersToCSV(w http.ResponseWriter) {
	w.Header().Set("Content-Disposition", "attachment;filename=2023_Rosters.csv")
	w.Header().Set("Transfer-Encoding", "chunked")
	// Initialize writer
	writer := csv.NewWriter(w)

	HeaderRow := []string{
		"Team", "First Name", "Last Name", "Position",
		"Archetype", "Year", "Age", "Stars",
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
		csvModel := models.MapPlayerToCSVModel(player)
		playerRow := []string{
			player.TeamAbbr, csvModel.FirstName, csvModel.LastName, csvModel.Position,
			csvModel.Archetype, csvModel.Year, strconv.Itoa(player.Age), strconv.Itoa(player.Stars),
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
	w.Header().Set("Content-Disposition", "attachment;filename="+team.TeamName+".csv")
	w.Header().Set("Transfer-Encoding", "chunked")
	// Initialize writer
	writer := csv.NewWriter(w)

	// Get Players
	players := GetAllCollegePlayersByTeamId(TeamID)

	HeaderRow := []string{
		"Team", "First Name", "Last Name", "Position",
		"Archetype", "Year", "Age", "Stars",
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
		csvModel := models.MapPlayerToCSVModel(player)
		playerRow := []string{
			team.TeamName, csvModel.FirstName, csvModel.LastName, csvModel.Position,
			csvModel.Archetype, csvModel.Year, strconv.Itoa(player.Age), strconv.Itoa(player.Stars),
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
	w.Header().Set("Content-Disposition", "attachment;filename="+team.TeamName+".csv")
	w.Header().Set("Transfer-Encoding", "chunked")
	// Initialize writer
	writer := csv.NewWriter(w)

	// Get Players
	players := GetNFLPlayersWithContractsByTeamID(TeamID)

	HeaderRow := []string{
		"Team", "First Name", "Last Name", "Position",
		"Archetype", "Year", "Age",
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
		csvModel := models.MapNFLPlayerToCSVModel(player)
		playerRow := []string{
			team.TeamName, csvModel.FirstName, csvModel.LastName, csvModel.Position,
			csvModel.Archetype, csvModel.Year, strconv.Itoa(player.Age),
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

func ExportCrootsToCSV(w http.ResponseWriter) {
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
			strconv.Itoa(croot.PlayerID), croot.FirstName, croot.LastName, croot.Position,
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
	w.Header().Set("Content-Disposition", "attachment;filename=2024SimNFLDraftClass.csv")
	w.Header().Set("Transfer-Encoding", "chunked")
	// Initialize writer
	writer := csv.NewWriter(w)

	// Get NFL Draft Class
	draftees := GetAllNFLDraftees()

	HeaderRow := []string{
		"PlayerID", "First Name", "Last Name", "Position",
		"Archetype", "Age", "Stars", "College",
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
			strconv.Itoa(player.PlayerID), player.FirstName, player.LastName, player.Position,
			player.Archetype, strconv.Itoa(player.Age), strconv.Itoa(player.Stars), player.College,
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

func ExportPlayerStatsToCSV(cp []models.CollegePlayerResponse, w http.ResponseWriter) {
	w.Header().Set("Content-Disposition", "attachment;filename=special_season_stats.csv")
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
		"FG Blocked", "Snaps", "Pancakes", "Ready for Mars?",
	}

	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	for _, p := range cp {
		Year, RedshirtStatus := util.GetYearAndRedshirtStatus(p.Year, p.IsRedshirt)
		seasonStats := p.SeasonStats

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
