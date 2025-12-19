package handlers

import (
	"dragon-core/internal/repositories"
	"dragon-core/internal/services"
	"github.com/gofiber/fiber/v2"
)

// 1. GET /api/question
// دالة لجلب سؤال (حالياً نجلب السؤال رقم 1 دائماً للتجربة)
func GetQuestion(c *fiber.Ctx) error {
	// مستقبلاً سنجعل الرقم عشوائياً
	// حالياً نطلب السؤال رقم 1 الذي خزنّاه في الفصل السابق
	question, err := repositories.GetQuestionCached(1)
	
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status": "error", 
			"message": "Question not found, maybe Dragon Balls are inert?",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"data": question,
	})
}

func SubmitAnswer(c *fiber.Ctx) error {
    // 1. استخراج هوية المستخدم المؤمنة من الـ Middleware
    // (يجب تحويلها لأن Locals تخزنها كـ interface)
    userIDFromToken := c.Locals("userID").(uint)

    type RequestWithTime struct {
        // لم نعد نحتاج UserID هنا، نحذفه من الهيكل أو نتجاهله
        QuestionID uint   `json:"question_id"`
        Selected   string `json:"selected"`
        TimeTaken  int    `json:"time_taken"`
    }

    var req RequestWithTime
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid input"})
    }

    // 2. استخدام ID الموثوق به بدلاً من req.UserID
    response, err := services.ProcessAnswer(userIDFromToken, req.QuestionID, req.Selected, req.TimeTaken)
	
	if err != nil {
		// معالجة الأخطاء (مثل نفاذ الطاقة)
		return c.Status(400).JSON(fiber.Map{
			"status": "error",
			"message": err.Error(),
		})
	}

	// 3. الرد
	return c.Status(200).JSON(response)
}