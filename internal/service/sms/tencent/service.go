package tencent

/**
 * @Description
 * @Date 2024/3/11 0:47
 **/
import (
	"context"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	sms2 "project/internal/service/sms"
)

type SmsService struct {
	// 腾讯云客户端
	client *sms.Client
	// appid  详见腾讯云官网
	appId string
	// 签名
	signature string
}

func NewSmsService(client *sms.Client, appId string, signature string) sms2.Service {
	return &SmsService{
		client:    client,
		appId:     appId,
		signature: signature,
	}

}

func (s *SmsService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	request := sms.NewSendSmsRequest()
	request.SetContext(ctx)
	/* 短信应用ID: 短信SdkAppId在 [短信控制台] 添加应用后生成的实际SdkAppId，示例如1400006666 */
	request.SmsSdkAppId = common.StringPtr(s.appId)
	/* 短信签名内容: 使用 UTF-8 编码，必须填写已审核通过的签名，签名信息可登录 [短信控制台] 查看 */
	request.SignName = common.StringPtr(s.signature)
	/* 用户的 session 内容: 可以携带用户侧 ID 等上下文信息，server 会原样返回 */
	request.SessionContext = common.StringPtr("")
	/* 模板参数: 若无模板参数，则设置为空*/
	request.TemplateParamSet = common.StringPtrs(args)
	/* 模板 ID: 必须填写已审核通过的模板 ID。模板ID可登录 [短信控制台] 查看 */
	request.TemplateId = common.StringPtr(tplId)
	/* 下发手机号码，采用 E.164 标准，+[国家或地区码][手机号]
	 * 示例如：+8613711112222， 其中前面有一个+号 ，86为国家码，13711112222为手机号，最多不要超过200个手机号*/
	request.PhoneNumberSet = common.StringPtrs(numbers)
	// 通过client对象调用想要访问的接口，需要传入请求对象
	response, err := s.client.SendSms(request)
	// 处理异常  直接抛出 不用处理  正常的话可以对错误进行断言

	if err != nil {
		fmt.Printf("An API error has returned: %s", err)
		return err
	}
	/* 官方的话直接将结构体进行序列化返回，这里可以继续对状态进行判断 */
	//b, _ := json.Marshal(response.Response)
	//// 打印返回的json字符串
	//fmt.Printf("%s", b)

	for _, statusPtr := range response.Response.SendStatusSet {
		// 因为SendStatusSet是一个切片,故需要遍历最后根据结构体的状态码判断是否成功
		if statusPtr == nil {
			// 不可能进来这里
			continue
		}
		status := *statusPtr
		if status.Code == nil || *(status.Code) != "Ok" {
			// 发送失败
			return fmt.Errorf("发送短信失败 code：%s,msg:%s", *status.Code, *status.Message)
		}

	}
	return nil

}
