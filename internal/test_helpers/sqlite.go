package testhelpers

import (
	"go-ecommerce/internal/adapters/storage/database/postgres/models"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// to unit tests
// sqlite allows create a database in memory
func NewSQLiteTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := "file::memory:?cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	// foreign keys on
	require.NoError(t, db.Exec("PRAGMA foreign_keys = ON;").Error)

	// automigrate all models
	require.NoError(t, db.AutoMigrate(
		&models.UserModel{},
		&models.ProductModel{},
		&models.CategoryModel{},
		&models.OrderModel{},
		&models.OrderProductModel{},
	))
	return db
}
