package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	Pool *pgxpool.Pool
}

func NewPostgresDB(ctx context.Context, connString string) (*PostgresDB, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database config: %w", err)
	}

	cfg.MaxConns = 25
	cfg.MinConns = 5
	cfg.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	db := &PostgresDB{Pool: pool}
	if err := db.Migrate(ctx); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

func (db *PostgresDB) Migrate(ctx context.Context) error {
	schema := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	CREATE TABLE IF NOT EXISTS topics (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL,
		title VARCHAR(255) NOT NULL,
		category VARCHAR(100) NOT NULL DEFAULT 'General',
		tone VARCHAR(50) NOT NULL DEFAULT 'orange',
		notes TEXT NOT NULL DEFAULT '',
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS srs_items (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		topic_id UUID NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
		user_id UUID NOT NULL,
		repetitions INT NOT NULL DEFAULT 0,
		interval_days INT NOT NULL DEFAULT 1,
		ease_factor NUMERIC(4, 2) NOT NULL DEFAULT 2.50,
		last_reviewed_at TIMESTAMPTZ,
		next_review_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		CONSTRAINT unique_topic_srs UNIQUE (topic_id)
	);

	CREATE INDEX IF NOT EXISTS idx_srs_due_date ON srs_items (user_id, next_review_at);

	CREATE TABLE IF NOT EXISTS review_logs (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		srs_item_id UUID NOT NULL REFERENCES srs_items(id) ON DELETE CASCADE,
		topic_id UUID NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
		user_id UUID NOT NULL,
		rating INT NOT NULL,
		interval_before INT NOT NULL,
		interval_after INT NOT NULL,
		ease_factor_after NUMERIC(4, 2) NOT NULL,
		reviewed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_review_logs_user_date ON review_logs (user_id, reviewed_at);
	`

	_, err := db.Pool.Exec(ctx, schema)
	if err != nil {
		return err
	}

	log.Println("[Postgres] Database schema verified and up-to-date")
	return nil
}

func (db *PostgresDB) Close() {
	db.Pool.Close()
}
