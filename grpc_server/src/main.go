package main

import (
	"github.com/maei/golang_grpc_bidirect_streaming/grpc_server/src/server"
	"github.com/maei/shared_utils_go/logger"
)

func main() {
	logger.Info("Starting GRPC-Server")
	server.StartGRPCServer()

}
