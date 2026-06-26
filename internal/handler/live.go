package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/gin-gonic/gin"

	"godan/internal/middleware"
	"godan/internal/model"
	"godan/internal/pkg/errcode"
	"godan/internal/pkg/response"
	"godan/internal/service"
)

var liveUpgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

// --- Live Chat Hub ---

type liveClient struct {
	conn   *websocket.Conn
	roomID uint64
	send   chan []byte
}

type LiveHub struct {
	rooms      map[uint64]map[*liveClient]bool
	register   chan *liveClient
	unregister chan *liveClient
	broadcast  chan *model.LiveMsg
	mu         sync.RWMutex
}

func NewLiveHub() *LiveHub {
	h := &LiveHub{
		rooms:      make(map[uint64]map[*liveClient]bool),
		register:   make(chan *liveClient),
		unregister: make(chan *liveClient),
		broadcast:  make(chan *model.LiveMsg, 256),
	}
	go h.run()
	return h
}

func (h *LiveHub) run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			if h.rooms[c.roomID] == nil {
				h.rooms[c.roomID] = make(map[*liveClient]bool)
			}
			h.rooms[c.roomID][c] = true
			h.mu.Unlock()

		case c := <-h.unregister:
			h.mu.Lock()
			if room, ok := h.rooms[c.roomID]; ok {
				delete(room, c)
				close(c.send)
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			data, _ := json.Marshal(msg)
			h.mu.RLock()
			room := h.rooms[msg.RoomID]
			h.mu.RUnlock()
			for c := range room {
				select {
				case c.send <- data:
				default:
				}
			}
		}
	}
}

func (h *LiveHub) HandleWebSocket(c *gin.Context) {
	roomID, _ := strconv.ParseUint(c.Query("room_id"), 10, 64)
	conn, _ := liveUpgrader.Upgrade(c.Writer, c.Request, nil)

	client := &liveClient{conn: conn, roomID: roomID, send: make(chan []byte, 64)}
	h.register <- client

	go func() {
		defer func() {
			h.unregister <- client
			conn.Close()
		}()
		for msg := range client.send {
			conn.WriteMessage(websocket.TextMessage, msg)
		}
	}()

	go func() {
		defer func() {
			h.unregister <- client
			conn.Close()
		}()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}
			var msg model.LiveMsg
			if json.Unmarshal(message, &msg) != nil {
				continue
			}
			msg.RoomID = roomID
			h.broadcast <- &msg
		}
	}()
}

// --- Live Handler ---

type LiveHandler struct {
	svc *service.LiveService
	hub *LiveHub
}

func NewLiveHandler(svc *service.LiveService, hub *LiveHub) *LiveHandler {
	return &LiveHandler{svc: svc, hub: hub}
}

func (h *LiveHandler) CreateRoom(c *gin.Context) {
	var req struct {
		Title    string `json:"title" binding:"required"`
		CoverURL string `json:"cover_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}
	userID := middleware.GetUserID(c)
	id, ec := h.svc.CreateRoom(userID, req.Title, req.CoverURL)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, gin.H{"room_id": id})
}

func (h *LiveHandler) StartLive(c *gin.Context) {
	roomID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID := middleware.GetUserID(c)

	rtmpURL, playURL, ec := h.svc.StartLive(userID, roomID)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, gin.H{"rtmp_url": rtmpURL, "play_url": playURL})
}

func (h *LiveHandler) StopLive(c *gin.Context) {
	roomID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID := middleware.GetUserID(c)
	ec := h.svc.StopLive(userID, roomID)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, nil)
}

func (h *LiveHandler) GetRoomInfo(c *gin.Context) {
	roomID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	r, ec := h.svc.GetRoomInfo(roomID)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, r)
}

func (h *LiveHandler) GetLiveList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	list, total, ec := h.svc.GetLiveList(page, pageSize)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, gin.H{"list": list, "total": total, "page": page, "page_size": pageSize})
}

func (h *LiveHandler) UpdateRoomInfo(c *gin.Context) {
	roomID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Title    string `json:"title"`
		CoverURL string `json:"cover_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}
	userID := middleware.GetUserID(c)
	ec := h.svc.UpdateRoomInfo(userID, roomID, req.Title, req.CoverURL)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, nil)
}

// --- Gift ---

func (h *LiveHandler) GetGiftList(c *gin.Context) {
	gifts, ec := h.svc.GetGiftList()
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, gin.H{"list": gifts})
}

func (h *LiveHandler) SendGift(c *gin.Context) {
	var req struct {
		RoomID uint64 `json:"room_id" binding:"required"`
		GiftID uint64 `json:"gift_id" binding:"required"`
		Count  int    `json:"count" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}
	userID := middleware.GetUserID(c)
	msg, ec := h.svc.SendGift(userID, req.RoomID, req.GiftID, req.Count)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	// broadcast liveMsg for gift
	h.hub.broadcast <- &model.LiveMsg{
		RoomID: req.RoomID, UserID: userID,
		Content: msg, Color: "#FFD700",
	}
	response.Success(c, gin.H{"message": "gift sent"})
}

func (h *LiveHandler) GetGiftRank(c *gin.Context) {
	roomID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	ranks, ec := h.svc.GetGiftRank(roomID)
	if ec != nil {
		response.Error(c, ec)
		return
	}
	response.Success(c, gin.H{"list": ranks})
}
