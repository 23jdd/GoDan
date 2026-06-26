package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"godan/internal/pkg/errcode"
	"godan/internal/pkg/response"
	"godan/internal/pkg/redis"
)

// RateLimit: fixed-window counter using Redis INCR + TTL.
func RateLimit(limit int, window time.Duration, mode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var key string
		switch mode {
		case "user":
			key = fmt.Sprintf("rl:user:%d:%d", GetUserID(c), time.Now().Unix()/int64(window.Seconds()))
		case "ip":
			key = fmt.Sprintf("rl:ip:%s:%d", c.ClientIP(), time.Now().Unix()/int64(window.Seconds()))
		default:
			key = fmt.Sprintf("rl:ip:%s:%d", c.ClientIP(), time.Now().Unix()/int64(window.Seconds()))
		}

		ctx := context.Background()
		count, err := redis.RDB.Incr(ctx, key).Result()
		if err == nil {
			redis.RDB.Expire(ctx, key, window)
		}

		if count > int64(limit) {
			response.Error(c, &errcode.ErrorCode{Code: 10429, Message: "too many requests"})
			c.Abort()
			return
		}

		c.Next()
	}
}
