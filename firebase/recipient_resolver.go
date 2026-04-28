package firebase

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// ResolveUIDsByUsernames queries the Firestore "users" collection for each
// username and returns the matching document IDs (i.e. Firebase Auth UIDs).
// Usernames that are not found are silently skipped and a warning is logged.
func ResolveUIDsByUsernames(ctx context.Context, usernames []string) []string {
	if len(usernames) == 0 {
		return nil
	}

	client := GetFirestoreClient()
	uids := make([]string, 0, len(usernames))

	for _, username := range usernames {
		if username == "" {
			continue
		}
		uid, err := resolveUID(ctx, client, username)
		if err != nil {
			log.Printf("firebase: could not resolve UID for username %q: %v", username, err)
			continue
		}
		uids = append(uids, uid)
	}

	return uids
}

func resolveUID(ctx context.Context, client *firestore.Client, username string) (string, error) {
	iter := client.Collection("users").
		Where("username", "==", username).
		Limit(1).
		Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err == iterator.Done {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return doc.Ref.ID, nil
}

// GetAllUsers returns every document in the Firestore "users" collection as a
// slice of UserRecord.  Each record's UID field is set from the document ID.
// Errors on individual documents are logged and skipped.
func GetAllUsers(ctx context.Context) ([]UserRecord, error) {
	client := GetFirestoreClient()

	iter := client.Collection("users").Documents(ctx)
	defer iter.Stop()

	var users []UserRecord
	for {
		docSnap, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var rec UserRecord
		if err := docSnap.DataTo(&rec); err != nil {
			log.Printf("firebase: could not decode user document %s: %v", docSnap.Ref.ID, err)
			continue
		}
		rec.UID = docSnap.Ref.ID
		users = append(users, rec)
	}

	return users, nil
}

// ResetMediaPointsForCFBUsers queries all Firestore "users" documents where
// teamId > 0 (i.e. the user has a college football team) and sets the
// SimCFBMediaPoints field to 0.  Errors on individual documents are logged
// and skipped so that a single bad record cannot abort the whole sweep.
func ResetMediaPointsForCFBUsers(ctx context.Context) error {
	client := GetFirestoreClient()

	iter := client.Collection("users").
		Where("teamId", ">", 0).
		Documents(ctx)
	defer iter.Stop()

	for {
		docSnap, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		_, err = docSnap.Ref.Update(ctx, []firestore.Update{
			{Path: "SimCFBMediaPoints", Value: 0},
		})
		if err != nil {
			log.Printf("firebase: failed to reset SimCFBMediaPoints for user %s: %v", docSnap.Ref.ID, err)
		}
	}

	return nil
}
