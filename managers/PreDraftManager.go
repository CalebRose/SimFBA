package managers

import (
	"encoding/csv"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/models"
	config "github.com/CalebRose/SimFBA/secrets"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
)

func RunPreDraftEvents() {
	_ = dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	draftees := GetAllNFLDraftees()

	playerMap := make(map[uint]models.NFLDraftee)
	for _, d := range draftees {
		playerMap[uint(d.PlayerID)] = d
	}

	eventList := GenerateTypicalListOfEvents()
	eventList = AddParticipants(util.GetParticipantIDS(), eventList, draftees)

	globalEventResults := []models.EventResults{}

	for _, event := range eventList {
		for _, player := range event.Participants {
			hidePerformance := ShouldHidePerformance()
			playerEvents := GenerateEvent(player, event, ts)
			playerEvents = RunEvents(player, hidePerformance, playerEvents)
			globalEventResults = append(globalEventResults, playerEvents)
		}
	}

	fmt.Println("--- DRY RUN: EXPORTING FULL SCOUTING RESULTS ---")
	ExportResultsToCSV(globalEventResults, playerMap)
}

func ExportResultsToCSV(results []models.EventResults, playerMap map[uint]models.NFLDraftee) {
	file, err := os.Create("full_scouting_results.csv")
	if err != nil {
		fmt.Println("Error creating CSV:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{
		"PlayerID", "First Name", "Last Name", "School", "Pos", "Pos2", "Archetype", "Arch2", "OVR_Grade", "Height", "Weight", "Event",
		"40Yard", "Bench", "3Cone", "Shuttle", "Vert", "Broad", "Wonderlic",
		"ThrowDist", "ThrowAcc", "InsideRun", "OutsideRun", "Catching", "RouteRun",
		"PassBlock", "RunBlock", "PassRush", "RunStop", "ManCov", "ZoneCov", "LBCov",
		"Kickoff", "FG", "PuntDist", "CoffinPunt",
		"G_Speed", "G_Agility", "G_Strength", "G_ThrowPower", "G_ThrowAcc", "G_Catch",
		"G_Carry", "G_RouteRun", "G_RunBlock", "G_PassBlock", "G_Tackle", "G_RunDef",
		"G_PassRush", "G_ManCov", "G_ZoneCov", "G_KickPow", "G_KickAcc", "G_PuntPow", "G_PuntAcc", "G_IQ",
	}
	writer.Write(header)

	for _, r := range results {
		p := playerMap[r.PlayerID]
		row := []string{
			strconv.Itoa(int(r.PlayerID)),
			p.FirstName,
			p.LastName,
			p.College,
			p.Position,
			p.PositionTwo,
			p.Archetype,
			p.ArchetypeTwo,
			p.OverallGrade,
			strconv.Itoa(p.Height),
			strconv.Itoa(p.Weight),
			r.Name,
			fmt.Sprintf("%.2f", r.FourtyYardDash),
			strconv.Itoa(int(r.BenchPress)),
			fmt.Sprintf("%.2f", r.ThreeCone),
			fmt.Sprintf("%.2f", r.Shuttle),
			strconv.Itoa(int(r.VerticalJump)),
			strconv.Itoa(int(r.BroadJump)),
			strconv.Itoa(int(r.Wonderlic)),
			fmt.Sprintf("%.2f", r.ThrowingDistance),
			fmt.Sprintf("%.2f", r.ThrowingAccuracy),
			fmt.Sprintf("%.2f", r.InsideRun),
			fmt.Sprintf("%.2f", r.OutsideRun),
			fmt.Sprintf("%.2f", r.Catching),
			fmt.Sprintf("%.2f", r.RouteRunning),
			fmt.Sprintf("%.2f", r.PassBlocking),
			fmt.Sprintf("%.2f", r.RunBlocking),
			fmt.Sprintf("%.2f", r.PassRush),
			fmt.Sprintf("%.2f", r.RunStop),
			fmt.Sprintf("%.2f", r.ManCoverage),
			fmt.Sprintf("%.2f", r.ZoneCoverage),
			fmt.Sprintf("%.2f", r.LBCoverage),
			fmt.Sprintf("%.2f", r.Kickoff),
			fmt.Sprintf("%.2f", r.Fieldgoal),
			fmt.Sprintf("%.2f", r.PuntDistance),
			fmt.Sprintf("%.2f", r.CoffinPunt),
			p.SpeedGrade, p.AgilityGrade, p.StrengthGrade, p.ThrowPowerGrade, p.ThrowAccuracyGrade,
			p.CatchingGrade, p.CarryingGrade, p.RouteRunningGrade, p.RunBlockGrade, p.PassBlockGrade,
			p.TackleGrade, p.RunDefenseGrade, p.PassRushGrade, p.ManCoverageGrade, p.ZoneCoverageGrade,
			p.KickPowerGrade, p.KickAccuracyGrade, p.PuntPowerGrade, p.PuntAccuracyGrade, p.FootballIQGrade,
		}
		writer.Write(row)
	}
}

func ShouldHidePerformance() bool {
	return rand.Intn(100) < 10
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
	}
	return (position1 + "/" + position2)
}

func CombineArchetypeStr(arch1 string, arch2 string) string {
	if len(arch2) == 0 {
		return arch1
	}
	return (arch1 + "/" + arch2)
}

func RunPositionEvents(draftee models.NFLDraftee, shouldHidePerformance bool, event models.EventResults) models.EventResults {
	position := CombinePositionStr(draftee.Position, draftee.PositionTwo)
	archetype := CombineArchetypeStr(draftee.Archetype, draftee.ArchetypeTwo)

	if strings.Contains(strings.ToLower(position), "qb") {
		event.ThrowingDistance = RunQBDistance(draftee.ThrowPower, event.IsCombine)
		event.ThrowingAccuracy = RunQBAccuracy(draftee.ThrowAccuracy, event.IsCombine)
	}
	if strings.Contains(strings.ToLower(position), "rb") {
		event.InsideRun = RunInsideRun(draftee.Speed, draftee.Strength, event.IsCombine)
		event.OutsideRun = RunOutsideRun(draftee.Speed, draftee.Agility, event.IsCombine)
		event.Catching = RunCatching(draftee.Catching, event.IsCombine)
		event.RouteRunning = RunRouteRunning(draftee.RouteRunning, event.IsCombine)
	}
	if strings.Contains(strings.ToLower(position), "wr") {
		event.Catching = RunCatching(draftee.Catching, event.IsCombine)
		event.RouteRunning = RunRouteRunning(draftee.RouteRunning, event.IsCombine)
	}
	if strings.Contains(strings.ToLower(position), "te") {
		event.Catching = RunCatching(draftee.Catching, event.IsCombine)
		event.RouteRunning = RunRouteRunning(draftee.RouteRunning, event.IsCombine)
		event.RunBlocking = RunRunBlocking(draftee.RunBlock, event.IsCombine, position)
	}
	if strings.Contains(strings.ToLower(position), "fb") {
		event.InsideRun = RunInsideRun(draftee.Speed, draftee.Strength, event.IsCombine)
		event.OutsideRun = RunOutsideRun(draftee.Speed, draftee.Agility, event.IsCombine)
		event.Catching = RunCatching(draftee.Catching, event.IsCombine)
		event.RouteRunning = RunRouteRunning(draftee.RouteRunning, event.IsCombine)
		event.RunBlocking = RunRunBlocking(draftee.RunBlock, event.IsCombine, position)
	}
	if strings.Contains(strings.ToLower(position), "ot") || strings.Contains(strings.ToLower(position), "og") || (strings.Contains(strings.ToLower(position), "c") && !strings.Contains(strings.ToLower(position), "cb")) {
		event.RunBlocking = RunRunBlocking(draftee.RunBlock, event.IsCombine, position)
		event.PassBlocking = RunPassBlocking(draftee.PassBlock, event.IsCombine, position)
	}
	if strings.Contains(strings.ToLower(position), "dt") || strings.Contains(strings.ToLower(position), "de") || (strings.Contains(strings.ToLower(position), "olb") && strings.Contains(strings.ToLower(archetype), "pass rush")) {
		event.RunStop = RunRunStop(draftee.RunDefense, event.IsCombine)
		event.PassRush = RunPassRush(draftee.PassRush, event.IsCombine)
	}
	if (strings.Contains(strings.ToLower(position), "olb") && !strings.Contains(strings.ToLower(archetype), "pass rush")) || strings.Contains(strings.ToLower(position), "ilb") {
		event.RunStop = RunRunStop(draftee.RunDefense, event.IsCombine)
		event.LBCoverage = RunLBCoverage(draftee.ManCoverage, draftee.ZoneCoverage, event.IsCombine)
	}
	if strings.Contains(strings.ToLower(position), "cb") || strings.Contains(strings.ToLower(position), "fs") || strings.Contains(strings.ToLower(position), "ss") {
		event.ManCoverage = RunManCoverage(draftee.ManCoverage, event.IsCombine)
		event.ZoneCoverage = RunZoneCoverage(draftee.ZoneCoverage, event.IsCombine)
	}
	if strings.Contains(strings.ToLower(position), "k") || strings.Contains(strings.ToLower(position), "p") {
		event.Kickoff = RunKickoffDrill(draftee.KickPower, draftee.PuntPower, event.IsCombine)
		event.Fieldgoal = RunFieldGoalDrill(draftee.KickPower, draftee.KickAccuracy, event.IsCombine)
		event.PuntDistance = RunPuntDistance(draftee.PuntPower, event.IsCombine)
		event.CoffinPunt = RunCoffinPunt(draftee.PuntAccuracy, event.IsCombine)
	}

	return event
}

func GetDummyDraftee(orginalDraftee models.NFLDraftee) models.NFLDraftee {
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
	default:
		return -2
	}
}

func Run40YardDash(speed int, isCombine bool) float32 {
	delta := GetDelta(isCombine)
	temp := math.Min(float64(speed)+delta, 99.0)
	res := math.Pow(100-temp, 2)/4000 + 4.3
	return float32(res)
}

func RunBenchPress(strength int, isCombine bool, position string) uint8 {
	delta := GetDelta(isCombine)
	temp := float64(strength) + delta
	if strings.Contains(strings.ToLower(position), "fb") {
		temp -= 10.0
	}
	temp = math.Min(temp, 99.0)
	res := (math.Pow(185-temp, 2)/600)*-1.0 + 66.0
	return uint8(math.Max(0, res))
}

func RunShuttle(agility int, isCombine bool) float32 {
	delta := GetDelta(isCombine)
	temp := math.Min(float64(agility)+delta, 99.0)
	res := math.Pow(100.0-temp, 2)/6000.0 + 3.7
	return float32(res)
}

func Run3Cone(agility int, isCombine bool) float32 {
	delta := GetDelta(isCombine)
	temp := math.Min(float64(agility)+delta, 99.0)
	res := math.Pow(100.0-temp, 2)/4200.0 + 6.28
	return float32(res)
}

func RunVertJump(agility int, strength int, weight int, isCombine bool) uint8 {
	s := math.Min(float64(strength)+GetDelta(isCombine), 99.0)
	a := math.Min(float64(agility)+GetDelta(isCombine), 99.0)
	res := math.Sqrt(((s+a)/float64(weight))*1500.0) + 11.0
	return uint8(math.Max(0, res))
}

func RunBroadJump(agility int, strength int, weight int, isCombine bool) uint8 {
	s := math.Min(float64(strength)+GetDelta(isCombine), 99.0)
	a := math.Min(float64(agility)+GetDelta(isCombine), 99.0)
	res := (math.Sqrt(((s+a)/float64(weight))*20000.0) / 2.0) + 79.0
	return uint8(math.Max(0, res))
}

func RunWonderlic(fbIQ int) uint8 {
	temp := math.Min(float64(fbIQ)+GetDelta(true), 99.0)
	res := math.Pow(temp-130.0, 3)/25000.0 + 51.0
	// Clamp result between 0 and 50 to prevent overflow wrap-around
	if res < 0 {
		return 0
	}
	if res > 50 {
		return 50
	}
	return uint8(res)
}

func RunQBAccuracy(throwAcc int, isCombine bool) float32 {
	temp := math.Max(0, math.Min(float64(throwAcc)+GetDelta(isCombine), 80.0))
	return float32((temp / 80.0) * 10.0)
}

func RunQBDistance(throwPow int, isCombine bool) float32 {
	temp := math.Max(0, math.Min(float64(throwPow)+GetDelta(isCombine), 85.0))
	return float32((temp / 85.0) * 10.0)
}

func RunInsideRun(speed int, strength int, isCombine bool) float32 {
	temp := math.Max(0, math.Min((float64(speed)+GetDelta(isCombine))+(float64(strength)+GetDelta(isCombine)), 170.0))
	return float32((temp / 170.0) * 10.0)
}

func RunOutsideRun(speed int, agility int, isCombine bool) float32 {
	temp := math.Max(0, math.Min((float64(speed)+GetDelta(isCombine))+(float64(agility)+GetDelta(isCombine)), 180.0))
	return float32((temp / 180.0) * 10.0)
}

func RunCatching(catching int, isCombine bool) float32 {
	temp := math.Max(0, math.Min(float64(catching)+GetDelta(isCombine), 80.0))
	return float32((temp / 80.0) * 10.0)
}

func RunRouteRunning(routeRun int, isCombine bool) float32 {
	temp := math.Max(0, math.Min(float64(routeRun)+GetDelta(isCombine), 65.0))
	return float32((temp / 65.0) * 10.0)
}

func RunRunBlocking(runBlock int, isCombine bool, position string) float32 {
	temp := float64(runBlock) + GetDelta(isCombine)
	if strings.Contains(strings.ToLower(position), "fb") {
		temp -= 15.0
	}
	if strings.Contains(strings.ToLower(position), "te") {
		temp -= 8.0
	}
	temp = math.Max(0, math.Min(temp, 85.0))
	return float32((temp / 85.0) * 10.0)
}

func RunPassBlocking(passBlock int, isCombine bool, position string) float32 {
	temp := float64(passBlock) + GetDelta(isCombine)
	if strings.Contains(strings.ToLower(position), "fb") {
		temp -= 15.0
	}
	if strings.Contains(strings.ToLower(position), "te") {
		temp -= 8.0
	}
	temp = math.Max(0, math.Min(temp, 85.0))
	return float32((temp / 85.0) * 10.0)
}

func RunRunStop(runDef int, isCombine bool) float32 {
	temp := math.Max(0, math.Min(float64(runDef)+GetDelta(isCombine), 85.0))
	return float32((temp / 85.0) * 10.0)
}

func RunPassRush(PassRush int, isCombine bool) float32 {
	temp := math.Max(0, math.Min(float64(PassRush)+GetDelta(isCombine), 80.0))
	return float32((temp / 80.0) * 10.0)
}

func RunLBCoverage(manCov int, zonCov int, isCombine bool) float32 {
	temp := math.Max(0, math.Min((float64(manCov)+GetDelta(isCombine))+(float64(zonCov)+GetDelta(isCombine)), 150.0))
	return float32((temp / 150.0) * 10.0)
}

func RunManCoverage(manCov int, isCombine bool) float32 {
	temp := math.Max(0, math.Min(float64(manCov)+GetDelta(isCombine), 90.0))
	return float32((temp / 90.0) * 10.0)
}

func RunZoneCoverage(zonCov int, isCombine bool) float32 {
	temp := math.Max(0, math.Min(float64(zonCov)+GetDelta(isCombine), 80.0))
	return float32((temp / 80.0) * 10.0)
}

func RunKickoffDrill(kickPow int, puntPow int, isCombine bool) float32 {
	temp := math.Max(0, math.Min(math.Max(float64(kickPow), float64(puntPow))+GetDelta(isCombine), 75.0))
	return float32((temp / 75.0) * 10.0)
}

func RunFieldGoalDrill(kickPow int, kickAcc int, isCombine bool) float32 {
	temp := math.Max(0, math.Min((float64(kickPow)+GetDelta(isCombine))+(float64(kickAcc)+GetDelta(isCombine)), 155.0))
	return float32((temp / 155.0) * 10.0)
}

func RunPuntDistance(puntPow int, isCombine bool) float32 {
	temp := math.Max(0, math.Min(float64(puntPow)+GetDelta(isCombine), 60.0))
	return float32((temp / 60.0) * 10.0)
}

func RunCoffinPunt(puntAcc int, isCombine bool) float32 {
	temp := math.Max(0, math.Min(float64(puntAcc)+GetDelta(isCombine), 66.0))
	return float32((temp / 66.0) * 10.0)
}

func GetDelta(isCombine bool) float64 {
	min, max := -5, 15
	if isCombine {
		min, max = -10, 10
	}
	return float64(rand.Intn(max-min) + min)
}

func AddParticipants(json map[string][]uint, events []models.PreDraftEvent, players []models.NFLDraftee) []models.PreDraftEvent {
	for i, event := range events {
		for _, playerID := range json[event.Name] {
			if participant, found := FindParticipant(playerID, players); found {
				events[i].Participants = append(events[i].Participants, participant)
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
	names := []string{"AAC Pro Day", "ACC Pro Day", "Big Ten Pro Day", "Big XII Pro Day", "C-USA Pro Day", "Independents Pro Day", "MAC Pro Day", "MWC Pro Day", "PAC-12 Pro Day", "SEC Pro Day", "Sun Belt Pro Day", "FCS Pro Day", "NFL Combine"}
	var tempList []models.PreDraftEvent
	for _, name := range names {
		tempList = append(tempList, models.PreDraftEvent{Name: name, IsCombine: name == "NFL Combine"})
	}
	return tempList
}
