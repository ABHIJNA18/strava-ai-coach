//Concert all databse activities into protobuf format

package coach

import (
	"time"

	"github.com/ABHIJNA18/strava-ai-coach/internal/database"
	coachpb "github.com/ABHIJNA18/strava-ai-coach/proto/generated/coach"
)

func mapActivitiesToProto(activities []database.Activity) []*coachpb.Activity {

	protoActivities := make([]*coachpb.Activity, 0, len(activities))

	for _, activity := range activities {
		protoActivities = append(protoActivities, &coachpb.Activity{
			Id:                 activity.StravaActivityID,
			Name:               activity.Name,
			SportType:          activity.SportType,
			Distance:           float32(activity.Distance),
			MovingTime:         int32(activity.MovingTime),
			AverageSpeed:       float32(activity.AverageSpeed),
			AverageHeartrate:   float32(activity.AverageHeartrate),
			MaxHeartrate:       float32(activity.MaxHeartrate),
			AverageCadence:     float32(activity.AverageCadence),
			TotalElevationGain: int32(activity.TotalElevationGain),
			StartDate:          activity.StartDate.Format(time.RFC3339),
		})

	}
	return protoActivities

}
