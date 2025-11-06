package structs

type CollegePlayByPlay struct {
	PlayByPlay
}

type NFLPlayByPlay struct {
	PlayByPlay
}

type PlayersInvoledInPlay struct {
	Qb       uint
	Back1    uint
	Back2    uint
	Back3    uint
	Slot1    uint
	Slot2    uint
	Le       uint
	Re       uint
	Lt       uint
	Lg       uint
	C        uint
	Rg       uint
	Rt       uint
	Rde      uint
	Rdt      uint
	Nt       uint
	Ldt      uint
	Lde      uint
	Rolb     uint
	Rilb     uint
	Mlb      uint
	Lilb     uint
	Lolb     uint
	Rcb      uint
	Extradb1 uint
	Extradb2 uint
	Extradb3 uint
	Fs       uint
	Ss       uint
	Fcb      uint
	Lcb      uint
	Blitzer1 uint
	Blitzer2 uint
	Blitzer3 uint
}

type PlayByPlay struct {
	ID                   uint
	GameID               uint
	WeekID               uint
	SeasonID             uint
	HomeTeamID           uint
	HomeTeamScore        uint8
	AwayTeamID           uint
	AwayTeamScore        uint8
	Quarter              uint8
	HomeHasBall          bool
	TimeRemaining        string // Will always be at a maximum of 5 characters: "00:00" format
	Down                 uint8
	Distance             uint8
	LineOfScrimmage      int8  // 0 is the team with the ball, 100 is the other side of the field
	PlayTypeID           uint8 // Enum using specific IDs as reference
	PlayNameID           uint8
	OffFormationID       uint8
	OffensivePoA         uint8
	DefensiveFormationID uint8
	DefensiveTendency    uint8
	BlitzNumber          uint
	LBCoverage           bool  // 0 == Man, 1 == Zone
	CBCoverage           bool  // 0 == Man, 1 == Zone
	SCoverage            bool  // 0 == Man, 1 == Zone
	QBPlayerID           uint  // Action player is the player invoking an action. In these instances, Kickers, Punters, QBs, RBs, and FBs. WRs in some instances
	BallCarrierID        uint  // The player that attains from the result of the acting player. RBs, WRs, DBs, LBs
	Tackler1ID           uint  // ID of the player doing the tackle
	Tackler2ID           uint  // ID of the player doing the tackle
	ResultYards          int8  // Must be int due to negative values
	TurnoverPlayerID     uint  // ID of player turning over the ball
	PenaltyID            uint8 //If greater than 0, there was a penalty on the play. Could also be used as reference for the type of penalty
	PenaltyPlayerID      uint  // The player from the resulting penalty
	PenaltyYards         int8  // Resulting yards from the penalty
	InjuredPlayerID      uint  // If a player was injured, show the player
	InjuryType           uint8 // Show injury type. The UINT stands for the enum
	InjurySeverity       int8  // In the event that there are different categories of an injury type, this column is here
	InjuryDuration       int8  // Duration should by a number
	IsTouchdown          bool
	IsTurnover           bool
	IsOutOfBounds        bool
	IsFumble             bool
	IsSafety             bool
	IsINT                bool
	IsComplete           bool
	IsScramble           bool
	IsSacked             bool
	IsTouchback          bool
	IsFairCatch          bool
	IsKneel              bool
	IsBlocked            bool
	IsGood               bool
	IsLeft               bool
	IsRight              bool
	IsOffUpright         bool
	IsShort              bool
	OnOffense            bool
	KickDistance         int8
	PressureID           uint
	PlayersInvoledInPlay
}

func (p *PlayByPlay) Map(play PlayByPlayDTO) {
	p.GameID = play.GameID
	p.WeekID = play.WeekID
	p.SeasonID = play.SeasonID
	p.HomeTeamID = play.HomeTeamID
	p.AwayTeamID = play.AwayTeamID
	p.HomeTeamScore = uint8(play.HomeTeamScore)
	p.AwayTeamScore = uint8(play.AwayTeamScore)
	p.Quarter = uint8(play.Quarter)
	p.TimeRemaining = play.TimeRemaining
	p.Down = uint8(play.Down)
	p.Distance = uint8(play.Distance)
	p.LineOfScrimmage = int8(play.LineOfScrimmage)
	p.PlayTypeID = uint8(play.PlayTypeID)
	p.PlayNameID = uint8(play.PlayNameID)
	p.OffFormationID = uint8(play.OffFormationID)
	p.OffensivePoA = uint8(play.OffensivePoA)
	p.DefensiveFormationID = uint8(play.DefensiveFormationID)
	p.DefensiveTendency = uint8(play.DefensiveTendency)
	p.BlitzNumber = play.BlitzNumber
	p.LBCoverage = play.LBCoverage
	p.CBCoverage = play.CBCoverage
	p.SCoverage = play.SCoverage
	p.QBPlayerID = play.QBPlayerID
	p.BallCarrierID = play.BallCarrierID
	p.Tackler1ID = play.Tackler1ID
	p.Tackler2ID = play.Tackler2ID
	p.ResultYards = int8(play.ResultYards)
	p.TurnoverPlayerID = play.TurnoverPlayerID
	p.PenaltyID = uint8(play.PenaltyID)
	p.PenaltyPlayerID = play.PenaltyPlayerID
	p.PenaltyYards = int8(play.PenaltyYards)
	p.InjuredPlayerID = play.InjuredPlayerID
	p.InjuryType = uint8(play.InjuryType)
	p.InjurySeverity = int8(play.InjurySeverity)
	p.InjuryDuration = int8(play.InjuryDuration)
	p.IsTouchdown = play.IsTouchdown
	p.IsTurnover = play.IsTurnover
	p.IsOutOfBounds = play.IsOutOfBounds
	p.IsFumble = play.IsFumble
	p.IsSafety = play.IsSafety
	p.IsINT = play.IsINT
	p.IsScramble = play.IsScramble
	p.IsSacked = play.IsSacked
	p.IsTouchback = play.IsTouchback
	p.IsFairCatch = play.IsFairCatch
	p.IsKneel = play.IsKneel
	p.IsBlocked = play.IsBlocked
	p.IsComplete = play.IsComplete
	p.IsGood = play.IsGood
	p.IsLeft = play.IsLeft
	p.IsLeft = play.IsRight
	p.IsOffUpright = play.IsOffUpright
	p.IsShort = play.IsShort
	p.KickDistance = int8(play.KickDistance)
	p.OnOffense = play.OnOffense
	p.HomeHasBall = play.HomeHasBall
	p.PressureID = uint(play.PressureID)
	p.PlayersInvoledInPlay = play.PlayersInvoledInPlay
}

type PlayByPlayDTO struct {
	GameID               uint
	WeekID               uint
	SeasonID             uint
	HomeTeamID           uint
	HomeTeamScore        uint
	AwayTeamID           uint
	AwayTeamScore        uint
	Quarter              uint
	TimeRemaining        string // Will always be at a maximum of 5 characters: "00:00" format
	HomeHasBall          bool
	Down                 uint
	Distance             uint
	LineOfScrimmage      uint // Could make into a uint? Keep track of which side of field? Maybe as a bool? Enum?
	PlayTypeID           uint // Enum using specific IDs as reference
	PlayNameID           uint
	OffFormationID       uint
	OffensivePoA         uint
	DefensiveFormationID uint
	DefensiveTendency    uint
	BlitzNumber          uint
	LBCoverage           bool // 0 == Man, 1 == Zone
	CBCoverage           bool // 0 == Man, 1 == Zone
	SCoverage            bool // 0 == Man, 1 == Zone
	QBPlayerID           uint // Action player is the player invoking an action. In these instances, Kickers, Punters, QBs, RBs, and FBs. WRs in some instances
	BallCarrierID        uint // The player that attains from the result of the acting player. RBs, WRs, DBs, LBs
	Tackler1ID           uint // ID of the player doing the tackle
	Tackler2ID           uint // ID of the player doing the tackle
	ResultYards          int  // Must be int due to negative values
	TurnoverPlayerID     uint // ID of player who gets the ball from the defensive side
	PenaltyID            uint // If greater than 0, there was a penalty on the play. Could also be used as reference for the type of penalty
	PenaltyPlayerID      uint // The player from the resulting penalty
	PenaltyYards         int  // Resulting yards from the penalty
	InjuredPlayerID      uint // If a player was injured, show the player
	InjuryType           uint // Show injury type. The UINT stands for the enum
	InjurySeverity       uint // In the event that there are different categories of an injury type, this column is here
	InjuryDuration       int  // Duration should by a number
	IsTouchdown          bool
	IsTurnover           bool
	IsOutOfBounds        bool
	IsFumble             bool
	IsSafety             bool
	IsINT                bool
	IsComplete           bool
	IsScramble           bool
	IsSacked             bool
	IsTouchback          bool
	IsFairCatch          bool
	IsKneel              bool
	IsBlocked            bool
	IsGood               bool
	IsLeft               bool
	IsRight              bool
	IsOffUpright         bool
	IsShort              bool
	OnOffense            bool
	KickDistance         int
	PressureID           int
	PlayersInvoledInPlay
}

type PlayByPlayResponse struct {
	PlayNumber         uint
	HomeTeamID         uint
	HomeTeamScore      uint8
	AwayTeamID         uint
	AwayTeamScore      uint8
	Quarter            uint8
	Possession         string
	TimeRemaining      string
	Down               uint8
	Distance           uint8
	LineOfScrimmage    string
	PlayType           string
	PlayName           string
	OffensiveFormation string
	DefensiveFormation string
	PointOfAttack      string
	DefensiveTendency  string
	BlitzNumber        uint
	LBCoverage         string
	CBCoverage         string
	SCoverage          string
	QBPlayerID         uint
	BallCarrierID      uint // The player that attains from the result of the acting player. RBs, WRs, DBs, LBs
	Tackler1ID         uint // ID of the player doing the tackle
	Tackler2ID         uint
	PresureID          uint
	ResultYards        int8
	Result             string
	StreamResult       []string
	PlayersInvoledInPlay
}

func (p *PlayByPlayResponse) AddPlayInformation(playType, playName, offFormation, defFormation, poa string) {
	p.PlayType = playType
	p.PlayName = playName
	p.OffensiveFormation = offFormation
	p.DefensiveFormation = defFormation
	p.PointOfAttack = poa
}

func (p *PlayByPlayResponse) AddResult(result []string, isStream bool) {
	if isStream {
		p.StreamResult = result
	} else {
		p.Result = result[0]
	}
}
