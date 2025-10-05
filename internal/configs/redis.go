package configs

import (
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

func InitRedis() *redis.Client {
	// data redish on .env
	RdbUser := os.Getenv("RDB_USER")
	RdbPass := os.Getenv("RDB_PWD")
	RdbHost := os.Getenv("RDB_HOST")
	RdbPort := os.Getenv("RDB_PORT")

	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", RdbHost, RdbPort),
		Username: RdbUser,
		Password: RdbPass,
	})
}
