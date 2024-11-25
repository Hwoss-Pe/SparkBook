package auth

import (
	sms "Webook/sms/service"
	"context"
	"github.com/golang-jwt/jwt/v5"
)

type SMSAuthService struct {
	sms sms.Service
	key []byte
}

func (s *SMSAuthService) Send(ctx context.Context, tplToken string, args []string, numbers ...string) error {
	var c Claims
	_, err := jwt.ParseWithClaims(tplToken, &c, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})
	if err != nil {
		return err
	}
	return s.sms.Send(ctx, tplToken, args, numbers...)
}

type Claims struct {
	jwt.RegisteredClaims
	Tpl string
}
