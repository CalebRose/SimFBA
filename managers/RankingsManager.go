package managers

import (
	"math"
	"math/rand"
	"time"

	"github.com/CalebRose/SimFBA/dbprovider"
	config "github.com/CalebRose/SimFBA/secrets"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
)

func AssignAllRecruitRanks() {
	db := dbprovider.GetInstance().GetDB()

	rand.Seed(time.Now().UnixNano())

	var recruits []structs.Recruit

	// var recruitsToSync []structs.Recruit

	db.Order("overall desc").Find(&recruits)

	rivalsModifiers := config.RivalsModifiers()

	for idx, croot := range recruits {
		// 247 Rankings
		rank247 := Get247Ranking(croot)
		// ESPN Rankings
		espnRank := GetESPNRanking(croot)

		// Rivals Ranking
		var rivalsRank float64 = 0
		if idx <= 249 {
			rivalsBonus := rivalsModifiers[idx]

			rivalsRank = GetRivalsRanking(croot.Stars, rivalsBonus)
		}

		var r float64 = croot.TopRankModifier

		if croot.TopRankModifier == 0 || croot.TopRankModifier < 0.95 || croot.TopRankModifier > 1.05 {
			r = 0.95 + rand.Float64()*(1.05-0.95)
		}

		croot.AssignRankValues(rank247, espnRank, rivalsRank, r)

		recruitingModifier := getRecruitingModifier(croot.Stars)

		croot.AssignRecruitingModifier(recruitingModifier)

		db.Save(&croot)

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

func GetRivalsRanking(stars int, bonus int) float64 {
	return GetRivalsStarModifier(stars) + float64(bonus)
}

func GetESPNStarRank(star int) int {
	if star == 5 {
		return 95
	} else if star == 4 {
		return 85
	} else if star == 3 {
		return 75
	} else if star == 2 {
		return 65
	}
	return 55
}

func GetArchetypeModifier(arch string) int {
	if arch == "Coverage" ||
		arch == "Run Stopper" ||
		arch == "Ball Hawk" ||
		arch == "Man Coverage" ||
		arch == "Pass Rusher" ||
		arch == "Rushing" {
		return 1
	} else if arch == "Possession" ||
		arch == "Field General" ||
		arch == "Nose Tackle" ||
		arch == "Blocking" ||
		arch == "Line Captain" {
		return -1
	} else if arch == "Speed Rusher" ||
		arch == "Pass Rush" || arch == "Scrambler" ||
		arch == "Vertical Threat" ||
		arch == "Speed" {
		return 2
	}
	return 0
}

func Get247PotentialModifier(pg string) float64 {
	if pg == "A+" {
		return 5.83
	} else if pg == "A" {
		return 5.06
	} else if pg == "A-" {
		return 4.77
	} else if pg == "B+" {
		return 4.33
	} else if pg == "B" {
		return 4.04
	} else if pg == "B-" {
		return 3.87
	} else if pg == "C+" {
		return 3.58
	} else if pg == "C" {
		return 3.43
	} else if pg == "C-" {
		return 3.31
	} else if pg == "D+" {
		return 3.03
	} else if pg == "D" {
		return 2.77
	} else if pg == "D-" {
		return 2.67
	}
	return 2.3
}

func GetESPNPotentialModifier(pg string) float64 {
	if pg == "A+" {
		return 1
	} else if pg == "A" {
		return 0.9
	} else if pg == "A-" {
		return 0.8
	} else if pg == "B+" {
		return 0.6
	} else if pg == "B" {
		return 0.4
	} else if pg == "B-" {
		return 0.2
	} else if pg == "C+" {
		return 0
	} else if pg == "C" {
		return -0.15
	} else if pg == "C-" {
		return -0.3
	} else if pg == "D+" {
		return -0.6
	} else if pg == "D" {
		return -0.75
	} else if pg == "D-" {
		return -0.9
	}
	return -1
}

func GetPredictiveOverall(r structs.Recruit) int {
	currentOverall := r.Overall

	var potentialProg int

	if r.PotentialGrade == "B+" ||
		r.PotentialGrade == "A-" ||
		r.PotentialGrade == "A" ||
		r.PotentialGrade == "A+" {
		potentialProg = 7
	} else if r.PotentialGrade == "B" ||
		r.PotentialGrade == "B-" ||
		r.PotentialGrade == "C+" {
		potentialProg = 5
	} else {
		potentialProg = 4
	}

	return currentOverall + (potentialProg * 3)
}

func GetRivalsStarModifier(stars int) float64 {
	rand.Seed(time.Now().UnixNano())
	if stars == 5 {
		return 6.1
	} else if stars == 4 {
		return RoundToFixedDecimalPlace(rand.Float64()*((6.0-5.8)+5.8), 1)
	} else if stars == 3 {
		return RoundToFixedDecimalPlace(rand.Float64()*((5.7-5.5)+5.5), 1)
	} else if stars == 2 {
		return RoundToFixedDecimalPlace(rand.Float64()*((5.4-5.2)+5.2), 1)
	} else {
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

func getRecruitingModifier(stars int) float64 {
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
