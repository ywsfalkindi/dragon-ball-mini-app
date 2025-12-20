package main

import (
	// "context" <--- تم حذف هذا السطر لأنه غير مستخدم ويسبب خطأ
	"log"

	"dragon-core/internal/config"
	"dragon-core/internal/database"
	// "dragon-core/internal/domain" <--- سنستبدل هذا بـ models
	"dragon-core/internal/models" // <--- الجديد: هنا توجد الجداول (User, Question, Score)
	"dragon-core/internal/handlers"
	"dragon-core/internal/middleware"
	"dragon-core/internal/repository" // تأكد أن ملفات user_repo.go موجودة هنا

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// 1. Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Config error: %v", err)
	}

	// 2. Database
	database.ConnectDB(cfg)

	// 3. Redis
	database.ConnectRedis()

	// 4. Migrations
	// التصحيح: نستخدم models بدلاً من domain لأن الجداول معرفة هناك
	err = database.DB.AutoMigrate(&models.User{}, &models.Match{}, &models.Question{}, &models.Score{})
	if err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}
	log.Println("✅ Database tables migrated successfully")

	// 5. Repositories
	// ملاحظة هامة: تأكد أن ملف 'internal/repository/user_repo.go' موجود وفيه دالة NewUserRepo
	userRepo := repository.NewUserRepo(database.DB)
	leaderboardRepo := repository.NewLeaderboardRepo(database.RDB)

	_ = userRepo
	_ = leaderboardRepo

	// 6. Fiber App
	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Routes
	api := app.Group("/api")
	api.Get("/health", handlers.HealthCheck)
	api.Post("/auth/login", handlers.Login)
	api.Post("/auth/refresh", handlers.RefreshToken)

	protected := api.Group("/protected")
	protected.Use(middleware.Protected())
	
	protected.Get("/question", handlers.GetQuestion)
	protected.Post("/answer", handlers.SubmitAnswer)

	log.Printf("🚀 Server running on port %s", cfg.AppPort)
	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatalf("❌ Server error: %v", err)
	}
}