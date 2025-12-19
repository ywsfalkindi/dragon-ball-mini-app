package handlers

import (
	"dragon-core/internal/auth"
	"dragon-core/internal/models"
	"dragon-core/internal/repositories"
	"encoding/json"
	"net/url"
	"os"

	"github.com/gofiber/fiber/v2"
)

// LoginRequest: الهيكل الذي نتوقعه من الـ Frontend
type LoginRequest struct {
	InitData string `json:"init_data"` // نص تليجرام الطويل
}

func Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// 1. التحقق من صحة البيانات القادمة من تليجرام
	botToken := os.Getenv("BOT_TOKEN") // ⚠️ يجب أن يكون في متغيرات البيئة
	if botToken == "" {
        // تنبيه: للتجربة فقط سنستخدم التوكن القديم إذا لم نجد متغير بيئة
		botToken = "8561338309:AAG1WFHGJgsh4ZkKMWviAhUhJHK2qWKOdJg"
	}

	isValid, err := auth.ValidateWebAppData(req.InitData, botToken)
	if err != nil || !isValid {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid Telegram data: " + err.Error()})
	}

	// 2. استخراج بيانات المستخدم من النص
	// data=...&user={"id":123,"first_name":"Goku",...}&...
	parsedData, _ := url.ParseQuery(req.InitData)
	userJSON := parsedData.Get("user")
	
	type TelegramUser struct {
		ID        int64  `json:"id"`
		FirstName string `json:"first_name"`
		Username  string `json:"username"`
	}
	var tgUser TelegramUser
	json.Unmarshal([]byte(userJSON), &tgUser)

	// 3. البحث عن المستخدم في الداتابيز أو إنشاؤه
	user, err := repositories.GetUserByTelegramID(tgUser.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	var userID uint
	if user == nil {
		// مستخدم جديد! لنقم بتسجيله
		newUser := models.User{
			TelegramID: tgUser.ID,
			Username:   tgUser.Username,
			// Energy: 10 (افتراضي من الداتابيز) [cite: 36]
		}
		if err := repositories.CreateUser(&newUser); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Could not create user"})
		}
		userID = newUser.ID
	} else {
		userID = user.ID
	}

	// 4. إنشاء الـ JWT (جواز السفر)
	token, err := auth.GenerateToken(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not generate token"})
	}

	// 5. الرد بالتوكن
	return c.JSON(fiber.Map{
		"token": token,
		"user":  user, // نعيد بيانات المستخدم أيضاً
	})
}