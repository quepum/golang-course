package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"repo-stat/collector/config"
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
	log.Info("starting collector server...", "addr", cfg.GRPC.Address)
	log.Debug("debug messages are enabled")

	pingUseCase := usecase.NewPing()
	collectorServer := grpccontroller.NewServer(log, pingUseCase)

	srv, err := grpcserver.New(cfg.GRPC.Address)
	if err != nil {
		return fmt.Errorf("create grpc server: %w", err)
	}

	collectorv1.RegisterCollectorServer(srv.GRPC(), collectorServer)

	if err := srv.Run(ctx); err != nil {
		return fmt.Errorf("run grpc server: %w", err)
	}
	return nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	if err := run(ctx); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
