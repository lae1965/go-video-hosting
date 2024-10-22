package main

import (
	"go-video-hosting/internal/server"
	"go-video-hosting/pkg/database"
	"go-video-hosting/pkg/handler"
	"go-video-hosting/pkg/service"
	"go-video-hosting/pkg/validator"
	"os"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("Error initialization config: %s", err.Error())
	}

	envVars := []string{"DB_PASSWORD", "ACCESS_KEY", "REFRESH_KEY", "MAIL_PASSWORD"}
	for _, envVar := range envVars {
		if os.Getenv(envVar) == "" {
			logrus.Fatalf("Requied environment variable %s is not set", envVar)
		}
	}

	dbSql, err := database.Connection()
	if err != nil {
		logrus.Fatalf("Failed to installation db: %s", err.Error())
	}

	validate := validator.NewValidator()
	db := database.NewDatabase(dbSql)
	service := service.NewService(db)
	handlers := handler.NewHandler(service, validate)
	srv := new(server.Server)

	port := viper.GetString("port")
	if port == "" {
		port = "8080"
	}
	if err := srv.Run(port, handlers.InitRoutes()); err != nil {
		logrus.Fatalf("Error occured while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
