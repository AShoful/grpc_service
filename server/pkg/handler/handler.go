package handler

import (
	"grpc/server/pkg/service"
)

type Handler struct {
	AuthHandler *AuthHandler
	BookHandler *BookHandler
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		AuthHandler: NewAuthHandler(services.Authorization),
		BookHandler: NewBookHandler(services.Book),
	}
}
