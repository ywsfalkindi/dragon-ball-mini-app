package middleware

import (
	"dragon-core/internal/auth"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// استقبال التوكن من الـ Header
		// الشكل المتوقع: "Bearer eyJhbGciOiJIUzI1NiIsIn..."
		authHeader := c.Get("Authorization")

		// --- التعامل مع الـ Backdoor بأمان ---
		// نقرأ متغير البيئة لنعرف هل نحن في وضع التطوير؟
		// APP_ENV يجب أن يكون "dev" في جهازك، و "prod" عند الرفع
		appEnv := os.Getenv("APP_ENV") 
		
		if appEnv == "dev" && authHeader == "test-token-for-goku" {
			// في وضع التطوير فقط نسمح بالمرور
			// ونفترض أن المستخدم هو رقم 1 (لأغراض التست)
			c.Locals("userID", uint(1))
			return c.Next()
		}
		// ----------------------------------

		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"status": "error", "message": "Missing Authorization Header"})
		}

		// تنظيف التوكن (حذف كلمة Bearer إذا وجدت)
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// التحقق من صحة التوكن
		userID, err := auth.ValidateToken(tokenString)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"status": "error", "message": "Invalid or Expired Token"})
		}

		// أهم خطوة: تخزين رقم المستخدم في الـ Context
		// لكي نستخدمه لاحقاً في ProcessAnswer بدون البحث عنه مجدداً
		c.Locals("userID", userID)

		return c.Next()
	}
}