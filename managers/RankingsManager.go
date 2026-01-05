package managers

import (
	"math"
	"math/rand"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/repository"
	config "github.com/CalebRose/SimFBA/secrets"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
)

func AssignAllRecruitRanks() {
	db := dbprovider.GetInstance().GetDB()

	var recruits []structs.Recruit

	// var recruitsToSync []structs.Recruit

	db.Order("overall desc").Where("stars > 0").Find(&recruits)

	rivalsMod := 100.0

	for _, croot := range recruits {
		// 247 Rankings
		rank247 := Get247Ranking(croot)
		// ESPN Rankings
		espnRank := GetESPNRanking(croot)

		// Rivals Ranking
		var rivalsRank float64 = 0
		rivalsBonus := rivalsMod
		rivalsRank = GetRivalsRanking(croot.Stars, rivalsBonus)

		var r float64 = croot.TopRankModifier

		if croot.TopRankModifier == 0 || croot.TopRankModifier < 0.95 || croot.TopRankModifier > 1.05 {
			r = 0.95 + rand.Float64()*(1.05-0.95)
		}

		if croot.Stars == 0 {
			rank247 = 0.001
			espnRank = 0.001
			rivalsRank = 0.001
			r = 1
		}

		croot.AssignRankValues(rank247, espnRank, rivalsRank, r)

		recruitingModifier := getRecruitingModifier()

		croot.AssignRecruitingModifier(recruitingModifier)
		shotgunVal := getShotgunVal()
		clutchVal := getClutchValue()
		croot.AssignNewAttributes(shotgunVal, clutchVal)

		repository.SaveRecruitRecord(croot, db)
		if rivalsMod > 0.1 {
			rivalsMod -= 0.1
		}

		// recruitsToSync = append(recruitsToSync, croot)
	}
}

func Get247Ranking(r structs.Recruit) float64 {
	ovr := r.Overall

	potentialGrade := Get247PotentialModifier(r.PotentialGrade)

	return float64(ovr) + (potentialGrade * 2)
}

func GetESPNRanking(r structs.Recruit) float64 {
	// ESPN Ranking = Star Rank + Archetype Modifier + weight difference + height difference
	// + potential val, and then round.

	starRank := GetESPNStarRank(r.Stars)
	archMod := GetArchetypeModifier(r.Archetype)
	potentialMod := GetESPNPotentialModifier(r.PotentialGrade)

	espnPositionMap := config.ESPNModifiers()
	heightMod := float64(r.Height) / espnPositionMap[r.Position]["Height"]
	weightMod := float64(r.Weight) / espnPositionMap[r.Position]["Weight"]
	espnRanking := math.Round(float64(starRank) + float64(archMod) + potentialMod + heightMod + weightMod)

	return espnRanking
}

func GetRivalsRanking(stars int, bonus float64) float64 {
	return GetRivalsStarModifier(stars) + bonus
}

func GetESPNStarRank(star int) int {
	switch star {
	case 5:
		return 95
	case 4:
		return 85
	case 3:
		return 75
	case 2:
		return 65
	}
	return 55
}

func GetArchetypeModifier(arch string) int {
	switch arch {
	case "Coverage", "Run Stopper", "Ball Hawk", "Man Coverage", "Pass Rusher", "Rushing", "Weakside", "Bandit", "Return Specialist", "Soccer Player", "Wingback":
		return 1
	case "Possession", "Field General", "Nose Tackle", "Lineman", "Blocking", "Line Captain":
		return -1
	case "Speed Rusher", "Pass Rush", "Scrambler", "Strongside", "Vertical Threat", "Triple-Threat", "Slotback", "Speed":
		return 2
	}
	return 0
}

func Get247PotentialModifier(pg string) float64 {
	switch pg {
	case "A+":
		return 5.83
	case "A":
		return 5.06
	case "A-":
		return 4.77
	case "B+":
		return 4.33
	case "B":
		return 4.04
	case "B-":
		return 3.87
	case "C+":
		return 3.58
	case "C":
		return 3.43
	case "C-":
		return 3.31
	case "D+":
		return 3.03
	case "D":
		return 2.77
	case "D-":
		return 2.67
	}
	return 2.3
}

func GetESPNPotentialModifier(pg string) float64 {
	switch pg {
	case "A+":
		return 1
	case "A":
		return 0.9
	case "A-":
		return 0.8
	case "B+":
		return 0.6
	case "B":
		return 0.4
	case "B-":
		return 0.2
	case "C+":
		return 0
	case "C":
		return -0.15
	case "C-":
		return -0.3
	case "D+":
		return -0.6
	case "D":
		return -0.75
	case "D-":
		return -0.9
	}
	return -1
}

func GetPredictiveOverall(r structs.Recruit) int {
	currentOverall := r.Overall

	var potentialProg int

	switch r.PotentialGrade {
	case "B+", "A-", "A", "A+":
		potentialProg = 7
	case "B", "B-", "C+":
		potentialProg = 5
	default:
		potentialProg = 4
	}

	return currentOverall + (potentialProg * 3)
}

func GetRivalsStarModifier(stars int) float64 {
	switch stars {
	case 5:
		return 6.1
	case 4:
		return RoundToFixedDecimalPlace(rand.Float64()*((6.0-5.8)+5.8), 1)
	case 3:
		return RoundToFixedDecimalPlace(rand.Float64()*((5.7-5.5)+5.5), 1)
	case 2:
		return RoundToFixedDecimalPlace(rand.Float64()*((5.4-5.2)+5.2), 1)
	default:
		return 5
	}
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func RoundToFixedDecimalPlace(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func Get247TeamRanking(rp structs.RecruitingTeamProfile, signedCroots []structs.Recruit) float64 {
	stddev := 10

	var Rank247 float64 = 0

	for idx, croot := range signedCroots {

		rank := float64((idx - 1) / stddev)

		expo := (-0.5 * (math.Pow(rank, 2)))

		weightedScore := (croot.Rank247 - 20) * math.Pow(math.E, expo)

		Rank247 += (weightedScore)
	}

	return Rank247
}

func GetESPNTeamRanking(rp structs.RecruitingTeamProfile, signedCroots []structs.Recruit) float64 {

	var espnRank float64 = 0

	for _, croot := range signedCroots {
		espnRank += croot.ESPNRank
	}

	return espnRank
}

func GetRivalsTeamRanking(rp structs.RecruitingTeamProfile, signedCroots []structs.Recruit) float64 {

	var rivalsRank float64 = 0

	for _, croot := range signedCroots {
		rivalsRank += croot.RivalsRank
	}

	return rivalsRank
}

func getRecruitingModifier() float64 {
	diceRoll := util.GenerateFloatFromRange(1, 20)
	if diceRoll == 1 {
		return 0.02
	}
	num := util.GenerateIntFromRange(1, 100)
	mod := 0.0
	if num < 11 {
		mod = util.GenerateFloatFromRange(1.80, 2.00)
	} else if num < 31 {
		mod = util.GenerateFloatFromRange(1.50, 1.69)
	} else if num < 71 {
		mod = util.GenerateFloatFromRange(1.15, 1.49)
	} else if num < 91 {
		mod = util.GenerateFloatFromRange(0.90, 1.14)
	} else {
		mod = util.GenerateFloatFromRange(0.80, 0.89)
	}

	return mod
}
