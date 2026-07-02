package domain

import (
	"errors"
	"time"
)

type AccessToken struct {
	Value         string
	ID            string
	UserID        int
	DeviceID      string
	ExpiredAt     time.Time
	GlobalVersion int
	DeviceVersion int
	SessionID     string
}

type RefreshToken struct {
	Value     string
	UserID    int
	DeviceID  string
	ExpiredAt time.Time
	Revoked   bool
	Used      bool
	SessionID string
}

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

var (
	// validation errors 400
	ErrUserNotFound    = errors.New("user not found") // 404
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrEmptyPassword   = errors.New("password is empty")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidRequest  = errors.New("invalid request") // общая ошибка для обозначения неверных запросов
	// access token errors 400
	ErrInvalidAccessToken        = errors.New("invalid access token") // общая ошибка, нужна для ситуаций по типу когда невозможно отвалидировать подпись и тп
	ErrAccessTokenExpired        = errors.New("access token expired")
	ErrAccessTokenRevoked        = errors.New("access token revoked") // возникает при несовпадении версий токена с актуальнымиы
	ErrAccessTokenAlreadyUsed    = errors.New("access token already used")
	ErrAccessTokenAlreadyRevoked = errors.New("access token already revoked")
	ErrInvalidSignature          = errors.New("err invalid sign")
	// refresh token errors 400
	ErrRefreshTokenExpired     = errors.New("refresh token expired")
	ErrRefreshTokenRevoked     = errors.New("refresh token revoked")
	ErrRefreshTokenNotFound    = errors.New("refresh token not found") // 404
	ErrRefreshTokenAlreadyUsed = errors.New("refresh token alreaddy used")
	ErrRefreshTokenNotUnique   = errors.New("refresh token not unique")
	ErrNotEnoughArgs           = errors.New("not enough arguments")
	ErrInternal                = errors.New("internal error")
)
