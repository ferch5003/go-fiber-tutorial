package main

import (
	"context"
	"fmt"
	"github.com/ferch5003/go-fiber-tutorial/cmd/api/bootstrap"
	"github.com/ferch5003/go-fiber-tutorial/cmd/api/router"
	"github.com/ferch5003/go-fiber-tutorial/config"
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

	app := fx.New(
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

		// Start web server.
		fx.Invoke(bootstrap.Start),
	)

	defer close(server.ErrChan)

	select {
	case <-server.ErrChan:
	default:
	}

	if err := app.Start(ctx); err != nil {
		panic(err)
	}

	select {
	case <-ctx.Done():
		fmt.Println("Application terminated successfully!")

		if err := mySQLContainer.CleanContainer(); err != nil {
			fmt.Println("Error cleaning MySQL container: ", err)
		}

		return
	case err := <-server.ErrChan:
		fmt.Println(err)

		if err = app.Stop(ctx); err != nil {
			fmt.Println("Error stopping the app..")
		}

		if err := mySQLContainer.CleanContainer(); err != nil {
			fmt.Println("Error cleaning MySQL container: ", err)
		}

		return
	}
}
