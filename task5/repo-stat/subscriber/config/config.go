package config

import (
	"repo-stat/platform/env"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
)

type App struct {
	AppName string `yaml:"app_name" env:"APP_NAME" env-default:"repo-stat-subscriber"`
}

type Services struct {
	API string `yaml:"api" env:"API_ADDRESS" env-default:"localhost:8080"`
}

type DB struct {
	DSN string `yaml:"dsn" env:"DB_DSN" env-default:"postgres://user:password@localhost:5432/subscriptions_db?sslmode=disable"`
}

type Config struct {
	App      App               `yaml:"app"`
	Services Services          `yaml:"services"`
	GRPC     grpcserver.Config `yaml:"grpc"`
	Logger   logger.Config     `yaml:"logger"`
	DB       DB                `yaml:"db"`
}

func MustLoad(path string) Config {
	var cfg Config
	env.MustLoad(path, &cfg)
	return cfg
}
