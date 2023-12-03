package models

import (
	config "github.com/CalebRose/SimFBA/secrets"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
	"github.com/jinzhu/gorm"
)

type NFLDraftee struct {
	gorm.Model
	structs.BasePlayer
	PlayerID           int
	HighSchool         string
	CollegeID          uint
	College            string
	DraftedTeamID      uint
	DraftedTeam        string
	DraftedRound       uint
	DraftPickID        uint
	DraftedPick        uint
	City               string
	State              string
	OverallGrade       string
	StaminaGrade       string
	InjuryGrade        string
	FootballIQGrade    string
	SpeedGrade         string
	CarryingGrade      string
	AgilityGrade       string
	CatchingGrade      string
	RouteRunningGrade  string
	ZoneCoverageGrade  string
	ManCoverageGrade   string
	StrengthGrade      string
	TackleGrade        string
	PassBlockGrade     string
	RunBlockGrade      string
	PassRushGrade      string
	RunDefenseGrade    string
	ThrowPowerGrade    string
	ThrowAccuracyGrade string
	KickAccuracyGrade  string
	KickPowerGrade     string
	PuntAccuracyGrade  string
	PuntPowerGrade     string
	PrimeAge           uint
}

func (n *NFLDraftee) Map(cp structs.CollegePlayer) {
	attributeMeans := config.NFLAttributeMeans()
	n.ID = cp.ID
	n.PlayerID = cp.PlayerID
	n.HighSchool = cp.HighSchool
	n.College = cp.TeamAbbr
	n.City = cp.City
	n.State = cp.State
	n.FirstName = cp.FirstName
	n.LastName = cp.LastName
	n.Position = cp.Position
	n.Archetype = cp.Archetype
	n.Height = cp.Height
	n.Weight = cp.Weight
	n.Age = cp.Age
	n.Stars = cp.Stars
	n.Overall = cp.Overall
	n.Stamina = cp.Stamina
	n.Injury = cp.Injury
	n.FootballIQ = cp.FootballIQ
	n.Speed = cp.Speed
	n.Carrying = cp.Carrying
	n.Agility = cp.Agility
	n.Catching = cp.Catching
	n.RouteRunning = cp.RouteRunning
	n.ZoneCoverage = cp.ZoneCoverage
	n.ManCoverage = cp.ManCoverage
	n.Strength = cp.Strength
	n.Tackle = cp.Tackle
	n.PassBlock = cp.PassBlock
	n.RunBlock = cp.RunBlock
	n.PassRush = cp.PassRush
	n.RunDefense = cp.RunDefense
	n.ThrowPower = cp.ThrowPower
	n.ThrowAccuracy = cp.ThrowAccuracy
	n.KickAccuracy = cp.KickAccuracy
	n.KickPower = cp.KickPower
	n.PuntAccuracy = cp.PuntAccuracy
	n.PuntPower = cp.PuntPower
	n.Progression = cp.Progression
	n.Discipline = cp.Discipline
	n.PotentialGrade = cp.PotentialGrade
	n.FreeAgency = cp.FreeAgency
	n.Personality = cp.Personality
	n.RecruitingBias = cp.RecruitingBias
	n.WorkEthic = cp.WorkEthic
	n.AcademicBias = cp.AcademicBias
	rangeNum := 7
	OverallGrade := util.GetNFLOverallGrade(util.GenerateIntFromRange(cp.Overall-rangeNum, cp.Overall+rangeNum))
	StaminaGrade := util.GetLetterGrade(util.GenerateIntFromRange(cp.Stamina-rangeNum, cp.Stamina+rangeNum), attributeMeans["Stamina"][cp.Position]["mean"], attributeMeans["Stamina"][cp.Position]["stddev"], cp.Year)
	InjuryGrade := util.GetLetterGrade(util.GenerateIntFromRange(cp.Injury-rangeNum, cp.Injury+rangeNum), attributeMeans["Injury"][cp.Position]["mean"], attributeMeans["Injury"][cp.Position]["stddev"], cp.Year)
	SpeedGrade := util.GetLetterGrade(util.GenerateIntFromRange(cp.Speed-rangeNum, cp.Speed+rangeNum), attributeMeans["Speed"][cp.Position]["mean"], attributeMeans["Speed"][cp.Position]["stddev"], cp.Year)
	FootballIQGrade := util.GetLetterGrade(util.GenerateIntFromRange(cp.FootballIQ-rangeNum, cp.FootballIQ+rangeNum), attributeMeans["FootballIQ"][cp.Position]["mean"], attributeMeans["FootballIQ"][cp.Position]["stddev"], cp.Year)
	AgilityGrade := util.GetLetterGrade(util.GenerateIntFromRange(cp.Agility-rangeNum, cp.Agility+rangeNum), attributeMeans["Agility"][cp.Position]["mean"], attributeMeans["Agility"][cp.Position]["stddev"], cp.Year)
	CarryingGrade := util.GetLetterGrade(util.GenerateIntFromRange(cp.Carrying-rangeNum, cp.Carrying+rangeNum), attributeMeans["Carrying"][cp.Position]["mean"], attributeMeans["Carrying"][cp.Position]["stddev"], cp.Year)
	CatchingGrade := util.GetLetterGrade(util.GenerateIntFromRange(cp.Catching-rangeNum, cp.Catching+rangeNum), attributeMeans["Catching"][cp.Position]["mean"], attributeMeans["Catching"][cp.Position]["stddev"], cp.Year)
	RouteRunningGrade := util.GetLetterGrade(util.GenerateIntFromRange(cp.RouteRunning-rangeNum, cp.RouteRunning+rangeNum), attributeMeans["RouteRunning"][cp.Position]["mean"], attributeMeans["RouteRunning"][cp.Position]["stddev"], cp.Year)
	ZoneCoverageGrade := util.GetLetterGrade(util.GenerateIntFromRange(cp.ZoneCoverage-rangeNum, cp.ZoneCoverage+rangeNum), attributeMeans["ZoneCoverage"][cp.Position]["mean"], attributeMeans["ZoneCoverage"][cp.Position]["stddev"], cp.Year)
	ManCoverageGrade := util.GetLetterGrade(util.GenerateIntFromRange(cp.ManCoverage-rangeNum, cp.ManCoverage+rangeNum), attributeMeans["ManCoverage"][cp.Position]["mean"], attributeMeans["ManCoverage"][cp.Position]["stddev"], cp.Year)
	StrengthGrade := util.GetLetterGrade(util.GenerateIntFromRange(cp.Strength-rangeNum, cp.Strength+rangeNum), attributeMeans["Strength"][cp.Position]["mean"], attributeMeans["Strength"][cp.Position]["stddev"], cp.Year)
	TackleGrade := util.GetLetterGrade(util.GenerateIntFromRange(cp.Tackle-rangeNum, cp.Tackle+rangeNum), attributeMeans["Tackle"][cp.Position]["mean"], attributeMeans["Tackle"][cp.Position]["stddev"], cp.Year)
	PassBlockGrade := util.GetLetterGrade(util.GenerateIntFromRange(cp.PassBlock-rangeNum, cp.PassBlock+rangeNum), attributeMeans["PassBlock"][cp.Position]["mean"], attributeMeans["PassBlock"][cp.Position]["stddev"], cp.Year)
	RunBlockGrade := util.GetLetterGrade(util.GenerateIntFromRange(cp.RunBlock-rangeNum, cp.RunBlock+rangeNum), attributeMeans["RunBlock"][cp.Position]["mean"], attributeMeans["RunBlock"][cp.Position]["stddev"], cp.Year)
	PassRushGrade := util.GetLetterGrade(util.GenerateIntFromRange(cp.PassRush-rangeNum, cp.PassRush+rangeNum), attributeMeans["PassRush"][cp.Position]["mean"], attributeMeans["PassRush"][cp.Position]["stddev"], cp.Year)
	RunDefenseGrade := util.GetLetterGrade(util.GenerateIntFromRange(cp.RunDefense-rangeNum, cp.RunDefense+rangeNum), attributeMeans["RunDefense"][cp.Position]["mean"], attributeMeans["RunDefense"][cp.Position]["stddev"], cp.Year)
	ThrowPowerGrade := util.GetLetterGrade(util.GenerateIntFromRange(cp.ThrowPower-rangeNum, cp.ThrowPower+rangeNum), attributeMeans["ThrowPower"][cp.Position]["mean"], attributeMeans["ThrowPower"][cp.Position]["stddev"], cp.Year)
	ThrowAccuracyGrade := util.GetLetterGrade(util.GenerateIntFromRange(cp.ThrowAccuracy-rangeNum, cp.ThrowAccuracy+rangeNum), attributeMeans["ThrowAccuracy"][cp.Position]["mean"], attributeMeans["ThrowAccuracy"][cp.Position]["stddev"], cp.Year)
	KickPowerGrade := util.GetLetterGrade(util.GenerateIntFromRange(cp.KickPower-rangeNum, cp.KickPower+rangeNum), attributeMeans["KickPower"][cp.Position]["mean"], attributeMeans["KickPower"][cp.Position]["stddev"], cp.Year)
	KickAccuracyGrade := util.GetLetterGrade(util.GenerateIntFromRange(cp.KickAccuracy-rangeNum, cp.KickAccuracy+rangeNum), attributeMeans["KickAccuracy"][cp.Position]["mean"], attributeMeans["KickAccuracy"][cp.Position]["stddev"], cp.Year)
	PuntPowerGrade := util.GetLetterGrade(util.GenerateIntFromRange(cp.PuntAccuracy-rangeNum, cp.PuntPower+rangeNum), attributeMeans["PuntPower"][cp.Position]["mean"], attributeMeans["PuntPower"][cp.Position]["stddev"], cp.Year)
	PuntAccuracyGrade := util.GetLetterGrade(cp.PuntAccuracy, attributeMeans["PuntAccuracy"][cp.Position]["mean"], attributeMeans["PuntAccuracy"][cp.Position]["stddev"], cp.Year)
	n.OverallGrade = OverallGrade
	n.StaminaGrade = StaminaGrade
	n.InjuryGrade = InjuryGrade
	n.FootballIQGrade = FootballIQGrade
	n.SpeedGrade = SpeedGrade
	n.AgilityGrade = AgilityGrade
	n.CarryingGrade = CarryingGrade
	n.CatchingGrade = CatchingGrade
	n.RouteRunningGrade = RouteRunningGrade
	n.ZoneCoverageGrade = ZoneCoverageGrade
	n.ManCoverageGrade = ManCoverageGrade
	n.StrengthGrade = StrengthGrade
	n.TackleGrade = TackleGrade
	n.PassBlockGrade = PassBlockGrade
	n.RunBlockGrade = RunBlockGrade
	n.PassRushGrade = PassRushGrade
	n.RunDefenseGrade = RunDefenseGrade
	n.ThrowPowerGrade = ThrowPowerGrade
	n.ThrowAccuracyGrade = ThrowAccuracyGrade
	n.KickPowerGrade = KickPowerGrade
	n.KickAccuracyGrade = KickAccuracyGrade
	n.PuntPowerGrade = PuntPowerGrade
	n.PuntAccuracyGrade = PuntAccuracyGrade

}

func (n *NFLDraftee) MapUnsignedPlayer(up structs.UnsignedPlayer) {
	attributeMeans := config.NFLAttributeMeans()
	n.ID = up.ID
	n.PlayerID = int(up.PlayerID)
	n.HighSchool = up.HighSchool
	n.College = up.TeamAbbr
	n.City = up.City
	n.State = up.State
	n.FirstName = up.FirstName
	n.LastName = up.LastName
	n.Position = up.Position
	n.Archetype = up.Archetype
	n.Height = up.Height
	n.Weight = up.Weight
	n.Age = up.Age
	n.Stars = up.Stars
	n.Overall = up.Overall
	n.Stamina = up.Stamina
	n.Injury = up.Injury
	n.FootballIQ = up.FootballIQ
	n.Speed = up.Speed
	n.Carrying = up.Carrying
	n.Agility = up.Agility
	n.Catching = up.Catching
	n.RouteRunning = up.RouteRunning
	n.ZoneCoverage = up.ZoneCoverage
	n.ManCoverage = up.ManCoverage
	n.Strength = up.Strength
	n.Tackle = up.Tackle
	n.PassBlock = up.PassBlock
	n.RunBlock = up.RunBlock
	n.PassRush = up.PassRush
	n.RunDefense = up.RunDefense
	n.ThrowPower = up.ThrowPower
	n.ThrowAccuracy = up.ThrowAccuracy
	n.KickAccuracy = up.KickAccuracy
	n.KickPower = up.KickPower
	n.PuntAccuracy = up.PuntAccuracy
	n.PuntPower = up.PuntPower
	n.Progression = up.Progression
	n.Discipline = up.Discipline
	n.PotentialGrade = up.PotentialGrade
	n.FreeAgency = up.FreeAgency
	n.Personality = up.Personality
	n.RecruitingBias = up.RecruitingBias
	n.WorkEthic = up.WorkEthic
	n.AcademicBias = up.AcademicBias
	rangeNum := 7
	OverallGrade := util.GetNFLOverallGrade(util.GenerateIntFromRange(up.Overall-rangeNum, up.Overall+rangeNum))
	StaminaGrade := util.GetLetterGrade(util.GenerateIntFromRange(up.Stamina-rangeNum, up.Stamina+rangeNum), attributeMeans["Stamina"][up.Position]["mean"], attributeMeans["Stamina"][up.Position]["stddev"], up.Year)
	InjuryGrade := util.GetLetterGrade(util.GenerateIntFromRange(up.Injury-rangeNum, up.Injury+rangeNum), attributeMeans["Injury"][up.Position]["mean"], attributeMeans["Injury"][up.Position]["stddev"], up.Year)
	SpeedGrade := util.GetLetterGrade(util.GenerateIntFromRange(up.Speed-rangeNum, up.Speed+rangeNum), attributeMeans["Speed"][up.Position]["mean"], attributeMeans["Speed"][up.Position]["stddev"], up.Year)
	FootballIQGrade := util.GetLetterGrade(util.GenerateIntFromRange(up.FootballIQ-rangeNum, up.FootballIQ+rangeNum), attributeMeans["FootballIQ"][up.Position]["mean"], attributeMeans["FootballIQ"][up.Position]["stddev"], up.Year)
	AgilityGrade := util.GetLetterGrade(util.GenerateIntFromRange(up.Agility-rangeNum, up.Agility+rangeNum), attributeMeans["Agility"][up.Position]["mean"], attributeMeans["Agility"][up.Position]["stddev"], up.Year)
	CarryingGrade := util.GetLetterGrade(util.GenerateIntFromRange(up.Carrying-rangeNum, up.Carrying+rangeNum), attributeMeans["Carrying"][up.Position]["mean"], attributeMeans["Carrying"][up.Position]["stddev"], up.Year)
	CatchingGrade := util.GetLetterGrade(util.GenerateIntFromRange(up.Catching-rangeNum, up.Catching+rangeNum), attributeMeans["Catching"][up.Position]["mean"], attributeMeans["Catching"][up.Position]["stddev"], up.Year)
	RouteRunningGrade := util.GetLetterGrade(util.GenerateIntFromRange(up.RouteRunning-rangeNum, up.RouteRunning+rangeNum), attributeMeans["RouteRunning"][up.Position]["mean"], attributeMeans["RouteRunning"][up.Position]["stddev"], up.Year)
	ZoneCoverageGrade := util.GetLetterGrade(util.GenerateIntFromRange(up.ZoneCoverage-rangeNum, up.ZoneCoverage+rangeNum), attributeMeans["ZoneCoverage"][up.Position]["mean"], attributeMeans["ZoneCoverage"][up.Position]["stddev"], up.Year)
	ManCoverageGrade := util.GetLetterGrade(util.GenerateIntFromRange(up.ManCoverage-rangeNum, up.ManCoverage+rangeNum), attributeMeans["ManCoverage"][up.Position]["mean"], attributeMeans["ManCoverage"][up.Position]["stddev"], up.Year)
	StrengthGrade := util.GetLetterGrade(util.GenerateIntFromRange(up.Strength-rangeNum, up.Strength+rangeNum), attributeMeans["Strength"][up.Position]["mean"], attributeMeans["Strength"][up.Position]["stddev"], up.Year)
	TackleGrade := util.GetLetterGrade(util.GenerateIntFromRange(up.Tackle-rangeNum, up.Tackle+rangeNum), attributeMeans["Tackle"][up.Position]["mean"], attributeMeans["Tackle"][up.Position]["stddev"], up.Year)
	PassBlockGrade := util.GetLetterGrade(util.GenerateIntFromRange(up.PassBlock-rangeNum, up.PassBlock+rangeNum), attributeMeans["PassBlock"][up.Position]["mean"], attributeMeans["PassBlock"][up.Position]["stddev"], up.Year)
	RunBlockGrade := util.GetLetterGrade(util.GenerateIntFromRange(up.RunBlock-rangeNum, up.RunBlock+rangeNum), attributeMeans["RunBlock"][up.Position]["mean"], attributeMeans["RunBlock"][up.Position]["stddev"], up.Year)
	PassRushGrade := util.GetLetterGrade(util.GenerateIntFromRange(up.PassRush-rangeNum, up.PassRush+rangeNum), attributeMeans["PassRush"][up.Position]["mean"], attributeMeans["PassRush"][up.Position]["stddev"], up.Year)
	RunDefenseGrade := util.GetLetterGrade(util.GenerateIntFromRange(up.RunDefense-rangeNum, up.RunDefense+rangeNum), attributeMeans["RunDefense"][up.Position]["mean"], attributeMeans["RunDefense"][up.Position]["stddev"], up.Year)
	ThrowPowerGrade := util.GetLetterGrade(util.GenerateIntFromRange(up.ThrowPower-rangeNum, up.ThrowPower+rangeNum), attributeMeans["ThrowPower"][up.Position]["mean"], attributeMeans["ThrowPower"][up.Position]["stddev"], up.Year)
	ThrowAccuracyGrade := util.GetLetterGrade(util.GenerateIntFromRange(up.ThrowAccuracy-rangeNum, up.ThrowAccuracy+rangeNum), attributeMeans["ThrowAccuracy"][up.Position]["mean"], attributeMeans["ThrowAccuracy"][up.Position]["stddev"], up.Year)
	KickPowerGrade := util.GetLetterGrade(util.GenerateIntFromRange(up.KickPower-rangeNum, up.KickPower+rangeNum), attributeMeans["KickPower"][up.Position]["mean"], attributeMeans["KickPower"][up.Position]["stddev"], up.Year)
	KickAccuracyGrade := util.GetLetterGrade(util.GenerateIntFromRange(up.KickAccuracy-rangeNum, up.KickAccuracy+rangeNum), attributeMeans["KickAccuracy"][up.Position]["mean"], attributeMeans["KickAccuracy"][up.Position]["stddev"], up.Year)
	PuntPowerGrade := util.GetLetterGrade(util.GenerateIntFromRange(up.PuntAccuracy-rangeNum, up.PuntPower+rangeNum), attributeMeans["PuntPower"][up.Position]["mean"], attributeMeans["PuntPower"][up.Position]["stddev"], up.Year)
	PuntAccuracyGrade := util.GetLetterGrade(up.PuntAccuracy, attributeMeans["PuntAccuracy"][up.Position]["mean"], attributeMeans["PuntAccuracy"][up.Position]["stddev"], up.Year)
	n.OverallGrade = OverallGrade
	n.StaminaGrade = StaminaGrade
	n.InjuryGrade = InjuryGrade
	n.FootballIQGrade = FootballIQGrade
	n.SpeedGrade = SpeedGrade
	n.AgilityGrade = AgilityGrade
	n.CarryingGrade = CarryingGrade
	n.CatchingGrade = CatchingGrade
	n.RouteRunningGrade = RouteRunningGrade
	n.ZoneCoverageGrade = ZoneCoverageGrade
	n.ManCoverageGrade = ManCoverageGrade
	n.StrengthGrade = StrengthGrade
	n.TackleGrade = TackleGrade
	n.PassBlockGrade = PassBlockGrade
	n.RunBlockGrade = RunBlockGrade
	n.PassRushGrade = PassRushGrade
	n.RunDefenseGrade = RunDefenseGrade
	n.ThrowPowerGrade = ThrowPowerGrade
	n.ThrowAccuracyGrade = ThrowAccuracyGrade
	n.KickPowerGrade = KickPowerGrade
	n.KickAccuracyGrade = KickAccuracyGrade
	n.PuntPowerGrade = PuntPowerGrade
	n.PuntAccuracyGrade = PuntAccuracyGrade
}

func (n *NFLDraftee) MapNewOverallGrade(grade string) {
	n.OverallGrade = grade
}

func (n *NFLDraftee) AssignDraftedTeam(num uint, pickID uint, teamID uint, team string) {
	n.DraftedPick = num
	n.DraftPickID = pickID
	n.DraftedTeamID = teamID
	n.DraftedTeam = team
}
