package testhelpers

import (
	"context"
	"fmt"
	"go-ecommerce/internal/adapters/config"
	"go-ecommerce/internal/adapters/storage/database/postgres/models"
	"go-ecommerce/internal/core/ports"
	"log/slog"
	"testing"

	"go-ecommerce/internal/adapters/storage/cache/redis"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// storage test containers
type PostgresContainer struct {
	Container testcontainers.Container
	DB        *gorm.DB
}

type RedisContainer struct {
	Container testcontainers.Container
	Client    ports.CacheRepository
}

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

// create a postgres container based an image
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

	// mapping host and port of container with postgres db
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

	// connect to database
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

// create a postgres redis based an image
func NewRedisContainer(t *testing.T) (*RedisContainer, error) {
	t.Helper()
	ctx := context.Background()

	contConfig := testcontainers.ContainerRequest{
		Image:        "redis:8.2.1-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForLog("Ready to accept connections"),
			wait.ForListeningPort("6379/tcp"),
		),
	}

	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: contConfig,
		Started:          true,
	})
	require.NoError(t, err)

	// mapping host and port of cotainer to redis
	port, err := redisC.MappedPort(ctx, "6379/tcp")
	require.NoError(t, err)

	host, err := redisC.Host(ctx)
	require.NoError(t, err)

	// set redis address
	addr := fmt.Sprintf("%s:%s", host, port.Port())

	// create redis client
	client, err := redis.New(ctx, &config.Redis{
		Addr: addr,
		DB:   0,
	})
	require.NoError(t, err)

	return &RedisContainer{
		Container: redisC,
		Client:    client,
	}, nil
}
