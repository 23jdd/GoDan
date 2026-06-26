package middleware

import (
	"github.com/gin-gonic/gin"

	"godan/internal/dao"
	"godan/internal/pkg/errcode"
	"godan/internal/pkg/response"
)

func Admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		if userID == 0 {
			response.Error(c, errcode.ErrUnauthorized)
			c.Abort()
			return
		}

		user, _ := dao.GetUserByID(userID)
		if user == nil || user.Role < 1 {
			response.Error(c, errcode.ErrForbidden)
			c.Abort()
			return
		}

		c.Set("user_role", user.Role)
		c.Next()
	}
}
