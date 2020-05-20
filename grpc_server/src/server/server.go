package server

import (
	"github.com/maei/golang_grpc_bidirect_streaming/grpc_server/src/domain/greetpb"
	"github.com/maei/shared_utils_go/logger"
	"google.golang.org/grpc"
	"io"
	"net"
	"time"
)

type server struct{}

var (
	s = grpc.NewServer()
)

func streamWorker() {

}

func (*server) GetGreeting(stream greetpb.GreetService_GetGreetingServer) error {
	logger.Info("gRPC greet-streaming started")

	// jobs channel to distribute incoming work
	jobs := make(chan string)
	// channel to close the jobs channel
	done := make(chan bool)

	// go routine to receive work from grpc-client
	go func() {
		for {
			req, reqErr := stream.Recv()
			if reqErr == io.EOF {
				break
			}
			if reqErr != nil {
				logger.Error("Error while fetich GRPC-Client request", reqErr)
				break
			}
			time.Sleep(5 * time.Second)
			jobs <- req.GetGreet().GetFirstName()
		}
		close(jobs)

	}()

	// go routine to catch work from job channel and send it back to grpc-client
	go func() {
		/*		for {
				j, more := <-jobs
				if more {
					log.Printf("worker gets string %s", j)
					res := &greetpb.GreetResponse{
						Result: fmt.Sprintf("Hello %v from the GRPC-Server", j),
					}
					streamErr := stream.Send(res)
					if streamErr != nil {
						logger.Error("Error while streaming data to GRPC-Client", streamErr)
						break
					}
				} else {
					log.Println("received all jobs")
					done <- true
					return
				}
			}*/
	}()

	logger.Info("gRPC-Server: All jobs done")
	<-done
	return nil

}

func StartGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Error("error while listening gRPC Server", err)
	}

	greetpb.RegisterGreetServiceServer(s, &server{})

	errServer := s.Serve(lis)
	if errServer != nil {
		logger.Error("error while serve gRPC Server", errServer)
	}
}
