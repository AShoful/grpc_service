package service

import (
	"errors"
	"grpc/server/models"
	"os"

	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockAuthRepo struct {
	mock.Mock
}

func (m *MockAuthRepo) CreateUser(user models.User) (uint, error) {
	args := m.Called(user)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockAuthRepo) GetUser(username string) (models.User, error) {
	args := m.Called(username)
	return args.Get(0).(models.User), args.Error(1)
}

func TestMain(m *testing.M) {
	os.Setenv("JWT_SECRET", "TEST_SECRET_KEY")
	code := m.Run()
	os.Exit(code)
}

func TestAuthService_CreateUser(t *testing.T) {
	mockRepo := new(MockAuthRepo)
	service := NewAuthService(mockRepo)

	user := models.User{
		Username: "test",
		Password: "12345",
	}

	mockRepo.On("CreateUser", mock.AnythingOfType("models.User")).
		Return(uint(1), nil)

	id, err := service.CreateUser(user)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), id)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_GenerateToken_Success(t *testing.T) {
	mockRepo := new(MockAuthRepo)
	service := NewAuthService(mockRepo)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	user := models.User{
		ID:       10,
		Username: "user",
		Password: string(hashed),
	}

	mockRepo.On("GetUser", "user").Return(user, nil)

	token, err := service.GenerateToken("user", "password123")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_GenerateToken_InvalidPassword(t *testing.T) {
	mockRepo := new(MockAuthRepo)
	service := NewAuthService(mockRepo)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)

	mockRepo.On("GetUser", "user").Return(models.User{
		ID:       1,
		Username: "user",
		Password: string(hashed),
	}, nil)

	token, err := service.GenerateToken("user", "wrong")

	assert.Error(t, err)
	assert.Equal(t, "invalid password", err.Error())
	assert.Empty(t, token)
}

func TestAuthService_GenerateToken_UserNotFound(t *testing.T) {
	mockRepo := new(MockAuthRepo)
	service := NewAuthService(mockRepo)

	mockRepo.On("GetUser", "ghost").
		Return(models.User{}, errors.New("user not found"))

	token, err := service.GenerateToken("ghost", "123")

	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestAuthService_ParseToken_Success(t *testing.T) {
	mockRepo := new(MockAuthRepo)
	service := NewAuthService(mockRepo)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": float64(42),
		"exp":     float64(time.Now().Add(time.Hour).Unix()),
		"iat":     float64(time.Now().Unix()),
	})

	tokenStr, _ := token.SignedString([]byte("TEST_SECRET_KEY"))

	id, err := service.ParseToken(tokenStr)

	assert.NoError(t, err)
	assert.Equal(t, uint(42), id)
}

func TestAuthService_ParseToken_InvalidSignature(t *testing.T) {
	mockRepo := new(MockAuthRepo)
	service := NewAuthService(mockRepo)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": float64(1),
	})

	tokenStr, _ := token.SignedString([]byte("WRONG_KEY"))

	_, err := service.ParseToken(tokenStr)

	assert.Error(t, err)
	assert.Equal(t, "invalid token", err.Error())
}

func TestAuthService_ParseToken_NoUserID(t *testing.T) {
	mockRepo := new(MockAuthRepo)
	service := NewAuthService(mockRepo)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": float64(time.Now().Add(time.Hour).Unix()),
	})

	tokenStr, _ := token.SignedString([]byte("TEST_SECRET_KEY"))

	_, err := service.ParseToken(tokenStr)

	assert.Error(t, err)
	assert.Equal(t, "user_id not found in token", err.Error())
}
