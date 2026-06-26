package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"godan/internal/middleware"
	"godan/internal/pkg/errcode"
	"godan/internal/pkg/response"
	"godan/internal/service"
)

type ActivityHandler struct {
	svc *service.ActivityService
}

func NewActivityHandler(svc *service.ActivityService) *ActivityHandler {
	return &ActivityHandler{svc: svc}
}

// GetTimeline godoc
// @Summary 关注动态流
// @Tags activity
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码"
// @Param page_size query int false "每页"
// @Success 200 {object} response.Response
// @Router /api/v1/timeline [get]
func (h *ActivityHandler) GetTimeline(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	list, total, ec := h.svc.GetTimeline(userID, page, pageSize)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, gin.H{"list": list, "total": total, "page": page, "page_size": pageSize})
}

// GetUserActivities godoc
// @Summary 个人动态列表
// @Tags activity
// @Produce json
// @Param id path int true "用户ID"
// @Param page query int false "页码"
// @Param page_size query int false "每页"
// @Success 200 {object} response.Response
// @Router /api/v1/user/{id}/activities [get]
func (h *ActivityHandler) GetUserActivities(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	list, total, ec := h.svc.GetUserActivities(userID, page, pageSize)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, gin.H{"list": list, "total": total, "page": page, "page_size": pageSize})
}
