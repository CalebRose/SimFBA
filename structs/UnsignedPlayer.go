package structs

import "github.com/jinzhu/gorm"

type UnsignedPlayer struct {
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
	TransferStatus     int // 1 == Intends, 2 == Is Transferring
	TransferLikeliness string
	Stats              []CollegePlayerStats     `gorm:"foreignKey:CollegePlayerID"`
	SeasonStats        CollegePlayerSeasonStats `gorm:"foreignKey:CollegePlayerID"`
	HasProgressed      bool
	WillDeclare        bool
	LegacyID           uint // Either a legacy school or a legacy coach
}

func (up *UnsignedPlayer) GraduatePlayer() {
	up.HasGraduated = true
}

func (up *UnsignedPlayer) Progress(attr CollegePlayerProgressions) {
	up.Age++
	up.Year++
	up.Agility = int8(attr.Agility)
	up.Speed = int8(attr.Speed)
	up.FootballIQ = int8(attr.FootballIQ)
	up.Carrying = int8(attr.Carrying)
	up.Catching = int8(attr.Catching)
	up.RouteRunning = int8(attr.RouteRunning)
	up.PassBlock = int8(attr.PassBlock)
	up.RunBlock = int8(attr.RunBlock)
	up.PassRush = int8(attr.PassRush)
	up.RunDefense = int8(attr.RunDefense)
	up.Tackle = int8(attr.Tackle)
	up.Strength = int8(attr.Strength)
	up.ManCoverage = int8(attr.ManCoverage)
	up.ZoneCoverage = int8(attr.ZoneCoverage)
	up.KickAccuracy = int8(attr.KickAccuracy)
	up.KickPower = int8(attr.KickPower)
	up.PuntAccuracy = int8(attr.PuntAccuracy)
	up.PuntPower = int8(attr.PuntPower)
	up.ThrowAccuracy = int8(attr.ThrowAccuracy)
	up.ThrowPower = int8(attr.ThrowPower)
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
	up.Year = int(r.Age) - 17
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
