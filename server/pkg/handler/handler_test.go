package handler_test

import (
	"testing"

	"grpc/server/pkg/handler"
	"grpc/server/pkg/service"
	mock_service "grpc/server/pkg/service/mocks"

	"github.com/golang/mock/gomock"
)

func TestNewHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMock := mock_service.NewMockAuthorization(ctrl)
	bookMock := mock_service.NewMockBook(ctrl)

	svc := &service.Service{
		Authorization: authMock,
		Book:          bookMock,
	}

	h := handler.NewHandler(svc)

	if h == nil {
		t.Fatal("expected handler to be created, got nil")
	}
	if h.AuthHandler == nil {
		t.Error("expected AuthHandler to be initialized, got nil")
	}
	if h.BookHandler == nil {
		t.Error("expected BookHandler to be initialized, got nil")
	}
}

func TestUserIDKey(t *testing.T) {
	k1 := handler.UserIDKey()
	k2 := handler.UserIDKey()

	if k1 == nil {
		t.Fatal("expected UserIDKey to return non-nil value")
	}
	if k1 != k2 {
		t.Error("expected UserIDKey to return a stable consistent key")
	}
}
