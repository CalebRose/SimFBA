package managers

import (
	"context"
	"fmt"
	"log"
	"sort"
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
	homePlayerStats []structs.CollegePlayerStats,
	awayPlayerStats []structs.CollegePlayerStats,
	collegePlayersMap map[uint]structs.CollegePlayer,

) {
	ctx := context.Background()

	gameID := strconv.Itoa(int(game.ID))
	eventKey := fmt.Sprintf("postgame_thread:cfb:season%d:game%s", game.SeasonID, gameID)

	title := buildPostGameThreadTitle(game.AwayTeam, game.HomeTeam, game.GameTitle)
	nodes := buildCFBPostGameNodes(game, homeTeamStats, awayTeamStats, homePlayerStats, awayPlayerStats, collegePlayersMap)
	bodyText := nodesToPlainText(nodes)
	richBody := buildRichDoc(nodes)

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
	homePlayerStats []structs.NFLPlayerStats,
	awayPlayerStats []structs.NFLPlayerStats,
	proPlayersMap map[uint]structs.NFLPlayer,
) {
	ctx := context.Background()

	gameID := strconv.Itoa(int(game.ID))
	eventKey := fmt.Sprintf("postgame_thread:nfl:season%d:game%s", game.SeasonID, gameID)

	title := buildPostGameThreadTitle(game.AwayTeam, game.HomeTeam, game.GameTitle)
	nodes := buildNFLPostGameNodes(game, homeTeamStats, awayTeamStats, homePlayerStats, awayPlayerStats, proPlayersMap)
	bodyText := nodesToPlainText(nodes)
	richBody := buildRichDoc(nodes)

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

func buildCFBPostGameNodes(
	game structs.CollegeGame,
	home structs.CollegeTeamStats,
	away structs.CollegeTeamStats,
	homePlayerStats []structs.CollegePlayerStats,
	awayPlayerStats []structs.CollegePlayerStats,
	collegePlayersMap map[uint]structs.CollegePlayer,
) []map[string]interface{} {
	nodes := buildPostGameNodes(
		game.AwayTeam, game.HomeTeam,
		game.AwayTeamScore, game.HomeTeamScore,
		game.Stadium, game.City, game.State,
		game.GameTemp, game.WindSpeed, game.WindCategory, game.Precip,
		game.IsDomed,
		game.MVP,
		away.BaseTeamStats, home.BaseTeamStats,
	)
	awayRows := toCFBPlayerStatRows(awayPlayerStats, collegePlayersMap)
	homeRows := toCFBPlayerStatRows(homePlayerStats, collegePlayersMap)
	nodes = appendPlayerStatTables(nodes, game.AwayTeam, awayRows, game.HomeTeam, homeRows)
	nodes = append(nodes, rtParagraph("Postgame discussion is open. Share your reactions below."))
	return nodes
}

func buildNFLPostGameNodes(
	game structs.NFLGame,
	home structs.NFLTeamStats,
	away structs.NFLTeamStats,
	homePlayerStats []structs.NFLPlayerStats,
	awayPlayerStats []structs.NFLPlayerStats,
	proPlayersMap map[uint]structs.NFLPlayer,
) []map[string]interface{} {
	nodes := buildPostGameNodes(
		game.AwayTeam, game.HomeTeam,
		game.AwayTeamScore, game.HomeTeamScore,
		game.Stadium, game.City, game.State,
		game.GameTemp, game.WindSpeed, game.WindCategory, game.Precip,
		game.IsDomed,
		game.MVP,
		away.BaseTeamStats, home.BaseTeamStats,
	)
	awayRows := toNFLPlayerStatRows(awayPlayerStats, proPlayersMap)
	homeRows := toNFLPlayerStatRows(homePlayerStats, proPlayersMap)
	nodes = appendPlayerStatTables(nodes, game.AwayTeam, awayRows, game.HomeTeam, homeRows)
	nodes = append(nodes, rtParagraph("Postgame discussion is open. Share your reactions below."))
	return nodes
}

// buildPostGameNodes constructs the ordered list of ProseMirror content nodes
// shared by both CFB and NFL forum post bodies.
func buildPostGameNodes(
	awayTeam, homeTeam string,
	awayScore, homeScore int,
	stadium, city, state string,
	gameTemp, windSpeed float64, windCategory, precip string,
	isDomed bool,
	mvp string,
	away, home structs.BaseTeamStats,
) []map[string]interface{} {
	nodes := []map[string]interface{}{}

	// ── Final score ──────────────────────────────────────────────────────────
	nodes = append(nodes, rtBoldParagraph(fmt.Sprintf(
		"FINAL: %s %d, %s %d",
		awayTeam, awayScore, homeTeam, homeScore,
	)))

	// ── Quarter scoring table ─────────────────────────────────────────────────
	nodes = append(nodes, rtHeading(3, "Scoring by Quarter"))
	awayTotal := away.Score1Q + away.Score2Q + away.Score3Q + away.Score4Q + away.ScoreOT
	homeTotal := home.Score1Q + home.Score2Q + home.Score3Q + home.Score4Q + home.ScoreOT
	nodes = append(nodes, rtTableNode(
		[]string{"Team", "Q1", "Q2", "Q3", "Q4", "OT", "Total"},
		[][]string{
			{awayTeam, fmt.Sprintf("%d", away.Score1Q), fmt.Sprintf("%d", away.Score2Q), fmt.Sprintf("%d", away.Score3Q), fmt.Sprintf("%d", away.Score4Q), fmt.Sprintf("%d", away.ScoreOT), fmt.Sprintf("%d", awayTotal)},
			{homeTeam, fmt.Sprintf("%d", home.Score1Q), fmt.Sprintf("%d", home.Score2Q), fmt.Sprintf("%d", home.Score3Q), fmt.Sprintf("%d", home.Score4Q), fmt.Sprintf("%d", home.ScoreOT), fmt.Sprintf("%d", homeTotal)},
		},
	))

	// ── Venue ─────────────────────────────────────────────────────────────────
	nodes = append(nodes, rtParagraph(fmt.Sprintf("Venue: %s — %s, %s", stadium, city, state)))

	// ── Weather ───────────────────────────────────────────────────────────────
	if !isDomed {
		weatherLine := fmt.Sprintf("Weather: %.0f°F", gameTemp)
		if windSpeed > 0 {
			weatherLine += fmt.Sprintf("  |  Wind: %.0f mph (%s)", windSpeed, windCategory)
		}
		if precip != "" && precip != "None" && precip != "Clear" {
			weatherLine += fmt.Sprintf("  |  %s", precip)
		}
		nodes = append(nodes, rtParagraph(weatherLine))
	}

	// ── MVP ───────────────────────────────────────────────────────────────────
	if mvp != "" {
		nodes = append(nodes, rtParagraph(fmt.Sprintf("MVP: %s", mvp)))
	}

	// ── Offense table ─────────────────────────────────────────────────────────
	nodes = append(nodes, rtHeading(3, "Offense"))
	nodes = append(nodes, rtTableNode(
		[]string{"Team", "Pass Yds", "Pass TD", "INT", "Rush Yds", "Rush TD"},
		[][]string{
			{awayTeam, fmt.Sprintf("%d", away.PassingYards), fmt.Sprintf("%d", away.PassingTouchdowns), fmt.Sprintf("%d", away.PassingInterceptions), fmt.Sprintf("%d", away.RushingYards), fmt.Sprintf("%d", away.RushingTouchdowns)},
			{homeTeam, fmt.Sprintf("%d", home.PassingYards), fmt.Sprintf("%d", home.PassingTouchdowns), fmt.Sprintf("%d", home.PassingInterceptions), fmt.Sprintf("%d", home.RushingYards), fmt.Sprintf("%d", home.RushingTouchdowns)},
		},
	))

	// ── Defense table ─────────────────────────────────────────────────────────
	nodes = append(nodes, rtHeading(3, "Defense & Turnovers"))
	nodes = append(nodes, rtTableNode(
		[]string{"Team", "Sacks", "INTs", "TFL", "Forced Fmbl"},
		[][]string{
			{awayTeam, fmt.Sprintf("%.0f", away.DefensiveSacks), fmt.Sprintf("%d", away.DefensiveInterceptions), fmt.Sprintf("%.0f", away.TacklesForLoss), fmt.Sprintf("%d", away.ForcedFumbles)},
			{homeTeam, fmt.Sprintf("%.0f", home.DefensiveSacks), fmt.Sprintf("%d", home.DefensiveInterceptions), fmt.Sprintf("%.0f", home.TacklesForLoss), fmt.Sprintf("%d", home.ForcedFumbles)},
		},
	))

	// ── Special teams table ───────────────────────────────────────────────────
	if away.FieldGoalsAttempted > 0 || home.FieldGoalsAttempted > 0 {
		nodes = append(nodes, rtHeading(3, "Special Teams"))
		nodes = append(nodes, rtTableNode(
			[]string{"Team", "FG", "FGA", "Long", "XP", "XPA"},
			[][]string{
				{awayTeam, fmt.Sprintf("%d", away.FieldGoalsMade), fmt.Sprintf("%d", away.FieldGoalsAttempted), fmt.Sprintf("%d", away.LongestFieldGoal), fmt.Sprintf("%d", away.ExtraPointsMade), fmt.Sprintf("%d", away.ExtraPointsAttempted)},
				{homeTeam, fmt.Sprintf("%d", home.FieldGoalsMade), fmt.Sprintf("%d", home.FieldGoalsAttempted), fmt.Sprintf("%d", home.LongestFieldGoal), fmt.Sprintf("%d", home.ExtraPointsMade), fmt.Sprintf("%d", home.ExtraPointsAttempted)},
			},
		))
	}

	return nodes
}

// ─────────────────────────────────────────────
// Player stat helpers
// ─────────────────────────────────────────────

// playerStatRow pairs a display label with a player's per-game stats.
type playerStatRow struct {
	Label string
	structs.BasePlayerStats
}

// toCFBPlayerStatRows converts a slice of college player stats into playerStatRows.
func toCFBPlayerStatRows(stats []structs.CollegePlayerStats, playerMap map[uint]structs.CollegePlayer) []playerStatRow {
	rows := make([]playerStatRow, 0, len(stats))
	for _, s := range stats {
		p, ok := playerMap[uint(s.CollegePlayerID)]
		if !ok {
			continue
		}
		label := fmt.Sprintf("[%d] %s %s %s %s", p.ID, p.TeamAbbr, p.Position, p.FirstName, p.LastName)
		rows = append(rows, playerStatRow{Label: label, BasePlayerStats: s.BasePlayerStats})
	}
	return rows
}

// toNFLPlayerStatRows converts a slice of NFL player stats into playerStatRows.
func toNFLPlayerStatRows(stats []structs.NFLPlayerStats, playerMap map[uint]structs.NFLPlayer) []playerStatRow {
	rows := make([]playerStatRow, 0, len(stats))
	for _, s := range stats {
		p, ok := playerMap[uint(s.NFLPlayerID)]
		if !ok {
			continue
		}
		label := fmt.Sprintf("[%d] %s %s %s %s", p.ID, p.TeamAbbr, p.Position, p.FirstName, p.LastName)
		rows = append(rows, playerStatRow{Label: label, BasePlayerStats: s.BasePlayerStats})
	}
	return rows
}

// appendPlayerStatTables appends per-player stat section tables to nodes.
// Away rows are listed before home rows within each table.
func appendPlayerStatTables(
	nodes []map[string]interface{},
	awayTeam string, awayRows []playerStatRow,
	homeTeam string, homeRows []playerStatRow,
) []map[string]interface{} {
	allRows := append(awayRows, homeRows...)

	// ── Passing stats ─────────────────────────────────────────────────────────
	var passRows [][]string
	passers := []playerStatRow{}

	for _, r := range allRows {
		if r.PassAttempts > 0 {
			passers = append(passers, r)
		}
	}
	sort.Slice(passers, func(i, j int) bool {
		// Sort by Team ID and then by Passing Yards DESC within each team
		return passers[i].TeamID < passers[j].TeamID && (passers[i].TeamID == passers[j].TeamID && passers[i].PassingYards > passers[j].PassingYards)
	})

	for _, p := range passers {
		if p.PassAttempts > 0 {
			passRows = append(passRows, []string{
				p.Label,
				fmt.Sprintf("%d/%d", p.PassCompletions, p.PassAttempts),
				fmt.Sprintf("%d", p.PassingYards),
				fmt.Sprintf("%d", p.PassingTDs),
				fmt.Sprintf("%d", p.Interceptions),
				fmt.Sprintf("%d", p.LongestPass),
				fmt.Sprintf("%d", p.Sacks),
			})
		}
	}
	if len(passRows) > 0 {
		nodes = append(nodes, rtHeading(3, "Passing Stats"))
		nodes = append(nodes, rtTableNode(
			[]string{"Player", "C/ATT", "Yds", "TD", "INT", "Long", "Sacked"},
			passRows,
		))
	}

	// ── Rushing stats ─────────────────────────────────────────────────────────
	var rushRows [][]string
	rushers := []playerStatRow{}
	for _, r := range allRows {
		if r.RushAttempts > 0 {

			rushers = append(rushers, r)
		}
	}

	sort.Slice(rushers, func(i, j int) bool {
		// Sort by Team ID and then by Rushing Yards DESC within each team
		return rushers[i].TeamID < rushers[j].TeamID && (rushers[i].TeamID == rushers[j].TeamID && rushers[i].RushingYards > rushers[j].RushingYards)
	})

	for _, r := range rushers {
		if r.RushAttempts > 0 {
			rushRows = append(rushRows, []string{
				r.Label,
				fmt.Sprintf("%d", r.RushAttempts),
				fmt.Sprintf("%d", r.RushingYards),
				fmt.Sprintf("%d", r.RushingTDs),
				fmt.Sprintf("%d", r.LongestRush),
				fmt.Sprintf("%d", r.Fumbles),
			})
		}
	}
	if len(rushRows) > 0 {
		nodes = append(nodes, rtHeading(3, "Rushing Stats"))
		nodes = append(nodes, rtTableNode(
			[]string{"Player", "Att", "Yds", "TD", "Long", "Fmbl"},
			rushRows,
		))
	}

	// ── Receiving stats ───────────────────────────────────────────────────────
	var recRows [][]string
	var receivers []playerStatRow
	for _, r := range allRows {
		if r.Targets > 0 {
			receivers = append(receivers, r)
		}
	}

	sort.Slice(receivers, func(i, j int) bool {
		// Sort by Team ID and then by Receiving Yards DESC within each team
		return receivers[i].TeamID < receivers[j].TeamID && (receivers[i].TeamID == receivers[j].TeamID && receivers[i].ReceivingYards > receivers[j].ReceivingYards)
	})

	for _, r := range receivers {
		if r.Targets > 0 {
			recRows = append(recRows, []string{
				r.Label,
				fmt.Sprintf("%d", r.Targets),
				fmt.Sprintf("%d", r.Catches),
				fmt.Sprintf("%d", r.ReceivingYards),
				fmt.Sprintf("%d", r.ReceivingTDs),
				fmt.Sprintf("%d", r.LongestReception),
			})
		}
	}
	if len(recRows) > 0 {
		nodes = append(nodes, rtHeading(3, "Receiving Stats"))
		nodes = append(nodes, rtTableNode(
			[]string{"Player", "Tgt", "Rec", "Yds", "TD", "Long"},
			recRows,
		))
	}

	// ── Defensive stats ───────────────────────────────────────────────────────
	var defRows [][]string
	var defenders []playerStatRow
	// Filter to players with at least 1 solo or assisted tackle, then sort by Team ID and total tackles DESC within each team
	for _, r := range allRows {
		if r.SoloTackles+r.AssistedTackles > 0 {
			defenders = append(defenders, r)
		}
	}

	sort.Slice(defenders, func(i, j int) bool {
		// Sort by Team ID and then by Total Tackles DESC within each team
		totalTacklesI := defenders[i].SoloTackles + defenders[i].AssistedTackles
		totalTacklesJ := defenders[j].SoloTackles + defenders[j].AssistedTackles
		return defenders[i].TeamID < defenders[j].TeamID && (defenders[i].TeamID == defenders[j].TeamID && totalTacklesI > totalTacklesJ)
	})

	for _, r := range defenders {
		if r.SoloTackles+r.AssistedTackles > 0 {
			defRows = append(defRows, []string{
				r.Label,
				fmt.Sprintf("%.0f", r.SoloTackles),
				fmt.Sprintf("%.0f", r.AssistedTackles),
				fmt.Sprintf("%.0f", r.TacklesForLoss),
				fmt.Sprintf("%.0f", r.SacksMade),
				fmt.Sprintf("%d", r.InterceptionsCaught),
				fmt.Sprintf("%d", r.PassDeflections),
				fmt.Sprintf("%d", r.ForcedFumbles),
			})
		}
	}
	if len(defRows) > 0 {
		nodes = append(nodes, rtHeading(3, "Defensive Stats"))
		nodes = append(nodes, rtTableNode(
			[]string{"Player", "Solo", "Ast", "TFL", "Sacks", "INT", "PD", "FF"},
			defRows,
		))
	}

	// ── Field Goals & Extra Points ────────────────────────────────────────────
	var fgRows [][]string
	var kickers []playerStatRow
	for _, r := range allRows {
		if r.FGAttempts > 0 || r.ExtraPointsAttempted > 0 {
			kickers = append(kickers, r)
		}
	}

	sort.Slice(kickers, func(i, j int) bool {
		// Sort by Team ID and then by FG Attempts DESC within each team
		return kickers[i].TeamID < kickers[j].TeamID && (kickers[i].TeamID == kickers[j].TeamID && kickers[i].FGAttempts > kickers[j].FGAttempts)
	})

	for _, r := range kickers {
		if r.FGAttempts > 0 || r.ExtraPointsAttempted > 0 {
			fgRows = append(fgRows, []string{
				r.Label,
				fmt.Sprintf("%d/%d", r.FGMade, r.FGAttempts),
				fmt.Sprintf("%d", r.LongestFG),
				fmt.Sprintf("%d/%d", r.ExtraPointsMade, r.ExtraPointsAttempted),
			})
		}
	}
	if len(fgRows) > 0 {
		nodes = append(nodes, rtHeading(3, "Field Goals & Extra Points"))
		nodes = append(nodes, rtTableNode(
			[]string{"Player", "FG M/A", "Long", "XP M/A"},
			fgRows,
		))
	}

	// ── Punting stats ─────────────────────────────────────────────────────────
	var puntRows [][]string
	var punters []playerStatRow
	for _, r := range allRows {
		if r.Punts > 0 {
			punters = append(punters, r)
		}
	}

	sort.Slice(punters, func(i, j int) bool {
		// Sort by Team ID and then by Punts DESC within each team
		return punters[i].TeamID < punters[j].TeamID && (punters[i].TeamID == punters[j].TeamID && punters[i].Punts > punters[j].Punts)
	})

	for _, r := range punters {
		if r.Punts > 0 {
			puntRows = append(puntRows, []string{
				r.Label,
				fmt.Sprintf("%d", r.Punts),
				fmt.Sprintf("%d", r.GrossPuntDistance),
				fmt.Sprintf("%d", r.NetPuntDistance),
				fmt.Sprintf("%d", r.PuntsInside20),
				fmt.Sprintf("%d", r.PuntTouchbacks),
			})
		}
	}
	if len(puntRows) > 0 {
		nodes = append(nodes, rtHeading(3, "Punting Stats"))
		nodes = append(nodes, rtTableNode(
			[]string{"Player", "Punts", "Gross Yds", "Net Yds", "In20", "TB"},
			puntRows,
		))
	}

	_ = awayTeam
	_ = homeTeam
	return nodes
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

// buildRichDoc wraps a slice of content nodes into a top-level ProseMirror doc.
func buildRichDoc(nodes []map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type":    "doc",
		"content": nodes,
	}
}

// rtParagraph creates a plain paragraph node.
func rtParagraph(text string) map[string]interface{} {
	return map[string]interface{}{
		"type": "paragraph",
		"content": []map[string]interface{}{
			{"type": "text", "text": text},
		},
	}
}

// rtBoldParagraph creates a paragraph with bold-marked text.
func rtBoldParagraph(text string) map[string]interface{} {
	return map[string]interface{}{
		"type": "paragraph",
		"content": []map[string]interface{}{
			{
				"type":  "text",
				"text":  text,
				"marks": []map[string]interface{}{{"type": "bold"}},
			},
		},
	}
}

// rtHeading creates a heading node at the given level (1–6).
func rtHeading(level int, text string) map[string]interface{} {
	return map[string]interface{}{
		"type":  "heading",
		"attrs": map[string]interface{}{"level": level, "textAlign": "left"},
		"content": []map[string]interface{}{
			{"type": "text", "text": text},
		},
	}
}

// rtTableCell creates a single table header or data cell wrapping text in a paragraph.
func rtTableCell(text string, isHeader bool) map[string]interface{} {
	cellType := "tableCell"
	if isHeader {
		cellType = "tableHeader"
	}
	return map[string]interface{}{
		"type":  cellType,
		"attrs": map[string]interface{}{"colspan": 1, "rowspan": 1, "colwidth": nil},
		"content": []map[string]interface{}{
			{
				"type":  "paragraph",
				"attrs": map[string]interface{}{"textAlign": nil},
				"content": []map[string]interface{}{
					{"type": "text", "text": text},
				},
			},
		},
	}
}

// rtTableNode builds a TipTap-compatible table node from header strings and row data.
// The first row is rendered as tableHeader cells; all subsequent rows as tableCell.
func rtTableNode(headers []string, rows [][]string) map[string]interface{} {
	tableRows := []map[string]interface{}{}

	headerCells := make([]map[string]interface{}, len(headers))
	for i, h := range headers {
		headerCells[i] = rtTableCell(h, true)
	}
	tableRows = append(tableRows, map[string]interface{}{
		"type":    "tableRow",
		"content": headerCells,
	})

	for _, row := range rows {
		cells := make([]map[string]interface{}, len(row))
		for i, cell := range row {
			cells[i] = rtTableCell(cell, false)
		}
		tableRows = append(tableRows, map[string]interface{}{
			"type":    "tableRow",
			"content": cells,
		})
	}

	return map[string]interface{}{
		"type":    "table",
		"content": tableRows,
	}
}

// nodesToPlainText extracts readable plain text from rich nodes for the bodyText field.
func nodesToPlainText(nodes []map[string]interface{}) string {
	var lines []string
	for _, node := range nodes {
		switch node["type"] {
		case "paragraph", "heading":
			if text := extractInlineText(node); text != "" {
				lines = append(lines, text)
			}
		case "table":
			if rows, ok := node["content"].([]map[string]interface{}); ok {
				for _, row := range rows {
					if cells, ok := row["content"].([]map[string]interface{}); ok {
						var cellTexts []string
						for _, cell := range cells {
							cellTexts = append(cellTexts, extractCellPlainText(cell))
						}
						lines = append(lines, strings.Join(cellTexts, "  |  "))
					}
				}
			}
		}
	}
	return strings.Join(lines, "\n\n")
}

func extractInlineText(node map[string]interface{}) string {
	if content, ok := node["content"].([]map[string]interface{}); ok {
		var texts []string
		for _, child := range content {
			if t, ok := child["text"].(string); ok {
				texts = append(texts, t)
			}
		}
		return strings.Join(texts, "")
	}
	return ""
}

func extractCellPlainText(cell map[string]interface{}) string {
	if content, ok := cell["content"].([]map[string]interface{}); ok {
		for _, para := range content {
			return extractInlineText(para)
		}
	}
	return ""
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

	paragraphs = append(paragraphs,
		[]string{"A friendly reminder on the rules of the transfer portal:",
			"Teams may place between 0 to 10 points on a transfer portal player. Each team has 100 points they can spread across any number of players.",
			"Once the board saves and if there are any points previously placed on a player (If it's past the first sync), teams cannot reduce the number of points on a player. They can only increase.",
			"The only time a team can place zero points on a transfer portal player is by removing them from the board entirely. Think of this as a game of chicken run.",
			"Teams can make a promise to a transfer portal player from their portal board. A promise made applied a multiplier onto the player's points.",
			"Once a player commits to a new team, they are removed from the transfer portal board and if a promise is made by the team signing the player, the promise becomes committed.",
			"Teams are responsible for fulfilling promises made to players. If a team fails to fulfill a promise, the player will re-enter the transfer portal. If a player graduates before a promise is committed, the team's portal reputation takes a larger hit.",
			"Promise wisely, and good luck to all teams in navigating this transfer portal cycle.",
		}...,
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
