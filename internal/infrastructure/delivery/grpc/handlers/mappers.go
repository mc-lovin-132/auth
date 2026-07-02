package delivery

import (
	"errors"
	"net/mail"

	"github.com/mc-lovin-132/auth/internal/domain"
	"github.com/mc-lovin-132/auth/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func errorMapper(err error) error {
	if errors.Is(err, domain.ErrInvalidEmail) ||
		errors.Is(err, domain.ErrEmptyPassword) ||
		errors.Is(err, domain.ErrInvalidPassword) ||
		errors.Is(err, domain.ErrInvalidRequest) ||
		errors.Is(err, domain.ErrInvalidAccessToken) ||
		errors.Is(err, domain.ErrAccessTokenExpired) ||
		errors.Is(err, domain.ErrAccessTokenRevoked) ||
		errors.Is(err, domain.ErrAccessTokenAlreadyRevoked) ||
		errors.Is(err, domain.ErrAccessTokenAlreadyUsed) ||
		errors.Is(err, domain.ErrInvalidSignature) ||
		errors.Is(err, domain.ErrRefreshTokenExpired) ||
		errors.Is(err, domain.ErrRefreshTokenRevoked) ||
		errors.Is(err, domain.ErrRefreshTokenAlreadyUsed) ||
		errors.Is(err, domain.ErrNotEnoughArgs) {
		return status.Error(codes.InvalidArgument, err.Error())
	} else if errors.Is(err, domain.ErrUserNotFound) || errors.Is(err, domain.ErrRefreshTokenNotFound) {
		return status.Error(codes.NotFound, err.Error())
	} else if errors.Is(err, domain.ErrRefreshTokenNotUnique) {
		return status.Error(codes.AlreadyExists, err.Error())
	} else if errors.Is(err, domain.ErrInternal) {
		return status.Error(codes.Internal, err.Error())
	} else {
		return status.Error(codes.Internal, err.Error())
	}
}

func isEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func validateLoginData(in *pb.LoginRequest) error {
	if !isEmailValid(in.Email) {
		return domain.ErrInvalidEmail
	}
	if in.Password == "" {
		return domain.ErrEmptyPassword
	}
	return nil
}
