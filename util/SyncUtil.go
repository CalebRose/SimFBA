package util

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/CalebRose/SimFBA/structs"
)

func GetWinsAndLossesForCollegeGames(games []structs.CollegeGame, TeamID int, ConferenceCheck bool) (int, int) {
	wins := 0
	losses := 0

	for _, game := range games {
		if !game.GameComplete {
			continue
		}
		if ConferenceCheck && !game.IsConference {
			continue
		}
		if (game.HomeTeamID == TeamID && game.HomeTeamWin) ||
			(game.AwayTeamID == TeamID && game.AwayTeamWin) {
			wins += 1
		} else {
			losses += 1
		}
	}

	return wins, losses
}

func GetConferenceChampionshipWeight(games []structs.CollegeGame, TeamID int) float64 {
	var weight float64 = 0

	for _, game := range games {
		if !game.IsConference {
			continue
		}
		if (game.HomeTeamID == TeamID && game.HomeTeamScore > game.AwayTeamScore) ||
			(game.AwayTeamID == TeamID && game.AwayTeamScore > game.HomeTeamScore) {
			weight = 1
		} else {
			weight = 0.5
		}
	}

	return weight
}

func GetPostSeasonWeight(games []structs.CollegeGame, TeamID int) float64 {
	for _, game := range games {
		if !game.IsPlayoffGame || !game.IsBowlGame {
			continue
		}
		return 1
	}
	return 0
}

func FilterOutRecruitingProfile(profiles []structs.RecruitPlayerProfile, ID int) []structs.RecruitPlayerProfile {
	var rp []structs.RecruitPlayerProfile

	for _, profile := range profiles {
		if profile.ProfileID != ID {
			rp = append(rp, profile)
		}
	}

	return rp
}

func GetTeamPointsMap() map[string]float64 {
	return map[string]float64{
		"USAF": 0,
		"AKRN": 0,
		"BAMA": 0,
		"APST": 0,
		"ZONA": 0,
		"AZST": 0,
		"ARK":  0,
		"ARST": 0,
		"ARMY": 0,
		"AUB":  0,
		"BALL": 0,
		"BAYL": 0,
		"BOIS": 0,
		"BC":   0,
		"BGSU": 0,
		"BUFF": 0,
		"BYU":  0,
		"CAL":  0,
		"CMU":  0,
		"CHAR": 0,
		"CINC": 0,
		"CLEM": 0,
		"CCU":  0,
		"COLO": 0,
		"CSU":  0,
		"DUKE": 0,
		"ECU":  0,
		"EMU":  0,
		"FIU":  0,
		"FLA":  0,
		"FAU":  0,
		"FSU":  0,
		"FRES": 0,
		"UGA":  0,
		"GASO": 0,
		"GAST": 0,
		"GT":   0,
		"HAWI": 0,
		"UHOU": 0,
		"ILLI": 0,
		"IND":  0,
		"IOWA": 0,
		"IAST": 0,
		"KANS": 0,
		"KSST": 0,
		"KENT": 0,
		"UKEN": 0,
		"LU":   0,
		"ULL":  0,
		"ULM":  0,
		"LT":   0,
		"LOU":  0,
		"LSU":  0,
		"MRSH": 0,
		"UMD":  0,
		"MEMP": 0,
		"MIAF": 0,
		"MIAO": 0,
		"MICH": 0,
		"MIST": 0,
		"MTSU": 0,
		"MINN": 0,
		"MSST": 0,
		"MIZZ": 0,
		"NAVY": 0,
		"NCST": 0,
		"NEB":  0,
		"NEV":  0,
		"UNM":  0,
		"NMSU": 0,
		"UNC":  0,
		"UNT":  0,
		"NIU":  0,
		"NW":   0,
		"ND":   0,
		"OHIO": 0,
		"OHST": 0,
		"OKLA": 0,
		"OKST": 0,
		"ODU":  0,
		"MISS": 0,
		"OREG": 0,
		"ORST": 0,
		"PNST": 0,
		"PITT": 0,
		"PURD": 0,
		"RICE": 0,
		"RUTG": 0,
		"SDSU": 0,
		"SJSU": 0,
		"SMU":  0,
		"USA":  0,
		"SOCA": 0,
		"USF":  0,
		"USM":  0,
		"STAN": 0,
		"CUSE": 0,
		"TCU":  0,
		"TEMP": 0,
		"TENN": 0,
		"TEX":  0,
		"TAMU": 0,
		"TXST": 0,
		"TTU":  0,
		"TLDO": 0,
		"TROY": 0,
		"TLNE": 0,
		"TULS": 0,
		"UAB":  0,
		"UCF":  0,
		"UCLA": 0,
		"CONN": 0,
		"MASS": 0,
		"UNLV": 0,
		"USC":  0,
		"UTEP": 0,
		"UTSA": 0,
		"UTAH": 0,
		"UTST": 0,
		"VAND": 0,
		"UVA":  0,
		"VT":   0,
		"WAKE": 0,
		"WASH": 0,
		"WAST": 0,
		"WVU":  0,
		"WKU":  0,
		"WMU":  0,
		"WISC": 0,
		"WYOM": 0,
	}
}

func ConvertStringToInt(num string) int {
	val, err := strconv.Atoi(num)
	if err != nil {
		log.Fatalln("Could not convert string to int")
	}

	return val
}

func ConvertStringToFloat(num string) float64 {
	floatNum, error := strconv.ParseFloat(num, 64)
	if error != nil {
		fmt.Println("Could not convert string to float 64, resetting as 0.")
		return 0
	}
	return floatNum
}

// Reads specific CSV values as Boolean. If the value is "0" or "FALSE" or "False", it will be read as false. Anything else is considered True.
func ConvertStringToBool(str string) bool {
	if str == "0" || str == "FALSE" || str == "False" {
		return false
	}
	return true
}

func IsAITeamContendingForCroot(profiles []structs.RecruitPlayerProfile) float64 {
	if len(profiles) == 0 {
		return 0
	}
	var leadingVal float64 = 0
	for _, profile := range profiles {
		if profile.TotalPoints != 0 && profile.TotalPoints > float64(leadingVal) {
			leadingVal = profile.TotalPoints
		}
	}

	return leadingVal
}

func ReadJson(path string) []byte {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	return content
}

func ReadLocalPath(subpath string) string {
	path, err := filepath.Abs(subpath)
	if err != nil {
		log.Fatal(err)
	}

	return path
}

func ReadCSV(path string) [][]string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal("Unable to read input file "+path, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	rows, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+path, err)
	}

	return rows
}

func GetStateRegionMatcher() map[string]map[string]string {
	path := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "CalebRose", "SimFBA", "data", "regionMatcher.json")
	content := ReadJson(path)

	var payload map[string]map[string]string

	err := json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatalln("Error during unmarshal: ", err)
	}

	return payload
}

func GetStateMatcher() map[string][]string {
	path := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "CalebRose", "SimFBA", "data", "stateMatcher.json")
	content := ReadJson(path)

	var payload map[string][]string

	err := json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatalln("Error during unmarshal: ", err)
	}

	return payload
}

func IsCrootCloseToHome(crootState string, crootCity string, teamState string, abbr string, stateMatcher map[string][]string, regionMatcher map[string]map[string]string) bool {
	if crootState == teamState {
		return true
	}
	state := crootState
	if crootState == "TX" || crootState == "CA" || crootState == "FL" {
		state = regionMatcher[crootState][crootCity]
	}

	for _, s := range stateMatcher[state] {
		if s == abbr {
			return true
		}
	}
	return false
}
