package test_containers

import (
	"context"
	"fmt"
	"go-ecommerce/internal/adapters/config"
	"go-ecommerce/internal/adapters/storage/cache/redis"
	"go-ecommerce/internal/core/ports"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type RedisContainer struct {
	Container testcontainers.Container
	Client    ports.CacheRepository
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
