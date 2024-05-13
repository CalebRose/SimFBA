package structs

type BaseTeamStats struct {
	PointsScored                  int
	PointsAgainst                 int
	TwoPointTries                 int
	TwoPointSucceed               int
	PassingYards                  int
	PassingAttempts               int
	PassingCompletions            int
	LongestPass                   int
	PassingTouchdowns             int
	PassingInterceptions          int
	QBRating                      int
	QBSacks                       int
	RushAttempts                  int
	RushingYards                  int
	RushingYardsPerAttempt        float64
	LongestRush                   int
	RushingTouchdowns             int
	RushingFumbles                int
	ReceivingTargets              int
	ReceivingCatches              int
	ReceivingYards                int
	YardsPerCatch                 float64
	ReceivingTouchdowns           int
	ReceivingFumbles              int
	PassingYardsAllowed           int
	PassingTDsAllowed             int
	PassingCompletionsAllowed     int
	RushingYardsAllowed           int
	RushingTDsAllowed             int
	RushingYardsPerAttemptAllowed float64
	SoloTackles                   int
	AssistedTackles               int
	TacklesForLoss                float64
	DefensiveSacks                float64
	ForcedFumbles                 int
	FumblesRecovered              int
	DefensiveInterceptions        int
	TurnoverYards                 int
	DefensiveTDs                  int
	Safeties                      int
	ExtraPointsMade               int
	ExtraPointsAttempted          int
	ExtraPointPercentage          float64
	FieldGoalsMade                int
	FieldGoalsAttempted           int
	FieldGoalsPercentage          float64
	LongestFieldGoal              int
	KickoffTBs                    int
	PuntTBs                       int
	Punts                         int
	PuntYards                     int
	PuntsInside20                 int
	PuntAverage                   float64
	Inside20YardLine              int
	KickReturnYards               int
	KickReturnTDs                 int
	PuntReturnYards               int
	PuntReturnTDs                 int
	OffensivePenalties            int
	DefensivePenalties            int
	OffPenaltyYards               int
	DefPenaltyYards               int
	Score1Q                       int
	Score2Q                       int
	Score3Q                       int
	Score4Q                       int
	Score5Q                       int
	Score6Q                       int
	Score7Q                       int
	ScoreOT                       int
	OffensiveScheme               string
	DefensiveScheme               string
	RevealResults                 bool
}
