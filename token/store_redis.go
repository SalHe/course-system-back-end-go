package token

import (
	"github.com/se2022-qiaqia/course-system/cache"
	"github.com/se2022-qiaqia/course-system/config"
	"time"
)

type RedisTokenStorage struct{}

func (r *RedisTokenStorage) Get(token string) (string, bool) {
	result, err := cache.RedisClient.Get(cache.Ctx(), key(token)).Result()
	if err != nil {
		return "", false
	}
	return result, true
}

func (r *RedisTokenStorage) Set(token string, jwt string) {
	cache.RedisClient.Set(
		cache.Ctx(),
		key(token),
		jwt,
		config.Config.Token.ExpireDuration(),
	)
}

func key(token string) string {
	return cache.Key(cache.PrefixToken, token)
}

func (r *RedisTokenStorage) SetExpire(token string, expire int64) {
	cache.RedisClient.Expire(
		cache.Ctx(),
		key(token),
		time.Duration(expire)*time.Second,
	)
}

func (r *RedisTokenStorage) Delete(token string) {
	cache.RedisClient.Del(
		cache.Ctx(),
		key(token),
	)
}
