package managers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	fbsvc "github.com/CalebRose/SimFBA/firebase"
	"github.com/CalebRose/SimFBA/structs"
)

// ForumManager handles operations related to the forum system within the application.
// Creates threads & posts in Firebase Firestore to facilitate post-game discussions.

// PostGameForumID is the Firestore document ID of the forum category used for
// post-game discussion threads.  Override by changing this constant once the
// forum is set up in Firestore.
const PostGameForumID = "postgame-discussions"

// CreatePostGameDiscussionThreadForCFBGame creates a system-generated postgame
// discussion thread in Firestore for a completed college football game.
// homeTeamStats and awayTeamStats are the per-game box-score records already
// loaded inside SyncTimeslot.
// The operation is idempotent: calling it twice for the same game has no effect.
func CreatePostGameDiscussionThreadForCFBGame(
	game structs.CollegeGame,
	homeTeamStats structs.CollegeTeamStats,
	awayTeamStats structs.CollegeTeamStats,
) {
	ctx := context.Background()

	gameID := strconv.Itoa(int(game.ID))
	eventKey := fmt.Sprintf("postgame_thread:cfb:season%d:game%s", game.SeasonID, gameID)

	title := buildPostGameThreadTitle(game.AwayTeam, game.HomeTeam, game.GameTitle)
	paragraphs := buildCFBPostGameParagraphs(game, homeTeamStats, awayTeamStats)
	bodyText := strings.Join(paragraphs, "\n\n")
	richBody := buildRichPostBody(paragraphs)

	input := fbsvc.CreateForumThreadInput{
		ForumID:           PostGameForumID + "-simcfb",
		ForumPath:         []string{PostGameForumID, "simcfb"},
		Title:             title,
		AuthorUID:         "system",
		AuthorUsername:    "SimSN",
		AuthorDisplayName: "SimSN System",
		CreatedByType:     fbsvc.CreatedBySystem,
		ThreadType:        fbsvc.ThreadTypeGameReference,
		FirstPostBodyText: bodyText,
		FirstPostBody:     richBody,
		ReferencedGameID:  gameID,
		ReferencedLeague:  "cfb",
		ExternalEventKey:  eventKey,
	}

	thread, err := fbsvc.CreateThread(ctx, input)
	if err != nil {
		log.Printf("ForumManager: failed to create CFB postgame thread for game %s: %v", gameID, err)
		return
	}

	log.Printf("ForumManager: created CFB postgame thread %s for game %s (%s)", thread.ID, gameID, title)
}

// CreatePostGameDiscussionThreadForNFLGame creates a system-generated postgame
// discussion thread in Firestore for a completed NFL game.
// homeTeamStats and awayTeamStats are the per-game box-score records already
// loaded inside SyncTimeslot.
// The operation is idempotent: calling it twice for the same game has no effect.
func CreatePostGameDiscussionThreadForNFLGame(
	game structs.NFLGame,
	homeTeamStats structs.NFLTeamStats,
	awayTeamStats structs.NFLTeamStats,
) {
	ctx := context.Background()

	gameID := strconv.Itoa(int(game.ID))
	eventKey := fmt.Sprintf("postgame_thread:nfl:season%d:game%s", game.SeasonID, gameID)

	title := buildPostGameThreadTitle(game.AwayTeam, game.HomeTeam, game.GameTitle)
	paragraphs := buildNFLPostGameParagraphs(game, homeTeamStats, awayTeamStats)
	bodyText := strings.Join(paragraphs, "\n\n")
	richBody := buildRichPostBody(paragraphs)

	input := fbsvc.CreateForumThreadInput{
		ForumID:           PostGameForumID + "-simnfl",
		ForumPath:         []string{PostGameForumID, "simnfl"},
		Title:             title,
		AuthorUID:         "system",
		AuthorUsername:    "SimSN",
		AuthorDisplayName: "SimSN System",
		CreatedByType:     fbsvc.CreatedBySystem,
		ThreadType:        fbsvc.ThreadTypeGameReference,
		FirstPostBodyText: bodyText,
		FirstPostBody:     richBody,
		ReferencedGameID:  gameID,
		ReferencedLeague:  "nfl",
		ExternalEventKey:  eventKey,
	}

	thread, err := fbsvc.CreateThread(ctx, input)
	if err != nil {
		log.Printf("ForumManager: failed to create NFL postgame thread for game %s: %v", gameID, err)
		return
	}

	log.Printf("ForumManager: created NFL postgame thread %s for game %s (%s)", thread.ID, gameID, title)
}

// ─────────────────────────────────────────────
// Body builders
// ─────────────────────────────────────────────

func buildCFBPostGameParagraphs(
	game structs.CollegeGame,
	home structs.CollegeTeamStats,
	away structs.CollegeTeamStats,
) []string {
	return buildPostGameParagraphs(
		game.AwayTeam, game.HomeTeam,
		game.AwayTeamScore, game.HomeTeamScore,
		game.Stadium, game.City, game.State,
		game.GameTemp, game.WindSpeed, game.WindCategory, game.Precip,
		game.IsDomed,
		game.MVP,
		away.BaseTeamStats, home.BaseTeamStats,
	)
}

func buildNFLPostGameParagraphs(
	game structs.NFLGame,
	home structs.NFLTeamStats,
	away structs.NFLTeamStats,
) []string {
	return buildPostGameParagraphs(
		game.AwayTeam, game.HomeTeam,
		game.AwayTeamScore, game.HomeTeamScore,
		game.Stadium, game.City, game.State,
		game.GameTemp, game.WindSpeed, game.WindCategory, game.Precip,
		game.IsDomed,
		game.MVP,
		away.BaseTeamStats, home.BaseTeamStats,
	)
}

// buildPostGameParagraphs constructs the ordered list of paragraph strings
// shared by both CFB and NFL forum post bodies.
func buildPostGameParagraphs(
	awayTeam, homeTeam string,
	awayScore, homeScore int,
	stadium, city, state string,
	gameTemp, windSpeed float64, windCategory, precip string,
	isDomed bool,
	mvp string,
	away, home structs.BaseTeamStats,
) []string {
	paras := []string{}

	// ── Final score ──────────────────────────────────────────────────────────
	paras = append(paras, fmt.Sprintf(
		"FINAL: %s %d, %s %d",
		awayTeam, awayScore, homeTeam, homeScore,
	))

	// ── Quarter-by-quarter scoring ───────────────────────────────────────────
	awayQs := formatQuarters(awayTeam, away)
	homeQs := formatQuarters(homeTeam, home)
	paras = append(paras, "SCORING BY QUARTER:\n"+awayQs+"\n"+homeQs)

	// ── Offensive stats ──────────────────────────────────────────────────────
	offLines := []string{
		"OFFENSE:",
		fmt.Sprintf("  %-20s  PASS: %4d yds  %dTD  %d INT   |   RUSH: %4d yds  %dTD",
			awayTeam,
			away.PassingYards, away.PassingTouchdowns, away.PassingInterceptions,
			away.RushingYards, away.RushingTouchdowns,
		),
		fmt.Sprintf("  %-20s  PASS: %4d yds  %dTD  %d INT   |   RUSH: %4d yds  %dTD",
			homeTeam,
			home.PassingYards, home.PassingTouchdowns, home.PassingInterceptions,
			home.RushingYards, home.RushingTouchdowns,
		),
	}
	paras = append(paras, strings.Join(offLines, "\n"))

	// ── Defensive / turnover stats ───────────────────────────────────────────
	defLines := []string{
		"DEFENSE & TURNOVERS:",
		fmt.Sprintf("  %-20s  SACKS: %.0f   INTs: %d   TFL: %.0f   FORCED FMBL: %d",
			awayTeam,
			away.DefensiveSacks, away.DefensiveInterceptions, away.TacklesForLoss, away.ForcedFumbles,
		),
		fmt.Sprintf("  %-20s  SACKS: %.0f   INTs: %d   TFL: %.0f   FORCED FMBL: %d",
			homeTeam,
			home.DefensiveSacks, home.DefensiveInterceptions, home.TacklesForLoss, home.ForcedFumbles,
		),
	}
	paras = append(paras, strings.Join(defLines, "\n"))

	// ── Special teams ────────────────────────────────────────────────────────
	if away.FieldGoalsAttempted > 0 || home.FieldGoalsAttempted > 0 {
		stLines := []string{
			"SPECIAL TEAMS:",
			fmt.Sprintf("  %-20s  FG: %d/%d (long %d)   XP: %d/%d",
				awayTeam,
				away.FieldGoalsMade, away.FieldGoalsAttempted, away.LongestFieldGoal,
				away.ExtraPointsMade, away.ExtraPointsAttempted,
			),
			fmt.Sprintf("  %-20s  FG: %d/%d (long %d)   XP: %d/%d",
				homeTeam,
				home.FieldGoalsMade, home.FieldGoalsAttempted, home.LongestFieldGoal,
				home.ExtraPointsMade, home.ExtraPointsAttempted,
			),
		}
		paras = append(paras, strings.Join(stLines, "\n"))
	}

	// ── Game info ────────────────────────────────────────────────────────────
	paras = append(paras, fmt.Sprintf("VENUE: %s — %s, %s", stadium, city, state))

	if !isDomed {
		weatherLine := fmt.Sprintf("WEATHER: %.0f°F", gameTemp)
		if windSpeed > 0 {
			weatherLine += fmt.Sprintf("  |  Wind: %.0f mph (%s)", windSpeed, windCategory)
		}
		if precip != "" && precip != "None" && precip != "Clear" {
			weatherLine += fmt.Sprintf("  |  %s", precip)
		}
		paras = append(paras, weatherLine)
	}

	// ── MVP ──────────────────────────────────────────────────────────────────
	if mvp != "" {
		paras = append(paras, fmt.Sprintf("MVP: %s", mvp))
	}

	// ── Discussion prompt ────────────────────────────────────────────────────
	paras = append(paras, "Postgame discussion is open. Share your reactions below.")

	return paras
}

// formatQuarters returns a single line showing per-quarter scoring for a team.
func formatQuarters(team string, s structs.BaseTeamStats) string {
	line := fmt.Sprintf("  %-20s  Q1: %2d  Q2: %2d  Q3: %2d  Q4: %2d",
		team, s.Score1Q, s.Score2Q, s.Score3Q, s.Score4Q)
	if s.ScoreOT > 0 {
		line += fmt.Sprintf("  OT: %2d", s.ScoreOT)
	}
	line += fmt.Sprintf("  TOTAL: %2d", s.Score1Q+s.Score2Q+s.Score3Q+s.Score4Q+s.ScoreOT)
	return line
}

// ─────────────────────────────────────────────
// Rich text helpers
// ─────────────────────────────────────────────

// buildRichPostBody converts a slice of paragraph strings into a ProseMirror
// document compatible with the frontend's RichTextDocument interface.
func buildRichPostBody(paragraphs []string) map[string]interface{} {
	content := make([]map[string]interface{}, 0, len(paragraphs))
	for _, p := range paragraphs {
		content = append(content, map[string]interface{}{
			"type": "paragraph",
			"content": []map[string]interface{}{
				{"type": "text", "text": p},
			},
		})
	}
	return map[string]interface{}{
		"type":    "doc",
		"content": content,
	}
}

// ─────────────────────────────────────────────
// Shared title helper
// ─────────────────────────────────────────────

func buildPostGameThreadTitle(awayTeam, homeTeam, gameTitle string) string {
	if gameTitle != "" {
		return fmt.Sprintf("Postgame Thread: %s", gameTitle)
	}
	return fmt.Sprintf("Postgame Thread: %s at %s", awayTeam, homeTeam)
}

// ─────────────────────────────────────────────
// Transfer portal helpers
// ─────────────────────────────────────────────

// TransferIntentionsSummary bundles all the counters produced by the transfer
// intentions run so they can be passed to the forum-thread creator without a
// long argument list.
type TransferIntentionsSummary struct {
	Season                 int
	TransferCount          int
	FreshmanCount          int
	RedshirtFreshmanCount  int
	SophomoreCount         int
	RedshirtSophomoreCount int
	JuniorCount            int
	RedshirtJuniorCount    int
	SeniorCount            int
	RedshirtSeniorCount    int
	LowCount               int
	MediumCount            int
	HighCount              int
}

// CreateTransferIntentionsForumThread creates a system-generated forum thread in
// the "daily" forum summarising the transfer intentions run for the given season.
// The operation is idempotent: calling it twice for the same season has no effect.
func CreateTransferIntentionsForumThread(summary TransferIntentionsSummary) {
	ctx := context.Background()

	title := fmt.Sprintf("SimCFB: Season %d Transfer Intentions", summary.Season)
	eventKey := fmt.Sprintf("transfer_intentions_thread:cfb:season%d", summary.Season)

	paragraphs := buildTransferIntentionsParagraphs(summary)
	bodyText := strings.Join(paragraphs, "\n\n")
	richBody := buildRichPostBody(paragraphs)

	input := fbsvc.CreateForumThreadInput{
		ForumID:           "media-simcfb",
		ForumPath:         []string{"media", "simcfb"},
		Title:             title,
		AuthorUID:         "system",
		AuthorUsername:    "SimSN",
		AuthorDisplayName: "SimSN System",
		CreatedByType:     fbsvc.CreatedBySystem,
		ThreadType:        fbsvc.ThreadTypeStandard,
		FirstPostBodyText: bodyText,
		FirstPostBody:     richBody,
		ReferencedLeague:  "cfb",
		ExternalEventKey:  eventKey,
	}

	thread, err := fbsvc.CreateThread(ctx, input)
	if err != nil {
		log.Printf("ForumManager: failed to create transfer intentions thread for season %d: %v", summary.Season, err)
		return
	}

	log.Printf("ForumManager: created transfer intentions thread %s for season %d", thread.ID, summary.Season)
}

func buildTransferIntentionsParagraphs(s TransferIntentionsSummary) []string {
	var paragraphs []string

	paragraphs = append(paragraphs,
		fmt.Sprintf(
			"Transfer season is underway for Season %d. A total of %d players have announced their intention to enter the transfer portal. Teams have one week to submit promises to retain their players.",
			s.Season, s.TransferCount,
		),
	)

	// Year-by-year breakdown
	paragraphs = append(paragraphs,
		fmt.Sprintf(
			"Class breakdown — Freshmen: %d | RS Freshmen: %d | Sophomores: %d | RS Sophomores: %d | Juniors: %d | RS Juniors: %d | Seniors: %d | RS Seniors: %d.",
			s.FreshmanCount, s.RedshirtFreshmanCount,
			s.SophomoreCount, s.RedshirtSophomoreCount,
			s.JuniorCount, s.RedshirtJuniorCount,
			s.SeniorCount, s.RedshirtSeniorCount,
		),
	)

	// Likeliness breakdown
	paragraphs = append(paragraphs,
		fmt.Sprintf(
			"Transfer likeliness — Low: %d | Medium: %d | High: %d.",
			s.LowCount, s.MediumCount, s.HighCount,
		),
	)

	paragraphs = append(paragraphs,
		"Which transfers are you keeping an eye on this season? Share your thoughts below!",
	)

	return paragraphs
}

// ─────────────────────────────────────────────
// Transfer portal sync thread
// ─────────────────────────────────────────────

// CreateTransferPortalSyncForumThread creates a system-generated forum thread in
// the "media-cfb" subforum summarising the signings from a single transfer portal
// sync round. signings is a list of human-readable player labels for every player
// that signed with a new team during the sync.
// The operation is idempotent: calling it twice for the same season/round has no
// effect.
func CreateTransferPortalSyncForumThread(season, round int, signings []string) {
	ctx := context.Background()

	title := fmt.Sprintf("SimCFB: Season %d Transfer Portal — Round %d Results", season, round)
	eventKey := fmt.Sprintf("transfer_portal_sync:cfb:season%d:round%d", season, round)

	paragraphs := buildTransferPortalSyncParagraphs(season, round, signings)
	bodyText := strings.Join(paragraphs, "\n\n")
	richBody := buildRichPostBody(paragraphs)

	input := fbsvc.CreateForumThreadInput{
		ForumID:           "media-simcfb",
		ForumPath:         []string{"media", "simcfb"},
		Title:             title,
		AuthorUID:         "system",
		AuthorUsername:    "SimSN",
		AuthorDisplayName: "SimSN System",
		CreatedByType:     fbsvc.CreatedBySystem,
		ThreadType:        fbsvc.ThreadTypeStandard,
		FirstPostBodyText: bodyText,
		FirstPostBody:     richBody,
		ReferencedLeague:  "cfb",
		ExternalEventKey:  eventKey,
	}

	thread, err := fbsvc.CreateThread(ctx, input)
	if err != nil {
		log.Printf("ForumManager: failed to create transfer portal sync thread for season %d round %d: %v", season, round, err)
		return
	}

	log.Printf("ForumManager: created transfer portal sync thread %s for season %d round %d", thread.ID, season, round)
}

func buildTransferPortalSyncParagraphs(season, round int, signings []string) []string {
	var paragraphs []string

	count := len(signings)
	if count == 0 {
		paragraphs = append(paragraphs,
			fmt.Sprintf(
				"Transfer Portal Round %d is complete for Season %d. No players signed with new programs this round.",
				round, season,
			),
		)
	} else {
		paragraphs = append(paragraphs,
			fmt.Sprintf(
				"Transfer Portal Round %d results are in for Season %d. A total of %d player(s) have signed with new programs this round.",
				round, season, count,
			),
		)
		for _, label := range signings {
			paragraphs = append(paragraphs, label)
		}
	}

	paragraphs = append(paragraphs, "Discuss the latest transfer portal news below!")

	return paragraphs
}

// ─────────────────────────────────────────────
// Transfer portal open thread
// ─────────────────────────────────────────────

// CreateTransferPortalOpenForumThread creates a system-generated forum thread in
// the "media-simcfb" subforum announcing the transfer portal is open for the
// given season, with one paragraph per player entering the portal.
// playerLabels is a list of human-readable labels built before WillTransfer()
// clears each player’s TeamAbbr.
// The operation is idempotent: calling it twice for the same season has no effect.
func CreateTransferPortalOpenForumThread(season int, playerLabels []string) {
	ctx := context.Background()

	title := fmt.Sprintf("SimCFB: Season %d Transfer Portal is Now Open", season)
	eventKey := fmt.Sprintf("transfer_portal_open:cfb:season%d", season)

	paragraphs := buildTransferPortalOpenParagraphs(season, playerLabels)
	bodyText := strings.Join(paragraphs, "\n\n")
	richBody := buildRichPostBody(paragraphs)

	input := fbsvc.CreateForumThreadInput{
		ForumID:           "media-simcfb",
		ForumPath:         []string{"media", "simcfb"},
		Title:             title,
		AuthorUID:         "system",
		AuthorUsername:    "SimSN",
		AuthorDisplayName: "SimSN System",
		CreatedByType:     fbsvc.CreatedBySystem,
		ThreadType:        fbsvc.ThreadTypeStandard,
		FirstPostBodyText: bodyText,
		FirstPostBody:     richBody,
		ReferencedLeague:  "cfb",
		ExternalEventKey:  eventKey,
	}

	thread, err := fbsvc.CreateThread(ctx, input)
	if err != nil {
		log.Printf("ForumManager: failed to create transfer portal open thread for season %d: %v", season, err)
		return
	}

	log.Printf("ForumManager: created transfer portal open thread %s for season %d", thread.ID, season)
}

func buildTransferPortalOpenParagraphs(season int, playerLabels []string) []string {
	var paragraphs []string

	count := len(playerLabels)
	paragraphs = append(paragraphs,
		fmt.Sprintf(
			"The SimCFB Transfer Portal is now open for Season %d. A total of %d player(s) have officially entered the portal and are seeking a new home.",
			season, count,
		),
	)

	if count > 0 {
		paragraphs = append(paragraphs, "The following players have entered the transfer portal:")
		for _, label := range playerLabels {
			paragraphs = append(paragraphs, label)
		}
	}

	paragraphs = append(paragraphs, "Which players are you targeting this transfer portal cycle? Share your thoughts below!")

	return paragraphs
}

// ─────────────────────────────────────────────
// Recruiting sync thread
// ─────────────────────────────────────────────

// CreateRecruitingSyncForumThread creates a system-generated weekly thread in
// the "media-simcfb" subforum listing every recruit that signed with a program
// during the sync. signings is a list of human-readable labels built at the
// moment each recruit commits (before any state is mutated further).
// The operation is idempotent: calling it twice for the same season/week has no
// effect.
func CreateRecruitingSyncForumThread(season, week int, signings []string) {
	ctx := context.Background()

	title := fmt.Sprintf("SimCFB: Season %d Week %d Recruiting Commitments", season, week)
	eventKey := fmt.Sprintf("recruiting_sync:cfb:season%d:week%d", season, week)

	paragraphs := buildRecruitingSyncParagraphs(season, week, signings)
	bodyText := strings.Join(paragraphs, "\n\n")
	richBody := buildRichPostBody(paragraphs)

	input := fbsvc.CreateForumThreadInput{
		ForumID:           "media-simcfb",
		ForumPath:         []string{"media", "simcfb"},
		Title:             title,
		AuthorUID:         "system",
		AuthorUsername:    "SimSN",
		AuthorDisplayName: "SimSN System",
		CreatedByType:     fbsvc.CreatedBySystem,
		ThreadType:        fbsvc.ThreadTypeStandard,
		FirstPostBodyText: bodyText,
		FirstPostBody:     richBody,
		ReferencedLeague:  "cfb",
		ExternalEventKey:  eventKey,
	}

	thread, err := fbsvc.CreateThread(ctx, input)
	if err != nil {
		log.Printf("ForumManager: failed to create recruiting sync thread for season %d week %d: %v", season, week, err)
		return
	}

	log.Printf("ForumManager: created recruiting sync thread %s for season %d week %d", thread.ID, season, week)
}

func buildRecruitingSyncParagraphs(season, week int, signings []string) []string {
	var paragraphs []string

	count := len(signings)
	if count == 0 {
		paragraphs = append(paragraphs,
			fmt.Sprintf(
				"Week %d recruiting is complete for Season %d. No recruits signed with a program this week.",
				week, season,
			),
		)
	} else {
		paragraphs = append(paragraphs,
			fmt.Sprintf(
				"Week %d recruiting results are in for Season %d. A total of %d recruit(s) have committed to a program this week.",
				week, season, count,
			),
		)
		for _, label := range signings {
			paragraphs = append(paragraphs, label)
		}
	}

	paragraphs = append(paragraphs, "React to the latest commitments and discuss your team\u2019s recruiting class below!")

	return paragraphs
}
