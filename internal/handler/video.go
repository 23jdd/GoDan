package handler

import (
	"io"
	"strconv"

	"github.com/gin-gonic/gin"

	"godan/internal/middleware"
	"godan/internal/pkg/errcode"
	"godan/internal/pkg/response"
	"godan/internal/service"
)

type VideoHandler struct {
	svc *service.VideoService
}

func NewVideoHandler(svc *service.VideoService) *VideoHandler {
	return &VideoHandler{svc: svc}
}

// InitUpload godoc
// @Summary 初始化分片上传
// @Tags video
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body InitUploadReq true "初始化参数"
// @Success 200 {object} response.Response
// @Router /api/v1/video/upload/init [post]
func (h *VideoHandler) InitUpload(c *gin.Context) {
	var req struct {
		Filename string `json:"filename" binding:"required"`
		FileSize int64  `json:"file_size" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}

	userID := middleware.GetUserID(c)
	uploadID, chunkCount, chunkSize, ec := h.svc.InitUpload(userID, req.Filename, req.FileSize)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, gin.H{
		"upload_id":   uploadID,
		"chunk_count": chunkCount,
		"chunk_size":  chunkSize,
	})
}

// UploadChunk godoc
// @Summary 上传分片
// @Tags video
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param upload_id formData string true "上传ID"
// @Param chunk_index formData int true "分片序号(0开始)"
// @Param file formData file true "分片数据"
// @Success 200 {object} response.Response
// @Router /api/v1/video/upload/chunk [post]
func (h *VideoHandler) UploadChunk(c *gin.Context) {
	uploadID := c.PostForm("upload_id")
	chunkIndex, err := strconv.Atoi(c.PostForm("chunk_index"))
	if err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}
	defer file.Close()

	chunkData, err := io.ReadAll(file)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}

	ec := h.svc.UploadChunk(uploadID, chunkIndex, chunkData)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, gin.H{
		"chunk_index": chunkIndex,
		"uploaded":    true,
	})
}

// CompleteUpload godoc
// @Summary 完成上传
// @Tags video
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body CompleteUploadReq true "完成参数"
// @Success 200 {object} response.Response
// @Router /api/v1/video/upload/complete [post]
func (h *VideoHandler) CompleteUpload(c *gin.Context) {
	var req struct {
		UploadID    string `json:"upload_id" binding:"required"`
		Title       string `json:"title" binding:"required,min=1,max=200"`
		Description string `json:"description"`
		CategoryID  int    `json:"category_id"`
		Tags        string `json:"tags"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}

	userID := middleware.GetUserID(c)
	videoID, videoURL, ec := h.svc.CompleteUpload(userID, req.UploadID, req.Title, req.Description, req.CategoryID, req.Tags)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, gin.H{
		"video_id":  videoID,
		"video_url": videoURL,
	})
}

// UploadStatus godoc
// @Summary 查询上传进度（断点续传）
// @Tags video
// @Produce json
// @Security ApiKeyAuth
// @Param upload_id query string true "上传ID"
// @Success 200 {object} response.Response
// @Router /api/v1/video/upload/status [get]
func (h *VideoHandler) UploadStatus(c *gin.Context) {
	uploadID := c.Query("upload_id")
	if uploadID == "" {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}

	uploaded, total, ec := h.svc.UploadStatus(uploadID)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, gin.H{
		"uploaded_chunks": uploaded,
		"total_chunks":    total,
	})
}

// AbortUpload godoc
// @Summary 取消上传
// @Tags video
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body AbortUploadReq true "取消参数"
// @Success 200 {object} response.Response
// @Router /api/v1/video/upload/abort [post]
func (h *VideoHandler) AbortUpload(c *gin.Context) {
	var req struct {
		UploadID string `json:"upload_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}

	ec := h.svc.AbortUpload(req.UploadID)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, nil)
}

// GetVideoDetail godoc
// @Summary 获取视频详情
// @Tags video
// @Produce json
// @Param id path int true "视频ID"
// @Success 200 {object} response.Response
// @Router /api/v1/video/{id} [get]
func (h *VideoHandler) GetVideoDetail(c *gin.Context) {
	videoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}

	v, ec := h.svc.GetVideoDetail(videoID)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, v)
}

// UpdateCover godoc
// @Summary 更新视频封面
// @Tags video
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "视频ID"
// @Param body body UpdateCoverReq true "封面URL"
// @Success 200 {object} response.Response
// @Router /api/v1/video/{id}/cover [put]
func (h *VideoHandler) UpdateCover(c *gin.Context) {
	videoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}
	var req struct {
		CoverURL string `json:"cover_url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}

	userID := middleware.GetUserID(c)
	ec := h.svc.UpdateVideoCover(userID, videoID, req.CoverURL)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, nil)
}

// PublishVideo godoc
// @Summary 发布视频（待审→已发布）
// @Tags video
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "视频ID"
// @Success 200 {object} response.Response
// @Router /api/v1/video/{id}/publish [post]
func (h *VideoHandler) PublishVideo(c *gin.Context) {
	videoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}

	userID := middleware.GetUserID(c)
	ec := h.svc.PublishVideo(userID, videoID)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, nil)
}

// DeleteVideo godoc
// @Summary 删除视频（下架）
// @Tags video
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "视频ID"
// @Success 200 {object} response.Response
// @Router /api/v1/video/{id} [delete]
func (h *VideoHandler) DeleteVideo(c *gin.Context) {
	videoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}

	userID := middleware.GetUserID(c)
	ec := h.svc.DeleteVideo(userID, videoID)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, nil)
}

// GetUserVideos godoc
// @Summary 获取用户投稿列表
// @Tags video
// @Produce json
// @Param id path int true "用户ID"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response
// @Router /api/v1/user/{id}/videos [get]
func (h *VideoHandler) GetUserVideos(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	videos, total, ec := h.svc.GetUserVideos(userID, page, pageSize)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, gin.H{"list": videos, "total": total, "page": page, "page_size": pageSize})
}

// GetCategoryVideos godoc
// @Summary 获取分区视频列表
// @Tags video
// @Produce json
// @Param category_id query int true "分区ID"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response
// @Router /api/v1/videos [get]
func (h *VideoHandler) GetCategoryVideos(c *gin.Context) {
	categoryID, _ := strconv.Atoi(c.DefaultQuery("category_id", "0"))
	sort := c.DefaultQuery("sort", "new") // new / hot
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	videos, total, ec := h.svc.GetCategoryVideos(categoryID, sort, page, pageSize)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, gin.H{"list": videos, "total": total, "page": page, "page_size": pageSize})
}

// GetMyVideos godoc
// @Summary 获取我的投稿列表
// @Tags video
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response
// @Router /api/v1/user/videos [get]
func (h *VideoHandler) GetMyVideos(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	videos, total, ec := h.svc.GetUserVideos(userID, page, pageSize)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, gin.H{"list": videos, "total": total, "page": page, "page_size": pageSize})
}

// GetHotVideos godoc
// @Summary 首页热门推荐
// @Tags video
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response
// @Router /api/v1/videos/hot [get]
func (h *VideoHandler) GetHotVideos(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	videos, total, ec := h.svc.GetHotVideos(page, pageSize)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, gin.H{"list": videos, "total": total, "page": page, "page_size": pageSize})
}

// SearchVideos godoc
// @Summary 视频搜索
// @Tags video
// @Produce json
// @Param q query string true "搜索关键词"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response
// @Router /api/v1/videos/search [get]
func (h *VideoHandler) SearchVideos(c *gin.Context) {
	q := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	videos, total, ec := h.svc.SearchVideos(q, page, pageSize)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, gin.H{"list": videos, "total": total, "page": page, "page_size": pageSize})
}

// GetRelatedVideos godoc
// @Summary 相关视频推荐
// @Tags video
// @Produce json
// @Param id path int true "视频ID"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response
// @Router /api/v1/video/{id}/related [get]
func (h *VideoHandler) GetRelatedVideos(c *gin.Context) {
	videoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	videos, ec := h.svc.GetRelatedVideos(videoID, page, pageSize)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, gin.H{"list": videos, "page": page, "page_size": pageSize})
}
