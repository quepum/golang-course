package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
	subscriberpb "repo-stat/proto/subscriber"
	"repo-stat/subscriber/config"
	"repo-stat/subscriber/internal/adapter/github"
	"repo-stat/subscriber/internal/adapter/storage"
	grpccontroller "repo-stat/subscriber/internal/controller/grpc"
	"repo-stat/subscriber/internal/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
)

func run(ctx context.Context) error {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	log := logger.MustMakeLogger(cfg.Logger.LogLevel)
	log.Info("starting subscriber server...")
	log.Debug("debug messages are enabled")

	log.Info("connecting to database...")
	pool, err := pgxpool.New(ctx, cfg.DB.DSN)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer pool.Close()
	log.Info("connected to database")

	initSQL, err := os.ReadFile("/migrations/init.sql")
	if err != nil {
		return fmt.Errorf("failed to read init.sql from /migrations/: %w", err)
	}

	_, err = pool.Exec(ctx, string(initSQL))
	if err != nil {
		return fmt.Errorf("failed to execute init.sql: %w", err)
	}
	log.Info("database schema initialized from init.sql")

	queries := storage.New(pool)
	repo := storage.NewRepository(queries)

	ghClient := github.NewClient()

	pingUseCase := usecase.NewPing()
	subscriptionUseCase := usecase.NewSubscriptionManager(repo, ghClient)

	server := grpccontroller.NewServer(log, pingUseCase, subscriptionUseCase)

	srv, err := grpcserver.New(cfg.GRPC.Address)
	if err != nil {
		return fmt.Errorf("create grpc server: %w", err)
	}

	subscriberpb.RegisterSubscriberServer(srv.GRPC(), server)

	if err := srv.Run(ctx); err != nil {
		return fmt.Errorf("run grpc server: %w", err)
	}

	return nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	if err := run(ctx); err != nil {
		_, err = fmt.Fprintln(os.Stderr, err)
		if err != nil {
			fmt.Printf("launching server error: %s\n", err)
		}
		cancel()
		os.Exit(1)
	}
}
