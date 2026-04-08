package firebase

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// ─────────────────────────────────────────────
// Notification Service
// ─────────────────────────────────────────────

// NotifyPlayerInjured creates one notification document per recipient for a
// player-injury event.  Idempotent: if a notification with the same
// SourceEventKey already exists for a recipient it is skipped.
func NotifyPlayerInjured(ctx context.Context, input PlayerInjuryNotificationInput) error {
	if len(input.RecipientUIDs) == 0 {
		return nil
	}

	message := fmt.Sprintf(
		"%s suffered %s and is expected to miss %d game(s).",
		input.PlayerName, input.InjuryName, input.GamesMissed,
	)
	linkTo := BuildTeamRosterRoute(input.League, input.TeamID)

	return writeNotificationsIfNew(ctx, input.RecipientUIDs, ForumNotification{
		Type:           NotificationTypeInjury,
		Domain:         input.Domain,
		LinkTo:         linkTo,
		Message:        message,
		ActorUsername:  "SimSN",
		IsRead:         false,
		SourceEventKey: input.SourceEventKey,
	})
}

// NotifyGameplanIssue sends a depth-chart / gameplan penalty notification to
// the coach or owner of a team.  Idempotent via SourceEventKey.
func NotifyGameplanIssue(ctx context.Context, input GameplanNotificationInput) error {
	if len(input.RecipientUIDs) == 0 {
		return nil
	}
	linkTo := BuildTeamGameplanRoute(input.League, input.TeamID)

	return writeNotificationsIfNew(ctx, input.RecipientUIDs, ForumNotification{
		Type:           NotificationTypeGameplan,
		Domain:         input.Domain,
		LinkTo:         linkTo,
		Message:        input.Message,
		ActorUsername:  "SimSN",
		IsRead:         false,
		SourceEventKey: input.SourceEventKey,
	})
}

// NotifyRecruitingSyncMissed sends a notification to a coach when they failed to
// allocate any recruiting points during the weekly sync.  Idempotent via SourceEventKey.
func NotifyRecruitingSyncMissed(ctx context.Context, input RecruitingSyncMissedNotificationInput) error {
	if len(input.RecipientUIDs) == 0 {
		return nil
	}
	linkTo := BuildTeamRecruitingRoute("cfb", input.TeamID)

	return writeNotificationsIfNew(ctx, input.RecipientUIDs, ForumNotification{
		Type:           NotificationTypeRecruiting,
		Domain:         DomainCFB,
		LinkTo:         linkTo,
		Message:        input.Message,
		ActorUsername:  "SimSN",
		IsRead:         false,
		SourceEventKey: input.SourceEventKey,
	})
}

// NotifyPracticeSquadOffer notifies an NFL team's owner and GM that another team
// has placed a practice squad offer on one of their players.  Idempotent via SourceEventKey.
func NotifyPracticeSquadOffer(ctx context.Context, input PracticeSquadOfferNotificationInput) error {
	if len(input.RecipientUIDs) == 0 {
		return nil
	}

	message := fmt.Sprintf(
		"%s have placed an offer on %s %s to pick up from your practice squad.",
		input.OfferingTeam, input.Position, input.PlayerName,
	)
	linkTo := BuildTeamRosterRoute("nfl", input.OwnerTeamID)

	return writeNotificationsIfNew(ctx, input.RecipientUIDs, ForumNotification{
		Type:           NotificationTypeFreeAgency,
		Domain:         DomainNFL,
		LinkTo:         linkTo,
		Message:        message,
		ActorUsername:  "SimSN",
		IsRead:         false,
		SourceEventKey: input.SourceEventKey,
	})
}

// NotifyTransferIntention notifies a coach that one of their players has declared
// an intention to enter the transfer portal.  Idempotent via SourceEventKey.
func NotifyTransferIntention(ctx context.Context, input TransferIntentionNotificationInput) error {
	if len(input.RecipientUIDs) == 0 {
		return nil
	}

	message := fmt.Sprintf(
		"%d star %s %s has a %s likeliness of entering the transfer portal. Please navigate to the Roster page to submit a promise.",
		input.Stars, input.Position, input.PlayerName, input.TransferLikeliness,
	)
	linkTo := BuildTeamRosterRoute("cfb", input.TeamID)

	return writeNotificationsIfNew(ctx, input.RecipientUIDs, ForumNotification{
		Type:           NotificationTypeTransfer,
		Domain:         DomainCFB,
		LinkTo:         linkTo,
		Message:        message,
		ActorUsername:  "SimSN",
		IsRead:         false,
		SourceEventKey: input.SourceEventKey,
	})
}

// NotifyTeamInjury notifies a team's coaches or owners that a player was injured
// during a game. The link leads to the team's roster page.
// Idempotent via SourceEventKey (keyed per player per game).
func NotifyTeamInjury(ctx context.Context, input TeamInjuryNotificationInput) error {
	if len(input.RecipientUIDs) == 0 {
		return nil
	}

	weeksStr := "1 week"
	if input.WeeksOfRecovery > 1 {
		weeksStr = fmt.Sprintf("%d weeks", input.WeeksOfRecovery)
	}
	message := fmt.Sprintf(
		"%s (%s) suffered %s and is expected to miss %s. Check the roster for details.",
		input.PlayerName, input.Position, input.InjuryType, weeksStr,
	)
	linkTo := BuildTeamRosterRoute(input.League, input.TeamID)

	return writeNotificationsIfNew(ctx, input.RecipientUIDs, ForumNotification{
		Type:           NotificationTypeInjury,
		Domain:         input.Domain,
		LinkTo:         linkTo,
		Message:        message,
		ActorUsername:  "SimSN",
		IsRead:         false,
		SourceEventKey: input.SourceEventKey,
	})
}

// NotifyRecruitSigned creates one notification document per recipient when a
// recruit commits to a team.  Idempotent via SourceEventKey.
func NotifyRecruitSigned(ctx context.Context, input RecruitSignedNotificationInput) error {
	if len(input.RecipientUIDs) == 0 {
		return nil
	}

	message := fmt.Sprintf("%s has signed with %s.", input.RecruitName, input.TeamName)
	linkTo := BuildTeamRecruitingRoute(input.League, input.TeamID)

	return writeNotificationsIfNew(ctx, input.RecipientUIDs, ForumNotification{
		Type:           NotificationTypeRecruiting,
		Domain:         input.Domain,
		LinkTo:         linkTo,
		Message:        message,
		ActorUsername:  "SimSN",
		IsRead:         false,
		SourceEventKey: input.SourceEventKey,
	})
}

// ─────────────────────────────────────────────
// Internal helpers
// ─────────────────────────────────────────────

// writeNotificationsIfNew writes one notification document per recipient UID,
// skipping any recipient that already has a document with the same
// SourceEventKey (idempotency guard).
func writeNotificationsIfNew(
	ctx context.Context,
	recipientUIDs []string,
	template ForumNotification,
) error {
	client := GetFirestoreClient()
	col := client.Collection("notifications")
	now := time.Now().UTC()

	for _, uid := range recipientUIDs {
		if uid == "" {
			continue
		}

		// Idempotency: skip if this event was already delivered to this recipient.
		if template.SourceEventKey != "" {
			exists, err := notificationExists(ctx, col, uid, template.SourceEventKey)
			if err != nil {
				log.Printf("firebase: idempotency check failed for uid=%s key=%s: %v", uid, template.SourceEventKey, err)
			}
			if exists {
				continue
			}
		}

		ref := col.NewDoc()
		n := template
		n.ID = ref.ID
		n.UID = uid
		n.CreatedAt = now

		if _, err := ref.Set(ctx, n); err != nil {
			log.Printf("firebase: failed to write notification for uid=%s: %v", uid, err)
		}
	}

	return nil
}

// notificationExists returns true when a notification doc already exists for
// the given uid and sourceEventKey.
func notificationExists(
	ctx context.Context,
	col *firestore.CollectionRef,
	uid string,
	sourceEventKey string,
) (bool, error) {
	iter := col.
		Where("uid", "==", uid).
		Where("sourceEventKey", "==", sourceEventKey).
		Limit(1).
		Documents(ctx)
	defer iter.Stop()

	_, err := iter.Next()
	if err == iterator.Done {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
