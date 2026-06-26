// @title           GoDan API
// @version         1.0
// @description     仿B站视频分享平台 API 文档
// @contact.name    GoDan
// @host            localhost:8080
// @BasePath        /api/v1

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"godan/internal/config"
	"godan/internal/dao"
	"godan/internal/pkg/database"
	"godan/internal/pkg/logger"
	"godan/internal/pkg/mongodb"
	"godan/internal/pkg/redis"
	"godan/internal/router"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	if err := logger.Init(&cfg.Log); err != nil {
		panic(fmt.Sprintf("failed to init logger: %v", err))
	}
	defer logger.Log.Sync()

	if err := database.InitMySQL(&cfg.MySQL); err != nil {
		logger.Log.Fatal("failed to init mysql", zap.Error(err))
	}
	defer database.Close()
	logger.Log.Info("mysql connected")

	if err := dao.AutoMigrate(); err != nil {
		logger.Log.Fatal("failed to run migration", zap.Error(err))
	}
	logger.Log.Info("migration completed")

	if err := redis.Init(&cfg.Redis); err != nil {
		logger.Log.Warn("redis init failed (non-fatal)", zap.Error(err))
	} else {
		defer redis.Close()
		logger.Log.Info("redis connected")
	}

	if err := mongodb.Init(&cfg.MongoDB); err != nil {
		logger.Log.Warn("mongodb init failed (non-fatal)", zap.Error(err))
	} else {
		defer mongodb.Close()
		logger.Log.Info("mongodb connected")
	}

	r := router.Setup(cfg)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Log.Info("server starting", zap.String("addr", addr))

	go func() {
		if err := r.Run(addr); err != nil {
			logger.Log.Fatal("failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("shutting down server...")
}
