package inmemory

import (
	"context"
	"time"
)

type Redis interface {
	GetRedis(ctx context.Context, key string) (string, error)
	SetRedis(ctx context.Context, key string, value string, ttl time.Duration) error
}