package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// متغير عام سنستخدمه في كل مكان للوصول للداتابيز
var DB *gorm.DB

// دالة الاتصال: تفتح الخط مع التنين
func ConnectDB() {
	// 1. تجهيز بيانات الاتصال (DSN)
	// host=localhost: لأن الداتابيز في دوكر على نفس الجهاز
	// user=postgres: المستخدم الافتراضي
	// password=mysecretpassword: كلمة السر التي وضعناها في الفصل 2
	// dbname=postgres: اسم قاعدة البيانات الافتراضية
	// port=5432: المنفذ الذي فتحناه في دوكر
	dsn := "host=localhost user=postgres password=123456 dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Riyadh"

	// 2. محاولة الاتصال
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("🔥 Failed to connect to the database! Is Docker running?", err)
	}

	fmt.Println("🐉 Connection to PostgreSQL established successfully!")
}