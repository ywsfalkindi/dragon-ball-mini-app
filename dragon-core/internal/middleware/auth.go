package middleware

import (
	"dragon-core/internal/auth"
	"github.com/gofiber/fiber/v2"
)

// التوكن الخاص بك أصبح جاهزاً للعمل
const BOT_TOKEN = "8561338309:AAG1WFHGJgsh4ZkKMWviAhUhJHK2qWKOdJg" 

func Protected() fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")

        // --- إضافة وضع المطور (Backdoor) ---
        // إذا كان التوكن هو "test-token-for-goku"، اسمح بالمرور فوراً
        // هذا مفيد جداً للتجربة في المتصفح دون تعقيدات تليجرام
        if authHeader == "test-token-for-goku" {
            return c.Next()
        }
        // ----------------------------------

        if authHeader == "" {
            return c.Status(401).JSON(fiber.Map{"status": "error", "message": "Who are you? No ID found! 🕵️‍♂️"})
        }

        // ... بقية كود التحقق من تليجرام ...
        isValid, err := auth.ValidateWebAppData(authHeader, BOT_TOKEN)
        if err != nil || !isValid {
            return c.Status(403).JSON(fiber.Map{"status": "error", "message": "Fake Saiyan Detected! Access Denied! 🚫"})
        }

        return c.Next()
    }
}