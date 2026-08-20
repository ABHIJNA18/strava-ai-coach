// owns the coaching pipeline from Go’s side
//run the database.GetRecentActivitiesByType()->Map activities-->Create AnalyzeActivitiesRequest where you Call Python and Return summary

package coach

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ABHIJNA18/strava-ai-coach/internal/database"
	coachpb "github.com/ABHIJNA18/strava-ai-coach/proto/generated/coach"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	//get runs from the last 30 days
	since := time.Now().AddDate(0, 0, -30)
	runs, err := database.GetActivitiesByTypeSince(
		s.db,
		athleteID,
		"Run",
		since,
	)
	if err != nil {
		fmt.Println("Error getting runs from the last 30 days in service.go")
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
		if status.Code(err) == codes.Canceled {
			fmt.Println("AnalyzeActivities request was canceled by the HTTP client")
			return "", err
		}
		fmt.Println("Error getting response from AnalyzeActivities in service.go:", err)
		return "", err
	}

	fmt.Println("response.summary in service.go :", response.GetSummary())
	return response.GetSummary(), nil

}

// GenerateCoaching fetches the athlete's recent runs and requests personalized
// coaching from the Python Coach Service.
func (s *Service) GenerateCoaching(ctx context.Context, athleteID int64, goal string) (string, error) {
	since := time.Now().AddDate(0, 0, -60)

	runs, err := database.GetActivitiesByTypeSince(
		s.db,
		athleteID,
		"Run",
		since,
	)
	if err != nil {
		fmt.Println(
			"Error getting runs from the last 60 days:",
			err,
		)
		return "", err
	}

	if len(runs) == 0 {
		return "There isn't enough running data from the last 60 days to provide personalized coaching recommendations yet. Keep logging your runs and check back once more activity is available.", nil
	}

	request := &coachpb.GenerateCoachingRequest{
		AthleteId: athleteID,
		Goal:      goal,
		Activities: mapActivitiesToProto(
			runs,
		),
	}

	response, err := s.client.GenerateCoaching(
		ctx,
		request,
	)
	if err != nil {
		fmt.Println(
			"Error getting response from GenerateCoaching:",
			err,
		)
		return "", err
	}

	return response.GetCoaching(), nil
}
