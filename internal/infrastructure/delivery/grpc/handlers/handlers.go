package delivery

import (
	"context"
	"time"

	"github.com/mc-lovin-132/auth/internal/domain"
	"github.com/mc-lovin-132/auth/pb"
)

type service interface {
	Login(ctx context.Context, email, password string) (*domain.AccessToken, *domain.RefreshToken, error)
	Refresh(ctx context.Context, token string) (*domain.AccessToken, *domain.RefreshToken, error)
	Access(ctx context.Context, token string) (*domain.AccessToken, error)
	RevokeByUser(ctx context.Context, userID int) error
	RevokeByDevice(ctx context.Context, userID int, deviceID string) error
	RevokeConcrete(ctx context.Context, token string) error
}

type Handler struct {
	pb.UnimplementedAuthServiceServer
	service service
}

func New(service service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	err := validateLoginData(in)
	if err != nil {
		return nil, errorMapper(err)
	}
	accessTokenObject, refreshTokenObject, err := h.service.Login(ctx, in.Email, in.Password)
	if err != nil {
		return nil, errorMapper(err)
	}
	return &pb.LoginResponse{
		AccessToken:  accessTokenObject.Value,
		RefreshToken: refreshTokenObject.Value,
	}, nil
}
func (h *Handler) Refresh(ctx context.Context, in *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	accessTokenObject, refreshTokenObject, err := h.service.Refresh(ctx, in.RefreshToken)
	if err != nil {
		return nil, errorMapper(err)
	}
	return &pb.RefreshResponse{
		AccessToken:  accessTokenObject.Value,
		RefreshToken: refreshTokenObject.Value,
	}, nil
}
func (h *Handler) Access(ctx context.Context, in *pb.AccessRequest) (*pb.AccessResponse, error) {
	accessTokenObject, err := h.service.Access(ctx, in.Token)
	if err != nil {
		return nil, errorMapper(err)
	}
	return &pb.AccessResponse{
		UserId:    int64(accessTokenObject.UserID),
		ExpiredAt: accessTokenObject.ExpiredAt.Format(time.RFC3339),
	}, nil
}
func (h *Handler) Revoke(ctx context.Context, in *pb.RevokeRequest) (*pb.RevokeResponse, error) {
	if in.GetAll() != nil {
		data := in.GetAll()
		err := h.service.RevokeByUser(ctx, int(data.UserId))
		if err != nil {
			return nil, errorMapper(err)
		}
	} else if in.GetByDevice() != nil {
		data := in.GetByDevice()
		err := h.service.RevokeByDevice(ctx, int(data.UserId), data.DeviceId)
		if err != nil {
			return nil, errorMapper(err)
		}
	} else if in.GetConcrete() != nil {
		data := in.GetConcrete()
		err := h.service.RevokeConcrete(ctx, data.Token)
		if err != nil {
			return nil, errorMapper(err)
		}
	} else {
		return nil, errorMapper(domain.ErrNotEnoughArgs)
	}
	return &pb.RevokeResponse{}, nil
}
