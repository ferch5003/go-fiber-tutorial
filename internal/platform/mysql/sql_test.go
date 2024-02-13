package mysql

import (
	"context"
	"github.com/ferch5003/go-fiber-tutorial/config"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewMySQLConnection_Successful(t *testing.T) {
	// Given
	configs := &config.EnvVars{
		MySQLUsername: "root",
		MySQLPassword: "root",
		MySQLDB:       "fiber_example",
	}
	container := NewMySQLContainer(context.Background())

	err := container.CreateOrUseContainer(configs)
	require.NoError(t, err)

	// When
	conn, err := NewMySQLConnection(configs)

	// Then
	require.NoError(t, err)
	require.NotEmpty(t, conn)

	// Necessary to clean containers.
	err = container.CleanContainer()
	require.NoError(t, err)
}
