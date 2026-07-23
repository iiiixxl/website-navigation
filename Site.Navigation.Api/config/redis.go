package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         GetRedisAddress(),
		Password:     GetRedisPassword(),
		DB:           GetRedisDB(),
		PoolSize:     GetRedisPoolSize(),
		MinIdleConns: GetRedisMinIdleConns(),
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis 连接失败: %w", err)
	}

	log.Println("Redis 连接成功")
	return nil
}

func CloseRedis() {
	if RedisClient == nil {
		return
	}
	if err := RedisClient.Close(); err != nil {
		log.Printf("关闭 Redis 连接失败: %v", err)
		return
	}
	log.Println("Redis 连接已关闭")
}
