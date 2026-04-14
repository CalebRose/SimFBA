package managers

import (
	"sort"
	"strconv"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
	"gorm.io/gorm"
)

func UpdateCollegeAIDepthChartsTEST() {
	db := dbprovider.GetInstance().GetDB()
	teams := GetAllCollegeTeams()
	for _, team := range teams {

		teamID := strconv.Itoa(int(team.ID))
		gp := repository.GetGameplanTESTByTeamID(teamID)
		ReAlignCollegeDepthChartTEST(db, teamID, gp)
	}

	ts := GetTimestamp()
	ts.ToggleAIDepthCharts()
	repository.SaveTimestamp(ts, db)
}

func GetTestOffensiveSchemesByTeamID(id uint) string {
	if id == 1 || id == 59 || id == 65 || id == 77 || id == 107 {
		return "West Coast"
	}
	if id == 7 || id == 54 || id == 98 {
		return "Power Run"
	}
	if id == 10 || id == 123 || id == 125 {
		return "Double Wing"
	}
	if id == 13 || id == 12 || id == 47 || id == 118 {
		return "Spread Option"
	}
	if id == 15 || id == 34 || id == 78 || id == 80 {
		return "Wing-T"
	}
	if id == 19 || id == 44 || id == 45 {
		return "Flexbone"
	}
	if id == 23 || id == 37 || id == 109 {
		return "Air Raid"
	}
	if id == 55 || id == 88 {
		return "I Option"
	}
	if id == 56 || id == 86 || id == 93 || id == 100 || id == 115 {
		return "Vertical"
	}
	if id == 63 || id == 99 || id == 108 || id == 62 || id == 39 {
		return "Pistol"
	}
	if id == 75 || id == 96 || id == 122 {
		return "Run and Shoot"
	}
	if id == 94 || id == 97 || id == 127 {
		return "Wishbone"
	}
	return ""
}

func GetTestDefensiveSchemesByTeamID(id uint) string {
	if id == 10 || id == 13 || id == 54 || id == 77 || id == 86 || id == 93 || id == 94 || id == 97 || id == 107 {
		return "Old School"
	}
	if id == 15 || id == 19 || id == 44 || id == 55 || id == 56 || id == 63 || id == 98 || id == 118 {
		return "2-Gap"
	}
	if id == 1 || id == 12 || id == 34 || id == 47 || id == 80 || id == 108 || id == 109 || id == 122 || id == 127 {
		return "4-Man Front Spread Stopper"
	}
	if id == 23 || id == 65 || id == 99 || id == 123 {
		return "3-Man Front Spread Stopper"
	}
	if id == 37 || id == 45 || id == 75 || id == 78 || id == 88 || id == 96 || id == 100 || id == 39 || id == 62 {
		return "Speed"
	}
	if id == 7 || id == 59 || id == 115 || id == 125 {
		return "Multiple"
	}
	return ""
}

func MassUpdateGameplanSchemesTEST(off, def string) {
	db := dbprovider.GetInstance().GetDB()
	teams := GetAllCollegeTeams()
	offensiveSchemes := GetOffensiveDefaultSchemes()
	defensiveSchemes := GetDefensiveDefaultSchemes()
	for _, team := range teams {
		teamID := strconv.Itoa(int(team.ID))
		gp := repository.GetGameplanTESTByTeamID(teamID)
		gp.UpdateSchemes(off, def)
		// offe := GetTestOffensiveSchemesByTeamID(id)
		// defe := GetTestDefensiveSchemesByTeamID(id)
		// Map Default Scheme for offense & defense
		offFormations := offensiveSchemes[off]
		defFormations := defensiveSchemes[def][off]

		dto := structs.CollegeGameplanTEST{
			TeamID: int(team.ID),
			BaseGameplan: structs.BaseGameplan{
				OffensiveScheme:    off,
				DefensiveScheme:    def,
				OffensiveFormation: offFormations,
				DefensiveFormation: defFormations,
				BlitzSafeties:      gp.BlitzSafeties,
				BlitzCorners:       gp.BlitzCorners,
				LinebackerCoverage: gp.LinebackerCoverage,
				MaximumFGDistance:  gp.MaximumFGDistance,
				GoFor4AndShort:     gp.GoFor4AndShort,
				GoFor4AndLong:      gp.GoFor4AndLong,
				DefaultOffense:     gp.DefaultOffense,
				DefaultDefense:     gp.DefaultDefense,
				PrimaryHB:          75,
				PitchFocus:         50,
				DiveFocus:          50,
			},
		}

		gp.UpdateCollegeGameplanTEST(dto)

		// Autosort Depth Chart
		ReAlignCollegeDepthChartTEST(db, teamID, gp)

		db.Save(&gp)
	}
}

func UpdateIndividualGameplanSchemeTEST(teamID, off, def string) {
	db := dbprovider.GetInstance().GetDB()
	offensiveSchemes := GetOffensiveDefaultSchemes()
	defensiveSchemes := GetDefensiveDefaultSchemes()

	gp := repository.GetGameplanTESTByTeamID(teamID)
	gp.UpdateSchemes(off, def)
	// Map Default Scheme for offense & defense
	offFormations := offensiveSchemes[off]
	defFormations := defensiveSchemes[def][off]

	dto := structs.CollegeGameplanTEST{
		TeamID: gp.TeamID,
		BaseGameplan: structs.BaseGameplan{
			OffensiveScheme:    off,
			DefensiveScheme:    def,
			OffensiveFormation: offFormations,
			DefensiveFormation: defFormations,
			BlitzSafeties:      gp.BlitzSafeties,
			BlitzCorners:       gp.BlitzCorners,
			LinebackerCoverage: gp.LinebackerCoverage,
			MaximumFGDistance:  gp.MaximumFGDistance,
			GoFor4AndShort:     gp.GoFor4AndShort,
			GoFor4AndLong:      gp.GoFor4AndLong,
			DefaultOffense:     gp.DefaultOffense,
			DefaultDefense:     gp.DefaultDefense,
			PrimaryHB:          75,
			PitchFocus:         50,
			DiveFocus:          50,
		},
	}

	gp.UpdateCollegeGameplanTEST(dto)

	// Autosort Depth Chart
	// ReAlignCollegeDepthChartTEST(db, teamID, gp)

	db.Save(&gp)

}

func ReAlignCollegeDepthChartTEST(db *gorm.DB, teamID string, gp structs.CollegeGameplanTEST) {
	roster := GetAllCollegePlayersByTeamIdWithoutRedshirts(teamID)
	dcPositions := repository.GetDepthChartPositionPlayersByDepthchartIDTEST(teamID)
	sort.Sort(structs.ByOverall(roster))
	positionMap := make(map[string][]structs.DepthChartPositionDTO)
	starterMap := make(map[uint]bool)
	backupMap := make(map[uint]bool)
	stuMap := make(map[uint]bool)
	offScheme := gp.OffensiveScheme
	defScheme := gp.DefensiveScheme
	isLT := true
	isLG := true
	isLE := true
	isLOLB := true

	goodFits := GetFitsByScheme(offScheme, false)
	badFits := GetFitsByScheme(defScheme, false)
	bonus := 5

	// Allocate the Position Map
	for _, cp := range roster {
		if cp.IsInjured || cp.IsRedshirting {
			continue
		}
		pos := cp.Position
		arch := cp.Archetype
		player := arch + " " + pos
		isGoodFit := CheckPlayerFits(player, goodFits)
		isBadFit := CheckPlayerFits(player, badFits)

		// Add to QB List
		if pos == "QB" || pos == "RB" || pos == "FB" || pos == "ATH" {
			score := 0
			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}
			if pos == "QB" {
				score += 75
			} else if pos == "ATH" && (arch == "Triple-Threat" || arch == "Field General") {
				score += 50
			}
			// score += ((cp.ThrowAccuracy + cp.ThrowPower) / 2)
			score += cp.Overall

			dcpObj := structs.DepthChartPositionDTO{
				Position:      pos,
				Archetype:     arch,
				Score:         score,
				CollegePlayer: cp,
			}
			positionMap["QB"] = append(positionMap["QB"], dcpObj)
		}
		// Add to RB List
		if pos == "RB" || pos == "FB" || pos == "WR" || pos == "TE" || pos == "ATH" {
			score := 0
			if pos == "RB" {
				score += 100
			} else if pos == "FB" {
				score += 25
			} else if pos == "ATH" && (arch == "Wingback" || arch == "Soccer Player" || arch == "Triple-Threat") {
				score += 50
			}
			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}

			score += ((cp.Speed + cp.Agility + cp.Strength + cp.Carrying) / 4)

			dcpObj := structs.DepthChartPositionDTO{
				Position:      pos,
				Archetype:     arch,
				Score:         score,
				CollegePlayer: cp,
			}
			positionMap["RB"] = append(positionMap["RB"], dcpObj)
		}

		// Add to FB List
		if pos == "FB" || pos == "TE" || pos == "RB" || pos == "ATH" {
			score := 0
			if pos == "FB" {
				score += 100
			} else if pos == "ATH" && (arch == "Wingback") {
				score += 50
			}

			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}
			score += ((cp.Strength + cp.Carrying + cp.PassBlock + cp.RunBlock) / 4)

			dcpObj := structs.DepthChartPositionDTO{
				Position:      pos,
				Archetype:     arch,
				Score:         score,
				CollegePlayer: cp,
			}
			positionMap["FB"] = append(positionMap["FB"], dcpObj)
		}

		// Add to TE List
		if pos == "FB" || pos == "TE" || pos == "ATH" {
			score := 0
			if pos == "TE" {
				score += 100
			} else if pos == "ATH" && (arch == "Slotback") {
				score += 50
			}

			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}

			score += int(float64(cp.Overall)*0.5) + int(float64(cp.RunBlock)*0.125) + int(float64(cp.PassBlock)*0.125) + int(float64(cp.Catching)*0.125) + int(float64(cp.Strength)*0.125)

			dcpObj := structs.DepthChartPositionDTO{
				Position:      pos,
				Archetype:     arch,
				Score:         score,
				CollegePlayer: cp,
			}
			positionMap["TE"] = append(positionMap["TE"], dcpObj)
		}
		// Add to WR List
		if pos == "WR" || pos == "TE" || pos == "RB" || pos == "ATH" {
			score := 0
			if pos == "WR" {
				score += 100
			} else if pos == "ATH" && (arch == "Wingback" || arch == "Slotback") {
				score += 50
			}

			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}

			score += int(float64(cp.Overall)*0.4) +
				int(float64(cp.Speed)*0.12) +
				int(float64(cp.Agility)*0.12) +
				int(float64(cp.Catching)*0.12) +
				int(float64(cp.Strength)*0.12) +
				int(float64(cp.RouteRunning)*0.12)

			dcpObj := structs.DepthChartPositionDTO{
				Position:      pos,
				Archetype:     arch,
				Score:         score,
				CollegePlayer: cp,
			}
			positionMap["WR"] = append(positionMap["WR"], dcpObj)
		}
		// Add to LT and RT List
		if pos == "OT" || pos == "OG" || pos == "C" || pos == "ATH" {
			score := 0
			if pos == "OT" {
				score += 100
			} else if pos == "OG" {
				score += 25
			} else if pos == "ATH" && (arch == "Lineman") {
				score += 50
			}

			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}

			score += int(float64(cp.Overall)*0.7) +
				int(float64(cp.Strength)*0.10) +
				int(float64(cp.RunBlock)*0.75) +
				int(float64(cp.PassBlock)*0.75) +
				int(float64(cp.Agility)*0.05)

			dcpObj := structs.DepthChartPositionDTO{
				Position:      pos,
				Archetype:     arch,
				Score:         score,
				CollegePlayer: cp,
			}
			if isLT {
				positionMap["LT"] = append(positionMap["LT"], dcpObj)
			} else {
				positionMap["RT"] = append(positionMap["RT"], dcpObj)
			}
			isLT = !isLT
		}
		// Add to LG and RG List
		if pos == "OT" || pos == "OG" || pos == "C" || pos == "ATH" {
			score := 0
			if pos == "OG" {
				score += 100
			} else if pos == "C" {
				score += 25
			} else if pos == "ATH" && (arch == "Lineman") {
				score += 50
			}

			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}

			score += int(float64(cp.Overall)*0.7) +
				int(float64(cp.Strength)*0.10) +
				int(float64(cp.RunBlock)*0.75) +
				int(float64(cp.PassBlock)*0.75) +
				int(float64(cp.Agility)*0.05)
			dcpObj := structs.DepthChartPositionDTO{
				Position:      pos,
				Archetype:     arch,
				Score:         score,
				CollegePlayer: cp,
			}
			if isLG {
				positionMap["LG"] = append(positionMap["LG"], dcpObj)
			} else {
				positionMap["RG"] = append(positionMap["RG"], dcpObj)
			}
			isLG = !isLG
		}
		// Add to C List
		if pos == "OT" || pos == "OG" || pos == "C" || pos == "ATH" {
			score := 0
			if pos == "C" {
				score += 100
			} else if pos == "OG" {
				score += 15
			} else if pos == "ATH" && (arch == "Lineman") {
				score += 50
			}

			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}

			score += int(float64(cp.Overall)*0.7) +
				int(float64(cp.Strength)*0.10) +
				int(float64(cp.RunBlock)*0.75) +
				int(float64(cp.PassBlock)*0.75) +
				int(float64(cp.Agility)*0.05)
			dcpObj := structs.DepthChartPositionDTO{
				Position:      pos,
				Archetype:     arch,
				Score:         score,
				CollegePlayer: cp,
			}
			positionMap["C"] = append(positionMap["C"], dcpObj)
		}

		// Add to LE List
		if pos == "DE" || pos == "DT" || pos == "OLB" || pos == "ATH" {
			score := 0
			if pos == "DE" {
				score += 100
			} else if pos == "OLB" {
				score += 25
			} else if pos == "DT" {
				score += 3
			} else if pos == "ATH" && (arch == "Lineman" || arch == "Strongside") {
				score += 50
			}

			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}

			score += int(float64(cp.Overall)*0.7) +
				int(float64(cp.Strength)*0.05) +
				int(float64(cp.Tackle)*0.05) +
				int(float64(cp.PassRush)*0.75) +
				int(float64(cp.RunDefense)*0.75) +
				int(float64(cp.Agility)*0.05)

			dcpObj := structs.DepthChartPositionDTO{
				Position:      pos,
				Archetype:     arch,
				Score:         score,
				CollegePlayer: cp,
			}
			if isLE {
				positionMap["LE"] = append(positionMap["LE"], dcpObj)
			} else {
				positionMap["RE"] = append(positionMap["RE"], dcpObj)
			}
			isLE = !isLE
		}

		// Add to DT list
		if pos == "DE" || pos == "DT" || pos == "OLB" || pos == "ATH" {
			score := 0
			if pos == "DT" {
				score += 100
			} else if pos == "DE" {
				score += 25
			} else if pos == "OLB" {
				score += 12
			} else if pos == "ATH" && (arch == "Lineman" || arch == "Strongside") {
				score += 50
			}

			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}

			score += int(float64(cp.Overall)*0.7) +
				int(float64(cp.Strength)*0.05) +
				int(float64(cp.Tackle)*0.05) +
				int(float64(cp.PassRush)*0.75) +
				int(float64(cp.RunDefense)*0.75) +
				int(float64(cp.Agility)*0.05)

			dcpObj := structs.DepthChartPositionDTO{
				Position:      pos,
				Archetype:     arch,
				Score:         score,
				CollegePlayer: cp,
			}
			positionMap["DT"] = append(positionMap["DT"], dcpObj)
		}

		// Add to OLB list
		if pos == "OLB" || pos == "DE" || pos == "ILB" || pos == "SS" || pos == "FS" || pos == "ATH" {
			score := 0
			if pos == "OLB" {
				score += 100
			} else if pos == "DE" {
				score += 10
			} else if pos == "ILB" {
				score += 25
			} else if pos == "SS" {
				score += 3
			} else if pos == "ATH" && (arch == "Weakside" || arch == "Strongside" || arch == "Bandit") {
				score += 50
			}

			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}

			score += int(float64(cp.Overall)*0.6) +
				int(float64(cp.Strength)*0.025) +
				int(float64(cp.Tackle)*0.055) +
				int(float64(cp.PassRush)*0.0755) +
				int(float64(cp.RunDefense)*0.0755) +
				int(float64(cp.ManCoverage)*0.075) +
				int(float64(cp.ZoneCoverage)*0.075) +
				int(float64(cp.Agility)*0.025)

			dcpObj := structs.DepthChartPositionDTO{
				Position:      pos,
				Archetype:     arch,
				Score:         score,
				CollegePlayer: cp,
			}
			if isLOLB {
				positionMap["LOLB"] = append(positionMap["LOLB"], dcpObj)
			} else {
				positionMap["ROLB"] = append(positionMap["ROLB"], dcpObj)
			}
			isLOLB = !isLOLB
		}

		// Add to ILB list
		if pos == "OLB" || pos == "DE" || pos == "ILB" || pos == "SS" || pos == "FS" || pos == "ATH" {
			score := 0
			if pos == "ILB" {
				score += 100
			} else if pos == "OLB" {
				score += 25
			} else if pos == "SS" {
				score += 8
			} else if pos == "DE" {
				score += 3
			} else if pos == "ATH" && (arch == "Weakside" || arch == "Bandit" || arch == "Field General") {
				score += 50
			}

			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}

			score += int(float64(cp.Overall)*0.6) +
				int(float64(cp.Strength)*0.025) +
				int(float64(cp.Tackle)*0.055) +
				int(float64(cp.PassRush)*0.0755) +
				int(float64(cp.RunDefense)*0.0755) +
				int(float64(cp.ManCoverage)*0.075) +
				int(float64(cp.ZoneCoverage)*0.075) +
				int(float64(cp.Agility)*0.025)

			dcpObj := structs.DepthChartPositionDTO{
				Position:      pos,
				Archetype:     arch,
				Score:         score,
				CollegePlayer: cp,
			}
			positionMap["MLB"] = append(positionMap["MLB"], dcpObj)
		}

		// Add to CB List
		if pos == "CB" || pos == "FS" || pos == "SS" || pos == "ATH" {
			score := 0
			if pos == "CB" {
				score += 100
			} else if pos == "FS" {
				score += 10
			} else if pos == "SS" {
				score += 8
			} else if pos == "ATH" && (arch == "Triple-Threat" || arch == "Bandit" || arch == "Weakside") {
				score += 50
			}

			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}

			score += int(float64(cp.Overall)*0.5) +
				int(float64(cp.Tackle)*0.05) +
				int(float64(cp.Agility)*0.1) +
				int(float64(cp.Catching)*0.1) +
				int(float64(cp.ManCoverage)*0.01) +
				int(float64(cp.ZoneCoverage)*0.01) +
				int(float64(cp.Speed)*0.05)

			dcpObj := structs.DepthChartPositionDTO{
				Position:      pos,
				Archetype:     arch,
				Score:         score,
				CollegePlayer: cp,
			}
			positionMap["CB"] = append(positionMap["CB"], dcpObj)
		}

		// Add to FS list
		if pos == "CB" || pos == "FS" || pos == "SS" || pos == "ATH" {
			score := 0
			if pos == "FS" {
				score += 100
			} else if pos == "CB" {
				score += 25
			} else if pos == "SS" {
				score += 12
			} else if pos == "ATH" && (arch == "Bandit" || arch == "Weakside") {
				score += 50
			}

			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}

			score += int(float64(cp.Overall)*0.5) +
				int(float64(cp.Tackle)*0.05) +
				int(float64(cp.Agility)*0.1) +
				int(float64(cp.Catching)*0.1) +
				int(float64(cp.ManCoverage)*0.01) +
				int(float64(cp.ZoneCoverage)*0.01) +
				int(float64(cp.Speed)*0.05)

			dcpObj := structs.DepthChartPositionDTO{
				Position:      pos,
				Archetype:     arch,
				Score:         score,
				CollegePlayer: cp,
			}
			positionMap["FS"] = append(positionMap["FS"], dcpObj)
		}

		// Add to SS list
		if pos == "CB" || pos == "FS" || pos == "SS" || pos == "ATH" {
			score := 0
			if pos == "SS" {
				score += 100
			} else if pos == "FS" {
				score += 25
			} else if pos == "CB" {
				score += 12
			} else if pos == "ATH" && (arch == "Bandit" || arch == "Weakside") {
				score += 50
			}

			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}

			score += int(float64(cp.Overall)*0.5) +
				int(float64(cp.Tackle)*0.05) +
				int(float64(cp.Agility)*0.1) +
				int(float64(cp.Catching)*0.1) +
				int(float64(cp.ManCoverage)*0.01) +
				int(float64(cp.ZoneCoverage)*0.01) +
				int(float64(cp.Speed)*0.05)

			dcpObj := structs.DepthChartPositionDTO{
				Position:      pos,
				Archetype:     arch,
				Score:         score,
				CollegePlayer: cp,
			}
			positionMap["SS"] = append(positionMap["SS"], dcpObj)
		}

		// Add to P list
		if pos == "K" || pos == "P" || pos == "QB" || pos == "ATH" {
			score := 0
			if pos == "P" {
				score += 100
			} else if pos == "K" {
				score += 25
			} else if pos == "ATH" && (arch == "Soccer Player") {
				score += 50
			}

			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}

			score += cp.PuntAccuracy + cp.PuntPower

			dcpObj := structs.DepthChartPositionDTO{
				Position:      pos,
				Archetype:     arch,
				Score:         score,
				CollegePlayer: cp,
			}
			positionMap["P"] = append(positionMap["P"], dcpObj)
		}

		// Add to K list (Field Goal)
		if pos == "K" || pos == "P" || pos == "QB" || pos == "ATH" {
			score := 0
			if pos == "K" {
				score += 100
			} else if pos == "P" {
				score += 25
			} else if pos == "ATH" && (arch == "Soccer Player") {
				score += 50
			}

			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}
			score += cp.KickAccuracy + cp.KickPower
			dcpObj := structs.DepthChartPositionDTO{
				Position:      pos,
				Archetype:     arch,
				Score:         score,
				CollegePlayer: cp,
			}
			positionMap["K"] = append(positionMap["K"], dcpObj)
		}

		// FG List
		if pos == "K" || pos == "P" || pos == "QB" || pos == "ATH" {
			score := 0
			if pos == "K" {
				score += 100
			} else if pos == "P" {
				score += 25
			} else if pos == "ATH" && (arch == "Soccer Player") {
				score += 50
			}

			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}

			score += cp.KickAccuracy + cp.KickPower

			dcpObj := structs.DepthChartPositionDTO{
				Position:      pos,
				Archetype:     arch,
				Score:         score,
				CollegePlayer: cp,
			}
			positionMap["FG"] = append(positionMap["FG"], dcpObj)
		}

		// PR
		if pos == "WR" || pos == "RB" || pos == "FS" || pos == "SS" || pos == "CB" || pos == "ATH" {
			score := 0
			if pos == "ATH" && arch == "Return Specialist" {
				score += 50
			} else if pos == "WR" || pos == "RB" {
				score += 25
			}
			score += cp.Agility

			dcpObj := structs.DepthChartPositionDTO{
				Position:      pos,
				Archetype:     arch,
				Score:         score,
				CollegePlayer: cp,
			}
			positionMap["PR"] = append(positionMap["PR"], dcpObj)
		}

		// KR
		if pos == "WR" || pos == "RB" || pos == "FS" || pos == "SS" || pos == "CB" || pos == "ATH" {
			score := 0
			if pos == "ATH" && arch == "Return Specialist" {
				score += 50
			} else if pos == "WR" || pos == "RB" {
				score += 25
			}
			score += cp.Speed

			dcpObj := structs.DepthChartPositionDTO{
				Position:      pos,
				Archetype:     arch,
				Score:         score,
				CollegePlayer: cp,
			}
			positionMap["KR"] = append(positionMap["KR"], dcpObj)
		}

		// STU
		if pos == "FB" || pos == "TE" || pos == "ILB" || pos == "OLB" || pos == "RB" || pos == "CB" || pos == "FS" || pos == "SS" || pos == "WR" || pos == "ATH" {
			score := 0
			if cp.Year == 2 || cp.Year == 1 {
				score += 50
			} else if cp.Year == 3 && cp.IsRedshirt {
				score += 25
			}

			score += cp.Tackle
			dcpObj := structs.DepthChartPositionDTO{
				Position:      pos,
				Archetype:     arch,
				Score:         score,
				CollegePlayer: cp,
			}
			positionMap["STU"] = append(positionMap["STU"], dcpObj)
		}
	}

	// Sort Each DC Position
	sort.Sort(structs.ByDCPosition(positionMap["QB"]))
	sort.Sort(structs.ByDCPosition(positionMap["RB"]))
	sort.Sort(structs.ByDCPosition(positionMap["FB"]))
	sort.Sort(structs.ByDCPosition(positionMap["WR"]))
	sort.Sort(structs.ByDCPosition(positionMap["TE"]))
	sort.Sort(structs.ByDCPosition(positionMap["LT"]))
	sort.Sort(structs.ByDCPosition(positionMap["RT"]))
	sort.Sort(structs.ByDCPosition(positionMap["LG"]))
	sort.Sort(structs.ByDCPosition(positionMap["RG"]))
	sort.Sort(structs.ByDCPosition(positionMap["C"]))
	sort.Sort(structs.ByDCPosition(positionMap["DT"]))
	sort.Sort(structs.ByDCPosition(positionMap["LE"]))
	sort.Sort(structs.ByDCPosition(positionMap["RE"]))
	sort.Sort(structs.ByDCPosition(positionMap["LOLB"]))
	sort.Sort(structs.ByDCPosition(positionMap["ROLB"]))
	sort.Sort(structs.ByDCPosition(positionMap["MLB"]))
	sort.Sort(structs.ByDCPosition(positionMap["CB"]))
	sort.Sort(structs.ByDCPosition(positionMap["FS"]))
	sort.Sort(structs.ByDCPosition(positionMap["SS"]))
	sort.Sort(structs.ByDCPosition(positionMap["P"]))
	sort.Sort(structs.ByDCPosition(positionMap["K"]))
	sort.Sort(structs.ByDCPosition(positionMap["PR"]))
	sort.Sort(structs.ByDCPosition(positionMap["KR"]))
	sort.Sort(structs.ByDCPosition(positionMap["FG"]))
	sort.Sort(structs.ByDCPosition(positionMap["STU"]))

	for _, dcp := range dcPositions {
		positionList := positionMap[dcp.Position]
		for _, pos := range positionList {
			if starterMap[pos.CollegePlayer.ID] &&
				dcp.Position != "FG" {
				continue
			}
			if backupMap[pos.CollegePlayer.ID] && dcp.PositionLevel != "1" && dcp.Position != "STU" {
				continue
			}
			if dcp.Position == "STU" && stuMap[pos.CollegePlayer.ID] {
				continue
			}

			if dcp.Position == "WR" {
				runnerDistPostition := gp.RunnerDistributionWRPosition
				positionLabel := dcp.Position + "" + dcp.PositionLevel
				if runnerDistPostition == positionLabel {
					gp.AssignRunnerWRID(dcp.CollegePlayer.ID)
				}
			}

			if dcp.Position == "STU" {
				stuMap[pos.CollegePlayer.ID] = true
			} else if dcp.PositionLevel == "1" && !starterMap[pos.CollegePlayer.ID] {
				starterMap[pos.CollegePlayer.ID] = true
			} else {
				backupMap[pos.CollegePlayer.ID] = true
			}
			dto := structs.CollegeDepthChartPositionTEST{
				DepthChartID:     dcp.DepthChartID,
				PlayerID:         int(pos.CollegePlayer.ID),
				FirstName:        pos.CollegePlayer.FirstName,
				LastName:         pos.CollegePlayer.LastName,
				OriginalPosition: pos.CollegePlayer.Position,
			}
			dto.AssignID(dcp.ID)
			dcp.UpdateDepthChartPosition(dto)
			db.Save(&dcp)
			break
		}
	}
}

func MigrateCFBGameplansAndDepthChartsForRemainingFCSTeams() {
	db := dbprovider.GetInstance().GetDB()

	teamPath := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\dc_positions_migration.csv"

	dcPositionsCSV := util.ReadCSV(teamPath)
	var gameplansList []structs.CollegeGameplan
	var testgameplansList []structs.CollegeGameplanTEST
	var dcList []structs.CollegeTeamDepthChart
	var testDCList []structs.CollegeTeamDepthChartTEST
	var dcPList []structs.CollegeDepthChartPosition
	var testDCPList []structs.CollegeDepthChartPositionTEST
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

func UpdateNFLAIDepthChartsTEST() {
	db := dbprovider.GetInstance().GetDB()
	teams := GetAllNFLTeams()
	for _, team := range teams {
		teamID := strconv.Itoa(int(team.ID))
		gp := repository.GetNFLGameplanTESTByTeamID(teamID)
		ReAlignNFLDepthChartTEST(db, teamID, gp)
	}

	ts := GetTimestamp()
	ts.ToggleAIDepthCharts()
	repository.SaveTimestamp(ts, db)
}

func MassUpdateNFLGameplanSchemesTEST(off, def string) {
	db := dbprovider.GetInstance().GetDB()
	teams := GetAllNFLTeams()
	offensiveSchemes := GetOffensiveDefaultSchemes()
	defensiveSchemes := GetDefensiveDefaultSchemes()
	for _, team := range teams {
		teamID := strconv.Itoa(int(team.ID))
		gp := repository.GetNFLGameplanTESTByTeamID(teamID)
		gp.UpdateSchemes(off, def)
		offFormations := offensiveSchemes[off]
		defFormations := defensiveSchemes[def][off]

		dto := structs.NFLGameplan{
			TeamID: team.ID,
			BaseGameplan: structs.BaseGameplan{
				OffensiveScheme:    off,
				DefensiveScheme:    def,
				OffensiveFormation: offFormations,
				DefensiveFormation: defFormations,
				BlitzSafeties:      gp.BlitzSafeties,
				BlitzCorners:       gp.BlitzCorners,
				LinebackerCoverage: gp.LinebackerCoverage,
				MaximumFGDistance:  gp.MaximumFGDistance,
				GoFor4AndShort:     gp.GoFor4AndShort,
				GoFor4AndLong:      gp.GoFor4AndLong,
				DefaultOffense:     gp.DefaultOffense,
				DefaultDefense:     gp.DefaultDefense,
				PrimaryHB:          75,
				PitchFocus:         50,
				DiveFocus:          50,
			},
		}

		gp.UpdateNFLGameplan(dto)

		// Autosort Depth Chart
		ReAlignNFLDepthChartTEST(db, teamID, gp)

		db.Save(&gp)
	}
}

func UpdateIndividualNFLGameplanSchemeTEST(teamID, off, def string) {
	db := dbprovider.GetInstance().GetDB()
	offensiveSchemes := GetOffensiveDefaultSchemes()
	defensiveSchemes := GetDefensiveDefaultSchemes()

	gp := repository.GetNFLGameplanTESTByTeamID(teamID)
	gp.UpdateSchemes(off, def)
	offFormations := offensiveSchemes[off]
	defFormations := defensiveSchemes[def][off]

	dto := structs.NFLGameplan{
		TeamID: gp.TeamID,
		BaseGameplan: structs.BaseGameplan{
			OffensiveScheme:    off,
			DefensiveScheme:    def,
			OffensiveFormation: offFormations,
			DefensiveFormation: defFormations,
			BlitzSafeties:      gp.BlitzSafeties,
			BlitzCorners:       gp.BlitzCorners,
			LinebackerCoverage: gp.LinebackerCoverage,
			MaximumFGDistance:  gp.MaximumFGDistance,
			GoFor4AndShort:     gp.GoFor4AndShort,
			GoFor4AndLong:      gp.GoFor4AndLong,
			DefaultOffense:     gp.DefaultOffense,
			DefaultDefense:     gp.DefaultDefense,
			PrimaryHB:          75,
			PitchFocus:         50,
			DiveFocus:          50,
		},
	}

	gp.UpdateNFLGameplan(dto)

	// Autosort Depth Chart
	// ReAlignNFLDepthChartTEST(db, teamID, gp)

	db.Save(&gp)
}

func ReAlignNFLDepthChartTEST(db *gorm.DB, teamID string, gp structs.NFLGameplanTEST) {
	roster := GetNFLPlayersWithContractsByTeamID(teamID)
	dcPositions := repository.GetNFLDepthChartPositionPlayersByDepthchartIDTEST(teamID)
	positionMap := make(map[string][]structs.DepthChartPositionDTO)
	starterMap := make(map[uint]bool)
	backupMap := make(map[uint]bool)
	stuMap := make(map[uint]bool)
	offScheme := gp.OffensiveScheme
	defScheme := gp.DefensiveScheme
	isLT := true
	isLG := true
	isLE := true
	isLOLB := true

	goodFits := GetFitsByScheme(offScheme, false)
	badFits := GetFitsByScheme(defScheme, false)
	bonus := 5

	// Allocate the Position Map
	for _, cp := range roster {
		if cp.IsInjured || cp.IsPracticeSquad || cp.WeeksOfRecovery > 0 {
			continue
		}
		pos := cp.Position
		arch := cp.Archetype
		player := arch + " " + pos
		isGoodFit := CheckPlayerFits(player, goodFits)
		isBadFit := CheckPlayerFits(player, badFits)

		// Add to QB List
		if pos == "QB" || pos == "WR" || pos == "TE" || pos == "RB" || pos == "FB" || pos == "K" || pos == "P" {
			score := 0
			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}
			if pos == "QB" {
				score += 75
			} else if pos == "ATH" && (arch == "Triple-Threat" || arch == "Field General") {
				score += 50
			}
			score += cp.Overall

			dcpObj := structs.DepthChartPositionDTO{
				Position:  pos,
				Archetype: arch,
				Score:     score,
				NFLPlayer: cp,
			}
			positionMap["QB"] = append(positionMap["QB"], dcpObj)
		}
		// Add to RB List
		if pos == "RB" || pos == "FB" || pos == "WR" || pos == "TE" {
			score := 0
			if pos == "RB" {
				score += 100
			} else if pos == "FB" {
				score += 25
			} else if pos == "ATH" && (arch == "Wingback" || arch == "Soccer Player" || arch == "Triple-Threat") {
				score += 50
			}
			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}

			score += ((cp.Speed + cp.Agility + cp.Strength + cp.Carrying) / 4)
			dcpObj := structs.DepthChartPositionDTO{
				Position:  pos,
				Archetype: arch,
				Score:     score,
				NFLPlayer: cp,
			}
			positionMap["RB"] = append(positionMap["RB"], dcpObj)
		}

		// Add to FB List
		if pos == "FB" || pos == "TE" || pos == "RB" || pos == "ILB" || pos == "OLB" {
			score := 0
			if pos == "FB" {
				score += 100
			} else if pos == "ATH" && (arch == "Wingback") {
				score += 50
			}

			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}
			score += ((cp.Strength + cp.Carrying + cp.PassBlock + cp.RunBlock) / 4)

			dcpObj := structs.DepthChartPositionDTO{
				Position:  pos,
				Archetype: arch,
				Score:     score,
				NFLPlayer: cp,
			}
			positionMap["FB"] = append(positionMap["FB"], dcpObj)
		}

		// Add to TE List
		if pos == "FB" || pos == "TE" || pos == "WR" {
			score := 0
			if pos == "TE" {
				score += 100
			} else if pos == "ATH" && (arch == "Slotback") {
				score += 50
			}

			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}

			score += int(float64(cp.Overall)*0.5) + int(float64(cp.RunBlock)*0.125) + int(float64(cp.PassBlock)*0.125) + int(float64(cp.Catching)*0.125) + int(float64(cp.Strength)*0.125)

			dcpObj := structs.DepthChartPositionDTO{
				Position:  pos,
				Archetype: arch,
				Score:     score,
				NFLPlayer: cp,
			}
			positionMap["TE"] = append(positionMap["TE"], dcpObj)
		}
		// Add to WR List
		if pos == "WR" || pos == "TE" || pos == "RB" || pos == "CB" {
			score := 0
			if pos == "WR" {
				score += 100
			} else if pos == "ATH" && (arch == "Wingback" || arch == "Slotback") {
				score += 50
			}

			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}

			score += int(float64(cp.Overall)*0.4) +
				int(float64(cp.Speed)*0.12) +
				int(float64(cp.Agility)*0.12) +
				int(float64(cp.Catching)*0.12) +
				int(float64(cp.Strength)*0.12) +
				int(float64(cp.RouteRunning)*0.12)

			dcpObj := structs.DepthChartPositionDTO{
				Position:  pos,
				Archetype: arch,
				Score:     score,
				NFLPlayer: cp,
			}
			positionMap["WR"] = append(positionMap["WR"], dcpObj)
		}
		// Add to LT and RT List
		if pos == "OT" || pos == "OG" || pos == "C" {
			score := 0
			if pos == "OT" {
				score += 100
			} else if pos == "OG" {
				score += 15
			} else if pos == "ATH" && (arch == "Lineman") {
				score += 50
			}

			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}

			score += int(float64(cp.Overall)*0.7) +
				int(float64(cp.Strength)*0.10) +
				int(float64(cp.RunBlock)*0.75) +
				int(float64(cp.PassBlock)*0.75) +
				int(float64(cp.Agility)*0.05)

			dcpObj := structs.DepthChartPositionDTO{
				Position:  pos,
				Archetype: arch,
				Score:     score,
				NFLPlayer: cp,
			}
			if isLT {
				positionMap["LT"] = append(positionMap["LT"], dcpObj)
			} else {
				positionMap["RT"] = append(positionMap["RT"], dcpObj)
			}
			isLT = !isLT
		}
		// Add to LG and RG List
		if pos == "OT" || pos == "OG" || pos == "C" {
			score := 0
			if pos == "OG" {
				score += 100
			} else if pos == "C" {
				score += 25
			} else if pos == "ATH" && (arch == "Lineman") {
				score += 50
			}

			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}

			score += int(float64(cp.Overall)*0.7) +
				int(float64(cp.Strength)*0.10) +
				int(float64(cp.RunBlock)*0.75) +
				int(float64(cp.PassBlock)*0.75) +
				int(float64(cp.Agility)*0.05)
			dcpObj := structs.DepthChartPositionDTO{
				Position:  pos,
				Archetype: arch,
				Score:     score,
				NFLPlayer: cp,
			}
			if isLG {
				positionMap["LG"] = append(positionMap["LG"], dcpObj)
			} else {
				positionMap["RG"] = append(positionMap["RG"], dcpObj)
			}
			isLG = !isLG
		}
		// Add to C List
		if pos == "OT" || pos == "OG" || pos == "C" {
			score := 0
			if pos == "C" {
				score += 100
			} else if pos == "OG" {
				score += 15
			} else if pos == "ATH" && (arch == "Lineman") {
				score += 50
			}

			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}

			score += int(float64(cp.Overall)*0.7) +
				int(float64(cp.Strength)*0.10) +
				int(float64(cp.RunBlock)*0.75) +
				int(float64(cp.PassBlock)*0.75) +
				int(float64(cp.Agility)*0.05)

			dcpObj := structs.DepthChartPositionDTO{
				Position:  pos,
				Archetype: arch,
				Score:     score,
				NFLPlayer: cp,
			}
			positionMap["C"] = append(positionMap["C"], dcpObj)
		}

		// Add to LE List
		if pos == "DE" || pos == "DT" || pos == "OLB" {
			score := 0
			if pos == "DE" {
				score += 100
			} else if pos == "OLB" {
				score += 25
			} else if pos == "DT" {
				score += 3
			} else if pos == "ATH" && (arch == "Lineman" || arch == "Strongside") {
				score += 50
			}

			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}

			score += int(float64(cp.Overall)*0.7) +
				int(float64(cp.Strength)*0.05) +
				int(float64(cp.Tackle)*0.05) +
				int(float64(cp.PassRush)*0.75) +
				int(float64(cp.RunDefense)*0.75) +
				int(float64(cp.Agility)*0.05)

			dcpObj := structs.DepthChartPositionDTO{
				Position:  pos,
				Archetype: arch,
				Score:     score,
				NFLPlayer: cp,
			}
			if isLE {
				positionMap["LE"] = append(positionMap["LE"], dcpObj)
			} else {
				positionMap["RE"] = append(positionMap["RE"], dcpObj)
			}
			isLE = !isLE
		}

		// Add to DT list
		if pos == "DE" || pos == "DT" || pos == "OLB" {
			score := 0
			if pos == "DT" {
				score += 100
			} else if pos == "DE" {
				score += 25
			} else if pos == "OLB" {
				score += 12
			} else if pos == "ATH" && (arch == "Lineman" || arch == "Strongside") {
				score += 50
			}

			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}

			score += int(float64(cp.Overall)*0.7) +
				int(float64(cp.Strength)*0.05) +
				int(float64(cp.Tackle)*0.05) +
				int(float64(cp.PassRush)*0.75) +
				int(float64(cp.RunDefense)*0.75) +
				int(float64(cp.Agility)*0.05)

			dcpObj := structs.DepthChartPositionDTO{
				Position:  pos,
				Archetype: arch,
				Score:     score,
				NFLPlayer: cp,
			}
			positionMap["DT"] = append(positionMap["DT"], dcpObj)
		}

		// Add to OLB list
		if pos == "OLB" || pos == "DE" || pos == "ILB" || pos == "SS" || pos == "FS" {
			score := 0
			if pos == "OLB" {
				score += 100
			} else if pos == "DE" {
				score += 10
			} else if pos == "ILB" {
				score += 25
			} else if pos == "SS" {
				score += 3
			} else if pos == "ATH" && (arch == "Weakside" || arch == "Strongside" || arch == "Bandit") {
				score += 50
			}

			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}

			score += int(float64(cp.Overall)*0.6) +
				int(float64(cp.Strength)*0.025) +
				int(float64(cp.Tackle)*0.055) +
				int(float64(cp.PassRush)*0.0755) +
				int(float64(cp.RunDefense)*0.0755) +
				int(float64(cp.ManCoverage)*0.075) +
				int(float64(cp.ZoneCoverage)*0.075) +
				int(float64(cp.Agility)*0.025)

			dcpObj := structs.DepthChartPositionDTO{
				Position:  pos,
				Archetype: arch,
				Score:     score,
				NFLPlayer: cp,
			}
			if isLOLB {
				positionMap["LOLB"] = append(positionMap["LOLB"], dcpObj)
			} else {
				positionMap["ROLB"] = append(positionMap["ROLB"], dcpObj)
			}
			isLOLB = !isLOLB
		}

		// Add to ILB list
		if pos == "OLB" || pos == "DE" || pos == "ILB" || pos == "SS" || pos == "FS" {
			score := 0
			if pos == "ILB" {
				score += 100
			} else if pos == "OLB" {
				score += 25
			} else if pos == "SS" {
				score += 8
			} else if pos == "DE" {
				score += 3
			} else if pos == "ATH" && (arch == "Weakside" || arch == "Bandit" || arch == "Field General") {
				score += 50
			}

			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}

			score += int(float64(cp.Overall)*0.6) +
				int(float64(cp.Strength)*0.025) +
				int(float64(cp.Tackle)*0.055) +
				int(float64(cp.PassRush)*0.0755) +
				int(float64(cp.RunDefense)*0.0755) +
				int(float64(cp.ManCoverage)*0.075) +
				int(float64(cp.ZoneCoverage)*0.075) +
				int(float64(cp.Agility)*0.025)

			dcpObj := structs.DepthChartPositionDTO{
				Position:  pos,
				Archetype: arch,
				Score:     score,
				NFLPlayer: cp,
			}
			positionMap["ILB"] = append(positionMap["ILB"], dcpObj)
		}

		// Add to CB List
		if pos == "CB" || pos == "FS" || pos == "SS" {
			score := 0
			if pos == "CB" {
				score += 100
			} else if pos == "FS" {
				score += 10
			} else if pos == "SS" {
				score += 8
			} else if pos == "ATH" && (arch == "Triple-Threat" || arch == "Bandit" || arch == "Weakside") {
				score += 50
			}

			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}

			score += int(float64(cp.Overall)*0.5) +
				int(float64(cp.Tackle)*0.05) +
				int(float64(cp.Agility)*0.1) +
				int(float64(cp.Catching)*0.1) +
				int(float64(cp.ManCoverage)*0.01) +
				int(float64(cp.ZoneCoverage)*0.01) +
				int(float64(cp.Speed)*0.05)

			dcpObj := structs.DepthChartPositionDTO{
				Position:  pos,
				Archetype: arch,
				Score:     score,
				NFLPlayer: cp,
			}
			positionMap["CB"] = append(positionMap["CB"], dcpObj)
		}

		// Add to FS list
		if pos == "CB" || pos == "FS" || pos == "SS" {
			score := 0
			if pos == "FS" {
				score += 100
			} else if pos == "CB" {
				score += 25
			} else if pos == "SS" {
				score += 12
			} else if pos == "ATH" && (arch == "Bandit" || arch == "Weakside") {
				score += 50
			}

			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}

			score += int(float64(cp.Overall)*0.5) +
				int(float64(cp.Tackle)*0.05) +
				int(float64(cp.Agility)*0.1) +
				int(float64(cp.Catching)*0.1) +
				int(float64(cp.ManCoverage)*0.01) +
				int(float64(cp.ZoneCoverage)*0.01) +
				int(float64(cp.Speed)*0.05)

			dcpObj := structs.DepthChartPositionDTO{
				Position:  pos,
				Archetype: arch,
				Score:     score,
				NFLPlayer: cp,
			}
			positionMap["FS"] = append(positionMap["FS"], dcpObj)
		}

		// Add to SS list
		if pos == "CB" || pos == "FS" || pos == "SS" {
			score := 0
			if pos == "SS" {
				score += 100
			} else if pos == "FS" {
				score += 25
			} else if pos == "CB" {
				score += 12
			} else if pos == "ATH" && (arch == "Bandit" || arch == "Weakside") {
				score += 50
			}

			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}

			score += int(float64(cp.Overall)*0.5) +
				int(float64(cp.Tackle)*0.05) +
				int(float64(cp.Agility)*0.1) +
				int(float64(cp.Catching)*0.1) +
				int(float64(cp.ManCoverage)*0.01) +
				int(float64(cp.ZoneCoverage)*0.01) +
				int(float64(cp.Speed)*0.05)

			dcpObj := structs.DepthChartPositionDTO{
				Position:  pos,
				Archetype: arch,
				Score:     score,
				NFLPlayer: cp,
			}
			positionMap["SS"] = append(positionMap["SS"], dcpObj)
		}

		// Add to P list
		if pos == "K" || pos == "P" || pos == "QB" {
			score := 0
			if pos == "P" {
				score += 100
			} else if pos == "K" {
				score += 25
			} else if pos == "ATH" && (arch == "Soccer Player") {
				score += 50
			}

			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}

			score += cp.PuntAccuracy + cp.PuntPower

			dcpObj := structs.DepthChartPositionDTO{
				Position:  pos,
				Archetype: arch,
				Score:     score,
				NFLPlayer: cp,
			}
			positionMap["P"] = append(positionMap["P"], dcpObj)
		}

		// Add to K list (Field Goal)
		if pos == "K" || pos == "P" || pos == "QB" {
			score := 0
			if pos == "K" {
				score += 100
			} else if pos == "P" {
				score += 25
			} else if pos == "ATH" && (arch == "Soccer Player") {
				score += 50
			}

			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}
			score += cp.KickAccuracy + cp.KickPower

			dcpObj := structs.DepthChartPositionDTO{
				Position:  pos,
				Archetype: arch,
				Score:     score,
				NFLPlayer: cp,
			}
			positionMap["K"] = append(positionMap["K"], dcpObj)
		}

		// FG List
		if pos == "K" || pos == "P" || pos == "QB" {
			score := 0
			if pos == "K" {
				score += 100
			} else if pos == "P" {
				score += 25
			} else if pos == "ATH" && (arch == "Soccer Player") {
				score += 50
			}

			if isGoodFit && !isBadFit {
				score += bonus
			} else if isBadFit && !isGoodFit {
				score -= bonus
			}

			score += cp.KickAccuracy + cp.KickPower

			dcpObj := structs.DepthChartPositionDTO{
				Position:  pos,
				Archetype: arch,
				Score:     score,
				NFLPlayer: cp,
			}
			positionMap["FG"] = append(positionMap["FG"], dcpObj)
		}

		// PR
		if pos == "WR" || pos == "RB" || pos == "FS" || pos == "SS" || pos == "CB" {
			score := 0
			if pos == "ATH" && arch == "Return Specialist" {
				score += 50
			} else if pos == "WR" || pos == "RB" {
				score += 25
			}
			score += cp.Agility

			dcpObj := structs.DepthChartPositionDTO{
				Position:  pos,
				Archetype: arch,
				Score:     score,
				NFLPlayer: cp,
			}
			positionMap["PR"] = append(positionMap["PR"], dcpObj)
		}

		// KR
		if pos == "WR" || pos == "RB" || pos == "FS" || pos == "SS" || pos == "CB" {
			score := 0
			if pos == "ATH" && arch == "Return Specialist" {
				score += 50
			} else if pos == "WR" || pos == "RB" {
				score += 25
			}
			score += cp.Speed

			dcpObj := structs.DepthChartPositionDTO{
				Position:  pos,
				Archetype: arch,
				Score:     score,
				NFLPlayer: cp,
			}
			positionMap["KR"] = append(positionMap["KR"], dcpObj)
		}

		// STU
		if pos == "FB" || pos == "TE" || pos == "ILB" || pos == "OLB" || pos == "RB" || pos == "CB" || pos == "FS" || pos == "SS" {
			score := 0
			switch cp.Experience {
			case 2:
				score += 50
			case 1:
				score += 45
			case 3:
				score += 15
			}

			score += cp.Tackle
			dcpObj := structs.DepthChartPositionDTO{
				Position:  pos,
				Archetype: arch,
				Score:     score,
				NFLPlayer: cp,
			}
			positionMap["STU"] = append(positionMap["STU"], dcpObj)
		}
	}

	// Sort Each DC Position
	sort.Sort(structs.ByDCPosition(positionMap["QB"]))
	sort.Sort(structs.ByDCPosition(positionMap["RB"]))
	sort.Sort(structs.ByDCPosition(positionMap["FB"]))
	sort.Sort(structs.ByDCPosition(positionMap["WR"]))
	sort.Sort(structs.ByDCPosition(positionMap["TE"]))
	sort.Sort(structs.ByDCPosition(positionMap["LT"]))
	sort.Sort(structs.ByDCPosition(positionMap["RT"]))
	sort.Sort(structs.ByDCPosition(positionMap["LG"]))
	sort.Sort(structs.ByDCPosition(positionMap["RG"]))
	sort.Sort(structs.ByDCPosition(positionMap["C"]))
	sort.Sort(structs.ByDCPosition(positionMap["DT"]))
	sort.Sort(structs.ByDCPosition(positionMap["LE"]))
	sort.Sort(structs.ByDCPosition(positionMap["RE"]))
	sort.Sort(structs.ByDCPosition(positionMap["LOLB"]))
	sort.Sort(structs.ByDCPosition(positionMap["ROLB"]))
	sort.Sort(structs.ByDCPosition(positionMap["ILB"]))
	sort.Sort(structs.ByDCPosition(positionMap["CB"]))
	sort.Sort(structs.ByDCPosition(positionMap["FS"]))
	sort.Sort(structs.ByDCPosition(positionMap["SS"]))
	sort.Sort(structs.ByDCPosition(positionMap["P"]))
	sort.Sort(structs.ByDCPosition(positionMap["K"]))
	sort.Sort(structs.ByDCPosition(positionMap["PR"]))
	sort.Sort(structs.ByDCPosition(positionMap["KR"]))
	sort.Sort(structs.ByDCPosition(positionMap["FG"]))
	sort.Sort(structs.ByDCPosition(positionMap["STU"]))

	for _, dcp := range dcPositions {
		positionList := positionMap[dcp.Position]
		for _, pos := range positionList {
			if starterMap[pos.NFLPlayer.ID] &&
				dcp.Position != "FG" {
				continue
			}
			if backupMap[pos.NFLPlayer.ID] && dcp.PositionLevel != "1" && dcp.Position != "STU" {
				continue
			}
			if dcp.Position == "STU" && stuMap[pos.NFLPlayer.ID] {
				continue
			}

			if dcp.Position == "WR" {
				runnerDistPostition := gp.RunnerDistributionWRPosition
				positionLabel := dcp.Position + "" + dcp.PositionLevel
				if runnerDistPostition == positionLabel {
					gp.AssignRunnerWRID(dcp.NFLPlayer.ID)
				}
			}

			if dcp.Position == "STU" {
				stuMap[pos.NFLPlayer.ID] = true
			} else if dcp.PositionLevel == "1" && !starterMap[pos.NFLPlayer.ID] {
				starterMap[pos.NFLPlayer.ID] = true
			} else {
				backupMap[pos.NFLPlayer.ID] = true
			}
			dto := structs.NFLDepthChartPositionTEST{
				DepthChartID:     dcp.DepthChartID,
				PlayerID:         pos.NFLPlayer.ID,
				FirstName:        pos.NFLPlayer.FirstName,
				LastName:         pos.NFLPlayer.LastName,
				OriginalPosition: pos.NFLPlayer.Position,
			}
			dto.AssignID(dcp.ID)
			dcp.UpdateDepthChartPosition(dto)
			db.Save(&dcp)
			break
		}
	}

	db.Save(&gp)
}

func SetupNFLTestDataStructs() {
	db := dbprovider.GetInstance().GetDB()

	nflTeams := GetAllNFLTeams()
	var testgameplansList []structs.NFLGameplanTEST
	var testDCList []structs.NFLDepthChartTEST
	var testDCPList []structs.NFLDepthChartPositionTEST

	for _, team := range nflTeams {
		teamID := strconv.Itoa(int(team.ID))
		gp := GetNFLGameplanByTeamID(teamID)
		dc := GetNFLDepthchartByTeamID(teamID)
		dcp := GetNFLDepthChartPositionsByDepthchartID(teamID)

		testGP := structs.NFLGameplanTEST{
			TeamID:       team.ID,
			BaseGameplan: gp.BaseGameplan,
		}
		testgameplansList = append(testgameplansList, testGP)
		testDC := structs.NFLDepthChartTEST{
			TeamID: dc.TeamID,
		}
		testDCList = append(testDCList, testDC)
		for _, p := range dcp {
			dcpt := structs.NFLDepthChartPositionTEST{
				DepthChartID:     p.DepthChartID,
				PlayerID:         p.PlayerID,
				Position:         p.Position,
				PositionLevel:    p.PositionLevel,
				OriginalPosition: p.OriginalPosition,
				FirstName:        p.FirstName,
				LastName:         p.LastName,
			}
			testDCPList = append(testDCPList, dcpt)
		}
	}

	repository.CreateNFLGameplansTESTRecordsBatch(db, testgameplansList, 200)
	repository.CreateNFLDepthChartTESTBatch(db, testDCList, 200)
	repository.CreateNFLDepthChartPositionsTESTBatch(db, testDCPList, 2000)
}
