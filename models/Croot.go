package models

import (
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
)

type Croot struct {
	ID             uint
	PlayerID       int
	TeamID         int
	College        string
	FirstName      string
	LastName       string
	Position       string
	Archetype      string
	Height         int
	Weight         int
	Stars          int
	PotentialGrade string
	Personality    string
	RecruitingBias string
	AcademicBias   string
	WorkEthic      string
	HighSchool     string
	City           string
	State          string
	AffinityOne    string
	AffinityTwo    string
	IsSigned       bool
	OverallGrade   string
	LeadingTeams   []LeadingTeams
}

type LeadingTeams struct {
	TeamName string
	TeamAbbr string
	Odds     float64
}

func (c *Croot) Map(r structs.Recruit) {
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
	c.OverallGrade = util.GetOverallGrade(r.Overall)
	c.IsSigned = r.IsSigned

	var totalPoints float64 = 0
	var runningThreshold float64 = 0

	for idx, recruitProfile := range r.RecruitPlayerProfiles {
		if idx == 0 {
			runningThreshold = float64(recruitProfile.TotalPoints) / 2
		}
		if recruitProfile.TotalPoints >= runningThreshold {
			totalPoints += float64(recruitProfile.TotalPoints)
		}
	}

	for i := 0; i < len(r.RecruitPlayerProfiles); i++ {
		var odds float64 = 0
		if r.RecruitPlayerProfiles[i].TotalPoints >= runningThreshold {
			odds = float64(r.RecruitPlayerProfiles[i].TotalPoints) / totalPoints
		}
		leadingTeam := LeadingTeams{
			TeamAbbr: r.RecruitPlayerProfiles[i].TeamAbbreviation,
			Odds:     odds,
		}
		c.LeadingTeams = append(c.LeadingTeams, leadingTeam)
	}
}
