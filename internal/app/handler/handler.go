package handler

import (
	"main/internal/app/config"
	"main/internal/app/ds"
	"main/internal/app/pkg/hash"
	"main/internal/app/redis"
	"main/internal/app/repository"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	router.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	// услуги
	api.GET("/banknotes", h.WithoutJWTError(ds.Buyer, ds.Moderator, ds.Admin), h.BanknotesList) // ?
	api.GET("/banknotes/:id", h.GetBanknotesById)                                               // ?
	api.POST("/banknotes", h.WithAuthCheck(ds.Moderator, ds.Admin), h.AddBanknote)
	api.PUT("/banknotes", h.WithAuthCheck(ds.Moderator, ds.Admin), h.BanknoteUpdate)
	api.PUT("/banknotes/upload-image", h.WithAuthCheck(ds.Moderator, ds.Admin), h.AddImage)
	api.DELETE("/banknotes", h.WithAuthCheck(ds.Moderator, ds.Admin), h.DeleteBanknote)
	api.POST("/banknotes/request", h.WithAuthCheck(ds.Buyer, ds.Moderator, ds.Admin), h.AddBanknoteToRequest)
	api.Use(cors.Default()).DELETE("/banknotes/delete/:id", h.DeleteBanknote)

	// заявки
	api.GET("/operations", h.WithAuthCheck(ds.Buyer, ds.Moderator, ds.Admin), h.OperationList)
	api.GET("/operations/:id", h.WithAuthCheck(ds.Buyer, ds.Moderator, ds.Admin), h.GetOperationById)
	// api.POST("/operations/", h.CreateDraft)
	api.PUT("/operations", h.WithAuthCheck(ds.Buyer, ds.Moderator, ds.Admin), h.UpdateOperation)

	// статусы
	api.PUT("/operations/form", h.WithAuthCheck(ds.Buyer, ds.Moderator, ds.Admin), h.FormOperationRequest)
	api.PUT("/operations/updateStatus", h.WithAuthCheck(ds.Moderator, ds.Admin), h.UpdateStatusOperationRequest)
	//api.PUT("/operations/finish/:id", h.WithAuthCheck([]ds.Role{ds.Admin}), h.FinishTenderRequest)

	api.DELETE("/operations", h.WithAuthCheck(ds.Buyer, ds.Moderator, ds.Admin), h.DeleteOperation)

	// m-m
	api.DELETE("/operation-request-company", h.WithoutJWTError(ds.Buyer, ds.Moderator, ds.Admin), h.DeleteBanknoteFromRequest)
	api.PUT("/operation-request-company", h.WithoutJWTError(ds.Buyer, ds.Moderator, ds.Admin), h.UpdateOperationBanknote)
	registerStatic(router)

	// auth && reg
	api.POST("/user/signIn", h.Login)
	api.POST("/user/signUp", h.Register)
	api.POST("/user/logout", h.Logout)

	// асинхронный сервис
	api.PUT("/tenders/user-form-start", h.WithoutJWTError(ds.Buyer, ds.Moderator, ds.Admin), h.UserRequest) // обращение к асинхронному сервису
	api.PUT("/tenders/user-form-finish", h.FinishUserRequest)
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
