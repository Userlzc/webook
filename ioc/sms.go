package ioc

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentSMS "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"os"
	sms2 "project/internal/service/sms"
	"project/internal/service/sms/tencent"
)

/**
 * @Description
 * @Date 2024/3/12 19:08
 **/

func InitSmsService() sms2.Service {
	return nil

}
func initTencentSMSService() sms2.Service {
	secretId, ok := os.LookupEnv("SMS_SECRET_ID")
	if !ok {
		panic("找不到腾讯的secret_id")
	}
	secretKey, ok := os.LookupEnv("SMS_SECRET_KEY")
	if !ok {
		panic("找不到腾讯的secret_key")
	}
	c, err := tencentSMS.NewClient(common.NewCredential(secretId, secretKey), "api-nanjing", profile.NewClientProfile())
	if err != nil {
		panic(err)
	}
	return tencent.NewSmsService(c, "14088888", "测试科技")

}
