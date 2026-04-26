package managers

import (
	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
)

func UpdateDraftPicks() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	path := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\2026\\2026_draft_picks_upload.csv"

	draftPickCSV := util.ReadCSV(path)

	draftPicks := GetAllCurrentSeasonDraftPicks()
	pickMap := make(map[uint]structs.NFLDraftPick)
	latestID := uint(1135) // Latest ID from Draft Pick table in DB

	for _, pick := range draftPicks {
		pickMap[pick.ID] = pick
	}

	for idx, row := range draftPickCSV {
		if idx == 0 {
			continue
		}

		draftPickID := util.ConvertStringToInt(row[0])
		draftRound := util.ConvertStringToInt(row[3])
		overallNumber := util.ConvertStringToInt(row[4])
		teamID := util.ConvertStringToInt(row[6])
		team := row[7]
		draftValue := util.ConvertStringToFloat(row[12])
		isCompensation := util.ConvertStringToBool(row[17])
		if isCompensation {
			draftPickID = int(latestID)
			latestID += 1
		}
		isVoid := util.ConvertStringToBool(row[18])
		draftPick := structs.NFLDraftPick{
			SeasonID: uint(ts.NFLSeasonID),
			Season:   uint(ts.Season),
		}
		if !isCompensation {
			draftPick = pickMap[uint(draftPickID)]
		}
		draftPick.MapValuesToDraftPick(uint(draftPickID), uint(draftRound), uint(overallNumber), uint(teamID), team, draftValue, isCompensation, isVoid)

		db.Save(&draftPick)
	}
}
