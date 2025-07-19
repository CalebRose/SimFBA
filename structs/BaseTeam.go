package structs

import "time"

type JerseyDTO struct {
	TeamID     string
	JerseyType uint8
}

type BaseTeam struct {
	TeamName          string
	Mascot            string
	TeamAbbr          string
	Coach             string
	City              string
	State             string
	Country           string
	StadiumID         uint
	Stadium           string
	StadiumCapacity   int
	RecordAttendance  int
	Enrollment        int
	FirstPlayed       int
	ColorOne          string
	ColorTwo          string
	ColorThree        string
	DiscordID         string
	OverallGrade      string
	OffenseGrade      string
	DefenseGrade      string
	SpecialTeamsGrade string
	PenaltyMarks      uint8
	JerseyType        uint8
	LastLogin         time.Time
}

func (bt *BaseTeam) UpdateLatestInstance() {
	bt.LastLogin = time.Now()
}

func (bt *BaseTeam) RemoveUserFromTeam() {
	bt.Coach = "AI"
}

func (bt *BaseTeam) AssignUserToTeam(user string) {
	bt.Coach = user
}

func (bt *BaseTeam) AssignDiscordID(id string) {
	bt.DiscordID = id
}

func (bt *BaseTeam) MarkTeamForPenalty() {
	bt.PenaltyMarks += 1
}

func (bt *BaseTeam) ResetMarks() {
	bt.PenaltyMarks = 0
}

func (bt *BaseTeam) AdjustJerseyType(jersey uint8) {
	bt.JerseyType = jersey
}

func (bt *BaseTeam) AssignTeamGrades(ovr, off, def, st string) {
	bt.OverallGrade = ovr
	bt.OffenseGrade = off
	bt.DefenseGrade = def
	bt.SpecialTeamsGrade = st
}
