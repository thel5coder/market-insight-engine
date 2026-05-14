package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"market-insight-engine/internal/config"
	"time"
)

// NewClient initializes an optimized Redis client with a managed connection pool.
func NewClient(ctx context.Context, cfg config.RedisConfig) (*redis.Client, error) {
	opts := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0,

		PoolSize:     cfg.MaxActiveConnection,
		MinIdleConns: cfg.MaxIdleConnection,
		PoolTimeout:  cfg.Timeout + (1 * time.Second),

		DialTimeout:  5 * time.Second,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	}

	rdb := redis.NewClient(opts)

	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := rdb.Ping(pingCtx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	slog.InfoContext(ctx, "Redis client initialized",
		slog.String("host", cfg.Host),
		slog.Int("pool_size", cfg.MaxActiveConnection),
	)

	return rdb, nil
}
