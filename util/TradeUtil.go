package util

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

func GetRoundAbbreviation(str string) string {
	switch str {
	case "1":
		return "1st Round"
	case "2":
		return "2nd Round"
	case "3":
		return "3rd Round"
	case "4":
		return "4th Round"
	case "5":
		return "5th Round"
	case "6":
		return "6th Round"
	}
	return "7th Round"
}

func GetTagObject() map[string]map[string]float64 {
	return map[string]map[string]float64{
		"QB": {
			"Franchise":  21.26,
			"Transition": 16.83,
			"Playtime":   12.77,
			"Basic":      10.98,
		},
		"RB": {
			"Franchise":  6.92,
			"Transition": 6.92,
			"Playtime":   6.34,
			"Basic":      5.74,
		},
		"FB": {
			"Franchise":  2.29,
			"Transition": 2.29,
			"Playtime":   1.8,
			"Basic":      1.6,
		},
		"TE": {
			"Franchise":  6.98,
			"Transition": 6.83,
			"Playtime":   5.45,
			"Basic":      5.01,
		},
		"WR": {
			"Franchise":  11,
			"Transition": 11,
			"Playtime":   8.88,
			"Basic":      8.25,
		},
		"OT": {
			"Franchise":  15.03,
			"Transition": 14.11,
			"Playtime":   11.91,
			"Basic":      11.09,
		},
		"OG": {
			"Franchise":  16.60,
			"Transition": 12.45,
			"Playtime":   7.98,
			"Basic":      7.44,
		},
		"C": {
			"Franchise":  16.60,
			"Transition": 8.24,
			"Playtime":   6.21,
			"Basic":      5.33,
		},
		"DE": {
			"Franchise":  13.75,
			"Transition": 13.59,
			"Playtime":   10.42,
			"Basic":      9.71,
		},
		"DT": {
			"Franchise":  11.90,
			"Transition": 10.15,
			"Playtime":   8.16,
			"Basic":      7.57,
		},
		"OLB": {
			"Franchise":  11.83,
			"Transition": 11.83,
			"Playtime":   9.63,
			"Basic":      8.64,
		},
		"ILB": {
			"Franchise":  11.83,
			"Transition": 11.83,
			"Playtime":   7.02,
			"Basic":      6.35,
		},
		"CB": {
			"Franchise":  11.67,
			"Transition": 10.12,
			"Playtime":   7.9,
			"Basic":      7.44,
		},
		"FS": {
			"Franchise":  7.09,
			"Transition": 7.06,
			"Playtime":   4.64,
			"Basic":      4.07,
		},
		"SS": {
			"Franchise":  14.4,
			"Transition": 8.76,
			"Playtime":   5.47,
			"Basic":      4.76,
		},
		"P": {
			"Franchise":  1.68,
			"Transition": 1.68,
			"Playtime":   1.2,
			"Basic":      1.07,
		},
		"K": {
			"Franchise":  2.02,
			"Transition": 2.02,
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
