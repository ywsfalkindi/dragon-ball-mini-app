package database

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

// متغير عالمي نستخدمه في أي مكان
var RDB *redis.Client

// سياق العمل (مطلوب في مكتبة Redis الجديدة)
var Ctx = context.Background()

func ConnectRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // عنوان دوكر
		Password: "",               // لا توجد كلمة مرور افتراضياً
		DB:       0,                // قاعدة البيانات الافتراضية
	})

	// تجربة الاتصال (Ping)
	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("🔴 Failed to connect to Redis:", err)
	}

	fmt.Println("⚡ Redis (Ultra Instinct) is ready!")
}