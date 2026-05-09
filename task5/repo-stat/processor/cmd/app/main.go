package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
	"repo-stat/processor/config"
	"repo-stat/processor/internal/adapter/kafka"
	"repo-stat/processor/internal/adapter/storage"
	grpccontroller "repo-stat/processor/internal/controller/grpc"
	"repo-stat/processor/internal/usecase"
	processorv1 "repo-stat/proto/processor"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
)

func run(ctx context.Context) error {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)
	log := logger.MustMakeLogger(cfg.Logger.LogLevel)
	log.Info("loaded config", "kafka_broker", cfg.Kafka.Broker)

	log.Info("starting processor server...", "addr", cfg.GRPC.Address)

	pool, err := pgxpool.New(ctx, cfg.DB.DSN)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer pool.Close()
	log.Info("connected to database")

	migrateURL := "file:///migrations"
	m, err := migrate.New(migrateURL, cfg.DB.DSN)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	log.Info("database migrations applied successfully")

	repoStorage := storage.NewRepository(pool)

	kafkaAdapter := kafka.NewAdapter(log, cfg.Kafka.Broker, cfg.Kafka.TopicRequest, cfg.Kafka.TopicResult)
	defer kafkaAdapter.Close()

	go func() {
		if err := kafkaAdapter.StartConsumingResults(ctx, repoStorage); err != nil {
			log.Error("kafka consumer stopped", "error", err)
		}
	}()

	pingUseCase := usecase.NewPing()
	repoInfoUseCase := usecase.NewRepoInfo(repoStorage, kafkaAdapter)

	server := grpccontroller.NewServer(log, pingUseCase, repoInfoUseCase)

	srv, err := grpcserver.New(cfg.GRPC.Address)
	if err != nil {
		return fmt.Errorf("create grpc server: %w", err)
	}

	processorv1.RegisterProcessorServer(srv.GRPC(), server)

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
