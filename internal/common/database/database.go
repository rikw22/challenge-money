package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewConnection(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	// Parse the connection string
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Set common pool configurations
	// TODO: Get this from config file or environment variable
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnIdleTime = time.Minute * 30
	config.MaxConnLifetime = time.Hour * 2
	config.HealthCheckPeriod = time.Minute

	// Create the connection pool
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Ping the database to verify the connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close() // Close the pool if ping fails
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("Successfully connected to the database!")
	return pool, nil
}
