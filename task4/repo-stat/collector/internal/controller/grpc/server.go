package grpc

import (
	"context"
	"log/slog"
	"repo-stat/collector/internal/dto"
	"repo-stat/collector/internal/usecase"
	collectorv1 "repo-stat/proto/collector"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	repo, err := s.repoInfo.GetRepoInfo(ctx, req.Owner, req.Repo)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "repository not found") {
			return nil, status.Error(codes.NotFound, errMsg)
		}
		if strings.Contains(errMsg, "rate limit") {
			return nil, status.Error(codes.ResourceExhausted, errMsg)
		}
		return nil, status.Error(codes.Internal, errMsg)
	}

	repoDTO := dto.RepoDTO{
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

func (s *Server) GetSubscribedStats(ctx context.Context, req *collectorv1.SubscribedStatsRequest) (*collectorv1.SubscribedStatsResponse, error) {
	s.log.Info("GetSubscribedStats requested")

	domainRepos, err := s.repoInfo.GetSubscribedStats(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var dtos []*dto.RepoDTO
	for _, repo := range domainRepos {
		dtos = append(dtos, &dto.RepoDTO{
			FullName:    repo.FullName,
			Description: repo.Description,
			Stars:       int32(repo.Stars),
			Forks:       int32(repo.Forks),
			CreatedAt:   repo.CreatedAt,
		})
	}

	var protoStats []*collectorv1.RepoInfoResponse
	for _, repo := range dtos {
		protoStats = append(protoStats, &collectorv1.RepoInfoResponse{
			FullName:    repo.FullName,
			Description: repo.Description,
			Stars:       repo.Stars,
			Forks:       repo.Forks,
			CreatedAt:   repo.CreatedAt,
		})
	}

	return &collectorv1.SubscribedStatsResponse{
		Stats: protoStats,
	}, nil
}
