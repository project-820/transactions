package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type PostgresParams struct {
	DSN string
}

func NewPostgresDB(ctx context.Context, params PostgresParams) (*bun.DB, error) {
	cfg, err := pgxpool.ParseConfig(params.DSN)
	if err != nil {
		return nil, fmt.Errorf("pgxpool parse config: %w", err)
	}

	dbPool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("pgxpool new with config: %w", err)
	}

	sqldb := stdlib.OpenDBFromPool(dbPool)

	db := bun.NewDB(sqldb, pgdialect.New())
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
