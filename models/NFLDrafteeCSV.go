package models

import (
	config "github.com/CalebRose/SimFBA/secrets"
	"github.com/CalebRose/SimFBA/util"
)

type NFLDrafteeCSV struct {
	PlayerID           int
	FirstName          string
	LastName           string
	Position           string
	Archetype          string
	Age                int
	Stars              int
	College            string
	HighSchool         string
	City               string
	State              string
	Height             int
	Weight             int
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
}

func MapNFLDrafteeToModel(player NFLDraftee) NFLDrafteeCSV {

	attributeMeans := config.AttributeMeans()
	OverallGrade := util.GetNFLOverallGrade(player.Overall)
	StaminaGrade := util.GetLetterGrade(player.Stamina, attributeMeans["Stamina"][player.Position]["mean"], attributeMeans["Stamina"][player.Position]["stddev"], 4)
	InjuryGrade := util.GetLetterGrade(player.Injury, attributeMeans["Injury"][player.Position]["mean"], attributeMeans["Injury"][player.Position]["stddev"], 4)
	SpeedGrade := util.GetLetterGrade(player.Speed, attributeMeans["Speed"][player.Position]["mean"], attributeMeans["Speed"][player.Position]["stddev"], 4)
	FootballIQGrade := util.GetLetterGrade(player.FootballIQ, attributeMeans["FootballIQ"][player.Position]["mean"], attributeMeans["FootballIQ"][player.Position]["stddev"], 4)
	AgilityGrade := util.GetLetterGrade(player.Agility, attributeMeans["Agility"][player.Position]["mean"], attributeMeans["Agility"][player.Position]["stddev"], 4)
	CarryingGrade := util.GetLetterGrade(player.Carrying, attributeMeans["Carrying"][player.Position]["mean"], attributeMeans["Carrying"][player.Position]["stddev"], 4)
	CatchingGrade := util.GetLetterGrade(player.Catching, attributeMeans["Catching"][player.Position]["mean"], attributeMeans["Catching"][player.Position]["stddev"], 4)
	RouteRunningGrade := util.GetLetterGrade(player.RouteRunning, attributeMeans["RouteRunning"][player.Position]["mean"], attributeMeans["RouteRunning"][player.Position]["stddev"], 4)
	ZoneCoverageGrade := util.GetLetterGrade(player.ZoneCoverage, attributeMeans["ZoneCoverage"][player.Position]["mean"], attributeMeans["ZoneCoverage"][player.Position]["stddev"], 4)
	ManCoverageGrade := util.GetLetterGrade(player.ManCoverage, attributeMeans["ManCoverage"][player.Position]["mean"], attributeMeans["ManCoverage"][player.Position]["stddev"], 4)
	StrengthGrade := util.GetLetterGrade(player.Strength, attributeMeans["Strength"][player.Position]["mean"], attributeMeans["Strength"][player.Position]["stddev"], 4)
	TackleGrade := util.GetLetterGrade(player.Tackle, attributeMeans["Tackle"][player.Position]["mean"], attributeMeans["Tackle"][player.Position]["stddev"], 4)
	PassBlockGrade := util.GetLetterGrade(player.PassBlock, attributeMeans["PassBlock"][player.Position]["mean"], attributeMeans["PassBlock"][player.Position]["stddev"], 4)
	RunBlockGrade := util.GetLetterGrade(player.RunBlock, attributeMeans["RunBlock"][player.Position]["mean"], attributeMeans["RunBlock"][player.Position]["stddev"], 4)
	PassRushGrade := util.GetLetterGrade(player.PassRush, attributeMeans["PassRush"][player.Position]["mean"], attributeMeans["PassRush"][player.Position]["stddev"], 4)
	RunDefenseGrade := util.GetLetterGrade(player.RunDefense, attributeMeans["RunDefense"][player.Position]["mean"], attributeMeans["RunDefense"][player.Position]["stddev"], 4)
	ThrowPowerGrade := util.GetLetterGrade(player.ThrowPower, attributeMeans["ThrowPower"][player.Position]["mean"], attributeMeans["ThrowPower"][player.Position]["stddev"], 4)
	ThrowAccuracyGrade := util.GetLetterGrade(player.ThrowAccuracy, attributeMeans["ThrowAccuracy"][player.Position]["mean"], attributeMeans["ThrowAccuracy"][player.Position]["stddev"], 4)
	KickPowerGrade := util.GetLetterGrade(player.KickPower, attributeMeans["KickPower"][player.Position]["mean"], attributeMeans["KickPower"][player.Position]["stddev"], 4)
	KickAccuracyGrade := util.GetLetterGrade(player.KickAccuracy, attributeMeans["KickAccuracy"][player.Position]["mean"], attributeMeans["KickAccuracy"][player.Position]["stddev"], 4)
	PuntPowerGrade := util.GetLetterGrade(player.PuntPower, attributeMeans["PuntPower"][player.Position]["mean"], attributeMeans["PuntPower"][player.Position]["stddev"], 4)
	PuntAccuracyGrade := util.GetLetterGrade(player.PuntAccuracy, attributeMeans["PuntAccuracy"][player.Position]["mean"], attributeMeans["PuntAccuracy"][player.Position]["stddev"], 4)

	return NFLDrafteeCSV{
		PlayerID:           player.PlayerID,
		FirstName:          player.FirstName,
		LastName:           player.LastName,
		Position:           player.Position,
		Archetype:          player.Archetype,
		Age:                player.Age,
		Stars:              player.Stars,
		College:            player.College,
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
	}
}
