package mysql

import (
	"context"
	"fmt"
	"github.com/ferch5003/go-fiber-tutorial/config"
	"github.com/ferch5003/go-fiber-tutorial/internal/platform"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

type mysSQLContainer struct {
	ctx       context.Context
	container *mysql.MySQLContainer
}

func NewMySQLContainer(ctx context.Context) platform.Container {
	return &mysSQLContainer{
		ctx: ctx,
	}
}

func (c *mysSQLContainer) CreateOrUseContainer(config *config.EnvVars) error {
	mysqlContainer, err := mysql.RunContainer(c.ctx,
		testcontainers.WithImage("mysql:latest"),
		mysql.WithDatabase(config.MySQLDB),
		mysql.WithUsername(config.MySQLUsername),
		mysql.WithPassword(config.MySQLPassword),
	)
	if err != nil {
		return fmt.Errorf("failed to start container: %s", err)
	}

	c.container = mysqlContainer

	connectionString, err := c.container.ConnectionString(c.ctx, "tls=skip-verify")
	if err != nil {
		return fmt.Errorf("failed to obtain connection string: %s", err)
	}

	config.MySQLDSN = connectionString

	return nil
}

func (c *mysSQLContainer) CleanContainer() error {
	return c.container.Terminate(c.ctx)
}
