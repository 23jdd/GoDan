package dao

import (
	"godan/internal/pkg/database"
)

func AddBlock(userID, blockedUserID uint64) error {
	_, err := database.DB.Exec(
		"INSERT INTO user_blocklists (user_id, blocked_user_id) VALUES (?, ?)",
		userID, blockedUserID,
	)
	return err
}

func RemoveBlock(userID, blockedUserID uint64) error {
	_, err := database.DB.Exec(
		"DELETE FROM user_blocklists WHERE user_id = ? AND blocked_user_id = ?",
		userID, blockedUserID,
	)
	return err
}

func IsBlocked(userID, targetUserID uint64) (bool, error) {
	var count int
	err := database.DB.Get(&count,
		"SELECT COUNT(*) FROM user_blocklists WHERE user_id = ? AND blocked_user_id = ?",
		userID, targetUserID,
	)
	return count > 0, err
}
