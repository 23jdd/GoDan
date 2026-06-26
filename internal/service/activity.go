package service

import (
	"fmt"

	"go.uber.org/zap"

	"godan/internal/dao"
	"godan/internal/model"
	"godan/internal/pkg/errcode"
	"godan/internal/pkg/logger"
)

type ActivityService struct{}

func NewActivityService() *ActivityService {
	return &ActivityService{}
}

func (s *ActivityService) CreateUploadActivity(userID, videoID uint64) {
	_ = dao.CreateActivity(&model.Activity{
		UserID:     userID,
		Type:       model.ActivityUpload,
		TargetID:   videoID,
		TargetType: 0,
	})
}

func (s *ActivityService) CreateLikeActivity(userID, videoID uint64) {
	_ = dao.CreateActivity(&model.Activity{
		UserID: userID, Type: model.ActivityLike, TargetID: videoID,
	})
}

func (s *ActivityService) CreateCoinActivity(userID, videoID uint64) {
	_ = dao.CreateActivity(&model.Activity{
		UserID: userID, Type: model.ActivityCoin, TargetID: videoID,
	})
}

func (s *ActivityService) CreateFavActivity(userID, videoID uint64) {
	_ = dao.CreateActivity(&model.Activity{
		UserID: userID, Type: model.ActivityFav, TargetID: videoID,
	})
}

func (s *ActivityService) CreateShareActivity(userID, videoID uint64) {
	_ = dao.CreateActivity(&model.Activity{
		UserID: userID, Type: model.ActivityShare, TargetID: videoID,
	})
}

func (s *ActivityService) GetTimeline(userID uint64, page, pageSize int) ([]model.ActivityWithUser, int64, *errcode.ErrorCode) {
	offset, limit := paginate(page, pageSize)
	list, total, err := dao.GetFollowTimeline(userID, offset, limit)
	if err != nil {
		logger.Log.Error("get timeline failed", zap.Error(err))
		return nil, 0, errcode.ErrInternal
	}
	if list == nil {
		list = []model.ActivityWithUser{}
	}
	return list, total, nil
}

func (s *ActivityService) GetUserActivities(userID uint64, page, pageSize int) ([]model.Activity, int64, *errcode.ErrorCode) {
	offset, limit := paginate(page, pageSize)
	list, total, err := dao.GetUserActivities(userID, offset, limit)
	if err != nil {
		logger.Log.Error("get user activities failed", zap.Error(err))
		return nil, 0, errcode.ErrInternal
	}
	if list == nil {
		list = []model.Activity{}
	}
	return list, total, nil
}

// --- Notification ---

type NotificationService struct {
	hub *NotificationHub
}

func NewNotificationService() *NotificationService {
	svc := &NotificationService{}
	svc.hub = NewNotificationHub()
	return svc
}

func (s *NotificationService) Hub() *NotificationHub {
	return s.hub
}

func (s *NotificationService) Send(userID uint64, notifType int8, title, content string, targetID uint64) {
	n := &model.Notification{
		UserID:   userID,
		Type:     notifType,
		Title:    title,
		Content:  content,
		TargetID: targetID,
	}
	id, err := dao.CreateNotification(n)
	if err != nil {
		logger.Log.Error("create notification failed", zap.Error(err))
		return
	}
	n.ID = id

	count, _ := dao.GetUnreadNotificationCount(userID)
	s.hub.Push(userID, fmt.Sprintf(`{"type":%d,"title":"%s","content":"%s","unread":%d}`, notifType, title, content, count))
}

func (s *NotificationService) GetNotifications(userID uint64, page, pageSize int) ([]model.Notification, int64, *errcode.ErrorCode) {
	offset, limit := paginate(page, pageSize)
	list, total, err := dao.GetNotifications(userID, offset, limit)
	if err != nil {
		logger.Log.Error("get notifications failed", zap.Error(err))
		return nil, 0, errcode.ErrInternal
	}
	if list == nil {
		list = []model.Notification{}
	}
	return list, total, nil
}

func (s *NotificationService) GetUnreadCount(userID uint64) (int64, *errcode.ErrorCode) {
	count, err := dao.GetUnreadNotificationCount(userID)
	if err != nil {
		return 0, errcode.ErrInternal
	}
	return count, nil
}

func (s *NotificationService) MarkRead(userID, notifID uint64) *errcode.ErrorCode {
	if err := dao.MarkNotificationRead(notifID, userID); err != nil {
		return errcode.ErrInternal
	}
	return nil
}

func (s *NotificationService) MarkAllRead(userID uint64) *errcode.ErrorCode {
	if err := dao.MarkAllNotificationsRead(userID); err != nil {
		return errcode.ErrInternal
	}
	return nil
}
