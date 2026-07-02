package repo

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// TODO: подумать как наиболее эффективно удалять устаревшие версии
type accessDeviceVersionRepo struct {
	rdb          *redis.Client
	startVersion int
	ttl          time.Duration
}

func NewDeviceVersion(rdb *redis.Client, startVersion int, ttl time.Duration) *accessDeviceVersionRepo {
	return &accessDeviceVersionRepo{
		rdb:          rdb,
		startVersion: startVersion,
		ttl:          ttl,
	}
}

func (r *accessDeviceVersionRepo) GetVersion(ctx context.Context, userID int, deviceID string) (int, error) {
	key := fmt.Sprintf("%d:%s", userID, deviceID)
	val, err := r.rdb.Get(ctx, key).Result()
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
func (r *accessDeviceVersionRepo) UpdateVersion(ctx context.Context, userID int, deviceID string) error {
	lastVersion, err := r.GetVersion(ctx, userID, deviceID)
	if err != nil {
		return redisErrorMappeer(err)
	}
	key := fmt.Sprintf("%d:%s", userID, deviceID)
	err = r.rdb.Set(ctx, key, strconv.Itoa(lastVersion+1), r.ttl).Err()
	if err != nil {
		return redisErrorMappeer(err)
	}
	return nil
}
