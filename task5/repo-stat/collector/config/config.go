package config

import (
	"repo-stat/platform/env"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
)

type App struct {
	AppName string `yaml:"app_name" env-default:"repo-stat-collector"`
}

type Services struct {
	Subscriber string `yaml:"subscriber" env:"SUBSCRIBER_ADDRESS" env-default:"subscriber:8080"`
}

type Kafka struct {
	Broker       string `yaml:"broker" env:"KAFKA_BROKER" env-default:"kafka:9092"`
	TopicRequest string `yaml:"topic_request" env-default:"repo-fetch-request"`
	TopicResult  string `yaml:"topic_result" env-default:"repo-fetch-result"`
	GroupID      string `yaml:"group_id" env-default:"collector-group"`
}

type Config struct {
	App      App               `yaml:"app"`
	Services Services          `yaml:"services"`
	GRPC     grpcserver.Config `yaml:"grpc"`
	Logger   logger.Config     `yaml:"logger"`
	Kafka    Kafka             `yaml:"kafka"`
}

func MustLoad(path string) Config {
	var cfg Config
	env.MustLoad(path, &cfg)
	return cfg
}
