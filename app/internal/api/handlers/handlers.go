package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Services struct {
	Id          int
	Name        string
	Description string
	Img         string
}

var services = []Services{
	{Id: 0, Name: "credit", Description: "Lorem ipsum dolor sit amet, consectetur adipisicing elit", Img: "/static/credit.jpg"},
	{Id: 1, Name: "deposite", Description: "Lorem ipsum dolor sit amet, consectetur adipisicing elit", Img: "/static/deposite.jpg"},
	{Id: 2, Name: "transfer", Description: "Lorem ipsum dolor sit amet, consectetur adipisicing elit", Img: "/static/transfer.jpg"},
	{Id: 3, Name: "open account", Description: "Lorem ipsum dolor sit amet, consectetur adipisicing elit", Img: "/static/accaunt.jpg"},
	{Id: 4, Name: "exchange", Description: "Lorem ipsum dolor sit amet, consectetur adipisicing elit", Img: "/static/change.png"},
}

func PingPongHandler(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"User-Agent": context.Request.UserAgent(),
		"URL":        context.Request.URL,
		"Accepted":   context.Accepted,
		"message":    "pong",
	})
}

func HomePageHandler(context *gin.Context) {
	context.HTML(http.StatusOK, "index.html", gin.H{
		"title":     "RIP 💀 project",
		"services":  services,
		"styleFile": "static/styles.css",
	})
}

func ServicesHandler(context *gin.Context) {
	query := context.DefaultQuery("query", "")
	var result []Services

	if query != "" {
		for _, data := range services {
			if strings.HasPrefix(data.Name, query) {
				result = append(result, data)
			}
		}
	}

	if len(result) == 0 {
		result = append(result, services...)
	}

	context.HTML(http.StatusOK, "index.html", gin.H{
		"title":     "RIP 💀 project",
		"services":  result,
		"styleFile": "static/styles.css",
	})
}

func ProductPageHandler(context *gin.Context) {
	idGet := context.Param("id")
	id, _ := strconv.Atoi(idGet)

	if len(services) < id {
		context.String(http.StatusNotFound, "404 ---- Not Found")
		return
	}

	context.HTML(http.StatusOK, "productPage.html", gin.H{
		"services":  services[id],
		"styleFile": "../static/styles.css",
	})
}
