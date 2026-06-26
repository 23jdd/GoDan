package model

import "time"

type Video struct {
	ID          uint64    `db:"id" json:"id"`
	UserID      uint64    `db:"user_id" json:"user_id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	CoverURL    string    `db:"cover_url" json:"cover_url"`
	VideoURL    string    `db:"video_url" json:"video_url"`
	Duration    int       `db:"duration" json:"duration"`
	CategoryID  int       `db:"category_id" json:"category_id"`
	Tags        string    `db:"tags" json:"tags"`
	FileSize    int64     `db:"file_size" json:"file_size"`
	Status      int8      `db:"status" json:"status"`       // 0:pending 1:published 2:rejected 3:removed
	PlayCount   int64     `db:"play_count" json:"play_count"`
	LikeCount   int64     `db:"like_count" json:"like_count"`
	CoinCount   int64     `db:"coin_count" json:"coin_count"`
	FavCount    int64     `db:"fav_count" json:"fav_count"`
	ShareCount  int64     `db:"share_count" json:"share_count"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

const (
	VideoStatusPending   = 0
	VideoStatusPublished = 1
	VideoStatusRejected  = 2
	VideoStatusRemoved   = 3
)
