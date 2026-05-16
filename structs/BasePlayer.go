package structs

type PlayerPreferences struct {
	ProgramPref        uint8
	ProfDevPref        uint8
	TraditionsPref     uint8
	FacilitiesPref     uint8
	AtmospherePref     uint8
	AcademicsPref      uint8
	ConferencePref     uint8
	CoachPref          uint8
	SeasonMomentumPref uint8
	CampusLifePref     uint8
	ReligionPref       uint8
	ServiceAcademyPref uint8
	SmallTownPref      uint8
	BigCityPref        uint8
	MediaSpotlightPref uint8
}

type BasePlayer struct {
	FirstName       string
	LastName        string
	Position        string
	Archetype       string
	PreviousTeamID  uint
	PreviousTeam    string
	Height          int8
	Weight          int8
	Age             int8
	Stars           int8
	Overall         int8
	Stamina         int8
	Injury          int8
	FootballIQ      int8
	Speed           int8
	Carrying        int8
	Agility         int8
	Catching        int8
	RouteRunning    int8
	ZoneCoverage    int8
	ManCoverage     int8
	Strength        int8
	Tackle          int8
	PassBlock       int8
	RunBlock        int8
	PassRush        int8
	RunDefense      int8
	ThrowPower      int8
	ThrowAccuracy   int8
	KickAccuracy    int8
	KickPower       int8
	PuntAccuracy    int8
	PuntPower       int8
	Progression     int8
	Discipline      int8
	PotentialGrade  string
	FreeAgency      string
	Personality     string
	RecruitingBias  string
	WorkEthic       string
	AcademicBias    string
	IsInjured       bool
	InjuryName      string
	InjuryType      string
	WeeksOfRecovery uint
	InjuryReserve   bool
	PrimeAge        uint
	Clutch          int // -1 is choker, 0 is normal, 1 is clutch, 2 is very clutch
	Shotgun         int // -1 is Under Center, 0 Balanced, 1 Shotgun
	PositionTwo     string
	ArchetypeTwo    string
	RelativeID      uint
	RelativeType    uint
	Notes           string
	PlayerPreferences
}

func (cp *BasePlayer) GetOverall() {
	var ovr float64 = 0
	switch cp.Position {
	case "QB":
		ovr = (0.1 * float64(cp.Agility)) + (0.25 * float64(cp.ThrowPower)) + (0.25 * float64(cp.ThrowAccuracy)) + (0.1 * float64(cp.Speed)) + (0.2 * float64(cp.FootballIQ)) + (0.1 * float64(cp.Strength))
		cp.Overall = int8(ovr)
	case "RB":
		ovr = (0.2 * float64(cp.Agility)) + (0.05 * float64(cp.PassBlock)) +
			(0.1 * float64(cp.Carrying)) + (0.25 * float64(cp.Speed)) +
			(0.15 * float64(cp.FootballIQ)) + (0.2 * float64(cp.Strength)) +
			(0.05 * float64(cp.Catching))
		cp.Overall = int8(ovr)
	case "FB":
		ovr = (0.1 * float64(cp.Agility)) + (0.1 * float64(cp.PassBlock)) +
			(0.1 * float64(cp.Carrying)) + (0.05 * float64(cp.Speed)) +
			(0.15 * float64(cp.FootballIQ)) + (0.2 * float64(cp.Strength)) +
			(0.05 * float64(cp.Catching)) + (0.25 * float64(cp.RunBlock))
		cp.Overall = int8(ovr)
	case "WR":
		ovr = (0.15 * float64(cp.FootballIQ)) + (0.2 * float64(cp.Speed)) +
			(0.1 * float64(cp.Agility)) + (0.05 * float64(cp.Carrying)) +
			(0.05 * float64(cp.Strength)) + (0.25 * float64(cp.Catching)) +
			(0.2 * float64(cp.RouteRunning))
		cp.Overall = int8(ovr)
	case "TE":
		ovr = (0.15 * float64(cp.FootballIQ)) + (0.1 * float64(cp.Speed)) +
			(0.1 * float64(cp.Agility)) + (0.05 * float64(cp.Carrying)) +
			(0.05 * float64(cp.PassBlock)) + (0.15 * float64(cp.RunBlock)) +
			(0.1 * float64(cp.Strength)) + (0.20 * float64(cp.Catching)) +
			(0.1 * float64(cp.RouteRunning))
		cp.Overall = int8(ovr)
	case "OT", "OG":
		ovr = (0.15 * float64(cp.FootballIQ)) + (0.05 * float64(cp.Agility)) +
			(0.3 * float64(cp.RunBlock)) + (0.2 * float64(cp.Strength)) +
			(0.3 * float64(cp.PassBlock))
		cp.Overall = int8(ovr)
	case "C":
		ovr = (0.2 * float64(cp.FootballIQ)) + (0.05 * float64(cp.Agility)) +
			(0.3 * float64(cp.RunBlock)) + (0.15 * float64(cp.Strength)) +
			(0.3 * float64(cp.PassBlock))
		cp.Overall = int8(ovr)
	case "DT":
		ovr = (0.15 * float64(cp.FootballIQ)) + (0.05 * float64(cp.Agility)) +
			(0.25 * float64(cp.RunDefense)) + (0.2 * float64(cp.Strength)) +
			(0.15 * float64(cp.PassRush)) + (0.2 * float64(cp.Tackle))
		cp.Overall = int8(ovr)
	case "DE":
		ovr = (0.15 * float64(cp.FootballIQ)) + (0.1 * float64(cp.Speed)) +
			(0.15 * float64(cp.RunDefense)) + (0.1 * float64(cp.Strength)) +
			(0.2 * float64(cp.PassRush)) + (0.2 * float64(cp.Tackle)) +
			(0.1 * float64(cp.Agility))
		cp.Overall = int8(ovr)
	case "ILB":
		ovr = (0.2 * float64(cp.FootballIQ)) + (0.1 * float64(cp.Speed)) +
			(0.15 * float64(cp.RunDefense)) + (0.1 * float64(cp.Strength)) +
			(0.1 * float64(cp.PassRush)) + (0.15 * float64(cp.Tackle)) +
			(0.1 * float64(cp.ZoneCoverage)) + (0.05 * float64(cp.ManCoverage)) +
			(0.05 * float64(cp.Agility))
		cp.Overall = int8(ovr)
	case "OLB":
		ovr = (0.15 * float64(cp.FootballIQ)) + (0.1 * float64(cp.Speed)) +
			(0.15 * float64(cp.RunDefense)) + (0.1 * float64(cp.Strength)) +
			(0.15 * float64(cp.PassRush)) + (0.15 * float64(cp.Tackle)) +
			(0.1 * float64(cp.ZoneCoverage)) + (0.05 * float64(cp.ManCoverage)) +
			(0.05 * float64(cp.Agility))
		cp.Overall = int8(ovr)
	case "CB":
		ovr = (0.15 * float64(cp.FootballIQ)) + (0.25 * float64(cp.Speed)) +
			(0.05 * float64(cp.Tackle)) + (0.05 * float64(cp.Strength)) +
			(0.15 * float64(cp.Agility)) + (0.15 * float64(cp.ZoneCoverage)) +
			(0.15 * float64(cp.ManCoverage)) + (0.05 * float64(cp.Catching))
		cp.Overall = int8(ovr)
	case "FS":
		ovr = (0.2 * float64(cp.FootballIQ)) + (0.2 * float64(cp.Speed)) +
			(0.05 * float64(cp.RunDefense)) + (0.05 * float64(cp.Strength)) +
			(0.05 * float64(cp.Catching)) + (0.05 * float64(cp.Tackle)) +
			(0.15 * float64(cp.ZoneCoverage)) + (0.15 * float64(cp.ManCoverage)) +
			(0.1 * float64(cp.Agility))
		cp.Overall = int8(ovr)
	case "SS":
		ovr = (0.15 * float64(cp.FootballIQ)) + (0.2 * float64(cp.Speed)) +
			(0.05 * float64(cp.RunDefense)) + (0.05 * float64(cp.Strength)) +
			(0.05 * float64(cp.Catching)) + (0.1 * float64(cp.Tackle)) +
			(0.15 * float64(cp.ZoneCoverage)) + (0.15 * float64(cp.ManCoverage)) +
			(0.1 * float64(cp.Agility))
		cp.Overall = int8(ovr)
	case "K":
		ovr = (0.2 * float64(cp.FootballIQ)) + (0.45 * float64(cp.KickPower)) +
			(0.45 * float64(cp.KickAccuracy))
		cp.Overall = int8(ovr)
	case "P":
		ovr = (0.2 * float64(cp.FootballIQ)) + (0.45 * float64(cp.PuntPower)) +
			(0.45 * float64(cp.PuntAccuracy))
		cp.Overall = int8(ovr)
	case "ATH":
		switch cp.Archetype {
		case "Field General":
			ovr = (.20 * float64(cp.FootballIQ)) + (.1 * float64(cp.ZoneCoverage)) + (.1 * float64(cp.ManCoverage)) + (.1 * float64(cp.RunDefense)) + (.1 * float64(cp.Speed)) + (.1 * float64(cp.Strength)) + (.1 * float64(cp.Tackle)) + (.1 * float64(cp.ThrowPower)) + (.1 * float64(cp.ThrowAccuracy))
		case "Triple-Threat":
			ovr = (.10 * float64(cp.FootballIQ)) + (.2 * float64(cp.Agility)) + (.1 * float64(cp.Carrying)) + (.1 * float64(cp.Catching)) + (.2 * float64(cp.Speed)) + (.1 * float64(cp.RouteRunning)) + (.1 * float64(cp.ThrowPower)) + (.1 * float64(cp.ThrowAccuracy))
		case "Wingback":
			ovr = (.1 * float64(cp.FootballIQ)) + (.2 * float64(cp.Agility)) + (.1 * float64(cp.Carrying)) + (.2 * float64(cp.Catching)) + (.2 * float64(cp.Speed)) + (.1 * float64(cp.RouteRunning)) + (.1 * float64(cp.RunBlock))
		case "Slotback":
			ovr = (.1 * float64(cp.FootballIQ)) + (.1 * float64(cp.Agility)) + (.1 * float64(cp.Carrying)) + (.1 * float64(cp.Catching)) + (.2 * float64(cp.Speed)) + (.1 * float64(cp.RouteRunning)) + (.1 * float64(cp.RunBlock)) + (.1 * float64(cp.PassBlock)) + (.1 * float64(cp.Strength))
		case "Lineman":
			ovr = (.1 * float64(cp.FootballIQ)) + (.1 * float64(cp.Agility)) + (.1 * float64(cp.RunBlock)) + (.1 * float64(cp.PassBlock)) + (.3 * float64(cp.Strength)) + (.1 * float64(cp.PassRush)) + (.1 * float64(cp.RunDefense)) + (.1 * float64(cp.Tackle))
		case "Strongside":
			ovr = (.1 * float64(cp.FootballIQ)) + (.1 * float64(cp.Agility)) + (.1 * float64(cp.ZoneCoverage)) + (.1 * float64(cp.ManCoverage)) + (.2 * float64(cp.Strength)) + (.1 * float64(cp.PassRush)) + (.1 * float64(cp.RunDefense)) + (.1 * float64(cp.Tackle)) + (.1 * float64(cp.Speed))
		case "Weakside":
			ovr = (.1 * float64(cp.FootballIQ)) + (.1 * float64(cp.Agility)) + (.1 * float64(cp.ZoneCoverage)) + (.1 * float64(cp.ManCoverage)) + (.1 * float64(cp.Strength)) + (.1 * float64(cp.PassRush)) + (.1 * float64(cp.RunDefense)) + (.1 * float64(cp.Tackle)) + (.2 * float64(cp.Speed))
		case "Bandit":
			ovr = (.1 * float64(cp.FootballIQ)) + (.1 * float64(cp.Agility)) + (.1 * float64(cp.ZoneCoverage)) + (.1 * float64(cp.ManCoverage)) + (.1 * float64(cp.Strength)) + (.1 * float64(cp.PassRush)) + (.1 * float64(cp.RunDefense)) + (.1 * float64(cp.Tackle)) + (.2 * float64(cp.Speed))
		case "Return Specialist":
			ovr = (.20 * float64(cp.FootballIQ)) + (.10 * float64(cp.Strength)) + (.20 * float64(cp.Speed)) + (.20 * float64(cp.Agility)) + (.20 * float64(cp.Catching)) + (.1 * float64(cp.Tackle))
		case "Soccer Player":
			ovr = (.10 * float64(cp.FootballIQ)) + (.10 * float64(cp.Agility)) + (.2 * float64(cp.KickPower)) + (.2 * float64(cp.KickAccuracy)) + (.2 * float64(cp.PuntPower)) + (.2 * float64(cp.PuntAccuracy))
		}
		cp.Overall = int8(ovr)
	}
}

func (cp *BasePlayer) SetIsInjured(isInjured bool, injuryType string, weeksOfRecovery uint) {
	cp.IsInjured = isInjured
	cp.InjuryType = injuryType
	cp.WeeksOfRecovery = weeksOfRecovery
}

func (cp *BasePlayer) ResetInjuryStatus() {
	cp.InjuryName = ""
	cp.InjuryType = ""
	cp.IsInjured = false
}

func (cp *BasePlayer) RecoveryCheck() {
	// Resolves Data Type issues
	var roof uint = 100000000
	cp.WeeksOfRecovery--
	if cp.WeeksOfRecovery == 0 || cp.WeeksOfRecovery > roof {
		cp.ResetInjuryStatus()
	}

}

func (cp *BasePlayer) AssignNewAttributes(shotgun, clutch int) {
	cp.Shotgun = shotgun
	cp.Clutch = clutch
}

func (cp *BasePlayer) ToggleInjuryReserve() {
	cp.InjuryReserve = !cp.InjuryReserve
}

func (cp *BasePlayer) MapProgression(prog int, letter string) {
	cp.Progression = int8(prog)
	cp.PotentialGrade = letter
}

func (cp *BasePlayer) SetRecruitingBias(bias string) {
	cp.RecruitingBias = bias
}

func (cp *BasePlayer) AssignPrimeAge(age uint) {
	cp.PrimeAge = age
}

func (bp *BasePlayer) DesignateNewPosition(pos, arch string) {
	bp.PositionTwo = pos
	bp.ArchetypeTwo = arch
}

func (bp *BasePlayer) ApplyFixedATHAttributes(fIQ, spe, agi, car, cat, rr, zc, mc, str, tack, pb, rb, pr, rd, tp, ta, ka, kp, pa, pp int) {
	bp.FootballIQ = int8(fIQ)
	bp.Speed = int8(spe)
	bp.Agility = int8(agi)
	bp.Carrying = int8(car)
	bp.Catching = int8(cat)
	bp.RouteRunning = int8(rr)
	bp.ZoneCoverage = int8(zc)
	bp.ManCoverage = int8(mc)
	bp.Strength = int8(str)
	bp.Tackle = int8(tack)
	bp.PassBlock = int8(pb)
	bp.RunBlock = int8(rb)
	bp.PassRush = int8(pr)
	bp.RunDefense = int8(rd)
	bp.ThrowPower = int8(tp)
	bp.ThrowAccuracy = int8(ta)
	bp.KickAccuracy = int8(ka)
	bp.KickPower = int8(kp)
	bp.PuntAccuracy = int8(pa)
	bp.PuntPower = int8(pp)
}

func (bp *BasePlayer) RevertAge() {
	bp.Age--
	if bp.Age < 18 {
		bp.Age = 18
	}
}

func (r *BasePlayer) AssignPreferences(pref PlayerPreferences) {
	r.PlayerPreferences = pref
}
