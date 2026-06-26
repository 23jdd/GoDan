package dao

import (
	"godan/internal/model"
	"godan/internal/pkg/database"
)

const activityLimit = 50

func CreateActivity(a *model.Activity) error {
	_, err := database.DB.Exec(
		"INSERT INTO activities (user_id, type, target_id, target_type) VALUES (?, ?, ?, ?)",
		a.UserID, a.Type, a.TargetID, a.TargetType,
	)
	return err
}

func GetFollowTimeline(userID uint64, offset, limit int) ([]model.ActivityWithUser, int64, error) {
	var total int64
	database.DB.Get(&total, `
		SELECT COUNT(*) FROM activities a
		INNER JOIN follows f ON f.followee_id = a.user_id
		WHERE f.follower_id = ?`, userID)

	var list []model.ActivityWithUser
	err := database.DB.Select(&list, `
		SELECT a.*, u.username, u.avatar FROM activities a
		INNER JOIN follows f ON f.followee_id = a.user_id
		INNER JOIN users u ON u.id = a.user_id
		WHERE f.follower_id = ?
		ORDER BY a.created_at DESC LIMIT ? OFFSET ?`, userID, limit, offset)

	return list, total, err
}

func GetUserActivities(userID uint64, offset, limit int) ([]model.Activity, int64, error) {
	var total int64
	database.DB.Get(&total, "SELECT COUNT(*) FROM activities WHERE user_id = ?", userID)

	var list []model.Activity
	err := database.DB.Select(&list,
		"SELECT * FROM activities WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?",
		userID, limit, offset)
	return list, total, err
}

// --- Notification ---

func CreateNotification(n *model.Notification) (uint64, error) {
	result, err := database.DB.Exec(
		`INSERT INTO notifications (user_id, type, title, content, target_id) VALUES (?, ?, ?, ?, ?)`,
		n.UserID, n.Type, n.Title, n.Content, n.TargetID,
	)
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	return uint64(id), nil
}

func GetNotifications(userID uint64, offset, limit int) ([]model.Notification, int64, error) {
	var total int64
	database.DB.Get(&total, "SELECT COUNT(*) FROM notifications WHERE user_id = ?", userID)

	var list []model.Notification
	err := database.DB.Select(&list,
		"SELECT * FROM notifications WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?",
		userID, limit, offset)
	return list, total, err
}

func GetUnreadNotificationCount(userID uint64) (int64, error) {
	var count int64
	err := database.DB.Get(&count,
		"SELECT COUNT(*) FROM notifications WHERE user_id = ? AND is_read = 0", userID)
	return count, err
}

func MarkNotificationRead(notifID, userID uint64) error {
	_, err := database.DB.Exec(
		"UPDATE notifications SET is_read = 1 WHERE id = ? AND user_id = ?", notifID, userID)
	return err
}

func MarkAllNotificationsRead(userID uint64) error {
	_, err := database.DB.Exec(
		"UPDATE notifications SET is_read = 1 WHERE user_id = ?", userID)
	return err
}
