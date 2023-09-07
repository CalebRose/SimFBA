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

func (ss *NFLTeamSeasonStats) ResetStats() {
	ss.GamesPlayed = 0
	ss.PassingYards = 0
	ss.PassingAttempts = 0
	ss.PassingCompletions = 0
	ss.PassingTouchdowns = 0
	ss.PassingInterceptions = 0
	ss.QBSacks = 0
	ss.LongestPass = 0
	ss.RushAttempts = 0
	ss.RushingTouchdowns = 0
	ss.RushingYards = 0
	ss.Fumbles = 0
	ss.ReceivingTargets = 0
	ss.ReceivingCatches = 0
	ss.ReceivingTouchdowns = 0
	ss.ReceivingYards = 0
	ss.AssistedTackles = 0
	ss.Tackles = 0
	ss.SoloTackles = 0
	ss.FumblesRecovered = 0
	ss.DefensiveSacks = 0
	ss.ForcedFumbles = 0
	ss.TacklesForLoss = 0
	ss.DefensiveInterceptions = 0
	ss.Safeties = 0
	ss.DefensiveTDs = 0
	ss.FieldGoalsMade = 0
	ss.FieldGoalsAttempted = 0
	ss.LongestFieldGoal = 0
	ss.ExtraPointsAttempted = 0
	ss.ExtraPointsMade = 0
	ss.KickoffTBs = 0
	ss.Punts = 0
	ss.PuntTBs = 0
	ss.PuntsInside20 = 0
	ss.KickReturnYards = 0
	ss.KickReturnTDs = 0
	ss.PuntReturnYards = 0
	ss.PuntReturnTDs = 0
	ss.PointsScored = 0
	ss.PointsAgainst = 0
	ss.TwoPointTries = 0
	ss.TwoPointSucceed = 0
	ss.PassingCompletionsAllowed = 0
	ss.PassingTDsAllowed = 0
	ss.PassingYardsAllowed = 0
	ss.RushingTDsAllowed = 0
	ss.RushingYardsAllowed = 0
	ss.Turnovers = 0
	ss.OffensivePenalties = 0
	ss.DefensivePenalties = 0
	ss.TotalOffensiveYards = 0
	ss.TotalYardsAllowed = 0
	ss.QBRating = 0
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
