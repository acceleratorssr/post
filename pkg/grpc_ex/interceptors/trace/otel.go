package trace

import (
	"context"
	"errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"post/pkg/grpc_ex/interceptors"
)

type OTELInterceptorBuilder struct {
	tracer     trace.Tracer
	propagator propagation.TextMapPropagator
	interceptors.Builder
	serviceName string
}

func NewOTELInterceptorBuilder(
	serviceName string,
	tracer trace.Tracer,
	propagator propagation.TextMapPropagator) *OTELInterceptorBuilder {
	return &OTELInterceptorBuilder{tracer: tracer,
		serviceName: serviceName, propagator: propagator}
}

func (b *OTELInterceptorBuilder) BuildUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	tracer := b.tracer
	if tracer == nil {
		tracer = otel.Tracer("post/pkg/grpc_ex")
	}

	propagator := b.propagator
	if propagator == nil {
		propagator = otel.GetTextMapPropagator() // 全局的propagator
	}

	attrs := []attribute.KeyValue{
		semconv.RPCSystemKey.String("grpc"),
		attribute.Key("rpc.grpc.kind").String("unary"),
		attribute.Key("rpc.component").String("server"),
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (reply interface{}, err error) {
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

				var grpcErr interface{ GRPCStatus() *status.Status }
				if errors.As(err, &grpcErr) {
					span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int64(int64(grpcErr.GRPCStatus().Code())))
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
		tracer = otel.GetTracerProvider().
			Tracer("post/pkg/grpcx")
	}

	propagator := b.propagator
	if propagator == nil {
		propagator = otel.GetTextMapPropagator()
	}

	attrs := []attribute.KeyValue{
		semconv.RPCSystemKey.String("grpc"),
		attribute.Key("rpc.grpc.kind").String("unary"),
		attribute.Key("rpc.component").String("client"),
	}

	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
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

				var grpcErr interface{ GRPCStatus() *status.Status }
				if errors.As(err, &grpcErr) {
					span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int64(int64(grpcErr.GRPCStatus().Code())))
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

// 获取客户端的链路元数据
func extract(ctx context.Context, propagators propagation.TextMapPropagator) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}

	// propagators.Extract 用于从传输介质（如 gRPC 头部）中提取上下文信息
	// GrpcHeaderCarrier 是一个适配器，用于将 gRPC 元数据包装成 TextMapCarrier，（实现它的方法）
	// 使得 propagators 能够从中提取信息
	// 最后注入到ctx内
	return propagators.Extract(ctx, GrpcHeaderCarrier(md))
}

func inject(ctx context.Context, propagators propagation.TextMapPropagator) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}

	// 使用 propagators.Inject 将信息注入到 metadata.MD 中
	// 具体注入的消息和具体的trace有关
	propagators.Inject(ctx, GrpcHeaderCarrier(md))

	// 创建一个新的上下文，将注入了上下文信息的元数据 md 作为传出的元数据附加到这个上下文中
	return metadata.NewOutgoingContext(ctx, md)
}

// GrpcHeaderCarrier ...
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
