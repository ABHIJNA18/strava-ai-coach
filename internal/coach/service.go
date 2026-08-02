// owns the coaching pipeline from Go’s side
//run the database.GetRecentActivitiesByType()->Map activities-->Create AnalyzeActivitiesRequest where you Call Python and Return summary

package coach

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ABHIJNA18/strava-ai-coach/internal/database"
	coachpb "github.com/ABHIJNA18/strava-ai-coach/proto/generated/coach"
)

type Service struct {
	db     *sql.DB
	client coachpb.CoachServiceClient
}

func NewService(db *sql.DB, client coachpb.CoachServiceClient) *Service {
	return &Service{
		db:     db,
		client: client,
	}
}

//Define a method attached to Service, later it will grow to AnalyseMonthlyRuns() etc

func (s *Service) AnalyzeRecentRuns(ctx context.Context, athleteID int64) (string, error) {

	//get runs
	runs, err := database.GetRecentActivitiesByType(
		s.db,
		athleteID,
		"Run",
		10,
	)
	if err != nil {
		fmt.Println("Error getting recent runs in service.go")
		return "", err
	}

	//create a request while mapping the runs
	//defined in coach.proto
	request := &coachpb.AnalyzeActivitiesRequest{
		AthleteId:  athleteID,
		Activities: mapActivitiesToProto(runs),
	}

	//make the request
	response, err := s.client.AnalyzeActivities(ctx, request)
	if err != nil {
		fmt.Println("Error getting response from AnalyzeActivities in service.go")
		return "", err
	}

	fmt.Println("response.summary in service.go :", response.GetSummary())
	return response.GetSummary(), nil

}
