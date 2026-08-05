// create Client to interact with python servcie
package coach

import (

	//"context"
	//"fmt"
	//"log"
	coachpb "github.com/ABHIJNA18/strava-ai-coach/proto/generated/coach"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient(address string) (coachpb.CoachServiceClient, *grpc.ClientConn, error) {

	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, err
	}
	client := coachpb.NewCoachServiceClient(conn)
	return client, conn, err
}
