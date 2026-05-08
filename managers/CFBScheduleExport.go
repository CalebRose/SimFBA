package managers

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/CalebRose/SimFBA/structs"
)

// ExportScheduleToCSV writes collegeGamesUpload data to a CSV file at filePath
// so the schedule can be reviewed before committing to the database.
func ExportScheduleToCSV(games []structs.CollegeGame, filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	header := []string{
		"SeasonID", "WeekID", "Week",
		"HomeTeamID", "HomeTeam",
		"AwayTeamID", "AwayTeam",
		"IsConference", "IsDivisional", "IsRivalry", "IsNeutral",
		"StadiumID", "Stadium", "City", "State",
		"TimeSlot",
	}
	if err := w.Write(header); err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	for _, g := range games {
		row := []string{
			strconv.Itoa(g.SeasonID),
			strconv.Itoa(g.WeekID),
			strconv.Itoa(g.Week),
			strconv.Itoa(g.HomeTeamID),
			g.HomeTeam,
			strconv.Itoa(g.AwayTeamID),
			g.AwayTeam,
			strconv.FormatBool(g.IsConference),
			strconv.FormatBool(g.IsDivisional),
			strconv.FormatBool(g.IsRivalryGame),
			strconv.FormatBool(g.IsNeutral),
			strconv.Itoa(int(g.StadiumID)),
			g.Stadium,
			g.City,
			g.State,
			g.TimeSlot,
		}
		if err := w.Write(row); err != nil {
			return fmt.Errorf("write row: %w", err)
		}
	}

	return w.Error()
}
