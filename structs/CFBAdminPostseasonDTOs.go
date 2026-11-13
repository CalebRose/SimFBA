package structs

type CFBConferenceChampionshipDTO struct {
	GameID       uint   `json:"GameID"`
	ConferenceID uint   `json:"ConferenceID"`
	Conference   string `json:"Conference"`
	HomeTeamID   uint   `json:"HomeTeamID"`
	AwayTeamID   uint   `json:"AwayTeamID"`
	TimeSlot     string `json:"TimeSlot"`
	StadiumID    uint   `json:"StadiumID"`
	Stadium      string `json:"Stadium"`
	City         string `json:"City"`
	State        string `json:"State"`
	IsDomed      bool   `json:"IsDomed"`
	IsNeutral    bool   `json:"IsNeutral"`
	GameTitle    string `json:"GameTitle"`
	Week         int    `json:"Week"`
}

type CFBBowlFixtureDTO struct {
	GameID     uint   `json:"GameID"`
	Week       int    `json:"Week"`
	TimeSlot   string `json:"TimeSlot"`
	GameTitle  string `json:"GameTitle"`
	HomeTeam   string `json:"HomeTeam"`
	AwayTeam   string `json:"AwayTeam"`
	HomeTeamID uint   `json:"HomeTeamID"`
	AwayTeamID uint   `json:"AwayTeamID"`
	StadiumID  uint   `json:"StadiumID"`
	Stadium    string `json:"Stadium"`
	City       string `json:"City"`
	State      string `json:"State"`
	IsDomed    bool   `json:"IsDomed"`
	IsNeutral  bool   `json:"IsNeutral"`
}

type CFBBowlUploadRequest struct {
	Fixtures          []CFBBowlFixtureDTO `json:"Fixtures"`
	OverrideWeekCheck bool                `json:"OverrideWeekCheck"`
}

type CFBPlayoffSeedDTO struct {
	Seed     uint   `json:"Seed"`
	TeamID   uint   `json:"TeamID"`
	TeamAbbr string `json:"TeamAbbr"`
}

type CFBPlayoffScheduleRequest struct {
	Seeds []CFBPlayoffSeedDTO `json:"Seeds"`
}
