package models

import "github.com/CalebRose/SimFBA/structs"

type SimTeamBoardResponse struct {
	ID                        uint
	TeamID                    int
	Team                      string
	TeamAbbreviation          string
	State                     string
	ScholarshipsAvailable     int
	WeeklyPoints              float64
	SpentPoints               float64
	TotalCommitments          int
	RecruitClassSize          int
	BaseEfficiencyScore       float64
	RecruitingEfficiencyScore float64
	PreviousOverallWinPer     float64
	PreviousConferenceWinPer  float64
	CurrentOverallWinPer      float64
	CurrentConferenceWinPer   float64
	ESPNScore                 float64
	RivalsScore               float64
	Rank247Score              float64
	CompositeScore            float64
	IsAI                      bool
	IsUserTeam                bool
	BattlesWon                int
	BattlesLost               int
	AIMinThreshold            int
	AIMaxThreshold            int
	AIStarMin                 int
	AIStarMax                 int
	AIAutoOfferscholarships   bool
	OffensiveScheme           string
	DefensiveScheme           string
	Recruiter                 string
	RecruitingClassRank       int
	Recruits                  []CrootProfile
	Affinities                []structs.ProfileAffinity `gorm:"foreignKey:ProfileID"`
}

func (stbr *SimTeamBoardResponse) Map(rtp structs.RecruitingTeamProfile, c []CrootProfile) {
	stbr.ID = rtp.ID
	stbr.TeamID = rtp.TeamID
	stbr.Team = rtp.Team
	stbr.IsAI = rtp.IsAI
	stbr.TeamAbbreviation = rtp.TeamAbbreviation
	stbr.State = rtp.State
	stbr.ScholarshipsAvailable = rtp.ScholarshipsAvailable
	stbr.WeeklyPoints = rtp.WeeklyPoints
	stbr.SpentPoints = rtp.SpentPoints
	stbr.TotalCommitments = rtp.TotalCommitments
	stbr.BaseEfficiencyScore = rtp.BaseEfficiencyScore
	stbr.RecruitingEfficiencyScore = rtp.RecruitingEfficiencyScore
	stbr.ESPNScore = rtp.ESPNScore
	stbr.RivalsScore = rtp.RivalsScore
	stbr.Rank247Score = rtp.Rank247Score
	stbr.CompositeScore = rtp.CompositeScore
	stbr.RecruitingClassRank = rtp.RecruitingClassRank
	stbr.Affinities = rtp.Affinities
	stbr.Recruits = c
	stbr.RecruitClassSize = rtp.RecruitClassSize
	stbr.IsUserTeam = rtp.IsUserTeam
	stbr.BattlesWon = rtp.BattlesWon
	stbr.BattlesLost = rtp.BattlesLost
	stbr.AIMinThreshold = rtp.AIMinThreshold
	stbr.AIMaxThreshold = rtp.AIMaxThreshold
	stbr.AIStarMin = rtp.AIStarMin
	stbr.AIStarMax = rtp.AIStarMax
	stbr.OffensiveScheme = rtp.OffensiveScheme
	stbr.DefensiveScheme = rtp.DefensiveScheme
	stbr.Recruiter = rtp.Recruiter
}
