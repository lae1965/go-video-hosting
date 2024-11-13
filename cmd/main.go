package main

import (
	"go-video-hosting/gRPC/client"
	"go-video-hosting/internal/database"
	"go-video-hosting/internal/handler"
	"go-video-hosting/internal/server"
	"go-video-hosting/internal/service"
	"go-video-hosting/internal/validator"
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

	grpcServer, err := grpcclient.NewFilesGRPCServer()
	if err != nil {
		logrus.Fatalf("Failed to connect to gRPC: %s", err.Error())
	}
	defer grpcServer.Connection.Close()

	grpcClient := grpcclient.NewFilesGRPCClient(grpcServer)
	validate := validator.NewValidator()
	db := database.NewDatabase(dbSql)
	service := service.NewService(db, *grpcClient)
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
