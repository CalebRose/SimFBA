package structs

type NFLPlayoffSeedDTO struct {
	Seed       uint   `json:"Seed"`
	TeamID     uint   `json:"TeamID"`
	TeamAbbr   string `json:"TeamAbbr"`
	Conference string `json:"Conference"` // "AFC" or "NFC"
}

type NFLPlayoffScheduleRequest struct {
	Seeds []NFLPlayoffSeedDTO `json:"Seeds"`
}
