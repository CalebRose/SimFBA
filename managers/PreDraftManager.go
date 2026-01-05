package managers

import (
	"fmt"
	"math"
	"math/rand"
	"strings"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/models"
	"github.com/CalebRose/SimFBA/repository"
	config "github.com/CalebRose/SimFBA/secrets"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
)

func RunPreDraftEvents() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	// SET ALL EVENTS TO THIS SEASON ID

	draftees := GetAllNFLDraftees()

	// Create a list of events (Pro Days and Combines)
	eventList := GenerateTypicalListOfEvents()

	// Add the participants to each list
	eventList = AddParticipants(util.GetParticipantIDS(), eventList, draftees)

	globalEventResults := []models.EventResults{}

	// For each event, create a result for each player
	for _, event := range eventList {
		for _, player := range event.Participants {
			// For some % of draftees, create results based on their advertised grades, not their real grades.
			hidePerformance := ShouldHidePerformance()

			playerEvents := GenerateEvent(player, event, ts)

			// Run events on them
			playerEvents = RunEvents(player, hidePerformance, playerEvents)

			// Append the results to the global event list that will be written to the database
			globalEventResults = append(globalEventResults, playerEvents)
		}
	}

	// Export the results
	repository.CreatePreDraftEventResultsBatch(db, globalEventResults, 200)
}

// Set % of draftees that only perform based on their advertised grades, not real grades
func ShouldHidePerformance() bool {
	// 10% chance
	chance := 10

	roll := rand.Intn(100)

	if roll < chance {
		return true
	} else {
		return false
	}
}

func GenerateEvent(draftee models.NFLDraftee, event models.PreDraftEvent, timestamp structs.Timestamp) models.EventResults {
	var newEvent models.EventResults
	newEvent.SeasonID = uint(timestamp.NFLSeasonID)
	newEvent.PlayerID = uint(draftee.PlayerID)
	newEvent.IsCombine = event.IsCombine
	newEvent.Name = event.Name
	return newEvent
}

func RunEvents(draftee models.NFLDraftee, shouldHidePerformance bool, event models.EventResults) models.EventResults {
	// If we should hide the performance, we should create a dummy draftee based on the original draftee, but change their actual attributes to the correct values
	// for a player of that letter grade
	if shouldHidePerformance {
		dummy := GetDummyDraftee(draftee)
		event = RunUniversalEvents(dummy, shouldHidePerformance, event)
		event = RunPositionEvents(dummy, shouldHidePerformance, event)
	} else {
		event = RunUniversalEvents(draftee, shouldHidePerformance, event)
		event = RunPositionEvents(draftee, shouldHidePerformance, event)
	}
	return event
}

func RunUniversalEvents(draftee models.NFLDraftee, shouldHidePerformance bool, event models.EventResults) models.EventResults {
	event.FourtyYardDash = Run40YardDash(draftee.Speed, event.IsCombine)
	event.BenchPress = RunBenchPress(draftee.Strength, event.IsCombine, draftee.Position)
	event.Shuttle = RunShuttle(draftee.Agility, event.IsCombine)
	event.ThreeCone = Run3Cone(draftee.Agility, event.IsCombine)
	event.VerticalJump = RunVertJump(draftee.Agility, draftee.Strength, draftee.Weight, event.IsCombine)
	event.BroadJump = RunBroadJump(draftee.Agility, draftee.Strength, draftee.Weight, event.IsCombine)

	if event.IsCombine {
		event.Wonderlic = RunWonderlic(draftee.FootballIQ)
	}

	return event
}

func CombinePositionStr(position1 string, position2 string) string {
	if len(position2) == 0 {
		return position1
	} else {
		return (position1 + "/" + position2)
	}
}

func CombineArchetypeStr(arch1 string, arch2 string) string {
	if len(arch2) == 0 {
		return arch1
	} else {
		return (arch1 + "/" + arch2)
	}
}

func RunPositionEvents(draftee models.NFLDraftee, shouldHidePerformance bool, event models.EventResults) models.EventResults {
	// Combine Positions into one string for dual position players
	position := CombinePositionStr(draftee.Position, draftee.PositionTwo)
	archetype := CombineArchetypeStr(draftee.Archetype, draftee.ArchetypeTwo)

	// Must handle whether player's true attributes should be hidden
	if strings.Contains(strings.ToLower(position), strings.ToLower("QB")) {
		event.ThrowingDistance = RunQBDistance(draftee.ThrowPower, event.IsCombine)
		event.ThrowingAccuracy = RunQBAccuracy(draftee.ThrowAccuracy, event.IsCombine)
	}
	if strings.Contains(strings.ToLower(position), strings.ToLower("RB")) {
		event.InsideRun = RunInsideRun(draftee.Speed, draftee.Strength, event.IsCombine)
		event.OutsideRun = RunOutsideRun(draftee.Speed, draftee.Agility, event.IsCombine)
		event.Catching = RunCatching(draftee.Catching, event.IsCombine)
		event.RouteRunning = RunRouteRunning(draftee.RouteRunning, event.IsCombine)
	}
	if strings.Contains(strings.ToLower(position), strings.ToLower("WR")) {
		event.Catching = RunCatching(draftee.Catching, event.IsCombine)
		event.RouteRunning = RunRouteRunning(draftee.RouteRunning, event.IsCombine)
	}
	if strings.Contains(strings.ToLower(position), strings.ToLower("TE")) {
		event.Catching = RunCatching(draftee.Catching, event.IsCombine)
		event.RouteRunning = RunRouteRunning(draftee.RouteRunning, event.IsCombine)
		event.RunBlocking = RunRunBlocking(draftee.RunBlock, event.IsCombine, position)
	}
	if strings.Contains(strings.ToLower(position), strings.ToLower("FB")) {
		event.InsideRun = RunInsideRun(draftee.Speed, draftee.Strength, event.IsCombine)
		event.OutsideRun = RunOutsideRun(draftee.Speed, draftee.Agility, event.IsCombine)
		event.Catching = RunCatching(draftee.Catching, event.IsCombine)
		event.RouteRunning = RunRouteRunning(draftee.RouteRunning, event.IsCombine)
		event.RunBlocking = RunRunBlocking(draftee.RunBlock, event.IsCombine, position)
	}
	if strings.Contains(strings.ToLower(position), strings.ToLower("TE")) {
		event.Catching = RunCatching(draftee.Catching, event.IsCombine)
		event.RouteRunning = RunRouteRunning(draftee.RouteRunning, event.IsCombine)
		event.RunBlocking = RunRunBlocking(draftee.RunBlock, event.IsCombine, position)
	}
	if strings.Contains(strings.ToLower(position), strings.ToLower("OT")) {
		event.RunBlocking = RunRunBlocking(draftee.RunBlock, event.IsCombine, position)
		event.PassBlocking = RunPassBlocking(draftee.PassBlock, event.IsCombine, position)
	}
	if strings.Contains(strings.ToLower(position), strings.ToLower("OG")) {
		event.RunBlocking = RunRunBlocking(draftee.RunBlock, event.IsCombine, position)
		event.PassBlocking = RunPassBlocking(draftee.PassBlock, event.IsCombine, position)
	}
	if strings.Contains(strings.ToLower(position), strings.ToLower("C")) && !strings.Contains(strings.ToLower(position), strings.ToLower("CB")) { // Special case so we don't get CBs in here.
		event.RunBlocking = RunRunBlocking(draftee.RunBlock, event.IsCombine, position)
		event.PassBlocking = RunPassBlocking(draftee.PassBlock, event.IsCombine, position)
	}
	if strings.Contains(strings.ToLower(position), strings.ToLower("DT")) {
		event.RunStop = RunRunStop(draftee.RunDefense, event.IsCombine)
		event.PassRush = RunPassRush(draftee.PassRush, event.IsCombine)
	}
	if strings.Contains(strings.ToLower(position), strings.ToLower("DE")) {
		event.RunStop = RunRunStop(draftee.RunDefense, event.IsCombine)
		event.PassRush = RunPassRush(draftee.PassRush, event.IsCombine)
	}
	if strings.Contains(strings.ToLower(position), strings.ToLower("OLB")) && strings.Contains(strings.ToLower(archetype), strings.ToLower("Pass Rush")) {
		event.RunStop = RunRunStop(draftee.RunDefense, event.IsCombine)
		event.PassRush = RunPassRush(draftee.PassRush, event.IsCombine)
	}
	if strings.Contains(strings.ToLower(position), strings.ToLower("OLB")) && !strings.Contains(strings.ToLower(archetype), strings.ToLower("Pass Rush")) {
		event.RunStop = RunRunStop(draftee.RunDefense, event.IsCombine)
		event.LBCoverage = RunLBCoverage(draftee.ManCoverage, draftee.ZoneCoverage, event.IsCombine)
	}
	if strings.Contains(strings.ToLower(position), strings.ToLower("ILB")) {
		event.RunStop = RunRunStop(draftee.RunDefense, event.IsCombine)
		event.LBCoverage = RunLBCoverage(draftee.ManCoverage, draftee.ZoneCoverage, event.IsCombine)
	}
	if strings.Contains(strings.ToLower(position), strings.ToLower("CB")) {
		event.ManCoverage = RunManCoverage(draftee.ManCoverage, event.IsCombine)
		event.ZoneCoverage = RunZoneCoverage(draftee.ZoneCoverage, event.IsCombine)
	}
	if strings.Contains(strings.ToLower(position), strings.ToLower("FS")) {
		event.ManCoverage = RunManCoverage(draftee.ManCoverage, event.IsCombine)
		event.ZoneCoverage = RunZoneCoverage(draftee.ZoneCoverage, event.IsCombine)
	}
	if strings.Contains(strings.ToLower(position), strings.ToLower("SS")) {
		event.ManCoverage = RunManCoverage(draftee.ManCoverage, event.IsCombine)
		event.ZoneCoverage = RunZoneCoverage(draftee.ZoneCoverage, event.IsCombine)
	}
	if strings.Contains(strings.ToLower(position), strings.ToLower("K")) {
		event.Kickoff = RunKickoffDrill(draftee.KickPower, draftee.PuntPower, event.IsCombine)
		event.Fieldgoal = RunFieldGoalDrill(draftee.KickPower, draftee.KickAccuracy, event.IsCombine)
		event.PuntDistance = RunPuntDistance(draftee.PuntPower, event.IsCombine)
		event.CoffinPunt = RunCoffinPunt(draftee.PuntAccuracy, event.IsCombine)
	}
	if strings.Contains(strings.ToLower(position), strings.ToLower("P")) {
		event.Kickoff = RunKickoffDrill(draftee.KickPower, draftee.PuntPower, event.IsCombine)
		event.Fieldgoal = RunFieldGoalDrill(draftee.KickPower, draftee.KickAccuracy, event.IsCombine)
		event.PuntDistance = RunPuntDistance(draftee.PuntPower, event.IsCombine)
		event.CoffinPunt = RunCoffinPunt(draftee.PuntAccuracy, event.IsCombine)
	}

	return event
}

func GetDummyDraftee(orginalDraftee models.NFLDraftee) models.NFLDraftee {
	// Attribute means
	// Use standard letter grade conversion to know how many standard deviations from the mean
	// With the mean and standard deviations, we can calculate a number to run drills for all of the skills
	attributeMeans := config.AttributeMeans()

	tempDraftee := orginalDraftee
	tempDraftee.Speed = int(GetNewAttributeRating(tempDraftee.SpeedGrade, attributeMeans, "Speed", (tempDraftee.Position)))
	tempDraftee.Agility = int(GetNewAttributeRating(tempDraftee.AgilityGrade, attributeMeans, "Agility", (tempDraftee.Position)))
	tempDraftee.Strength = int(GetNewAttributeRating(tempDraftee.StrengthGrade, attributeMeans, "Strength", (tempDraftee.Position)))
	tempDraftee.ThrowPower = int(GetNewAttributeRating(tempDraftee.ThrowPowerGrade, attributeMeans, "ThrowPower", (tempDraftee.Position)))
	tempDraftee.ThrowAccuracy = int(GetNewAttributeRating(tempDraftee.ThrowAccuracyGrade, attributeMeans, "ThrowAccuracy", (tempDraftee.Position)))
	tempDraftee.Catching = int(GetNewAttributeRating(tempDraftee.CarryingGrade, attributeMeans, "Catching", (tempDraftee.Position)))
	tempDraftee.RouteRunning = int(GetNewAttributeRating(tempDraftee.RouteRunningGrade, attributeMeans, "RouteRunning", (tempDraftee.Position)))
	tempDraftee.RunBlock = int(GetNewAttributeRating(tempDraftee.RunBlockGrade, attributeMeans, "RunBlock", (tempDraftee.Position)))
	tempDraftee.PassBlock = int(GetNewAttributeRating(tempDraftee.PassBlockGrade, attributeMeans, "PassBlock", (tempDraftee.Position)))
	tempDraftee.RunDefense = int(GetNewAttributeRating(tempDraftee.RunDefenseGrade, attributeMeans, "RunDefense", (tempDraftee.Position)))
	tempDraftee.PassRush = int(GetNewAttributeRating(tempDraftee.PassRushGrade, attributeMeans, "PassRush", (tempDraftee.Position)))
	tempDraftee.ManCoverage = int(GetNewAttributeRating(tempDraftee.ManCoverageGrade, attributeMeans, "ManCoverage", (tempDraftee.Position)))
	tempDraftee.ZoneCoverage = int(GetNewAttributeRating(tempDraftee.ZoneCoverageGrade, attributeMeans, "ZoneCoverage", (tempDraftee.Position)))
	tempDraftee.KickPower = int(GetNewAttributeRating(tempDraftee.KickPowerGrade, attributeMeans, "KickPower", (tempDraftee.Position)))
	tempDraftee.KickAccuracy = int(GetNewAttributeRating(tempDraftee.KickAccuracyGrade, attributeMeans, "KickAccuracy", (tempDraftee.Position)))
	tempDraftee.PuntPower = int(GetNewAttributeRating(tempDraftee.PuntPowerGrade, attributeMeans, "PuntPower", (tempDraftee.Position)))
	tempDraftee.PuntAccuracy = int(GetNewAttributeRating(tempDraftee.PuntAccuracyGrade, attributeMeans, "PuntAccuracy", (tempDraftee.Position)))
	tempDraftee.FootballIQ = int(GetNewAttributeRating(tempDraftee.FootballIQGrade, attributeMeans, "FootballIQ", (tempDraftee.Position)))
	tempDraftee.Tackle = int(GetNewAttributeRating(tempDraftee.FootballIQGrade, attributeMeans, "Tackle", (tempDraftee.Position)))
	tempDraftee.Carrying = int(GetNewAttributeRating(tempDraftee.FootballIQGrade, attributeMeans, "Carrying", (tempDraftee.Position)))

	return tempDraftee
}

func GetMeanForAttribute(mapping map[string]map[string]map[string]float32, attribute string, position string) float32 {
	return mapping[attribute][position]["mean"]
}

func GetStdDevForAttribute(mapping map[string]map[string]map[string]float32, attribute string, position string) float32 {
	return mapping[attribute][position]["stddev"]
}

func GetNewAttributeRating(grade string, mapping map[string]map[string]map[string]float32, attribute string, position string) uint {
	mean := GetMeanForAttribute(mapping, attribute, position)
	stddev := GetStdDevForAttribute(mapping, attribute, position)

	return uint(mean + (stddev * TranslateLetterGradeToStdDevs(grade)))
}

func TranslateLetterGradeToStdDevs(grade string) float32 {
	switch grade {
	case "A+":
		return 2.5
	case "A":
		return 2
	case "A-":
		return 1.75
	case "B+":
		return 1.5
	case "B":
		return 1
	case "B-":
		return 0.75
	case "C+":
		return 0.5
	case "C":
		return 0
	case "C-":
		return -0.5
	case "D+":
		return -0.75
	case "D":
		return -1
	case "D-":
		return -1.5
	default: // F
		return -2
	}
}

func Run40YardDash(speed int, isCombine bool) float32 {
	delta := GetDelta(isCombine)

	temp := float64(speed) + delta

	if temp > 99.0 {
		temp = 99.0
	}

	temp = 100 - temp
	temp = math.Pow(temp, 2)
	temp = temp / 4000
	temp = temp + 4.3

	return float32(temp)
}

func RunBenchPress(strength int, isCombine bool, position string) uint8 {
	delta := GetDelta(isCombine)

	temp := 0.0

	if strings.Contains(strings.ToLower(position), strings.ToLower("FB")) {
		temp = float64(strength) + delta - 10.0
	} else {
		temp = float64(strength) + delta
	}
	if temp > 99.0 {
		temp = 99.0
	}
	temp = 185.0 - temp
	temp = math.Pow(temp, 2)
	temp = temp / 600.0
	temp = temp * -1.0
	temp = temp + 66.0
	return uint8(temp)
}

func RunShuttle(agility int, isCombine bool) float32 {
	delta := GetDelta(isCombine)

	temp := float64(agility) + delta
	if temp > 99.0 {
		temp = 99.0
	}
	temp = 100.0 - temp
	temp = math.Pow(temp, 2)
	temp = temp / 6000.0
	temp = temp + 3.7
	return float32(temp)
}

func Run3Cone(agility int, isCombine bool) float32 {
	delta := GetDelta(isCombine)

	temp := float64(agility) + delta
	if temp > 99.0 {
		temp = 99.0
	}
	temp = 100.0 - temp
	temp = math.Pow(temp, 2)
	temp = temp / 4200.0
	temp = temp + 6.28
	return float32(temp)
}

func RunVertJump(agility int, strength int, weight int, isCombine bool) uint8 {
	delta := GetDelta(isCombine)
	newStrength := float64(strength) + delta
	if newStrength > 99.0 {
		newStrength = 99.0
	}
	delta = GetDelta(isCombine)
	newAgility := float64(agility) + delta
	if newAgility > 99.0 {
		newAgility = 99.0
	}
	temp := ((newStrength + newAgility) / float64(weight))
	temp = 1500 * temp
	temp = math.Sqrt(temp)
	temp = temp + 11.0
	return uint8(temp)
}

func RunBroadJump(agility int, strength int, weight int, isCombine bool) uint8 {
	delta := GetDelta(isCombine)
	newStrength := float64(strength) + delta
	if newStrength > 99.0 {
		newStrength = 99.0
	}
	delta = GetDelta(isCombine)
	newAgility := float64(agility) + delta
	if newAgility > 99.0 {
		newAgility = 99.0
	}
	temp := ((newStrength + newAgility) / float64(weight))
	temp = 20000 * temp
	temp = math.Sqrt(temp)
	temp = temp / 2.0
	temp = temp + 79.0
	return uint8(temp)
}

func RunWonderlic(fbIQ int) uint8 {
	delta := GetDelta(true)
	temp := float64(fbIQ) + delta
	if temp > 99.0 {
		temp = 99.0
	}
	temp = temp - 130.0
	temp = math.Pow(temp, 3)
	temp = temp / 25000.0
	temp = temp + 51.0
	return uint8(temp)
}

func RunQBAccuracy(throwAcc int, isCombine bool) float32 {
	delta := GetDelta(isCombine)

	temp := float64(throwAcc) + delta
	// Make sure they can't score more than 10.00
	if temp > 80.0 {
		temp = 80.0
	}
	temp = temp / 80.0
	temp = temp * 10.0
	return float32(temp)
}

func RunQBDistance(throwPow int, isCombine bool) float32 {
	delta := GetDelta(isCombine)

	temp := float64(throwPow) + delta
	// Make sure they can't score more than 10.00
	if temp > 85.0 {
		temp = 85.0
	}
	temp = temp / 85.0
	temp = temp * 10.0
	return float32(temp)
}

func RunInsideRun(speed int, strength int, isCombine bool) float32 {
	delta := GetDelta(isCombine)
	newSpeed := float64(speed) + delta
	delta = GetDelta(isCombine)
	newStrength := float64(strength) + delta
	temp := newSpeed + newStrength
	// Make sure they can't score more than 10.00
	if temp > 170.0 {
		temp = 170.0
	}
	temp = temp / 170.0
	temp = temp * 10.0
	return float32(temp)
}

func RunOutsideRun(speed int, agility int, isCombine bool) float32 {
	delta := GetDelta(isCombine)
	newSpeed := float64(speed) + delta
	delta = GetDelta(isCombine)
	newAgility := float64(agility) + delta
	temp := newSpeed + newAgility
	// Make sure they can't score more than 10.00
	if temp > 180.0 {
		temp = 180.0
	}
	temp = temp / 180.0
	temp = temp * 10.0
	return float32(temp)
}

func RunCatching(catching int, isCombine bool) float32 {
	delta := GetDelta(isCombine)
	temp := float64(catching) + delta
	// Make sure they can't score more than 10.00
	if temp > 80.0 {
		temp = 80.0
	}
	temp = temp / 80.0
	temp = temp * 10.0
	return float32(temp)
}

func RunRouteRunning(routeRun int, isCombine bool) float32 {
	delta := GetDelta(isCombine)
	temp := float64(routeRun) + delta
	// Make sure they can't score more than 10.00
	if temp > 65.0 {
		temp = 65.0
	}
	temp = temp / 65.0
	temp = temp * 10.0
	return float32(temp)
}

func RunRunBlocking(runBlock int, isCombine bool, position string) float32 {
	delta := GetDelta(isCombine)
	temp := 0.0
	// Make sure they can't score more than 10.00
	if strings.Contains(strings.ToLower(position), strings.ToLower("FB")) {
		temp = float64(runBlock) + delta - 15.0
	} else if strings.Contains(strings.ToLower(position), strings.ToLower("TE")) {
		temp = float64(runBlock) + delta - 8.0
	} else {
		temp = float64(runBlock) + delta
	}
	if temp > 85.0 {
		temp = 85.0
	}
	temp = temp / 85.0
	temp = temp * 10.0
	return float32(temp)
}

func RunPassBlocking(passBlock int, isCombine bool, position string) float32 {
	delta := GetDelta(isCombine)
	temp := 0.0
	// Make sure they can't score more than 10.00
	if strings.Contains(strings.ToLower(position), strings.ToLower("FB")) {
		temp = float64(passBlock) + delta - 15.0
	} else if strings.Contains(strings.ToLower(position), strings.ToLower("TE")) {
		temp = float64(passBlock) + delta - 8.0
	} else {
		temp = float64(passBlock) + delta
	}
	if temp > 85.0 {
		temp = 85.0
	}
	temp = temp / 85.0
	temp = temp * 10.0
	return float32(temp)
}

func RunRunStop(runDef int, isCombine bool) float32 {
	delta := GetDelta(isCombine)
	temp := float64(runDef) + delta
	// Make sure they can't score more than 10.00
	if temp > 85.0 {
		temp = 85.0
	}
	temp = temp / 85.0
	temp = temp * 10.0
	return float32(temp)
}

func RunPassRush(PassRush int, isCombine bool) float32 {
	delta := GetDelta(isCombine)
	temp := float64(PassRush) + delta
	// Make sure they can't score more than 10.00
	if temp > 80.0 {
		temp = 80.0
	}
	temp = temp / 80.0
	temp = temp * 10.0
	return float32(temp)
}

func RunLBCoverage(manCov int, zonCov int, isCombine bool) float32 {
	delta := GetDelta(isCombine)
	newManCov := float64(manCov) + delta
	delta = GetDelta(isCombine)
	newZonCov := float64(zonCov) + delta
	temp := newManCov + newZonCov
	if temp > 150.0 {
		temp = 150.0
	}
	temp = temp / 150.0
	temp = temp * 10
	return float32(temp)
}

func RunManCoverage(manCov int, isCombine bool) float32 {
	delta := GetDelta(isCombine)
	temp := float64(manCov) + delta
	if temp > 90.0 {
		temp = 90.0
	}
	temp = temp / 90.0
	temp = temp * 10
	return float32(temp)
}

func RunZoneCoverage(zonCov int, isCombine bool) float32 {
	delta := GetDelta(isCombine)
	temp := float64(zonCov) + delta
	if temp > 80.0 {
		temp = 80.0
	}
	temp = temp / 80.0
	temp = temp * 10
	return float32(temp)
}

func RunKickoffDrill(kickPow int, puntPow int, isCombine bool) float32 {
	delta := GetDelta(isCombine)

	// Get larger of the two kicking values
	temp := math.Max(float64(kickPow), float64(puntPow)) + delta
	if temp > 75.0 {
		temp = 75.0
	}
	temp = temp / 75.0
	temp = temp * 10
	return float32(temp)
}

func RunFieldGoalDrill(kickPow int, kickAcc int, isCombine bool) float32 {
	delta := GetDelta(isCombine)
	newKickPow := float64(kickPow) + delta
	delta = GetDelta(isCombine)
	newKickAcc := float64(kickAcc) + delta
	temp := newKickPow + newKickAcc
	if temp > 155.0 {
		temp = 155.0
	}
	temp = temp / 155.0
	temp = temp * 10
	return float32(temp)
}

func RunPuntDistance(puntPow int, isCombine bool) float32 {
	delta := GetDelta(isCombine)

	// Get larger of the two kicking values
	temp := float64(puntPow) + delta
	if temp > 60.0 {
		temp = 60.0
	}
	temp = temp / 60.0
	temp = temp * 10
	return float32(temp)
}

func RunCoffinPunt(puntAcc int, isCombine bool) float32 {
	delta := GetDelta(isCombine)

	// Get larger of the two kicking values
	temp := float64(puntAcc) + delta
	if temp > 66.0 {
		temp = 66.0
	}
	temp = temp / 66.0
	temp = temp * 10
	return float32(temp)
}

func GetDelta(isCombine bool) float64 {
	min := 0
	max := 0

	if isCombine {
		min = -10
		max = 10
	} else {
		min = -5
		max = 15
	}

	return float64(rand.Intn((max - min) + min))
}

func AddParticipants(json map[string][]uint, events []models.PreDraftEvent, players []models.NFLDraftee) []models.PreDraftEvent {
	// For each event, find that events list of players, and get that player from the draftee list and add to participants
	for i, event := range events {
		for _, player := range json[event.Name] {
			participant, found := FindParticipant(player, players)

			if found {
				events[i].Participants = append(events[i].Participants, participant)
			} else {
				// Participant not found
				fmt.Println("EVENT PARTICIPANT NOT FOUND!!!")
			}
		}
	}

	return events
}

func FindParticipant(x uint, list []models.NFLDraftee) (models.NFLDraftee, bool) {
	for _, n := range list {
		if x == uint(n.PlayerID) {
			return n, true
		}
	}
	return models.NFLDraftee{}, false
}

func GenerateTypicalListOfEvents() []models.PreDraftEvent {
	var tempList []models.PreDraftEvent
	var tempEvent models.PreDraftEvent
	tempEvent.Name = "AAC Pro Day"
	tempEvent.IsCombine = false
	tempList = append(tempList, tempEvent)

	tempEvent.Name = "ACC Pro Day"
	tempEvent.IsCombine = false
	tempList = append(tempList, tempEvent)

	tempEvent.Name = "Big Ten Pro Day"
	tempEvent.IsCombine = false
	tempList = append(tempList, tempEvent)

	tempEvent.Name = "Big XII Pro Day"
	tempEvent.IsCombine = false
	tempList = append(tempList, tempEvent)

	tempEvent.Name = "C-USA Pro Day"
	tempEvent.IsCombine = false
	tempList = append(tempList, tempEvent)

	tempEvent.Name = "MAC Pro Day"
	tempEvent.IsCombine = false
	tempList = append(tempList, tempEvent)

	tempEvent.Name = "MWC Pro Day"
	tempEvent.IsCombine = false
	tempList = append(tempList, tempEvent)

	tempEvent.Name = "SEC Pro Day"
	tempEvent.IsCombine = false
	tempList = append(tempList, tempEvent)

	tempEvent.Name = "Sun Belt Pro Day"
	tempEvent.IsCombine = false
	tempList = append(tempList, tempEvent)

	tempEvent.Name = "FCS Pro Day"
	tempEvent.IsCombine = false
	tempList = append(tempList, tempEvent)

	tempEvent.Name = "NFL Combine"
	tempEvent.IsCombine = true
	tempList = append(tempList, tempEvent)

	return tempList
}
