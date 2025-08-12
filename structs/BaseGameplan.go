package structs

type GamePlanResponse struct {
	CollegeGP         CollegeGameplan
	NFLGP             NFLGameplan
	CollegeDC         CollegeTeamDepthChart
	NFLDC             NFLDepthChart
	OpponentScheme    string
	CollegeOppPlayers []CollegePlayer
	NFLOppPlayers     []NFLPlayer
}

type BaseGameplan struct {
	OffensiveScheme   string
	OffRunToPassRatio int
	OffensiveFormation
	DefensiveScheme string
	DefensiveFormation
	BlitzSafeties        bool
	BlitzCorners         bool
	LinebackerCoverage   string
	CornersCoverage      string
	SafetiesCoverage     string
	DiveFocus            int
	PitchFocus           int
	PrimaryHB            int
	MaximumFGDistance    int
	GoFor4AndShort       int
	GoFor4AndLong        int
	HasSchemePenalty     bool
	OffenseSchemePenalty uint
	DefenseSchemePenalty uint
	DefaultOffense       bool
	DefaultDefense       bool
	PreviousWeekBye      bool
	FocusPlays           string
	DoubleTeam           int
}

func (cg *BaseGameplan) ApplySchemePenalty(IsOffense bool, diff int) {
	cg.HasSchemePenalty = true
	if IsOffense {
		cg.OffenseSchemePenalty = 4 - uint(diff)
		if cg.OffenseSchemePenalty > 3 {
			cg.OffenseSchemePenalty = 3
		}
	} else {
		cg.DefenseSchemePenalty = 4 - uint(diff)
		if cg.DefenseSchemePenalty > 3 {
			cg.DefenseSchemePenalty = 3
		}
	}
}

func (cg *BaseGameplan) LowerPenalty() {
	if cg.OffenseSchemePenalty > 0 {
		cg.OffenseSchemePenalty--
	}
	if cg.DefenseSchemePenalty > 0 {
		cg.DefenseSchemePenalty--
	}
	if cg.OffenseSchemePenalty == 0 && cg.DefenseSchemePenalty == 0 {
		cg.HasSchemePenalty = false
	}
}

func (bg *BaseGameplan) UpdateCollegeGameplan(dto CollegeGameplan) {
	// Validation is done in UI, so we're just passing data along in API
	bg.OffensiveScheme = dto.OffensiveScheme
	bg.OffRunToPassRatio = dto.OffRunToPassRatio
	bg.OffFormation1Name = dto.OffFormation1Name
	bg.OffForm1Weight = dto.OffForm1Weight
	bg.OffForm1TraditionalRun = dto.OffForm1TraditionalRun
	bg.OffForm1OptionRun = dto.OffForm1OptionRun
	bg.OffForm1RPO = dto.OffForm1RPO
	bg.OffForm1Pass = dto.OffForm1Pass
	bg.OffFormation2Name = dto.OffFormation2Name
	bg.OffForm2Weight = dto.OffForm2Weight
	bg.OffForm2TraditionalRun = dto.OffForm2TraditionalRun
	bg.OffForm2OptionRun = dto.OffForm2OptionRun
	bg.OffForm2RPO = dto.OffForm2RPO
	bg.OffForm2Pass = dto.OffForm2Pass
	bg.OffFormation3Name = dto.OffFormation3Name
	bg.OffForm3Weight = dto.OffForm3Weight
	bg.OffForm3TraditionalRun = dto.OffForm3TraditionalRun
	bg.OffForm3OptionRun = dto.OffForm3OptionRun
	bg.OffForm3RPO = dto.OffForm3RPO
	bg.OffForm3Pass = dto.OffForm3Pass
	bg.OffFormation4Name = dto.OffFormation4Name
	bg.OffForm4Weight = dto.OffForm4Weight
	bg.OffForm4TraditionalRun = dto.OffForm4TraditionalRun
	bg.OffForm4OptionRun = dto.OffForm4OptionRun
	bg.OffForm4RPO = dto.OffForm4RPO
	bg.OffForm4Pass = dto.OffForm4Pass
	bg.OffFormation5Name = dto.OffFormation5Name
	bg.OffForm5Weight = dto.OffForm5Weight
	bg.OffForm5TraditionalRun = dto.OffForm5TraditionalRun
	bg.OffForm5OptionRun = dto.OffForm5OptionRun
	bg.OffForm5RPO = dto.OffForm5RPO
	bg.OffForm5Pass = dto.OffForm5Pass
	bg.RunnerDistributionQB = dto.RunnerDistributionQB
	bg.RunnerDistributionRB1 = dto.RunnerDistributionRB1
	bg.RunnerDistributionRB2 = dto.RunnerDistributionRB2
	bg.RunnerDistributionRB3 = dto.RunnerDistributionRB3
	bg.RunnerDistributionFB1 = dto.RunnerDistributionFB1
	bg.RunnerDistributionFB2 = dto.RunnerDistributionFB2
	bg.RunnerDistributionWR = dto.RunnerDistributionWR
	bg.RunnerDistributionWRPosition = dto.RunnerDistributionWRPosition
	bg.RunOutsideLeft = dto.RunOutsideLeft
	bg.RunOutsideRight = dto.RunOutsideRight
	bg.RunInsideLeft = dto.RunInsideLeft
	bg.RunInsideRight = dto.RunInsideRight
	bg.RunPowerLeft = dto.RunPowerLeft
	bg.RunPowerRight = dto.RunPowerRight
	bg.RunDrawLeft = dto.RunDrawLeft
	bg.RunDrawRight = dto.RunDrawRight
	bg.ReadOptionLeft = dto.ReadOptionLeft
	bg.ReadOptionRight = dto.ReadOptionRight
	bg.SpeedOptionLeft = dto.SpeedOptionLeft
	bg.SpeedOptionRight = dto.SpeedOptionRight
	bg.InvertedOptionLeft = dto.InvertedOptionLeft
	bg.InvertedOptionRight = dto.InvertedOptionRight
	bg.TripleOptionLeft = dto.TripleOptionLeft
	bg.TripleOptionRight = dto.TripleOptionRight
	bg.PassQuick = dto.PassQuick
	bg.PassShort = dto.PassShort
	bg.PassLong = dto.PassLong
	bg.PassDeep = dto.PassDeep
	bg.PassScreen = dto.PassScreen
	bg.PassPAShort = dto.PassPAShort
	bg.PassPALong = dto.PassPALong
	bg.PassPADeep = dto.PassPADeep
	bg.LeftVsRight = dto.LeftVsRight
	bg.ChoiceOutside = dto.ChoiceOutside
	bg.ChoiceInside = dto.ChoiceInside
	bg.ChoicePower = dto.ChoicePower
	bg.PeekOutside = dto.PeekOutside
	bg.PeekInside = dto.PeekInside
	bg.PeekPower = dto.PeekPower
	bg.TargetingWR1 = dto.TargetingWR1
	bg.TargetDepthWR1 = dto.TargetDepthWR1
	bg.TargetingWR2 = dto.TargetingWR2
	bg.TargetDepthWR2 = dto.TargetDepthWR2
	bg.TargetingWR3 = dto.TargetingWR3
	bg.TargetDepthWR3 = dto.TargetDepthWR3
	bg.TargetingWR4 = dto.TargetingWR4
	bg.TargetDepthWR4 = dto.TargetDepthWR4
	bg.TargetingWR5 = dto.TargetingWR5
	bg.TargetDepthWR5 = dto.TargetDepthWR5
	bg.TargetingTE1 = dto.TargetingTE1
	bg.TargetDepthTE1 = dto.TargetDepthTE1
	bg.TargetingTE2 = dto.TargetingTE2
	bg.TargetDepthTE2 = dto.TargetDepthTE2
	bg.TargetingTE3 = dto.TargetingTE3
	bg.TargetDepthTE3 = dto.TargetDepthTE3
	bg.TargetingRB1 = dto.TargetingRB1
	bg.TargetDepthRB1 = dto.TargetDepthRB1
	bg.TargetingRB2 = dto.TargetingRB2
	bg.TargetDepthRB2 = dto.TargetDepthRB2
	bg.TargetingFB1 = dto.TargetingFB1
	bg.TargetDepthFB1 = dto.TargetDepthFB1
	bg.DefensiveScheme = dto.DefensiveScheme
	bg.DefFormation1 = dto.DefFormation1
	bg.DefFormation1RunToPass = dto.DefFormation1RunToPass
	bg.DefFormation1BlitzWeight = dto.DefFormation1BlitzWeight
	bg.DefFormation1BlitzAggression = dto.DefFormation1BlitzAggression
	bg.DefFormation2 = dto.DefFormation2
	bg.DefFormation2RunToPass = dto.DefFormation2RunToPass
	bg.DefFormation2BlitzWeight = dto.DefFormation2BlitzWeight
	bg.DefFormation2BlitzAggression = dto.DefFormation2BlitzAggression
	bg.DefFormation3 = dto.DefFormation3
	bg.DefFormation3RunToPass = dto.DefFormation3RunToPass
	bg.DefFormation3BlitzWeight = dto.DefFormation3BlitzWeight
	bg.DefFormation3BlitzAggression = dto.DefFormation3BlitzAggression
	bg.DefFormation4 = dto.DefFormation4
	bg.DefFormation4RunToPass = dto.DefFormation4RunToPass
	bg.DefFormation4BlitzWeight = dto.DefFormation4BlitzWeight
	bg.DefFormation4BlitzAggression = dto.DefFormation4BlitzAggression
	bg.DefFormation5 = dto.DefFormation5
	bg.DefFormation5RunToPass = dto.DefFormation5RunToPass
	bg.DefFormation5BlitzWeight = dto.DefFormation5BlitzWeight
	bg.DefFormation5BlitzAggression = dto.DefFormation5BlitzAggression
	bg.BlitzSafeties = dto.BlitzSafeties
	bg.BlitzCorners = dto.BlitzCorners
	if len(dto.LinebackerCoverage) > 0 {
		bg.LinebackerCoverage = dto.LinebackerCoverage
	}
	if len(dto.CornersCoverage) > 0 {
		bg.CornersCoverage = dto.CornersCoverage
	}
	if len(dto.SafetiesCoverage) > 0 {
		bg.SafetiesCoverage = dto.SafetiesCoverage
	}
	bg.PrimaryHB = dto.PrimaryHB
	bg.DiveFocus = dto.DiveFocus
	bg.PitchFocus = dto.PitchFocus
	bg.MaximumFGDistance = dto.MaximumFGDistance
	bg.GoFor4AndLong = dto.GoFor4AndLong
	bg.GoFor4AndShort = dto.GoFor4AndShort
	bg.DefaultOffense = dto.DefaultOffense
	bg.DefaultDefense = dto.DefaultDefense
	bg.PreviousWeekBye = dto.PreviousWeekBye
	bg.FocusPlays = dto.FocusPlays
	bg.DoubleTeam = dto.DoubleTeam
}

func (bg *BaseGameplan) UpdateNFLGameplan(dto NFLGameplan) {
	// Validation is done in UI, so we're just passing data along in API
	bg.OffensiveScheme = dto.OffensiveScheme
	bg.OffRunToPassRatio = dto.OffRunToPassRatio
	bg.OffFormation1Name = dto.OffFormation1Name
	bg.OffForm1Weight = dto.OffForm1Weight
	bg.OffForm1TraditionalRun = dto.OffForm1TraditionalRun
	bg.OffForm1OptionRun = dto.OffForm1OptionRun
	bg.OffForm1RPO = dto.OffForm1RPO
	bg.OffForm1Pass = dto.OffForm1Pass
	bg.OffFormation2Name = dto.OffFormation2Name
	bg.OffForm2Weight = dto.OffForm2Weight
	bg.OffForm2TraditionalRun = dto.OffForm2TraditionalRun
	bg.OffForm2OptionRun = dto.OffForm2OptionRun
	bg.OffForm2RPO = dto.OffForm2RPO
	bg.OffForm2Pass = dto.OffForm2Pass
	bg.OffFormation3Name = dto.OffFormation3Name
	bg.OffForm3Weight = dto.OffForm3Weight
	bg.OffForm3TraditionalRun = dto.OffForm3TraditionalRun
	bg.OffForm3OptionRun = dto.OffForm3OptionRun
	bg.OffForm3RPO = dto.OffForm3RPO
	bg.OffForm3Pass = dto.OffForm3Pass
	bg.OffFormation4Name = dto.OffFormation4Name
	bg.OffForm4Weight = dto.OffForm4Weight
	bg.OffForm4TraditionalRun = dto.OffForm4TraditionalRun
	bg.OffForm4OptionRun = dto.OffForm4OptionRun
	bg.OffForm4RPO = dto.OffForm4RPO
	bg.OffForm4Pass = dto.OffForm4Pass
	bg.OffFormation5Name = dto.OffFormation5Name
	bg.OffForm5Weight = dto.OffForm5Weight
	bg.OffForm5TraditionalRun = dto.OffForm5TraditionalRun
	bg.OffForm5OptionRun = dto.OffForm5OptionRun
	bg.OffForm5RPO = dto.OffForm5RPO
	bg.OffForm5Pass = dto.OffForm5Pass
	bg.RunnerDistributionQB = dto.RunnerDistributionQB
	bg.RunnerDistributionRB1 = dto.RunnerDistributionRB1
	bg.RunnerDistributionRB2 = dto.RunnerDistributionRB2
	bg.RunnerDistributionRB3 = dto.RunnerDistributionRB3
	bg.RunnerDistributionFB1 = dto.RunnerDistributionFB1
	bg.RunnerDistributionFB2 = dto.RunnerDistributionFB2
	bg.RunnerDistributionWR = dto.RunnerDistributionWR
	bg.RunnerDistributionWRPosition = dto.RunnerDistributionWRPosition
	bg.RunOutsideLeft = dto.RunOutsideLeft
	bg.RunOutsideRight = dto.RunOutsideRight
	bg.RunInsideLeft = dto.RunInsideLeft
	bg.RunInsideRight = dto.RunInsideRight
	bg.RunPowerLeft = dto.RunPowerLeft
	bg.RunPowerRight = dto.RunPowerRight
	bg.RunDrawLeft = dto.RunDrawLeft
	bg.RunDrawRight = dto.RunDrawRight
	bg.ReadOptionLeft = dto.ReadOptionLeft
	bg.ReadOptionRight = dto.ReadOptionRight
	bg.SpeedOptionLeft = dto.SpeedOptionLeft
	bg.SpeedOptionRight = dto.SpeedOptionRight
	bg.InvertedOptionLeft = dto.InvertedOptionLeft
	bg.InvertedOptionRight = dto.InvertedOptionRight
	bg.TripleOptionLeft = dto.TripleOptionLeft
	bg.TripleOptionRight = dto.TripleOptionRight
	bg.PassQuick = dto.PassQuick
	bg.PassShort = dto.PassShort
	bg.PassLong = dto.PassLong
	bg.PassDeep = dto.PassDeep
	bg.PassScreen = dto.PassScreen
	bg.PassPAShort = dto.PassPAShort
	bg.PassPALong = dto.PassPALong
	bg.PassPADeep = dto.PassPADeep
	bg.LeftVsRight = dto.LeftVsRight
	bg.ChoiceOutside = dto.ChoiceOutside
	bg.ChoiceInside = dto.ChoiceInside
	bg.ChoicePower = dto.ChoicePower
	bg.PeekOutside = dto.PeekOutside
	bg.PeekInside = dto.PeekInside
	bg.PeekPower = dto.PeekPower
	bg.TargetingWR1 = dto.TargetingWR1
	bg.TargetDepthWR1 = dto.TargetDepthWR1
	bg.TargetingWR2 = dto.TargetingWR2
	bg.TargetDepthWR2 = dto.TargetDepthWR2
	bg.TargetingWR3 = dto.TargetingWR3
	bg.TargetDepthWR3 = dto.TargetDepthWR3
	bg.TargetingWR4 = dto.TargetingWR4
	bg.TargetDepthWR4 = dto.TargetDepthWR4
	bg.TargetingWR5 = dto.TargetingWR5
	bg.TargetDepthWR5 = dto.TargetDepthWR5
	bg.TargetingTE1 = dto.TargetingTE1
	bg.TargetDepthTE1 = dto.TargetDepthTE1
	bg.TargetingTE2 = dto.TargetingTE2
	bg.TargetDepthTE2 = dto.TargetDepthTE2
	bg.TargetingTE3 = dto.TargetingTE3
	bg.TargetDepthTE3 = dto.TargetDepthTE3
	bg.TargetingRB1 = dto.TargetingRB1
	bg.TargetDepthRB1 = dto.TargetDepthRB1
	bg.TargetingRB2 = dto.TargetingRB2
	bg.TargetDepthRB2 = dto.TargetDepthRB2
	bg.TargetingFB1 = dto.TargetingFB1
	bg.TargetDepthFB1 = dto.TargetDepthFB1
	bg.DefensiveScheme = dto.DefensiveScheme
	bg.DefFormation1 = dto.DefFormation1
	bg.DefFormation1RunToPass = dto.DefFormation1RunToPass
	bg.DefFormation1BlitzWeight = dto.DefFormation1BlitzWeight
	bg.DefFormation1BlitzAggression = dto.DefFormation1BlitzAggression
	bg.DefFormation2 = dto.DefFormation2
	bg.DefFormation2RunToPass = dto.DefFormation2RunToPass
	bg.DefFormation2BlitzWeight = dto.DefFormation2BlitzWeight
	bg.DefFormation2BlitzAggression = dto.DefFormation2BlitzAggression
	bg.DefFormation3 = dto.DefFormation3
	bg.DefFormation3RunToPass = dto.DefFormation3RunToPass
	bg.DefFormation3BlitzWeight = dto.DefFormation3BlitzWeight
	bg.DefFormation3BlitzAggression = dto.DefFormation3BlitzAggression
	bg.DefFormation4 = dto.DefFormation4
	bg.DefFormation4RunToPass = dto.DefFormation4RunToPass
	bg.DefFormation4BlitzWeight = dto.DefFormation4BlitzWeight
	bg.DefFormation4BlitzAggression = dto.DefFormation4BlitzAggression
	bg.DefFormation5 = dto.DefFormation5
	bg.DefFormation5RunToPass = dto.DefFormation5RunToPass
	bg.DefFormation5BlitzWeight = dto.DefFormation5BlitzWeight
	bg.DefFormation5BlitzAggression = dto.DefFormation5BlitzAggression
	bg.BlitzSafeties = dto.BlitzSafeties
	bg.BlitzCorners = dto.BlitzCorners
	bg.LinebackerCoverage = dto.LinebackerCoverage
	bg.CornersCoverage = dto.CornersCoverage
	bg.SafetiesCoverage = dto.SafetiesCoverage
	bg.PrimaryHB = dto.PrimaryHB
	bg.DiveFocus = dto.DiveFocus
	bg.PitchFocus = dto.PitchFocus
	bg.MaximumFGDistance = dto.MaximumFGDistance
	bg.GoFor4AndLong = dto.GoFor4AndLong
	bg.GoFor4AndShort = dto.GoFor4AndShort
	bg.DefaultOffense = dto.DefaultOffense
	bg.DefaultDefense = dto.DefaultDefense
	bg.DoubleTeam = dto.DoubleTeam
	bg.FocusPlays = dto.FocusPlays
}

func (bg *BaseGameplan) UpdateSchemes(off, def string) {
	bg.OffensiveScheme = off
	bg.DefensiveScheme = def
	bg.DefaultOffense = true
	bg.DefaultDefense = true
}

func (bg *BaseGameplan) AssignRunnerWRID(id uint) {
	bg.RunnerDistributionWRID = id
}

type OffensiveFormation struct {
	OffFormation1Name            string
	OffForm1Weight               int
	OffForm1TraditionalRun       int
	OffForm1OptionRun            int
	OffForm1Pass                 int
	OffForm1RPO                  int
	OffFormation2Name            string
	OffForm2Weight               int
	OffForm2TraditionalRun       int
	OffForm2OptionRun            int
	OffForm2Pass                 int
	OffForm2RPO                  int
	OffFormation3Name            string
	OffForm3Weight               int
	OffForm3TraditionalRun       int
	OffForm3OptionRun            int
	OffForm3Pass                 int
	OffForm3RPO                  int
	OffFormation4Name            string
	OffForm4Weight               int
	OffForm4TraditionalRun       int
	OffForm4OptionRun            int
	OffForm4Pass                 int
	OffForm4RPO                  int
	OffFormation5Name            string
	OffForm5Weight               int
	OffForm5TraditionalRun       int
	OffForm5OptionRun            int
	OffForm5Pass                 int
	OffForm5RPO                  int
	RunnerDistributionQB         int
	RunnerDistributionRB1        int
	RunnerDistributionRB2        int
	RunnerDistributionRB3        int
	RunnerDistributionFB1        int
	RunnerDistributionFB2        int
	RunnerDistributionWR         int    // Jet Sweep stuff
	RunnerDistributionWRPosition string // WR1, WR2, WR3, WR4, WR5
	RunnerDistributionWRID       uint
	RunOutsideLeft               int
	RunOutsideRight              int
	RunInsideLeft                int
	RunInsideRight               int
	RunPowerLeft                 int
	RunPowerRight                int
	RunDrawLeft                  int
	RunDrawRight                 int
	ReadOptionLeft               int
	ReadOptionRight              int
	SpeedOptionLeft              int
	SpeedOptionRight             int
	InvertedOptionLeft           int
	InvertedOptionRight          int
	TripleOptionLeft             int
	TripleOptionRight            int
	PassQuick                    int // Now Pass Short
	PassShort                    int // Now Pass Medium
	PassLong                     int
	PassDeep                     int
	PassScreen                   int
	PassPAShort                  int
	PassPAMedium                 int
	PassPALong                   int
	PassPADeep                   int
	LeftVsRight                  int
	ChoiceOutside                int
	ChoiceInside                 int
	ChoicePower                  int
	PeekOutside                  int
	PeekInside                   int
	PeekPower                    int
	TargetingWR1                 int
	TargetDepthWR1               string // Quick, Short, Long, None
	TargetingWR2                 int
	TargetDepthWR2               string
	TargetingWR3                 int
	TargetDepthWR3               string
	TargetingWR4                 int
	TargetDepthWR4               string
	TargetingWR5                 int
	TargetDepthWR5               string
	TargetingTE1                 int
	TargetDepthTE1               string
	TargetingTE2                 int
	TargetDepthTE2               string
	TargetingTE3                 int
	TargetDepthTE3               string
	TargetingRB1                 int
	TargetDepthRB1               string
	TargetingRB2                 int
	TargetDepthRB2               string
	TargetingFB1                 int
	TargetDepthFB1               string
}

type DefensiveFormation struct {
	DefFormation1                string
	DefFormation1RunToPass       int
	DefFormation1BlitzWeight     int
	DefFormation1BlitzAggression string
	DefFormation2                string
	DefFormation2RunToPass       int
	DefFormation2BlitzWeight     int
	DefFormation2BlitzAggression string
	DefFormation3                string
	DefFormation3RunToPass       int
	DefFormation3BlitzWeight     int
	DefFormation3BlitzAggression string
	DefFormation4                string
	DefFormation4RunToPass       int
	DefFormation4BlitzWeight     int
	DefFormation4BlitzAggression string
	DefFormation5                string
	DefFormation5RunToPass       int
	DefFormation5BlitzWeight     int
	DefFormation5BlitzAggression string
}

func (bg *BaseGameplan) UpdateCollegeGameplanTEST(dto CollegeGameplanTEST) {
	// Validation is done in UI, so we're just passing data along in API
	bg.OffensiveScheme = dto.OffensiveScheme
	bg.OffRunToPassRatio = dto.OffRunToPassRatio
	bg.OffFormation1Name = dto.OffFormation1Name
	bg.OffForm1Weight = dto.OffForm1Weight
	bg.OffForm1TraditionalRun = dto.OffForm1TraditionalRun
	bg.OffForm1OptionRun = dto.OffForm1OptionRun
	bg.OffForm1RPO = dto.OffForm1RPO
	bg.OffForm1Pass = dto.OffForm1Pass
	bg.OffFormation2Name = dto.OffFormation2Name
	bg.OffForm2Weight = dto.OffForm2Weight
	bg.OffForm2TraditionalRun = dto.OffForm2TraditionalRun
	bg.OffForm2OptionRun = dto.OffForm2OptionRun
	bg.OffForm2RPO = dto.OffForm2RPO
	bg.OffForm2Pass = dto.OffForm2Pass
	bg.OffFormation3Name = dto.OffFormation3Name
	bg.OffForm3Weight = dto.OffForm3Weight
	bg.OffForm3TraditionalRun = dto.OffForm3TraditionalRun
	bg.OffForm3OptionRun = dto.OffForm3OptionRun
	bg.OffForm3RPO = dto.OffForm3RPO
	bg.OffForm3Pass = dto.OffForm3Pass
	bg.OffFormation4Name = dto.OffFormation4Name
	bg.OffForm4Weight = dto.OffForm4Weight
	bg.OffForm4TraditionalRun = dto.OffForm4TraditionalRun
	bg.OffForm4OptionRun = dto.OffForm4OptionRun
	bg.OffForm4RPO = dto.OffForm4RPO
	bg.OffForm4Pass = dto.OffForm4Pass
	bg.OffFormation5Name = dto.OffFormation5Name
	bg.OffForm5Weight = dto.OffForm5Weight
	bg.OffForm5TraditionalRun = dto.OffForm5TraditionalRun
	bg.OffForm5OptionRun = dto.OffForm5OptionRun
	bg.OffForm5RPO = dto.OffForm5RPO
	bg.OffForm5Pass = dto.OffForm5Pass
	bg.RunnerDistributionQB = dto.RunnerDistributionQB
	bg.RunnerDistributionRB1 = dto.RunnerDistributionRB1
	bg.RunnerDistributionRB2 = dto.RunnerDistributionRB2
	bg.RunnerDistributionRB3 = dto.RunnerDistributionRB3
	bg.RunnerDistributionFB1 = dto.RunnerDistributionFB1
	bg.RunnerDistributionFB2 = dto.RunnerDistributionFB2
	bg.RunnerDistributionWR = dto.RunnerDistributionWR
	bg.RunnerDistributionWRPosition = dto.RunnerDistributionWRPosition
	bg.RunOutsideLeft = dto.RunOutsideLeft
	bg.RunOutsideRight = dto.RunOutsideRight
	bg.RunInsideLeft = dto.RunInsideLeft
	bg.RunInsideRight = dto.RunInsideRight
	bg.RunPowerLeft = dto.RunPowerLeft
	bg.RunPowerRight = dto.RunPowerRight
	bg.RunDrawLeft = dto.RunDrawLeft
	bg.RunDrawRight = dto.RunDrawRight
	bg.ReadOptionLeft = dto.ReadOptionLeft
	bg.ReadOptionRight = dto.ReadOptionRight
	bg.SpeedOptionLeft = dto.SpeedOptionLeft
	bg.SpeedOptionRight = dto.SpeedOptionRight
	bg.InvertedOptionLeft = dto.InvertedOptionLeft
	bg.InvertedOptionRight = dto.InvertedOptionRight
	bg.TripleOptionLeft = dto.TripleOptionLeft
	bg.TripleOptionRight = dto.TripleOptionRight
	bg.PassQuick = dto.PassQuick
	bg.PassShort = dto.PassShort
	bg.PassLong = dto.PassLong
	bg.PassDeep = dto.PassDeep
	bg.PassScreen = dto.PassScreen
	bg.PassPAShort = dto.PassPAShort
	bg.PassPAMedium = dto.PassPAMedium
	bg.PassPALong = dto.PassPALong
	bg.PassPADeep = dto.PassPADeep
	bg.LeftVsRight = dto.LeftVsRight
	bg.ChoiceOutside = dto.ChoiceOutside
	bg.ChoiceInside = dto.ChoiceInside
	bg.ChoicePower = dto.ChoicePower
	bg.PeekOutside = dto.PeekOutside
	bg.PeekInside = dto.PeekInside
	bg.PeekPower = dto.PeekPower
	bg.TargetingWR1 = dto.TargetingWR1
	bg.TargetDepthWR1 = dto.TargetDepthWR1
	bg.TargetingWR2 = dto.TargetingWR2
	bg.TargetDepthWR2 = dto.TargetDepthWR2
	bg.TargetingWR3 = dto.TargetingWR3
	bg.TargetDepthWR3 = dto.TargetDepthWR3
	bg.TargetingWR4 = dto.TargetingWR4
	bg.TargetDepthWR4 = dto.TargetDepthWR4
	bg.TargetingWR5 = dto.TargetingWR5
	bg.TargetDepthWR5 = dto.TargetDepthWR5
	bg.TargetingTE1 = dto.TargetingTE1
	bg.TargetDepthTE1 = dto.TargetDepthTE1
	bg.TargetingTE2 = dto.TargetingTE2
	bg.TargetDepthTE2 = dto.TargetDepthTE2
	bg.TargetingTE3 = dto.TargetingTE3
	bg.TargetDepthTE3 = dto.TargetDepthTE3
	bg.TargetingRB1 = dto.TargetingRB1
	bg.TargetDepthRB1 = dto.TargetDepthRB1
	bg.TargetingRB2 = dto.TargetingRB2
	bg.TargetDepthRB2 = dto.TargetDepthRB2
	bg.TargetingFB1 = dto.TargetingFB1
	bg.TargetDepthFB1 = dto.TargetDepthFB1
	bg.DefensiveScheme = dto.DefensiveScheme
	bg.DefFormation1 = dto.DefFormation1
	bg.DefFormation1RunToPass = dto.DefFormation1RunToPass
	bg.DefFormation1BlitzWeight = dto.DefFormation1BlitzWeight
	bg.DefFormation1BlitzAggression = dto.DefFormation1BlitzAggression
	bg.DefFormation2 = dto.DefFormation2
	bg.DefFormation2RunToPass = dto.DefFormation2RunToPass
	bg.DefFormation2BlitzWeight = dto.DefFormation2BlitzWeight
	bg.DefFormation2BlitzAggression = dto.DefFormation2BlitzAggression
	bg.DefFormation3 = dto.DefFormation3
	bg.DefFormation3RunToPass = dto.DefFormation3RunToPass
	bg.DefFormation3BlitzWeight = dto.DefFormation3BlitzWeight
	bg.DefFormation3BlitzAggression = dto.DefFormation3BlitzAggression
	bg.DefFormation4 = dto.DefFormation4
	bg.DefFormation4RunToPass = dto.DefFormation4RunToPass
	bg.DefFormation4BlitzWeight = dto.DefFormation4BlitzWeight
	bg.DefFormation4BlitzAggression = dto.DefFormation4BlitzAggression
	bg.DefFormation5 = dto.DefFormation5
	bg.DefFormation5RunToPass = dto.DefFormation5RunToPass
	bg.DefFormation5BlitzWeight = dto.DefFormation5BlitzWeight
	bg.DefFormation5BlitzAggression = dto.DefFormation5BlitzAggression
	bg.BlitzSafeties = dto.BlitzSafeties
	bg.BlitzCorners = dto.BlitzCorners
	if len(dto.LinebackerCoverage) > 0 {
		bg.LinebackerCoverage = dto.LinebackerCoverage
	}
	if len(dto.CornersCoverage) > 0 {
		bg.CornersCoverage = dto.CornersCoverage
	}
	if len(dto.SafetiesCoverage) > 0 {
		bg.SafetiesCoverage = dto.SafetiesCoverage
	}
	bg.PrimaryHB = dto.PrimaryHB
	bg.DiveFocus = dto.DiveFocus
	bg.PitchFocus = dto.PitchFocus
	bg.MaximumFGDistance = dto.MaximumFGDistance
	bg.GoFor4AndLong = dto.GoFor4AndLong
	bg.GoFor4AndShort = dto.GoFor4AndShort
	bg.DefaultOffense = dto.DefaultOffense
	bg.DefaultDefense = dto.DefaultDefense
	bg.PreviousWeekBye = dto.PreviousWeekBye
	bg.FocusPlays = dto.FocusPlays
	bg.DoubleTeam = dto.DoubleTeam
}

type SchemeCount struct {
	PowerRun      int
	Vertical      int
	WestCoast     int
	IOption       int
	RunAndShoot   int
	AirRaid       int
	Pistol        int
	SpreadOption  int
	WingT         int
	DoubleWing    int
	Flexbone      int
	Wishbone      int
	OldSchool     int
	TwoGap        int
	FourManFront  int
	ThreeManFront int
	Speed         int
	Multiple      int
}

func (sc *SchemeCount) IncrementScheme(scheme string, val int) {
	if scheme == "Power Run" {
		sc.PowerRun += val
	}
	if scheme == "Vertical" {
		sc.Vertical += val
	}
	if scheme == "West Coast" {
		sc.WestCoast += val
	}
	if scheme == "I Option" {
		sc.IOption += val
	}
	if scheme == "Run and Shoot" {
		sc.RunAndShoot += val
	}
	if scheme == "Air Raid" {
		sc.AirRaid += val
	}
	if scheme == "Pistol" {
		sc.Pistol += val
	}
	if scheme == "Spread Option" {
		sc.SpreadOption += val
	}
	if scheme == "Wing-T" {
		sc.WingT += val
	}
	if scheme == "Double Wing" {
		sc.DoubleWing += val
	}
	if scheme == "Flexbone" {
		sc.Flexbone += val
	}
	if scheme == "Wishbone" {
		sc.Wishbone += val
	}
	if scheme == "Old School" {
		sc.OldSchool += val
	}
	if scheme == "2-Gap" {
		sc.TwoGap += val
	}
	if scheme == "4-Man Front Spread Stopper" {
		sc.FourManFront += val
	}
	if scheme == "3-Man Front Spread Stopper" {
		sc.ThreeManFront += val
	}
	if scheme == "Speed" {
		sc.Speed += val
	}
	if scheme == "Multiple" {
		sc.Multiple += val
	}
}

func (sc *SchemeCount) DecrementScheme(scheme string, val int) {
	if scheme == "Power Run" {
		sc.PowerRun -= val
	}
	if scheme == "Vertical" {
		sc.Vertical -= val
	}
	if scheme == "West Coast" {
		sc.WestCoast -= val
	}
	if scheme == "I Option" {
		sc.IOption -= val
	}
	if scheme == "Run and Shoot" {
		sc.RunAndShoot -= val
	}
	if scheme == "Air Raid" {
		sc.AirRaid -= val
	}
	if scheme == "Pistol" {
		sc.Pistol -= val
	}
	if scheme == "Spread Option" {
		sc.SpreadOption -= val
	}
	if scheme == "Wing-T" {
		sc.WingT -= val
	}
	if scheme == "Double Wing" {
		sc.DoubleWing -= val
	}
	if scheme == "Flexbone" {
		sc.Flexbone -= val
	}
	if scheme == "Wishbone" {
		sc.Wishbone -= val
	}
	if scheme == "Old School" {
		sc.OldSchool -= val
	}
	if scheme == "2-Gap" {
		sc.TwoGap -= val
	}
	if scheme == "4-Man Front Spread Stopper" {
		sc.FourManFront -= val
	}
	if scheme == "3-Man Front Spread Stopper" {
		sc.ThreeManFront -= val
	}
	if scheme == "Speed" {
		sc.Speed -= val
	}
	if scheme == "Multiple" {
		sc.Multiple -= val
	}
}
