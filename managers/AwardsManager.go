package managers

import (
	"sort"
	"strconv"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
	"gorm.io/gorm"
)

type AwardsModel struct {
	TeamID    int
	FirstName string
	LastName  string
	Position  string
	Archetype string
	School    string
	Score     float64
	Games     int
}

// Sorting Funcs
type ByScore []AwardsModel

func (h ByScore) Len() int      { return len(h) }
func (h ByScore) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h ByScore) Less(i, j int) bool {
	return h[i].Score > h[j].Score
}

type AwardsList struct {
	HeismanList        []AwardsModel // Best Player Overall
	CoachOfTheYearList []AwardsModel // Best Coach
	DaveyOBrienList    []AwardsModel // Best QB
	DoakWalkerList     []AwardsModel // Best RB
	BiletnikoffList    []AwardsModel // Best WR
	MackeyList         []AwardsModel // Best TE
	RimingtonList      []AwardsModel // Best C
	OutlandList        []AwardsModel // Best Interior Lineman
	JoeMooreList       []AwardsModel // Best Offensive Line
	NagurskiList       []AwardsModel // Best Defensive Player
	HendricksList      []AwardsModel // Best Defensive Lineman
	ThorpeList         []AwardsModel // Best DB
	ButkusList         []AwardsModel // Best Linebacker
	LouGrozaList       []AwardsModel // Best Kicker
	RayGuyList         []AwardsModel // Best Punter
	JetAward           []AwardsModel // Best Return Specialist
}

// AwardConfig defines the configuration for each award
type AwardConfig struct {
	Name           string
	PositionFilter []string                                     // positions eligible for this award (nil = all positions)
	ScoreModifier  func(score float64, position string) float64 // custom scoring modifier
	MaxEntries     int                                          // max entries in final list
}

// ScoreCalculator handles the score calculation logic
type ScoreCalculator struct {
	weightMap      map[string]float64
	homeTeamMapper map[int]string
}

// Standard award configurations
var awardConfigs = map[string]AwardConfig{
	"Heisman": {
		Name:           "Heisman Trophy",
		PositionFilter: nil, // All positions eligible
		ScoreModifier:  nil, // Use standard scoring
		MaxEntries:     25,
	},
	"DaveyOBrien": {
		Name:           "Davey O'Brien Award",
		PositionFilter: []string{"QB"},
		ScoreModifier:  nil,
		MaxEntries:     25,
	},
	"DoakWalker": {
		Name:           "Doak Walker Award",
		PositionFilter: []string{"RB"},
		ScoreModifier:  nil,
		MaxEntries:     25,
	},
	"Biletnikoff": {
		Name:           "Biletnikoff Award",
		PositionFilter: []string{"WR"},
		ScoreModifier:  nil,
		MaxEntries:     25,
	},
	"Mackey": {
		Name:           "John Mackey Award",
		PositionFilter: []string{"TE"},
		ScoreModifier:  nil,
		MaxEntries:     25,
	},
	"Rimington": {
		Name:           "Rimington Trophy",
		PositionFilter: []string{"C"},
		ScoreModifier:  enhanceOffensiveLineScoring,
		MaxEntries:     25,
	},
	"Outland": {
		Name:           "Outland Trophy",
		PositionFilter: []string{"OT", "OG", "DT", "DE"},
		ScoreModifier:  enhanceLinemanScoring,
		MaxEntries:     25,
	},
	"JoeMoore": {
		Name:           "Joe Moore Award",
		PositionFilter: []string{"OT", "OG", "C"},
		ScoreModifier:  enhanceOffensiveLineScoring,
		MaxEntries:     25,
	},
	"Nagurski": {
		Name:           "Bronko Nagurski Trophy",
		PositionFilter: []string{"DT", "DE", "ILB", "OLB", "CB", "FS", "SS"},
		ScoreModifier:  enhanceDefensiveScoring,
		MaxEntries:     25,
	},
	"Hendricks": {
		Name:           "Ted Hendricks Award",
		PositionFilter: []string{"DE", "DT"},
		ScoreModifier:  enhanceDefensiveLineScoring,
		MaxEntries:     25,
	},
	"Thorpe": {
		Name:           "Jim Thorpe Award",
		PositionFilter: []string{"CB", "FS", "SS"},
		ScoreModifier:  enhanceSecondaryScoring,
		MaxEntries:     25,
	},
	"Butkus": {
		Name:           "Dick Butkus Award",
		PositionFilter: []string{"ILB", "OLB"},
		ScoreModifier:  enhanceLinebackerScoring,
		MaxEntries:     25,
	},
	"LouGroza": {
		Name:           "Lou Groza Award",
		PositionFilter: []string{"K"},
		ScoreModifier:  enhanceKickerScoring,
		MaxEntries:     25,
	},
	"RayGuy": {
		Name:           "Ray Guy Award",
		PositionFilter: []string{"P"},
		ScoreModifier:  enhancePunterScoring,
		MaxEntries:     25,
	},
	"Jet": {
		Name:           "Jet Award",
		PositionFilter: nil, // Return specialist can be any position
		ScoreModifier:  enhanceReturnSpecialistScoring,
		MaxEntries:     25,
	},
}

// NewScoreCalculator creates a new score calculator with the provided mappings
func NewScoreCalculator(weightMap map[string]float64, homeTeamMapper map[int]string) *ScoreCalculator {
	return &ScoreCalculator{
		weightMap:      weightMap,
		homeTeamMapper: homeTeamMapper,
	}
}

// CalculatePlayerScore calculates the score for a player using unified scoring logic
func (sc *ScoreCalculator) CalculatePlayerScore(cp structs.CollegePlayer, games []structs.CollegeGame, modifier func(score float64, position string) float64) float64 {
	var score float64 = 0
	homeTeam := sc.homeTeamMapper[cp.TeamID]
	stats := cp.Stats
	var totalMod float64 = 0

	for idx, stat := range stats {
		if idx > 12 {
			continue
		}
		var statScore float64 = 0

		opposingTeam := stat.OpposingTeam

		// Unified scoring logic - extract common calculations
		statScore += sc.calculatePassingScore(stat)
		statScore += sc.calculateRushingScore(stat, cp.Position)
		statScore += sc.calculateReceivingScore(stat, cp.Position)
		statScore += sc.calculateDefensiveScore(stat)
		statScore += sc.calculateSpecialTeamsScore(stat)

		// Apply game context
		game := GetCollegeGameStructByGameID(games, stat.GameID)
		opposingTeamWeight := sc.weightMap[opposingTeam]

		if (game.HomeTeamWin && cp.TeamID != game.HomeTeamID) || (game.AwayTeamWin && cp.TeamID != game.AwayTeamID) {
			opposingTeamWeight *= -0.4125
		}

		totalMod += opposingTeamWeight
		score += statScore
	}

	// Normalize score
	score = (score / float64(len(stats))) * (float64(len(stats)) / 12)
	score = score * (1 + totalMod + sc.weightMap[homeTeam])

	// Apply position-specific modifier if provided
	if modifier != nil {
		score = modifier(score, cp.Position)
	}

	return score
}

// Individual scoring components
func (sc *ScoreCalculator) calculatePassingScore(stat structs.CollegePlayerStats) float64 {
	score := float64(stat.PassingYards) * 0.068
	score += float64(stat.PassingTDs) * 4
	score -= float64(stat.Interceptions) * 2.25
	score -= float64(stat.Sacks) * 2.25
	return score
}

func (sc *ScoreCalculator) calculateRushingScore(stat structs.CollegePlayerStats, position string) float64 {
	var score float64
	if position == "RB" || position == "FB" {
		score += float64(stat.RushingYards) * 0.1125
		score += float64(stat.RushingTDs) * 6
	} else {
		score += float64(stat.RushingYards) * 0.1
		score += float64(stat.RushingTDs) * 5
	}
	score -= float64(stat.Fumbles) * 5.75
	return score
}

func (sc *ScoreCalculator) calculateReceivingScore(stat structs.CollegePlayerStats, position string) float64 {
	drops := float64(stat.Targets - stat.Catches)
	var score float64

	if position == "WR" || position == "TE" {
		score += float64(stat.Catches) * 0.525
		score += float64(stat.ReceivingYards) * 0.1125
		score += float64(stat.ReceivingTDs) * 6
	} else {
		score += float64(stat.Catches) * 0.25
		score += float64(stat.ReceivingYards) * 0.05
		score += float64(stat.ReceivingTDs) * 4
	}
	score -= drops * 0.75
	return score
}

func (sc *ScoreCalculator) calculateDefensiveScore(stat structs.CollegePlayerStats) float64 {
	score := float64(stat.SoloTackles) * 1
	score += float64(stat.STSoloTackles) * 1
	score += float64(stat.AssistedTackles) * 0.9
	score += float64(stat.STAssistedTackles) * 0.9
	score += float64(stat.TacklesForLoss) * 6.25
	score += float64(stat.SacksMade) * 7.125
	score += float64(stat.PassDeflections) * 6.25
	score += float64(stat.ForcedFumbles) * 8
	score += float64(stat.RecoveredFumbles) * 6
	score += float64(stat.InterceptionsCaught) * 15
	score += float64(stat.PuntsBlocked) * 10
	score += float64(stat.Safeties) * 10
	score += float64(stat.DefensiveTDs) * 20
	return score
}

func (sc *ScoreCalculator) calculateSpecialTeamsScore(stat structs.CollegePlayerStats) float64 {
	score := float64(stat.KickReturnTDs) * 14
	score += float64(stat.PuntReturnTDs) * 14
	score += float64(stat.FGBlocked) * 12
	score += float64(stat.FGMade) * 0.3
	score += float64(stat.ExtraPointsMade) * 0.1
	return score
}

// Position-specific scoring modifiers
func enhanceDefensiveScoring(score float64, position string) float64 {
	return score * 1.2 // Boost defensive stats
}

func enhanceDefensiveLineScoring(score float64, position string) float64 {
	return score * 1.15 // Slight boost for defensive linemen
}

func enhanceSecondaryScoring(score float64, position string) float64 {
	return score * 1.1 // Boost for secondary players
}

func enhanceLinebackerScoring(score float64, position string) float64 {
	return score * 1.1 // Boost for linebackers
}

func enhanceOffensiveLineScoring(score float64, position string) float64 {
	return score * 0.8 // O-line typically gets less statistical credit
}

func enhanceLinemanScoring(score float64, position string) float64 {
	return score * 1.05 // General lineman boost
}

func enhanceKickerScoring(score float64, position string) float64 {
	return score * 3.0 // Heavily weight kicking stats
}

func enhancePunterScoring(score float64, position string) float64 {
	return score * 2.5 // Heavily weight punting stats
}

func enhanceReturnSpecialistScoring(score float64, position string) float64 {
	return score * 1.8 // Boost return specialist stats
}

// GenerateAwardsList creates an awards list based on the provided configuration
func GenerateAwardsList(awardKey string, collegePlayers []structs.CollegePlayer, calculator *ScoreCalculator, teamGameMapper map[int][]structs.CollegeGame) []AwardsModel {
	config, exists := awardConfigs[awardKey]
	if !exists {
		return []AwardsModel{}
	}

	var candidates []AwardsModel
	teamCountMap := make(map[int]int)

	for _, cp := range collegePlayers {
		// Filter by position if required
		if config.PositionFilter != nil && !contains(config.PositionFilter, cp.Position) {
			continue
		}

		if len(cp.Stats) == 0 {
			continue
		}

		score := calculator.CalculatePlayerScore(cp, teamGameMapper[cp.TeamID], config.ScoreModifier)

		candidate := AwardsModel{
			TeamID:    cp.TeamID,
			FirstName: cp.FirstName,
			LastName:  cp.LastName,
			Position:  cp.Position,
			Archetype: cp.Archetype,
			School:    cp.TeamAbbr,
			Score:     score,
			Games:     len(cp.Stats),
		}

		candidates = append(candidates, candidate)
	}

	// Sort by score
	sort.Sort(ByScore(candidates))

	// Limit to max entries and enforce team count limits
	officialList := []AwardsModel{}
	count := 0

	for _, candidate := range candidates {
		if count == config.MaxEntries {
			break
		}
		teamCount := teamCountMap[candidate.TeamID]
		if teamCount > 1 { // Max 2 players per team per award
			continue
		}
		count++
		teamCountMap[candidate.TeamID]++
		officialList = append(officialList, candidate)
	}

	return officialList
}

// GetAllPostSeasonAwardsLists generates all award lists with the refactored approach
func GetAllPostSeasonAwardsLists() AwardsList {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()
	seasonID := strconv.Itoa(ts.CollegeSeasonID)
	collegeTeams := GetAllCollegeTeams()
	collegeStandings := GetAllCollegeStandingsBySeasonID(seasonID)
	cfbStandingsMap := MakeCollegeStandingsMapByTeamID(collegeStandings)
	collegeGames := repository.FindCollegeGamesRecords(seasonID, false)
	cfbGameMap := MakeCollegeGameMapByTeamID(collegeGames)

	var collegePlayers []structs.CollegePlayer
	var teamWeight = make(map[string]float64)
	var homeTeamMapper = make(map[int]string)
	var teamGameMapper = make(map[int][]structs.CollegeGame)

	// Initialize team data (same as original)
	for _, team := range collegeTeams {
		homeTeamMapper[int(team.ID)] = team.TeamAbbr
		games := cfbGameMap[team.ID]
		currentYearStandings := cfbStandingsMap[team.ID]

		if len(games) == 0 || currentYearStandings.ID == 0 {
			continue
		}

		teamGameMapper[int(team.ID)] = games

		var weight float64 = 0
		if currentYearStandings.TotalLosses+currentYearStandings.TotalWins > 0 {
			weight = float64(currentYearStandings.TotalWins) / 100
		}

		teamWeight[team.TeamAbbr] = weight
	}

	// Get player data (same as original)
	var distinctCollegeStats []structs.CollegePlayerStats
	db.Distinct("college_player_id").Where("snaps > 0").Find(&distinctCollegeStats)
	distinctCollegePlayerIDs := GetCollegePlayerIDs(distinctCollegeStats)

	db.Preload("Stats", func(db *gorm.DB) *gorm.DB {
		return db.Where("snaps > 0 and season_id = ? and week_id < ? and game_type = ?", seasonID, strconv.Itoa(ts.CollegeWeekID), "2")
	}).Where("id IN ?", distinctCollegePlayerIDs).Find(&collegePlayers)

	// Create score calculator
	calculator := NewScoreCalculator(teamWeight, homeTeamMapper)

	// Generate all award lists using the unified approach
	return AwardsList{
		HeismanList:        GenerateAwardsList("Heisman", collegePlayers, calculator, teamGameMapper),
		CoachOfTheYearList: []AwardsModel{}, // This requires different logic
		DaveyOBrienList:    GenerateAwardsList("DaveyOBrien", collegePlayers, calculator, teamGameMapper),
		DoakWalkerList:     GenerateAwardsList("DoakWalker", collegePlayers, calculator, teamGameMapper),
		BiletnikoffList:    GenerateAwardsList("Biletnikoff", collegePlayers, calculator, teamGameMapper),
		MackeyList:         GenerateAwardsList("Mackey", collegePlayers, calculator, teamGameMapper),
		RimingtonList:      GenerateAwardsList("Rimington", collegePlayers, calculator, teamGameMapper),
		OutlandList:        GenerateAwardsList("Outland", collegePlayers, calculator, teamGameMapper),
		JoeMooreList:       GenerateAwardsList("JoeMoore", collegePlayers, calculator, teamGameMapper),
		NagurskiList:       GenerateAwardsList("Nagurski", collegePlayers, calculator, teamGameMapper),
		HendricksList:      GenerateAwardsList("Hendricks", collegePlayers, calculator, teamGameMapper),
		ThorpeList:         GenerateAwardsList("Thorpe", collegePlayers, calculator, teamGameMapper),
		ButkusList:         GenerateAwardsList("Butkus", collegePlayers, calculator, teamGameMapper),
		LouGrozaList:       GenerateAwardsList("LouGroza", collegePlayers, calculator, teamGameMapper),
		RayGuyList:         GenerateAwardsList("RayGuy", collegePlayers, calculator, teamGameMapper),
		JetAward:           GenerateAwardsList("Jet", collegePlayers, calculator, teamGameMapper),
	}
}

// Helper function to check if slice contains value
func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
