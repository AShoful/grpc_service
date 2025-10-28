package main

import (
	grpcserver "grpc/server"
	"grpc/server/models"
	"grpc/server/pkg/handler"
	"grpc/server/pkg/repository"
	"grpc/server/pkg/service"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	db := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	db.AutoMigrate(&models.Book{})
	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	handler := handler.NewHandler(service)

	grpcserver.RunServer(handler, service)
}

func initConfig() error {
	viper.AddConfigPath("server/configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
