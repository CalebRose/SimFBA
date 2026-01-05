package managers

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
)

func GetCollegeGameStructByGameID(games []structs.CollegeGame, gameID int) structs.CollegeGame {
	for _, game := range games {
		if int(game.ID) == gameID {
			return game
		}
	}

	return structs.CollegeGame{}
}

func GetCollegePlayerIDsBySeasonStats(cps []structs.CollegePlayerSeasonStats) []string {
	var list []string

	for _, cp := range cps {
		list = append(list, strconv.Itoa(int(cp.CollegePlayerID)))
	}

	return list
}

func GetNFLPlayerIDsBySeasonStats(cps []structs.NFLPlayerSeasonStats) []string {
	var list []string

	for _, cp := range cps {
		list = append(list, strconv.Itoa(int(cp.NFLPlayerID)))
	}

	return list
}

func GetCollegePlayerIDs(cps []structs.CollegePlayerStats) []string {
	var list []string

	for _, cp := range cps {
		list = append(list, strconv.Itoa(int(cp.CollegePlayerID)))
	}

	return list
}

func GetOffensiveDefaultSchemes() map[string]structs.OffensiveFormation {
	path := filepath.Join(os.Getenv("ROOT"), "data", "defaultOffensiveSchemes.json")
	content := util.ReadJson(path)

	var payload map[string]structs.OffensiveFormation

	err := json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatalln("Error during unmarshal: ", err)
	}

	return payload
}

func GetDefensiveDefaultSchemes() map[string]map[string]structs.DefensiveFormation {
	path := filepath.Join(os.Getenv("ROOT"), "data", "defaultDefensiveSchemes.json")
	content := util.ReadJson(path)

	var payload map[string]map[string]structs.DefensiveFormation

	err := json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatalln("Error during unmarshal: ", err)
	}

	return payload
}

func IsAITeamContendingForCroot(profiles []structs.RecruitPlayerProfile) float64 {
	if len(profiles) == 0 {
		return 0
	}
	var leadingVal float64 = 0
	for _, profile := range profiles {
		if profile.TotalPoints != 0 && profile.TotalPoints > float64(leadingVal) {
			leadingVal = profile.TotalPoints
		}
	}

	return leadingVal
}

func GetWinsAndLossesForCollegeGames(games []structs.CollegeGame, TeamID int, ConferenceCheck bool) (int, int) {
	wins := 0
	losses := 0

	for _, game := range games {
		if !game.GameComplete {
			continue
		}
		if ConferenceCheck && !game.IsConference {
			continue
		}
		if (game.HomeTeamID == TeamID && game.HomeTeamWin) ||
			(game.AwayTeamID == TeamID && game.AwayTeamWin) {
			wins += 1
		} else {
			losses += 1
		}
	}

	return wins, losses
}

func GetConferenceChampionshipWeight(games []structs.CollegeGame, TeamID int) float64 {
	var weight float64 = 0

	for _, game := range games {
		if !game.IsConference {
			continue
		}
		if (game.HomeTeamID == TeamID && game.HomeTeamScore > game.AwayTeamScore) ||
			(game.AwayTeamID == TeamID && game.AwayTeamScore > game.HomeTeamScore) {
			weight = 1
		} else {
			weight = 0.5
		}
	}

	return weight
}

func GetPostSeasonWeight(games []structs.CollegeGame, TeamID int) float64 {
	for _, game := range games {
		if !game.IsPlayoffGame || !game.IsBowlGame {
			continue
		}
		return 1
	}
	return 0
}

func FilterOutRecruitingProfile(profiles []structs.RecruitPlayerProfile, ID int) []structs.RecruitPlayerProfile {
	var rp []structs.RecruitPlayerProfile

	for _, profile := range profiles {
		if profile.ProfileID != ID {
			rp = append(rp, profile)
		}
	}

	return rp
}

func IsAITeamContendingForPortalPlayer(profiles []structs.TransferPortalProfile) int {
	if len(profiles) == 0 {
		return 0
	}
	leadingVal := 0
	for _, profile := range profiles {
		if profile.TotalPoints != 0 && profile.TotalPoints > float64(leadingVal) {
			leadingVal = int(profile.TotalPoints)
		}
	}

	return leadingVal
}

func FilterOutPortalProfile(profiles []structs.TransferPortalProfile, ID uint) []structs.TransferPortalProfile {
	var rp []structs.TransferPortalProfile

	for _, profile := range profiles {
		if profile.ProfileID != ID {
			rp = append(rp, profile)
		}
	}

	return rp
}
