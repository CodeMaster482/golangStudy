package main

import (
	"main/internal/api"
	"main/internal/config"
	"main/internal/dsn"
	"main/internal/handler"
	"main/internal/minio"
	"main/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	router := gin.Default()

	miniocli := minio.NewMinioClient(logger)

	conf, err := config.NewConfig(logger)
	if err != nil {
		logger.Fatalf("Error conf reading: %v", err)
	}

	dbConnStr := dsn.FromEnv()

	repo, err := repository.NewRepository(dbConnStr, logger)
	if err != nil {
		logger.Fatalf("Error repo creating: %v", err)
	}

	h := handler.NewHandler(repo, miniocli, logger)

	app := api.NewApiServer(conf, router, logger, h)

	app.RunApi()

	// api.StartServer()
}
