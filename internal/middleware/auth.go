package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"godan/internal/config"
	pkgjwt "godan/internal/pkg/jwt"
	"godan/internal/pkg/response"
)

func Auth(cfg *config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Auth token
		authHeader := c.GetHeader("Authorization")
		// NotFound
		if authHeader == "" {
			response.Error(c, pkgjwt.ErrTokenFromResponse(nil))
			c.Abort()
			return
		}
		// bearer token format
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Error(c, pkgjwt.ErrTokenFromResponse(nil))
			c.Abort()
			return
		}
		// Get claim
		claims, err := pkgjwt.ParseToken(parts[1], cfg.AccessSecret)
		if err != nil {
			response.Error(c, pkgjwt.ErrTokenFromResponse(err))
			c.Abort()
			return
		}
		//
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

// OptionalAuth is repeated
func OptionalAuth(cfg *config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.Next()
			return
		}

		claims, err := pkgjwt.ParseToken(parts[1], cfg.AccessSecret)
		if err == nil {

			c.Set("user_id", claims.UserID)
		}

		c.Next()
	}
}

func GetUserID(c *gin.Context) uint64 {
	id, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	uid, ok := id.(uint64)
	if !ok {
		return 0
	}
	return uid
}
