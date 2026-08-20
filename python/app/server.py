# Receives the request, calls the coach orchestrator, and returns the protobuf response.
# This file runs the Python gRPC Coach Service.
# It receives protobuf requests and returns protobuf responses.

import json
import grpc

from concurrent import futures

from python.generated import coach_pb2
from python.generated import coach_pb2_grpc
from python.app.coach import (
    generate_coaching_summary,
    generate_personalized_coaching,
)


class CoachService(
    coach_pb2_grpc.CoachServiceServicer
):
    # Handles the existing 30-day summary request.
    def AnalyzeActivities(self,request,context,):

        print("Analysing activities in server.py")

        try:
            summary = generate_coaching_summary(request.activities)

            return coach_pb2.AnalyzeActivitiesResponse(summary=summary)

        except Exception as error:

            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(error))

            print("Error analysing activities:",error,)

            return coach_pb2.AnalyzeActivitiesResponse(summary="")

    # Handles personalized goal-based coaching.
    def GenerateCoaching(self,request,context,):

        print("Generating personalized coaching in server.py")

        try:
            coaching = generate_personalized_coaching(request.goal,request.activities,)

            return coach_pb2.GenerateCoachingResponse(
                coaching=json.dumps(
                    coaching,
                    ensure_ascii=True,
                )
            )

        except Exception as error:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(error))

            print(
                "Error generating personalized coaching:",
                repr(error),
            )

            return coach_pb2.GenerateCoachingResponse(
                coaching=""
            )

def serve():
    server = grpc.server(
        futures.ThreadPoolExecutor(
            max_workers=10
        )
    )

    coach_pb2_grpc.add_CoachServiceServicer_to_server(
        CoachService(),
        server,
    )

    server.add_insecure_port(
        "[::]:50051"
    )

    print(
        "Python Coach Service listening on port 50051..."
    )

    server.start()
    server.wait_for_termination()


if __name__ == "__main__":
    serve()