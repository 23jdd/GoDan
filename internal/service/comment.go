package service

import (
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.uber.org/zap"

	"godan/internal/dao"
	"godan/internal/model"
	"godan/internal/pkg/errcode"
	"godan/internal/pkg/logger"
	"godan/internal/pkg/sensitive"
)

type CommentService struct {
	filter *sensitive.DFAFilter
}

func NewCommentService(words []string) *CommentService {
	return &CommentService{
		filter: sensitive.NewDFAFilter(words),
	}
}

func (s *CommentService) CreateComment(userID, videoID uint64, content, parentID, rootID string, replyToUID uint64) (*model.Comment, *errcode.ErrorCode) {
	content = strings.TrimSpace(content)
	if len(content) == 0 || len(content) > 1000 {
		return nil, &errcode.ErrorCode{Code: 40000, Message: "comment content must be 1-1000 chars"}
	}

	if s.filter.Check(content) {
		return nil, &errcode.ErrorCode{Code: 40001, Message: "comment contains inappropriate content"}
	}

	user, _ := dao.GetUserByID(userID)
	if user == nil {
		return nil, errcode.ErrUserNotFound
	}

	c := &model.Comment{
		VideoID:    videoID,
		UserID:     userID,
		Username:   user.Username,
		Avatar:     user.Avatar,
		Content:    content,
		ParentID:   parentID,
		RootID:     rootID,
		ReplyToUID: replyToUID,
	}

	if parentID != "" {
		if rootID == "" {
			c.RootID = parentID
		}
		dao.IncrReplyCount(c.RootID)
	}

	id, err := dao.CreateComment(c)
	if err != nil {
		logger.Log.Error("create comment failed", zap.Error(err))
		return nil, errcode.ErrInternal
	}
	c.ID = bsonID(id)
	return c, nil
}

func (s *CommentService) DeleteComment(userID uint64, commentID string) *errcode.ErrorCode {
	c, err := dao.GetCommentByID(commentID)
	if err != nil {
		return &errcode.ErrorCode{Code: 40002, Message: "comment not found"}
	}

	if c.UserID != userID {
		return errcode.ErrForbidden
	}

	if err := dao.DeleteComment(commentID); err != nil {
		logger.Log.Error("delete comment failed", zap.Error(err))
		return errcode.ErrInternal
	}
	return nil
}

func (s *CommentService) GetRootComments(videoID uint64, sort string, page, pageSize int) ([]model.Comment, int64, *errcode.ErrorCode) {
	offset, limit := paginate(page, pageSize)
	comments, total, err := dao.GetRootComments(videoID, sort, int64(offset), int64(limit))
	if err != nil {
		logger.Log.Error("get comments failed", zap.Error(err))
		return nil, 0, errcode.ErrInternal
	}
	return comments, total, nil
}

func (s *CommentService) GetReplies(rootID string, page, pageSize int) ([]model.Comment, int64, *errcode.ErrorCode) {
	offset, limit := paginate(page, pageSize)
	replies, total, err := dao.GetReplies(rootID, int64(offset), int64(limit))
	if err != nil {
		logger.Log.Error("get replies failed", zap.Error(err))
		return nil, 0, errcode.ErrInternal
	}
	return replies, total, nil
}

func (s *CommentService) LikeComment(userID uint64, commentID string) *errcode.ErrorCode {
	liked, _ := dao.IsCommentLikedByUser(commentID, userID)
	if liked {
		return &errcode.ErrorCode{Code: 40003, Message: "already liked"}
	}

	if err := dao.AddCommentLike(commentID, userID); err != nil {
		logger.Log.Error("like comment failed", zap.Error(err))
		return errcode.ErrInternal
	}
	dao.IncrCommentLikeCount(commentID, 1)
	return nil
}

func (s *CommentService) UnlikeComment(userID uint64, commentID string) *errcode.ErrorCode {
	liked, _ := dao.IsCommentLikedByUser(commentID, userID)
	if !liked {
		return &errcode.ErrorCode{Code: 40004, Message: "not liked yet"}
	}

	if err := dao.RemoveCommentLike(commentID, userID); err != nil {
		logger.Log.Error("unlike comment failed", zap.Error(err))
		return errcode.ErrInternal
	}
	dao.IncrCommentLikeCount(commentID, -1)
	return nil
}

func bsonID(v string) bson.ObjectID {
	id, _ := bson.ObjectIDFromHex(v)
	return id
}
