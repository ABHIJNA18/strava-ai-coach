// This file defines the JSON model used for Strava webhook events.

package webhook

type Event struct {
	ObjectType     string                 `json:"object_type"`
	ObjectID       int64                  `json:"object_id"`
	AspectType     string                 `json:"aspect_type"`
	OwnerID        int64                  `json:"owner_id"`
	SubscriptionID int64                  `json:"subscription_id"`
	EventTime      int64                  `json:"event_time"`
	Updates        map[string]interface{} `json:"updates"`
}