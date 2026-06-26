package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"godan/internal/middleware"
	"godan/internal/pkg/errcode"
	"godan/internal/pkg/response"
	"godan/internal/service"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// Register godoc
// @Summary 用户注册
// @Tags user
// @Accept json
// @Produce json
// @Param body body RegisterReq true "注册参数"
// @Success 200 {object} response.Response
// @Router /api/v1/user/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=2,max=30"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}

	user, ec := h.svc.Register(req.Username, req.Email, req.Phone, req.Password)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, gin.H{"user_id": user.ID})
}

// Login godoc
// @Summary 用户登录
// @Tags user
// @Accept json
// @Produce json
// @Param body body LoginReq true "登录参数"
// @Success 200 {object} response.Response
// @Router /api/v1/user/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Account  string `json:"account" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}

	user, accessToken, refreshToken, ec := h.svc.Login(req.Account, req.Password)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, gin.H{
		"user_id":       user.ID,
		"username":      user.Username,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    h.svc.Cfg().JWT.AccessExpire,
	})
}

// RefreshToken godoc
// @Summary 刷新AccessToken
// @Tags user
// @Accept json
// @Produce json
// @Param body body RefreshTokenReq true "刷新参数"
// @Success 200 {object} response.Response
// @Router /api/v1/user/refresh [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}

	accessToken, ec := h.svc.RefreshAccessToken(req.RefreshToken)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, gin.H{
		"access_token": accessToken,
		"expires_in":   h.svc.Cfg().JWT.AccessExpire,
	})
}

// GetProfile godoc
// @Summary 获取当前用户信息
// @Tags user
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Router /api/v1/user/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	profile, ec := h.svc.GetProfile(userID)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, profile)
}

// GetUserProfile godoc
// @Summary 获取指定用户主页
// @Tags user
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response
// @Router /api/v1/user/profile/{id} [get]
func (h *UserHandler) GetUserProfile(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}

	profile, ec := h.svc.GetProfile(userID)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, profile)
}

// UpdateProfile godoc
// @Summary 修改个人资料
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body UpdateProfileReq true "修改参数"
// @Success 200 {object} response.Response
// @Router /api/v1/user/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req struct {
		Username string     `json:"username"`
		Avatar   string     `json:"avatar"`
		Bio      string     `json:"bio"`
		Birthday *time.Time `json:"birthday"`
		Gender   int8       `json:"gender" binding:"omitempty,oneof=0 1 2"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}

	userID := middleware.GetUserID(c)
	ec := h.svc.UpdateProfile(userID, req.Username, req.Avatar, req.Bio, req.Birthday, req.Gender)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, nil)
}

// ChangePassword godoc
// @Summary 修改密码
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body ChangePasswordReq true "修改密码参数"
// @Success 200 {object} response.Response
// @Router /api/v1/user/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6,max=128"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}

	userID := middleware.GetUserID(c)
	ec := h.svc.ChangePassword(userID, req.OldPassword, req.NewPassword)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, nil)
}

// BindEmail godoc
// @Summary 绑定邮箱
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body BindEmailReq true "绑定邮箱参数"
// @Success 200 {object} response.Response
// @Router /api/v1/user/bind/email [post]
func (h *UserHandler) BindEmail(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}

	userID := middleware.GetUserID(c)
	ec := h.svc.BindEmail(userID, req.Email, req.Code)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, nil)
}

// BindPhone godoc
// @Summary 绑定手机号
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body BindPhoneReq true "绑定手机参数"
// @Success 200 {object} response.Response
// @Router /api/v1/user/bind/phone [post]
func (h *UserHandler) BindPhone(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}

	userID := middleware.GetUserID(c)
	ec := h.svc.BindPhone(userID, req.Phone, req.Code)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, nil)
}

// SendVerificationCode godoc
// @Summary 发送验证码
// @Tags user
// @Accept json
// @Produce json
// @Param body body SendCodeReq true "发送验证码参数"
// @Success 200 {object} response.Response
// @Router /api/v1/user/code/send [post]
func (h *UserHandler) SendVerificationCode(c *gin.Context) {
	var req struct {
		Target string `json:"target" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return
	}

	ec := h.svc.SendVerificationCode(req.Target)
	if ec != nil {
		response.Error(c, ec)
		return
	}

	response.Success(c, gin.H{"message": "verification code sent"})
}
