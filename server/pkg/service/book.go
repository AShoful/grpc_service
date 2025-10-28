package service

import (
	"grpc/server/models"
	"grpc/server/pkg/repository"
)

type BookService struct {
	repo repository.Book
}

func NewBookService(repo repository.Book) *BookService {
	return &BookService{repo: repo}
}

func (s *BookService) Create(book models.Book) (uint, error) {
	return s.repo.Create(book)
}

func (s *BookService) GetAll() ([]models.Book, error) {
	return s.repo.GetAll()
}

func (s *BookService) GetById(bookId uint) (models.Book, error) {
	return s.repo.GetById(bookId)
}

func (s *BookService) Delete(userid, bookId uint) error {
	return s.repo.Delete(userid, bookId)
}

func (s *BookService) Update(userId, bookId uint, book models.UpdateBook) error {
	return s.repo.Update(userId, bookId, book)
}
