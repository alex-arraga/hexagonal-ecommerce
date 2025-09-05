package main

import (
	"context"
	"go-ecommerce/internal/adapters/config"
	"go-ecommerce/internal/adapters/logger"
	"go-ecommerce/internal/adapters/storage/cache/redis"
	"go-ecommerce/internal/adapters/storage/database/postgres"
	"log/slog"
	"os"
)

func main() {
	config, err := config.New()
	if err != nil {
		slog.Error("Error loading environment variables", "error", err)
		os.Exit(1)
	}

	// Set logger
	logger.Set(config.App)
	slog.Info("Starting application", "app", config.App.Name, "env", config.App.Env)

	// Init database
	ctx := context.Background()
	db, err := postgres.New(ctx, config.DB)
	if err != nil {
		slog.Error("Error initializing database connection", "error", err)
		os.Exit(1)
	}
	defer postgres.Close()

	// execute database migrations
	err = postgres.Migrate(db)
	if err != nil {
		slog.Error("Error executing database migrations", "error", err)
		os.Exit(1)
	}
	slog.Info("Successfully migrated the database")

	// Init cache service
	cache, err := redis.New(ctx, config.Redis)
	if err != nil {
		slog.Error("Error initializing cache connection", "error", err)
		os.Exit(1)
	}
	defer cache.Close()
	slog.Info("Successfully connected to the cache server")

}
