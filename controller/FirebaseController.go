package controller

import (
	"context"
	"encoding/json"
	"net/http"

	fbsvc "github.com/CalebRose/SimFBA/firebase"
)

// TestNotificationToTuscan sends a test notification to the user "TuscanSota".
// GET /firebase/test/notification/
func TestNotificationToTuscan(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	uids := fbsvc.ResolveUIDsByUsernames(ctx, []string{"TuscanSota"})
	if len(uids) == 0 {
		http.Error(w, "Could not resolve UID for TuscanSota", http.StatusNotFound)
		return
	}

	eventKey := fbsvc.BuildSourceEventKey("test_notification", "tuscan")
	err := fbsvc.NotifyGameplanIssue(ctx, fbsvc.GameplanNotificationInput{
		League:         "cfb",
		Domain:         fbsvc.DomainSystem,
		TeamID:         0,
		TeamName:       "Test Team",
		TeamAbbr:       "TST",
		Message:        "This is a test notification from the SimFBA API. If you can see this, Firebase notifications are working correctly!",
		RecipientUIDs:  uids,
		SourceEventKey: eventKey,
	})
	if err != nil {
		http.Error(w, "Failed to send notification: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":      true,
		"message": "Test notification sent to TuscanSota",
		"uids":    uids,
	})
}

// TestForumPost creates a test thread with an initial post in the "daily" forum.
// GET /firebase/test/forum/
func TestForumPost(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	input := fbsvc.CreateForumThreadInput{
		ForumID:           "daily",
		ForumPath:         []string{"daily"},
		Title:             "API Test Thread",
		AuthorUID:         "system",
		AuthorUsername:    "SimSN",
		AuthorDisplayName: "SimSN System",
		CreatedByType:     fbsvc.CreatedBySystem,
		ThreadType:        fbsvc.ThreadTypeStandard,
		FirstPostBodyText: "This is a test thread created by the SimFBA API. If you can see this, forum thread creation is working correctly!",
		ExternalEventKey:  fbsvc.BuildSourceEventKey("test_forum_thread", "daily"),
	}

	thread, err := fbsvc.CreateThread(ctx, input)
	if err != nil {
		http.Error(w, "Failed to create forum thread: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":       true,
		"message":  "Test thread created in daily forum",
		"threadId": thread.ID,
		"title":    thread.Title,
	})
}
