package grpc

import (
	"context"
	"log/slog"
	"repo-stat/collector/internal/dto"
	"repo-stat/collector/internal/usecase"
	collectorv1 "repo-stat/proto/collector"
)

type Server struct {
	collectorv1.UnimplementedCollectorServer
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

func (s *Server) Ping(ctx context.Context, req *collectorv1.PingRequest) (*collectorv1.PingResponse, error) {
	s.log.Debug("Collector ping request received")
	return &collectorv1.PingResponse{
		Reply: s.ping.Execute(ctx),
	}, nil
}

func (s *Server) GetRepoInfo(ctx context.Context, req *collectorv1.RepoInfoRequest) (*collectorv1.RepoInfoResponse, error) {
	s.log.Info("GetRepoInfo requested", "owner", req.Owner, "repo", req.Repo)
	repo, err := s.repoInfo.Execute(ctx, req.Owner, req.Repo)
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

	return &collectorv1.RepoInfoResponse{
		FullName:    repoDTO.FullName,
		Description: repoDTO.Description,
		Stars:       repoDTO.Stars,
		Forks:       repoDTO.Forks,
		CreatedAt:   repoDTO.CreatedAt,
	}, nil
}
