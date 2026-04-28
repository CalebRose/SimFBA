package structs

import "github.com/jinzhu/gorm"

type Recruit struct {
	gorm.Model
	PlayerID int
	TeamID   int
	BasePlayer
	HighSchool            string
	City                  string
	State                 string
	AffinityOne           string
	AffinityTwo           string
	IsSigned              bool
	IsCustomCroot         bool
	CustomCrootFor        string
	College               string
	OverallRank           float64
	RivalsRank            float64
	ESPNRank              float64
	Rank247               float64
	TopRankModifier       float64
	RecruitingModifier    float64
	RecruitingStatus      string
	RecruitPlayerProfiles []RecruitPlayerProfile `gorm:"foreignKey:RecruitID"`
	// RecruitPoints         []RecruitPointAllocation `gorm:"foreignKey:RecruitID"`
}

func (r *Recruit) UpdatePlayerID(id int) {
	r.PlayerID = id
}

func (r *Recruit) UpdateTeamID(id int) {
	r.TeamID = id
	if id > 0 {
		r.IsSigned = true
	}
}

func (r *Recruit) AssignCollege(abbr string) {
	r.College = abbr
}

func (r *Recruit) ApplyRecruitingStatus(num float64, threshold float64) {
	percentage := num / threshold

	if threshold == 0 || num == 0 || percentage < 0.26 {
		r.RecruitingStatus = "Not Ready"
	} else if percentage < 0.51 {
		r.RecruitingStatus = "Hearing Offers"
	} else if percentage < 0.76 {
		r.RecruitingStatus = "Visiting Schools"
	} else if percentage < 0.96 {
		r.RecruitingStatus = "Finalizing Decisions"
	} else if percentage < 1 {
		r.RecruitingStatus = "Ready to Sign"
	} else {
		r.RecruitingStatus = "Signed"
	}
}

func (r *Recruit) UpdateSigningStatus() {
	r.IsSigned = true
}

func (r *Recruit) Map(createRecruitDTO CreateRecruitDTO, lastPlayerID uint) {
	r.ID = lastPlayerID
	r.FirstName = createRecruitDTO.FirstName
	r.LastName = createRecruitDTO.LastName
	r.Position = createRecruitDTO.Position
	r.Archetype = createRecruitDTO.Archetype
	r.Age = createRecruitDTO.Age
	r.Height = createRecruitDTO.Height
	r.Weight = createRecruitDTO.Weight
	r.Stars = createRecruitDTO.Stars
	r.Overall = createRecruitDTO.Overall
	r.Stamina = createRecruitDTO.Stamina
	r.Injury = createRecruitDTO.Injury
	r.FootballIQ = createRecruitDTO.FootballIQ
	r.WorkEthic = createRecruitDTO.WorkEthic
	r.Speed = createRecruitDTO.Speed
	r.Carrying = createRecruitDTO.Carrying
	r.Agility = createRecruitDTO.Agility
	r.Catching = createRecruitDTO.Catching
	r.RouteRunning = createRecruitDTO.RouteRunning
	r.ZoneCoverage = createRecruitDTO.ZoneCoverage
	r.ManCoverage = createRecruitDTO.ManCoverage
	r.Strength = createRecruitDTO.Strength
	r.Tackle = createRecruitDTO.Tackle
	r.PassBlock = createRecruitDTO.PassBlock
	r.RunBlock = createRecruitDTO.RunBlock
	r.PassRush = createRecruitDTO.PassRush
	r.RunDefense = createRecruitDTO.RunDefense
	r.ThrowPower = createRecruitDTO.ThrowPower
	r.ThrowAccuracy = createRecruitDTO.ThrowAccuracy
	r.KickAccuracy = createRecruitDTO.KickAccuracy
	r.KickPower = createRecruitDTO.KickPower
	r.PuntAccuracy = createRecruitDTO.PuntAccuracy
	r.PuntPower = createRecruitDTO.PuntPower
	r.Progression = createRecruitDTO.Progression
	r.Discipline = createRecruitDTO.Discipline
	r.PotentialGrade = createRecruitDTO.PotentialGrade
	r.HighSchool = createRecruitDTO.HighSchool
	r.City = createRecruitDTO.City
	r.State = createRecruitDTO.State
	r.AffinityOne = createRecruitDTO.AffinityOne
	r.AffinityTwo = createRecruitDTO.AffinityTwo
	r.FreeAgency = createRecruitDTO.FreeAgency
	r.Personality = createRecruitDTO.Personality
	r.RecruitingBias = createRecruitDTO.RecruitingBias
	r.AcademicBias = createRecruitDTO.AcademicBias
	r.IsSigned = false
	r.CustomCrootFor = createRecruitDTO.CreatedFor
	r.IsCustomCroot = true
}

func (r *Recruit) AssignPlayerID(ID int) {
	r.PlayerID = ID
}

func (r *Recruit) AssignID(ID int) {
	r.ID = uint(ID)
}

func (r *Recruit) AssignRankValues(rank247 float64, espnRank float64, rivalsRank float64, modifier float64) {
	r.Rank247 = rank247
	r.ESPNRank = espnRank
	r.RivalsRank = rivalsRank
	r.TopRankModifier = modifier
}

func (r *Recruit) AssignRecruitingModifier(recruitingMod float64) {
	r.RecruitingModifier = recruitingMod
}

func (r *Recruit) ProgressUnsignedRecruit(attr CollegePlayerProgressions) {
	r.Age++
	r.Agility = attr.Agility
	r.Speed = attr.Speed
	r.FootballIQ = attr.FootballIQ
	r.Carrying = attr.Carrying
	r.Catching = attr.Catching
	r.RouteRunning = attr.RouteRunning
	r.PassBlock = attr.PassBlock
	r.RunBlock = attr.RunBlock
	r.PassRush = attr.PassRush
	r.RunDefense = attr.RunDefense
	r.Tackle = attr.Tackle
	r.ManCoverage = attr.ManCoverage
	r.ZoneCoverage = attr.ZoneCoverage
	r.KickAccuracy = attr.KickAccuracy
	r.KickPower = attr.KickPower
	r.PuntAccuracy = attr.PuntAccuracy
	r.PuntPower = attr.PuntPower
	r.ThrowAccuracy = attr.ThrowAccuracy
	r.ThrowPower = attr.ThrowPower
}

func (r *Recruit) GetOverall() {
	var ovr float64 = 0
	switch r.Position {
	case "QB":
		ovr = (0.1 * float64(r.Agility)) + (0.25 * float64(r.ThrowPower)) + (0.25 * float64(r.ThrowAccuracy)) + (0.1 * float64(r.Speed)) + (0.2 * float64(r.FootballIQ)) + (0.1 * float64(r.Strength))
		r.Overall = int(ovr)
	case "RB":
		ovr = (0.2 * float64(r.Agility)) + (0.05 * float64(r.PassBlock)) +
			(0.1 * float64(r.Carrying)) + (0.25 * float64(r.Speed)) +
			(0.15 * float64(r.FootballIQ)) + (0.2 * float64(r.Strength)) +
			(0.05 * float64(r.Catching))
		r.Overall = int(ovr)
	case "FB":
		ovr = (0.1 * float64(r.Agility)) + (0.1 * float64(r.PassBlock)) +
			(0.1 * float64(r.Carrying)) + (0.05 * float64(r.Speed)) +
			(0.15 * float64(r.FootballIQ)) + (0.2 * float64(r.Strength)) +
			(0.05 * float64(r.Catching)) + (0.25 * float64(r.RunBlock))
		r.Overall = int(ovr)
	case "WR":
		ovr = (0.15 * float64(r.FootballIQ)) + (0.2 * float64(r.Speed)) +
			(0.1 * float64(r.Agility)) + (0.05 * float64(r.Carrying)) +
			(0.05 * float64(r.Strength)) + (0.25 * float64(r.Catching)) +
			(0.2 * float64(r.RouteRunning))
		r.Overall = int(ovr)
	case "TE":
		ovr = (0.15 * float64(r.FootballIQ)) + (0.1 * float64(r.Speed)) +
			(0.1 * float64(r.Agility)) + (0.05 * float64(r.Carrying)) +
			(0.05 * float64(r.PassBlock)) + (0.15 * float64(r.RunBlock)) +
			(0.1 * float64(r.Strength)) + (0.20 * float64(r.Catching)) +
			(0.1 * float64(r.RouteRunning))
		r.Overall = int(ovr)
	case "OT", "OG":
		ovr = (0.15 * float64(r.FootballIQ)) + (0.05 * float64(r.Agility)) +
			(0.3 * float64(r.RunBlock)) + (0.2 * float64(r.Strength)) +
			(0.3 * float64(r.PassBlock))
		r.Overall = int(ovr)
	case "C":
		ovr = (0.2 * float64(r.FootballIQ)) + (0.05 * float64(r.Agility)) +
			(0.3 * float64(r.RunBlock)) + (0.15 * float64(r.Strength)) +
			(0.3 * float64(r.PassBlock))
		r.Overall = int(ovr)
	case "DT":
		ovr = (0.15 * float64(r.FootballIQ)) + (0.05 * float64(r.Agility)) +
			(0.25 * float64(r.RunDefense)) + (0.2 * float64(r.Strength)) +
			(0.15 * float64(r.PassRush)) + (0.2 * float64(r.Tackle)) +
			(0.1 * float64(r.Agility))
		r.Overall = int(ovr)
	case "DE":
		ovr = (0.15 * float64(r.FootballIQ)) + (0.1 * float64(r.Speed)) +
			(0.15 * float64(r.RunDefense)) + (0.1 * float64(r.Strength)) +
			(0.2 * float64(r.PassRush)) + (0.2 * float64(r.Tackle)) +
			(0.1 * float64(r.Agility))
		r.Overall = int(ovr)
	case "ILB":
		ovr = (0.2 * float64(r.FootballIQ)) + (0.1 * float64(r.Speed)) +
			(0.15 * float64(r.RunDefense)) + (0.1 * float64(r.Strength)) +
			(0.1 * float64(r.PassRush)) + (0.15 * float64(r.Tackle)) +
			(0.1 * float64(r.ZoneCoverage)) + (0.05 * float64(r.ManCoverage)) +
			(0.05 * float64(r.Agility))
		r.Overall = int(ovr)
	case "OLB":
		ovr = (0.15 * float64(r.FootballIQ)) + (0.1 * float64(r.Speed)) +
			(0.15 * float64(r.RunDefense)) + (0.1 * float64(r.Strength)) +
			(0.15 * float64(r.PassRush)) + (0.15 * float64(r.Tackle)) +
			(0.1 * float64(r.ZoneCoverage)) + (0.05 * float64(r.ManCoverage)) +
			(0.05 * float64(r.Agility))
		r.Overall = int(ovr)
	case "CB":
		ovr = (0.15 * float64(r.FootballIQ)) + (0.25 * float64(r.Speed)) +
			(0.05 * float64(r.Tackle)) + (0.05 * float64(r.Strength)) +
			(0.15 * float64(r.Agility)) + (0.15 * float64(r.ZoneCoverage)) +
			(0.15 * float64(r.ManCoverage)) + (0.05 * float64(r.Catching))
		r.Overall = int(ovr)
	case "FS":
		ovr = (0.2 * float64(r.FootballIQ)) + (0.2 * float64(r.Speed)) +
			(0.05 * float64(r.RunDefense)) + (0.05 * float64(r.Strength)) +
			(0.05 * float64(r.Catching)) + (0.05 * float64(r.Tackle)) +
			(0.15 * float64(r.ZoneCoverage)) + (0.15 * float64(r.ManCoverage)) +
			(0.1 * float64(r.Agility))
		r.Overall = int(ovr)
	case "SS":
		ovr = (0.15 * float64(r.FootballIQ)) + (0.2 * float64(r.Speed)) +
			(0.05 * float64(r.RunDefense)) + (0.05 * float64(r.Strength)) +
			(0.05 * float64(r.Catching)) + (0.1 * float64(r.Tackle)) +
			(0.15 * float64(r.ZoneCoverage)) + (0.15 * float64(r.ManCoverage)) +
			(0.1 * float64(r.Agility))
		r.Overall = int(ovr)
	case "K":
		ovr = (0.2 * float64(r.FootballIQ)) + (0.45 * float64(r.KickPower)) +
			(0.45 * float64(r.KickAccuracy))
		r.Overall = int(ovr)
	case "P":
		ovr = (0.2 * float64(r.FootballIQ)) + (0.45 * float64(r.PuntPower)) +
			(0.45 * float64(r.PuntAccuracy))
		r.Overall = int(ovr)
	case "ATH":
		ovr = (float64(r.FootballIQ) + float64(r.Speed) + float64(r.Agility) +
			float64(r.Carrying) + float64(r.Catching) + float64(r.RouteRunning) +
			float64(r.RunBlock) + float64(r.PassBlock) + float64(r.PassRush) +
			float64(r.RunDefense) + float64(r.Tackle) + float64(r.Strength) +
			float64(r.ZoneCoverage) + float64(r.ManCoverage) + float64(r.ThrowAccuracy) +
			float64(r.ThrowPower) + float64(r.PuntAccuracy) + float64(r.PuntPower) +
			float64(r.KickAccuracy) + float64(r.KickPower)) / 20
		r.Overall = int(ovr)
	}
}

func (r *Recruit) AssignJUCOSchool(school string) {
	r.HighSchool = school
}

func (r *Recruit) AssignWalkon(abbr string, teamID int, id uint) {
	r.College = abbr
	r.TeamID = teamID
	r.PlayerID = int(id)
	r.ID = id
}

func (r *Recruit) AssignNewRecruitingBiases(bias, af1, af2 string) {
	r.RecruitingBias = bias
	r.AffinityOne = af1
	r.AffinityTwo = af2
}

func (r *Recruit) AssignRelativeData(rID, rType, teamID uint, team, notes string) {
	r.RelativeID = rID
	r.RelativeType = rType
	r.Notes = notes
	if teamID > 0 {
		r.UpdateTeamID(int(teamID))
		r.AssignCollege(team)
	}
}

func (r *Recruit) AssignTwinData(lastName, city, state, highschool string) {
	r.LastName = lastName
	r.City = city
	r.State = state
	r.HighSchool = highschool
}
