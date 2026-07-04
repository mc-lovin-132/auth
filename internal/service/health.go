package service

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

const (
	await = "await"
	alive = "alive"
	dead  = "dead"
)

type healthService struct {
	db  *sqlx.DB
	rdb *redis.Client
}

func NewHealthService(db *sqlx.DB, rdb *redis.Client) *healthService {
	return &healthService{db: db, rdb: rdb}
}

func (h *healthService) Status(ctx context.Context) (string, error) {
	if h.db == nil {
		return await, nil
	}
	if h.rdb == nil {
		return await, nil
	}
	err := h.db.PingContext(ctx)
	if err != nil {
		return dead, err
	}
	_, err = h.rdb.Ping(ctx).Result()
	if err != nil {
		return dead, err
	}
	return alive, nil
}
