package managers

import (
	"github.com/CalebRose/SimFBA/structs"
)

func MakeCollegeInjuryList(players []structs.CollegePlayer) []structs.CollegePlayer {
	injuryList := []structs.CollegePlayer{}

	for _, p := range players {
		if p.IsInjured {
			injuryList = append(injuryList, p)
		}
	}
	return injuryList
}

func MakeCollegePortalList(players []structs.CollegePlayer) []structs.CollegePlayer {
	portalList := []structs.CollegePlayer{}

	for _, p := range players {
		if p.TransferStatus > 0 {
			portalList = append(portalList, p)
		}
	}
	return portalList
}

func MakeProInjuryList(players []structs.NFLPlayer) []structs.NFLPlayer {
	injuryList := []structs.NFLPlayer{}

	for _, p := range players {
		if p.IsInjured {
			injuryList = append(injuryList, p)
		}
	}
	return injuryList
}

func MakePracticeSquadList(players []structs.NFLPlayer) []structs.NFLPlayer {
	playerList := []structs.NFLPlayer{}

	for _, p := range players {
		if p.IsPracticeSquad {
			playerList = append(playerList, p)
		}
	}
	return playerList
}

func MakeGeneralFreeAgentList(players []structs.NFLPlayer, phase uint) []structs.NFLPlayer {
	playerList := []structs.NFLPlayer{}

	for _, p := range players {
		if p.IsFreeAgent && p.Experience > 1 && phase >= 6 {
			playerList = append(playerList, p)
		}
	}
	return playerList
}

func MakeUDFAList(players []structs.NFLPlayer) []structs.NFLPlayer {
	playerList := []structs.NFLPlayer{}

	for _, p := range players {
		if p.IsFreeAgent && p.Experience <= 1 {
			playerList = append(playerList, p)
		}
	}
	return playerList
}

func MakeGameResultsPlayerListFromCFB(stats []structs.CollegePlayerStats, players []structs.CollegePlayer) []structs.GameResultsPlayer {
	var matchRows []structs.GameResultsPlayer
	statMap := make(map[uint]structs.CollegePlayerStats)
	for _, s := range stats {
		statMap[uint(s.CollegePlayerID)] = s
	}
	for _, p := range players {
		s := statMap[p.ID]
		if s.ID == 0 || s.Snaps == 0 {
			continue
		}

		row := structs.GameResultsPlayer{
			ID:                   p.ID,
			FirstName:            p.FirstName,
			LastName:             p.LastName,
			Position:             p.Position,
			Archetype:            p.Archetype,
			Year:                 uint(p.Year),
			TeamAbbr:             p.TeamAbbr,
			League:               "CFB",
			Snaps:                int(s.Snaps),
			PassingYards:         int(s.PassingYards),
			PassAttempts:         int(s.PassAttempts),
			PassCompletions:      int(s.PassCompletions),
			PassingTDs:           int(s.PassingTDs),
			Interceptions:        int(s.Interceptions),
			LongestPass:          int(s.LongestPass),
			Sacks:                int(s.Sacks),
			RushAttempts:         int(s.RushAttempts),
			RushingYards:         int(s.RushingYards),
			RushingTDs:           int(s.RushingTDs),
			Fumbles:              int(s.Fumbles),
			LongestRush:          int(s.LongestRush),
			Targets:              int(s.Targets),
			Catches:              int(s.Catches),
			ReceivingYards:       int(s.ReceivingYards),
			ReceivingTDs:         int(s.ReceivingTDs),
			LongestReception:     int(s.LongestReception),
			SoloTackles:          s.SoloTackles,
			AssistedTackles:      s.AssistedTackles,
			TacklesForLoss:       s.TacklesForLoss,
			SacksMade:            s.SacksMade,
			ForcedFumbles:        int(s.ForcedFumbles),
			RecoveredFumbles:     int(s.RecoveredFumbles),
			PassDeflections:      int(s.PassDeflections),
			InterceptionsCaught:  int(s.InterceptionsCaught),
			Safeties:             int(s.Safeties),
			DefensiveTDs:         int(s.DefensiveTDs),
			FGMade:               int(s.FGMade),
			FGAttempts:           int(s.FGAttempts),
			LongestFG:            int(s.LongestFG),
			ExtraPointsMade:      int(s.ExtraPointsMade),
			ExtraPointsAttempted: int(s.ExtraPointsAttempted),
			KickoffTouchbacks:    int(s.KickoffTouchbacks),
			Punts:                int(s.Punts),
			PuntTouchbacks:       int(s.PuntTouchbacks),
			PuntsInside20:        int(s.PuntsInside20),
			KickReturns:          int(s.KickReturns),
			KickReturnTDs:        int(s.KickReturnTDs),
			KickReturnYards:      int(s.KickReturnYards),
			PuntReturns:          int(s.PuntReturns),
			PuntReturnTDs:        int(s.PuntReturnTDs),
			PuntReturnYards:      int(s.PuntReturnYards),
			STSoloTackles:        s.STSoloTackles,
			STAssistedTackles:    s.STAssistedTackles,
			PuntsBlocked:         int(s.PuntsBlocked),
			FGBlocked:            int(s.FGBlocked),
			Pancakes:             int(s.Pancakes),
			SacksAllowed:         int(s.SacksAllowed),
			DefensivePressures:   int(s.DefensivePressures),
			Hurries:              int(s.Hurries),
			PassRushSnaps:        int(s.PassRushSnaps),
			PassRushWins:         int(s.PassRushWins),
			PressuresAllowed:     int(s.PressuresAllowed),
			PassBlockSnaps:       int(s.PassBlockSnaps),
			PassBlockWins:        int(s.PassBlockWins),
			PlayedGame:           int(s.PlayedGame),
			StartedGame:          int(s.StartedGame),
			WasInjured:           s.WasInjured,
			WeeksOfRecovery:      s.WeeksOfRecovery,
			InjuryType:           s.InjuryType,
		}

		matchRows = append(matchRows, row)
	}

	return matchRows
}

func MakeGameResultsPlayerListFromNFL(stats []structs.NFLPlayerStats, players []structs.NFLPlayer) []structs.GameResultsPlayer {
	var matchRows []structs.GameResultsPlayer
	statMap := make(map[uint]structs.NFLPlayerStats)
	for _, s := range stats {
		statMap[uint(s.NFLPlayerID)] = s
	}
	for _, p := range players {
		s := statMap[p.ID]
		if s.ID == 0 || s.Snaps == 0 {
			continue
		}

		row := structs.GameResultsPlayer{
			ID:                   p.ID,
			FirstName:            p.FirstName,
			LastName:             p.LastName,
			Position:             p.Position,
			Archetype:            p.Archetype,
			Year:                 uint(p.Experience),
			TeamAbbr:             p.TeamAbbr,
			League:               "NFL",
			Snaps:                int(s.Snaps),
			PassingYards:         int(s.PassingYards),
			PassAttempts:         int(s.PassAttempts),
			PassCompletions:      int(s.PassCompletions),
			PassingTDs:           int(s.PassingTDs),
			Interceptions:        int(s.Interceptions),
			LongestPass:          int(s.LongestPass),
			Sacks:                int(s.Sacks),
			RushAttempts:         int(s.RushAttempts),
			RushingYards:         int(s.RushingYards),
			RushingTDs:           int(s.RushingTDs),
			Fumbles:              int(s.Fumbles),
			LongestRush:          int(s.LongestRush),
			Targets:              int(s.Targets),
			Catches:              int(s.Catches),
			ReceivingYards:       int(s.ReceivingYards),
			ReceivingTDs:         int(s.ReceivingTDs),
			LongestReception:     int(s.LongestReception),
			SoloTackles:          s.SoloTackles,
			AssistedTackles:      s.AssistedTackles,
			TacklesForLoss:       s.TacklesForLoss,
			SacksMade:            s.SacksMade,
			ForcedFumbles:        int(s.ForcedFumbles),
			RecoveredFumbles:     int(s.RecoveredFumbles),
			PassDeflections:      int(s.PassDeflections),
			InterceptionsCaught:  int(s.InterceptionsCaught),
			Safeties:             int(s.Safeties),
			DefensiveTDs:         int(s.DefensiveTDs),
			FGMade:               int(s.FGMade),
			FGAttempts:           int(s.FGAttempts),
			LongestFG:            int(s.LongestFG),
			ExtraPointsMade:      int(s.ExtraPointsMade),
			ExtraPointsAttempted: int(s.ExtraPointsAttempted),
			KickoffTouchbacks:    int(s.KickoffTouchbacks),
			Punts:                int(s.Punts),
			PuntTouchbacks:       int(s.PuntTouchbacks),
			PuntsInside20:        int(s.PuntsInside20),
			KickReturns:          int(s.KickReturns),
			KickReturnTDs:        int(s.KickReturnTDs),
			KickReturnYards:      int(s.KickReturnYards),
			PuntReturns:          int(s.PuntReturns),
			PuntReturnTDs:        int(s.PuntReturnTDs),
			PuntReturnYards:      int(s.PuntReturnYards),
			STSoloTackles:        s.STSoloTackles,
			STAssistedTackles:    s.STAssistedTackles,
			PuntsBlocked:         int(s.PuntsBlocked),
			FGBlocked:            int(s.FGBlocked),
			Pancakes:             int(s.Pancakes),
			SacksAllowed:         int(s.SacksAllowed),
			DefensivePressures:   int(s.DefensivePressures),
			Hurries:              int(s.Hurries),
			PassRushSnaps:        int(s.PassRushSnaps),
			PassRushWins:         int(s.PassRushWins),
			PressuresAllowed:     int(s.PressuresAllowed),
			PassBlockSnaps:       int(s.PassBlockSnaps),
			PassBlockWins:        int(s.PassBlockWins),
			PlayedGame:           int(s.PlayedGame),
			StartedGame:          int(s.StartedGame),
			WasInjured:           s.WasInjured,
			WeeksOfRecovery:      s.WeeksOfRecovery,
			InjuryType:           s.InjuryType,
		}

		matchRows = append(matchRows, row)
	}

	return matchRows
}
