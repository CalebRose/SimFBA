package structs

type CompensatoryDraftPickDTO struct {
	TeamID      uint `json:"TeamID"`
	Season      uint `json:"Season"`
	SeasonID    uint `json:"SeasonID"`
	DraftRound  uint `json:"DraftRound"`
	DraftNumber uint `json:"DraftNumber"`
}

type ForfeitDraftPickDTO struct {
	DraftPickID uint   `json:"DraftPickID"`
	Reason      string `json:"Reason"`
}

type RunCombineDTO struct {
	PlayerIDs []uint `json:"PlayerIDs"`
}

type ProDayNominationDTO struct {
	TeamID    uint   `json:"TeamID"`
	PlayerIDs []uint `json:"PlayerIDs"`
}
