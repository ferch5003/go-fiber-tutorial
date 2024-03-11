package bootstrap

import (
	"context"
	"github.com/ferch5003/go-fiber-tutorial/cmd/api/router"
	"github.com/ferch5003/go-fiber-tutorial/config"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"sync"
	"testing"
)

type mockUserRouter struct {
	mock.Mock
}

func (m *mockUserRouter) Register() {
	m.Called()
}

func TestStart_Successful(t *testing.T) {
	// Given
	server := &Server{
		ErrChan: make(chan error),
		Wg:      &sync.WaitGroup{},
		Mutex:   &sync.Mutex{},
	}

	defer close(server.ErrChan)

	mur := new(mockUserRouter)
	mur.On("Register")

	app := fx.New(
		fx.Supply(
			fx.Annotate(
				mur,
				fx.As(new(router.Router))),
		),
		fx.Provide(router.NewRouter),
		fx.Provide(zap.NewDevelopment),
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
		Wg:      &sync.WaitGroup{},
		Mutex:   &sync.Mutex{},
	}

	defer close(server.ErrChan)

	mur := new(mockUserRouter)
	mur.On("Register")

	app := fx.New(
		fx.Supply(
			fx.Annotate(
				mur,
				fx.As(new(router.Router))),
		),
		fx.Provide(router.NewRouter),
		fx.Provide(zap.NewDevelopment),
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
