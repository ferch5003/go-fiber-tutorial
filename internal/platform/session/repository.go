package session

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type Repository interface {
	SetSession(ctx context.Context, token string, claims map[string]any) error
	GetSession(ctx context.Context, token string) (map[string]string, error)
}

type repository struct {
	conn *redis.Client
}

func NewRepository(conn *redis.Client) Repository {
	return &repository{conn: conn}
}

func (r repository) SetSession(ctx context.Context, token string, claims map[string]any) error {
	sessionExists := r.conn.Exists(ctx, fmt.Sprintf("user:%s", token)).Val()
	if sessionExists == 0 {
		err := r.conn.HSet(ctx, fmt.Sprintf("user:%s", token), claims).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r repository) GetSession(ctx context.Context, token string) (map[string]string, error) {
	return r.conn.HGetAll(ctx, fmt.Sprintf("user:%s", token)).Result()
}
