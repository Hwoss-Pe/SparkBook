package trace

import (
	"Webook/pkg/grpcx/interceptor"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

//这里用openTelemetry和Zipkin结合做链路追踪

func (b *OTELInterceptorBuilder) BuildUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	tracer := b.tracer
	if tracer == nil {
		tracer = otel.Tracer("Webook/pkg/grpcx")
	}
	propagator := b.propagator
	if propagator == nil {
		propagator = otel.GetTextMapPropagator()
	}
	attrs := []attribute.KeyValue{
		semconv.RPCSystemKey.String("grpc"),
		attribute.Key("rpc.grpc.kind").String("unary"),
		attribute.Key("rpc.component").String("server"),
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx = extract(ctx, propagator)
		ctx, span := tracer.Start(ctx, info.FullMethod,
			trace.WithAttributes(attrs...),
			trace.WithSpanKind(trace.SpanKindServer))
		defer func() {
			span.End()
		}()
		span.SetAttributes(
			semconv.RPCMethodKey.String(info.FullMethod),
			semconv.NetPeerNameKey.String(b.PeerName(ctx)),
			attribute.Key("net.peer.ip").String(b.PeerIP(ctx)),
		)
		defer func() {
			if err != nil {
				span.RecordError(err)
				if e := errors.FromError(err); e != nil {
					span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int64(int64(e.Code)))
				}
				span.SetStatus(codes.Error, err.Error())
			} else {
				span.SetStatus(codes.Ok, "OK")
			}
		}()
		return handler(ctx, req)
	}
}
func (b *OTELInterceptorBuilder) BuildUnaryClientInterceptor() grpc.UnaryClientInterceptor {

	tracer := b.tracer
	if tracer == nil {
		tracer = otel.GetTracerProvider().Tracer("Webook/pkg/grpcx")
	}
	propagator := b.propagator
	if propagator != nil {
		propagator = otel.GetTextMapPropagator()
	}
	attrs := []attribute.KeyValue{
		semconv.RPCSystemKey.String("grpc"),
		attribute.Key("rpc.grpc.kind").String("unary"),
		attribute.Key("rpc.component").String("client"),
	}

	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		newAttrs := append(attrs,
			semconv.RPCMethodKey.String(method),
			semconv.NetPeerNameKey.String(b.serviceName))
		ctx, span := tracer.Start(ctx, method,
			trace.WithSpanKind(trace.SpanKindClient),
			trace.WithAttributes(newAttrs...))
		ctx = inject(ctx, propagator)
		defer func() {
			if err != nil {
				span.RecordError(err)
				if e := errors.FromError(err); e != nil {
					span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int64(int64(e.Code)))
				}
				span.SetStatus(codes.Error, err.Error())
			} else {
				span.SetStatus(codes.Ok, "OK")
			}
			span.End()
		}()
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

type OTELInterceptorBuilder struct {
	tracer      trace.Tracer
	propagator  propagation.TextMapPropagator
	serviceName string
	interceptor.Builder
}

func NewOTELInterceptorBuilder(
	serviceName string,
	tracer trace.Tracer,
	propagator propagation.TextMapPropagator) *OTELInterceptorBuilder {
	return &OTELInterceptorBuilder{tracer: tracer,
		serviceName: serviceName, propagator: propagator}
}

func extract(ctx context.Context, propagators propagation.TextMapPropagator) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}

	return propagators.Extract(ctx, GrpcHeaderCarrier(md))
}
func inject(ctx context.Context, propagators propagation.TextMapPropagator) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	propagators.Inject(ctx, GrpcHeaderCarrier(md))
	return metadata.NewOutgoingContext(ctx, md)
}

// GrpcHeaderCarrier 的设计主要用于支持 gRPC 元数据的操作metadata.MD进行操作
type GrpcHeaderCarrier metadata.MD

// Get returns the value associated with the passed key.
func (mc GrpcHeaderCarrier) Get(key string) string {
	vals := metadata.MD(mc).Get(key)
	if len(vals) > 0 {
		return vals[0]
	}
	return ""
}

// Set stores the key-value pair.
func (mc GrpcHeaderCarrier) Set(key string, value string) {
	metadata.MD(mc).Set(key, value)
}

// Keys lists the keys stored in this carrier.
func (mc GrpcHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(mc))
	for k := range metadata.MD(mc) {
		keys = append(keys, k)
	}
	return keys
}
