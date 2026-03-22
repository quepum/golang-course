package main

import (
	"log"
	"net/http"
	"os"
	_ "task2/docs/swagger"
	"task2/internal/gateway/adapter"
	"task2/internal/gateway/handler"
	"task2/internal/gateway/usecase"

	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	port := getEnv("HTTP_PORT", ":8080")
	collectorAddress := getEnv("COLLECTOR_ADDR", "localhost:50051")

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
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	log.Printf("Gateway starting on %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
