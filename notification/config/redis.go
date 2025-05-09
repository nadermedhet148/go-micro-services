package config

import (
	"os"

	goredislib "github.com/go-redis/redis"
	"github.com/go-redsync/redsync/v4/redis"
	"github.com/go-redsync/redsync/v4/redis/goredis"
)

func ConnectRedis() redis.Pool {
	client := goredislib.NewClient(&goredislib.Options{
		Addr: os.Getenv("REDIS_HOST"),
	})
	return goredis.NewPool(client)
}
