package utils

import (
	"context"
	"encoding/json"
	"time"

	"sitenavigation/config"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisHelper struct {
	client *redis.Client
}

func NewRedisHelper() *RedisHelper {
	return &RedisHelper{client: config.RedisClient}
}

func (r *RedisHelper) Available() bool {
	return r != nil && r.client != nil
}

func (r *RedisHelper) Set(key string, value interface{}, expiration time.Duration) error {
	if !r.Available() {
		return nil
	}
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisHelper) Get(key string) (string, error) {
	if !r.Available() {
		return "", redis.Nil
	}
	return r.client.Get(ctx, key).Result()
}

func (r *RedisHelper) Del(keys ...string) error {
	if !r.Available() {
		return nil
	}
	return r.client.Del(ctx, keys...).Err()
}

func (r *RedisHelper) SetJSON(key string, value interface{}, expiration time.Duration) error {
	if !r.Available() {
		return nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, expiration).Err()
}

func (r *RedisHelper) GetJSON(key string, dest interface{}) error {
	if !r.Available() {
		return redis.Nil
	}
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}
