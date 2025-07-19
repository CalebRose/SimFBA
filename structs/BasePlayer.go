package structs

type BasePlayer struct {
	FirstName       string
	LastName        string
	Position        string
	Archetype       string
	PreviousTeamID  uint
	PreviousTeam    string
	Height          int
	Weight          int
	Age             int
	Stars           int
	Overall         int
	Stamina         int
	Injury          int
	FootballIQ      int
	Speed           int
	Carrying        int
	Agility         int
	Catching        int
	RouteRunning    int
	ZoneCoverage    int
	ManCoverage     int
	Strength        int
	Tackle          int
	PassBlock       int
	RunBlock        int
	PassRush        int
	RunDefense      int
	ThrowPower      int
	ThrowAccuracy   int
	KickAccuracy    int
	KickPower       int
	PuntAccuracy    int
	PuntPower       int
	Progression     int
	Discipline      int
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
}

func (cp *BasePlayer) GetOverall() {
	var ovr float64 = 0
	if cp.Position == "QB" {
		ovr = (0.1 * float64(cp.Agility)) + (0.25 * float64(cp.ThrowPower)) + (0.25 * float64(cp.ThrowAccuracy)) + (0.1 * float64(cp.Speed)) + (0.2 * float64(cp.FootballIQ)) + (0.1 * float64(cp.Strength))
		cp.Overall = int(ovr)
	} else if cp.Position == "RB" {
		ovr = (0.2 * float64(cp.Agility)) + (0.05 * float64(cp.PassBlock)) +
			(0.1 * float64(cp.Carrying)) + (0.25 * float64(cp.Speed)) +
			(0.15 * float64(cp.FootballIQ)) + (0.2 * float64(cp.Strength)) +
			(0.05 * float64(cp.Catching))
		cp.Overall = int(ovr)
	} else if cp.Position == "FB" {
		ovr = (0.1 * float64(cp.Agility)) + (0.1 * float64(cp.PassBlock)) +
			(0.1 * float64(cp.Carrying)) + (0.05 * float64(cp.Speed)) +
			(0.15 * float64(cp.FootballIQ)) + (0.2 * float64(cp.Strength)) +
			(0.05 * float64(cp.Catching)) + (0.25 * float64(cp.RunBlock))
		cp.Overall = int(ovr)
	} else if cp.Position == "WR" {
		ovr = (0.15 * float64(cp.FootballIQ)) + (0.2 * float64(cp.Speed)) +
			(0.1 * float64(cp.Agility)) + (0.05 * float64(cp.Carrying)) +
			(0.05 * float64(cp.Strength)) + (0.25 * float64(cp.Catching)) +
			(0.2 * float64(cp.RouteRunning))
		cp.Overall = int(ovr)
	} else if cp.Position == "TE" {
		ovr = (0.15 * float64(cp.FootballIQ)) + (0.1 * float64(cp.Speed)) +
			(0.1 * float64(cp.Agility)) + (0.05 * float64(cp.Carrying)) +
			(0.05 * float64(cp.PassBlock)) + (0.15 * float64(cp.RunBlock)) +
			(0.1 * float64(cp.Strength)) + (0.20 * float64(cp.Catching)) +
			(0.1 * float64(cp.RouteRunning))
		cp.Overall = int(ovr)
	} else if cp.Position == "OT" || cp.Position == "OG" {
		ovr = (0.15 * float64(cp.FootballIQ)) + (0.05 * float64(cp.Agility)) +
			(0.3 * float64(cp.RunBlock)) + (0.2 * float64(cp.Strength)) +
			(0.3 * float64(cp.PassBlock))
		cp.Overall = int(ovr)
	} else if cp.Position == "C" {
		ovr = (0.2 * float64(cp.FootballIQ)) + (0.05 * float64(cp.Agility)) +
			(0.3 * float64(cp.RunBlock)) + (0.15 * float64(cp.Strength)) +
			(0.3 * float64(cp.PassBlock))
		cp.Overall = int(ovr)
	} else if cp.Position == "DT" {
		ovr = (0.15 * float64(cp.FootballIQ)) + (0.05 * float64(cp.Agility)) +
			(0.25 * float64(cp.RunDefense)) + (0.2 * float64(cp.Strength)) +
			(0.15 * float64(cp.PassRush)) + (0.2 * float64(cp.Tackle))
		cp.Overall = int(ovr)
	} else if cp.Position == "DE" {
		ovr = (0.15 * float64(cp.FootballIQ)) + (0.1 * float64(cp.Speed)) +
			(0.15 * float64(cp.RunDefense)) + (0.1 * float64(cp.Strength)) +
			(0.2 * float64(cp.PassRush)) + (0.2 * float64(cp.Tackle)) +
			(0.1 * float64(cp.Agility))
		cp.Overall = int(ovr)
	} else if cp.Position == "ILB" {
		ovr = (0.2 * float64(cp.FootballIQ)) + (0.1 * float64(cp.Speed)) +
			(0.15 * float64(cp.RunDefense)) + (0.1 * float64(cp.Strength)) +
			(0.1 * float64(cp.PassRush)) + (0.15 * float64(cp.Tackle)) +
			(0.1 * float64(cp.ZoneCoverage)) + (0.05 * float64(cp.ManCoverage)) +
			(0.05 * float64(cp.Agility))
		cp.Overall = int(ovr)
	} else if cp.Position == "OLB" {
		ovr = (0.15 * float64(cp.FootballIQ)) + (0.1 * float64(cp.Speed)) +
			(0.15 * float64(cp.RunDefense)) + (0.1 * float64(cp.Strength)) +
			(0.15 * float64(cp.PassRush)) + (0.15 * float64(cp.Tackle)) +
			(0.1 * float64(cp.ZoneCoverage)) + (0.05 * float64(cp.ManCoverage)) +
			(0.05 * float64(cp.Agility))
		cp.Overall = int(ovr)
	} else if cp.Position == "CB" {
		ovr = (0.15 * float64(cp.FootballIQ)) + (0.25 * float64(cp.Speed)) +
			(0.05 * float64(cp.Tackle)) + (0.05 * float64(cp.Strength)) +
			(0.15 * float64(cp.Agility)) + (0.15 * float64(cp.ZoneCoverage)) +
			(0.15 * float64(cp.ManCoverage)) + (0.05 * float64(cp.Catching))
		cp.Overall = int(ovr)
	} else if cp.Position == "FS" {
		ovr = (0.2 * float64(cp.FootballIQ)) + (0.2 * float64(cp.Speed)) +
			(0.05 * float64(cp.RunDefense)) + (0.05 * float64(cp.Strength)) +
			(0.05 * float64(cp.Catching)) + (0.05 * float64(cp.Tackle)) +
			(0.15 * float64(cp.ZoneCoverage)) + (0.15 * float64(cp.ManCoverage)) +
			(0.1 * float64(cp.Agility))
		cp.Overall = int(ovr)
	} else if cp.Position == "SS" {
		ovr = (0.15 * float64(cp.FootballIQ)) + (0.2 * float64(cp.Speed)) +
			(0.05 * float64(cp.RunDefense)) + (0.05 * float64(cp.Strength)) +
			(0.05 * float64(cp.Catching)) + (0.1 * float64(cp.Tackle)) +
			(0.15 * float64(cp.ZoneCoverage)) + (0.15 * float64(cp.ManCoverage)) +
			(0.1 * float64(cp.Agility))
		cp.Overall = int(ovr)
	} else if cp.Position == "K" {
		ovr = (0.2 * float64(cp.FootballIQ)) + (0.45 * float64(cp.KickPower)) +
			(0.45 * float64(cp.KickAccuracy))
		cp.Overall = int(ovr)
	} else if cp.Position == "P" {
		ovr = (0.2 * float64(cp.FootballIQ)) + (0.45 * float64(cp.PuntPower)) +
			(0.45 * float64(cp.PuntAccuracy))
		cp.Overall = int(ovr)
	} else if cp.Position == "ATH" {
		if cp.Archetype == "Field General" {
			ovr = (.20 * float64(cp.FootballIQ)) + (.1 * float64(cp.ZoneCoverage)) + (.1 * float64(cp.ManCoverage)) + (.1 * float64(cp.RunDefense)) + (.1 * float64(cp.Speed)) + (.1 * float64(cp.Strength)) + (.1 * float64(cp.Tackle)) + (.1 * float64(cp.ThrowPower)) + (.1 * float64(cp.ThrowAccuracy))
		} else if cp.Archetype == "Triple-Threat" {
			ovr = (.10 * float64(cp.FootballIQ)) + (.2 * float64(cp.Agility)) + (.1 * float64(cp.Carrying)) + (.1 * float64(cp.Catching)) + (.2 * float64(cp.Speed)) + (.1 * float64(cp.RouteRunning)) + (.1 * float64(cp.ThrowPower)) + (.1 * float64(cp.ThrowAccuracy))
		} else if cp.Archetype == "Wingback" {
			ovr = (.1 * float64(cp.FootballIQ)) + (.2 * float64(cp.Agility)) + (.1 * float64(cp.Carrying)) + (.2 * float64(cp.Catching)) + (.2 * float64(cp.Speed)) + (.1 * float64(cp.RouteRunning)) + (.1 * float64(cp.RunBlock))
		} else if cp.Archetype == "Slotback" {
			ovr = (.1 * float64(cp.FootballIQ)) + (.1 * float64(cp.Agility)) + (.1 * float64(cp.Carrying)) + (.1 * float64(cp.Catching)) + (.2 * float64(cp.Speed)) + (.1 * float64(cp.RouteRunning)) + (.1 * float64(cp.RunBlock)) + (.1 * float64(cp.PassBlock)) + (.1 * float64(cp.Strength))
		} else if cp.Archetype == "Lineman" {
			ovr = (.1 * float64(cp.FootballIQ)) + (.1 * float64(cp.Agility)) + (.1 * float64(cp.RunBlock)) + (.1 * float64(cp.PassBlock)) + (.3 * float64(cp.Strength)) + (.1 * float64(cp.PassRush)) + (.1 * float64(cp.RunDefense)) + (.1 * float64(cp.Tackle))
		} else if cp.Archetype == "Strongside" {
			ovr = (.1 * float64(cp.FootballIQ)) + (.1 * float64(cp.Agility)) + (.1 * float64(cp.ZoneCoverage)) + (.1 * float64(cp.ManCoverage)) + (.2 * float64(cp.Strength)) + (.1 * float64(cp.PassRush)) + (.1 * float64(cp.RunDefense)) + (.1 * float64(cp.Tackle)) + (.1 * float64(cp.Speed))
		} else if cp.Archetype == "Weakside" {
			ovr = (.1 * float64(cp.FootballIQ)) + (.1 * float64(cp.Agility)) + (.1 * float64(cp.ZoneCoverage)) + (.1 * float64(cp.ManCoverage)) + (.1 * float64(cp.Strength)) + (.1 * float64(cp.PassRush)) + (.1 * float64(cp.RunDefense)) + (.1 * float64(cp.Tackle)) + (.2 * float64(cp.Speed))
		} else if cp.Archetype == "Bandit" {
			ovr = (.1 * float64(cp.FootballIQ)) + (.1 * float64(cp.Agility)) + (.1 * float64(cp.ZoneCoverage)) + (.1 * float64(cp.ManCoverage)) + (.1 * float64(cp.Strength)) + (.1 * float64(cp.PassRush)) + (.1 * float64(cp.RunDefense)) + (.1 * float64(cp.Tackle)) + (.2 * float64(cp.Speed))
		} else if cp.Archetype == "Return Specialist" {
			ovr = (.20 * float64(cp.FootballIQ)) + (.10 * float64(cp.Strength)) + (.20 * float64(cp.Speed)) + (.20 * float64(cp.Agility)) + (.20 * float64(cp.Catching)) + (.1 * float64(cp.Tackle))
		} else if cp.Archetype == "Soccer Player" {
			ovr = (.10 * float64(cp.FootballIQ)) + (.10 * float64(cp.Agility)) + (.2 * float64(cp.KickPower)) + (.2 * float64(cp.KickAccuracy)) + (.2 * float64(cp.PuntPower)) + (.2 * float64(cp.PuntAccuracy))
		}
		cp.Overall = int(ovr)
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
	cp.Progression = prog
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
	bp.FootballIQ = fIQ
	bp.Speed = spe
	bp.Agility = agi
	bp.Carrying = car
	bp.Catching = cat
	bp.RouteRunning = rr
	bp.ZoneCoverage = zc
	bp.ManCoverage = mc
	bp.Strength = str
	bp.Tackle = tack
	bp.PassBlock = pb
	bp.RunBlock = rb
	bp.PassRush = pr
	bp.RunDefense = rd
	bp.ThrowPower = tp
	bp.ThrowAccuracy = ta
	bp.KickAccuracy = ka
	bp.KickPower = kp
	bp.PuntAccuracy = pa
	bp.PuntPower = pp
}

func (bp *BasePlayer) RevertAge() {
	bp.Age--
	if bp.Age < 18 {
		bp.Age = 18
	}
}
