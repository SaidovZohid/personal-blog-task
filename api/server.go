package api

import (
	v1 "github.com/SaidovZohid/personal-blog-task/api/v1"
	"github.com/SaidovZohid/personal-blog-task/config"
	"github.com/SaidovZohid/personal-blog-task/storage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

type RouterOptions struct {
	Cfg      *config.Config
	Storage  storage.StorageI
	InMemory storage.InMemoryStorageI
	Logger   *zap.Logger
}

// @title           Swagger for personal blog api
// @version         1.0
// @description     This is personal blog api
// @BasePath  /v1
func New(opt *RouterOptions) *gin.Engine {
	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "*")
	router.Use(cors.New(corsConfig))

	handlerV1 := v1.New(&v1.HandlerV1Options{
		Cfg:      opt.Cfg,
		Storage:  opt.Storage,
		InMemory: opt.InMemory,
		Logger:   opt.Logger,
	})

	_ = handlerV1

	apiV1 := router.Group("/v1")
	{
		apiV1.POST("/auth/signup", handlerV1.SignUp)
		apiV1.POST("/auth/verify", handlerV1.VerifyEmail)
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
