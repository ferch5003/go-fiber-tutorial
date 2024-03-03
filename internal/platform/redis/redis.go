package redis

import (
	"fmt"
	"github.com/ferch5003/go-fiber-tutorial/config"
	"github.com/redis/go-redis/v9"
)

func NewConnection(config *config.EnvVars) (*redis.Client, error) {
	opt, err := redis.ParseURL(config.RedisConnection)
	if err != nil {
		return &redis.Client{}, err
	}

	client := redis.NewClient(opt)

	fmt.Println("Redis Connected!")

	return client, nil
}
