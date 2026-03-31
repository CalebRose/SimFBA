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
			"Transition": 16.97,
			"Playtime":   10.43,
			"Basic":      9.22,
		},
		"RB": {
			"Franchise":  6.92,
			"Transition": 5.87,
			"Playtime":   4.54,
			"Basic":      4.12,
		},
		"FB": {
			"Franchise":  2.29,
			"Transition": 1.65,
			"Playtime":   0.86,
			"Basic":      0.73,
		},
		"TE": {
			"Franchise":  6.98,
			"Transition": 6.0,
			"Playtime":   4.63,
			"Basic":      4.24,
		},
		"WR": {
			"Franchise":  11.0,
			"Transition": 9.79,
			"Playtime":   7.92,
			"Basic":      7.37,
		},
		"OT": {
			"Franchise":  14.53,
			"Transition": 13.07,
			"Playtime":   10.82,
			"Basic":      10.17,
		},
		"OG": {
			"Franchise":  16.60,
			"Transition": 13.42,
			"Playtime":   9.27,
			"Basic":      8.77,
		},
		"C": {
			"Franchise":  16.60,
			"Transition": 13.42,
			"Playtime":   9.27,
			"Basic":      8.77,
		},
		"DE": {
			"Franchise":  13.75,
			"Transition": 12.5,
			"Playtime":   9.13,
			"Basic":      7.89,
		},
		"DT": {
			"Franchise":  11.9,
			"Transition": 10.77,
			"Playtime":   8.7,
			"Basic":      8.12,
		},
		"OLB": {
			"Franchise":  11.83,
			"Transition": 10.36,
			"Playtime":   8.37,
			"Basic":      7.81,
		},
		"ILB": {
			"Franchise":  11.83,
			"Transition": 10.36,
			"Playtime":   8.37,
			"Basic":      7.81,
		},
		"CB": {
			"Franchise":  11.67,
			"Transition": 10.10,
			"Playtime":   7.92,
			"Basic":      7.31,
		},
		"FS": {
			"Franchise":  7.09,
			"Transition": 5.99,
			"Playtime":   3.83,
			"Basic":      3.36,
		},
		"SS": {
			"Franchise":  14.4,
			"Transition": 10.68,
			"Playtime":   5.51,
			"Basic":      4.71,
		},
		"P": {
			"Franchise":  1.68,
			"Transition": 1.27,
			"Playtime":   0.7,
			"Basic":      0.6,
		},
		"K": {
			"Franchise":  2.02,
			"Transition": 1.29,
			"Playtime":   0.56,
			"Basic":      0.49,
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
