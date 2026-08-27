package config

import (
	"repo-stat/platform/env"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
)

type App struct {
	AppName string `yaml:"app_name" env-default:"repo-stat-processor"`
}

type Services struct {
	Collector string `yaml:"collector" env:"COLLECTOR_ADDRESS" env-default:"localhost:8083"`
}

type DB struct {
	DSN string `yaml:"dsn" env:"DB_DSN" env-default:"postgres://user:password@localhost:5432/subscriptions_db?sslmode=disable"`
}

type Kafka struct {
	Broker       string `yaml:"broker" env:"KAFKA_BROKER" env-default:"kafka:9092"`
	TopicRequest string `yaml:"topic_request" env-default:"repo-fetch-request"`
	TopicResult  string `yaml:"topic_result" env-default:"repo-fetch-result"`
}

type Config struct {
	App      App               `yaml:"app"`
	Services Services          `yaml:"services"`
	GRPC     grpcserver.Config `yaml:"grpc"`
	Logger   logger.Config     `yaml:"logger"`
	DB       DB                `yaml:"db"`
	Kafka    Kafka             `yaml:"kafka"`
}

func MustLoad(path string) Config {
	var cfg Config
	env.MustLoad(path, &cfg)
	return cfg
}
