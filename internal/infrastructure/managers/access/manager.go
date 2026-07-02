package managers

import (
	"context"
	"time"

	"github.com/mc-lovin-132/auth/internal/domain"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type AccessManager struct {
	secretKey string
	tokenTTL  time.Duration
}

func New(secretKey string, tokenTTL time.Duration) *AccessManager {
	return &AccessManager{
		secretKey: secretKey,
		tokenTTL:  tokenTTL,
	}
}

func (m *AccessManager) Gen(ctx context.Context, data *domain.AccessToken) (*domain.AccessToken, error) {
	expiredAt := time.Now().Add(m.tokenTTL)
	tokenID := uuid.New().String()

	// payload
	claims := fromDomain(expiredAt, tokenID, data)

	// создаем новый токен и указываем метод хеширования (заголовки добавляются автоматически)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// добавляем токену подпись на основе секретного ключа
	tokenValue, err := token.SignedString([]byte(m.secretKey))
	if err != nil {
		return nil, accessManagerErrorMapper(err)
	}
	return toDomain(tokenValue, claims), nil
}

// type Keyfunc func(*Token) (interface{}, error)
func (m *AccessManager) Parse(ctx context.Context, value string) (*domain.AccessToken, error) {
	// функция которая должна возвращать секретный ключ для сверки подписей
	keyFunc := func(token *jwt.Token) (any, error) { return m.secretKey, nil }
	// парсинг пэйлоада из токена и валидация его
	// (проверка что подписи совпадают, вызов метода Valid())
	token, err := jwt.ParseWithClaims(value, &Claims{}, keyFunc)
	if err != nil {
		return nil, accessManagerErrorMapper(err)
	}
	// приводим пейлоад к нашему кастомному типо и проверяем что токен был отвалидирован корректно
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return toDomain(value, claims), nil
	}

	return nil, domain.ErrInvalidAccessToken

}
