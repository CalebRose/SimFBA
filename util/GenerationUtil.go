package util

import (
	"fmt"
	"math/rand"
	"strings"
)

func GetPositionList() []string {
	return []string{
		"QB", "RB", "FB", "TE", "WR", "OT", "OG", "C",
		"DT", "DE", "ILB", "OLB", "CB", "FS", "SS", "P", "K",
		"ATH",
	}
}

func GenerateIntFromRange(min int, max int) int {
	diff := max - min + 1
	if diff < 0 {
		diff = 1
	}
	return rand.Intn(diff) + min
}

func GenerateFloatFromRange(min float64, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func GenerateNormalizedIntFromRange(min int, max int) int {
	mean := float64(min+max) / 2.0
	stdDev := float64(max-min) / 6.0 // This approximates the 3-sigma rule

	for {
		// Generate a number using normal distribution around the mean
		num := rand.NormFloat64()*stdDev + mean
		// Round to nearest integer and convert to int type
		intNum := int(num + 0.5) // Adding 0.5 before truncating simulates rounding
		// Check if the generated number is within bounds
		if intNum >= min && intNum <= max {
			return intNum
		}
		// If not within bounds, loop again
	}
}

func GenerateNormalizedIntFromMeanStdev(mean, stdDev float64) float64 {
	num := rand.NormFloat64()*stdDev + mean
	// Round to nearest integer and convert to int type
	intNum := int(num + 0.5) // Adding 0.5 before truncating simulates rounding
	return float64(intNum)
}

func RegressValue(val, min, max int) int {
	newVal := val - GenerateNormalizedIntFromRange(min, max)

	if newVal < 1 {
		return 1
	}
	return newVal
}

func PickFromStringList(list []string) string {
	if len(list) == 0 {
		return ""
	}
	return list[rand.Intn(len(list))]
}

func GetProgressionRating() int {
	weight := GenerateIntFromRange(1, 10000)

	val := weight / 100

	return val
}

func GeneratePotential() int {
	num := GenerateIntFromRange(1, 100)

	if num < 10 {
		return GenerateIntFromRange(1, 20)
	} else if num < 20 {
		return GenerateIntFromRange(21, 40)
	} else if num < 80 {
		return GenerateIntFromRange(41, 55)
	} else if num < 85 {
		return GenerateIntFromRange(56, 65)
	} else if num < 90 {
		return GenerateIntFromRange(66, 75)
	} else if num < 95 {
		return GenerateIntFromRange(76, 85)
	} else {
		return GenerateIntFromRange(86, 99)
	}
}

func GenerateNFLPotential(pot int) int {
	floor := pot - 20
	ceil := pot + 20
	if floor < 0 {
		diff := 0 - floor
		floor = 0
		ceil += diff
	}
	if ceil > 100 {
		diff := ceil - 100
		ceil = 100
		floor += diff
	}
	return GenerateIntFromRange(floor, ceil)
}

func GetWeightedPotentialGrade(rating int) string {
	weightedRating := GenerateIntFromRange(rating-15, rating+15)
	if weightedRating > 100 {
		weightedRating = 99
	} else if weightedRating < 0 {
		weightedRating = 0
	}

	if weightedRating > 88 {
		return "A+"
	}
	if weightedRating > 80 {
		return "A"
	}
	if weightedRating > 74 {
		return "A-"
	}
	if weightedRating > 68 {
		return "B+"
	}
	if weightedRating > 62 {
		return "B"
	}
	if weightedRating > 56 {
		return "B-"
	}
	if weightedRating > 50 {
		return "C+"
	}
	if weightedRating > 44 {
		return "C"
	}
	if weightedRating > 38 {
		return "C-"
	}
	if weightedRating > 32 {
		return "D+"
	}
	if weightedRating > 26 {
		return "D"
	}
	if weightedRating > 20 {
		return "D-"
	}
	return "F"
}

func GetPrimeAge(pos, arch string) int {
	venerable := false
	vDiceRoll := GenerateIntFromRange(1, 10000)
	chance := getVenerableChance(pos)
	if vDiceRoll < chance {
		venerable = true
	}

	mean, stddev := getPositionMean(pos, venerable)

	age := GenerateNormalizedIntFromMeanStdev(mean, stddev)
	return int(age)
}

func getPositionMean(pos string, venerable bool) (float64, float64) {
	meanMap := getPositionMeanMap()
	return meanMap[venerable][pos][0], meanMap[venerable][pos][1]
}

func getPositionMeanMap() map[bool]map[string][]float64 {
	return map[bool]map[string][]float64{
		true: {
			"QB":  []float64{39, 2},
			"RB":  []float64{32, 1},
			"FB":  []float64{35, 1},
			"WR":  []float64{35, 1},
			"TE":  []float64{35, 1},
			"OT":  []float64{35, 1},
			"OG":  []float64{35, 1},
			"C":   []float64{35, 1},
			"DE":  []float64{36, 1},
			"DT":  []float64{35, 1},
			"ILB": []float64{35, 1},
			"OLB": []float64{35, 1},
			"CB":  []float64{35, 1},
			"FS":  []float64{35, 1},
			"SS":  []float64{35, 1},
			"K":   []float64{40, 1},
			"P":   []float64{38, 1},
			"ATH": []float64{36, 1},
		},
		false: {
			"QB":  []float64{32, 2},
			"RB":  []float64{26, 0.67},
			"FB":  []float64{26, 0.67},
			"WR":  []float64{29, 1.33},
			"TE":  []float64{28, 1.33},
			"OT":  []float64{30, 1.4},
			"OG":  []float64{30, 0.67},
			"C":   []float64{30, 1.33},
			"DE":  []float64{30, 0.67},
			"DT":  []float64{30, 1.33},
			"ILB": []float64{29, 0.67},
			"OLB": []float64{29, 0.67},
			"CB":  []float64{29, 0.67},
			"FS":  []float64{29, 1.33},
			"SS":  []float64{29, 1.33},
			"K":   []float64{34, 0.73},
			"P":   []float64{31, 2},
			"ATH": []float64{30, 1.33},
		},
	}
}

func getVenerableChance(pos string) int {
	if pos == "QB" || pos == "K" || pos == "P" {
		return 20
	}
	if pos == "RB" {
		return 5
	}
	return 10
}

func GetPersonality() string {
	chance := GenerateIntFromRange(1, 3)
	if chance < 3 {
		return "Average"
	}
	list := []string{"Reserved",
		"Eccentric",
		"Motivational",
		"Disloyal",
		"Cooperative",
		"Irrational",
		"Focused",
		"Book Worm",
		"Motivation",
		"Abrasive",
		"Absent Minded",
		"Uncooperative",
		"Introvert",
		"Disruptive",
		"Outgoing",
		"Tough",
		"Paranoid",
		"Chill",
		"Stoic",
		"Dramatic",
		"Extroverted",
		"Selfish",
		"Impatient",
		"Reliable",
		"Frail",
		"Relaxed",
		"Average",
		"Flamboyant",
		"Perfectionist",
		"Popular",
		"Jokester",
		"Narcissist",
		"Laid Back"}

	return PickFromStringList(list)
}

func GetAcademicBias() string {
	chance := GenerateIntFromRange(1, 3)
	if chance < 3 {
		return "Average"
	}
	list := []string{"Takes AP classes",
		"Sits at the front of the class",
		"Seeks out tutoring",
		"Tutor",
		"Wants to finish degree",
		"Teacher's Pet",
		"Sits at the back of the class",
		"Values academics",
		"Studious",
		"Frequent visits to the principal",
		"Class Clown",
		"More likely to get academic probation",
		"Has other priorities",
		"Distracted",
		"Loves Learning",
		"Studies hard",
		"Less likely to get academic probation",
		"Never Studies",
		"Average",
		"Naturally book smart",
		"Borderline failing",
		"Skips classes often",
		"Didn't come here to play school",
		"Spends more time on Tiktok than focusing in class"}

	return PickFromStringList(list)
}

func GetRecruitingBias() string {
	chance := GenerateIntFromRange(1, 3)
	if chance < 3 {
		return "Average"
	}
	list := []string{"Prefers to play in a different state",
		"Prefers to play for an up-and-coming team",
		"Open-Minded",
		"Prefers to play for a team where he can start immediately",
		"Prefers to be close to home",
		"Prefers to play for a national championship contender",
		"Prefers to play for a specific coach",
		"Average",
		"Legacy",
		"Prefers to play for a team with a rich history",
	}

	return PickFromStringList(list)
}

func GetWorkEthic() string {
	chance := GenerateIntFromRange(1, 3)
	if chance < 3 {
		return "Average"
	}
	list := []string{"Persistant",
		"Lazy",
		"Footwork king",
		"Hard-working",
		"Complacent",
		"Skips Leg Day",
		"Working-Class mentality",
		"Film Room Genius",
		"Focuses on Max Weight",
		"Track Athlete",
		"Average",
		"Center of Attention",
		"Gym Rat",
		"Focuses on Max Reps",
		"Loud",
		"Quiet",
		"Streams too much",
		"Trolls on Discord"}
	return PickFromStringList(list)
}

func GetFreeAgencyBias() string {
	chance := GenerateIntFromRange(1, 3)
	if chance < 3 {
		return "Average"
	}
	list := []string{
		"I'm the starter",
		"Market-driven",
		"Wants extensions",
		"Drafted team discount",
		"Highest bidder",
		"Championship seeking",
		"Loyal",
		"Average",
		"Hometown hero",
		"Money motivated",
		"Hates Tags",
		"Adversarial",
		"Considering retirement"}

	return PickFromStringList(list)
}

func GetNewArchetypeMap() map[string]map[string]map[string]string {
	return map[string]map[string]map[string]string{
		"QB": {
			"Balanced": {
				"RB": "Balanced",
				"WR": "Possession",
				"TE": "Receiving",
			},
			"Scrambler": {
				"RB": "Speed",
				"WR": "Speed",
				"TE": "Vertical Threat",
			},
			"Pocket": {
				"RB": "Power",
				"WR": "Possession",
				"TE": "Receiving",
			},
			"Field General": {
				"RB": "Balanced",
				"WR": "Possession",
				"TE": "Receiving",
			},
		},
		"RB": {
			"Balanced": {
				"FB": "Balanced",
				"WR": "Possession",
				"TE": "Blocking",
			},
			"Speed": {
				"FB": "Rushing",
				"WR": "Speed",
				"TE": "Vertical Threat",
			},
			"Power": {
				"FB": "Balanced",
				"WR": "Possession",
				"TE": "Blocking",
			},
			"Receiving": {
				"FB": "Receiving",
				"WR": "Route Runner",
				"TE": "Vertical Threat",
			},
		},
		"FB": {
			"Balanced": {
				"RB": "Power",
				"WR": "Possession",
				"TE": "Receiving",
				"OT": "Balanced",
				"OG": "Balanced",
			},
			"Blocking": {
				"RB": "Power",
				"WR": "Possession",
				"TE": "Blocking",
				"OT": "Balanced",
				"OG": "Balanced",
			},
			"Rushing": {
				"RB": "Balanced",
				"WR": "Speed",
				"TE": "Receiving",
				"OT": "Balanced",
				"OG": "Balanced",
			},
			"Receiving": {
				"RB": "Receiving",
				"WR": "Route Runner",
				"TE": "Receiving",
				"OT": "Balanced",
				"OG": "Balanced",
			},
		},
		"TE": {
			"Vertical Threat": {
				"RB": "Receiving",
				"WR": "Possession",
				"FB": "Receiving",
				"OT": "Balanced",
				"OG": "Balanced",
			},
			"Blocking": {
				"RB": "Power",
				"WR": "Red Zone Threat",
				"FB": "Blocking",
				"OT": "Balanced",
				"OG": "Balanced",
			},
			"Receiving": {
				"RB": "Receiving",
				"WR": "Red Zone Threat",
				"FB": "Receiving",
				"OT": "Balanced",
				"OG": "Balanced",
			},
		},
		"WR": {
			"Route Runner": {
				"FB": "Receiving",
				"RB": "Receiving",
				"TE": "Receiving",
			},
			"Speed": {
				"FB": "Receiving",
				"RB": "Speed",
				"TE": "Vertical Threat",
			},
			"Red Zone Threat": {
				"FB": "Receiving",
				"RB": "Balanced",
				"TE": "Receiving",
			},
			"Possession": {
				"FB": "Receiving",
				"RB": "Receiving",
				"TE": "Receiving",
			},
		},
		"OG": {
			"Balanced": {
				"OT": "Balanced",
				"C":  "Balanced",
				"TE": "Blocking",
			},
			"Pass Blocking": {
				"OT": "Pass Blocking",
				"C":  "Pass Blocking",
				"TE": "Blocking",
			},
			"Run Blocking": {
				"OT": "Run Blocking",
				"C":  "Run Blocking",
				"TE": "Blocking",
			},
		},
		"OT": {
			"Balanced": {
				"OG": "Balanced",
				"C":  "Balanced",
				"TE": "Blocking",
			},
			"Pass Blocking": {
				"OG": "Pass Blocking",
				"C":  "Pass Blocking",
				"TE": "Blocking",
			},
			"Run Blocking": {
				"OG": "Run Blocking",
				"C":  "Run Blocking",
				"TE": "Blocking",
			},
		},
		"C": {
			"Line Captain": {
				"OT": "Balanced",
				"OG": "Balanced",
				"TE": "Blocking",
			},
			"Balanced": {
				"OT": "Balanced",
				"OG": "Balanced",
				"TE": "Blocking",
			},
			"Pass Blocking": {
				"OT": "Pass Blocking",
				"OG": "Pass Blocking",
				"TE": "Blocking",
			},
			"Run Blocking": {
				"OT": "Run Blocking",
				"OG": "Run Blocking",
				"TE": "Blocking",
			},
		},
		"DE": {
			"Balanced": {
				"DT":  "Balanced",
				"OLB": "Run Stopper",
				"ILB": "Run Stopper",
			},
			"Speed Rush": {
				"DT":  "Pass Rusher",
				"OLB": "Pass Rush",
				"ILB": "Speed",
			},
			"Run Stopper": {
				"DT":  "Balanced",
				"OLB": "Run Stopper",
				"ILB": "Run Stopper",
			},
		},
		"DT": {
			"Balanced": {
				"DE":  "Balanced",
				"OLB": "Run Stopper",
			},
			"Pass Rusher": {
				"DE":  "Speed Rush",
				"OLB": "Pass Rush",
			},
			"Nose Tackle": {
				"DE":  "Run Stopper",
				"OLB": "Run Stopper",
			},
		},
		"OLB": {
			"Coverage": {
				"DE":  "Balanced",
				"ILB": "Coverage",
				"SS":  "Man Coverage",
			},
			"Pass Rush": {
				"DE":  "Speed Rush",
				"ILB": "Speed",
				"SS":  "Man Coverage",
			},
			"Run Stopper": {
				"DE":  "Run Stopper",
				"ILB": "Run Stopper",
				"SS":  "Run Stopper",
			},
			"Speed": {
				"DE":  "Speed Rush",
				"ILB": "Speed",
				"SS":  "Zone Coverage",
			},
		},
		"ILB": {
			"Coverage": {
				"DE":  "Balanced",
				"OLB": "Coverage",
				"SS":  "Man Coverage",
			},
			"Field General": {
				"DE":  "Balanced",
				"OLB": "Coverage",
				"SS":  "Run Stopper",
			},
			"Run Stopper": {
				"DE":  "Run Stopper",
				"OLB": "Run Stopper",
				"SS":  "Run Stopper",
			},
			"Speed": {
				"DE":  "Speed Rush",
				"OLB": "Speed",
				"SS":  "Run Stopper",
			},
		},
		"SS": {
			"Man Coverage": {
				"FS":  "Man Coverage",
				"OLB": "Coverage",
				"ILB": "Coverage",
				"CB":  "Man Coverage",
				"DE":  "Run Stopper",
			},
			"Zone Coverage": {
				"FS":  "Zone Coverage",
				"OLB": "Coverage",
				"ILB": "Coverage",
				"CB":  "Zone Coverage",
				"DE":  "Run Stopper",
			},
			"Ball Hawk": {
				"FS":  "Ball Hawk",
				"OLB": "Speed",
				"ILB": "Speed",
				"CB":  "Ball Hawk",
				"DE":  "Run Stopper",
			},
			"Run Stopper": {
				"FS":  "Run Stopper",
				"OLB": "Run Stopper",
				"ILB": "Run Stopper",
				"CB":  "Man Coverage",
				"DE":  "Run Stopper",
			},
		},
		"FS": {
			"Man Coverage": {
				"SS":  "Man Coverage",
				"OLB": "Coverage",
				"ILB": "Coverage",
				"CB":  "Man Coverage",
				"DE":  "Run Stopper",
			},
			"Zone Coverage": {
				"SS":  "Zone Coverage",
				"OLB": "Coverage",
				"ILB": "Coverage",
				"CB":  "Zone Coverage",
				"DE":  "Run Stopper",
			},
			"Ball Hawk": {
				"SS":  "Ball Hawk",
				"OLB": "Speed",
				"ILB": "Speed",
				"CB":  "Ball Hawk",
				"DE":  "Run Stopper",
			},
			"Run Stopper": {
				"SS":  "Run Stopper",
				"OLB": "Run Stopper",
				"ILB": "Run Stopper",
				"CB":  "Zone Coverage",
				"DE":  "Run Stopper",
			},
		},
		"CB": {
			"Man Coverage": {
				"SS":  "Man Coverage",
				"OLB": "Coverage",
				"ILB": "Coverage",
				"FS":  "Man Coverage",
			},
			"Zone Coverage": {
				"SS":  "Zone Coverage",
				"OLB": "Coverage",
				"ILB": "Coverage",
				"FS":  "Zone Coverage",
			},
			"Ball Hawk": {
				"SS":  "Ball Hawk",
				"OLB": "Speed",
				"ILB": "Speed",
				"FS":  "Ball Hawk",
			},
		},
		"K": {
			"Power": {
				"P":  "Power",
				"QB": "Balanced",
			},
			"Accuracy": {
				"P":  "Accuracy",
				"QB": "Balanced",
			},
		},
		"P": {
			"Power": {
				"K":  "Power",
				"QB": "Balanced",
			},
			"Accuracy": {
				"K":  "Accuracy",
				"QB": "Balanced",
			},
		},
		"ATH": {
			"Field General": {
				"ILB": "Field General",
				"QB":  "Field General",
			},
			"Triple-Threat": {
				"K":  "Accuracy",
				"QB": PickFromStringList([]string{"Balanced", "Scrambler"}),
				"P":  "Accuracy",
				"WR": "Speed",
				"RB": "Speed",
			},
			"Wingback": {
				"RB": PickFromStringList([]string{"Balanced", "Speed"}),
				"WR": "Speed",
				"TE": "Receiving",
			},
			"Slotback": {
				"RB": "Power",
				"WR": "Red Zone Threat",
				"TE": "Vertical Threat",
			},
			"Lineman": {
				"OT": "Balanced",
				"OG": "Balanced",
				"C":  "Balanced",
				"DE": "Balanced",
				"DT": "Balanced",
			},
			"Strongside": {
				"DE":  PickFromStringList([]string{"Balanced", "Run Stopper"}),
				"OLB": PickFromStringList([]string{"Coverage", "Run Stopper"}),
			},
			"Weakside": {
				"ILB": "Coverage",
				"OLB": "Coverage",
				"FS":  PickFromStringList([]string{"Zone Coverage", "Man Coverage"}),
				"CB":  PickFromStringList([]string{"Zone Coverage", "Man Coverage"}),
				"SS":  PickFromStringList([]string{"Zone Coverage", "Man Coverage"}),
			},
			"Bandit": {
				"SS":  "Run Stopper",
				"CB":  "Ball Hawk",
				"FS":  "Run Stopper",
				"OLB": "Run Stopper",
			},
			"Return Specialist": {
				"K":  "Accuracy",
				"QB": "Balanced",
				"WR": "Route Runner",
				"RB": "Balanced",
			},
			"Soccer Player": {
				"K":  "Accuracy",
				"P":  "Accuracy",
				"WR": "Speed",
				"RB": "Speed",
			},
		},
	}
}

func GetStarRating() int {
	roll := GenerateIntFromRange(1, 10000)
	if roll < 67 {
		return 5
	}
	if roll < 865 {
		return 4
	}
	if roll < 3683 {
		return 3
	}
	if roll < 7367 {
		return 2
	}
	return 1
}

func PickState() string {
	diceRoll := GenerateIntFromRange(1, 4748)
	if diceRoll < 672 {
		return "FL"
	}
	if diceRoll < 1332 {
		return "TX"
	}
	if diceRoll < 1780 {
		return "CA"
	}
	if diceRoll < 2200 {
		return "GA"
	}
	if diceRoll < 2442 {
		return "OH"
	}
	if diceRoll < 2632 {
		return "LA"
	}
	if diceRoll < 2814 {
		return "AL"
	}
	if diceRoll < 2980 {
		return "NC"
	}
	if diceRoll < 3124 {
		return "MI"
	}
	if diceRoll < 3256 {
		return "IL"
	}
	if diceRoll < 3372 {
		return "VA"
	}
	if diceRoll < 3484 {
		return "PA"
	}
	if diceRoll < 3584 {
		return "MD"
	}
	if diceRoll < 3680 {
		return "NJ"
	}
	if diceRoll < 3768 {
		return "MS"
	}
	if diceRoll < 3844 {
		return "IN"
	}
	if diceRoll < 3916 {
		return "SC"
	}
	if diceRoll < 3988 {
		return "TN"
	}
	if diceRoll < 4058 {
		return "WA"
	}
	if diceRoll < 4124 {
		return "OK"
	}
	if diceRoll < 4190 {
		return "UT"
	}
	if diceRoll < 4244 {
		return "MO"
	}
	if diceRoll < 4288 {
		return "AZ"
	}
	if diceRoll < 4328 {
		return "KY"
	}
	if diceRoll < 4366 {
		return "NY"
	}
	if diceRoll < 4402 {
		return "HI"
	}
	if diceRoll < 4438 {
		return "MN"
	}
	if diceRoll < 4468 {
		return "AR"
	}
	if diceRoll < 4498 {
		return "CO"
	}
	if diceRoll < 4528 {
		return "KS"
	}
	if diceRoll < 4558 {
		return "WI"
	}
	if diceRoll < 4586 {
		return "OR"
	}
	if diceRoll < 4612 {
		return "NV"
	}
	if diceRoll < 4636 {
		return "IA"
	}
	if diceRoll < 4658 {
		return "MA"
	}
	if diceRoll < 4676 {
		return "CT"
	}
	if diceRoll < 4688 {
		return "NE"
	}
	if diceRoll < 4700 {
		return "NM"
	}
	if diceRoll < 4712 {
		return "DC"
	}
	if diceRoll < 4722 {
		return "WV"
	}
	if diceRoll < 4728 {
		return "ID"
	}
	if diceRoll < 4732 {
		return "DE"
	}
	if diceRoll < 4736 {
		return "WY"
	}
	if diceRoll < 4738 {
		return "ND"
	}
	if diceRoll < 4740 {
		return "RI"
	}
	if diceRoll < 4742 {
		return "SD"
	}
	if diceRoll < 4743 {
		return "AK"
	}
	if diceRoll < 4744 {
		return "ME"
	}
	if diceRoll < 4745 {
		return "MT"
	}
	if diceRoll < 4746 {
		return "NH"
	}
	return "VT"
}

// getStateAbbreviation returns the two-letter state abbreviation for a given state name.
func GetStateAbbreviation(state string) (string, error) {
	// Map of state names to their two-letter abbreviations
	stateAbbreviations := map[string]string{
		"Alabama":        "AL",
		"Alaska":         "AK",
		"Arizona":        "AZ",
		"Arkansas":       "AR",
		"California":     "CA",
		"Colorado":       "CO",
		"Connecticut":    "CT",
		"Delaware":       "DE",
		"Florida":        "FL",
		"Georgia":        "GA",
		"Hawai'i":        "HI",
		"Idaho":          "ID",
		"Illinois":       "IL",
		"Indiana":        "IN",
		"Iowa":           "IA",
		"Kansas":         "KS",
		"Kentucky":       "KY",
		"Louisiana":      "LA",
		"Maine":          "ME",
		"Maryland":       "MD",
		"Massachusetts":  "MA",
		"Michigan":       "MI",
		"Minnesota":      "MN",
		"Mississippi":    "MS",
		"Missouri":       "MO",
		"Montana":        "MT",
		"Nebraska":       "NE",
		"Nevada":         "NV",
		"New Hampshire":  "NH",
		"New Jersey":     "NJ",
		"New Mexico":     "NM",
		"New York":       "NY",
		"North Carolina": "NC",
		"North Dakota":   "ND",
		"Ohio":           "OH",
		"Oklahoma":       "OK",
		"Oregon":         "OR",
		"Pennsylvania":   "PA",
		"Rhode Island":   "RI",
		"South Carolina": "SC",
		"South Dakota":   "SD",
		"Tennessee":      "TN",
		"Texas":          "TX",
		"Utah":           "UT",
		"Vermont":        "VT",
		"Virginia":       "VA",
		"Washington":     "WA",
		"West Virginia":  "WV",
		"Wisconsin":      "WI",
		"Wyoming":        "WY",
	}

	// Normalize the input by trimming spaces and capitalizing the first letter of each word
	normalizedState := strings.Title(strings.ToLower(strings.TrimSpace(state)))

	// Check if the state exists in the map
	if abbreviation, ok := stateAbbreviations[normalizedState]; ok {
		return abbreviation, nil
	}

	return "", fmt.Errorf("state not found: %s", state)
}

func PickPosition() string {
	roll := GenerateIntFromRange(1, 100000)
	if roll < 6681 {
		return "QB"
	}
	if roll < 14780 {
		return "RB"
	}
	if roll < 16645 {
		return "FB"
	}
	if roll < 26952 {
		return "WR"
	}
	if roll < 31123 {
		return "TE"
	}
	if roll < 39979 {
		return "OT"
	}
	if roll < 45622 {
		return "OG"
	}
	if roll < 48569 {
		return "C"
	}
	if roll < 58346 {
		return "DE"
	}
	if roll < 64356 {
		return "DT"
	}
	if roll < 69219 {
		return "ILB"
	}
	if roll < 75683 {
		return "OLB"
	}
	if roll < 84097 {
		return "CB"
	}
	if roll < 88180 {
		return "FS"
	}
	if roll < 91171 {
		return "SS"
	}
	if roll < 93035 {
		return "K"
	}
	if roll < 94900 {
		return "P"
	}
	return "ATH"
}

func PickAffinity(stars int, af1 string, pickingAf2 bool) string {
	if af1 == "" && pickingAf2 {
		return ""
	}
	returnAnAffinity := GenerateIntFromRange(1, 2)
	if returnAnAffinity == 2 {
		return ""
	}
	roll := GenerateIntFromRange(1, 100)
	if roll < 40 && af1 != "Close to Home" {
		return "Close to Home"
	}

	list := []string{}

	if af1 != "Academics" {
		list = append(list, "Academics")
	}

	if af1 != "Service" {
		list = append(list, "Service")
	}

	if af1 != "Religion" {
		list = append(list, "Religion")
	}

	if af1 != "Small School" {
		list = append(list, "Small School")
	}

	coinFlip := GenerateIntFromRange(1, 2)
	if coinFlip == 1 && af1 != "Small Town" && af1 != "Big City" {
		list = append(list, "Small Town")
	} else if coinFlip == 2 && af1 != "Small Town" && af1 != "Big City" {
		list = append(list, "Big City")
	}

	if stars >= 2 {
		list = append(list, "Media Spotlight", "Rising Stars")
	}
	if stars >= 3 {
		list = append(list, "Frontrunner", "Large Crowds")
	}

	return PickFromStringList(list)
}
