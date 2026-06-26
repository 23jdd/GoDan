package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"godan/internal/middleware"
	"godan/internal/pkg/errcode"
	"godan/internal/pkg/response"
	"godan/internal/service"
)

type InteractionHandler struct {
	svc *service.InteractionService
}

func NewInteractionHandler(svc *service.InteractionService) *InteractionHandler {
	return &InteractionHandler{svc: svc}
}

// LikeVideo godoc
// @Summary 点赞视频
// @Tags interaction
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "视频ID"
// @Success 200 {object} response.Response
// @Router /api/v1/video/{id}/like [post]
func (h *InteractionHandler) LikeVideo(c *gin.Context) {
	videoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}
	userID := middleware.GetUserID(c)
	ec := h.svc.LikeVideo(userID, videoID)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, nil)
}

// CancelLike godoc
// @Summary 取消点赞/点踩
// @Tags interaction
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "视频ID"
// @Success 200 {object} response.Response
// @Router /api/v1/video/{id}/like [delete]
func (h *InteractionHandler) CancelLike(c *gin.Context) {
	videoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}
	userID := middleware.GetUserID(c)
	ec := h.svc.CancelLikeDislike(userID, videoID)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, nil)
}

// LikeStatus godoc
// @Summary 查询点赞状态
// @Tags interaction
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "视频ID"
// @Success 200 {object} response.Response
// @Router /api/v1/video/{id}/like/status [get]
func (h *InteractionHandler) LikeStatus(c *gin.Context) {
	videoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}
	userID := middleware.GetUserID(c)
	t, _ := h.svc.GetUserLikeStatus(userID, videoID)
	response.Success(c, gin.H{"type": t})
}

// GiveCoin godoc
// @Summary 投币
// @Tags interaction
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "视频ID"
// @Param body body CoinReq true "投币数量"
// @Success 200 {object} response.Response
// @Router /api/v1/video/{id}/coin [post]
func (h *InteractionHandler) GiveCoin(c *gin.Context) {
	videoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}
	var req struct {
		Count int `json:"count" binding:"required,min=1,max=2"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}
	userID := middleware.GetUserID(c)
	ec := h.svc.GiveCoin(userID, videoID, req.Count)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, nil)
}

// RemainingCoins godoc
// @Summary 今日剩余投币数
// @Tags interaction
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Router /api/v1/user/coins/remaining [get]
func (h *InteractionHandler) RemainingCoins(c *gin.Context) {
	userID := middleware.GetUserID(c)
	remaining, ec := h.svc.GetUserDailyRemainingCoins(userID)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, gin.H{"remaining": remaining})
}

// ShareVideo godoc
// @Summary 分享视频
// @Tags interaction
// @Produce json
// @Param id path int true "视频ID"
// @Success 200 {object} response.Response
// @Router /api/v1/video/{id}/share [post]
func (h *InteractionHandler) ShareVideo(c *gin.Context) {
	videoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}
	link, ec := h.svc.ShareVideo(videoID)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, gin.H{"link": link})
}

// --- Favorite ---

// CreateFolder godoc
// @Summary 创建收藏夹
// @Tags favorite
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body CreateFolderReq true "收藏夹参数"
// @Success 200 {object} response.Response
// @Router /api/v1/favorite/folder [post]
func (h *InteractionHandler) CreateFolder(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required,min=1,max=50"`
		Description string `json:"description"`
		IsPublic    int8   `json:"is_public"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}
	userID := middleware.GetUserID(c)
	id, ec := h.svc.CreateFolder(userID, req.Name, req.Description, req.IsPublic)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, gin.H{"folder_id": id})
}

// UpdateFolder godoc
// @Summary 编辑收藏夹
// @Tags favorite
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "收藏夹ID"
// @Param body body CreateFolderReq true "修改参数"
// @Success 200 {object} response.Response
// @Router /api/v1/favorite/folder/{id} [put]
func (h *InteractionHandler) UpdateFolder(c *gin.Context) {
	folderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}
	var req struct {
		Name        string `json:"name" binding:"required,min=1,max=50"`
		Description string `json:"description"`
		IsPublic    int8   `json:"is_public"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}
	userID := middleware.GetUserID(c)
	ec := h.svc.UpdateFolder(userID, folderID, req.Name, req.Description, req.IsPublic)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, nil)
}

// DeleteFolder godoc
// @Summary 删除收藏夹
// @Tags favorite
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "收藏夹ID"
// @Success 200 {object} response.Response
// @Router /api/v1/favorite/folder/{id} [delete]
func (h *InteractionHandler) DeleteFolder(c *gin.Context) {
	folderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}
	userID := middleware.GetUserID(c)
	ec := h.svc.DeleteFolder(userID, folderID)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, nil)
}

// GetFolders godoc
// @Summary 获取收藏夹列表
// @Tags favorite
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Router /api/v1/user/folders [get]
func (h *InteractionHandler) GetFolders(c *gin.Context) {
	userID := middleware.GetUserID(c)
	folders, ec := h.svc.GetUserFolders(userID)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, gin.H{"list": folders})
}

// AddToFolder godoc
// @Summary 收藏视频到收藏夹
// @Tags favorite
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body AddToFolderReq true "收藏参数"
// @Success 200 {object} response.Response
// @Router /api/v1/favorite/add [post]
func (h *InteractionHandler) AddToFolder(c *gin.Context) {
	var req struct {
		FolderID uint64 `json:"folder_id" binding:"required"`
		VideoID  uint64 `json:"video_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}
	userID := middleware.GetUserID(c)
	ec := h.svc.AddToFolder(userID, req.FolderID, req.VideoID)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, nil)
}

// RemoveFromFolder godoc
// @Summary 从收藏夹移除视频
// @Tags favorite
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body AddToFolderReq true "移除参数"
// @Success 200 {object} response.Response
// @Router /api/v1/favorite/remove [post]
func (h *InteractionHandler) RemoveFromFolder(c *gin.Context) {
	var req struct {
		FolderID uint64 `json:"folder_id" binding:"required"`
		VideoID  uint64 `json:"video_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}
	userID := middleware.GetUserID(c)
	ec := h.svc.RemoveFromFolder(userID, req.FolderID, req.VideoID)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, nil)
}

// GetFolderItems godoc
// @Summary 获取收藏夹内容
// @Tags favorite
// @Produce json
// @Param id path int true "收藏夹ID"
// @Param page query int false "页码"
// @Param page_size query int false "每页"
// @Success 200 {object} response.Response
// @Router /api/v1/favorite/folder/{id}/items [get]
func (h *InteractionHandler) GetFolderItems(c *gin.Context) {
	folderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	videos, total, ec := h.svc.GetFolderItems(int(folderID), page, pageSize)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, gin.H{"list": videos, "total": total, "page": page, "page_size": pageSize})
}
