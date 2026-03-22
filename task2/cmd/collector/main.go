package main

import (
	"log"
	"net"
	"os"
	"task2/internal/collector/adapter"
	"task2/internal/collector/handler"
	"task2/internal/collector/usecase"
	pb "task2/pkg/proto/v1"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables or defaults")
	}

	port := getEnv("COLLECTOR_PORT", ":50051")

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
