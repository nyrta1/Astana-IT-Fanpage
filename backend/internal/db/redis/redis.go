package redis

import (
	"aitu-funpage/backend/internal/config"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var ctx = context.Background()

func SaveVerificationCodeToRedis(ctx context.Context, rdbConfig config.RedisConfig, username, verificationCode string) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     rdbConfig.Addr,
		Password: rdbConfig.Password,
		DB:       rdbConfig.DB,
	})
	err := rdb.Set(ctx, fmt.Sprintf("verification_code:%s", username), verificationCode, time.Minute*10).Err()
	if err != nil {
		return err
	}
	return nil
}

func CheckVerificationCode(ctx context.Context, rdbConfig config.RedisConfig, username, verificationCode string) (bool, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     rdbConfig.Addr,
		Password: rdbConfig.Password,
		DB:       rdbConfig.DB,
	})

	val, err := rdb.Get(ctx, fmt.Sprintf("verification_code:%s", username)).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}

	if val == verificationCode {
		_ = rdb.Del(ctx, fmt.Sprintf("verification_code:%s", username))
		return true, nil
	}

	return false, nil
}
