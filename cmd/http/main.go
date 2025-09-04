package main

import (
	"context"
	"go-ecommerce/internal/adapters/config"
	"go-ecommerce/internal/adapters/logger"
	"go-ecommerce/internal/adapters/storage/database/postgres/connection"
	"go-ecommerce/internal/adapters/storage/database/postgres/migrations"
	"log/slog"
	"os"
)

func main() {
	config, err := config.New()
	if err != nil {
		slog.Error("error loading environment variables", "error", err)
		os.Exit(1)
	}

	// Set logger
	logger.Set(config.App)
	slog.Info("Starting application", "app", config.App.Name, "env", config.App.Env)

	// Init database
	ctx := context.Background()
	db, err := connection.New(ctx, config.DB)
	if err != nil {
		slog.Error("Error initializing database connection", "error", err)
		os.Exit(1)
	}
	defer connection.Close()

	// execute database migrations
	err = migrations.ExecMigrations(db)
	if err != nil {
		slog.Error("Error executing database migrations", "error", err)
		os.Exit(1)
	}

}
