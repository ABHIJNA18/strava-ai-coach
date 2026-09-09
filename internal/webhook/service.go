// This file processes validated Strava webhook events asynchronously.
//Unsupported but structurally valid events are ignored instead of causing Strava retries
// Supported events are activity create/update/delete and athlete deauthorization. Activity create/update events trigger a Strava API call to fetch the activity details and store them in PostgreSQL. Activity delete events remove the activity from PostgreSQL. Athlete deauthorization events remove the OAuth token and revoke all sessions for the athlete.

package webhook

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/ABHIJNA18/strava-ai-coach/internal/database"
	"github.com/ABHIJNA18/strava-ai-coach/internal/strava"
)

type WebhookService struct {
	db          *sql.DB
	syncService *strava.SyncService
}

func NewWebhookService(
	db *sql.DB,
	syncService *strava.SyncService,
) *WebhookService {
	return &WebhookService{
		db:          db,
		syncService: syncService,
	}
}

func (s *WebhookService) ProcessEvent(event Event) {
	switch {
	case event.ObjectType == "activity" &&
		(event.AspectType == "create" ||
			event.AspectType == "update"):
		s.processActivityUpsert(event)

	case event.ObjectType == "activity" &&
		event.AspectType == "delete":
		s.processActivityDelete(event)

	case event.ObjectType == "athlete" &&
		event.AspectType == "update" &&
		isDeauthorization(event):
		s.processDeauthorization(event)

	default:
		log.Printf(
			"Ignoring unsupported webhook event: object_type=%s aspect_type=%s",
			event.ObjectType,
			event.AspectType,
		)
	}
}

func (s *WebhookService) processActivityUpsert(event Event) {
	localAthleteID, err := database.GetAthleteByStravaAthleteID(
		s.db,
		event.OwnerID,
	)
	if err != nil {
		log.Printf(
			"Could not find local athlete for Strava owner %d: %v",
			event.OwnerID,
			err,
		)
		return
	}

	err = s.syncService.SyncActivity(
		localAthleteID,
		event.ObjectID,
	)
	if err != nil {
		log.Printf(
			"Failed to synchronize Strava activity %d: %v",
			event.ObjectID,
			err,
		)
		return
	}

	log.Printf(
		"Activity %d synchronized for local athlete %d",
		event.ObjectID,
		localAthleteID,
	)
}

func (s *WebhookService) processActivityDelete(event Event) {
	localAthleteID, err := database.GetAthleteByStravaAthleteID(
		s.db,
		event.OwnerID,
	)
	if err != nil {
		log.Printf(
			"Could not find local athlete for Strava owner %d: %v",
			event.OwnerID,
			err,
		)
		return
	}

	err = database.DeleteActivityByStravaID(
		s.db,
		localAthleteID,
		event.ObjectID,
	)
	if err != nil {
		log.Printf(
			"Failed to delete activity %d: %v",
			event.ObjectID,
			err,
		)
		return
	}

	log.Printf(
		"Activity %d deleted for local athlete %d",
		event.ObjectID,
		localAthleteID,
	)
}

func (s *WebhookService) processDeauthorization(event Event) {
	localAthleteID, err := database.GetAthleteByStravaAthleteID(
		s.db,
		event.OwnerID,
	)
	if err != nil {
		log.Printf(
			"Could not find local athlete for Strava owner %d: %v",
			event.OwnerID,
			err,
		)
		return
	}

	if err := database.DeleteOAuthTokenByAthleteID(
		s.db,
		localAthleteID,
	); err != nil {
		log.Printf(
			"Failed to delete OAuth token for athlete %d: %v",
			localAthleteID,
			err,
		)
		return
	}

	if err := database.RevokeAllSessionsForAthlete(
		s.db,
		localAthleteID,
	); err != nil {
		log.Printf(
			"Failed to revoke sessions for athlete %d: %v",
			localAthleteID,
			err,
		)
		return
	}

	log.Printf(
		"Strava authorization revoked for local athlete %d",
		localAthleteID,
	)
}

/*Look for "authorized" in the event.

If it doesn't exist:

    → not deauthorization

If it exists and says false:

    → deauthorization

Otherwise:

    → not deauthorization
*/

func isDeauthorization(event Event) bool {

	//if its deauthorisation, Updates map will have a key "authorized" with value false or "false"
	//if the key which is "authorized" exists, then continue,if not, then its not deauthorisation
	value, exists := event.Updates["authorized"]

	//if key doesn't exist, not deauth
	if !exists {
		return false //this isn’t deauthorization
	}

	// Because value is an interface{}, Go doesn’t know what type it contains
	switch authorized := value.(type) {
	case bool:
		return !authorized // if auth = true, not deauthorisation and if auth = false, then it is deauthorisation

	case string:
		parsed, err := strconv.ParseBool(authorized)
		return err == nil && !parsed

	default:
		return false
	}
}

func validateEvent(event Event) error {
	if event.ObjectType == "" {
		return fmt.Errorf("missing object_type")
	}

	if event.AspectType == "" {
		return fmt.Errorf("missing aspect_type")
	}

	if event.ObjectID <= 0 {
		return fmt.Errorf("object_id must be positive")
	}

	if event.OwnerID <= 0 {
		return fmt.Errorf("owner_id must be positive")
	}

	return nil
}

func isSupportedEvent(event Event) bool {
	if event.ObjectType == "activity" {
		return event.AspectType == "create" ||
			event.AspectType == "update" ||
			event.AspectType == "delete"
	}

	if event.ObjectType == "athlete" &&
		event.AspectType == "update" {
		return isDeauthorization(event)
	}

	return false
}