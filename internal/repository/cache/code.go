package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

/**
 * @Description
 * @Date 2024/3/11 18:14
 **/
var (
	//go:embed lua/set_code.lua
	luaSetCode string
	//go:embed lua/verify_code.lua
	luaVerifyCode       string
	ErrCodeSendToMany   = errors.New("发送太频繁")
	ErrCodeVerifyToMany = errors.New("验证太频繁")
)

type CodeCache interface {
	Set(ctx context.Context, biz string, phone string, code string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}

type RedisCodeCache struct {
	cmd redis.Cmdable
}

func NewCodeCache(cmd redis.Cmdable) CodeCache {
	return &RedisCodeCache{
		cmd: cmd,
	}

}

func (c *RedisCodeCache) Set(ctx context.Context, biz string, phone string, code string) error {
	// 调用lua脚本
	res, err := c.cmd.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case -2:
		return errors.New("验证码存在 但是没有过期时间")
	case -1:
		return ErrCodeSendToMany
	default:
		return nil // 发送成功
	}
}
func (c *RedisCodeCache) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	key := c.key(biz, phone)
	res, err := c.cmd.Eval(ctx, luaVerifyCode, []string{key}, inputCode).Int()
	if err != nil {
		return false, err
	}
	switch res {
	case -2:
		return false, nil // 这里-2代表的是输入错误可重新输入 只允许3次
	case -1:
		return false, ErrCodeVerifyToMany // -1则代表被攻击了输入的次数太多了
	default:
		return true, nil // 发送成功
	}

}
func (c *RedisCodeCache) key(biz, phone string) string {

	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
