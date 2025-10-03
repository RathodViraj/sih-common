package database

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client
var ctx = context.Background()

func InitRedis() *redis.Client {
	redisURL := os.Getenv("REDIS_ADDR")
	if redisURL == "" {
		log.Fatal("No REDIS_ADDR set, skipping Redis init")
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}

	rdb = redis.NewClient(opt)

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis!")
	return rdb
}
