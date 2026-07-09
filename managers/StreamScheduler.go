package managers

import (
	"context"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"
	time "time"

	fbsvc "github.com/CalebRose/SimFBA/firebase"
	"github.com/CalebRose/SimFBA/repository"
	"github.com/CalebRose/SimFBA/structs"
)

const maxStreamSlots = 8

// cron guards — prevent duplicate streaming goroutines per league.
var (
	cfbCronMu     sync.Mutex
	cfbCronCancel context.CancelFunc
	nflCronMu     sync.Mutex
	nflCronCancel context.CancelFunc
)

// PendingGame is a lightweight descriptor for a game waiting to enter a slot.
type PendingGame struct {
	GameID        uint
	HomeTeamID    uint
	AwayTeamID    uint
	HomeTeam      string
	AwayTeam      string
	IsUserGame    bool // true if either team is user-coached / user-owned
	HomeTeamRank  int
	AwayTeamRank  int
	HomeTeamCoach string
	AwayTeamCoach string
	Arena         string
	City          string
	State         string
	Country       string
	TimeSlot      string
}

// GameStream represents one active streaming slot.
type GameStream struct {
	GameID    uint
	StartTime time.Time
	EndTime   time.Time
	League    string
}

// StreamScheduler manages up to maxStreamSlots concurrent game streams and an
// ordered queue of pending games for a single league.
type StreamScheduler struct {
	mu               sync.Mutex
	ActiveSlots      [maxStreamSlots]*GameStream
	Queue            []PendingGame
	League           string // "chl" or "phl"
	isCollege        bool
	CFBPlayByPlayMap map[uint][]structs.CollegePlayByPlay
	NFLPlayByPlayMap map[uint][]structs.NFLPlayByPlay
}

// parseMMSS converts a clock string in "MM:SS" format to the equivalent number
// of seconds.  Returns 0 for any malformed input.
func parseMMSS(s string) int {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return 0
	}
	mins, err1 := strconv.Atoi(parts[0])
	secs, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil || mins < 0 || secs < 0 || secs >= 60 {
		return 0
	}
	return mins*60 + secs
}

// computeStreamTimes sums SecondsConsumed across a play-by-play slice and
// returns a start time of now, the corresponding end time, and the total seconds.
func computeStreamTimes(totalSecs int) (start, end time.Time) {
	start = time.Now().UTC()
	end = start.Add(time.Duration(totalSecs) * time.Second)
	return
}

// loadTotalSeconds calculates the total game-clock seconds that elapsed during
// a game by inspecting the last play-by-play record.
//
// Formula: (quartersPlayed × 900) − timeRemainingInLastQuarter
//
// A regulation game (4 quarters, ending at 0:00) yields 3600 s.  Overtime
// quarters each add another 900 s; a walk-off ending mid-quarter subtracts the
// unused clock via parseMMSS.  Returns 0 if no records are found.
func loadTotalSeconds(gameID uint, cfbPlayByPlayMap map[uint][]structs.CollegePlayByPlay, nflPlayByPlayMap map[uint][]structs.NFLPlayByPlay, isCollege bool) int {
	const quarterSecs = 15 * 60 // 900 seconds per quarter

	if isCollege {
		plays := cfbPlayByPlayMap[gameID]
		if len(plays) == 0 {
			return 0
		}
		last := plays[len(plays)-1]
		timeRemainingInLastQuarter := parseMMSS(last.TimeRemaining)
		return int(last.Quarter)*quarterSecs - timeRemainingInLastQuarter
	}

	plays := nflPlayByPlayMap[gameID]
	if len(plays) == 0 {
		return 0
	}
	last := plays[len(plays)-1]
	timeRemainingInLastQuarter := parseMMSS(last.TimeRemaining)
	return int(last.Quarter)*quarterSecs - timeRemainingInLastQuarter
}

// loadTotalPlays returns the number of play-by-play records for a game.
func loadTotalPlays(gameID uint, cfbPlayByPlayMap map[uint][]structs.CollegePlayByPlay, nflPlayByPlayMap map[uint][]structs.NFLPlayByPlay, isCollege bool) int {
	if isCollege {
		return len(cfbPlayByPlayMap[gameID])
	}
	return len(nflPlayByPlayMap[gameID])
}

// dequeue pops the first item from the queue.
func (s *StreamScheduler) dequeue() (PendingGame, bool) {
	if len(s.Queue) == 0 {
		return PendingGame{}, false
	}
	next := s.Queue[0]
	s.Queue = s.Queue[1:]
	return next, true
}

// InitQueue loads all complete, unrevealed games for the current matchup and
// sorts them with user-coached/owned games first, then by GameID.
func (s *StreamScheduler) InitQueue(weekID, seasonID, gameDay string, isPreseason bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var userGames, aiGames []PendingGame
	cfbPlayByPlayMap := make(map[uint][]structs.CollegePlayByPlay)
	nflPlayByPlayMap := make(map[uint][]structs.NFLPlayByPlay)

	if s.isCollege {
		games := GetCollegeGamesByWeekIdAndSeasonID(weekID, seasonID, isPreseason)
		gameIDs := make([]string, len(games))
		for i, g := range games {
			gameIDs[i] = strconv.Itoa(int(g.ID))
		}
		playByPlays := GetCFBPlayByPlaysByGameIDs(gameIDs)
		for _, play := range playByPlays {
			cfbPlayByPlayMap[uint(play.GameID)] = append(cfbPlayByPlayMap[uint(play.GameID)], play)
		}
		teamMap := GetCollegeTeamMap()
		for _, g := range games {
			if !g.GameComplete || g.IsRevealed {
				continue
			}
			if gameDay == "Thursday Night" && g.TimeSlot != "Thursday Night" {
				continue
			}
			if gameDay == "Friday Night" && g.TimeSlot != "Friday Night" {
				continue
			}
			// Queue up all of the saturday games together
			if gameDay == "Saturday Morning" && (g.TimeSlot == "Thursday Night" || g.TimeSlot == "Friday Night") {
				continue
			}
			homeTeam := teamMap[uint(g.HomeTeamID)]
			awayTeam := teamMap[uint(g.AwayTeamID)]
			pg := PendingGame{
				GameID:       g.ID,
				HomeTeamID:   uint(g.HomeTeamID),
				AwayTeamID:   uint(g.AwayTeamID),
				HomeTeam:     homeTeam.TeamAbbr,
				AwayTeam:     awayTeam.TeamAbbr,
				IsUserGame:   homeTeam.Coach != "AI" || awayTeam.Coach != "AI",
				HomeTeamRank: int(g.HomeTeamRank),
				AwayTeamRank: int(g.AwayTeamRank),
				Arena:        g.Stadium,
				City:         g.City,
				State:        g.State,
				Country:      "",
				TimeSlot:     g.TimeSlot,
			}
			if pg.IsUserGame {
				userGames = append(userGames, pg)
			} else {
				aiGames = append(aiGames, pg)
			}
		}
	} else {
		preseason := "N"
		if isPreseason {
			preseason = "Y"
		}
		games := repository.FindNFLGamesRecords(repository.GamesQuery{
			SeasonID:        seasonID,
			WeekID:          weekID,
			IsPreseasonGame: preseason,
		})
		gameIDs := make([]string, len(games))
		for i, g := range games {
			gameIDs[i] = strconv.Itoa(int(g.ID))
		}
		playByPlays := GetNFLPlayByPlaysByGameIDs(gameIDs)
		for _, play := range playByPlays {
			nflPlayByPlayMap[uint(play.GameID)] = append(nflPlayByPlayMap[uint(play.GameID)], play)
		}
		nflTeams := GetAllNFLTeams()
		nflTeamMap := MakeNFLTeamMap(nflTeams)
		for _, g := range games {
			if !g.GameComplete || g.IsRevealed {
				continue
			}
			if gameDay == "Thursday Night Football" && g.TimeSlot != "Thursday Night Football" {
				continue
			}
			if gameDay == "Monday Night Football" && g.TimeSlot != "Monday Night Football" {
				continue
			}
			homeTeam := nflTeamMap[uint(g.HomeTeamID)]
			awayTeam := nflTeamMap[uint(g.AwayTeamID)]
			isUser := homeTeam.NFLOwnerName != "" || awayTeam.NFLOwnerName != "" ||
				homeTeam.NFLGMName != "" || awayTeam.NFLGMName != ""
			pg := PendingGame{
				GameID:       g.ID,
				HomeTeamID:   uint(g.HomeTeamID),
				AwayTeamID:   uint(g.AwayTeamID),
				HomeTeam:     homeTeam.TeamAbbr,
				AwayTeam:     awayTeam.TeamAbbr,
				IsUserGame:   isUser,
				HomeTeamRank: 0,
				AwayTeamRank: 0,
				Arena:        g.Stadium,
				City:         g.City,
				State:        g.State,
				Country:      "",
				TimeSlot:     g.TimeSlot,
			}
			if pg.IsUserGame {
				userGames = append(userGames, pg)
			} else {
				aiGames = append(aiGames, pg)
			}
		}
	}

	s.CFBPlayByPlayMap = cfbPlayByPlayMap
	s.NFLPlayByPlayMap = nflPlayByPlayMap

	// User-coached games fill the front of the queue; AI games follow.
	s.Queue = append(userGames, aiGames...)
	// sort Queue by timeslot (Thursday Night, Thursday Night Football, Friday Night, Saturday Morning, Saturday Afternoon, Saturday Evening, Saturday Night, Sunday Morning, Sunday Noon, Sunday Night Football, Monday Night Football, all else)
	sort.Slice(s.Queue, func(i, j int) bool {
		timeslotOrder := map[string]int{
			"Thursday Night":          1,
			"Thursday Night Football": 2,
			"Friday Night":            3,
			"Saturday Morning":        4,
			"Saturday Afternoon":      5,
			"Saturday Evening":        6,
			"Saturday Night":          7,
			"Sunday Morning":          8,
			"Sunday Noon":             9,
			"Sunday Night Football":   10,
			"Monday Night Football":   11,
		}
		iOrder, iOk := timeslotOrder[s.Queue[i].TimeSlot]
		jOrder, jOk := timeslotOrder[s.Queue[j].TimeSlot]
		if !iOk {
			iOrder = 12
		}
		if !jOk {
			jOrder = 12
		}
		return iOrder < jOrder
	})

	log.Printf("StreamScheduler(%s): queued %d games (%d user, %d AI)",
		s.League, len(s.Queue), len(userGames), len(aiGames))
}

// Tick is called by the cron on every interval.  It marks completed game slots
// as revealed in Firebase, then promotes pending games into freed slots.
func (s *StreamScheduler) Tick(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now().UTC()

	// 1. Mark completed game slots revealed and free them.
	for i, slot := range s.ActiveSlots {
		if slot == nil {
			continue
		}
		// If the slot has not yet ended, skip it.
		if now.Before(slot.EndTime) {
			continue
		}
		gameID := strconv.Itoa(int(slot.GameID))
		if slot.League == "chl" {
			RevealCFBGameOnInterface(gameID)
		} else {
			RevealNFLGameOnInterface(gameID)
		}
		// Slot has elapsed — mark revealed and clear.
		go func(gameID uint, league string) {
			writeCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			if err := fbsvc.SetGameRevealed(writeCtx, gameID, league); err != nil {
				log.Printf("StreamScheduler: SetGameRevealed(gameID=%d, league=%s): %v", gameID, league, err)
			}
		}(slot.GameID, slot.League)
		s.ActiveSlots[i] = nil
	}

	// 2. Fill freed slots from the queue.
	for i, slot := range s.ActiveSlots {
		if slot != nil || len(s.Queue) == 0 {
			continue
		}
		next, ok := s.dequeue()
		if !ok {
			break
		}

		totalSecs := loadTotalSeconds(next.GameID, s.CFBPlayByPlayMap, s.NFLPlayByPlayMap, s.isCollege)
		if totalSecs == 0 {
			log.Printf("StreamScheduler(%s): skipping game %d — no PbP records found", s.League, next.GameID)
			i--
			continue
		}
		totalPlays := loadTotalPlays(next.GameID, s.CFBPlayByPlayMap, s.NFLPlayByPlayMap, s.isCollege)
		start, end := computeStreamTimes(totalSecs)
		record := fbsvc.LiveGameRecord{
			GameID:          int(next.GameID),
			HomeTeamID:      int(next.HomeTeamID),
			AwayTeamID:      int(next.AwayTeamID),
			HomeTeam:        next.HomeTeam,
			AwayTeam:        next.AwayTeam,
			League:          s.League,
			StreamStartTime: start,
			StreamEndTime:   end,
			TotalPlays:      totalPlays,
			IsRevealed:      false,
			HomeTeamRank:    next.HomeTeamRank,
			AwayTeamRank:    next.AwayTeamRank,
			Arena:           next.Arena,
			City:            next.City,
			State:           next.State,
			Country:         next.Country,
		}
		go func(rec fbsvc.LiveGameRecord, league string) {
			if err := fbsvc.UploadLiveGame(ctx, rec, league); err != nil {
				writeCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				defer cancel()
				if err := fbsvc.UploadLiveGame(writeCtx, rec, league); err != nil {
					log.Printf("StreamScheduler: UploadLiveGame(gameID=%d, league=%s): %v", rec.GameID, league, err)
				}
			}
		}(record, s.League)

		s.ActiveSlots[i] = &GameStream{
			GameID:    next.GameID,
			StartTime: start,
			EndTime:   end,
			League:    s.League,
		}
		log.Printf("StreamScheduler(%s): activated game %d (ends at %s)", s.League, next.GameID, end.Format(time.RFC3339))
	}
}

// IsIdle returns true when all slots are empty and the queue is exhausted.
func (s *StreamScheduler) IsIdle() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.Queue) > 0 {
		return false
	}
	for _, slot := range s.ActiveSlots {
		if slot != nil {
			return false
		}
	}
	return true
}

// StartCFBLiveStreamingCron initialises a CFB StreamScheduler, fills its queue,
// and runs Tick on a 5-second interval until all games are revealed.
// A second call cancels any in-progress cron before starting a new one.
func StartCFBLiveStreamingCron() {
	ts := GetTimestamp()
	if !ts.RunCron || ts.IsOffSeason || ts.CollegeSeasonOver {
		return
	}

	cfbCronMu.Lock()
	if cfbCronCancel != nil {
		cfbCronCancel() // stop previous run
	}
	ctx, cancel := context.WithCancel(context.Background())
	cfbCronCancel = cancel
	cfbCronMu.Unlock()

	if err := fbsvc.PurgeStaleLiveGames(ctx, "cfb"); err != nil {
		log.Printf("RunGames: PurgeStaleLiveGames(cfb): %v", err)
	}

	gameDay := "Thursday Night"
	if ts.ThursdayGames {
		gameDay = "Friday Night"
	}
	if ts.FridayGames {
		gameDay = "Saturday Morning"
	}

	scheduler := &StreamScheduler{League: "cfb", isCollege: true}
	scheduler.InitQueue(
		strconv.Itoa(int(ts.CollegeWeekID)),
		strconv.Itoa(int(ts.CollegeSeasonID)),
		gameDay,
		ts.CFBSpringGames,
	)
	if len(scheduler.Queue) == 0 {
		log.Println("StreamScheduler(cfb): no games to stream")
		cancel()
		return
	}

	scheduler.Tick(ctx) // fill initial slots immediately

	ticker := time.NewTicker(5 * time.Second)
	go func() {
		defer ticker.Stop()
		defer cancel()
		for {
			select {
			case <-ctx.Done():
				log.Println("StreamScheduler(cfb): context cancelled, stopping")
				return
			case <-ticker.C:
				scheduler.Tick(ctx)
				if scheduler.IsIdle() {
					log.Println("StreamScheduler(cfb): all games complete, stopping")
					return
				}
			}
		}
	}()
}

// StartPHLLiveStreamingCron initialises a PHL StreamScheduler and runs it.
// A second call cancels any in-progress cron before starting a new one.
func StartNFLLiveStreamingCron() {
	ts := GetTimestamp()
	if !ts.RunCron || ts.IsOffSeason {
		return
	}

	nflCronMu.Lock()
	if nflCronCancel != nil {
		nflCronCancel()
	}
	ctx, cancel := context.WithCancel(context.Background())
	nflCronCancel = cancel
	nflCronMu.Unlock()

	if err := fbsvc.PurgeStaleLiveGames(ctx, "nfl"); err != nil {
		log.Printf("RunGames: PurgeStaleLiveGames(nfl): %v", err)
	}

	gameDay := "Thursday Night Football"
	if ts.NFLThursday {
		gameDay = "Sunday Noon"
	}
	if ts.NFLSundayNoon {
		gameDay = "Sunday Afternoon"
	}
	if ts.NFLSundayAfternoon {
		gameDay = "Sunday Evening"
	}
	if ts.NFLSundayEvening {
		gameDay = "Monday Night Football"
	}

	scheduler := &StreamScheduler{League: "nfl", isCollege: false}
	scheduler.InitQueue(
		strconv.Itoa(int(ts.NFLWeekID)),
		strconv.Itoa(int(ts.NFLSeasonID)),
		gameDay,
		ts.NFLPreseason,
	)
	if len(scheduler.Queue) == 0 {
		log.Println("StreamScheduler(nfl): no games to stream")
		cancel()
		return
	}

	scheduler.Tick(ctx)

	ticker := time.NewTicker(5 * time.Second)
	go func() {
		defer ticker.Stop()
		defer cancel()
		for {
			select {
			case <-ctx.Done():
				log.Println("StreamScheduler(nfl): context cancelled, stopping")
				return
			case <-ticker.C:
				scheduler.Tick(ctx)
				if scheduler.IsIdle() {
					log.Println("StreamScheduler(nfl): all games complete, stopping")
					return
				}
			}
		}
	}()
}
