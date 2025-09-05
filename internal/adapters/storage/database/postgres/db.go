package postgres

import (
	"context"
	"go-ecommerce/internal/adapters/config"
	"go-ecommerce/internal/adapters/storage/database/postgres/models"
	"log/slog"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gorm_logger "gorm.io/gorm/logger"
)

var DB *gorm.DB

// create a new connections or instance of database
func New(ctx context.Context, config *config.DB) (*gorm.DB, error) {
	// Connection to db using GORM
	conn, err := gorm.Open(postgres.Open(config.DSN), &gorm.Config{
		Logger: gorm_logger.Default.LogMode(gorm_logger.Info),
	})
	if err != nil {
		slog.Error("error connecting database", "error", err)
		return nil, err
	}

	// Low level connection
	sqlDB, err := conn.DB()
	if err != nil {
		slog.Error("error returning sql DB pointer", "error", err)
		return nil, err
	}

	// Pool connections config
	maxOpenConns, err := strconv.Atoi(config.MaxOpenConnections)
	if err != nil {
		slog.Error("error obtaining max open connections database config", "error", err)
		return nil, err
	}

	maxIdleConns, err := strconv.Atoi(config.MaxIdleConnections)
	if err != nil {
		slog.Error("error obtaining max idle connections database config", "error", err)
		return nil, err
	}

	maxLifeTime, err := strconv.Atoi(config.MaxLifeTime)
	if err != nil {
		slog.Error("error obtaining max life time database config", "error", err)
		return nil, err
	}

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(maxLifeTime) * time.Minute)

	DB = conn

	return DB, nil
}

// close database connections
func Close() {
	if DB == nil {
		return
	}

	sqlDB, err := DB.DB()
	if err != nil {
		slog.Error("failed to get sql.DB for closing", "error", err)
		return
	}

	if err := sqlDB.Close(); err != nil {
		slog.Error("failed to close database connection", "error", err)
		return
	}

	slog.Info("database connection closed successfully")
}

// execute migrations
func Migrate(db *gorm.DB) error {
	err := automigrateSchemas(db,
		&models.UserModel{},
	)
	if err != nil {
		return err
	}
	return nil
}

// loop for all migrations and execute
func automigrateSchemas(db *gorm.DB, models ...interface{}) error {
	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			slog.Error("Error executing migrations", "migrations", "database")
			return err
		}
	}
	slog.Info("Migrations successfully executed")
	return nil
}
