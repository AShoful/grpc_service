package handler_test

import (
	"context"
	"errors"
	"testing"

	"grpc/proto"
	"grpc/server/models"
	"grpc/server/pkg/handler"
	mock_service "grpc/server/pkg/service/mocks"

	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAuthHandler_SignUp_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mock_service.NewMockAuthorization(ctrl)
	h := handler.NewAuthHandler(mockAuth)

	req := &proto.User{
		Name:     "John",
		Username: "john123",
		Password: "pass123",
	}

	expectedID := uint(10)

	mockAuth.
		EXPECT().
		CreateUser(models.User{
			Name:     "John",
			Username: "john123",
			Password: "pass123",
		}).
		Return(expectedID, nil)

	resp, err := h.SignUp(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Id != uint32(expectedID) {
		t.Fatalf("expected id %d, got %d", expectedID, resp.Id)
	}
}

func TestAuthHandler_SignUp_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mock_service.NewMockAuthorization(ctrl)
	h := handler.NewAuthHandler(mockAuth)

	req := &proto.User{
		Name:     "",
		Username: "user",
		Password: "pass",
	}

	_, err := h.SignUp(context.Background(), req)
	if err == nil {
		t.Fatal("expected validation error")
	}

	st, _ := status.FromError(err)
	if st.Code() != codes.InvalidArgument {
		t.Fatalf("expected InvalidArgument, got: %v", st.Code())
	}
}

func TestAuthHandler_SignUp_CreateUserError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mock_service.NewMockAuthorization(ctrl)
	h := handler.NewAuthHandler(mockAuth)

	req := &proto.User{
		Name:     "John",
		Username: "john123",
		Password: "pass123",
	}

	mockAuth.
		EXPECT().
		CreateUser(gomock.Any()).
		Return(uint(0), errors.New("db error"))

	_, err := h.SignUp(context.Background(), req)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAuthHandler_SignIn_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mock_service.NewMockAuthorization(ctrl)
	h := handler.NewAuthHandler(mockAuth)

	req := &proto.SignInRequest{
		Username: "john123",
		Password: "pass123",
	}

	mockAuth.
		EXPECT().
		GenerateToken("john123", "pass123").
		Return("token_abc", nil)

	resp, err := h.SignIn(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Token != "token_abc" {
		t.Fatalf("expected token token_abc, got %s", resp.Token)
	}
}

func TestAuthHandler_SignIn_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mock_service.NewMockAuthorization(ctrl)
	h := handler.NewAuthHandler(mockAuth)

	req := &proto.SignInRequest{
		Username: "",
		Password: "pass",
	}

	_, err := h.SignIn(context.Background(), req)
	if err == nil {
		t.Fatal("expected validation error")
	}

	st, _ := status.FromError(err)
	if st.Code() != codes.InvalidArgument {
		t.Fatalf("expected InvalidArgument, got %v", st.Code())
	}
}

func TestAuthHandler_SignIn_GenerateTokenError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mock_service.NewMockAuthorization(ctrl)
	h := handler.NewAuthHandler(mockAuth)

	req := &proto.SignInRequest{
		Username: "john",
		Password: "pass123",
	}

	mockAuth.
		EXPECT().
		GenerateToken("john", "pass123").
		Return("", errors.New("invalid credentials"))

	_, err := h.SignIn(context.Background(), req)
	if err == nil {
		t.Fatal("expected error")
	}
}
