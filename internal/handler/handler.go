package handler

import (
	"main/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go"
	"github.com/sirupsen/logrus"
)

const (
	creatorID   = 1
	moderatorID = 1
)

type Handler struct {
	Logger     *logrus.Logger
	Repository *repository.Repository
	MinioCli   *minio.Client
}

func NewHandler(r *repository.Repository, cli *minio.Client, l *logrus.Logger) *Handler {
	return &Handler{
		Repository: r,
		MinioCli:   cli,
		Logger:     l,
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	api := router.Group("/api")

	// service
	api.GET("/banknotes", h.BanknotesList)
	api.GET("/banknotes/:id", h.BanknoteById)
	// api.POST("/banknotes/:id", h.BanknoteUpdate)
	api.DELETE("/banknotes/:id", h.DeleteBanknote)

	api.POST("/banknotes/request/:id", h.AddBanknoteToRequest)

	//api.POST("/banknotes/request/:id", h.AddBanknoteToRequest)

	// application
	//
	//api.GET("/operations", h.OperationList)
	//api.GET("/operations/:id", h.OperationById)
	//api.PUT("/oprations", h.OprationUpdate)
	//api.PUT("/opration/form/:id", h.FormOprationApplication)
	//api.PUT("/oprations/reject/:id", h.RejectOperationApplication)
	//api.PUT("/opration/finish/:id", h.FinishOperationApplication)
	//
	//api.DELETE("/banknoteOperation", h.DeleteBanknoteFromApplication)
	//api.PUT("/banknoteOperation", h.UpdateOperationBnaknote)

	// router.GET("/banknotes", h.BanknotesList)
	// router.GET("/banknotes/:id", h.BanknoteById)
	// router.POST("api/deleteBanknote", h.DeleteBanknote)

	registerStatic(router)
}

func registerStatic(router *gin.Engine) {
	router.LoadHTMLGlob("./static/templates/*")
	router.Static("/static", "./static")
	router.Static("/css", "./static")
	router.Static("/img", "./static")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	h.Logger.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}

func (h *Handler) successHandler(ctx *gin.Context, key string, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		key:      data,
	})
}

func (h *Handler) successAddHandler(ctx *gin.Context, key string, data interface{}) {
	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		key:      data,
	})
}
