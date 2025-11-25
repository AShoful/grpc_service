package service

import (
	"grpc/server/models"
	mock_repository "grpc/server/pkg/repository/mocks"

	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestBookService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBook := mock_repository.NewMockBook(ctrl)
	service := NewBookService(mockBook)

	book := models.Book{Title: "Test Book"}

	mockBook.EXPECT().Create(book).Return(uint(1), nil)

	id, err := service.Create(book)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), id)
}

func TestBookService_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBook := mock_repository.NewMockBook(ctrl)
	service := NewBookService(mockBook)

	books := []models.Book{
		{ID: 1, Title: "Book1"},
		{ID: 2, Title: "Book2"},
	}

	mockBook.EXPECT().GetAll().Return(books, nil)

	result, err := service.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, books, result)
}

func TestBookService_GetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBook := mock_repository.NewMockBook(ctrl)
	service := NewBookService(mockBook)

	book := models.Book{ID: 1, Title: "Book1"}

	mockBook.EXPECT().GetById(uint(1)).Return(book, nil)

	result, err := service.GetById(1)
	assert.NoError(t, err)
	assert.Equal(t, book, result)
}

func TestBookService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBook := mock_repository.NewMockBook(ctrl)
	service := NewBookService(mockBook)

	mockBook.EXPECT().Delete(uint(1), uint(2)).Return(nil)

	err := service.Delete(1, 2)
	assert.NoError(t, err)
}

func TestBookService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBook := mock_repository.NewMockBook(ctrl)
	service := NewBookService(mockBook)

	title := "Updated"
	update := models.UpdateBook{Title: &title}

	mockBook.EXPECT().Update(uint(1), uint(2), update).Return(nil)

	err := service.Update(1, 2, update)
	assert.NoError(t, err)
}
