package dao

import (
	"database/sql"
	"errors"
	"fmt"

	"godan/internal/model"
	"godan/internal/pkg/database"
)

func GetUserByID(id uint64) (*model.User, error) {
	var u model.User
	err := database.DB.Get(&u, "SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("GetUserByID: %w", err)
	}
	return &u, nil
}

func GetUserByEmail(email string) (*model.User, error) {
	if email == "" {
		return nil, nil
	}
	var u model.User
	err := database.DB.Get(&u, "SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("GetUserByEmail: %w", err)
	}
	return &u, nil
}

func GetUserByPhone(phone string) (*model.User, error) {
	if phone == "" {
		return nil, nil
	}
	var u model.User
	err := database.DB.Get(&u, "SELECT * FROM users WHERE phone = ?", phone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("GetUserByPhone: %w", err)
	}
	return &u, nil
}

func CreateUser(u *model.User) (uint64, error) {
	result, err := database.DB.Exec(
		`INSERT INTO users (username, email, phone, password_hash, avatar, bio, birthday, gender, status)
		 VALUES (?, NULLIF(?, ''), NULLIF(?, ''), ?, ?, ?, ?, ?, ?)`,
		u.Username, u.Email, u.Phone, u.PasswordHash, u.Avatar, u.Bio, u.Birthday, u.Gender, u.Status,
	)
	if err != nil {
		return 0, fmt.Errorf("CreateUser: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreateUser LastInsertId: %w", err)
	}
	return uint64(id), nil
}

func UpdateUser(u *model.User) error {
	_, err := database.DB.Exec(
		`UPDATE users SET username=?, avatar=?, bio=?, birthday=?, gender=?, updated_at=NOW()
		 WHERE id=?`,
		u.Username, u.Avatar, u.Bio, u.Birthday, u.Gender, u.ID,
	)
	return err
}

func UpdateUserPassword(userID uint64, passwordHash string) error {
	_, err := database.DB.Exec(
		"UPDATE users SET password_hash=?, updated_at=NOW() WHERE id=?",
		passwordHash, userID,
	)
	return err
}

func UpdateUserEmail(userID uint64, email string) error {
	_, err := database.DB.Exec(
		"UPDATE users SET email=NULLIF(?, ''), updated_at=NOW() WHERE id=?",
		email, userID,
	)
	return err
}

func UpdateUserPhone(userID uint64, phone string) error {
	_, err := database.DB.Exec(
		"UPDATE users SET phone=NULLIF(?, ''), updated_at=NOW() WHERE id=?",
		phone, userID,
	)
	return err
}

func GetUserProfile(userID uint64) (*model.UserProfile, error) {
	var p model.UserProfile
	err := database.DB.Get(&p, `
		SELECT u.id, u.username, u.avatar, u.bio, u.birthday, u.gender, u.created_at,
		       COALESCE(v.cnt, 0) AS video_count,
		       COALESCE(f1.cnt, 0) AS follower_count,
		       COALESCE(f2.cnt, 0) AS followee_count
		FROM users u
		LEFT JOIN (SELECT user_id, COUNT(*) AS cnt FROM videos GROUP BY user_id) v ON v.user_id = u.id
		LEFT JOIN (SELECT followee_id, COUNT(*) AS cnt FROM follows GROUP BY followee_id) f1 ON f1.followee_id = u.id
		LEFT JOIN (SELECT follower_id, COUNT(*) AS cnt FROM follows GROUP BY follower_id) f2 ON f2.follower_id = u.id
		WHERE u.id = ?`, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("GetUserProfile: %w", err)
	}
	return &p, nil
}
