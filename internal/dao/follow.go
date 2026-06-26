package dao

import (
	"errors"
	"fmt"

	"godan/internal/model"
	"godan/internal/pkg/database"
)

func CreateFollow(followerID, followeeID uint64) error {
	_, err := database.DB.Exec(
		"INSERT INTO follows (follower_id, followee_id) VALUES (?, ?)",
		followerID, followeeID,
	)
	return err
}

func DeleteFollow(followerID, followeeID uint64) error {
	_, err := database.DB.Exec(
		"DELETE FROM follows WHERE follower_id = ? AND followee_id = ?",
		followerID, followeeID,
	)
	return err
}

func IsFollowing(followerID, followeeID uint64) (bool, error) {
	var count int
	err := database.DB.Get(&count,
		"SELECT COUNT(*) FROM follows WHERE follower_id = ? AND followee_id = ?",
		followerID, followeeID,
	)
	return count > 0, err
}

func GetFollowerList(userID uint64, offset, limit int) ([]model.User, int64, error) {
	var total int64
	if err := database.DB.Get(&total,
		"SELECT COUNT(*) FROM follows WHERE followee_id = ?", userID); err != nil {
		return nil, 0, err
	}

	var users []model.User
	err := database.DB.Select(&users, `
		SELECT u.* FROM users u
		INNER JOIN follows f ON f.follower_id = u.id
		WHERE f.followee_id = ?
		ORDER BY f.created_at DESC
		LIMIT ? OFFSET ?`, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("GetFollowerList: %w", err)
	}
	return users, total, nil
}

func GetFolloweeList(userID uint64, offset, limit int) ([]model.User, int64, error) {
	var total int64
	if err := database.DB.Get(&total,
		"SELECT COUNT(*) FROM follows WHERE follower_id = ?", userID); err != nil {
		return nil, 0, err
	}

	var users []model.User
	err := database.DB.Select(&users, `
		SELECT u.* FROM users u
		INNER JOIN follows f ON f.followee_id = u.id
		WHERE f.follower_id = ?
		ORDER BY f.created_at DESC
		LIMIT ? OFFSET ?`, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("GetFolloweeList: %w", err)
	}
	return users, total, nil
}

func GetMutualFollows(userID uint64, offset, limit int) ([]model.User, int64, error) {
	var total int64
	if err := database.DB.Get(&total, `
		SELECT COUNT(*) FROM follows f1
		INNER JOIN follows f2 ON f1.followee_id = f2.follower_id AND f1.follower_id = f2.followee_id
		WHERE f1.follower_id = ?`, userID); err != nil {
		return nil, 0, err
	}

	var users []model.User
	err := database.DB.Select(&users, `
		SELECT u.* FROM users u
		INNER JOIN follows f1 ON f1.followee_id = u.id
		INNER JOIN follows f2 ON f2.follower_id = f1.followee_id AND f2.followee_id = f1.follower_id
		WHERE f1.follower_id = ?
		ORDER BY f1.created_at DESC
		LIMIT ? OFFSET ?`, userID, limit, offset)
	if errors.Is(err, nil) {
		return users, total, nil
	}
	return nil, 0, fmt.Errorf("GetMutualFollows: %w", err)
}
