package postgres

import (
	"1337bo4rd/internal/adapter/config"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/lib/pq"
)

// Singleton Pattern: SQL database connection
var postgresDB *sql.DB

func OpenDB(config *config.DB) (*sql.DB, error) {
	// init varibales
	dbName := config.Name
	host := config.Host
	userName := config.User
	password := config.Password
	port := config.Port

	// postgres://YourUserName:YourPassword@YourHostname:5432/YourDatabaseName
	// postgres://username:password@host/dbname
	var dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", userName, password, host, port, dbName)

	if postgresDB != nil {
		return postgresDB, nil
	}

	slog.Info("Trying to connect to PostgreSQL database...")

	// Retry logic: attempt to connect multiple times
	maxRetries := 3                  // Try 6 times (30 seconds total if we wait 5 seconds between retries)
	retryInterval := 5 * time.Second // Retry every 5 seconds
	var db *sql.DB
	var err error

	// Try connecting up to maxRetries times
	for range maxRetries {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			// Create a context with a timeout for the Ping operation
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Try to ping the database to check if it's available
			err = db.PingContext(ctx)
			if err == nil {
				// If the ping is successful, use the connection
				postgresDB = db
				slog.Info("PostgreSQL database connection established")
				return postgresDB, nil
			}
		}

		// If any error occurs, log it and retry after a delay
		//slog.Errorf("Failed to connect to PostgreSQL, retrying in %v... (attempt %d/%d)", retryInterval, i+1, maxRetries)
		time.Sleep(retryInterval)
	}

	slog.Error("Failed to connect to PostgreSQL after", "attempts", maxRetries, "interval", retryInterval, "error", err)
	return nil, err
}
