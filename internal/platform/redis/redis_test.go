package redis

import (
	"context"
	"github.com/ferch5003/go-fiber-tutorial/config"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewConnection_Successful(t *testing.T) {
	t.Parallel()

	// Given
	configs := &config.EnvVars{
		RedisConnection: "redis://localhost:6379",
	}
	container := NewRedisContainer(context.Background())

	err := container.CreateOrUseContainer(configs)
	require.NoError(t, err)

	// When
	conn, err := NewConnection(configs)

	// Then
	require.NoError(t, err)
	require.NotEmpty(t, conn)

	// Necessary to clean containers.
	err = container.CleanContainer()
	require.NoError(t, err)
}
