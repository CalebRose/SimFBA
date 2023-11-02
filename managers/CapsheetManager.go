package managers

import (
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/structs"
)

func GetCapsheetByTeamID(TeamID string) structs.NFLCapsheet {
	db := dbprovider.GetInstance().GetDB()

	capSheet := structs.NFLCapsheet{}

	err := db.Where("nfl_team_id = ?", TeamID).Find(&capSheet).Error
	if err != nil {
		fmt.Println("Could not find capsheet, returning new one")
		return structs.NFLCapsheet{}
	}

	return capSheet
}

func AllocateCapsheets() {
	db := dbprovider.GetInstance().GetDB()

	teams := GetAllNFLTeams()

	for _, team := range teams {
		TeamID := strconv.Itoa(int(team.ID))

		players := GetNFLPlayersWithContractsByTeamID(TeamID)

		Capsheet := GetCapsheetByTeamID(TeamID)

		if Capsheet.ID == 0 {
			Capsheet.AssignCapsheet(team.ID)
		}

		Capsheet.ResetCapsheet()

		sort.Sort(structs.ByTotalContract(players))
		window := 50
		for idx, player := range players {
			if idx > window {
				break
			}
			if player.IsPracticeSquad {
				window += 1
				continue
			}
			contract := player.Contract

			Capsheet.AddContractToCapsheet(contract)
		}

		db.Save(&Capsheet)
	}
}

func GetContractByPlayerID(PlayerID string) structs.NFLContract {
	db := dbprovider.GetInstance().GetDB()

	contract := structs.NFLContract{}

	err := db.Where("nfl_player_id = ? AND is_active = ?", PlayerID, true).Find(&contract).Error
	if err != nil {
		fmt.Println("No active contract for " + PlayerID)
		return structs.NFLContract{}
	}

	return contract
}

func GetAllContracts() []structs.NFLContract {
	db := dbprovider.GetInstance().GetDB()

	contracts := []structs.NFLContract{}

	err := db.Where("is_active = ?", true).Find(&contracts).Error
	if err != nil {
		log.Fatalln("Could not find all active contracts")
	}

	return contracts
}

func CalculateContractValues() {
	db := dbprovider.GetInstance().GetDB()

	contracts := GetAllContracts()

	for _, c := range contracts {
		c.CalculateContract()
		db.Save(&c)
	}
}

func AllocateRetiredContracts() {
	db := dbprovider.GetInstance().GetDB()
	// Get All Capsheets
	nflTeams := GetAllNFLTeams()
	capsheetMap := make(map[uint]*structs.NFLCapsheet)

	// Iterate & Map each capsheet
	for i := 0; i < len(nflTeams); i++ {
		team := nflTeams[i]
		TeamID := strconv.Itoa(int(team.ID))

		Capsheet := GetCapsheetByTeamID(TeamID)
		Capsheet.ProgressCapsheet()

		capsheetMap[team.ID] = &Capsheet
	}

	// Get All Retired Players
	retiredContract := GetRetiredContracts()
	// Iterate
	for _, contract := range retiredContract {
		// Add dead cap to next year
		capsheet := capsheetMap[uint(contract.TeamID)]
		capsheet.CutPlayerFromCapsheet(contract)
		db.Delete(&contract)
	}
	// If still active contract, add to dead cap
	// Then delete contract

	for _, team := range nflTeams {
		cap := capsheetMap[team.ID]

		db.Save(&cap)
	}
}
