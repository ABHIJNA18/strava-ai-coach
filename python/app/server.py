import grpc
from concurrent import futures

from python.generated import coach_pb2
from python.generated import coach_pb2_grpc


class CoachService(coach_pb2_grpc.CoachServiceServicer):

    def Ping(self,request,context):
        return coach_pb2.PingResponse( # creating an instance of PingResponse class
            message="Hello from Python!"
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




