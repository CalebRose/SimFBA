package structs

import "github.com/jinzhu/gorm"

type CollegePlayerSeasonStats struct {
	gorm.Model
	CollegePlayerID uint
	TeamID          uint
	SeasonID        uint
	Year            uint
	BasePlayerStats
	GamesPlayed  int
	QBRating     float64
	Tackles      float64
	RushingAvg   float64
	PassingAvg   float64
	ReceivingAvg float64
	Completion   float64
}

func (ss *CollegePlayerSeasonStats) ResetStats() {
	ss.GamesPlayed = 0
	ss.PassingYards = 0
	ss.PassAttempts = 0
	ss.PassCompletions = 0
	ss.PassingTDs = 0
	ss.Interceptions = 0
	ss.Sacks = 0
	ss.LongestPass = 0
	ss.RushAttempts = 0
	ss.RushingTDs = 0
	ss.RushingYards = 0
	ss.LongestRush = 0
	ss.Fumbles = 0
	ss.Targets = 0
	ss.Catches = 0
	ss.ReceivingTDs = 0
	ss.ReceivingYards = 0
	ss.LongestReception = 0
	ss.Tackles = 0
	ss.AssistedTackles = 0
	ss.SoloTackles = 0
	ss.RecoveredFumbles = 0
	ss.SacksMade = 0
	ss.ForcedFumbles = 0
	ss.TacklesForLoss = 0
	ss.PassDeflections = 0
	ss.InterceptionsCaught = 0
	ss.Safeties = 0
	ss.DefensiveTDs = 0
	ss.FGMade = 0
	ss.FGAttempts = 0
	ss.LongestFG = 0
	ss.ExtraPointsAttempted = 0
	ss.ExtraPointsMade = 0
	ss.KickoffTouchbacks = 0
	ss.Punts = 0
	ss.PuntTouchbacks = 0
	ss.PuntsInside20 = 0
	ss.KickReturns = 0
	ss.KickReturnYards = 0
	ss.KickReturnTDs = 0
	ss.PuntReturns = 0
	ss.PuntReturnYards = 0
	ss.PuntReturnTDs = 0
	ss.STSoloTackles = 0
	ss.STAssistedTackles = 0
	ss.PuntsBlocked = 0
	ss.FGBlocked = 0
	ss.Snaps = 0
	ss.Pancakes = 0
	ss.SacksAllowed = 0
	ss.PlayedGame = 0
	ss.StartedGame = 0
	ss.QBRating = 0
	ss.PassingAvg = 0
	ss.Completion = 0
	ss.RushingAvg = 0
	ss.ReceivingAvg = 0
}

func (ss *CollegePlayerSeasonStats) MapStats(stats []CollegePlayerStats) {
	for _, stat := range stats {
		ss.GamesPlayed++
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
		ss.Tackles = float64(ss.SoloTackles) + (float64(ss.AssistedTackles) / 2) + ss.STSoloTackles + (float64(ss.AssistedTackles) / 2)
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
	if ss.PassAttempts > 0 {
		numerator := passingYards + passingTDs + passComps - ints
		ss.QBRating = numerator / float64(ss.PassAttempts)
		ss.PassingAvg = float64(ss.PassingYards) / float64(ss.GamesPlayed)
		ss.Completion = float64(ss.PassCompletions) / float64(ss.PassAttempts)
	}

	if ss.RushAttempts > 0 {
		ss.RushingAvg = float64(ss.RushingYards) / float64(ss.RushAttempts)

	}

	if ss.Catches > 0 {
		ss.ReceivingAvg = float64(ss.ReceivingYards) / float64(ss.Catches)
	}
}

func (ss *CollegePlayerSeasonStats) ReduceStats(stats []CollegePlayerStats) {
	for _, stat := range stats {
		ss.GamesPlayed--
		ss.PassingYards = ss.PassingYards - stat.PassingYards
		ss.PassAttempts = ss.PassAttempts - stat.PassAttempts
		ss.PassCompletions = ss.PassCompletions - stat.PassCompletions
		ss.PassingTDs = ss.PassingTDs - stat.PassingTDs
		ss.Interceptions = ss.Interceptions - stat.Interceptions
		ss.Sacks = ss.Sacks - stat.Sacks
		if stat.LongestPass > ss.LongestPass {
			ss.LongestPass = stat.LongestPass
		}
		ss.RushAttempts = ss.RushAttempts - stat.RushAttempts
		ss.RushingTDs = ss.RushingTDs - stat.RushingTDs
		ss.RushingYards = ss.RushingYards - stat.RushingYards
		if stat.LongestRush > ss.LongestRush {
			ss.LongestRush = stat.LongestRush
		}
		ss.Fumbles = ss.Fumbles - stat.Fumbles
		ss.Targets = ss.Targets - stat.Targets
		ss.Catches = ss.Catches - stat.Catches
		ss.ReceivingTDs = ss.ReceivingTDs - stat.ReceivingTDs
		ss.ReceivingYards = ss.ReceivingYards - stat.ReceivingYards
		if stat.LongestReception > ss.LongestReception {
			ss.LongestReception = stat.LongestReception
		}
		ss.Tackles = float64(ss.SoloTackles) - (float64(ss.AssistedTackles) / 2) - ss.STSoloTackles - (float64(ss.AssistedTackles) / 2)
		ss.AssistedTackles = ss.AssistedTackles - stat.AssistedTackles
		ss.SoloTackles = ss.SoloTackles - stat.SoloTackles
		ss.RecoveredFumbles = ss.RecoveredFumbles - stat.RecoveredFumbles
		ss.SacksMade = ss.SacksMade - stat.SacksMade
		ss.ForcedFumbles = ss.ForcedFumbles - stat.ForcedFumbles
		ss.TacklesForLoss = ss.TacklesForLoss - stat.TacklesForLoss
		ss.PassDeflections = ss.PassDeflections - stat.PassDeflections
		ss.InterceptionsCaught = ss.InterceptionsCaught - stat.InterceptionsCaught
		ss.Safeties = ss.Safeties - stat.Safeties
		ss.DefensiveTDs = ss.DefensiveTDs - stat.DefensiveTDs
		ss.FGMade = ss.FGMade - stat.FGMade
		ss.FGAttempts = ss.FGAttempts - stat.FGAttempts
		if stat.LongestFG > ss.LongestFG {
			ss.LongestFG = stat.LongestFG
		}
		ss.ExtraPointsAttempted = ss.ExtraPointsAttempted - stat.ExtraPointsAttempted
		ss.ExtraPointsMade = ss.ExtraPointsMade - stat.ExtraPointsMade
		ss.KickoffTouchbacks = ss.KickoffTouchbacks - stat.KickoffTouchbacks
		ss.Punts = ss.Punts - stat.Punts
		ss.PuntTouchbacks = ss.PuntTouchbacks - stat.PuntTouchbacks
		ss.PuntsInside20 = ss.PuntsInside20 - stat.PuntsInside20
		ss.KickReturns = ss.KickReturns - stat.KickReturns
		ss.KickReturnYards = ss.KickReturnYards - stat.KickReturnYards
		ss.KickReturnTDs = ss.KickReturnTDs - stat.KickReturnTDs
		ss.PuntReturns = ss.PuntReturns - stat.PuntReturns
		ss.PuntReturnYards = ss.PuntReturnYards - stat.PuntReturnYards
		ss.PuntReturnTDs = ss.PuntReturnTDs - stat.PuntReturnTDs
		ss.STSoloTackles = ss.STSoloTackles - stat.STSoloTackles
		ss.STAssistedTackles = ss.STAssistedTackles - stat.STAssistedTackles
		ss.PuntsBlocked = ss.PuntsBlocked - stat.PuntsBlocked
		ss.FGBlocked = ss.FGBlocked - stat.FGBlocked
		ss.Snaps = ss.Snaps - stat.Snaps
		ss.Pancakes = ss.Pancakes - stat.Pancakes
		ss.SacksAllowed = ss.SacksAllowed - stat.SacksAllowed
		ss.PlayedGame = ss.PlayedGame - stat.PlayedGame
		ss.StartedGame = ss.StartedGame - stat.StartedGame
	}
	passingYards := float64(8.4) * float64(ss.PassingYards)
	passingTDs := float64(330) * float64(ss.PassingTDs)
	passComps := float64(100) * float64(ss.PassCompletions)
	ints := float64(200) * float64(ss.Interceptions)
	if ss.PassAttempts > 0 {
		numerator := passingYards + passingTDs + passComps - ints
		ss.QBRating = numerator / float64(ss.PassAttempts)
		ss.PassingAvg = float64(ss.PassingYards) / float64(ss.GamesPlayed)
		ss.Completion = float64(ss.PassCompletions) / float64(ss.PassAttempts)
	}

	if ss.RushAttempts > 0 {
		ss.RushingAvg = float64(ss.RushingYards) / float64(ss.RushAttempts)

	}

	if ss.Catches > 0 {
		ss.ReceivingAvg = float64(ss.ReceivingYards) / float64(ss.Catches)
	}
}
