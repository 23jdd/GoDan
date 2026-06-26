package service

import (
	"errors"

	"github.com/go-sql-driver/mysql"

	"godan/internal/dao"
	"godan/internal/model"
	"godan/internal/pkg/errcode"
)

type FollowService struct {
	notifSvc *NotificationService
	actSvc   *ActivityService
}

func NewFollowService(notifSvc *NotificationService, actSvc *ActivityService) *FollowService {
	return &FollowService{notifSvc: notifSvc, actSvc: actSvc}
}

func (s *FollowService) Follow(followerID, followeeID uint64) *errcode.ErrorCode {
	if followerID == followeeID {
		return errcode.ErrSelfFollow
	}

	followee, _ := dao.GetUserByID(followeeID)
	if followee == nil {
		return errcode.ErrUserNotFound
	}

	if err := dao.CreateFollow(followerID, followeeID); err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return errcode.ErrAlreadyFollowed
		}
		return errcode.ErrInternal
	}

	if follower, _ := dao.GetUserByID(followerID); follower != nil {
		s.notifSvc.Send(followeeID, model.NotifFollow,
			"关注通知",
			follower.Username+" 关注了你",
			0,
		)
	}

	return nil
}

func (s *FollowService) Unfollow(followerID, followeeID uint64) *errcode.ErrorCode {
	if err := dao.DeleteFollow(followerID, followeeID); err != nil {
		return errcode.ErrInternal
	}
	return nil
}

func (s *FollowService) GetFollowerList(userID uint64, page, pageSize int) ([]model.User, int64, *errcode.ErrorCode) {
	offset, limit := paginate(page, pageSize)
	users, total, err := dao.GetFollowerList(userID, offset, limit)
	if err != nil {
		return nil, 0, errcode.ErrInternal
	}
	if users == nil {
		users = []model.User{}
	}
	return users, total, nil
}

func (s *FollowService) GetFolloweeList(userID uint64, page, pageSize int) ([]model.User, int64, *errcode.ErrorCode) {
	offset, limit := paginate(page, pageSize)
	users, total, err := dao.GetFolloweeList(userID, offset, limit)
	if err != nil {
		return nil, 0, errcode.ErrInternal
	}
	if users == nil {
		users = []model.User{}
	}
	return users, total, nil
}

func (s *FollowService) GetMutualFollows(userID uint64, page, pageSize int) ([]model.User, int64, *errcode.ErrorCode) {
	offset, limit := paginate(page, pageSize)
	users, total, err := dao.GetMutualFollows(userID, offset, limit)
	if err != nil {
		return nil, 0, errcode.ErrInternal
	}
	if users == nil {
		users = []model.User{}
	}
	return users, total, nil
}

func (s *FollowService) Block(userID, blockedUserID uint64) *errcode.ErrorCode {
	if userID == blockedUserID {
		return errcode.ErrSelfBlock
	}

	if err := dao.AddBlock(userID, blockedUserID); err != nil {
		return errcode.ErrInternal
	}
	return nil
}

func (s *FollowService) Unblock(userID, blockedUserID uint64) *errcode.ErrorCode {
	if err := dao.RemoveBlock(userID, blockedUserID); err != nil {
		return errcode.ErrInternal
	}
	return nil
}
