package interceptor

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
	"user-service/env"
	"user-service/src/common/dto/request"
	"user-service/src/common/errors"
	"user-service/src/common/log"
	"user-service/src/factory"
	"user-service/src/publisher"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Unary struct {
	grpcLogPublisher *publisher.Kafka
}

func NewUnary(f *factory.Factory) *Unary {
	return &Unary{
		grpcLogPublisher: f.GrpcLogPublisher,
	}
}

func (u *Unary) Error(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	res, err := handler(ctx, req)

	if err != nil {
		if e, ok := err.(*errors.Response); ok {
			return nil, status.Errorf(e.GRPCCode, "%s", e.Message)
		}

		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	return res, nil
}

func (u *Unary) Recovery(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Logger.Error(fmt.Sprintf("panic: %v", r))

			resp = nil
			err = status.Errorf(codes.Internal, "internal server error")
		}
	}()
	res, err := handler(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u *Unary) BasicAuthValidator(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing authorization token")
	}

	authHeader := values[0]
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, status.Errorf(codes.Unauthenticated, "invalid authorization token")
	}
	encodedToken := strings.TrimPrefix(authHeader, "Bearer ")
	decoded, err := base64.StdEncoding.DecodeString(encodedToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid base64 token")
	}

	if string(decoded) != env.Conf.CurrentApp.GrpcBasicAuth {
		return nil, status.Errorf(codes.Unauthenticated, "invalid basic auth token")
	}
	return handler(ctx, req)
}

func (u *Unary) Logging(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	now := time.Now()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Internal, "missing metadata")
	}
	res, err := handler(ctx, req)

	var errMessage string
	if err != nil {
		errMessage = err.Error()
	}

	var requestID *string
	if requestIDs, ok := md["request_id"]; ok && len(requestIDs) > 0 {
		requestID = &requestIDs[0]
	}

	var userID *string
	if userIDs, ok := md["user_id"]; ok && len(userIDs) > 0 {
		userID = &userIDs[0]
	}
	var appID *string
	if appIDs, ok := md["app_id"]; ok && len(appIDs) > 0 {
		appID = &appIDs[0]
	}

	st, _ := status.FromError(err)

	u.grpcLogPublisher.Publish(request.GrpcLog{
		Timestamp:  time.Now().Format(time.RFC3339),
		RequestID:  requestID,
		UserID:     userID,
		AppID:      appID,
		Method:     info.FullMethod,
		Body:       req,
		StatusCode: st.Code(),
		Latency:    float64(time.Since(now).Milliseconds()),
		Error:      errMessage,
	})
	return res, err
}
