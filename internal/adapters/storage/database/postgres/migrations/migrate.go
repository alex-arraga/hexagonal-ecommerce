package migrations

import (
	"log/slog"

	"gorm.io/gorm"
)

func migrate(db *gorm.DB, models ...interface{}) error {
	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			slog.Error("Error executing migrations", "migrations", "database")
			return err
		}
	}
	slog.Info("Migrations successfully executed")
	return nil
}

func ExecMigrations(db *gorm.DB) error {
	err := migrate(db)
	if err != nil {
		return err
	}
	return nil
}
