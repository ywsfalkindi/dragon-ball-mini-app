package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config يحمل جميع إعدادات التطبيق في مكان واحد
type Config struct {
	AppPort       string
	DBHost        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBPort        string
	JwtSecret     string
	TelegramToken string
	AppEnv        string // dev or prod
}

// LoadConfig يقوم بقراءة ملف.env وتحميل المتغيرات
func LoadConfig() (*Config, error) {
	// تحميل المتغيرات من ملف.env إذا كان موجوداً
	// نتجاهل الخطأ هنا لأننا قد نعتمد على متغيرات النظام مباشرة في بيئة الإنتاج (Docker/Cloud)
	_ = godotenv.Load()

	cfg := &Config{
		AppPort:       getEnv("APP_PORT", "3000"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", "password"), // غيرها في الإنتاج!
		DBName:        getEnv("DB_NAME", "dragonball"),
		DBPort:        getEnv("DB_PORT", "5432"),
		JwtSecret:     getEnv("JWT_SECRET", "super_secret_genki_dama_key"),
		TelegramToken: getEnv("BOT_TOKEN", ""),
		AppEnv:        getEnv("APP_ENV", "dev"),
	}

	// التحقق من القيم الحرجة (Validation)
	if cfg.TelegramToken == "" && cfg.AppEnv == "prod" {
		return nil, fmt.Errorf("FATAL: BOT_TOKEN is missing in production environment!")
	}

	return cfg, nil
}

// getEnv دالة مساعدة لجلب القيمة أو استخدام البديل الافتراضي
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}