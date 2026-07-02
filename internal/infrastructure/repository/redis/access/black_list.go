package repo

import (
	"context"
	"time"

	"github.com/mc-lovin-132/auth/internal/domain"

	"github.com/redis/go-redis/v9"
)

type AccessTokenBlackList struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewBlackList(rdb *redis.Client, ttl time.Duration) *AccessTokenBlackList {
	return &AccessTokenBlackList{rdb: rdb, ttl: ttl}
}

func (l *AccessTokenBlackList) Add(ctx context.Context, jti, value string) error {
	exists, err := l.Exists(ctx, jti)
	if err != nil {
		return redisErrorMappeer(err)
	}
	if exists {
		return domain.ErrAccessTokenAlreadyRevoked
	}
	err = l.rdb.Set(ctx, jti, value, l.ttl).Err()
	if err != nil {
		return redisErrorMappeer(err)
	}
	return nil
}
func (l *AccessTokenBlackList) Exists(ctx context.Context, jti string) (bool, error) {
	_, err := l.rdb.Get(ctx, jti).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, redisErrorMappeer(err)
	}
	return true, nil
}
