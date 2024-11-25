package domain

type AsyncSms struct {
	id      int64
	TplId   int64
	Args    []string
	Numbers []string
	//重试配置
	RetryMax int
}
