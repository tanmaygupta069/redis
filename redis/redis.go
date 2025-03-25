package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tanmaygupta069/redis/config"
	"github.com/tanmaygupta069/redis/internal"
)

var once sync.Once

type RedisInterface interface {
	Get(key string) (string, error)
	Set(key string, value string, exp int) error
	Delete(key string) (int64, error)
	Exists(key string) (int64, error)
	Decrement(key string, field string) (int64, error)
	SetNx(key string, value string) (bool, error)
}

type RedisImplementation struct {
	redisClient *redis.Client
}

func NewRedisClient(redisConfig *config.RedisConfig) RedisInterface {
	return &RedisImplementation{
		redisClient: internal.GetRedisClient(redisConfig),
	}
}

func (r *RedisImplementation) Get(key string) (string, error) {
	val := r.redisClient.Get(context.Background(), key).Val()
	if val == "" || val == redis.Nil.Error() {
		return "", fmt.Errorf("key not found in cache")
	}
	return val, nil
}

func (r *RedisImplementation) Set(key string, value string, exp int) error {
	err := r.redisClient.Set(context.Background(), key, value, time.Duration(exp)*time.Minute).Err()
	if err != nil {
		return err
	}
	return r.redisClient.Expire(context.Background(), key, time.Duration(exp)*time.Minute).Err()
}

func (r *RedisImplementation) Delete(key string) (int64, error) {
	return r.redisClient.Del(context.Background(), key).Result()
}

func (r *RedisImplementation) Exists(key string) (int64, error) {
	return r.redisClient.Exists(context.Background(), key).Result()
}

func (r *RedisImplementation) Decrement(key string, field string) (int64, error) {
	return r.redisClient.HIncrBy(context.Background(), key, field, -1).Result()
}

func (r *RedisImplementation) SetNx(key string, value string) (bool, error) {
	return r.redisClient.SetNX(context.Background(), key, value, 0).Result()
}
