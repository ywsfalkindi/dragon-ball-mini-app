package handlers

import (
	"github.com/gofiber/fiber/v2"
	"dragon-core/internal/models" // استدعاء المودل الذي صنعناه
)

// دالة تفحص صحة السيرفر
// Ctx = Context (سياق الطلب - يحمل كل معلومات الطلب والرد)
func HealthCheck(c *fiber.Ctx) error {
	// الرد باستخدام المودل المرتب
	response := models.JSend{
		Status:  "success",
		Message: "Senzu Bean eaten! Server is full power! 💊",
	}

	// إرسال الرد بصيغة JSON مع كود 200 (OK)
	return c.Status(200).JSON(response)
}