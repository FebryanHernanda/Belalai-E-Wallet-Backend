package utils

import (
	"context"
	"encoding/json"
	"fmt"
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

// get redis data return as slice of model
func RedisGetData[M any](reqCntxt context.Context, rdb redis.Client, rediskey string) (*M, error) {
	// Store unmarshalling result on generic type
	var result M
	// cache-aside pattern
	// cek data redis first
	cmd := rdb.Get(reqCntxt, rediskey)
	if err := cmd.Err(); err != nil {
		if err == redis.Nil {
			log.Printf("Redis key %s not found\n", rediskey)
			return nil, nil // cache miss
		}
		log.Println("Redis Error.\nCause:", err.Error())
		return nil, err
	} else {
		// cache hit
		cmdByte, err := cmd.Bytes()
		if err != nil {
			log.Println("Error reading Redis bytes.\nCause:", err.Error())
			return nil, err
		} else {
			if err := json.Unmarshal(cmdByte, &result); err != nil {
				log.Println("Error unmarshalling Redis data.\nCause:", err.Error())
				return nil, err
			}
		}
	}
	// Return value, and error nil if not error
	return &result, nil
}

// Renew cache redis
func RedisRenewData[m any](reqCntxt context.Context, redc redis.Client, rediskey string, anyModel m, tt time.Duration) error {
	// convert any model into byte
	bt, err := json.Marshal(anyModel)
	if err != nil {
		log.Println("Internal Server Error.\nCause: ", err.Error())
	} else {
		if err := redc.Set(reqCntxt, rediskey, string(bt), tt).Err(); err != nil {
			log.Println("Redis Error.\nCause: ", err.Error())
		}
	}
	// return nil nil, if not error
	return nil
}

// redis invalidation
// delete when update some data
func DeleteAllCache(reqContxt context.Context, rdb redis.Client) error {
	rdbKeys := []string{
		"Belalai-E-wallet:filter-user",
	}
	cmd := rdb.Del(reqContxt, rdbKeys...)
	deletedCount, err := cmd.Result()
	if err != nil {
		log.Println("Redis Error.\nCause:", err.Error())
		return err
	}
	if deletedCount == 0 {
		log.Println("No keys were deleted.")
	} else {
		log.Printf("Successfully deleted %d keys.\n", deletedCount)
	}
	// return error nill if success
	return nil
}

// Cache new user data after registration
func CacheNewUser[U any](ctx context.Context, rdb redis.Client, userID int64, userData U) error {
	userKey := fmt.Sprintf("Belalai-E-wallet:user:%d", userID)
	return RedisRenewData(ctx, rdb, userKey, userData, 1*time.Hour)
}

// Invalidate user list caches after registration (since new user added)
func InvalidateUserListCache(ctx context.Context, rdb redis.Client) error {
	keysToDelete := []string{
		"Belalai-E-wallet:filter-user",
		"Belalai-E-wallet:user-list",
		"Belalai-E-wallet:all-users",
	}

	cmd := rdb.Del(ctx, keysToDelete...)
	deletedCount, err := cmd.Result()
	if err != nil {
		log.Println("Redis Error when invalidating user list cache.\nCause:", err.Error())
		return err
	}

	if deletedCount > 0 {
		log.Printf("Successfully invalidated %d user list cache keys after registration.\n", deletedCount)
	}

	return nil
}

// Complete cache operations after successful registration
func HandleRegistrationCache[U any](ctx context.Context, rdb redis.Client, userID int64, userData U) error {
	// Cache the new user data
	if err := CacheNewUser(ctx, rdb, userID, userData); err != nil {
		log.Println("Warning: Failed to cache new user data:", err)
		// Don't return error, continue with invalidation
	}

	// Invalidate user list caches
	if err := InvalidateUserListCache(ctx, rdb); err != nil {
		log.Println("Warning: Failed to invalidate user list cache:", err)
		return err
	}

	return nil
}

// update profile

func InvalidateUserProfileCache(ctx context.Context, rdb redis.Client, userID int64) error {
	keysToDelete := []string{
		fmt.Sprintf("Belalai-E-wallet:user:%d", userID),
		fmt.Sprintf("Belalai-E-wallet:user-profile:%d", userID),
		fmt.Sprintf("Belalai-E-wallet:profile:%d", userID),
		"Belalai-E-wallet:filter-user", // User list might show profile info
	}

	cmd := rdb.Del(ctx, keysToDelete...)
	deletedCount, err := cmd.Result()
	if err != nil {
		log.Println("Redis Error when invalidating user profile cache.\nCause:", err.Error())
		return err
	}

	if deletedCount > 0 {
		log.Printf("Successfully invalidated %d user profile cache keys after update.\n", deletedCount)
	}

	return nil
}

// Cache updated profile data
func CacheUpdatedProfile[P any](ctx context.Context, rdb redis.Client, userID int64, profileData P) error {
	profileKey := fmt.Sprintf("Belalai-E-wallet:user-profile:%d", userID)
	return RedisRenewData(ctx, rdb, profileKey, profileData, 1*time.Hour)
}

// Complete cache operations after successful profile update
func HandleUpdateProfileCache[P any](ctx context.Context, rdb redis.Client, userID int64, profileData P) error {

	if err := InvalidateUserProfileCache(ctx, rdb, userID); err != nil {
		log.Println("Warning: Failed to invalidate profile cache:", err)

	}

	// Cache the updated profile data
	if err := CacheUpdatedProfile(ctx, rdb, userID, profileData); err != nil {
		log.Println("Warning: Failed to cache updated profile data:", err)
		return err
	}

	return nil
}
