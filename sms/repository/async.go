package repository

import (
	"Webook/sms/domain"
	"Webook/sms/repository/dao"
	"context"
)

var ErrWaitingSMSNotFound = dao.ErrWaitingSMSNotFound

//go:generate mockgen -source=./async.go -package=repomocks -destination=mocks/async_sms_repository.mock.go AsyncSmsRepository
type AsyncSmsRepository interface {
	//先抢占，后进行发送，查询结果
	Add(ctx context.Context, s domain.AsyncSms) error
	PreemptWaitingSMS(ctx context.Context) (domain.AsyncSms, error)
	ReportScheduleResult(ctx context.Context, id int64, success bool) error
}
