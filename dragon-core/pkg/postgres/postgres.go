package postgres

import (
	"fmt"
	"time"

	"dragon-core/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewConnection ينشئ اتصالاً جديداً ويعيده ككائن
func NewConnection(cfg *config.Config) (*gorm.DB, error) {
	// 1. تجهيز نص الاتصال (DSN)
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Riyadh",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)

	// 2. إعدادات السجل (Logging)
	// في وضع التطوير نريد رؤية كل استعلام SQL، وفي الإنتاج نريد الأخطاء فقط
	var logLevel logger.LogLevel
	if cfg.AppEnv == "dev" {
		logLevel = logger.Info
	} else {
		logLevel = logger.Error
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}

	// 3. فتح الاتصال
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err!= nil {
		return nil, err
	}

	// 4. ضبط مسبح الاتصالات (Connection Pool) لتحمل الضغط العالي
	sqlDB, err := db.DB()
	if err!= nil {
		return nil, err
	}

	// أقصى عدد للاتصالات الخاملة (الجاهزة للاستخدام فوراً)
	sqlDB.SetMaxIdleConns(10)

	// أقصى عدد للاتصالات المفتوحة في نفس الوقت (لمنع انهيار قاعدة البيانات)
	sqlDB.SetMaxOpenConns(100)

	// مدة حياة الاتصال قبل أن يتم إغلاقه وتجديده (لسلامة الشبكة)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}