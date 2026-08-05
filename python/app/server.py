# Receives the request, calls the coach orchestrator, and returns the protobuf response.
# This file runs the Python gRPC Coach Service.
# It receives protobuf requests and returns protobuf responses.

import grpc
from concurrent import futures

from python.generated import coach_pb2
from python.generated import coach_pb2_grpc
from python.app.coach import generate_coaching_summary

class CoachService(coach_pb2_grpc.CoachServiceServicer): 
    #CoachService is a child class of CoachServiceServicer

    def AnalyzeActivities(self, request, context):
        print("Analysing activities in server.py")

        try:
            summary =generate_coaching_summary(request.activities)
            return coach_pb2.AnalyzeActivitiesResponse(
                summary=summary
            )
        except Exception as error:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(error))

            print("Error analysing activities in server.py")

            return coach_pb2.AnalyzeActivitiesResponse(
                summary=""
            )


def serve():

    server = grpc.server(
        futures.ThreadPoolExecutor(max_workers=10)
    )

    coach_pb2_grpc.add_CoachServiceServicer_to_server(
        CoachService(),
        server,
    )

    server.add_insecure_port("[::]:50051")
    print("Python Coach Service listening on port 50051...")
    server.start()
    server.wait_for_termination()
 

if __name__=="__main__":
    serve()




