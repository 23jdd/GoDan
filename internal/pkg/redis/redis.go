package redis

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"godan/internal/config"
)

var RDB *goredis.Client

func Init(cfg *config.RedisConfig) error {
	client := goredis.NewClient(&goredis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to ping redis: %w", err)
	}

	RDB = client
	return nil
}

func Close() {
	if RDB != nil {
		RDB.Close()
	}
}

func Get(ctx context.Context, key string) (string, error) {
	return RDB.Get(ctx, key).Result()
}

func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return RDB.Set(ctx, key, value, expiration).Err()
}

func Del(ctx context.Context, keys ...string) error {
	return RDB.Del(ctx, keys...).Err()
}

func Exists(ctx context.Context, key string) (bool, error) {
	n, err := RDB.Exists(ctx, key).Result()
	return n > 0, err
}
