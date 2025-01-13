package domain

import (
	"github.com/robfig/cron/v3"
	"time"
)

type CronJob struct {
	Id   int64
	Name string // 全局唯一
	//	用什么来运行
	Executor   string
	Cfg        string
	Expression string
	NextTime   time.Time
	//是否放弃抢占
	CancelFunc func()
}

func (j CronJob) Next(t time.Time) time.Time {
	// 这个地方 Expression 必须不能出错，这需要用户在注册的时候确保
	expr := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom |
		cron.Month | cron.Dow |
		cron.Descriptor)
	s, _ := expr.Parse(j.Expression)
	//利用解析器可以去创建对应的Scheduler，然后根据这个去查找下一触发时间点
	return s.Next(t)
}
