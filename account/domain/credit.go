package domain

//这里的就是一次成交，会携带多个用户支付，某个业务某个uid对应的账号对应金额对应类型

type Credit struct {
	Biz   string
	BizId int64
	Items []CreditItem
}

type CreditItem struct {
	Uid         int64
	Account     int64
	Amt         int64
	Currency    string
	AccountType AccountType
}
type AccountType uint8

func (a AccountType) AsUint8() uint8 {
	return uint8(a)
}

const (
	AccountTypeUnknown = iota
	AccountTypeReward
	AccountTypeSystem
)
