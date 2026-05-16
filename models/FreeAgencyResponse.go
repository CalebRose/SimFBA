package models

import "github.com/CalebRose/SimFBA/structs"

type FreeAgencyResponse struct {
	FreeAgents    []FreeAgentResponse
	WaiverPlayers []WaiverWirePlayerResponse
	PracticeSquad []FreeAgentResponse
	TeamOffers    []structs.FreeAgencyOffer
	RosterCount   uint
}

type FreeAgentResponse struct {
	ID                uint
	PlayerID          int
	FirstName         string
	LastName          string
	Position          string
	PositionTwo       string
	Archetype         string
	ArchetypeTwo      string
	Height            int8
	Weight            int8
	Age               int8
	Overall           int8
	FootballIQ        int8
	Speed             int8
	Carrying          int8
	Agility           int8
	Catching          int8
	RouteRunning      int8
	ZoneCoverage      int8
	ManCoverage       int8
	Strength          int8
	Tackle            int8
	PassBlock         int8
	RunBlock          int8
	PassRush          int8
	RunDefense        int8
	ThrowPower        int8
	ThrowAccuracy     int8
	KickAccuracy      int8
	KickPower         int8
	PuntAccuracy      int8
	PuntPower         int8
	InjuryRating      int8
	Stamina           int8
	PotentialGrade    string
	FreeAgency        string
	Personality       string
	RecruitingBias    string
	WorkEthic         string
	AcademicBias      string
	PreviousTeamID    uint
	PreviousTeam      string
	TeamID            int
	College           string
	TeamAbbr          string
	Hometown          string
	State             string
	Experience        uint
	IsActive          bool
	IsFreeAgent       bool
	IsWaived          bool
	MinimumValue      float64
	DraftedTeam       string
	ShowLetterGrade   bool
	IsPracticeSquad   bool
	IsAcceptingOffers bool
	IsNegotiating     bool
	AAV               float64
	Shotgun           int // -1 is Under Center, 0 Balanced, 1 Shotgun
	SeasonStats       structs.NFLPlayerSeasonStats
	Offers            []structs.FreeAgencyOffer
}

type WaiverWirePlayerResponse struct {
	ID       uint
	PlayerID int
	structs.BasePlayer
	TeamID            int
	College           string
	TeamAbbr          string
	Hometown          string
	State             string
	Experience        uint
	IsActive          bool
	IsFreeAgent       bool
	IsWaived          bool
	IsAcceptingOffers bool
	IsNegotiating     bool
	MinimumValue      float64
	PreviousTeam      string
	DraftedTeam       string
	ShowLetterGrade   bool
	IsPracticeSquad   bool
	SeasonStats       structs.NFLPlayerSeasonStats
	WaiverOffers      []structs.NFLWaiverOffer
	Contract          structs.NFLContract
}
