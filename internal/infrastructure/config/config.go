package config

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConfig struct {
    Host     string
    Port     string
    Database string
    User     string
    Password string
    SSLMode  string
}

func NewDBConfig() DBConfig {
    return DBConfig{
        Host:     getEnv("DB_HOST", "localhost"),
        Port:     getEnv("DB_PORT", "5432"),
        Database: getEnv("DB_NAME", "budget_book_db"),
        User:     getEnv("DB_USER", "budget_book_user"),
        Password: getEnv("DB_PASSWORD", "budget_book_pass"),
        SSLMode:  getEnv("DB_SSLMODE", "disable"),
    }
}

func NewDBPool(cfg DBConfig) (*pgxpool.Pool, error) {
    dsn := fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s?sslmode=%s&timezone=Asia/Tokyo",
        cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.SSLMode,
    )

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("DB接続プールの作成に失敗しました: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("DBへの接続確認に失敗しました: %w", err)
	}

	return pool, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}