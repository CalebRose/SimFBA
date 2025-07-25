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
	Height            int
	Weight            int
	Age               int
	Overall           int
	FootballIQ        int
	Speed             int
	Carrying          int
	Agility           int
	Catching          int
	RouteRunning      int
	ZoneCoverage      int
	ManCoverage       int
	Strength          int
	Tackle            int
	PassBlock         int
	RunBlock          int
	PassRush          int
	RunDefense        int
	ThrowPower        int
	ThrowAccuracy     int
	KickAccuracy      int
	KickPower         int
	PuntAccuracy      int
	PuntPower         int
	InjuryRating      int
	Stamina           int
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
