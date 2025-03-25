package internal

import (
	"context"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/tanmaygupta069/redis/config"
)

var once sync.Once

var redisClient *redis.Client

func InitializeRedisClient(redisConfig *config.RedisConfig) {
	once.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", redisConfig.Host, redisConfig.Port),
			Password: redisConfig.Password, // no password set
			DB:       redisConfig.Db,       // use default DB
		})

		_, err := client.Ping(context.Background()).Result()
		if err != nil {
			fmt.Printf("error in pinging redis client : %v", err)
		}
		redisClient = client
	})
}

func GetRedisClient(redisConfig *config.RedisConfig) *redis.Client {
	if redisClient == nil {
		InitializeRedisClient(redisConfig)
	}
	return redisClient
}
