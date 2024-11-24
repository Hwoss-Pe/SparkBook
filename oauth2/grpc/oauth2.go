package grpc

import (
	oauth2v1 "Webook/api/proto/gen/api/proto/oauth2/v1"
	"Webook/oauth2/service/wechat"
	"context"
	"google.golang.org/grpc"
)

type Oauth2ServiceServer struct {
	oauth2v1.UnimplementedOauth2ServiceServer
	service wechat.Service
}

func (o *Oauth2ServiceServer) AuthURL(ctx context.Context, request *oauth2v1.AuthURLRequest) (*oauth2v1.AuthURLResponse, error) {
	url, err := o.service.AuthURL(ctx, request.State)
	if err != nil {
		return nil, err
	}
	return &oauth2v1.AuthURLResponse{
		Url: url,
	}, nil
}

func (o *Oauth2ServiceServer) VerifyCode(ctx context.Context, request *oauth2v1.VerifyCodeRequest) (*oauth2v1.VerifyCodeResponse, error) {
	info, err := o.service.VerifyCode(ctx, request.Code)
	if err != nil {
		return nil, err
	}
	return &oauth2v1.VerifyCodeResponse{
		OpenId:  info.OpenId,
		UnionId: info.UnionId,
	}, nil

}

func (o *Oauth2ServiceServer) Register(server grpc.ServiceRegistrar) {
	oauth2v1.RegisterOauth2ServiceServer(server, o)
}
func NewOauth2ServiceServer(svc wechat.Service) *Oauth2ServiceServer {
	return &Oauth2ServiceServer{
		service: svc,
	}
}
