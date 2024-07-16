package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	R *redis.Client
}

var ctx = context.Background()

func NewRedisClient() *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	return &RedisClient{
		R: client,
	}
}

func(rdb *RedisClient) BlacklistToken(token string, expirationTime time.Duration) error {
	
	err := rdb.R.Set(ctx, token, "blacklisted", expirationTime).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rdb *RedisClient) IsTokenBlacklisted(token string) (bool, error) {
    val, err := rdb.R.Get(ctx, token).Result()
    if err == redis.Nil {
        return false, nil
    } else if err != nil {
        return false, err
    }
    return val == "blacklisted", nil
}