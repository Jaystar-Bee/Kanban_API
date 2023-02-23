package initializer

import (
	"fmt"
	"os"

	"github.com/go-redis/redis"
)

func RedisConnect() *redis.Client {

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URI"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	status := redisClient.Ping()
	fmt.Println(status)

	return redisClient
}
