package util

func GetNumericalSortValueByLetterGrade(grade string) int {
	if grade == "A+" {
		return 1
	} else if grade == "A" {
		return 2
	} else if grade == "A-" {
		return 3
	} else if grade == "B+" {
		return 4
	} else if grade == "B" {
		return 5
	} else if grade == "B-" {
		return 6
	} else if grade == "C+" {
		return 7
	} else if grade == "C" {
		return 8
	} else if grade == "C-" {
		return 9
	} else if grade == "D+" {
		return 10
	} else if grade == "D" {
		return 11
	} else if grade == "D-" {
		return 12
	}
	return 13
}

func GetDrafteeSalary(pick, year, round uint, isSalary bool) float64 {
	var valueMap [38][4]float64
	if isSalary {
		valueMap = getSalaryMap()
	} else {
		valueMap = getBonusMap()
	}
	if pick >= 1 && pick <= 32 && year >= 1 && year <= 4 {
		return valueMap[pick-1][year-1]
	}
	return valueMap[31+round][year-1]
}

func getSalaryMap() [38][4]float64 {
	return [38][4]float64{
		{3.25, 3.25, 3.25, 3.25},
		{2.75, 2.75, 2.75, 2.75},
		{2.75, 2.75, 2.75, 2.75},
		{2.75, 2.75, 2.75, 2.75},
		{2.75, 2.75, 2.75, 2.75},
		{2.25, 2.25, 2.25, 2.25},
		{2.25, 2.25, 2.25, 2.25},
		{2.25, 2.25, 2.25, 2.25},
		{2.25, 2.25, 2.25, 2.25},
		{2.25, 2.25, 2.25, 2.25},
		{1.90, 1.90, 1.90, 1.90},
		{1.90, 1.90, 1.90, 1.90},
		{1.90, 1.90, 1.90, 1.90},
		{1.90, 1.90, 1.90, 1.90},
		{1.90, 1.90, 1.90, 1.90},
		{1.90, 1.90, 1.90, 1.90},
		{1.50, 1.50, 1.50, 1.50},
		{1.50, 1.50, 1.50, 1.50},
		{1.50, 1.50, 1.50, 1.50},
		{1.50, 1.50, 1.50, 1.50},
		{1.50, 1.50, 1.50, 1.50},
		{1.50, 1.50, 1.50, 1.50},
		{1.50, 1.50, 1.50, 1.50},
		{1.50, 1.50, 1.50, 1.50},
		{1.25, 1.25, 1.25, 1.25},
		{1.25, 1.25, 1.25, 1.25},
		{1.25, 1.25, 1.25, 1.25},
		{1.25, 1.25, 1.25, 1.25},
		{1.25, 1.25, 1.25, 1.25},
		{1.25, 1.25, 1.25, 1.25},
		{1.25, 1.25, 1.25, 1.25},
		{1.25, 1.25, 1.25, 1.25},
		{1, 1, 1, 1},
		{0.75, 0.75, 0.75, 0.75},
		{0.9, 0.9, 0.9, 0.9},
		{0.75, 0.75, 0.75, 0.75},
		{0.9, 0.9, 0.9, 0.9},
		{0.8, 0.8, 0.8, 0.8},
	}
}

func getBonusMap() [38][4]float64 {
	return [38][4]float64{
		{3.25, 3.25, 3.25, 3.25},
		{2.75, 2.75, 2.75, 2.75},
		{2.75, 2.75, 2.75, 2.75},
		{2.75, 2.75, 2.75, 2.75},
		{2.75, 2.75, 2.75, 2.75},
		{2.25, 2.25, 2.25, 2.25},
		{2.25, 2.25, 2.25, 2.25},
		{2.25, 2.25, 2.25, 2.25},
		{2.25, 2.25, 2.25, 2.25},
		{2.25, 2.25, 2.25, 2.25},
		{1.90, 1.90, 1.90, 1.90},
		{1.90, 1.90, 1.90, 1.90},
		{1.90, 1.90, 1.90, 1.90},
		{1.90, 1.90, 1.90, 1.90},
		{1.90, 1.90, 1.90, 1.90},
		{1.90, 1.90, 1.90, 1.90},
		{1.50, 1.50, 1.50, 1.50},
		{1.50, 1.50, 1.50, 1.50},
		{1.50, 1.50, 1.50, 1.50},
		{1.50, 1.50, 1.50, 1.50},
		{1.50, 1.50, 1.50, 1.50},
		{1.50, 1.50, 1.50, 1.50},
		{1.50, 1.50, 1.50, 1.50},
		{1.50, 1.50, 1.50, 1.50},
		{1.25, 1.25, 1.25, 1.25},
		{1.25, 1.25, 1.25, 1.25},
		{1.25, 1.25, 1.25, 1.25},
		{1.25, 1.25, 1.25, 1.25},
		{1.25, 1.25, 1.25, 1.25},
		{1.25, 1.25, 1.25, 1.25},
		{1.25, 1.25, 1.25, 1.25},
		{1.25, 1.25, 1.25, 1.25},
		{1, 1, 1, 1},
		{0.75, 0.75, 0.75, 0.75},
		{0.3, 0.3, 0.3, 0.3},
		{0.25, 0.25, 0.25, 0.25},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}
}