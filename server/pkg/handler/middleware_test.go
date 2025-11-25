package handler_test

import (
	"context"
	"errors"
	"grpc/server/pkg/handler"
	"grpc/server/pkg/service"
	mock_service "grpc/server/pkg/service/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ctxWithMetadata(token string) context.Context {
	return metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+token))
}

func fakeHandler(ctx context.Context, req interface{}) (interface{}, error) {
	return ctx.Value("ok"), nil
}

func TestUnaryAuthInterceptor_SignUpPassThrough(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mock_service.NewMockAuthorization(ctrl)
	srv := &service.Service{Authorization: mockAuth}

	interceptor := handler.UnaryAuthInterceptor(srv)

	info := &grpc.UnaryServerInfo{FullMethod: "/proto.UserService/SignUp"}
	ctx := context.WithValue(context.Background(), "ok", true)

	resp, err := interceptor(ctx, nil, info, fakeHandler)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != true {
		t.Fatal("expected handler to run")
	}
}

func TestUnaryAuthInterceptor_SignInPassThrough(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mock_service.NewMockAuthorization(ctrl)
	srv := &service.Service{Authorization: mockAuth}

	interceptor := handler.UnaryAuthInterceptor(srv)

	info := &grpc.UnaryServerInfo{FullMethod: "/proto.UserService/SignIn"}
	ctx := context.WithValue(context.Background(), "ok", true)

	resp, err := interceptor(ctx, nil, info, fakeHandler)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != true {
		t.Fatal("expected handler to run")
	}
}

func TestUnaryAuthInterceptor_MissingMetadata(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mock_service.NewMockAuthorization(ctrl)
	srv := &service.Service{Authorization: mockAuth}

	interceptor := handler.UnaryAuthInterceptor(srv)
	info := &grpc.UnaryServerInfo{FullMethod: "/proto.BookService/GetBook"}

	_, err := interceptor(context.Background(), nil, info, fakeHandler)
	if err == nil || err.Error() != "missing metadata" {
		t.Fatalf("expected missing metadata error, got %v", err)
	}
}

func TestUnaryAuthInterceptor_MissingToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mock_service.NewMockAuthorization(ctrl)
	srv := &service.Service{Authorization: mockAuth}

	interceptor := handler.UnaryAuthInterceptor(srv)
	info := &grpc.UnaryServerInfo{FullMethod: "/proto.BookService/GetBook"}

	ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{}) // нет токена

	_, err := interceptor(ctx, nil, info, fakeHandler)
	if err == nil || err.Error() != "missing token" {
		t.Fatalf("expected missing token error, got %v", err)
	}
}

func TestUnaryAuthInterceptor_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mock_service.NewMockAuthorization(ctrl)
	srv := &service.Service{Authorization: mockAuth}

	interceptor := handler.UnaryAuthInterceptor(srv)
	info := &grpc.UnaryServerInfo{FullMethod: "/proto.BookService/GetBook"}

	ctx := ctxWithMetadata("badtoken")

	mockAuth.EXPECT().
		ParseToken("badtoken").
		Return(uint(0), errors.New("invalid"))

	_, err := interceptor(ctx, nil, info, fakeHandler)
	if err == nil || err.Error() != "invalid token: invalid" {
		t.Fatalf("expected invalid token error, got %v", err)
	}
}

func TestUnaryAuthInterceptor_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mock_service.NewMockAuthorization(ctrl)
	srv := &service.Service{Authorization: mockAuth}

	interceptor := handler.UnaryAuthInterceptor(srv)
	info := &grpc.UnaryServerInfo{FullMethod: "/proto.BookService/GetBook"}

	ctx := ctxWithMetadata("goodtoken")

	mockAuth.EXPECT().
		ParseToken("goodtoken").
		Return(uint(42), nil)

	handlerFn := func(ctx context.Context, req interface{}) (interface{}, error) {
		// Проверяем, что userID положен в контекст
		userID := ctx.Value(handler.UserIDKey())
		if userID != uint(42) {
			t.Fatalf("expected userID=42 in context, got %v", userID)
		}
		return "ok", nil
	}

	resp, err := interceptor(ctx, nil, info, handlerFn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != "ok" {
		t.Fatalf("expected 'ok', got %v", resp)
	}
}
