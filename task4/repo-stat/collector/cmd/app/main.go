package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"repo-stat/collector/config"
	"repo-stat/collector/internal/adapter/github"
	"repo-stat/collector/internal/adapter/subscriber"
	grpccontroller "repo-stat/collector/internal/controller/grpc"
	"repo-stat/collector/internal/usecase"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
	collectorv1 "repo-stat/proto/collector"
)

func run(ctx context.Context) error {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)
	log := logger.MustMakeLogger(cfg.Logger.LogLevel)

	githubClient := github.NewClient()

	subClient, err := subscriber.NewClient(cfg.Services.Subscriber, log)
	if err != nil {
		return fmt.Errorf("create subscriber client: %w", err)
	}
	defer subClient.Close()

	repoInfoUseCase := usecase.NewRepoInfo(githubClient, subClient, log)
	pingUseCase := usecase.NewPing()

	collectorServer := grpccontroller.NewServer(log, pingUseCase, repoInfoUseCase)

	srv, err := grpcserver.New(cfg.GRPC.Address)
	if err != nil {
		return fmt.Errorf("create grpc server: %w", err)
	}

	collectorv1.RegisterCollectorServer(srv.GRPC(), collectorServer)

	log.Info("starting collector server...", "addr", cfg.GRPC.Address)
	log.Debug("debug messages are enabled")

	if err := srv.Run(ctx); err != nil {
		return fmt.Errorf("run grpc server: %w", err)
	}
	return nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)

	if err := run(ctx); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		cancel()
		os.Exit(1)
	}
}
