package repository

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/kevinsuu/OrderManagerSystem/product-service/internal/config"
)

// RedisRepository 定義了 Redis 儲存庫的接口
type RedisRepository interface {
	Close() error
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, key string) error
}

// redisRepository 是 RedisRepository 接口的具體實現
type redisRepository struct {
	client *redis.Client
}

func NewRedisRepository(cfg config.RedisConfig) RedisRepository {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 測試連接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}

	return &redisRepository{
		client: client,
	}
}

func (r *redisRepository) Close() error {
	return r.client.Close()
}

func (r *redisRepository) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *redisRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *redisRepository) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
