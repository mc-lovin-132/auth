package repo

import (
	"context"
	"time"

	"github.com/mc-lovin-132/auth/internal/domain"

	"github.com/redis/go-redis/v9"
)

type AccessTokenUsedRepo struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewUsed(rdb *redis.Client, ttl time.Duration) *AccessTokenUsedRepo {
	return &AccessTokenUsedRepo{rdb: rdb, ttl: ttl}
}

func (l *AccessTokenUsedRepo) Add(ctx context.Context, jti, value string) error {
	exists, err := l.Exists(ctx, jti)
	if err != nil {
		return redisErrorMappeer(err)
	}
	if exists {
		return domain.ErrAccessTokenAlreadyUsed
	}
	err = l.rdb.Set(ctx, jti, value, l.ttl).Err()
	if err != nil {
		return redisErrorMappeer(err)
	}
	return nil
}
func (l *AccessTokenUsedRepo) Exists(ctx context.Context, jti string) (bool, error) {
	_, err := l.rdb.Get(ctx, jti).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, redisErrorMappeer(err)
	}
	return true, nil
}
