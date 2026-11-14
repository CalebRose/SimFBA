package managers

import (
	"math/rand"
	"os"
	"reflect"
	"testing"

	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
)

// Go 1.24+ makes rand.Seed a no-op by default. These tests exercise the existing
// global-RNG progression code with repeatable, in-memory players; no DB is used.
func enableProgressionTestSeed(t *testing.T) {
	t.Helper()
	t.Setenv("GODEBUG", os.Getenv("GODEBUG")+",randseednop=0")
}

func TestSpecialistSnapModifierBoundaries(t *testing.T) {
	for _, tc := range []struct {
		snaps    int
		min, max float64
	}{
		{0, 0.5, 0.75}, {1, 0.5, 0.75}, {11, 0.5, 0.75},
		{12, 0.9, 1.1}, {23, 0.9, 1.1},
		{24, 1.0, 1.25}, {25, 1.0, 1.25}, {600, 1.0, 1.25},
	} {
		for i := 0; i < 32; i++ {
			got := getSpecialistSnapModifier(tc.snaps)
			if got < tc.min || got >= tc.max {
				t.Fatalf("%d season snaps: modifier %v outside [%v, %v)", tc.snaps, got, tc.min, tc.max)
			}
		}
	}
}

func TestSpecialistProgressionRouting(t *testing.T) {
	enableProgressionTestSeed(t)
	for _, tc := range []struct{ position, attribute string }{
		{"K", "Kick Power"}, {"P", "Punt Power"},
	} {
		rand.Seed(1)
		got := primarySpecialistProgression(70, 40, tc.position, "Balanced", 24, 0, tc.attribute, false)
		rand.Seed(1)
		legacy := PrimaryProgression(70, 40, tc.position, "Balanced", 0, tc.attribute, false)
		if got <= legacy {
			t.Errorf("%s: 24 season snaps must beat the unchanged zero-SPG path: %d <= %d", tc.position, got, legacy)
		}
		for _, total := range []int{0, 11} {
			rand.Seed(2)
			got = primarySpecialistProgression(70, 40, tc.position, "Balanced", total, 0, tc.attribute, false)
			rand.Seed(2)
			want := PrimaryProgression(70, 40, tc.position, "Balanced", 0, tc.attribute, false)
			if got != want {
				t.Errorf("%s: %d season snaps changed the baseline: %d != %d", tc.position, total, got, want)
			}
		}
		// Redshirts and practice-squad players retain their existing separate formula.
		rand.Seed(3)
		got = primarySpecialistProgression(70, 40, tc.position, "Balanced", 24, 0, tc.attribute, true)
		rand.Seed(3)
		want := PrimaryProgression(70, 40, tc.position, "Balanced", 0, tc.attribute, true)
		if got != want {
			t.Errorf("%s: redshirt/practice-squad behavior changed: %d != %d", tc.position, got, want)
		}
	}
	for _, position := range []string{"QB", "RB", "WR", "DE", "ATH", "KR", "PR", "ST"} {
		for _, spg := range []int{0, 10, 11, 20, 21, 30, 31, 45, 46} {
			rand.Seed(4)
			got := primarySpecialistProgression(70, 40, position, "Balanced", 1000, spg, "Speed", false)
			rand.Seed(4)
			want := PrimaryProgression(70, 40, position, "Balanced", spg, "Speed", false)
			if got != want {
				t.Errorf("non-specialist %s at %d SPG used specialist tiers: %d != %d", position, spg, got, want)
			}
		}
	}
	if got := primarySpecialistProgression(70, 0, "K", "Power", 24, 0, "Kick Power", false); got != 1 {
		t.Errorf("zero attribute behavior changed: %d", got)
	}
	if got := primarySpecialistProgression(70, 99, "P", "Power", 24, 0, "Punt Power", false); got != 99 {
		t.Errorf("attribute cap changed: %d", got)
	}
}

func TestSecondaryArchetypeMappingCorrections(t *testing.T) {
	for _, tc := range []struct{ position, archetype, destination, want string }{
		{"DE", "Speed Rusher", "DT", "Pass Rusher"},
		{"DE", "Speed Rusher", "OLB", "Pass Rush"},
		{"DE", "Speed Rusher", "ILB", "Speed"},
		{"DT", "Pass Rusher", "DE", "Speed Rusher"},
		{"OLB", "Pass Rush", "DE", "Speed Rusher"},
		{"OLB", "Speed", "DE", "Speed Rusher"},
		{"ILB", "Speed", "DE", "Speed Rusher"},
		{"K", "Balanced", "P", "Balanced"},
		{"P", "Balanced", "K", "Balanced"},
	} {
		got, ok := getNewArchetype(tc.position, tc.archetype, tc.destination)
		if !ok || got != tc.want {
			t.Errorf("%s/%s -> %s: got %q, %v; want %q", tc.position, tc.archetype, tc.destination, got, ok, tc.want)
		}
	}
	for _, position := range []string{"K", "P"} {
		if _, ok := getNewArchetype(position, "Balanced", "QB"); ok {
			t.Errorf("Balanced %s unexpectedly gained a new QB conversion", position)
		}
	}
	for position, archetypes := range util.GetNewArchetypeMap() {
		for archetype, destinations := range archetypes {
			if position == "DE" && archetype == "Speed Rush" {
				t.Error("noncanonical DE source key remains")
			}
			if destinations["DE"] == "Speed Rush" {
				t.Errorf("%s/%s still emits the noncanonical DE archetype", position, archetype)
			}
		}
	}
}

func TestBalancedRBPassBlockingBonus(t *testing.T) {
	if got := GetArchetypeMod("RB", "Balanced", "Pass Blocking"); got != 0.05 {
		t.Errorf("Balanced RB pass-blocking bonus = %v; want 0.05", got)
	}
	if got := GetArchetypeMod("RB", "Balanced", "Agility"); got != 0.05 {
		t.Errorf("Balanced RB agility bonus changed: %v", got)
	}
}

func TestOffseasonSpecialistsUseRawSeasonSnaps(t *testing.T) {
	enableProgressionTestSeed(t)
	for _, position := range []string{"K", "P"} {
		player := progressionTestPlayer(position, "Balanced")
		// Only 4 positional snaps remain after ST exclusion. All 24 generic
		// season snaps must still count toward specialist attribute progression.
		lowSnaps := progressionTestSnaps(position, 4)
		lowSnaps.STSnaps = 20
		highSnaps := progressionTestSnaps(position, 600)
		for _, league := range []string{"CFB", "NFL"} {
			rand.Seed(10)
			got := runProgressionFixture(league, player, 24, lowSnaps)
			rand.Seed(10)
			want := runProgressionFixture(league, player, 600, highSnaps)
			if !reflect.DeepEqual(progressionTestAttributes(got), progressionTestAttributes(want)) {
				t.Errorf("%s %s: raw 24 snaps did not use the same top tier as 600 snaps", league, position)
			}
		}
	}
}

func TestOffseasonUsesNewSecondaryArchetype(t *testing.T) {
	enableProgressionTestSeed(t)
	for _, tc := range []struct{ position, archetype, destination, newArchetype string }{
		{"RB", "Receiving", "WR", "Route Runner"},
		{"DE", "Speed Rusher", "DT", "Pass Rusher"},
		{"K", "Balanced", "P", "Balanced"},
	} {
		for _, league := range []string{"CFB", "NFL"} {
			player := progressionTestPlayer(tc.position, tc.archetype)
			snaps := progressionTestSnaps(tc.destination, 500)
			rand.Seed(20)
			got := runProgressionFixture(league, player, 500, snaps)
			if got.Position != tc.position || got.Archetype != tc.archetype ||
				got.PositionTwo != tc.destination || got.ArchetypeTwo != tc.newArchetype {
				t.Fatalf("%s %s: primary/secondary identity not preserved: %s/%s + %s/%s", league, tc.position, got.Position, got.Archetype, got.PositionTwo, got.ArchetypeTwo)
			}
			// A player already at the destination provides the same attribute
			// calculation. Account for RNG used to construct the conversion map.
			player.Position, player.Archetype = tc.destination, tc.newArchetype
			rand.Seed(20)
			util.GetNewArchetypeMap()
			want := runProgressionFixture(league, player, 500, snaps)
			if !reflect.DeepEqual(progressionTestAttributes(got), progressionTestAttributes(want)) {
				t.Errorf("%s %s/%s -> %s: progression did not use %s archetype", league, tc.position, tc.archetype, tc.destination, tc.newArchetype)
			}
		}
	}
}

func TestSecondaryPositionThresholdAndExistingSlotUnchanged(t *testing.T) {
	for _, league := range []string{"CFB", "NFL"} {
		for _, count := range []int{80, 81} {
			player := progressionTestPlayer("RB", "Receiving")
			snaps := progressionTestSnaps("WR", count)
			snaps.RBSnaps = uint16(100 - count)
			got := runProgressionFixture(league, player, 100, snaps)
			if (got.PositionTwo == "WR") != (count == 81) {
				t.Errorf("%s: secondary-position >80%% threshold changed at %d snaps", league, count)
			}
			player.PositionTwo, player.ArchetypeTwo = "FB", "Receiving"
			got = runProgressionFixture(league, player, 100, snaps)
			if got.PositionTwo != "FB" || got.ArchetypeTwo != "Receiving" {
				t.Errorf("%s: an existing secondary position was replaced", league)
			}
		}
	}
}

func progressionTestPlayer(position, archetype string) structs.BasePlayer {
	return structs.BasePlayer{
		Position: position, Archetype: archetype, Progression: 70, Age: 22, PrimeAge: 30,
		Agility: 40, ThrowPower: 40, ThrowAccuracy: 40, Speed: 40, FootballIQ: 40,
		Strength: 40, RunBlock: 40, PassBlock: 40, RunDefense: 40, PassRush: 40,
		Carrying: 40, Tackle: 40, RouteRunning: 40, Catching: 40,
		PuntPower: 40, PuntAccuracy: 40, KickPower: 40, KickAccuracy: 40,
		ManCoverage: 40, ZoneCoverage: 40,
	}
}

func progressionTestSnaps(position string, count int) structs.BasePlayerSeasonSnaps {
	game := structs.BasePlayerGameSnaps{}
	game.MapSnapsToPosition(position, count)
	season := structs.BasePlayerSeasonSnaps{}
	season.AddToSeason(game)
	return season
}

func runProgressionFixture(league string, player structs.BasePlayer, total int, snaps structs.BasePlayerSeasonSnaps) structs.BasePlayer {
	if league == "CFB" {
		stats := []structs.CollegePlayerStats{{BasePlayerStats: structs.BasePlayerStats{Snaps: int16(total)}}}
		result := ProgressCollegePlayer(structs.CollegePlayer{BasePlayer: player, Year: 2}, "7", stats,
			structs.CollegePlayerSeasonSnaps{BasePlayerSeasonSnaps: snaps})
		return result.BasePlayer
	}
	result := ProgressNFLPlayer(structs.NFLPlayer{BasePlayer: player}, "7", total, total/17,
		structs.NFLPlayerSeasonSnaps{BasePlayerSeasonSnaps: snaps})
	return result.BasePlayer
}

func progressionTestAttributes(p structs.BasePlayer) []int8 {
	return []int8{p.Agility, p.ThrowPower, p.ThrowAccuracy, p.Speed, p.FootballIQ,
		p.Strength, p.RunBlock, p.PassBlock, p.RunDefense, p.PassRush, p.Carrying,
		p.Tackle, p.RouteRunning, p.Catching, p.PuntPower, p.PuntAccuracy,
		p.KickPower, p.KickAccuracy, p.ManCoverage, p.ZoneCoverage}
}
