package dto

import "repo-stat/api/internal/domain"

type PingResponse struct {
	Status   string                     `json:"status"`
	Services []domain.ServicePingResult `json:"services"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
