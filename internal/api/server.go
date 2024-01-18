package api

import (
	"fmt"

	"main/internal/api/handlers"
	"main/internal/config"
	"main/internal/handler"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ApiServer struct {
	Config  *config.Config
	Logger  *logrus.Logger
	Router  *gin.Engine
	Handler *handler.Handler
}

func NewApiServer(cfg *config.Config, r *gin.Engine, log *logrus.Logger, h *handler.Handler) *ApiServer {
	return &ApiServer{
		Config:  cfg,
		Logger:  log,
		Router:  r,
		Handler: h,
	}
}

func (api *ApiServer) RunApi() {
	api.Logger.Info("Server start up")
	api.Handler.RegisterHandler(api.Router)

	serverAddress := fmt.Sprintf("%s:%d", api.Config.ServiceHost, api.Config.ServicePort)
	if err := api.Router.Run(serverAddress); err != nil {
		api.Logger.Fatalln(err)
	}
	api.Logger.Info("Server down")
}

func StartServer() {
	log.Println("Server start up")

	log.Println("Started on: http://localhost:8080")

	router := gin.Default()

	router.LoadHTMLGlob("templates/*")

	router.Static("/image", "./resources")

	router.GET("/banknotes", handlers.ServicesHandler)

	router.GET("/banknotes/:id", handlers.ProductPageHandler)

	router.GET("/home", handlers.HomePageHandler)

	router.Run()

	log.Println("Server down")
}
