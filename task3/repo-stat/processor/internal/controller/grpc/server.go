package grpc

import (
	"context"
	"log/slog"
	"repo-stat/processor/internal/usecase"
	processorv1 "repo-stat/proto/processor"
)

type Server struct {
	processorv1.UnimplementedProcessorServer
	log  *slog.Logger
	ping *usecase.Ping
}

func NewServer(log *slog.Logger, ping *usecase.Ping) *Server {
	return &Server{
		log:  log,
		ping: ping,
	}
}

func (s *Server) Ping(ctx context.Context, req *processorv1.PingRequest) (*processorv1.PingResponse, error) {
	s.log.Debug("Processor ping request received")
	return &processorv1.PingResponse{
		Reply: s.ping.Execute(ctx),
	}, nil
}

func (s *Server) GetRepoInfo(ctx context.Context, req *processorv1.RepoInfoRequest) (*processorv1.RepoInfoResponse, error) {
	return nil, nil
}
