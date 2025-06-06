package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Connect() error {
	// Get database configuration from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Validate required environment variables
	if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
		return fmt.Errorf("missing required database environment variables")
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)

	// Try to connect with retries
	var pool *pgxpool.Pool
	var err error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		pool, err = pgxpool.New(ctx, dsn)
		cancel()

		if err == nil {
			break
		}

		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			time.Sleep(time.Second * 2)
		}
	}

	if err != nil {
		return fmt.Errorf("failed to connect to database after %d attempts: %v", maxRetries, err)
	}

	DB = pool

	// Run migrations
	if err := runMigrations(dsn); err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}

	return nil
}

func IsConnected() bool {
	if DB == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := DB.Ping(ctx)
	return err == nil
}

func runMigrations(dsn string) error {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %v", err)
	}

	// Try different possible paths for migrations
	migrationPaths := []string{
		filepath.Join(wd, "migrations"),
		filepath.Join(wd, "db", "migrations"),
		"/app/migrations",
		"/app/db/migrations",
	}

	var migrationPath string
	for _, path := range migrationPaths {
		if _, err := os.Stat(path); err == nil {
			migrationPath = path
			break
		}
	}

	if migrationPath == "" {
		return fmt.Errorf("migrations directory not found in any of the expected locations")
	}

	log.Printf("Using migrations from: %s", migrationPath)

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationPath),
		dsn,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migration: %v", err)
	}

	return nil
}
