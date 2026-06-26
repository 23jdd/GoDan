package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"godan/internal/middleware"
	"godan/internal/pkg/errcode"
	"godan/internal/pkg/response"
	"godan/internal/service"
)

type CommentHandler struct {
	svc *service.CommentService
}

func NewCommentHandler(svc *service.CommentService) *CommentHandler {
	return &CommentHandler{svc: svc}
}

// CreateComment godoc
// @Summary 发表评论
// @Tags comment
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body CreateCommentReq true "评论参数"
// @Success 200 {object} response.Response
// @Router /api/v1/comment [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {
	var req struct {
		VideoID    uint64 `json:"video_id" binding:"required"`
		Content    string `json:"content" binding:"required"`
		ParentID   string `json:"parent_id"`
		RootID     string `json:"root_id"`
		ReplyToUID uint64 `json:"reply_to_uid"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}

	userID := middleware.GetUserID(c)
	comment, ec := h.svc.CreateComment(userID, req.VideoID, req.Content, req.ParentID, req.RootID, req.ReplyToUID)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, comment)
}

// DeleteComment godoc
// @Summary 删除评论
// @Tags comment
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "评论ID"
// @Success 200 {object} response.Response
// @Router /api/v1/comment/{id} [delete]
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	commentID := c.Param("id")
	userID := middleware.GetUserID(c)

	ec := h.svc.DeleteComment(userID, commentID)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, nil)
}

// GetRootComments godoc
// @Summary 获取一级评论列表
// @Tags comment
// @Produce json
// @Param video_id query int true "视频ID"
// @Param sort query string false "排序(new/hot)"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response
// @Router /api/v1/comments [get]
func (h *CommentHandler) GetRootComments(c *gin.Context) {
	videoID, _ := strconv.ParseUint(c.Query("video_id"), 10, 64)
	sort := c.DefaultQuery("sort", "new")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	comments, total, ec := h.svc.GetRootComments(videoID, sort, page, pageSize)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, gin.H{"list": comments, "total": total, "page": page, "page_size": pageSize})
}

// GetReplies godoc
// @Summary 获取楼中楼回复
// @Tags comment
// @Produce json
// @Param root_id query string true "根评论ID"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response
// @Router /api/v1/comments/replies [get]
func (h *CommentHandler) GetReplies(c *gin.Context) {
	rootID := c.Query("root_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	replies, total, ec := h.svc.GetReplies(rootID, page, pageSize)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, gin.H{"list": replies, "total": total, "page": page, "page_size": pageSize})
}

// LikeComment godoc
// @Summary 点赞评论
// @Tags comment
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "评论ID"
// @Success 200 {object} response.Response
// @Router /api/v1/comment/{id}/like [post]
func (h *CommentHandler) LikeComment(c *gin.Context) {
	commentID := c.Param("id")
	userID := middleware.GetUserID(c)

	ec := h.svc.LikeComment(userID, commentID)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, nil)
}

// UnlikeComment godoc
// @Summary 取消点赞评论
// @Tags comment
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "评论ID"
// @Success 200 {object} response.Response
// @Router /api/v1/comment/{id}/like [delete]
func (h *CommentHandler) UnlikeComment(c *gin.Context) {
	commentID := c.Param("id")
	userID := middleware.GetUserID(c)

	ec := h.svc.UnlikeComment(userID, commentID)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, nil)
}
