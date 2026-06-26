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
	uploader, err := storage.NewMultipartUploader(&cfg.Storage)
	if err != nil {
		panic("failed to init multipart uploader: " + err.Error())
	}

	userSvc := service.NewUserService(cfg)
	followSvc := service.NewFollowService()
	videoSvc := service.NewVideoService(store, uploader, cfg)
	interactionSvc := service.NewInteractionService(cfg)

	userH := handler.NewUserHandler(userSvc)
	followH := handler.NewFollowHandler(followSvc)
	uploadH := handler.NewUploadHandler(store)
	videoH := handler.NewVideoHandler(videoSvc)
	interactionH := handler.NewInteractionHandler(interactionSvc)
	danmakuH := handler.NewDanmakuHandler(interactionSvc)
	danmakuHub := handler.NewDanmakuHub(interactionSvc)

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
		api.POST("/user/register/code", userH.RegisterWithCode)
		api.POST("/user/login", userH.Login)
		api.POST("/user/login/code", userH.LoginByCode)
		api.POST("/user/refresh", userH.RefreshToken)
		api.POST("/user/code/send", userH.SendVerificationCode)

		// 弹幕 WebSocket
		api.GET("/danmakus/ws", func(c *gin.Context) {
			danmakuHub.HandleWebSocket(c)
		})

		// 可选认证
		optional := api.Group("")
		optional.Use(middleware.OptionalAuth(&cfg.JWT))
		{
			optional.GET("/user/profile/:id", userH.GetUserProfile)
			optional.GET("/user/:id/followers", followH.GetFollowers)
			optional.GET("/user/:id/followees", followH.GetFollowees)
			optional.GET("/user/:id/videos", videoH.GetUserVideos)
			optional.GET("/video/:id", videoH.GetVideoDetail)
			optional.GET("/video/:id/related", videoH.GetRelatedVideos)
			optional.GET("/videos", videoH.GetCategoryVideos)
			optional.GET("/videos/hot", videoH.GetHotVideos)
			optional.GET("/videos/search", videoH.SearchVideos)
			optional.GET("/danmakus", danmakuH.GetDanmakus)
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

			auth.POST("/video/upload/init", videoH.InitUpload)
			auth.POST("/video/upload/chunk", videoH.UploadChunk)
			auth.POST("/video/upload/complete", videoH.CompleteUpload)
			auth.GET("/video/upload/status", videoH.UploadStatus)
			auth.POST("/video/upload/abort", videoH.AbortUpload)

			auth.POST("/video/cover/upload", uploadH.UploadCover)
			auth.PUT("/video/:id/cover", videoH.UpdateCover)
			auth.POST("/video/:id/publish", videoH.PublishVideo)
			auth.DELETE("/video/:id", videoH.DeleteVideo)
			auth.GET("/user/videos", videoH.GetMyVideos)

			// 互动
			auth.POST("/video/:id/like", interactionH.LikeVideo)
			auth.DELETE("/video/:id/like", interactionH.CancelLike)
			auth.GET("/video/:id/like/status", interactionH.LikeStatus)
			auth.POST("/video/:id/coin", interactionH.GiveCoin)
			auth.GET("/user/coins/remaining", interactionH.RemainingCoins)
			auth.POST("/video/:id/share", interactionH.ShareVideo)

			// 收藏夹
			auth.POST("/favorite/folder", interactionH.CreateFolder)
			auth.PUT("/favorite/folder/:id", interactionH.UpdateFolder)
			auth.DELETE("/favorite/folder/:id", interactionH.DeleteFolder)
			auth.GET("/user/folders", interactionH.GetFolders)
			auth.POST("/favorite/add", interactionH.AddToFolder)
			auth.POST("/favorite/remove", interactionH.RemoveFromFolder)
			auth.GET("/favorite/folder/:id/items", interactionH.GetFolderItems)
		}
	}

	return r
}
