package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"godan/internal/model"
	"godan/internal/pkg/logger"
	"godan/internal/pkg/response"
	"godan/internal/service"
)

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type danmakuClient struct {
	conn    *websocket.Conn
	videoID uint64
	send    chan []byte
}

type DanmakuHub struct {
	clients    map[uint64]map[*danmakuClient]bool
	broadcast  chan *model.DanmakuMsg
	register   chan *danmakuClient
	unregister chan *danmakuClient
	svc        *service.InteractionService
}

func NewDanmakuHub(svc *service.InteractionService) *DanmakuHub {
	hub := &DanmakuHub{
		clients:    make(map[uint64]map[*danmakuClient]bool),
		broadcast:  make(chan *model.DanmakuMsg, 256),
		register:   make(chan *danmakuClient),
		unregister: make(chan *danmakuClient),
		svc:        svc,
	}
	go hub.run()
	return hub
}

func (h *DanmakuHub) run() {
	for {
		select {
		case client := <-h.register:
			if h.clients[client.videoID] == nil {
				h.clients[client.videoID] = make(map[*danmakuClient]bool)
			}
			h.clients[client.videoID][client] = true

		case client := <-h.unregister:
			if clients, ok := h.clients[client.videoID]; ok {
				delete(clients, client)
				close(client.send)
			}

		case msg := <-h.broadcast:
			data, _ := json.Marshal(msg)
			if clients, ok := h.clients[msg.VideoID]; ok {
				for client := range clients {
					select {
					case client.send <- data:
					default:
						close(client.send)
						delete(clients, client)
					}
				}
			}
		}
	}
}

func (h *DanmakuHub) HandleWebSocket(c *gin.Context) {
	videoID, err := strconv.ParseUint(c.Query("video_id"), 10, 64)
	if err != nil {
		return
	}

	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Log.Error("ws upgrade failed", zap.Error(err))
		return
	}

	client := &danmakuClient{
		conn:    conn,
		videoID: videoID,
		send:    make(chan []byte, 64),
	}

	h.register <- client

	go client.writePump()
	go client.readPump(h)
}

func (c *danmakuClient) readPump(hub *DanmakuHub) {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		var msg model.DanmakuMsg
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}
		msg.VideoID = c.videoID

		d, ec := hub.svc.SendDanmaku(msg.UserID, msg)
		if ec != nil {
			continue
		}

		hub.broadcast <- &model.DanmakuMsg{
			VideoID:  d.VideoID,
			UserID:   d.UserID,
			Content:  d.Content,
			Color:    d.Color,
			Type:     d.Type,
			Position: d.Position,
		}
	}
}

func (c *danmakuClient) writePump() {
	defer c.conn.Close()

	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
}

// --- REST Danmaku ---

type DanmakuHandler struct {
	svc *service.InteractionService
}

func NewDanmakuHandler(svc *service.InteractionService) *DanmakuHandler {
	return &DanmakuHandler{svc: svc}
}

// GetDanmakus godoc
// @Summary 获取分时段弹幕
// @Tags danmaku
// @Produce json
// @Param video_id query int true "视频ID"
// @Param start query int true "起始毫秒"
// @Param end query int true "结束毫秒"
// @Success 200 {object} response.Response
// @Router /api/v1/danmakus [get]
func (h *DanmakuHandler) GetDanmakus(c *gin.Context) {
	videoID, _ := strconv.ParseUint(c.Query("video_id"), 10, 64)
	start, _ := strconv.Atoi(c.DefaultQuery("start", "0"))
	end, _ := strconv.Atoi(c.DefaultQuery("end", "60000"))

	list, ec := h.svc.GetDanmakus(videoID, start, end)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, gin.H{"list": list})
}
