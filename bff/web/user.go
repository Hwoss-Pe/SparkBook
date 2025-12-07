package web

import (
	codev1 "Webook/api/proto/gen/api/proto/code/v1"
	userv1 "Webook/api/proto/gen/api/proto/user/v1"
	jwt2 "Webook/bff/web/jwt"
	"Webook/pkg/ginx"
	"Webook/user/errs"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"time"
)

const (
	emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$` //密码必须包含数字、特殊字符，并且长度不能小于 8 位
	userIdKey            = "userId"                                                           //给session的
	bizLogin             = "login"
)

type UserHandler struct {
	userSvc          userv1.UsersServiceClient
	codeSvc          codev1.CodeServiceClient
	emailRegexExp    *regexp.Regexp
	passwordRegexExp *regexp.Regexp
	jwt2.Handler
}

func (c *UserHandler) RegisterRoute(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", ginx.WrapReq[SignUpReq](c.Signup))
	// session 机制
	//ug.POST("/login", c.Login)
	ug.POST("/login", ginx.WrapReq[LoginReq](c.LoginJWT))
	ug.POST("/logout", c.Logout)
	ug.POST("/edit", c.Edit)
	//ug.GET("/profile", c.Profile)
	ug.GET("/profile", c.ProfileJWT)
	ug.POST("/login_sms/code/send", c.SendSMSLoginCode)
	ug.POST("/login_sms", c.LoginSMS)
	ug.POST("/refresh_token", c.RefreshToken)
}

func (c *UserHandler) Signup(ctx *gin.Context, req SignUpReq) (ginx.Result, error) {
	isEmail, err := c.emailRegexExp.MatchString(req.Email)
	if err != nil {
		return Result{
			Code: errs.UserInternalServerError,
			Msg:  "系统错误",
		}, err
	}
	if !isEmail {
		return Result{
			Code: errs.UserInvalidInput,
			Msg:  "邮箱输入错误",
		}, nil
	}
	if req.Password != req.ConfirmPassword {
		return Result{
			Code: errs.UserInvalidInput,
			Msg:  "两次输入密码不对",
		}, nil
	}
	isPassword, err := c.passwordRegexExp.MatchString(req.Password)
	if err != nil {
		return Result{
			Code: errs.UserInvalidInput,
			Msg:  "系统错误",
		}, err
	}
	if !isPassword {
		return Result{
			Code: errs.UserInvalidInput,
			Msg:  "密码必须包含数字、特殊字符，并且长度不能小于 8 位",
		}, nil
	}
	_, err = c.userSvc.Signup(ctx.Request.Context(), &userv1.SignupRequest{
		User: &userv1.User{
			Email:    req.Email,
			Password: req.Password,
		},
	})
	if err != nil {
		return Result{
			Code: errs.UserInternalServerError,
			Msg:  "系统错误",
		}, err
	}
	return Result{
		Msg: "OK",
	}, nil
}

// LoginJWT 使用的是 JWT，
func (c *UserHandler) LoginJWT(ctx *gin.Context, req LoginReq) (ginx.Result, error) {
	u, err := c.userSvc.Login(ctx.Request.Context(), &userv1.LoginRequest{
		Email: req.Email, Password: req.Password,
	})
	if err != nil {
		return ginx.Result{}, err
	}
	err = c.SetLoginToken(ctx, u.User.Id) //登录后设置token
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{Msg: "登录成功"}, nil
}

func (c *UserHandler) Logout(ctx *gin.Context) {
	err := c.ClearToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
}

func NewUserHandler(userSvc userv1.UsersServiceClient, codeSvc codev1.CodeServiceClient, jhl jwt2.Handler) *UserHandler {
	return &UserHandler{
		Handler:          jhl,
		userSvc:          userSvc,
		codeSvc:          codeSvc,
		emailRegexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
	}
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type SignUpReq struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

// Login 这个用的是session机制
func (c *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	// 调用 Bind 方法的时候，如果有问题，Bind 方法已经直接写响应回去了
	if err := ctx.Bind(&req); err != nil {
		return
	}
	u, err := c.userSvc.Login(ctx.Request.Context(), &userv1.LoginRequest{
		Email: req.Email, Password: req.Password})
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	sess := sessions.Default(ctx)
	sess.Set(userIdKey, u.User.Id)
	sess.Options(sessions.Options{
		// 60 秒过期
		MaxAge: 60,
	})
	err = sess.Save()
	if err != nil {
		ctx.String(http.StatusOK, "服务器异常")
		return
	}
	ctx.String(http.StatusOK, "登录成功")
}

func (c *UserHandler) Edit(ctx *gin.Context) {
	type Req struct {
		// 注意，其它字段，尤其是密码、邮箱和手机，
		// 修改都要通过别的手段
		// 邮箱和手机都要验证
		// 密码更加不用多说了
		Nickname string `json:"nickname"`
		// 2023-01-01
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
		Avatar   string `json:"avatar"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	if req.Nickname == "" {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "昵称不能为空"})
		return
	}

	if len(req.AboutMe) > 1024 {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "关于 过长"})
		return
	}
	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "日期格式不对"})
		return
	}
	//	校验token
	uc := ctx.MustGet("user").(jwt2.UserClaims)
	_, err = c.userSvc.UpdateNonSensitiveInfo(ctx,
		&userv1.UpdateNonSensitiveInfoRequest{
			User: &userv1.User{
				Id:       uc.Id,
				Nickname: req.Nickname,
				AboutMe:  req.AboutMe,
				Avatar:   req.Avatar,
				Birthday: timestamppb.New(birthday),
			},
		})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	ctx.JSON(http.StatusOK, Result{Msg: "OK"})
}

// ProfileJWT jwt版本
func (c *UserHandler) ProfileJWT(ctx *gin.Context) {
	type Profile struct {
		Email    string
		Phone    string
		Nickname string
		Birthday string
		AboutMe  string
		Avatar   string
	}
	uc := ctx.MustGet("user").(jwt2.UserClaims)
	resp, err := c.userSvc.Profile(ctx, &userv1.ProfileRequest{Id: uc.Id})
	if err != nil {
		// 按照道理来说，这边 id 对应的数据肯定存在，所以要是没找到，
		// 那就说明是系统出了问题。
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	u := resp.User
	ctx.JSON(http.StatusOK, Profile{
		Email:    u.Email,
		Phone:    u.Phone,
		Nickname: u.Nickname,
		Birthday: u.Birthday.AsTime().Format(time.DateOnly),
		AboutMe:  u.AboutMe,
		Avatar:   u.Avatar,
	})
}
func (c *UserHandler) Profile(ctx *gin.Context) {
	type Profile struct {
		Email string
	}
	sess := sessions.Default(ctx)
	id := sess.Get(userIdKey).(int64)
	u, err := c.userSvc.Profile(ctx, &userv1.ProfileRequest{
		Id: id,
	})
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, Profile{
		Email: u.User.Email,
	})
}

// SendSMSLoginCode 发送验证码接口
func (c *UserHandler) SendSMSLoginCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "请输入手机号码"})
		return
	}
	_, err := c.codeSvc.Send(ctx, &codev1.CodeSendRequest{
		Biz: bizLogin, Phone: req.Phone,
	})
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{Msg: "发送成功"})
	//case .ErrCodeSendTooMany:
	//	ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "短信发送太频繁，请稍后再试"})
	default:
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
}

func (c *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	resp, err := c.codeSvc.Verify(ctx, &codev1.VerifyRequest{
		Biz: bizLogin, Phone: req.Phone, InputCode: req.Code,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统异常"})
		zap.L().Error("用户手机号码登录失败", zap.Error(err))
		return
	}
	if !resp.Answer {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "验证码错误"})
		return
	}

	u, err := c.userSvc.FindOrCreate(ctx, &userv1.FindOrCreateRequest{
		Phone: req.Phone,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "系统错误"})
		return
	}
	// 用 uuid 来标识这一次会话
	ssid := uuid.New().String()
	err = c.SetJWTToken(ctx, ssid, u.User.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Msg: "系统错误"})
		return
	}
	ctx.JSON(http.StatusOK, Result{Msg: "登录成功"})
}

func (c *UserHandler) RefreshToken(ctx *gin.Context) {
	// 长 token 也放在这里
	tokenStr := c.ExtractTokenString(ctx)
	var rc jwt2.RefreshClaims
	token, err := jwt.ParseWithClaims(tokenStr, &rc, func(token *jwt.Token) (interface{}, error) {
		return jwt2.RefreshTokenKey, nil
	})
	if err != nil || token == nil || !token.Valid {
		ctx.JSON(http.StatusUnauthorized, Result{Code: 4, Msg: "请登录"})
		return
	}
	// 校验 ssid
	err = c.CheckSession(ctx, rc.Ssid)
	if err != nil {
		// 系统错误或者用户已经主动退出登录了
		// 这里也可以考虑说，如果在 Redis 已经崩溃的时候，
		// 就不要去校验是不是已经主动退出登录了。
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	err = c.SetJWTToken(ctx, rc.Ssid, rc.Id)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, Result{Code: 4, Msg: "请登录"})
		return
	}
	ctx.JSON(http.StatusOK, Result{Msg: "刷新成功"})
}
