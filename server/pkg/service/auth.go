package service

import (
	"errors"
	"fmt"
	"grpc/server/models"
	"grpc/server/pkg/repository"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

const (
	tokenTTL = 12 * time.Hour
)

type AuthService struct {
	repo      repository.Authorization
	jwtSecret []byte
}

type tokenClaims struct {
	jwt.RegisteredClaims
	UserID uint `json:"user_id"`
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found (ok in prod)")
	}
}

func NewAuthService(repo repository.Authorization) *AuthService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is not set in environment variables")
	}

	return &AuthService{
		repo:      repo,
		jwtSecret: []byte(secret),
	}
}

func (s *AuthService) CreateUser(user models.User) (uint, error) {
	hashedPassword, _ := generatePasswordHash(user.Password)
	user.Password = hashedPassword
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {

	user, err := s.repo.GetUser(username)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: user.ID,
	})

	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) ParseToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("cannot parse claims")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("user_id not found in token")
	}

	return uint(userIDFloat), nil
}

func generatePasswordHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
