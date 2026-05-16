package structs

import "github.com/jinzhu/gorm"

type NFLPlayer struct {
	gorm.Model
	BasePlayer
	PlayerID             int
	TeamID               int
	CollegeID            uint
	College              string
	TeamAbbr             string
	Experience           uint
	HighSchool           string
	Hometown             string
	State                string
	IsActive             bool
	IsPracticeSquad      bool
	IsFreeAgent          bool
	IsWaived             bool
	IsOnTradeBlock       bool
	IsAcceptingOffers    bool
	IsNegotiating        bool
	NegotiationRound     uint
	SigningRound         uint
	MinimumValue         float64
	OriginalMinimumValue float64
	AAV                  float64
	OriginalAAV          float64
	DraftedTeamID        uint
	DraftedTeam          string
	DraftedRound         uint
	DraftPickID          uint
	DraftedPick          uint
	ShowLetterGrade      bool
	HasProgressed        bool
	Rejections           int
	ProBowls             uint8
	TagType              uint8                // 0 == Basic, 1 == Franchise, 2 == Transition, 3 == Playtime
	Stats                []NFLPlayerStats     `gorm:"foreignKey:NFLPlayerID"`
	SeasonStats          NFLPlayerSeasonStats `gorm:"foreignKey:NFLPlayerID"`
	Contract             NFLContract          `gorm:"foreignKey:NFLPlayerID"`
	Offers               []FreeAgencyOffer    `gorm:"foreignKey:NFLPlayerID"`
	WaiverOffers         []NFLWaiverOffer     `gorm:"foreignKey:NFLPlayerID"`
	Extensions           []NFLExtensionOffer  `gorm:"foreignKey:NFLPlayerID"`
}

// Sorting Funcs
type ByTotalContract []NFLPlayer

func (rp ByTotalContract) Len() int      { return len(rp) }
func (rp ByTotalContract) Swap(i, j int) { rp[i], rp[j] = rp[j], rp[i] }
func (rp ByTotalContract) Less(i, j int) bool {
	p1 := rp[i].Contract
	p2 := rp[j].Contract
	p1Total := p1.Y1BaseSalary + p1.Y1Bonus
	p2Total := p2.Y1BaseSalary + p2.Y1Bonus
	return p1Total > p2Total
}

func (np *NFLPlayer) SetID(val uint) {
	np.ID = val
}

func (np *NFLPlayer) AssignMissingValues(pr int, aca string, fa string, per string, rec string, we string) {
	np.Progression = int8(pr)
	np.AcademicBias = aca
	np.FreeAgency = fa
	np.Personality = per
	np.WorkEthic = we
	np.RecruitingBias = rec
}

func (np *NFLPlayer) AssignMinimumValue(val, aav float64) {
	np.MinimumValue = val
	np.AAV = aav
}

func (np *NFLPlayer) ShowRealAttributeValue() {
	np.ShowLetterGrade = false
}

func (np *NFLPlayer) ToggleIsFreeAgent() {
	np.PreviousTeamID = uint(np.TeamID)
	np.PreviousTeam = np.TeamAbbr
	np.IsFreeAgent = true
	np.TeamID = 0
	np.TeamAbbr = ""
	np.IsAcceptingOffers = true
	np.IsNegotiating = false
	np.IsOnTradeBlock = false
	np.IsPracticeSquad = false
	np.Rejections = 0
	np.IsWaived = false
	np.Rejections = 0
}

func (np *NFLPlayer) SignPlayer(TeamID int, Abbr string) {
	np.IsFreeAgent = false
	np.IsWaived = false
	np.TeamID = TeamID
	np.TeamAbbr = Abbr
	np.IsAcceptingOffers = false
	np.IsNegotiating = false
	np.IsPracticeSquad = false
	np.Rejections = 0
}

func (np *NFLPlayer) ToggleIsPracticeSquad() {
	np.IsPracticeSquad = !np.IsPracticeSquad
	np.IsNegotiating = false
	if np.IsPracticeSquad {
		np.IsAcceptingOffers = true
	}
}

func (np *NFLPlayer) ToggleTradeBlock() {
	np.IsOnTradeBlock = !np.IsOnTradeBlock
}

func (np *NFLPlayer) RemoveFromTradeBlock() {
	np.IsOnTradeBlock = false
}

func (np *NFLPlayer) WaivePlayer() {
	np.PreviousTeamID = uint(np.TeamID)
	np.PreviousTeam = np.TeamAbbr
	np.TeamID = 0
	np.TeamAbbr = ""
	np.RemoveFromTradeBlock()
	np.IsWaived = true
}

func (np *NFLPlayer) ConvertWaivedPlayerToFA() {
	np.IsWaived = false
	np.IsFreeAgent = true
	np.IsAcceptingOffers = true
}

func (np *NFLPlayer) ToggleIsNegotiating() {
	np.IsNegotiating = true
	np.IsAcceptingOffers = false
}

func (np *NFLPlayer) WaitUntilAfterDraft() {
	np.IsNegotiating = false
	np.IsAcceptingOffers = false
}

func (np *NFLPlayer) AssignWorkEthic(we string) {
	np.WorkEthic = we
}

func (np *NFLPlayer) AssignPersonality(we string) {
	np.Personality = we
}

func (np *NFLPlayer) AssignFreeAgency(we string) {
	np.FreeAgency = we
}

func (np *NFLPlayer) AssignFAPreferences(negotiation uint, signing uint) {
	np.NegotiationRound = negotiation
	np.SigningRound = signing
}

func (np *NFLPlayer) TradePlayer(id uint, team string) {
	np.PreviousTeam = np.TeamAbbr
	np.PreviousTeamID = uint(np.TeamID)
	np.TeamID = int(id)
	np.TeamAbbr = team
	np.IsOnTradeBlock = false
}

func (f *NFLPlayer) DeclineOffer(week int) {
	f.Rejections += 1
	if week >= 23 {
		f.Rejections += 2
	}
}

func (f *NFLPlayer) DecreaseMinimumValue() {
	if f.Age < 24 || f.MinimumValue < 1 {
		return
	}

	// Decrease by 1% each time, but never go below 70% of the original minimum value or 70% overall
	f.MinimumValue = float64(f.MinimumValue) * 0.99
	if f.MinimumValue < 0.7 {
		f.MinimumValue = 0.7
	}
	// If the minimum value is less than 70% of the original minimum value, set it to 70% of the original minimum value
	if f.MinimumValue < f.OriginalMinimumValue*0.7 {
		f.MinimumValue = f.OriginalMinimumValue * 0.7
	}
	f.AAV = float64(f.AAV) * 0.99
	if f.AAV < f.OriginalAAV*0.7 {
		f.AAV = f.OriginalAAV * 0.7
	}
	if f.AAV < 0.7 {
		f.AAV = 0.7
	}
}

func (f *NFLPlayer) ToggleHasProgressed() {
	f.HasProgressed = true
}

func (np *NFLPlayer) Progress(attr CollegePlayerProgressions) {
	np.Age++
	np.Experience++
	np.Agility = int8(AddAttribute(int(attr.Agility)))
	np.Speed = int8(AddAttribute(int(attr.Speed)))
	np.FootballIQ = int8(AddAttribute(int(attr.FootballIQ)))
	np.Carrying = int8(AddAttribute(int(attr.Carrying)))
	np.Catching = int8(AddAttribute(int(attr.Catching)))
	np.RouteRunning = int8(AddAttribute(int(attr.RouteRunning)))
	np.PassBlock = int8(AddAttribute(int(attr.PassBlock)))
	np.RunBlock = int8(AddAttribute(int(attr.RunBlock)))
	np.PassRush = int8(AddAttribute(int(attr.PassRush)))
	np.RunDefense = int8(AddAttribute(int(attr.RunDefense)))
	np.Tackle = int8(AddAttribute(int(attr.Tackle)))
	np.Strength = int8(AddAttribute(int(attr.Strength)))
	np.ManCoverage = int8(AddAttribute(int(attr.ManCoverage)))
	np.ZoneCoverage = int8(AddAttribute(int(attr.ZoneCoverage)))
	np.KickAccuracy = int8(AddAttribute(int(attr.KickAccuracy)))
	np.KickPower = int8(AddAttribute(int(attr.KickPower)))
	np.PuntAccuracy = int8(AddAttribute(int(attr.PuntAccuracy)))
	np.PuntPower = int8(AddAttribute(int(attr.PuntPower)))
	np.ThrowAccuracy = int8(AddAttribute(int(attr.ThrowAccuracy)))
	np.ThrowPower = int8(AddAttribute(int(attr.ThrowPower)))
	np.HasProgressed = true
	np.ShowLetterGrade = false
	np.IsInjured = false
	np.WeeksOfRecovery = 0
	np.InjuryType = ""
	np.InjuryReserve = false
	if len(attr.PotentialGrade) > 0 {
		np.PotentialGrade = attr.PotentialGrade
	}
	np.Rejections = 0
}

func (f *NFLPlayer) MapSeasonStats(seasonStats NFLPlayerSeasonStats) {
	f.SeasonStats = seasonStats
}

func (f *NFLPlayer) AddTagType(tagType uint8) {
	f.TagType = tagType
}

func (np *NFLPlayer) ApplyTrainingCampInfo(attr CollegePlayerProgressions) {
	np.Agility = int8(AddAttribute(int(attr.Agility)))
	np.Speed = int8(AddAttribute(int(attr.Speed)))
	np.FootballIQ = int8(AddAttribute(int(attr.FootballIQ)))
	np.Carrying = int8(AddAttribute(int(attr.Carrying)))
	np.Catching = int8(AddAttribute(int(attr.Catching)))
	np.RouteRunning = int8(AddAttribute(int(attr.RouteRunning)))
	np.PassBlock = int8(AddAttribute(int(attr.PassBlock)))
	np.RunBlock = int8(AddAttribute(int(attr.RunBlock)))
	np.PassRush = int8(AddAttribute(int(attr.PassRush)))
	np.RunDefense = int8(AddAttribute(int(attr.RunDefense)))
	np.Tackle = int8(AddAttribute(int(attr.Tackle)))
	np.Strength = int8(AddAttribute(int(attr.Strength)))
	np.ManCoverage = int8(AddAttribute(int(attr.ManCoverage)))
	np.ZoneCoverage = int8(AddAttribute(int(attr.ZoneCoverage)))
	np.KickAccuracy = int8(AddAttribute(int(attr.KickAccuracy)))
	np.KickPower = int8(AddAttribute(int(attr.KickPower)))
	np.PuntAccuracy = int8(AddAttribute(int(attr.PuntAccuracy)))
	np.PuntPower = int8(AddAttribute(int(attr.PuntPower)))
	np.ThrowAccuracy = int8(AddAttribute(int(attr.ThrowAccuracy)))
	np.ThrowPower = int8(AddAttribute(int(attr.ThrowPower)))
	np.IsInjured = attr.WeeksOfRecovery > 0
	np.WeeksOfRecovery = uint(attr.WeeksOfRecovery)
	np.InjuryName = attr.InjuryText
	np.InjuryType = attr.InjuryText
	np.InjuryReserve = false
}

func (cp *NFLPlayer) AddSeasonStats(seasonStats NFLPlayerSeasonStats) {
	cp.SeasonStats = seasonStats
}

func AddAttribute(attr int) int {
	newVal := attr
	if newVal > 99 {
		newVal = 99
	} else if newVal < 1 {
		newVal = 1
	}
	return newVal
}
