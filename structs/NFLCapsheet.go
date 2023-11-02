package structs

import "github.com/jinzhu/gorm"

type NFLCapsheet struct {
	gorm.Model
	NFLTeamID uint
	Y1Bonus   float64
	Y1Salary  float64
	Y1CapHit  float64
	Y2Bonus   float64
	Y2Salary  float64
	Y2CapHit  float64
	Y3Bonus   float64
	Y3Salary  float64
	Y3CapHit  float64
	Y4Bonus   float64
	Y4Salary  float64
	Y4CapHit  float64
	Y5Bonus   float64
	Y5Salary  float64
	Y5CapHit  float64
}

func (nc *NFLCapsheet) AssignCapsheet(id uint) {
	nc.ID = id
	nc.NFLTeamID = id
}

func (nc *NFLCapsheet) ResetCapsheet() {
	nc.Y1Bonus = 0
	nc.Y1Salary = 0
	nc.Y2Bonus = 0
	nc.Y2Salary = 0
	nc.Y3Bonus = 0
	nc.Y3Salary = 0
	nc.Y4Bonus = 0
	nc.Y4Salary = 0
	nc.Y5Bonus = 0
	nc.Y5Salary = 0
}

func (nc *NFLCapsheet) AddContractToCapsheet(contract NFLContract) {
	nc.Y1Bonus += contract.Y1Bonus
	nc.Y1Salary += contract.Y1BaseSalary
	nc.Y2Bonus += contract.Y2Bonus
	nc.Y2Salary += contract.Y2BaseSalary
	nc.Y3Bonus += contract.Y3Bonus
	nc.Y3Salary += contract.Y3BaseSalary
	nc.Y4Bonus += contract.Y4Bonus
	nc.Y4Salary += contract.Y4BaseSalary
	nc.Y5Bonus += contract.Y5Bonus
	nc.Y5Salary += contract.Y5BaseSalary
}

func (nc *NFLCapsheet) SubtractFromCapsheet(contract NFLContract) {
	nc.Y1CapHit += contract.Y1Bonus
	nc.Y1Bonus -= contract.Y1Bonus
	nc.Y2Bonus -= contract.Y2Bonus
	nc.Y3Bonus -= contract.Y3Bonus
	nc.Y4Bonus -= contract.Y4Bonus
	nc.Y5Bonus -= contract.Y5Bonus
	nc.Y1Salary -= contract.Y1BaseSalary
	nc.Y2Salary -= contract.Y2BaseSalary
	nc.Y3Salary -= contract.Y3BaseSalary
	nc.Y4Salary -= contract.Y4BaseSalary
	nc.Y5Salary -= contract.Y5BaseSalary
}

func (nc *NFLCapsheet) CutPlayerFromCapsheet(contract NFLContract) {
	nc.Y1CapHit += contract.Y1Bonus + contract.Y2Bonus + contract.Y3Bonus + contract.Y4Bonus + contract.Y5Bonus
	nc.Y1Bonus -= contract.Y1Bonus
	nc.Y2Bonus -= contract.Y2Bonus
	nc.Y3Bonus -= contract.Y3Bonus
	nc.Y4Bonus -= contract.Y4Bonus
	nc.Y5Bonus -= contract.Y5Bonus
	nc.Y1Salary -= contract.Y1BaseSalary
	nc.Y2Salary -= contract.Y2BaseSalary
	nc.Y3Salary -= contract.Y3BaseSalary
	nc.Y4Salary -= contract.Y4BaseSalary
	nc.Y5Salary -= contract.Y5BaseSalary
}

func (nc *NFLCapsheet) SubtractFromCapsheetViaTrade(contract NFLContract) {
	nc.Y1CapHit += contract.Y1Bonus
	nc.Y1Bonus -= contract.Y1Bonus
	nc.Y2Bonus -= contract.Y2Bonus
	nc.Y3Bonus -= contract.Y3Bonus
	nc.Y4Bonus -= contract.Y4Bonus
	nc.Y5Bonus -= contract.Y5Bonus
	nc.Y2Salary -= contract.Y2BaseSalary
	nc.Y3Salary -= contract.Y3BaseSalary
	nc.Y4Salary -= contract.Y4BaseSalary
	nc.Y5Salary -= contract.Y5BaseSalary
}

func (nc *NFLCapsheet) NegotiateSalaryDifference(SalaryDifference float64, CapHit float64) {
	nc.Y1Salary -= SalaryDifference
	nc.Y1CapHit += CapHit
}

func (nc *NFLCapsheet) AddContractViaTrade(contract NFLContract, differenceValue float64) {
	// nc.Y1Bonus += contract.Y1Bonus
	nc.Y1Salary += differenceValue
	nc.Y2Bonus += contract.Y2Bonus
	nc.Y2Salary += contract.Y2BaseSalary
	nc.Y3Bonus += contract.Y3Bonus
	nc.Y3Salary += contract.Y3BaseSalary
	nc.Y4Bonus += contract.Y4Bonus
	nc.Y4Salary += contract.Y4BaseSalary
	nc.Y5Bonus += contract.Y5Bonus
	nc.Y5Salary += contract.Y5BaseSalary
}

func (nc *NFLCapsheet) ProgressCapsheet() {
	nc.Y1Salary = nc.Y2Salary
	nc.Y1Bonus = nc.Y2Bonus
	nc.Y1CapHit = nc.Y2CapHit
	nc.Y2Salary = nc.Y3Salary
	nc.Y2Bonus = nc.Y3Bonus
	nc.Y2CapHit = nc.Y3CapHit
	nc.Y3Salary = nc.Y4Salary
	nc.Y3Bonus = nc.Y4Bonus
	nc.Y3CapHit = nc.Y4CapHit
	nc.Y4Salary = nc.Y5Salary
	nc.Y4Bonus = nc.Y5Bonus
	nc.Y4CapHit = nc.Y5CapHit
	nc.Y5Salary = 0
	nc.Y5Bonus = 0
	nc.Y5CapHit = 0
}
