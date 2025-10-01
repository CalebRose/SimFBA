package managers

import (
	"bufio"
	"encoding/csv"
	"errors"
	"math/rand/v2"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
	"gorm.io/gorm"
)

func RunTrainingCamps(year string) error {
	db := dbprovider.GetInstance().GetDB()

	readPath := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\" + year + "\\trainingcamp.csv"
	writePath := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\" + year + "\\trainingcamp_results.csv"

	_, err := os.Stat(readPath)
	if errors.Is(err, os.ErrNotExist) {
		return errors.New("Training camp CSV not found for year " + year)
	} else if err != nil {
		return errors.New("Error checking for training camp CSV for " + year + ": " + err.Error())
	}

	_, err = os.Stat(writePath)
	if err == nil {
		return errors.New("Training camp output CSV already exists for " + year + ". Training camp may have already applied.")
	} else if !errors.Is(err, os.ErrNotExist) {
		return errors.New("Something weird happened. Training camp output may already exist for " + year + ", but an unrelated error occurred while checking for it: " + err.Error())
	}

	drillSelectionsCSV := util.ReadCSV(readPath)

	drillResultsCSV, err := os.Create(writePath)
	if err != nil {
		panic(err)
	}

	csvWriter := csv.NewWriter(bufio.NewWriter(drillResultsCSV))

	csvWriter.Write([]string{"PlayerID", "Team", "DrillPosition", "Archetype", "FirstName", "LastName", "Age", "PositionDrill", "PositionDrillAttribute", "PositionDrillResult", "TeamDrill", "TeamDrillAttribute",
		"TeamDrillResult", "EventText", "InjuryText", "WeeksOut"})

	defer drillResultsCSV.Close()
	defer csvWriter.Flush()

	for idx, row := range drillSelectionsCSV {
		if idx == 0 {
			continue
		}

		teamId := row[0]

		players := GetNFLPlayersWithContractsByTeamID(teamId)

		positionOverrides := getPositionOverrides(strings.ToLower(row[13]))

		for _, player := range players {
			drillPosition := player.Position
			drillArchetype := player.Archetype
			if slices.Contains(positionOverrides, strconv.Itoa(player.PlayerID)) || player.Position == "ATH" {
				drillPosition = player.PositionTwo
				drillArchetype = player.ArchetypeTwo
			}

			positionDrill := getPositionDrill(drillPosition, row, player)
			teamDrill := row[12]

			runDrills(player, drillPosition, drillArchetype, positionDrill, teamDrill, csvWriter, db)
		}
	}
	return nil
}

func runDrills(player structs.NFLPlayer, drillPosition string, drillArchetype string, positionDrill string, teamDrill string, csvWriter *csv.Writer, db *gorm.DB) {
	positionDrillAttribute := getAttribute(drillPosition, drillArchetype, positionDrill)
	teamDrillAttribute := getAttribute(drillPosition, drillArchetype, teamDrill)

	eventModifier := getEventModifier(player)
	eventText := ""
	if eventModifier != 0 {
		eventText = events[0][eventModifier][rand.IntN(len(events[0][eventModifier]))]
	}

	injuryText, injuryWeeks := checkInjury(player)

	changedAttrs := &structs.CollegePlayerProgressions{
		FootballIQ:      player.FootballIQ,
		Speed:           player.Speed,
		Carrying:        player.Carrying,
		Agility:         player.Agility,
		Catching:        player.Catching,
		RouteRunning:    player.RouteRunning,
		ZoneCoverage:    player.ZoneCoverage,
		ManCoverage:     player.ManCoverage,
		Strength:        player.Strength,
		Tackle:          player.Tackle,
		PassBlock:       player.PassBlock,
		RunBlock:        player.RunBlock,
		PassRush:        player.PassRush,
		RunDefense:      player.RunDefense,
		ThrowPower:      player.ThrowPower,
		ThrowAccuracy:   player.ThrowAccuracy,
		KickAccuracy:    player.KickAccuracy,
		KickPower:       player.KickPower,
		PuntAccuracy:    player.PuntAccuracy,
		PuntPower:       player.PuntPower,
		InjuryText:      injuryText,
		WeeksOfRecovery: injuryWeeks,
	}

	positionDrillResult := getDrillResult(player, eventModifier)
	teamDrillResult := getDrillResult(player, eventModifier)

	applyDrillResult(changedAttrs, positionDrillAttribute, positionDrillResult)
	applyDrillResult(changedAttrs, teamDrillAttribute, teamDrillResult)

	csvWriter.Write([]string{strconv.Itoa(player.PlayerID), player.TeamAbbr, drillPosition, drillArchetype, player.FirstName, player.LastName, strconv.Itoa(player.Age), positionDrill, positionDrillAttribute,
		strconv.Itoa(positionDrillResult), teamDrill, teamDrillAttribute, strconv.Itoa(teamDrillResult), eventText, injuryText,
		strconv.Itoa(injuryWeeks)},
	)
	player.ApplyTrainingCampInfo(*changedAttrs)
	player.GetOverall()

	repository.SaveNFLPlayer(player, db)
}

func getPositionOverrides(overrides string) []string {
	if overrides == "" {
		return []string{}
	}
	// Player names should ideally be formatted as FirstnameLastname and spaces in between, so they can be treated as part of the same column of the CSV.
	return strings.Split(overrides, " ")
}

func getPositionDrill(drillPosition string, row []string, player structs.NFLPlayer) string {
	switch drillPosition {
	case "QB":
		return row[1]
	case "RB":
		return row[2]
	case "FB":
		return row[3]
	case "TE":
		return row[4]
	case "WR":
		return row[5]
	case "OT", "OG", "C":
		return row[6]
	case "DT", "DE":
		return row[7]
	case "ILB":
		return row[8]
	case "OLB":
		if player.Archetype == "Pass Rush" {
			return row[7]
		} else {
			return row[8]
		}
	case "CB":
		return row[9]
	case "FS", "SS":
		return row[10]
	case "K", "P":
		return row[11]
	default:
		return "tackle"
	}
}

func getEventModifier(player structs.NFLPlayer) int {
	discipline := float32(player.Discipline)

	negative := (-.6 * discipline) + 60
	positive := (.6 * discipline) + negative

	// Older veterans are more acclimated to the NFL and are less likely to have camp events, positive or negative.
	if player.Age > 24 {
		negative = negative / 2
		positive = positive / 2
	} else if player.Age > 27 {
		negative = negative / 4
		positive = positive / 4
	}

	eventRoll := float32(rand.IntN(100))
	if eventRoll < negative {
		eventSeverity := rand.IntN(100)
		if eventSeverity < 60 {
			return -1
		} else if eventSeverity < 90 {
			return -2
		} else {
			return -3
		}
	} else if eventRoll < positive {
		eventSeverity := rand.IntN(100)
		if eventSeverity < 60 {
			return 1
		} else if eventSeverity < 90 {
			return 2
		} else {
			return 3
		}
	}
	// no event
	return 0
}

// Build global dictionary of events
var events = []map[int][]string{
	{
		3: {"Dominates in drills, consistently outperforming expectations.",
			"Takes the lead in drills, setting the pace for others to follow.",
			"Emerges as an early standout in minicamp reports.",
			"Impresses coaching staff with natural leadership skills during team huddles.",
			"Consistently executes plays at a high level, earning trust from the coaching staff."},
		2: {"Receives praise in team meetings for attention to detail.",
			"Quickly picks up the playbook and shows understanding in practice.",
			"Earns a shout-out in a press conference from the head coach.",
			"Demonstrates unexpected versatility by excelling in an unfamiliar role.",
			"Sets a strong example for fellow rookies with a positive attitude and work ethic."},
		1: {"Impresses position coach with consistent effort.",
			"Completes all conditioning drills without issue.",
			"Makes a solid play during a scrimmage.",
			"Responds well to coaching, quickly implementing feedback.",
			"Consistently shows good sportsmanship and camaraderie with teammates."},
		-1: {"Shows up late to a team meeting.",
			"Struggles with conditioning drills.",
			"Misses a minor assignment in a scrimmage and gets yelled at by a coach over it.",
			"Gets unusually frustrated after losing a rep in individual drills.",
			"Position coaches completly forget his name."},
		-2: {"Blows a key play during a scrimmage in front of the coaching staff.",
			"Gets called out in team meetings for lack of focus.",
			"Struggles with communication and timing on the field.",
			"Repeatedly forgets assignments, slowing down drills for others.",
			"Performs poorly in conditioning, noticeably lagging behind teammates."},
		-3: {"Misses multiple meetings or practices without an excuse, causing a formal warning from the team.",
			"Has an on-field meltdown during a scrimmage, leading to being pulled from practice.",
			"Publicly criticizes coaching decisions during interviews, damaging relationships with the staff.",
			"Consistently fails to execute assignments, prompting questions about future roster status.",
			"Demonstrates reckless behavior off the field, sparking disciplinary actions."},
	},
}

func checkInjury(player structs.NFLPlayer) (string, int) {
	injuryCheck := rand.IntN(1000)
	if injuryCheck < 25 {
		return getInjuryDetails(player)
	}
	return "None", 0
}

func getInjuryDetails(player structs.NFLPlayer) (string, int) {
	injuryRoll := getInjuryRoll()
	weeksOut := 0

	if injuryRoll == 10 {
		weeksOut = rand.IntN(20) + 1
	} else {
		// results in a range from 0-15, spread roughly evenly across injury ratings from 0-100
		injuryModifier := int((100.0 - float32(player.Injury)) / 6.67)
		weeksOut = injuryTimes[injuryRoll][15-injuryModifier]
	}

	severity := getInjurySeverity(weeksOut)

	return getInjuryText(severity), weeksOut
}

func getInjuryRoll() int {
	roll := rand.IntN(1000)

	if roll < 21 {
		return 0
	} else if roll < 83 {
		return 1
	} else if roll < 166 {
		return 2
	} else if roll < 270 {
		return 3
	} else if roll < 416 {
		return 4
	} else if roll < 583 {
		return 5
	} else if roll < 729 {
		return 6
	} else if roll < 833 {
		return 7
	} else if roll < 916 {
		return 8
	} else if roll < 978 {
		return 9
	} else {
		return 10
	}
}

var injuryTimes = [][16]int{
	{15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0},
	{13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 0, 0},
	{11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 1, 0, 0, 0, 0},
	{9, 8, 7, 6, 5, 4, 3, 2, 1, 1, 0, 0, 0, 0, 0, 0},
	{7, 7, 6, 5, 5, 4, 3, 2, 1, 1, 0, 0, 0, 0, 0, 0},
	{5, 5, 4, 4, 3, 3, 2, 2, 1, 0, 0, 0, 0, 0, 0, 0},
	{6, 6, 5, 4, 4, 3, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0},
	{8, 7, 6, 5, 4, 3, 2, 1, 1, 0, 0, 0, 0, 0, 0, 0},
	{10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 1, 0, 0, 0, 0, 0},
	{12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 1, 0, 0, 0},
}

func getInjurySeverity(weeksOut int) string {
	if weeksOut < 3 {
		return "Minor"
	} else if weeksOut < 7 {
		return "Moderate"
	} else if weeksOut < 13 {
		return "Severe"
	} else {
		return "Season Ending"
	}
}

func getInjuryText(severity string) string {
	switch severity {
	case "Minor":
		return minorInjuries[rand.IntN(len(minorInjuries))]
	case "Moderate":
		return moderateInjuries[rand.IntN(len(moderateInjuries))]
	case "Severe":
		return severeInjuries[rand.IntN(len(severeInjuries))]
	case "Season Ending":
		return seasonEndingInjuries[rand.IntN(len(seasonEndingInjuries))]
	}
	return "Unknown Injury"
}

func getDrillResult(player structs.NFLPlayer, eventModifier int) int {
	badDrillThreshold := 5
	neutralDrillThreshold := 15
	okayDrillThreshold := 50
	goodDrillThreshold := 80

	// Older veterans are less likely to have major breakthroughs at camp
	if player.Age > 27 {
		badDrillThreshold = 5
		neutralDrillThreshold = 40
		okayDrillThreshold = 80
		goodDrillThreshold = 95
	}

	result := 0
	drillRoll := rand.IntN(100)
	if drillRoll < badDrillThreshold {
		result = -1
	} else if drillRoll < neutralDrillThreshold {
		result = 0
	} else if drillRoll < okayDrillThreshold {
		result = 1
	} else if drillRoll < goodDrillThreshold {
		result = 2
	} else {
		// great drill result
		result = 3
	}

	return result + eventModifier
}

func applyDrillResult(progression *structs.CollegePlayerProgressions, attribute string, modifier int) {
	if attribute == "football_iq" {
		if progression.FootballIQ + modifier > 99 {
			progression.FootballIQ = 99
		} else {
			progression.FootballIQ += modifier
		}
	}
	if attribute == "speed" {
		if progression.Speed + modifier > 99 {
			progression.Speed = 99
		} else {
			progression.Speed += modifier
		}
	}
	if attribute == "carrying" {
		if progression.Carrying + modifier > 99 {
			progression.Carrying = 99
		} else {
			progression.Carrying += modifier
		}
	}
	if attribute == "agility" {
		if progression.Agility + modifier > 99 {
			progression.Agility = 99
		} else {
			progression.Agility += modifier
		}
	}
	if attribute == "catching" {
		if progression.Catching + modifier > 99 {
			progression.Catching = 99
		} else {
			progression.Catching += modifier
		}
	}
	if attribute == "route_running" {
		if progression.RouteRunning + modifier > 99 {
			progression.RouteRunning = 99
		} else {
			progression.RouteRunning += modifier
		}
	}
	if attribute == "zone_coverage" {
		if progression.ZoneCoverage + modifier > 99 {
			progression.ZoneCoverage = 99
		} else {
			progression.ZoneCoverage += modifier
		}
	}
	if attribute == "man_coverage" {
		if progression.ManCoverage + modifier > 99 {
			progression.ManCoverage = 99
		} else {
			progression.ManCoverage += modifier
		}
	}
	if attribute == "strength" {
		if progression.Strength + modifier > 99 {
			progression.Strength = 99
		} else {
			progression.Strength += modifier
		}
	}
	if attribute == "tackle" {
		if progression.Tackle + modifier > 99 {
			progression.Tackle = 99
		} else {
			progression.Tackle += modifier
		}
	}
	if attribute == "pass_block" {
		if progression.PassBlock + modifier > 99 {
			progression.PassBlock = 99
		} else {
			progression.PassBlock += modifier
		}
	}
	if attribute == "run_block" {
		if progression.RunBlock + modifier > 99 {
			progression.RunBlock = 99
		} else {
			progression.RunBlock += modifier
		}
	}
	if attribute == "pass_rush" {
		if progression.PassRush + modifier > 99 {
			progression.PassRush = 99
		} else {
			progression.PassRush += modifier
		}
	}
	if attribute == "run_defense" {
		if progression.RunDefense + modifier > 99 {
			progression.RunDefense = 99
		} else {
			progression.RunDefense += modifier
		}
	}
	if attribute == "throw_power" {
		if progression.ThrowPower + modifier > 99 {
			progression.ThrowPower = 99
		} else {
			progression.ThrowPower += modifier
		}
	}
	if attribute == "throw_accuracy" {
		if progression.ThrowAccuracy + modifier > 99 {
			progression.ThrowAccuracy = 99
		} else {
			progression.ThrowAccuracy += modifier
		}
	}
	if attribute == "kick_accuracy" {
		if progression.KickAccuracy + modifier > 99 {
			progression.KickAccuracy = 99
		} else {
			progression.KickAccuracy += modifier
		}
	}
	if attribute == "kick_power" {
		if progression.KickPower + modifier > 99 {
			progression.KickPower = 99
		} else {
			progression.KickPower += modifier
		}
	}
	if attribute == "punt_accuracy" {
		if progression.PuntAccuracy + modifier > 99 {
			progression.PuntAccuracy = 99
		} else {
			progression.PuntAccuracy += modifier
		}
	}
	if attribute == "punt_power" {
		if progression.PuntPower + modifier > 99 {
			progression.PuntPower = 99
		} else {
			progression.PuntPower += modifier
		}
	}
}

// Returns which attribute to change
func getAttribute(position string, archetype string, drill string) string {
	if position == "K" {
		switch drill {
		case "accuracy":
			return "kick_accuracy"
		case "power":
			return "kick_power"
		default:
			// for team drills, pick a kicking attribute at random. K/P progressions could use the extra bonus.
			if rand.IntN(2) == 0 {
				return "kick_accuracy"
			} else {
				return "kick_power"
			}
		}
	} else if position == "P" {
		switch drill {
		case "accuracy":
			return "punt_accuracy"
		case "power":
			return "punt_power"
		default:
			// for team drills, pick a punting attribute at random. K/P progressions could use the extra bonus.
			if rand.IntN(2) == 0 {
				return "punt_accuracy"
			} else {
				return "punt_power"
			}
		}
	} else if drill == "speed" {
		return "speed"
	} else if drill == "lift" {
		return "strength"
	} else if drill == "film" {
		return "football_iq"
	} else if drill == "plyometrics" {
		return "agility"
	} else if position == "QB" {
		if drill == "dropback" {
			return "throw_power"
		} else if drill == "screen" {
			return "throw_accuracy"
		} else if strings.Contains(drill, "pass") {
			if archetype == "Pocket" || archetype == "Balanced" {
				return "throw_power"
			} else {
				return "throw_accuracy"
			}
		} else {
			if archetype == "Pocket" || archetype == "Balanced" {
				return "throw_accuracy"
			} else {
				return "throw_power"
			}
		}
	} else if position == "RB" {
		if drill == "gauntlet" {
			return "carrying"
		} else if drill == "square" {
			return "route_running"
		} else if drill == "blitzpickup" {
			return "pass_block"
		} else if drill == "jugs" {
			return "catching"
		} else if strings.Contains(drill, "pass") {
			switch archetype {
			case "Speed", "Balanced":
				return "catching"
			case "Power":
				return "pass_block"
			default:
				return "route_running"
			}
		} else {
			switch archetype {
			case "Power":
				return "strength"
			case "Receiving":
				return "agility"
			default:
				chance := rand.IntN(2)
				if chance == 0 {
					return "strength"
				}
				return "agility"
			}
		}
	} else if position == "FB" {
		if drill == "gauntlet" {
			return "carrying"
		} else if drill == "square" {
			return "route_running"
		} else if drill == "blitzpickup" {
			return "pass_block"
		} else if drill == "jugs" {
			return "catching"
		} else if drill == "lead" {
			return "run_block"
		} else if strings.Contains(drill, "pass") {
			if archetype == "Blocking" {
				return "pass_block"
			} else {
				return "catching"
			}
		} else {
			switch archetype {
			case "Blocking":
				return "run_block"
			case "Rushing":
				return "strength"
			case "Receiving":
				return "agility"
			default:
				chance := rand.IntN(3)
				if chance == 0 {
					return "strength"
				}
				if chance == 1 {
					return "run_block"
				}
				return "agility"
			}
		}
	} else if position == "TE" {
		if drill == "gauntlet" {
			return "carrying"
		} else if drill == "square" {
			return "route_running"
		} else if drill == "blitzpickup" {
			return "pass_block"
		} else if drill == "jugs" {
			return "catching"
		} else if drill == "arc" {
			return "run_block"
		} else if strings.Contains(drill, "pass") {
			switch archetype {
			case "Blocking":
				return "pass_block"
			case "Vertical Threat":
				return "speed"
			default:
				return "catching"
			}
		} else {
			return "run_block"
		}
	} else if position == "WR" {
		if drill == "gauntlet" {
			return "carrying"
		} else if drill == "square" {
			return "route_running"
		} else if drill == "jugs" {
			return "catching"
		} else if drill == "screenblock" {
			return "run_block"
		} else if strings.Contains(drill, "pass") {
			switch archetype {
			case "Possession", "Red Zone Threat":
				return "catching"
			case "Speed":
				return "speed"
			default:
				return "route_running"
			}
		} else {
			return "run_block"
		}
	} else if position == "OT" {
		if drill == "mirror" {
			return "pass_block"
		} else if drill == "sled" {
			return "run_block"
		} else if strings.Contains(drill, "pass") {
			return "pass_block"
		} else {
			return "run_block"
		}
	} else if position == "OG" {
		if drill == "mirror" {
			return "pass_block"
		} else if drill == "sled" {
			return "run_block"
		} else if strings.Contains(drill, "pass") {
			return "pass_block"
		} else {
			return "run_block"
		}
	} else if position == "C" {
		if drill == "mirror" {
			return "pass_block"
		} else if drill == "sled" {
			return "run_block"
		} else if strings.Contains(drill, "pass") {
			return "pass_block"
		} else {
			return "run_block"
		}
	} else if position == "DT" {
		if drill == "rip" {
			return "pass_rush"
		} else if drill == "shed" {
			return "run_defense"
		} else if strings.Contains(drill, "pass") {
			return "pass_rush"
		} else {
			return "run_defense"
		}
	} else if position == "DE" {
		if drill == "rip" {
			return "pass_rush"
		} else if drill == "shed" {
			return "run_defense"
		} else if strings.Contains(drill, "pass") {
			return "pass_rush"
		} else {
			return "run_defense"
		}
	} else if position == "OLB" {
		// EDGE
		if archetype == "Pass Rush" || archetype == "Run Stopper" {
			if drill == "rip" {
				return "pass_rush"
			} else if drill == "shed" {
				return "run_defense"
			} else if strings.Contains(drill, "pass") {
				return "pass_rush"
			} else {
				return "run_defense"
			}
			// Off Ball
		} else {
			if drill == "runfit" {
				return "run_defense"
			} else if drill == "rushlane" {
				return "pass_rush"
			} else if drill == "zonedrop" {
				return "zone_coverage"
			} else if drill == "hipturn" {
				return "man_coverage"
			} else if strings.Contains(drill, "pass") {
				chance := rand.IntN(2)
				if chance == 0 {
					return "zone_coverage"
				}
				return "man_coverage"
			} else {
				return "run_defense"
			}
		}
	} else if position == "ILB" {
		if drill == "runfit" {
			return "run_defense"
		} else if drill == "rushlane" {
			return "pass_rush"
		} else if drill == "zonedrop" {
			return "zone_coverage"
		} else if drill == "hipturn" {
			return "man_coverage"
		} else if strings.Contains(drill, "pass") {
			chance := rand.IntN(2)
			if chance == 0 {
				return "zone_coverage"
			}
			return "man_coverage"
		} else {
			return "run_defense"
		}
	} else if position == "CB" {
		if drill == "zonedrop" {
			return "zone_coverage"
		} else if drill == "hipturn" {
			return "man_coverage"
		} else if drill == "jugs" {
			return "catching"
		} else if strings.Contains(drill, "pass") {
			switch archetype {
			case "Man Coverage":
				return "man_coverage"
			case "Zone Coverage":
				return "zone_coverage"
			default:
				chance := rand.IntN(2)
				if chance == 0 {
					return "zone_coverage"
				}
				return "man_coverage"
			}
		} else {
			return "tackle"
		}
	} else if position == "FS" {
		if drill == "centerfield" {
			return "zone_coverage"
		} else if drill == "match" {
			return "man_coverage"
		} else if drill == "jugs" {
			return "catching"
		} else if drill == "alley" {
			return "run_defense"
		} else if drill == "handcombat" {
			return "pass_rush"
		} else if strings.Contains(drill, "pass") {
			switch archetype {
			case "Man Coverage":
				return "man_coverage"
			case "Zone Coverage":
				return "zone_coverage"
			default:
				chance := rand.IntN(2)
				if chance == 0 {
					return "zone_coverage"
				}
				return "man_coverage"
			}
		} else {
			return "tackle"
		}
	} else if position == "SS" {
		if drill == "centerfield" {
			return "zone_coverage"
		} else if drill == "match" {
			return "man_coverage"
		} else if drill == "jugs" {
			return "catching"
		} else if drill == "alley" {
			return "run_defense"
		} else if drill == "handcombat" {
			return "pass_rush"
		} else if strings.Contains(drill, "pass") {
			switch archetype {
			case "Man Coverage":
				return "man_coverage"
			case "Zone Coverage":
				return "zone_coverage"
			default:
				chance := rand.IntN(2)
				if chance == 0 {
					return "zone_coverage"
				}
				return "man_coverage"
			}
		} else {
			return "tackle"
		}
	}
	return "bad position"
}

var minorInjuries = []string{
	"Illness",
	"Stinger",
	"Strained Biceps",
	"Strained Triceps",
	"Bruised Foot",
	"Sprained Foot",
	"Bruised Hip",
	"Strained Hip",
	"Strained Groin",
	"Strained Calf",
	"Strained Quadriceps",
	"Sprained Wrist",
	"Elbow Tendonitis",
	"Bruised Elbow",
	"Sprained Elbow",
	"Strained Back",
	"Hyperextended Back",
	"Ankle Bruise",
	"Ankle Sprain",
	"Strained Shoulder",
	"Shoulder Tendonitis",
	"Separated Shoulder",
	"Sprained Thumb",
	"Sprained Knee",
	"Concussion",
}

var moderateInjuries = []string{
	"Lacerated Spleen",
	"Fractured Toe",
	"Dislocated Toe",
	"Bruised Toe",
	"Sprained Toe",
	"Turf Toe",
	"Wrist Bruise",
	"Sprained Wrist",
	"Hip Strain",
	"Fractured Ribs",
	"Achilles Tendonitis",
	"Bruised Achilles",
	"Dislocated Elbow",
	"Elbow Tendonitis",
	"Sprained Elbow",
	"Strained Groin",
	"Pulled Groin",
	"Strained Calf",
	"Pulled Calf",
	"Bruised Thumb",
	"Sprained Thumb",
	"Dislocated Thumb",
	"Fractured Thumb",
	"MCL Bruise",
	"PCL Bruise",
	"Patellar Tendon Bruise",
	"Strained Quadriceps",
	"Pulled Quadriceps",
	"Concussion",
	"Strained Biceps",
	"Pulled Biceps",
	"Strained Triceps",
	"Pulled Triceps",
	"Strained Pectoral",
	"Pulled Pectoral",
	"High Ankle Sprain",
	"Bruised Hamstring",
	"Pulled Hamstring",
	"Neck Bruise",
	"Sprained Neck",
	"ACL Bruise",
	"Dislocated Shoulder",
	"Shoulder Tendonitis",
	"Separated Shoulder",
	"Sprained Rotator Cuff",
	"Dislocated Ankle",
	"Dislocated Foot",
	"Sprained Foot",
	"Back Disk Tear",
}

var severeInjuries = []string{
	"Biceps Tear",
	"Triceps Tear",
	"Quadriceps Tear",
	"MCL Bruise",
	"MCL Tendonitis",
	"PCL Bruise",
	"PCL Tendonitis",
	"Patellar Tendon Bruise",
	"Patellar Tendonitis",
	"Knee Meniscus Bruise",
	"Knee Meniscus Tear",
	"Achilles Tendonitis",
	"Hamstring Tendonitis",
	"Fractured Wrist",
	"Fractured Jaw",
	"ACL Bruise",
	"ACL Tendonitis",
	"Calf Tear",
	"Groin Tear",
	"Pulled Pectoral",
	"Pectoral Tear",
	"Fractured Ribs",
	"Sprained Neck",
	"Back Disk Tear",
	"Fractured Ankle",
	"Fractured Foot",
	"Strained Rotator Cuff",
}

var seasonEndingInjuries = []string{
	"Ruptured Hamstring",
	"Patellar Tendon Tear",
	"Knee Meniscus Tear",
	"Fractured Foot",
	"Ruptured Achilles",
	"MCL Tear",
	"PCL Tear",
	"Fractured Hip",
	"Fractured Spine",
	"Fractured Ankle",
	"ACL Tear",
	"Rotator Cuff Tear",
}
