package middleware

import (
	"encoding/gob"
	"github.com/ecodeclub/ekit/set"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	publicPaths set.Set[string]
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	s := set.NewMapSet[string](4)
	s.Add("users/signup")
	s.Add("users/login_sms/code/send")
	s.Add("users/login_sms")
	s.Add("users/login")
	return &LoginMiddlewareBuilder{
		publicPaths: s,
	}
}

func (l *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	gob.Register(time.Time{})
	return func(ctx *gin.Context) {
		if l.publicPaths.Exist(ctx.Request.URL.Path) {
			return
		}
		sess := sessions.Default(ctx)
		if sess.Get("userId") == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		const timeKey = "update_time"
		val := sess.Get(timeKey)
		updateTime, ok := val.(time.Time)
		if val == nil || ok && time.Now().Sub(updateTime) > time.Second*10 {
			sess.Options(sessions.Options{
				MaxAge: 60,
			})
			sess.Set(timeKey, time.Now())
			if err := sess.Save(); err != nil {
				panic(err)
			}
		}
	}
}
