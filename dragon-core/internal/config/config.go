package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort          string
	DBHost           string
	DBUser           string
	DBPassword       string
	DBName           string
	DBPort           string
	JwtAccessSecret  string // جديد: مفتاح التوكن السريع
	JwtRefreshSecret string // جديد: مفتاح التوكن الطويل
	TelegramToken    string
	AppEnv           string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppPort:          getEnv("APP_PORT", "3000"),
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBUser:           getEnv("DB_USER", "postgres"),
		DBPassword:       getEnv("DB_PASSWORD", "password"),
		DBName:           getEnv("DB_NAME", "dragonball"),
		DBPort:           getEnv("DB_PORT", "5432"),
		// مفاتيح منفصلة للأمان العالي
		JwtAccessSecret:  getEnv("JWT_ACCESS_SECRET", "default_access_secret"),
		JwtRefreshSecret: getEnv("JWT_REFRESH_SECRET", "default_refresh_secret"),
		TelegramToken:    getEnv("BOT_TOKEN", ""),
		AppEnv:           getEnv("APP_ENV", "dev"),
	}

	if cfg.TelegramToken == "" && cfg.AppEnv == "prod" {
		return nil, fmt.Errorf("FATAL: BOT_TOKEN is missing in production environment!")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}