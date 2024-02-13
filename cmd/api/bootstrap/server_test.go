package bootstrap

import (
	"context"
	"github.com/ferch5003/go-fiber-tutorial/cmd/api/router"
	"github.com/ferch5003/go-fiber-tutorial/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"testing"
)

func TestStart_Successful(t *testing.T) {
	// Given
	server := &Server{
		ErrChan: make(chan error),
	}

	app := fx.New(
		fx.Provide(router.NewRouter),
		fx.Provide(config.NewConfigurations),
		fx.Supply(server),
		fx.Provide(NewFiberServer),

		fx.Invoke(Start),
	)

	ctx := context.Background()

	// When
	err := app.Start(ctx)
	require.NoError(t, err)

	// Then
	err = app.Stop(ctx)
	require.NoError(t, err)
}

func TestStart_FailsDueToInvalidConfiguration(t *testing.T) {
	// Given
	server := &Server{
		ErrChan: make(chan error),
	}

	app := fx.New(
		fx.Provide(router.NewRouter),
		fx.Supply(&config.EnvVars{Host: "bad_host", Port: "bad_port"}),
		fx.Supply(server),
		fx.Provide(NewFiberServer),

		fx.Invoke(Start),
	)
	ctx := context.Background()

	// When
	err := app.Start(ctx)
	require.NoError(t, err)

	// Then
	// Application waits if an error occurs.
	err = <-server.ErrChan
	require.ErrorContains(t, err, "failed to listen")
	require.ErrorContains(t, err, "unknown port")
}
