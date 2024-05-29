package db

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/zer0day88/tinder/config"
)

func InitRedis() (*redis.Client, error) {
	var (
		host     = config.Key.Database.Redis.Host
		password = config.Key.Database.Redis.Password
		port     = config.Key.Database.Redis.Port
	)

	var rdb *redis.Client

	redisURL := fmt.Sprintf("redis://%s:%s@%s:%d", "default", password, host, port)

	if config.Key.Database.Redis.UseTLS {
		redisURL = fmt.Sprintf("rediss://%s:%s@%s:%d", "default", password, host, port)

	}

	addr, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	rdb = redis.NewClient(addr)

	if result := rdb.Ping(context.TODO()); result.Err() != nil {
		return nil, result.Err()
	}

	return rdb, nil
}
