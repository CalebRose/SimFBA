package structs

type BasePlayerGameSnaps struct {
	ID       uint
	SeasonID uint
	PlayerID uint
	GameID   uint
	WeekID   uint
	QBSnaps  uint16
	RBSnaps  uint16
	FBSnaps  uint16
	WRSnaps  uint16
	TESnaps  uint16
	OTSnaps  uint16
	OGSnaps  uint16
	CSnaps   uint16
	DTSnaps  uint16
	DESnaps  uint16
	OLBSnaps uint16
	ILBSnaps uint16
	CBSnaps  uint16
	FSSnaps  uint16
	SSSnaps  uint16
	KSnaps   uint16
	PSnaps   uint16
	STSnaps  uint16
	KRSnaps  uint16
	PRSnaps  uint16
	KOSSnaps uint16
}

func (g *BasePlayerGameSnaps) MapSnapsToPosition(pos string, snaps int) {
	switch pos {
	case "QB":
		g.QBSnaps += uint16(snaps)
	case "RB":
		g.RBSnaps += uint16(snaps)
	case "FB":
		g.FBSnaps += uint16(snaps)
	case "WR":
		g.WRSnaps += uint16(snaps)
	case "TE":
		g.TESnaps += uint16(snaps)
	case "OT":
		g.OTSnaps += uint16(snaps)
	case "OG":
		g.OGSnaps += uint16(snaps)
	case "C":
		g.CSnaps += uint16(snaps)
	case "DE":
		g.DESnaps += uint16(snaps)
	case "DT":
		g.DTSnaps += uint16(snaps)
	case "OLB":
		g.OLBSnaps += uint16(snaps)
	case "ILB":
		g.ILBSnaps += uint16(snaps)
	case "CB":
		g.CBSnaps += uint16(snaps)
	case "FS":
		g.FSSnaps += uint16(snaps)
	case "SS":
		g.SSSnaps += uint16(snaps)
	case "P":
		g.PSnaps += uint16(snaps)
	case "K":
		g.KSnaps += uint16(snaps)
	case "ST":
		g.STSnaps += uint16(snaps)
	case "KR":
		g.KRSnaps += uint16(snaps)
	case "PR":
		g.PRSnaps += uint16(snaps)
	case "KOS":
		g.KOSSnaps += uint16(snaps)
	}
}

type BasePlayerSeasonSnaps struct {
	ID       uint
	SeasonID uint
	PlayerID uint
	QBSnaps  uint16
	RBSnaps  uint16
	FBSnaps  uint16
	WRSnaps  uint16
	TESnaps  uint16
	OTSnaps  uint16
	OGSnaps  uint16
	CSnaps   uint16
	DTSnaps  uint16
	DESnaps  uint16
	OLBSnaps uint16
	ILBSnaps uint16
	CBSnaps  uint16
	FSSnaps  uint16
	SSSnaps  uint16
	KSnaps   uint16
	PSnaps   uint16
	STSnaps  uint16
	KRSnaps  uint16
	PRSnaps  uint16
	KOSSnaps uint16
}

func (s *BasePlayerSeasonSnaps) Reset() {
	s.QBSnaps = 0
	s.RBSnaps = 0
	s.FBSnaps = 0
	s.WRSnaps = 0
	s.TESnaps = 0
	s.OTSnaps = 0
	s.OGSnaps = 0
	s.CSnaps = 0
	s.DTSnaps = 0
	s.DESnaps = 0
	s.OLBSnaps = 0
	s.ILBSnaps = 0
	s.CBSnaps = 0
	s.FSSnaps = 0
	s.SSSnaps = 0
	s.PSnaps = 0
	s.KSnaps = 0
	s.STSnaps = 0
	s.PRSnaps = 0
	s.KRSnaps = 0
	s.KOSSnaps = 0
}

func (s *BasePlayerSeasonSnaps) AddToSeason(g BasePlayerGameSnaps) {
	s.QBSnaps += g.QBSnaps
	s.RBSnaps += g.RBSnaps
	s.FBSnaps += g.FBSnaps
	s.WRSnaps += g.WRSnaps
	s.TESnaps += g.TESnaps
	s.OTSnaps += g.OTSnaps
	s.OGSnaps += g.OGSnaps
	s.CSnaps += g.CSnaps
	s.DTSnaps += g.DTSnaps
	s.DESnaps += g.DESnaps
	s.OLBSnaps += g.OLBSnaps
	s.ILBSnaps += g.ILBSnaps
	s.CBSnaps += g.CBSnaps
	s.FSSnaps += g.FSSnaps
	s.SSSnaps += g.SSSnaps
	s.PSnaps += g.PSnaps
	s.KSnaps += g.KSnaps
	s.STSnaps += g.STSnaps
	s.PRSnaps += g.PRSnaps
	s.KRSnaps += g.KRSnaps
	s.KOSSnaps += g.KOSSnaps
}

func (s *BasePlayerSeasonSnaps) GetTotalSnaps() int {
	return int(s.QBSnaps + s.RBSnaps + s.FBSnaps + s.WRSnaps +
		s.TESnaps + s.OTSnaps + s.OGSnaps + s.CSnaps + s.DTSnaps +
		s.DESnaps + s.OLBSnaps + s.ILBSnaps + s.CBSnaps + s.FSSnaps +
		s.SSSnaps + s.PSnaps + s.KSnaps + s.STSnaps + s.PRSnaps +
		s.KRSnaps + s.KOSSnaps)
}

type CollegePlayerSeasonSnaps struct {
	BasePlayerSeasonSnaps
}

type NFLPlayerSeasonSnaps struct {
	BasePlayerSeasonSnaps
}

type CollegePlayerGameSnaps struct {
	BasePlayerGameSnaps
}

type NFLPlayerGameSnaps struct {
	BasePlayerGameSnaps
}
