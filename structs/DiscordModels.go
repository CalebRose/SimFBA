package structs

type CFBComparisonModel struct {
	TeamOneID      uint
	TeamOne        string
	TeamOneWins    uint
	TeamOneLosses  uint
	TeamOneStreak  uint
	TeamOneMSeason uint
	TeamOneMScore  string
	TeamTwoID      uint
	TeamTwo        string
	TeamTwoWins    uint
	TeamTwoLosses  uint
	TeamTwoStreak  uint
	TeamTwoMSeason uint
	TeamTwoMScore  string
	CurrentStreak  uint
	LatestWin      string
}
