package jwt

import (
	"errors"
	"time"

	gjwt "github.com/golang-jwt/jwt/v5"

	"godan/internal/pkg/errcode"
)

var (
	ErrTokenExpired = errors.New("token expired")
	ErrTokenInvalid = errors.New("token invalid")
)

type Claims struct {
	UserID uint64 `json:"user_id"`
	gjwt.RegisteredClaims
}

func GenerateToken(userID uint64, secret string, expireSeconds int) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: gjwt.RegisteredClaims{
			ExpiresAt: gjwt.NewNumericDate(time.Now().Add(time.Duration(expireSeconds) * time.Second)),
			IssuedAt:  gjwt.NewNumericDate(time.Now()),
		},
	}

	token := gjwt.NewWithClaims(gjwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenString string, secret string) (*Claims, error) {
	token, err := gjwt.ParseWithClaims(tokenString, &Claims{}, func(token *gjwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		if errors.Is(err, gjwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}

func ErrTokenFromResponse(err error) *errcode.ErrorCode {
	if errors.Is(err, ErrTokenExpired) {
		return errcode.ErrTokenExpired
	}
	return errcode.ErrTokenInvalid
}
