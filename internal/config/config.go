package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Env               string
	Port              string
	JWTSecret         []byte
	JWTExpiry         time.Duration
	DatabasePath      string
	CookieDomain      string
	SecureCookie      bool
	PlaygroundEnabled bool
	LogLevel          string

	// SQL 连接池配置（可选）
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// Load 从 .env 文件和环境变量加载配置
func Load() *Config {
	// 加载 .env 文件（如果存在）
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	cfg := &Config{
		Env:          getEnv("ENV", "development"),
		Port:         getEnv("PORT", "8080"),
		DatabasePath: getEnv("DATABASE_PATH", "database.sqlite"),
		CookieDomain: getEnv("COOKIE_DOMAIN", "localhost"),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
	}

	// JWT 配置
	secret := getEnv("JWT_SECRET", "fallback-super-secret-key-change-me-in-production")
	cfg.JWTSecret = []byte(secret)

	// JWT 过期时间（默认 24 小时）
	expiryHours, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	cfg.JWTExpiry = time.Duration(expiryHours) * time.Hour

	// Cookie Secure 标志
	cfg.SecureCookie = getEnv("HTTPS_ENABLED", "false") == "true"

	// Playground 开关
	cfg.PlaygroundEnabled = getEnv("PLAYGROUND_ENABLED", "true") == "true"

	// SQLite 连接池（可选）
	cfg.MaxOpenConns, _ = strconv.Atoi(getEnv("SQL_MAX_OPEN_CONNS", "25"))
	cfg.MaxIdleConns, _ = strconv.Atoi(getEnv("SQL_MAX_IDLE_CONNS", "25"))
	lifetime, _ := time.ParseDuration(getEnv("SQL_CONN_MAX_LIFETIME", "5m"))
	cfg.ConnMaxLifetime = lifetime
	idletime, _ := time.ParseDuration(getEnv("SQL_CONN_MAX_IDLETIME", "5m"))
	cfg.ConnMaxIdleTime = idletime

	return cfg
}

// 辅助函数：获取环境变量，支持默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
