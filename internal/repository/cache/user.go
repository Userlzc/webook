package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"project/internal/domain"
	"time"
)

/**
 * @Description
 * @Date 2024/3/8 22:56
 **/

type UserCache interface {
	Get(ctx context.Context, id int64) (domain.User, error)
	Set(ctx context.Context, u domain.User) error
}

type RedisUserCache struct {
	cmd        redis.Cmdable // 操作redis
	expiration time.Duration // 过期时间
}

func NewUserCache(cmd redis.Cmdable) UserCache {
	return &RedisUserCache{
		cmd:        cmd,
		expiration: time.Minute * 15, // 也可以从外边传
	}

}
func (cache *RedisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.key(id)
	data, err := cache.cmd.Get(ctx, key).Result()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	// 反序列化
	err = json.Unmarshal([]byte(data), &u)
	return u, err

}
func (cache *RedisUserCache) Set(ctx context.Context, u domain.User) error {
	key := cache.key(u.Id)
	// value是结构体数据 redis不能存，需序列化
	data, err := json.Marshal(u)
	if err != nil {
		return err
	}
	return cache.cmd.Set(ctx, key, data, cache.expiration).Err()

}

func (cache *RedisUserCache) key(uid int64) string {
	return fmt.Sprintf("user-info-%d", uid)

}
