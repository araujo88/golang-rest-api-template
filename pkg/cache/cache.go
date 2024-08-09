package cache

import (
	"os"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":6379", // Redis server address (change to localhost when running local)
		Password: "",                                // Password, leave empty if none
		DB:       0,                                 // Default DB
	})
}
