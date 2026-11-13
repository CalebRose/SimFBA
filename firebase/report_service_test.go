package firebase

import (
	"context"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/CalebRose/SimFBA/structs"
)

func TestTradeTransactionReport(t *testing.T) {
	report, err := TradeTransactionReport(42, "Trade processed", structs.Timestamp{NFLWeekID: 102, CollegeWeekID: 999, Phase: 4})
	if err != nil {
		t.Fatal(err)
	}

	if report.League != "nfl" || report.TeamID != 42 || report.ReportType != "trade_transaction" ||
		report.Message != "Trade processed" || report.ResponseType != AdminReportTradeTransaction ||
		report.Status != "pending" || report.Phase != 4 || report.WeekID != 102 || report.ReportedBy != "system" {
		t.Fatalf("unexpected report: %+v", report)
	}
	if _, err := TradeTransactionReport(42, "trade", structs.Timestamp{NFLWeekID: -1}); err == nil {
		t.Error("negative NFL week ID should be rejected")
	}
	if strconv.IntSize == 64 {
		tooLarge := uint64(math.MaxInt64) + 1
		if _, err := TradeTransactionReport(uint(tooLarge), "trade", structs.Timestamp{}); err == nil {
			t.Error("team ID outside the Firestore integer range should be rejected")
		}
	}

}

func TestCreateAdminReportReturnsClientError(t *testing.T) {
	if os.Getenv("TEST_TRADE_REPORT_CLIENT_ERROR") == "1" {
		if err := CreateAdminReport(context.Background(), "trade-proposal-1", AdminReport{}); err == nil {
			t.Fatal("expected invalid credentials to return an error")
		}

		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=^TestCreateAdminReportReturnsClientError$")
	cmd.Env = append(os.Environ(),
		"TEST_TRADE_REPORT_CLIENT_ERROR=1",
		"FIREBASE_PRIVATE_KEY=",
		"FIREBASE_SERVICE_ACCOUNT_KEY="+filepath.Join(t.TempDir(), "missing.json"),
	)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("report client failure terminated the process: %v\n%s", err, output)
	}
}

func TestBrokenGameplanReport(t *testing.T) {
	report, err := BrokenGameplanReport("cfb", 7, "QB1 is injured", structs.Timestamp{CollegeWeekID: 31, NFLWeekID: 42, Phase: 3})
	if err != nil {
		t.Fatal(err)
	}
	if report.League != "cfb" || report.TeamID != 7 || report.ReportType != "broken_gameplan" ||
		report.Message != "QB1 is injured" || report.ResponseType != AdminReportBrokenGameplan ||
		report.Status != "pending" || report.Phase != 3 || report.WeekID != 31 || report.ReportedBy != "system" {
		t.Fatalf("unexpected report: %+v", report)
	}
	if _, err := BrokenGameplanReport("cfb", 7, "broken", structs.Timestamp{CollegeWeekID: -1}); err == nil {
		t.Error("negative college week ID should be rejected")
	}
	if strconv.IntSize == 64 {
		tooLarge := uint64(math.MaxInt64) + 1
		if _, err := BrokenGameplanReport("cfb", uint(tooLarge), "broken", structs.Timestamp{}); err == nil {
			t.Error("team ID outside the Firestore integer range should be rejected")
		}
	}
	nfl, err := BrokenGameplanReport("nfl", 7, "injured", structs.Timestamp{CollegeWeekID: 31, NFLWeekID: 42, Phase: 3})
	if err != nil || nfl.League != "nfl" || nfl.TeamID != 7 || nfl.WeekID != 42 || nfl.Phase != 3 ||
		nfl.ReportType != "broken_gameplan" || nfl.Message != "injured" ||
		nfl.ResponseType != AdminReportBrokenGameplan || nfl.Status != "pending" || nfl.ReportedBy != "system" {
		t.Fatalf("unexpected NFL report: %+v, %v", nfl, err)
	}
	if _, err := BrokenGameplanReport("nfl", 7, "broken", structs.Timestamp{NFLWeekID: -1}); err == nil {
		t.Error("negative NFL week ID should be rejected")
	}
	if _, err := BrokenGameplanReport("invalid", 7, "broken", structs.Timestamp{}); err == nil {
		t.Error("unknown league should be rejected")
	}
}
