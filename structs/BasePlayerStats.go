package structs

type BasePlayerStats struct {
	PassingYards         int16
	PassAttempts         int16
	PassCompletions      int16
	PassingTDs           int16
	Interceptions        int16
	LongestPass          int16
	Sacks                int16
	RushAttempts         int16
	RushingYards         int16
	RushingTDs           int16
	Fumbles              int16
	LongestRush          int16
	Targets              int16
	Catches              int16
	ReceivingYards       int16
	ReceivingTDs         int16
	LongestReception     int16
	SoloTackles          float64
	AssistedTackles      float64
	TacklesForLoss       float64
	SacksMade            float64
	ForcedFumbles        int16
	RecoveredFumbles     int16
	PassDeflections      int16
	InterceptionsCaught  int16
	Safeties             int16
	DefensiveTDs         int16
	FGMade               int16
	FGAttempts           int16
	LongestFG            int16
	ExtraPointsMade      int16
	ExtraPointsAttempted int16
	KickoffTouchbacks    int16
	Punts                int16
	GrossPuntDistance    int16
	NetPuntDistance      int16
	PuntTouchbacks       int16
	PuntsInside20        int16
	KickReturns          int16
	KickReturnTDs        int16
	KickReturnYards      int16
	PuntReturns          int16
	PuntReturnTDs        int16
	PuntReturnYards      int16
	STSoloTackles        float64
	STAssistedTackles    float64
	PuntsBlocked         int16
	FGBlocked            int16
	Snaps                int16
	Pancakes             int16
	SacksAllowed         int16
	PlayedGame           int16
	StartedGame          int16
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
