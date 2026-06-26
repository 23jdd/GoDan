package router

import (
	"github.com/gin-gonic/gin"

	"godan/internal/config"
	"godan/internal/handler"
	"godan/internal/middleware"
	"godan/internal/pkg/response"
	"godan/internal/pkg/storage"
	"godan/internal/service"
)

func Setup(cfg *config.Config) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)

	r := gin.New()
	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	store, err := storage.New(&cfg.Storage)
	if err != nil {
		panic("failed to init storage: " + err.Error())
	}

	userSvc := service.NewUserService(cfg)
	followSvc := service.NewFollowService()

	userH := handler.NewUserHandler(userSvc)
	followH := handler.NewFollowHandler(followSvc)
	uploadH := handler.NewUploadHandler(store)

	// 本地存储模式：静态文件服务
	if cfg.Storage.Type == "local" {
		r.Static(cfg.Storage.Local.URLPrefix, cfg.Storage.Local.Path)
	}

	r.GET("/ping", func(c *gin.Context) {
		response.Success(c, gin.H{"ping": "pong"})
	})

	api := r.Group("/api/v1")
	{
		// 公开接口
		api.POST("/user/register", userH.Register)
		api.POST("/user/login", userH.Login)
		api.POST("/user/refresh", userH.RefreshToken)
		api.POST("/user/code/send", userH.SendVerificationCode)

		// 可选认证
		optional := api.Group("")
		optional.Use(middleware.OptionalAuth(&cfg.JWT))
		{
			optional.GET("/user/profile/:id", userH.GetUserProfile)
			optional.GET("/user/:id/followers", followH.GetFollowers)
			optional.GET("/user/:id/followees", followH.GetFollowees)
		}

		// 需要认证
		auth := api.Group("")
		auth.Use(middleware.Auth(&cfg.JWT))
		{
			auth.POST("/upload/avatar", uploadH.UploadAvatar)

			auth.GET("/user/profile", userH.GetProfile)
			auth.PUT("/user/profile", userH.UpdateProfile)
			auth.PUT("/user/password", userH.ChangePassword)

			auth.POST("/user/bind/email", userH.BindEmail)
			auth.POST("/user/bind/phone", userH.BindPhone)

			auth.POST("/user/follow", followH.Follow)
			auth.POST("/user/unfollow", followH.Unfollow)
			auth.GET("/user/mutual-follows", followH.GetMutualFollows)

			auth.POST("/user/block", followH.Block)
			auth.POST("/user/unblock", followH.Unblock)
		}
	}

	return r
}
