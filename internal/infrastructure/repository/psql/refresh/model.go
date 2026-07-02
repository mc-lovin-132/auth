package repo

import (
	"time"

	"github.com/mc-lovin-132/auth/internal/domain"
)

type RefreshTokenModel struct {
	Value     string    `db:"value"`
	UserID    int       `db:"user_id"`
	DeviceID  string    `db:"device_id"`
	ExpiredAt time.Time `db:"expired_at"`
	Revoked   bool      `db:"revoked"`
	Used      bool      `db:"used"`
	SessionID string    `db:"session_id"`
}

func fromDomain(token *domain.RefreshToken) *RefreshTokenModel {
	return &RefreshTokenModel{
		Value:     token.Value,
		UserID:    token.UserID,
		DeviceID:  token.DeviceID,
		ExpiredAt: token.ExpiredAt,
		Revoked:   token.Revoked,
		Used:      token.Used,
		SessionID: token.SessionID,
	}
}

func toDomain(token *RefreshTokenModel) *domain.RefreshToken {
	return &domain.RefreshToken{
		Value:     token.Value,
		UserID:    token.UserID,
		DeviceID:  token.DeviceID,
		ExpiredAt: token.ExpiredAt,
		Revoked:   token.Revoked,
		Used:      token.Used,
		SessionID: token.SessionID,
	}
}
