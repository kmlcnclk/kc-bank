package healthcheck

import "context"

type HealthCheckRequest struct {
}

type HealthCheckResponse struct {
	Status string `json:"status"`
}

type HealthCheckHandler struct {
}

func NewHealthCheckHandler() *HealthCheckHandler {
	return &HealthCheckHandler{}
}

func (h *HealthCheckHandler) Handle(ctx context.Context, req *HealthCheckRequest) (*HealthCheckResponse, error) {
	return &HealthCheckResponse{Status: "OK"}, nil
}
