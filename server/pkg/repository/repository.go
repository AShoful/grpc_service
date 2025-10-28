package repository

import (
	"grpc/server/models"

	"gorm.io/gorm"
)

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

type BookRepo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Book:          NewBookPostgres(db),
	}
}

func NewBookRepo(db *gorm.DB) *BookRepo {
	return &BookRepo{db: db}
}
