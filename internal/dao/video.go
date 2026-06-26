package dao

import (
	"database/sql"
	"errors"

	"godan/internal/model"
	"godan/internal/pkg/database"
)

func CreateVideo(v *model.Video) (uint64, error) {
	result, err := database.DB.Exec(
		`INSERT INTO videos (user_id, title, description, cover_url, video_url, duration, category_id, tags, file_size, status)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		v.UserID, v.Title, v.Description, v.CoverURL, v.VideoURL,
		v.Duration, v.CategoryID, v.Tags, v.FileSize, v.Status,
	)
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	return uint64(id), nil
}

func GetVideoByID(id uint64) (*model.Video, error) {
	var v model.Video
	err := database.DB.Get(&v, "SELECT * FROM videos WHERE id = ?", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &v, nil
}

func GetVideosByUserID(userID uint64, offset, limit int) ([]model.Video, int64, error) {
	var total int64
	if err := database.DB.Get(&total, "SELECT COUNT(*) FROM videos WHERE user_id = ? AND status = 1", userID); err != nil {
		return nil, 0, err
	}

	var videos []model.Video
	err := database.DB.Select(&videos,
		"SELECT * FROM videos WHERE user_id = ? AND status = 1 ORDER BY created_at DESC LIMIT ? OFFSET ?",
		userID, limit, offset,
	)
	return videos, total, err
}

func GetVideosByCategory(categoryID int, offset, limit int) ([]model.Video, int64, error) {
	var total int64
	if err := database.DB.Get(&total, "SELECT COUNT(*) FROM videos WHERE category_id = ? AND status = 1", categoryID); err != nil {
		return nil, 0, err
	}

	var videos []model.Video
	err := database.DB.Select(&videos,
		"SELECT * FROM videos WHERE category_id = ? AND status = 1 ORDER BY created_at DESC LIMIT ? OFFSET ?",
		categoryID, limit, offset,
	)
	return videos, total, err
}

func UpdateVideoStatus(id uint64, status int8) error {
	_, err := database.DB.Exec("UPDATE videos SET status = ?, updated_at = NOW() WHERE id = ?", status, id)
	return err
}

func IncrVideoPlayCount(id uint64) error {
	_, err := database.DB.Exec("UPDATE videos SET play_count = play_count + 1 WHERE id = ?", id)
	return err
}
