package structs

import (
	"math/rand"

	"github.com/jinzhu/gorm"
)

type CollegePlayer struct {
	gorm.Model
	BasePlayer
	PlayerID      int
	TeamID        int
	TeamAbbr      string
	HighSchool    string
	City          string
	State         string
	Year          int
	IsRedshirt    bool
	IsRedshirting bool
	HasGraduated  bool
	Stats         []CollegePlayerStats     `gorm:"foreignKey:CollegePlayerID"`
	SeasonStats   CollegePlayerSeasonStats `gorm:"foreignKey:CollegePlayerID"`
	HasProgressed bool
}

type ByOverall []CollegePlayer

func (rp ByOverall) Len() int      { return len(rp) }
func (rp ByOverall) Swap(i, j int) { rp[i], rp[j] = rp[j], rp[i] }
func (rp ByOverall) Less(i, j int) bool {
	return rp[i].Overall > rp[j].Overall
}

func (p *CollegePlayer) SetRedshirtingStatus() {
	if !p.IsRedshirt && !p.IsRedshirting {
		p.IsRedshirting = true
	}
}

func (p *CollegePlayer) SetRedshirtStatus() {
	if p.IsRedshirting {
		p.IsRedshirting = false
		p.IsRedshirt = true
	}
}

func (cp *CollegePlayer) GetOverall() {
	var ovr float64 = 0
	if cp.Position == "QB" {
		ovr = (0.1 * float64(cp.Agility)) + (0.25 * float64(cp.ThrowPower)) + (0.25 * float64(cp.ThrowAccuracy)) + (0.1 * float64(cp.Speed)) + (0.2 * float64(cp.FootballIQ)) + (0.1 * float64(cp.Strength))
		cp.Overall = int(ovr)
	} else if cp.Position == "RB" {
		ovr = (0.2 * float64(cp.Agility)) + (0.05 * float64(cp.PassBlock)) +
			(0.1 * float64(cp.Carrying)) + (0.25 * float64(cp.Speed)) +
			(0.15 * float64(cp.FootballIQ)) + (0.2 * float64(cp.Strength)) +
			(0.05 * float64(cp.Catching))
		cp.Overall = int(ovr)
	} else if cp.Position == "FB" {
		ovr = (0.1 * float64(cp.Agility)) + (0.1 * float64(cp.PassBlock)) +
			(0.1 * float64(cp.Carrying)) + (0.05 * float64(cp.Speed)) +
			(0.15 * float64(cp.FootballIQ)) + (0.2 * float64(cp.Strength)) +
			(0.05 * float64(cp.Catching)) + (0.25 * float64(cp.RunBlock))
		cp.Overall = int(ovr)
	} else if cp.Position == "WR" {
		ovr = (0.15 * float64(cp.FootballIQ)) + (0.2 * float64(cp.Speed)) +
			(0.1 * float64(cp.Agility)) + (0.05 * float64(cp.Carrying)) +
			(0.05 * float64(cp.Strength)) + (0.25 * float64(cp.Catching)) +
			(0.2 * float64(cp.RouteRunning))
		cp.Overall = int(ovr)
	} else if cp.Position == "TE" {
		ovr = (0.15 * float64(cp.FootballIQ)) + (0.1 * float64(cp.Speed)) +
			(0.1 * float64(cp.Agility)) + (0.05 * float64(cp.Carrying)) +
			(0.05 * float64(cp.PassBlock)) + (0.15 * float64(cp.RunBlock)) +
			(0.1 * float64(cp.Strength)) + (0.20 * float64(cp.Catching)) +
			(0.1 * float64(cp.RouteRunning))
		cp.Overall = int(ovr)
	} else if cp.Position == "OT" || cp.Position == "OG" {
		ovr = (0.15 * float64(cp.FootballIQ)) + (0.05 * float64(cp.Agility)) +
			(0.3 * float64(cp.RunBlock)) + (0.2 * float64(cp.Strength)) +
			(0.3 * float64(cp.PassBlock))
		cp.Overall = int(ovr)
	} else if cp.Position == "C" {
		ovr = (0.2 * float64(cp.FootballIQ)) + (0.05 * float64(cp.Agility)) +
			(0.3 * float64(cp.RunBlock)) + (0.15 * float64(cp.Strength)) +
			(0.3 * float64(cp.PassBlock))
		cp.Overall = int(ovr)
	} else if cp.Position == "DT" {
		ovr = (0.15 * float64(cp.FootballIQ)) + (0.05 * float64(cp.Agility)) +
			(0.25 * float64(cp.RunDefense)) + (0.2 * float64(cp.Strength)) +
			(0.15 * float64(cp.PassRush)) + (0.2 * float64(cp.Tackle)) +
			(0.1 * float64(cp.Agility))
		cp.Overall = int(ovr)
	} else if cp.Position == "DE" {
		ovr = (0.15 * float64(cp.FootballIQ)) + (0.1 * float64(cp.Speed)) +
			(0.15 * float64(cp.RunDefense)) + (0.1 * float64(cp.Strength)) +
			(0.2 * float64(cp.PassRush)) + (0.2 * float64(cp.Tackle)) +
			(0.1 * float64(cp.Agility))
		cp.Overall = int(ovr)
	} else if cp.Position == "ILB" {
		ovr = (0.2 * float64(cp.FootballIQ)) + (0.1 * float64(cp.Speed)) +
			(0.15 * float64(cp.RunDefense)) + (0.1 * float64(cp.Strength)) +
			(0.1 * float64(cp.PassRush)) + (0.15 * float64(cp.Tackle)) +
			(0.1 * float64(cp.ZoneCoverage)) + (0.05 * float64(cp.ManCoverage)) +
			(0.05 * float64(cp.Agility))
		cp.Overall = int(ovr)
	} else if cp.Position == "OLB" {
		ovr = (0.15 * float64(cp.FootballIQ)) + (0.1 * float64(cp.Speed)) +
			(0.15 * float64(cp.RunDefense)) + (0.1 * float64(cp.Strength)) +
			(0.15 * float64(cp.PassRush)) + (0.15 * float64(cp.Tackle)) +
			(0.1 * float64(cp.ZoneCoverage)) + (0.05 * float64(cp.ManCoverage)) +
			(0.05 * float64(cp.Agility))
		cp.Overall = int(ovr)
	} else if cp.Position == "CB" {
		ovr = (0.15 * float64(cp.FootballIQ)) + (0.25 * float64(cp.Speed)) +
			(0.05 * float64(cp.Tackle)) + (0.05 * float64(cp.Strength)) +
			(0.15 * float64(cp.Agility)) + (0.15 * float64(cp.ZoneCoverage)) +
			(0.15 * float64(cp.ManCoverage)) + (0.05 * float64(cp.Catching))
		cp.Overall = int(ovr)
	} else if cp.Position == "FS" {
		ovr = (0.2 * float64(cp.FootballIQ)) + (0.2 * float64(cp.Speed)) +
			(0.05 * float64(cp.RunDefense)) + (0.05 * float64(cp.Strength)) +
			(0.05 * float64(cp.Catching)) + (0.05 * float64(cp.Tackle)) +
			(0.15 * float64(cp.ZoneCoverage)) + (0.15 * float64(cp.ManCoverage)) +
			(0.1 * float64(cp.Agility))
		cp.Overall = int(ovr)
	} else if cp.Position == "SS" {
		ovr = (0.15 * float64(cp.FootballIQ)) + (0.2 * float64(cp.Speed)) +
			(0.05 * float64(cp.RunDefense)) + (0.05 * float64(cp.Strength)) +
			(0.05 * float64(cp.Catching)) + (0.1 * float64(cp.Tackle)) +
			(0.15 * float64(cp.ZoneCoverage)) + (0.15 * float64(cp.ManCoverage)) +
			(0.1 * float64(cp.Agility))
		cp.Overall = int(ovr)
	} else if cp.Position == "K" {
		ovr = (0.2 * float64(cp.FootballIQ)) + (0.45 * float64(cp.KickPower)) +
			(0.45 * float64(cp.KickAccuracy))
		cp.Overall = int(ovr)
	} else if cp.Position == "P" {
		ovr = (0.2 * float64(cp.FootballIQ)) + (0.45 * float64(cp.PuntPower)) +
			(0.45 * float64(cp.PuntAccuracy))
		cp.Overall = int(ovr)
	} else if cp.Position == "ATH" {
		ovr = (float64(cp.FootballIQ) + float64(cp.Speed) + float64(cp.Agility) +
			float64(cp.Carrying) + float64(cp.Catching) + float64(cp.RouteRunning) +
			float64(cp.RunBlock) + float64(cp.PassBlock) + float64(cp.PassRush) +
			float64(cp.RunDefense) + float64(cp.Tackle) + float64(cp.Strength) +
			float64(cp.ZoneCoverage) + float64(cp.ManCoverage) + float64(cp.ThrowAccuracy) +
			float64(cp.ThrowPower) + float64(cp.PuntAccuracy) + float64(cp.PuntPower) +
			float64(cp.KickAccuracy) + float64(cp.KickPower)) / 20
		cp.Overall = int(ovr)
	}
}

func (cp *CollegePlayer) GetPotentialGrade() {
	adjust := rand.Intn(20) - 10
	if adjust == 0 {
		test := rand.Intn(2000) - 1000

		if test > 0 {
			adjust += 1
		} else if test < 0 {
			adjust -= 1
		} else {
			adjust = 0
		}
	}
	potential := cp.Progression + adjust
	if potential > 80 {
		cp.PotentialGrade = "A+"
	} else if potential > 70 {
		cp.PotentialGrade = "A"
	} else if potential > 65 {
		cp.PotentialGrade = "A-"
	} else if potential > 60 {
		cp.PotentialGrade = "B+"
	} else if potential > 55 {
		cp.PotentialGrade = "B"
	} else if potential > 50 {
		cp.PotentialGrade = "B-"
	} else if potential > 40 {
		cp.PotentialGrade = "C+"
	} else if potential > 30 {
		cp.PotentialGrade = "C"
	} else if potential > 25 {
		cp.PotentialGrade = "C-"
	} else if potential > 20 {
		cp.PotentialGrade = "D+"
	} else if potential > 15 {
		cp.PotentialGrade = "D"
	} else if potential > 10 {
		cp.PotentialGrade = "D-"
	} else {
		cp.PotentialGrade = "F"
	}
}

func (cp *CollegePlayer) Progress(attr CollegePlayerProgressions) {
	cp.Age++
	cp.Year++
	cp.Agility = attr.Agility
	cp.Speed = attr.Speed
	cp.FootballIQ = attr.FootballIQ
	cp.Carrying = attr.Carrying
	cp.Catching = attr.Catching
	cp.RouteRunning = attr.RouteRunning
	cp.PassBlock = attr.PassBlock
	cp.RunBlock = attr.RunBlock
	cp.PassRush = attr.PassRush
	cp.RunDefense = attr.RunDefense
	cp.Tackle = attr.Tackle
	cp.Strength = attr.Strength
	cp.ManCoverage = attr.ManCoverage
	cp.ZoneCoverage = attr.ZoneCoverage
	cp.KickAccuracy = attr.KickAccuracy
	cp.KickPower = attr.KickPower
	cp.PuntAccuracy = attr.PuntAccuracy
	cp.PuntPower = attr.PuntPower
	cp.ThrowAccuracy = attr.ThrowAccuracy
	cp.ThrowPower = attr.ThrowPower
	cp.HasProgressed = true
}

func (cp *CollegePlayer) GraduatePlayer() {
	cp.HasGraduated = true
}

func (cp *CollegePlayer) MapFromRecruit(r Recruit, t CollegeTeam) {
	cp.ID = r.ID
	cp.TeamID = int(t.ID)
	cp.TeamAbbr = t.TeamAbbr
	cp.PlayerID = r.PlayerID
	cp.HighSchool = r.HighSchool
	cp.City = r.City
	cp.State = r.State
	cp.Year = r.Age - 17
	cp.IsRedshirt = false
	cp.IsRedshirting = false
	cp.HasGraduated = false
	cp.Age = r.Age + 1
	cp.FirstName = r.FirstName
	cp.LastName = r.LastName
	cp.Position = r.Position
	cp.Archetype = r.Archetype
	cp.Height = r.Height
	cp.Weight = r.Weight
	cp.Age = r.Age
	cp.Stars = r.Stars
	cp.Overall = r.Overall
	cp.Stamina = r.Stamina
	cp.Injury = r.Injury
	cp.FootballIQ = r.FootballIQ
	cp.Speed = r.Speed
	cp.Carrying = r.Carrying
	cp.Agility = r.Agility
	cp.Catching = r.Catching
	cp.RouteRunning = r.RouteRunning
	cp.ZoneCoverage = r.ZoneCoverage
	cp.ManCoverage = r.ManCoverage
	cp.Strength = r.Strength
	cp.Tackle = r.Tackle
	cp.PassBlock = r.PassBlock
	cp.RunBlock = r.RunBlock
	cp.PassRush = r.PassRush
	cp.RunDefense = r.RunDefense
	cp.ThrowPower = r.ThrowPower
	cp.ThrowAccuracy = r.ThrowAccuracy
	cp.KickAccuracy = r.KickAccuracy
	cp.KickPower = r.KickPower
	cp.PuntAccuracy = r.PuntAccuracy
	cp.PuntPower = r.PuntPower
	cp.Progression = r.Progression
	cp.Discipline = r.Discipline
	cp.PotentialGrade = r.PotentialGrade
	cp.FreeAgency = r.FreeAgency
	cp.Personality = r.Personality
	cp.RecruitingBias = r.RecruitingBias
	cp.WorkEthic = r.WorkEthic
	cp.AcademicBias = r.AcademicBias
}

func (cp *CollegePlayer) AssignTeamValues(t CollegeTeam) {
	cp.TeamID = int(t.ID)
	cp.TeamAbbr = t.TeamAbbr
}
