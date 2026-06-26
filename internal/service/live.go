package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"godan/internal/dao"
	"godan/internal/model"
	"godan/internal/pkg/errcode"
	"godan/internal/pkg/logger"
	"godan/internal/pkg/redis"
)

type LiveService struct{}

func NewLiveService() *LiveService {
	return &LiveService{}
}

func (s *LiveService) CreateRoom(userID uint64, title, coverURL string) (uint64, *errcode.ErrorCode) {
	r := &model.LiveRoom{UserID: userID, Title: title, CoverURL: coverURL}
	id, err := dao.CreateLiveRoom(r)
	if err != nil {
		logger.Log.Error("create room failed", zap.Error(err))
		return 0, errcode.ErrInternal
	}
	return id, nil
}

func (s *LiveService) GetRoomInfo(roomID uint64) (*model.LiveRoom, *errcode.ErrorCode) {
	r, err := dao.GetLiveRoomByID(roomID)
	if err != nil {
		return nil, errcode.ErrVideoNotFound
	}
	return r, nil
}

func (s *LiveService) StartLive(userID, roomID uint64) (string, string, *errcode.ErrorCode) {
	r, err := dao.GetLiveRoomByID(roomID)
	if err != nil || r.UserID != userID {
		return "", "", errcode.ErrForbidden
	}

	key, _ := dao.GetRoomStreamKey(roomID, userID)
	dao.UpdateRoomStatus(roomID, model.RoomLive)

	rtmpURL := fmt.Sprintf("rtmp://127.0.0.1:1935/live/%s", key)
	playURL := fmt.Sprintf("/live/%d/play.m3u8", roomID)

	return rtmpURL, playURL, nil
}

func (s *LiveService) StopLive(userID, roomID uint64) *errcode.ErrorCode {
	r, err := dao.GetLiveRoomByID(roomID)
	if err != nil || r.UserID != userID {
		return errcode.ErrForbidden
	}
	dao.UpdateRoomStatus(roomID, model.RoomOff)
	return nil
}

func (s *LiveService) GetLiveList(page, pageSize int) ([]model.LiveRoom, int64, *errcode.ErrorCode) {
	offset, limit := paginate(page, pageSize)
	list, total, err := dao.GetLiveRoomList(offset, limit)
	if err != nil {
		return nil, 0, errcode.ErrInternal
	}
	if list == nil {
		list = []model.LiveRoom{}
	}
	return list, total, nil
}

func (s *LiveService) UpdateRoomInfo(userID, roomID uint64, title, coverURL string) *errcode.ErrorCode {
	r, err := dao.GetLiveRoomByID(roomID)
	if err != nil || r.UserID != userID {
		return errcode.ErrForbidden
	}
	dao.UpdateRoomInfo(roomID, title, coverURL)
	return nil
}

// --- Online Count ---

func (s *LiveService) JoinRoom(roomID uint64, userID string) int {
	redis.Set(context.Background(), "live:viewer:"+userID, roomID, 0)
	redis.Set(context.Background(), fmt.Sprintf("live:room:%d:viewer:%s", roomID, userID), "1", 0)
	return s.countViewers(roomID)
}

func (s *LiveService) LeaveRoom(roomID uint64, userID string) int {
	redis.Del(context.Background(), fmt.Sprintf("live:room:%d:viewer:%s", roomID, userID))
	return s.countViewers(roomID)
}

func (s *LiveService) countViewers(roomID uint64) int {
	ctx := context.Background()
	keys, _ := redis.RDB.Keys(ctx, fmt.Sprintf("live:room:%d:viewer:*", roomID)).Result()
	return len(keys)
}

// --- Gift ---

func (s *LiveService) GetGiftList() ([]model.Gift, *errcode.ErrorCode) {
	gifts, err := dao.GetGiftList()
	if err != nil {
		return nil, errcode.ErrInternal
	}
	if gifts == nil {
		gifts = []model.Gift{}
	}
	return gifts, nil
}

func (s *LiveService) SendGift(userID, roomID, giftID uint64, count int) (string, *errcode.ErrorCode) {
	gift, _ := dao.GetGiftByID(giftID)
	if gift == nil || gift.ID == 0 {
		return "", &errcode.ErrorCode{Code: 60001, Message: "gift not found"}
	}

	totalCoin := gift.Price * count

	record := &model.GiftRecord{
		UserID:    userID,
		RoomID:    roomID,
		GiftID:    giftID,
		Count:     count,
		TotalCoin: totalCoin,
	}
	if err := dao.CreateGiftRecord(record); err != nil {
		return "", errcode.ErrInternal
	}

	msg := fmt.Sprintf(`{"user_id":%d,"gift":"%s","count":%d,"total_coin":%d}`, userID, gift.Name, count, totalCoin)
	return msg, nil
}

func (s *LiveService) GetGiftRank(roomID uint64) (any, *errcode.ErrorCode) {
	ranks, err := dao.GetGiftRank(roomID, 20)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	return ranks, nil
}
