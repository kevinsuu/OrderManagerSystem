package repository

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/kevinsuu/OrderManagerSystem/order-service/internal/config"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(cfg config.RedisConfig) *RedisRepository {
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

	return &RedisRepository{
		client: client,
	}
}

func (r *RedisRepository) Close() error {
	return r.client.Close()
}

func (r *RedisRepository) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisRepository) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
