package service

import (
	"context"

	"github.com/mc-lovin-132/auth/internal/domain"
)

type service interface {
	Login(ctx context.Context, email, password string) (*domain.AccessToken, *domain.RefreshToken, error)
	Refresh(ctx context.Context, token string) (*domain.AccessToken, *domain.RefreshToken, error)
	Access(ctx context.Context, token string) (*domain.AccessToken, error)
	RevokeByUser(ctx context.Context, userID int) error
	RevokeByDevice(ctx context.Context, userID int, deviceID string) error
	RevokeConcrete(ctx context.Context, token string) error
}

// ПРИМЕЧАНИЕ: маппинг ошибок тут нигде не нужен так как высокоуровневый Service
// работает с более низкоуровневыми сервисами, которые возвращают доменные ошибки
type Service struct {
	userClient            userClient
	accessTokenGlobalRepo accessTokenGlobalRepo
	accessTokenDeviceRepo accessTokenDeviceRepo
	accessTokenBlackList  accessTokenBlackList
	accessTokenUsedRepo   accessTokenUsedRepo
	refreshTokenRepo      refreshTokenRepo
	deviceFingerprinter   deviceFingerprinter
	sessioner             sessioner
	refreshManager        refreshManager
	accessManager         accessManager
}

func New(
	userClient userClient,
	accessTokenGlobalRepo accessTokenGlobalRepo,
	accessTokenDeviceRepo accessTokenDeviceRepo,
	accessTokenBlackList accessTokenBlackList,
	accessTokenUsedRepo accessTokenUsedRepo,
	refreshTokenRepo refreshTokenRepo,
	deviceFingerprinter deviceFingerprinter,
	sessioner sessioner,
	refreshManager refreshManager,
	accessManager accessManager,
) *Service {
	return &Service{
		userClient:            userClient,
		accessTokenGlobalRepo: accessTokenGlobalRepo,
		accessTokenDeviceRepo: accessTokenDeviceRepo,
		accessTokenBlackList:  accessTokenBlackList,
		accessTokenUsedRepo:   accessTokenUsedRepo,
		refreshTokenRepo:      refreshTokenRepo,
		deviceFingerprinter:   deviceFingerprinter,
		sessioner:             sessioner,
		refreshManager:        refreshManager,
		accessManager:         accessManager,
	}
}

// client
type userClient interface {
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

// redis
type accessTokenGlobalRepo interface {
	GetVersion(ctx context.Context, userID int) (int, error)
	UpdateVersion(ctx context.Context, userID int) error
}

type accessTokenDeviceRepo interface {
	GetVersion(ctx context.Context, userID int, deviceID string) (int, error)
	UpdateVersion(ctx context.Context, userID int, deviceID string) error
}

type accessTokenBlackList interface {
	Add(ctx context.Context, jti, value string) error
	Exists(ctx context.Context, jti string) (bool, error)
}

type accessTokenUsedRepo interface {
	Add(ctx context.Context, jti, value string) error
	Exists(ctx context.Context, jti string) (bool, error)
}

// psql
type refreshTokenRepo interface {
	Get(ctx context.Context, token string) (*domain.RefreshToken, error)
	GetBySessionID(ctx context.Context, sessionID string) (*domain.RefreshToken, error)
	Create(ctx context.Context, token *domain.RefreshToken) error
	MarkAsUsed(ctx context.Context, value string) error // value токена по совместительству является его id
	MarkAsRevokedByDevice(ctx context.Context, userID int, deviceID string) error
	MarkAsRevokedByUser(ctx context.Context, userID int) error
	MarkAsRevokedByConcrete(ctx context.Context, value string) error
}

// generators
type deviceFingerprinter interface {
	// на основе фингерпринта вычисляем deviceID
	// данные необходимые для определения устройства определим позже
	Detect(ctx context.Context) (string, error)
}
type sessioner interface {
	Gen(ctx context.Context) string
}

// managers
// общие данные по типу deviceID, и внешние даные типа UserID получают извне
// остальное генерируют сами (без запросов к сторонним сервисам)
// знают секретный ключ, время жизни токенов
type refreshManager interface {
	Gen(ctx context.Context, userID int, deviceID, sessionID string) *domain.RefreshToken
}
type accessManager interface {
	Gen(ctx context.Context, data *domain.AccessToken) (*domain.AccessToken, error)
	// если подпись невалидна ошибка
	// парсе проверяет не истек ли и валидна ли подпись
	Parse(ctx context.Context, value string) (*domain.AccessToken, error)
}
