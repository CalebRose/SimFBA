package managers

import (
	"log"
	"math"
	"sort"
	"strconv"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
)

// =============================================================================
// Constants
// =============================================================================

const snapPositionFloor uint16 = 50

// starterSnapThreshold is the minimum snaps-per-season to qualify as a starter
// for playtime tag type determination.
const starterSnapThreshold = 765

// =============================================================================
// Internal types
// =============================================================================

type playerGroupEntry struct {
	Player   structs.NFLPlayer
	Contract structs.NFLContract
}

type rankedGroupEntry struct {
	Entry      playerGroupEntry
	AdjOverall int // overall with +4 adjustment for players <= age 25
}

// =============================================================================
// Position-group helpers
// =============================================================================

// getPositionGroup returns the calculation group for a player.
// EDGE grouping is always determined by the stored position + archetype;
// snap-based position is used for all other positions.
func getPositionGroup(storedPos, archetype, snapPos string) string {
	// EDGE: Speed Rusher DEs and Pass Rush OLBs
	if storedPos == "DE" && archetype == "Speed Rusher" {
		return "EDGE"
	}
	if storedPos == "OLB" && archetype == "Pass Rush" {
		return "EDGE"
	}

	effectivePos := snapPos
	if effectivePos == "" {
		effectivePos = storedPos
	}

	switch effectivePos {
	case "DE", "DT":
		return "DL"
	case "OLB", "ILB":
		return "LB"
	default:
		return effectivePos // QB, RB, FB, WR, TE, OT, OG, C, CB, FS, SS, K, P
	}
}

// determineSnapPosition returns the position where the player played the most
// snaps this season, provided it meets the minimum floor. Falls back to
// storedPos if no position reaches the floor.
func determineSnapPosition(storedPos string, snaps structs.NFLPlayerSeasonSnaps) string {
	snapsByPos := map[string]uint16{
		"QB":  snaps.QBSnaps,
		"RB":  snaps.RBSnaps,
		"FB":  snaps.FBSnaps,
		"WR":  snaps.WRSnaps,
		"TE":  snaps.TESnaps,
		"OT":  snaps.OTSnaps,
		"OG":  snaps.OGSnaps,
		"C":   snaps.CSnaps,
		"DE":  snaps.DESnaps,
		"DT":  snaps.DTSnaps,
		"OLB": snaps.OLBSnaps,
		"ILB": snaps.ILBSnaps,
		"CB":  snaps.CBSnaps,
		"FS":  snaps.FSSnaps,
		"SS":  snaps.SSSnaps,
		"K":   snaps.KSnaps,
		"P":   snaps.PSnaps,
	}

	var bestPos string
	var maxSnaps uint16
	for pos, count := range snapsByPos {
		if count >= snapPositionFloor && count > maxSnaps {
			maxSnaps = count
			bestPos = pos
		}
	}
	if bestPos == "" {
		return storedPos
	}
	return bestPos
}

// getTopTierCount returns the number of players in the "elite" tier for a group.
func getTopTierCount(group string) int {
	switch group {
	case "OT", "OG", "DL", "EDGE", "LB", "WR", "CB":
		return 10
	default: // QB, RB, FB, TE, C, FS, SS, K, P
		return 5
	}
}

// getAgeExclusionThreshold returns the minimum age (inclusive) at which a
// player is excluded from the custom mid-tier comparison set (step 5).
func getAgeExclusionThreshold(group string) int {
	switch group {
	case "QB", "K":
		return 34
	case "RB", "FB":
		return 27
	default:
		return 29
	}
}

// getAgeAdjustmentFactor returns the decimal adjustment factor for a player's
// age and position group (e.g., +0.15 means +15%, -0.10 means −10%).
func getAgeAdjustmentFactor(group string, age int) float64 {
	switch group {
	case "QB", "K":
		if age >= 36 {
			return -0.90
		}
		factors := map[int]float64{
			23: 0.15, 24: 0.10, 25: 0.05, 26: 0.00,
			27: -0.05, 28: -0.10, 29: -0.15, 30: -0.20,
			31: -0.25, 32: -0.30, 33: -0.45, 34: -0.60, 35: -0.75,
		}
		if f, ok := factors[age]; ok {
			return f
		}
		return 0.00

	case "RB", "FB":
		if age >= 32 {
			return -0.90
		}
		factors := map[int]float64{
			23: 0.10, 24: 0.05, 25: 0.00, 26: -0.10,
			27: -0.20, 28: -0.30, 29: -0.45, 30: -0.60, 31: -0.75,
		}
		if f, ok := factors[age]; ok {
			return f
		}
		return 0.00

	default:
		if age >= 34 {
			return -0.90
		}
		factors := map[int]float64{
			23: 0.15, 24: 0.10, 25: 0.05, 26: 0.00,
			27: -0.10, 28: -0.20, 29: -0.30, 30: -0.40,
			31: -0.50, 32: -0.60, 33: -0.75,
		}
		if f, ok := factors[age]; ok {
			return f
		}
		return 0.00
	}
}

func avgFloats(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	var total float64
	for _, v := range vals {
		total += v
	}
	return total / float64(len(vals))
}

// =============================================================================
// Extension value calculation (per group)
// =============================================================================

// computeGroupExtensionValues calculates the expected minimum contract value
// and AAV for every player in a position group, following the steps in the
// TechDocs/extension_value_and_tag_calculations.md document.
func computeGroupExtensionValues(group string, entries []playerGroupEntry) []structs.NFLPlayer {
	if len(entries) == 0 {
		return nil
	}

	topTierCount := getTopTierCount(group)
	ageThreshold := getAgeExclusionThreshold(group)

	// Step 2: compute adjusted overall (age <= 25 gets +4 bonus for ranking)
	ranked := make([]rankedGroupEntry, len(entries))
	for i, e := range entries {
		adj := int(e.Player.Overall)
		if int(e.Player.Age) <= 25 {
			adj += 4
		}
		ranked[i] = rankedGroupEntry{Entry: e, AdjOverall: adj}
	}

	// Step 3: sort by adjusted overall desc, age asc
	sort.Slice(ranked, func(i, j int) bool {
		if ranked[i].AdjOverall != ranked[j].AdjOverall {
			return ranked[i].AdjOverall > ranked[j].AdjOverall
		}
		return ranked[i].Entry.Player.Age < ranked[j].Entry.Player.Age
	})

	// Determine the adjusted-overall floor for the top tier.
	// All players tied at the Nth place are included (per "tied at 5th" rule).
	topTierMinOverall := math.MinInt32 // default: all players are top tier
	if topTierCount < len(ranked) {
		topTierMinOverall = ranked[topTierCount-1].AdjOverall
	}

	// Step 4a: collect active contract signing values for top-tier reference
	var groupSigningValues []float64
	for _, e := range entries {
		if e.Contract.IsActive {
			groupSigningValues = append(groupSigningValues, e.Contract.SigningValue)
		}
	}
	sort.Slice(groupSigningValues, func(i, j int) bool {
		return groupSigningValues[i] > groupSigningValues[j]
	})

	var maxContractValue float64
	if len(groupSigningValues) > 0 {
		maxContractValue = groupSigningValues[0]
		// Caruso-Bordewyk Rule: if the highest value is >= 150% of the second
		// highest, exclude it to avoid distorting the top-tier reference.
		if len(groupSigningValues) > 1 && maxContractValue >= groupSigningValues[1]*1.5 {
			maxContractValue = groupSigningValues[1]
		}
	}
	topTierExpectedValue := maxContractValue * 1.10

	// First pass: assign raw (pre-smoothing) expected values
	rawValues := make([]float64, len(ranked))
	for i, rp := range ranked {
		if rp.AdjOverall >= topTierMinOverall {
			// Step 4: elite tier — max contract value + 10%
			rawValues[i] = topTierExpectedValue
		} else {
			// Steps 5–6: mid-tier — highest signing value among players with
			// equal-or-lower actual overall in the filtered custom set
			actualOverall := int(rp.Entry.Player.Overall)
			var bestValue float64
			for _, e := range entries {
				c := e.Contract
				if !c.IsActive {
					continue
				}
				// Exclude rookies, UDFAs, and age-ineligible players
				if c.ContractType == "Rookie" || c.ContractType == "UDFA" {
					continue
				}
				if int(e.Player.Age) >= ageThreshold {
					continue
				}
				if int(e.Player.Overall) <= actualOverall && c.SigningValue > bestValue {
					bestValue = c.SigningValue
				}
			}
			rawValues[i] = bestValue
		}
	}

	// Step 7: smoothing — average values of players at (overall+1) and (overall-1)
	adjOverallValMap := make(map[int][]float64)
	for i, rp := range ranked {
		adjOverallValMap[rp.AdjOverall] = append(adjOverallValMap[rp.AdjOverall], rawValues[i])
	}

	result := make([]structs.NFLPlayer, len(ranked))
	for i, rp := range ranked {
		rawVal := rawValues[i]

		// Build smoothing components from adjacent overall levels
		var smoothComponents []float64
		if higherVals, ok := adjOverallValMap[rp.AdjOverall+1]; ok {
			smoothComponents = append(smoothComponents, avgFloats(higherVals))
		}
		if lowerVals, ok := adjOverallValMap[rp.AdjOverall-1]; ok {
			smoothComponents = append(smoothComponents, avgFloats(lowerVals))
		}
		smoothedVal := avgFloats(smoothComponents)

		// Step 8: take the higher of raw and smoothed values
		bestVal := rawVal
		if smoothedVal > bestVal {
			bestVal = smoothedVal
		}

		// Step 9: apply age adjustment
		ageFactor := getAgeAdjustmentFactor(group, int(rp.Entry.Player.Age))
		adjustedVal := bestVal * (1.0 + ageFactor)
		if adjustedVal < 0 {
			adjustedVal = 0
		}

		// Step 10: AAV = 40% of value expectation
		aav := adjustedVal * 0.40

		p := rp.Entry.Player
		p.AssignCalculatedValues(adjustedVal, aav)
		result[i] = p
	}

	return result
}

// =============================================================================
// CalculatePlayerMinimumAndAAVValues
// =============================================================================

// CalculatePlayerMinimumAndAAVValues calculates and assigns the minimum
// contract (extension) value and AAV for every active NFL player.
// Run once per offseason. Also updates OriginalMinimumValue / OriginalAAV so
// the week-15 reset restores to this freshly computed baseline.
func CalculatePlayerMinimumAndAAVValues() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID := strconv.Itoa(ts.NFLSeasonID)

	nflPlayers := GetAllNFLPlayers()
	contracts := repository.FindAllActiveNFLContracts()
	seasonSnapsMap := GetNFLPlayerSeasonSnapMap(seasonID)

	// Build contract lookup keyed by NFL player ID
	contractMap := make(map[int]structs.NFLContract, len(contracts))
	for _, c := range contracts {
		contractMap[c.NFLPlayerID] = c
	}

	// Group all players by their position group
	groupEntries := make(map[string][]playerGroupEntry)
	for _, p := range nflPlayers {
		snapPos := determineSnapPosition(p.Position, seasonSnapsMap[p.ID])
		group := getPositionGroup(p.Position, p.Archetype, snapPos)
		groupEntries[group] = append(groupEntries[group], playerGroupEntry{
			Player:   p,
			Contract: contractMap[int(p.ID)],
		})
	}

	// Calculate extension values per group, then save
	saved := 0
	for group, entries := range groupEntries {
		updated := computeGroupExtensionValues(group, entries)
		for _, p := range updated {
			if err := db.Model(&p).Updates(map[string]interface{}{
				"minimum_value":          p.MinimumValue,
				"original_minimum_value": p.OriginalMinimumValue,
				"aav":                    p.AAV,
				"original_aav":           p.OriginalAAV,
			}).Error; err != nil {
				log.Printf("ValuationManager: failed to save player %d (%s): %v", p.ID, group, err)
				continue
			}
			saved++
		}
	}

	log.Printf("CalculatePlayerMinimumAndAAVValues: updated %d players", saved)
}

// =============================================================================
// CalculateTagValues
// =============================================================================

// CalculateTagValues computes tag amounts for every position group based on
// current-year pay (Y1BaseSalary + Y1Bonus) and saves them to the NFLTagData
// table. It also assigns TagType to each player (Franchise / Transition /
// Playtime / Basic) based on Pro Bowl selections and snap history.
func CalculateTagValues() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID := strconv.Itoa(ts.NFLSeasonID)

	nflPlayers := GetAllNFLPlayers()
	contracts := repository.FindAllActiveNFLContracts()
	allSeasonSnaps := repository.FindAllNFLPlayerSeasonSnaps()
	seasonSnapsMap := GetNFLPlayerSeasonSnapMap(seasonID)

	// Contract lookup by player ID
	contractMap := make(map[int]structs.NFLContract, len(contracts))
	for _, c := range contracts {
		contractMap[c.NFLPlayerID] = c
	}

	// Historical snap lookup: playerID -> seasonID -> snaps
	histSnapsMap := make(map[uint]map[uint]structs.NFLPlayerSeasonSnaps)
	for _, s := range allSeasonSnaps {
		if histSnapsMap[s.PlayerID] == nil {
			histSnapsMap[s.PlayerID] = make(map[uint]structs.NFLPlayerSeasonSnaps)
		}
		histSnapsMap[s.PlayerID][s.SeasonID] = s
	}

	// Group players (with active contracts) by position group, collecting pay
	type payEntry struct {
		Player   structs.NFLPlayer
		Pay      float64 // Y1BaseSalary + Y1Bonus
		Contract structs.NFLContract
	}
	groupPayEntries := make(map[string][]payEntry)

	for _, p := range nflPlayers {
		c, ok := contractMap[int(p.ID)]
		if !ok || !c.IsActive {
			continue
		}
		snapPos := determineSnapPosition(p.Position, seasonSnapsMap[p.ID])
		group := getPositionGroup(p.Position, p.Archetype, snapPos)
		groupPayEntries[group] = append(groupPayEntries[group], payEntry{
			Player:   p,
			Pay:      c.Y1BaseSalary + c.Y1Bonus,
			Contract: c,
		})
	}

	// Compute and persist tag amounts for each group
	existingTagData := repository.FindAllNFLTagData()
	tagDataByPos := make(map[string]structs.NFLTagData, len(existingTagData))
	for _, td := range existingTagData {
		tagDataByPos[td.Position] = td
	}

	for group, entries := range groupPayEntries {
		// Sort by current-year pay descending
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Pay > entries[j].Pay
		})

		pays := make([]float64, len(entries))
		for i, e := range entries {
			pays[i] = e.Pay
		}

		td := tagDataByPos[group]
		td.Position = group
		td.Franchise = avgFloats(topN(pays, 5))
		td.Transition = avgFloats(topN(pays, 10))
		td.Playtime = avgFloats(rangeN(pays, 2, 19)) // ranks 3–20 (0-indexed 2–19)
		td.Basic = avgFloats(rangeN(pays, 2, 24))    // ranks 3–25 (0-indexed 2–24)

		repository.SaveNFLTagData(td, db)
	}

	// Assign TagType to each player based on Pro Bowls and first-3-season snaps
	for _, p := range nflPlayers {
		newTagType := determineTagType(p, ts.NFLSeasonID, histSnapsMap)
		if newTagType == p.TagType {
			continue
		}
		if err := db.Model(&p).Update("tag_type", newTagType).Error; err != nil {
			log.Printf("ValuationManager: failed to update tag_type for player %d: %v", p.ID, err)
		}
	}

	log.Printf("CalculateTagValues: tag amounts computed for %d position groups", len(groupPayEntries))
}

// topN returns the first n elements of an already-sorted (descending) slice.
func topN(sorted []float64, n int) []float64 {
	if n > len(sorted) {
		n = len(sorted)
	}
	return sorted[:n]
}

// rangeN returns elements at 0-based indices [start, end] of a sorted slice.
// Both bounds are inclusive. If the slice is shorter than needed, returns what
// is available within the range.
func rangeN(sorted []float64, start, end int) []float64 {
	if start >= len(sorted) {
		return nil
	}
	if end >= len(sorted) {
		end = len(sorted) - 1
	}
	return sorted[start : end+1]
}

// determineTagType returns the appropriate TagType constant for a player:
//
//	0 == Basic
//	1 == Franchise  (multiple Pro Bowls)
//	2 == Transition (one Pro Bowl)
//	3 == Playtime   (starter snaps in 2 of first 3 seasons)
func determineTagType(p structs.NFLPlayer, currentSeasonID int, histSnaps map[uint]map[uint]structs.NFLPlayerSeasonSnaps) uint8 {
	if p.ProBowls >= 2 {
		return 1 // Franchise
	}
	if p.ProBowls == 1 {
		return 2 // Transition
	}

	// Playtime: >= 765 snaps in at least 2 of the player's first 3 NFL seasons
	if meetsPlaytimeCriteria(p, currentSeasonID, histSnaps) {
		return 3 // Playtime
	}

	return 0 // Basic
}

// meetsPlaytimeCriteria returns true if the player had starter-level snaps
// (>= 765) in at least 2 of their first 3 NFL seasons.
func meetsPlaytimeCriteria(p structs.NFLPlayer, currentSeasonID int, histSnaps map[uint]map[uint]structs.NFLPlayerSeasonSnaps) bool {
	if p.Experience == 0 {
		return false
	}
	firstSeasonID := uint(currentSeasonID) - p.Experience

	qualifyingSeasons := 0
	for offset := uint(0); offset < 3; offset++ {
		sid := firstSeasonID + offset
		playerSeasons, ok := histSnaps[p.ID]
		if !ok {
			continue
		}
		snaps, ok := playerSeasons[sid]
		if !ok {
			continue
		}
		if snaps.GetTotalSnaps() >= starterSnapThreshold {
			qualifyingSeasons++
		}
	}
	return qualifyingSeasons >= 2
}

// =============================================================================
// ResetNFLPlayerMinimumValues
// =============================================================================

// ResetNFLPlayerMinimumValues resets every player's working MinimumValue and
// AAV back to their stored originals. Called around week 15 of the regular
// season so that players enter the offseason with fresh, freshly-calculated
// baselines rather than values degraded by DecreaseMinimumValue calls.
func ResetNFLPlayerMinimumValues() {
	db := dbprovider.GetInstance().GetDB()

	nflPlayers := GetAllNFLPlayers()

	reset := 0
	for _, p := range nflPlayers {
		if p.MinimumValue == p.OriginalMinimumValue && p.AAV == p.OriginalAAV {
			continue // already at baseline, skip the write
		}
		p.ResetMinimumAndAAVValues()
		if err := db.Model(&p).Updates(map[string]interface{}{
			"minimum_value": p.MinimumValue,
			"aav":           p.AAV,
		}).Error; err != nil {
			log.Printf("ValuationManager: failed to reset player %d: %v", p.ID, err)
			continue
		}
		reset++
	}

	log.Printf("ResetNFLPlayerMinimumValues: reset %d players", reset)
}

// =============================================================================
// SeedNFLTagDataFromJSON
// =============================================================================

// SeedNFLTagDataFromJSON populates the NFLTagData table from the legacy
// tagData.json file. Call this once after running the AutoMigrate to create
// the table; subsequent updates should go through CalculateTagValues().
//
// NOTE: the JSON stores data by individual position (DE, OLB, ILB, DT, …).
// For grouped positions (EDGE, DL, LB) the seeding logic takes the arithmetic
// average of the constituent positions as a reasonable starting point.
func SeedNFLTagDataFromJSON() {
	db := dbprovider.GetInstance().GetDB()

	raw := util.GetTagData()

	// Positions that map 1-to-1 with groups
	directGroups := []string{"QB", "RB", "FB", "WR", "TE", "OT", "OG", "C", "CB", "FS", "SS", "K", "P"}
	for _, pos := range directGroups {
		vals, ok := raw[pos]
		if !ok {
			continue
		}
		td := repository.FindNFLTagDataByPosition(pos)
		td.Position = pos
		td.Franchise = vals["Franchise"]
		td.Transition = vals["Transition"]
		td.Playtime = vals["Playtime"]
		td.Basic = vals["Basic"]
		repository.SaveNFLTagData(td, db)
	}

	// EDGE = average of DE and OLB entries from JSON
	edge := seedNFLTagGroup(raw, "EDGE", []string{"DE", "OLB"})
	existing := repository.FindNFLTagDataByPosition("EDGE")
	edge.Model = existing.Model // preserve existing DB row if present
	repository.SaveNFLTagData(edge, db)

	// DL = average of DT and DE entries
	dl := seedNFLTagGroup(raw, "DL", []string{"DT", "DE"})
	existingDL := repository.FindNFLTagDataByPosition("DL")
	dl.Model = existingDL.Model
	repository.SaveNFLTagData(dl, db)

	// LB = average of ILB and OLB entries
	lb := seedNFLTagGroup(raw, "LB", []string{"ILB", "OLB"})
	existingLB := repository.FindNFLTagDataByPosition("LB")
	lb.Model = existingLB.Model
	repository.SaveNFLTagData(lb, db)

	log.Printf("SeedNFLTagDataFromJSON: seeding complete")
}

// seedNFLTagGroup builds an NFLTagData record for a merged group by averaging
// the constituent positions from the raw JSON map.
func seedNFLTagGroup(raw map[string]map[string]float64, group string, positions []string) structs.NFLTagData {
	var fran, trans, play, basic float64
	count := 0
	for _, pos := range positions {
		vals, ok := raw[pos]
		if !ok {
			continue
		}
		fran += vals["Franchise"]
		trans += vals["Transition"]
		play += vals["Playtime"]
		basic += vals["Basic"]
		count++
	}
	if count == 0 {
		return structs.NFLTagData{Position: group}
	}
	n := float64(count)
	return structs.NFLTagData{
		Position:   group,
		Franchise:  fran / n,
		Transition: trans / n,
		Playtime:   play / n,
		Basic:      basic / n,
	}
}

// GetTagAmountsForGroup returns the tag dollar amounts for a given position
// group (e.g., "EDGE", "QB"). Prefers DB data; falls back to JSON.
func GetTagAmountsForGroup(group string) structs.NFLTagData {
	td := repository.FindNFLTagDataByPosition(group)
	if td.ID != 0 {
		return td
	}
	// If not in DB yet, build from JSON
	raw := util.GetTagData()
	if vals, ok := raw[group]; ok {
		return structs.NFLTagData{
			Position:   group,
			Franchise:  vals["Franchise"],
			Transition: vals["Transition"],
			Playtime:   vals["Playtime"],
			Basic:      vals["Basic"],
		}
	}
	return structs.NFLTagData{Position: group}
}
