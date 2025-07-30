package structs

import (
	"math/rand"

	"github.com/jinzhu/gorm"
)

type CollegePlayer struct {
	gorm.Model
	BasePlayer
	PlayerID           int
	TeamID             int
	TeamAbbr           string
	HighSchool         string
	City               string
	State              string
	Year               int
	IsRedshirt         bool
	IsRedshirting      bool
	HasGraduated       bool
	TransferStatus     int                      // 1 == Intends, 2 == Is Transferring
	TransferLikeliness string                   // Low, Medium, High
	Stats              []CollegePlayerStats     `gorm:"foreignKey:CollegePlayerID"`
	SeasonStats        CollegePlayerSeasonStats `gorm:"foreignKey:CollegePlayerID"`
	HasProgressed      bool
	WillDeclare        bool
	LegacyID           uint // Either a legacy school or a legacy coach
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
			(0.15 * float64(cp.PassRush)) + (0.2 * float64(cp.Tackle))
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
		if cp.Archetype == "Field General" {
			ovr = (.20 * float64(cp.FootballIQ)) + (.1 * float64(cp.ZoneCoverage)) + (.1 * float64(cp.ManCoverage)) + (.1 * float64(cp.RunDefense)) + (.1 * float64(cp.Speed)) + (.1 * float64(cp.Strength)) + (.1 * float64(cp.Tackle)) + (.1 * float64(cp.ThrowPower)) + (.1 * float64(cp.ThrowAccuracy))
		} else if cp.Archetype == "Triple-Threat" {
			ovr = (.10 * float64(cp.FootballIQ)) + (.2 * float64(cp.Agility)) + (.1 * float64(cp.Carrying)) + (.1 * float64(cp.Catching)) + (.2 * float64(cp.Speed)) + (.1 * float64(cp.RouteRunning)) + (.1 * float64(cp.ThrowPower)) + (.1 * float64(cp.ThrowAccuracy))
		} else if cp.Archetype == "Wingback" {
			ovr = (.1 * float64(cp.FootballIQ)) + (.2 * float64(cp.Agility)) + (.1 * float64(cp.Carrying)) + (.2 * float64(cp.Catching)) + (.2 * float64(cp.Speed)) + (.1 * float64(cp.RouteRunning)) + (.1 * float64(cp.RunBlock))
		} else if cp.Archetype == "Slotback" {
			ovr = (.1 * float64(cp.FootballIQ)) + (.1 * float64(cp.Agility)) + (.1 * float64(cp.Carrying)) + (.1 * float64(cp.Catching)) + (.2 * float64(cp.Speed)) + (.1 * float64(cp.RouteRunning)) + (.1 * float64(cp.RunBlock)) + (.1 * float64(cp.PassBlock)) + (.1 * float64(cp.Strength))
		} else if cp.Archetype == "Lineman" {
			ovr = (.1 * float64(cp.FootballIQ)) + (.1 * float64(cp.Agility)) + (.1 * float64(cp.RunBlock)) + (.1 * float64(cp.PassBlock)) + (.3 * float64(cp.Strength)) + (.1 * float64(cp.PassRush)) + (.1 * float64(cp.RunDefense)) + (.1 * float64(cp.Tackle))
		} else if cp.Archetype == "Strongside" {
			ovr = (.1 * float64(cp.FootballIQ)) + (.1 * float64(cp.Agility)) + (.1 * float64(cp.ZoneCoverage)) + (.1 * float64(cp.ManCoverage)) + (.2 * float64(cp.Strength)) + (.1 * float64(cp.PassRush)) + (.1 * float64(cp.RunDefense)) + (.1 * float64(cp.Tackle)) + (.1 * float64(cp.Speed))
		} else if cp.Archetype == "Weakside" {
			ovr = (.1 * float64(cp.FootballIQ)) + (.1 * float64(cp.Agility)) + (.1 * float64(cp.ZoneCoverage)) + (.1 * float64(cp.ManCoverage)) + (.1 * float64(cp.Strength)) + (.1 * float64(cp.PassRush)) + (.1 * float64(cp.RunDefense)) + (.1 * float64(cp.Tackle)) + (.2 * float64(cp.Speed))
		} else if cp.Archetype == "Bandit" {
			ovr = (.1 * float64(cp.FootballIQ)) + (.1 * float64(cp.Agility)) + (.1 * float64(cp.ZoneCoverage)) + (.1 * float64(cp.ManCoverage)) + (.1 * float64(cp.Strength)) + (.1 * float64(cp.PassRush)) + (.1 * float64(cp.RunDefense)) + (.1 * float64(cp.Tackle)) + (.2 * float64(cp.Speed))
		} else if cp.Archetype == "Return Specialist" {
			ovr = (.20 * float64(cp.FootballIQ)) + (.10 * float64(cp.Strength)) + (.20 * float64(cp.Speed)) + (.20 * float64(cp.Agility)) + (.20 * float64(cp.Catching)) + (.1 * float64(cp.Tackle))
		} else if cp.Archetype == "Soccer Player" {
			ovr = (.10 * float64(cp.FootballIQ)) + (.10 * float64(cp.Agility)) + (.2 * float64(cp.KickPower)) + (.2 * float64(cp.KickAccuracy)) + (.2 * float64(cp.PuntPower)) + (.2 * float64(cp.PuntAccuracy))
		}
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

func (cp *CollegePlayer) Progress(attr CollegePlayerProgressions, isBoomBust bool) {
	if !isBoomBust {
		cp.Age++
		cp.Year++
	}
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

func (cp *CollegePlayer) DeclareTransferIntention(weight int) {
	cp.TransferStatus = 1
	if weight < 30 {
		cp.TransferLikeliness = "Low"
	} else if weight < 70 {
		cp.TransferLikeliness = "Medium"
	} else {
		cp.TransferLikeliness = "High"
	}
}

func (cp *CollegePlayer) WillStay() {
	cp.TransferStatus = 0
	cp.WillDeclare = false
}

func (cp *CollegePlayer) WillNotTransfer() {
	cp.TransferStatus = 0
}

func (cp *CollegePlayer) WillTransfer() {
	cp.TransferStatus = 2
	cp.PreviousTeam = cp.TeamAbbr
	cp.PreviousTeamID = uint(cp.TeamID)
	cp.TeamAbbr = ""
	cp.TeamID = 0
}

func (cp *CollegePlayer) WillReturn() {
	cp.TransferStatus = 0
	cp.TeamAbbr = cp.PreviousTeam
	cp.TeamID = int(cp.PreviousTeamID)
	cp.PreviousTeam = ""
	cp.PreviousTeamID = 0
}

func (b *CollegePlayer) DismissFromTeam() {
	b.PreviousTeamID = uint(b.TeamID)
	b.PreviousTeam = b.TeamAbbr
	b.TeamID = 0
	b.TeamAbbr = ""
	b.TransferStatus = 2
}

func (cp *CollegePlayer) SignWithNewTeam(teamID int, teamAbbr string) {
	cp.TransferStatus = 0
	cp.TeamAbbr = teamAbbr
	cp.TeamID = teamID
	cp.TransferLikeliness = ""
}

func (cp *CollegePlayer) AddSeasonStats(seasonStats CollegePlayerSeasonStats) {
	cp.SeasonStats = seasonStats
}

func (cp *CollegePlayer) RevertRedshirting(isRS bool) {
	cp.IsRedshirt = false
	cp.IsRedshirting = isRS
}

func (cp *CollegePlayer) RevertYearProgression() {
	cp.RevertAge()
	cp.Year--
}
