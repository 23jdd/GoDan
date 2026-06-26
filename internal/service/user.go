package service

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"unicode"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"godan/internal/config"
	"godan/internal/dao"
	"godan/internal/model"
	"godan/internal/pkg/errcode"
	pkgjwt "godan/internal/pkg/jwt"
	"godan/internal/pkg/logger"
	"godan/internal/pkg/redis"
)

type UserService struct {
	cfg *config.Config
}

func NewUserService(cfg *config.Config) *UserService {
	return &UserService{cfg: cfg}
}

func (s *UserService) Register(username, email, phone, password string) (*model.User, *errcode.ErrorCode) {
	if strings.TrimSpace(username) == "" {
		return nil, errcode.ErrInvalidParams
	}
	if email == "" && phone == "" {
		return nil, errcode.ErrInvalidParams
	}
	if !s.isPasswordValid(password) {
		return nil, &errcode.ErrorCode{Code: 20000, Message: "password must be at least 6 characters with letters and numbers"}
	}

	if email != "" {
		existing, _ := dao.GetUserByEmail(email)
		if existing != nil {
			return nil, errcode.ErrUserExists
		}
	}
	if phone != "" {
		existing, _ := dao.GetUserByPhone(phone)
		if existing != nil {
			return nil, errcode.ErrUserExists
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Error("bcrypt hash failed", zap.Error(err))
		return nil, errcode.ErrInternal
	}

	user := &model.User{
		Username:     username,
		Email:        email,
		Phone:        phone,
		PasswordHash: string(hash),
		Status:       1,
	}

	id, err := dao.CreateUser(user)
	if err != nil {
		logger.Log.Error("create user failed", zap.Error(err))
		return nil, errcode.ErrInternal
	}

	user.ID = id
	return user, nil
}

func (s *UserService) RegisterWithCode(username, email, phone, password, code string) (*model.User, *errcode.ErrorCode) {
	if strings.TrimSpace(username) == "" {
		return nil, errcode.ErrInvalidParams
	}
	if email == "" && phone == "" {
		return nil, errcode.ErrInvalidParams
	}
	if code == "" {
		return nil, errcode.ErrInvalidParams
	}
	if !s.isPasswordValid(password) {
		return nil, &errcode.ErrorCode{Code: 20000, Message: "password must be at least 6 characters with letters and numbers"}
	}

	if email != "" {
		if !s.verifyCode("email:"+email, code) {
			return nil, errcode.ErrCodeIncorrect
		}
		existing, _ := dao.GetUserByEmail(email)
		if existing != nil {
			return nil, errcode.ErrUserExists
		}
	}
	if phone != "" {
		if !s.verifyCode("phone:"+phone, code) {
			return nil, errcode.ErrCodeIncorrect
		}
		existing, _ := dao.GetUserByPhone(phone)
		if existing != nil {
			return nil, errcode.ErrUserExists
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Error("bcrypt hash failed", zap.Error(err))
		return nil, errcode.ErrInternal
	}

	user := &model.User{
		Username:     username,
		Email:        email,
		Phone:        phone,
		PasswordHash: string(hash),
		Status:       1,
	}

	id, err := dao.CreateUser(user)
	if err != nil {
		logger.Log.Error("create user failed", zap.Error(err))
		return nil, errcode.ErrInternal
	}

	user.ID = id
	return user, nil
}

func (s *UserService) Login(account, password string) (*model.User, string, string, *errcode.ErrorCode) {
	var user *model.User
	if s.isEmail(account) {
		user, _ = dao.GetUserByEmail(account)
	} else {
		user, _ = dao.GetUserByPhone(account)
	}

	if user == nil {
		return nil, "", "", errcode.ErrUserNotFound
	}
	if user.Status != 1 {
		return nil, "", "", errcode.ErrUserBanned
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", "", errcode.ErrPasswordWrong
	}

	accessToken, err := pkgjwt.GenerateToken(user.ID, s.cfg.JWT.AccessSecret, s.cfg.JWT.AccessExpire)
	if err != nil {
		logger.Log.Error("generate access token failed", zap.Error(err))
		return nil, "", "", errcode.ErrInternal
	}

	refreshToken, err := pkgjwt.GenerateToken(user.ID, s.cfg.JWT.RefreshSecret, s.cfg.JWT.RefreshExpire)
	if err != nil {
		logger.Log.Error("generate refresh token failed", zap.Error(err))
		return nil, "", "", errcode.ErrInternal
	}

	return user, accessToken, refreshToken, nil
}

func (s *UserService) LoginByCode(account, code string) (*model.User, string, string, *errcode.ErrorCode) {
	prefix := "email"
	if !s.isEmail(account) {
		prefix = "phone"
	}

	if !s.verifyCode(prefix+":"+account, code) {
		return nil, "", "", errcode.ErrCodeIncorrect
	}

	var user *model.User
	if prefix == "email" {
		user, _ = dao.GetUserByEmail(account)
	} else {
		user, _ = dao.GetUserByPhone(account)
	}

	if user == nil {
		return nil, "", "", errcode.ErrUserNotFound
	}
	if user.Status != 1 {
		return nil, "", "", errcode.ErrUserBanned
	}

	accessToken, err := pkgjwt.GenerateToken(user.ID, s.cfg.JWT.AccessSecret, s.cfg.JWT.AccessExpire)
	if err != nil {
		logger.Log.Error("generate access token failed", zap.Error(err))
		return nil, "", "", errcode.ErrInternal
	}

	refreshToken, err := pkgjwt.GenerateToken(user.ID, s.cfg.JWT.RefreshSecret, s.cfg.JWT.RefreshExpire)
	if err != nil {
		logger.Log.Error("generate refresh token failed", zap.Error(err))
		return nil, "", "", errcode.ErrInternal
	}

	return user, accessToken, refreshToken, nil
}

func (s *UserService) RefreshAccessToken(refreshToken string) (string, *errcode.ErrorCode) {
	claims, err := pkgjwt.ParseToken(refreshToken, s.cfg.JWT.RefreshSecret)
	if err != nil {
		return "", errcode.ErrTokenInvalid
	}

	user, _ := dao.GetUserByID(claims.UserID)
	if user == nil || user.Status != 1 {
		return "", errcode.ErrUserNotFound
	}

	accessToken, err := pkgjwt.GenerateToken(user.ID, s.cfg.JWT.AccessSecret, s.cfg.JWT.AccessExpire)
	if err != nil {
		logger.Log.Error("generate access token failed", zap.Error(err))
		return "", errcode.ErrInternal
	}

	return accessToken, nil
}

func (s *UserService) GetProfile(userID uint64) (*model.UserProfile, *errcode.ErrorCode) {
	profile, err := dao.GetUserProfile(userID)
	if err != nil {
		logger.Log.Error("get profile failed", zap.Error(err))
		return nil, errcode.ErrInternal
	}
	if profile == nil {
		return nil, errcode.ErrUserNotFound
	}
	return profile, nil
}

func (s *UserService) UpdateProfile(userID uint64, username, avatar, bio string, birthday *time.Time, gender int8) *errcode.ErrorCode {
	user, _ := dao.GetUserByID(userID)
	if user == nil {
		return errcode.ErrUserNotFound
	}

	if username != "" {
		user.Username = username
	}
	if avatar != "" {
		user.Avatar = avatar
	}
	if bio != "" {
		user.Bio = bio
	}
	if birthday != nil {
		user.Birthday = birthday
	}
	if gender >= 0 && gender <= 2 {
		user.Gender = gender
	}

	if err := dao.UpdateUser(user); err != nil {
		logger.Log.Error("update user failed", zap.Error(err))
		return errcode.ErrInternal
	}
	return nil
}

func (s *UserService) ChangePassword(userID uint64, oldPassword, newPassword string) *errcode.ErrorCode {
	user, _ := dao.GetUserByID(userID)
	if user == nil {
		return errcode.ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return errcode.ErrPasswordWrong
	}

	if !s.isPasswordValid(newPassword) {
		return &errcode.ErrorCode{Code: 20000, Message: "new password format invalid"}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Error("bcrypt hash failed", zap.Error(err))
		return errcode.ErrInternal
	}

	if err := dao.UpdateUserPassword(userID, string(hash)); err != nil {
		logger.Log.Error("update password failed", zap.Error(err))
		return errcode.ErrInternal
	}
	return nil
}

func (s *UserService) BindEmail(userID uint64, email, code string) *errcode.ErrorCode {
	if !s.verifyCode("email:"+email, code) {
		return errcode.ErrCodeIncorrect
	}

	existing, _ := dao.GetUserByEmail(email)
	if existing != nil && existing.ID != userID {
		return errcode.ErrUserExists
	}

	if err := dao.UpdateUserEmail(userID, email); err != nil {
		logger.Log.Error("bind email failed", zap.Error(err))
		return errcode.ErrInternal
	}
	return nil
}

func (s *UserService) BindPhone(userID uint64, phone, code string) *errcode.ErrorCode {
	if !s.verifyCode("phone:"+phone, code) {
		return errcode.ErrCodeIncorrect
	}

	existing, _ := dao.GetUserByPhone(phone)
	if existing != nil && existing.ID != userID {
		return errcode.ErrUserExists
	}

	if err := dao.UpdateUserPhone(userID, phone); err != nil {
		logger.Log.Error("bind phone failed", zap.Error(err))
		return errcode.ErrInternal
	}
	return nil
}

func (s *UserService) SendVerificationCode(target string) (string, *errcode.ErrorCode) {
	prefix := "email"
	if !s.isEmail(target) {
		prefix = "phone"
	}

	key := fmt.Sprintf("code:send:%s:%s", prefix, target)

	ctx := context.Background()
	exists, _ := redis.Exists(ctx, key)
	if exists {
		return "", errcode.ErrCodeSendTooFast
	}

	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	codeKey := fmt.Sprintf("code:%s:%s", prefix, target)
	if err := redis.Set(ctx, codeKey, code, time.Duration(s.cfg.Code.Expire)*time.Second); err != nil {
		logger.Log.Error("redis set code failed", zap.Error(err))
		return "", errcode.ErrInternal
	}

	if err := redis.Set(ctx, key, "1", time.Duration(s.cfg.Code.SendInterval)*time.Second); err != nil {
		logger.Log.Error("redis set send interval failed", zap.Error(err))
	}

	logger.Log.Info("verification code sent",
		zap.String("target", target),
		zap.String("code", code),
	)

	return code, nil
}

func (s *UserService) verifyCode(key, code string) bool {
	ctx := context.Background()
	stored, err := redis.Get(ctx, "code:"+key)
	if err != nil || stored != code {
		return false
	}
	redis.Del(ctx, "code:"+key)
	return true
}

func (s *UserService) isEmail(v string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(v)
}

func (s *UserService) Cfg() *config.Config {
	return s.cfg
}

func (s *UserService) isPasswordValid(password string) bool {
	if len(password) < 6 || len(password) > 128 {
		return false
	}
	hasLetter := false
	hasDigit := false
	for _, ch := range password {
		switch {
		case unicode.IsLetter(ch):
			hasLetter = true
		case unicode.IsDigit(ch):
			hasDigit = true
		}
	}
	return hasLetter && hasDigit
}
