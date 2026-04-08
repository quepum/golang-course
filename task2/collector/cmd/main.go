package main

import (
	"log"
	"net"
	"task2/collector/internal/adapter"
	"task2/collector/internal/handler"
	"task2/collector/internal/usecase"
	pb "task2/pkg/proto/v1"
	"task2/pkg/utils"

	"google.golang.org/grpc"
)

func main() {
	utils.LoadEnv()

	port := utils.GetEnv(utils.EnvCollectorPort, utils.DefaultCollectorPort)

	repoAdapter := adapter.NewGitHubRepo()
	uc := usecase.NewRepoUseCase(repoAdapter)
	grpcHandler := handler.NewGrpcServer(uc)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to listen: ", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRepoServiceServer(grpcServer, grpcHandler)

	log.Println("Starting server on port", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve on %s: %v", port, err)
	}
}
