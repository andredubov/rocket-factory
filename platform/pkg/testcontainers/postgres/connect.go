package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

func connectToPostgresDB(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, errors.Errorf("failed to parse DSN: %v", err)
	}

	connPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, errors.Errorf("failed to connect to database: %v", err)
	}

	err = connPool.Ping(ctx)
	if err != nil {
		return nil, errors.Errorf("failed to ping postgres: %v", err)
	}

	return connPool, nil
}
