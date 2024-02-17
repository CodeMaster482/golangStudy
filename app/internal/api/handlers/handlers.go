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
	{Id: 0, Name: "10", Description: "Lorem ipsum dolor sit amet, consectetur adipisicing elit", Img: "/static/ten.jpg"},
	{Id: 1, Name: "50", Description: "Lorem ipsum dolor sit amet, consectetur adipisicing elit", Img: "/static/fifty.jpg"},
	{Id: 2, Name: "100", Description: "Lorem ipsum dolor sit amet, consectetur adipisicing elit", Img: "/static/handred.jpg"},
	{Id: 3, Name: "500", Description: "Lorem ipsum dolor sit amet, consectetur adipisicing elit", Img: "/static/fivehandred.jpg"},
	{Id: 4, Name: "1000", Description: "Lorem ipsum dolor sit amet, consectetur adipisicing elit", Img: "/static/thouthend.jpg"},
}

func PingPongHandler(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"User-Agent": context.Request.UserAgent(),
		"URL":        context.Request.URL,
		"Accepted":   context.Accepted,
		"message":    "pong",
	})
}

func ServicesHandler(context *gin.Context) {
	query := context.DefaultQuery("banknote", "")
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
		"Query":     query,
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
