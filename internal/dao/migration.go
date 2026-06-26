package dao

import "godan/internal/pkg/database"

func AutoMigrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
			username VARCHAR(50) NOT NULL DEFAULT '',
			email VARCHAR(100) NOT NULL DEFAULT '',
			phone VARCHAR(20) NOT NULL DEFAULT '',
			password_hash VARCHAR(255) NOT NULL DEFAULT '',
			avatar VARCHAR(500) NOT NULL DEFAULT '',
			bio VARCHAR(500) NOT NULL DEFAULT '',
			birthday DATE NULL,
			gender TINYINT NOT NULL DEFAULT 0,
			status TINYINT NOT NULL DEFAULT 1,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			UNIQUE KEY uk_email (email),
			UNIQUE KEY uk_phone (phone)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		`CREATE TABLE IF NOT EXISTS follows (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
			follower_id BIGINT UNSIGNED NOT NULL,
			followee_id BIGINT UNSIGNED NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			UNIQUE KEY uk_follower_followee (follower_id, followee_id),
			INDEX idx_followee_id (followee_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		`CREATE TABLE IF NOT EXISTS user_blocklists (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
			user_id BIGINT UNSIGNED NOT NULL,
			blocked_user_id BIGINT UNSIGNED NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			UNIQUE KEY uk_user_blocked (user_id, blocked_user_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
	}

	for _, q := range queries {
		if _, err := database.DB.Exec(q); err != nil {
			return err
		}
	}
	return nil
}
