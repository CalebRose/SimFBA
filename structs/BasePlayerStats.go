package structs

type BasePlayerStats struct {
	PassingYards         int
	PassAttempts         int
	PassCompletions      int
	PassingTDs           int
	Interceptions        int
	LongestPass          int
	Sacks                int
	RushAttempts         int
	RushingYards         int
	RushingTDs           int
	Fumbles              int
	LongestRush          int
	Targets              int
	Catches              int
	ReceivingYards       int
	ReceivingTDs         int
	LongestReception     int
	SoloTackles          float64
	AssistedTackles      float64
	TacklesForLoss       float64
	SacksMade            float64
	ForcedFumbles        int
	RecoveredFumbles     int
	PassDeflections      int
	InterceptionsCaught  int
	Safeties             int
	DefensiveTDs         int
	FGMade               int
	FGAttempts           int
	LongestFG            int
	ExtraPointsMade      int
	ExtraPointsAttempted int
	KickoffTouchbacks    int
	Punts                int
	GrossPuntDistance    int
	NetPuntDistance      int
	PuntTouchbacks       int
	PuntsInside20        int
	KickReturns          int
	KickReturnTDs        int
	KickReturnYards      int
	PuntReturns          int
	PuntReturnTDs        int
	PuntReturnYards      int
	STSoloTackles        float64
	STAssistedTackles    float64
	PuntsBlocked         int
	FGBlocked            int
	Snaps                int
	Pancakes             int
	SacksAllowed         int
	PlayedGame           int
	StartedGame          int
	WasInjured           bool
	WeeksOfRecovery      uint
	InjuryType           string
	RevealResults        bool
	TeamID               uint
	Team                 string
	PreviousTeamID       uint
	PreviousTeam         string
	GameType             uint8
}

func (bp *BasePlayerStats) MapTeamInfo(teamID uint, team string) {
	bp.TeamID = teamID
	bp.Team = team
}

func (b *BasePlayerStats) AddGameType(gameType uint8) {
	b.GameType = gameType
}
