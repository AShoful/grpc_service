package service

import (
	"grpc/server/models"
	"grpc/server/pkg/repository"

	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeAuthRepo struct{}

func (f fakeAuthRepo) CreateUser(user models.User) (uint, error) {
	return 0, nil
}

func (f fakeAuthRepo) GetUser(username string) (models.User, error) {
	return models.User{}, nil
}

type fakeBookRepo struct{}

func (f fakeBookRepo) Create(book models.Book) (uint, error) {
	return 0, nil
}

func (f fakeBookRepo) GetAll() ([]models.Book, error) {
	return nil, nil
}

func (f fakeBookRepo) GetById(bookId uint) (models.Book, error) {
	return models.Book{}, nil
}

func (f fakeBookRepo) Delete(userId, bookId uint) error {
	return nil
}

func (f fakeBookRepo) Update(userId, bookId uint, book models.UpdateBook) error {
	return nil
}

func TestNewService(t *testing.T) {
	repos := &repository.Repository{
		Authorization: fakeAuthRepo{},
		Book:          fakeBookRepo{},
	}

	svc := NewService(repos)

	assert.NotNil(t, svc)

	_, ok := svc.Authorization.(*AuthService)
	assert.True(t, ok, "Authorization must be *AuthService")

	_, ok = svc.Book.(*BookService)
	assert.True(t, ok, "Book must be *BookService")

	authService := svc.Authorization.(*AuthService)
	bookService := svc.Book.(*BookService)

	assert.NotNil(t, authService.repo)
	assert.NotNil(t, bookService.repo)
}
