package domain

type Amount struct {
	//货币单位
	Currency string
	Total    int64
}
type Payment struct {
	Amt         Amount
	BizTradeNO  string
	Description string
	Status      PaymentStatus
	// 第三方那边返回的 ID
	TxnID string
}

func (s PaymentStatus) AsUint8() uint8 {
	return uint8(s)
}

type PaymentStatus uint8

// 支付成功，失败，退款
const (
	PaymentStatusUnknown = iota
	PaymentStatusInit
	PaymentStatusSuccess
	PaymentStatusFailed
	PaymentStatusRefund
)
