package usecase

import (
	"context"
	"repo-stat/api/internal/domain"
)

type Pinger interface {
	Ping(ctx context.Context) domain.PingStatus
	GetName() string
}

type Ping struct {
	pingers []Pinger
}

func NewPing(pingers ...Pinger) *Ping {
	return &Ping{
		pingers: pingers,
	}
}

func (u *Ping) Execute(ctx context.Context) domain.PingResult {
	var services []domain.ServicePingResult
	allUp := true

	for _, pinger := range u.pingers {
		status := pinger.Ping(ctx)
		services = append(services, domain.ServicePingResult{
			Name:   pinger.GetName(),
			Status: status,
		})

		if status != domain.PingStatusUp {
			allUp = false
		}
	}

	overallStatus := "ok"
	if !allUp {
		overallStatus = "fail"
	}

	return domain.PingResult{
		Status:   overallStatus,
		Services: services,
	}
}
