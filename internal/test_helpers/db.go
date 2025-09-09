package testhelpers

import (
	"context"
	"fmt"
	"go-ecommerce/internal/adapters/storage/database/postgres/models"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
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
	require.NoError(t, db.AutoMigrate(&models.UserModel{}))
	return db
}

type PostgresContainer struct {
	Container testcontainers.Container
	DB        *gorm.DB
}

// to integration test
// create a postgres container based image
func NewPostgresContainerDB(t *testing.T) (*PostgresContainer, error) {
	t.Helper()
	ctx := context.Background()

	// container config
	dbUser := "test"
	dbPassword := "secret"
	dbName := "dbtest"

	contConfig := testcontainers.ContainerRequest(testcontainers.ContainerRequest{
		Image:        "postgres:17-alpine",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForListeningPort("5432/tcp"),
		),
		Env: map[string]string{
			"POSTGRES_USER":     dbUser,
			"POSTGRES_PASSWORD": dbPassword,
			"POSTGRES_DB":       dbName,
		},
	})

	// create container
	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: contConfig,
		Started:          true,
	})
	if err != nil {
		slog.Error("error creating postgres container", "error", err)
		return nil, err
	}

	port, err := postgresC.MappedPort(ctx, "5432/tcp")
	if err != nil {
		slog.Error("error obtaining port", "error", err)
		return nil, err
	}

	host, err := postgresC.Host(ctx)
	if err != nil {
		slog.Error("error obtaining host", "error", err)
		return nil, err
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port.Port(), dbUser, dbPassword, dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("error connecting to database", "error", err)
		return nil, err
	}

	// exec migrations
	require.NoError(t, db.AutoMigrate(&models.UserModel{}))

	return &PostgresContainer{
		Container: postgresC,
		DB:        db,
	}, nil
}
