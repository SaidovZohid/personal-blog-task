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
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Insert jwt access token
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

	apiV1 := router.Group("/v1")
	{
		// auth
		apiV1.POST("/auth/signup", handlerV1.SignUp)
		apiV1.POST("/auth/verify", handlerV1.VerifyEmail)
		apiV1.POST("/auth/login", handlerV1.Login)
		{
			// post
			apiV1.POST("/posts", handlerV1.AuthMiddleWare, handlerV1.CreatePost)
			apiV1.GET("/posts/:id", handlerV1.GetPost)
			apiV1.PUT("/posts/:id", handlerV1.AuthMiddleWare, handlerV1.UpdatePost)
			apiV1.DELETE("/posts/:id", handlerV1.AuthMiddleWare, handlerV1.DeletePost)
			apiV1.GET("/posts", handlerV1.GetAllPosts)

			// comment
			apiV1.POST("/comments", handlerV1.AuthMiddleWare, handlerV1.CreateComment)
			apiV1.PUT("/comments/:id", handlerV1.AuthMiddleWare, handlerV1.UpdateComment)
			apiV1.DELETE("/comments/:id", handlerV1.AuthMiddleWare, handlerV1.DeleteComment)
			apiV1.GET("/comments", handlerV1.GetAllComments)

			// reply
			apiV1.POST("/replies", handlerV1.AuthMiddleWare, handlerV1.CreateReply)
			apiV1.PUT("/replies/:id", handlerV1.AuthMiddleWare, handlerV1.UpdateReply)
			apiV1.DELETE("/replies/:id", handlerV1.AuthMiddleWare, handlerV1.DeleteReply)
			apiV1.GET("/replies", handlerV1.GetAllReplies)
		}
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
