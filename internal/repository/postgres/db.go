package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"market-insight-engine/internal/config"
	"net/url"
	"time"
)

// NewDB creates a highly optimized connection pool using pgxpool.
func NewDB(ctx context.Context, cfg config.PostgresConfig) (*pgxpool.Pool, error) {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, cfg.Password),
		Host:   fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Path:   cfg.Name,
	}

	q := u.Query()
	q.Set("sslmode", cfg.SSLMode)
	if cfg.ConnectionTimeout > 0 {
		q.Set("connect_timeout", fmt.Sprintf("%d", cfg.ConnectionTimeout))
	}
	u.RawQuery = q.Encode()

	poolConfig, err := pgxpool.ParseConfig(u.String())
	if err != nil {
		return nil, fmt.Errorf("failed to parse constructed postgres DSN: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.PoolSize)
	poolConfig.MinConns = int32(cfg.MinIdleConn)
	poolConfig.MaxConnLifetime = cfg.MaxLifetime
	poolConfig.MaxConnIdleTime = cfg.IdleTimeout
	
	poolConfig.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("postgres ping failed after pool creation: %w", err)
	}

	slog.InfoContext(ctx, "Postgres engine initialized",
		slog.String("db_name", cfg.Name),
		slog.Int("pool_limit", cfg.PoolSize),
	)

	return pool, nil
}
