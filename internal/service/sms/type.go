package sms

import "context"

/**
 * @Description
 * @Date 2024/3/10 22:25
 **/

//发送短信的抽象
// 屏蔽不同供应商的区别

type Service interface {
	Send(ctx context.Context, tplId string, args []string, numbers ...string) error
}

/*	Send(ctx context.Context, number string, appId string, signature string, tplId string,params []any)
	其中的appid  和signature是固定的，一个生产线这个东西是固定不动的*/
