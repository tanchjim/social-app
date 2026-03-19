package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	ServerPort string
	GinMode    string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// JWT
	JWTSecret            string
	JWTExpireHours       int
	JWTRefreshExpireDays int

	// Tencent COS
	COSSecretID  string
	COSSecretKey string
	COSBucket    string
	COSRegion    string

	// Content Moderation
	ModerationAPIKey  string
	ModerationAPIURL  string
}

func Load() *Config {
	// Load .env file if exists (ignore error if file not found)
	_ = godotenv.Load()

	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		GinMode:    getEnv("GIN_MODE", "debug"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "social_app"),

		JWTSecret:            getEnv("JWT_SECRET", "your_jwt_secret_key"),
		JWTExpireHours:       getEnvInt("JWT_EXPIRE_HOURS", 24),
		JWTRefreshExpireDays: getEnvInt("JWT_REFRESH_EXPIRE_DAYS", 7),

		COSSecretID:  getEnv("COS_SECRET_ID", ""),
		COSSecretKey: getEnv("COS_SECRET_KEY", ""),
		COSBucket:    getEnv("COS_BUCKET", ""),
		COSRegion:    getEnv("COS_REGION", "ap-guangzhou"),

		ModerationAPIKey: getEnv("MODERATION_API_KEY", ""),
		ModerationAPIURL: getEnv("MODERATION_API_URL", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}
