package grpc

import (
	"context"
	"log/slog"
	subscriberpb "repo-stat/proto/subscriber"
	"repo-stat/subscriber/internal/domain"
	"repo-stat/subscriber/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	subscriberpb.UnimplementedSubscriberServer
	log                 *slog.Logger
	ping                *usecase.Ping
	subscriptionManager *usecase.SubscriptionManager
}

func NewServer(log *slog.Logger, ping *usecase.Ping, subscriptionManager *usecase.SubscriptionManager) *Server {
	return &Server{
		log:                 log,
		ping:                ping,
		subscriptionManager: subscriptionManager,
	}
}

func (s *Server) Ping(ctx context.Context, _ *subscriberpb.PingRequest) (*subscriberpb.PingResponse, error) {
	s.log.Debug("subscriber ping request received") // Исправил опечатку subscriberp -> subscriber
	return &subscriberpb.PingResponse{
		Reply: s.ping.Execute(ctx),
	}, nil
}

func (s *Server) AddSubscription(ctx context.Context, req *subscriberpb.AddSubscriptionRequest) (*subscriberpb.AddSubscriptionResponse, error) {
	domainSub := &domain.Subscription{
		Owner: req.Owner,
		Repo:  req.Repo,
	}

	err := s.subscriptionManager.AddSubscription(ctx, domainSub)

	if err != nil {
		if err.Error() == "repository does not exist on GitHub" {
			return &subscriberpb.AddSubscriptionResponse{
				Status:  false,
				Message: err.Error(),
			}, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &subscriberpb.AddSubscriptionResponse{
		Status:  true,
		Message: "subscription added successfully",
	}, nil
}

func (s *Server) RemoveSubscription(ctx context.Context, req *subscriberpb.RemoveSubscriptionRequest) (*subscriberpb.RemoveSubscriptionResponse, error) {
	err := s.subscriptionManager.RemoveSubscription(ctx, req.Owner, req.Repo)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &subscriberpb.RemoveSubscriptionResponse{
		Status: true,
	}, nil
}

func (s *Server) ListSubscriptions(ctx context.Context, req *subscriberpb.ListSubscriptionsRequest) (*subscriberpb.ListSubscriptionsResponse, error) {
	domainSubs, err := s.subscriptionManager.ListSubscriptions(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var protoSubs []*subscriberpb.Subscription
	for _, sub := range domainSubs {
		protoSubs = append(protoSubs, &subscriberpb.Subscription{
			Owner: sub.Owner,
			Repo:  sub.Repo,
		})
	}

	return &subscriberpb.ListSubscriptionsResponse{
		Subscriptions: protoSubs,
	}, nil
}
