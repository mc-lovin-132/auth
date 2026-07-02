package repo

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type accessGlobalVersionRepo struct {
	rdb          *redis.Client
	startVersion int
	ttl          time.Duration
}

// подумать как наиболее эффективно удалять устаревшие версии
func NewGlobalVersion(rdb *redis.Client, startVersion int, ttl time.Duration) *accessGlobalVersionRepo {
	return &accessGlobalVersionRepo{
		rdb:          rdb,
		startVersion: startVersion,
		ttl:          ttl,
	}
}

func (r *accessGlobalVersionRepo) GetVersion(ctx context.Context, userID int) (int, error) {
	val, err := r.rdb.Get(ctx, strconv.Itoa(userID)).Result()
	if err != nil {
		if err == redis.Nil {
			return r.startVersion, nil
		}
		return intZeroValue, redisErrorMappeer(err)
	}
	v, err := strconv.Atoi(val)
	if err != nil {
		return intZeroValue, redisErrorMappeer(err)
	}
	return v, nil

}
func (r *accessGlobalVersionRepo) UpdateVersion(ctx context.Context, userID int) error {
	lastVersion, err := r.GetVersion(ctx, userID)
	if err != nil {
		return redisErrorMappeer(err)
	}
	err = r.rdb.Set(ctx, strconv.Itoa(userID), strconv.Itoa(lastVersion+1), r.ttl).Err()
	if err != nil {
		return redisErrorMappeer(err)
	}
	return nil
}
