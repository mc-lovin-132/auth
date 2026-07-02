package managers

import (
	"time"

	"github.com/mc-lovin-132/auth/internal/domain"
)

type Claims struct {
	ID            string    `json:"jti"`
	UserID        int       `json:"user_id"`
	DeviceID      string    `json:"device_id"`
	ExpiredAt     time.Time `json:"exp"`
	GlobalVersion int       `json:"global_version"`
	DeviceVersion int       `json:"device_version"`
	SessionID     string    `json:"session_id"`
}

func (c Claims) Valid() error {
	if c.ExpiredAt.Before(time.Now()) {
		return domain.ErrAccessTokenExpired
	}
	return nil
}

func toDomain(value string, c *Claims) *domain.AccessToken {
	return &domain.AccessToken{
		Value:         value,
		ID:            c.ID,
		UserID:        c.UserID,
		DeviceID:      c.DeviceID,
		ExpiredAt:     c.ExpiredAt,
		GlobalVersion: c.GlobalVersion,
		DeviceVersion: c.DeviceVersion,
		SessionID:     c.SessionID,
	}
}

func fromDomain(expiredAt time.Time, tokenID string, d *domain.AccessToken) *Claims {
	return &Claims{
		ID:            tokenID,
		UserID:        d.UserID,
		DeviceID:      d.DeviceID,
		ExpiredAt:     expiredAt,
		GlobalVersion: d.GlobalVersion,
		DeviceVersion: d.DeviceVersion,
		SessionID:     d.SessionID,
	}
}
