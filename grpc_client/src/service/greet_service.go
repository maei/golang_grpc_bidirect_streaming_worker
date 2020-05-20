package service

import (
	"context"
	"fmt"
	"github.com/maei/golang_grpc_bidirect_streaming/grpc_client/src/client"
	"github.com/maei/golang_grpc_bidirect_streaming/grpc_client/src/domain/greetpb"
	"github.com/maei/shared_utils_go/logger"
	"io"
	"time"
)

var GreetService greetServiceInterface = &greetService{}

type greetServiceInterface interface {
	Greeting()
}

type greetService struct{}

func generateName(name string) *greetpb.GreetRequest {
	res := &greetpb.GreetRequest{
		Greet: &greetpb.Greeting{FirstName: name},
	}
	return res
}

func (*greetService) Greeting() {
	logger.Info("Start GRPC-call")

	conn, clientErr := client.GRPCClient.SetClient()
	if clientErr != nil {
		logger.Error("Error while creating GRPC-Client", clientErr)
	}

	c := greetpb.NewGreetServiceClient(conn)

	stream, streamErr := c.GetGreeting(context.Background())
	if streamErr != nil {
		logger.Error("Error while streaming to GRPC-Server", streamErr)
	}

	// create some random names
	names := []string{"Matthias", "Sonia", "Heidi", "Jochen"}

	// create a channel of type struct
	//waitc := make(chan struct{})
	waitc := make(chan bool)

	// send messages to gRPC-Server in a go-routine
	go func() {
		for _, name := range names {
			res := generateName(name)
			logger.Info(fmt.Sprintf("Sending Request %v to gRPC Server", res))
			stream.Send(res)
			time.Sleep(2 * time.Second)
		}
		stream.CloseSend()
	}()

	// receive messages from the gRPC-Server in a go-routine
	go func() {
		for {
			res, resErr := stream.Recv()
			if resErr == io.EOF {
				break
			}
			if resErr != nil {
				logger.Error("Error while receiving Data from gRPC-Server", resErr)
				break
			}
			fmt.Println(fmt.Sprintf("Result from gRPC-Server: %v", res.GetResult()))
		}
		waitc <- true
		//close(waitc)

	}()
	//block until everything is done
	<-waitc
}
