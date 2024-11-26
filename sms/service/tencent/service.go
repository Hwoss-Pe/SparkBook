package tencent

import (
	"Webook/sms/service"
	"context"
	"fmt"
	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/slice"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"go.uber.org/zap"
)

type TencentService struct {
	client   *sms.Client
	appId    *string
	signName *string
}

func (s *TencentService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	req := sms.NewSendSmsRequest()
	req.PhoneNumberSet = toStringPtrSlice(numbers)
	req.SmsSdkAppId = s.appId

	req.SetContext(ctx)
	req.TemplateParamSet = toStringPtrSlice(args)
	req.TemplateId = ekit.ToPtr[string](tplId)
	req.SignName = s.signName
	resp, err := s.client.SendSms(req)
	zap.L().Debug("调用腾讯短信服务",
		zap.Any("req", req),
		zap.Any("resp", resp))
	if err != nil {
		return err
	}

	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "OK" {
		}
		return fmt.Errorf("发送失败，code: %s, 原因：%s",
			*status.Code, *status.Message)
	}

	return nil
}

func NewService(client *sms.Client, appId string, signName string) service.Service {
	return &TencentService{client: client, appId: &appId, signName: &signName}
}

func toStringPtrSlice(src []string) []*string {
	return slice.Map[string, *string](src, func(idx int, src string) *string {
		return &src
	})
}
