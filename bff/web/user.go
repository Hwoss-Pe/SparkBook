package web

import (
	codev1 "Webook/api/proto/gen/api/proto/code/v1"
	userv1 "Webook/api/proto/gen/api/proto/user/v1"
	"regexp"
)

const (
	emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	userIdKey            = "userId"
	bizLogin             = "login"
)

type UserHandler struct {
	userSvc          userv1.UsersServiceClient
	codeSvc          codev1.CodeServiceClient
	emailRegexExp    *regexp.Regexp
	passwordRegexExp *regexp.Regexp
}
