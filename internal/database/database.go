package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"task-scheduler/internal/config"
)

var DB *sql.DB

func Init(cfg config.DatabaseConfig) error {
	var err error

	DB, err = sql.Open("mysql", cfg.DSN())
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	DB.SetMaxOpenConns(cfg.MaxOpenConns)
	DB.SetMaxIdleConns(cfg.MaxIdleConns)
	DB.SetConnMaxLifetime(time.Hour)

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %w", err)
	}

	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
