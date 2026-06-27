package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MaxBody  Is A Middleware for limit Body Length
func MaxBody(maxsize int64) gin.HandlerFunc {
	return func(context *gin.Context) {
		length := context.Request.ContentLength
		if length > maxsize || length <= 0 {
			context.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
	}
}
