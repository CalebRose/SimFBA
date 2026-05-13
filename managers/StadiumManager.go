package managers

import (
	"fmt"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/structs"
)

func GetAllStadiums() []structs.Stadium {
	db := dbprovider.GetInstance().GetDB()

	var stadiums []structs.Stadium

	db.Find(&stadiums)

	return stadiums
}

func GetStadiumByStadiumID(id string) structs.Stadium {
	db := dbprovider.GetInstance().GetDB()

	var stadium structs.Stadium

	db.Where("id = ?", id).Find(&stadium)

	return stadium
}

// EnsureNFLStadiumsExist checks every NFL team and creates a Stadium record for any
// team that doesn't already have one, using the address and capacity values stored
// on the NFLTeam struct itself.
func EnsureNFLStadiumsExist() {
	db := dbprovider.GetInstance().GetDB()

	allTeams := GetAllNFLTeams()

	// Build a map of teamID -> Stadium for all existing NFL stadiums.
	stadiums := GetAllStadiums()
	stadiumByTeam := make(map[uint]structs.Stadium, len(stadiums))
	for _, s := range stadiums {
		if s.LeagueName != "NFL" {
			continue
		}
		stadiumByTeam[s.TeamID] = s
	}

	// For each NFL team without a stadium record, create one from the team's own fields.
	for _, t := range allTeams {
		if _, exists := stadiumByTeam[t.ID]; exists {
			continue
		}
		stadium := structs.Stadium{
			StadiumName:      t.Stadium,
			TeamID:           t.ID,
			TeamAbbr:         t.TeamAbbr,
			City:             t.City,
			State:            t.State,
			Country:          t.Country,
			Capacity:         uint(t.StadiumCapacity),
			RecordAttendance: uint(t.RecordAttendance),
			LeagueID:         3,
			LeagueName:       "NFL",
		}
		db.Create(&stadium)
		fmt.Printf("[INFO] Created stadium %q for %s %s (teamID=%d)\n",
			stadium.StadiumName, t.TeamName, t.Mascot, t.ID)
	}
}
