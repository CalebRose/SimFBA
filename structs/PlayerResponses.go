package structs

import (
	"sort"
	"strconv"

	config "github.com/CalebRose/SimFBA/secrets"
	"github.com/CalebRose/SimFBA/util"
)

type CollegePlayerResponse struct {
	ID int
	BasePlayer
	TeamID       int
	TeamAbbr     string
	City         string
	State        string
	Year         int
	IsRedshirt   bool
	ConferenceID int
	Conference   string
	Stats        CollegePlayerStats
	SeasonStats  CollegePlayerSeasonStats
}

type NFLPlayerResponse struct {
	ID int
	BasePlayer
	TeamID       int
	TeamAbbr     string
	City         string
	State        string
	Year         int
	ConferenceID int
	Conference   string
	DivisionID   int
	Division     string
	Stats        NFLPlayerStats
	SeasonStats  NFLPlayerSeasonStats
}

type DiscordPlayerResponse struct {
	Player       CollegePlayerCSV
	CollegeStats CollegePlayerSeasonStats
	NFLStats     NFLPlayerSeasonStats
}

type CollegePlayerCSV struct {
	FirstName          string
	LastName           string
	Position           string
	Archetype          string
	PositionTwo        string
	ArchetypeTwo       string
	Year               string
	Team               string
	PreviousTeam       string
	Age                int
	Stars              int
	HighSchool         string
	City               string
	State              string
	College            string
	Height             int
	Weight             int
	Shotgun            int
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
	PotentialGrade     string
	RedshirtStatus     string
	Stats              []CollegePlayerStats
}

func MapPlayerForStats(player CollegePlayer) CollegePlayerCSV {
	attributeMeans := config.AttributeMeans()
	Year, RedShirtStatus := util.GetYearAndRedshirtStatus(player.Year, player.IsRedshirt)
	OverallGrade := util.GetOverallGrade(player.Overall, player.Year)
	StaminaGrade := util.GetLetterGrade(player.Stamina, attributeMeans["Stamina"][player.Position]["mean"], attributeMeans["Stamina"][player.Position]["stddev"], player.Year)
	InjuryGrade := util.GetLetterGrade(player.Injury, attributeMeans["Injury"][player.Position]["mean"], attributeMeans["Injury"][player.Position]["stddev"], player.Year)
	SpeedGrade := util.GetLetterGrade(player.Speed, attributeMeans["Speed"][player.Position]["mean"], attributeMeans["Speed"][player.Position]["stddev"], player.Year)
	FootballIQGrade := util.GetLetterGrade(player.FootballIQ, attributeMeans["FootballIQ"][player.Position]["mean"], attributeMeans["FootballIQ"][player.Position]["stddev"], player.Year)
	AgilityGrade := util.GetLetterGrade(player.Agility, attributeMeans["Agility"][player.Position]["mean"], attributeMeans["Agility"][player.Position]["stddev"], player.Year)
	CarryingGrade := util.GetLetterGrade(player.Carrying, attributeMeans["Carrying"][player.Position]["mean"], attributeMeans["Carrying"][player.Position]["stddev"], player.Year)
	CatchingGrade := util.GetLetterGrade(player.Catching, attributeMeans["Catching"][player.Position]["mean"], attributeMeans["Catching"][player.Position]["stddev"], player.Year)
	RouteRunningGrade := util.GetLetterGrade(player.RouteRunning, attributeMeans["RouteRunning"][player.Position]["mean"], attributeMeans["RouteRunning"][player.Position]["stddev"], player.Year)
	ZoneCoverageGrade := util.GetLetterGrade(player.ZoneCoverage, attributeMeans["ZoneCoverage"][player.Position]["mean"], attributeMeans["ZoneCoverage"][player.Position]["stddev"], player.Year)
	ManCoverageGrade := util.GetLetterGrade(player.ManCoverage, attributeMeans["ManCoverage"][player.Position]["mean"], attributeMeans["ManCoverage"][player.Position]["stddev"], player.Year)
	StrengthGrade := util.GetLetterGrade(player.Strength, attributeMeans["Strength"][player.Position]["mean"], attributeMeans["Strength"][player.Position]["stddev"], player.Year)
	TackleGrade := util.GetLetterGrade(player.Tackle, attributeMeans["Tackle"][player.Position]["mean"], attributeMeans["Tackle"][player.Position]["stddev"], player.Year)
	PassBlockGrade := util.GetLetterGrade(player.PassBlock, attributeMeans["PassBlock"][player.Position]["mean"], attributeMeans["PassBlock"][player.Position]["stddev"], player.Year)
	RunBlockGrade := util.GetLetterGrade(player.RunBlock, attributeMeans["RunBlock"][player.Position]["mean"], attributeMeans["RunBlock"][player.Position]["stddev"], player.Year)
	PassRushGrade := util.GetLetterGrade(player.PassRush, attributeMeans["PassRush"][player.Position]["mean"], attributeMeans["PassRush"][player.Position]["stddev"], player.Year)
	RunDefenseGrade := util.GetLetterGrade(player.RunDefense, attributeMeans["RunDefense"][player.Position]["mean"], attributeMeans["RunDefense"][player.Position]["stddev"], player.Year)
	ThrowPowerGrade := util.GetLetterGrade(player.ThrowPower, attributeMeans["ThrowPower"][player.Position]["mean"], attributeMeans["ThrowPower"][player.Position]["stddev"], player.Year)
	ThrowAccuracyGrade := util.GetLetterGrade(player.ThrowAccuracy, attributeMeans["ThrowAccuracy"][player.Position]["mean"], attributeMeans["ThrowAccuracy"][player.Position]["stddev"], player.Year)
	KickPowerGrade := util.GetLetterGrade(player.KickPower, attributeMeans["KickPower"][player.Position]["mean"], attributeMeans["KickPower"][player.Position]["stddev"], player.Year)
	KickAccuracyGrade := util.GetLetterGrade(player.KickAccuracy, attributeMeans["KickAccuracy"][player.Position]["mean"], attributeMeans["KickAccuracy"][player.Position]["stddev"], player.Year)
	PuntPowerGrade := util.GetLetterGrade(player.PuntPower, attributeMeans["PuntPower"][player.Position]["mean"], attributeMeans["PuntPower"][player.Position]["stddev"], player.Year)
	PuntAccuracyGrade := util.GetLetterGrade(player.PuntAccuracy, attributeMeans["PuntAccuracy"][player.Position]["mean"], attributeMeans["PuntAccuracy"][player.Position]["stddev"], player.Year)

	return CollegePlayerCSV{
		FirstName:          player.FirstName,
		LastName:           player.LastName,
		Position:           player.Position,
		Archetype:          player.Archetype,
		PositionTwo:        player.PositionTwo,
		ArchetypeTwo:       player.ArchetypeTwo,
		Year:               Year,
		Age:                player.Age,
		Stars:              player.Stars,
		HighSchool:         player.HighSchool,
		City:               player.City,
		State:              player.State,
		Height:             player.Height,
		Weight:             player.Weight,
		OverallGrade:       OverallGrade,
		StaminaGrade:       StaminaGrade,
		InjuryGrade:        InjuryGrade,
		FootballIQGrade:    FootballIQGrade,
		SpeedGrade:         SpeedGrade,
		CarryingGrade:      CarryingGrade,
		AgilityGrade:       AgilityGrade,
		CatchingGrade:      CatchingGrade,
		RouteRunningGrade:  RouteRunningGrade,
		ZoneCoverageGrade:  ZoneCoverageGrade,
		ManCoverageGrade:   ManCoverageGrade,
		StrengthGrade:      StrengthGrade,
		TackleGrade:        TackleGrade,
		PassBlockGrade:     PassBlockGrade,
		RunBlockGrade:      RunBlockGrade,
		PassRushGrade:      PassRushGrade,
		RunDefenseGrade:    RunDefenseGrade,
		ThrowPowerGrade:    ThrowPowerGrade,
		ThrowAccuracyGrade: ThrowAccuracyGrade,
		KickAccuracyGrade:  KickAccuracyGrade,
		KickPowerGrade:     KickPowerGrade,
		PuntAccuracyGrade:  PuntAccuracyGrade,
		PuntPowerGrade:     PuntPowerGrade,
		PotentialGrade:     player.PotentialGrade,
		RedshirtStatus:     RedShirtStatus,
		Stats:              player.Stats,
	}
}

func MapPlayerToCSVModel(player CollegePlayer) CollegePlayerCSV {

	attributeMeans := config.AttributeMeans()
	Year, RedShirtStatus := util.GetYearAndRedshirtStatus(player.Year, player.IsRedshirt)
	OverallGrade := util.GetOverallGrade(player.Overall, player.Year)
	StaminaGrade := ""
	InjuryGrade := ""
	SpeedGrade := ""
	FootballIQGrade := ""
	AgilityGrade := ""
	CarryingGrade := ""
	CatchingGrade := ""
	RouteRunningGrade := ""
	ZoneCoverageGrade := ""
	ManCoverageGrade := ""
	StrengthGrade := ""
	TackleGrade := ""
	PassBlockGrade := ""
	RunBlockGrade := ""
	PassRushGrade := ""
	RunDefenseGrade := ""
	ThrowPowerGrade := ""
	ThrowAccuracyGrade := ""
	KickPowerGrade := ""
	KickAccuracyGrade := ""
	PuntPowerGrade := ""
	PuntAccuracyGrade := ""

	if player.Position == "ATH" {
		if player.Archetype == "Field General" {
			StaminaGrade = util.GetLetterGrade(player.Stamina, attributeMeans["Stamina"]["QB"]["mean"], attributeMeans["Stamina"][player.Position]["stddev"], player.Year)
			InjuryGrade = util.GetLetterGrade(player.Injury, attributeMeans["Injury"]["QB"]["mean"], attributeMeans["Injury"][player.Position]["stddev"], player.Year)
			SpeedGrade = util.GetLetterGrade(player.Speed, attributeMeans["Speed"]["ILB"]["mean"], attributeMeans["Speed"][player.Position]["stddev"], player.Year)
			FootballIQGrade = util.GetLetterGrade(player.FootballIQ, attributeMeans["FootballIQ"]["ILB"]["mean"], attributeMeans["FootballIQ"][player.Position]["stddev"], player.Year)
			AgilityGrade = util.GetLetterGrade(player.Agility, attributeMeans["Agility"]["QB"]["mean"], attributeMeans["Agility"][player.Position]["stddev"], player.Year)
			CarryingGrade = util.GetLetterGrade(player.Carrying, attributeMeans["Carrying"]["QB"]["mean"], attributeMeans["Carrying"][player.Position]["stddev"], player.Year)
			CatchingGrade = util.GetLetterGrade(player.Catching, attributeMeans["Catching"]["QB"]["mean"], attributeMeans["Catching"][player.Position]["stddev"], player.Year)
			RouteRunningGrade = util.GetLetterGrade(player.RouteRunning, attributeMeans["RouteRunning"]["QB"]["mean"], attributeMeans["RouteRunning"][player.Position]["stddev"], player.Year)
			ZoneCoverageGrade = util.GetLetterGrade(player.ZoneCoverage, attributeMeans["ZoneCoverage"]["ILB"]["mean"], attributeMeans["ZoneCoverage"][player.Position]["stddev"], player.Year)
			ManCoverageGrade = util.GetLetterGrade(player.ManCoverage, attributeMeans["ManCoverage"]["ILB"]["mean"], attributeMeans["ManCoverage"][player.Position]["stddev"], player.Year)
			StrengthGrade = util.GetLetterGrade(player.Strength, attributeMeans["Strength"]["ILB"]["mean"], attributeMeans["Strength"][player.Position]["stddev"], player.Year)
			TackleGrade = util.GetLetterGrade(player.Tackle, attributeMeans["Tackle"]["ILB"]["mean"], attributeMeans["Tackle"][player.Position]["stddev"], player.Year)
			PassBlockGrade = util.GetLetterGrade(player.PassBlock, attributeMeans["PassBlock"]["QB"]["mean"], attributeMeans["PassBlock"][player.Position]["stddev"], player.Year)
			RunBlockGrade = util.GetLetterGrade(player.RunBlock, attributeMeans["RunBlock"]["QB"]["mean"], attributeMeans["RunBlock"][player.Position]["stddev"], player.Year)
			PassRushGrade = util.GetLetterGrade(player.PassRush, attributeMeans["PassRush"]["ILB"]["mean"], attributeMeans["PassRush"][player.Position]["stddev"], player.Year)
			RunDefenseGrade = util.GetLetterGrade(player.RunDefense, attributeMeans["RunDefense"]["ILB"]["mean"], attributeMeans["RunDefense"][player.Position]["stddev"], player.Year)
			ThrowPowerGrade = util.GetLetterGrade(player.ThrowPower, attributeMeans["ThrowPower"]["QB"]["mean"], attributeMeans["ThrowPower"][player.Position]["stddev"], player.Year)
			ThrowAccuracyGrade = util.GetLetterGrade(player.ThrowAccuracy, attributeMeans["ThrowAccuracy"]["QB"]["mean"], attributeMeans["ThrowAccuracy"][player.Position]["stddev"], player.Year)
			KickPowerGrade = util.GetLetterGrade(player.KickPower, attributeMeans["KickPower"]["QB"]["mean"], attributeMeans["KickPower"][player.Position]["stddev"], player.Year)
			KickAccuracyGrade = util.GetLetterGrade(player.KickAccuracy, attributeMeans["KickAccuracy"]["QB"]["mean"], attributeMeans["KickAccuracy"][player.Position]["stddev"], player.Year)
			PuntPowerGrade = util.GetLetterGrade(player.PuntPower, attributeMeans["PuntPower"]["QB"]["mean"], attributeMeans["PuntPower"][player.Position]["stddev"], player.Year)
			PuntAccuracyGrade = util.GetLetterGrade(player.PuntAccuracy, attributeMeans["PuntAccuracy"]["QB"]["mean"], attributeMeans["PuntAccuracy"][player.Position]["stddev"], player.Year)
		} else if player.Archetype == "Triple-Threat" {
			StaminaGrade = util.GetLetterGrade(player.Stamina, attributeMeans["Stamina"]["QB"]["mean"], attributeMeans["Stamina"][player.Position]["stddev"], player.Year)
			InjuryGrade = util.GetLetterGrade(player.Injury, attributeMeans["Injury"]["QB"]["mean"], attributeMeans["Injury"][player.Position]["stddev"], player.Year)
			SpeedGrade = util.GetLetterGrade(player.Speed, attributeMeans["Speed"]["RB"]["mean"], attributeMeans["Speed"][player.Position]["stddev"], player.Year)
			FootballIQGrade = util.GetLetterGrade(player.FootballIQ, attributeMeans["FootballIQ"]["QB"]["mean"], attributeMeans["FootballIQ"][player.Position]["stddev"], player.Year)
			AgilityGrade = util.GetLetterGrade(player.Agility, attributeMeans["Agility"]["RB"]["mean"], attributeMeans["Agility"][player.Position]["stddev"], player.Year)
			CarryingGrade = util.GetLetterGrade(player.Carrying, attributeMeans["Carrying"]["RB"]["mean"], attributeMeans["Carrying"][player.Position]["stddev"], player.Year)
			CatchingGrade = util.GetLetterGrade(player.Catching, attributeMeans["Catching"]["WR"]["mean"], attributeMeans["Catching"][player.Position]["stddev"], player.Year)
			RouteRunningGrade = util.GetLetterGrade(player.RouteRunning, attributeMeans["RouteRunning"]["WR"]["mean"], attributeMeans["RouteRunning"][player.Position]["stddev"], player.Year)
			ZoneCoverageGrade = util.GetLetterGrade(player.ZoneCoverage, attributeMeans["ZoneCoverage"]["CB"]["mean"], attributeMeans["ZoneCoverage"][player.Position]["stddev"], player.Year)
			ManCoverageGrade = util.GetLetterGrade(player.ManCoverage, attributeMeans["ManCoverage"]["CB"]["mean"], attributeMeans["ManCoverage"][player.Position]["stddev"], player.Year)
			StrengthGrade = util.GetLetterGrade(player.Strength, attributeMeans["Strength"]["RB"]["mean"], attributeMeans["Strength"][player.Position]["stddev"], player.Year)
			TackleGrade = util.GetLetterGrade(player.Tackle, attributeMeans["Tackle"]["WR"]["mean"], attributeMeans["Tackle"][player.Position]["stddev"], player.Year)
			PassBlockGrade = util.GetLetterGrade(player.PassBlock, attributeMeans["PassBlock"]["RB"]["mean"], attributeMeans["PassBlock"][player.Position]["stddev"], player.Year)
			RunBlockGrade = util.GetLetterGrade(player.RunBlock, attributeMeans["RunBlock"]["RB"]["mean"], attributeMeans["RunBlock"][player.Position]["stddev"], player.Year)
			PassRushGrade = util.GetLetterGrade(player.PassRush, attributeMeans["PassRush"]["RB"]["mean"], attributeMeans["PassRush"][player.Position]["stddev"], player.Year)
			RunDefenseGrade = util.GetLetterGrade(player.RunDefense, attributeMeans["RunDefense"]["RB"]["mean"], attributeMeans["RunDefense"][player.Position]["stddev"], player.Year)
			ThrowPowerGrade = util.GetLetterGrade(player.ThrowPower, attributeMeans["ThrowPower"]["QB"]["mean"], attributeMeans["ThrowPower"][player.Position]["stddev"], player.Year)
			ThrowAccuracyGrade = util.GetLetterGrade(player.ThrowAccuracy, attributeMeans["ThrowAccuracy"]["QB"]["mean"], attributeMeans["ThrowAccuracy"][player.Position]["stddev"], player.Year)
			KickPowerGrade = util.GetLetterGrade(player.KickPower, attributeMeans["KickPower"]["QB"]["mean"], attributeMeans["KickPower"][player.Position]["stddev"], player.Year)
			KickAccuracyGrade = util.GetLetterGrade(player.KickAccuracy, attributeMeans["KickAccuracy"]["QB"]["mean"], attributeMeans["KickAccuracy"][player.Position]["stddev"], player.Year)
			PuntPowerGrade = util.GetLetterGrade(player.PuntPower, attributeMeans["PuntPower"]["QB"]["mean"], attributeMeans["PuntPower"][player.Position]["stddev"], player.Year)
			PuntAccuracyGrade = util.GetLetterGrade(player.PuntAccuracy, attributeMeans["PuntAccuracy"]["QB"]["mean"], attributeMeans["PuntAccuracy"][player.Position]["stddev"], player.Year)
		} else if player.Archetype == "Wingback" {
			StaminaGrade = util.GetLetterGrade(player.Stamina, attributeMeans["Stamina"]["QB"]["mean"], attributeMeans["Stamina"][player.Position]["stddev"], player.Year)
			InjuryGrade = util.GetLetterGrade(player.Injury, attributeMeans["Injury"]["QB"]["mean"], attributeMeans["Injury"][player.Position]["stddev"], player.Year)
			SpeedGrade = util.GetLetterGrade(player.Speed, attributeMeans["Speed"]["RB"]["mean"], attributeMeans["Speed"][player.Position]["stddev"], player.Year)
			FootballIQGrade = util.GetLetterGrade(player.FootballIQ, attributeMeans["FootballIQ"]["RB"]["mean"], attributeMeans["FootballIQ"][player.Position]["stddev"], player.Year)
			AgilityGrade = util.GetLetterGrade(player.Agility, attributeMeans["Agility"]["RB"]["mean"], attributeMeans["Agility"][player.Position]["stddev"], player.Year)
			CarryingGrade = util.GetLetterGrade(player.Carrying, attributeMeans["Carrying"]["RB"]["mean"], attributeMeans["Carrying"][player.Position]["stddev"], player.Year)
			CatchingGrade = util.GetLetterGrade(player.Catching, attributeMeans["Catching"]["WR"]["mean"], attributeMeans["Catching"][player.Position]["stddev"], player.Year)
			RouteRunningGrade = util.GetLetterGrade(player.RouteRunning, attributeMeans["RouteRunning"]["WR"]["mean"], attributeMeans["RouteRunning"][player.Position]["stddev"], player.Year)
			ZoneCoverageGrade = util.GetLetterGrade(player.ZoneCoverage, attributeMeans["ZoneCoverage"]["CB"]["mean"], attributeMeans["ZoneCoverage"][player.Position]["stddev"], player.Year)
			ManCoverageGrade = util.GetLetterGrade(player.ManCoverage, attributeMeans["ManCoverage"]["CB"]["mean"], attributeMeans["ManCoverage"][player.Position]["stddev"], player.Year)
			StrengthGrade = util.GetLetterGrade(player.Strength, attributeMeans["Strength"]["RB"]["mean"], attributeMeans["Strength"][player.Position]["stddev"], player.Year)
			TackleGrade = util.GetLetterGrade(player.Tackle, attributeMeans["Tackle"]["WR"]["mean"], attributeMeans["Tackle"][player.Position]["stddev"], player.Year)
			PassBlockGrade = util.GetLetterGrade(player.PassBlock, attributeMeans["PassBlock"]["RB"]["mean"], attributeMeans["PassBlock"][player.Position]["stddev"], player.Year)
			RunBlockGrade = util.GetLetterGrade(player.RunBlock, attributeMeans["RunBlock"]["RB"]["mean"], attributeMeans["RunBlock"][player.Position]["stddev"], player.Year)
			PassRushGrade = util.GetLetterGrade(player.PassRush, attributeMeans["PassRush"]["CB"]["mean"], attributeMeans["PassRush"][player.Position]["stddev"], player.Year)
			RunDefenseGrade = util.GetLetterGrade(player.RunDefense, attributeMeans["RunDefense"]["CB"]["mean"], attributeMeans["RunDefense"][player.Position]["stddev"], player.Year)
			ThrowPowerGrade = util.GetLetterGrade(player.ThrowPower, attributeMeans["ThrowPower"]["QB"]["mean"], attributeMeans["ThrowPower"][player.Position]["stddev"], player.Year)
			ThrowAccuracyGrade = util.GetLetterGrade(player.ThrowAccuracy, attributeMeans["ThrowAccuracy"]["QB"]["mean"], attributeMeans["ThrowAccuracy"][player.Position]["stddev"], player.Year)
			KickPowerGrade = util.GetLetterGrade(player.KickPower, attributeMeans["KickPower"]["QB"]["mean"], attributeMeans["KickPower"][player.Position]["stddev"], player.Year)
			KickAccuracyGrade = util.GetLetterGrade(player.KickAccuracy, attributeMeans["KickAccuracy"]["QB"]["mean"], attributeMeans["KickAccuracy"][player.Position]["stddev"], player.Year)
			PuntPowerGrade = util.GetLetterGrade(player.PuntPower, attributeMeans["PuntPower"]["QB"]["mean"], attributeMeans["PuntPower"][player.Position]["stddev"], player.Year)
			PuntAccuracyGrade = util.GetLetterGrade(player.PuntAccuracy, attributeMeans["PuntAccuracy"]["QB"]["mean"], attributeMeans["PuntAccuracy"][player.Position]["stddev"], player.Year)
		} else if player.Archetype == "Slotback" {
			StaminaGrade = util.GetLetterGrade(player.Stamina, attributeMeans["Stamina"]["QB"]["mean"], attributeMeans["Stamina"][player.Position]["stddev"], player.Year)
			InjuryGrade = util.GetLetterGrade(player.Injury, attributeMeans["Injury"]["QB"]["mean"], attributeMeans["Injury"][player.Position]["stddev"], player.Year)
			SpeedGrade = util.GetLetterGrade(player.Speed, attributeMeans["Speed"]["WR"]["mean"], attributeMeans["Speed"][player.Position]["stddev"], player.Year)
			FootballIQGrade = util.GetLetterGrade(player.FootballIQ, attributeMeans["FootballIQ"]["WR"]["mean"], attributeMeans["FootballIQ"][player.Position]["stddev"], player.Year)
			AgilityGrade = util.GetLetterGrade(player.Agility, attributeMeans["Agility"]["WR"]["mean"], attributeMeans["Agility"][player.Position]["stddev"], player.Year)
			CarryingGrade = util.GetLetterGrade(player.Carrying, attributeMeans["Carrying"]["RB"]["mean"], attributeMeans["Carrying"][player.Position]["stddev"], player.Year)
			CatchingGrade = util.GetLetterGrade(player.Catching, attributeMeans["Catching"]["WR"]["mean"], attributeMeans["Catching"][player.Position]["stddev"], player.Year)
			RouteRunningGrade = util.GetLetterGrade(player.RouteRunning, attributeMeans["RouteRunning"]["WR"]["mean"], attributeMeans["RouteRunning"][player.Position]["stddev"], player.Year)
			ZoneCoverageGrade = util.GetLetterGrade(player.ZoneCoverage, attributeMeans["ZoneCoverage"]["CB"]["mean"], attributeMeans["ZoneCoverage"][player.Position]["stddev"], player.Year)
			ManCoverageGrade = util.GetLetterGrade(player.ManCoverage, attributeMeans["ManCoverage"]["CB"]["mean"], attributeMeans["ManCoverage"][player.Position]["stddev"], player.Year)
			StrengthGrade = util.GetLetterGrade(player.Strength, attributeMeans["Strength"]["TE"]["mean"], attributeMeans["Strength"][player.Position]["stddev"], player.Year)
			TackleGrade = util.GetLetterGrade(player.Tackle, attributeMeans["Tackle"]["WR"]["mean"], attributeMeans["Tackle"][player.Position]["stddev"], player.Year)
			PassBlockGrade = util.GetLetterGrade(player.PassBlock, attributeMeans["PassBlock"]["FB"]["mean"], attributeMeans["PassBlock"][player.Position]["stddev"], player.Year)
			RunBlockGrade = util.GetLetterGrade(player.RunBlock, attributeMeans["RunBlock"]["FB"]["mean"], attributeMeans["RunBlock"][player.Position]["stddev"], player.Year)
			PassRushGrade = util.GetLetterGrade(player.PassRush, attributeMeans["PassRush"]["TE"]["mean"], attributeMeans["PassRush"][player.Position]["stddev"], player.Year)
			RunDefenseGrade = util.GetLetterGrade(player.RunDefense, attributeMeans["RunDefense"]["TE"]["mean"], attributeMeans["RunDefense"][player.Position]["stddev"], player.Year)
			ThrowPowerGrade = util.GetLetterGrade(player.ThrowPower, attributeMeans["ThrowPower"]["QB"]["mean"], attributeMeans["ThrowPower"][player.Position]["stddev"], player.Year)
			ThrowAccuracyGrade = util.GetLetterGrade(player.ThrowAccuracy, attributeMeans["ThrowAccuracy"]["QB"]["mean"], attributeMeans["ThrowAccuracy"][player.Position]["stddev"], player.Year)
			KickPowerGrade = util.GetLetterGrade(player.KickPower, attributeMeans["KickPower"]["QB"]["mean"], attributeMeans["KickPower"][player.Position]["stddev"], player.Year)
			KickAccuracyGrade = util.GetLetterGrade(player.KickAccuracy, attributeMeans["KickAccuracy"]["QB"]["mean"], attributeMeans["KickAccuracy"][player.Position]["stddev"], player.Year)
			PuntPowerGrade = util.GetLetterGrade(player.PuntPower, attributeMeans["PuntPower"]["QB"]["mean"], attributeMeans["PuntPower"][player.Position]["stddev"], player.Year)
			PuntAccuracyGrade = util.GetLetterGrade(player.PuntAccuracy, attributeMeans["PuntAccuracy"]["QB"]["mean"], attributeMeans["PuntAccuracy"][player.Position]["stddev"], player.Year)
		} else if player.Archetype == "Lineman" {
			StaminaGrade = util.GetLetterGrade(player.Stamina, attributeMeans["Stamina"]["OG"]["mean"], attributeMeans["Stamina"][player.Position]["stddev"], player.Year)
			InjuryGrade = util.GetLetterGrade(player.Injury, attributeMeans["Injury"]["OG"]["mean"], attributeMeans["Injury"][player.Position]["stddev"], player.Year)
			SpeedGrade = util.GetLetterGrade(player.Speed, attributeMeans["Speed"]["DE"]["mean"], attributeMeans["Speed"][player.Position]["stddev"], player.Year)
			FootballIQGrade = util.GetLetterGrade(player.FootballIQ, attributeMeans["FootballIQ"]["C"]["mean"], attributeMeans["FootballIQ"][player.Position]["stddev"], player.Year)
			AgilityGrade = util.GetLetterGrade(player.Agility, attributeMeans["Agility"]["DE"]["mean"], attributeMeans["Agility"][player.Position]["stddev"], player.Year)
			CarryingGrade = util.GetLetterGrade(player.Carrying, attributeMeans["Carrying"]["RB"]["mean"], attributeMeans["Carrying"][player.Position]["stddev"], player.Year)
			CatchingGrade = util.GetLetterGrade(player.Catching, attributeMeans["Catching"]["RB"]["mean"], attributeMeans["Catching"][player.Position]["stddev"], player.Year)
			RouteRunningGrade = util.GetLetterGrade(player.RouteRunning, attributeMeans["RouteRunning"]["OG"]["mean"], attributeMeans["RouteRunning"][player.Position]["stddev"], player.Year)
			ZoneCoverageGrade = util.GetLetterGrade(player.ZoneCoverage, attributeMeans["ZoneCoverage"]["CB"]["mean"], attributeMeans["ZoneCoverage"][player.Position]["stddev"], player.Year)
			ManCoverageGrade = util.GetLetterGrade(player.ManCoverage, attributeMeans["ManCoverage"]["CB"]["mean"], attributeMeans["ManCoverage"][player.Position]["stddev"], player.Year)
			StrengthGrade = util.GetLetterGrade(player.Strength, attributeMeans["Strength"]["DT"]["mean"], attributeMeans["Strength"][player.Position]["stddev"], player.Year)
			TackleGrade = util.GetLetterGrade(player.Tackle, attributeMeans["Tackle"]["DT"]["mean"], attributeMeans["Tackle"][player.Position]["stddev"], player.Year)
			PassBlockGrade = util.GetLetterGrade(player.PassBlock, attributeMeans["PassBlock"]["OG"]["mean"], attributeMeans["PassBlock"][player.Position]["stddev"], player.Year)
			RunBlockGrade = util.GetLetterGrade(player.RunBlock, attributeMeans["RunBlock"]["OG"]["mean"], attributeMeans["RunBlock"][player.Position]["stddev"], player.Year)
			PassRushGrade = util.GetLetterGrade(player.PassRush, attributeMeans["PassRush"]["DE"]["mean"], attributeMeans["PassRush"][player.Position]["stddev"], player.Year)
			RunDefenseGrade = util.GetLetterGrade(player.RunDefense, attributeMeans["RunDefense"]["DE"]["mean"], attributeMeans["RunDefense"][player.Position]["stddev"], player.Year)
			ThrowPowerGrade = util.GetLetterGrade(player.ThrowPower, attributeMeans["ThrowPower"]["QB"]["mean"], attributeMeans["ThrowPower"][player.Position]["stddev"], player.Year)
			ThrowAccuracyGrade = util.GetLetterGrade(player.ThrowAccuracy, attributeMeans["ThrowAccuracy"]["QB"]["mean"], attributeMeans["ThrowAccuracy"][player.Position]["stddev"], player.Year)
			KickPowerGrade = util.GetLetterGrade(player.KickPower, attributeMeans["KickPower"]["K"]["mean"], attributeMeans["KickPower"][player.Position]["stddev"], player.Year)
			KickAccuracyGrade = util.GetLetterGrade(player.KickAccuracy, attributeMeans["KickAccuracy"]["K"]["mean"], attributeMeans["KickAccuracy"][player.Position]["stddev"], player.Year)
			PuntPowerGrade = util.GetLetterGrade(player.PuntPower, attributeMeans["PuntPower"]["P"]["mean"], attributeMeans["PuntPower"][player.Position]["stddev"], player.Year)
			PuntAccuracyGrade = util.GetLetterGrade(player.PuntAccuracy, attributeMeans["PuntAccuracy"]["P"]["mean"], attributeMeans["PuntAccuracy"][player.Position]["stddev"], player.Year)
		} else if player.Archetype == "Strongside" {
			StaminaGrade = util.GetLetterGrade(player.Stamina, attributeMeans["Stamina"]["OG"]["mean"], attributeMeans["Stamina"][player.Position]["stddev"], player.Year)
			InjuryGrade = util.GetLetterGrade(player.Injury, attributeMeans["Injury"]["OG"]["mean"], attributeMeans["Injury"][player.Position]["stddev"], player.Year)
			SpeedGrade = util.GetLetterGrade(player.Speed, attributeMeans["Speed"]["DE"]["mean"], attributeMeans["Speed"][player.Position]["stddev"], player.Year)
			FootballIQGrade = util.GetLetterGrade(player.FootballIQ, attributeMeans["FootballIQ"]["OLB"]["mean"], attributeMeans["FootballIQ"][player.Position]["stddev"], player.Year)
			AgilityGrade = util.GetLetterGrade(player.Agility, attributeMeans["Agility"]["DE"]["mean"], attributeMeans["Agility"][player.Position]["stddev"], player.Year)
			CarryingGrade = util.GetLetterGrade(player.Carrying, attributeMeans["Carrying"]["OLB"]["mean"], attributeMeans["Carrying"][player.Position]["stddev"], player.Year)
			CatchingGrade = util.GetLetterGrade(player.Catching, attributeMeans["Catching"]["OLB"]["mean"], attributeMeans["Catching"][player.Position]["stddev"], player.Year)
			RouteRunningGrade = util.GetLetterGrade(player.RouteRunning, attributeMeans["RouteRunning"]["OG"]["mean"], attributeMeans["RouteRunning"][player.Position]["stddev"], player.Year)
			ZoneCoverageGrade = util.GetLetterGrade(player.ZoneCoverage, attributeMeans["ZoneCoverage"]["OLB"]["mean"], attributeMeans["ZoneCoverage"][player.Position]["stddev"], player.Year)
			ManCoverageGrade = util.GetLetterGrade(player.ManCoverage, attributeMeans["ManCoverage"]["OLB"]["mean"], attributeMeans["ManCoverage"][player.Position]["stddev"], player.Year)
			StrengthGrade = util.GetLetterGrade(player.Strength, attributeMeans["Strength"]["DE"]["mean"], attributeMeans["Strength"][player.Position]["stddev"], player.Year)
			TackleGrade = util.GetLetterGrade(player.Tackle, attributeMeans["Tackle"]["DE"]["mean"], attributeMeans["Tackle"][player.Position]["stddev"], player.Year)
			PassBlockGrade = util.GetLetterGrade(player.PassBlock, attributeMeans["PassBlock"]["OG"]["mean"], attributeMeans["PassBlock"][player.Position]["stddev"], player.Year)
			RunBlockGrade = util.GetLetterGrade(player.RunBlock, attributeMeans["RunBlock"]["OG"]["mean"], attributeMeans["RunBlock"][player.Position]["stddev"], player.Year)
			PassRushGrade = util.GetLetterGrade(player.PassRush, attributeMeans["PassRush"]["DE"]["mean"], attributeMeans["PassRush"][player.Position]["stddev"], player.Year)
			RunDefenseGrade = util.GetLetterGrade(player.RunDefense, attributeMeans["RunDefense"]["DE"]["mean"], attributeMeans["RunDefense"][player.Position]["stddev"], player.Year)
			ThrowPowerGrade = util.GetLetterGrade(player.ThrowPower, attributeMeans["ThrowPower"]["QB"]["mean"], attributeMeans["ThrowPower"][player.Position]["stddev"], player.Year)
			ThrowAccuracyGrade = util.GetLetterGrade(player.ThrowAccuracy, attributeMeans["ThrowAccuracy"]["QB"]["mean"], attributeMeans["ThrowAccuracy"][player.Position]["stddev"], player.Year)
			KickPowerGrade = util.GetLetterGrade(player.KickPower, attributeMeans["KickPower"]["K"]["mean"], attributeMeans["KickPower"][player.Position]["stddev"], player.Year)
			KickAccuracyGrade = util.GetLetterGrade(player.KickAccuracy, attributeMeans["KickAccuracy"]["K"]["mean"], attributeMeans["KickAccuracy"][player.Position]["stddev"], player.Year)
			PuntPowerGrade = util.GetLetterGrade(player.PuntPower, attributeMeans["PuntPower"]["P"]["mean"], attributeMeans["PuntPower"][player.Position]["stddev"], player.Year)
			PuntAccuracyGrade = util.GetLetterGrade(player.PuntAccuracy, attributeMeans["PuntAccuracy"]["P"]["mean"], attributeMeans["PuntAccuracy"][player.Position]["stddev"], player.Year)
		} else if player.Archetype == "Weakside" {
			StaminaGrade = util.GetLetterGrade(player.Stamina, attributeMeans["Stamina"]["OG"]["mean"], attributeMeans["Stamina"][player.Position]["stddev"], player.Year)
			InjuryGrade = util.GetLetterGrade(player.Injury, attributeMeans["Injury"]["OG"]["mean"], attributeMeans["Injury"][player.Position]["stddev"], player.Year)
			SpeedGrade = util.GetLetterGrade(player.Speed, attributeMeans["Speed"]["ILB"]["mean"], attributeMeans["Speed"][player.Position]["stddev"], player.Year)
			FootballIQGrade = util.GetLetterGrade(player.FootballIQ, attributeMeans["FootballIQ"]["ILB"]["mean"], attributeMeans["FootballIQ"][player.Position]["stddev"], player.Year)
			AgilityGrade = util.GetLetterGrade(player.Agility, attributeMeans["Agility"]["ILB"]["mean"], attributeMeans["Agility"][player.Position]["stddev"], player.Year)
			CarryingGrade = util.GetLetterGrade(player.Carrying, attributeMeans["Carrying"]["OLB"]["mean"], attributeMeans["Carrying"][player.Position]["stddev"], player.Year)
			CatchingGrade = util.GetLetterGrade(player.Catching, attributeMeans["Catching"]["OLB"]["mean"], attributeMeans["Catching"][player.Position]["stddev"], player.Year)
			RouteRunningGrade = util.GetLetterGrade(player.RouteRunning, attributeMeans["RouteRunning"]["OG"]["mean"], attributeMeans["RouteRunning"][player.Position]["stddev"], player.Year)
			ZoneCoverageGrade = util.GetLetterGrade(player.ZoneCoverage, attributeMeans["ZoneCoverage"]["ILB"]["mean"], attributeMeans["ZoneCoverage"][player.Position]["stddev"], player.Year)
			ManCoverageGrade = util.GetLetterGrade(player.ManCoverage, attributeMeans["ManCoverage"]["ILB"]["mean"], attributeMeans["ManCoverage"][player.Position]["stddev"], player.Year)
			StrengthGrade = util.GetLetterGrade(player.Strength, attributeMeans["Strength"]["OLB"]["mean"], attributeMeans["Strength"][player.Position]["stddev"], player.Year)
			TackleGrade = util.GetLetterGrade(player.Tackle, attributeMeans["Tackle"]["OLB"]["mean"], attributeMeans["Tackle"][player.Position]["stddev"], player.Year)
			PassBlockGrade = util.GetLetterGrade(player.PassBlock, attributeMeans["PassBlock"]["OG"]["mean"], attributeMeans["PassBlock"][player.Position]["stddev"], player.Year)
			RunBlockGrade = util.GetLetterGrade(player.RunBlock, attributeMeans["RunBlock"]["OG"]["mean"], attributeMeans["RunBlock"][player.Position]["stddev"], player.Year)
			PassRushGrade = util.GetLetterGrade(player.PassRush, attributeMeans["PassRush"]["ILB"]["mean"], attributeMeans["PassRush"][player.Position]["stddev"], player.Year)
			RunDefenseGrade = util.GetLetterGrade(player.RunDefense, attributeMeans["RunDefense"]["ILB"]["mean"], attributeMeans["RunDefense"][player.Position]["stddev"], player.Year)
			ThrowPowerGrade = util.GetLetterGrade(player.ThrowPower, attributeMeans["ThrowPower"]["QB"]["mean"], attributeMeans["ThrowPower"][player.Position]["stddev"], player.Year)
			ThrowAccuracyGrade = util.GetLetterGrade(player.ThrowAccuracy, attributeMeans["ThrowAccuracy"]["QB"]["mean"], attributeMeans["ThrowAccuracy"][player.Position]["stddev"], player.Year)
			KickPowerGrade = util.GetLetterGrade(player.KickPower, attributeMeans["KickPower"]["K"]["mean"], attributeMeans["KickPower"][player.Position]["stddev"], player.Year)
			KickAccuracyGrade = util.GetLetterGrade(player.KickAccuracy, attributeMeans["KickAccuracy"]["K"]["mean"], attributeMeans["KickAccuracy"][player.Position]["stddev"], player.Year)
			PuntPowerGrade = util.GetLetterGrade(player.PuntPower, attributeMeans["PuntPower"]["P"]["mean"], attributeMeans["PuntPower"][player.Position]["stddev"], player.Year)
			PuntAccuracyGrade = util.GetLetterGrade(player.PuntAccuracy, attributeMeans["PuntAccuracy"]["P"]["mean"], attributeMeans["PuntAccuracy"][player.Position]["stddev"], player.Year)
		} else if player.Archetype == "Bandit" {
			StaminaGrade = util.GetLetterGrade(player.Stamina, attributeMeans["Stamina"]["OG"]["mean"], attributeMeans["Stamina"][player.Position]["stddev"], player.Year)
			InjuryGrade = util.GetLetterGrade(player.Injury, attributeMeans["Injury"]["OG"]["mean"], attributeMeans["Injury"][player.Position]["stddev"], player.Year)
			SpeedGrade = util.GetLetterGrade(player.Speed, attributeMeans["Speed"]["SS"]["mean"], attributeMeans["Speed"][player.Position]["stddev"], player.Year)
			FootballIQGrade = util.GetLetterGrade(player.FootballIQ, attributeMeans["FootballIQ"]["ILB"]["mean"], attributeMeans["FootballIQ"][player.Position]["stddev"], player.Year)
			AgilityGrade = util.GetLetterGrade(player.Agility, attributeMeans["Agility"]["SS"]["mean"], attributeMeans["Agility"][player.Position]["stddev"], player.Year)
			CarryingGrade = util.GetLetterGrade(player.Carrying, attributeMeans["Carrying"]["SS"]["mean"], attributeMeans["Carrying"][player.Position]["stddev"], player.Year)
			CatchingGrade = util.GetLetterGrade(player.Catching, attributeMeans["Catching"]["SS"]["mean"], attributeMeans["Catching"][player.Position]["stddev"], player.Year)
			RouteRunningGrade = util.GetLetterGrade(player.RouteRunning, attributeMeans["RouteRunning"]["SS"]["mean"], attributeMeans["RouteRunning"][player.Position]["stddev"], player.Year)
			ZoneCoverageGrade = util.GetLetterGrade(player.ZoneCoverage, attributeMeans["ZoneCoverage"]["SS"]["mean"], attributeMeans["ZoneCoverage"][player.Position]["stddev"], player.Year)
			ManCoverageGrade = util.GetLetterGrade(player.ManCoverage, attributeMeans["ManCoverage"]["SS"]["mean"], attributeMeans["ManCoverage"][player.Position]["stddev"], player.Year)
			StrengthGrade = util.GetLetterGrade(player.Strength, attributeMeans["Strength"]["OLB"]["mean"], attributeMeans["Strength"][player.Position]["stddev"], player.Year)
			TackleGrade = util.GetLetterGrade(player.Tackle, attributeMeans["Tackle"]["OLB"]["mean"], attributeMeans["Tackle"][player.Position]["stddev"], player.Year)
			PassBlockGrade = util.GetLetterGrade(player.PassBlock, attributeMeans["PassBlock"]["OG"]["mean"], attributeMeans["PassBlock"][player.Position]["stddev"], player.Year)
			RunBlockGrade = util.GetLetterGrade(player.RunBlock, attributeMeans["RunBlock"]["OG"]["mean"], attributeMeans["RunBlock"][player.Position]["stddev"], player.Year)
			PassRushGrade = util.GetLetterGrade(player.PassRush, attributeMeans["PassRush"]["OLB"]["mean"], attributeMeans["PassRush"][player.Position]["stddev"], player.Year)
			RunDefenseGrade = util.GetLetterGrade(player.RunDefense, attributeMeans["RunDefense"]["OLB"]["mean"], attributeMeans["RunDefense"][player.Position]["stddev"], player.Year)
			ThrowPowerGrade = util.GetLetterGrade(player.ThrowPower, attributeMeans["ThrowPower"]["QB"]["mean"], attributeMeans["ThrowPower"][player.Position]["stddev"], player.Year)
			ThrowAccuracyGrade = util.GetLetterGrade(player.ThrowAccuracy, attributeMeans["ThrowAccuracy"]["QB"]["mean"], attributeMeans["ThrowAccuracy"][player.Position]["stddev"], player.Year)
			KickPowerGrade = util.GetLetterGrade(player.KickPower, attributeMeans["KickPower"]["K"]["mean"], attributeMeans["KickPower"][player.Position]["stddev"], player.Year)
			KickAccuracyGrade = util.GetLetterGrade(player.KickAccuracy, attributeMeans["KickAccuracy"]["K"]["mean"], attributeMeans["KickAccuracy"][player.Position]["stddev"], player.Year)
			PuntPowerGrade = util.GetLetterGrade(player.PuntPower, attributeMeans["PuntPower"]["P"]["mean"], attributeMeans["PuntPower"][player.Position]["stddev"], player.Year)
			PuntAccuracyGrade = util.GetLetterGrade(player.PuntAccuracy, attributeMeans["PuntAccuracy"]["P"]["mean"], attributeMeans["PuntAccuracy"][player.Position]["stddev"], player.Year)
		} else if player.Archetype == "Return Specialist" {
			StaminaGrade = util.GetLetterGrade(player.Stamina, attributeMeans["Stamina"]["OG"]["mean"], attributeMeans["Stamina"][player.Position]["stddev"], player.Year)
			InjuryGrade = util.GetLetterGrade(player.Injury, attributeMeans["Injury"]["OG"]["mean"], attributeMeans["Injury"][player.Position]["stddev"], player.Year)
			SpeedGrade = util.GetLetterGrade(player.Speed, attributeMeans["Speed"]["RB"]["mean"], attributeMeans["Speed"][player.Position]["stddev"], player.Year)
			FootballIQGrade = util.GetLetterGrade(player.FootballIQ, attributeMeans["FootballIQ"]["ILB"]["mean"], attributeMeans["FootballIQ"][player.Position]["stddev"], player.Year)
			AgilityGrade = util.GetLetterGrade(player.Agility, attributeMeans["Agility"]["RB"]["mean"], attributeMeans["Agility"][player.Position]["stddev"], player.Year)
			CarryingGrade = util.GetLetterGrade(player.Carrying, attributeMeans["Carrying"]["RB"]["mean"], attributeMeans["Carrying"][player.Position]["stddev"], player.Year)
			CatchingGrade = util.GetLetterGrade(player.Catching, attributeMeans["Catching"]["WR"]["mean"], attributeMeans["Catching"][player.Position]["stddev"], player.Year)
			RouteRunningGrade = util.GetLetterGrade(player.RouteRunning, attributeMeans["RouteRunning"]["WR"]["mean"], attributeMeans["RouteRunning"][player.Position]["stddev"], player.Year)
			ZoneCoverageGrade = util.GetLetterGrade(player.ZoneCoverage, attributeMeans["ZoneCoverage"]["SS"]["mean"], attributeMeans["ZoneCoverage"][player.Position]["stddev"], player.Year)
			ManCoverageGrade = util.GetLetterGrade(player.ManCoverage, attributeMeans["ManCoverage"]["SS"]["mean"], attributeMeans["ManCoverage"][player.Position]["stddev"], player.Year)
			StrengthGrade = util.GetLetterGrade(player.Strength, attributeMeans["Strength"]["RB"]["mean"], attributeMeans["Strength"][player.Position]["stddev"], player.Year)
			TackleGrade = util.GetLetterGrade(player.Tackle, attributeMeans["Tackle"]["RB"]["mean"], attributeMeans["Tackle"][player.Position]["stddev"], player.Year)
			PassBlockGrade = util.GetLetterGrade(player.PassBlock, attributeMeans["PassBlock"]["RB"]["mean"], attributeMeans["PassBlock"][player.Position]["stddev"], player.Year)
			RunBlockGrade = util.GetLetterGrade(player.RunBlock, attributeMeans["RunBlock"]["RB"]["mean"], attributeMeans["RunBlock"][player.Position]["stddev"], player.Year)
			PassRushGrade = util.GetLetterGrade(player.PassRush, attributeMeans["PassRush"]["OLB"]["mean"], attributeMeans["PassRush"][player.Position]["stddev"], player.Year)
			RunDefenseGrade = util.GetLetterGrade(player.RunDefense, attributeMeans["RunDefense"]["OLB"]["mean"], attributeMeans["RunDefense"][player.Position]["stddev"], player.Year)
			ThrowPowerGrade = util.GetLetterGrade(player.ThrowPower, attributeMeans["ThrowPower"]["QB"]["mean"], attributeMeans["ThrowPower"][player.Position]["stddev"], player.Year)
			ThrowAccuracyGrade = util.GetLetterGrade(player.ThrowAccuracy, attributeMeans["ThrowAccuracy"]["QB"]["mean"], attributeMeans["ThrowAccuracy"][player.Position]["stddev"], player.Year)
			KickPowerGrade = util.GetLetterGrade(player.KickPower, attributeMeans["KickPower"]["K"]["mean"], attributeMeans["KickPower"][player.Position]["stddev"], player.Year)
			KickAccuracyGrade = util.GetLetterGrade(player.KickAccuracy, attributeMeans["KickAccuracy"]["K"]["mean"], attributeMeans["KickAccuracy"][player.Position]["stddev"], player.Year)
			PuntPowerGrade = util.GetLetterGrade(player.PuntPower, attributeMeans["PuntPower"]["P"]["mean"], attributeMeans["PuntPower"][player.Position]["stddev"], player.Year)
			PuntAccuracyGrade = util.GetLetterGrade(player.PuntAccuracy, attributeMeans["PuntAccuracy"]["P"]["mean"], attributeMeans["PuntAccuracy"][player.Position]["stddev"], player.Year)
		} else if player.Archetype == "Soccer Player" {
			StaminaGrade = util.GetLetterGrade(player.Stamina, attributeMeans["Stamina"]["RB"]["mean"], attributeMeans["Stamina"][player.Position]["stddev"], player.Year)
			InjuryGrade = util.GetLetterGrade(player.Injury, attributeMeans["Injury"]["RB"]["mean"], attributeMeans["Injury"][player.Position]["stddev"], player.Year)
			SpeedGrade = util.GetLetterGrade(player.Speed, attributeMeans["Speed"]["RB"]["mean"], attributeMeans["Speed"][player.Position]["stddev"], player.Year)
			FootballIQGrade = util.GetLetterGrade(player.FootballIQ, attributeMeans["FootballIQ"]["ILB"]["mean"], attributeMeans["FootballIQ"][player.Position]["stddev"], player.Year)
			AgilityGrade = util.GetLetterGrade(player.Agility, attributeMeans["Agility"]["RB"]["mean"], attributeMeans["Agility"][player.Position]["stddev"], player.Year)
			CarryingGrade = util.GetLetterGrade(player.Carrying, attributeMeans["Carrying"]["RB"]["mean"], attributeMeans["Carrying"][player.Position]["stddev"], player.Year)
			CatchingGrade = util.GetLetterGrade(player.Catching, attributeMeans["Catching"]["WR"]["mean"], attributeMeans["Catching"][player.Position]["stddev"], player.Year)
			RouteRunningGrade = util.GetLetterGrade(player.RouteRunning, attributeMeans["RouteRunning"]["WR"]["mean"], attributeMeans["RouteRunning"][player.Position]["stddev"], player.Year)
			ZoneCoverageGrade = util.GetLetterGrade(player.ZoneCoverage, attributeMeans["ZoneCoverage"]["SS"]["mean"], attributeMeans["ZoneCoverage"][player.Position]["stddev"], player.Year)
			ManCoverageGrade = util.GetLetterGrade(player.ManCoverage, attributeMeans["ManCoverage"]["SS"]["mean"], attributeMeans["ManCoverage"][player.Position]["stddev"], player.Year)
			StrengthGrade = util.GetLetterGrade(player.Strength, attributeMeans["Strength"]["RB"]["mean"], attributeMeans["Strength"][player.Position]["stddev"], player.Year)
			TackleGrade = util.GetLetterGrade(player.Tackle, attributeMeans["Tackle"]["RB"]["mean"], attributeMeans["Tackle"][player.Position]["stddev"], player.Year)
			PassBlockGrade = util.GetLetterGrade(player.PassBlock, attributeMeans["PassBlock"]["RB"]["mean"], attributeMeans["PassBlock"][player.Position]["stddev"], player.Year)
			RunBlockGrade = util.GetLetterGrade(player.RunBlock, attributeMeans["RunBlock"]["RB"]["mean"], attributeMeans["RunBlock"][player.Position]["stddev"], player.Year)
			PassRushGrade = util.GetLetterGrade(player.PassRush, attributeMeans["PassRush"]["OLB"]["mean"], attributeMeans["PassRush"][player.Position]["stddev"], player.Year)
			RunDefenseGrade = util.GetLetterGrade(player.RunDefense, attributeMeans["RunDefense"]["OLB"]["mean"], attributeMeans["RunDefense"][player.Position]["stddev"], player.Year)
			ThrowPowerGrade = util.GetLetterGrade(player.ThrowPower, attributeMeans["ThrowPower"]["QB"]["mean"], attributeMeans["ThrowPower"][player.Position]["stddev"], player.Year)
			ThrowAccuracyGrade = util.GetLetterGrade(player.ThrowAccuracy, attributeMeans["ThrowAccuracy"]["QB"]["mean"], attributeMeans["ThrowAccuracy"][player.Position]["stddev"], player.Year)
			KickPowerGrade = util.GetLetterGrade(player.KickPower, attributeMeans["KickPower"]["K"]["mean"], attributeMeans["KickPower"][player.Position]["stddev"], player.Year)
			KickAccuracyGrade = util.GetLetterGrade(player.KickAccuracy, attributeMeans["KickAccuracy"]["K"]["mean"], attributeMeans["KickAccuracy"][player.Position]["stddev"], player.Year)
			PuntPowerGrade = util.GetLetterGrade(player.PuntPower, attributeMeans["PuntPower"]["P"]["mean"], attributeMeans["PuntPower"][player.Position]["stddev"], player.Year)
			PuntAccuracyGrade = util.GetLetterGrade(player.PuntAccuracy, attributeMeans["PuntAccuracy"]["P"]["mean"], attributeMeans["PuntAccuracy"][player.Position]["stddev"], player.Year)
		}
	} else {
		StaminaGrade = util.GetLetterGrade(player.Stamina, attributeMeans["Stamina"][player.Position]["mean"], attributeMeans["Stamina"][player.Position]["stddev"], player.Year)
		InjuryGrade = util.GetLetterGrade(player.Injury, attributeMeans["Injury"][player.Position]["mean"], attributeMeans["Injury"][player.Position]["stddev"], player.Year)
		SpeedGrade = util.GetLetterGrade(player.Speed, attributeMeans["Speed"][player.Position]["mean"], attributeMeans["Speed"][player.Position]["stddev"], player.Year)
		FootballIQGrade = util.GetLetterGrade(player.FootballIQ, attributeMeans["FootballIQ"][player.Position]["mean"], attributeMeans["FootballIQ"][player.Position]["stddev"], player.Year)
		AgilityGrade = util.GetLetterGrade(player.Agility, attributeMeans["Agility"][player.Position]["mean"], attributeMeans["Agility"][player.Position]["stddev"], player.Year)
		CarryingGrade = util.GetLetterGrade(player.Carrying, attributeMeans["Carrying"][player.Position]["mean"], attributeMeans["Carrying"][player.Position]["stddev"], player.Year)
		CatchingGrade = util.GetLetterGrade(player.Catching, attributeMeans["Catching"][player.Position]["mean"], attributeMeans["Catching"][player.Position]["stddev"], player.Year)
		RouteRunningGrade = util.GetLetterGrade(player.RouteRunning, attributeMeans["RouteRunning"][player.Position]["mean"], attributeMeans["RouteRunning"][player.Position]["stddev"], player.Year)
		ZoneCoverageGrade = util.GetLetterGrade(player.ZoneCoverage, attributeMeans["ZoneCoverage"][player.Position]["mean"], attributeMeans["ZoneCoverage"][player.Position]["stddev"], player.Year)
		ManCoverageGrade = util.GetLetterGrade(player.ManCoverage, attributeMeans["ManCoverage"][player.Position]["mean"], attributeMeans["ManCoverage"][player.Position]["stddev"], player.Year)
		StrengthGrade = util.GetLetterGrade(player.Strength, attributeMeans["Strength"][player.Position]["mean"], attributeMeans["Strength"][player.Position]["stddev"], player.Year)
		TackleGrade = util.GetLetterGrade(player.Tackle, attributeMeans["Tackle"][player.Position]["mean"], attributeMeans["Tackle"][player.Position]["stddev"], player.Year)
		PassBlockGrade = util.GetLetterGrade(player.PassBlock, attributeMeans["PassBlock"][player.Position]["mean"], attributeMeans["PassBlock"][player.Position]["stddev"], player.Year)
		RunBlockGrade = util.GetLetterGrade(player.RunBlock, attributeMeans["RunBlock"][player.Position]["mean"], attributeMeans["RunBlock"][player.Position]["stddev"], player.Year)
		PassRushGrade = util.GetLetterGrade(player.PassRush, attributeMeans["PassRush"][player.Position]["mean"], attributeMeans["PassRush"][player.Position]["stddev"], player.Year)
		RunDefenseGrade = util.GetLetterGrade(player.RunDefense, attributeMeans["RunDefense"][player.Position]["mean"], attributeMeans["RunDefense"][player.Position]["stddev"], player.Year)
		ThrowPowerGrade = util.GetLetterGrade(player.ThrowPower, attributeMeans["ThrowPower"][player.Position]["mean"], attributeMeans["ThrowPower"][player.Position]["stddev"], player.Year)
		ThrowAccuracyGrade = util.GetLetterGrade(player.ThrowAccuracy, attributeMeans["ThrowAccuracy"][player.Position]["mean"], attributeMeans["ThrowAccuracy"][player.Position]["stddev"], player.Year)
		KickPowerGrade = util.GetLetterGrade(player.KickPower, attributeMeans["KickPower"][player.Position]["mean"], attributeMeans["KickPower"][player.Position]["stddev"], player.Year)
		KickAccuracyGrade = util.GetLetterGrade(player.KickAccuracy, attributeMeans["KickAccuracy"][player.Position]["mean"], attributeMeans["KickAccuracy"][player.Position]["stddev"], player.Year)
		PuntPowerGrade = util.GetLetterGrade(player.PuntPower, attributeMeans["PuntPower"][player.Position]["mean"], attributeMeans["PuntPower"][player.Position]["stddev"], player.Year)
		PuntAccuracyGrade = util.GetLetterGrade(player.PuntAccuracy, attributeMeans["PuntAccuracy"][player.Position]["mean"], attributeMeans["PuntAccuracy"][player.Position]["stddev"], player.Year)
	}

	return CollegePlayerCSV{
		FirstName:          player.FirstName,
		LastName:           player.LastName,
		Position:           player.Position,
		Archetype:          player.Archetype,
		PositionTwo:        player.PositionTwo,
		ArchetypeTwo:       player.ArchetypeTwo,
		Year:               Year,
		Age:                player.Age,
		Stars:              player.Stars,
		HighSchool:         player.HighSchool,
		City:               player.City,
		State:              player.State,
		Height:             player.Height,
		Weight:             player.Weight,
		OverallGrade:       OverallGrade,
		StaminaGrade:       StaminaGrade,
		InjuryGrade:        InjuryGrade,
		FootballIQGrade:    FootballIQGrade,
		SpeedGrade:         SpeedGrade,
		CarryingGrade:      CarryingGrade,
		AgilityGrade:       AgilityGrade,
		CatchingGrade:      CatchingGrade,
		RouteRunningGrade:  RouteRunningGrade,
		ZoneCoverageGrade:  ZoneCoverageGrade,
		ManCoverageGrade:   ManCoverageGrade,
		StrengthGrade:      StrengthGrade,
		TackleGrade:        TackleGrade,
		PassBlockGrade:     PassBlockGrade,
		RunBlockGrade:      RunBlockGrade,
		PassRushGrade:      PassRushGrade,
		RunDefenseGrade:    RunDefenseGrade,
		ThrowPowerGrade:    ThrowPowerGrade,
		ThrowAccuracyGrade: ThrowAccuracyGrade,
		KickAccuracyGrade:  KickAccuracyGrade,
		KickPowerGrade:     KickPowerGrade,
		PuntAccuracyGrade:  PuntAccuracyGrade,
		PuntPowerGrade:     PuntPowerGrade,
		PotentialGrade:     player.PotentialGrade,
		RedshirtStatus:     RedShirtStatus,
		Shotgun:            player.Shotgun,
		Team:               player.TeamAbbr,
	}
}

func MapPortalPlayerToCSVModel(player TransferPlayerResponse) CollegePlayerCSV {

	attributeMeans := config.AttributeMeans()
	Year, RedShirtStatus := util.GetYearAndRedshirtStatus(player.Year, player.IsRedshirt)
	OverallGrade := util.GetOverallGrade(player.Overall, player.Year)
	StaminaGrade := ""
	InjuryGrade := ""
	SpeedGrade := ""
	FootballIQGrade := ""
	AgilityGrade := ""
	CarryingGrade := ""
	CatchingGrade := ""
	RouteRunningGrade := ""
	ZoneCoverageGrade := ""
	ManCoverageGrade := ""
	StrengthGrade := ""
	TackleGrade := ""
	PassBlockGrade := ""
	RunBlockGrade := ""
	PassRushGrade := ""
	RunDefenseGrade := ""
	ThrowPowerGrade := ""
	ThrowAccuracyGrade := ""
	KickPowerGrade := ""
	KickAccuracyGrade := ""
	PuntPowerGrade := ""
	PuntAccuracyGrade := ""

	if player.Position == "ATH" {
		if player.Archetype == "Field General" {
			StaminaGrade = util.GetLetterGrade(player.Stamina, attributeMeans["Stamina"]["QB"]["mean"], attributeMeans["Stamina"][player.Position]["stddev"], player.Year)
			InjuryGrade = util.GetLetterGrade(player.Injury, attributeMeans["Injury"]["QB"]["mean"], attributeMeans["Injury"][player.Position]["stddev"], player.Year)
			SpeedGrade = util.GetLetterGrade(player.Speed, attributeMeans["Speed"]["ILB"]["mean"], attributeMeans["Speed"][player.Position]["stddev"], player.Year)
			FootballIQGrade = util.GetLetterGrade(player.FootballIQ, attributeMeans["FootballIQ"]["ILB"]["mean"], attributeMeans["FootballIQ"][player.Position]["stddev"], player.Year)
			AgilityGrade = util.GetLetterGrade(player.Agility, attributeMeans["Agility"]["QB"]["mean"], attributeMeans["Agility"][player.Position]["stddev"], player.Year)
			CarryingGrade = util.GetLetterGrade(player.Carrying, attributeMeans["Carrying"]["QB"]["mean"], attributeMeans["Carrying"][player.Position]["stddev"], player.Year)
			CatchingGrade = util.GetLetterGrade(player.Catching, attributeMeans["Catching"]["QB"]["mean"], attributeMeans["Catching"][player.Position]["stddev"], player.Year)
			RouteRunningGrade = util.GetLetterGrade(player.RouteRunning, attributeMeans["RouteRunning"]["QB"]["mean"], attributeMeans["RouteRunning"][player.Position]["stddev"], player.Year)
			ZoneCoverageGrade = util.GetLetterGrade(player.ZoneCoverage, attributeMeans["ZoneCoverage"]["ILB"]["mean"], attributeMeans["ZoneCoverage"][player.Position]["stddev"], player.Year)
			ManCoverageGrade = util.GetLetterGrade(player.ManCoverage, attributeMeans["ManCoverage"]["ILB"]["mean"], attributeMeans["ManCoverage"][player.Position]["stddev"], player.Year)
			StrengthGrade = util.GetLetterGrade(player.Strength, attributeMeans["Strength"]["ILB"]["mean"], attributeMeans["Strength"][player.Position]["stddev"], player.Year)
			TackleGrade = util.GetLetterGrade(player.Tackle, attributeMeans["Tackle"]["ILB"]["mean"], attributeMeans["Tackle"][player.Position]["stddev"], player.Year)
			PassBlockGrade = util.GetLetterGrade(player.PassBlock, attributeMeans["PassBlock"]["QB"]["mean"], attributeMeans["PassBlock"][player.Position]["stddev"], player.Year)
			RunBlockGrade = util.GetLetterGrade(player.RunBlock, attributeMeans["RunBlock"]["QB"]["mean"], attributeMeans["RunBlock"][player.Position]["stddev"], player.Year)
			PassRushGrade = util.GetLetterGrade(player.PassRush, attributeMeans["PassRush"]["ILB"]["mean"], attributeMeans["PassRush"][player.Position]["stddev"], player.Year)
			RunDefenseGrade = util.GetLetterGrade(player.RunDefense, attributeMeans["RunDefense"]["ILB"]["mean"], attributeMeans["RunDefense"][player.Position]["stddev"], player.Year)
			ThrowPowerGrade = util.GetLetterGrade(player.ThrowPower, attributeMeans["ThrowPower"]["QB"]["mean"], attributeMeans["ThrowPower"][player.Position]["stddev"], player.Year)
			ThrowAccuracyGrade = util.GetLetterGrade(player.ThrowAccuracy, attributeMeans["ThrowAccuracy"]["QB"]["mean"], attributeMeans["ThrowAccuracy"][player.Position]["stddev"], player.Year)
			KickPowerGrade = util.GetLetterGrade(player.KickPower, attributeMeans["KickPower"]["QB"]["mean"], attributeMeans["KickPower"][player.Position]["stddev"], player.Year)
			KickAccuracyGrade = util.GetLetterGrade(player.KickAccuracy, attributeMeans["KickAccuracy"]["QB"]["mean"], attributeMeans["KickAccuracy"][player.Position]["stddev"], player.Year)
			PuntPowerGrade = util.GetLetterGrade(player.PuntPower, attributeMeans["PuntPower"]["QB"]["mean"], attributeMeans["PuntPower"][player.Position]["stddev"], player.Year)
			PuntAccuracyGrade = util.GetLetterGrade(player.PuntAccuracy, attributeMeans["PuntAccuracy"]["QB"]["mean"], attributeMeans["PuntAccuracy"][player.Position]["stddev"], player.Year)
		} else if player.Archetype == "Triple-Threat" {
			StaminaGrade = util.GetLetterGrade(player.Stamina, attributeMeans["Stamina"]["QB"]["mean"], attributeMeans["Stamina"][player.Position]["stddev"], player.Year)
			InjuryGrade = util.GetLetterGrade(player.Injury, attributeMeans["Injury"]["QB"]["mean"], attributeMeans["Injury"][player.Position]["stddev"], player.Year)
			SpeedGrade = util.GetLetterGrade(player.Speed, attributeMeans["Speed"]["RB"]["mean"], attributeMeans["Speed"][player.Position]["stddev"], player.Year)
			FootballIQGrade = util.GetLetterGrade(player.FootballIQ, attributeMeans["FootballIQ"]["QB"]["mean"], attributeMeans["FootballIQ"][player.Position]["stddev"], player.Year)
			AgilityGrade = util.GetLetterGrade(player.Agility, attributeMeans["Agility"]["RB"]["mean"], attributeMeans["Agility"][player.Position]["stddev"], player.Year)
			CarryingGrade = util.GetLetterGrade(player.Carrying, attributeMeans["Carrying"]["RB"]["mean"], attributeMeans["Carrying"][player.Position]["stddev"], player.Year)
			CatchingGrade = util.GetLetterGrade(player.Catching, attributeMeans["Catching"]["WR"]["mean"], attributeMeans["Catching"][player.Position]["stddev"], player.Year)
			RouteRunningGrade = util.GetLetterGrade(player.RouteRunning, attributeMeans["RouteRunning"]["WR"]["mean"], attributeMeans["RouteRunning"][player.Position]["stddev"], player.Year)
			ZoneCoverageGrade = util.GetLetterGrade(player.ZoneCoverage, attributeMeans["ZoneCoverage"]["CB"]["mean"], attributeMeans["ZoneCoverage"][player.Position]["stddev"], player.Year)
			ManCoverageGrade = util.GetLetterGrade(player.ManCoverage, attributeMeans["ManCoverage"]["CB"]["mean"], attributeMeans["ManCoverage"][player.Position]["stddev"], player.Year)
			StrengthGrade = util.GetLetterGrade(player.Strength, attributeMeans["Strength"]["RB"]["mean"], attributeMeans["Strength"][player.Position]["stddev"], player.Year)
			TackleGrade = util.GetLetterGrade(player.Tackle, attributeMeans["Tackle"]["WR"]["mean"], attributeMeans["Tackle"][player.Position]["stddev"], player.Year)
			PassBlockGrade = util.GetLetterGrade(player.PassBlock, attributeMeans["PassBlock"]["RB"]["mean"], attributeMeans["PassBlock"][player.Position]["stddev"], player.Year)
			RunBlockGrade = util.GetLetterGrade(player.RunBlock, attributeMeans["RunBlock"]["RB"]["mean"], attributeMeans["RunBlock"][player.Position]["stddev"], player.Year)
			PassRushGrade = util.GetLetterGrade(player.PassRush, attributeMeans["PassRush"]["RB"]["mean"], attributeMeans["PassRush"][player.Position]["stddev"], player.Year)
			RunDefenseGrade = util.GetLetterGrade(player.RunDefense, attributeMeans["RunDefense"]["RB"]["mean"], attributeMeans["RunDefense"][player.Position]["stddev"], player.Year)
			ThrowPowerGrade = util.GetLetterGrade(player.ThrowPower, attributeMeans["ThrowPower"]["QB"]["mean"], attributeMeans["ThrowPower"][player.Position]["stddev"], player.Year)
			ThrowAccuracyGrade = util.GetLetterGrade(player.ThrowAccuracy, attributeMeans["ThrowAccuracy"]["QB"]["mean"], attributeMeans["ThrowAccuracy"][player.Position]["stddev"], player.Year)
			KickPowerGrade = util.GetLetterGrade(player.KickPower, attributeMeans["KickPower"]["QB"]["mean"], attributeMeans["KickPower"][player.Position]["stddev"], player.Year)
			KickAccuracyGrade = util.GetLetterGrade(player.KickAccuracy, attributeMeans["KickAccuracy"]["QB"]["mean"], attributeMeans["KickAccuracy"][player.Position]["stddev"], player.Year)
			PuntPowerGrade = util.GetLetterGrade(player.PuntPower, attributeMeans["PuntPower"]["QB"]["mean"], attributeMeans["PuntPower"][player.Position]["stddev"], player.Year)
			PuntAccuracyGrade = util.GetLetterGrade(player.PuntAccuracy, attributeMeans["PuntAccuracy"]["QB"]["mean"], attributeMeans["PuntAccuracy"][player.Position]["stddev"], player.Year)
		} else if player.Archetype == "Wingback" {
			StaminaGrade = util.GetLetterGrade(player.Stamina, attributeMeans["Stamina"]["QB"]["mean"], attributeMeans["Stamina"][player.Position]["stddev"], player.Year)
			InjuryGrade = util.GetLetterGrade(player.Injury, attributeMeans["Injury"]["QB"]["mean"], attributeMeans["Injury"][player.Position]["stddev"], player.Year)
			SpeedGrade = util.GetLetterGrade(player.Speed, attributeMeans["Speed"]["RB"]["mean"], attributeMeans["Speed"][player.Position]["stddev"], player.Year)
			FootballIQGrade = util.GetLetterGrade(player.FootballIQ, attributeMeans["FootballIQ"]["RB"]["mean"], attributeMeans["FootballIQ"][player.Position]["stddev"], player.Year)
			AgilityGrade = util.GetLetterGrade(player.Agility, attributeMeans["Agility"]["RB"]["mean"], attributeMeans["Agility"][player.Position]["stddev"], player.Year)
			CarryingGrade = util.GetLetterGrade(player.Carrying, attributeMeans["Carrying"]["RB"]["mean"], attributeMeans["Carrying"][player.Position]["stddev"], player.Year)
			CatchingGrade = util.GetLetterGrade(player.Catching, attributeMeans["Catching"]["WR"]["mean"], attributeMeans["Catching"][player.Position]["stddev"], player.Year)
			RouteRunningGrade = util.GetLetterGrade(player.RouteRunning, attributeMeans["RouteRunning"]["WR"]["mean"], attributeMeans["RouteRunning"][player.Position]["stddev"], player.Year)
			ZoneCoverageGrade = util.GetLetterGrade(player.ZoneCoverage, attributeMeans["ZoneCoverage"]["CB"]["mean"], attributeMeans["ZoneCoverage"][player.Position]["stddev"], player.Year)
			ManCoverageGrade = util.GetLetterGrade(player.ManCoverage, attributeMeans["ManCoverage"]["CB"]["mean"], attributeMeans["ManCoverage"][player.Position]["stddev"], player.Year)
			StrengthGrade = util.GetLetterGrade(player.Strength, attributeMeans["Strength"]["RB"]["mean"], attributeMeans["Strength"][player.Position]["stddev"], player.Year)
			TackleGrade = util.GetLetterGrade(player.Tackle, attributeMeans["Tackle"]["WR"]["mean"], attributeMeans["Tackle"][player.Position]["stddev"], player.Year)
			PassBlockGrade = util.GetLetterGrade(player.PassBlock, attributeMeans["PassBlock"]["RB"]["mean"], attributeMeans["PassBlock"][player.Position]["stddev"], player.Year)
			RunBlockGrade = util.GetLetterGrade(player.RunBlock, attributeMeans["RunBlock"]["RB"]["mean"], attributeMeans["RunBlock"][player.Position]["stddev"], player.Year)
			PassRushGrade = util.GetLetterGrade(player.PassRush, attributeMeans["PassRush"]["CB"]["mean"], attributeMeans["PassRush"][player.Position]["stddev"], player.Year)
			RunDefenseGrade = util.GetLetterGrade(player.RunDefense, attributeMeans["RunDefense"]["CB"]["mean"], attributeMeans["RunDefense"][player.Position]["stddev"], player.Year)
			ThrowPowerGrade = util.GetLetterGrade(player.ThrowPower, attributeMeans["ThrowPower"]["QB"]["mean"], attributeMeans["ThrowPower"][player.Position]["stddev"], player.Year)
			ThrowAccuracyGrade = util.GetLetterGrade(player.ThrowAccuracy, attributeMeans["ThrowAccuracy"]["QB"]["mean"], attributeMeans["ThrowAccuracy"][player.Position]["stddev"], player.Year)
			KickPowerGrade = util.GetLetterGrade(player.KickPower, attributeMeans["KickPower"]["QB"]["mean"], attributeMeans["KickPower"][player.Position]["stddev"], player.Year)
			KickAccuracyGrade = util.GetLetterGrade(player.KickAccuracy, attributeMeans["KickAccuracy"]["QB"]["mean"], attributeMeans["KickAccuracy"][player.Position]["stddev"], player.Year)
			PuntPowerGrade = util.GetLetterGrade(player.PuntPower, attributeMeans["PuntPower"]["QB"]["mean"], attributeMeans["PuntPower"][player.Position]["stddev"], player.Year)
			PuntAccuracyGrade = util.GetLetterGrade(player.PuntAccuracy, attributeMeans["PuntAccuracy"]["QB"]["mean"], attributeMeans["PuntAccuracy"][player.Position]["stddev"], player.Year)
		} else if player.Archetype == "Slotback" {
			StaminaGrade = util.GetLetterGrade(player.Stamina, attributeMeans["Stamina"]["QB"]["mean"], attributeMeans["Stamina"][player.Position]["stddev"], player.Year)
			InjuryGrade = util.GetLetterGrade(player.Injury, attributeMeans["Injury"]["QB"]["mean"], attributeMeans["Injury"][player.Position]["stddev"], player.Year)
			SpeedGrade = util.GetLetterGrade(player.Speed, attributeMeans["Speed"]["WR"]["mean"], attributeMeans["Speed"][player.Position]["stddev"], player.Year)
			FootballIQGrade = util.GetLetterGrade(player.FootballIQ, attributeMeans["FootballIQ"]["WR"]["mean"], attributeMeans["FootballIQ"][player.Position]["stddev"], player.Year)
			AgilityGrade = util.GetLetterGrade(player.Agility, attributeMeans["Agility"]["WR"]["mean"], attributeMeans["Agility"][player.Position]["stddev"], player.Year)
			CarryingGrade = util.GetLetterGrade(player.Carrying, attributeMeans["Carrying"]["RB"]["mean"], attributeMeans["Carrying"][player.Position]["stddev"], player.Year)
			CatchingGrade = util.GetLetterGrade(player.Catching, attributeMeans["Catching"]["WR"]["mean"], attributeMeans["Catching"][player.Position]["stddev"], player.Year)
			RouteRunningGrade = util.GetLetterGrade(player.RouteRunning, attributeMeans["RouteRunning"]["WR"]["mean"], attributeMeans["RouteRunning"][player.Position]["stddev"], player.Year)
			ZoneCoverageGrade = util.GetLetterGrade(player.ZoneCoverage, attributeMeans["ZoneCoverage"]["CB"]["mean"], attributeMeans["ZoneCoverage"][player.Position]["stddev"], player.Year)
			ManCoverageGrade = util.GetLetterGrade(player.ManCoverage, attributeMeans["ManCoverage"]["CB"]["mean"], attributeMeans["ManCoverage"][player.Position]["stddev"], player.Year)
			StrengthGrade = util.GetLetterGrade(player.Strength, attributeMeans["Strength"]["TE"]["mean"], attributeMeans["Strength"][player.Position]["stddev"], player.Year)
			TackleGrade = util.GetLetterGrade(player.Tackle, attributeMeans["Tackle"]["WR"]["mean"], attributeMeans["Tackle"][player.Position]["stddev"], player.Year)
			PassBlockGrade = util.GetLetterGrade(player.PassBlock, attributeMeans["PassBlock"]["FB"]["mean"], attributeMeans["PassBlock"][player.Position]["stddev"], player.Year)
			RunBlockGrade = util.GetLetterGrade(player.RunBlock, attributeMeans["RunBlock"]["FB"]["mean"], attributeMeans["RunBlock"][player.Position]["stddev"], player.Year)
			PassRushGrade = util.GetLetterGrade(player.PassRush, attributeMeans["PassRush"]["TE"]["mean"], attributeMeans["PassRush"][player.Position]["stddev"], player.Year)
			RunDefenseGrade = util.GetLetterGrade(player.RunDefense, attributeMeans["RunDefense"]["TE"]["mean"], attributeMeans["RunDefense"][player.Position]["stddev"], player.Year)
			ThrowPowerGrade = util.GetLetterGrade(player.ThrowPower, attributeMeans["ThrowPower"]["QB"]["mean"], attributeMeans["ThrowPower"][player.Position]["stddev"], player.Year)
			ThrowAccuracyGrade = util.GetLetterGrade(player.ThrowAccuracy, attributeMeans["ThrowAccuracy"]["QB"]["mean"], attributeMeans["ThrowAccuracy"][player.Position]["stddev"], player.Year)
			KickPowerGrade = util.GetLetterGrade(player.KickPower, attributeMeans["KickPower"]["QB"]["mean"], attributeMeans["KickPower"][player.Position]["stddev"], player.Year)
			KickAccuracyGrade = util.GetLetterGrade(player.KickAccuracy, attributeMeans["KickAccuracy"]["QB"]["mean"], attributeMeans["KickAccuracy"][player.Position]["stddev"], player.Year)
			PuntPowerGrade = util.GetLetterGrade(player.PuntPower, attributeMeans["PuntPower"]["QB"]["mean"], attributeMeans["PuntPower"][player.Position]["stddev"], player.Year)
			PuntAccuracyGrade = util.GetLetterGrade(player.PuntAccuracy, attributeMeans["PuntAccuracy"]["QB"]["mean"], attributeMeans["PuntAccuracy"][player.Position]["stddev"], player.Year)
		} else if player.Archetype == "Lineman" {
			StaminaGrade = util.GetLetterGrade(player.Stamina, attributeMeans["Stamina"]["OG"]["mean"], attributeMeans["Stamina"][player.Position]["stddev"], player.Year)
			InjuryGrade = util.GetLetterGrade(player.Injury, attributeMeans["Injury"]["OG"]["mean"], attributeMeans["Injury"][player.Position]["stddev"], player.Year)
			SpeedGrade = util.GetLetterGrade(player.Speed, attributeMeans["Speed"]["DE"]["mean"], attributeMeans["Speed"][player.Position]["stddev"], player.Year)
			FootballIQGrade = util.GetLetterGrade(player.FootballIQ, attributeMeans["FootballIQ"]["C"]["mean"], attributeMeans["FootballIQ"][player.Position]["stddev"], player.Year)
			AgilityGrade = util.GetLetterGrade(player.Agility, attributeMeans["Agility"]["DE"]["mean"], attributeMeans["Agility"][player.Position]["stddev"], player.Year)
			CarryingGrade = util.GetLetterGrade(player.Carrying, attributeMeans["Carrying"]["RB"]["mean"], attributeMeans["Carrying"][player.Position]["stddev"], player.Year)
			CatchingGrade = util.GetLetterGrade(player.Catching, attributeMeans["Catching"]["RB"]["mean"], attributeMeans["Catching"][player.Position]["stddev"], player.Year)
			RouteRunningGrade = util.GetLetterGrade(player.RouteRunning, attributeMeans["RouteRunning"]["OG"]["mean"], attributeMeans["RouteRunning"][player.Position]["stddev"], player.Year)
			ZoneCoverageGrade = util.GetLetterGrade(player.ZoneCoverage, attributeMeans["ZoneCoverage"]["CB"]["mean"], attributeMeans["ZoneCoverage"][player.Position]["stddev"], player.Year)
			ManCoverageGrade = util.GetLetterGrade(player.ManCoverage, attributeMeans["ManCoverage"]["CB"]["mean"], attributeMeans["ManCoverage"][player.Position]["stddev"], player.Year)
			StrengthGrade = util.GetLetterGrade(player.Strength, attributeMeans["Strength"]["DT"]["mean"], attributeMeans["Strength"][player.Position]["stddev"], player.Year)
			TackleGrade = util.GetLetterGrade(player.Tackle, attributeMeans["Tackle"]["DT"]["mean"], attributeMeans["Tackle"][player.Position]["stddev"], player.Year)
			PassBlockGrade = util.GetLetterGrade(player.PassBlock, attributeMeans["PassBlock"]["OG"]["mean"], attributeMeans["PassBlock"][player.Position]["stddev"], player.Year)
			RunBlockGrade = util.GetLetterGrade(player.RunBlock, attributeMeans["RunBlock"]["OG"]["mean"], attributeMeans["RunBlock"][player.Position]["stddev"], player.Year)
			PassRushGrade = util.GetLetterGrade(player.PassRush, attributeMeans["PassRush"]["DE"]["mean"], attributeMeans["PassRush"][player.Position]["stddev"], player.Year)
			RunDefenseGrade = util.GetLetterGrade(player.RunDefense, attributeMeans["RunDefense"]["DE"]["mean"], attributeMeans["RunDefense"][player.Position]["stddev"], player.Year)
			ThrowPowerGrade = util.GetLetterGrade(player.ThrowPower, attributeMeans["ThrowPower"]["QB"]["mean"], attributeMeans["ThrowPower"][player.Position]["stddev"], player.Year)
			ThrowAccuracyGrade = util.GetLetterGrade(player.ThrowAccuracy, attributeMeans["ThrowAccuracy"]["QB"]["mean"], attributeMeans["ThrowAccuracy"][player.Position]["stddev"], player.Year)
			KickPowerGrade = util.GetLetterGrade(player.KickPower, attributeMeans["KickPower"]["K"]["mean"], attributeMeans["KickPower"][player.Position]["stddev"], player.Year)
			KickAccuracyGrade = util.GetLetterGrade(player.KickAccuracy, attributeMeans["KickAccuracy"]["K"]["mean"], attributeMeans["KickAccuracy"][player.Position]["stddev"], player.Year)
			PuntPowerGrade = util.GetLetterGrade(player.PuntPower, attributeMeans["PuntPower"]["P"]["mean"], attributeMeans["PuntPower"][player.Position]["stddev"], player.Year)
			PuntAccuracyGrade = util.GetLetterGrade(player.PuntAccuracy, attributeMeans["PuntAccuracy"]["P"]["mean"], attributeMeans["PuntAccuracy"][player.Position]["stddev"], player.Year)
		} else if player.Archetype == "Strongside" {
			StaminaGrade = util.GetLetterGrade(player.Stamina, attributeMeans["Stamina"]["OG"]["mean"], attributeMeans["Stamina"][player.Position]["stddev"], player.Year)
			InjuryGrade = util.GetLetterGrade(player.Injury, attributeMeans["Injury"]["OG"]["mean"], attributeMeans["Injury"][player.Position]["stddev"], player.Year)
			SpeedGrade = util.GetLetterGrade(player.Speed, attributeMeans["Speed"]["DE"]["mean"], attributeMeans["Speed"][player.Position]["stddev"], player.Year)
			FootballIQGrade = util.GetLetterGrade(player.FootballIQ, attributeMeans["FootballIQ"]["OLB"]["mean"], attributeMeans["FootballIQ"][player.Position]["stddev"], player.Year)
			AgilityGrade = util.GetLetterGrade(player.Agility, attributeMeans["Agility"]["DE"]["mean"], attributeMeans["Agility"][player.Position]["stddev"], player.Year)
			CarryingGrade = util.GetLetterGrade(player.Carrying, attributeMeans["Carrying"]["OLB"]["mean"], attributeMeans["Carrying"][player.Position]["stddev"], player.Year)
			CatchingGrade = util.GetLetterGrade(player.Catching, attributeMeans["Catching"]["OLB"]["mean"], attributeMeans["Catching"][player.Position]["stddev"], player.Year)
			RouteRunningGrade = util.GetLetterGrade(player.RouteRunning, attributeMeans["RouteRunning"]["OG"]["mean"], attributeMeans["RouteRunning"][player.Position]["stddev"], player.Year)
			ZoneCoverageGrade = util.GetLetterGrade(player.ZoneCoverage, attributeMeans["ZoneCoverage"]["OLB"]["mean"], attributeMeans["ZoneCoverage"][player.Position]["stddev"], player.Year)
			ManCoverageGrade = util.GetLetterGrade(player.ManCoverage, attributeMeans["ManCoverage"]["OLB"]["mean"], attributeMeans["ManCoverage"][player.Position]["stddev"], player.Year)
			StrengthGrade = util.GetLetterGrade(player.Strength, attributeMeans["Strength"]["DE"]["mean"], attributeMeans["Strength"][player.Position]["stddev"], player.Year)
			TackleGrade = util.GetLetterGrade(player.Tackle, attributeMeans["Tackle"]["DE"]["mean"], attributeMeans["Tackle"][player.Position]["stddev"], player.Year)
			PassBlockGrade = util.GetLetterGrade(player.PassBlock, attributeMeans["PassBlock"]["OG"]["mean"], attributeMeans["PassBlock"][player.Position]["stddev"], player.Year)
			RunBlockGrade = util.GetLetterGrade(player.RunBlock, attributeMeans["RunBlock"]["OG"]["mean"], attributeMeans["RunBlock"][player.Position]["stddev"], player.Year)
			PassRushGrade = util.GetLetterGrade(player.PassRush, attributeMeans["PassRush"]["DE"]["mean"], attributeMeans["PassRush"][player.Position]["stddev"], player.Year)
			RunDefenseGrade = util.GetLetterGrade(player.RunDefense, attributeMeans["RunDefense"]["DE"]["mean"], attributeMeans["RunDefense"][player.Position]["stddev"], player.Year)
			ThrowPowerGrade = util.GetLetterGrade(player.ThrowPower, attributeMeans["ThrowPower"]["QB"]["mean"], attributeMeans["ThrowPower"][player.Position]["stddev"], player.Year)
			ThrowAccuracyGrade = util.GetLetterGrade(player.ThrowAccuracy, attributeMeans["ThrowAccuracy"]["QB"]["mean"], attributeMeans["ThrowAccuracy"][player.Position]["stddev"], player.Year)
			KickPowerGrade = util.GetLetterGrade(player.KickPower, attributeMeans["KickPower"]["K"]["mean"], attributeMeans["KickPower"][player.Position]["stddev"], player.Year)
			KickAccuracyGrade = util.GetLetterGrade(player.KickAccuracy, attributeMeans["KickAccuracy"]["K"]["mean"], attributeMeans["KickAccuracy"][player.Position]["stddev"], player.Year)
			PuntPowerGrade = util.GetLetterGrade(player.PuntPower, attributeMeans["PuntPower"]["P"]["mean"], attributeMeans["PuntPower"][player.Position]["stddev"], player.Year)
			PuntAccuracyGrade = util.GetLetterGrade(player.PuntAccuracy, attributeMeans["PuntAccuracy"]["P"]["mean"], attributeMeans["PuntAccuracy"][player.Position]["stddev"], player.Year)
		} else if player.Archetype == "Weakside" {
			StaminaGrade = util.GetLetterGrade(player.Stamina, attributeMeans["Stamina"]["OG"]["mean"], attributeMeans["Stamina"][player.Position]["stddev"], player.Year)
			InjuryGrade = util.GetLetterGrade(player.Injury, attributeMeans["Injury"]["OG"]["mean"], attributeMeans["Injury"][player.Position]["stddev"], player.Year)
			SpeedGrade = util.GetLetterGrade(player.Speed, attributeMeans["Speed"]["ILB"]["mean"], attributeMeans["Speed"][player.Position]["stddev"], player.Year)
			FootballIQGrade = util.GetLetterGrade(player.FootballIQ, attributeMeans["FootballIQ"]["ILB"]["mean"], attributeMeans["FootballIQ"][player.Position]["stddev"], player.Year)
			AgilityGrade = util.GetLetterGrade(player.Agility, attributeMeans["Agility"]["ILB"]["mean"], attributeMeans["Agility"][player.Position]["stddev"], player.Year)
			CarryingGrade = util.GetLetterGrade(player.Carrying, attributeMeans["Carrying"]["OLB"]["mean"], attributeMeans["Carrying"][player.Position]["stddev"], player.Year)
			CatchingGrade = util.GetLetterGrade(player.Catching, attributeMeans["Catching"]["OLB"]["mean"], attributeMeans["Catching"][player.Position]["stddev"], player.Year)
			RouteRunningGrade = util.GetLetterGrade(player.RouteRunning, attributeMeans["RouteRunning"]["OG"]["mean"], attributeMeans["RouteRunning"][player.Position]["stddev"], player.Year)
			ZoneCoverageGrade = util.GetLetterGrade(player.ZoneCoverage, attributeMeans["ZoneCoverage"]["ILB"]["mean"], attributeMeans["ZoneCoverage"][player.Position]["stddev"], player.Year)
			ManCoverageGrade = util.GetLetterGrade(player.ManCoverage, attributeMeans["ManCoverage"]["ILB"]["mean"], attributeMeans["ManCoverage"][player.Position]["stddev"], player.Year)
			StrengthGrade = util.GetLetterGrade(player.Strength, attributeMeans["Strength"]["OLB"]["mean"], attributeMeans["Strength"][player.Position]["stddev"], player.Year)
			TackleGrade = util.GetLetterGrade(player.Tackle, attributeMeans["Tackle"]["OLB"]["mean"], attributeMeans["Tackle"][player.Position]["stddev"], player.Year)
			PassBlockGrade = util.GetLetterGrade(player.PassBlock, attributeMeans["PassBlock"]["OG"]["mean"], attributeMeans["PassBlock"][player.Position]["stddev"], player.Year)
			RunBlockGrade = util.GetLetterGrade(player.RunBlock, attributeMeans["RunBlock"]["OG"]["mean"], attributeMeans["RunBlock"][player.Position]["stddev"], player.Year)
			PassRushGrade = util.GetLetterGrade(player.PassRush, attributeMeans["PassRush"]["ILB"]["mean"], attributeMeans["PassRush"][player.Position]["stddev"], player.Year)
			RunDefenseGrade = util.GetLetterGrade(player.RunDefense, attributeMeans["RunDefense"]["ILB"]["mean"], attributeMeans["RunDefense"][player.Position]["stddev"], player.Year)
			ThrowPowerGrade = util.GetLetterGrade(player.ThrowPower, attributeMeans["ThrowPower"]["QB"]["mean"], attributeMeans["ThrowPower"][player.Position]["stddev"], player.Year)
			ThrowAccuracyGrade = util.GetLetterGrade(player.ThrowAccuracy, attributeMeans["ThrowAccuracy"]["QB"]["mean"], attributeMeans["ThrowAccuracy"][player.Position]["stddev"], player.Year)
			KickPowerGrade = util.GetLetterGrade(player.KickPower, attributeMeans["KickPower"]["K"]["mean"], attributeMeans["KickPower"][player.Position]["stddev"], player.Year)
			KickAccuracyGrade = util.GetLetterGrade(player.KickAccuracy, attributeMeans["KickAccuracy"]["K"]["mean"], attributeMeans["KickAccuracy"][player.Position]["stddev"], player.Year)
			PuntPowerGrade = util.GetLetterGrade(player.PuntPower, attributeMeans["PuntPower"]["P"]["mean"], attributeMeans["PuntPower"][player.Position]["stddev"], player.Year)
			PuntAccuracyGrade = util.GetLetterGrade(player.PuntAccuracy, attributeMeans["PuntAccuracy"]["P"]["mean"], attributeMeans["PuntAccuracy"][player.Position]["stddev"], player.Year)
		} else if player.Archetype == "Bandit" {
			StaminaGrade = util.GetLetterGrade(player.Stamina, attributeMeans["Stamina"]["OG"]["mean"], attributeMeans["Stamina"][player.Position]["stddev"], player.Year)
			InjuryGrade = util.GetLetterGrade(player.Injury, attributeMeans["Injury"]["OG"]["mean"], attributeMeans["Injury"][player.Position]["stddev"], player.Year)
			SpeedGrade = util.GetLetterGrade(player.Speed, attributeMeans["Speed"]["SS"]["mean"], attributeMeans["Speed"][player.Position]["stddev"], player.Year)
			FootballIQGrade = util.GetLetterGrade(player.FootballIQ, attributeMeans["FootballIQ"]["ILB"]["mean"], attributeMeans["FootballIQ"][player.Position]["stddev"], player.Year)
			AgilityGrade = util.GetLetterGrade(player.Agility, attributeMeans["Agility"]["SS"]["mean"], attributeMeans["Agility"][player.Position]["stddev"], player.Year)
			CarryingGrade = util.GetLetterGrade(player.Carrying, attributeMeans["Carrying"]["SS"]["mean"], attributeMeans["Carrying"][player.Position]["stddev"], player.Year)
			CatchingGrade = util.GetLetterGrade(player.Catching, attributeMeans["Catching"]["SS"]["mean"], attributeMeans["Catching"][player.Position]["stddev"], player.Year)
			RouteRunningGrade = util.GetLetterGrade(player.RouteRunning, attributeMeans["RouteRunning"]["SS"]["mean"], attributeMeans["RouteRunning"][player.Position]["stddev"], player.Year)
			ZoneCoverageGrade = util.GetLetterGrade(player.ZoneCoverage, attributeMeans["ZoneCoverage"]["SS"]["mean"], attributeMeans["ZoneCoverage"][player.Position]["stddev"], player.Year)
			ManCoverageGrade = util.GetLetterGrade(player.ManCoverage, attributeMeans["ManCoverage"]["SS"]["mean"], attributeMeans["ManCoverage"][player.Position]["stddev"], player.Year)
			StrengthGrade = util.GetLetterGrade(player.Strength, attributeMeans["Strength"]["OLB"]["mean"], attributeMeans["Strength"][player.Position]["stddev"], player.Year)
			TackleGrade = util.GetLetterGrade(player.Tackle, attributeMeans["Tackle"]["OLB"]["mean"], attributeMeans["Tackle"][player.Position]["stddev"], player.Year)
			PassBlockGrade = util.GetLetterGrade(player.PassBlock, attributeMeans["PassBlock"]["OG"]["mean"], attributeMeans["PassBlock"][player.Position]["stddev"], player.Year)
			RunBlockGrade = util.GetLetterGrade(player.RunBlock, attributeMeans["RunBlock"]["OG"]["mean"], attributeMeans["RunBlock"][player.Position]["stddev"], player.Year)
			PassRushGrade = util.GetLetterGrade(player.PassRush, attributeMeans["PassRush"]["OLB"]["mean"], attributeMeans["PassRush"][player.Position]["stddev"], player.Year)
			RunDefenseGrade = util.GetLetterGrade(player.RunDefense, attributeMeans["RunDefense"]["OLB"]["mean"], attributeMeans["RunDefense"][player.Position]["stddev"], player.Year)
			ThrowPowerGrade = util.GetLetterGrade(player.ThrowPower, attributeMeans["ThrowPower"]["QB"]["mean"], attributeMeans["ThrowPower"][player.Position]["stddev"], player.Year)
			ThrowAccuracyGrade = util.GetLetterGrade(player.ThrowAccuracy, attributeMeans["ThrowAccuracy"]["QB"]["mean"], attributeMeans["ThrowAccuracy"][player.Position]["stddev"], player.Year)
			KickPowerGrade = util.GetLetterGrade(player.KickPower, attributeMeans["KickPower"]["K"]["mean"], attributeMeans["KickPower"][player.Position]["stddev"], player.Year)
			KickAccuracyGrade = util.GetLetterGrade(player.KickAccuracy, attributeMeans["KickAccuracy"]["K"]["mean"], attributeMeans["KickAccuracy"][player.Position]["stddev"], player.Year)
			PuntPowerGrade = util.GetLetterGrade(player.PuntPower, attributeMeans["PuntPower"]["P"]["mean"], attributeMeans["PuntPower"][player.Position]["stddev"], player.Year)
			PuntAccuracyGrade = util.GetLetterGrade(player.PuntAccuracy, attributeMeans["PuntAccuracy"]["P"]["mean"], attributeMeans["PuntAccuracy"][player.Position]["stddev"], player.Year)
		} else if player.Archetype == "Return Specialist" {
			StaminaGrade = util.GetLetterGrade(player.Stamina, attributeMeans["Stamina"]["OG"]["mean"], attributeMeans["Stamina"][player.Position]["stddev"], player.Year)
			InjuryGrade = util.GetLetterGrade(player.Injury, attributeMeans["Injury"]["OG"]["mean"], attributeMeans["Injury"][player.Position]["stddev"], player.Year)
			SpeedGrade = util.GetLetterGrade(player.Speed, attributeMeans["Speed"]["RB"]["mean"], attributeMeans["Speed"][player.Position]["stddev"], player.Year)
			FootballIQGrade = util.GetLetterGrade(player.FootballIQ, attributeMeans["FootballIQ"]["ILB"]["mean"], attributeMeans["FootballIQ"][player.Position]["stddev"], player.Year)
			AgilityGrade = util.GetLetterGrade(player.Agility, attributeMeans["Agility"]["RB"]["mean"], attributeMeans["Agility"][player.Position]["stddev"], player.Year)
			CarryingGrade = util.GetLetterGrade(player.Carrying, attributeMeans["Carrying"]["RB"]["mean"], attributeMeans["Carrying"][player.Position]["stddev"], player.Year)
			CatchingGrade = util.GetLetterGrade(player.Catching, attributeMeans["Catching"]["WR"]["mean"], attributeMeans["Catching"][player.Position]["stddev"], player.Year)
			RouteRunningGrade = util.GetLetterGrade(player.RouteRunning, attributeMeans["RouteRunning"]["WR"]["mean"], attributeMeans["RouteRunning"][player.Position]["stddev"], player.Year)
			ZoneCoverageGrade = util.GetLetterGrade(player.ZoneCoverage, attributeMeans["ZoneCoverage"]["SS"]["mean"], attributeMeans["ZoneCoverage"][player.Position]["stddev"], player.Year)
			ManCoverageGrade = util.GetLetterGrade(player.ManCoverage, attributeMeans["ManCoverage"]["SS"]["mean"], attributeMeans["ManCoverage"][player.Position]["stddev"], player.Year)
			StrengthGrade = util.GetLetterGrade(player.Strength, attributeMeans["Strength"]["RB"]["mean"], attributeMeans["Strength"][player.Position]["stddev"], player.Year)
			TackleGrade = util.GetLetterGrade(player.Tackle, attributeMeans["Tackle"]["RB"]["mean"], attributeMeans["Tackle"][player.Position]["stddev"], player.Year)
			PassBlockGrade = util.GetLetterGrade(player.PassBlock, attributeMeans["PassBlock"]["RB"]["mean"], attributeMeans["PassBlock"][player.Position]["stddev"], player.Year)
			RunBlockGrade = util.GetLetterGrade(player.RunBlock, attributeMeans["RunBlock"]["RB"]["mean"], attributeMeans["RunBlock"][player.Position]["stddev"], player.Year)
			PassRushGrade = util.GetLetterGrade(player.PassRush, attributeMeans["PassRush"]["OLB"]["mean"], attributeMeans["PassRush"][player.Position]["stddev"], player.Year)
			RunDefenseGrade = util.GetLetterGrade(player.RunDefense, attributeMeans["RunDefense"]["OLB"]["mean"], attributeMeans["RunDefense"][player.Position]["stddev"], player.Year)
			ThrowPowerGrade = util.GetLetterGrade(player.ThrowPower, attributeMeans["ThrowPower"]["QB"]["mean"], attributeMeans["ThrowPower"][player.Position]["stddev"], player.Year)
			ThrowAccuracyGrade = util.GetLetterGrade(player.ThrowAccuracy, attributeMeans["ThrowAccuracy"]["QB"]["mean"], attributeMeans["ThrowAccuracy"][player.Position]["stddev"], player.Year)
			KickPowerGrade = util.GetLetterGrade(player.KickPower, attributeMeans["KickPower"]["K"]["mean"], attributeMeans["KickPower"][player.Position]["stddev"], player.Year)
			KickAccuracyGrade = util.GetLetterGrade(player.KickAccuracy, attributeMeans["KickAccuracy"]["K"]["mean"], attributeMeans["KickAccuracy"][player.Position]["stddev"], player.Year)
			PuntPowerGrade = util.GetLetterGrade(player.PuntPower, attributeMeans["PuntPower"]["P"]["mean"], attributeMeans["PuntPower"][player.Position]["stddev"], player.Year)
			PuntAccuracyGrade = util.GetLetterGrade(player.PuntAccuracy, attributeMeans["PuntAccuracy"]["P"]["mean"], attributeMeans["PuntAccuracy"][player.Position]["stddev"], player.Year)
		} else if player.Archetype == "Soccer Player" {
			StaminaGrade = util.GetLetterGrade(player.Stamina, attributeMeans["Stamina"]["RB"]["mean"], attributeMeans["Stamina"][player.Position]["stddev"], player.Year)
			InjuryGrade = util.GetLetterGrade(player.Injury, attributeMeans["Injury"]["RB"]["mean"], attributeMeans["Injury"][player.Position]["stddev"], player.Year)
			SpeedGrade = util.GetLetterGrade(player.Speed, attributeMeans["Speed"]["RB"]["mean"], attributeMeans["Speed"][player.Position]["stddev"], player.Year)
			FootballIQGrade = util.GetLetterGrade(player.FootballIQ, attributeMeans["FootballIQ"]["ILB"]["mean"], attributeMeans["FootballIQ"][player.Position]["stddev"], player.Year)
			AgilityGrade = util.GetLetterGrade(player.Agility, attributeMeans["Agility"]["RB"]["mean"], attributeMeans["Agility"][player.Position]["stddev"], player.Year)
			CarryingGrade = util.GetLetterGrade(player.Carrying, attributeMeans["Carrying"]["RB"]["mean"], attributeMeans["Carrying"][player.Position]["stddev"], player.Year)
			CatchingGrade = util.GetLetterGrade(player.Catching, attributeMeans["Catching"]["WR"]["mean"], attributeMeans["Catching"][player.Position]["stddev"], player.Year)
			RouteRunningGrade = util.GetLetterGrade(player.RouteRunning, attributeMeans["RouteRunning"]["WR"]["mean"], attributeMeans["RouteRunning"][player.Position]["stddev"], player.Year)
			ZoneCoverageGrade = util.GetLetterGrade(player.ZoneCoverage, attributeMeans["ZoneCoverage"]["SS"]["mean"], attributeMeans["ZoneCoverage"][player.Position]["stddev"], player.Year)
			ManCoverageGrade = util.GetLetterGrade(player.ManCoverage, attributeMeans["ManCoverage"]["SS"]["mean"], attributeMeans["ManCoverage"][player.Position]["stddev"], player.Year)
			StrengthGrade = util.GetLetterGrade(player.Strength, attributeMeans["Strength"]["RB"]["mean"], attributeMeans["Strength"][player.Position]["stddev"], player.Year)
			TackleGrade = util.GetLetterGrade(player.Tackle, attributeMeans["Tackle"]["RB"]["mean"], attributeMeans["Tackle"][player.Position]["stddev"], player.Year)
			PassBlockGrade = util.GetLetterGrade(player.PassBlock, attributeMeans["PassBlock"]["RB"]["mean"], attributeMeans["PassBlock"][player.Position]["stddev"], player.Year)
			RunBlockGrade = util.GetLetterGrade(player.RunBlock, attributeMeans["RunBlock"]["RB"]["mean"], attributeMeans["RunBlock"][player.Position]["stddev"], player.Year)
			PassRushGrade = util.GetLetterGrade(player.PassRush, attributeMeans["PassRush"]["OLB"]["mean"], attributeMeans["PassRush"][player.Position]["stddev"], player.Year)
			RunDefenseGrade = util.GetLetterGrade(player.RunDefense, attributeMeans["RunDefense"]["OLB"]["mean"], attributeMeans["RunDefense"][player.Position]["stddev"], player.Year)
			ThrowPowerGrade = util.GetLetterGrade(player.ThrowPower, attributeMeans["ThrowPower"]["QB"]["mean"], attributeMeans["ThrowPower"][player.Position]["stddev"], player.Year)
			ThrowAccuracyGrade = util.GetLetterGrade(player.ThrowAccuracy, attributeMeans["ThrowAccuracy"]["QB"]["mean"], attributeMeans["ThrowAccuracy"][player.Position]["stddev"], player.Year)
			KickPowerGrade = util.GetLetterGrade(player.KickPower, attributeMeans["KickPower"]["K"]["mean"], attributeMeans["KickPower"][player.Position]["stddev"], player.Year)
			KickAccuracyGrade = util.GetLetterGrade(player.KickAccuracy, attributeMeans["KickAccuracy"]["K"]["mean"], attributeMeans["KickAccuracy"][player.Position]["stddev"], player.Year)
			PuntPowerGrade = util.GetLetterGrade(player.PuntPower, attributeMeans["PuntPower"]["P"]["mean"], attributeMeans["PuntPower"][player.Position]["stddev"], player.Year)
			PuntAccuracyGrade = util.GetLetterGrade(player.PuntAccuracy, attributeMeans["PuntAccuracy"]["P"]["mean"], attributeMeans["PuntAccuracy"][player.Position]["stddev"], player.Year)
		}
	} else {
		StaminaGrade = util.GetLetterGrade(player.Stamina, attributeMeans["Stamina"][player.Position]["mean"], attributeMeans["Stamina"][player.Position]["stddev"], player.Year)
		InjuryGrade = util.GetLetterGrade(player.Injury, attributeMeans["Injury"][player.Position]["mean"], attributeMeans["Injury"][player.Position]["stddev"], player.Year)
		SpeedGrade = util.GetLetterGrade(player.Speed, attributeMeans["Speed"][player.Position]["mean"], attributeMeans["Speed"][player.Position]["stddev"], player.Year)
		FootballIQGrade = util.GetLetterGrade(player.FootballIQ, attributeMeans["FootballIQ"][player.Position]["mean"], attributeMeans["FootballIQ"][player.Position]["stddev"], player.Year)
		AgilityGrade = util.GetLetterGrade(player.Agility, attributeMeans["Agility"][player.Position]["mean"], attributeMeans["Agility"][player.Position]["stddev"], player.Year)
		CarryingGrade = util.GetLetterGrade(player.Carrying, attributeMeans["Carrying"][player.Position]["mean"], attributeMeans["Carrying"][player.Position]["stddev"], player.Year)
		CatchingGrade = util.GetLetterGrade(player.Catching, attributeMeans["Catching"][player.Position]["mean"], attributeMeans["Catching"][player.Position]["stddev"], player.Year)
		RouteRunningGrade = util.GetLetterGrade(player.RouteRunning, attributeMeans["RouteRunning"][player.Position]["mean"], attributeMeans["RouteRunning"][player.Position]["stddev"], player.Year)
		ZoneCoverageGrade = util.GetLetterGrade(player.ZoneCoverage, attributeMeans["ZoneCoverage"][player.Position]["mean"], attributeMeans["ZoneCoverage"][player.Position]["stddev"], player.Year)
		ManCoverageGrade = util.GetLetterGrade(player.ManCoverage, attributeMeans["ManCoverage"][player.Position]["mean"], attributeMeans["ManCoverage"][player.Position]["stddev"], player.Year)
		StrengthGrade = util.GetLetterGrade(player.Strength, attributeMeans["Strength"][player.Position]["mean"], attributeMeans["Strength"][player.Position]["stddev"], player.Year)
		TackleGrade = util.GetLetterGrade(player.Tackle, attributeMeans["Tackle"][player.Position]["mean"], attributeMeans["Tackle"][player.Position]["stddev"], player.Year)
		PassBlockGrade = util.GetLetterGrade(player.PassBlock, attributeMeans["PassBlock"][player.Position]["mean"], attributeMeans["PassBlock"][player.Position]["stddev"], player.Year)
		RunBlockGrade = util.GetLetterGrade(player.RunBlock, attributeMeans["RunBlock"][player.Position]["mean"], attributeMeans["RunBlock"][player.Position]["stddev"], player.Year)
		PassRushGrade = util.GetLetterGrade(player.PassRush, attributeMeans["PassRush"][player.Position]["mean"], attributeMeans["PassRush"][player.Position]["stddev"], player.Year)
		RunDefenseGrade = util.GetLetterGrade(player.RunDefense, attributeMeans["RunDefense"][player.Position]["mean"], attributeMeans["RunDefense"][player.Position]["stddev"], player.Year)
		ThrowPowerGrade = util.GetLetterGrade(player.ThrowPower, attributeMeans["ThrowPower"][player.Position]["mean"], attributeMeans["ThrowPower"][player.Position]["stddev"], player.Year)
		ThrowAccuracyGrade = util.GetLetterGrade(player.ThrowAccuracy, attributeMeans["ThrowAccuracy"][player.Position]["mean"], attributeMeans["ThrowAccuracy"][player.Position]["stddev"], player.Year)
		KickPowerGrade = util.GetLetterGrade(player.KickPower, attributeMeans["KickPower"][player.Position]["mean"], attributeMeans["KickPower"][player.Position]["stddev"], player.Year)
		KickAccuracyGrade = util.GetLetterGrade(player.KickAccuracy, attributeMeans["KickAccuracy"][player.Position]["mean"], attributeMeans["KickAccuracy"][player.Position]["stddev"], player.Year)
		PuntPowerGrade = util.GetLetterGrade(player.PuntPower, attributeMeans["PuntPower"][player.Position]["mean"], attributeMeans["PuntPower"][player.Position]["stddev"], player.Year)
		PuntAccuracyGrade = util.GetLetterGrade(player.PuntAccuracy, attributeMeans["PuntAccuracy"][player.Position]["mean"], attributeMeans["PuntAccuracy"][player.Position]["stddev"], player.Year)
	}

	return CollegePlayerCSV{
		FirstName:          player.FirstName,
		LastName:           player.LastName,
		Position:           player.Position,
		Archetype:          player.Archetype,
		PositionTwo:        player.PositionTwo,
		ArchetypeTwo:       player.ArchetypeTwo,
		Year:               Year,
		Age:                player.Age,
		Stars:              player.Stars,
		State:              player.State,
		Height:             player.Height,
		Weight:             player.Weight,
		OverallGrade:       OverallGrade,
		StaminaGrade:       StaminaGrade,
		InjuryGrade:        InjuryGrade,
		FootballIQGrade:    FootballIQGrade,
		SpeedGrade:         SpeedGrade,
		CarryingGrade:      CarryingGrade,
		AgilityGrade:       AgilityGrade,
		CatchingGrade:      CatchingGrade,
		RouteRunningGrade:  RouteRunningGrade,
		ZoneCoverageGrade:  ZoneCoverageGrade,
		ManCoverageGrade:   ManCoverageGrade,
		StrengthGrade:      StrengthGrade,
		TackleGrade:        TackleGrade,
		PassBlockGrade:     PassBlockGrade,
		RunBlockGrade:      RunBlockGrade,
		PassRushGrade:      PassRushGrade,
		RunDefenseGrade:    RunDefenseGrade,
		ThrowPowerGrade:    ThrowPowerGrade,
		ThrowAccuracyGrade: ThrowAccuracyGrade,
		KickAccuracyGrade:  KickAccuracyGrade,
		KickPowerGrade:     KickPowerGrade,
		PuntAccuracyGrade:  PuntAccuracyGrade,
		PuntPowerGrade:     PuntPowerGrade,
		PotentialGrade:     player.PotentialGrade,
		RedshirtStatus:     RedShirtStatus,
		Team:               player.PreviousTeam,
	}
}

func MapNFLPlayerToCSVModel(player NFLPlayer) CollegePlayerCSV {

	attributeMeans := config.AttributeMeans()
	Year := util.GetNFLYear(player.Experience)
	OverallGrade := strconv.Itoa(player.Overall)
	StaminaGrade := strconv.Itoa(player.Stamina)
	InjuryGrade := strconv.Itoa(player.Injury)
	SpeedGrade := strconv.Itoa(player.Speed)
	FootballIQGrade := strconv.Itoa(player.FootballIQ)
	AgilityGrade := strconv.Itoa(player.Agility)
	CarryingGrade := strconv.Itoa(player.Carrying)
	CatchingGrade := strconv.Itoa(player.Catching)
	RouteRunningGrade := strconv.Itoa(player.RouteRunning)
	ZoneCoverageGrade := strconv.Itoa(player.ZoneCoverage)
	ManCoverageGrade := strconv.Itoa(player.ManCoverage)
	StrengthGrade := strconv.Itoa(player.Strength)
	TackleGrade := strconv.Itoa(player.Tackle)
	PassBlockGrade := strconv.Itoa(player.PassBlock)
	RunBlockGrade := strconv.Itoa(player.RunBlock)
	PassRushGrade := strconv.Itoa(player.PassRush)
	RunDefenseGrade := strconv.Itoa(player.RunDefense)
	ThrowPowerGrade := strconv.Itoa(player.ThrowPower)
	ThrowAccuracyGrade := strconv.Itoa(player.ThrowAccuracy)
	KickPowerGrade := strconv.Itoa(player.KickPower)
	KickAccuracyGrade := strconv.Itoa(player.KickAccuracy)
	PuntPowerGrade := strconv.Itoa(player.PuntPower)
	PuntAccuracyGrade := strconv.Itoa(player.PuntAccuracy)

	if player.Experience < 2 || player.ShowLetterGrade {
		OverallGrade = util.GetOverallGrade(player.Overall, int(player.Experience))
		StaminaGrade = util.GetLetterGrade(player.Stamina, attributeMeans["Stamina"][player.Position]["mean"], attributeMeans["Stamina"][player.Position]["stddev"], int(player.Experience))
		InjuryGrade = util.GetLetterGrade(player.Injury, attributeMeans["Injury"][player.Position]["mean"], attributeMeans["Injury"][player.Position]["stddev"], int(player.Experience))
		SpeedGrade = util.GetLetterGrade(player.Speed, attributeMeans["Speed"][player.Position]["mean"], attributeMeans["Speed"][player.Position]["stddev"], int(player.Experience))
		FootballIQGrade = util.GetLetterGrade(player.FootballIQ, attributeMeans["FootballIQ"][player.Position]["mean"], attributeMeans["FootballIQ"][player.Position]["stddev"], int(player.Experience))
		AgilityGrade = util.GetLetterGrade(player.Agility, attributeMeans["Agility"][player.Position]["mean"], attributeMeans["Agility"][player.Position]["stddev"], int(player.Experience))
		CarryingGrade = util.GetLetterGrade(player.Carrying, attributeMeans["Carrying"][player.Position]["mean"], attributeMeans["Carrying"][player.Position]["stddev"], int(player.Experience))
		CatchingGrade = util.GetLetterGrade(player.Catching, attributeMeans["Catching"][player.Position]["mean"], attributeMeans["Catching"][player.Position]["stddev"], int(player.Experience))
		RouteRunningGrade = util.GetLetterGrade(player.RouteRunning, attributeMeans["RouteRunning"][player.Position]["mean"], attributeMeans["RouteRunning"][player.Position]["stddev"], int(player.Experience))
		ZoneCoverageGrade = util.GetLetterGrade(player.ZoneCoverage, attributeMeans["ZoneCoverage"][player.Position]["mean"], attributeMeans["ZoneCoverage"][player.Position]["stddev"], int(player.Experience))
		ManCoverageGrade = util.GetLetterGrade(player.ManCoverage, attributeMeans["ManCoverage"][player.Position]["mean"], attributeMeans["ManCoverage"][player.Position]["stddev"], int(player.Experience))
		StrengthGrade = util.GetLetterGrade(player.Strength, attributeMeans["Strength"][player.Position]["mean"], attributeMeans["Strength"][player.Position]["stddev"], int(player.Experience))
		TackleGrade = util.GetLetterGrade(player.Tackle, attributeMeans["Tackle"][player.Position]["mean"], attributeMeans["Tackle"][player.Position]["stddev"], int(player.Experience))
		PassBlockGrade = util.GetLetterGrade(player.PassBlock, attributeMeans["PassBlock"][player.Position]["mean"], attributeMeans["PassBlock"][player.Position]["stddev"], int(player.Experience))
		RunBlockGrade = util.GetLetterGrade(player.RunBlock, attributeMeans["RunBlock"][player.Position]["mean"], attributeMeans["RunBlock"][player.Position]["stddev"], int(player.Experience))
		PassRushGrade = util.GetLetterGrade(player.PassRush, attributeMeans["PassRush"][player.Position]["mean"], attributeMeans["PassRush"][player.Position]["stddev"], int(player.Experience))
		RunDefenseGrade = util.GetLetterGrade(player.RunDefense, attributeMeans["RunDefense"][player.Position]["mean"], attributeMeans["RunDefense"][player.Position]["stddev"], int(player.Experience))
		ThrowPowerGrade = util.GetLetterGrade(player.ThrowPower, attributeMeans["ThrowPower"][player.Position]["mean"], attributeMeans["ThrowPower"][player.Position]["stddev"], int(player.Experience))
		ThrowAccuracyGrade = util.GetLetterGrade(player.ThrowAccuracy, attributeMeans["ThrowAccuracy"][player.Position]["mean"], attributeMeans["ThrowAccuracy"][player.Position]["stddev"], int(player.Experience))
		KickPowerGrade = util.GetLetterGrade(player.KickPower, attributeMeans["KickPower"][player.Position]["mean"], attributeMeans["KickPower"][player.Position]["stddev"], int(player.Experience))
		KickAccuracyGrade = util.GetLetterGrade(player.KickAccuracy, attributeMeans["KickAccuracy"][player.Position]["mean"], attributeMeans["KickAccuracy"][player.Position]["stddev"], int(player.Experience))
		PuntPowerGrade = util.GetLetterGrade(player.PuntPower, attributeMeans["PuntPower"][player.Position]["mean"], attributeMeans["PuntPower"][player.Position]["stddev"], int(player.Experience))
		PuntAccuracyGrade = util.GetLetterGrade(player.PuntAccuracy, attributeMeans["PuntAccuracy"][player.Position]["mean"], attributeMeans["PuntAccuracy"][player.Position]["stddev"], int(player.Experience))
	}

	return CollegePlayerCSV{
		FirstName:          player.FirstName,
		LastName:           player.LastName,
		Position:           player.Position,
		Archetype:          player.Archetype,
		PositionTwo:        player.PositionTwo,
		ArchetypeTwo:       player.ArchetypeTwo,
		PreviousTeam:       player.PreviousTeam,
		Team:               player.TeamAbbr,
		Year:               Year,
		Age:                player.Age,
		Stars:              player.Stars,
		HighSchool:         player.HighSchool,
		College:            player.College,
		City:               player.Hometown,
		State:              player.State,
		Height:             player.Height,
		Weight:             player.Weight,
		OverallGrade:       OverallGrade,
		StaminaGrade:       StaminaGrade,
		InjuryGrade:        InjuryGrade,
		FootballIQGrade:    FootballIQGrade,
		SpeedGrade:         SpeedGrade,
		CarryingGrade:      CarryingGrade,
		AgilityGrade:       AgilityGrade,
		CatchingGrade:      CatchingGrade,
		RouteRunningGrade:  RouteRunningGrade,
		ZoneCoverageGrade:  ZoneCoverageGrade,
		ManCoverageGrade:   ManCoverageGrade,
		StrengthGrade:      StrengthGrade,
		TackleGrade:        TackleGrade,
		PassBlockGrade:     PassBlockGrade,
		RunBlockGrade:      RunBlockGrade,
		PassRushGrade:      PassRushGrade,
		RunDefenseGrade:    RunDefenseGrade,
		ThrowPowerGrade:    ThrowPowerGrade,
		ThrowAccuracyGrade: ThrowAccuracyGrade,
		KickAccuracyGrade:  KickAccuracyGrade,
		KickPowerGrade:     KickPowerGrade,
		PuntAccuracyGrade:  PuntAccuracyGrade,
		PuntPowerGrade:     PuntPowerGrade,
		PotentialGrade:     player.PotentialGrade,
	}
}

type Croot struct {
	ID               uint
	PlayerID         int
	TeamID           int
	College          string
	FirstName        string
	LastName         string
	Position         string
	Archetype        string
	Height           int
	Weight           int
	Stars            int
	PotentialGrade   string
	Personality      string
	RecruitingBias   string
	AcademicBias     string
	WorkEthic        string
	HighSchool       string
	City             string
	State            string
	AffinityOne      string
	AffinityTwo      string
	RecruitingStatus string
	RecruitModifier  float64
	IsCustomCroot    bool
	CustomCrootFor   string
	IsSigned         bool
	OverallGrade     string
	TotalRank        float64
	LeadingTeams     []LeadingTeams
	PlayerPreferences
}

type LeadingTeams struct {
	TeamID         uint
	TeamName       string
	TeamAbbr       string
	Odds           float64
	HasScholarship bool
}

// Sorting Funcs
type ByLeadingPoints []LeadingTeams

func (rp ByLeadingPoints) Len() int      { return len(rp) }
func (rp ByLeadingPoints) Swap(i, j int) { rp[i], rp[j] = rp[j], rp[i] }
func (rp ByLeadingPoints) Less(i, j int) bool {
	return rp[i].Odds > rp[j].Odds
}

func (c *Croot) Map(r Recruit) {
	c.ID = r.ID
	c.PlayerID = r.PlayerID
	c.TeamID = r.TeamID
	c.FirstName = r.FirstName
	c.LastName = r.LastName
	c.Position = r.Position
	c.Archetype = r.Archetype
	c.Height = r.Height
	c.Weight = r.Weight
	c.Stars = r.Stars
	c.PotentialGrade = r.PotentialGrade
	c.Personality = r.Personality
	c.RecruitingBias = r.RecruitingBias
	c.AcademicBias = r.AcademicBias
	c.WorkEthic = r.WorkEthic
	c.HighSchool = r.HighSchool
	c.City = r.City
	c.State = r.State
	c.AffinityOne = r.AffinityOne
	c.AffinityTwo = r.AffinityTwo
	c.College = r.College
	c.OverallGrade = util.GetOverallGrade(r.Overall, 1)
	c.IsSigned = r.IsSigned
	c.RecruitingStatus = r.RecruitingStatus
	c.RecruitModifier = r.RecruitingModifier
	c.IsCustomCroot = r.IsCustomCroot
	c.CustomCrootFor = r.CustomCrootFor
	c.PlayerPreferences = r.PlayerPreferences

	mod := r.TopRankModifier
	if mod == 0 {
		mod = 1
	}
	c.TotalRank = (r.RivalsRank + r.ESPNRank + r.Rank247) / r.TopRankModifier

	var totalPoints float64 = 0
	var runningThreshold float64 = 0

	sortedProfiles := r.RecruitPlayerProfiles

	sort.Sort(ByPoints(sortedProfiles))

	for _, recruitProfile := range sortedProfiles {
		if recruitProfile.TeamReachedMax {
			continue
		}
		if runningThreshold == 0 {
			runningThreshold = float64(recruitProfile.TotalPoints) * 0.66
		}

		if recruitProfile.TotalPoints >= runningThreshold {
			totalPoints += float64(recruitProfile.TotalPoints)
		}

	}

	for i := 0; i < len(sortedProfiles); i++ {
		if sortedProfiles[i].TeamReachedMax || sortedProfiles[i].RemovedFromBoard {
			continue
		}
		var odds float64 = 0

		if sortedProfiles[i].TotalPoints >= runningThreshold && runningThreshold > 0 {
			odds = float64(sortedProfiles[i].TotalPoints) / totalPoints
		}
		leadingTeam := LeadingTeams{
			TeamID:         uint(r.RecruitPlayerProfiles[i].ProfileID),
			TeamAbbr:       r.RecruitPlayerProfiles[i].TeamAbbreviation,
			Odds:           odds,
			HasScholarship: r.RecruitPlayerProfiles[i].Scholarship,
		}
		c.LeadingTeams = append(c.LeadingTeams, leadingTeam)
	}
	sort.Sort(ByLeadingPoints(c.LeadingTeams))
}

type ByCrootRank []Croot

func (c ByCrootRank) Len() int      { return len(c) }
func (c ByCrootRank) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c ByCrootRank) Less(i, j int) bool {
	return c[i].TotalRank > c[j].TotalRank || c[i].Stars > c[j].Stars
}
