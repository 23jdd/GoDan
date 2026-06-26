package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"godan/internal/middleware"
	"godan/internal/pkg/errcode"
	"godan/internal/pkg/response"
	"godan/internal/service"
)

type AdminHandler struct {
	svc *service.AdminService
}

func NewAdminHandler(svc *service.AdminService) *AdminHandler {
	return &AdminHandler{svc: svc}
}

// --- User ---

func (h *AdminHandler) GetUserList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	users, total, ec := h.svc.GetUserList(page, pageSize)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, gin.H{"list": users, "total": total, "page": page, "page_size": pageSize})
}

func (h *AdminHandler) BanUser(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	adminID := middleware.GetUserID(c)
	ec := h.svc.BanUser(adminID, userID)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, nil)
}

func (h *AdminHandler) UnbanUser(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	ec := h.svc.UnbanUser(userID)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, nil)
}

func (h *AdminHandler) SetRole(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Role int8 `json:"role" binding:"required,min=0,max=2"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}
	adminID := middleware.GetUserID(c)
	ec := h.svc.SetRole(adminID, userID, req.Role)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, nil)
}

// --- Review ---

func (h *AdminHandler) GetPendingVideos(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	videos, total, ec := h.svc.GetPendingVideos(page, pageSize)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, gin.H{"list": videos, "total": total, "page": page, "page_size": pageSize})
}

func (h *AdminHandler) ApproveVideo(c *gin.Context) {
	videoID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	ec := h.svc.ApproveVideo(videoID)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, nil)
}

func (h *AdminHandler) RejectVideo(c *gin.Context) {
	videoID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	ec := h.svc.RejectVideo(videoID)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, nil)
}

// --- Category ---

func (h *AdminHandler) CreateCategory(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
		Sort int    `json:"sort"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}
	id, ec := h.svc.CreateCategory(req.Name, req.Sort)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, gin.H{"id": id})
}

func (h *AdminHandler) UpdateCategory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Name string `json:"name" binding:"required"`
		Sort int    `json:"sort"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}
	ec := h.svc.UpdateCategory(id, req.Name, req.Sort)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, nil)
}

func (h *AdminHandler) DeleteCategory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	ec := h.svc.DeleteCategory(id)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, nil)
}

func (h *AdminHandler) GetCategories(c *gin.Context) {
	list, ec := h.svc.GetCategories()
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, gin.H{"list": list})
}

// --- Dashboard ---

func (h *AdminHandler) GetDashboard(c *gin.Context) {
	stats, ec := h.svc.GetDashboard()
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, stats)
}

// --- Report ---

func (h *AdminHandler) CreateReport(c *gin.Context) {
	var req struct {
		TargetType int8   `json:"target_type" binding:"required"`
		TargetID   string `json:"target_id" binding:"required"`
		Reason     string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}
	userID := middleware.GetUserID(c)
	ec := h.svc.CreateReport(userID, req.TargetType, req.TargetID, req.Reason)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, gin.H{"message": "report submitted"})
}

func (h *AdminHandler) GetReports(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	list, total, ec := h.svc.GetReports(page, pageSize)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, gin.H{"list": list, "total": total, "page": page, "page_size": pageSize})
}

func (h *AdminHandler) HandleReport(c *gin.Context) {
	reportID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Status int8 `json:"status" binding:"required,oneof=1 2"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}
	ec := h.svc.HandleReport(reportID, req.Status)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, nil)
}
