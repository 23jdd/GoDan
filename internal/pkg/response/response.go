package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"godan/internal/pkg/errcode"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    errcode.Success.Code,
		Message: errcode.Success.Message,
		Data:    data,
	})
}

func Error(c *gin.Context, e *errcode.ErrorCode) {
	c.JSON(e.HTTPStatus(), Response{
		Code:    e.Code,
		Message: e.Message,
	})
}

func ErrorWithMsg(c *gin.Context, e *errcode.ErrorCode, msg string) {
	c.JSON(e.HTTPStatus(), Response{
		Code:    e.Code,
		Message: msg,
	})
}
