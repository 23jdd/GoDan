package dao

import (
	"fmt"

	"github.com/google/uuid"

	"godan/internal/model"
	"godan/internal/pkg/database"
)

// --- Live Room ---

func CreateLiveRoom(r *model.LiveRoom) (uint64, error) {
	r.StreamKey = fmt.Sprintf("live_%d_%s", r.UserID, uuid.New().String()[:8])
	result, err := database.DB.Exec(
		"INSERT INTO live_rooms (user_id, title, cover_url, stream_key, status) VALUES (?, ?, ?, ?, ?)",
		r.UserID, r.Title, r.CoverURL, r.StreamKey, model.RoomOff,
	)
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	return uint64(id), nil
}

func GetLiveRoomByID(id uint64) (*model.LiveRoom, error) {
	var r model.LiveRoom
	err := database.DB.Get(&r,
		`SELECT lr.*, u.username, u.avatar FROM live_rooms lr
		 INNER JOIN users u ON u.id = lr.user_id WHERE lr.id = ?`, id)
	return &r, err
}

func GetLiveRoomByUserID(userID uint64) (*model.LiveRoom, error) {
	var r model.LiveRoom
	err := database.DB.Get(&r, "SELECT * FROM live_rooms WHERE user_id = ?", userID)
	return &r, err
}

func UpdateRoomStatus(roomID uint64, status int8) error {
	_, err := database.DB.Exec(
		"UPDATE live_rooms SET status = ?, updated_at = NOW() WHERE id = ?", status, roomID)
	return err
}

func UpdateRoomInfo(roomID uint64, title, coverURL string) error {
	_, err := database.DB.Exec(
		"UPDATE live_rooms SET title = ?, cover_url = ?, updated_at = NOW() WHERE id = ?",
		title, coverURL, roomID)
	return err
}

func GetLiveRoomList(offset, limit int) ([]model.LiveRoom, int64, error) {
	var total int64
	database.DB.Get(&total, "SELECT COUNT(*) FROM live_rooms WHERE status = ?", model.RoomLive)

	var list []model.LiveRoom
	err := database.DB.Select(&list,
		`SELECT lr.*, u.username, u.avatar FROM live_rooms lr
		 INNER JOIN users u ON u.id = lr.user_id
		 WHERE lr.status = ? ORDER BY lr.viewer_count DESC LIMIT ? OFFSET ?`,
		model.RoomLive, limit, offset)
	return list, total, err
}

func GetRoomStreamKey(roomID, userID uint64) (string, error) {
	var key string
	err := database.DB.Get(&key,
		"SELECT stream_key FROM live_rooms WHERE id = ? AND user_id = ?", roomID, userID)
	return key, err
}

// --- Gift System ---

func GetGiftList() ([]model.Gift, error) {
	var gifts []model.Gift
	err := database.DB.Select(&gifts, "SELECT * FROM gifts ORDER BY price ASC")
	return gifts, err
}

func GetGiftByID(id uint64) (*model.Gift, error) {
	var g model.Gift
	err := database.DB.Get(&g, "SELECT * FROM gifts WHERE id = ?", id)
	return &g, err
}

func CreateGiftRecord(r *model.GiftRecord) error {
	_, err := database.DB.Exec(
		"INSERT INTO gift_records (user_id, room_id, gift_id, count, total_coin) VALUES (?, ?, ?, ?, ?)",
		r.UserID, r.RoomID, r.GiftID, r.Count, r.TotalCoin,
	)
	return err
}

func GetGiftRank(roomID uint64, limit int) ([]struct {
	UserID uint64 `db:"user_id"`
	Name   string `db:"name"`
	Total  int    `db:"total_coin"`
}, error) {
	var ranks []struct {
		UserID uint64 `db:"user_id"`
		Name   string `db:"name"`
		Total  int    `db:"total_coin"`
	}
	err := database.DB.Select(&ranks,
		`SELECT gr.user_id, u.username AS name, SUM(gr.total_coin) AS total_coin
		 FROM gift_records gr
		 INNER JOIN users u ON u.id = gr.user_id
		 WHERE gr.room_id = ?
		 GROUP BY gr.user_id ORDER BY total_coin DESC LIMIT ?`, roomID, limit)
	return ranks, err
}

// --- Seed Default Gifts ---

func SeedGifts() {
	var count int
	database.DB.Get(&count, "SELECT COUNT(*) FROM gifts")
	if count > 0 {
		return
	}
	gifts := []struct{ name, icon string; price int }{
		{"辣条", "🌶", 1},
		{"棒棒糖", "🍭", 5},
		{"小星星", "⭐", 10},
		{"火箭", "🚀", 100},
		{"嘉年华", "🎪", 500},
	}
	for _, g := range gifts {
		database.DB.Exec("INSERT INTO gifts (name, icon, price) VALUES (?, ?, ?)", g.name, g.icon, g.price)
	}
}
