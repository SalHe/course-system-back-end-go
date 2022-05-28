package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/se2022-qiaqia/course-system/config"
)

var RedisClient *redis.Client
var ctx = context.Background()

const (
	PrefixToken = "TOKEN"
)

func Init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Addr(),
		Password: config.Config.Redis.Password,
		DB:       config.Config.Redis.Db,
	})
}

func Key(prefix, key string) string {
	return config.Config.Redis.Prefix + "_" + prefix + "_" + key
}

func Ctx() context.Context {
	return ctx
}
