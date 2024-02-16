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

	router.Use(corsMiddleware())

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


func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
		}

		c.Next()
	}
}
