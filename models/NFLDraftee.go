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
	BoomOrBust         bool
	BoomOrBustStatus   string
}

func (n *NFLDraftee) Map(cp structs.CollegePlayer) {
	n.ID = cp.ID
	n.PlayerID = cp.PlayerID
	n.HighSchool = cp.HighSchool
	n.College = cp.TeamAbbr
	n.CollegeID = uint(cp.TeamID)
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
	n.PositionTwo = cp.PositionTwo
	n.ArchetypeTwo = cp.ArchetypeTwo
	n.PrimeAge = cp.PrimeAge
}

func (n *NFLDraftee) GetLetterGrades() {
	attributeMeans := config.NFLAttributeMeans()
	rangeNum := 7
	OverallGrade := util.GetNFLOverallGrade(util.GenerateIntFromRange(int(n.Overall)-rangeNum, int(n.Overall)+rangeNum))
	StaminaGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(n.Stamina)-rangeNum, int(n.Stamina)+rangeNum), attributeMeans["Stamina"][n.Position]["mean"], attributeMeans["Stamina"][n.Position]["stddev"], 5)
	InjuryGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(n.Injury)-rangeNum, int(n.Injury)+rangeNum), attributeMeans["Injury"][n.Position]["mean"], attributeMeans["Injury"][n.Position]["stddev"], 5)
	SpeedGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(n.Speed)-rangeNum, int(n.Speed)+rangeNum), attributeMeans["Speed"][n.Position]["mean"], attributeMeans["Speed"][n.Position]["stddev"], 5)
	FootballIQGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(n.FootballIQ)-rangeNum, int(n.FootballIQ)+rangeNum), attributeMeans["FootballIQ"][n.Position]["mean"], attributeMeans["FootballIQ"][n.Position]["stddev"], 5)
	AgilityGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(n.Agility)-rangeNum, int(n.Agility)+rangeNum), attributeMeans["Agility"][n.Position]["mean"], attributeMeans["Agility"][n.Position]["stddev"], 5)
	CarryingGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(n.Carrying)-rangeNum, int(n.Carrying)+rangeNum), attributeMeans["Carrying"][n.Position]["mean"], attributeMeans["Carrying"][n.Position]["stddev"], 5)
	CatchingGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(n.Catching)-rangeNum, int(n.Catching)+rangeNum), attributeMeans["Catching"][n.Position]["mean"], attributeMeans["Catching"][n.Position]["stddev"], 5)
	RouteRunningGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(n.RouteRunning)-rangeNum, int(n.RouteRunning)+rangeNum), attributeMeans["RouteRunning"][n.Position]["mean"], attributeMeans["RouteRunning"][n.Position]["stddev"], 5)
	ZoneCoverageGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(n.ZoneCoverage)-rangeNum, int(n.ZoneCoverage)+rangeNum), attributeMeans["ZoneCoverage"][n.Position]["mean"], attributeMeans["ZoneCoverage"][n.Position]["stddev"], 5)
	ManCoverageGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(n.ManCoverage)-rangeNum, int(n.ManCoverage)+rangeNum), attributeMeans["ManCoverage"][n.Position]["mean"], attributeMeans["ManCoverage"][n.Position]["stddev"], 5)
	StrengthGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(n.Strength)-rangeNum, int(n.Strength)+rangeNum), attributeMeans["Strength"][n.Position]["mean"], attributeMeans["Strength"][n.Position]["stddev"], 5)
	TackleGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(n.Tackle)-rangeNum, int(n.Tackle)+rangeNum), attributeMeans["Tackle"][n.Position]["mean"], attributeMeans["Tackle"][n.Position]["stddev"], 5)
	PassBlockGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(n.PassBlock)-rangeNum, int(n.PassBlock)+rangeNum), attributeMeans["PassBlock"][n.Position]["mean"], attributeMeans["PassBlock"][n.Position]["stddev"], 5)
	RunBlockGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(n.RunBlock)-rangeNum, int(n.RunBlock)+rangeNum), attributeMeans["RunBlock"][n.Position]["mean"], attributeMeans["RunBlock"][n.Position]["stddev"], 5)
	PassRushGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(n.PassRush)-rangeNum, int(n.PassRush)+rangeNum), attributeMeans["PassRush"][n.Position]["mean"], attributeMeans["PassRush"][n.Position]["stddev"], 5)
	RunDefenseGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(n.RunDefense)-rangeNum, int(n.RunDefense)+rangeNum), attributeMeans["RunDefense"][n.Position]["mean"], attributeMeans["RunDefense"][n.Position]["stddev"], 5)
	ThrowPowerGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(n.ThrowPower)-rangeNum, int(n.ThrowPower)+rangeNum), attributeMeans["ThrowPower"][n.Position]["mean"], attributeMeans["ThrowPower"][n.Position]["stddev"], 5)
	ThrowAccuracyGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(n.ThrowAccuracy)-rangeNum, int(n.ThrowAccuracy)+rangeNum), attributeMeans["ThrowAccuracy"][n.Position]["mean"], attributeMeans["ThrowAccuracy"][n.Position]["stddev"], 5)
	KickPowerGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(n.KickPower)-rangeNum, int(n.KickPower)+rangeNum), attributeMeans["KickPower"][n.Position]["mean"], attributeMeans["KickPower"][n.Position]["stddev"], 5)
	KickAccuracyGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(n.KickAccuracy)-rangeNum, int(n.KickAccuracy)+rangeNum), attributeMeans["KickAccuracy"][n.Position]["mean"], attributeMeans["KickAccuracy"][n.Position]["stddev"], 5)
	PuntPowerGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(n.PuntPower)-rangeNum, int(n.PuntPower)+rangeNum), attributeMeans["PuntPower"][n.Position]["mean"], attributeMeans["PuntPower"][n.Position]["stddev"], 5)
	PuntAccuracyGrade := util.GetLetterGrade(int(n.PuntAccuracy), attributeMeans["PuntAccuracy"][n.Position]["mean"], attributeMeans["PuntAccuracy"][n.Position]["stddev"], 5)
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
	OverallGrade := util.GetNFLOverallGrade(util.GenerateIntFromRange(int(up.Overall)-rangeNum, int(up.Overall)+rangeNum))
	StaminaGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(up.Stamina)-rangeNum, int(up.Stamina)+rangeNum), attributeMeans["Stamina"][up.Position]["mean"], attributeMeans["Stamina"][up.Position]["stddev"], up.Year)
	InjuryGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(up.Injury)-rangeNum, int(up.Injury)+rangeNum), attributeMeans["Injury"][up.Position]["mean"], attributeMeans["Injury"][up.Position]["stddev"], up.Year)
	SpeedGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(up.Speed)-rangeNum, int(up.Speed)+rangeNum), attributeMeans["Speed"][up.Position]["mean"], attributeMeans["Speed"][up.Position]["stddev"], up.Year)
	FootballIQGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(up.FootballIQ)-rangeNum, int(up.FootballIQ)+rangeNum), attributeMeans["FootballIQ"][up.Position]["mean"], attributeMeans["FootballIQ"][up.Position]["stddev"], up.Year)
	AgilityGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(up.Agility)-rangeNum, int(up.Agility)+rangeNum), attributeMeans["Agility"][up.Position]["mean"], attributeMeans["Agility"][up.Position]["stddev"], up.Year)
	CarryingGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(up.Carrying)-rangeNum, int(up.Carrying)+rangeNum), attributeMeans["Carrying"][up.Position]["mean"], attributeMeans["Carrying"][up.Position]["stddev"], up.Year)
	CatchingGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(up.Catching)-rangeNum, int(up.Catching)+rangeNum), attributeMeans["Catching"][up.Position]["mean"], attributeMeans["Catching"][up.Position]["stddev"], up.Year)
	RouteRunningGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(up.RouteRunning)-rangeNum, int(up.RouteRunning)+rangeNum), attributeMeans["RouteRunning"][up.Position]["mean"], attributeMeans["RouteRunning"][up.Position]["stddev"], up.Year)
	ZoneCoverageGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(up.ZoneCoverage)-rangeNum, int(up.ZoneCoverage)+rangeNum), attributeMeans["ZoneCoverage"][up.Position]["mean"], attributeMeans["ZoneCoverage"][up.Position]["stddev"], up.Year)
	ManCoverageGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(up.ManCoverage)-rangeNum, int(up.ManCoverage)+rangeNum), attributeMeans["ManCoverage"][up.Position]["mean"], attributeMeans["ManCoverage"][up.Position]["stddev"], up.Year)
	StrengthGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(up.Strength)-rangeNum, int(up.Strength)+rangeNum), attributeMeans["Strength"][up.Position]["mean"], attributeMeans["Strength"][up.Position]["stddev"], up.Year)
	TackleGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(up.Tackle)-rangeNum, int(up.Tackle)+rangeNum), attributeMeans["Tackle"][up.Position]["mean"], attributeMeans["Tackle"][up.Position]["stddev"], up.Year)
	PassBlockGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(up.PassBlock)-rangeNum, int(up.PassBlock)+rangeNum), attributeMeans["PassBlock"][up.Position]["mean"], attributeMeans["PassBlock"][up.Position]["stddev"], up.Year)
	RunBlockGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(up.RunBlock)-rangeNum, int(up.RunBlock)+rangeNum), attributeMeans["RunBlock"][up.Position]["mean"], attributeMeans["RunBlock"][up.Position]["stddev"], up.Year)
	PassRushGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(up.PassRush)-rangeNum, int(up.PassRush)+rangeNum), attributeMeans["PassRush"][up.Position]["mean"], attributeMeans["PassRush"][up.Position]["stddev"], up.Year)
	RunDefenseGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(up.RunDefense)-rangeNum, int(up.RunDefense)+rangeNum), attributeMeans["RunDefense"][up.Position]["mean"], attributeMeans["RunDefense"][up.Position]["stddev"], up.Year)
	ThrowPowerGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(up.ThrowPower)-rangeNum, int(up.ThrowPower)+rangeNum), attributeMeans["ThrowPower"][up.Position]["mean"], attributeMeans["ThrowPower"][up.Position]["stddev"], up.Year)
	ThrowAccuracyGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(up.ThrowAccuracy)-rangeNum, int(up.ThrowAccuracy)+rangeNum), attributeMeans["ThrowAccuracy"][up.Position]["mean"], attributeMeans["ThrowAccuracy"][up.Position]["stddev"], up.Year)
	KickPowerGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(up.KickPower)-rangeNum, int(up.KickPower)+rangeNum), attributeMeans["KickPower"][up.Position]["mean"], attributeMeans["KickPower"][up.Position]["stddev"], up.Year)
	KickAccuracyGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(up.KickAccuracy)-rangeNum, int(up.KickAccuracy)+rangeNum), attributeMeans["KickAccuracy"][up.Position]["mean"], attributeMeans["KickAccuracy"][up.Position]["stddev"], up.Year)
	PuntPowerGrade := util.GetLetterGrade(util.GenerateIntFromRange(int(up.PuntAccuracy)-rangeNum, int(up.PuntPower)+rangeNum), attributeMeans["PuntPower"][up.Position]["mean"], attributeMeans["PuntPower"][up.Position]["stddev"], up.Year)
	PuntAccuracyGrade := util.GetLetterGrade(int(up.PuntAccuracy), attributeMeans["PuntAccuracy"][up.Position]["mean"], attributeMeans["PuntAccuracy"][up.Position]["stddev"], up.Year)
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

func (n *NFLDraftee) AssignDraftedTeam(num, round uint, pickID uint, teamID uint, team string) {
	n.DraftedPick = num
	n.DraftedRound = round
	n.DraftPickID = pickID
	n.DraftedTeamID = teamID
	n.DraftedTeam = team
}

func (n *NFLDraftee) AssignBoomBustStatus(status string) {
	n.BoomOrBust = true
	n.BoomOrBustStatus = status
}

func (np *NFLDraftee) Progress(attr structs.CollegePlayerProgressions) {
	np.Agility = int8(attr.Agility)
	np.Speed = int8(attr.Speed)
	np.FootballIQ = int8(attr.FootballIQ)
	np.Carrying = int8(attr.Carrying)
	np.Catching = int8(attr.Catching)
	np.RouteRunning = int8(attr.RouteRunning)
	np.PassBlock = int8(attr.PassBlock)
	np.RunBlock = int8(attr.RunBlock)
	np.PassRush = int8(attr.PassRush)
	np.RunDefense = int8(attr.RunDefense)
	np.Tackle = int8(attr.Tackle)
	np.Strength = int8(attr.Strength)
	np.ManCoverage = int8(attr.ManCoverage)
	np.ZoneCoverage = int8(attr.ZoneCoverage)
	np.KickAccuracy = int8(attr.KickAccuracy)
	np.KickPower = int8(attr.KickPower)
	np.PuntAccuracy = int8(attr.PuntAccuracy)
	np.PuntPower = int8(attr.PuntPower)
	np.ThrowAccuracy = int8(attr.ThrowAccuracy)
	np.ThrowPower = int8(attr.ThrowPower)
}
