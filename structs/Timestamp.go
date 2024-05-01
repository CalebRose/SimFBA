package structs

import "github.com/jinzhu/gorm"

type Timestamp struct {
	gorm.Model
	CollegeWeekID                 int
	CollegeWeek                   int
	CollegeSeasonID               int
	Season                        int
	NFLSeasonID                   int
	NFLWeekID                     int
	NFLWeek                       int
	CFBSpringGames                bool
	ThursdayGames                 bool
	FridayGames                   bool
	SaturdayMorning               bool
	SaturdayNoon                  bool
	SaturdayEvening               bool
	SaturdayNight                 bool
	CollegePollRan                bool
	NFLThursday                   bool
	NFLSundayNoon                 bool
	NFLSundayAfternoon            bool
	NFLSundayEvening              bool
	NFLMondayEvening              bool
	NFLTradingAllowed             bool
	NFLPreseason                  bool
	RecruitingEfficiencySynced    bool
	RecruitingSynced              bool
	GMActionsCompleted            bool
	IsOffSeason                   bool
	IsNFLOffSeason                bool
	IsRecruitingLocked            bool
	AIDepthchartsSync             bool
	AIRecruitingBoardsSynced      bool
	IsFreeAgencyLocked            bool
	IsDraftTime                   bool
	RunGames                      bool
	Y1Capspace                    float64
	Y2Capspace                    float64
	Y3Capspace                    float64
	Y4Capspace                    float64
	Y5Capspace                    float64
	DeadCapLimit                  float64
	FreeAgencyRound               uint
	RunCron                       bool
	IsTesting                     bool
	CollegeSeasonOver             bool
	NFLSeasonOver                 bool
	CrootsGenerated               bool
	ProgressedCollegePlayers      bool
	ProgressedProfessionalPlayers bool
	TransferPortalPhase           uint
	TransferPortalRound           uint
}

func (t *Timestamp) MoveUpWeekCollege() {
	t.CollegeWeekID++
	t.CollegeWeek++
	if t.CollegeWeek > 3 && t.CFBSpringGames {
		t.CollegeWeek = 0
		t.CFBSpringGames = false
		t.CollegeWeekID -= 4
	}
}

func (t *Timestamp) MoveUpWeekNFL() {
	t.NFLWeekID++
	t.NFLWeek++
	if t.NFLPreseason && t.NFLWeek > 3 {
		t.NFLWeek = 1
		t.NFLPreseason = false
	}
}

func (t *Timestamp) MoveUpFreeAgencyRound() {
	t.FreeAgencyRound++
	if t.FreeAgencyRound > 10 {
		t.FreeAgencyRound = 0
		t.IsFreeAgencyLocked = true
		t.IsDraftTime = true
	}
}

func (t *Timestamp) DraftIsOver() {
	t.IsDraftTime = false
	t.IsNFLOffSeason = false
	t.NFLPreseason = true
	t.IsOffSeason = false
	t.IsFreeAgencyLocked = false
}

func (t *Timestamp) MoveUpSeason() {
	t.CollegeSeasonID++
	t.Season++
	t.CollegeWeek = 0
	t.NFLWeek = 0
	t.NFLSeasonID++
	t.Y1Capspace = t.Y2Capspace
	t.Y2Capspace = t.Y3Capspace
	t.Y3Capspace = t.Y4Capspace
	t.Y4Capspace = t.Y5Capspace
	t.Y5Capspace += 5
}

func (t *Timestamp) ToggleRES() {
	t.RecruitingEfficiencySynced = !t.RecruitingEfficiencySynced
}

func (t *Timestamp) ToggleRecruiting() {
	t.RecruitingSynced = false
	t.IsRecruitingLocked = false
}

func (t *Timestamp) ToggleGMActions() {
	t.GMActionsCompleted = !t.GMActionsCompleted
}

func (t *Timestamp) ToggleLockRecruiting() {
	t.IsRecruitingLocked = !t.IsRecruitingLocked
}

func (t *Timestamp) ToggleFALock() {
	t.IsFreeAgencyLocked = !t.IsFreeAgencyLocked
}

func (t *Timestamp) SyncToNextWeek() {
	t.MoveUpWeekCollege()
	t.MoveUpWeekNFL()
	// Reset Toggles
	t.ThursdayGames = false
	t.FridayGames = false
	t.NFLThursday = false
	t.SaturdayNoon = false
	t.SaturdayMorning = false
	t.SaturdayEvening = false
	t.SaturdayNight = false
	t.NFLSundayNoon = false
	t.NFLSundayAfternoon = false
	t.NFLSundayEvening = false
	t.NFLMondayEvening = false
	t.AIDepthchartsSync = false
	t.AIRecruitingBoardsSynced = false
	t.RunGames = false
	// t.ToggleRES()
	t.ToggleRecruiting()
	t.ToggleGMActions()

	// Migrate game results ?
}

func (t *Timestamp) ToggleTimeSlot(ts string) {
	if ts == "Thursday Night" {
		t.ThursdayGames = true
	} else if ts == "Thursday Night Football" {
		t.NFLThursday = true
	} else if ts == "Friday Night" {
		t.FridayGames = true
	} else if ts == "Saturday Morning" {
		t.SaturdayMorning = true
	} else if ts == "Saturday Afternoon" {
		t.SaturdayNoon = true
	} else if ts == "Saturday Evening" {
		t.SaturdayEvening = true
	} else if ts == "Saturday Night" {
		t.SaturdayNight = true
	} else if ts == "Sunday Noon" {
		t.NFLSundayNoon = true
	} else if ts == "Sunday Afternoon" {
		t.NFLSundayAfternoon = true
	} else if ts == "Sunday Night Football" {
		t.NFLSundayEvening = true
	} else if ts == "Monday Night Football" {
		t.NFLMondayEvening = true
	}
}

func (t *Timestamp) ToggleRunGames() {
	t.RunGames = !t.RunGames
}

func (t *Timestamp) ToggleAIrecruitingSync() {
	t.AIRecruitingBoardsSynced = !t.AIRecruitingBoardsSynced
}

func (t *Timestamp) ToggleAIDepthCharts() {
	t.AIDepthchartsSync = !t.AIDepthchartsSync
}

func (t *Timestamp) ToggleDraftTime() {
	t.IsDraftTime = !t.IsDraftTime
	// t.IsNBAOffseason = false
}

func (t *Timestamp) TogglePollRan() {
	t.CollegePollRan = !t.CollegePollRan
}

func (t *Timestamp) EndTheCollegeSeason() {
	t.IsOffSeason = true
	t.TransferPortalPhase = 1
	t.CollegeSeasonOver = true
}

func (t *Timestamp) ClosePortal() {
	t.TransferPortalPhase = 0
}

func (t *Timestamp) EnactPromisePhase() {
	t.TransferPortalPhase = 2
}

func (t *Timestamp) EnactPortalPhase() {
	t.TransferPortalPhase = 3
}

func (t *Timestamp) IncrementTransferPortalRound() {
	t.IsRecruitingLocked = false
	if t.TransferPortalRound < 10 {
		t.TransferPortalRound += 1
	}
}

func (t *Timestamp) EndTheProfessionalSeason() {
	t.IsNFLOffSeason = true
	t.FreeAgencyRound = 1
	t.IsDraftTime = false
	t.IsFreeAgencyLocked = true
	t.NFLSeasonOver = true
}
