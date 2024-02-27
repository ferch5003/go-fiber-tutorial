package main

import (
	"context"
	"fmt"
	"github.com/ferch5003/go-fiber-tutorial/cmd/api/bootstrap"
	"github.com/ferch5003/go-fiber-tutorial/cmd/api/router"
	"github.com/ferch5003/go-fiber-tutorial/config"
	"github.com/ferch5003/go-fiber-tutorial/db/seeds"
	"github.com/ferch5003/go-fiber-tutorial/internal/platform/console"
	"github.com/ferch5003/go-fiber-tutorial/internal/platform/mysql"
	"go.uber.org/fx"

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

	app := fx.New(
		// Clear terminal/console
		fx.Invoke(cmd.Clear),

		// creates: config.EnvVars
		fx.Provide(config.NewConfigurations),
		// creates: *bootstrap.Server
		fx.Supply(server),
		// creates: *fiber.Router
		fx.Provide(router.NewRouter),
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

		// Provide seeders
		fx.Provide(seeds.NewSeed),

		fx.Invoke(seeds.Execute),

		// Start web server.
		fx.Invoke(bootstrap.Start),
	)

	defer func() {
		select {
		case _, ok := <-server.ErrChan:
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
	case <-ctx.Done():
		if err := app.Stop(ctx); err != nil {
			fmt.Println("Error stopping the app...", err)
		}

		fmt.Println("Application terminated successfully!")

		if err := mySQLContainer.CleanContainer(); err != nil {
			fmt.Println("Error cleaning MySQL container: ", err)
		}

		return
	case err := <-server.ErrChan:
		fmt.Println(err)

		if err = app.Stop(ctx); err != nil {
			fmt.Println("Error stopping the app...", err)
		}

		if err := mySQLContainer.CleanContainer(); err != nil {
			fmt.Println("Error cleaning MySQL container: ", err)
		}

		return
	}
}
