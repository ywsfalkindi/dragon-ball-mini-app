package database

import (
	"dragon-core/internal/config" // استيراد ملف الكونفيج
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// متغير عام سنستخدمه في كل مكان للوصول للداتابيز
var DB *gorm.DB

// ConnectDB: تأخذ الإعدادات كمدخل (Parameter)
func ConnectDB(cfg *config.Config) {
	// 1. تجهيز بيانات الاتصال (DSN) ديناميكياً
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Riyadh",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	// 2. محاولة الاتصال
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("🔥 Failed to connect to database! Error: %v", err)
	}

	fmt.Println("🐉 Connection to PostgreSQL established successfully!")
}