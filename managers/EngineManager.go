package managers

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"reflect"
	"sort"

	"github.com/CalebRose/SimFBA/structs"
)

/*
For engine related testing and development

*/

// Attributes to report on. Names must match BasePlayer field names exactly
// (used via reflection so this list is the single source of truth).
var attributeFields = []string{
	"Speed", "Agility", "FootballIQ", "Strength", "Carrying", "Catching",
	"RouteRunning", "ZoneCoverage", "ManCoverage", "Tackle", "PassBlock",
	"RunBlock", "PassRush", "RunDefense", "ThrowPower", "ThrowAccuracy",
	"KickAccuracy", "KickPower", "PuntAccuracy", "PuntPower",
}

// Minimum group size before we trust a percentile read. Below this we
// still compute it, but flag it in the "Reliable" column.
const minReliableSampleSize = 10

// --- Percentile math ---------------------------------------------------

// linearInterpPercentile matches numpy's default ("linear") method, which
// is the standard approach and how R-7 / Excel PERCENTILE.INC behave.
func linearInterpPercentile(sorted []float64, p float64) float64 {
	n := len(sorted)
	if n == 0 {
		return 0
	}
	if n == 1 {
		return sorted[0]
	}
	idx := p / 100 * float64(n-1)
	lo := int(math.Floor(idx))
	hi := int(math.Ceil(idx))
	if lo == hi {
		return sorted[lo]
	}
	weight := idx - float64(lo)
	return sorted[lo]*(1-weight) + sorted[hi]*weight
}

type stats struct {
	Count    int
	Min      float64
	P10      float64
	P25      float64
	Median   float64
	P75      float64
	P90      float64
	P95      float64
	P99      float64
	Max      float64
	Reliable bool
}

func computeStats(values []int8) (stats, bool) {
	if len(values) == 0 {
		return stats{}, false
	}
	// Skip attributes that are entirely unused/default for this group
	// (e.g. KickPower for a WR group) — all-zero means "not applicable".
	allZero := true
	for _, v := range values {
		if v != 0 {
			allZero = false
			break
		}
	}
	if allZero {
		return stats{}, false
	}

	f := make([]float64, len(values))
	for i, v := range values {
		f[i] = float64(v)
	}
	sort.Float64s(f)

	return stats{
		Count:    len(f),
		Min:      f[0],
		P10:      linearInterpPercentile(f, 10),
		P25:      linearInterpPercentile(f, 25),
		Median:   linearInterpPercentile(f, 50),
		P75:      linearInterpPercentile(f, 75),
		P90:      linearInterpPercentile(f, 90),
		P95:      linearInterpPercentile(f, 95),
		P99:      linearInterpPercentile(f, 99),
		Max:      f[len(f)-1],
		Reliable: len(f) >= minReliableSampleSize,
	}, true
}

// --- Grouping ------------------------------------------------------------

type groupKey struct {
	Scope     string // "Position" or "Position+Archetype"
	Position  string
	Archetype string // empty for Position-only scope
}

func getAttrValue(bp structs.BasePlayer, field string) int8 {
	v := reflect.ValueOf(bp).FieldByName(field)
	return int8(v.Int())
}

func buildCFBGroups(players []structs.CollegePlayer, byArchetype bool) map[groupKey]map[string][]int8 {
	groups := make(map[groupKey]map[string][]int8)

	for _, p := range players {
		// Adjust this filter to match what population you actually want
		// (e.g. exclude waived/free agents, or include everyone).

		var key groupKey
		if byArchetype {
			key = groupKey{Scope: "Position+Archetype", Position: p.Position, Archetype: p.Archetype}
		} else {
			key = groupKey{Scope: "Position", Position: p.Position}
		}

		if _, ok := groups[key]; !ok {
			groups[key] = make(map[string][]int8)
		}
		for _, field := range attributeFields {
			val := getAttrValue(p.BasePlayer, field)
			groups[key][field] = append(groups[key][field], val)
		}
	}
	return groups
}

func buildNFLGroups(players []structs.NFLPlayer, byArchetype bool) map[groupKey]map[string][]int8 {
	groups := make(map[groupKey]map[string][]int8)

	for _, p := range players {
		// Adjust this filter to match what population you actually want
		// (e.g. exclude waived/free agents, or include everyone).
		if p.IsWaived {
			continue
		}

		var key groupKey
		if byArchetype {
			key = groupKey{Scope: "Position+Archetype", Position: p.Position, Archetype: p.Archetype}
		} else {
			key = groupKey{Scope: "Position", Position: p.Position}
		}

		if _, ok := groups[key]; !ok {
			groups[key] = make(map[string][]int8)
		}
		for _, field := range attributeFields {
			val := getAttrValue(p.BasePlayer, field)
			groups[key][field] = append(groups[key][field], val)
		}
	}
	return groups
}

// --- Report generation ---------------------------------------------------

func writeReport(w *csv.Writer, groups map[groupKey]map[string][]int8) {
	for key, attrs := range groups {
		for _, field := range attributeFields {
			values := attrs[field]
			s, ok := computeStats(values)
			if !ok {
				continue // attribute not applicable to this group
			}
			w.Write([]string{
				key.Scope,
				key.Position,
				key.Archetype,
				field,
				fmt.Sprintf("%d", s.Count),
				fmt.Sprintf("%.2f", s.Min),
				fmt.Sprintf("%.2f", s.P10),
				fmt.Sprintf("%.2f", s.P25),
				fmt.Sprintf("%.2f", s.Median),
				fmt.Sprintf("%.2f", s.P75),
				fmt.Sprintf("%.2f", s.P90),
				fmt.Sprintf("%.2f", s.P95),
				fmt.Sprintf("%.2f", s.P99),
				fmt.Sprintf("%.2f", s.Max),
				fmt.Sprintf("%t", s.Reliable),
			})
		}
	}
}

func NFLAttributePercentilesReport() {
	players := GetAllNFLPlayers()

	byArchetype := buildNFLGroups(players, true)
	byPosition := buildNFLGroups(players, false)

	f, err := os.Create("nfl_attribute_percentiles.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	w.Write([]string{
		"Scope", "Position", "Archetype", "Attribute", "Count",
		"Min", "P10", "P25", "Median", "P75", "P90", "P95", "P99", "Max",
		"Reliable",
	})

	writeReport(w, byPosition)
	writeReport(w, byArchetype)

	fmt.Println("Wrote attribute_percentiles.csv")
}

func CFBAttributePercentilesReport() {
	players := GetAllCollegePlayers()

	byArchetype := buildCFBGroups(players, true)
	byPosition := buildCFBGroups(players, false)

	f, err := os.Create("cfb_attribute_percentiles.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	w.Write([]string{
		"Scope", "Position", "Archetype", "Attribute", "Count",
		"Min", "P10", "P25", "Median", "P75", "P90", "P95", "P99", "Max",
		"Reliable",
	})

	writeReport(w, byPosition)
	writeReport(w, byArchetype)

	fmt.Println("Wrote attribute_percentiles.csv")
}
