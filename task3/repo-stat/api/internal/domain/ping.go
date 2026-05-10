package domain

type PingStatus string

const (
	PingStatusUp   PingStatus = "up"
	PingStatusDown PingStatus = "down"
)

type ServicePingResult struct {
	Name   string     `json:"name"`
	Status PingStatus `json:"status"`
}

type PingResult struct {
	Status   string              `json:"status"`
	Services []ServicePingResult `json:"services"`
}
