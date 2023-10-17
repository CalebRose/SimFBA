package structs

type BasePlayer struct {
	FirstName       string
	LastName        string
	Position        string
	Archetype       string
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
			(0.15 * float64(cp.PassRush)) + (0.2 * float64(cp.Tackle)) +
			(0.1 * float64(cp.Agility))
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
		ovr = (float64(cp.FootballIQ) + float64(cp.Speed) + float64(cp.Agility) +
			float64(cp.Carrying) + float64(cp.Catching) + float64(cp.RouteRunning) +
			float64(cp.RunBlock) + float64(cp.PassBlock) + float64(cp.PassRush) +
			float64(cp.RunDefense) + float64(cp.Tackle) + float64(cp.Strength) +
			float64(cp.ZoneCoverage) + float64(cp.ManCoverage) + float64(cp.ThrowAccuracy) +
			float64(cp.ThrowPower) + float64(cp.PuntAccuracy) + float64(cp.PuntPower) +
			float64(cp.KickAccuracy) + float64(cp.KickPower)) / 20
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

func (cp *BasePlayer) ToggleInjuryReserve() {
	cp.InjuryReserve = !cp.InjuryReserve
}

func (cp *BasePlayer) MapProgression(prog int, letter string) {
	cp.Progression = prog
	cp.PotentialGrade = letter
}
