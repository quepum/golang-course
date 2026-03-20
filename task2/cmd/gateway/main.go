package main

import (
	"log"
	"net/http"
	"task2/internal/gateway/adapter"
	"task2/internal/gateway/handler"
	"task2/internal/gateway/usecase"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	collectorAddress = "localhost:50051"
	httpPort         = ":8080"
)

func main() {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(collectorAddress, opts...)
	if err != nil {
		log.Fatalf("failed to create gRPC client: %v", err)
	}

	defer conn.Close()
	log.Printf("Connected to collector at %s", collectorAddress)

	repoAdapter := adapter.NewCollectorClient(conn)
	uc := usecase.NewRepoUseCase(repoAdapter)
	httpHandler := handler.NewHTTPServer(uc)

	http.HandleFunc("/repos/", httpHandler.GetRepoInfo)
	log.Printf("Gateway starting on %s", httpPort)
	if err := http.ListenAndServe(httpPort, nil); err != nil {
		log.Fatal(err)
	}
}
