package interceptor

import (
	"bag-service/env"
	"bag-service/src/common/constant"
	"context"
	"encoding/base64"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Unary struct{}

func NewUnary() *Unary {
	return &Unary{}
}

func (u *Unary) AddBasicAuth(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(map[string]string{})
	}
	auth := base64.StdEncoding.EncodeToString([]byte(env.Conf.ProductService.GrpcAuth))
	md.Append("authorization", "Bearer "+auth)
	ctx = metadata.NewOutgoingContext(ctx, md)
	return invoker(ctx, method, req, reply, cc, opts...)
}
func (u *Unary) AddMetaData(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(map[string]string{})
	}
	if requestID, ok := ctx.Value(constant.RequestID).(string); ok {
		md.Append("request_id", requestID)
	}
	if userID, ok := ctx.Value(constant.UserID).(string); ok {
		md.Append("user_id", userID)
	}

	md.Append("app_id", "order-service")

	ctx = metadata.NewOutgoingContext(ctx, md)
	return invoker(ctx, method, req, reply, cc, opts...)
}
