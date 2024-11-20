package domain

import (
	"time"
)

type User struct {
	Id       int64
	Email    string
	Nickname string
	Password string
	Phone    string
	AboutMe  string
	Ctime    time.Time
	Birthday time.Time
	//	可以使用微信扫码登录，这里封装一下对应的登录api所用参数
	WechatInfo WechatInfo
}

type WechatInfo struct {
	//对应的应用内唯一
	OpenId string
	//公司账号内唯一
	UnionId string
}
