package model

import "time"

type VideoLike struct {
	ID        uint64    `db:"id" json:"id"`
	UserID    uint64    `db:"user_id" json:"user_id"`
	VideoID   uint64    `db:"video_id" json:"video_id"`
	Type      int8      `db:"type" json:"type"` // 1=like, -1=dislike
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type VideoCoin struct {
	ID        uint64    `db:"id" json:"id"`
	UserID    uint64    `db:"user_id" json:"user_id"`
	VideoID   uint64    `db:"video_id" json:"video_id"`
	Count     int       `db:"count" json:"count"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type FavoriteFolder struct {
	ID          uint64    `db:"id" json:"id"`
	UserID      uint64    `db:"user_id" json:"user_id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	IsPublic    int8      `db:"is_public" json:"is_public"`
	Count       int       `db:"count" json:"count"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type FavoriteItem struct {
	ID        uint64    `db:"id" json:"id"`
	FolderID  uint64    `db:"folder_id" json:"folder_id"`
	VideoID   uint64    `db:"video_id" json:"video_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Danmaku struct {
	ID        uint64    `db:"id" json:"id"`
	VideoID   uint64    `db:"video_id" json:"video_id"`
	UserID    uint64    `db:"user_id" json:"user_id"`
	Content   string    `db:"content" json:"content"`
	Color     string    `db:"color" json:"color"`
	Type      int8      `db:"type" json:"type"`       // 0:scroll, 1:top, 2:bottom
	Position  int       `db:"position" json:"position"` // milliseconds in video
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type DanmakuMsg struct {
	VideoID  uint64 `json:"video_id"`
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	Content  string `json:"content"`
	Color    string `json:"color"`
	Type     int8   `json:"type"`
	Position int    `json:"position"`
}
