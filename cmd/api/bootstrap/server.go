package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ferch5003/go-fiber-tutorial/cmd/api/router"
	"github.com/ferch5003/go-fiber-tutorial/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/fx"
)

const (
	_defaultHost = "localhost"
	_defaultPort = "3000"
)

type Server struct {
	ErrChan chan error
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
	router *router.Router) {
	host := _defaultHost // Default Host
	if cfg != nil && cfg.Host != "" {
		host = cfg.Host
	}

	port := _defaultPort // Default Port
	if cfg != nil && cfg.Port != "" {
		port = cfg.Port
	}

	// Log all requests.
	app.Use(logger.New())

	app.Use(recover.New())

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Printf("Starting fiber server on %s:%s\n", host, port)

			router.Register()

			go func() {
				fmt.Println("Starting...")

				if err := app.Listen(fmt.Sprintf("%s:%s", host, port)); err != nil {
					server.ErrChan <- err
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Closing server...")

			return app.Shutdown()
		},
	})
}
