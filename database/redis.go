package database

import (
	"fmt"
	"os"
	"sync"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
)

var (
	redisClient *redis.Client
	once        sync.Once
)

func GetRedisClient() *redis.Client {
	once.Do(func() {
		// Load .env file
		err := godotenv.Load()
		if err != nil {
			fmt.Println("Error loading .env file")
		}

		// Read Redis host and port from environment variables
		redisHost := os.Getenv("REDIS_HOST")
		redisPort := os.Getenv("REDIS_PORT")

		if redisHost == "" {
			redisHost = "localhost"
		}
		if redisPort == "" {
			redisPort = "6379"
		}

		redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

		redisClient = redis.NewClient(&redis.Options{
			Addr: redisAddr,
		})
	})
	return redisClient
}
