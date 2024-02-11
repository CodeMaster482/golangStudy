package main

import (
	"context"
	"main/internal/api"
	"main/internal/app/config"
	"main/internal/app/dsn"
	"main/internal/app/handler"
	"main/internal/app/redis"
	"main/internal/app/repository"

	Minio "main/internal/app/minio"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	minioClient := Minio.NewMinioClient(logger)

	router := gin.Default()

	conf, err := config.NewConfig(logger)
	if err != nil {
		logger.Fatalf("Error conf reading: %v", err)
	}

	ctx := context.Background()
	redisClient, errRedis := redis.New(ctx, conf.Redis)
	if errRedis != nil {
		logger.Fatalf("Errof with redis connect: %s", err)
	}

	dbConnStr, err := dsn.FromEnv()
	if err != nil {
		logger.Fatalf("Error dsn reading: %v", err)
	}

	repo, err := repository.NewRepository(dbConnStr, logger)
	if err != nil {
		logger.Fatalf("Error repo creating: %v", err)
	}

	h := handler.NewHandler(conf, repo, minioClient, redisClient, logger)

	app := api.NewApiServer(conf, router, logger, h)

	app.RunApi()

	// api.StartServer()
}
