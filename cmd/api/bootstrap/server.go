package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ferch5003/go-fiber-tutorial/cmd/api/router"
	"github.com/ferch5003/go-fiber-tutorial/config"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"sync"
)

const (
	_defaultHost = "localhost"
	_defaultPort = "3000"
)

type Server struct {
	ErrChan chan error
	Wg      *sync.WaitGroup
	Mutex   *sync.Mutex
}

func NewFiberServer() *fiber.App {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	return app
}

func Start(
	lc fx.Lifecycle,
	cfg *config.EnvVars,
	server *Server,
	app *fiber.App,
	router *router.GeneralRouter,
	logger *zap.Logger) {
	host := _defaultHost // Default Host
	if cfg != nil && cfg.Host != "" {
		host = cfg.Host
	}

	port := _defaultPort // Default Port
	if cfg != nil && cfg.Port != "" {
		port = cfg.Port
	}

	// Log all requests.
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger,
	}))

	app.Use(recover.New())

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info(fmt.Sprintf("Starting fiber server on %s:%s", host, port))

			router.Register()

			server.Wg.Add(1)
			go func() {
				logger.Info("Starting...")

				if err := app.Listen(fmt.Sprintf("%s:%s", host, port)); err != nil {
					server.ErrChan <- err
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing server...")

			return app.Shutdown()
		},
	})
}
