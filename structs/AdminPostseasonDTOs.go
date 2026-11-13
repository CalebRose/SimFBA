package structs

type CompensatoryDraftPickDTO struct {
	TeamID      uint `json:"TeamID"`
	Season      uint `json:"Season"`
	SeasonID    uint `json:"SeasonID"`
	DraftRound  uint `json:"DraftRound"`
	DraftNumber uint `json:"DraftNumber"`
}
