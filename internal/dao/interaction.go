package dao

import (
	"godan/internal/model"
	"godan/internal/pkg/database"
)

// --- Like/Dislike ---

func UpsertVideoLike(userID, videoID uint64, likeType int8) error {
	_, err := database.DB.Exec(
		`INSERT INTO video_likes (user_id, video_id, type) VALUES (?, ?, ?)
		 ON DUPLICATE KEY UPDATE type = VALUES(type)`,
		userID, videoID, likeType,
	)
	return err
}

func DeleteVideoLike(userID, videoID uint64) error {
	_, err := database.DB.Exec("DELETE FROM video_likes WHERE user_id = ? AND video_id = ?", userID, videoID)
	return err
}

func GetUserLikeStatus(userID, videoID uint64) (int8, error) {
	var t int8
	err := database.DB.Get(&t, "SELECT type FROM video_likes WHERE user_id = ? AND video_id = ?", userID, videoID)
	return t, err
}

func UpdateVideoLikeCount(videoID uint64, delta int) error {
	_, err := database.DB.Exec("UPDATE videos SET like_count = like_count + ? WHERE id = ?", delta, videoID)
	return err
}

// --- Coin ---

func GetUserDailyCoins(userID uint64, date string) (int, error) {
	var total int
	err := database.DB.Get(&total,
		"SELECT COALESCE(SUM(count), 0) FROM video_coins WHERE user_id = ? AND DATE(created_at) = ?",
		userID, date,
	)
	return total, err
}

func AddVideoCoin(userID, videoID uint64, count int) error {
	_, err := database.DB.Exec(
		"INSERT INTO video_coins (user_id, video_id, count) VALUES (?, ?, ?)",
		userID, videoID, count,
	)
	return err
}

func UpdateVideoCoinCount(videoID uint64, delta int) error {
	_, err := database.DB.Exec("UPDATE videos SET coin_count = coin_count + ? WHERE id = ?", delta, videoID)
	return err
}

// --- Favorite ---

func CreateFavoriteFolder(f *model.FavoriteFolder) (uint64, error) {
	result, err := database.DB.Exec(
		"INSERT INTO favorite_folders (user_id, name, description, is_public) VALUES (?, ?, ?, ?)",
		f.UserID, f.Name, f.Description, f.IsPublic,
	)
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	return uint64(id), nil
}

func UpdateFavoriteFolder(folderID, userID uint64, name, description string, isPublic int8) error {
	_, err := database.DB.Exec(
		`UPDATE favorite_folders SET name=?, description=?, is_public=?, updated_at=NOW()
		 WHERE id=? AND user_id=?`, name, description, isPublic, folderID, userID,
	)
	return err
}

func DeleteFavoriteFolder(folderID, userID uint64) error {
	_, err := database.DB.Exec("DELETE FROM favorite_folders WHERE id=? AND user_id=?", folderID, userID)
	return err
}

func GetFavoriteFoldersByUser(userID uint64) ([]model.FavoriteFolder, error) {
	var folders []model.FavoriteFolder
	err := database.DB.Select(&folders,
		"SELECT * FROM favorite_folders WHERE user_id=? ORDER BY created_at DESC", userID,
	)
	return folders, err
}

func GetFavoriteFolderByID(folderID, userID uint64) (*model.FavoriteFolder, error) {
	var f model.FavoriteFolder
	err := database.DB.Get(&f,
		"SELECT * FROM favorite_folders WHERE id=? AND user_id=?", folderID, userID,
	)
	return &f, err
}

func AddFavoriteItem(folderID, videoID uint64) error {
	_, err := database.DB.Exec(
		`INSERT INTO favorite_items (folder_id, video_id) VALUES (?, ?)
		 ON DUPLICATE KEY UPDATE id=id`, folderID, videoID,
	)
	if err != nil {
		return err
	}
	_, _ = database.DB.Exec("UPDATE favorite_folders SET count = count + 1 WHERE id=?", folderID)
	_, _ = database.DB.Exec("UPDATE videos SET fav_count = fav_count + 1 WHERE id=?", videoID)
	return nil
}

func RemoveFavoriteItem(folderID, videoID uint64) error {
	_, err := database.DB.Exec("DELETE FROM favorite_items WHERE folder_id=? AND video_id=?", folderID, videoID)
	if err != nil {
		return err
	}
	_, _ = database.DB.Exec("UPDATE favorite_folders SET count = GREATEST(count - 1, 0) WHERE id=?", folderID)
	_, _ = database.DB.Exec("UPDATE videos SET fav_count = GREATEST(fav_count - 1, 0) WHERE id=?", videoID)
	return nil
}

func GetFavoriteItems(folderID uint64, offset, limit int) ([]model.Video, int64, error) {
	var total int64
	database.DB.Get(&total, "SELECT COUNT(*) FROM favorite_items WHERE folder_id=?", folderID)

	var videos []model.Video
	err := database.DB.Select(&videos,
		`SELECT v.* FROM videos v
		 INNER JOIN favorite_items fi ON fi.video_id = v.id
		 WHERE fi.folder_id = ? ORDER BY fi.created_at DESC LIMIT ? OFFSET ?`,
		folderID, limit, offset,
	)
	return videos, total, err
}

func IsVideoFavorited(folderID, videoID uint64) (bool, error) {
	var count int
	err := database.DB.Get(&count,
		"SELECT COUNT(*) FROM favorite_items WHERE folder_id=? AND video_id=?", folderID, videoID,
	)
	return count > 0, err
}

// --- Danmaku ---

func CreateDanmaku(d *model.Danmaku) (uint64, error) {
	result, err := database.DB.Exec(
		`INSERT INTO danmakus (video_id, user_id, content, color, type, position)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		d.VideoID, d.UserID, d.Content, d.Color, d.Type, d.Position,
	)
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	return uint64(id), nil
}

func GetDanmakusByTimeRange(videoID uint64, start, end int) ([]model.Danmaku, error) {
	var list []model.Danmaku
	err := database.DB.Select(&list,
		`SELECT d.id, d.video_id, d.user_id, d.content, d.color, d.type, d.position, d.created_at
		 FROM danmakus d
		 WHERE d.video_id = ? AND d.position BETWEEN ? AND ?
		 ORDER BY d.position ASC LIMIT 500`,
		videoID, start, end,
	)
	return list, err
}

func GetDanmakuCount(videoID uint64) (int64, error) {
	var count int64
	err := database.DB.Get(&count, "SELECT COUNT(*) FROM danmakus WHERE video_id = ?", videoID)
	return count, err
}

func IncrVideoShareCount(videoID uint64) error {
	_, err := database.DB.Exec("UPDATE videos SET share_count = share_count + 1 WHERE id = ?", videoID)
	return err
}
