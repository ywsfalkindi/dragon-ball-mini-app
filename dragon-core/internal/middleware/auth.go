package middleware

import (
	"dragon-core/internal/auth"
	"github.com/gofiber/fiber/v2"
)

// التوكن الخاص بك أصبح جاهزاً للعمل
const BOT_TOKEN = "8561338309:AAG1WFHGJgsh4ZkKMWviAhUhJHK2qWKOdJg" 

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. البحث عن التصريح في الـ Header
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"status": "error", "message": "Who are you? No ID found! 🕵️‍♂️"})
		}

		// استخدمنا القيمة مباشرة دون الحاجة لمكتبة strings
		initData := authHeader

		// 2. التحقق من صحة البيانات
		isValid, err := auth.ValidateWebAppData(initData, BOT_TOKEN)

		if err != nil || !isValid {
			return c.Status(403).JSON(fiber.Map{"status": "error", "message": "Fake Saiyan Detected! Access Denied! 🚫"})
		}

		return c.Next()
	}
}