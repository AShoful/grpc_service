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

func ctxWithUserID(id uint) context.Context {
	return context.WithValue(context.Background(), handler.UserIDKey(), id)
}

func TestBookHandler_CreateBook_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBook := mock_service.NewMockBook(ctrl)
	h := handler.NewBookHandler(mockBook)

	req := &proto.Book{
		Title:  "Go in Action",
		Author: "John",
	}

	expectedID := uint(10)
	userId := uint(1)

	mockBook.
		EXPECT().
		Create(models.Book{
			Title:  "Go in Action",
			Author: "John",
			UserId: userId,
		}).
		Return(expectedID, nil)

	ctx := ctxWithUserID(userId)

	resp, err := h.CreateBook(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Id != uint32(expectedID) {
		t.Fatalf("expected id %d, got %d", expectedID, resp.Id)
	}
}

func TestBookHandler_CreateBook_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBook := mock_service.NewMockBook(ctrl)
	h := handler.NewBookHandler(mockBook)

	req := &proto.Book{
		Title:  "",
		Author: "Author",
	}

	_, err := h.CreateBook(context.Background(), req)
	if err == nil {
		t.Fatal("expected validation error")
	}

	st, _ := status.FromError(err)
	if st.Code() != codes.InvalidArgument {
		t.Fatalf("expected InvalidArgument, got %v", st.Code())
	}
}

func TestBookHandler_GetBook_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBook := mock_service.NewMockBook(ctrl)
	h := handler.NewBookHandler(mockBook)

	expected := models.Book{
		ID:     5,
		Title:  "Go",
		Author: "Rob",
		UserId: 1,
	}

	mockBook.
		EXPECT().
		GetById(uint(5)).
		Return(expected, nil)

	resp, err := h.GetBook(context.Background(), &proto.BookId{Id: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Title != expected.Title || resp.Author != expected.Author {
		t.Fatalf("data mismatch")
	}
}

func TestBookHandler_GetBook_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBook := mock_service.NewMockBook(ctrl)
	h := handler.NewBookHandler(mockBook)

	mockBook.
		EXPECT().
		GetById(uint(99)).
		Return(models.Book{}, errors.New("not found"))

	_, err := h.GetBook(context.Background(), &proto.BookId{Id: 99})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestBookHandler_GetBooks_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBook := mock_service.NewMockBook(ctrl)
	h := handler.NewBookHandler(mockBook)

	books := []models.Book{
		{ID: 1, Title: "A", Author: "B", UserId: 1},
		{ID: 2, Title: "C", Author: "D", UserId: 2},
	}

	mockBook.
		EXPECT().
		GetAll().
		Return(books, nil)

	resp, err := h.GetBooks(context.Background(), &proto.Empty{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Books) != 2 {
		t.Fatalf("expected 2 books, got %d", len(resp.Books))
	}
}

func TestBookHandler_UpdateBook_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBook := mock_service.NewMockBook(ctrl)
	h := handler.NewBookHandler(mockBook)

	userId := uint(1)

	req := &proto.Book{
		Id:     10,
		Title:  "Updated",
		Author: "New Author",
	}

	mockBook.
		EXPECT().
		Update(userId, uint(10), models.UpdateBook{
			Title:  &req.Title,
			Author: &req.Author,
		}).
		Return(nil)

	ctx := ctxWithUserID(userId)

	resp, err := h.UpdateBook(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Title != req.Title {
		t.Fatalf("update not applied")
	}
}

func TestBookHandler_UpdateBook_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBook := mock_service.NewMockBook(ctrl)
	h := handler.NewBookHandler(mockBook)

	req := &proto.Book{
		Id:     10,
		Title:  "",
		Author: "",
	}

	_, err := h.UpdateBook(context.Background(), req)
	if err == nil {
		t.Fatal("expected validation error")
	}

	st, _ := status.FromError(err)
	if st.Code() != codes.InvalidArgument {
		t.Fatalf("expected InvalidArgument, got %v", st.Code())
	}
}

func TestBookHandler_DeleteBook_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBook := mock_service.NewMockBook(ctrl)
	h := handler.NewBookHandler(mockBook)

	userId := uint(1)

	mockBook.
		EXPECT().
		Delete(userId, uint(10)).
		Return(nil)

	ctx := ctxWithUserID(userId)

	_, err := h.DeleteBook(ctx, &proto.BookId{Id: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBookHandler_DeleteBook_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBook := mock_service.NewMockBook(ctrl)
	h := handler.NewBookHandler(mockBook)

	userId := uint(1)

	mockBook.
		EXPECT().
		Delete(userId, uint(5)).
		Return(errors.New("delete error"))

	ctx := ctxWithUserID(userId)

	_, err := h.DeleteBook(ctx, &proto.BookId{Id: 5})
	if err == nil {
		t.Fatal("expected error")
	}
}
