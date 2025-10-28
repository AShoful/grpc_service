package repository

import (
	"errors"
	"fmt"
	"grpc/server/models"

	"gorm.io/gorm"
)

type BookPostgres struct {
	db *gorm.DB
}

func NewBookPostgres(db *gorm.DB) *BookPostgres {
	return &BookPostgres{db: db}
}

func (r *BookPostgres) Create(book models.Book) (uint, error) {
	if err := r.db.Create(&book).Error; err != nil {
		return 0, fmt.Errorf("failed to create book: %w", err)
	}
	return book.ID, nil
}

func (r *BookPostgres) GetAll() ([]models.Book, error) {
	var books []models.Book
	if err := r.db.Find(&books).Error; err != nil {
		return nil, fmt.Errorf("failed to get all books: %w", err)
	}
	return books, nil
}

func (r *BookPostgres) GetById(bookId uint) (models.Book, error) {
	var book models.Book
	err := r.db.First(&book, bookId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Book{}, fmt.Errorf("book with id %d not found", bookId)
		}
		return models.Book{}, fmt.Errorf("failed to find book with id %d: %w", bookId, err)
	}
	return book, nil
}

func (r *BookPostgres) Delete(userId, bookId uint) error {
	var book models.Book

	if err := r.db.First(&book, bookId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("book not found")
		}
		return err
	}

	if userId != book.UserId {
		return fmt.Errorf("user does not have permission to delete this book")
	}

	if err := r.db.Delete(&book).Error; err != nil {
		return fmt.Errorf("failed to delete book: %w", err)
	}

	return nil
}

func (r *BookPostgres) Update(userId, bookId uint, input models.UpdateBook) error {
	var book models.Book
	if err := r.db.First(&book, bookId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("book not found")
		}
		return err
	}

	if userId != book.UserId {
		return fmt.Errorf("user does not have permission to delete this book")
	}

	if input.Title != nil {
		book.Title = *input.Title
	}

	if input.Author != nil {
		book.Author = *input.Author
	}

	if err := r.db.Save(&book).Error; err != nil {
		return fmt.Errorf("failed to save book: %w", err)
	}

	return nil
}
