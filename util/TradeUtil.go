package util

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

func GetRoundAbbreviation(str string) string {
	if str == "1" {
		return "1st Round"
	} else if str == "2" {
		return "2nd Round"
	} else if str == "3" {
		return "3rd Round"
	} else if str == "4" {
		return "4th Round"
	} else if str == "5" {
		return "5th Round"
	} else if str == "6" {
		return "6th Round"
	}
	return "7th Round"
}

func GetTagObject() map[string]map[string]float64 {
	return map[string]map[string]float64{
		"QB": {
			"Franchise":  18.94,
			"Transition": 16.83,
			"Playtime":   12.77,
			"Basic":      10.98,
		},
		"RB": {
			"Franchise":  10.03,
			"Transition": 8.59,
			"Playtime":   6.34,
			"Basic":      5.74,
		},
		"FB": {
			"Franchise":  4.27,
			"Transition": 3.1,
			"Playtime":   1.8,
			"Basic":      1.6,
		},
		"TE": {
			"Franchise":  7.85,
			"Transition": 6.83,
			"Playtime":   5.45,
			"Basic":      5.01,
		},
		"WR": {
			"Franchise":  12.89,
			"Transition": 11.15,
			"Playtime":   8.88,
			"Basic":      8.25,
		},
		"OT": {
			"Franchise":  15.66,
			"Transition": 14.11,
			"Playtime":   11.91,
			"Basic":      11.09,
		},
		"OG": {
			"Franchise":  16.72,
			"Transition": 12.45,
			"Playtime":   7.98,
			"Basic":      7.44,
		},
		"C": {
			"Franchise":  9.42,
			"Transition": 8.24,
			"Playtime":   6.21,
			"Basic":      5.33,
		},
		"DE": {
			"Franchise":  15.43,
			"Transition": 13.59,
			"Playtime":   10.42,
			"Basic":      9.71,
		},
		"DT": {
			"Franchise":  11.47,
			"Transition": 10.15,
			"Playtime":   8.16,
			"Basic":      7.57,
		},
		"OLB": {
			"Franchise":  13.66,
			"Transition": 12.34,
			"Playtime":   9.63,
			"Basic":      8.64,
		},
		"ILB": {
			"Franchise":  9.84,
			"Transition": 8.87,
			"Playtime":   7.02,
			"Basic":      6.35,
		},
		"CB": {
			"Franchise":  11.64,
			"Transition": 10.12,
			"Playtime":   7.9,
			"Basic":      7.44,
		},
		"FS": {
			"Franchise":  8.24,
			"Transition": 7.06,
			"Playtime":   4.64,
			"Basic":      4.07,
		},
		"SS": {
			"Franchise":  10.9,
			"Transition": 8.76,
			"Playtime":   5.47,
			"Basic":      4.76,
		},
		"P": {
			"Franchise":  2.92,
			"Transition": 2.07,
			"Playtime":   1.2,
			"Basic":      1.07,
		},
		"K": {
			"Franchise":  3.48,
			"Transition": 2.69,
			"Playtime":   1.56,
			"Basic":      1.38,
		},
	}
}

func GetTagData() map[string]map[string]float64 {
	path := filepath.Join(os.Getenv("ROOT"), "data", "tagData.json")
	content := ReadJson(path)

	var payload map[string]map[string]float64

	err := json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatalln("Error during unmarshal: ", err)
	}

	return payload
}
