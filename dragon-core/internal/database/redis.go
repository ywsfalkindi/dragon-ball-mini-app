package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

// متغير عالمي نستخدمه في أي مكان
var RDB *redis.Client

// سياق العمل (مطلوب في مكتبة Redis الجديدة)
var Ctx = context.Background()

func ConnectRedis() {
	// نقرأ عنوان Redis من البيئة، وإذا لم نجد نستخدم اللوكال
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	RDB = redis.NewClient(&redis.Options{
		Addr:     redisAddr, 
		Password: "", 
		DB:       0,
	})

	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("🔴 Failed to connect to Redis:", err)
	}

	fmt.Println("⚡ Redis is ready! Connected to:", redisAddr)
}