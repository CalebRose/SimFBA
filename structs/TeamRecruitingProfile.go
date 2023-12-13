package structs

import "github.com/jinzhu/gorm"

// RecruitingTeamProfile - The profile for a team for recruiting
type RecruitingTeamProfile struct {
	gorm.Model
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
	RecruitingClassRank       int
	CaughtCheating            bool
	IsFBS                     bool
	IsAI                      bool
	IsUserTeam                bool
	AIBehavior                string
	AIQuality                 string
	WeeksMissed               int
	BattlesWon                int
	BattlesLost               int
	AIMinThreshold            int
	AIMaxThreshold            int
	AIStarMin                 int
	AIStarMax                 int
	AIAutoOfferscholarships   bool
	OffensiveScheme           string
	DefensiveScheme           string
	Recruiter                 string                 // Username
	Recruits                  []RecruitPlayerProfile `gorm:"foreignKey:ProfileID"`
	Affinities                []ProfileAffinity      `gorm:"foreignKey:ProfileID"`
}

type SaveProfile interface{}

func (r *RecruitingTeamProfile) AssignRES(res float64) {
	r.RecruitingEfficiencyScore = res
}

func (r *RecruitingTeamProfile) ResetSpentPoints() {
	if r.SpentPoints == 0 {
		r.WeeksMissed += 1
	} else {
		r.WeeksMissed = 0
	}
	r.SpentPoints = 0
}

func (r *RecruitingTeamProfile) SubtractScholarshipsAvailable() {
	if r.ScholarshipsAvailable > 0 {
		r.ScholarshipsAvailable--
	}
}

func (r *RecruitingTeamProfile) ReallocateScholarship() {
	if r.ScholarshipsAvailable < 40 {
		r.ScholarshipsAvailable++
	}
}

func (r *RecruitingTeamProfile) ResetScholarshipCount() {
	r.ScholarshipsAvailable = 40
}

func (r *RecruitingTeamProfile) AllocateSpentPoints(points float64) {
	r.SpentPoints = points
}

func (r *RecruitingTeamProfile) AIAllocateSpentPoints(points float64) {
	r.SpentPoints += points
}

func (r *RecruitingTeamProfile) ResetWeeklyPoints(points float64) {
	r.WeeklyPoints = points
}

func (r *RecruitingTeamProfile) AddRecruitsToProfile(croots []RecruitPlayerProfile) {
	r.Recruits = croots
}

func (r *RecruitingTeamProfile) AssignRivalsRank(score float64) {
	r.RivalsScore = score
}

func (r *RecruitingTeamProfile) Assign247Rank(score float64) {
	r.Rank247Score = score
}

func (r *RecruitingTeamProfile) AssignESPNRank(score float64) {
	r.ESPNScore = score
}

func (r *RecruitingTeamProfile) AssignCompositeRank(score float64) {
	r.CompositeScore = score
}

func (r *RecruitingTeamProfile) UpdateTotalSignedRecruits(num int) {
	r.TotalCommitments = num
}

func (r *RecruitingTeamProfile) IncreaseCommitCount() {
	r.TotalCommitments++
}

func (r *RecruitingTeamProfile) ApplyCaughtCheating() {
	r.CaughtCheating = true
}

func (r *RecruitingTeamProfile) ActivateAI() {
	r.IsAI = true
	r.IsUserTeam = false
}

func (r *RecruitingTeamProfile) DeactivateAI() {
	r.IsAI = false
	r.IsUserTeam = true
}

func (r *RecruitingTeamProfile) ToggleAIBehavior() {
	r.IsAI = !r.IsAI
}

func (r *RecruitingTeamProfile) UpdateAIBehavior(isAi, autoOffer bool, starMax, starMin, min, max int, offScheme, defScheme string) {
	r.IsAI = isAi
	r.AIAutoOfferscholarships = autoOffer
	r.AIStarMax = starMax
	r.AIStarMin = starMin
	r.AIMinThreshold = min
	r.AIMaxThreshold = max
	r.OffensiveScheme = offScheme
	r.DefensiveScheme = defScheme
}

func (r *RecruitingTeamProfile) SetRecruitingClassSize(val int) {
	if val > 25 && r.IsFBS {
		r.RecruitClassSize = 25
	} else if val > 20 && !r.IsFBS {
		r.RecruitClassSize = 20
	} else {
		r.RecruitClassSize = val
	}

}

func (r *RecruitingTeamProfile) AddBattleWon() {
	r.BattlesWon += 1
}

func (r *RecruitingTeamProfile) AddBattleLost() {
	r.BattlesLost += 1
}

func (r *RecruitingTeamProfile) AssignRecruiter(name string) {
	r.Recruiter = name
}
