package repository

import (
	"grpc/server/models"

	"gorm.io/gorm"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type Authorization interface {
	CreateUser(user models.User) (uint, error)
	GetUser(username string) (models.User, error)
}

type Book interface {
	Create(book models.Book) (uint, error)
	GetAll() ([]models.Book, error)
	GetById(bookId uint) (models.Book, error)
	Delete(userId, bookId uint) error
	Update(userId, bookId uint, book models.UpdateBook) error
}

type Repository struct {
	Authorization
	Book
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Book:          NewBookPostgres(db),
	}
}
