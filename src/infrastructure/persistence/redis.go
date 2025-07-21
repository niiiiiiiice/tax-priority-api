package persistence

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func NewRedisConfig() *RedisConfig {
	return &RedisConfig{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
	}
}

func ConnectRedis(config *RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password:        config.Password,
		DB:              config.DB,
		PoolSize:        10,
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: time.Minute * 30,
		ReadTimeout:     time.Second * 3,
		WriteTimeout:    time.Second * 3,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Printf("Successfully connected to Redis at %s:%d", config.Host, config.Port)
	return client, nil
}

type RedisKeys struct {
	FAQByID       string
	FAQByCategory string
	FAQActive     string
	FAQCount      string
	FAQCategories string
	FAQSearch     string
}

func NewRedisKeys() *RedisKeys {
	return &RedisKeys{
		FAQByID:       "faq:id:%s",
		FAQByCategory: "faq:category:%s",
		FAQActive:     "faq:active",
		FAQCount:      "faq:count",
		FAQCategories: "faq:categories",
		FAQSearch:     "faq:search:%s",
	}
}

func (r *RedisKeys) GetFAQByIDKey(id string) string {
	return fmt.Sprintf(r.FAQByID, id)
}

func (r *RedisKeys) GetFAQByCategoryKey(category string) string {
	return fmt.Sprintf(r.FAQByCategory, category)
}

func (r *RedisKeys) GetFAQSearchKey(query string) string {
	return fmt.Sprintf(r.FAQSearch, query)
}
