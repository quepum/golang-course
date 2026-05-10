package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
	"repo-stat/processor/config"
	"repo-stat/processor/internal/adapter/collector"
	grpccontroller "repo-stat/processor/internal/controller/grpc"
	"repo-stat/processor/internal/usecase"
	processorv1 "repo-stat/proto/processor"
)

func run(ctx context.Context) error {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	log := logger.MustMakeLogger(cfg.Logger.LogLevel)
	log.Info("starting processor server...", "addr", cfg.GRPC.Address)
	log.Debug("debug messages are enabled")

	pingUseCase := usecase.NewPing()

	collectorClient, err := collector.NewClient(cfg.Services.Collector, log)
	if err != nil {
		log.Error("failed to create collector client", "error", err)
		return fmt.Errorf("create collector client: %w", err)
	}

	repoInfoUseCase := usecase.NewRepoInfo(collectorClient)

	serverHandler := grpccontroller.NewServer(log, pingUseCase, repoInfoUseCase)

	srv, err := grpcserver.New(cfg.GRPC.Address)
	if err != nil {
		return fmt.Errorf("create grpc server: %w", err)
	}

	processorv1.RegisterProcessorServer(srv.GRPC(), serverHandler)

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
