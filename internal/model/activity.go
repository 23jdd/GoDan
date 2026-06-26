package model

import "time"

type Activity struct {
	ID         uint64    `db:"id" json:"id"`
	UserID     uint64    `db:"user_id" json:"user_id"`
	Type       int8      `db:"type" json:"type"` // 1:upload 2:like 3:coin 4:fav 5:share
	TargetID   uint64    `db:"target_id" json:"target_id"`
	TargetType int8      `db:"target_type" json:"target_type"` // 0:video
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

const (
	ActivityUpload = 1
	ActivityLike   = 2
	ActivityCoin   = 3
	ActivityFav    = 4
	ActivityShare  = 5
)

type ActivityWithUser struct {
	Activity
	Username string `db:"username" json:"username"`
	Avatar   string `db:"avatar" json:"avatar"`
}

type Notification struct {
	ID        uint64    `db:"id" json:"id"`
	UserID    uint64    `db:"user_id" json:"user_id"`
	Type      int8      `db:"type" json:"type"` // 1:like 2:comment 3:reply 4:follow 5:system
	Title     string    `db:"title" json:"title"`
	Content   string    `db:"content" json:"content"`
	TargetID  uint64    `db:"target_id" json:"target_id"`
	IsRead    int8      `db:"is_read" json:"is_read"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

const (
	NotifLike    = 1
	NotifComment = 2
	NotifReply   = 3
	NotifFollow  = 4
	NotifSystem  = 5
)
