package redis

import (
	"context"
	"github.com/GoncharovFyodor/hezzltest/internal/config"
	"github.com/redis/go-redis/v9"
	"log"
)

func New(ctx context.Context, cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host,
		Password: cfg.Redis.Password,
		DB:       0,
	})

	err := client.Ping(ctx).Err()

	if err != nil {
		log.Fatal("Не удалось подключиться к redis: ", err)
	}

	return client
}
