package redis

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
	once   sync.Once
)

type Config struct {
	Host     string
	Port     string
	Password string
}

func GetConfig() Config {
	return Config{
		Host:     getEnvOrDefault("REDIS_HOST", "localhost"),
		Port:     getEnvOrDefault("REDIS_PORT", "6379"),
		Password: getEnvOrDefault("REDIS_PASSWORD", ""),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetClient() *redis.Client {
	once.Do(func() {
		config := GetConfig()
		client = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
			Password: config.Password,
			DB:       0,
		})
	})
	return client
}

func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}

func Ping(ctx context.Context) error {
	return GetClient().Ping(ctx).Err()
}
