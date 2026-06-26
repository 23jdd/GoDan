package service

import (
	"go.uber.org/zap"

	"godan/internal/dao"
	"godan/internal/model"
	"godan/internal/pkg/errcode"
	"godan/internal/pkg/logger"
)

type AdminService struct{}

func NewAdminService() *AdminService {
	return &AdminService{}
}

// --- User ---

func (s *AdminService) GetUserList(page, pageSize int) ([]model.User, int64, *errcode.ErrorCode) {
	offset, limit := paginate(page, pageSize)
	users, total, err := dao.GetUserList(offset, limit)
	if err != nil {
		logger.Log.Error("get user list failed", zap.Error(err))
		return nil, 0, errcode.ErrInternal
	}
	if users == nil {
		users = []model.User{}
	}
	return users, total, nil
}

func (s *AdminService) BanUser(adminID, userID uint64) *errcode.ErrorCode {
	if adminID == userID {
		return &errcode.ErrorCode{Code: 70001, Message: "cannot ban yourself"}
	}
	u, _ := dao.GetUserByID(userID)
	if u == nil {
		return errcode.ErrUserNotFound
	}
	if u.Role >= 1 {
		return &errcode.ErrorCode{Code: 70002, Message: "cannot ban admin"}
	}
	if err := dao.UpdateUserStatus(userID, 0); err != nil {
		return errcode.ErrInternal
	}
	return nil
}

func (s *AdminService) UnbanUser(userID uint64) *errcode.ErrorCode {
	if err := dao.UpdateUserStatus(userID, 1); err != nil {
		return errcode.ErrInternal
	}
	return nil
}

func (s *AdminService) SetRole(adminID uint64, userID uint64, role int8) *errcode.ErrorCode {
	if adminID == userID {
		return &errcode.ErrorCode{Code: 70001, Message: "cannot change your own role"}
	}
	if role < 0 || role > 2 {
		return &errcode.ErrorCode{Code: 70003, Message: "invalid role"}
	}
	if err := dao.UpdateUserRole(userID, role); err != nil {
		return errcode.ErrInternal
	}
	return nil
}

// --- Review ---

func (s *AdminService) GetPendingVideos(page, pageSize int) ([]model.Video, int64, *errcode.ErrorCode) {
	offset, limit := paginate(page, pageSize)
	videos, total, err := dao.GetPendingVideos(offset, limit)
	if err != nil {
		return nil, 0, errcode.ErrInternal
	}
	if videos == nil {
		videos = []model.Video{}
	}
	return videos, total, nil
}

func (s *AdminService) ApproveVideo(videoID uint64) *errcode.ErrorCode {
	if err := dao.UpdateVideoStatus(videoID, model.VideoStatusPublished); err != nil {
		return errcode.ErrInternal
	}
	return nil
}

func (s *AdminService) RejectVideo(videoID uint64) *errcode.ErrorCode {
	if err := dao.UpdateVideoStatus(videoID, model.VideoStatusRejected); err != nil {
		return errcode.ErrInternal
	}
	return nil
}

// --- Category ---

func (s *AdminService) CreateCategory(name string, sort int) (uint64, *errcode.ErrorCode) {
	id, err := dao.CreateCategory(name, sort)
	if err != nil {
		return 0, errcode.ErrInternal
	}
	return id, nil
}

func (s *AdminService) UpdateCategory(id uint64, name string, sort int) *errcode.ErrorCode {
	if err := dao.UpdateCategory(id, name, sort); err != nil {
		return errcode.ErrInternal
	}
	return nil
}

func (s *AdminService) DeleteCategory(id uint64) *errcode.ErrorCode {
	if err := dao.DeleteCategory(id); err != nil {
		return errcode.ErrInternal
	}
	return nil
}

func (s *AdminService) GetCategories() ([]model.Category, *errcode.ErrorCode) {
	list, err := dao.GetCategoryList()
	if err != nil {
		return nil, errcode.ErrInternal
	}
	if list == nil {
		list = []model.Category{}
	}
	return list, nil
}

// --- Dashboard ---

func (s *AdminService) GetDashboard() (map[string]int64, *errcode.ErrorCode) {
	stats, err := dao.GetDashboardStats()
	if err != nil {
		return nil, errcode.ErrInternal
	}
	return stats, nil
}

// --- Report ---

func (s *AdminService) CreateReport(userID uint64, targetType int8, targetID, reason string) *errcode.ErrorCode {
	if err := dao.CreateReport(userID, targetType, targetID, reason); err != nil {
		return errcode.ErrInternal
	}
	return nil
}

func (s *AdminService) GetReports(page, pageSize int) ([]model.Report, int64, *errcode.ErrorCode) {
	offset, limit := paginate(page, pageSize)
	list, total, err := dao.GetReportList(offset, limit)
	if err != nil {
		return nil, 0, errcode.ErrInternal
	}
	if list == nil {
		list = []model.Report{}
	}
	return list, total, nil
}

func (s *AdminService) HandleReport(reportID uint64, status int8) *errcode.ErrorCode {
	if err := dao.UpdateReportStatus(reportID, status); err != nil {
		return errcode.ErrInternal
	}
	return nil
}
