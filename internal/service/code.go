package service

import (
	"context"
	"fmt"
	"math/rand"
	"project/internal/repository"
	"project/internal/service/sms"
)

/**
 * @Description
 * @Date 2024/3/11 16:38
 **/

const tplId = ""

var ErrCodeSendTooMany = repository.ErrCodeVerifyTooMany

type CodeService interface {
	Send(ctx context.Context, biz string, phone string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}
type codeService struct {
	sms  sms.Service
	repo repository.CodeRepository
}

func NewCodeService(sms sms.Service, repo repository.CodeRepository) CodeService {
	return &codeService{
		sms:  sms,
		repo: repo,
	}

}

// Send  短信发送功能  biz代表业务（验证码可用于登录，风险验证等）这都属于不同的业务 code 需要自己生成一个
func (c *codeService) Send(ctx context.Context, biz string, phone string) error {
	// 生成一个6位数的验证码以（0开头，06%d）
	code := c.generate()
	err := c.repo.Send(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	return c.sms.Send(ctx, tplId, []string{code}, phone)

}
func (c *codeService) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	ok, err := c.repo.Verify(ctx, biz, phone, inputCode)
	if err == repository.ErrCodeVerifyTooMany {
		// 这里相当于屏蔽了验证次数太多的错误 我们就是告诉调用者这个不对
		return false, nil
	}
	return ok, nil

}
func (c *codeService) generate() string {
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}
