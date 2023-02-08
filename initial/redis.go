package initializer

import (
	"fmt"

	"github.com/go-redis/redis"
)

func RedisConnect() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	status := redisClient.Ping()
	fmt.Println(status)

	return redisClient
}
