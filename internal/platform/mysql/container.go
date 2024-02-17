package mysql

import (
	"context"
	"fmt"
	"github.com/ferch5003/go-fiber-tutorial/config"
	"github.com/ferch5003/go-fiber-tutorial/internal/platform"
	"github.com/ferch5003/go-fiber-tutorial/internal/platform/files"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"io/fs"
	"path/filepath"
	"regexp"
)

// _sqlFilesRegex only detects *.sql files on given filepath.
const _sqlFilesRegex = `[\w.-]+\.sql$`

type mysSQLContainer struct {
	ctx       context.Context
	container *mysql.MySQLContainer
}

func NewMySQLContainer(ctx context.Context) platform.Container {
	return &mysSQLContainer{
		ctx: ctx,
	}
}

func getSQLFiles(dir string) ([]string, error) {
	// slice with only *.sql files.
	sqlFiles := make([]string, 0)

	// This regex only accepts *.sql files.
	re := regexp.MustCompile(_sqlFilesRegex)

	walk := func(path string, info fs.FileInfo, err error) error {
		if !re.MatchString(path) {
			return nil
		}

		if !info.IsDir() {
			sqlFiles = append(sqlFiles, path)
		}

		return nil
	}

	err := filepath.Walk(dir, walk)
	if err != nil {
		return []string{}, err
	}

	return sqlFiles, nil
}

func (c *mysSQLContainer) CreateOrUseContainer(config *config.EnvVars) error {
	migrationsDir, err := files.GetDir("migrations")
	if err != nil {
		return err
	}

	migrationFiles, err := getSQLFiles(migrationsDir)
	if err != nil {
		return err
	}

	mysqlContainer, err := mysql.RunContainer(c.ctx,
		testcontainers.WithImage("mysql:latest"),
		mysql.WithDatabase(config.MySQLDB),
		mysql.WithUsername(config.MySQLUsername),
		mysql.WithPassword(config.MySQLPassword),
		mysql.WithScripts(migrationFiles...),
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
