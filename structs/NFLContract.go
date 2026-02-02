package structs

import "github.com/jinzhu/gorm"

type NFLContract struct {
	gorm.Model
	PlayerID        int
	NFLPlayerID     int
	TeamID          uint
	Team            string
	OriginalTeamID  uint
	OriginalTeam    string
	ContractLength  int
	Y1BaseSalary    float64
	Y1Bonus         float64
	Y2BaseSalary    float64
	Y2Bonus         float64
	Y3BaseSalary    float64
	Y3Bonus         float64
	Y4BaseSalary    float64
	Y4Bonus         float64
	Y5BaseSalary    float64
	Y5Bonus         float64
	BonusPercentage float64
	ContractType    string // Pro Bowl, Starter, Veteran, New ?
	ContractValue   float64
	SigningValue    float64
	IsActive        bool
	IsComplete      bool
	IsExtended      bool
	HasProgressed   bool
	PlayerRetired   bool
	TagType         uint8
	IsTagged        bool
	IsCut           bool
}

func (c *NFLContract) DeactivateContract() {
	c.IsActive = false
}

func (c *NFLContract) CutContract() {
	c.IsActive = false
	c.IsCut = true
}

func (c *NFLContract) ReassignTeam(TeamID uint, Team string) {
	c.TeamID = TeamID
	c.Team = Team
}

func (c *NFLContract) TradePlayer(TeamID uint, Team string, percentage float64) {
	c.TeamID = TeamID
	c.Team = Team
	c.Y1BaseSalary = c.Y1BaseSalary * percentage
	c.Y1Bonus = 0
	c.Y2Bonus = 0
	c.Y3Bonus = 0
	c.Y4Bonus = 0
	c.Y5Bonus = 0
}

func (c *NFLContract) ProgressContract() {
	c.Y1BaseSalary = c.Y2BaseSalary
	c.Y1Bonus = c.Y2Bonus
	c.Y2BaseSalary = c.Y3BaseSalary
	c.Y2Bonus = c.Y3Bonus
	c.Y3BaseSalary = c.Y4BaseSalary
	c.Y3Bonus = c.Y4Bonus
	c.Y4BaseSalary = c.Y5BaseSalary
	c.Y4Bonus = c.Y5Bonus
	c.Y5BaseSalary = 0
	c.Y5Bonus = 0
	c.ContractLength -= 1
	c.CalculateContract()
	c.HasProgressed = true

	if c.Y1BaseSalary == 0 && c.Y1Bonus == 0 && c.ContractLength == 0 {
		c.IsComplete = true
		c.DeactivateContract()
	} else if c.Y1BaseSalary == 0 && c.Y1Bonus == 0 && c.ContractLength > 0 {
		c.Y1Bonus = 0.7
	}
}

func (c *NFLContract) CalculateContract() {
	// Calculate Value
	y1SalaryVal := c.Y1BaseSalary * 0.8
	y1BonusVal := c.Y1Bonus * 1
	y2SalaryVal := c.Y2BaseSalary * 0.4
	y2BonusVal := c.Y2Bonus * 0.9
	y3SalaryVal := c.Y3BaseSalary * 0.2
	y3BonusVal := c.Y3Bonus * 0.8
	y4SalaryVal := c.Y4BaseSalary * 0.1
	y4BonusVal := c.Y4Bonus * 0.7
	y5SalaryVal := c.Y5BaseSalary * 0.05
	y5BonusVal := c.Y5Bonus * 0.6
	c.ContractValue = y1SalaryVal + y1BonusVal + y2SalaryVal + y2BonusVal + y3SalaryVal + y3BonusVal + y4SalaryVal + y4BonusVal + y5SalaryVal + y5BonusVal
}

func (c *NFLContract) MapExtension(e NFLExtensionOffer) {
	c.ContractLength = e.ContractLength
	c.Y1BaseSalary = e.Y1BaseSalary
	c.Y1Bonus = e.Y1Bonus
	c.Y2BaseSalary = e.Y2BaseSalary
	c.Y2Bonus = e.Y2Bonus
	c.Y3BaseSalary = e.Y3BaseSalary
	c.Y3Bonus = e.Y3Bonus
	c.Y4BaseSalary = e.Y4BaseSalary
	c.Y4Bonus = e.Y4Bonus
	c.Y5BaseSalary = e.Y5BaseSalary
	c.Y5Bonus = e.Y5Bonus
	c.BonusPercentage = e.BonusPercentage
	c.CalculateContract()
	c.SigningValue = c.ContractValue
	c.IsActive = true
	c.IsComplete = false
	c.IsExtended = true
	c.IsTagged = e.IsTag
}

func (c *NFLContract) MapPracticeSquadOffer(f FreeAgencyOffer) {
	c.OriginalTeam = c.Team
	c.OriginalTeamID = c.TeamID
	c.TeamID = f.TeamID
	c.Team = f.Team
	c.ContractLength = f.ContractLength
	c.Y1BaseSalary = f.Y1BaseSalary
	c.Y1Bonus = f.Y1Bonus
	c.Y2BaseSalary = f.Y2BaseSalary
	c.Y2Bonus = f.Y2Bonus
	c.Y3BaseSalary = f.Y3BaseSalary
	c.Y3Bonus = f.Y3Bonus
	c.Y4BaseSalary = f.Y4BaseSalary
	c.Y4Bonus = f.Y4Bonus
	c.Y5BaseSalary = f.Y5BaseSalary
	c.Y5Bonus = f.Y5Bonus
	c.BonusPercentage = f.BonusPercentage
	c.CalculateContract()
	c.IsActive = true
	c.IsComplete = false
	c.IsExtended = false
	c.IsTagged = false
}

func (c *NFLContract) FixContract(length int, y1s, y1b, y2s, y2b, y3s, y3b, y4s, y4b, y5s, y5b, bonus float64) {
	c.ContractLength = length
	c.Y1BaseSalary = y1s
	c.Y1Bonus = y1b
	c.Y2BaseSalary = y2s
	c.Y2Bonus = y2b
	c.Y3BaseSalary = y3s
	c.Y3Bonus = y3b
	c.Y4BaseSalary = y4s
	c.Y4Bonus = y4b
	c.Y5BaseSalary = y5s
	c.Y5Bonus = y5b
	c.BonusPercentage = bonus
	c.CalculateContract()
	c.IsActive = true
	c.IsExtended = true
}

func (c *NFLContract) ToggleRetirement() {
	c.PlayerRetired = true
}

func (c *NFLContract) TagContract(tagType uint8, salary, bonus float64) {
	if c.ContractLength == 1 {
		c.IsTagged = true
		c.ContractLength += 1
		c.Y2BaseSalary = salary
		c.Y2Bonus = bonus
	}
}

type NFLTagDTO struct {
	PlayerID uint
	TagType  uint8
	Position string
}
