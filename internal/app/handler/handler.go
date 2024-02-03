package handler

import (
	"main/internal/app/config"
	"main/internal/app/pkg/hash"
	"main/internal/app/redis"
	"main/internal/app/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
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
	Config     *config.Config
	Redis      *redis.Client
	Hasher     hash.PasswordHasher
	//TokenManager auth.TokenManager
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
	api.GET("/banknotes", h.BanknotesList)         // услуги
	api.GET("/banknotes/:id", h.BanknoteById)      // конкретная
	api.PUT("/banknotes/:id", h.BanknoteUpdate)    // изменить
	api.DELETE("/banknotes/:id", h.DeleteBanknote) // удалить

	api.POST("/banknotes/request/:id", h.AddBanknoteToRequest) // добавить учлуги к заявке

	// application
	api.GET("/operations", h.OperationList) // все заявки
	api.GET("/operations/:id", h.GetOperationById)
	api.PUT("/operations", h.UpdateOperation)
	api.PUT("/operation/form/:id", h.FormOperationRequest)
	api.PUT("/operations/reject/:id", h.RejectOperationRequest)
	api.PUT("/operation/finish/:id", h.FinishOperationRequest)

	api.DELETE("/banknoteOperation", h.DeleteBanknoteFromRequest)
	api.PUT("/banknoteOperation", h.UpdateOperationBanknote)

	// router.GET("/banknotes", h.BanknotesList)
	// router.GET("/banknotes/:id", h.BanknoteById)
	// router.POST("api/deleteBanknote", h.DeleteBanknote)

	registerStatic(router)
}

func registerStatic(router *gin.Engine) {
	//router.LoadHTMLGlob("./static/templates/*")
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
