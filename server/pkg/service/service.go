package service

import (
	"grpc/server/models"
	"grpc/server/pkg/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Service struct {
	Authorization
	Book
}

type Authorization interface {
	CreateUser(user models.User) (uint, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (uint, error)
}

type Book interface {
	Create(book models.Book) (uint, error)
	GetAll() ([]models.Book, error)
	GetById(bookId uint) (models.Book, error)
	Delete(userId, bookId uint) error
	Update(userId, bookId uint, book models.UpdateBook) error
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Book:          NewBookService(repos.Book),
	}
}
