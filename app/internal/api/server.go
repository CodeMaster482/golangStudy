package api

import (
	"go-web-bmstu/internal/api/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	log.Println("Server start up")

	log.Println("Started on: http://localhost:8080")

	router := gin.Default()

	router.LoadHTMLGlob("../../templates/*")

	router.Static("/static", "../../resources")

	router.GET("/ping", handlers.PingPongHandler)

	router.GET("/service", handlers.ServicesHandler)

	router.GET("/service/:id", handlers.ProductPageHandler)

	router.GET("/home", handlers.HomePageHandler)

	router.Run()

	log.Println("Server down")
}
