package otel

import (
	sms "Webook/sms/service"
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	svc    sms.Service
	tracer trace.Tracer
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	ctx, span := s.tracer.Start(ctx, "sms_send")
	defer span.End()
	span.SetAttributes(attribute.String("tplId", tplId))
	err := s.svc.Send(ctx, tplId, args, numbers...)
	if err != nil {
		span.RecordError(err)
	}
	return err
}

func NewService(svc sms.Service) *Service {
	return &Service{
		svc:    svc,
		tracer: otel.GetTracerProvider().Tracer("sms_service"),
	}
}
