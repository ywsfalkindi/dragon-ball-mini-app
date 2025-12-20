package main

import (
	"log"

	"dragon-core/internal/config"
	"dragon-core/internal/database"
	"dragon-core/internal/handlers"
	"dragon-core/internal/middleware"
	"dragon-core/internal/models"
	"dragon-core/internal/repository"
	"dragon-core/internal/worker"

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
	err = database.DB.AutoMigrate(&models.User{}, &models.Match{}, &models.Question{}, &models.Score{})
	if err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}
	log.Println("✅ Database tables migrated successfully")

	// 👇👇👇 (جديد) زرع الأسئلة إذا كانت القاعدة فارغة 👇👇👇
	seedQuestions()
	// 👆👆👆 ------------------------------------ 👆👆👆

	log.Println("📥 Loading questions into Redis Cache...")
	if err := repository.CacheAllQuestionIDs(); err != nil {
		log.Printf("⚠️ Warning: Failed to cache questions: %v", err)
	}

	// ⚡ تشغيل العامل في الخلفية
	log.Println("👷 Starting Background Worker...")
	worker.StartSyncWorker()

	// 5. Repositories
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

// دالة لزرع الأسئلة الأولية
func seedQuestions() {
	var count int64
	database.DB.Model(&models.Question{}).Count(&count)
	if count == 0 {
		log.Println("🌱 Database is empty. Seeding initial questions...")
		questions := []models.Question{
			{QuestionText: "من هو أول سوبر سايان ظهر في السلسلة؟", OptionA: "فيجيتا", OptionB: "غوكو", OptionC: "جوهان", OptionD: "برولي", CorrectOption: "B", Difficulty: 1},
			{QuestionText: "ما هو اسم والد غوكو؟", OptionA: "راديتز", OptionB: "نابا", OptionC: "باردوك", OptionD: "كينغ فيجيتا", CorrectOption: "C", Difficulty: 1},
			{QuestionText: "كم عدد كرات التنين؟", OptionA: "5", OptionB: "6", OptionC: "7", OptionD: "8", CorrectOption: "C", Difficulty: 1},
			{QuestionText: "من قام بتدمير كوكب فيجيتا؟", OptionA: "فريزا", OptionB: "سيل", OptionC: "ماجين بو", OptionD: "بيروس", CorrectOption: "A", Difficulty: 1},
			{QuestionText: "ما هي التقنية التي تعلمها غوكو من الكاي الشمالي؟", OptionA: "كاميهاميها", OptionB: "كايكين", OptionC: "فاينل فلاش", OptionD: "ماسينكو", CorrectOption: "B", Difficulty: 2},
			{QuestionText: "ما هو لون تحول غوكو في غريزة السوبر (Ultra Instinct)؟", OptionA: "أحمر", OptionB: "أزرق", OptionC: "فضي", OptionD: "ذهبي", CorrectOption: "C", Difficulty: 3},
		}
		database.DB.Create(&questions)
		log.Println("✅ Added initial questions to the database.")
	} else {
		log.Println("ℹ️ Database already has questions.")
	}
}