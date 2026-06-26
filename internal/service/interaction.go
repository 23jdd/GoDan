package service

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"godan/internal/config"
	"godan/internal/dao"
	"godan/internal/model"
	"godan/internal/pkg/errcode"
	"godan/internal/pkg/logger"
	"godan/internal/pkg/redis"
)

const dailyCoinLimit = 5

type InteractionService struct {
	cfg *config.Config
}

func NewInteractionService(cfg *config.Config) *InteractionService {
	return &InteractionService{cfg: cfg}
}

// --- Like ---

func (s *InteractionService) LikeVideo(userID, videoID uint64) *errcode.ErrorCode {
	v, _ := dao.GetVideoByID(videoID)
	if v == nil || v.Status != model.VideoStatusPublished {
		return errcode.ErrVideoNotFound
	}

	if err := dao.UpsertVideoLike(userID, videoID, 1); err != nil {
		logger.Log.Error("like video failed", zap.Error(err))
		return errcode.ErrInternal
	}
	dao.UpdateVideoLikeCount(videoID, 1)
	return nil
}

func (s *InteractionService) DislikeVideo(userID, videoID uint64) *errcode.ErrorCode {
	if err := dao.UpsertVideoLike(userID, videoID, -1); err != nil {
		logger.Log.Error("dislike video failed", zap.Error(err))
		return errcode.ErrInternal
	}
	return nil
}

func (s *InteractionService) CancelLikeDislike(userID, videoID uint64) *errcode.ErrorCode {
	prev, _ := dao.GetUserLikeStatus(userID, videoID)
	if prev == 1 {
		dao.UpdateVideoLikeCount(videoID, -1)
	}
	if err := dao.DeleteVideoLike(userID, videoID); err != nil {
		logger.Log.Error("cancel like failed", zap.Error(err))
		return errcode.ErrInternal
	}
	return nil
}

func (s *InteractionService) GetUserLikeStatus(userID, videoID uint64) (int8, *errcode.ErrorCode) {
	t, _ := dao.GetUserLikeStatus(userID, videoID)
	return t, nil
}

// --- Coin ---

func (s *InteractionService) GiveCoin(userID, videoID uint64, count int) *errcode.ErrorCode {
	if count < 1 || count > 2 {
		return &errcode.ErrorCode{Code: 30010, Message: "coin count must be 1 or 2"}
	}

	today := time.Now().Format("2006-01-02")
	dailyCoins, _ := dao.GetUserDailyCoins(userID, today)
	if dailyCoins+count > dailyCoinLimit {
		return &errcode.ErrorCode{Code: 30011, Message: "daily coin limit exceeded"}
	}

	if err := dao.AddVideoCoin(userID, videoID, count); err != nil {
		logger.Log.Error("give coin failed", zap.Error(err))
		return errcode.ErrInternal
	}
	dao.UpdateVideoCoinCount(videoID, count)
	return nil
}

func (s *InteractionService) GetUserDailyRemainingCoins(userID uint64) (int, *errcode.ErrorCode) {
	today := time.Now().Format("2006-01-02")
	dailyCoins, _ := dao.GetUserDailyCoins(userID, today)
	remaining := dailyCoinLimit - dailyCoins
	if remaining < 0 {
		remaining = 0
	}
	return remaining, nil
}

// --- Favorite Folder ---

func (s *InteractionService) CreateFolder(userID uint64, name, description string, isPublic int8) (uint64, *errcode.ErrorCode) {
	f := &model.FavoriteFolder{
		UserID:      userID,
		Name:        name,
		Description: description,
		IsPublic:    isPublic,
	}
	id, err := dao.CreateFavoriteFolder(f)
	if err != nil {
		logger.Log.Error("create folder failed", zap.Error(err))
		return 0, errcode.ErrInternal
	}
	return id, nil
}

func (s *InteractionService) UpdateFolder(userID, folderID uint64, name, description string, isPublic int8) *errcode.ErrorCode {
	if err := dao.UpdateFavoriteFolder(folderID, userID, name, description, isPublic); err != nil {
		logger.Log.Error("update folder failed", zap.Error(err))
		return errcode.ErrInternal
	}
	return nil
}

func (s *InteractionService) DeleteFolder(userID, folderID uint64) *errcode.ErrorCode {
	if err := dao.DeleteFavoriteFolder(folderID, userID); err != nil {
		logger.Log.Error("delete folder failed", zap.Error(err))
		return errcode.ErrInternal
	}
	return nil
}

func (s *InteractionService) GetUserFolders(userID uint64) ([]model.FavoriteFolder, *errcode.ErrorCode) {
	folders, err := dao.GetFavoriteFoldersByUser(userID)
	if err != nil {
		logger.Log.Error("get folders failed", zap.Error(err))
		return nil, errcode.ErrInternal
	}
	if folders == nil {
		folders = []model.FavoriteFolder{}
	}
	return folders, nil
}

func (s *InteractionService) AddToFolder(userID, folderID, videoID uint64) *errcode.ErrorCode {
	f, _ := dao.GetFavoriteFolderByID(folderID, userID)
	if f == nil || f.ID == 0 {
		return &errcode.ErrorCode{Code: 30012, Message: "folder not found"}
	}
	if err := dao.AddFavoriteItem(folderID, videoID); err != nil {
		logger.Log.Error("add to folder failed", zap.Error(err))
		return errcode.ErrInternal
	}
	return nil
}

func (s *InteractionService) RemoveFromFolder(userID, folderID, videoID uint64) *errcode.ErrorCode {
	f, _ := dao.GetFavoriteFolderByID(folderID, userID)
	if f == nil || f.ID == 0 {
		return &errcode.ErrorCode{Code: 30012, Message: "folder not found"}
	}
	if err := dao.RemoveFavoriteItem(folderID, videoID); err != nil {
		logger.Log.Error("remove from folder failed", zap.Error(err))
		return errcode.ErrInternal
	}
	return nil
}

func (s *InteractionService) GetFolderItems(folderID, page, pageSize int) ([]model.Video, int64, *errcode.ErrorCode) {
	offset, limit := paginate(page, pageSize)
	videos, total, err := dao.GetFavoriteItems(uint64(folderID), offset, limit)
	if err != nil {
		logger.Log.Error("get folder items failed", zap.Error(err))
		return nil, 0, errcode.ErrInternal
	}
	if videos == nil {
		videos = []model.Video{}
	}
	return videos, total, nil
}

// --- Share ---

func (s *InteractionService) ShareVideo(videoID uint64) (string, *errcode.ErrorCode) {
	v, _ := dao.GetVideoByID(videoID)
	if v == nil || v.Status != model.VideoStatusPublished {
		return "", errcode.ErrVideoNotFound
	}
	dao.IncrVideoShareCount(videoID)
	return fmt.Sprintf("/video/%d", videoID), nil
}

// --- Play Count De-duplication ---

func (s *InteractionService) RecordPlay(userID uint64, videoID uint64, ip string) bool {
	ctx := context.Background()
	key := fmt.Sprintf("play:%d:%s", videoID, ip)
	if userID > 0 {
		key = fmt.Sprintf("play:%d:%d", videoID, userID)
	}
	exists, _ := redis.Exists(ctx, key)
	if exists {
		return false
	}
	redis.Set(ctx, key, "1", 5*time.Minute)
	return true
}

// --- Danmaku ---

func (s *InteractionService) SendDanmaku(userID uint64, msg model.DanmakuMsg) (*model.Danmaku, *errcode.ErrorCode) {
	if len(msg.Content) == 0 || len(msg.Content) > 200 {
		return nil, &errcode.ErrorCode{Code: 40010, Message: "danmaku content must be 1-200 chars"}
	}

	d := &model.Danmaku{
		VideoID:  msg.VideoID,
		UserID:   userID,
		Content:  msg.Content,
		Color:    msg.Color,
		Type:     msg.Type,
		Position: msg.Position,
	}
	if d.Color == "" {
		d.Color = "#FFFFFF"
	}

	id, err := dao.CreateDanmaku(d)
	if err != nil {
		logger.Log.Error("create danmaku failed", zap.Error(err))
		return nil, errcode.ErrInternal
	}
	d.ID = id
	return d, nil
}

func (s *InteractionService) GetDanmakus(videoID uint64, start, end int) ([]model.Danmaku, *errcode.ErrorCode) {
	list, err := dao.GetDanmakusByTimeRange(videoID, start, end)
	if err != nil {
		logger.Log.Error("get danmakus failed", zap.Error(err))
		return nil, errcode.ErrInternal
	}
	if list == nil {
		list = []model.Danmaku{}
	}
	return list, nil
}
