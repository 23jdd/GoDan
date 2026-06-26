package handler

import (
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"godan/internal/middleware"
	"godan/internal/pkg/errcode"
	"godan/internal/pkg/logger"
	"godan/internal/pkg/response"
	"godan/internal/service"
)

type NotificationHandler struct {
	svc  *service.NotificationService
	hub  *service.NotificationHub
	up   websocket.Upgrader
}

func NewNotificationHandler(svc *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		svc: svc,
		hub: svc.Hub(),
		up:  websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
	}
}

// HandleWebSocket godoc
// @Summary 通知 WebSocket
// @Tags notification
// @Param token query string true "JWT Token"
// @Router /api/v1/notifications/ws [get]
func (h *NotificationHandler) HandleWebSocket(c *gin.Context) {
	userID := middleware.GetUserID(c)

	conn, err := h.up.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Log.Error("notif ws upgrade failed", zap.Error(err))
		return
	}

	client := h.hub.AddClient(userID, conn)
	go client.WritePump()
}

// GetNotifications godoc
// @Summary 通知列表
// @Tags notification
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码"
// @Param page_size query int false "每页"
// @Success 200 {object} response.Response
// @Router /api/v1/notifications [get]
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	list, total, ec := h.svc.GetNotifications(userID, page, pageSize)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, gin.H{"list": list, "total": total, "page": page, "page_size": pageSize})
}

// GetUnreadCount godoc
// @Summary 未读通知数
// @Tags notification
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Router /api/v1/notifications/unread [get]
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID := middleware.GetUserID(c)
	count, ec := h.svc.GetUnreadCount(userID)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, gin.H{"unread": count})
}

// MarkRead godoc
// @Summary 标记已读
// @Tags notification
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "通知ID"
// @Success 200 {object} response.Response
// @Router /api/v1/notifications/{id}/read [post]
func (h *NotificationHandler) MarkRead(c *gin.Context) {
	notifID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}
	userID := middleware.GetUserID(c)
	ec := h.svc.MarkRead(userID, notifID)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, nil)
}

// MarkAllRead godoc
// @Summary 全部已读
// @Tags notification
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Router /api/v1/notifications/read-all [post]
func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	userID := middleware.GetUserID(c)
	ec := h.svc.MarkAllRead(userID)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, nil)
}
