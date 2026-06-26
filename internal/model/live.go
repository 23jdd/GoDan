package model

import "time"

type LiveRoom struct {
	ID          uint64    `db:"id" json:"id"`
	UserID      uint64    `db:"user_id" json:"user_id"`
	Title       string    `db:"title" json:"title"`
	CoverURL    string    `db:"cover_url" json:"cover_url"`
	StreamKey   string    `db:"stream_key" json:"-"` // RTMP 推流密钥，不暴露
	Status      int8      `db:"status" json:"status"` // 0:off, 1:live
	ViewerCount int64     `db:"viewer_count" json:"viewer_count"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`

	// 非DB字段
	Username string `db:"-" json:"username"`
	Avatar   string `db:"-" json:"avatar"`
}

type LiveMsg struct {
	RoomID   uint64 `json:"room_id"`
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	Content  string `json:"content"`
	Color    string `json:"color"`
}

type Gift struct {
	ID   uint64 `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Icon string `db:"icon" json:"icon"`
	Price int    `db:"price" json:"price"`
}

type GiftRecord struct {
	ID        uint64    `db:"id" json:"id"`
	UserID    uint64    `db:"user_id" json:"user_id"`
	RoomID    uint64    `db:"room_id" json:"room_id"`
	GiftID    uint64    `db:"gift_id" json:"gift_id"`
	Count     int       `db:"count" json:"count"`
	TotalCoin int       `db:"total_coin" json:"total_coin"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

const (
	RoomOff  = 0
	RoomLive = 1
)
