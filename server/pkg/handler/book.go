package handler

import (
	"context"
	"fmt"
	"grpc/proto"
	"grpc/server/models"
	"grpc/server/pkg/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BookHandler struct {
	proto.UnimplementedBookServiceServer
	bookService service.Book
}

func NewBookHandler(bookService service.Book) *BookHandler {
	return &BookHandler{bookService: bookService}
}

func (h *BookHandler) CreateBook(ctx context.Context, req *proto.Book) (*proto.BookId, error) {
	book := models.Book{
		Title:  req.Title,
		Author: req.Author,
	}

	if err := validate.Struct(book); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userId, _ := UserIDFromContext(ctx)

	book.UserId = userId

	id, err := h.bookService.Create(book)
	if err != nil {
		return nil, err
	}
	return &proto.BookId{Id: uint32(id)}, nil
}

func (h *BookHandler) GetBook(ctx context.Context, req *proto.BookId) (*proto.Book, error) {
	book, err := h.bookService.GetById(uint(req.Id))
	if err != nil {
		return nil, err
	}
	return &proto.Book{
		Id:     uint32(book.ID),
		Title:  book.Title,
		Author: book.Author,
		Userid: uint32(book.UserId),
	}, nil
}

func (h *BookHandler) GetBooks(ctx context.Context, req *proto.Empty) (*proto.BookList, error) {
	books, err := h.bookService.GetAll()
	if err != nil {
		return nil, err
	}
	var pbBooks []*proto.Book
	for _, b := range books {
		pbBooks = append(pbBooks, &proto.Book{
			Id:     uint32(b.ID),
			Title:  b.Title,
			Author: b.Author,
			Userid: uint32(b.UserId),
		})
	}
	return &proto.BookList{Books: pbBooks}, nil
}

func (h *BookHandler) UpdateBook(ctx context.Context, req *proto.Book) (*proto.Book, error) {
	updateBook := models.UpdateBook{
		Title:  &req.Title,
		Author: &req.Author,
	}

	if err := validate.Struct(updateBook); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userId, _ := UserIDFromContext(ctx)

	err := h.bookService.Update(uint(userId), uint(req.Id), updateBook)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (h *BookHandler) DeleteBook(ctx context.Context, req *proto.BookId) (*proto.Empty, error) {
	userId, _ := UserIDFromContext(ctx)
	if err := h.bookService.Delete(uint(userId), uint(req.Id)); err != nil {
		return nil, err
	}
	return &proto.Empty{}, nil
}

func UserIDFromContext(ctx context.Context) (uint, error) {
	id, ok := ctx.Value(userIDKey).(uint)
	if !ok {
		return 0, fmt.Errorf("user_id not found in context")
	}

	return id, nil
}
