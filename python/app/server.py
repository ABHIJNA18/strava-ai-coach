import grpc
from concurrent import futures

from python.generated import coach_pb2
from python.generated import coach_pb2_grpc
from python.app.metrics.running import (
    calculate_run_count,
    calculate_total_distance,
    calculate_average_heartrate,
    calculate_average_pace,
)
def format_pace(seconds_per_km):
    if seconds_per_km <= 0:
        return "N/A"

    minutes = int(seconds_per_km // 60)
    seconds = int(round(seconds_per_km % 60))

    if seconds == 60:
        minutes += 1
        seconds = 0

    return f"{minutes}:{seconds:02d} min/km"


class CoachService(coach_pb2_grpc.CoachServiceServicer): 
    #CoachService is a child class of CoachServiceServicer

    def AnalyzeActivities(self, request, context):
        print("Caculating metrics in python server.py")
        run_count = calculate_run_count(request.activities)
        total_distance_meters = calculate_total_distance(request.activities)
        average_heartrate = calculate_average_heartrate(request.activities)
        average_pace = calculate_average_pace(request.activities)

        total_distance_km = total_distance_meters / 1000

        summary = (
            f"You completed {run_count} runs covering {total_distance_km:.1f} km.\n"
            f"Average heart rate was {average_heartrate:.0f} bpm.\n"
            f"Average pace was {format_pace(average_pace)}.\n"
            "Great consistency this week."
        )

        return coach_pb2.AnalyzeActivitiesResponse(
            summary=summary
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




