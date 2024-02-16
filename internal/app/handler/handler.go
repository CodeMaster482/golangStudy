package handler

import (
	"main/internal/app/config"
	"main/internal/app/ds"
	"main/internal/app/pkg/hash"
	"main/internal/app/redis"
	"main/internal/app/repository"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "main/docs"
)

type Handler struct {
	Config     *config.Config
	Repository *repository.Repository
	MinioCli   *minio.Client
	Redis      *redis.Client
	Hasher     hash.PasswordHasher
	Logger     *logrus.Logger
	//TokenManager auth.TokenManager
}

func NewHandler(cfg *config.Config,
	r *repository.Repository,
	mcli *minio.Client,
	rcli *redis.Client,
	log *logrus.Logger,
) *Handler {
	return &Handler{
		Config:     cfg,
		Repository: r,
		MinioCli:   mcli,
		Redis:      rcli,
		Hasher:     hash.NewSHA256Hasher(os.Getenv("SALT")),
		Logger:     log,
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")

	// услуги
	api.GET("/banknote", h.WithoutJWTError(ds.Buyer, ds.Moderator, ds.Admin), h.BanknotesList) // ?
	api.GET("/banknotes/:id", h.GetBanknotesById)                                              // ?
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
	api.DELETE("/operation-request-banknote", h.WithoutJWTError(ds.Buyer, ds.Moderator, ds.Admin), h.DeleteBanknoteFromRequest)
	api.PUT("/operation-request-banknote", h.WithoutJWTError(ds.Buyer, ds.Moderator, ds.Admin), h.UpdateOperationBanknote)
	registerStatic(router)

	// auth && reg
	api.POST("/user/signIn", h.Login)
	api.POST("/user/signUp", h.Register)
	api.POST("/user/logout", h.Logout)

	// асинхронный сервис
	api.PUT("/operations/user-form-start", h.WithoutJWTError(ds.Buyer, ds.Moderator, ds.Admin), h.UserRequest) // обращение к асинхронному сервису
	api.PUT("/operations/user-form-finish", h.FinishUserRequest)
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
