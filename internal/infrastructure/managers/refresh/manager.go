package service

import (
	"context"
	"time"

	"github.com/mc-lovin-132/auth/internal/domain"

	"github.com/google/uuid"
)

type RefreshManager struct {
	tokenTTL time.Duration
}

func New(tokenTTL time.Duration) *RefreshManager {
	return &RefreshManager{tokenTTL: tokenTTL}
}

func (m *RefreshManager) Gen(ctx context.Context, userID int, deviceID, sessionID string) *domain.RefreshToken {
	expiredAt := time.Now().Add(m.tokenTTL)
	tokenValue := uuid.New().String()

	return &domain.RefreshToken{
		Value:     tokenValue,
		UserID:    userID,
		DeviceID:  deviceID,
		ExpiredAt: expiredAt,
		Revoked:   false,
		Used:      false,
		SessionID: sessionID,
	}
}
