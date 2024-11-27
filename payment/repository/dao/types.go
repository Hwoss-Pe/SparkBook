package dao

import (
	"Webook/payment/domain"
	"context"
	"database/sql"
	"time"
)

type PaymentDAO interface {
	Insert(ctx context.Context, pmt Payment) error
	UpdateTxnIDAndStatus(ctx context.Context, bizTradeNo string, txnID string, status domain.PaymentStatus) error
	FindExpiredPayment(ctx context.Context, offset int, limit int, t time.Time) ([]Payment, error)
	GetPayment(ctx context.Context, bizTradeNO string) (Payment, error)
}

type Payment struct {
	Id          int64 `gorm:"primaryKey,autoIncrement" bson:"id,omitempty"`
	Amt         int64
	Currency    string
	Description string `gorm:"description"`

	// 业务方传过来的
	BizTradeNO string `gorm:"column:biz_trade_no;type:varchar(256);unique"`

	// 第三方支付平台的事务 ID，唯一的
	TxnID sql.NullString `gorm:"column:txn_id;type:varchar(128);unique"`

	Status uint8
	Utime  int64
	Ctime  int64
}
