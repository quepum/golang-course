package grpc

import (
	"context"
	"log/slog"
	"repo-stat/collector/internal/usecase"
	collectorv1 "repo-stat/proto/collector"
)

type Server struct {
	collectorv1.UnimplementedCollectorServer
	log  *slog.Logger
	ping *usecase.Ping
}

func NewServer(log *slog.Logger, ping *usecase.Ping) *Server {
	return &Server{
		log:  log,
		ping: ping,
	}
}

func (s *Server) Ping(ctx context.Context, req *collectorv1.PingRequest) (*collectorv1.PingResponse, error) {
	s.log.Debug("Collector ping request received")
	return &collectorv1.PingResponse{
		Reply: s.ping.Execute(ctx),
	}, nil
}

func (s *Server) GetRepoInfo(ctx context.Context, req *collectorv1.RepoInfoRequest) (*collectorv1.RepoInfoResponse, error) {
	s.log.Info("GetRepoInfo requested", "owner", req.Owner, "repo", req.Repo)
	return nil, nil
}
