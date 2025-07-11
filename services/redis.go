package services

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisService struct {
	Client *redis.Client
}

func NewRedisService() (*RedisService, error) {
	// Redis connect
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0, // Use default DB
	})
	// Ping the Redis server to check if the connection is successful
	err := client.Ping(context.Background()).Err()
	if err != nil {
		fmt.Println("Failed to connect to Redis:", err)
		return nil, err
	}

	return &RedisService{Client: client}, nil
}

func (r *RedisService) CheckKeyRedis(key string) (bool, error) {
	val, err := r.Client.Get(context.Background(), key).Result()
	if err != nil {
		return false, err
	}
	if val == "" {
		return false, nil
	}
	return true, nil
}

func (r *RedisService) GetKeyRedis(key string) (string, error) {
	val, err := r.Client.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *RedisService) AddKeyRedis(key, value string, interval time.Duration) error {
	err := r.Client.Set(context.Background(), key, value, interval).Err()
	if err != nil {
		return err
	}
	return nil
}
