package service

import (
	smsv1 "Webook/api/proto/gen/api/proto/sms/v1"
	"Webook/code/repository"
	"context"
	"errors"
	"fmt"
	"math/rand"
)

var ErrCodeSendTooMany = repository.ErrCodeSendTooMany

// 假设这个是验证码模板id
const codeTplId = "1877556"

//go:generate mockgen -source=./code.go -package=svcmocks -destination=mocks/code.mock.go CodeService
type CodeService interface {
	Send(ctx context.Context, biz string, phone string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}

type SMSCodeService struct {
	sms  smsv1.SmsServiceClient
	repo repository.CodeRepository
}

func (s *SMSCodeService) Send(ctx context.Context, biz string, phone string) error {
	//	实际上就是通知短信给某个手机号发送生成好验证码
	code := s.generate()
	err := s.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	//返回数据的容错逻辑应该在sms那边处理好，不在这里关心
	_, err = s.sms.Send(ctx, &smsv1.SmsSendRequest{
		TplId:   codeTplId,
		Args:    []string{code},
		Numbers: []string{phone},
	})
	return err
}

func (s *SMSCodeService) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	ok, err := s.repo.Verify(ctx, biz, phone, inputCode)
	if errors.Is(err, repository.ErrCodeVerifyTooManyTimes) {
		// 不正常校验验证码，这里可以告警也可以直接让他稍后在输入
		// 待完善
		return false, nil
	}
	return ok, err
}

func NewSMSCodeService(svc smsv1.SmsServiceClient, repo repository.CodeRepository) CodeService {
	return &SMSCodeService{
		sms:  svc,
		repo: repo,
	}
}
func (s *SMSCodeService) generate() string {
	// 用随机数生成一个
	num := rand.Intn(999999)
	return fmt.Sprintf("%6d", num)
}
