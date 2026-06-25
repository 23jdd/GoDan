package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"godan/internal/config"
	"godan/internal/pkg/database"
	"godan/internal/pkg/logger"
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

	logger.Log.Info("config loaded", zap.Any("config", cfg))

	if err := database.InitMySQL(&cfg.MySQL); err != nil {
		logger.Log.Fatal("failed to init mysql", zap.Error(err))
	}
	defer database.Close()
	logger.Log.Info("mysql connected")

	r := router.Setup(cfg.Server.Mode)

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
