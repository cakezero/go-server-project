package utils

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/redis/go-redis/v9"
)

func GetRedisClient() *redis.Client {

	redisClient := redis.NewClient(&redis.Options{
		Addr: REDIS_URI,
		Username: REDIS_USERNAME,
		Password: REDIS_PASSWORD,
		DB: 0,
	})

	return redisClient
}

func DB () error {
	err := mgm.SetDefaultConfig(nil, "go_calculator", options.Client().ApplyURI(DB_URI))

	return err
}
