package structs

import "github.com/jinzhu/gorm"

type FreeAgencyOfferDTO struct {
	ID             uint
	NFLPlayerID    uint
	TeamID         uint
	Team           string
	ContractLength int
	Y1BaseSalary   float64
	Y1Bonus        float64
	Y2BaseSalary   float64
	Y2Bonus        float64
	Y3BaseSalary   float64
	Y3Bonus        float64
	Y4BaseSalary   float64
	Y4Bonus        float64
	Y5BaseSalary   float64
	Y5Bonus        float64
}

type FreeAgencyOffer struct {
	gorm.Model
	NFLPlayerID     uint
	TeamID          uint
	Team            string
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
	TotalBonus      float64
	TotalSalary     float64
	ContractValue   float64
	BonusPercentage float64
	IsActive        bool
}

func (f *FreeAgencyOffer) CalculateOffer(offer FreeAgencyOfferDTO) {
	f.NFLPlayerID = offer.NFLPlayerID
	f.TeamID = offer.TeamID
	f.Team = offer.Team
	f.ContractLength = offer.ContractLength
	f.Y1BaseSalary = offer.Y1BaseSalary
	f.Y1Bonus = offer.Y1Bonus
	f.Y2BaseSalary = offer.Y2BaseSalary
	f.Y2Bonus = offer.Y2Bonus
	f.Y3BaseSalary = offer.Y3BaseSalary
	f.Y3Bonus = offer.Y3Bonus
	f.Y4BaseSalary = offer.Y4BaseSalary
	f.Y4Bonus = offer.Y4Bonus
	f.Y5BaseSalary = offer.Y5BaseSalary
	f.Y5Bonus = offer.Y5Bonus
	f.IsActive = true

	// Calculate Value
	y1SalaryVal := f.Y1BaseSalary * 0.8
	y1BonusVal := f.Y1Bonus * 1
	y2SalaryVal := f.Y2BaseSalary * 0.4
	y2BonusVal := f.Y2Bonus * 0.9
	y3SalaryVal := f.Y3BaseSalary * 0.2
	y3BonusVal := f.Y3Bonus * 0.8
	y4SalaryVal := f.Y4BaseSalary * 0.1
	y4BonusVal := f.Y4Bonus * 0.7
	y5SalaryVal := f.Y5BaseSalary * 0.05
	y5BonusVal := f.Y5Bonus * 0.6
	f.ContractValue = y1SalaryVal + y1BonusVal + y2SalaryVal + y2BonusVal + y3SalaryVal + y3BonusVal + y4SalaryVal + y4BonusVal + y5SalaryVal + y5BonusVal
	f.TotalBonus = f.Y1Bonus + f.Y2Bonus + f.Y3Bonus + f.Y4Bonus + f.Y5Bonus
	f.TotalSalary = f.Y1BaseSalary + f.Y2BaseSalary + f.Y3BaseSalary + f.Y4BaseSalary + f.Y5BaseSalary
	total := f.TotalBonus + f.TotalSalary
	f.BonusPercentage = f.TotalBonus / (total)
}

func (f *FreeAgencyOffer) CancelOffer() {
	f.IsActive = false
}

func (f *FreeAgencyOffer) AssignID(id uint) {
	f.ID = id
}

// Sorting Funcs
type ByContractValue []FreeAgencyOffer

func (fo ByContractValue) Len() int      { return len(fo) }
func (fo ByContractValue) Swap(i, j int) { fo[i], fo[j] = fo[j], fo[i] }
func (fo ByContractValue) Less(i, j int) bool {
	return fo[i].ContractValue > fo[j].ContractValue
}

// Table for storing Extensions for contracted players
type NFLExtensionOffer struct {
	gorm.Model
	NFLPlayerID     uint
	TeamID          uint
	SeasonID        uint
	Team            string
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
	TotalBonus      float64
	TotalSalary     float64
	ContractValue   float64
	BonusPercentage float64
	IsAccepted      bool
	IsActive        bool
}

func (f *NFLExtensionOffer) AssignID(id uint) {
	f.ID = id
}

func (f *NFLExtensionOffer) CalculateOffer(offer FreeAgencyOfferDTO) {
	f.NFLPlayerID = offer.NFLPlayerID
	f.TeamID = offer.TeamID
	f.Team = offer.Team
	f.ContractLength = offer.ContractLength
	f.Y1BaseSalary = offer.Y1BaseSalary
	f.Y1Bonus = offer.Y1Bonus
	f.Y2BaseSalary = offer.Y2BaseSalary
	f.Y2Bonus = offer.Y2Bonus
	f.Y3BaseSalary = offer.Y3BaseSalary
	f.Y3Bonus = offer.Y3Bonus
	f.Y4BaseSalary = offer.Y4BaseSalary
	f.Y4Bonus = offer.Y4Bonus
	f.Y5BaseSalary = offer.Y5BaseSalary
	f.Y5Bonus = offer.Y5Bonus
	f.IsActive = true

	// Calculate Value
	y1SalaryVal := f.Y1BaseSalary * 0.8
	y1BonusVal := f.Y1Bonus * 1
	y2SalaryVal := f.Y2BaseSalary * 0.4
	y2BonusVal := f.Y2Bonus * 0.9
	y3SalaryVal := f.Y3BaseSalary * 0.2
	y3BonusVal := f.Y3Bonus * 0.8
	y4SalaryVal := f.Y4BaseSalary * 0.1
	y4BonusVal := f.Y4Bonus * 0.7
	y5SalaryVal := f.Y5BaseSalary * 0.05
	y5BonusVal := f.Y5Bonus * 0.6
	f.ContractValue = y1SalaryVal + y1BonusVal + y2SalaryVal + y2BonusVal + y3SalaryVal + y3BonusVal + y4SalaryVal + y4BonusVal + y5SalaryVal + y5BonusVal
	f.TotalBonus = f.Y1Bonus + f.Y2Bonus + f.Y3Bonus + f.Y4Bonus + f.Y5Bonus
	f.TotalSalary = f.Y1BaseSalary + f.Y2BaseSalary + f.Y3BaseSalary + f.Y4BaseSalary + f.Y5BaseSalary
	total := f.TotalBonus + f.TotalSalary
	f.BonusPercentage = f.TotalBonus / (total)
}

func (f *NFLExtensionOffer) AcceptOffer() {
	f.IsAccepted = true
	f.CancelOffer()
}

func (f *NFLExtensionOffer) CancelOffer() {
	f.IsActive = false
}
