package structs

import "github.com/jinzhu/gorm"

type NFLPlayer struct {
	gorm.Model
	BasePlayer
	PlayerID          int
	TeamID            int
	College           string
	TeamAbbr          string
	Experience        uint
	HighSchool        string
	Hometown          string
	State             string
	IsActive          bool
	IsPracticeSquad   bool
	IsFreeAgent       bool
	IsWaived          bool
	IsOnTradeBlock    bool
	IsAcceptingOffers bool
	IsNegotiating     bool
	NegotiationRound  uint
	SigningRound      uint
	MinimumValue      float64
	PreviousTeamID    uint
	PreviousTeam      string
	DraftedTeamID     uint
	DraftedTeam       string
	DraftedRound      uint
	DraftedPick       uint
	ShowLetterGrade   bool
	Rejections        int
	Stats             []NFLPlayerStats     `gorm:"foreignKey:NFLPlayerID"`
	SeasonStats       NFLPlayerSeasonStats `gorm:"foreignKey:NFLPlayerID"`
	Contract          NFLContract          `gorm:"foreignKey:NFLPlayerID"`
	Offers            []FreeAgencyOffer    `gorm:"foreignKey:NFLPlayerID"`
	WaiverOffers      []NFLWaiverOffer     `gorm:"foreignKey:NFLPlayerID"`
	Extensions        []NFLExtensionOffer  `gorm:"foreignKey:NFLPlayerID"`
}

// Sorting Funcs
type ByTotalContract []NFLPlayer

func (rp ByTotalContract) Len() int      { return len(rp) }
func (rp ByTotalContract) Swap(i, j int) { rp[i], rp[j] = rp[j], rp[i] }
func (rp ByTotalContract) Less(i, j int) bool {
	p1 := rp[i].Contract
	p2 := rp[j].Contract
	p1Total := p1.Y1BaseSalary + p1.Y1Bonus + p1.Y2BaseSalary + p1.Y2Bonus + p1.Y3BaseSalary + p1.Y3Bonus + p1.Y4BaseSalary + p1.Y4Bonus + p1.Y5BaseSalary + p1.Y5Bonus
	p2Total := p2.Y1BaseSalary + p2.Y1Bonus + p2.Y2BaseSalary + p2.Y2Bonus + p2.Y3BaseSalary + p2.Y3Bonus + p2.Y4BaseSalary + p2.Y4Bonus + p2.Y5BaseSalary + p2.Y5Bonus
	return p1Total > p2Total
}

func (np *NFLPlayer) AssignMissingValues(pr int, aca string, fa string, per string, rec string, we string) {
	np.Progression = pr
	np.AcademicBias = aca
	np.FreeAgency = fa
	np.Personality = per
	np.WorkEthic = we
	np.RecruitingBias = rec
}

func (np *NFLPlayer) AssignMinimumValue(val float64) {
	np.MinimumValue = val
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
}

func (np *NFLPlayer) SignPlayer(TeamID int, Abbr string) {
	np.IsFreeAgent = false
	np.IsWaived = false
	np.TeamID = TeamID
	np.TeamAbbr = Abbr
	np.IsAcceptingOffers = false
	np.IsNegotiating = false
	np.IsPracticeSquad = false
}

func (np *NFLPlayer) ToggleIsPracticeSquad() {
	np.IsPracticeSquad = true
	np.IsNegotiating = false
	np.IsAcceptingOffers = true
	np.PreviousTeam = np.TeamAbbr
	np.PreviousTeamID = uint(np.TeamID)
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

func (f *NFLPlayer) DeclineOffer() {
	f.Rejections += 1
}
