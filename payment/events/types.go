package events

// PaymentEvent 也是最简设计
// 后续如果要接入大数据之类的，那么就可以考虑提供 payment 详情
type PaymentEvent struct {
	BizTradeNO string
	Status     uint8
}

func (PaymentEvent) Topic() string {
	return "payment_events"
}
