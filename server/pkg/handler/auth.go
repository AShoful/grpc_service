package handler

import (
	"context"
	"fmt"
	"grpc/proto"
	"grpc/server/models"
	"grpc/server/pkg/service"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	proto.UnimplementedUserServiceServer
	userService service.Authorization
}

func NewAuthHandler(userService service.Authorization) *AuthHandler {
	return &AuthHandler{userService: userService}
}

var validate = validator.New()

func (h *AuthHandler) SignUp(ctx context.Context, req *proto.User) (*proto.UserId, error) {
	user := models.User{
		Name:     req.Name,
		Username: req.Username,
		Password: req.Password,
	}

	if err := validate.Struct(user); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	id, err := h.userService.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return &proto.UserId{Id: uint32(id)}, nil
}

func (h *AuthHandler) SignIn(ctx context.Context, req *proto.SignInRequest) (*proto.AuthResponse, error) {

	signInInput := models.SignInInput{
		Username: req.Username,
		Password: req.Password,
	}

	fmt.Println(signInInput)

	if err := validate.Struct(signInInput); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := h.userService.GenerateToken(req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	return &proto.AuthResponse{Token: token}, nil
}
