package grpc

import (
	userv1 "Webook/api/proto/gen/api/proto/user/v1"
	"Webook/user/domain"
	"Webook/user/service"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserServiceServer struct {
	service service.UserService
	userv1.UnimplementedUsersServiceServer
}

func (u *UserServiceServer) Signup(ctx context.Context, request *userv1.SignupRequest) (*userv1.SignupResponse, error) {
	user := convertToDomain(request.User)
	err := u.service.Signup(ctx, user)
	return &userv1.SignupResponse{}, err
}

func (u *UserServiceServer) FindOrCreate(ctx context.Context, request *userv1.FindOrCreateRequest) (*userv1.FindOrCreateResponse, error) {
	user, err := u.service.FindOrCreate(ctx, request.Phone)
	return &userv1.FindOrCreateResponse{
		User: convertToV(user),
	}, err
}

func (u *UserServiceServer) Login(ctx context.Context, request *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	user, err := u.service.Login(ctx, request.GetEmail(), request.GetPassword())
	return &userv1.LoginResponse{
		User: convertToV(user),
	}, err
}

func (u *UserServiceServer) Profile(ctx context.Context, request *userv1.ProfileRequest) (*userv1.ProfileResponse, error) {
	user, err := u.service.Profile(ctx, request.GetId())
	return &userv1.ProfileResponse{
		User: convertToV(user),
	}, err
}

func (u *UserServiceServer) UpdateNonSensitiveInfo(ctx context.Context, request *userv1.UpdateNonSensitiveInfoRequest) (*userv1.UpdateNonSensitiveInfoResponse, error) {
	err := u.service.UpdateNonSensitiveInfo(ctx, convertToDomain(request.User))
	return &userv1.UpdateNonSensitiveInfoResponse{}, err
}

func (u *UserServiceServer) FindOrCreateByWechat(ctx context.Context, request *userv1.FindOrCreateByWechatRequest) (*userv1.FindOrCreateByWechatResponse, error) {
	user, err := u.service.FindOrCreateByWechat(ctx, domain.WechatInfo{
		OpenId:  request.GetInfo().GetOpenId(),
		UnionId: request.GetInfo().GetUnionId(),
	})
	return &userv1.FindOrCreateByWechatResponse{
		User: convertToV(user),
	}, err
}

func NewUserServiceServer(userService service.UserService) *UserServiceServer {
	return &UserServiceServer{service: userService}
}

func (u *UserServiceServer) Register(server grpc.ServiceRegistrar) {
	//注册的时候需要先实现生成的UserServiceServer
	userv1.RegisterUsersServiceServer(server, u)
}

func convertToDomain(u *userv1.User) domain.User {
	du := domain.User{}
	du.Id = u.GetId()
	du.Nickname = u.GetNickname()
	du.Password = u.GetPassword()
	du.Phone = u.GetPhone()
	du.AboutMe = u.GetAboutMe()
	du.Ctime = u.GetCtime().AsTime()
	du.WechatInfo = domain.WechatInfo{
		OpenId:  u.GetWechatInfo().GetOpenId(),
		UnionId: u.GetWechatInfo().GetUnionId(),
	}
	return du
}

func convertToV(user domain.User) *userv1.User {
	vUser := &userv1.User{
		Id:       user.Id,
		Email:    user.Email,
		Nickname: user.Nickname,
		Password: user.Password,
		Phone:    user.Phone,
		AboutMe:  user.AboutMe,
		Ctime:    timestamppb.New(user.Ctime),
		Birthday: timestamppb.New(user.Birthday),
		WechatInfo: &userv1.WechatInfo{
			OpenId:  user.WechatInfo.OpenId,
			UnionId: user.WechatInfo.UnionId,
		},
	}
	return vUser
}
