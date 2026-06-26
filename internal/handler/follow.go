package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"godan/internal/middleware"
	"godan/internal/pkg/errcode"
	"godan/internal/pkg/response"
	"godan/internal/service"
)

type FollowHandler struct {
	svc *service.FollowService
}

func NewFollowHandler(svc *service.FollowService) *FollowHandler {
	return &FollowHandler{svc: svc}
}

// Follow godoc
// @Summary 关注用户
// @Tags follow
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body FollowReq true "关注参数"
// @Success 200 {object} response.Response
// @Router /api/v1/user/follow [post]
func (h *FollowHandler) Follow(c *gin.Context) {
	var req struct {
		UserID uint64 `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}

	followerID := middleware.GetUserID(c)
	ec := h.svc.Follow(followerID, req.UserID)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, nil)
}

// Unfollow godoc
// @Summary 取消关注
// @Tags follow
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body FollowReq true "取消关注参数"
// @Success 200 {object} response.Response
// @Router /api/v1/user/unfollow [post]
func (h *FollowHandler) Unfollow(c *gin.Context) {
	var req struct {
		UserID uint64 `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}

	followerID := middleware.GetUserID(c)
	ec := h.svc.Unfollow(followerID, req.UserID)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, nil)
}

// GetFollowers godoc
// @Summary 获取粉丝列表
// @Tags follow
// @Produce json
// @Param id path int true "用户ID"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response
// @Router /api/v1/user/{id}/followers [get]
func (h *FollowHandler) GetFollowers(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, total, ec := h.svc.GetFollowerList(userID, page, pageSize)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, gin.H{
		"list":      users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetFollowees godoc
// @Summary 获取关注列表
// @Tags follow
// @Produce json
// @Param id path int true "用户ID"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response
// @Router /api/v1/user/{id}/followees [get]
func (h *FollowHandler) GetFollowees(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, total, ec := h.svc.GetFolloweeList(userID, page, pageSize)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, gin.H{
		"list":      users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetMutualFollows godoc
// @Summary 获取互相关注列表
// @Tags follow
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response
// @Router /api/v1/user/mutual-follows [get]
func (h *FollowHandler) GetMutualFollows(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, total, ec := h.svc.GetMutualFollows(userID, page, pageSize)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, gin.H{
		"list":      users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// Block godoc
// @Summary 拉黑用户
// @Tags follow
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body BlockReq true "拉黑参数"
// @Success 200 {object} response.Response
// @Router /api/v1/user/block [post]
func (h *FollowHandler) Block(c *gin.Context) {
	var req struct {
		UserID uint64 `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}

	userID := middleware.GetUserID(c)
	ec := h.svc.Block(userID, req.UserID)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, nil)
}

// Unblock godoc
// @Summary 取消拉黑
// @Tags follow
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body BlockReq true "取消拉黑参数"
// @Success 200 {object} response.Response
// @Router /api/v1/user/unblock [post]
func (h *FollowHandler) Unblock(c *gin.Context) {
	var req struct {
		UserID uint64 `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}

	userID := middleware.GetUserID(c)
	ec := h.svc.Unblock(userID, req.UserID)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, nil)
}
