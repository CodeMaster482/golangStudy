package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Services struct {
	Id          int
	Name        string
	Description string
	Img         string
}

func StartServer() {
	log.Println("Server start up")

	runner := gin.Default()

	services := []Services{
		{Id: 0, Name: "credit", Description: "Lorem ipsum dolor sit amet, consectetur adipisicing elit", Img: "/image/black.png"},
		{Id: 1, Name: "deposite", Description: "Lorem ipsum dolor sit amet, consectetur adipisicing elit", Img: "/image/white.png"},
		{Id: 2, Name: "scam", Description: "Lorem ipsum dolor sit amet, consectetur adipisicing elit", Img: "/image/credit.png"},
	}

	runner.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	runner.LoadHTMLGlob("../../templates/*")

	runner.GET("/home", func(context *gin.Context) {
		// context.HTML(http.StatusOK, "productCard.html", gin.H{
		//	"services": services,
		// })

		context.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":    "RIP 💀 project",
			"services": services,
		})
	})

	runner.Static("/image", "../../resources")

	runner.Run()

	log.Println("Server down")
}
