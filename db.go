package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectDB() (*sql.DB, error) {
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		return nil, fmt.Errorf("Database URL not found")
	}

	db, errDB := sql.Open("pgx", databaseURL)
	if errDB != nil {
		return nil, fmt.Errorf("database connection failed: %w", errDB)
	}

	ctxBackground := context.Background()
	ctx, cancel := context.WithTimeout(ctxBackground, 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	return db, nil
}
