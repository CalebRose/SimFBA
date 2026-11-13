package managers

import (
	"strings"
	"testing"

	"github.com/CalebRose/SimFBA/structs"
	"github.com/jinzhu/gorm"
)

func TestTradeReportIdentityAndMessage(t *testing.T) {
	proposal := structs.NFLTradeProposal{
		Model: gorm.Model{ID: 17}, NFLTeamID: 1, NFLTeam: "AAA",
		RecepientTeamID: 2, RecepientTeam: "BBB",
	}
	player := structs.NFLPlayer{Model: gorm.Model{ID: 8}, BasePlayer: structs.BasePlayer{FirstName: "Player", LastName: "One"}}
	pick := structs.NFLDraftPick{Model: gorm.Model{ID: 9}, Season: 2030, DraftRound: 2, OriginalTeam: "CCC"}
	details := []string{tradePlayerDetail(player, "AAA", "BBB", 25), tradePickDetail(pick, "BBB", "AAA")}
	id, message, err := tradeReportIdentityAndMessage(proposal, details, false)
	if err != nil || id != "trade-proposal-nfl-17" {
		t.Fatalf("unexpected proposal report ID %q: %v", id, err)
	}
	for _, part := range []string{"#17", "AAA (team 1)", "BBB (team 2)", "Player One (player 8): AAA to BBB", "25.00%", "2030 round 2 pick (pick 9, originally CCC): BBB to AAA"} {
		if !strings.Contains(message, part) {
			t.Errorf("message %q does not contain %q", message, part)
		}
	}
	id, _, err = tradeReportIdentityAndMessage(proposal, details, true)
	if err != nil || id != "draft-trade-nfl-17" {
		t.Fatalf("unexpected draft report ID %q: %v", id, err)
	}
	proposal.ID = 0
	first, _, err := tradeReportIdentityAndMessage(proposal, details, true)
	if err != nil {
		t.Fatal(err)
	}
	second, _, err := tradeReportIdentityAndMessage(proposal, details, true)
	if err != nil || first == second || !strings.HasPrefix(first, "draft-trade-nfl-") {
		t.Fatalf("drafts without IDs need distinct report IDs: %q, %q, %v", first, second, err)
	}
	if _, _, err := tradeReportIdentityAndMessage(proposal, details, false); err == nil {
		t.Error("expected missing proposal ID to be rejected")
	}
	proposal.RecepientTeamID = proposal.NFLTeamID
	if _, _, err := tradeReportIdentityAndMessage(proposal, details, true); err == nil {
		t.Error("expected identical teams to be rejected")
	}
}
