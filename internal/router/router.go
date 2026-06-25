package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"godan/internal/middleware"
	"godan/internal/pkg/response"
)

func Setup(mode string) *gin.Engine {
	gin.SetMode(mode)

	r := gin.New()

	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	r.GET("/ping", func(c *gin.Context) {
		response.Success(c, gin.H{"ping": "pong"})
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    10005,
			"message": "route not found",
		})
	})

	return r
}
