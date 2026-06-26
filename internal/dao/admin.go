package dao

import (
	"time"

	"godan/internal/model"
	"godan/internal/pkg/database"
)

// --- User Admin ---

func GetUserList(offset, limit int) ([]model.User, int64, error) {
	var total int64
	database.DB.Get(&total, "SELECT COUNT(*) FROM users")

	var users []model.User
	err := database.DB.Select(&users,
		"SELECT * FROM users ORDER BY created_at DESC LIMIT ? OFFSET ?", limit, offset)
	return users, total, err
}

func UpdateUserRole(userID uint64, role int8) error {
	_, err := database.DB.Exec("UPDATE users SET role = ? WHERE id = ?", role, userID)
	return err
}

func UpdateUserStatus(userID uint64, status int8) error {
	_, err := database.DB.Exec("UPDATE users SET status = ? WHERE id = ?", status, userID)
	return err
}

// --- Video Review ---

func GetPendingVideos(offset, limit int) ([]model.Video, int64, error) {
	var total int64
	database.DB.Get(&total, "SELECT COUNT(*) FROM videos WHERE status = ?", model.VideoStatusPending)

	var videos []model.Video
	err := database.DB.Select(&videos,
		"SELECT * FROM videos WHERE status = ? ORDER BY created_at DESC LIMIT ? OFFSET ?",
		model.VideoStatusPending, limit, offset)
	return videos, total, err
}

// --- Category ---

func CreateCategory(name string, sort int) (uint64, error) {
	result, err := database.DB.Exec("INSERT INTO categories (name, sort) VALUES (?, ?)", name, sort)
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	return uint64(id), nil
}

func UpdateCategory(id uint64, name string, sort int) error {
	_, err := database.DB.Exec("UPDATE categories SET name=?, sort=? WHERE id=?", name, sort, id)
	return err
}

func DeleteCategory(id uint64) error {
	_, err := database.DB.Exec("DELETE FROM categories WHERE id=?", id)
	return err
}

func GetCategoryList() ([]model.Category, error) {
	var list []model.Category
	err := database.DB.Select(&list, "SELECT * FROM categories ORDER BY sort ASC")
	return list, err
}

// --- Stats ---

func GetDailyStats(date string) (*model.DailyStats, error) {
	var s model.DailyStats
	s.Date = date

	database.DB.Get(&s.NewUsers,
		"SELECT COUNT(*) FROM users WHERE DATE(created_at) = ?", date)
	database.DB.Get(&s.NewVideos,
		"SELECT COUNT(*) FROM videos WHERE DATE(created_at) = ?", date)
	database.DB.Get(&s.TotalPlays,
		`SELECT COALESCE(SUM(play_count), 0) FROM videos WHERE status = 1`)

	return &s, nil
}

// --- Report ---

func CreateReport(userID uint64, targetType int8, targetID, reason string) error {
	_, err := database.DB.Exec(
		"INSERT INTO reports (user_id, target_type, target_id, reason) VALUES (?, ?, ?, ?)",
		userID, targetType, targetID, reason,
	)
	return err
}

func GetReportList(offset, limit int) ([]model.Report, int64, error) {
	var total int64
	database.DB.Get(&total, "SELECT COUNT(*) FROM reports")

	var list []model.Report
	err := database.DB.Select(&list,
		"SELECT * FROM reports ORDER BY created_at DESC LIMIT ? OFFSET ?", limit, offset)
	return list, total, err
}

func UpdateReportStatus(id uint64, status int8) error {
	_, err := database.DB.Exec("UPDATE reports SET status = ? WHERE id = ?", status, id)
	return err
}

func GetReportByID(id uint64) (*model.Report, error) {
	var r model.Report
	err := database.DB.Get(&r, "SELECT * FROM reports WHERE id = ?", id)
	return &r, err
}

// --- Dashboard Today ---

func GetDashboardStats() (map[string]int64, error) {
	var totalUsers, todayUsers int64
	var totalVideos, pendingVideos int64
	var totalPlays int64

	today := time.Now().Format("2006-01-02")

	database.DB.Get(&totalUsers, "SELECT COUNT(*) FROM users")
	database.DB.Get(&todayUsers, "SELECT COUNT(*) FROM users WHERE DATE(created_at) = ?", today)
	database.DB.Get(&totalVideos, "SELECT COUNT(*) FROM videos WHERE status = 1")
	database.DB.Get(&pendingVideos, "SELECT COUNT(*) FROM videos WHERE status = 0")
	database.DB.Get(&totalPlays, "SELECT COALESCE(SUM(play_count), 0) FROM videos WHERE status = 1")

	return map[string]int64{
		"total_users":    totalUsers,
		"today_users":    todayUsers,
		"total_videos":   totalVideos,
		"pending_videos": pendingVideos,
		"total_plays":    totalPlays,
	}, nil
}
