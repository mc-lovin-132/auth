package delivery

import (
	"context"

	"github.com/mc-lovin-132/auth/pb"
)

type healthService interface {
	Status(ctx context.Context) (string, error)
}

type HealthHandler struct {
	pb.UnimplementedHealthServiceServer
	service healthService
}

func NewHealthHandler(service healthService) *HealthHandler {
	return &HealthHandler{service: service}
}

func (h *HealthHandler) Health(ctx context.Context, in *pb.HealthRequest) (*pb.HealthResponse, error) {
	status, err := h.service.Status(ctx)
	if err != nil {
		return nil, errorMapper(err)
	}
	return &pb.HealthResponse{Status: status}, nil
}
func (h *HealthHandler) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingResponse, error) {
	select {
	case <-ctx.Done():
		return nil, errorMapper(ctx.Err())
	default:
	}
	return &pb.PingResponse{Message: "pong"}, nil
}
