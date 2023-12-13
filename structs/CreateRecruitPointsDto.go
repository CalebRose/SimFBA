package structs

type CreateRecruitProfileDto struct {
	PlayerID            int
	SeasonID            int
	RecruitID           int
	ProfileID           int
	Team                string
	RES                 float64
	AffinityOneEligible bool
	AffinityTwoEligible bool
	PlayerRecruit       Recruit
	Recruiter           string
}
