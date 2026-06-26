package model

import "time"

type User struct {
	ID           uint64     `db:"id" json:"id"`
	Username     string     `db:"username" json:"username"`
	Email        string     `db:"email" json:"-"`
	Phone        string     `db:"phone" json:"-"`
	PasswordHash string     `db:"password_hash" json:"-"`
	Avatar       string     `db:"avatar" json:"avatar"`
	Bio          string     `db:"bio" json:"bio"`
	Birthday     *time.Time `db:"birthday" json:"birthday"`
	Gender       int8       `db:"gender" json:"gender"`
	Role         int8       `db:"role" json:"role"` // 0:user, 1:admin, 2:super_admin
	Status       int8       `db:"status" json:"status"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
}

type UserProfile struct {
	ID            uint64     `json:"id"`
	Username      string     `json:"username"`
	Avatar        string     `json:"avatar"`
	Bio           string     `json:"bio"`
	Birthday      *time.Time `json:"birthday"`
	Gender        int8       `json:"gender"`
	VideoCount    int64      `json:"video_count"`
	FollowerCount int64      `json:"follower_count"`
	FolloweeCount int64      `json:"followee_count"`
	CreatedAt     time.Time  `json:"created_at"`
}

type Follow struct {
	ID         uint64    `db:"id" json:"id"`
	FollowerID uint64    `db:"follower_id" json:"follower_id"`
	FolloweeID uint64    `db:"followee_id" json:"followee_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type Blocklist struct {
	ID            uint64    `db:"id" json:"id"`
	UserID        uint64    `db:"user_id" json:"user_id"`
	BlockedUserID uint64    `db:"blocked_user_id" json:"blocked_user_id"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}
