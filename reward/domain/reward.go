package domain

type Target struct {
	//具体什么业务
	Biz   string
	BizId int64
	//打赏的名称
	BizName string
	//打赏给谁
	Uid int64
}

type Reward struct {
	Id     int64
	Uid    int64
	Target Target
	Amt    int64
	Status RewardStatus
}
type CodeURL struct {
	URL string
	Rid int64
}

func (r RewardStatus) AsUint8() uint8 {
	return uint8(r)
}

// Completed 判断是否该订单已经结束
func (r Reward) Completed() bool {
	return r.Status == RewardStatusFailed || r.Status == RewardStatusPayed
}

type RewardStatus uint8

const (
	RewardStatusUnknown = iota
	RewardStatusInit
	RewardStatusPayed
	RewardStatusFailed
)
