package ioc

import "github.com/redis/go-redis/v9"

/**
 * @Description
 * @Date 2024/3/12 19:01
 **/

func InitRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return redisClient
}
