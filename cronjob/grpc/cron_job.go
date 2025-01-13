package grpc2

import (
	cronjobv1 "Webook/api/proto/gen/api/proto/cornjob/v1"
	"Webook/cronjob/domain"
	"Webook/cronjob/service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CronJobServiceServer struct {
	svc service.CronJobService
	cronjobv1.UnimplementedCronJobServiceServer
}

func (c *CronJobServiceServer) Preempt(ctx context.Context, request *cronjobv1.PreemptRequest) (*cronjobv1.PreemptResponse, error) {
	job, err := c.svc.Preempt(ctx)
	return &cronjobv1.PreemptResponse{
		Cronjob: convertToV(job),
	}, err
}
func (c *CronJobServiceServer) Register(server grpc.ServiceRegistrar) {
	cronjobv1.RegisterCronJobServiceServer(server, c)
}

func (c *CronJobServiceServer) ResetNextTime(ctx context.Context, request *cronjobv1.ResetNextTimeRequest) (*cronjobv1.ResetNextTimeResponse, error) {
	err := c.svc.ResetNextTime(ctx, convertToDomain(request.Job))
	return &cronjobv1.ResetNextTimeResponse{}, err
}

func (c *CronJobServiceServer) AddJob(ctx context.Context, request *cronjobv1.AddJobRequest) (*cronjobv1.AddJobResponse, error) {
	err := c.svc.AddJob(ctx, convertToDomain(request.Job))
	return &cronjobv1.AddJobResponse{}, err
}

func NewCronJobServiceServer(svc service.CronJobService) *CronJobServiceServer {
	return &CronJobServiceServer{
		svc: svc,
	}
}

func convertToV(domainCronJob domain.CronJob) *cronjobv1.CronJob {
	return &cronjobv1.CronJob{
		Id:         domainCronJob.Id,
		Name:       domainCronJob.Name,
		Executor:   domainCronJob.Executor,
		Cfg:        domainCronJob.Cfg,
		Expression: domainCronJob.Expression,
		NextTime:   timestamppb.New(domainCronJob.NextTime),
	}
}

func convertToDomain(vCronJob *cronjobv1.CronJob) domain.CronJob {
	return domain.CronJob{
		Id:         vCronJob.GetId(),
		Name:       vCronJob.GetName(),
		Executor:   vCronJob.GetExecutor(),
		Cfg:        vCronJob.GetCfg(),
		Expression: vCronJob.GetExpression(),
		NextTime:   vCronJob.GetNextTime().AsTime(),
	}
}
