package managers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

type CrootGenerator struct {
	firstNameMap      map[string][][]string
	lastNameMap       map[string][][]string
	collegePlayerList []structs.CollegePlayer
	coachList         []structs.CollegeCoach
	teamMap           map[uint]structs.CollegeTeam
	positionList      []string
	CrootList         []structs.Recruit
	FacesList         []structs.FaceData
	GlobalList        []structs.Player
	attributeBlob     map[string]map[string]map[string]map[string]interface{}
	crootLocations    map[string][]structs.CrootLocation
	faceDataBlob      map[string][]string
	newID             uint
	count             int
	requiredPlayers   int
	qbCount           int
	rbCount           int
	fbCount           int
	wrCount           int
	teCount           int
	otCount           int
	ogCount           int
	cCount            int
	dtCount           int
	deCount           int
	ilbCount          int
	olbCount          int
	cbCount           int
	fsCount           int
	ssCount           int
	pCount            int
	kCount            int
	athCount          int
	star5             int
	star4             int
	star3             int
	star2             int
	star1             int
	highestOvr        int
	lowestOvr         int
	pickedEthnicity   string
	caser             cases.Caser
}

func (pg *CrootGenerator) GenerateRecruits() {
	for pg.count < pg.requiredPlayers {
		player, globalPlayer := pg.generatePlayer()
		pg.CrootList = append(pg.CrootList, player)
		pg.GlobalList = append(pg.GlobalList, globalPlayer)
		pg.updateStatistics(player) // A method to update player counts and statistics
		if player.RelativeType == 5 {
			twinPlayer, twinGlobal := pg.generateTwin(&player)
			pg.updateStatistics(twinPlayer)
			pg.CrootList = append(pg.CrootList, twinPlayer)
			pg.GlobalList = append(pg.GlobalList, twinGlobal)
			pg.count++
		}
		pg.count++
		pg.newID++
	}
}

func (pg *CrootGenerator) generatePlayer() (structs.Recruit, structs.Player) {
	cpLen := len(pg.collegePlayerList) - 1
	coachLen := len(pg.coachList)
	relativeType := 0
	relativeID := 0
	coachTeamID := 0
	coachTeamAbbr := ""
	notes := ""
	star := util.GetStarRating()
	state := util.PickState()
	pg.pickedEthnicity = pickEthnicity()
	firstNameList := pg.firstNameMap[pg.pickedEthnicity]
	lastNameList := pg.lastNameMap[pg.pickedEthnicity]
	fName := getName(firstNameList)
	firstName := pg.caser.String(strings.ToLower(fName))
	lastName := ""
	roof := 100
	relativeRoll := util.GenerateIntFromRange(1, roof)
	relativeIdx := 0
	if relativeRoll == roof {
		relativeType = getRelativeType()
		if relativeType == 2 {
			// Brother of college player
			fmt.Println("BROTHER")
			relativeIdx = util.GenerateIntFromRange(0, cpLen)
			if relativeIdx < 0 || relativeIdx > len(pg.collegePlayerList) {
				relativeIdx = util.GenerateIntFromRange(0, cpLen)
			}
			cp := pg.collegePlayerList[relativeIdx]
			relativeID = int(cp.ID)
			lastName = cp.LastName
			state = cp.State
			notes = "Brother of " + cp.TeamAbbr + " " + cp.Position + " " + cp.FirstName + " " + cp.LastName
		} else if relativeType == 3 {
			fmt.Println("COUSIN")
			// Cousin
			relativeIdx = util.GenerateIntFromRange(0, cpLen)
			if relativeIdx < 0 || relativeIdx > len(pg.collegePlayerList) {
				relativeIdx = util.GenerateIntFromRange(0, cpLen)
			}
			cp := pg.collegePlayerList[relativeIdx]
			relativeID = int(cp.ID)
			coinFlip := util.GenerateIntFromRange(1, 2)
			if coinFlip == 1 {
				lastName = cp.LastName
			} else {
				lName := getName(lastNameList)
				lastName = pg.caser.String(strings.ToLower(lName))
			}
			state = cp.State
			notes = "Cousin of " + cp.TeamAbbr + " " + cp.Position + " " + cp.FirstName + " " + cp.LastName
		} else if relativeType == 4 {
			// Half Brother
			fmt.Println("HALF BROTHER GENERATED")
			relativeIdx = util.GenerateIntFromRange(0, cpLen)
			if relativeIdx < 0 || relativeIdx > len(pg.collegePlayerList) {
				relativeIdx = util.GenerateIntFromRange(0, cpLen)
			}
			cp := pg.collegePlayerList[relativeIdx]
			relativeID = int(cp.ID)
			coinFlip := util.GenerateIntFromRange(1, 3)
			if coinFlip < 3 {
				lastName = cp.LastName
			} else {
				lName := getName(lastNameList)
				lastName = pg.caser.String(strings.ToLower(lName))
			}
			state = cp.State
			notes = "Half-Brother of " + cp.TeamAbbr + " " + cp.Position + " " + cp.FirstName + " " + cp.LastName
		} else if relativeType == 5 {
			// Twin
			relativeType = 5
			relativeID = int(pg.newID)
		} else if relativeType == 6 {
			// Coach's Son
			fmt.Println("COACH'S SON")
			relativeIdx = util.GenerateIntFromRange(0, coachLen)
			if relativeIdx < 0 || relativeIdx > len(pg.coachList) {
				relativeIdx = util.GenerateIntFromRange(0, coachLen)
			}
			coach := pg.coachList[relativeIdx]
			relativeID = int(coach.ID)
			lastName = getCoachLastName(coach.CoachName)
			team := pg.teamMap[coach.TeamID]
			if len(team.State) > 2 {
				stateLabel, err := util.GetStateAbbreviation(team.State)
				if err != nil {
					log.Panicln("WRONG STATE INPUT")
				}
				state = stateLabel
			} else {
				state = team.State
			}
			notes = "Son of Coach " + coach.CoachName + " of " + coach.Team
			coachTeamID = int(coach.TeamID)
			coachTeamAbbr = coach.Team
		} else if relativeType == 7 {
			// Coach's Nephew
			fmt.Println("COACH'S NEPHEW")
			relativeIdx = util.GenerateIntFromRange(0, coachLen)
			if relativeIdx < 0 || relativeIdx > len(pg.coachList) {
				relativeIdx = util.GenerateIntFromRange(0, coachLen)
			}
			coach := pg.coachList[relativeIdx]
			relativeID = int(coach.ID)
			coachLastName := getCoachLastName(coach.CoachName)
			coinFlip := util.GenerateIntFromRange(1, 2)
			if coinFlip == 1 {
				lastName = coachLastName
			} else {
				lName := getName(lastNameList)
				lastName = pg.caser.String(strings.ToLower(lName))
			}
			team := pg.teamMap[coach.TeamID]
			if len(team.State) > 2 {
				stateLabel, err := util.GetStateAbbreviation(team.State)
				if err != nil {
					log.Panicln("WRONG STATE INPUT")
				}
				state = stateLabel
			} else {
				state = team.State
			}
			notes = "Nephew of Coach " + coach.CoachName + " of " + coach.Team
			coachTeamID = int(coach.TeamID)
			coachTeamAbbr = coach.Team
		}
	} else {
		relativeType = 1
	}
	if relativeType == 1 || relativeType == 5 {
		lName := getName(lastNameList)
		lastName = pg.caser.String(strings.ToLower(lName))
	}

	pickedPosition := util.PickPosition()
	player := createRecruit(pickedPosition, star, firstName, lastName, pg.attributeBlob, state, pg.crootLocations[state])
	player.AssignRelativeData(uint(relativeID), uint(relativeType), uint(coachTeamID), coachTeamAbbr, notes)
	globalPlayer := structs.Player{
		CollegePlayerID: int(pg.newID),
		RecruitID:       int(pg.newID),
		NFLPlayerID:     int(pg.newID),
	}

	globalPlayer.AssignID(pg.newID)

	skinColor := getSkinColorByEthnicity(pg.pickedEthnicity)

	face := getFace(pg.newID, player.Weight, skinColor, "", pg.faceDataBlob)

	pg.FacesList = append(pg.FacesList, face)

	return player, globalPlayer
}

func (pg *CrootGenerator) generateTwin(player *structs.Recruit) (structs.Recruit, structs.Player) {
	fmt.Println("TWIN!!")
	// Generate Twin Record
	firstTwinRelativeID := int(pg.newID)
	pg.newID++
	// Twin being generated is secondTwin
	secondTwinRelativeID := pg.newID
	firstNameList := pg.firstNameMap[pg.pickedEthnicity]
	twinName := getName(firstNameList)
	twinN := pg.caser.String(strings.ToLower(twinName))
	twinPosition := util.PickFromStringList(pg.positionList)
	coinFlip := util.GenerateIntFromRange(1, 2)
	stars := util.GetStarRating()
	if coinFlip == 2 {
		twinPosition = player.Position
		stars = player.Stars
	}
	twinNotes := "Twin Brother of " + strconv.Itoa(player.Stars) + " Star Recruit " + player.Position + " " + player.FirstName + " " + player.LastName
	twinPlayer := createRecruit(twinPosition, stars, twinN, player.LastName, pg.attributeBlob, player.State, pg.crootLocations[player.State])
	twinPlayer.AssignRelativeData(uint(firstTwinRelativeID), 4, 0, "", twinNotes)
	twinPlayer.AssignTwinData(player.LastName, player.City, player.State, player.HighSchool)
	notes := "Twin Brother of " + strconv.Itoa(twinPlayer.Stars) + " Star Recruit " + twinPlayer.Position + " " + twinPlayer.FirstName + " " + twinPlayer.LastName
	player.AssignRelativeData(uint(secondTwinRelativeID), 4, 0, "", notes)
	globalTwinPlayer := structs.Player{
		CollegePlayerID: int(secondTwinRelativeID),
		RecruitID:       int(secondTwinRelativeID),
		NFLPlayerID:     int(secondTwinRelativeID),
	}
	globalTwinPlayer.AssignID(secondTwinRelativeID)
	globalPlayer := structs.Player{
		CollegePlayerID: firstTwinRelativeID,
		RecruitID:       firstTwinRelativeID,
		NFLPlayerID:     firstTwinRelativeID,
	}
	globalPlayer.AssignID(uint(firstTwinRelativeID))
	skinColor := getSkinColorByEthnicity(pg.pickedEthnicity)

	face := getFace(secondTwinRelativeID, twinPlayer.Weight, skinColor, "", pg.faceDataBlob)

	pg.FacesList = append(pg.FacesList, face)
	return twinPlayer, globalTwinPlayer
}

func (pg *CrootGenerator) updateStatistics(player structs.Recruit) {
	if player.Stars == 5 {
		pg.star5++
	} else if player.Stars == 4 {
		pg.star4++
	} else if player.Stars == 3 {
		pg.star3++
	} else if player.Stars == 2 {
		pg.star2++
	} else {
		pg.star1++
	}
	if player.Position == "QB" {
		pg.qbCount++
	} else if player.Position == "RB" {
		pg.rbCount++
	} else if player.Position == "FB" {
		pg.fbCount++
	} else if player.Position == "WR" {
		pg.wrCount++
	} else if player.Position == "TE" {
		pg.teCount++
	} else if player.Position == "OT" {
		pg.otCount++
	} else if player.Position == "OG" {
		pg.ogCount++
	} else if player.Position == "C" {
		pg.cCount++
	} else if player.Position == "DT" {
		pg.dtCount++
	} else if player.Position == "DE" {
		pg.deCount++
	} else if player.Position == "ILB" {
		pg.ilbCount++
	} else if player.Position == "OLB" {
		pg.olbCount++
	} else if player.Position == "CB" {
		pg.cbCount++
	} else if player.Position == "FS" {
		pg.fsCount++
	} else if player.Position == "SS" {
		pg.ssCount++
	} else if player.Position == "K" {
		pg.kCount++
	} else if player.Position == "P" {
		pg.pCount++
	} else if player.Position == "ATH" {
		pg.athCount++
	}

	if player.Overall > pg.highestOvr {
		pg.highestOvr = player.Overall
	}
	if player.Overall < pg.lowestOvr {
		pg.lowestOvr = player.Overall
	}
}

func (pg *CrootGenerator) OutputRecruitStats() {
	// Croot Stats
	fmt.Println("Total Recruit Count: ", pg.count)
	fmt.Println("Total 5 Star  Count: ", pg.star5)
	fmt.Println("Total 4 Star  Count: ", pg.star4)
	fmt.Println("Total 3 Star  Count: ", pg.star3)
	fmt.Println("Total 2 Star  Count: ", pg.star2)
	fmt.Println("Total 1 Star  Count: ", pg.star1)
	fmt.Println("Total QB  Count: ", pg.qbCount)
	fmt.Println("Total RB  Count: ", pg.rbCount)
	fmt.Println("Total FB  Count: ", pg.fbCount)
	fmt.Println("Total WR  Count: ", pg.wrCount)
	fmt.Println("Total TE  Count: ", pg.teCount)
	fmt.Println("Total OT  Count: ", pg.otCount)
	fmt.Println("Total OG  Count: ", pg.ogCount)
	fmt.Println("Total C  Count: ", pg.cCount)
	fmt.Println("Total DT  Count: ", pg.dtCount)
	fmt.Println("Total DE  Count: ", pg.deCount)
	fmt.Println("Total ILB  Count: ", pg.ilbCount)
	fmt.Println("Total OLB  Count: ", pg.olbCount)
	fmt.Println("Total CB  Count: ", pg.cbCount)
	fmt.Println("Total FS  Count: ", pg.fsCount)
	fmt.Println("Total SS  Count: ", pg.ssCount)
	fmt.Println("Total K  Count: ", pg.kCount)
	fmt.Println("Total P  Count: ", pg.pCount)
	fmt.Println("Total ATH  Count: ", pg.athCount)

	fmt.Println("Highest Recruit Ovr: ", pg.highestOvr)
	fmt.Println("Lowest  Recruit Ovr: ", pg.lowestOvr)
}

func GenerateCroots() {
	db := dbprovider.GetInstance().GetDB()
	var lastPlayerRecord structs.Player
	ts := GetTimestamp()

	if ts.CrootsGenerated {
		return
	}

	err := db.Last(&lastPlayerRecord).Error
	if err != nil {
		log.Fatalln("Could not grab last player record from players table...")
	}

	fNameMap, lNameMap := getNameMaps()
	generator := CrootGenerator{
		firstNameMap:      fNameMap,
		lastNameMap:       lNameMap,
		collegePlayerList: GetAllCollegePlayers(),
		coachList:         GetAllAICollegeCoaches(),
		teamMap:           GetCollegeTeamMap(),
		crootLocations:    getCrootLocations(),
		attributeBlob:     getAttributeBlob(),
		positionList:      util.GetPositionList(),
		newID:             lastPlayerRecord.ID + 1,
		requiredPlayers:   util.GenerateIntFromRange(6400, 6601),
		faceDataBlob:      getFaceDataBlob(),
		count:             0,
		star5:             0,
		star4:             0,
		star3:             0,
		star2:             0,
		star1:             0,
		highestOvr:        0,
		lowestOvr:         100000,
		CrootList:         []structs.Recruit{},
		GlobalList:        []structs.Player{},
		FacesList:         []structs.FaceData{},
		caser:             cases.Title(language.English),
		pickedEthnicity:   "",
	}

	// The plan is to ensure that every recruit is signed
	generator.GenerateRecruits()
	// Croot Stats
	generator.OutputRecruitStats()

	repository.CreateCFBRecruitRecordsBatch(db, generator.CrootList, 500)
	repository.CreateGlobalPlayerRecordsBatch(db, generator.GlobalList, 500)
	repository.CreateFaceRecordsBatch(db, generator.FacesList, 500)
	ts.ToggleGeneratedCroots()
	repository.SaveTimestamp(ts, db)
	AssignAllRecruitRanks()
}

func GenerateWalkOns() {
	fmt.Println(time.Now().UnixNano())
	db := dbprovider.GetInstance().GetDB()
	AllTeams := GetRecruitingProfileForRecruitSync()
	ts := GetTimestamp()
	count := 0
	attributeBlob := getAttributeBlob()
	highSchoolBlob := getCrootLocations()
	faceDataBlob := getFaceDataBlob()
	faces := []structs.FaceData{}
	recruitBatchList := []structs.Recruit{}
	recruitProfileBatchList := []structs.RecruitPlayerProfile{}
	globalPlayerList := []structs.Player{}
	firstNameMap, lastNameMap := getNameMaps()
	collegePlayers := GetAllCollegePlayers()
	collegeRosterMap := MakeCollegePlayerMapByTeamID(collegePlayers, true)

	newID := getLatestRecord(db)

	for _, team := range AllTeams {
		roster := collegeRosterMap[team.ID]
		rosterCount := len(roster)
		for _, p := range roster {
			if p.Year == 5 || (p.Year == 4 && !p.IsRedshirt) {
				rosterCount--
			}
		}
		rosterCount += team.TotalCommitments
		rosterLimit := 105
		if !team.IsFBS {
			rosterLimit = 80
		}
		if rosterCount >= rosterLimit {
			continue
		}
		positionList := []string{"QB", "FB", "RB", "TE", "WR", "OT", "OG", "C", "DT", "DE", "OLB", "ILB", "CB", "FS", "SS", "K", "P", "ATH"}

		// Get Team Needs
		walkonLimit := rosterLimit - rosterCount

		// Randomize List
		rand.Shuffle(len(positionList), func(i, j int) {
			positionList[i], positionList[j] = positionList[j], positionList[i]
		})

		// Recruit Generation
		for _, pos := range positionList {
			if count >= walkonLimit {
				break
			}

			ethnicity := pickEthnicity()

			state := pickWalkonState(team.State)

			recruit := createWalkon(pos, firstNameMap[ethnicity], lastNameMap[ethnicity], attributeBlob, state, highSchoolBlob[state])

			recruit.AssignWalkon(team.TeamAbbreviation, int(team.ID), newID)

			recruitPlayerRecord := structs.RecruitPlayerProfile{
				ProfileID:        int(team.ID),
				RecruitID:        int(newID),
				IsSigned:         true,
				IsLocked:         true,
				TeamAbbreviation: team.TeamAbbreviation,
				SeasonID:         ts.CollegeSeasonID,
				TotalPoints:      1,
			}

			playerRecord := structs.Player{
				RecruitID:       int(newID),
				CollegePlayerID: int(newID),
				NFLPlayerID:     int(newID),
			}
			playerRecord.AssignID(newID)
			count++

			skinColor := getSkinColorByEthnicity(ethnicity)

			face := getFace(newID, recruit.Weight, skinColor, "", faceDataBlob)
			faces = append(faces, face)
			globalPlayerList = append(globalPlayerList, playerRecord)
			recruitBatchList = append(recruitBatchList, recruit)
			recruitProfileBatchList = append(recruitProfileBatchList, recruitPlayerRecord)
			newID++
			team.IncreaseCommitCount()
			repository.SaveRecruitingTeamProfile(team, db)

		}
		count = 0
		fmt.Println("Finished walkon generation for " + team.TeamAbbreviation)
	}
	repository.CreateFaceRecordsBatch(db, faces, 500)
	repository.CreateCFBRecruitRecordsBatch(db, recruitBatchList, 500)
	repository.CreateCFBRecruitProfileRecordsBatch(db, recruitProfileBatchList, 500)
	repository.CreateGlobalPlayerRecordsBatch(db, globalPlayerList, 500)
}

func CreateCustomCroots() {
	db := dbprovider.GetInstance().GetDB()
	path := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\2026\\2026_Custom_Croot_Class.csv"
	crootCSV := util.ReadCSV(path)
	attributeBlob := getAttributeBlob()
	latestID := getLatestRecord(db)

	crootList := []structs.Recruit{}
	globalList := []structs.Player{}
	faceData := []structs.FaceData{}
	faceDataBlob := getFaceDataBlob()

	for idx, row := range crootCSV {
		if idx < 1 {
			continue
		}
		if row[0] == "" {
			break
		}
		croot, face := createCustomCroot(row, latestID, attributeBlob, faceDataBlob)
		croot.AssignID(int(latestID))
		gp := structs.Player{
			CollegePlayerID: int(croot.ID),
			NFLPlayerID:     int(croot.ID),
			RecruitID:       int(croot.ID),
		}

		gp.AssignID(croot.ID)
		latestID++
		globalList = append(globalList, gp)
		crootList = append(crootList, croot)
		faceData = append(faceData, face)
	}

	repository.CreateFaceRecordsBatch(db, faceData, 100)
	repository.CreateCFBRecruitRecordsBatch(db, crootList, 100)
	repository.CreateGlobalPlayerRecordsBatch(db, globalList, 100)
}

func GenerateCoachesForAITeams() {
	db := dbprovider.GetInstance().GetDB()

	teams := GetOnlyAITeamRecruitingProfiles()
	firstNameMap, lastNameMap := getNameMaps()

	coachList := []structs.CollegeCoach{}
	allActiveCoaches := GetAllAICollegeCoaches()

	retiredPlayers := GetRetiredSimNFLPlayers()
	retireeMap := make(map[uint]bool)
	coachMap := make(map[uint]bool)

	for _, coach := range allActiveCoaches {
		if coach.FormerPlayerID > 0 {
			coachMap[coach.FormerPlayerID] = true
		}
	}

	for _, team := range teams {
		// Skip over teams currently controlled by a user
		if !team.IsAI || team.IsUserTeam {
			continue
		}

		pickedEthnicity := pickEthnicity()
		almaMater := pickAlmaMater(teams)
		coach := createCollegeCoach(team, almaMater.ID, almaMater.TeamAbbreviation, firstNameMap[pickedEthnicity], lastNameMap[pickedEthnicity], retiredPlayers, &retireeMap, &coachMap)
		team.UpdateAIBehavior(true, true, coach.StarMax, coach.StarMin, coach.PointMin, coach.PointMax, coach.OffensiveScheme, coach.DefensiveScheme)
		team.AssignRecruiter(coach.CoachName)
		coachList = append(coachList, coach)

		db.Save(&team)
	}

	for _, coach := range coachList {
		db.Create(&coach)
	}
}

func createRecruit(position string, stars int, firstName, lastName string, blob map[string]map[string]map[string]map[string]interface{}, state string, hsBlob []structs.CrootLocation) structs.Recruit {
	age := 18
	city, highSchool := getCityAndHighSchool(hsBlob)

	archetype := getArchetype(position)
	height := getAttributeValue(position, archetype, stars, "Height", blob)
	weight := getAttributeValue(position, archetype, stars, "Weight", blob)
	footballIQ := getAttributeValue(position, archetype, stars, "Football IQ", blob)
	speed := getAttributeValue(position, archetype, stars, "Speed", blob)
	agility := getAttributeValue(position, archetype, stars, "Agility", blob)
	carrying := getAttributeValue(position, archetype, stars, "Carrying", blob)
	catching := getAttributeValue(position, archetype, stars, "Catching", blob)
	routeRunning := getAttributeValue(position, archetype, stars, "Route Running", blob)
	zoneCoverage := getAttributeValue(position, archetype, stars, "Zone Coverage", blob)
	manCoverage := getAttributeValue(position, archetype, stars, "Man Coverage", blob)
	strength := getAttributeValue(position, archetype, stars, "Strength", blob)
	tackle := getAttributeValue(position, archetype, stars, "Tackle", blob)
	passBlock := getAttributeValue(position, archetype, stars, "Pass Block", blob)
	runBlock := getAttributeValue(position, archetype, stars, "Run Block", blob)
	passRush := getAttributeValue(position, archetype, stars, "Pass Rush", blob)
	runDefense := getAttributeValue(position, archetype, stars, "Run Defense", blob)
	throwPower := getAttributeValue(position, archetype, stars, "Throw Power", blob)
	throwAccuracy := getAttributeValue(position, archetype, stars, "Throw Accuracy", blob)
	kickAccuracy := getAttributeValue(position, archetype, stars, "Kick Accuracy", blob)
	kickPower := getAttributeValue(position, archetype, stars, "Kick Power", blob)
	puntAccuracy := getAttributeValue(position, archetype, stars, "Punt Accuracy", blob)
	puntPower := getAttributeValue(position, archetype, stars, "Punt Power", blob)
	injury := util.GenerateNormalizedIntFromMeanStdev(50, 15)
	stamina := util.GenerateNormalizedIntFromMeanStdev(50, 15)
	discipline := util.GenerateNormalizedIntFromMeanStdev(50, 15)
	progression := util.GenerateNormalizedIntFromMeanStdev(50, 15)

	freeAgency := util.GetFreeAgencyBias()
	personality := util.GetPersonality()
	recruitingBias := util.GetRecruitingBias()
	workEthic := util.GetWorkEthic()
	academicBias := util.GetAcademicBias()
	potentialGrade := util.GetWeightedPotentialGrade(int(progression))

	affinityOne := util.PickAffinity(stars, "", false)
	affinityTwo := util.PickAffinity(stars, affinityOne, true)

	basePlayer := structs.BasePlayer{
		FirstName:      firstName,
		LastName:       lastName,
		Position:       position,
		Archetype:      archetype,
		Age:            age,
		Stars:          stars,
		Height:         height,
		Weight:         weight,
		Stamina:        int(stamina),
		Injury:         int(injury),
		FootballIQ:     footballIQ,
		Speed:          speed,
		Carrying:       carrying,
		Agility:        agility,
		Catching:       catching,
		RouteRunning:   routeRunning,
		ZoneCoverage:   zoneCoverage,
		ManCoverage:    manCoverage,
		Strength:       strength,
		Tackle:         tackle,
		PassBlock:      passBlock,
		RunBlock:       runBlock,
		PassRush:       passRush,
		RunDefense:     runDefense,
		ThrowPower:     throwPower,
		ThrowAccuracy:  throwAccuracy,
		KickAccuracy:   kickAccuracy,
		KickPower:      kickPower,
		PuntAccuracy:   puntAccuracy,
		PuntPower:      puntPower,
		Progression:    int(progression),
		Discipline:     int(discipline),
		PotentialGrade: potentialGrade,
		FreeAgency:     freeAgency,
		Personality:    personality,
		RecruitingBias: recruitingBias,
		WorkEthic:      workEthic,
		AcademicBias:   academicBias,
	}

	basePlayer.GetOverall()

	return structs.Recruit{
		BasePlayer:  basePlayer,
		City:        city,
		HighSchool:  highSchool,
		State:       state,
		IsSigned:    false,
		AffinityOne: affinityOne,
		AffinityTwo: affinityTwo,
	}
}

func createWalkon(position string, firstNameList [][]string, lastNameList [][]string, blob map[string]map[string]map[string]map[string]interface{}, state string, hsBlob []structs.CrootLocation) structs.Recruit {
	fName := getName(firstNameList)
	lName := getName(lastNameList)
	firstName := strings.Title(strings.ToLower(fName))
	lastName := strings.Title(strings.ToLower(lName))
	age := 18
	city, highSchool := getCityAndHighSchool(hsBlob)

	archetype := getArchetype(position)
	stars := getWalkonStarRating()
	height := getAttributeValue(position, archetype, stars, "Height", blob)
	weight := getAttributeValue(position, archetype, stars, "Weight", blob)
	footballIQ := getAttributeValue(position, archetype, stars, "Football IQ", blob)
	speed := getAttributeValue(position, archetype, stars, "Speed", blob)
	agility := getAttributeValue(position, archetype, stars, "Agility", blob)
	carrying := getAttributeValue(position, archetype, stars, "Carrying", blob)
	catching := getAttributeValue(position, archetype, stars, "Catching", blob)
	routeRunning := getAttributeValue(position, archetype, stars, "Route Running", blob)
	zoneCoverage := getAttributeValue(position, archetype, stars, "Zone Coverage", blob)
	manCoverage := getAttributeValue(position, archetype, stars, "Man Coverage", blob)
	strength := getAttributeValue(position, archetype, stars, "Strength", blob)
	tackle := getAttributeValue(position, archetype, stars, "Tackle", blob)
	passBlock := getAttributeValue(position, archetype, stars, "Pass Block", blob)
	runBlock := getAttributeValue(position, archetype, stars, "Run Block", blob)
	passRush := getAttributeValue(position, archetype, stars, "Pass Rush", blob)
	runDefense := getAttributeValue(position, archetype, stars, "Run Defense", blob)
	throwPower := getAttributeValue(position, archetype, stars, "Throw Power", blob)
	throwAccuracy := getAttributeValue(position, archetype, stars, "Throw Accuracy", blob)
	kickAccuracy := getAttributeValue(position, archetype, stars, "Kick Accuracy", blob)
	kickPower := getAttributeValue(position, archetype, stars, "Kick Power", blob)
	puntAccuracy := getAttributeValue(position, archetype, stars, "Punt Accuracy", blob)
	puntPower := getAttributeValue(position, archetype, stars, "Punt Power", blob)
	injury := util.GenerateNormalizedIntFromMeanStdev(50, 15)
	stamina := util.GenerateNormalizedIntFromMeanStdev(50, 15)
	discipline := util.GenerateNormalizedIntFromMeanStdev(50, 15)
	progression := util.GenerateNormalizedIntFromMeanStdev(50, 15)

	freeAgency := util.GetFreeAgencyBias()
	personality := util.GetPersonality()
	recruitingBias := util.GetRecruitingBias()
	workEthic := util.GetWorkEthic()
	academicBias := util.GetAcademicBias()
	potentialGrade := util.GetWeightedPotentialGrade(int(progression))

	basePlayer := structs.BasePlayer{
		FirstName:      firstName,
		LastName:       lastName,
		Position:       position,
		Archetype:      archetype,
		Age:            age,
		Stars:          0,
		Height:         height,
		Weight:         weight,
		Stamina:        int(stamina),
		Injury:         int(injury),
		FootballIQ:     footballIQ,
		Speed:          speed,
		Carrying:       carrying,
		Agility:        agility,
		Catching:       catching,
		RouteRunning:   routeRunning,
		ZoneCoverage:   zoneCoverage,
		ManCoverage:    manCoverage,
		Strength:       strength,
		Tackle:         tackle,
		PassBlock:      passBlock,
		RunBlock:       runBlock,
		PassRush:       passRush,
		RunDefense:     runDefense,
		ThrowPower:     throwPower,
		ThrowAccuracy:  throwAccuracy,
		KickAccuracy:   kickAccuracy,
		KickPower:      kickPower,
		PuntAccuracy:   puntAccuracy,
		PuntPower:      puntPower,
		Progression:    int(progression),
		Discipline:     int(discipline),
		PotentialGrade: potentialGrade,
		FreeAgency:     freeAgency,
		Personality:    personality,
		RecruitingBias: recruitingBias,
		WorkEthic:      workEthic,
		AcademicBias:   academicBias,
	}

	basePlayer.GetOverall()

	return structs.Recruit{
		BasePlayer:      basePlayer,
		City:            city,
		HighSchool:      highSchool,
		State:           state,
		IsSigned:        true,
		TopRankModifier: 1,
	}
}

func createCustomCroot(croot []string, id uint, blob map[string]map[string]map[string]map[string]interface{}, faceDataBlob map[string][]string) (structs.Recruit, structs.FaceData) {
	firstName := croot[0]
	lastName := croot[1]
	position := croot[2]
	archetype := croot[3]
	// stars := 5
	stars := getCustomCrootStarRating()
	height := util.ConvertStringToInt(croot[4])
	weight := util.ConvertStringToInt(croot[5])
	city := croot[6]
	highSchool := croot[7]
	state := croot[8]
	crootFor := croot[9]
	relativeID := croot[10]
	relativeType := croot[11]
	notes := croot[12]
	affinityOne := croot[13]
	affinityTwo := croot[14]
	ethnicity := croot[15]
	if ethnicity == "" {
		ethnicity = pickEthnicity()
	}
	gender := croot[16]
	hasNoAffinities := affinityOne == "" && affinityTwo == ""
	age := 18
	footballIQ := getAttributeValue(position, archetype, stars, "Football IQ", blob)
	speed := getAttributeValue(position, archetype, stars, "Speed", blob)
	agility := getAttributeValue(position, archetype, stars, "Agility", blob)
	carrying := getAttributeValue(position, archetype, stars, "Carrying", blob)
	catching := getAttributeValue(position, archetype, stars, "Catching", blob)
	routeRunning := getAttributeValue(position, archetype, stars, "Route Running", blob)
	zoneCoverage := getAttributeValue(position, archetype, stars, "Zone Coverage", blob)
	manCoverage := getAttributeValue(position, archetype, stars, "Man Coverage", blob)
	strength := getAttributeValue(position, archetype, stars, "Strength", blob)
	tackle := getAttributeValue(position, archetype, stars, "Tackle", blob)
	passBlock := getAttributeValue(position, archetype, stars, "Pass Block", blob)
	runBlock := getAttributeValue(position, archetype, stars, "Run Block", blob)
	passRush := getAttributeValue(position, archetype, stars, "Pass Rush", blob)
	runDefense := getAttributeValue(position, archetype, stars, "Run Defense", blob)
	throwPower := getAttributeValue(position, archetype, stars, "Throw Power", blob)
	throwAccuracy := getAttributeValue(position, archetype, stars, "Throw Accuracy", blob)
	kickAccuracy := getAttributeValue(position, archetype, stars, "Kick Accuracy", blob)
	kickPower := getAttributeValue(position, archetype, stars, "Kick Power", blob)
	puntAccuracy := getAttributeValue(position, archetype, stars, "Punt Accuracy", blob)
	puntPower := getAttributeValue(position, archetype, stars, "Punt Power", blob)
	injury := util.GenerateIntFromRange(40, 100)
	stamina := util.GenerateIntFromRange(40, 100)
	discipline := util.GenerateIntFromRange(40, 100)
	progression := util.GenerateIntFromRange(40, 80)
	freeAgency := util.GetFreeAgencyBias()
	personality := util.GetPersonality()
	recruitingBias := util.GetRecruitingBias()
	workEthic := util.GetWorkEthic()
	academicBias := util.GetAcademicBias()
	potentialGrade := util.GetWeightedPotentialGrade(progression)
	if hasNoAffinities {
		affinityOne = util.PickAffinity(stars, "", false)
		affinityTwo = util.PickAffinity(stars, affinityOne, true)
	}

	basePlayer := structs.BasePlayer{
		FirstName:      firstName,
		LastName:       lastName,
		Position:       position,
		Archetype:      archetype,
		Age:            age,
		Stars:          stars,
		Height:         height,
		Weight:         weight,
		Stamina:        stamina,
		Injury:         injury,
		FootballIQ:     footballIQ,
		Speed:          speed,
		Carrying:       carrying,
		Agility:        agility,
		Catching:       catching,
		RouteRunning:   routeRunning,
		ZoneCoverage:   zoneCoverage,
		ManCoverage:    manCoverage,
		Strength:       strength,
		Tackle:         tackle,
		PassBlock:      passBlock,
		RunBlock:       runBlock,
		PassRush:       passRush,
		RunDefense:     runDefense,
		ThrowPower:     throwPower,
		ThrowAccuracy:  throwAccuracy,
		KickAccuracy:   kickAccuracy,
		KickPower:      kickPower,
		PuntAccuracy:   puntAccuracy,
		PuntPower:      puntPower,
		Progression:    progression,
		Discipline:     discipline,
		PotentialGrade: potentialGrade,
		FreeAgency:     freeAgency,
		Personality:    personality,
		RecruitingBias: recruitingBias,
		WorkEthic:      workEthic,
		AcademicBias:   academicBias,
		RelativeID:     uint(util.ConvertStringToInt(relativeID)),
		RelativeType:   uint(util.ConvertStringToInt(relativeType)),
		Notes:          notes,
	}

	basePlayer.GetOverall()

	skinColor := getSkinColorByEthnicity(ethnicity)

	faceData := getFace(id, weight, skinColor, gender, faceDataBlob)
	recruit := structs.Recruit{
		BasePlayer:     basePlayer,
		PlayerID:       int(id),
		City:           city,
		HighSchool:     highSchool,
		State:          state,
		IsSigned:       false,
		IsCustomCroot:  true,
		CustomCrootFor: crootFor,
		AffinityOne:    affinityOne,
		AffinityTwo:    affinityTwo,
	}

	return recruit, faceData
}

func createCollegeCoach(team structs.RecruitingTeamProfile, almaMaterID uint, almaMater string, firstNameList, lastNameList [][]string, retiredPlayers []structs.NFLRetiredPlayer, retireeMap, coachMap *map[uint]bool) structs.CollegeCoach {
	firstName := ""
	lastName := ""
	diceRoll := util.GenerateIntFromRange(1, 50)
	formerPlayerID := uint(0)
	almaID := almaMaterID
	alma := almaMater
	age := 32
	posOne := ""
	posTwo := ""
	posThree := ""
	if diceRoll == 50 {
		// Get a former player as a coach
		idx := util.GenerateIntFromRange(0, len(retiredPlayers)-1)
		retiree := retiredPlayers[idx]
		for (*retireeMap)[retiree.ID] || (*coachMap)[retiree.ID] {
			idx = util.GenerateIntFromRange(0, len(retiredPlayers)-1)
			retiree = retiredPlayers[idx]
		}
		(*retireeMap)[retiree.ID] = true
		(*coachMap)[retiree.ID] = true
		formerPlayerID = retiree.ID
		alma = retiree.College
		firstName = retiree.FirstName
		lastName = retiree.LastName
		posOne = retiree.Position
		age = retiree.Age + 1
	} else {
		fName := getName(firstNameList)
		lName := getName(lastNameList)
		caser := cases.Title(language.English)
		firstName = caser.String(strings.ToLower(fName))
		lastName = caser.String(strings.ToLower(lName))
		age = getCoachAge()
	}
	fullName := firstName + " " + lastName

	schoolQuality := team.AIQuality
	adminBehavior := team.AIBehavior
	goodHire := getGoodHire(schoolQuality, adminBehavior)
	starMin, starMax := getStarRange(schoolQuality, goodHire)
	pointMin, pointmax := getPointRange(schoolQuality, goodHire)
	odds1 := 0
	odds2 := 0
	odds3 := 0
	odds4 := 0
	odds5 := 0

	starList := make([]int, 5)
	for i := starMin; i <= starMax; i++ {
		starList = append(starList, i)
	}

	for _, star := range starList {
		if star == 1 {
			odds1 = 10
		} else if star == 2 {
			odds2 = 10
		} else if star == 3 {
			odds3 = 8
		} else if star == 4 {
			odds4 = 5
		} else if star == 5 {
			odds5 = 5
		}
	}

	offensiveSchemeList := []string{"Power Run", "Vertical", "West Coast", "I Option", "Run and Shoot", "Air Raid", "Pistol", "Spread Option", "Wing-T", "Double Wing", "Wishbone", "Flexbone"}
	offensiveScheme := util.PickFromStringList(offensiveSchemeList)
	defensiveSchemeList := []string{"Old School Front 7 Man", "2-Gap Zone", "4-man Front Spread Stopper Zone", "3-man Front Spread Stopper Zone", "Speed Man", "Multiple Man"}
	defensiveScheme := util.PickFromStringList(defensiveSchemeList)
	contractLength := util.GenerateIntFromRange(2, 5)
	startingPrestige := getStartingPrestige(goodHire)
	teamBuildingList := []string{"Recruiting", "Transfer", "Average"}
	teamBuildPref := util.PickFromStringList(teamBuildingList)
	careerPrefList := []string{"Average", "Prefers to Stay at Current Job", "Wants to coach Alma-Mater", "Wants a more competitive job", "Average"}
	careerPref := util.PickFromStringList(careerPrefList)
	promiseTendencyList := []string{"Average", "Under-Promise", "Over-Promise"}
	promiseTendency := util.PickFromStringList(promiseTendencyList)
	positionList := []string{"QB", "RB", "WR", "TE", "FB", "OT", "OG", "C", "DT", "DE", "ILB", "OLB", "FS", "SS", "CB", "P", "K", "ATH"}
	if posOne == "" {
		posOne = util.PickFromStringList(positionList)
	}
	for posTwo == "" || posTwo == posOne {
		posTwo = util.PickFromStringList(positionList)
	}
	for posThree == "" || posThree == posOne || posThree == posTwo {
		posThree = util.PickFromStringList(positionList)
	}
	if (careerPref == "Wants to coach at Alma Mater" && almaID == team.ID) || (schoolQuality == "Blue Blood" && careerPref == "Wants a more competitive job") {
		careerPref = "Prefers to Stay at Current Job"
	}
	if goodHire {
		fmt.Println("Good hire for " + team.TeamAbbreviation + "!")
	}
	formerPlayer := formerPlayerID > 0

	if formerPlayer {
		fmt.Println("Former SimNFL Player " + fullName + " is committing to coach for " + team.TeamAbbreviation + "!")
	}

	coach := structs.CollegeCoach{
		CoachName:              fullName,
		Age:                    age,
		TeamID:                 team.ID,
		Team:                   team.TeamAbbreviation,
		FormerPlayerID:         formerPlayerID,
		AlmaMaterID:            almaID,
		AlmaMater:              alma,
		Prestige:               startingPrestige,
		PointMin:               pointMin,
		PointMax:               pointmax,
		StarMin:                starMin,
		StarMax:                starMax,
		Odds1:                  odds1,
		Odds2:                  odds2,
		Odds3:                  odds3,
		Odds4:                  odds4,
		Odds5:                  odds5,
		OffensiveScheme:        offensiveScheme,
		DefensiveScheme:        defensiveScheme,
		TeambuildingPreference: teamBuildPref,
		CareerPreference:       careerPref,
		PromiseTendency:        promiseTendency,
		SchoolTenure:           0,
		CareerTenure:           0,
		ContractLength:         contractLength,
		YearsRemaining:         contractLength,
		IsRetired:              false,
		IsFormerPlayer:         formerPlayer,
		PortalReputation:       100,
		PositionOne:            posOne,
		PositionTwo:            posTwo,
		PositionThree:          posThree,
	}

	if startingPrestige > 1 {
		for i := 0; i < startingPrestige; i++ {
			selectStar := util.GenerateIntFromRange(starMin, starMax)
			coach.IncrementOdds(selectStar)
		}
	}

	return coach
}

func pickEthnicity() string {
	min := 0
	max := 10000
	num := util.GenerateIntFromRange(min, max)

	if num < 5000 {
		return "Caucasian"
	} else if num < 7800 {
		return "African"
	} else if num < 8900 {
		return "Hispanic"
	} else if num < 9975 {
		return "Asian"
	}
	return "NativeAmerican"
}

func pickWalkonState(state string) string {
	if state == "AL" {
		return util.PickFromStringList([]string{"AL", "LA", "MS", "TN", "GA", "FL"})
	} else if state == "AR" {
		return util.PickFromStringList([]string{"AR", "LA", "MO", "TN", "TX"})
	} else if state == "AZ" {
		return util.PickFromStringList([]string{"AZ", "NM", "CA"})
	} else if state == "CA" {
		return util.PickFromStringList([]string{"CA", "AZ", "HI"})
	} else if state == "CO" {
		return util.PickFromStringList([]string{"CO", "KS", "UT", "WY"})
	} else if state == "CT" {
		return util.PickFromStringList([]string{"CT", "NY", "NJ", "RI"})
	} else if state == "DC" {
		return util.PickFromStringList([]string{"DC", "MD", "VA"})
	} else if state == "FL" {
		return util.PickFromStringList([]string{"FL", "GA", "AL"})
	} else if state == "GA" {
		return util.PickFromStringList([]string{"GA", "FL", "SC", "AL"})
	} else if state == "HI" {
		return util.PickFromStringList([]string{"HI"})
	} else if state == "IA" {
		return util.PickFromStringList([]string{"IA", "MN", "WI", "NE"})
	} else if state == "ID" {
		return util.PickFromStringList([]string{"ID", "WA", "UT"})
	} else if state == "IN" {
		return util.PickFromStringList([]string{"IN", "IL", "OH", "MI", "AK"})
	} else if state == "IL" {
		return util.PickFromStringList([]string{"IL", "IN", "WI", "MI"})
	} else if state == "KS" {
		return util.PickFromStringList([]string{"KS", "MO", "NE"})
	} else if state == "KY" {
		return util.PickFromStringList([]string{"KY", "OH", "TN"})
	} else if state == "LA" {
		return util.PickFromStringList([]string{"LA", "TX", "MS"})
	} else if state == "MA" {
		return util.PickFromStringList([]string{"MA", "CT", "RI", "NH", "VT", "ME"})
	} else if state == "MD" {
		return util.PickFromStringList([]string{"DC", "MD", "VA", "DE"})
	} else if state == "MI" {
		return util.PickFromStringList([]string{"MI", "OH", "IN"})
	} else if state == "MN" {
		return util.PickFromStringList([]string{"MN", "WI", "IA"})
	} else if state == "MO" {
		return util.PickFromStringList([]string{"MO", "AR", "KS"})
	} else if state == "MS" {
		return util.PickFromStringList([]string{"MS", "LA", "AL"})
	} else if state == "MT" {
		return util.PickFromStringList([]string{"MT", "ID", "WY"})
	} else if state == "NC" {
		return util.PickFromStringList([]string{"NC", "SC", "VA"})
	} else if state == "ND" {
		return util.PickFromStringList([]string{"ND", "SD", "MN"})
	} else if state == "NE" {
		return util.PickFromStringList([]string{"NE", "KS", "SD", "IA"})
	} else if state == "NH" {
		return util.PickFromStringList([]string{"NH", "VT", "ME", "MA"})
	} else if state == "NJ" {
		return util.PickFromStringList([]string{"NJ", "DE", "NY", "CT", "PA"})
	} else if state == "NM" {
		return util.PickFromStringList([]string{"NM", "AZ", "TX"})
	} else if state == "NV" {
		return util.PickFromStringList([]string{"NV", "UT", "CA"})
	} else if state == "NY" {
		return util.PickFromStringList([]string{"NY", "NJ", "PA", "CT"})
	} else if state == "OH" {
		return util.PickFromStringList([]string{"OH", "KY", "MI", "PA"})
	} else if state == "OK" {
		return util.PickFromStringList([]string{"OK", "TX", "KS", "AR"})
	} else if state == "OR" {
		return util.PickFromStringList([]string{"OR", "WA", "CA"})
	} else if state == "PA" {
		return util.PickFromStringList([]string{"PA", "NJ", "DE", "OH", "WV"})
	} else if state == "RI" {
		return util.PickFromStringList([]string{"RI", "MA", "CT", "NY"})
	} else if state == "SC" {
		return util.PickFromStringList([]string{"SC", "NC", "GA"})
	} else if state == "SD" {
		return util.PickFromStringList([]string{"SD", "ND", "MN", "NE"})
	} else if state == "TN" {
		return util.PickFromStringList([]string{"TN", "KY", "GA", "AL", "AR"})
	} else if state == "TX" {
		return util.PickFromStringList([]string{"TX"})
	} else if state == "UT" {
		return util.PickFromStringList([]string{"UT", "CO", "ID", "AZ"})
	} else if state == "VA" {
		return util.PickFromStringList([]string{"VA", "WV", "DC", "MD"})
	} else if state == "WA" {
		return util.PickFromStringList([]string{"WA", "OR", "ID", "AK"})
	} else if state == "WI" {
		return util.PickFromStringList([]string{"WI", "MN", "IL", "MI"})
	} else if state == "WV" {
		return util.PickFromStringList([]string{"WV", "PA", "VA"})
	} else if state == "WY" {
		return util.PickFromStringList([]string{"WY", "CO", "UT", "MO", "ID"})
	}

	return "AK"
}

func getCrootLocations() map[string][]structs.CrootLocation {
	path := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\HS.json"

	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln("Error when opening file: ", err)
	}

	var payload map[string][]structs.CrootLocation
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during unmarshal: ", err)
	}

	return payload
}

func getAttributeBlob() map[string]map[string]map[string]map[string]interface{} {
	path := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\attributeBlob.json"

	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var payload map[string]map[string]map[string]map[string]interface{}
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during unmarshal: ", err)
	}

	return payload
}

func getNameList(ethnicity string, isFirstName bool) [][]string {
	path := filepath.Join(os.Getenv("ROOT"), "data")
	// path := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data"
	var fileName string
	if ethnicity == "Caucasian" {
		if isFirstName {
			fileName = "FNameW.csv"
		} else {
			fileName = "LNameW.csv"
		}
	} else if ethnicity == "African" {
		if isFirstName {
			fileName = "FNameB.csv"
		} else {
			fileName = "LNameB.csv"
		}
	} else if ethnicity == "Asian" {
		if isFirstName {
			fileName = "FNameA.csv"
		} else {
			fileName = "LNameA.csv"
		}
	} else if ethnicity == "NativeAmerican" {
		if isFirstName {
			fileName = "FNameN.csv"
		} else {
			fileName = "LNameN.csv"
		}
	} else {
		if isFirstName {
			fileName = "FNameH.csv"
		} else {
			fileName = "LNameH.csv"
		}
	}
	folderStr := "\\First Names\\"
	if !isFirstName {
		folderStr = "\\Last Names\\"
	}
	path = path + folderStr + fileName
	f, err := os.Open(path)
	if err != nil {
		log.Fatal("Unable to read input file "+path, err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+path, err)
	}

	return records
}

func getNameMaps() (map[string][][]string, map[string][][]string) {
	CaucasianFirstNameList := getNameList("Caucasian", true)
	CaucasianLastNameList := getNameList("Caucasian", false)
	AfricanFirstNameList := getNameList("African", true)
	AfricanLastNameList := getNameList("African", false)
	HispanicFirstNameList := getNameList("Hispanic", true)
	HispanicLastNameList := getNameList("Hispanic", false)
	NativeFirstNameList := getNameList("NativeAmerican", true)
	NativeLastNameList := getNameList("NativeAmerican", false)
	AsianFirstNameList := getNameList("Asian", true)
	AsianLastNameList := getNameList("Asian", false)

	firstNameMap := make(map[string][][]string)
	firstNameMap["Caucasian"] = CaucasianFirstNameList
	firstNameMap["African"] = AfricanFirstNameList
	firstNameMap["Hispanic"] = HispanicFirstNameList
	firstNameMap["NativeAmerican"] = NativeFirstNameList
	firstNameMap["Asian"] = AsianFirstNameList

	lastNameMap := make(map[string][][]string)
	lastNameMap["Caucasian"] = CaucasianLastNameList
	lastNameMap["African"] = AfricanLastNameList
	lastNameMap["Hispanic"] = HispanicLastNameList
	lastNameMap["NativeAmerican"] = NativeLastNameList
	lastNameMap["Asian"] = AsianLastNameList

	return (firstNameMap), (lastNameMap)
}

func getName(list [][]string) string {
	endOfListWeight, err := strconv.Atoi(list[len(list)-1][1])
	if err != nil {
		log.Fatalln("Could not convert number from string")
	}
	name := ""
	num := util.GenerateIntFromRange(1, endOfListWeight)
	for i := 1; i < len(list); i++ {
		weight, err := strconv.Atoi(list[i][1])
		if err != nil {
			log.Fatalln("Could not convert number from string in name generator")
		}
		if num < weight {
			name = list[i][0]
			break
		}
	}
	return name
}

func getArchetype(pos string) string {
	if pos == "QB" {
		return util.PickFromStringList([]string{"Balanced", "Pocket", "Scrambler", "Field General"})
	} else if pos == "RB" {
		return util.PickFromStringList([]string{"Balanced", "Power", "Speed", "Receiving"})
	} else if pos == "FB" {
		return util.PickFromStringList([]string{"Balanced", "Blocking", "Receiving", "Rushing"})
	} else if pos == "WR" {
		return util.PickFromStringList([]string{"Speed", "Possession", "Route Runner", "Red Zone Threat"})
	} else if pos == "TE" {
		return util.PickFromStringList([]string{"Blocking", "Receiving", "Vertical Threat"})
	} else if pos == "OT" || pos == "OG" {
		return util.PickFromStringList([]string{"Balanced", "Pass Blocking", "Run Blocking"})
	} else if pos == "C" {
		return util.PickFromStringList([]string{"Balanced", "Pass Blocking", "Run Blocking", "Line Captain"})
	} else if pos == "DT" {
		return util.PickFromStringList([]string{"Balanced", "Nose Tackle", "Pass Rusher"})
	} else if pos == "DE" {
		return util.PickFromStringList([]string{"Balanced", "Run Stopper", "Speed Rusher"})
	} else if pos == "ILB" {
		return util.PickFromStringList([]string{"Coverage", "Field General", "Run Stopper", "Speed"})
	} else if pos == "OLB" {
		return util.PickFromStringList([]string{"Coverage", "Pass Rush", "Run Stopper", "Speed"})
	} else if pos == "CB" {
		return util.PickFromStringList([]string{"Ball Hawk", "Man Coverage", "Zone Coverage"})
	} else if pos == "FS" || pos == "SS" {
		return util.PickFromStringList([]string{"Run Stopper", "Ball Hawk", "Man Coverage", "Zone Coverage"})
	} else if pos == "K" || pos == "P" {
		return util.PickFromStringList([]string{"Balanced", "Accuracy", "Power"})
	} else if pos == "ATH" {
		return util.PickFromStringList([]string{"Bandit", "Return Specialist", "Wingback", "Soccer Player", "Slotback", "Lineman", "Strongside", "Weakside", "Triple-Threat", "Field General"})
	}
	return ""

}

// Going to be honest, this should be a JSON file. This would be a huge blob of a map.
func getAttributeValue(pos string, arch string, star int, attr string, blob map[string]map[string]map[string]map[string]interface{}) int {
	starStr := strconv.Itoa(star)
	if pos == "QB" {
		if attr == "Catching" || attr == "Zone Coverage" || attr == "Man Coverage" || attr == "Tackle" {
			return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
		} else if attr == "Kick Accuracy" || attr == "Kick Power" || attr == "Punt Accuracy" || attr == "Punt Power" || attr == "Pass Block" || attr == "Run Block" || attr == "Pass Rush" || attr == "Run Defense" || attr == "Route Running" {
			return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
		}
		return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])

	} else if pos == "RB" {
		if attr == "Zone Coverage" || attr == "Man Coverage" || attr == "Tackle" || attr == "Throw Power" || attr == "Throw Accuracy" {
			return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
		} else if attr == "Kick Accuracy" || attr == "Kick Power" || attr == "Punt Accuracy" || attr == "Punt Power" || attr == "Pass Block" || attr == "Run Block" || attr == "Pass Rush" || attr == "Run Defense" {
			return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
		}
		return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
	} else if pos == "FB" {
		if attr == "Zone Coverage" || attr == "Man Coverage" || attr == "Tackle" || attr == "Throw Power" || attr == "Throw Accuracy" {
			return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
		} else if attr == "Kick Accuracy" || attr == "Kick Power" || attr == "Punt Accuracy" || attr == "Punt Power" || attr == "Pass Rush" || attr == "Run Defense" {
			return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
		}
		return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
	} else if pos == "WR" {
		if attr == "Zone Coverage" || attr == "Man Coverage" || attr == "Tackle" || attr == "Throw Power" || attr == "Throw Accuracy" {
			return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
		} else if attr == "Kick Accuracy" || attr == "Kick Power" || attr == "Punt Accuracy" || attr == "Punt Power" || attr == "Pass Block" || attr == "Pass Rush" || attr == "Run Defense" {
			return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
		}
		return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
	} else if pos == "TE" {
		if attr == "Zone Coverage" || attr == "Man Coverage" || attr == "Tackle" || attr == "Pass Rush" || attr == "Run Defense" || attr == "Throw Power" || attr == "Throw Accuracy" {
			return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
		} else if attr == "Kick Accuracy" || attr == "Kick Power" || attr == "Punt Accuracy" || attr == "Punt Power" {
			return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
		}
		return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
	} else if pos == "OT" || pos == "OG" {
		if attr == "Carrying" || attr == "Catching" || attr == "Zone Coverage" || attr == "Man Coverage" || attr == "Tackle" {
			return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
		} else if attr == "Kick Accuracy" || attr == "Kick Power" || attr == "Punt Accuracy" || attr == "Punt Power" || attr == "Pass Rush" || attr == "Run Defense" || attr == "Route Running" || attr == "Throw Power" || attr == "Throw Accuracy" {
			return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
		}
		return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
	} else if pos == "C" {
		if attr == "Carrying" || attr == "Catching" || attr == "Zone Coverage" || attr == "Man Coverage" || attr == "Tackle" {
			return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
		} else if attr == "Kick Accuracy" || attr == "Kick Power" || attr == "Punt Accuracy" || attr == "Punt Power" || attr == "Pass Rush" || attr == "Run Defense" || attr == "Route Running" || attr == "Throw Power" || attr == "Throw Accuracy" {
			return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
		}
		return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
	} else if pos == "DT" {
		if attr == "Carrying" || attr == "Catching" || attr == "Zone Coverage" || attr == "Man Coverage" || attr == "Throw Power" || attr == "Throw Accuracy" {
			return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
		} else if attr == "Kick Accuracy" || attr == "Kick Power" || attr == "Punt Accuracy" || attr == "Punt Power" || attr == "Pass Block" || attr == "Run Block" || attr == "Route Running" {
			return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
		}
		return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])

	} else if pos == "DE" {
		if attr == "Carrying" || attr == "Catching" || attr == "Zone Coverage" || attr == "Man Coverage" || attr == "Throw Power" || attr == "Throw Accuracy" {
			return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
		} else if attr == "Kick Accuracy" || attr == "Kick Power" || attr == "Punt Accuracy" || attr == "Punt Power" || attr == "Pass Block" || attr == "Run Block" || attr == "Route Running" {
			return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
		}
		return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
	} else if pos == "ILB" {
		if attr == "Carrying" || attr == "Catching" || attr == "Throw Power" || attr == "Throw Accuracy" {
			return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
		} else if attr == "Kick Accuracy" || attr == "Kick Power" || attr == "Punt Accuracy" || attr == "Punt Power" || attr == "Pass Block" || attr == "Run Block" || attr == "Route Running" {
			return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
		}
		return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
	} else if pos == "OLB" {
		if attr == "Carrying" || attr == "Catching" || attr == "Throw Power" || attr == "Throw Accuracy" {
			return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
		} else if attr == "Kick Accuracy" || attr == "Kick Power" || attr == "Punt Accuracy" || attr == "Punt Power" || attr == "Pass Block" || attr == "Run Block" || attr == "Route Running" {
			return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
		}
		return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
	} else if pos == "CB" {
		if attr == "Carrying" || attr == "Throw Power" || attr == "Throw Accuracy" || attr == "Route Running" {
			return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
		} else if attr == "Kick Accuracy" || attr == "Kick Power" || attr == "Punt Accuracy" || attr == "Punt Power" || attr == "Pass Block" || attr == "Run Block" || attr == "Pass Rush" || attr == "Run Defense" {
			return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
		}
		return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
	} else if pos == "FS" {
		if attr == "Carrying" || attr == "Throw Power" || attr == "Throw Accuracy" {
			return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
		} else if attr == "Kick Accuracy" || attr == "Kick Power" || attr == "Punt Accuracy" || attr == "Punt Power" || attr == "Pass Block" || attr == "Run Block" || attr == "Pass Rush" || attr == "Route Running" {
			return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
		}
		return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
	} else if pos == "SS" {
		if attr == "Carrying" || attr == "Throw Power" || attr == "Throw Accuracy" {
			return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
		} else if attr == "Kick Accuracy" || attr == "Kick Power" || attr == "Punt Accuracy" || attr == "Punt Power" || attr == "Pass Block" || attr == "Run Block" || attr == "Pass Rush" || attr == "Route Running" {
			return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
		}
		return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
	} else if pos == "K" {
		if attr == "Carrying" || attr == "Agility" || attr == "Catching" || attr == "Zone Coverage" || attr == "Man Coverage" || attr == "Punt Accuracy" || attr == "Punt Power" || attr == "Speed" || attr == "Throw Power" || attr == "Throw Accuracy" {
			return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
		} else if attr == "Pass Block" || attr == "Run Block" || attr == "Pass Rush" || attr == "Run Defense" || attr == "Route Running" || attr == "Tackle" || attr == "Strength" {
			return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
		}
		return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
	} else if pos == "P" {
		if attr == "Carrying" || attr == "Agility" || attr == "Catching" || attr == "Zone Coverage" || attr == "Man Coverage" || attr == "Kick Accuracy" || attr == "Kick Power" || attr == "Speed" || attr == "Throw Power" || attr == "Throw Accuracy" {
			return getValueFromInterfaceRange(starStr, blob["Default"]["Default"]["Default"])
		} else if attr == "Pass Block" || attr == "Run Block" || attr == "Pass Rush" || attr == "Run Defense" || attr == "Route Running" || attr == "Tackle" || attr == "Strength" {
			return getValueFromInterfaceRange(starStr, blob["Under"]["Under"]["Under"])
		}
		return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
	} else if pos == "ATH" {
		return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
	}
	return util.GenerateIntFromRange(5, 15)
}

func getCityAndHighSchool(schools []structs.CrootLocation) (string, string) {
	randInt := util.GenerateIntFromRange(0, len(schools)-1)

	return schools[randInt].City, schools[randInt].HighSchool
}

func getValueFromInterfaceRange(star string, starMap map[string]interface{}) int {
	u := starMap[star]
	// if ok {
	// 	fmt.Println("(Was able to get value)")
	// }

	minMax, ok := u.([]interface{})
	if !ok {
		fmt.Printf("This is not an int: " + star)
	}

	min, ok := minMax[0].(float64)
	if !ok {
		fmt.Printf("This is not an int: " + star)
	}

	max, ok := minMax[1].(float64)
	if !ok {
		fmt.Printf("This is not an int: " + star)
	}
	return util.GenerateIntFromRange(int(min), int(max))
}

func getCustomCrootStarRating() int {
	weightedRoll := util.GenerateIntFromRange(1, 100)

	if weightedRoll < 45 {
		return 3
	} else if weightedRoll < 80 {
		return 4
	}
	return 5
}

func getWalkonStarRating() int {
	weightedRoll := util.GenerateIntFromRange(0, 10000)

	if weightedRoll < 9001 {
		return 1
	} else if weightedRoll < 9601 {
		return 2
	} else if weightedRoll < 9901 {
		return 3
	} else if weightedRoll < 9976 {
		return 4
	}
	return 5
}

func getLatestRecord(db *gorm.DB) uint {
	var lastPlayerRecord structs.Player

	err := db.Last(&lastPlayerRecord).Error
	if err != nil {
		log.Fatalln("Could not grab last player record from players table...")
	}

	return lastPlayerRecord.ID + 1
}

func pickAlmaMater(teams []structs.RecruitingTeamProfile) structs.RecruitingTeamProfile {
	start := 0
	end := len(teams) - 1
	idx := util.GenerateIntFromRange(start, end)
	return teams[idx]
}

func getCoachAge() int {
	num := util.GenerateIntFromRange(1, 100)

	if num < 10 {
		return util.GenerateIntFromRange(32, 36)
	} else if num < 25 {
		return util.GenerateIntFromRange(37, 39)
	} else if num < 55 {
		return util.GenerateIntFromRange(40, 49)
	} else if num < 80 {
		return util.GenerateIntFromRange(50, 59)
	} else if num < 95 {
		return util.GenerateIntFromRange(60, 65)
	} else {
		return util.GenerateIntFromRange(66, 70)
	}
}

func getGoodHire(schoolQuality, adminBehavior string) bool {
	diceRoll := util.GenerateIntFromRange(1, 20)
	mod := 0
	if schoolQuality == "P6" || schoolQuality == "Playoff Buster" {
		mod += 1
	} else if schoolQuality == "Blue Blood" {
		mod += 3
	}
	if adminBehavior == "Aggressive" {
		mod += 3
	} else if adminBehavior == "Conservative" {
		mod -= 3
	}

	sum := diceRoll + mod
	goodHire := sum > 12
	return goodHire
}

func getStarRange(schoolQuality string, goodHire bool) (int, int) {

	if schoolQuality == "Blue Blood" {
		if goodHire {
			return 3, 5
		} else {
			return 3, 4
		}
	} else if schoolQuality == "Playoff Buster" {
		if goodHire {
			return 2, 4
		} else {
			return 2, 3
		}
	} else if schoolQuality == "Normal" {
		if goodHire {
			return 2, 4
		} else {
			return 2, 3
		}
	} else {
		if goodHire {
			return 1, 3
		} else {
			return 1, 2
		}
	}
}

func getPointRange(schoolQuality string, goodHire bool) (int, int) {
	min := 0
	max := 15
	if schoolQuality == "Blue Blood" {
		if goodHire {
			min = util.GenerateIntFromRange(7, 8)
			max = util.GenerateIntFromRange(12, 16)
		} else {
			min = util.GenerateIntFromRange(6, 7)
			max = util.GenerateIntFromRange(10, 13)
		}
	} else if schoolQuality == "Playoff Buster" {
		if goodHire {
			min = util.GenerateIntFromRange(5, 7)
			max = util.GenerateIntFromRange(10, 15)
		} else {
			min = util.GenerateIntFromRange(4, 6)
			max = util.GenerateIntFromRange(10, 12)
		}
	} else if schoolQuality == "Normal" {
		if goodHire {
			min = util.GenerateIntFromRange(5, 8)
			max = util.GenerateIntFromRange(10, 14)
		} else {
			min = util.GenerateIntFromRange(4, 6)
			max = util.GenerateIntFromRange(8, 12)
		}
	} else {
		if goodHire {
			min = util.GenerateIntFromRange(3, 6)
			max = util.GenerateIntFromRange(8, 12)
		} else {
			min = 4
			max = util.GenerateIntFromRange(6, 8)
		}
	}
	return min, max
}

func getStartingPrestige(goodHire bool) int {
	if goodHire {
		return util.GenerateIntFromRange(3, 7)
	}
	return util.GenerateIntFromRange(1, 5)
}

func getRelativeType() int {
	roll := util.GenerateIntFromRange(1, 1000)
	// Brother of existing player
	if roll < 600 {
		return 2
	}
	// Cousin of existing player
	if roll < 800 {
		return 3
	}
	// Half brother of existing player
	if roll < 850 {
		return 4
	}
	// Twin
	if roll < 900 {
		return 5
	}
	// Best friend of another recruit
	if roll < 925 {
		return 8
	}
	// Best friend of a college player
	if roll < 950 {
		return 8
	}
	// Coach's Son
	if roll < 985 {
		return 6
	}
	// Coach's Nephew
	return 7
}

func getCoachLastName(fullName string) string {
	nameSplit := strings.Split(fullName, " ")
	return nameSplit[1]
}
