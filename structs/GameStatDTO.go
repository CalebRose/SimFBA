package structs

import "encoding/json"

type GameStatDTO struct {
	GameID            uint
	HomeTeam          TeamStatDTO
	AwayTeam          TeamStatDTO
	HomePlayers       []PlayerStatDTO
	AwayPlayers       []PlayerStatDTO
	HomeScore         int
	AwayScore         int
	Plays             []PlayByPlayDTO
	PlayerSnapTracker PlayerSnapTracker
	HomeGameplan      json.RawMessage
	AwayGameplan      json.RawMessage
}

type PlayerSnapTracker struct {
	PlayerSnapCounts map[int]map[string]int
}
