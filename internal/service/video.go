package service

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"

	"godan/internal/config"
	"godan/internal/dao"
	"godan/internal/model"
	"godan/internal/pkg/errcode"
	"godan/internal/pkg/logger"
	"godan/internal/pkg/redis"
	"godan/internal/pkg/storage"
)

const videoChunkSize = 5 << 20 // 5MB

type uploadState struct {
	Key        string              `json:"key"`
	UploadID   string              `json:"upload_id"`
	FileSize   int64               `json:"file_size"`
	ChunkCount int                 `json:"chunk_count"`
	Parts      []storage.UploadPart `json:"parts"`
	CT         string              `json:"content_type"`
}

type VideoService struct {
	store    storage.Storage
	uploader storage.MultipartUploader
	cfg      *config.Config
}

func NewVideoService(store storage.Storage, uploader storage.MultipartUploader, cfg *config.Config) *VideoService {
	return &VideoService{store: store, uploader: uploader, cfg: cfg}
}

func (s *VideoService) InitUpload(userID uint64, filename string, fileSize int64) (string, int, int, *errcode.ErrorCode) {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" || !isVideoExt(ext) {
		return "", 0, 0, &errcode.ErrorCode{Code: 30002, Message: "unsupported video format"}
	}

	key := storage.GenVideoKey(userID, filename)
	ct := mimeByExt(ext)

	uploadID, err := s.uploader.InitUpload(context.Background(), key, ct)
	if err != nil {
		logger.Log.Error("init multipart upload failed", zap.Error(err))
		return "", 0, 0, errcode.ErrInternal
	}

	chunkCount := int((fileSize + videoChunkSize - 1) / videoChunkSize)

	state := uploadState{
		Key:        key,
		UploadID:   uploadID,
		FileSize:   fileSize,
		ChunkCount: chunkCount,
		CT:         ct,
	}

	data, _ := json.Marshal(state)
	stateKey := fmt.Sprintf("video:upload:%s", uploadID)
	if err := redis.Set(context.Background(), stateKey, string(data), 24*time.Hour); err != nil {
		logger.Log.Error("redis set upload state failed", zap.Error(err))
		return "", 0, 0, errcode.ErrInternal
	}

	return uploadID, chunkCount, videoChunkSize, nil
}

func (s *VideoService) UploadChunk(uploadID string, chunkIndex int, data []byte) *errcode.ErrorCode {
	state, ec := s.getUploadState(uploadID)
	if ec != nil {
		return ec
	}

	if chunkIndex < 0 || chunkIndex >= state.ChunkCount {
		return &errcode.ErrorCode{Code: 30003, Message: "invalid chunk index"}
	}

	partNumber := chunkIndex + 1
	err := s.uploader.UploadPart(context.Background(), state.Key, state.UploadID, partNumber, strings.NewReader(string(data)), int64(len(data)))
	if err != nil {
		logger.Log.Error("upload part failed", zap.Error(err))
		return errcode.ErrInternal
	}

	found := false
	for i, p := range state.Parts {
		if p.PartNumber == partNumber {
			state.Parts[i].PartSize = len(data)
			found = true
			break
		}
	}
	if !found {
		state.Parts = append(state.Parts, storage.UploadPart{PartNumber: partNumber, PartSize: len(data)})
	}

	s.saveUploadState(uploadID, state)
	return nil
}

func (s *VideoService) CompleteUpload(userID uint64, uploadID, title, description string, categoryID int, tags string) (uint64, string, *errcode.ErrorCode) {
	state, ec := s.getUploadState(uploadID)
	if ec != nil {
		return 0, "", ec
	}

	if len(state.Parts) != state.ChunkCount {
		return 0, "", &errcode.ErrorCode{Code: 30004, Message: "upload not complete, missing chunks"}
	}

	videoURL, err := s.uploader.CompleteUpload(context.Background(), state.Key, state.UploadID, state.Parts)
	if err != nil {
		logger.Log.Error("complete multipart upload failed", zap.Error(err))
		return 0, "", errcode.ErrInternal
	}

	video := &model.Video{
		UserID:      userID,
		Title:       title,
		Description: description,
		VideoURL:    videoURL,
		CategoryID:  categoryID,
		Tags:        tags,
		FileSize:    state.FileSize,
		Status:      model.VideoStatusPending,
	}

	id, err := dao.CreateVideo(video)
	if err != nil {
		logger.Log.Error("create video record failed", zap.Error(err))
		return 0, "", errcode.ErrInternal
	}

	s.delUploadState(uploadID)
	return id, videoURL, nil
}

func (s *VideoService) UploadStatus(uploadID string) ([]int, int, *errcode.ErrorCode) {
	state, ec := s.getUploadState(uploadID)
	if ec != nil {
		return nil, 0, ec
	}

	var uploaded []int
	for _, p := range state.Parts {
		uploaded = append(uploaded, p.PartNumber-1)
	}
	return uploaded, state.ChunkCount, nil
}

func (s *VideoService) AbortUpload(uploadID string) *errcode.ErrorCode {
	state, ec := s.getUploadState(uploadID)
	if ec != nil {
		return ec
	}

	if err := s.uploader.AbortUpload(context.Background(), state.Key, state.UploadID); err != nil {
		logger.Log.Error("abort upload failed", zap.Error(err))
	}
	s.delUploadState(uploadID)
	return nil
}

func (s *VideoService) getUploadState(uploadID string) (*uploadState, *errcode.ErrorCode) {
	stateKey := fmt.Sprintf("video:upload:%s", uploadID)
	val, err := redis.Get(context.Background(), stateKey)
	if err != nil || val == "" {
		return nil, &errcode.ErrorCode{Code: 30005, Message: "upload not found or expired"}
	}

	var state uploadState
	if err := json.Unmarshal([]byte(val), &state); err != nil {
		return nil, errcode.ErrInternal
	}
	return &state, nil
}

func (s *VideoService) saveUploadState(uploadID string, state *uploadState) {
	data, _ := json.Marshal(state)
	redis.Set(context.Background(), fmt.Sprintf("video:upload:%s", uploadID), string(data), 24*time.Hour)
}

func (s *VideoService) delUploadState(uploadID string) {
	redis.Del(context.Background(), fmt.Sprintf("video:upload:%s", uploadID))
}

func isVideoExt(ext string) bool {
	switch ext {
	case ".mp4", ".avi", ".mov", ".mkv", ".flv", ".wmv", ".webm", ".m4v":
		return true
	}
	return false
}

func mimeByExt(ext string) string {
	switch ext {
	case ".mp4":
		return "video/mp4"
	case ".webm":
		return "video/webm"
	case ".avi":
		return "video/x-msvideo"
	case ".mov":
		return "video/quicktime"
	case ".mkv":
		return "video/x-matroska"
	case ".flv":
		return "video/x-flv"
	case ".wmv":
		return "video/x-ms-wmv"
	default:
		return "application/octet-stream"
	}
}
