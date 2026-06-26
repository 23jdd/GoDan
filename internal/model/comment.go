package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Comment struct {
	ID         bson.ObjectID `bson:"_id,omitempty" json:"id"`
	VideoID    uint64        `bson:"video_id" json:"video_id"`
	UserID     uint64        `bson:"user_id" json:"user_id"`
	Username   string        `bson:"username" json:"username"`
	Avatar     string        `bson:"avatar" json:"avatar"`
	Content    string        `bson:"content" json:"content"`
	ParentID   string        `bson:"parent_id" json:"parent_id"`     // "" = root
	RootID     string        `bson:"root_id" json:"root_id"`         // root comment id
	ReplyToUID uint64        `bson:"reply_to_uid" json:"reply_to_uid"`
	ReplyCount int64         `bson:"reply_count" json:"reply_count"`
	LikeCount  int64         `bson:"like_count" json:"like_count"`
	Status     int8          `bson:"status" json:"status"` // 0:normal 1:deleted
	CreatedAt  time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time     `bson:"updated_at" json:"updated_at"`
}

const (
	CommentNormal  = 0
	CommentDeleted = 1
)

type CommentListResp struct {
	List     []Comment `json:"list"`
	Total    int64     `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"page_size"`
}
