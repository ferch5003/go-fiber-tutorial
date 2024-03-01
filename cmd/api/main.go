package main

import (
	"context"
	"github.com/ferch5003/go-fiber-tutorial/cmd/api/bootstrap"
	"github.com/ferch5003/go-fiber-tutorial/cmd/api/router"
	"github.com/ferch5003/go-fiber-tutorial/config"
	"github.com/ferch5003/go-fiber-tutorial/db/seeds"
	"github.com/ferch5003/go-fiber-tutorial/internal/platform/console"
	"github.com/ferch5003/go-fiber-tutorial/internal/platform/mysql"
	"go.uber.org/fx"
	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	server := &bootstrap.Server{
		ErrChan: make(chan error),
	}

	ctx := context.Background()
	mysqlCtx := context.Background()
	mySQLContainer := mysql.NewMySQLContainer(mysqlCtx)
	cmd := console.NewConsole()
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	app := fx.New(
		// Clear terminal/console
		fx.Invoke(cmd.Clear),

		// creates: config.EnvVars
		fx.Provide(config.NewConfigurations),
		// creates: *bootstrap.Server
		fx.Supply(server),
		// creates: *zap.Logger
		fx.Supply(logger),
		// creates: *fiber.Router
		fx.Provide(
			fx.Annotate(
				router.NewRouter,
				fx.ParamTags( // Equivalent to *fiber.App, config.Envars, []Router `group:"routers"` in constructor
					``,
					``,
					`group:"routers"`),
			),
		),
		// creates: *fiber.App
		fx.Provide(bootstrap.NewFiberServer),
		// creates: context.Context
		fx.Supply(ctx),

		// Create MysSQL Container
		fx.Invoke(mySQLContainer.CreateOrUseContainer),

		// creates: *sqlx.DB
		fx.Provide(mysql.NewMySQLConnection),

		// Provide modules
		router.NewUserModule,
		router.NewTodoModule,

		// Provide seeders
		fx.Provide(seeds.NewSeed),

		fx.Invoke(seeds.Execute),

		// Start web server.
		fx.Invoke(bootstrap.Start),
	)

	defer func() {
		select {
		case _, ok := <-(server.ErrChan):
			if ok {
				close(server.ErrChan)
			}
		default:
		}
	}()

	if err := app.Start(ctx); err != nil {
		panic(err)
	}

	select {
	case <-app.Done():
		if err := app.Stop(ctx); err != nil {
			logger.DPanic("Error stopping the app...", zap.Error(err))
		}

		logger.Info("Application terminated successfully!")

		if err := mySQLContainer.CleanContainer(); err != nil {
			logger.DPanic("Error cleaning MySQL container: ", zap.Error(err))
		}
	case err := <-server.ErrChan:
		logger.Info("", zap.Error(err))

		if err = app.Stop(ctx); err != nil {
			logger.DPanic("Error stopping the app...", zap.Error(err))
		}

		if err := mySQLContainer.CleanContainer(); err != nil {
			logger.DPanic("Error cleaning MySQL container: ", zap.Error(err))
		}
	}
}
