package structs

import "github.com/jinzhu/gorm"

type HistoricCollegePlayer struct {
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

func (hcp *HistoricCollegePlayer) MapUnsignedPlayer(up UnsignedPlayer) {
	hcp.ID = up.ID
	hcp.TeamID = 0
	hcp.TeamAbbr = up.TeamAbbr
	hcp.PlayerID = up.PlayerID
	hcp.HighSchool = up.HighSchool
	hcp.City = up.City
	hcp.State = up.State
	hcp.Year = up.Year
	hcp.IsRedshirt = false
	hcp.IsRedshirting = false
	hcp.HasGraduated = false
	hcp.Age = up.Age
	hcp.FirstName = up.FirstName
	hcp.LastName = up.LastName
	hcp.Position = up.Position
	hcp.Archetype = up.Archetype
	hcp.Height = up.Height
	hcp.Weight = up.Weight
	hcp.Age = up.Age
	hcp.Stars = up.Stars
	hcp.Overall = up.Overall
	hcp.Stamina = up.Stamina
	hcp.Injury = up.Injury
	hcp.FootballIQ = up.FootballIQ
	hcp.Speed = up.Speed
	hcp.Carrying = up.Carrying
	hcp.Agility = up.Agility
	hcp.Catching = up.Catching
	hcp.RouteRunning = up.RouteRunning
	hcp.ZoneCoverage = up.ZoneCoverage
	hcp.ManCoverage = up.ManCoverage
	hcp.Strength = up.Strength
	hcp.Tackle = up.Tackle
	hcp.PassBlock = up.PassBlock
	hcp.RunBlock = up.RunBlock
	hcp.PassRush = up.PassRush
	hcp.RunDefense = up.RunDefense
	hcp.ThrowPower = up.ThrowPower
	hcp.ThrowAccuracy = up.ThrowAccuracy
	hcp.KickAccuracy = up.KickAccuracy
	hcp.KickPower = up.KickPower
	hcp.PuntAccuracy = up.PuntAccuracy
	hcp.PuntPower = up.PuntPower
	hcp.Progression = up.Progression
	hcp.Discipline = up.Discipline
	hcp.PotentialGrade = up.PotentialGrade
	hcp.FreeAgency = up.FreeAgency
	hcp.Personality = up.Personality
	hcp.RecruitingBias = up.RecruitingBias
	hcp.WorkEthic = up.WorkEthic
	hcp.AcademicBias = up.AcademicBias
}
