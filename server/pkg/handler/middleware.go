package handler

import (
	"context"
	"fmt"
	"grpc/server/pkg/service"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type contextKey string

const userIDKey contextKey = "user_id"

func UnaryAuthInterceptor(service *service.Service) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		if info.FullMethod == "/proto.UserService/SignUp" ||
			info.FullMethod == "/proto.UserService/SignIn" {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, fmt.Errorf("missing metadata")
		}

		tokens := md.Get("authorization")
		if len(tokens) == 0 {
			return nil, fmt.Errorf("missing token")
		}

		tokenStr := strings.TrimPrefix(tokens[0], "Bearer ")

		userID, err := service.Authorization.ParseToken(tokenStr)
		if err != nil {
			return nil, fmt.Errorf("invalid token: %v", err)
		}

		newCtx := context.WithValue(ctx, userIDKey, userID)
		return handler(newCtx, req)
	}
}

func UserIDFromContext(ctx context.Context) (uint, error) {
	id, ok := ctx.Value(userIDKey).(uint)
	if !ok {
		return 0, fmt.Errorf("user_id not found in context")
	}
	return id, nil
}
