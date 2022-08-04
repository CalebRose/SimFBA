package models

import "github.com/CalebRose/SimFBA/structs"

type CollegePlayerResponse struct {
	ID int
	structs.BasePlayer
	TeamID       int
	TeamAbbr     string
	City         string
	State        string
	Year         int
	IsRedshirt   bool
	ConferenceID int
	Conference   string
	PlayerStats  []structs.CollegePlayerStats
	SeasonStats  PlayerSeasonStats
}

func (cpr *CollegePlayerResponse) MapSeasonalStats() {
	var ss PlayerSeasonStats
	for _, stat := range cpr.PlayerStats {
		ss.PassingYards = ss.PassingYards + stat.PassingYards
		ss.PassAttempts = ss.PassAttempts + stat.PassAttempts
		ss.PassCompletions = ss.PassCompletions + stat.PassCompletions
		ss.PassingTDs = ss.PassingTDs + stat.PassingTDs
		ss.Interceptions = ss.Interceptions + stat.Interceptions
		ss.Sacks = ss.Sacks + stat.Sacks
		if stat.LongestPass > ss.LongestPass {
			ss.LongestPass = stat.LongestPass
		}
		ss.RushAttempts = ss.RushAttempts + stat.RushAttempts
		ss.RushingTDs = ss.RushingTDs + stat.RushingTDs
		ss.RushingYards = ss.RushingYards + stat.RushingYards
		if stat.LongestRush > ss.LongestRush {
			ss.LongestRush = stat.LongestRush
		}
		ss.Fumbles = ss.Fumbles + stat.Fumbles
		ss.Targets = ss.Targets + stat.Targets
		ss.Catches = ss.Catches + stat.Catches
		ss.ReceivingTDs = ss.ReceivingTDs + stat.ReceivingTDs
		ss.ReceivingYards = ss.ReceivingYards + stat.ReceivingYards
		if stat.LongestReception > ss.LongestReception {
			ss.LongestReception = stat.LongestReception
		}
		ss.Tackles = float64(ss.SoloTackles) + (float64(ss.AssistedTackles) / 2)
		ss.AssistedTackles = ss.AssistedTackles + stat.AssistedTackles
		ss.SoloTackles = ss.SoloTackles + stat.SoloTackles
		ss.RecoveredFumbles = ss.RecoveredFumbles + stat.RecoveredFumbles
		ss.SacksMade = ss.SacksMade + stat.SacksMade
		ss.ForcedFumbles = ss.ForcedFumbles + stat.ForcedFumbles
		ss.TacklesForLoss = ss.TacklesForLoss + stat.TacklesForLoss
		ss.PassDeflections = ss.PassDeflections + stat.PassDeflections
		ss.InterceptionsCaught = ss.InterceptionsCaught + stat.InterceptionsCaught
		ss.Safeties = ss.Safeties + stat.Safeties
		ss.DefensiveTDs = ss.DefensiveTDs + stat.DefensiveTDs
		ss.FGMade = ss.FGMade + stat.FGMade
		ss.FGAttempts = ss.FGAttempts + stat.FGAttempts
		if stat.LongestFG > ss.LongestFG {
			ss.LongestFG = stat.LongestFG
		}
		ss.ExtraPointsAttempted = ss.ExtraPointsAttempted + stat.ExtraPointsAttempted
		ss.ExtraPointsMade = ss.ExtraPointsMade + stat.ExtraPointsMade
		ss.KickoffTouchbacks = ss.KickoffTouchbacks + stat.KickoffTouchbacks
		ss.Punts = ss.Punts + stat.Punts
		ss.PuntTouchbacks = ss.PuntTouchbacks + stat.PuntTouchbacks
		ss.PuntsInside20 = ss.PuntsInside20 + stat.PuntsInside20
		ss.KickReturns = ss.KickReturns + stat.KickReturns
		ss.KickReturnYards = ss.KickReturnYards + stat.KickReturnYards
		ss.KickReturnTDs = ss.KickReturnTDs + stat.KickReturnTDs
		ss.PuntReturns = ss.PuntReturns + stat.PuntReturns
		ss.PuntReturnYards = ss.PuntReturnYards + stat.PuntReturnYards
		ss.PuntReturnTDs = ss.PuntReturnTDs + stat.PuntReturnTDs
		ss.STSoloTackles = ss.STSoloTackles + stat.STSoloTackles
		ss.STAssistedTackles = ss.STAssistedTackles + stat.STAssistedTackles
		ss.PuntsBlocked = ss.PuntsBlocked + stat.PuntsBlocked
		ss.FGBlocked = ss.FGBlocked + stat.FGBlocked
		ss.Snaps = ss.Snaps + stat.Snaps
		ss.Pancakes = ss.Pancakes + stat.Pancakes
		ss.SacksAllowed = ss.SacksAllowed + stat.SacksAllowed
		ss.PlayedGame = ss.PlayedGame + stat.PlayedGame
		ss.StartedGame = ss.StartedGame + stat.StartedGame
	}
	passingYards := float64(8.4) * float64(ss.PassingYards)
	passingTDs := float64(330) * float64(ss.PassingTDs)
	passComps := float64(100) * float64(ss.PassCompletions)
	ints := float64(200) * float64(ss.Interceptions)
	if ss.PassAttempts != 0 {
		numerator := passingYards + passingTDs + passComps - ints
		ss.QBRating = numerator / float64(ss.PassAttempts)
	}

	cpr.SeasonStats = ss
}
