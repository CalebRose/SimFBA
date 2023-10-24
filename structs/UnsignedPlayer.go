package structs

import "github.com/jinzhu/gorm"

type UnsignedPlayer struct {
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

func (up *UnsignedPlayer) GraduatePlayer() {
	up.HasGraduated = true
}

func (up *UnsignedPlayer) Progress(attr CollegePlayerProgressions) {
	up.Age++
	up.Year++
	up.Agility = attr.Agility
	up.Speed = attr.Speed
	up.FootballIQ = attr.FootballIQ
	up.Carrying = attr.Carrying
	up.Catching = attr.Catching
	up.RouteRunning = attr.RouteRunning
	up.PassBlock = attr.PassBlock
	up.RunBlock = attr.RunBlock
	up.PassRush = attr.PassRush
	up.RunDefense = attr.RunDefense
	up.Tackle = attr.Tackle
	up.Strength = attr.Strength
	up.ManCoverage = attr.ManCoverage
	up.ZoneCoverage = attr.ZoneCoverage
	up.KickAccuracy = attr.KickAccuracy
	up.KickPower = attr.KickPower
	up.PuntAccuracy = attr.PuntAccuracy
	up.PuntPower = attr.PuntPower
	up.ThrowAccuracy = attr.ThrowAccuracy
	up.ThrowPower = attr.ThrowPower
	up.HasProgressed = true
}

func (up *UnsignedPlayer) MapFromRecruit(r Recruit) {
	up.ID = r.ID
	up.TeamID = 0
	up.TeamAbbr = ""
	up.PlayerID = r.PlayerID
	up.HighSchool = r.HighSchool
	up.City = r.City
	up.State = r.State
	up.Year = r.Age - 17
	up.IsRedshirt = false
	up.IsRedshirting = false
	up.HasGraduated = false
	up.Age = r.Age + 1
	up.FirstName = r.FirstName
	up.LastName = r.LastName
	up.Position = r.Position
	up.Archetype = r.Archetype
	up.Height = r.Height
	up.Weight = r.Weight
	up.Age = r.Age
	up.Stars = r.Stars
	up.Overall = r.Overall
	up.Stamina = r.Stamina
	up.Injury = r.Injury
	up.FootballIQ = r.FootballIQ
	up.Speed = r.Speed
	up.Carrying = r.Carrying
	up.Agility = r.Agility
	up.Catching = r.Catching
	up.RouteRunning = r.RouteRunning
	up.ZoneCoverage = r.ZoneCoverage
	up.ManCoverage = r.ManCoverage
	up.Strength = r.Strength
	up.Tackle = r.Tackle
	up.PassBlock = r.PassBlock
	up.RunBlock = r.RunBlock
	up.PassRush = r.PassRush
	up.RunDefense = r.RunDefense
	up.ThrowPower = r.ThrowPower
	up.ThrowAccuracy = r.ThrowAccuracy
	up.KickAccuracy = r.KickAccuracy
	up.KickPower = r.KickPower
	up.PuntAccuracy = r.PuntAccuracy
	up.PuntPower = r.PuntPower
	up.Progression = r.Progression
	up.Discipline = r.Discipline
	up.PotentialGrade = r.PotentialGrade
	up.FreeAgency = r.FreeAgency
	up.Personality = r.Personality
	up.RecruitingBias = r.RecruitingBias
	up.WorkEthic = r.WorkEthic
	up.AcademicBias = r.AcademicBias
}
