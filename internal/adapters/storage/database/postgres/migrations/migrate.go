package migrations

import (
	"go-ecommerce/internal/adapters/storage/database/postgres/connection"
	"log/slog"
)

func migrate(models ...interface{}) {
	for _, model := range models {
		if err := connection.DB.AutoMigrate(model); err != nil {
			slog.Error("Error executing migrations", "migrations", "database")
		}
	}
	slog.Info("Migrations successfully executed")
}

func ExecMigrations() {
	migrate()
}
