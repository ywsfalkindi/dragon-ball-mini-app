package handlers

import (
	"dragon-core/internal/database"
	"dragon-core/internal/repository"
	"dragon-core/internal/services"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GET /api/question
func GetQuestion(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	// 1. جلب السؤال (يمكن جعله عشوائياً لاحقاً)
	questionID := uint(1) 
	question, err := repository.GetQuestionCached(questionID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No questions found"})
	}

	// 2. 🛡️ Security: بدء العداد الزمني في السيرفر (Redis)
	// المفتاح: game:timer:{user_id}:{question_id}
	timerKey := fmt.Sprintf("game:timer:%d:%d", userID, question.ID)
	
	// نخزن وقت الآن بصيغة UnixNano (دقة عالية جداً)
	now := time.Now().UnixMilli()
	
	// مدة صلاحية المفتاح قصيرة (مثلاً دقيقة واحدة) لتنظيف الذاكرة تلقائياً
	database.RDB.Set(database.Ctx, timerKey, now, 2*time.Minute)

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   question,
	})
}

// POST /api/protected/answer
func SubmitAnswer(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	// الهيكل لم يعد يحتاج TimeTaken لأن السيرفر سيحسبه
	type AnswerRequest struct {
		QuestionID uint   `json:"question_id"`
		Selected   string `json:"selected"`
	}

	var req AnswerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid input"})
	}

	// استدعاء الخدمة لمعالجة الإجابة
	response, err := services.ProcessAnswer(userID, req.QuestionID, req.Selected)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(response)
}