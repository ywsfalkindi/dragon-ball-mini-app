package handlers

import (
	"dragon-core/internal/database"
	"dragon-core/internal/services"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GET /api/question
func GetQuestion(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	// التغيير: استخدام الدالة العشوائية بدلاً من ID ثابت
	question, err := services.GetRandomQuestion(userID)
	
	if err != nil {
		// إذا انتهت الأسئلة
		if err.Error() == "no more questions available! you beat the game" {
			return c.Status(404).JSON(fiber.Map{"status": "finished", "message": "You completed all training!"})
		}
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No questions found"})
	}

	// ... (باقي كود العداد الزمني Redis يبقى كما هو )
	timerKey := fmt.Sprintf("game:timer:%d:%d", userID, question.ID)
	now := time.Now().UnixMilli()
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