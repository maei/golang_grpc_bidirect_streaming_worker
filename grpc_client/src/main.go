package main

import (
	"github.com/maei/golang_grpc_bidirect_streaming/grpc_client/src/service"
	"github.com/maei/shared_utils_go/logger"
	"time"
)

func main() {
	logger.Info("Starting GRPC-Client")
	time.Sleep(4 * time.Second)
	service.GreetService.Greeting()

}
