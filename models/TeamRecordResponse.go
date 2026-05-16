package models

import (
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
)

type TeamRecordResponse struct {
	OverallWins             int
	OverallLosses           int
	CurrentSeasonWins       int
	CurrentSeasonLosses     int
	BowlWins                int
	BowlLosses              int
	ConferenceChampionships []string
	DivisionTitles          []string
	NationalChampionships   []string
	TopPlayers              []TopPlayer
}

func (t *TeamRecordResponse) AddTopPlayers(players []TopPlayer) {
	t.TopPlayers = players
}

type TopPlayer struct {
	PlayerID     uint
	FirstName    string
	LastName     string
	Position     string
	Archetype    string
	PositionTwo  string
	ArchetypeTwo string
	OverallGrade string
	Overall      int
}

func (t *TopPlayer) MapCollegePlayer(player structs.CollegePlayer) {
	t.PlayerID = player.ID
	t.FirstName = player.FirstName
	t.LastName = player.LastName
	t.Position = player.Position
	t.PositionTwo = player.PositionTwo
	t.Archetype = player.Archetype
	t.ArchetypeTwo = player.ArchetypeTwo
	t.Overall = int(player.Overall)
	t.OverallGrade = util.GetOverallGrade(int(player.Overall), player.Year)
}

func (t *TopPlayer) MapNFLPlayer(player structs.NFLPlayer) {
	t.PlayerID = player.ID
	t.FirstName = player.FirstName
	t.LastName = player.LastName
	t.Position = player.Position
	t.PositionTwo = player.PositionTwo
	t.Archetype = player.Archetype
	t.ArchetypeTwo = player.ArchetypeTwo
	t.Overall = int(player.Overall)
	t.OverallGrade = util.GetNFLOverallGrade(int(player.Overall))
}
