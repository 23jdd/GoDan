package handler

import (
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"godan/internal/middleware"
	"godan/internal/pkg/errcode"
	"godan/internal/pkg/logger"
	"godan/internal/pkg/response"
	"godan/internal/pkg/storage"
)

const maxAvatarSize = 5 << 20 // 5MB

var allowedTypes = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".gif":  "image/gif",
	".webp": "image/webp",
}

type UploadHandler struct {
	store storage.Storage
}

func NewUploadHandler(store storage.Storage) *UploadHandler {
	return &UploadHandler{store: store}
}

// UploadAvatar godoc
// @Summary 上传头像
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param file formData file true "图片文件"
// @Success 200 {object} response.Response
// @Router /api/v1/upload/avatar [post]
func (h *UploadHandler) UploadAvatar(c *gin.Context) {
	userID := middleware.GetUserID(c)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}
	defer file.Close()

	if header.Size > maxAvatarSize {
		response.ErrorWithMsg(c, errcode.ErrInvalidParams, "file size exceeds 5MB")
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	contentType, ok := allowedTypes[ext]
	if !ok {
		response.ErrorWithMsg(c, errcode.ErrInvalidParams, "unsupported format, allowed: jpg/jpeg/png/gif/webp")
		return
	}

	buf, err := io.ReadAll(file)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}

	detected := http.DetectContentType(buf)
	if !strings.HasPrefix(detected, "image/") {
		response.ErrorWithMsg(c, errcode.ErrInvalidParams, "invalid image file")
		return
	}

	key := storage.GenAvatarKey(userID, ext)
	url, err := h.store.Upload(c.Request.Context(), key, strings.NewReader(string(buf)), int64(len(buf)), contentType)
	if err != nil {
		logger.Log.Error("upload avatar failed", zap.Error(err))
		response.Error(c, errcode.ErrInternal)
		return
	}

	response.Success(c, gin.H{"url": url})
}
