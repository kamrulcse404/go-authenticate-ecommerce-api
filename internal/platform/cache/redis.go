package cache

import (
	"context"
	"ecommerce-api/internal/platform/config"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func Connect(cfg config.RedisConfig) (*redis.Client, error) {

	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.Ping(ctx).Err()
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return client, nil
}
