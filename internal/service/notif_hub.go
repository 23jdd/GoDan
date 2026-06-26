package service

import (
	"github.com/gorilla/websocket"
)

type notifClient struct {
	conn   *websocket.Conn
	userID uint64
	send   chan string
}

type NotificationHub struct {
	clients    map[uint64]*notifClient
	register   chan *notifClient
	unregister chan *notifClient
}

func NewNotificationHub() *NotificationHub {
	h := &NotificationHub{
		clients:    make(map[uint64]*notifClient),
		register:   make(chan *notifClient),
		unregister: make(chan *notifClient),
	}
	go h.run()
	return h
}

func (h *NotificationHub) run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c.userID] = c
		case c := <-h.unregister:
			if existing, ok := h.clients[c.userID]; ok && existing == c {
				delete(h.clients, c.userID)
				close(c.send)
			}
		}
	}
}

func (h *NotificationHub) Push(userID uint64, msg string) {
	if c, ok := h.clients[userID]; ok {
		select {
		case c.send <- msg:
		default:
		}
	}
}

func (h *NotificationHub) AddClient(userID uint64, conn *websocket.Conn) *notifClient {
	c := &notifClient{
		conn:   conn,
		userID: userID,
		send:   make(chan string, 32),
	}
	h.register <- c
	return c
}

func (h *NotificationHub) RemoveClient(userID uint64, c *notifClient) {
	h.unregister <- c
}

func (c *notifClient) WritePump() {
	defer c.conn.Close()
	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			break
		}
	}
}
