package model

import "time"

type Category struct {
	ID        uint64    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Sort      int       `db:"sort" json:"sort"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type DailyStats struct {
	Date       string `db:"date" json:"date"`
	NewUsers   int64  `db:"new_users" json:"new_users"`
	NewVideos  int64  `db:"new_videos" json:"new_videos"`
	TotalPlays int64  `db:"total_plays" json:"total_plays"`
}

type Report struct {
	ID         uint64    `db:"id" json:"id"`
	UserID     uint64    `db:"user_id" json:"user_id"`
	TargetType int8      `db:"target_type" json:"target_type"` // 0:video, 1:comment
	TargetID   string    `db:"target_id" json:"target_id"`
	Reason     string    `db:"reason" json:"reason"`
	Status     int8      `db:"status" json:"status"` // 0:pending, 1:resolved, 2:dismissed
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}
