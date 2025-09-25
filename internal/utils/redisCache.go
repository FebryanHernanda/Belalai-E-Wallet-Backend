package utils

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// All cache redis key should writen with start like this
// Belalai-E-wallet:<your argument>

func BlackListTokenRedish(reqCntxt context.Context, rdb redis.Client, token string) error {
	err := rdb.Set(reqCntxt, "Belalai-E-wallet:blacklist:"+token, "true", 30*time.Minute).Err()
	if err != nil {
		log.Println("Redis Error when blacklist token:", err)
		return err
	}
	// return error nil, if success
	return nil
}
