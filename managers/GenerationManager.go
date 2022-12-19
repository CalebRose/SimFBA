package managers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
)

func GenerateWalkOns() {
	fmt.Println(time.Now().UnixNano())
	rand.Seed(time.Now().UnixNano())
	db := dbprovider.GetInstance().GetDB()
	AllTeams := GetRecruitingProfileForRecruitSync()
	count := 0
	attributeBlob := getAttributeBlob()
	highSchoolBlob := getCrootLocations()

	firstNameMap, lastNameMap := getNameMaps()

	var lastPlayerRecord structs.Player

	err := db.Last(&lastPlayerRecord).Error
	if err != nil {
		log.Fatalln("Could not grab last player record from players table...")
	}

	newID := lastPlayerRecord.ID + 1

	for _, team := range AllTeams {
		id := strconv.Itoa(int(team.ID))
		signedRecruits := GetSignedRecruitsByTeamProfileID(id)
		if len(signedRecruits) == team.RecruitClassSize {
			continue
		}
		positionList := []string{}

		// Get Team Needs
		teamNeeds := GetRecruitingNeeds(id)
		limit := team.RecruitClassSize - len(signedRecruits)

		for _, recruit := range signedRecruits {
			if teamNeeds[recruit.Position] > 0 {
				teamNeeds[recruit.Position] -= 1
			}
		}

		// Get All Needed Positions into a list
		for k, v := range teamNeeds {
			i := v
			for i > 0 {
				positionList = append(positionList, k)
				i--
			}
		}

		// Randomize List
		rand.Shuffle(len(positionList), func(i, j int) {
			positionList[i], positionList[j] = positionList[j], positionList[i]
		})

		// Recruit Generation
		for _, pos := range positionList {
			if count >= limit {
				break
			}

			year := 1
			ethnicity := pickEthnicity()

			state := pickState(team.State)

			recruit := createRecruit(ethnicity, pos, year, firstNameMap[ethnicity], lastNameMap[ethnicity], newID, attributeBlob, state, highSchoolBlob[state])

			recruit.AssignWalkon(team.TeamAbbreviation, int(team.ID), newID)

			recruitPlayerRecord := structs.RecruitPlayerProfile{
				ProfileID:        int(team.ID),
				RecruitID:        int(newID),
				IsSigned:         true,
				IsLocked:         true,
				TeamAbbreviation: team.TeamAbbreviation,
				SeasonID:         2,
				TotalPoints:      1,
			}

			playerRecord := structs.Player{
				RecruitID:       int(newID),
				CollegePlayerID: int(newID),
				NFLPlayerID:     int(newID),
			}
			playerRecord.AssignID(newID)
			count++

			db.Create(&playerRecord)
			db.Create(&recruit)
			db.Create(&recruitPlayerRecord)
			newID++
			team.IncreaseCommitCount()
			db.Save(&team)
		}
		count = 0
		fmt.Println("Finished walkon generation for " + team.TeamAbbreviation)
	}
}

func createRecruit(ethnicity string, position string, year int, firstNameList [][]string, lastNameList [][]string, id uint, blob map[string]map[string]map[string]map[string]interface{}, state string, hsBlob []structs.CrootLocation) structs.Recruit {
	fName := getName(firstNameList)
	lName := getName(lastNameList)
	firstName := strings.Title(strings.ToLower(fName))
	lastName := strings.Title(strings.ToLower(lName))
	age := 18
	city, highSchool := getCityAndHighSchool(hsBlob)

	archetype := getArchetype(position)
	stars := getStarRating()
	if stars == 5 {
		fmt.Println("WE'VE GOT A FIVE STAR")
	}
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
	injury := util.GenerateIntFromRange(10, 100)
	stamina := util.GenerateIntFromRange(10, 100)
	discipline := util.GenerateIntFromRange(10, 100)
	progression := util.GenerateIntFromRange(1, 100)

	freeAgency := util.GetFreeAgencyBias()
	personality := util.GetPersonality()
	recruitingBias := util.GetRecruitingBias()
	workEthic := util.GetWorkEthic()
	academicBias := util.GetAcademicBias()
	potentialGrade := util.GetWeightedPotentialGrade(progression)

	basePlayer := structs.BasePlayer{
		FirstName:      firstName,
		LastName:       lastName,
		Position:       position,
		Archetype:      archetype,
		Age:            age,
		Stars:          0,
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
	}

	basePlayer.GetOverall()

	return structs.Recruit{
		BasePlayer: basePlayer,
		City:       city,
		HighSchool: highSchool,
		State:      state,
		IsSigned:   true,
	}
}

func pickEthnicity() string {
	min := 0
	max := 10000
	num := util.GenerateIntFromRange(min, max)

	if num < 6000 {
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

func pickState(state string) string {
	if state == "AL" {
		return util.PickFromStringList([]string{"AL", "LA", "MS", "TN", "GA", "FL"})
	} else if state == "AR" {
		return util.PickFromStringList([]string{"AR", "LA", "MO", "TN", "TX"})
	} else if state == "AZ" {
		return util.PickFromStringList([]string{"AZ", "NM", "CA"})
	} else if state == "CA" {
		return util.PickFromStringList([]string{"CA", "AZ"})
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
		return util.PickFromStringList([]string{"WA", "OR", "ID"})
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

	content, err := ioutil.ReadFile(path)
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
	path := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimNBA\\data"
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
	path = path + "\\" + fileName
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
		return util.PickFromStringList([]string{"Balanced", "Ball Hawk", "Man Coverage", "Zone Coverage"})
	} else if pos == "FS" || pos == "SS" {
		return util.PickFromStringList([]string{"Run Stopper", "Ball Hawk", "Man Coverage", "Zone Coverage"})
	} else if pos == "K" || pos == "P" {
		return util.PickFromStringList([]string{"Balanced", "Accuracy", "Power"})
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
	}
	return util.GenerateIntFromRange(5, 15)
}

func getCityAndHighSchool(schools []structs.CrootLocation) (string, string) {
	randInt := util.GenerateIntFromRange(0, len(schools)-1)

	return schools[randInt].City, schools[randInt].HighSchool
}

func getValueFromInterfaceRange(star string, starMap map[string]interface{}) int {
	u, ok := starMap[star]
	if ok {
		fmt.Println("Was able to get value)")
	}

	minMax, ok := u.([]interface{})
	if !ok {
		fmt.Printf("This is not an int")
	}

	min, ok := minMax[0].(float64)
	if !ok {
		fmt.Printf("This is not an int")
	}

	max, ok := minMax[1].(float64)
	if !ok {
		fmt.Printf("This is not an int")
	}
	return util.GenerateIntFromRange(int(min), int(max))
}

func getStarRating() int {
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
