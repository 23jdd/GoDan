package errcode

import "net/http"

type ErrorCode struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var (
	// 通用
	Success          = &ErrorCode{0, "success"}
	ErrInternal      = &ErrorCode{10001, "internal server error"}
	ErrInvalidParams = &ErrorCode{10002, "invalid parameters"}
	ErrUnauthorized  = &ErrorCode{10003, "unauthorized"}
	ErrForbidden     = &ErrorCode{10004, "forbidden"}
	ErrNotFound      = &ErrorCode{10005, "not found"}

	// 用户
	ErrUserExists      = &ErrorCode{20001, "user already exists"}
	ErrUserNotFound    = &ErrorCode{20002, "user not found"}
	ErrPasswordWrong   = &ErrorCode{20003, "wrong password"}
	ErrCodeIncorrect   = &ErrorCode{20004, "verification code incorrect"}
	ErrTokenExpired    = &ErrorCode{20005, "token expired"}
	ErrTokenInvalid    = &ErrorCode{20006, "token invalid"}
	ErrUserBanned      = &ErrorCode{20007, "user is banned"}
	ErrCodeSendTooFast = &ErrorCode{20008, "verification code sent too frequently"}
	ErrSelfFollow      = &ErrorCode{20009, "cannot follow yourself"}
	ErrAlreadyFollowed = &ErrorCode{20010, "already followed"}
	ErrNotFollowed     = &ErrorCode{20011, "not followed"}
	ErrSelfBlock       = &ErrorCode{20012, "cannot block yourself"}

	// 视频
	ErrVideoNotFound = &ErrorCode{30001, "video not found"}

	// 评论
	ErrCommentNotFound = &ErrorCode{40001, "comment not found"}
)

func (e *ErrorCode) HTTPStatus() int {
	switch e.Code {
	case 0:
		return http.StatusOK
	case 10003:
		return http.StatusUnauthorized
	case 10004:
		return http.StatusForbidden
	case 10005:
		return http.StatusNotFound
	default:
		if e.Code >= 20000 && e.Code < 30000 {
			return http.StatusBadRequest
		}
		return http.StatusInternalServerError
	}
}
