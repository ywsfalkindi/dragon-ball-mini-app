package main

import (
	"log"
	"dragon-core/internal/database"
	"dragon-core/internal/handlers"
	"dragon-core/internal/models"
	"dragon-core/internal/repositories"
	"dragon-core/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/cors" // المكتبة المستدعاة
)

func main() {
	// 1. الاتصال بقواعد البيانات
	database.ConnectDB()
	database.ConnectRedis()
	database.DB.AutoMigrate(&models.User{}, &models.Score{}, &models.Question{})

	// --- كود تجريبي لاختبار الكاش ---
	q1 := models.Question{
		ID:           1,
		QuestionText: "Who is Goku's father?",
		OptionA:      "Vegeta",
		OptionB:      "Bardock",
		OptionC:      "Nappa",
		OptionD:      "King Vegeta",
		CorrectOption: "B",
	}
	database.DB.Save(&q1) 

	log.Println("--- Testing Cache System ---")
	repositories.GetQuestionCached(1)
	repositories.GetQuestionCached(1)
	log.Println("---------------------------")

	// 2. إعداد سيرفر Fiber
	app := fiber.New(fiber.Config{
		AppName: "Dragon Ball Bot API",
	})

	// تفعيل الـ CORS (هذا السطر سيحل المشكلة)
	// يسمح للمتصفحات (وتطبيق تليجرام) بطلب البيانات من السيرفر
	app.Use(cors.New())

	app.Use(logger.New())

	api := app.Group("/api")
	api.Get("/health", handlers.HealthCheck)
	api.Get("/question", handlers.GetQuestion)
	
	protected := api.Group("/protected", middleware.Protected())
	protected.Post("/answer", handlers.SubmitAnswer)

	log.Println("🔥 Server is going Super Saiyan on port 3000...")
	
	err := app.Listen(":3000")
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}