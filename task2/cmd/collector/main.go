package main

import (
	"log"
	"net"
	"task2/internal/collector/adapter"
	"task2/internal/collector/handler"
	"task2/internal/collector/usecase"
	pb "task2/pkg/proto/v1"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func main() {

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
		log.Fatal("failed to serve: ", err)
	}
}
