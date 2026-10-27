package structs

import "github.com/jinzhu/gorm"

type Timestamp struct {
	gorm.Model
	Phase                         uint
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
	PreDraftPhase                 bool
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

func (t *Timestamp) MoveUpPhase() {
	t.Phase++
}

func (t *Timestamp) MoveUpWeekCollege() {
	t.CollegeWeekID++
	t.CollegeWeek++
	if t.CollegeWeek > 3 && t.CFBSpringGames {
		t.CollegeWeek = 0
		t.CollegeWeekID -= 4
		t.CFBSpringGames = false
	}
	if t.CollegeWeek > 20 {
		t.CollegeSeasonOver = true
	}
}

func (t *Timestamp) MoveUpWeekNFL() {
	t.NFLWeekID++
	t.NFLWeek++
	if t.NFLPreseason && t.NFLWeek > 3 {
		t.NFLWeek = 0
		t.NFLWeekID -= 4
		t.NFLPreseason = false
	}
	if t.NFLWeek > 22 {
		t.NFLSeasonOver = true
	}
}

func (t *Timestamp) MoveUpFreeAgencyRound() {
	t.FreeAgencyRound++
	if t.FreeAgencyRound > 24 {
		t.FreeAgencyRound = 0
		t.IsFreeAgencyLocked = true
		t.IsDraftTime = true
	}
}

func (t *Timestamp) DraftIsOver() {
	t.IsDraftTime = false
	t.IsNFLOffSeason = false
	t.NFLPreseason = true
	t.CFBSpringGames = true
	t.IsOffSeason = false
	t.IsFreeAgencyLocked = false
}

func (t *Timestamp) MoveUpSeason() {
	t.Phase = 1
	t.CollegeSeasonID++
	t.Season++
	t.CollegeWeek = 0
	t.NFLWeek = 0
	t.NFLSeasonID++
	t.CollegeWeekID = (((t.CollegeSeasonID + 2020) - 2000) * 100)
	t.NFLWeekID = (((t.CollegeSeasonID + 2020) - 2000) * 100)
	t.Y1Capspace = t.Y2Capspace
	t.Y2Capspace = t.Y3Capspace
	t.Y3Capspace = t.Y4Capspace
	t.Y4Capspace = t.Y5Capspace
	t.Y5Capspace += 5
	t.NFLSeasonOver = false
	t.CollegeSeasonOver = false
	t.ProgressedCollegePlayers = false
	t.ProgressedProfessionalPlayers = false
	t.FreeAgencyRound = 1
	t.PreDraftPhase = true
	t.IsNFLOffSeason = true
	t.IsOffSeason = true
	t.CrootsGenerated = false
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

func (t *Timestamp) AdminSetWeek(collegeWeek int, nflWeek int) {
	t.CollegeWeek = collegeWeek
	t.NFLWeek = nflWeek
	t.CollegeWeekID = (((t.CollegeSeasonID + 2020) - 2000) * 100) + t.CollegeWeek
	t.NFLWeekID = (((t.CollegeSeasonID + 2020) - 2000) * 100) + t.NFLWeek
}

func (t *Timestamp) AdminToggleTimeSlot(ts string) {
	switch ts {
	case "Thursday Night":
		t.ThursdayGames = !t.ThursdayGames
	case "Thursday Night Football":
		t.NFLThursday = !t.NFLThursday
	case "Friday Night":
		t.FridayGames = !t.FridayGames
	case "Saturday Morning":
		t.SaturdayMorning = !t.SaturdayMorning
	case "Saturday Afternoon":
		t.SaturdayNoon = !t.SaturdayNoon
	case "Saturday Evening":
		t.SaturdayEvening = !t.SaturdayEvening
	case "Saturday Night":
		t.SaturdayNight = !t.SaturdayNight
	case "Sunday Noon":
		t.NFLSundayNoon = !t.NFLSundayNoon
	case "Sunday Afternoon":
		t.NFLSundayAfternoon = !t.NFLSundayAfternoon
	case "Sunday Night Football":
		t.NFLSundayEvening = !t.NFLSundayEvening
	case "Monday Night Football":
		t.NFLMondayEvening = !t.NFLMondayEvening
	}
}

func (t *Timestamp) ToggleTimeSlot(ts string) {
	switch ts {
	case "Thursday Night":
		t.ThursdayGames = true
	case "Thursday Night Football":
		t.NFLThursday = true
	case "Friday Night":
		t.FridayGames = true
	case "Saturday Morning":
		t.SaturdayMorning = true
	case "Saturday Afternoon":
		t.SaturdayNoon = true
	case "Saturday Evening":
		t.SaturdayEvening = true
	case "Saturday Night":
		t.SaturdayNight = true
	case "Sunday Noon":
		t.NFLSundayNoon = true
	case "Sunday Afternoon":
		t.NFLSundayAfternoon = true
	case "Sunday Night Football":
		t.NFLSundayEvening = true
	case "Monday Night Football":
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
	if !t.IsDraftTime {
		t.PreDraftPhase = false
	}
	// t.IsNBAOffseason = false
}

func (t *Timestamp) TogglePollRan() {
	t.CollegePollRan = !t.CollegePollRan
}

func (t *Timestamp) EndTheCollegeSeason() {
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
	t.TransferPortalRound = 1
}

func (t *Timestamp) IncrementTransferPortalRound() {
	t.IsRecruitingLocked = false
	if t.TransferPortalRound <= 10 {
		t.TransferPortalRound += 1
	} else {
		t.TransferPortalPhase = 0
		t.TransferPortalRound = 0
	}
}

func (t *Timestamp) ToggleGeneratedCroots() {
	t.CrootsGenerated = !t.CrootsGenerated
}

func (t *Timestamp) ToggleCollegeProgression() {
	t.ProgressedCollegePlayers = true
	t.TransferPortalPhase = 1
}

func (t *Timestamp) ToggleProfessionalProgression() {
	t.ProgressedProfessionalPlayers = true
	t.IsFreeAgencyLocked = false
	t.IsDraftTime = true
}

func (t *Timestamp) GetNFLCurrentGameType() (int, string) {
	if t.NFLPreseason {
		return 1, "1"
	}
	if t.NFLWeek > 18 {
		return 3, "3"
	}
	return 2, "2"
}

func (t *Timestamp) GetCFBCurrentGameType() (int, string) {
	if t.CFBSpringGames {
		return 1, "1"
	}
	if t.CollegeWeek > 14 {
		return 3, "3"
	}
	return 2, "2"
}

func (t *Timestamp) ReadyForSeasonMigration() bool {
	return (t.NFLSeasonOver && t.CollegeSeasonOver && !t.IsNFLOffSeason && !t.IsOffSeason && t.ProgressedCollegePlayers && t.ProgressedProfessionalPlayers) || t.Phase > 34
}

func (t *Timestamp) RunRedshirtsAlgorithm() bool {
	return t.Phase == 14 || t.CollegeWeek == 4
}

func (t *Timestamp) RunGenerateCollegeRankings() bool {
	return t.CollegeWeek >= 2 && !t.CFBSpringGames && t.CollegeWeek < 21
}

func (t *Timestamp) RunSyncPromises() bool {
	return (t.CollegeSeasonOver && t.CollegeWeek == 21) || t.Phase == 31
}

func (t *Timestamp) RunExtensionAndTagValues() bool {
	return t.Phase == 19 || t.NFLWeek == 9
}

func (t *Timestamp) RunCollegeProgressions() bool {
	return (t.CollegeSeasonOver && !t.ProgressedCollegePlayers) || (t.Phase == 31 && !t.ProgressedCollegePlayers)
}

func (t *Timestamp) RunProfessionalProgressions() bool {
	return (t.NFLSeasonOver && !t.ProgressedProfessionalPlayers) || (t.Phase == 33 && !t.ProgressedProfessionalPlayers)
}

func (t *Timestamp) RunGenerateWeatherFunction() bool {
	return t.Phase == 1
}

func (t *Timestamp) RunGenerateOOCSchedule() bool {
	return t.Phase == 6
}

func (t *Timestamp) RunCollegeRankings() bool {
	return t.CollegeWeek >= 2 && !t.CFBSpringGames && t.CollegeWeek < 21
}

func (t *Timestamp) RunUDFAs() bool {
	return t.Phase == 4
}

func (t *Timestamp) RunTrainingCamps() bool {
	return t.Phase == 5
}

func (t *Timestamp) RunGenerateWalkonAlgorithm() bool {
	return t.Phase == 30 || t.CollegeWeek == 20
}

func (t *Timestamp) RunTransferIntentions() bool {
	return (t.RunCron && t.CollegeSeasonOver && t.TransferPortalPhase == 1) || t.Phase == 31
}

func (t *Timestamp) ReadyToActivatePortal() bool {
	return (t.RunCron && t.CollegeSeasonOver && t.TransferPortalPhase == 2) || t.Phase == 32
}

func (t *Timestamp) IsCollegeRegularSeason() bool {
	return t.RunCron && !t.IsOffSeason && !t.CollegeSeasonOver
}

func (t *Timestamp) IsNFLRegularSeason() bool {
	return t.RunCron && !t.IsNFLOffSeason && !t.NFLSeasonOver
}
