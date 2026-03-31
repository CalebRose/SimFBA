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
			"Franchise":  21.76,
			"Transition": 17.47,
			"Playtime":   10.93,
			"Basic":      9.72,
		},
		"RB": {
			"Franchise":  7.42,
			"Transition": 6.37,
			"Playtime":   5.04,
			"Basic":      4.62,
		},
		"FB": {
			"Franchise":  2.79,
			"Transition": 2.15,
			"Playtime":   1.36,
			"Basic":      1.23,
		},
		"TE": {
			"Franchise":  7.48,
			"Transition": 6.5,
			"Playtime":   5.13,
			"Basic":      4.74,
		},
		"WR": {
			"Franchise":  11.5,
			"Transition": 10.29,
			"Playtime":   8.42,
			"Basic":      7.87,
		},
		"OT": {
			"Franchise":  15.03,
			"Transition": 13.57,
			"Playtime":   11.32,
			"Basic":      10.67,
		},
		"OG": {
			"Franchise":  17.10,
			"Transition": 13.92,
			"Playtime":   9.77,
			"Basic":      9.27,
		},
		"C": {
			"Franchise":  17.10,
			"Transition": 13.92,
			"Playtime":   9.77,
			"Basic":      9.27,
		},
		"DE": {
			"Franchise":  14.25,
			"Transition": 13.0,
			"Playtime":   9.63,
			"Basic":      8.39,
		},
		"DT": {
			"Franchise":  12.4,
			"Transition": 11.27,
			"Playtime":   9.2,
			"Basic":      8.62,
		},
		"OLB": {
			"Franchise":  12.33,
			"Transition": 10.86,
			"Playtime":   8.87,
			"Basic":      8.31,
		},
		"ILB": {
			"Franchise":  12.33,
			"Transition": 10.86,
			"Playtime":   8.87,
			"Basic":      8.31,
		},
		"CB": {
			"Franchise":  12.17,
			"Transition": 10.60,
			"Playtime":   8.42,
			"Basic":      7.81,
		},
		"FS": {
			"Franchise":  7.59,
			"Transition": 6.49,
			"Playtime":   4.33,
			"Basic":      3.86,
		},
		"SS": {
			"Franchise":  14.9,
			"Transition": 11.18,
			"Playtime":   6.01,
			"Basic":      5.21,
		},
		"P": {
			"Franchise":  2.18,
			"Transition": 1.77,
			"Playtime":   1.2,
			"Basic":      1.1,
		},
		"K": {
			"Franchise":  2.52,
			"Transition": 1.79,
			"Playtime":   1.06,
			"Basic":      0.99,
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
