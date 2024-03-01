package redis

import (
	"context"
	"fmt"
	"github.com/ferch5003/go-fiber-tutorial/config"
	"github.com/ferch5003/go-fiber-tutorial/internal/platform"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

type redisContainer struct {
	ctx       context.Context
	container *redis.RedisContainer
}

func NewRedisContainer(ctx context.Context) platform.Container {
	return &redisContainer{
		ctx: ctx,
	}
}

func (r redisContainer) CreateOrUseContainer(config *config.EnvVars) (err error) {
	container, err := redis.RunContainer(r.ctx,
		testcontainers.WithImage("docker.io/redis:latest"),
		redis.WithSnapshotting(10, 2),
		redis.WithLogLevel(redis.LogLevelVerbose),
	)
	if err != nil {
		return fmt.Errorf("failed to start container: %s", err)
	}

	r.container = container

	connectionString, err := r.container.ConnectionString(r.ctx)
	if err != nil {
		return fmt.Errorf("failed to obtain connection string: %s", err)
	}

	config.RedisConnection = connectionString

	fmt.Println(connectionString)

	return nil
}

func (r redisContainer) CleanContainer() (err error) {
	return r.container.Terminate(r.ctx)
}
