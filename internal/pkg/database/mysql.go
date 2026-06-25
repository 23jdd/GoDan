package database

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"godan/internal/config"
)

var DB *sqlx.DB

func InitMySQL(cfg *config.MySQLConfig) error {
	db, err := sqlx.Connect("mysql", cfg.DSN())
	if err != nil {
		return fmt.Errorf("failed to connect mysql: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping mysql: %w", err)
	}

	DB = db
	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
