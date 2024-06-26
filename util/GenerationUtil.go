package util

import (
	"math/rand"
)

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
	min := 25
	max := 27
	mod := 0
	diceRoll := GenerateIntFromRange(1, 20)
	if pos == "QB" {
		if diceRoll == 100 {
			mod += GenerateIntFromRange(1, 5)
		}
		min = 26
		max = 32
	}
	if pos == "RB" {
		if diceRoll == 100 {
			mod += GenerateIntFromRange(1, 3)
		}
		min = 24
		max = 27
	}
	if pos == "FB" {
		if diceRoll == 100 {
			mod += GenerateIntFromRange(1, 3)
		}
		min = 24
		max = 29
	}
	if pos == "WR" {
		if diceRoll == 100 {
			mod += GenerateIntFromRange(1, 3)
		}
		min = 25
		max = 29
	}
	if pos == "TE" {
		if diceRoll == 100 {
			mod += GenerateIntFromRange(1, 3)
		}
		min = 26
		max = 29
	}
	if pos == "OG" {
		min = 27
		max = 30
	}
	if pos == "OT" {
		min = 27
		max = 30
	}
	if pos == "C" {
		min = 27
		max = 30
	}
	if pos == "DE" {
		min = 27
		max = 30
	}
	if pos == "DT" {
		min = 27
		max = 30
	}
	if pos == "ILB" {
		if diceRoll == 100 {
			mod += GenerateIntFromRange(1, 3)
		}
		min = 26
		max = 29
	}
	if pos == "OLB" {
		if diceRoll == 100 {
			mod += GenerateIntFromRange(1, 3)
		}
		min = 26
		max = 29
	}
	if pos == "CB" {
		if diceRoll == 100 {
			mod += GenerateIntFromRange(1, 3)
		}
		min = 25
		max = 29
	}
	if pos == "FS" {
		if diceRoll == 100 {
			mod += GenerateIntFromRange(1, 3)
		}
		min = 26
		max = 28
	}
	if pos == "SS" {
		if diceRoll == 100 {
			mod += GenerateIntFromRange(1, 3)
		}
		min = 26
		max = 28
	}
	if pos == "K" {
		min = 25
		max = 30
	}
	if pos == "P" {
		min = 25
		max = 30
	}
	if pos == "ATH" {
		if arch == "Field General" {
			if diceRoll == 100 {
				mod += GenerateIntFromRange(1, 3)
			}
			min = 25
			max = 30
		}
		if arch == "Triple-Threat" {
			if diceRoll == 100 {
				mod += GenerateIntFromRange(1, 3)
			}
			min = 24
			max = 27
		}
		if arch == "Wingback" {
			if diceRoll == 100 {
				mod += GenerateIntFromRange(1, 3)
			}
			min = 24
			max = 28
		}
		if arch == "Slotback" {
			if diceRoll == 100 {
				mod += GenerateIntFromRange(1, 3)
			}
			min = 24
			max = 28
		}
		if arch == "Lineman" {
			min = 27
			max = 30
		}
		if arch == "Strongside" {
			if diceRoll == 100 {
				mod += GenerateIntFromRange(1, 3)
			}
			min = 26
			max = 29
		}
		if arch == "Weakside" {
			if diceRoll == 100 {
				mod += GenerateIntFromRange(1, 3)
			}
			min = 26
			max = 29
		}
		if arch == "Bandit" {
			if diceRoll == 100 {
				mod += GenerateIntFromRange(1, 3)
			}
			min = 26
			max = 29
		}
		if arch == "Soccer Player" {
			min = 26
			max = 29
		}
	}
	return GenerateNormalizedIntFromRange(min, max) + mod
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
		"Narcissist"}

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
		"Didn't come here to play school"}

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
