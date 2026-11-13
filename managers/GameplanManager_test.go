package managers

import (
	"strings"
	"testing"
	"time"

	"github.com/CalebRose/SimFBA/structs"
	"github.com/jinzhu/gorm"
)

func TestCollegeDepthChartIssues(t *testing.T) {
	players := []structs.CollegePlayer{
		{Model: gorm.Model{ID: 10}, PlayerID: 110, TeamID: 1, BasePlayer: structs.BasePlayer{FirstName: "Healthy", LastName: "Player"}},
		{Model: gorm.Model{ID: 11}, PlayerID: 111, TeamID: 1, BasePlayer: structs.BasePlayer{FirstName: "Injured", LastName: "Player", IsInjured: true, WeeksOfRecovery: 3}},
		{Model: gorm.Model{ID: 12}, PlayerID: 112, TeamID: 2, BasePlayer: structs.BasePlayer{FirstName: "Transferred", LastName: "Player", IsInjured: true, WeeksOfRecovery: 1}},
		{Model: gorm.Model{ID: 13}, TeamID: 1, HasGraduated: true, BasePlayer: structs.BasePlayer{FirstName: "Graduated", LastName: "Player"}},
	}
	positions := []structs.CollegeDepthChartPosition{
		{Position: "QB", PositionLevel: "1", PlayerID: 10},
		{Position: "RB", PositionLevel: "1", PlayerID: 11},
		{Position: "WR", PositionLevel: "2", PlayerID: 12},
		{Position: "TE", PositionLevel: "1", PlayerID: 13},
		{Position: "CB", PositionLevel: "1", PlayerID: 14, FirstName: "Missing", LastName: "Player"},
		{Position: "K", PositionLevel: "1"},
		{Position: "QB", PositionLevel: "3"},
	}
	issues := collegeDepthChartIssues(1, positions, players)
	if len(issues) != 6 {
		t.Fatalf("expected 6 issues, got %d: %v", len(issues), issues)
	}
	text := strings.Join(issues, "\n")
	for _, expected := range []string{
		"RB1: Injured Player (player 11) is injured (3 weeks of recovery)",
		"WR2: Transferred Player (player 12) is injured",
		"WR2: Transferred Player (player 12) is not on the active team roster",
		"TE1: Graduated Player (player 13) is not on the active team roster",
		"CB1: Missing Player (player 14) is missing from player records",
		"K1: empty depth chart slot",
	} {
		if !strings.Contains(text, expected) {
			t.Errorf("issues %q do not contain %q", text, expected)
		}
	}
	if strings.Contains(text, "Healthy") || strings.Contains(text, "QB3") {
		t.Errorf("valid players and optional empty slots must not be reported: %s", text)
	}
}

func TestCollegeDepthChartIssuesUsesRecordID(t *testing.T) {
	positions := []structs.CollegeDepthChartPosition{{PlayerID: 110, Position: "QB", PositionLevel: "1"}}
	players := []structs.CollegePlayer{{Model: gorm.Model{ID: 10}, PlayerID: 110, TeamID: 1}}
	issues := collegeDepthChartIssues(1, positions, players)
	if len(issues) != 1 || !strings.Contains(issues[0], "missing from player records") {
		t.Fatalf("depth chart references must use college player record IDs: %v", issues)
	}

}

func TestCollegeDepthChartIssuesTreatsDeletedPlayerAsMissing(t *testing.T) {
	deleted := time.Now()
	positions := []structs.CollegeDepthChartPosition{{PlayerID: 10, Position: "QB", PositionLevel: "1"}}
	players := []structs.CollegePlayer{{Model: gorm.Model{ID: 10, DeletedAt: &deleted}, TeamID: 1}}
	issues := collegeDepthChartIssues(1, positions, players)
	if len(issues) != 1 || !strings.Contains(issues[0], "missing from player records") {
		t.Fatalf("deleted players should be treated as missing: %v", issues)
	}
}

func TestCollegeDepthChartIssuesIncludesRedshirts(t *testing.T) {
	positions := []structs.CollegeDepthChartPosition{{PlayerID: 10, Position: "QB", PositionLevel: "1"}}
	players := []structs.CollegePlayer{{Model: gorm.Model{ID: 10}, TeamID: 1, IsRedshirting: true, BasePlayer: structs.BasePlayer{FirstName: "Redshirt", LastName: "Player"}}}
	issues := collegeDepthChartIssues(1, positions, players)
	if len(issues) != 1 || !strings.Contains(issues[0], "QB1: Redshirt Player (player 10) is redshirting") {
		t.Fatalf("expected redshirt issue: %v", issues)
	}
}

func TestNFLDepthChartIssues(t *testing.T) {
	deleted := time.Now()
	players := []structs.NFLPlayer{
		{Model: gorm.Model{ID: 10}, PlayerID: 110, TeamID: 1, BasePlayer: structs.BasePlayer{FirstName: "Healthy", LastName: "Player"}},
		{Model: gorm.Model{ID: 11}, TeamID: 1, BasePlayer: structs.BasePlayer{FirstName: "Injured", LastName: "Player", IsInjured: true, WeeksOfRecovery: 3}},
		{Model: gorm.Model{ID: 12}, TeamID: 2, IsPracticeSquad: true, BasePlayer: structs.BasePlayer{FirstName: "Moved", LastName: "Player"}},
		{Model: gorm.Model{ID: 13, DeletedAt: &deleted}, TeamID: 1},
	}
	positions := []structs.NFLDepthChartPosition{
		{Position: "QB", PositionLevel: "1", PlayerID: 10},
		{Position: "RB", PositionLevel: "1", PlayerID: 11},
		{Position: "WR", PositionLevel: "2", PlayerID: 12},
		{Position: "TE", PositionLevel: "1", PlayerID: 13, FirstName: "Deleted", LastName: "Player"},
		{Position: "CB", PositionLevel: "1", PlayerID: 14, FirstName: "Missing", LastName: "Player"},
		{Position: "K", PositionLevel: "1"},
		{Position: "QB", PositionLevel: "3"},
	}
	issues := nflDepthChartIssues(1, positions, players)
	if len(issues) != 6 {
		t.Fatalf("expected 6 issues, got %d: %v", len(issues), issues)
	}
	text := strings.Join(issues, "\n")
	for _, expected := range []string{
		"RB1: Injured Player (player 11) is injured (3 weeks of recovery)",
		"WR2: Moved Player (player 12) is on the practice squad",
		"WR2: Moved Player (player 12) is not on the team roster",
		"TE1: Deleted Player (player 13) is missing from player records",
		"CB1: Missing Player (player 14) is missing from player records",
		"K1: empty depth chart slot",
	} {
		if !strings.Contains(text, expected) {
			t.Errorf("issues %q do not contain %q", text, expected)
		}
	}
	if strings.Contains(text, "Healthy") || strings.Contains(text, "QB3") {
		t.Errorf("valid players and optional empty slots must not be reported: %s", text)
	}
}
