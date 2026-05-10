package grpc

import (
	"context"
	"log/slog"
	"repo-stat/processor/internal/dto"
	"repo-stat/processor/internal/usecase"
	processorv1 "repo-stat/proto/processor"
)

type Server struct {
	processorv1.UnimplementedProcessorServer
	log      *slog.Logger
	ping     *usecase.Ping
	repoInfo *usecase.RepoInfo
}

func NewServer(log *slog.Logger, ping *usecase.Ping, repoInfo *usecase.RepoInfo) *Server {
	return &Server{
		log:      log,
		ping:     ping,
		repoInfo: repoInfo,
	}
}

func (s *Server) Ping(ctx context.Context, req *processorv1.PingRequest) (*processorv1.PingResponse, error) {
	s.log.Debug("Processor ping request received")
	return &processorv1.PingResponse{
		Reply: s.ping.Execute(ctx),
	}, nil
}

func (s *Server) GetRepoInfo(ctx context.Context, req *processorv1.RepoInfoRequest) (*processorv1.RepoInfoResponse, error) {
	s.log.Info("GetRepoInfo requested via Processor", "url", req.Url)

	repo, err := s.repoInfo.Execute(ctx, req.Url)
	if err != nil {
		return nil, err
	}

	repoDTO := dto.RepoResponse{
		FullName:    repo.FullName,
		Description: repo.Description,
		Stars:       int32(repo.Stars),
		Forks:       int32(repo.Forks),
		CreatedAt:   repo.CreatedAt,
	}

	return &processorv1.RepoInfoResponse{
		FullName:    repoDTO.FullName,
		Description: repoDTO.Description,
		Stars:       repoDTO.Stars,
		Forks:       repoDTO.Forks,
		CreatedAt:   repoDTO.CreatedAt,
	}, nil
}
