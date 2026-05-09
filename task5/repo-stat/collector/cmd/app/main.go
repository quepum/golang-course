package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"repo-stat/collector/config"
	"repo-stat/collector/internal/adapter/github"
	"repo-stat/collector/internal/adapter/kafka"
	"repo-stat/collector/internal/adapter/subscriber"
	grpccontroller "repo-stat/collector/internal/controller/grpc"
	"repo-stat/collector/internal/usecase"
	"repo-stat/platform/grpcserver"
	platformKafka "repo-stat/platform/kafka"
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

	kafkaAdapter := kafka.NewAdapter(log, cfg.Kafka.Broker, cfg.Kafka.TopicRequest, cfg.Kafka.TopicResult, cfg.Kafka.GroupID)
	defer kafkaAdapter.Close()

	repoInfoUseCase := usecase.NewRepoInfo(githubClient, subClient, log)

	go repoInfoUseCase.StartBackgroundUpdate(ctx, kafkaAdapter)
	go func() {
		log.Info("collector listening for kafka tasks...")
		err := kafkaAdapter.StartListening(ctx, func(req platformKafka.FetchRequest) {
			log.Debug("processing task", "owner", req.Owner, "repo", req.Repo)

			repo, err := repoInfoUseCase.GetRepoInfo(ctx, req.Owner, req.Repo)

			res := platformKafka.FetchResult{Owner: req.Owner, Repo: req.Repo}
			if err != nil {
				res.Error = err.Error()
			} else {
				res.FullName = repo.FullName
				res.Stars = repo.Stars
				res.Forks = repo.Forks
				res.Description = repo.Description
				res.CreatedAt = repo.CreatedAt
			}

			if err := kafkaAdapter.SendFetchResult(ctx, res); err != nil {
				log.Error("failed to send result to kafka", "error", err)
			}
		})
		if err != nil && !errors.Is(err, context.Canceled) {
			log.Error("kafka consumer fatal error", "error", err)
		}
	}()

	pingUseCase := usecase.NewPing()
	server := grpccontroller.NewServer(log, pingUseCase)
	srv, err := grpcserver.New(cfg.GRPC.Address)
	if err != nil {
		return fmt.Errorf("create grpc server: %w", err)
	}

	collectorv1.RegisterCollectorServer(srv.GRPC(), server)

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
