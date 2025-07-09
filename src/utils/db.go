package utils

import (
	// "fmt"

	"github.com/redis/go-redis/v9"
)

func GetRedisClient() *redis.Client {
	err := LoadEnv()
	if err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: DB_URI,
		Username: DB_USERNAME,
		Password: DB_PASSWORD,
		DB: 0,
	})

	return redisClient
}
