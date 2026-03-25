package structs

import "github.com/jinzhu/gorm"

type TeamRequest struct {
	gorm.Model
	TeamID     int
	Username   string
	IsApproved bool
}

func (r *TeamRequest) ApproveTeamRequest() {
	r.IsApproved = true
}

func (r *TeamRequest) RejectTeamRequest() {
	r.IsApproved = false
}

type TeamRequestsResponse struct {
	CollegeRequests []TeamRequest
	ProRequests     []NFLRequest
	AcceptedTrades  []NFLTradeProposal
	DraftPicks      []NFLDraftPick
}
