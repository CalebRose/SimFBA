package util

import "strconv"

func GetNFLOverallGrade(value int) string {
	if value > 63 {
		return "A+"
	} else if value > 61 {
		return "A"
	} else if value > 59 {
		return "A-"
	} else if value > 57 {
		return "B+"
	} else if value > 55 {
		return "B"
	} else if value > 53 {
		return "B-"
	} else if value > 51 {
		return "C+"
	} else if value > 49 {
		return "C"
	} else if value > 47 {
		return "C-"
	} else if value > 45 {
		return "D+"
	} else if value > 43 {
		return "D"
	} else if value > 41 {
		return "D-"
	}
	return "F"
}

func GetOverallGrade(value int, year int) string {
	if year < 3 {
		if value > 44 {
			return "A"
		} else if value > 34 {
			return "B"
		} else if value > 24 {
			return "C"
		} else if value > 14 {
			return "D"
		}
	} else {
		if value > 47 {
			return "A"
		} else if value > 44 {
			return "A-"
		} else if value > 40 {
			return "B+"
		} else if value > 37 {
			return "B"
		} else if value > 34 {
			return "B-"
		} else if value > 30 {
			return "C+"
		} else if value > 27 {
			return "C"
		} else if value > 24 {
			return "C-"
		} else if value > 20 {
			return "D+"
		} else if value > 17 {
			return "D"
		} else if value > 14 {
			return "D-"
		}
	}

	return "F"
}

func GetLetterGrade(Attribute int, mean float32, stddev float32, year int) string {
	if mean == 0 || stddev == 0 {
		return GetOverallGrade(Attribute, year)
	}

	val := float32(Attribute)
	dev := stddev * 2
	if year < 3 {
		if val > mean+dev {
			return "A"
		}
		dev = stddev
		if val > mean+dev {
			return "B"
		}
		if val > mean {
			return "C"
		}
		dev = stddev * -1
		if val > mean+dev {
			return "D"
		}
	} else {
		dev = stddev * 2.5
		if val > mean+dev {
			return "A+"
		}
		dev = stddev * 2
		if val > mean+dev {
			return "A"
		}
		dev = stddev * 1.75
		if val > mean+dev {
			return "A-"
		}
		dev = stddev * 1.5
		if val > mean+dev {
			return "B+"
		}
		dev = stddev * 1
		if val > mean+dev {
			return "B"
		}
		dev = stddev * 0.75
		if val > mean+dev {
			return "B-"
		}
		dev = stddev * 0.5
		if val > mean+dev {
			return "C+"
		}
		if val > mean {
			return "C"
		}
		dev = stddev * -0.5
		if val > mean+dev {
			return "C-"
		}
		dev = stddev * -0.75
		if val > mean+dev {
			return "D+"
		}
		dev = stddev * -1
		if val > mean+dev {
			return "D"
		}
		dev = stddev * -1.5
		if val > mean+dev {
			return "D-"
		}
	}
	return "F"
}

func GetYearAndRedshirtStatus(year int, redshirt bool) (string, string) {
	status := "Eligible"
	switch year {
	case 1:
		return "Fr", status
	case 2:
		return "So", status
	case 3:
		return "Jr", status
	case 4:
		return "Sr", status
	case 5:
		return "5th", status
	default:
		return "Grad", status
	}
}

func GetNFLYear(year uint) string {
	if year < 2 {
		return "R"
	}
	return strconv.Itoa(int(year))
}

func GetWinningVerb(score1, score2 int) string {
	if score2 == 0 {
		return " shuts out "
	}

	diff := score1 - score2
	if diff < 7 {
		return PickFromStringList([]string{" upsets ", " edges ", " sneaks an upset "})
	} else if diff < 14 {
		return PickFromStringList([]string{" upsets ", " outshines ", " knocks out ", " stuns "})
	}
	return PickFromStringList([]string{" upsets ", " silences ", " strikes down ", " knocks out "})
}
