package structs

import "github.com/jinzhu/gorm"

type NFLTeamSeasonStats struct {
	gorm.Model
	TeamID   uint
	SeasonID uint
	Year     int
	BaseTeamStats
	GamesPlayed         int
	TotalOffensiveYards int
	TotalYardsAllowed   int
	Fumbles             int
	QBRating            float64
	Tackles             float64
	Turnovers           int
}

type NFLTeamStats struct {
	gorm.Model
	TeamID       uint
	GameID       uint
	WeekID       uint
	SeasonID     uint
	OpposingTeam string
	BaseTeamStats
}

func (ss *NFLTeamSeasonStats) MapStats(stats []NFLTeamStats) {
	if ss.TeamID == 0 {
		ss.TeamID = stats[0].TeamID
		ss.SeasonID = 2
		ss.Year = 2023
	}
	for _, stat := range stats {
		ss.GamesPlayed++
		ss.PassingYards = ss.PassingYards + stat.PassingYards
		ss.PassingAttempts = ss.PassingAttempts + stat.PassingAttempts
		ss.PassingCompletions = ss.PassingCompletions + stat.PassingCompletions
		ss.PassingTouchdowns = ss.PassingTouchdowns + stat.PassingTouchdowns
		ss.PassingInterceptions = ss.PassingInterceptions + stat.PassingInterceptions
		ss.QBSacks = ss.QBSacks + stat.QBSacks
		if stat.LongestPass > ss.LongestPass {
			ss.LongestPass = stat.LongestPass
		}
		ss.RushAttempts = ss.RushAttempts + stat.RushAttempts
		ss.RushingTouchdowns = ss.RushingTouchdowns + stat.RushingTouchdowns
		ss.RushingYards = ss.RushingYards + stat.RushingYards
		ss.Fumbles = ss.Fumbles + stat.ReceivingFumbles + ss.RushingFumbles
		ss.ReceivingTargets = ss.ReceivingTargets + stat.ReceivingTargets
		ss.ReceivingCatches = ss.ReceivingCatches + stat.ReceivingCatches
		ss.ReceivingTouchdowns = ss.ReceivingTouchdowns + stat.ReceivingTouchdowns
		ss.ReceivingYards = ss.ReceivingYards + stat.ReceivingYards
		ss.AssistedTackles = ss.AssistedTackles + stat.AssistedTackles
		ss.Tackles = float64(ss.SoloTackles) + (float64(ss.AssistedTackles) / 2)
		ss.SoloTackles = ss.SoloTackles + stat.SoloTackles
		ss.FumblesRecovered = ss.FumblesRecovered + stat.FumblesRecovered
		ss.DefensiveSacks = ss.DefensiveSacks + stat.DefensiveSacks
		ss.ForcedFumbles = ss.ForcedFumbles + stat.ForcedFumbles
		ss.TacklesForLoss = ss.TacklesForLoss + stat.TacklesForLoss
		ss.DefensiveInterceptions = ss.DefensiveInterceptions + stat.DefensiveInterceptions
		ss.Safeties = ss.Safeties + stat.Safeties
		ss.DefensiveTDs = ss.DefensiveTDs + stat.DefensiveTDs
		ss.FieldGoalsMade = ss.FieldGoalsMade + stat.FieldGoalsMade
		ss.FieldGoalsAttempted = ss.FieldGoalsAttempted + stat.FieldGoalsAttempted
		if stat.LongestFieldGoal > ss.LongestFieldGoal {
			ss.LongestFieldGoal = stat.LongestFieldGoal
		}
		ss.ExtraPointsAttempted = ss.ExtraPointsAttempted + stat.ExtraPointsAttempted
		ss.ExtraPointsMade = ss.ExtraPointsMade + stat.ExtraPointsMade
		ss.KickoffTBs = ss.KickoffTBs + stat.KickoffTBs
		ss.Punts = ss.Punts + stat.Punts
		ss.PuntTBs = ss.PuntTBs + stat.PuntTBs
		ss.PuntsInside20 = ss.PuntsInside20 + stat.PuntsInside20
		ss.KickReturnYards = ss.KickReturnYards + stat.KickReturnYards
		ss.KickReturnTDs = ss.KickReturnTDs + stat.KickReturnTDs
		ss.PuntReturnYards = ss.PuntReturnYards + stat.PuntReturnYards
		ss.PuntReturnTDs = ss.PuntReturnTDs + stat.PuntReturnTDs
		ss.PointsScored = ss.PointsScored + stat.PointsScored
		ss.PointsAgainst = ss.PointsAgainst + stat.PointsAgainst
		ss.TwoPointTries = ss.TwoPointTries + stat.TwoPointTries
		ss.TwoPointSucceed = ss.TwoPointSucceed + stat.TwoPointSucceed
		ss.PassingCompletionsAllowed = ss.PassingCompletionsAllowed + stat.PassingCompletionsAllowed
		ss.PassingTDsAllowed = ss.PassingTDsAllowed + stat.PassingTDsAllowed
		ss.PassingYardsAllowed = ss.PassingYardsAllowed + stat.PassingYardsAllowed
		ss.RushingTDsAllowed = ss.RushingTDsAllowed + stat.RushingTDsAllowed
		ss.RushingYardsAllowed = ss.RushingYardsAllowed + stat.RushingYardsAllowed
		ss.Turnovers = ss.Turnovers + stat.DefensiveInterceptions + stat.FumblesRecovered
		ss.OffensivePenalties = ss.OffensivePenalties + stat.OffensivePenalties
		ss.DefensivePenalties = ss.DefensivePenalties + stat.DefensivePenalties
	}
	ss.TotalOffensiveYards = ss.PassingYards + ss.RushingYards
	ss.TotalYardsAllowed = ss.PassingYardsAllowed + ss.RushingYardsAllowed
	passingYards := float64(8.4) * float64(ss.PassingYards)
	passingTDs := float64(330) * float64(ss.PassingTouchdowns)
	passComps := float64(100) * float64(ss.PassingCompletions)
	ints := float64(200) * float64(ss.PassingInterceptions)
	if ss.PassingAttempts != 0 {
		numerator := passingYards + passingTDs + passComps - ints
		ss.QBRating = numerator / float64(ss.PassingAttempts)
	}
}

func (ss *NFLTeamSeasonStats) SubtractStats(stats []NFLTeamStats) {
	if ss.TeamID == 0 {
		ss.TeamID = stats[0].TeamID
		ss.SeasonID = 2
		ss.Year = 2023
	}
	for _, stat := range stats {
		ss.GamesPlayed--
		ss.PassingYards = ss.PassingYards - stat.PassingYards
		ss.PassingAttempts = ss.PassingAttempts - stat.PassingAttempts
		ss.PassingCompletions = ss.PassingCompletions - stat.PassingCompletions
		ss.PassingTouchdowns = ss.PassingTouchdowns - stat.PassingTouchdowns
		ss.PassingInterceptions = ss.PassingInterceptions - stat.PassingInterceptions
		ss.QBSacks = ss.QBSacks - stat.QBSacks
		if stat.LongestPass > ss.LongestPass {
			ss.LongestPass = stat.LongestPass
		}
		ss.RushAttempts = ss.RushAttempts - stat.RushAttempts
		ss.RushingTouchdowns = ss.RushingTouchdowns - stat.RushingTouchdowns
		ss.RushingYards = ss.RushingYards - stat.RushingYards
		ss.Fumbles = ss.Fumbles - stat.ReceivingFumbles - ss.RushingFumbles
		ss.ReceivingTargets = ss.ReceivingTargets - stat.ReceivingTargets
		ss.ReceivingCatches = ss.ReceivingCatches - stat.ReceivingCatches
		ss.ReceivingTouchdowns = ss.ReceivingTouchdowns - stat.ReceivingTouchdowns
		ss.ReceivingYards = ss.ReceivingYards - stat.ReceivingYards
		ss.AssistedTackles = ss.AssistedTackles - stat.AssistedTackles
		ss.Tackles = float64(ss.SoloTackles) - (float64(ss.AssistedTackles) / 2)
		ss.SoloTackles = ss.SoloTackles - stat.SoloTackles
		ss.FumblesRecovered = ss.FumblesRecovered - stat.FumblesRecovered
		ss.DefensiveSacks = ss.DefensiveSacks - stat.DefensiveSacks
		ss.ForcedFumbles = ss.ForcedFumbles - stat.ForcedFumbles
		ss.TacklesForLoss = ss.TacklesForLoss - stat.TacklesForLoss
		ss.DefensiveInterceptions = ss.DefensiveInterceptions - stat.DefensiveInterceptions
		ss.Safeties = ss.Safeties - stat.Safeties
		ss.DefensiveTDs = ss.DefensiveTDs - stat.DefensiveTDs
		ss.FieldGoalsMade = ss.FieldGoalsMade - stat.FieldGoalsMade
		ss.FieldGoalsAttempted = ss.FieldGoalsAttempted - stat.FieldGoalsAttempted
		if stat.LongestFieldGoal > ss.LongestFieldGoal {
			ss.LongestFieldGoal = stat.LongestFieldGoal
		}
		ss.ExtraPointsAttempted = ss.ExtraPointsAttempted - stat.ExtraPointsAttempted
		ss.ExtraPointsMade = ss.ExtraPointsMade - stat.ExtraPointsMade
		ss.KickoffTBs = ss.KickoffTBs - stat.KickoffTBs
		ss.Punts = ss.Punts - stat.Punts
		ss.PuntTBs = ss.PuntTBs - stat.PuntTBs
		ss.PuntsInside20 = ss.PuntsInside20 - stat.PuntsInside20
		ss.KickReturnYards = ss.KickReturnYards - stat.KickReturnYards
		ss.KickReturnTDs = ss.KickReturnTDs - stat.KickReturnTDs
		ss.PuntReturnYards = ss.PuntReturnYards - stat.PuntReturnYards
		ss.PuntReturnTDs = ss.PuntReturnTDs - stat.PuntReturnTDs
		ss.PointsScored = ss.PointsScored - stat.PointsScored
		ss.PointsAgainst = ss.PointsAgainst - stat.PointsAgainst
		ss.TwoPointTries = ss.TwoPointTries - stat.TwoPointTries
		ss.TwoPointSucceed = ss.TwoPointSucceed - stat.TwoPointSucceed
		ss.PassingCompletionsAllowed = ss.PassingCompletionsAllowed - stat.PassingCompletionsAllowed
		ss.PassingTDsAllowed = ss.PassingTDsAllowed - stat.PassingTDsAllowed
		ss.PassingYardsAllowed = ss.PassingYardsAllowed - stat.PassingYardsAllowed
		ss.RushingTDsAllowed = ss.RushingTDsAllowed - stat.RushingTDsAllowed
		ss.RushingYardsAllowed = ss.RushingYardsAllowed - stat.RushingYardsAllowed
		ss.Turnovers = ss.Turnovers - stat.DefensiveInterceptions - stat.FumblesRecovered
		ss.OffensivePenalties = ss.OffensivePenalties - stat.OffensivePenalties
		ss.DefensivePenalties = ss.DefensivePenalties - stat.DefensivePenalties
	}
	ss.TotalOffensiveYards = ss.PassingYards - ss.RushingYards
	ss.TotalYardsAllowed = ss.PassingYardsAllowed - ss.RushingYardsAllowed
	passingYards := float64(8.4) * float64(ss.PassingYards)
	passingTDs := float64(330) * float64(ss.PassingTouchdowns)
	passComps := float64(100) * float64(ss.PassingCompletions)
	ints := float64(200) * float64(ss.PassingInterceptions)
	if ss.PassingAttempts != 0 {
		numerator := passingYards + passingTDs + passComps - ints
		ss.QBRating = numerator / float64(ss.PassingAttempts)
	}
}
