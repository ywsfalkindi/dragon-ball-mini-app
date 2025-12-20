package main

import (
	"context"
	"log"

	"dragon-core/internal/config"
	"dragon-core/internal/domain"
	"dragon-core/internal/repository" // استيراد المستودعات الجديدة
	"dragon-core/pkg/postgres"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/redis/go-redis/v9" // استيراد عميل Redis
)

func main() {
	// 1. Config
	cfg, err := config.LoadConfig()
	if err!= nil {
		log.Fatalf("❌ Config error: %v", err)
	}

	// 2. Database
	db, err := postgres.NewConnection(cfg)
	if err!= nil {
		log.Fatalf("❌ DB error: %v", err)
	}

	// 3. Redis Setup (جديد)
	// نقوم بإنشاء اتصال Redis هنا
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // في الإنتاج، اجعل هذا في الـ config
		Password: "",               // لا توجد كلمة مرور افتراضياً
		DB:       0,                // استخدام الداتابيز الافتراضية
	})
	// اختبار اتصال Redis
	if _, err := rdb.Ping(context.Background()).Result(); err!= nil {
		log.Fatalf("❌ Redis connection failed: %v", err)
	}
	log.Println("✅ Redis connected (Ultra Instinct Ready)")

	// 4. Migrations (ترحيل الجداول)
	// الآن نقوم بإنشاء جدولي Users و Matches
	err = db.AutoMigrate(&domain.User{}, &domain.Match{})
	if err!= nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}
	log.Println("✅ Database tables migrated successfully")

	// 5. تهيئة المستودعات (Dependency Injection)
	// نجهز هذه المتغيرات لاستخدامها لاحقاً في الـ Handlers
	userRepo := repository.NewUserRepo(db)
	leaderboardRepo := repository.NewLeaderboardRepo(rdb)

	// (سنتجاهل تحذير "المتغير غير مستخدم" مؤقتاً لأننا سنستخدمهم في الفصل القادم)
	_ = userRepo
	_ = leaderboardRepo

	// 6. Fiber App
	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "db": "connected", "redis": "connected"})
	})

	log.Printf("🚀 Server running on port %s", cfg.AppPort)
	if err := app.Listen(":" + cfg.AppPort); err!= nil {
		log.Fatalf("❌ Server error: %v", err)
	}
}