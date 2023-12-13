package models

import "github.com/CalebRose/SimFBA/structs"

type RosterPageResponse struct {
	Teams   []structs.CollegeTeam
	Coaches []structs.CollegeCoach
}
