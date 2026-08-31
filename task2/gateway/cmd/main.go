package main

import (
	"log"
	"net/http"
	_ "task2/docs/swagger"
	"task2/gateway/internal/adapter"
	"task2/gateway/internal/handler"
	"task2/gateway/internal/usecase"
	"task2/pkg/utils"

	httpSwagger "github.com/swaggo/http-swagger/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	utils.LoadEnv()

	port := utils.GetEnv(utils.EnvHTTPPort, utils.DefaultHTTPPort)
	collectorAddress := utils.GetEnv(utils.EnvCollectorAddr, utils.DefaultCollectorAddr)

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

	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	http.HandleFunc("/repos/", httpHandler.GetRepoInfo)

	log.Printf("Gateway starting on %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
