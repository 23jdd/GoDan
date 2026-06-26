package dao

import (
	"database/sql"
	"errors"
	"fmt"

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

func GetVideoDetail(videoID uint64) (*model.VideoDetail, error) {
	var d model.VideoDetail
	err := database.DB.Get(&d, `
		SELECT v.id, v.user_id, v.title, v.description, v.cover_url, v.video_url,
		       v.duration, v.category_id, v.tags, v.file_size, v.status,
		       v.play_count, v.like_count, v.coin_count, v.fav_count, v.share_count,
		       v.created_at, v.updated_at,
		       u.id AS 'author.id', u.username AS 'author.username', u.avatar AS 'author.avatar'
		FROM videos v
		INNER JOIN users u ON u.id = v.user_id
		WHERE v.id = ?`, videoID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &d, nil
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

func GetVideosByCategory(categoryID int, sort string, offset, limit int) ([]model.Video, int64, error) {
	orderBy := orderClause(sort)

	var total int64
	var countSQL string
	if categoryID == 0 {
		countSQL = "SELECT COUNT(*) FROM videos WHERE status = 1"
	} else {
		countSQL = "SELECT COUNT(*) FROM videos WHERE category_id = ? AND status = 1"
	}
	if err := database.DB.Get(&total, countSQL, categoryID); err != nil {
		return nil, 0, err
	}

	var query string
	if categoryID == 0 {
		query = fmt.Sprintf("SELECT * FROM videos WHERE status = 1 %s LIMIT ? OFFSET ?", orderBy)
	} else {
		query = fmt.Sprintf("SELECT * FROM videos WHERE category_id = ? AND status = 1 %s LIMIT ? OFFSET ?", orderBy)
	}

	var videos []model.Video
	var err error
	if categoryID == 0 {
		err = database.DB.Select(&videos, query, limit, offset)
	} else {
		err = database.DB.Select(&videos, query, categoryID, limit, offset)
	}
	return videos, total, err
}

func GetHotVideos(offset, limit int) ([]model.Video, int64, error) {
	var total int64
	database.DB.Get(&total, "SELECT COUNT(*) FROM videos WHERE status = 1")

	var videos []model.Video
	err := database.DB.Select(&videos, `
		SELECT * FROM videos WHERE status = 1
		ORDER BY (play_count * 3 + like_count * 10 + coin_count * 20 + fav_count * 6 + share_count * 4)
		         / POW(TIMESTAMPDIFF(HOUR, created_at, NOW()) + 2, 1.8) DESC
		LIMIT ? OFFSET ?`, limit, offset)
	return videos, total, err
}

func SearchVideos(keyword string, offset, limit int) ([]model.Video, int64, error) {
	kw := "%" + keyword + "%"

	var total int64
	database.DB.Get(&total,
		"SELECT COUNT(*) FROM videos WHERE status = 1 AND (title LIKE ? OR description LIKE ? OR tags LIKE ?)",
		kw, kw, kw,
	)

	var videos []model.Video
	err := database.DB.Select(&videos,
		`SELECT * FROM videos WHERE status = 1
		 AND (title LIKE ? OR description LIKE ? OR tags LIKE ?)
		 ORDER BY play_count DESC LIMIT ? OFFSET ?`,
		kw, kw, kw, limit, offset,
	)
	return videos, total, err
}

func GetRelatedVideos(videoID uint64, categoryID int, tags string, offset, limit int) ([]model.Video, error) {
	var videos []model.Video
	err := database.DB.Select(&videos, `
		SELECT * FROM videos
		WHERE status = 1 AND id != ? AND category_id = ?
		ORDER BY (play_count * 3 + like_count * 10) / POW(TIMESTAMPDIFF(HOUR, created_at, NOW()) + 2, 1.5) DESC
		LIMIT ? OFFSET ?`, videoID, categoryID, limit, offset)
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func UpdateVideoStatus(id uint64, status int8) error {
	_, err := database.DB.Exec("UPDATE videos SET status = ?, updated_at = NOW() WHERE id = ?", status, id)
	return err
}

func UpdateVideoCover(id uint64, coverURL string) error {
	_, err := database.DB.Exec("UPDATE videos SET cover_url = ?, updated_at = NOW() WHERE id = ?", coverURL, id)
	return err
}

func IncrVideoPlayCount(id uint64) error {
	_, err := database.DB.Exec("UPDATE videos SET play_count = play_count + 1 WHERE id = ?", id)
	return err
}

func orderClause(sort string) string {
	switch sort {
	case "hot":
		return "ORDER BY (play_count * 3 + like_count * 10 + coin_count * 20) / POW(TIMESTAMPDIFF(HOUR, created_at, NOW()) + 2, 1.5) DESC"
	default:
		return "ORDER BY created_at DESC"
	}
}
