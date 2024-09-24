package main

import (
	"go-video-hosting/internal/server"
	"go-video-hosting/pkg/handler"
	"go-video-hosting/pkg/repository"
	"go-video-hosting/pkg/service"
	"log"

	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("Error initialization config: %s", err.Error())
	}

	repo := repository.NewRepository()
	service := service.NewService(repo)
	handlers := handler.NewHandler(service)
	srv := new(server.Server)

	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		log.Fatalf("Error occured while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
