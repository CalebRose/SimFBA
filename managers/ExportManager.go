package managers

import (
	"encoding/csv"
	"log"
	"net/http"
	"strconv"

	"github.com/CalebRose/SimFBA/util"
)

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
		csvModel := util.MapPlayerToCSVModel(player)
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
